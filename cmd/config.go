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
