package compositions

import (
	"time"

	"the-engine/internal/types"
)

// GetCompositions returns the list of available compositions
// It loads from YAML files in the compositions directory, falling back to hardcoded data if loading fails
func GetCompositions() []types.Composition {
	// Try to load from YAML files first
	compositions, err := LoadCompositionsFromYAML()
	if err == nil && len(compositions) > 0 {
		return compositions
	}

	// Fallback to hardcoded data if YAML loading fails
	return getHardcodedCompositions()
}

// getHardcodedCompositions returns the hardcoded list of compositions
func getHardcodedCompositions() []types.Composition {
	now := time.Now().Format(time.RFC3339)

	return []types.Composition{
		// AWS compositions
		{Name: "aws-compute", Provider: "aws", Type: "compute", Labels: map[string]string{"provider": "aws", "engine.io/composition": "compute"}, CreatedAt: now},
		{Name: "aws-networking", Provider: "aws", Type: "networking", Labels: map[string]string{"provider": "aws", "engine.io/composition": "networking"}, CreatedAt: now},
		{Name: "aws-loadbalancer", Provider: "aws", Type: "loadbalancer", Labels: map[string]string{"provider": "aws", "engine.io/composition": "loadbalancer"}, CreatedAt: now},
		{Name: "aws-storage", Provider: "aws", Type: "storage", Labels: map[string]string{"provider": "aws", "engine.io/composition": "storage"}, CreatedAt: now},
		{Name: "aws-database", Provider: "aws", Type: "database", Labels: map[string]string{"provider": "aws", "engine.io/composition": "database"}, CreatedAt: now},
		{Name: "aws-monitoring", Provider: "aws", Type: "monitoring", Labels: map[string]string{"provider": "aws", "engine.io/composition": "monitoring"}, CreatedAt: now},
		{Name: "aws-iam", Provider: "aws", Type: "iam", Labels: map[string]string{"provider": "aws", "engine.io/composition": "iam"}, CreatedAt: now},

		// Azure compositions
		{Name: "azure-compute", Provider: "azure", Type: "compute", Labels: map[string]string{"provider": "azure", "engine.io/composition": "compute"}, CreatedAt: now},
		{Name: "azure-networking", Provider: "azure", Type: "networking", Labels: map[string]string{"provider": "azure", "engine.io/composition": "networking"}, CreatedAt: now},
		{Name: "azure-loadbalancer", Provider: "azure", Type: "loadbalancer", Labels: map[string]string{"provider": "azure", "engine.io/composition": "loadbalancer"}, CreatedAt: now},
		{Name: "azure-storage", Provider: "azure", Type: "storage", Labels: map[string]string{"provider": "azure", "engine.io/composition": "storage"}, CreatedAt: now},
		{Name: "azure-database", Provider: "azure", Type: "database", Labels: map[string]string{"provider": "azure", "engine.io/composition": "database"}, CreatedAt: now},
		{Name: "azure-monitoring", Provider: "azure", Type: "monitoring", Labels: map[string]string{"provider": "azure", "engine.io/composition": "monitoring"}, CreatedAt: now},
		{Name: "azure-iam", Provider: "azure", Type: "iam", Labels: map[string]string{"provider": "azure", "engine.io/composition": "iam"}, CreatedAt: now},

		// GCP compositions
		{Name: "gcp-compute", Provider: "gcp", Type: "compute", Labels: map[string]string{"provider": "gcp", "engine.io/composition": "compute"}, CreatedAt: now},
		{Name: "gcp-networking", Provider: "gcp", Type: "networking", Labels: map[string]string{"provider": "gcp", "engine.io/composition": "networking"}, CreatedAt: now},
		{Name: "gcp-loadbalancer", Provider: "gcp", Type: "loadbalancer", Labels: map[string]string{"provider": "gcp", "engine.io/composition": "loadbalancer"}, CreatedAt: now},
		{Name: "gcp-storage", Provider: "gcp", Type: "storage", Labels: map[string]string{"provider": "gcp", "engine.io/composition": "storage"}, CreatedAt: now},
		{Name: "gcp-database", Provider: "gcp", Type: "database", Labels: map[string]string{"provider": "gcp", "engine.io/composition": "database"}, CreatedAt: now},
		{Name: "gcp-monitoring", Provider: "gcp", Type: "monitoring", Labels: map[string]string{"provider": "gcp", "engine.io/composition": "monitoring"}, CreatedAt: now},
		{Name: "gcp-iam", Provider: "gcp", Type: "iam", Labels: map[string]string{"provider": "gcp", "engine.io/composition": "iam"}, CreatedAt: now},

		// Hetzner compositions
		{Name: "hetzner-compute", Provider: "hetzner", Type: "compute", Labels: map[string]string{"provider": "hetzner", "engine.io/composition": "compute"}, CreatedAt: now},
		{Name: "hetzner-networking", Provider: "hetzner", Type: "networking", Labels: map[string]string{"provider": "hetzner", "engine.io/composition": "networking"}, CreatedAt: now},
		{Name: "hetzner-storage", Provider: "hetzner", Type: "storage", Labels: map[string]string{"provider": "hetzner", "engine.io/composition": "storage"}, CreatedAt: now},
		{Name: "hetzner-database", Provider: "hetzner", Type: "database", Labels: map[string]string{"provider": "hetzner", "engine.io/composition": "database"}, CreatedAt: now},
		{Name: "hetzner-monitoring", Provider: "hetzner", Type: "monitoring", Labels: map[string]string{"provider": "hetzner", "engine.io/composition": "monitoring"}, CreatedAt: now},

		// OVH compositions
		{Name: "ovh-compute", Provider: "ovh", Type: "compute", Labels: map[string]string{"provider": "ovh", "engine.io/composition": "compute"}, CreatedAt: now},
		{Name: "ovh-networking", Provider: "ovh", Type: "networking", Labels: map[string]string{"provider": "ovh", "engine.io/composition": "networking"}, CreatedAt: now},
		{Name: "ovh-loadbalancer", Provider: "ovh", Type: "loadbalancer", Labels: map[string]string{"provider": "ovh", "engine.io/composition": "loadbalancer"}, CreatedAt: now},
		{Name: "ovh-storage", Provider: "ovh", Type: "storage", Labels: map[string]string{"provider": "ovh", "engine.io/composition": "storage"}, CreatedAt: now},
		{Name: "ovh-database", Provider: "ovh", Type: "database", Labels: map[string]string{"provider": "ovh", "engine.io/composition": "database"}, CreatedAt: now},
		{Name: "ovh-monitoring", Provider: "ovh", Type: "monitoring", Labels: map[string]string{"provider": "ovh", "engine.io/composition": "monitoring"}, CreatedAt: now},

		// DigitalOcean compositions
		{Name: "digitalocean-compute", Provider: "digitalocean", Type: "compute", Labels: map[string]string{"provider": "digitalocean", "engine.io/composition": "compute"}, CreatedAt: now},
		{Name: "digitalocean-networking", Provider: "digitalocean", Type: "networking", Labels: map[string]string{"provider": "digitalocean", "engine.io/composition": "networking"}, CreatedAt: now},
		{Name: "digitalocean-loadbalancer", Provider: "digitalocean", Type: "loadbalancer", Labels: map[string]string{"provider": "digitalocean", "engine.io/composition": "loadbalancer"}, CreatedAt: now},
		{Name: "digitalocean-storage", Provider: "digitalocean", Type: "storage", Labels: map[string]string{"provider": "digitalocean", "engine.io/composition": "storage"}, CreatedAt: now},
		{Name: "digitalocean-database", Provider: "digitalocean", Type: "database", Labels: map[string]string{"provider": "digitalocean", "engine.io/composition": "database"}, CreatedAt: now},
		{Name: "digitalocean-monitoring", Provider: "digitalocean", Type: "monitoring", Labels: map[string]string{"provider": "digitalocean", "engine.io/composition": "monitoring"}, CreatedAt: now},

		// Shared compositions
		{Name: "shared-monitoring-stack", Provider: "shared", Type: "monitoring", Labels: map[string]string{"provider": "shared", "engine.io/composition": "monitoring"}, CreatedAt: now},
		{Name: "shared-database-compute", Provider: "shared", Type: "database", Labels: map[string]string{"provider": "shared", "engine.io/composition": "database"}, CreatedAt: now},
		{Name: "shared-vault", Provider: "shared", Type: "secrets", Labels: map[string]string{"provider": "shared", "engine.io/composition": "secrets"}, CreatedAt: now},
		{Name: "shared-dns-server", Provider: "shared", Type: "dns", Labels: map[string]string{"provider": "shared", "engine.io/composition": "dns"}, CreatedAt: now},
		{Name: "shared-loadbalancer-compute", Provider: "shared", Type: "loadbalancer", Labels: map[string]string{"provider": "shared", "engine.io/composition": "loadbalancer"}, CreatedAt: now},
	}
}
