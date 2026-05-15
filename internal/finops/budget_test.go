package finops

import (
	"context"
	"testing"
)

func TestGetCurrentSpend(t *testing.T) {
	spend := GetCurrentSpend("aws")
	if spend <= 0 {
		t.Error("Expected positive spend for aws")
	}

	spend = GetCurrentSpend("unknown")
	if spend != 0 {
		t.Error("Expected 0 for unknown provider")
	}
}

func TestGetTeamBudget(t *testing.T) {
	budget, err := GetTeamBudget(context.Background(), "platform")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if budget <= 0 {
		t.Error("Expected positive budget")
	}

	_, err = GetTeamBudget(context.Background(), "unknown")
	if err == nil {
		t.Error("Expected error for unknown team")
	}
}

func TestGetCurrentTeamSpend(t *testing.T) {
	spend, err := GetCurrentTeamSpend(context.Background(), "platform")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if spend <= 0 {
		t.Error("Expected positive spend")
	}

	_, err = GetCurrentTeamSpend(context.Background(), "unknown")
	if err == nil {
		t.Error("Expected error for unknown team")
	}
}

func TestEstimateCost(t *testing.T) {
	cost := EstimateCost("aws", "small")
	if cost <= 0 {
		t.Error("Expected positive cost")
	}

	cost = EstimateCost("unknown", "small")
	if cost != 0 {
		t.Error("Expected 0 for unknown provider")
	}
}

func TestIsDevelopmentEnvironment(t *testing.T) {
	result := IsDevelopmentEnvironment(context.Background(), "aws")
	if !result {
		t.Error("Expected true for aws")
	}

	result = IsDevelopmentEnvironment(context.Background(), "hetzner")
	if result {
		t.Error("Expected false for hetzner")
	}
}

func TestCheckBudget(t *testing.T) {
	// Test error case - deployment exceeds budget
	err := CheckBudget(context.Background(), "dev", 100.0)
	if err == nil {
		t.Error("Expected error when deployment exceeds budget")
	}
}

func TestGetCostRecommendations(t *testing.T) {
	recommendations := GetCostRecommendations(context.Background(), "platform")
	if len(recommendations) == 0 {
		t.Error("Expected some recommendations")
	}
}

func TestCalculateSavings(t *testing.T) {
	savings := CalculateSavings("aws", "pro", "small")
	if savings <= 0 {
		t.Error("Expected positive savings")
	}
}

func TestGetProviderRecommendation(t *testing.T) {
	provider, cost := GetProviderRecommendation("small")
	if provider == "" {
		t.Error("Expected non-empty provider")
	}
	if cost <= 0 {
		t.Error("Expected positive cost")
	}
}
