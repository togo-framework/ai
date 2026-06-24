package ai

import (
	"encoding/json"
	"net/http"

	"github.com/togo-framework/togo"
)

// Handler exposes the AI service over REST. Mount under /api/ai in your app:
//
//	mux.Handle("/api/ai/", http.StripPrefix("/api/ai", ai.Handler(k)))
//
// Routes: POST /chat (ChatRequest -> ChatResponse), POST /embed (EmbedRequest -> EmbedResponse).
func Handler(k *togo.Kernel) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /chat", func(w http.ResponseWriter, r *http.Request) {
		svc, ok := FromKernel(k)
		if !ok {
			http.Error(w, "ai not configured", http.StatusInternalServerError)
			return
		}
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := svc.Chat(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		writeJSON(w, resp)
	})
	mux.HandleFunc("POST /embed", func(w http.ResponseWriter, r *http.Request) {
		svc, ok := FromKernel(k)
		if !ok {
			http.Error(w, "ai not configured", http.StatusInternalServerError)
			return
		}
		var req EmbedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := svc.Embed(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		writeJSON(w, resp)
	})
	return mux
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
