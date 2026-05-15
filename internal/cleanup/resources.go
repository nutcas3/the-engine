package cleanup

import (
	"context"
	"log"
	"strings"
)

// checkNoActiveResources verifies that no active resources remain for the given environment across all
// registered providers. It consults the providers' ResourceCleaner interfaces to inspect resource metadata.
func (cm *CleanupManager) checkNoActiveResources(ctx context.Context, environmentName string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if len(cm.providers) == 0 {
		log.Printf("No resource providers registered; treating environment %s as inactive", environmentName)
		return true
	}

	targetEnv := strings.ToLower(environmentName)

	for provider, cleaner := range cm.providers {
		resources, err := cleaner.ListResources(ctx, targetEnv)
		if err != nil {
			log.Printf("Failed to list resources for provider %s in environment %s: %v", provider, targetEnv, err)
			return false
		}

		for _, resourceID := range resources {
			metadata, err := cleaner.GetResourceMetadata(ctx, resourceID)
			if err != nil {
				log.Printf("Failed to fetch metadata for resource %s from provider %s: %v", resourceID, provider, err)
				return false
			}

			if resourceIsActive(metadata) {
				log.Printf("Resource %s from provider %s remains active in environment %s", resourceID, provider, targetEnv)
				return false
			}
		}
	}

	return true
}

func resourceIsActive(metadata *ResourceMetadata) bool {
	if metadata == nil {
		return true
	}

	if metadata.Tags != nil {
		if status, ok := metadata.Tags["status"]; ok {
			switch strings.ToLower(status) {
			case "terminated", "stopped", "shut-down", "succeeded", "completed":
				return false
			case "":
				return true
			default:
				return true
			}
		}
	}

	return true
}
