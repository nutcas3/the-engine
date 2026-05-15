package cleanup

import (
	"context"
	"testing"
	"time"
)

func TestCleanupManager_checkTestsComplete(t *testing.T) {
	cm := &CleanupManager{}
	result := cm.checkTestsComplete(context.Background(), "test-policy")
	if result {
		t.Error("Expected false when GITHUB_TOKEN not set")
	}
}

func TestCleanupManager_checkNoActivePipelines(t *testing.T) {
	cm := &CleanupManager{}
	result := cm.checkNoActivePipelines(context.Background(), "test-policy")
	if result {
		t.Error("Expected false when ArgoCD credentials not configured")
	}
}

func TestCleanupManager_checkEnvironmentIdle(t *testing.T) {
	cm := &CleanupManager{}
	result := cm.checkEnvironmentIdle(context.Background(), "test-policy", 30*time.Minute)
	if result {
		t.Error("Expected false when PROMETHEUS_URL not configured")
	}
}

func TestCleanupManager_checkNoUserActivity(t *testing.T) {
	cm := &CleanupManager{}
	result := cm.checkNoUserActivity(context.Background(), "test-policy", 30*time.Minute)
	if result {
		t.Error("Expected false when no activity check configured")
	}
}
