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
	t.Skip("TODO: Test multiline task text handling")

	// Test edge cases:
	// - Tasks with embedded newlines
	// - How list displays multiline tasks
	// - How edit operations handle multiline tasks
}
