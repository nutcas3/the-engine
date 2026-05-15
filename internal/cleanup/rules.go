package cleanup

import (
	"context"
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
func (cm *CleanupManager) shouldShutdown(ctx context.Context, resourceID string, policy *CleanupPolicy) bool {
	if !policy.AutoShutdown || policy.ShutdownAfter == 0 {
		return false
	}

	if cm.shouldExclude(resourceID, policy.ExcludePatterns) {
		return false
	}

	// For test environments, check if tests are complete
	if policy.Environment == EnvironmentTest {
		// Check for test completion tag or status
		// This integrates with CI/CD systems or test runners
		testsComplete := cm.checkTestsComplete(ctx, policy.Name)
		return testsComplete
	}

	// For dev environments, check for idle/inactive status
	if policy.Environment == EnvironmentDev {
		// Check for recent activity (API calls, deployments, etc.)
		// This integrates with monitoring systems
		isIdle := cm.checkEnvironmentIdle(ctx, policy.Name, policy.ShutdownAfter)
		return isIdle
	}

	return false
}

func (cm *CleanupManager) shouldNuke(ctx context.Context, policy *CleanupPolicy) bool {
	if policy.NukeAfter == 0 {
		return false
	}

	switch policy.Environment {
	case EnvironmentTest:
		testsComplete := cm.checkTestsComplete(ctx, policy.Name)
		noActivePipelines := cm.checkNoActivePipelines(ctx, policy.Name)
		envIdle := cm.checkEnvironmentIdle(ctx, policy.Name, policy.NukeAfter)

		log.Printf("Evaluating nuke for test environment %s: testsComplete=%v noActivePipelines=%v envIdle=%v",
			policy.Name, testsComplete, noActivePipelines, envIdle)

		return testsComplete && noActivePipelines && envIdle

	case EnvironmentDev:
		noRecentDeployments := cm.checkNoRecentDeployments(ctx, policy.Name, policy.NukeAfter)
		noActiveResources := cm.checkNoActiveResources(ctx, policy.Name)
		noUserActivity := cm.checkNoUserActivity(ctx, policy.Name, policy.NukeAfter)

		log.Printf("Evaluating nuke for dev environment %s: noDeployments=%v noResources=%v noUserActivity=%v",
			policy.Name, noRecentDeployments, noActiveResources, noUserActivity)

		return noRecentDeployments && noActiveResources && noUserActivity

	case EnvironmentProd, EnvironmentStaging:
		return false

	default:
		log.Printf("Unknown environment type %s for policy %s; skipping nuke", policy.Environment, policy.Name)
		return false
	}
}
