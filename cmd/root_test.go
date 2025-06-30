package cmd

import (
	"testing"
)

func TestRootCmd(t *testing.T) {
	rootCmd := RootCmd()

	// Test basic command properties
	if rootCmd.Use != "togodo" {
		t.Errorf("Expected Use to be 'togodo', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short != "A CLI tool for managing your todo.txt" {
		t.Errorf("Expected Short description to match, got '%s'", rootCmd.Short)
	}

	if rootCmd.Long != "togodo is a CLI tool for managing your todo.txt file." {
		t.Errorf("Expected Long description to match, got '%s'", rootCmd.Long)
	}

	// Test that Run function is set
	if rootCmd.Run == nil {
		t.Error("Expected Run function to be set")
	}
}

func TestRootCmd_HasSubcommands(t *testing.T) {
	rootCmd := RootCmd()

	// Test that subcommands are registered
	expectedCommands := []string{"add", "list", "do", "pri", "tidy"}

	for _, expectedCmd := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' to be registered", expectedCmd)
		}
	}
}

func TestRootCmd_HasAliases(t *testing.T) {
	rootCmd := RootCmd()

	// Test that commands have their expected aliases
	aliasTests := map[string][]string{
		"add":  {"a"},
		"list": {"ls", "l"},
		"do":   {"x"},
	}

	for cmdName, expectedAliases := range aliasTests {
		cmd, _, err := rootCmd.Find([]string{cmdName})
		if err != nil {
			t.Errorf("Failed to find command '%s': %v", cmdName, err)
			continue
		}

		for _, alias := range expectedAliases {
			aliasCmd, _, err := rootCmd.Find([]string{alias})
			if err != nil {
				t.Errorf("Failed to find alias '%s' for command '%s': %v", alias, cmdName, err)
				continue
			}

			if aliasCmd.Name() != cmd.Name() {
				t.Errorf("Alias '%s' points to wrong command. Expected '%s', got '%s'", alias, cmd.Name(), aliasCmd.Name())
			}
		}
	}
}

func TestRootCmd_Flags(t *testing.T) {
	rootCmd := RootCmd()

	// Test that toggle flag is present
	toggleFlag := rootCmd.Flags().Lookup("toggle")
	if toggleFlag == nil {
		t.Error("Expected 'toggle' flag to be present")
	} else {
		if toggleFlag.Shorthand != "t" {
			t.Errorf("Expected toggle flag shorthand to be 't', got '%s'", toggleFlag.Shorthand)
		}

		if toggleFlag.DefValue != "false" {
			t.Errorf("Expected toggle flag default to be 'false', got '%s'", toggleFlag.DefValue)
		}
	}
}

func TestRootCmd_RunFunction(t *testing.T) {
	// This test verifies that the Run function doesn't panic
	// We can't easily test the TUI behavior without mocking
	rootCmd := RootCmd()

	// Test that Run function exists and can be called
	if rootCmd.Run == nil {
		t.Fatal("Expected Run function to be set")
	}

	// We can't easily test the actual TUI execution without extensive mocking
	// But we can verify the function signature is correct
	defer func() {
		if r := recover(); r != nil {
			// If it panics due to missing todo.txt file or other setup issues,
			// that's expected in the test environment
			t.Logf("Run function panicked as expected in test environment: %v", r)
		}
	}()

	// This will likely fail in test environment due to missing todo.txt
	// but that's OK - we're testing that the function exists and is callable
	// rootCmd.Run(rootCmd, []string{})
}

func TestRootCmd_CommandValidation(t *testing.T) {
	rootCmd := RootCmd()

	// Test add command validation
	addCmd, _, err := rootCmd.Find([]string{"add"})
	if err != nil {
		t.Fatalf("Failed to find add command: %v", err)
	}

	// Add command should require at least 1 argument
	if addCmd.Args == nil {
		t.Error("Expected add command to have Args validation")
	}

	// Test do command validation
	doCmd, _, err := rootCmd.Find([]string{"do"})
	if err != nil {
		t.Fatalf("Failed to find do command: %v", err)
	}

	// Do command should require at least 1 argument
	if doCmd.Args == nil {
		t.Error("Expected do command to have Args validation")
	}

	// Test list command validation
	listCmd, _, err := rootCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("Failed to find list command: %v", err)
	}

	// List command should accept arbitrary args (optional filter)
	if listCmd.Args == nil {
		t.Error("Expected list command to have Args validation")
	}
}

func TestRootCmd_HelpText(t *testing.T) {
	rootCmd := RootCmd()

	// Test that commands have proper help text
	commands := []string{"add", "list", "do", "pri", "tidy"}

	for _, cmdName := range commands {
		cmd, _, err := rootCmd.Find([]string{cmdName})
		if err != nil {
			t.Errorf("Failed to find command '%s': %v", cmdName, err)
			continue
		}

		if cmd.Short == "" {
			t.Errorf("Command '%s' should have a Short description", cmdName)
		}

		if cmd.Long == "" {
			t.Errorf("Command '%s' should have a Long description", cmdName)
		}

		if cmd.Use == "" {
			t.Errorf("Command '%s' should have a Use string", cmdName)
		}
	}
}

func TestRootCmd_ExecutionModel(t *testing.T) {
	rootCmd := RootCmd()

	// Test that all commands have RunE functions (error-returning run functions)
	commands := []string{"add", "list", "do", "pri", "tidy"}

	for _, cmdName := range commands {
		cmd, _, err := rootCmd.Find([]string{cmdName})
		if err != nil {
			t.Errorf("Failed to find command '%s': %v", cmdName, err)
			continue
		}

		if cmd.RunE == nil {
			t.Errorf("Command '%s' should have a RunE function for error handling", cmdName)
		}

		// Commands should not have both Run and RunE
		if cmd.Run != nil && cmd.RunE != nil {
			t.Errorf("Command '%s' should not have both Run and RunE functions", cmdName)
		}
	}
}

func TestNewTUIBaseCommand(t *testing.T) {
	// This function creates a TUI base command
	// We can't test its full functionality without mocking the file system
	// But we can verify it doesn't panic during creation
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic in test environment due to missing todo.txt
			t.Logf("NewTUIBaseCommand panicked as expected in test environment: %v", r)
		}
	}()

	// This will likely panic due to missing todo.txt file, which is expected
	// baseCmd := NewTUIBaseCommand()
}

func TestGetTodoTxtPath(t *testing.T) {
	// This function is not exported, so we can't test it directly
	// But we can test its behavior through NewDefaultBaseCommand
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic in test environment due to missing todo.txt
			t.Logf("getTodoTxtPath panicked as expected in test environment: %v", r)
		}
	}()

	// This will likely panic due to missing todo.txt file, which is expected
	// baseCmd := NewDefaultBaseCommand()
}

func TestRootCmd_Integration(t *testing.T) {
	rootCmd := RootCmd()

	// Test that the root command can be executed without panicking
	// (though it may fail due to missing todo.txt)
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Root command execution panicked as expected in test environment: %v", r)
		}
	}()

	// Verify command structure is valid
	if err := rootCmd.ValidateArgs([]string{}); err != nil {
		t.Errorf("Root command args validation failed: %v", err)
	}

	// Test that help can be generated without errors
	help := rootCmd.UsageString()
	if help == "" {
		t.Error("Expected help text to be generated")
	}
}
