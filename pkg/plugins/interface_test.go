package plugins

import (
	"context"
	"errors"
	"testing"
)

type mockPlugin struct {
	name        string
	version     string
	description string
	config      map[string]interface{}
}

func (m *mockPlugin) Name() string {
	return m.name
}

func (m *mockPlugin) Version() string {
	return m.version
}

func (m *mockPlugin) Description() string {
	return m.description
}

func (m *mockPlugin) Initialize(config map[string]interface{}) error {
	m.config = config
	return nil
}

func (m *mockPlugin) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"result": "success"}, nil
}

func (m *mockPlugin) Cleanup() error {
	return nil
}

func TestNewPluginManager(t *testing.T) {
	pm := NewPluginManager()
	if pm == nil {
		t.Fatal("Expected non-nil plugin manager")
	}
	if pm.plugins == nil {
		t.Error("Expected plugins map to be initialized")
	}
}

func TestPluginManager_Register(t *testing.T) {
	pm := NewPluginManager()
	plugin := &mockPlugin{name: "test", version: "1.0", description: "test plugin"}

	err := pm.Register(plugin)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test duplicate registration
	err = pm.Register(plugin)
	if !errors.Is(err, ErrPluginAlreadyExists) {
		t.Errorf("Expected ErrPluginAlreadyExists, got %v", err)
	}
}

func TestPluginManager_Unregister(t *testing.T) {
	pm := NewPluginManager()
	plugin := &mockPlugin{name: "test", version: "1.0", description: "test plugin"}

	pm.Register(plugin)

	err := pm.Unregister("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test unregister non-existent
	err = pm.Unregister("nonexistent")
	if !errors.Is(err, ErrPluginNotFound) {
		t.Errorf("Expected ErrPluginNotFound, got %v", err)
	}
}

func TestPluginManager_Get(t *testing.T) {
	pm := NewPluginManager()
	plugin := &mockPlugin{name: "test", version: "1.0", description: "test plugin"}

	pm.Register(plugin)

	retrieved, err := pm.Get("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrieved.Name() != "test" {
		t.Errorf("Expected test, got %s", retrieved.Name())
	}

	// Test get non-existent
	_, err = pm.Get("nonexistent")
	if !errors.Is(err, ErrPluginNotFound) {
		t.Errorf("Expected ErrPluginNotFound, got %v", err)
	}
}

func TestPluginManager_List(t *testing.T) {
	pm := NewPluginManager()
	plugin1 := &mockPlugin{name: "test1", version: "1.0", description: "test plugin 1"}
	plugin2 := &mockPlugin{name: "test2", version: "1.0", description: "test plugin 2"}

	pm.Register(plugin1)
	pm.Register(plugin2)

	plugins := pm.List()
	if len(plugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(plugins))
	}
}

func TestPluginManager_Execute(t *testing.T) {
	pm := NewPluginManager()
	plugin := &mockPlugin{name: "test", version: "1.0", description: "test plugin"}

	pm.Register(plugin)

	result, err := pm.Execute(context.Background(), "test", map[string]interface{}{"input": "test"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result["result"] != "success" {
		t.Errorf("Expected success, got %v", result["result"])
	}

	// Test execute non-existent
	_, err = pm.Execute(context.Background(), "nonexistent", map[string]interface{}{})
	if !errors.Is(err, ErrPluginNotFound) {
		t.Errorf("Expected ErrPluginNotFound, got %v", err)
	}
}
