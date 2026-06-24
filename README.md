# togo · ai

The **AI plugin** for [togo](https://to-go.dev) — a unified LLM interface with a pluggable provider driver registry. Install a provider plugin and select it with `AI_DRIVER`.

```bash
togo install togo-framework/ai
togo install togo-framework/ai-openai   # then set AI_DRIVER=openai, OPENAI_API_KEY=…
```

## Interface

```go
type Provider interface {
    Chat(ctx, ChatRequest) (ChatResponse, error)
    Embed(ctx, EmbedRequest) (EmbedResponse, error)
}
// optional: Streamer { ChatStream(ctx, req, onChunk) }
```

## Usage

```go
svc, _ := ai.FromKernel(k)
resp, _ := svc.Chat(ctx, ai.ChatRequest{Messages: []ai.Message{{Role: ai.RoleUser, Content: "Hi"}}})
fmt.Println(resp.Content, resp.Usage.TotalTokens)
```

REST: mount `ai.Handler(k)` under `/api/ai` → `POST /chat`, `POST /embed`.

## Providers

`ai-openai` · `ai-anthropic` · `ai-gemini` · `ai-ollama` · `ai-grok` · `ai-deepseek` · `ai-qwen` · `ai-adk` · `ai-agno`. Capabilities: `ai-rag`, `ai-agentops`. The built-in `echo` driver is the safe dev default (no external calls).

MIT
