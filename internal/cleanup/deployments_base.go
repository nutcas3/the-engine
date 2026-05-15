package cleanup

import (
	"context"
	"log"
	"os"
	"strings"
	"time"
)

func (cm *CleanupManager) checkNoRecentDeployments(ctx context.Context, environmentName string, threshold time.Duration) bool {
	system := strings.ToLower(strings.TrimSpace(os.Getenv("DEPLOYMENT_SYSTEM")))
	if system == "" {
		system = "argocd"
	}

	switch system {
	case "argocd":
		return cm.checkArgoCDDeployments(ctx, environmentName, threshold)
	case "flux":
		return cm.checkFluxDeployments(ctx, environmentName, threshold)
	case "kubernetes":
		return cm.checkKubernetesDeployments(ctx, environmentName, threshold)
	default:
		log.Printf("Unknown DEPLOYMENT_SYSTEM '%s'; assuming deployments are recent for environment %s", system, environmentName)
		return false
	}
}
