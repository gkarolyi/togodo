package tests

import (
	"testing"
)

// TestMultilineAdd tests adding multiple tasks in one command
// Ported from: t2000-multiline.sh
func TestMultilineAdd(t *testing.T) {
	env := SetupTestEnv(t)

	// Check file is initially empty
	initialContent := env.ReadTodoFile()
	t.Logf("Initial file content: '%s'", initialContent)

	// Add task with newlines should create multiple tasks
	output, code := env.RunCommand("add", "task one\ntask two\ntask three")
	t.Logf("Add command output: '%s'", output)
	if code != 0 {
		t.Errorf("Expected exit code 0, got %d", code)
	}

	// Should show each task was added
	if output == "" {
		t.Error("Expected non-empty output")
	}

	// Should indicate multiple tasks added
	if !containsSubstring(output, "tasks added") {
		t.Errorf("Expected 'tasks added' message, got:\n%s", output)
	}

	// Verify tasks in file
	content := env.ReadTodoFile()

	// Should have 3 separate lines
	lines := countNonEmptyLines(content)
	if lines != 3 {
		t.Errorf("Expected 3 tasks in file, got %d\nContent:\n%s", lines, content)
	}

	// Should contain all three tasks
	if !containsSubstring(content, "task one") {
		t.Errorf("Expected 'task one' in file, got:\n%s", content)
	}
	if !containsSubstring(content, "task two") {
		t.Errorf("Expected 'task two' in file, got:\n%s", content)
	}
	if !containsSubstring(content, "task three") {
		t.Errorf("Expected 'task three' in file, got:\n%s", content)
	}
}

// Helper functions
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func countNonEmptyLines(s string) int {
	lines := 0
	for _, line := range splitLines(s) {
		if len(trimSpace(line)) > 0 {
			lines++
		}
	}
	return lines
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// TestMultilineHandling tests how multiline text is handled in tasks
// Ported from: t2000-multiline.sh
func TestMultilineHandling(t *testing.T) {
	env := SetupTestEnv(t)

	// Test 1: Tasks with embedded newlines should be split into separate tasks
	// This is the behavior implemented in TestMultilineAdd
	// Here we verify the file format remains valid after multiline add
	env.WriteTodoFile("")
	_, code := env.RunCommand("add", "first line\nsecond line")
	if code != 0 {
		t.Errorf("Expected exit code 0, got %d", code)
	}

	// Verify file has two separate tasks
	content := env.ReadTodoFile()
	if !containsSubstring(content, "first line") {
		t.Errorf("Expected 'first line' in file")
	}
	if !containsSubstring(content, "second line") {
		t.Errorf("Expected 'second line' in file")
	}

	// Test 2: List displays each task on its own line
	listOutput, code := env.RunCommand("list")
	if code != 0 {
		t.Errorf("Expected exit code 0 for list, got %d", code)
	}

	// Each task should be on its own line in the output
	lines := splitLines(listOutput)
	taskLines := 0
	for _, line := range lines {
		if containsSubstring(line, "first line") || containsSubstring(line, "second line") {
			taskLines++
		}
	}
	if taskLines != 2 {
		t.Errorf("Expected 2 task lines in list output, got %d\nOutput:\n%s", taskLines, listOutput)
	}

	// Test 3: Edit operations work correctly on tasks added from multiline input
	// Append to the first task
	appendOutput, code := env.RunCommand("append", "1", "appended text")
	if code != 0 {
		t.Errorf("Expected exit code 0 for append, got %d", code)
	}

	// Verify the append worked
	if !containsSubstring(appendOutput, "appended text") {
		t.Errorf("Expected append confirmation in output")
	}

	// Verify file content is still valid
	finalContent := env.ReadTodoFile()
	finalLines := countNonEmptyLines(finalContent)
	if finalLines != 2 {
		t.Errorf("Expected 2 tasks after append, got %d", finalLines)
	}
}
