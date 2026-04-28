package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"the-engine/internal/finops"
)

type DeployCommand struct {
	Provider string  `json:"provider"`
	Tier     string  `json:"tier"`
	Region   string  `json:"region"`
	Budget   float64 `json:"budget"`
	Team     string  `json:"team"`
	Manual   bool    `json:"manual"`
}

type ListCommand struct {
	Provider string `json:"provider"`
	Tier     string `json:"tier"`
	Team     string `json:"team"`
}

type CostCommand struct {
	Team  string `json:"team"`
	Month string `json:"month"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "deploy":
		handleDeploy()
	case "list":
		handleList()
	case "cost":
		handleCost()
	case "drift-check":
		handleDriftCheck()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleDeploy() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: engine deploy --provider <provider> --tier <tier> --region <region> [--budget <amount>] [--team <team>] [--manual]")
		os.Exit(1)
	}

	cmd := parseDeployCommand()

	// Validate input
	if cmd.Provider == "" || cmd.Tier == "" {
		fmt.Println("Error: --provider and --tier are required")
		os.Exit(1)
	}

	// Check budget if specified
	if cmd.Budget > 0 {
		err := finops.CheckBudget(context.Background(), cmd.Team, finops.EstimateCost(cmd.Provider, cmd.Tier))
		if err != nil {
			fmt.Printf("[REJECTED]: %s\n", err.Error())
			fmt.Printf("Try --tier micro or request a budget increase\n")
			os.Exit(1)
		}
	}

	// Generate the XRD manifest
	manifest := generateXComputeManifest(cmd)

	fmt.Printf("Deploying %s tier on %s in %s\n", cmd.Tier, cmd.Provider, cmd.Region)
	fmt.Printf("Estimated monthly cost: $%.2f\n", finops.EstimateCost(cmd.Provider, cmd.Tier))
	fmt.Printf("\nGenerated manifest:\n")
	fmt.Println(manifest)
}

func handleList() {
	cmd := parseListCommand()

	fmt.Printf("Listing deployments")
	if cmd.Provider != "" {
		fmt.Printf(" on %s", cmd.Provider)
	}
	if cmd.Tier != "" {
		fmt.Printf(" with tier %s", cmd.Tier)
	}
	if cmd.Team != "" {
		fmt.Printf(" for team %s", cmd.Team)
	}

	// Mock listing - in production would query actual deployments
	deployments := []map[string]any{
		{
			"id":       "xcompute-123",
			"provider": "hetzner",
			"tier":     "micro",
			"region":   "nbg1",
			"status":   "running",
			"cost":     4.90,
			"created":  time.Now().Add(-2 * time.Hour),
		},
		{
			"id":       "xcompute-456",
			"provider": "azure",
			"tier":     "small",
			"region":   "westeurope",
			"status":   "running",
			"cost":     31.00,
			"created":  time.Now().Add(-24 * time.Hour),
		},
	}

	for _, deployment := range deployments {
		fmt.Printf("ID: %s | Provider: %s | Tier: %s | Region: %s | Status: %s | Cost: $%.2f\n",
			deployment["id"], deployment["provider"], deployment["tier"],
			deployment["region"], deployment["status"], deployment["cost"])
	}
}

func handleCost() {
	cmd := parseCostCommand()

	fmt.Printf("Cost report for team %s", cmd.Team)
	if cmd.Month != "" {
		fmt.Printf(" for %s", cmd.Month)
	}
	fmt.Println(":")

	// Get team budget and current spend
	budget, err := finops.GetTeamBudget(context.Background(), cmd.Team)
	if err != nil {
		fmt.Printf("Error getting budget: %v\n", err)
		os.Exit(1)
	}

	spend, err := finops.GetCurrentTeamSpend(context.Background(), cmd.Team)
	if err != nil {
		fmt.Printf("Error getting spend: %v\n", err)
		os.Exit(1)
	}

	percentage := (spend / budget) * 100

	fmt.Printf("Budget: $%.2f\n", budget)
	fmt.Printf("Current Spend: $%.2f (%.1f%%)\n", spend, percentage)

	if percentage > 90 {
		fmt.Printf("Status: CRITICAL - Budget nearly exceeded!\n")
	} else if percentage > 75 {
		fmt.Printf("Status: WARNING - Approaching budget limit\n")
	} else {
		fmt.Printf("Status: OK\n")
	}

	// Show cost recommendations
	recommendations := finops.GetCostRecommendations(context.Background(), cmd.Team)
	if len(recommendations) > 0 {
		fmt.Printf("\nRecommendations:\n")
		for i, rec := range recommendations {
			fmt.Printf("%d. %s\n", i+1, rec)
		}
	}
}

func handleDriftCheck() {
	fmt.Println("Checking for configuration drift...")

	// Mock drift detection - in production would check actual vs desired state
	driftEvents := []map[string]any{
		{
			"resource": "xcompute-123",
			"type":     "configuration_drift",
			"severity": "warning",
			"details":  "Instance size changed from micro to small",
		},
		{
			"resource": "xcompute-456",
			"type":     "security_drift",
			"severity": "critical",
			"details":  "SSH key added to instance",
		},
	}

	if len(driftEvents) == 0 {
		fmt.Println("No drift detected. All resources are in sync.")
		return
	}

	fmt.Printf("Found %d drift events:\n", len(driftEvents))
	for _, event := range driftEvents {
		fmt.Printf("Resource: %s | Type: %s | Severity: %s | Details: %s\n",
			event["resource"], event["type"], event["severity"], event["details"])
	}

	fmt.Printf("\nRun 'engine reconcile --all' to fix drift\n")
}

func parseDeployCommand() DeployCommand {
	cmd := DeployCommand{
		Region: "us-east-1",
		Team:   "dev",
	}

	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "--provider" && i+1 < len(os.Args):
			cmd.Provider = os.Args[i+1]
			i++
		case arg == "--tier" && i+1 < len(os.Args):
			cmd.Tier = os.Args[i+1]
			i++
		case arg == "--region" && i+1 < len(os.Args):
			cmd.Region = os.Args[i+1]
			i++
		case arg == "--budget" && i+1 < len(os.Args):
			var budget float64
			fmt.Sscanf(os.Args[i+1], "%f", &budget)
			cmd.Budget = budget
			i++
		case arg == "--team" && i+1 < len(os.Args):
			cmd.Team = os.Args[i+1]
			i++
		case arg == "--manual":
			cmd.Manual = true
		}
	}

	return cmd
}

func parseListCommand() ListCommand {
	cmd := ListCommand{}

	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "--provider" && i+1 < len(os.Args):
			cmd.Provider = os.Args[i+1]
			i++
		case arg == "--tier" && i+1 < len(os.Args):
			cmd.Tier = os.Args[i+1]
			i++
		case arg == "--team" && i+1 < len(os.Args):
			cmd.Team = os.Args[i+1]
			i++
		}
	}

	return cmd
}

func parseCostCommand() CostCommand {
	cmd := CostCommand{
		Month: "current",
		Team:  "platform",
	}

	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "--team" && i+1 < len(os.Args):
			cmd.Team = os.Args[i+1]
			i++
		case arg == "--month" && i+1 < len(os.Args):
			cmd.Month = os.Args[i+1]
			i++
		}
	}

	return cmd
}

func generateXComputeManifest(cmd DeployCommand) string {
	manifest := map[string]any{
		"apiVersion": "engine.io/v1alpha1",
		"kind":       "XCompute",
		"metadata": map[string]any{
			"name": fmt.Sprintf("%s-%s-%s", cmd.Team, cmd.Provider, cmd.Tier),
		},
		"spec": map[string]any{
			"provider":        cmd.Provider,
			"tier":            cmd.Tier,
			"region":          cmd.Region,
			"budget_max":      cmd.Budget,
			"manual_approval": cmd.Manual,
		},
	}

	jsonBytes, _ := json.MarshalIndent(manifest, "", "  ")
	return string(jsonBytes)
}

func printUsage() {
	fmt.Println("Sovereign Engine CLI - Multi-Cloud Infrastructure Platform")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  engine deploy --provider <provider> --tier <tier> [--region <region>] [--budget <amount>] [--team <team>] [--manual]")
	fmt.Println("  engine list [--provider <provider>] [--tier <tier>] [--team <team>]")
	fmt.Println("  engine cost --team <team> [--month <month>]")
	fmt.Println("  engine drift-check")
	fmt.Println("  engine help")
	fmt.Println()
	fmt.Println("Providers: aws, azure, gcp, hetzner, ovh, digitalocean")
	fmt.Println("Tiers: micro, small, pro")
	fmt.Println("Regions: us-east-1, us-west-2, europe-west1, etc.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  engine deploy --provider hetzner --tier micro --region nbg1")
	fmt.Println("  engine deploy --provider azure --tier pro --budget 500 --team platform")
	fmt.Println("  engine list --provider hetzner")
	fmt.Println("  engine cost --team platform")
	fmt.Println("  engine drift-check")
}
