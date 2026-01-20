package tests

import (
	"strings"
	"testing"
)

// TestHelp tests the help command
// Ported from: t2100-help.sh
func TestHelp(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("help command shows usage", func(t *testing.T) {
		output, code := env.RunCommand("help")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// Should show list of available commands
		if !strings.Contains(output, "add") {
			t.Errorf("Help should mention 'add' command, got:\n%s", output)
		}
		if !strings.Contains(output, "list") {
			t.Errorf("Help should mention 'list' command, got:\n%s", output)
		}
	})

	t.Run("--help flag shows usage", func(t *testing.T) {
		output, code := env.RunCommand("--help")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// Should show same help as 'help' command
		if !strings.Contains(output, "Usage") || !strings.Contains(output, "togodo") {
			t.Errorf("Help should show usage information, got:\n%s", output)
		}
	})
}

// TestCommandHelp tests help for specific commands
// Ported from: t2110-help-action.sh
func TestCommandHelp(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("help for add command", func(t *testing.T) {
		output, code := env.RunCommand("add", "--help")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// Should show add command specific help
		if !strings.Contains(output, "add") {
			t.Errorf("Add help should mention the command, got:\n%s", output)
		}
	})

	t.Run("help for list command", func(t *testing.T) {
		output, code := env.RunCommand("list", "--help")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// Should show list command specific help
		if !strings.Contains(output, "list") {
			t.Errorf("List help should mention the command, got:\n%s", output)
		}
	})
}

// TestShortHelp tests short/condensed help output
// Ported from: t2120-shorthelp.sh
func TestShortHelp(t *testing.T) {
	t.Skip("TODO: Implement shorthelp command for condensed help")

	// env := SetupTestEnv(t)
	//
	// output, code := env.RunCommand("shorthelp")
	// // Should show abbreviated command list
}
