package tests

import (
	"regexp"
	"testing"
)

// TestAddWithDate tests that adding tasks includes creation date when enabled
// Ported from: t1010-add-date.sh
func TestAddWithDate(t *testing.T) {
	env := SetupTestEnv(t)

	// Enable auto-dating via config
	_, code := env.RunCommand("config", "auto_add_creation_date", "true")
	if code != 0 {
		t.Fatalf("Failed to set config: exit code %d", code)
	}

	// Add a task without an explicit date
	output, code := env.RunCommand("add", "task with auto date")
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d: %s", code, output)
	}

	// Read the file to verify date was added
	fileContent := env.ReadTodoFile()

	// Verify the text contains a date in YYYY-MM-DD format
	datePattern := `\d{4}-\d{2}-\d{2}`
	if !regexp.MustCompile(datePattern).MatchString(fileContent) {
		t.Errorf("Expected date in task text, got: %s", fileContent)
	}

	// Verify the task text is present (after the date)
	if !regexp.MustCompile(`\d{4}-\d{2}-\d{2} task with auto date`).MatchString(fileContent) {
		t.Errorf("Expected task with date prefix, got: %s", fileContent)
	}
}

// TestAddWithoutDate tests that tasks don't get dates when auto-dating is disabled
func TestAddWithoutDate(t *testing.T) {
	env := SetupTestEnv(t)

	// Ensure auto-dating is disabled (default)
	_, code := env.RunCommand("config", "auto_add_creation_date", "false")
	if code != 0 {
		t.Fatalf("Failed to set config: exit code %d", code)
	}

	// Add a task
	output, code := env.RunCommand("add", "task without auto date")
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d: %s", code, output)
	}

	// Read the file to verify NO date was added
	fileContent := env.ReadTodoFile()

	// Verify the text does NOT contain a date in YYYY-MM-DD format
	datePattern := `\d{4}-\d{2}-\d{2}`
	if regexp.MustCompile(datePattern).MatchString(fileContent) {
		t.Errorf("Expected no date in task text, but found date in: %s", fileContent)
	}

	// Verify the task text is present (without date prefix)
	if !regexp.MustCompile(`task without auto date`).MatchString(fileContent) {
		t.Errorf("Expected task text 'task without auto date', got: %s", fileContent)
	}
}
