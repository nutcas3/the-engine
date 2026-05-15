package nuke

import (
	"context"
	"fmt"
	"log"
	"time"
)

func NewNukeManager() *NukeManager {
	return &NukeManager{
		operations: make(map[string]*NukeOperation),
		providers:  make(map[string]NukeProvider),
	}
}

func (nm *NukeManager) RegisterProvider(provider string, nukeProvider NukeProvider) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	nm.providers[provider] = nukeProvider
}

func (nm *NukeManager) NukeEnvironment(ctx context.Context, environment, provider string) (*NukeOperation, error) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	
	nukeProvider, exists := nm.providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", provider)
	}
	
	// Validate nuke operation
	if err := nukeProvider.ValidateNuke(ctx, environment); err != nil {
		return nil, fmt.Errorf("nuke validation failed: %w", err)
	}
	
	// Create operation
	operation := &NukeOperation{
		ID:          fmt.Sprintf("nuke-%d", time.Now().UnixNano()),
		Environment: environment,
		Provider:    provider,
		Status:      StatusPending,
		StartedAt:   time.Now(),
		Resources:   []string{},
		Errors:      []string{},
	}
	
	nm.operations[operation.ID] = operation
	
	// Execute nuke in background
	go nm.executeNuke(ctx, operation, nukeProvider)
	
	return operation, nil
}

func (nm *NukeManager) executeNuke(ctx context.Context, operation *NukeOperation, provider NukeProvider) {
	nm.mu.Lock()
	operation.Status = StatusRunning
	nm.mu.Unlock()
	
	log.Printf("Starting nuke operation %s for environment %s", operation.ID, operation.Environment)
	
	// List all resources
	resources, err := provider.ListResources(ctx, operation.Environment)
	if err != nil {
		nm.mu.Lock()
		operation.Status = StatusFailed
		operation.Errors = append(operation.Errors, err.Error())
		now := time.Now()
		operation.CompletedAt = &now
		nm.mu.Unlock()
		return
	}
	
	operation.Resources = resources
	
	// Delete each resource
	for _, resourceID := range resources {
		select {
		case <-ctx.Done():
			nm.mu.Lock()
			operation.Status = StatusCancelled
			now := time.Now()
			operation.CompletedAt = &now
			nm.mu.Unlock()
			return
		default:
			err := provider.DeleteResource(ctx, resourceID)
			if err != nil {
				nm.mu.Lock()
				operation.Errors = append(operation.Errors, 
					fmt.Sprintf("Failed to delete %s: %v", resourceID, err))
				nm.mu.Unlock()
				log.Printf("Failed to delete resource %s: %v", resourceID, err)
			} else {
				log.Printf("Deleted resource %s", resourceID)
			}
		}
	}
	
	// Mark as completed
	nm.mu.Lock()
	now := time.Now()
	operation.CompletedAt = &now
	if len(operation.Errors) > 0 {
		operation.Status = StatusFailed
	} else {
		operation.Status = StatusCompleted
	}
	nm.mu.Unlock()
	
	log.Printf("Nuke operation %s completed with status %s", operation.ID, operation.Status)
}

func (nm *NukeManager) GetOperation(operationID string) (*NukeOperation, error) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	
	operation, exists := nm.operations[operationID]
	if !exists {
		return nil, fmt.Errorf("operation not found")
	}
	return operation, nil
}

func (nm *NukeManager) GetOperations() []*NukeOperation {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	
	operations := make([]*NukeOperation, 0, len(nm.operations))
	for _, operation := range nm.operations {
		operations = append(operations, operation)
	}
	return operations
}

func (nm *NukeManager) CancelOperation(operationID string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	
	operation, exists := nm.operations[operationID]
	if !exists {
		return fmt.Errorf("operation not found")
	}
	
	if operation.Status != StatusPending && operation.Status != StatusRunning {
		return fmt.Errorf("operation cannot be cancelled in current state: %s", operation.Status)
	}
	
	operation.Status = StatusCancelled
	now := time.Now()
	operation.CompletedAt = &now
	
	return nil
}

func (nm *NukeManager) DryRun(ctx context.Context, environment, provider string) ([]string, error) {
	nm.mu.RLock()
	nukeProvider, exists := nm.providers[provider]
	nm.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", provider)
	}
	
	return nukeProvider.ListResources(ctx, environment)
}

func (nm *NukeManager) ValidateEnvironment(ctx context.Context, environment, provider string) error {
	nm.mu.RLock()
	nukeProvider, exists := nm.providers[provider]
	nm.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("provider not found: %s", provider)
	}
	
	return nukeProvider.ValidateNuke(ctx, environment)
}
