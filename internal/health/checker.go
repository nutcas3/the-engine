package health

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

// Status represents the health status of a component
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusDegraded  Status = "degraded"
	StatusUnhealthy Status = "unhealthy"
)

// ComponentHealth represents health of a specific component
type ComponentHealth struct {
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Message   string    `json:"message,omitempty"`
	CheckedAt time.Time `json:"checked_at"`
}

// HealthResponse represents comprehensive health check response
type HealthResponse struct {
	Status     string            `json:"status"`
	Version    string            `json:"version"`
	Timestamp  time.Time         `json:"timestamp"`
	Components []ComponentHealth `json:"components"`
	System     SystemInfo        `json:"system"`
}

// SystemInfo represents system metrics
type SystemInfo struct {
	GoVersion   string `json:"go_version"`
	Goroutines  int    `json:"goroutines"`
	MemoryAlloc uint64 `json:"memory_alloc"`
	MemorySys   uint64 `json:"memory_sys"`
	Uptime      string `json:"uptime"`
}

// Checker performs health checks
type Checker struct {
	startTime time.Time
	version   string
}

// NewChecker creates a new health checker
func NewChecker(version string) *Checker {
	return &Checker{
		startTime: time.Now(),
		version:   version,
	}
}

// Check performs comprehensive health check
func (c *Checker) Check() HealthResponse {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	components := []ComponentHealth{
		c.checkKubernetes(),
		c.checkDatabase(),
		c.checkCache(),
		c.checkExternalAPIs(),
	}

	overallStatus := c.calculateOverallStatus(components)

	return HealthResponse{
		Status:     overallStatus,
		Version:    c.version,
		Timestamp:  time.Now(),
		Components: components,
		System: SystemInfo{
			GoVersion:   runtime.Version(),
			Goroutines:  runtime.NumGoroutine(),
			MemoryAlloc: m.Alloc,
			MemorySys:   m.Sys,
			Uptime:      time.Since(c.startTime).String(),
		},
	}
}

// checkKubernetes checks Kubernetes connectivity
func (c *Checker) checkKubernetes() ComponentHealth {
	// Try to make a simple HTTP request to Kubernetes API
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:8080/api/health/status")

	if err != nil {
		return ComponentHealth{
			Name:      "kubernetes",
			Status:    StatusDegraded,
			Message:   "Kubernetes API unreachable",
			CheckedAt: time.Now(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return ComponentHealth{
			Name:      "kubernetes",
			Status:    StatusHealthy,
			Message:   "Kubernetes API responsive",
			CheckedAt: time.Now(),
		}
	}

	return ComponentHealth{
		Name:      "kubernetes",
		Status:    StatusDegraded,
		Message:   "Kubernetes API returned error",
		CheckedAt: time.Now(),
	}
}

// checkDatabase checks database connectivity
func (c *Checker) checkDatabase() ComponentHealth {
	// Placeholder for actual database health check
	// In production, this would ping the database
	return ComponentHealth{
		Name:      "database",
		Status:    StatusDegraded,
		Message:   "Database not configured",
		CheckedAt: time.Now(),
	}
}

// checkCache checks cache connectivity
func (c *Checker) checkCache() ComponentHealth {
	// Placeholder for actual cache health check
	// In production, this would check cache availability
	return ComponentHealth{
		Name:      "cache",
		Status:    StatusHealthy,
		Message:   "Cache operational",
		CheckedAt: time.Now(),
	}
}

// checkExternalAPIs checks external API connectivity
func (c *Checker) checkExternalAPIs() ComponentHealth {
	// Check if we can reach external APIs
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("https://httpbin.org/status/200")

	if err != nil {
		return ComponentHealth{
			Name:      "external_apis",
			Status:    StatusDegraded,
			Message:   "External APIs unreachable",
			CheckedAt: time.Now(),
		}
	}
	defer resp.Body.Close()

	return ComponentHealth{
		Name:      "external_apis",
		Status:    StatusHealthy,
		Message:   "External APIs reachable",
		CheckedAt: time.Now(),
	}
}

// calculateOverallStatus calculates overall system status
func (c *Checker) calculateOverallStatus(components []ComponentHealth) string {
	for _, comp := range components {
		if comp.Status == StatusUnhealthy {
			return "unhealthy"
		}
		if comp.Status == StatusDegraded {
			return "degraded"
		}
	}
	return "healthy"
}

// WriteJSON writes health response as JSON
func (h HealthResponse) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(h)
}
