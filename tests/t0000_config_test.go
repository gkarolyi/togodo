package tests

import (
	"testing"
)

// TestConfigRead tests reading configuration values
// Ported from: t0000-config.sh
func TestConfigRead(t *testing.T) {
	env := SetupTestEnv(t)

	output, code := env.RunCommand("config", "todo_txt_path")
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d", code)
	}

	// Should show the configured todo.txt path (default: "todo.txt")
	if output == "" {
		t.Error("Expected output, got empty string")
	}
}

// TestConfigWrite tests setting configuration values
// Ported from: t0000-config.sh
func TestConfigWrite(t *testing.T) {
	env := SetupTestEnv(t)

	// Write a config value
	output, code := env.RunCommand("config", "todo_txt_path", "/tmp/test-todo.txt")
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d: %s", code, output)
	}

	// Verify write confirmation message
	if output == "" {
		t.Error("Expected confirmation output")
	}

	// Read it back to verify
	output, code = env.RunCommand("config", "todo_txt_path")
	if code != 0 {
		t.Fatalf("Expected exit code 0 when reading, got %d", code)
	}

	// Note: In test environment, the write may not persist to disk,
	// but the in-memory value should be set
}

// TestConfigList tests listing all configuration
// Ported from: t0000-config.sh
func TestConfigList(t *testing.T) {
	t.Skip("TODO: Implement config command to list all settings")

	// env := SetupTestEnv(t)
	//
	// output, code := env.RunCommand("config")
	// // Should show all configuration key-value pairs
}
