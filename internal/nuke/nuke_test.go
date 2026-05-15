package nuke

import (
	"context"
	"testing"
)

type mockNukeProvider struct {
	resources []string
}

func (m *mockNukeProvider) ListResources(ctx context.Context, environment string) ([]string, error) {
	return m.resources, nil
}

func (m *mockNukeProvider) DeleteResource(ctx context.Context, resourceID string) error {
	return nil
}

func (m *mockNukeProvider) ValidateNuke(ctx context.Context, environment string) error {
	return nil
}

func TestNewNukeManager(t *testing.T) {
	nm := NewNukeManager()
	if nm == nil {
		t.Fatal("Expected non-nil nuke manager")
	}
	if nm.providers == nil {
		t.Error("Expected providers map to be initialized")
	}
	if nm.operations == nil {
		t.Error("Expected operations map to be initialized")
	}
}

func TestNukeManager_RegisterProvider(t *testing.T) {
	nm := NewNukeManager()
	provider := &mockNukeProvider{resources: []string{"resource-1", "resource-2"}}

	nm.RegisterProvider("test", provider)

	if len(nm.providers) != 1 {
		t.Errorf("Expected 1 provider, got %d", len(nm.providers))
	}
}

func TestNukeManager_NukeEnvironment(t *testing.T) {
	nm := NewNukeManager()
	provider := &mockNukeProvider{resources: []string{"resource-1"}}
	nm.RegisterProvider("test", provider)

	op, err := nm.NukeEnvironment(context.Background(), "test", "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if op == nil {
		t.Error("Expected non-nil operation")
	}
	if op.Environment != "test" {
		t.Errorf("Expected test, got %s", op.Environment)
	}
	if op.Provider != "test" {
		t.Errorf("Expected test, got %s", op.Provider)
	}
}

func TestNukeManager_GetOperation(t *testing.T) {
	nm := NewNukeManager()
	provider := &mockNukeProvider{resources: []string{}}
	nm.RegisterProvider("test", provider)

	op, _ := nm.NukeEnvironment(context.Background(), "test", "test")

	retrieved, err := nm.GetOperation(op.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrieved.ID != op.ID {
		t.Errorf("Expected %s, got %s", op.ID, retrieved.ID)
	}
}

func TestNukeManager_GetOperations(t *testing.T) {
	nm := NewNukeManager()
	provider := &mockNukeProvider{resources: []string{}}
	nm.RegisterProvider("test", provider)

	nm.NukeEnvironment(context.Background(), "test", "test")

	ops := nm.GetOperations()
	if len(ops) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(ops))
	}
}

func TestNukeManager_CancelOperation(t *testing.T) {
	nm := NewNukeManager()
	provider := &mockNukeProvider{resources: []string{}}
	nm.RegisterProvider("test", provider)

	op, _ := nm.NukeEnvironment(context.Background(), "test", "test")

	err := nm.CancelOperation(op.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
