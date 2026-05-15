package cleanup

import (
	"testing"
	"time"
)

func TestEnvironmentType(t *testing.T) {
	tests := []struct {
		name  string
		value EnvironmentType
		valid bool
	}{
		{
			name:  "dev environment",
			value: EnvironmentDev,
			valid: true,
		},
		{
			name:  "test environment",
			value: EnvironmentTest,
			valid: true,
		},
		{
			name:  "staging environment",
			value: EnvironmentStaging,
			valid: true,
		},
		{
			name:  "production environment",
			value: EnvironmentProd,
			valid: true,
		},
		{
			name:  "unknown environment",
			value: EnvironmentType("unknown"),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.valid && tt.value != "" {
				// Check if it's one of the known constants
				if tt.value != EnvironmentDev && tt.value != EnvironmentTest &&
					tt.value != EnvironmentStaging && tt.value != EnvironmentProd {
					// This is expected for unknown values
				}
			}
		})
	}
}

func TestCleanupPolicy(t *testing.T) {
	tests := []struct {
		name   string
		policy *CleanupPolicy
		valid  bool
	}{
		{
			name: "valid policy with all fields",
			policy: &CleanupPolicy{
				Name:            "test-policy",
				Environment:     EnvironmentDev,
				AutoShutdown:    true,
				ShutdownAfter:   3600 * time.Second,
				NukeAfter:        86400 * time.Second,
				Enabled:         true,
				ExcludePatterns: []string{"database-*", "essential-*"},
			},
			valid: true,
		},
		{
			name: "policy with empty name",
			policy: &CleanupPolicy{
				Name:            "",
				Environment:     EnvironmentDev,
				AutoShutdown:    true,
				ShutdownAfter:   3600 * time.Second,
				NukeAfter:        86400 * time.Second,
				Enabled:         true,
				ExcludePatterns: []string{},
			},
			valid: false,
		},
		{
			name: "policy with zero shutdown time",
			policy: &CleanupPolicy{
				Name:            "test-policy",
				Environment:     EnvironmentTest,
				AutoShutdown:    true,
				ShutdownAfter:   0,
				NukeAfter:        86400 * time.Second,
				Enabled:         true,
				ExcludePatterns: []string{},
			},
			valid: true,
		},
		{
			name: "policy disabled",
			policy: &CleanupPolicy{
				Name:            "test-policy",
				Environment:     EnvironmentStaging,
				AutoShutdown:    false,
				ShutdownAfter:   3600 * time.Second,
				NukeAfter:        86400 * time.Second,
				Enabled:         false,
				ExcludePatterns: []string{},
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				if tt.policy.Name == "" {
					t.Errorf("valid policy should have a name")
				}
			}
		})
	}
}

func TestResourceMetadata(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		metadata *ResourceMetadata
		valid    bool
	}{
		{
			name: "valid metadata",
			metadata: &ResourceMetadata{
				ID:          "resource-123",
				CreatedAt:   now,
				Environment: "dev",
				Tags:        map[string]string{"owner": "team-a", "purpose": "testing"},
			},
			valid: true,
		},
		{
			name: "metadata with empty ID",
			metadata: &ResourceMetadata{
				ID:          "",
				CreatedAt:   now,
				Environment: "test",
				Tags:        map[string]string{},
			},
			valid: false,
		},
		{
			name: "metadata with nil tags",
			metadata: &ResourceMetadata{
				ID:          "resource-123",
				CreatedAt:   now,
				Environment: "staging",
				Tags:        nil,
			},
			valid: true,
		},
		{
			name: "metadata with empty environment",
			metadata: &ResourceMetadata{
				ID:          "resource-123",
				CreatedAt:   now,
				Environment: "",
				Tags:        map[string]string{},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				if tt.metadata.ID == "" {
					t.Errorf("valid metadata should have an ID")
				}
				if tt.metadata.CreatedAt.IsZero() {
					t.Errorf("valid metadata should have a creation time")
				}
			}
		})
	}
}
