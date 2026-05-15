package handlers

import (
	"encoding/json"
	"net/http"

	"the-engine/internal/compositions"
	"the-engine/internal/types"
)

// HandleDeployments returns deployment data
func (h *Handlers) HandleDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if h.k8sClient != nil {
		realDeployments, err := h.k8sClient.GetDeployments()
		if err == nil && len(realDeployments) > 0 {
			json.NewEncoder(w).Encode(realDeployments)
			return
		}
	}

	json.NewEncoder(w).Encode([]types.Deployment{})
}

// HandleCompositions returns composition data
func (h *Handlers) HandleCompositions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if cached, found := h.cache.Get("compositions"); found {
		json.NewEncoder(w).Encode(cached)
		return
	}

	compositionList, err := compositions.GetCompositions()
	if err != nil {
		http.Error(w, "Failed to load compositions", http.StatusInternalServerError)
		return
	}

	h.cache.Set("compositions", compositionList)
	json.NewEncoder(w).Encode(compositionList)
}
