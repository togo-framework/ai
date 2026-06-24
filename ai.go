// Package ai is the togo AI plugin: a unified LLM interface with a pluggable
// driver registry. Provider plugins (ai-openai, ai-anthropic, ai-gemini,
// ai-ollama, …) register a driver via init(); select one with AI_DRIVER. Mirrors
// the mail/storage driver-plugin pattern. A safe "echo" driver is the dev default.
package ai

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/togo-framework/togo"
)

// Chat roles.
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"
)

// Message is a single chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// Tool describes a function the model may call.
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

// ToolCall is a model's request to call a tool (arguments are JSON).
type ToolCall struct {
	Name string `json:"name"`
	Args string `json:"arguments"`
}

// ChatRequest is a chat-completion request.
type ChatRequest struct {
	Model       string    `json:"model,omitempty"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Tools       []Tool    `json:"tools,omitempty"`
}

// Usage reports token consumption (consumed by the billing plugin).
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatResponse is a chat-completion result.
type ChatResponse struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Model     string     `json:"model,omitempty"`
	Usage     Usage      `json:"usage"`
}

// Chunk is a streamed delta.
type Chunk struct {
	Delta string `json:"delta"`
	Done  bool   `json:"done"`
}

// EmbedRequest requests embeddings for one or more inputs.
type EmbedRequest struct {
	Model  string   `json:"model,omitempty"`
	Inputs []string `json:"inputs"`
}

// EmbedResponse holds the resulting vectors.
type EmbedResponse struct {
	Vectors [][]float32 `json:"vectors"`
	Usage   Usage       `json:"usage"`
}

// Provider is the LLM driver interface every ai-* provider plugin implements.
type Provider interface {
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
	Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error)
}

// Streamer is an optional interface for streaming chat completions.
type Streamer interface {
	ChatStream(ctx context.Context, req ChatRequest, onChunk func(Chunk) error) error
}

// DriverFactory builds a Provider from the kernel/env.
type DriverFactory func(k *togo.Kernel) (Provider, error)

var (
	regMu   sync.RWMutex
	drivers = map[string]DriverFactory{}
)

// RegisterDriver registers an AI provider driver (called from a provider plugin's init()).
func RegisterDriver(name string, f DriverFactory) {
	regMu.Lock()
	drivers[name] = f
	regMu.Unlock()
}

// Drivers returns the registered driver names.
func Drivers() []string {
	regMu.RLock()
	defer regMu.RUnlock()
	out := make([]string, 0, len(drivers))
	for n := range drivers {
		out = append(out, n)
	}
	return out
}

func init() {
	RegisterDriver("echo", func(k *togo.Kernel) (Provider, error) { return &echo{}, nil })

	togo.RegisterProviderFunc("ai", togo.PriorityService, func(k *togo.Kernel) error {
		name := os.Getenv("AI_DRIVER")
		if name == "" {
			name = "echo" // dev default until a provider plugin is installed + configured
		}
		regMu.RLock()
		f, ok := drivers[name]
		regMu.RUnlock()
		if !ok {
			return fmt.Errorf("ai: unknown driver %q (install its plugin, e.g. togo install togo-framework/ai-openai)", name)
		}
		p, err := f(k)
		if err != nil {
			return err
		}
		k.Set("ai", &Service{provider: p, driver: name})
		return nil
	})
}

// Service is the kernel-bound AI service.
type Service struct {
	provider Provider
	driver   string
}

// Chat performs a chat completion via the configured provider.
func (s *Service) Chat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	return s.provider.Chat(ctx, req)
}

// Embed returns embeddings via the configured provider.
func (s *Service) Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error) {
	return s.provider.Embed(ctx, req)
}

// ChatStream streams a chat completion; falls back to a single chunk if the
// provider doesn't implement Streamer.
func (s *Service) ChatStream(ctx context.Context, req ChatRequest, onChunk func(Chunk) error) error {
	if st, ok := s.provider.(Streamer); ok {
		return st.ChatStream(ctx, req, onChunk)
	}
	resp, err := s.provider.Chat(ctx, req)
	if err != nil {
		return err
	}
	if e := onChunk(Chunk{Delta: resp.Content}); e != nil {
		return e
	}
	return onChunk(Chunk{Done: true})
}

// Driver returns the active driver name.
func (s *Service) Driver() string { return s.driver }

// FromKernel returns the AI service bound to the kernel.
func FromKernel(k *togo.Kernel) (*Service, bool) {
	v, ok := k.Get("ai")
	if !ok {
		return nil, false
	}
	s, ok := v.(*Service)
	return s, ok
}

// echo is the dev driver: it echoes the last user message, no external calls.
type echo struct{}

func (e *echo) Chat(_ context.Context, req ChatRequest) (ChatResponse, error) {
	last := ""
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == RoleUser {
			last = req.Messages[i].Content
			break
		}
	}
	return ChatResponse{Content: "echo: " + last, Model: "echo"}, nil
}

func (e *echo) Embed(_ context.Context, req EmbedRequest) (EmbedResponse, error) {
	vecs := make([][]float32, len(req.Inputs))
	for i := range vecs {
		vecs[i] = []float32{0, 0, 0}
	}
	return EmbedResponse{Vectors: vecs}, nil
}
