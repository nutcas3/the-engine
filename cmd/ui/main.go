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
	http.HandleFunc("/api/cost/estimate", h.HandleCostEstimate)
	http.HandleFunc("/api/health/status", h.HandleHealth)
	http.HandleFunc("/api/stream", h.HandleSSE)
	http.HandleFunc("/api/swagger", h.HandleSwagger)

	// Alert and cleanup endpoints
	http.HandleFunc("/api/alerts", h.HandleAlerts)
	http.HandleFunc("/api/cleanup/policies", h.HandleCleanupPolicies)
	http.HandleFunc("/api/cron/jobs", h.HandleCronJobs)
	http.HandleFunc("/api/nuke/operations", h.HandleNukeOperations)
	http.HandleFunc("/api/cleanup/shutdown", h.HandleManualShutdown)
	http.HandleFunc("/api/nuke/environment", h.HandleManualNuke)

	// Security configuration endpoints
	http.HandleFunc("/api/security/config", h.HandleSecurityConfigForm)
	http.HandleFunc("/api/security/docs", h.HandleSecurityDocs)
	http.HandleFunc("/api/security/save", h.HandleSecurityConfigSave)

	// Start server
	fmt.Println("Sovereign Engine UI Backend starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
