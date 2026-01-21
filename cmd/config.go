package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// ConfigReadResult contains the result of reading a config value
type ConfigReadResult struct {
	Key   string
	Value interface{}
	Found bool
}

// ConfigRead reads a configuration value by key
func ConfigRead(key string) (ConfigReadResult, error) {
	if !viper.IsSet(key) {
		return ConfigReadResult{
			Key:   key,
			Found: false,
		}, fmt.Errorf("configuration key '%s' not found", key)
	}

	return ConfigReadResult{
		Key:   key,
		Value: viper.Get(key),
		Found: true,
	}, nil
}

// ConfigWriteResult contains the result of writing a config value
type ConfigWriteResult struct {
	Key      string
	OldValue interface{}
	NewValue string
	Created  bool
}

// ConfigWrite writes a configuration value by key
func ConfigWrite(key string, value string) (ConfigWriteResult, error) {
	oldValue := viper.Get(key)
	created := !viper.IsSet(key)

	viper.Set(key, value)

	// Ensure config directory exists before writing
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// No config file loaded, determine the path
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ConfigWriteResult{}, fmt.Errorf("failed to get home directory: %w", err)
		}
		configFile = filepath.Join(homeDir, ".config", "togodo", "config.toml")
	}

	// Create parent directory if it doesn't exist
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return ConfigWriteResult{}, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save config to file
	if err := viper.WriteConfig(); err != nil {
		// If config file doesn't exist, create it
		if err := viper.SafeWriteConfig(); err != nil {
			return ConfigWriteResult{}, fmt.Errorf("failed to write config: %w", err)
		}
	}

	return ConfigWriteResult{
		Key:      key,
		OldValue: oldValue,
		NewValue: value,
		Created:  created,
	}, nil
}

// ConfigListResult contains all configuration settings
type ConfigListResult struct {
	Settings map[string]interface{}
}

// ConfigList lists all configuration settings
func ConfigList() (ConfigListResult, error) {
	settings := viper.AllSettings()

	return ConfigListResult{
		Settings: settings,
	}, nil
}
