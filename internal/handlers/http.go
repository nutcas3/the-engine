package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"the-engine/internal/cache"
	"the-engine/internal/compositions"
	"the-engine/internal/docs"
	"the-engine/internal/finops"
	"the-engine/internal/health"
	"the-engine/internal/kubernetes"
	"the-engine/internal/rate"
	"the-engine/internal/types"
)

// Handlers holds the dependencies for HTTP handlers
type Handlers struct {
	k8sClient     *kubernetes.Client
	healthChecker *health.Checker
	cache         *cache.Cache
	RateLimiter   *rate.RateLimiter
}

// NewHandlers creates a new Handlers instance
func NewHandlers(k8sClient *kubernetes.Client) *Handlers {
	return &Handlers{
		k8sClient:     k8sClient,
		healthChecker: health.NewChecker("1.0.0"),
		cache:         cache.NewCache(5 * time.Minute),
		RateLimiter:   rate.NewRateLimiter(100, 10), // 100 requests/sec, burst 10
	}
}

// HandleIndex serves the main HTML page
func (h *Handlers) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "web/index.html")
}

// HandleDeployments returns deployment data
func (h *Handlers) HandleDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Try to get real deployments from Kubernetes if available
	if h.k8sClient != nil {
		realDeployments, err := h.k8sClient.GetDeployments()
		if err == nil && len(realDeployments) > 0 {
			json.NewEncoder(w).Encode(realDeployments)
			return
		}
	}

	// Fallback to empty array if no data
	json.NewEncoder(w).Encode([]types.Deployment{})
}

// HandleCompositions returns composition data
func (h *Handlers) HandleCompositions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Try cache first
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

// HandleSwagger returns OpenAPI specification
func (h *Handlers) HandleSwagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	spec := docs.GetSpec()
	json.NewEncoder(w).Encode(spec)
}

// HandleCostEstimate returns cost estimation for provider and tier
func (h *Handlers) HandleCostEstimate(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	tier := r.URL.Query().Get("tier")

	if provider == "" || tier == "" {
		http.Error(w, "provider and tier parameters required", http.StatusBadRequest)
		return
	}

	cost := finops.EstimateCost(provider, tier)

	response := map[string]interface{}{
		"provider":     provider,
		"tier":         tier,
		"monthly_cost": cost,
		"hourly_cost":  cost / 730, // Approximate hours in a month
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleCostMonthly returns comprehensive cost data
func (h *Handlers) HandleCostMonthly(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	team := r.URL.Query().Get("team")
	if team == "" {
		team = "platform"
	}

	currentSpend, err := finops.GetCurrentTeamSpend(ctx, team)
	if err != nil {
		currentSpend = 0
	}

	budget, err := finops.GetTeamBudget(ctx, team)
	if err != nil {
		budget = 2000
	}

	utilization := 0.0
	if budget > 0 {
		utilization = (currentSpend / budget) * 100
	}

	// Get cost recommendations
	recommendations := finops.GetCostRecommendations(ctx, team)

	// Get provider cost comparison
	providers := []string{"aws", "azure", "gcp", "hetzner", "ovh", "digitalocean"}
	providerCosts := make(map[string]float64)
	for _, provider := range providers {
		providerCosts[provider] = finops.GetCurrentSpend(provider)
	}

	response := types.CostResponse{
		Team:            team,
		MonthlySpend:    currentSpend,
		Budget:          budget,
		Utilization:     utilization,
		LastUpdated:     time.Now().Format(time.RFC3339),
		Recommendations: recommendations,
		ProviderCosts:   providerCosts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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
