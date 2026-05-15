package types

import "time"

// Deployment represents a deployment in the system
type Deployment struct {
	ID        string    `json:"id"`
	Provider  string    `json:"provider"`
	Tier      string    `json:"tier"`
	Region    string    `json:"region"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Composition represents a Crossplane composition
type Composition struct {
	Name      string            `json:"name"`
	Provider  string            `json:"provider"`
	Type      string            `json:"type"`
	Labels    map[string]string `json:"labels"`
	CreatedAt string            `json:"created_at"`
}

// CostResponse represents cost information
type CostResponse struct {
	Team         string  `json:"team"`
	MonthlySpend float64 `json:"monthly_spend"`
	Budget       float64 `json:"budget"`
	Utilization  float64 `json:"utilization"`
	LastUpdated  string  `json:"last_updated"`
}

// HealthResponse represents health status
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}
