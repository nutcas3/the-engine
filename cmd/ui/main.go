package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"the-engine/internal/finops"

	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Deployment struct {
	ID        string    `json:"id"`
	Provider  string    `json:"provider"`
	Tier      string    `json:"tier"`
	Region    string    `json:"region"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CostResponse struct {
	Team         string  `json:"team"`
	MonthlySpend float64 `json:"monthly_spend"`
	Budget       float64 `json:"budget"`
	Utilization  float64 `json:"utilization"`
	LastUpdated  string  `json:"last_updated"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

type Composition struct {
	Name      string            `json:"name"`
	Provider  string            `json:"provider"`
	Type      string            `json:"type"`
	Labels    map[string]string `json:"labels"`
	CreatedAt string            `json:"created_at"`
}

var k8sClient *kubernetes.Clientset
var dynamicClient dynamic.Interface

func main() {
	// Initialize Kubernetes client
	var err error
	k8sClient, err = getK8sClient()
	dynamicClient = nil
	if err != nil {
		log.Printf("Warning: Could not connect to Kubernetes: %v (using mock data)", err)
		k8sClient = nil
	} else {
		log.Println("Connected to Kubernetes successfully")
	}

	// Setup routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/deployments", handleDeployments)
	http.HandleFunc("/api/compositions", handleCompositions)
	http.HandleFunc("/api/cost/monthly", handleCostMonthly)
	http.HandleFunc("/api/health/status", handleHealth)
	http.HandleFunc("/api/stream", handleSSE)

	// Start server
	fmt.Println("Sovereign Engine UI Backend starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getK8sClient() (*kubernetes.Clientset, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Sovereign Engine Dashboard</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        body { font-family: Arial, sans-serif; max-width: 1200px; margin: 0 auto; padding: 20px; }
        .section { margin: 20px 0; padding: 20px; border: 1px solid #ddd; border-radius: 8px; }
        .card { margin: 10px 0; padding: 15px; background: #f5f5f5; border-radius: 4px; }
        button { padding: 10px 20px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #0056b3; }
        select { padding: 10px; margin: 10px 0; min-width: 200px; }
    </style>
</head>
<body>
    <h1>Sovereign Engine Dashboard</h1>
    
    <div class="section">
        <h2>Composition Selection</h2>
        <div id="composition-selector">
            <select id="composition-select">
                <option value="">Select Composition</option>
            </select>
            <button hx-get="/api/compositions" hx-target="#composition-select" hx-swap="innerHTML">Load Compositions</button>
        </div>
    </div>
    
    <div class="section">
        <h2>Deployments</h2>
        <button hx-get="/api/deployments" hx-target="#deployments" hx-swap="innerHTML">Load Deployments</button>
        <div id="deployments"></div>
    </div>
    
    <div class="section">
        <h2>Cost Data</h2>
        <input type="text" id="team-input" placeholder="Enter team name" value="platform">
        <button hx-get="/api/cost/monthly?team=platform" hx-target="#cost" hx-swap="innerHTML" hx-include="#team-input">Load Cost Data</button>
        <div id="cost"></div>
    </div>
    
    <div class="section">
        <h2>Health Status</h2>
        <button hx-get="/api/health/status" hx-target="#health" hx-swap="innerHTML">Check Health</button>
        <div id="health"></div>
    </div>
    
    <div class="section">
        <h2>Real-time Updates</h2>
        <button onclick="connectSSE()">Connect to SSE Stream</button>
        <div id="sse-output"></div>
    </div>
    
    <script>
        function connectSSE() {
            const eventSource = new EventSource('/api/stream');
            eventSource.onmessage = function(event) {
                const data = JSON.parse(event.data);
                document.getElementById('sse-output').innerHTML += '<div class="card">' + JSON.stringify(data, null, 2) + '</div>';
            };
        }
    </script>
</body>
</html>
`)
}

func handleDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Try to get real deployments from Kubernetes if available
	if k8sClient != nil {
		realDeployments, err := getK8sDeployments()
		if err == nil && len(realDeployments) > 0 {
			json.NewEncoder(w).Encode(realDeployments)
			return
		}
	}

	// Fallback to empty array if no data
	json.NewEncoder(w).Encode([]Deployment{})
}

func handleCompositions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	compositions := []Composition{
		{Name: "aws-compute", Provider: "aws", Type: "compute", Labels: map[string]string{"provider": "aws", "engine.io/composition": "compute"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "azure-compute", Provider: "azure", Type: "compute", Labels: map[string]string{"provider": "azure", "engine.io/composition": "compute"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "gcp-compute", Provider: "gcp", Type: "compute", Labels: map[string]string{"provider": "gcp", "engine.io/composition": "compute"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "hetzner-compute", Provider: "hetzner", Type: "compute", Labels: map[string]string{"provider": "hetzner", "engine.io/composition": "compute"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "aws-networking", Provider: "aws", Type: "networking", Labels: map[string]string{"provider": "aws", "engine.io/composition": "networking"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "aws-loadbalancer", Provider: "aws", Type: "loadbalancer", Labels: map[string]string{"provider": "aws", "engine.io/composition": "loadbalancer"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "aws-storage", Provider: "aws", Type: "storage", Labels: map[string]string{"provider": "aws", "engine.io/composition": "storage"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "aws-database", Provider: "aws", Type: "database", Labels: map[string]string{"provider": "aws", "engine.io/composition": "database"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "aws-monitoring", Provider: "aws", Type: "monitoring", Labels: map[string]string{"provider": "aws", "engine.io/composition": "monitoring"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "shared-monitoring-stack", Provider: "shared", Type: "monitoring", Labels: map[string]string{"provider": "shared", "engine.io/composition": "monitoring"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "shared-database-compute", Provider: "shared", Type: "database", Labels: map[string]string{"provider": "shared", "engine.io/composition": "database"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "shared-vault", Provider: "shared", Type: "secrets", Labels: map[string]string{"provider": "shared", "engine.io/composition": "secrets"}, CreatedAt: time.Now().Format(time.RFC3339)},
		{Name: "shared-dns-server", Provider: "shared", Type: "dns", Labels: map[string]string{"provider": "shared", "engine.io/composition": "dns"}, CreatedAt: time.Now().Format(time.RFC3339)},
	}

	json.NewEncoder(w).Encode(compositions)
}

func getK8sDeployments() ([]Deployment, error) {
	if k8sClient == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	// Try to get XCompute resources from Crossplane namespace
	// This is a simplified version - in production you'd use the Crossplane CRD client
	namespaces, err := k8sClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var k8sDeployments []Deployment
	for _, ns := range namespaces.Items {
		pods, err := k8sClient.CoreV1().Pods(ns.Name).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, pod := range pods.Items {
			provider := pod.Labels["engine.io/provider"]
			tier := pod.Labels["engine.io/tier"]
			region := pod.Labels["topology.kubernetes.io/zone"]

			if provider == "" {
				provider = "unknown"
			}
			if tier == "" {
				tier = "unknown"
			}
			if region == "" {
				region = "default"
			}

			status := "running"
			if pod.Status.Phase != "Running" {
				status = string(pod.Status.Phase)
			}

			k8sDeployments = append(k8sDeployments, Deployment{
				ID:        pod.Name,
				Provider:  provider,
				Tier:      tier,
				Region:    region,
				Status:    status,
				CreatedAt: pod.CreationTimestamp.Time,
			})
		}
	}

	return k8sDeployments, nil
}

func handleCostMonthly(w http.ResponseWriter, r *http.Request) {
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

	response := CostResponse{
		Team:         team,
		MonthlySpend: currentSpend,
		Budget:       budget,
		Utilization:  utilization,
		LastUpdated:  time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
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
			health := HealthResponse{
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
