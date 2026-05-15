package health

import "time"

// checkDatabase checks database connectivity
func (c *Checker) checkDatabase() ComponentHealth {
	return ComponentHealth{
		Name:      "database",
		Status:    StatusDegraded,
		Message:   "Database not configured",
		CheckedAt: time.Now(),
	}
}
