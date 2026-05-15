package health

import (
	"testing"
	"time"
)

func TestNewChecker(t *testing.T) {
	checker := NewChecker("1.0.0")
	if checker == nil {
		t.Error("Expected checker to be created")
	}
	if checker.version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", checker.version)
	}
	if checker.startTime.IsZero() {
		t.Error("Expected start time to be set")
	}
}

func TestChecker_Check(t *testing.T) {
	checker := NewChecker("1.0.0")
	healthResponse := checker.Check()

	if healthResponse.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", healthResponse.Version)
	}
	if healthResponse.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
	if len(healthResponse.Components) == 0 {
		t.Error("Expected components to be present")
	}
	if healthResponse.System.GoVersion == "" {
		t.Error("Expected Go version to be set")
	}
	if healthResponse.System.Goroutines == 0 {
		t.Error("Expected goroutines count to be set")
	}
}

func TestChecker_calculateOverallStatus(t *testing.T) {
	checker := NewChecker("1.0.0")

	// Test all healthy
	healthyComponents := []ComponentHealth{
		{Name: "test1", Status: StatusHealthy, CheckedAt: time.Now()},
		{Name: "test2", Status: StatusHealthy, CheckedAt: time.Now()},
	}
	status := checker.calculateOverallStatus(healthyComponents)
	if status != "healthy" {
		t.Errorf("Expected 'healthy', got '%s'", status)
	}

	// Test degraded
	degradedComponents := []ComponentHealth{
		{Name: "test1", Status: StatusHealthy, CheckedAt: time.Now()},
		{Name: "test2", Status: StatusDegraded, CheckedAt: time.Now()},
	}
	status = checker.calculateOverallStatus(degradedComponents)
	if status != "degraded" {
		t.Errorf("Expected 'degraded', got '%s'", status)
	}

	// Test unhealthy
	unhealthyComponents := []ComponentHealth{
		{Name: "test1", Status: StatusHealthy, CheckedAt: time.Now()},
		{Name: "test2", Status: StatusUnhealthy, CheckedAt: time.Now()},
	}
	status = checker.calculateOverallStatus(unhealthyComponents)
	if status != "unhealthy" {
		t.Errorf("Expected 'unhealthy', got '%s'", status)
	}
}
