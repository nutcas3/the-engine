package alerts

import (
	"context"
	"testing"
	"time"
)

func TestNewAlertManager(t *testing.T) {
	am := NewAlertManager()
	if am == nil {
		t.Fatal("Expected non-nil alert manager")
	}
	if am.alerts == nil {
		t.Error("Expected alerts map to be initialized")
	}
	if am.alertChan == nil {
		t.Error("Expected alert channel to be initialized")
	}
	if am.shutdownChan == nil {
		t.Error("Expected shutdown channel to be initialized")
	}
}

func TestAlertManager_CreateAlert(t *testing.T) {
	am := NewAlertManager()
	ctx := context.Background()
	am.Start(ctx)
	defer am.Stop()

	alert := am.CreateAlert(AlertTypeCost, SeverityWarning, "test", "resource-1", "Test alert message")
	if alert == nil {
		t.Fatal("Expected non-nil alert")
	}
	if alert.ID == "" {
		t.Error("Expected alert ID to be set")
	}
	if alert.Type != AlertTypeCost {
		t.Errorf("Expected %s, got %s", AlertTypeCost, alert.Type)
	}
	if alert.Severity != SeverityWarning {
		t.Errorf("Expected %s, got %s", SeverityWarning, alert.Severity)
	}
	if alert.Environment != "test" {
		t.Errorf("Expected test, got %s", alert.Environment)
	}
	if alert.Resource != "resource-1" {
		t.Errorf("Expected resource-1, got %s", alert.Resource)
	}
	if alert.Resolved {
		t.Error("Expected alert to be unresolved")
	}
}

func TestAlertManager_GetAlerts(t *testing.T) {
	am := NewAlertManager()
	ctx := context.Background()
	am.Start(ctx)
	defer am.Stop()

	am.CreateAlert(AlertTypeCost, SeverityWarning, "test", "resource-1", "Test alert 1")
	am.CreateAlert(AlertTypeTTL, SeverityCritical, "test", "resource-2", "Test alert 2")

	// Give time for alerts to be processed
	time.Sleep(100 * time.Millisecond)

	alerts := am.GetAlerts()
	if len(alerts) != 2 {
		t.Errorf("Expected 2 alerts, got %d", len(alerts))
	}
}

func TestAlertManager_GetAlertsByEnvironment(t *testing.T) {
	am := NewAlertManager()
	ctx := context.Background()
	am.Start(ctx)
	defer am.Stop()

	am.CreateAlert(AlertTypeCost, SeverityWarning, "test", "resource-1", "Test alert")
	am.CreateAlert(AlertTypeTTL, SeverityWarning, "prod", "resource-2", "Test alert")

	time.Sleep(100 * time.Millisecond)

	testAlerts := am.GetAlertsByEnvironment("test")
	if len(testAlerts) != 1 {
		t.Errorf("Expected 1 alert for test environment, got %d", len(testAlerts))
	}

	prodAlerts := am.GetAlertsByEnvironment("prod")
	if len(prodAlerts) != 1 {
		t.Errorf("Expected 1 alert for prod environment, got %d", len(prodAlerts))
	}
}

func TestAlertManager_GetUnresolvedAlerts(t *testing.T) {
	am := NewAlertManager()
	ctx := context.Background()
	am.Start(ctx)
	defer am.Stop()

	alert1 := am.CreateAlert(AlertTypeCost, SeverityWarning, "test", "resource-1", "Test alert 1")
	am.CreateAlert(AlertTypeTTL, SeverityWarning, "test", "resource-2", "Test alert 2")

	time.Sleep(100 * time.Millisecond)

	am.ResolveAlert(alert1.ID)

	unresolved := am.GetUnresolvedAlerts()
	if len(unresolved) != 1 {
		t.Errorf("Expected 1 unresolved alert, got %d", len(unresolved))
	}
}

func TestAlertManager_ResolveAlert(t *testing.T) {
	am := NewAlertManager()
	ctx := context.Background()
	am.Start(ctx)
	defer am.Stop()

	alert := am.CreateAlert(AlertTypeCost, SeverityWarning, "test", "resource-1", "Test alert")
	time.Sleep(100 * time.Millisecond)

	err := am.ResolveAlert(alert.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	retrievedAlerts := am.GetAlerts()
	if len(retrievedAlerts) != 1 {
		t.Errorf("Expected 1 alert, got %d", len(retrievedAlerts))
	}

	if !retrievedAlerts[0].Resolved {
		t.Error("Expected alert to be resolved")
	}
	if retrievedAlerts[0].ResolvedAt == nil {
		t.Error("Expected ResolvedAt to be set")
	}

	// Resolve non-existent alert
	err = am.ResolveAlert("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent alert")
	}
}

func TestAlertManager_CheckCostAlert(t *testing.T) {
	am := NewAlertManager()
	ctx := context.Background()
	am.Start(ctx)
	defer am.Stop()

	// Test no alert when under threshold
	am.CheckCostAlert("test", 50.0, 100.0)

	// Test warning alert at 80%
	am.CheckCostAlert("test", 85.0, 100.0)

	// Test critical alert at 90%
	am.CheckCostAlert("test", 95.0, 100.0)

	time.Sleep(100 * time.Millisecond)

	alerts := am.GetAlertsByEnvironment("test")
	if len(alerts) < 2 {
		t.Errorf("Expected at least 2 alerts, got %d", len(alerts))
	}
}

func TestAlertManager_CheckTTLAlert(t *testing.T) {
	am := NewAlertManager()
	ctx := context.Background()
	am.Start(ctx)
	defer am.Stop()

	// Test no alert when within TTL
	createdAt := time.Now()
	am.CheckTTLAlert("test", "resource-1", createdAt, 1*time.Hour)

	// Test alert when TTL exceeded
	oldCreatedAt := time.Now().Add(-2 * time.Hour)
	am.CheckTTLAlert("test", "resource-2", oldCreatedAt, 1*time.Hour)

	time.Sleep(100 * time.Millisecond)

	alerts := am.GetAlertsByEnvironment("test")
	if len(alerts) < 1 {
		t.Errorf("Expected at least 1 alert, got %d", len(alerts))
	}
}
