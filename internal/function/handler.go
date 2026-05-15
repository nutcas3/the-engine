package function

import (
	"context"
	"fmt"

	"the-engine/internal/finops"
	providerpkg "the-engine/internal/provider"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
)

type Handler struct{}

func (h *Handler) RunFunction(ctx context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	// Extract values from the request
	var observed map[string]any
	if req.Observed != nil {
		if len(req.Observed.Resources) > 0 {
			observed = req.Observed.Resources["managed-compute"].Resource.AsMap()
		}
	}

	cloud := GetString(observed, "spec.provider")
	tier := GetString(observed, "spec.tier")
	region := GetString(observed, "spec.region")
	budgetLimit := GetFloat64(observed, "spec.budget_max")
	manualApproval := GetBool(observed, "spec.manual_approval")

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
	desiredResource := ProcessDeployment(ctx, cloud, tier, region, budgetLimit)
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

// ProcessDeployment maps the deployment to the appropriate provider
func ProcessDeployment(ctx context.Context, provider string, tier string, region string, budgetLimit float64) map[string]any {
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
