package plugins

import (
	"context"
	"errors"
)

var (
	ErrPluginAlreadyExists = errors.New("plugin already exists")
	ErrPluginNotFound      = errors.New("plugin not found")
)

type Plugin interface {
	Name() string

	Version() string

	Description() string

	Initialize(config map[string]any) error

	Execute(ctx context.Context, input map[string]any) (map[string]any, error)

	Cleanup() error
}

type PluginManager struct {
	plugins map[string]Plugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

func (pm *PluginManager) Register(plugin Plugin) error {
	if _, exists := pm.plugins[plugin.Name()]; exists {
		return ErrPluginAlreadyExists
	}
	pm.plugins[plugin.Name()] = plugin
	return nil
}

func (pm *PluginManager) Unregister(name string) error {
	plugin, exists := pm.plugins[name]
	if !exists {
		return ErrPluginNotFound
	}

	if err := plugin.Cleanup(); err != nil {
		return err
	}

	delete(pm.plugins, name)
	return nil
}

func (pm *PluginManager) Get(name string) (Plugin, error) {
	plugin, exists := pm.plugins[name]
	if !exists {
		return nil, ErrPluginNotFound
	}
	return plugin, nil
}

func (pm *PluginManager) List() []Plugin {
	plugins := make([]Plugin, 0, len(pm.plugins))
	for _, plugin := range pm.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

func (pm *PluginManager) Execute(ctx context.Context, name string, input map[string]any) (map[string]any, error) {
	plugin, err := pm.Get(name)
	if err != nil {
		return nil, err
	}
	return plugin.Execute(ctx, input)
}
