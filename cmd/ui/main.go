package main

import (
	"fmt"
	"log"
	"net/http"

	"the-engine/internal/handlers"
	"the-engine/internal/kubernetes"
)

func main() {
	// Initialize Kubernetes client
	k8sClient := kubernetes.NewClientOrMock()

	// Initialize handlers
	h := handlers.NewHandlers(k8sClient)

	// Setup routes
	http.HandleFunc("/", h.HandleIndex)
	http.HandleFunc("/api/deployments", h.HandleDeployments)
	http.HandleFunc("/api/compositions", h.HandleCompositions)
	http.HandleFunc("/api/cost/monthly", h.HandleCostMonthly)
	http.HandleFunc("/api/health/status", h.HandleHealth)
	http.HandleFunc("/api/stream", h.HandleSSE)

	// Start server
	fmt.Println("Sovereign Engine UI Backend starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
