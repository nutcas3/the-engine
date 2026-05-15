package compositions

import (
	"the-engine/internal/types"
)

// GetCompositions returns the list of available compositions
// It loads from YAML files in the compositions directory
func GetCompositions() ([]types.Composition, error) {
	return LoadCompositionsFromYAML()
}
