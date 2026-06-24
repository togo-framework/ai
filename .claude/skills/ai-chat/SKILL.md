---
name: ai-chat
description: Add an AI chat/completion feature to this togo app using the ai plugin and a provider.
---

Use the togo `ai` plugin to add chat/completions:

1. Install a provider: `togo install togo-framework/ai-openai` (or ai-anthropic / ai-ollama / ai-gemini…). Set its key in `.env` and `AI_DRIVER=openai`.
2. Call the kernel ai service: `ai.FromKernel(k).Chat(ctx, messages, opts)`.
3. For RAG: `togo install togo-framework/ai-rag` (+ `rag-postgres` for production), ingest documents, then retrieve + generate.

Run `togo serve` and hit `POST /api/ai/chat`.
