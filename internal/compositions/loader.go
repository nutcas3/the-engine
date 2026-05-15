package compositions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"the-engine/internal/types"

	"gopkg.in/yaml.v3"
)

// CompositionMetadata represents metadata extracted from composition YAML files
type CompositionMetadata struct {
	Name      string            `yaml:"name"`
	Provider  string            `yaml:"provider"`
	Type      string            `yaml:"type"`
	Labels    map[string]string `yaml:"labels"`
	CreatedAt string            `yaml:"createdAt"`
}

// LoadCompositionsFromYAML loads composition metadata from YAML files in the compositions directory
func LoadCompositionsFromYAML() ([]types.Composition, error) {
	compositionsDir := "compositions"
	var result []types.Composition

	err := filepath.Walk(compositionsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		// Read YAML file
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Parse YAML to extract metadata
		metadata, err := parseCompositionMetadata(data, path)
		if err != nil {
			// If parsing fails, create a basic entry from the filename
			metadata = createMetadataFromPath(path)
		}

		result = append(result, metadata)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk compositions directory: %w", err)
	}

	return result, nil
}

// parseCompositionMetadata extracts metadata from a composition YAML file
func parseCompositionMetadata(data []byte, path string) (types.Composition, error) {
	var metadata struct {
		Metadata struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
		Spec struct {
			CompositionRef struct {
				Name string `yaml:"name"`
			} `yaml:"compositionRef"`
		} `yaml:"spec"`
	}

	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return types.Composition{}, err
	}

	// Extract provider and type from path
	provider, compType := extractProviderAndType(path)

	return types.Composition{
		Name:      metadata.Metadata.Name,
		Provider:  provider,
		Type:      compType,
		Labels:    map[string]string{"engine.io/composition": compType, "provider": provider},
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// createMetadataFromPath creates composition metadata from the file path
func createMetadataFromPath(path string) types.Composition {
	provider, compType := extractProviderAndType(path)
	name := strings.TrimSuffix(filepath.Base(path), ".yaml")

	return types.Composition{
		Name:      name,
		Provider:  provider,
		Type:      compType,
		Labels:    map[string]string{"engine.io/composition": compType, "provider": provider},
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}

// extractProviderAndType extracts provider and type from a file path
func extractProviderAndType(path string) (string, string) {
	parts := strings.Split(path, string(filepath.Separator))
	
	var provider, compType string
	for i, part := range parts {
		if part == "compositions" && i+1 < len(parts) {
			provider = parts[i+1]
		}
		if strings.Contains(part, ".yaml") {
			compType = strings.TrimSuffix(part, ".yaml")
		}
	}

	if provider == "" {
		provider = "unknown"
	}
	if compType == "" {
		compType = "unknown"
	}

	return provider, compType
}
