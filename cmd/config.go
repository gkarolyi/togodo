package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	TodoTxtPath string `mapstructure:"todo_txt_path"`
}

// InitConfig initializes Viper configuration
func InitConfig() error {
	// Set config name and type
	viper.SetConfigName("togodo")
	viper.SetConfigType("toml")

	// Add config paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	viper.AddConfigPath(filepath.Join(homeDir, ".config"))    // ~/.config/togodo.toml
	viper.AddConfigPath(homeDir)                              // ~/.togodo.toml

	// Set default values
	viper.SetDefault("todo_txt_path", "todo.txt")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist, we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file %s: %w", viper.ConfigFileUsed(), err)
		}
	}

	return nil
}

// GetTodoTxtPath returns the configured todo.txt file path
func GetTodoTxtPath() string {
	path := viper.GetString("todo_txt_path")

	// Expand tilde to home directory if needed
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(homeDir, path[1:])
		}
	}

	return path
}
