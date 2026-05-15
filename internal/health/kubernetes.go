package health

import (
	"net/http"
	"time"
)

// checkKubernetes checks Kubernetes connectivity
func (c *Checker) checkKubernetes() ComponentHealth {
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
