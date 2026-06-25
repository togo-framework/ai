# ai — documentation

  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />

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

Environment variables read by this plugin (extracted from the source):

| Env var | Notes |
|---|---|
| `AI_DRIVER` | _see provider docs_ |
| `G` | _see provider docs_ |

## Usage

```go
provider := ai.FromKernel(k)
resp, err := provider.Chat(ctx, []ai.Message{{Role: "user", Content: "Hello"}}, ai.Options{})
// streaming + provider.Embed(ctx, texts) for vectors; resp.Usage carries token counts
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/ai
- README: ../README.md
