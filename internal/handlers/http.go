package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"the-engine/internal/alerts"
	"the-engine/internal/compositions"
	"the-engine/internal/docs"
	"the-engine/internal/finops"
	"the-engine/internal/types"
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

// HandleCostEstimate returns cost estimation for provider and tier
func (h *Handlers) HandleCostEstimate(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	tier := r.URL.Query().Get("tier")

	if provider == "" || tier == "" {
		http.Error(w, "provider and tier parameters required", http.StatusBadRequest)
		return
	}

	cost := finops.EstimateCost(provider, tier)

	response := map[string]any{
		"provider":     provider,
		"tier":         tier,
		"monthly_cost": cost,
		"hourly_cost":  cost / 730,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleCostMonthly returns comprehensive cost data as HTML
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

	recommendations := finops.GetCostRecommendations(ctx, team)

	providers := []string{"aws", "azure", "gcp", "hetzner", "ovh", "digitalocean"}
	providerCosts := make(map[string]float64)
	for _, provider := range providers {
		providerCosts[provider] = finops.GetCurrentSpend(provider)
	}

	var html strings.Builder
	html.WriteString(fmt.Sprintf(`
<div class="space-y-4">
    <div class="grid grid-cols-3 gap-4">
        <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
            <div class="text-gray-500 dark:text-gray-400 text-sm">Monthly Spend</div>
            <div class="text-2xl font-bold text-gray-900 dark:text-white">$%.2f</div>
        </div>
        <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
            <div class="text-gray-500 dark:text-gray-400 text-sm">Budget</div>
            <div class="text-2xl font-bold text-gray-900 dark:text-white">$%.2f</div>
        </div>
        <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
            <div class="text-gray-500 dark:text-gray-400 text-sm">Utilization</div>
            <div class="text-2xl font-bold text-gray-900 dark:text-white">%.1f%%</div>
        </div>
    </div>
    
    <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
        <h3 class="font-semibold text-gray-900 dark:text-white mb-2">Provider Costs</h3>
        <div class="space-y-2">
`, currentSpend, budget, utilization))

	for provider, cost := range providerCosts {
		html.WriteString(fmt.Sprintf(`            <div class="flex justify-between text-gray-700 dark:text-gray-300">
                <span class="capitalize">%s</span>
                <span>$%.2f</span>
            </div>`, provider, cost))
	}

	html.WriteString(`
        </div>
    </div>
    
    <div class="bg-yellow-50 dark:bg-yellow-900/20 p-4 rounded-lg border border-yellow-200 dark:border-yellow-800">
        <h3 class="font-semibold text-yellow-800 dark:text-yellow-200 mb-2">Recommendations</h3>
        <ul class="space-y-1 text-yellow-700 dark:text-yellow-300">
`)

	for _, rec := range recommendations {
		html.WriteString(fmt.Sprintf(`            <li>• %s</li>`, rec))
	}

	html.WriteString(`
        </ul>
    </div>
    
    <div class="text-gray-500 dark:text-gray-400 text-sm">
        Last updated: ` + time.Now().Format(time.RFC3339) + `
    </div>
</div>`)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html.String()))
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

// HandleAlerts returns alert data
func (h *Handlers) HandleAlerts(w http.ResponseWriter, r *http.Request) {
	environment := r.URL.Query().Get("environment")

	var alerts []*alerts.Alert
	if environment != "" {
		alerts = h.alertManager.GetAlertsByEnvironment(environment)
	} else {
		alerts = h.alertManager.GetAlerts()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

// HandleCleanupPolicies returns cleanup policies
func (h *Handlers) HandleCleanupPolicies(w http.ResponseWriter, r *http.Request) {
	policies := h.cleanupManager.GetPolicies()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policies)
}

// HandleCronJobs returns cron jobs
func (h *Handlers) HandleCronJobs(w http.ResponseWriter, r *http.Request) {
	jobs := h.cronScheduler.GetJobs()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// HandleNukeOperations returns nuke operations
func (h *Handlers) HandleNukeOperations(w http.ResponseWriter, r *http.Request) {
	operations := h.nukeManager.GetOperations()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operations)
}

// HandleManualShutdown manually triggers shutdown of a resource
func (h *Handlers) HandleManualShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Provider   string `json:"provider"`
		ResourceID string `json:"resource_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.cleanupManager.ManualShutdown(context.Background(), request.Provider, request.ResourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "shutdown initiated"})
}

// HandleManualNuke manually triggers nuking of an environment
func (h *Handlers) HandleManualNuke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Environment string `json:"environment"`
		Provider    string `json:"provider"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	operation, err := h.nukeManager.NukeEnvironment(context.Background(), request.Environment, request.Provider)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operation)
}
