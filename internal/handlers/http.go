package handlers

import (
	"encoding/json"
	"net/http"

	"the-engine/internal/docs"
)

// HandleIndex serves the main HTML page
func (h *Handlers) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "web/index.html")
}

// HandleSwagger returns OpenAPI specification
func (h *Handlers) HandleSwagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	spec := docs.GetSpec()
	json.NewEncoder(w).Encode(spec)
}
