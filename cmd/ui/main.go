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

	// Setup routes with rate limiting for API endpoints
	http.HandleFunc("/", h.HandleIndex)
	http.Handle("/api/deployments", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleDeployments)))
	http.Handle("/api/compositions", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleCompositions)))
	http.Handle("/api/cost/monthly", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleCostMonthly)))
	http.Handle("/api/cost/estimate", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleCostEstimate)))
	http.Handle("/api/health/status", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleHealth)))
	http.HandleFunc("/api/stream", h.HandleSSE)
	http.Handle("/api/swagger", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleSwagger)))

	// Alert and cleanup endpoints
	http.Handle("/api/alerts", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleAlerts)))
	http.Handle("/api/cleanup/policies", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleCleanupPolicies)))
	http.Handle("/api/cron/jobs", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleCronJobs)))
	http.Handle("/api/nuke/operations", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleNukeOperations)))
	http.Handle("/api/cleanup/shutdown", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleManualShutdown)))
	http.Handle("/api/nuke/environment", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleManualNuke)))

	// Security configuration endpoints
	http.HandleFunc("/api/security/config", h.HandleSecurityConfigForm)
	http.HandleFunc("/api/security/docs", h.HandleSecurityDocs)
	http.Handle("/api/security/config/save", h.RateLimiter.Middleware(http.HandlerFunc(h.HandleSecurityConfigSave)))

	// Start server
	fmt.Println("Sovereign Engine UI Backend starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
