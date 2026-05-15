package cleanup

import (
	"context"
	"sync"
	"time"
)

// EnvironmentType represents different environment types
type EnvironmentType string

const (
	EnvironmentDev     EnvironmentType = "dev"
	EnvironmentTest    EnvironmentType = "test"
	EnvironmentStaging EnvironmentType = "staging"
	EnvironmentProd    EnvironmentType = "prod"
)

// CleanupManager manages automated cleanup operations
type CleanupManager struct {
	policies     map[string]*CleanupPolicy
	mu           sync.RWMutex
	alertChan    chan string
	shutdownChan chan struct{}
	providers    map[string]ResourceCleaner
}

// CleanupPolicy defines cleanup rules for environments
type CleanupPolicy struct {
	Name            string          `json:"name"`
	Environment     EnvironmentType `json:"environment"`
	AutoShutdown    bool            `json:"auto_shutdown"`
	ShutdownAfter   time.Duration   `json:"shutdown_after"`
	NukeAfter       time.Duration   `json:"nuke_after"`
	Enabled         bool            `json:"enabled"`
	ExcludePatterns []string        `json:"exclude_patterns"`
}

// ResourceMetadata contains metadata about a resource
type ResourceMetadata struct {
	ID          string
	CreatedAt   time.Time
	Environment string
	Tags        map[string]string
}

// ResourceCleaner defines interface for cleaning resources from providers
type ResourceCleaner interface {
	Shutdown(ctx context.Context, resourceID string) error
	Nuke(ctx context.Context, environment string) error
	ListResources(ctx context.Context, environment string) ([]string, error)
	GetResourceMetadata(ctx context.Context, resourceID string) (*ResourceMetadata, error)
}
