package cmd

import (
	"fmt"

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

	// Only write to file if a config file path is configured
	// In tests, viper won't have a config file loaded, so we skip file writes
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		// Config file is loaded - write changes to disk
		if err := viper.WriteConfig(); err != nil {
			return ConfigWriteResult{}, fmt.Errorf("failed to write config: %w", err)
		}
	}
	// If no config file is loaded (e.g., in tests), just update in-memory values

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
