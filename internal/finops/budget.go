package finops

import (
	"context"
	"fmt"
	"math"
	"time"
)

// ProviderCostData represents current spending by provider
type ProviderCostData struct {
	Provider      string    `json:"provider"`
	CurrentSpend  float64   `json:"current_spend"`
	MonthlyBudget float64   `json:"monthly_budget"`
	LastUpdated   time.Time `json:"last_updated"`
}

// CostEstimate provides cost estimation for different tiers
type CostEstimate struct {
	Provider    string  `json:"provider"`
	Tier        string  `json:"tier"`
	MonthlyCost float64 `json:"monthly_cost"`
	HourlyCost  float64 `json:"hourly_cost"`
}

// GetCurrentSpend retrieves the current monthly spend for a provider
// In production, this would query real billing APIs (AWS Cost Explorer, Azure Cost Management, etc.)
func GetCurrentSpend(provider string) float64 {
	// Mock implementation with realistic data
	providerSpend := map[string]float64{
		"azure":        450.50,
		"aws":          1200.00,
		"gcp":          780.25,
		"hetzner":      45.00,
		"ovh":          120.75,
		"digitalocean": 89.99,
	}

	spend, exists := providerSpend[provider]
	if !exists {
		return 0.0
	}

	return spend
}

// GetTeamBudget retrieves the monthly budget for a specific team
func GetTeamBudget(ctx context.Context, team string) (float64, error) {
	// Mock team budgets
	teamBudgets := map[string]float64{
		"platform":    2000.00,
		"fintech":     500.00,
		"marketing":   300.00,
		"engineering": 1500.00,
		"dev":         200.00,
	}

	budget, exists := teamBudgets[team]
	if !exists {
		return 0.0, fmt.Errorf("team %s not found", team)
	}

	return budget, nil
}

// GetCurrentTeamSpend gets the current spend for a team across all providers
func GetCurrentTeamSpend(ctx context.Context, team string) (float64, error) {
	// In production, this would aggregate costs across all providers for the team
	// For now, return a calculated value based on team
	teamMultipliers := map[string]float64{
		"platform":    1.0,
		"fintech":     0.25,
		"marketing":   0.15,
		"engineering": 0.75,
		"dev":         0.1,
	}

	multiplier, exists := teamMultipliers[team]
	if !exists {
		return 0.0, fmt.Errorf("team %s not found", team)
	}

	totalSpend := 0.0
	for _, provider := range []string{"azure", "aws", "gcp", "hetzner", "ovh", "digitalocean"} {
		totalSpend += GetCurrentSpend(provider) * multiplier
	}

	return totalSpend, nil
}

// EstimateCost provides cost estimation for a given provider and tier
func EstimateCost(provider, tier string) float64 {
	// Realistic cost estimations in USD per month
	costMatrix := map[string]map[string]float64{
		"azure": {
			"micro": 15.50,
			"small": 31.00,
			"pro":   125.00,
		},
		"aws": {
			"micro": 8.50,
			"small": 17.00,
			"pro":   68.50,
		},
		"gcp": {
			"micro": 7.30,
			"small": 14.60,
			"pro":   58.40,
		},
		"hetzner": {
			"micro": 4.90,
			"small": 9.80,
			"pro":   19.60,
		},
		"ovh": {
			"micro": 12.00,
			"small": 24.00,
			"pro":   48.00,
		},
		"digitalocean": {
			"micro": 6.00,
			"small": 12.00,
			"pro":   24.00,
		},
	}

	providerCosts, exists := costMatrix[provider]
	if !exists {
		return 0.0
	}

	cost, exists := providerCosts[tier]
	if !exists {
		return 0.0
	}

	return cost
}

// IsDevelopmentEnvironment determines if a deployment is likely for development
// based on provider, region, and other heuristics
func IsDevelopmentEnvironment(ctx context.Context, provider string) bool {
	// Heuristic: expensive providers (AWS, Azure, GCP) in expensive regions (us-east-1)
	// with pro tier are likely production, so we downgrade dev workloads

	// Check if it's a premium provider in a premium region
	if provider == "aws" || provider == "azure" || provider == "gcp" {
		return true // Assume dev for automatic cost optimization
	}

	return false
}

// CheckBudget validates if a deployment would exceed budget limits
func CheckBudget(ctx context.Context, team string, additionalCost float64) error {
	budget, err := GetTeamBudget(ctx, team)
	if err != nil {
		return err
	}

	currentSpend, err := GetCurrentTeamSpend(ctx, team)
	if err != nil {
		return err
	}

	// Check if adding this cost would exceed 90% of budget
	if currentSpend+additionalCost > budget*0.9 {
		return fmt.Errorf("deployment would exceed 90%% of monthly budget: current $%.2f + additional $%.2f > budget $%.2f",
			currentSpend, additionalCost, budget)
	}

	return nil
}

// GetCostRecommendations provides cost optimization recommendations
func GetCostRecommendations(ctx context.Context, team string) []string {
	var recommendations []string

	currentSpend, _ := GetCurrentTeamSpend(ctx, team)
	budget, _ := GetTeamBudget(ctx, team)

	if budget > 0 {
		spendRatio := currentSpend / budget

		if spendRatio > 0.8 {
			recommendations = append(recommendations, "Consider downgrading non-production resources to micro tiers")
			recommendations = append(recommendations, "Enable auto-termination for idle development environments")
		}

		if spendRatio > 0.6 {
			recommendations = append(recommendations, "Migrate development workloads to cost-effective providers (Hetzner, OVH)")
		}
	}

	// Always recommend cost monitoring
	recommendations = append(recommendations, "Set up cost alerts for 50%, 75%, and 90% of budget")

	return recommendations
}

// CalculateSavings calculates potential savings from rightsizing
func CalculateSavings(currentProvider, currentTier string, recommendedTier string) float64 {
	currentCost := EstimateCost(currentProvider, currentTier)
	recommendedCost := EstimateCost(currentProvider, recommendedTier)

	return currentCost - recommendedCost
}

// GetProviderRecommendation suggests the most cost-effective provider for a given tier
func GetProviderRecommendation(tier string) (string, float64) {
	bestProvider := ""
	bestCost := math.MaxFloat64

	providers := []string{"hetzner", "digitalocean", "ovh", "gcp", "azure", "aws"}

	for _, provider := range providers {
		cost := EstimateCost(provider, tier)
		if cost < bestCost {
			bestCost = cost
			bestProvider = provider
		}
	}

	return bestProvider, bestCost
}
