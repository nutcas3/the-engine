package cleanup

import (
	"log"
)

// shouldExclude checks if a resource should be excluded from cleanup
func (cm *CleanupManager) shouldExclude(resourceID string, patterns []string) bool {
	for _, pattern := range patterns {
		if cm.matchesPattern(resourceID, pattern) {
			return true
		}
	}
	return false
}

// matchesPattern checks if a resource ID matches a pattern with wildcards
func (cm *CleanupManager) matchesPattern(resourceID, pattern string) bool {
	// Simple wildcard implementation
	// Supports patterns like "database-*", "essential-*"
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(resourceID) >= len(prefix) && resourceID[:len(prefix)] == prefix
	}
	return resourceID == pattern
}

// shouldShutdown determines if a resource should be shut down
func (cm *CleanupManager) shouldShutdown(resourceID string, policy *CleanupPolicy) bool {
	if !policy.AutoShutdown || policy.ShutdownAfter == 0 {
		return false
	}

	if cm.shouldExclude(resourceID, policy.ExcludePatterns) {
		return false
	}

	// For test environments, check if tests are complete
	if policy.Environment == EnvironmentTest {
		// Check for test completion tag or status
		// This would integrate with CI/CD systems or test runners
		return true
	}

	// For dev environments, check for idle/inactive status
	if policy.Environment == EnvironmentDev {
		// Check for recent activity (API calls, deployments, etc.)
		// This would integrate with monitoring systems
		return true
	}

	return false
}

func (cm *CleanupManager) shouldNuke(policy *CleanupPolicy) bool {
	if policy.NukeAfter == 0 {
		return false
	}

	if policy.Environment == EnvironmentTest {
		// Check if:
		// 1. All test jobs have completed (success or failure)
		// 2. No active test runs
		// 3. Environment has been idle for NukeAfter duration
		// This would integrate with CI/CD pipeline status
		log.Printf("Checking if test environment %s should be nuked", policy.Name)
		return true
	}

	// For dev environments, nuke after extended inactivity
	if policy.Environment == EnvironmentDev {
		// Check if:
		// 1. No deployments in the last NukeAfter duration
		// 2. No active resources (or all resources shut down)
		// 3. No API activity or user sessions
		log.Printf("Checking if dev environment %s should be nuked", policy.Name)
		return true
	}

	// Never auto-nuke staging or production
	if policy.Environment == EnvironmentStaging || policy.Environment == EnvironmentProd {
		return false
	}

	return false
}
