package nuke

import (
	"context"
	"sync"
	"time"
)

type NukeOperation struct {
	ID          string          `json:"id"`
	Environment string          `json:"environment"`
	Provider    string          `json:"provider"`
	Status      OperationStatus `json:"status"`
	StartedAt   time.Time       `json:"started_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	Resources   []string        `json:"resources"`
	Errors      []string        `json:"errors,omitempty"`
}

type OperationStatus string

const (
	StatusPending   OperationStatus = "pending"
	StatusRunning   OperationStatus = "running"
	StatusCompleted OperationStatus = "completed"
	StatusFailed    OperationStatus = "failed"
	StatusCancelled OperationStatus = "cancelled"
)

type NukeManager struct {
	operations map[string]*NukeOperation
	mu         sync.RWMutex
	providers  map[string]NukeProvider
}

type NukeProvider interface {
	ListResources(ctx context.Context, environment string) ([]string, error)
	DeleteResource(ctx context.Context, resourceID string) error
	ValidateNuke(ctx context.Context, environment string) error
}
