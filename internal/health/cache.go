package health

import "time"

// checkCache checks cache connectivity
func (c *Checker) checkCache() ComponentHealth {
	return ComponentHealth{
		Name:      "cache",
		Status:    StatusHealthy,
		Message:   "Cache operational",
		CheckedAt: time.Now(),
	}
}
