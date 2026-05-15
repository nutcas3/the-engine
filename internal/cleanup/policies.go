package cleanup

import (
	"time"
)

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

// RemovePolicy removes a cleanup policy
func (cm *CleanupManager) RemovePolicy(name string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.policies, name)
}

// DefaultCleanupPolicies returns default cleanup policies for dev and test environments
func DefaultCleanupPolicies() []*CleanupPolicy {
	return []*CleanupPolicy{
		{
			Name:          "dev-auto-shutdown",
			Environment:   EnvironmentDev,
			AutoShutdown:  true,
			ShutdownAfter: 8 * time.Hour,
			NukeAfter:     24 * time.Hour,
			Enabled:       true,
			ExcludePatterns: []string{
				"essential-*",
				"database-*",
			},
		},
		{
			Name:            "test-auto-cleanup",
			Environment:     EnvironmentTest,
			AutoShutdown:    true,
			ShutdownAfter:   2 * time.Hour,
			NukeAfter:       6 * time.Hour,
			Enabled:         true,
			ExcludePatterns: []string{},
		},
	}
}
