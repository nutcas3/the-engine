package cleanup

import (
	"context"
	"log"
	"os"
	"strings"
	"time"
)

// checkEnvironmentIdle evaluates whether an environment has been idle for at least the given duration
// by delegating to the configured monitoring backend. Supported backends are Prometheus, Datadog, and
// CloudWatch, controlled via the MONITORING_SYSTEM environment variable.
func (cm *CleanupManager) checkEnvironmentIdle(ctx context.Context, environmentName string, idleThreshold time.Duration) bool {
	system := strings.ToLower(strings.TrimSpace(os.Getenv("MONITORING_SYSTEM")))
	if system == "" {
		system = "prometheus"
	}

	switch system {
	case "prometheus":
		return cm.checkPrometheusIdle(ctx, environmentName, idleThreshold)
	case "datadog":
		return cm.checkDatadogIdle(ctx, environmentName, idleThreshold)
	case "cloudwatch":
		return cm.checkCloudWatchIdle(ctx, environmentName, idleThreshold)
	default:
		log.Printf("Unknown MONITORING_SYSTEM '%s'; treating environment %s as active", system, environmentName)
		return false
	}
}

// checkNoUserActivity leverages the same monitoring backend logic as checkEnvironmentIdle.
func (cm *CleanupManager) checkNoUserActivity(ctx context.Context, environmentName string, threshold time.Duration) bool {
	return cm.checkEnvironmentIdle(ctx, environmentName, threshold)
}
