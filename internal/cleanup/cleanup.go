package cleanup

import (
	"context"
	"fmt"
	"log"
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

// CleanupManager manages automated cleanup operations
type CleanupManager struct {
	policies     map[string]*CleanupPolicy
	mu           sync.RWMutex
	alertChan    chan string
	shutdownChan chan struct{}
	providers    map[string]ResourceCleaner
}

// ResourceCleaner defines interface for cleaning resources from providers
type ResourceCleaner interface {
	Shutdown(ctx context.Context, resourceID string) error
	Nuke(ctx context.Context, environment string) error
	ListResources(ctx context.Context, environment string) ([]string, error)
}

// NewCleanupManager creates a new cleanup manager
func NewCleanupManager() *CleanupManager {
	return &CleanupManager{
		policies:     make(map[string]*CleanupPolicy),
		alertChan:    make(chan string, 100),
		shutdownChan: make(chan struct{}),
		providers:    make(map[string]ResourceCleaner),
	}
}

// RegisterProvider registers a resource cleaner for a provider
func (cm *CleanupManager) RegisterProvider(provider string, cleaner ResourceCleaner) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.providers[provider] = cleaner
}

// Start begins the cleanup manager background processing
func (cm *CleanupManager) Start(ctx context.Context) {
	go cm.processCleanup(ctx)
	go cm.monitorResources(ctx)
}

// Stop stops the cleanup manager
func (cm *CleanupManager) Stop() {
	close(cm.shutdownChan)
}

// processCleanup handles cleanup operations
func (cm *CleanupManager) processCleanup(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-cm.shutdownChan:
			return
		case <-ticker.C:
			cm.checkAndCleanup(ctx)
		}
	}
}

// monitorResources monitors resources for cleanup triggers
func (cm *CleanupManager) monitorResources(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-cm.shutdownChan:
			return
		case <-ticker.C:
			cm.checkResourceTTLs(ctx)
		}
	}
}

// checkAndCleanup performs cleanup based on policies
func (cm *CleanupManager) checkAndCleanup(ctx context.Context) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, policy := range cm.policies {
		if !policy.Enabled || !policy.AutoShutdown {
			continue
		}

		for provider, cleaner := range cm.providers {
			resources, err := cleaner.ListResources(ctx, string(policy.Environment))
			if err != nil {
				log.Printf("Failed to list resources for %s: %v", provider, err)
				continue
			}

			for _, resourceID := range resources {
				if cm.shouldExclude(resourceID, policy.ExcludePatterns) {
					continue
				}

				// Check if resource should be shut down
				if cm.shouldShutdown(resourceID, policy) {
					log.Printf("Shutting down resource %s in %s environment", resourceID, policy.Environment)
					if err := cleaner.Shutdown(ctx, resourceID); err != nil {
						log.Printf("Failed to shutdown resource %s: %v", resourceID, err)
					}
				}

				// Check if environment should be nuked
				if cm.shouldNuke(policy) {
					log.Printf("Nuking %s environment", policy.Environment)
					if err := cleaner.Nuke(ctx, string(policy.Environment)); err != nil {
						log.Printf("Failed to nuke environment %s: %v", policy.Environment, err)
					}
				}
			}
		}
	}
}

// checkResourceTTLs checks if resources have exceeded their TTL
func (cm *CleanupManager) checkResourceTTLs(ctx context.Context) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, policy := range cm.policies {
		if !policy.Enabled || policy.ShutdownAfter == 0 {
			continue
		}

		for provider, cleaner := range cm.providers {
			resources, err := cleaner.ListResources(ctx, string(policy.Environment))
			if err != nil {
				log.Printf("Failed to list resources for %s: %v", provider, err)
				continue
			}

			// TODO: Implement resource age tracking and TTL-based shutdown
			// This would check resource metadata for creation time and compare against policy.ShutdownAfter
			_ = resources // Will be used when resource tracking is implemented
		}
	}
}

// shouldExclude checks if a resource should be excluded from cleanup
func (cm *CleanupManager) shouldExclude(_ string, patterns []string) bool {
	// Simple pattern matching - could be enhanced with regex
	// For now, return false as pattern matching would need actual implementation
	return len(patterns) > 0
}

// shouldShutdown determines if a resource should be shut down
func (cm *CleanupManager) shouldShutdown(resourceID string, policy *CleanupPolicy) bool {
	// This would check resource age, test completion status, etc.
	// For now, return false to be implemented based on actual resource tracking
	return false
}

// shouldNuke determines if an environment should be nuked
func (cm *CleanupManager) shouldNuke(policy *CleanupPolicy) bool {
	// This would check if all tests are complete, environment is idle, etc.
	// For now, return false to be implemented based on actual monitoring
	return false
}

// AddPolicy adds a cleanup policy
func (cm *CleanupManager) AddPolicy(policy *CleanupPolicy) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.policies[policy.Name] = policy
}

// GetPolicies returns all cleanup policies
func (cm *CleanupManager) GetPolicies() []*CleanupPolicy {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	policies := make([]*CleanupPolicy, 0, len(cm.policies))
	for _, policy := range cm.policies {
		policies = append(policies, policy)
	}
	return policies
}

// GetPolicy returns a specific cleanup policy
func (cm *CleanupManager) GetPolicy(name string) (*CleanupPolicy, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	policy, exists := cm.policies[name]
	if !exists {
		return nil, fmt.Errorf("policy not found")
	}
	return policy, nil
}

// DeletePolicy removes a cleanup policy
func (cm *CleanupManager) DeletePolicy(name string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	_, exists := cm.policies[name]
	if !exists {
		return fmt.Errorf("policy not found")
	}

	delete(cm.policies, name)
	return nil
}

// ManualShutdown manually triggers shutdown of a resource
func (cm *CleanupManager) ManualShutdown(ctx context.Context, provider, resourceID string) error {
	cm.mu.RLock()
	cleaner, exists := cm.providers[provider]
	cm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("provider not found")
	}

	return cleaner.Shutdown(ctx, resourceID)
}

// ManualNuke manually triggers nuking of an environment
func (cm *CleanupManager) ManualNuke(ctx context.Context, provider, environment string) error {
	cm.mu.RLock()
	cleaner, exists := cm.providers[provider]
	cm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("provider not found")
	}

	return cleaner.Nuke(ctx, environment)
}

// DefaultCleanupPolicies returns default cleanup policies for common environments
func DefaultCleanupPolicies() []*CleanupPolicy {
	return []*CleanupPolicy{
		{
			Name:            "dev-auto-shutdown",
			Environment:     EnvironmentDev,
			AutoShutdown:    true,
			ShutdownAfter:   8 * time.Hour,  // Shutdown after 8 hours
			NukeAfter:       24 * time.Hour, // Nuke after 24 hours
			Enabled:         true,
			ExcludePatterns: []string{"essential-*", "database-*"},
		},
		{
			Name:            "test-auto-cleanup",
			Environment:     EnvironmentTest,
			AutoShutdown:    true,
			ShutdownAfter:   2 * time.Hour, // Shutdown after 2 hours
			NukeAfter:       6 * time.Hour, // Nuke after 6 hours
			Enabled:         true,
			ExcludePatterns: []string{},
		},
	}
}
