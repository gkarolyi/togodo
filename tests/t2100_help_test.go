package tests

import (
	"strings"
	"testing"
)

// TestHelpOutput tests basic help command output
// Ported from: t2100-help.sh "help output"
func TestHelpOutput(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("help shows required sections", func(t *testing.T) {
		output, code := env.RunCommand("help")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Expected section headers from upstream
		requiredSections := []string{
			"Usage:",
			"Options:",
			"Built-in Actions:",
		}

		for _, section := range requiredSections {
			if !strings.Contains(output, section) {
				t.Errorf("Help output missing required section: %s\nGot:\n%s", section, output)
			}
		}

		// Should NOT have Environment variables section (that's only in -vv)
		if strings.Contains(output, "Environment variables:") {
			t.Errorf("Basic help should not show Environment variables section")
		}
	})
}

// TestVerboseHelpOutput tests verbose help command output
// Ported from: t2100-help.sh "verbose help output"
func TestVerboseHelpOutput(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("help -v shows required sections", func(t *testing.T) {
		output, code := env.RunCommand("-v", "help")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Expected section headers (same as basic help)
		requiredSections := []string{
			"Usage:",
			"Options:",
			"Built-in Actions:",
		}

		for _, section := range requiredSections {
			if !strings.Contains(output, section) {
				t.Errorf("Verbose help output missing required section: %s\nGot:\n%s", section, output)
			}
		}

		// Should NOT have Environment variables section (that's only in -vv)
		if strings.Contains(output, "Environment variables:") {
			t.Errorf("Verbose help should not show Environment variables section")
		}
	})
}

// TestVeryVerboseHelpOutput tests very verbose help command output
// Ported from: t2100-help.sh "very verbose help output"
func TestVeryVerboseHelpOutput(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("help -vv shows all sections including environment", func(t *testing.T) {
		output, code := env.RunCommand("-vv", "help")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Expected section headers (including Environment variables)
		requiredSections := []string{
			"Usage:",
			"Options:",
			"Environment variables:",
			"Built-in Actions:",
		}

		for _, section := range requiredSections {
			if !strings.Contains(output, section) {
				t.Errorf("Very verbose help output missing required section: %s\nGot:\n%s", section, output)
			}
		}
	})
}

// TestHelpOutputWithCustomAction tests help output when custom actions exist
// Ported from: t2100-help.sh "help output with custom action"
func TestHelpOutputWithCustomAction(t *testing.T) {
	t.Skip("TODO: Implement custom action support")

	env := SetupTestEnv(t)

	// TODO: Create custom action "foo"
	// In upstream, this is done via make_action function which creates
	// a script in the actions directory

	t.Run("help -v shows add-on actions section", func(t *testing.T) {
		output, code := env.RunCommand("-v", "help")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Expected section headers (including Add-on Actions)
		requiredSections := []string{
			"Usage:",
			"Options:",
			"Built-in Actions:",
			"Add-on Actions:",
		}

		for _, section := range requiredSections {
			if !strings.Contains(output, section) {
				t.Errorf("Help with custom action missing required section: %s\nGot:\n%s", section, output)
			}
		}
	})
}
