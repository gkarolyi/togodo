package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

func TestConfigRead(t *testing.T) {
	// Setup test config
	viper.Set("test_key", "test_value")
	defer viper.Set("test_key", nil)

	result, err := ConfigRead("test_key")
	if err != nil {
		t.Fatalf("ConfigRead failed: %v", err)
	}

	if !result.Found {
		t.Error("Expected key to be found")
	}

	if result.Value != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", result.Value)
	}
}

func TestConfigReadNotFound(t *testing.T) {
	result, err := ConfigRead("nonexistent_key")
	if err == nil {
		t.Error("Expected error for nonexistent key")
	}

	if result.Found {
		t.Error("Expected key not to be found")
	}
}
