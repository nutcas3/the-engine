package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"the-engine/internal/finops"
)

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
