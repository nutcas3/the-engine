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
	"the-engine/internal/types"
)

// Handlers holds the dependencies for HTTP handlers
type Handlers struct {
	k8sClient     *kubernetes.Client
	healthChecker *health.Checker
	cache         *cache.Cache
}

// NewHandlers creates a new Handlers instance
func NewHandlers(k8sClient *kubernetes.Client) *Handlers {
	return &Handlers{
		k8sClient:     k8sClient,
		healthChecker: health.NewChecker("1.0.0"),
		cache:         cache.NewCache(5 * time.Minute),
	}
}

// HandleIndex serves the main HTML page
func (h *Handlers) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Sovereign Engine Dashboard</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
</head>
<body class="bg-gray-100 min-h-screen p-8">
    <div class="max-w-6xl mx-auto">
        <h1 class="text-4xl font-bold text-gray-800 mb-8">Sovereign Engine Dashboard</h1>
        
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Composition Selection</h2>
            <div id="composition-selector" class="flex gap-4">
                <select id="composition-select" class="border border-gray-300 rounded-md px-4 py-2 flex-1">
                    <option value="">Select Composition</option>
                </select>
                <button hx-get="/api/compositions" hx-target="#composition-select" hx-swap="innerHTML" 
                    class="bg-blue-600 hover:bg-blue-700 text-white font-medium px-6 py-2 rounded-md transition-colors">
                    Load Compositions
                </button>
            </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Deployments</h2>
            <button hx-get="/api/deployments" hx-target="#deployments" hx-swap="innerHTML"
                class="bg-blue-600 hover:bg-blue-700 text-white font-medium px-6 py-2 rounded-md transition-colors mb-4">
                Load Deployments
            </button>
            <div id="deployments" class="mt-4"></div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Cost Data</h2>
            <div class="flex gap-4 mb-4">
                <input type="text" id="team-input" placeholder="Enter team name" value="platform"
                    class="border border-gray-300 rounded-md px-4 py-2 flex-1">
                <button hx-get="/api/cost/monthly?team=platform" hx-target="#cost" hx-swap="innerHTML" hx-include="#team-input"
                    class="bg-blue-600 hover:bg-blue-700 text-white font-medium px-6 py-2 rounded-md transition-colors">
                    Load Cost Data
                </button>
            </div>
            <div id="cost" class="mt-4"></div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Health Status</h2>
            <button hx-get="/api/health/status" hx-target="#health" hx-swap="innerHTML"
                class="bg-blue-600 hover:bg-blue-700 text-white font-medium px-6 py-2 rounded-md transition-colors mb-4">
                Check Health
            </button>
            <div id="health" class="mt-4"></div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Real-time Updates</h2>
            <button onclick="connectSSE()"
                class="bg-green-600 hover:bg-green-700 text-white font-medium px-6 py-2 rounded-md transition-colors mb-4">
                Connect to SSE Stream
            </button>
            <div id="sse-output" class="mt-4"></div>
        </div>
    </div>
    
    <script>
        function connectSSE() {
            const eventSource = new EventSource('/api/stream');
            eventSource.onmessage = function(event) {
                const data = JSON.parse(event.data);
                const output = document.getElementById('sse-output');
                output.innerHTML += '<div class="bg-gray-50 rounded-md p-4 mb-2 border border-gray-200">' + 
                    JSON.stringify(data, null, 2) + '</div>';
            };
        }
    </script>
</body>
</html>
`)
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

// HandleCostMonthly returns cost data
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

	response := types.CostResponse{
		Team:         team,
		MonthlySpend: currentSpend,
		Budget:       budget,
		Utilization:  utilization,
		LastUpdated:  time.Now().Format(time.RFC3339),
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
