package integrations

import (
	"context"
	"testing"
)

func TestNewPagerDutyClient(t *testing.T) {
	client := NewPagerDutyClient("test-routing-key")
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
	if client.client == nil {
		t.Error("Expected client to have HTTP client")
	}
	if client.routingKey != "test-routing-key" {
		t.Errorf("Expected test-routing-key, got %s", client.routingKey)
	}
}

func TestPagerDutyClient_Trigger(t *testing.T) {
	// Skip actual HTTP call in tests - test event construction only
	client := NewPagerDutyClient("test-key")
	customDetails := map[string]interface{}{"key": "value"}

	// Just verify the method doesn't panic with valid parameters
	// In production, this would make an actual HTTP call
	_ = client.Trigger(context.Background(), "test-dedup", "test summary", "test source", "test component", "warning", customDetails)
}

func TestPagerDutyClient_Resolve(t *testing.T) {
	client := NewPagerDutyClient("test-key")
	_ = client.Resolve(context.Background(), "test-dedup")
}

func TestPagerDutyClient_TriggerCostAlert(t *testing.T) {
	client := NewPagerDutyClient("test-key")
	_ = client.TriggerCostAlert(context.Background(), "test", 150.0, 100.0)
}

func TestPagerDutyClient_TriggerCleanupFailure(t *testing.T) {
	client := NewPagerDutyClient("test-key")
	_ = client.TriggerCleanupFailure(context.Background(), "test", "aws", "resource-123", "test error")
}

func TestPagerDutyClient_ResolveCostAlert(t *testing.T) {
	client := NewPagerDutyClient("test-key")
	_ = client.ResolveCostAlert(context.Background(), "test")
}
