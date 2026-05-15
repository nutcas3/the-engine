# Plugin System

The Sovereign Engine plugin system allows you to extend functionality with custom providers and integrations.

## Creating a Custom Plugin

1. Create a new file in `pkg/plugins/` (e.g., `my_provider.go`)

2. Implement the `Plugin` interface:

```go
package plugins

import (
    "context"
    "fmt"
)

type MyProviderPlugin struct {
    config map[string]interface{}
}

func (p *MyProviderPlugin) Name() string {
    return "my-provider"
}

func (p *MyProviderPlugin) Version() string {
    return "1.0.0"
}

func (p *MyProviderPlugin) Description() string {
    return "My custom provider plugin"
}

func (p *MyProviderPlugin) Initialize(config map[string]interface{}) error {
    p.config = config
    return nil
}

func (p *MyProviderPlugin) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    // Your plugin logic here
    return map[string]interface{}{
        "result": "success",
    }, nil
}

func (p *MyProviderPlugin) Cleanup() error {
    // Cleanup logic
    return nil
}

func NewMyProviderPlugin() Plugin {
    return &MyProviderPlugin{}
}
```

3. Register your plugin in the plugin manager:

```go
manager := plugins.NewPluginManager()
plugin := plugins.NewMyProviderPlugin()
err := manager.Register(plugin)
```

4. Execute your plugin:

```go
result, err := manager.Execute(ctx, "my-provider", input)
```

## Plugin Use Cases

- **Custom Cloud Providers**: Add support for providers not in the default set
- **Custom Notifications**: Send notifications to custom webhook endpoints
- **Custom Metrics**: Collect and report business-specific metrics
- **Custom Cleanup Logic**: Implement specialized cleanup rules for your environment

## Example Plugins

See `example_provider.go` for a complete example of a custom provider plugin.
