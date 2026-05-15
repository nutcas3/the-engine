package cleanup

import (
	"context"
	"fmt"
	"log"
	"time"
)

func NewCleanupManager() *CleanupManager {
	return &CleanupManager{
		policies:     make(map[string]*CleanupPolicy),
		alertChan:    make(chan string, 100),
		shutdownChan: make(chan struct{}),
		providers:    make(map[string]ResourceCleaner),
	}
}

func (cm *CleanupManager) RegisterProvider(provider string, cleaner ResourceCleaner) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.providers[provider] = cleaner
}

func (cm *CleanupManager) Start(ctx context.Context) {
	go cm.processCleanup(ctx)
	go cm.monitorResources(ctx)
}

func (cm *CleanupManager) Stop() {
	close(cm.shutdownChan)
}

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
				if cm.shouldShutdown(ctx, resourceID, policy) {
					log.Printf("Shutting down resource %s in %s environment", resourceID, policy.Environment)
					if err := cleaner.Shutdown(ctx, resourceID); err != nil {
						log.Printf("Failed to shutdown resource %s: %v", resourceID, err)
					}
				}

				// Check if environment should be nuked
				if cm.shouldNuke(ctx, policy) {
					log.Printf("Nuking %s environment", policy.Environment)
					if err := cleaner.Nuke(ctx, string(policy.Environment)); err != nil {
						log.Printf("Failed to nuke environment %s: %v", policy.Environment, err)
					}
				}
			}
		}
	}
}

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

			for _, resourceID := range resources {
				if cm.shouldExclude(resourceID, policy.ExcludePatterns) {
					continue
				}

				metadata, err := cleaner.GetResourceMetadata(ctx, resourceID)
				if err != nil {
					log.Printf("Failed to get metadata for resource %s: %v", resourceID, err)
					continue
				}

				resourceAge := time.Since(metadata.CreatedAt)
				if resourceAge > policy.ShutdownAfter {
					log.Printf("Resource %s in %s has exceeded TTL (%v > %v), shutting down",
						resourceID, policy.Environment, resourceAge, policy.ShutdownAfter)

					if err := cleaner.Shutdown(ctx, resourceID); err != nil {
						log.Printf("Failed to shutdown resource %s: %v", resourceID, err)
					} else {
						cm.alertChan <- fmt.Sprintf("Resource %s shutdown due to TTL expiration", resourceID)
					}
				}
			}
		}
	}
}

func (cm *CleanupManager) GetPolicy(name string) (*CleanupPolicy, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	policy, exists := cm.policies[name]
	if !exists {
		return nil, fmt.Errorf("policy not found")
	}
	return policy, nil
}

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

func (cm *CleanupManager) ManualShutdown(ctx context.Context, provider, resourceID string) error {
	cm.mu.RLock()
	cleaner, exists := cm.providers[provider]
	cm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("provider not found")
	}

	return cleaner.Shutdown(ctx, resourceID)
}

func (cm *CleanupManager) ManualNuke(ctx context.Context, provider, environment string) error {
	cm.mu.RLock()
	cleaner, exists := cm.providers[provider]
	cm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("provider not found")
	}

	return cleaner.Nuke(ctx, environment)
}
