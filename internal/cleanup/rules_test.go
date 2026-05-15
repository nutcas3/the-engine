package cleanup

import (
	"context"
	"testing"
)

func TestMatchesPattern(t *testing.T) {
	cm := &CleanupManager{}

	tests := []struct {
		name     string
		resource string
		pattern  string
		want     bool
	}{
		{
			name:     "exact match",
			resource: "database-primary",
			pattern:  "database-primary",
			want:     true,
		},
		{
			name:     "wildcard match",
			resource: "database-primary",
			pattern:  "database-*",
			want:     true,
		},
		{
			name:     "wildcard no match",
			resource: "cache-redis",
			pattern:  "database-*",
			want:     false,
		},
		{
			name:     "prefix match",
			resource: "essential-worker-1",
			pattern:  "essential-*",
			want:     true,
		},
		{
			name:     "no match",
			resource: "worker-1",
			pattern:  "essential-*",
			want:     false,
		},
		{
			name:     "empty pattern",
			resource: "test-resource",
			pattern:  "",
			want:     false,
		},
		{
			name:     "wildcard only",
			resource: "anything",
			pattern:  "*",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cm.matchesPattern(tt.resource, tt.pattern)
			if got != tt.want {
				t.Errorf("matchesPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldExclude(t *testing.T) {
	cm := &CleanupManager{}

	tests := []struct {
		name     string
		resource string
		patterns []string
		want     bool
	}{
		{
			name:     "no patterns",
			resource: "test-resource",
			patterns: []string{},
			want:     false,
		},
		{
			name:     "matching pattern",
			resource: "database-primary",
			patterns: []string{"database-*"},
			want:     true,
		},
		{
			name:     "no matching pattern",
			resource: "cache-redis",
			patterns: []string{"database-*", "essential-*"},
			want:     false,
		},
		{
			name:     "multiple patterns, one matches",
			resource: "essential-worker-1",
			patterns: []string{"database-*", "essential-*"},
			want:     true,
		},
		{
			name:     "exact match in patterns",
			resource: "production-db",
			patterns: []string{"production-db"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cm.shouldExclude(tt.resource, tt.patterns)
			if got != tt.want {
				t.Errorf("shouldExclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldShutdown(t *testing.T) {
	tests := []struct {
		name     string
		policy   *CleanupPolicy
		resource string
		want     bool
	}{
		{
			name: "auto shutdown disabled",
			policy: &CleanupPolicy{
				Name:          "test-policy",
				Environment:   EnvironmentTest,
				AutoShutdown:  false,
				ShutdownAfter: 3600,
			},
			resource: "test-resource",
			want:     false,
		},
		{
			name: "shutdown after is zero",
			policy: &CleanupPolicy{
				Name:          "test-policy",
				Environment:   EnvironmentTest,
				AutoShutdown:  true,
				ShutdownAfter: 0,
			},
			resource: "test-resource",
			want:     false,
		},
		{
			name: "resource excluded",
			policy: &CleanupPolicy{
				Name:            "test-policy",
				Environment:     EnvironmentTest,
				AutoShutdown:    true,
				ShutdownAfter:   3600,
				ExcludePatterns: []string{"database-*"},
			},
			resource: "database-primary",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CleanupManager{}
			got := cm.shouldShutdown(context.Background(), tt.resource, tt.policy)
			if got != tt.want {
				t.Errorf("shouldShutdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldNuke(t *testing.T) {
	tests := []struct {
		name   string
		policy *CleanupPolicy
		want   bool
	}{
		{
			name: "nuke after is zero",
			policy: &CleanupPolicy{
				Name:        "test-policy",
				Environment: EnvironmentTest,
				NukeAfter:   0,
			},
			want: false,
		},
		{
			name: "production environment",
			policy: &CleanupPolicy{
				Name:        "prod-policy",
				Environment: EnvironmentProd,
				NukeAfter:   86400,
			},
			want: false,
		},
		{
			name: "staging environment",
			policy: &CleanupPolicy{
				Name:        "staging-policy",
				Environment: EnvironmentStaging,
				NukeAfter:   86400,
			},
			want: false,
		},
		{
			name: "test environment",
			policy: &CleanupPolicy{
				Name:        "test-policy",
				Environment: EnvironmentTest,
				NukeAfter:   86400,
			},
			want: false, // Integration checks return false in CI without config
		},
		{
			name: "dev environment",
			policy: &CleanupPolicy{
				Name:        "dev-policy",
				Environment: EnvironmentDev,
				NukeAfter:   86400,
			},
			want: false, // Integration checks return false in CI without config
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CleanupManager{}
			got := cm.shouldNuke(context.Background(), tt.policy)
			if got != tt.want {
				t.Errorf("shouldNuke() = %v, want %v", got, tt.want)
			}
		})
	}
}
