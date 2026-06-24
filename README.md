<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/ai</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/ai"><img src="https://pkg.go.dev/badge/github.com/togo-framework/ai.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/ai
```

<!-- /togo-header -->

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

<!-- togo-sponsors -->
---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
