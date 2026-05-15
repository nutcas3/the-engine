package main

import (
	"context"
	"fmt"

	"the-engine/internal/finops"
	providerpkg "the-engine/internal/provider"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
)

type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer
}

func (f *Function) RunFunction(ctx context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	// Extract values from the request
	var observed map[string]any
	if req.Observed != nil {
		if len(req.Observed.Resources) > 0 {
			observed = req.Observed.Resources["managed-compute"].Resource.AsMap()
		}
	}

	// Extract intent from the XRD
	cloud := getString(observed, "spec.provider")
	tier := getString(observed, "spec.tier")
	region := getString(observed, "spec.region")
	budgetLimit := getFloat64(observed, "spec.budget_max")
	manualApproval := getBool(observed, "spec.manual_approval")

	// Set defaults if not provided
	if cloud == "" {
		cloud = "hetzner"
	}
	if tier == "" {
		tier = "micro"
	}
	if region == "" {
		region = "nbg1"
	}

	// 2. FinOps Guardrail: Real-time Budget Check
	if !manualApproval && budgetLimit > 0 {
		currentSpend := finops.GetCurrentSpend(cloud)
		if currentSpend > budgetLimit {
			return &fnv1beta1.RunFunctionResponse{
				Results: []*fnv1beta1.Result{
					{
						Message: fmt.Sprintf("budget exceeded for %s: current spend $%.2f > limit $%.2f", cloud, currentSpend, budgetLimit),
					},
				},
			}, nil
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
		return &fnv1beta1.RunFunctionResponse{
			Results: []*fnv1beta1.Result{
				{
					Message: fmt.Sprintf("unsupported cloud provider: %s", cloud),
				},
			},
		}, nil
	}

	// 5. Set desired resource in response
	desiredStruct, err := structpb.NewStruct(desiredResource)
	if err != nil {
		return nil, fmt.Errorf("failed to convert desired resource to struct: %w", err)
	}

	rsp := &fnv1beta1.RunFunctionResponse{
		Desired: &fnv1beta1.State{
			Resources: map[string]*fnv1beta1.Resource{
				"managed-compute": {
					Resource: desiredStruct,
				},
			},
		},
	}

	// 6. Add cost estimation message
	estimatedCost := finops.EstimateCost(cloud, tier)
	rsp.Results = []*fnv1beta1.Result{
		{
			Message: fmt.Sprintf("Estimated monthly cost: $%.2f", estimatedCost),
		},
	}

	return rsp, nil
}

// Helper functions to extract values from unstructured objects
func getString(obj map[string]any, path ...string) string {
	current := obj
	for i, key := range path {
		if i == len(path)-1 {
			if val, ok := current[key].(string); ok {
				return val
			}
			return ""
		}
		if next, ok := current[key].(map[string]any); ok {
			current = next
		} else {
			return ""
		}
	}
	return ""
}

func getFloat64(obj map[string]any, path ...string) float64 {
	current := obj
	for i, key := range path {
		if i == len(path)-1 {
			switch val := current[key].(type) {
			case float64:
				return val
			case float32:
				return float64(val)
			case int:
				return float64(val)
			case int64:
				return float64(val)
			default:
				return 0
			}
		}
		if next, ok := current[key].(map[string]any); ok {
			current = next
		} else {
			return 0
		}
	}
	return 0
}

func getBool(obj map[string]any, path ...string) bool {
	current := obj
	for i, key := range path {
		if i == len(path)-1 {
			if val, ok := current[key].(bool); ok {
				return val
			}
			return false
		}
		if next, ok := current[key].(map[string]any); ok {
			current = next
		} else {
			return false
		}
	}
	return false
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
