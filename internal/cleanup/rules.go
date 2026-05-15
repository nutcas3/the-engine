package cleanup

import (
	"context"
	"time"
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

// checkNoUserActivity checks if there has been no user activity
func (cm *CleanupManager) checkNoUserActivity(ctx context.Context, environmentName string, threshold time.Duration) bool {
	// In a real implementation, this would:
	// 1. Check API logs for recent requests
	// 2. Check authentication logs for recent logins
	// 3. Check application logs for user activity

	_ = environmentName
	_ = ctx
	_ = threshold

	// In production, integrate with your logging/monitoring system:
	// - CloudWatch Logs: Query for recent API calls
	// - Datadog: Check recent events
	// - Prometheus: Query request metrics
	// - Custom: Check your application logs

	// For now, return true to enable the nuke logic
	return true
}
