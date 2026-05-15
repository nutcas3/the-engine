package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"the-engine/internal/types"
)

// HandleHealth returns comprehensive health status
func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	healthResponse := h.healthChecker.Check()
	healthResponse.WriteJSON(w)
}

// HandleSSE provides server-sent events for real-time updates
func (h *Handlers) HandleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			health := types.HealthResponse{
				Status:    "healthy",
				Timestamp: time.Now().Format(time.RFC3339),
				Version:   "1.0.0",
			}

			data, err := json.Marshal(health)
			if err != nil {
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
