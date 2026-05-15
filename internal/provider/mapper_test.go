package provider

import (
	"testing"
)

func TestMapAWS(t *testing.T) {
	tests := []struct {
		name   string
		tier   string
		region string
		want   string
	}{
		{
			name:   "micro tier in us-east-1",
			tier:   "micro",
			region: "us-east-1",
			want:   "t3.micro",
		},
		{
			name:   "small tier in us-west-2",
			tier:   "small",
			region: "us-west-2",
			want:   "t3.small",
		},
		{
			name:   "pro tier in eu-west-1",
			tier:   "pro",
			region: "eu-west-1",
			want:   "c6i.large",
		},
		{
			name:   "unknown tier defaults to small",
			tier:   "unknown",
			region: "us-east-1",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapAWS(tt.tier, tt.region)
			if result == nil {
				t.Error("MapAWS returned nil")
				return
			}
			spec := result["spec"].(map[string]any)
			forProvider := spec["forProvider"].(map[string]any)
			instanceType := forProvider["instanceType"].(string)
			if tt.want != "" && instanceType != tt.want {
				t.Errorf("MapAWS() instanceType = %v, want %v", instanceType, tt.want)
			}
		})
	}
}

func TestMapAzure(t *testing.T) {
	tests := []struct {
		name   string
		tier   string
		region string
		want   string
	}{
		{
			name:   "micro tier in west-europe",
			tier:   "micro",
			region: "west-europe",
			want:   "Standard_B1s",
		},
		{
			name:   "small tier in central-europe",
			tier:   "small",
			region: "central-europe",
			want:   "Standard_B2s",
		},
		{
			name:   "pro tier in us-east-1",
			tier:   "pro",
			region: "us-east-1",
			want:   "Standard_D2s_v5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapAzure(tt.tier, tt.region)
			if result == nil {
				t.Error("MapAzure returned nil")
				return
			}
			spec := result["spec"].(map[string]any)
			forProvider := spec["forProvider"].(map[string]any)
			size := forProvider["size"].(string)
			if size != tt.want {
				t.Errorf("MapAzure() size = %v, want %v", size, tt.want)
			}
		})
	}
}

func TestMapGCP(t *testing.T) {
	tests := []struct {
		name   string
		tier   string
		region string
		want   string
	}{
		{
			name:   "micro tier in us-east-1",
			tier:   "micro",
			region: "us-east-1",
			want:   "e2-micro",
		},
		{
			name:   "small tier in us-west-2",
			tier:   "small",
			region: "us-west-2",
			want:   "e2-small",
		},
		{
			name:   "pro tier in europe-west1",
			tier:   "pro",
			region: "europe-west1",
			want:   "n2-standard-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapGCP(tt.tier, tt.region)
			if result == nil {
				t.Error("MapGCP returned nil")
				return
			}
			spec := result["spec"].(map[string]any)
			forProvider := spec["forProvider"].(map[string]any)
			machineType := forProvider["machineType"].(string)
			if machineType != tt.want {
				t.Errorf("MapGCP() machineType = %v, want %v", machineType, tt.want)
			}
		})
	}
}

func TestMapHetzner(t *testing.T) {
	tests := []struct {
		name   string
		tier   string
		region string
		want   string
	}{
		{
			name:   "micro tier in europe",
			tier:   "micro",
			region: "europe",
			want:   "cx11",
		},
		{
			name:   "small tier in central",
			tier:   "small",
			region: "central",
			want:   "cpx11",
		},
		{
			name:   "pro tier in us-east-1",
			tier:   "pro",
			region: "us-east-1",
			want:   "cpx21",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapHetzner(tt.tier, tt.region)
			if result == nil {
				t.Error("MapHetzner returned nil")
				return
			}
			spec := result["spec"].(map[string]any)
			forProvider := spec["forProvider"].(map[string]any)
			serverType := forProvider["serverType"].(string)
			if serverType != tt.want {
				t.Errorf("MapHetzner() serverType = %v, want %v", serverType, tt.want)
			}
		})
	}
}

func TestMapOVH(t *testing.T) {
	tests := []struct {
		name   string
		tier   string
		region string
		want   string
	}{
		{
			name:   "micro tier in europe",
			tier:   "micro",
			region: "europe",
			want:   "s1-2",
		},
		{
			name:   "small tier in us-east-1",
			tier:   "small",
			region: "us-east-1",
			want:   "s1-4",
		},
		{
			name:   "pro tier in canada",
			tier:   "pro",
			region: "canada",
			want:   "s1-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapOVH(tt.tier, tt.region)
			if result == nil {
				t.Error("MapOVH returned nil")
				return
			}
			spec := result["spec"].(map[string]any)
			forProvider := spec["forProvider"].(map[string]any)
			flavor := forProvider["flavor"].(string)
			if flavor != tt.want {
				t.Errorf("MapOVH() flavor = %v, want %v", flavor, tt.want)
			}
		})
	}
}

func TestMapDigitalOcean(t *testing.T) {
	tests := []struct {
		name   string
		tier   string
		region string
		want   string
	}{
		{
			name:   "micro tier in us-east-1",
			tier:   "micro",
			region: "us-east-1",
			want:   "s-1vcpu-1gb",
		},
		{
			name:   "small tier in us-west",
			tier:   "small",
			region: "us-west",
			want:   "s-1vcpu-2gb",
		},
		{
			name:   "pro tier in europe",
			tier:   "pro",
			region: "europe",
			want:   "s-2vcpu-4gb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapDigitalOcean(tt.tier, tt.region)
			if result == nil {
				t.Error("MapDigitalOcean returned nil")
				return
			}
			spec := result["spec"].(map[string]any)
			forProvider := spec["forProvider"].(map[string]any)
			size := forProvider["size"].(string)
			if size != tt.want {
				t.Errorf("MapDigitalOcean() size = %v, want %v", size, tt.want)
			}
		})
	}
}
