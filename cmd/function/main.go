package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"the-engine/internal/finops"
	providerpkg "the-engine/internal/provider"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"google.golang.org/grpc"
)

type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer
}

func (f *Function) RunFunction(ctx context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	// Extract values from the request using proto
	var cloud, tier, region string
	var budgetLimit float64
	var manualApproval bool

	// Try to extract from input struct
	if req.Input != nil {
		obj := req.Input.AsMap()
		if spec, ok := obj["spec"].(map[string]any); ok {
			if val, ok := spec["provider"].(string); ok {
				cloud = val
			}
			if val, ok := spec["tier"].(string); ok {
				tier = val
			}
			if val, ok := spec["region"].(string); ok {
				region = val
			}
			if val, ok := spec["budget_max"].(float64); ok {
				budgetLimit = val
			}
			if val, ok := spec["manual_approval"].(bool); ok {
				manualApproval = val
			}
		}
	}

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
				Meta: &fnv1beta1.ResponseMeta{
					Tag: "budget-exceeded",
				},
			}, fmt.Errorf("budget exceeded for %s: current spend $%.2f > limit $%.2f", cloud, currentSpend, budgetLimit)
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
			Meta: &fnv1beta1.ResponseMeta{
				Tag: "unsupported-provider",
			},
		}, fmt.Errorf("unsupported cloud provider: %s", cloud)
	}

	// 5. Create response with cost estimation
	estimatedCost := finops.EstimateCost(cloud, tier)
	_ = estimatedCost // Will be added to response when proper SDK integration is complete

	return &fnv1beta1.RunFunctionResponse{
		Meta: &fnv1beta1.ResponseMeta{
			Tag: "success",
		},
	}, nil
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

func main() {
	lis, err := net.Listen("tcp", ":9443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	fnv1beta1.RegisterFunctionRunnerServiceServer(s, &Function{})

	log.Printf("Sovereign Engine function server listening on :9443")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
