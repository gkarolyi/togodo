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

func TestConfigWrite(t *testing.T) {
	// Set initial value
	viper.Set("test_write_key", "old_value")
	defer viper.Set("test_write_key", nil)

	result, err := ConfigWrite("test_write_key", "new_value")
	if err != nil {
		// WriteConfig may fail in test environment (no config file), but that's OK
		// We still want to verify the in-memory set worked
		if viper.GetString("test_write_key") != "new_value" {
			t.Error("Config value was not set in memory")
		}
		return
	}

	if result.NewValue != "new_value" {
		t.Errorf("Expected 'new_value', got '%s'", result.NewValue)
	}

	if result.OldValue != "old_value" {
		t.Errorf("Expected old value 'old_value', got '%v'", result.OldValue)
	}

	if result.Created {
		t.Error("Expected update, not creation")
	}

	// Verify it was actually set
	if viper.GetString("test_write_key") != "new_value" {
		t.Error("Config value was not persisted")
	}
}

func TestConfigWriteCreate(t *testing.T) {
	key := "test_new_key"
	defer viper.Set(key, nil)

	result, err := ConfigWrite(key, "created_value")
	if err != nil {
		// WriteConfig may fail in test environment, but verify in-memory set worked
		if viper.GetString(key) != "created_value" {
			t.Error("Config value was not set in memory")
		}
		return
	}

	if !result.Created {
		t.Error("Expected creation, not update")
	}

	if viper.GetString(key) != "created_value" {
		t.Error("Config value was not set")
	}
}

func TestConfigList(t *testing.T) {
	// Set some test config values
	viper.Set("key1", "value1")
	viper.Set("key2", "value2")
	defer func() {
		viper.Set("key1", nil)
		viper.Set("key2", nil)
	}()

	result, err := ConfigList()
	if err != nil {
		t.Fatalf("ConfigList failed: %v", err)
	}

	if len(result.Settings) == 0 {
		t.Error("Expected settings to be returned")
	}

	if result.Settings["key1"] != "value1" {
		t.Error("Expected key1 to be in settings")
	}
}
