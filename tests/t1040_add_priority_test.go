package tests

import (
	"regexp"
	"testing"
)

// TestAddWithPriority tests adding tasks with priority in the text
// Ported from: t1040-add-priority.sh
func TestAddWithPriority(t *testing.T) {
	env := SetupTestEnv(t)

	// Add a task with priority
	output, code := env.RunCommand("add", "(A) high priority task")
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d: %s", code, output)
	}

	// Read the file to verify priority was preserved
	fileContent := env.ReadTodoFile()

	// Verify the task has priority A and contains expected text
	if !regexp.MustCompile(`\(A\) high priority task`).MatchString(fileContent) {
		t.Errorf("Expected '(A) high priority task', got: %s", fileContent)
	}

	// Verify priority marker is present
	if !regexp.MustCompile(`^\(A\)`).MatchString(fileContent) {
		t.Errorf("Expected priority (A) at start of task, got: %s", fileContent)
	}

	// Verify output format matches todo.txt-cli
	// Format: "1 (A) high priority task"
	expectedOutput := "1 (A) high priority task\nTODO: 1 added."
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

// TestAddWithPriorityAndDate tests adding priority tasks with auto-dating enabled
func TestAddWithPriorityAndDate(t *testing.T) {
	env := SetupTestEnv(t)

	// Enable auto-dating
	_, code := env.RunCommand("config", "auto_add_creation_date", "true")
	if code != 0 {
		t.Fatalf("Failed to set config: exit code %d", code)
	}

	// Add a task with priority
	output, code := env.RunCommand("add", "(A) high priority task")
	if code != 0 {
		t.Fatalf("Expected exit code 0, got %d: %s", code, output)
	}

	// Read the file to verify date was added after priority
	fileContent := env.ReadTodoFile()

	// Verify the task has format: (A) YYYY-MM-DD text
	datePattern := `^\(A\) \d{4}-\d{2}-\d{2} high priority task`
	if !regexp.MustCompile(datePattern).MatchString(fileContent) {
		t.Errorf("Expected priority followed by date, got: %s", fileContent)
	}

	// Verify priority comes before date (not after)
	if regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \(A\)`).MatchString(fileContent) {
		t.Error("Date should not come before priority")
	}
}
