package plugins

import (
	"context"
	"fmt"
)

// ExampleProviderPlugin is an example custom provider plugin
type ExampleProviderPlugin struct {
	config map[string]any
}

// Name returns the plugin name
func (p *ExampleProviderPlugin) Name() string {
	return "example-provider"
}

// Version returns the plugin version
func (p *ExampleProviderPlugin) Version() string {
	return "1.0.0"
}

// Description returns a description of what the plugin does
func (p *ExampleProviderPlugin) Description() string {
	return "Example provider plugin demonstrating how to create custom providers"
}

// Initialize sets up the plugin with configuration
func (p *ExampleProviderPlugin) Initialize(config map[string]any) error {
	p.config = config
	return nil
}

// Execute runs the plugin logic with the given context and input
func (p *ExampleProviderPlugin) Execute(ctx context.Context, input map[string]any) (map[string]any, error) {
	// Example: Map a custom provider's tier to instance types
	tier, ok := input["tier"].(string)
	if !ok {
		return nil, fmt.Errorf("tier not provided in input")
	}

	// Custom provider instance type mapping
	instanceTypes := map[string]string{
		"micro": "custom-micro",
		"small": "custom-small",
		"pro":   "custom-pro",
	}

	instanceType, ok := instanceTypes[tier]
	if !ok {
		instanceType = "custom-default"
	}

	return map[string]any{
		"instance_type": instanceType,
		"provider":      "custom-provider",
		"region":        p.config["region"],
	}, nil
}

// Cleanup performs any necessary cleanup when the plugin is unloaded
func (p *ExampleProviderPlugin) Cleanup() error {
	p.config = nil
	return nil
}

// NewExampleProviderPlugin creates a new example provider plugin
func NewExampleProviderPlugin() Plugin {
	return &ExampleProviderPlugin{}
}
