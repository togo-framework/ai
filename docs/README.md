# ai — documentation

togo AI plugin — unified LLM interface (chat/embed/tools/stream) with a pluggable provider driver registry (openai, anthropic, gemini, ollama, …)

## Overview

Package ai is the togo AI plugin: a unified LLM interface with a pluggable
driver registry. Provider plugins (ai-openai, ai-anthropic, ai-gemini,
ai-ollama, …) register a driver via init(); select one with AI_DRIVER. Mirrors
the mail/storage driver-plugin pattern. A safe "echo" driver is the dev default.

## Install

```bash
togo install togo-framework/ai
```

Set `AI_DRIVER=<provider>` and install a provider driver (ai-openai, ai-anthropic, …).

## Configuration

Environment variables read by this plugin (extracted from the source — see the gateway/provider docs for each value):

| Env var |
|---|
| `AI_DRIVER` |

## Usage

```go
provider := ai.FromKernel(k)
resp, err := provider.Chat(ctx, []ai.Message{{Role: "user", Content: "Hello"}}, ai.Options{})
// streaming + provider.Embed(ctx, texts) for vectors; resp.Usage carries token counts
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/ai
- Full README: ../README.md
