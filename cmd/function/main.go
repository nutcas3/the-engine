package main

import (
	"context"
	"fmt"

	"the-engine/internal/finops"
	"the-engine/internal/provider"

	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
)

type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer
}

func (f *Function) RunFunction(ctx context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	rsp := response.New()
	oxr, _ := request.GetObservedCompositeResource(req)

	// 1. Extract intent from the XRD
	cloud, _ := oxr.Resource.GetString("spec.provider")
	tier, _ := oxr.Resource.GetString("spec.tier")
	region, _ := oxr.Resource.GetString("spec.region")
	budgetLimit, _ := oxr.Resource.GetFloat64("spec.budget_max")
	manualApproval, _ := oxr.Resource.GetBool("spec.manual_approval")

	// 2. FinOps Guardrail: Real-time Budget Check (unless manual approval)
	if !manualApproval && budgetLimit > 0 {
		currentSpend := finops.GetCurrentSpend(cloud)
		if currentSpend > budgetLimit {
			return response.Fatal(rsp, fmt.Errorf("budget exceeded for %s: current spend $%.2f > limit $%.2f", cloud, currentSpend, budgetLimit)), nil
		}
	}

	// 3. Automated Thrift: Downgrade dev environments to cheaper clouds
	if tier == "pro" && region == "us-east-1" && !manualApproval {
		// Check if this is likely a dev environment and downgrade to micro on Hetzner
		if finops.IsDevelopmentEnvironment(ctx, cloud) {
			cloud = "hetzner"
			tier = "micro"
			region = "nbg1"
			response.Warning(rsp, "Automatically downgraded to micro tier on Hetzner for cost optimization")
		}
	}

	// 4. Multi-Cloud Mapping: Translate to Provider-specific resources
	var desiredResource any
	switch cloud {
	case "azure":
		desiredResource = provider.MapAzure(tier, region)
	case "aws":
		desiredResource = provider.MapAWS(tier, region)
	case "gcp":
		desiredResource = provider.MapGCP(tier, region)
	case "hetzner":
		desiredResource = provider.MapHetzner(tier, region)
	case "ovh":
		desiredResource = provider.MapOVH(tier, region)
	case "digitalocean":
		desiredResource = provider.MapDigitalOcean(tier, region)
	default:
		return response.Fatal(rsp, fmt.Errorf("unsupported cloud provider: %s", cloud)), nil
	}

	// 5. Set the desired composed resource
	response.SetDesiredComposedResource(rsp, "managed-compute", desiredResource)
	
	// 6. Add cost estimation to response
	estimatedCost := finops.EstimateCost(cloud, tier)
	response.SetCondition(rsp, "Ready", "Estimated monthly cost: $%.2f", estimatedCost)
	
	return rsp, nil
}
