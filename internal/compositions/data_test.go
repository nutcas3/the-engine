package compositions

import (
	"os"
	"testing"
)

func TestGetCompositions(t *testing.T) {
	// Skip test if compositions directory doesn't exist
	if _, err := os.Stat("compositions"); os.IsNotExist(err) {
		t.Skip("compositions directory not found, skipping test")
	}

	compositions, err := GetCompositions()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(compositions) == 0 {
		t.Error("Expected at least one composition")
	}
}

func TestGetCompositions_YAMLLoading(t *testing.T) {
	// Skip test if compositions directory doesn't exist
	if _, err := os.Stat("compositions"); os.IsNotExist(err) {
		t.Skip("compositions directory not found, skipping test")
	}

	compositions, err := GetCompositions()
	if err != nil {
		t.Errorf("Expected no error loading from YAML, got %v", err)
	}

	// Verify compositions have expected fields
	for _, comp := range compositions {
		if comp.Name == "" {
			t.Error("Composition name should not be empty")
		}
		if comp.Provider == "" {
			t.Error("Composition provider should not be empty")
		}
		if comp.Type == "" {
			t.Error("Composition type should not be empty")
		}
	}
}
