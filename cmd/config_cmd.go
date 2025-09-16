package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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

	viper.AddConfigPath(filepath.Join(homeDir, ".config")) // ~/.config/togodo.toml
	viper.AddConfigPath(homeDir)                           // ~/.togodo.toml

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


var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "View or set configuration options",
	Long: `View or set configuration options for togodo.

Examples:
  togodo config                    # Show all configuration
  togodo config todo_txt_path      # Show specific config value
  togodo config todo_txt_path ~/my-todos.txt  # Set config value

Configuration is stored in ~/.config/togodo.toml`,
	Args: cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		base := NewDefaultBaseCommand()
		return executeConfig(base, args)
	},
}

func executeConfig(base *BaseCommand, args []string) error {
	switch len(args) {
	case 0:
		return showAllConfig(base)
	case 1:
		return showConfig(base, args[0])
	case 2:
		return setConfig(base, args[0], args[1])
	default:
		return fmt.Errorf("too many arguments")
	}
}

func showAllConfig(base *BaseCommand) error {
	settings := viper.AllSettings()
	if len(settings) == 0 {
		base.Output.WriteLine("No configuration found")
		return nil
	}

	for key, value := range settings {
		base.Output.WriteLine(fmt.Sprintf("%s = %v", key, value))
	}
	return nil
}

func showConfig(base *BaseCommand, key string) error {
	if !viper.IsSet(key) {
		base.Output.WriteLine(fmt.Sprintf("Configuration key '%s' is not set", key))
		return nil
	}

	value := viper.Get(key)
	base.Output.WriteLine(fmt.Sprintf("%s = %v", key, value))
	return nil
}

func setConfig(base *BaseCommand, key, value string) error {
	// Validate the key (only allow known configuration keys)
	validKeys := map[string]bool{
		"todo_txt_path": true,
	}

	if !validKeys[key] {
		return fmt.Errorf("invalid configuration key '%s'. Valid keys: %s",
			key, strings.Join(getValidKeys(validKeys), ", "))
	}

	// Set the configuration value
	viper.Set(key, value)

	// Ensure config file exists
	if err := ensureConfigFile(); err != nil {
		return fmt.Errorf("error ensuring config file exists: %w", err)
	}

	// Write the configuration
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing configuration: %w", err)
	}

	base.Output.WriteLine(fmt.Sprintf("Set %s = %s", key, value))
	return nil
}

func getValidKeys(validKeys map[string]bool) []string {
	keys := make([]string, 0, len(validKeys))
	for key := range validKeys {
		keys = append(keys, key)
	}
	return keys
}

func ensureConfigFile() error {
	configFile := viper.ConfigFileUsed()

	// If no config file is currently being used, create one at ~/.config/togodo.toml
	if configFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}

		configDir := filepath.Join(homeDir, ".config")
		configFile = filepath.Join(configDir, "togodo.toml")

		// Create the .config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}

		// Set the config file path for viper
		viper.SetConfigFile(configFile)
	}

	// Check if config file exists, create it if it doesn't
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create empty config file
		file, err := os.Create(configFile)
		if err != nil {
			return fmt.Errorf("error creating config file: %w", err)
		}
		file.Close()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)
}
