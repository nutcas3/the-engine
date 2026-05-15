package plugins

import (
	"context"
	"testing"
)

func TestExampleProviderPlugin_Name(t *testing.T) {
	provider := &ExampleProviderPlugin{}
	if provider.Name() != "example-provider" {
		t.Errorf("Expected example-provider, got %s", provider.Name())
	}
}

func TestExampleProviderPlugin_Version(t *testing.T) {
	provider := &ExampleProviderPlugin{}
	if provider.Version() != "1.0.0" {
		t.Errorf("Expected 1.0.0, got %s", provider.Version())
	}
}

func TestExampleProviderPlugin_Description(t *testing.T) {
	provider := &ExampleProviderPlugin{}
	if provider.Description() == "" {
		t.Error("Expected non-empty description")
	}
}

func TestExampleProviderPlugin_Initialize(t *testing.T) {
	provider := &ExampleProviderPlugin{}
	config := map[string]interface{}{"api_key": "test"}
	err := provider.Initialize(config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestExampleProviderPlugin_Execute(t *testing.T) {
	provider := &ExampleProviderPlugin{}
	provider.Initialize(map[string]interface{}{"api_key": "test"})

	input := map[string]interface{}{"tier": "small"}
	result, err := provider.Execute(context.Background(), input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result")
	}
}

func TestExampleProviderPlugin_Cleanup(t *testing.T) {
	provider := &ExampleProviderPlugin{}
	err := provider.Cleanup()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
