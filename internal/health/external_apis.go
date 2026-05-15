package health

import (
	"net/http"
	"time"
)

// checkExternalAPIs checks external API connectivity
func (c *Checker) checkExternalAPIs() ComponentHealth {
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
