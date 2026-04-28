package main

import (
	"context"
	"fmt"

	"the-engine/internal/finops"
	providerpkg "the-engine/internal/provider"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
)

type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer
}

func (f *Function) RunFunction(ctx context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	// Default values - SDK integration pending
	cloud := "hetzner"
	tier := "micro"
	region := "nbg1"
	budgetLimit := 0.0
	manualApproval := false

	// 2. FinOps Guardrail: Real-time Budget Check
	if !manualApproval && budgetLimit > 0 {
		currentSpend := finops.GetCurrentSpend(cloud)
		if currentSpend > budgetLimit {
			return &fnv1beta1.RunFunctionResponse{}, fmt.Errorf("budget exceeded for %s: current spend $%.2f > limit $%.2f", cloud, currentSpend, budgetLimit)
		}
	}

	// 3. Automated Thrift: Downgrade dev environments
	if tier == "pro" && region == "us-east-1" && !manualApproval {
		if finops.IsDevelopmentEnvironment(ctx, cloud) {
			cloud = "hetzner"
			tier = "micro"
			region = "nbg1"
		}
	}

	// 4. Multi-Cloud Mapping
	desiredResource := processDeployment(ctx, cloud, tier, region, budgetLimit)
	if desiredResource == nil {
		return &fnv1beta1.RunFunctionResponse{}, fmt.Errorf("unsupported cloud provider: %s", cloud)
	}

	// 5. Cost estimation
	estimatedCost := finops.EstimateCost(cloud, tier)
	_ = estimatedCost
	_ = desiredResource

	return &fnv1beta1.RunFunctionResponse{}, nil
}

func processDeployment(ctx context.Context, provider string, tier string, region string, budgetLimit float64) map[string]any {
	switch provider {
	case "azure":
		return providerpkg.MapAzure(tier, region)
	case "aws":
		return providerpkg.MapAWS(tier, region)
	case "gcp":
		return providerpkg.MapGCP(tier, region)
	case "hetzner":
		return providerpkg.MapHetzner(tier, region)
	case "ovh":
		return providerpkg.MapOVH(tier, region)
	case "digitalocean":
		return providerpkg.MapDigitalOcean(tier, region)
	default:
		return nil
	}
}
