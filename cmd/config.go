package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/internal/injector"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "View or set configuration options",
	Long: `View or set configuration options for togodo.\n\nExamples:\n  togodo config                    # Show all configuration\n  togodo config todo_txt_path      			  # Show specific config value\n  togodo config todo_txt_path ~/my-todos.txt  # Set config value\n\nConfiguration is stored in ~/.config/togodo.toml`,
	Args: cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		presenter := injector.CreateCLIPresenter()
		return executeConfig(presenter, args)
	},
}

func executeConfig(presenter *cli.Presenter, args []string) error {
	switch len(args) {
	case 0:
		return showAllConfig(presenter)
	case 1:
		return showConfig(presenter, args[0])
	case 2:
		return setConfig(presenter, args[0], args[1])
	default:
		return fmt.Errorf("too many arguments")
	}
}

func showAllConfig(presenter *cli.Presenter) error {
	settings := viper.AllSettings()
	if len(settings) == 0 {
		presenter.WriteLine("No configuration found")
		return nil
	}

	for key, value := range settings {
		presenter.WriteLine(fmt.Sprintf("%s = %v", key, value))
	}
	return nil
}

func showConfig(presenter *cli.Presenter, key string) error {
	if !viper.IsSet(key) {
		presenter.WriteLine(fmt.Sprintf("Configuration key '%s' is not set", key))
		return nil
	}

	value := viper.Get(key)
	presenter.WriteLine(fmt.Sprintf("%s = %v", key, value))
	return nil
}

func setConfig(presenter *cli.Presenter, key, value string) error {
	// Validate the key (only allow known configuration keys)
	validKeys := map[string]bool{
		"todo_txt_path": true,
	}

	if !validKeys[key] {
		return fmt.Errorf("invalid configuration key '%s'. Valid keys: %s",
			key,
			strings.Join(getValidKeys(validKeys), ", "))
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

	presenter.WriteLine(fmt.Sprintf("Set %s = %s", key, value))
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
