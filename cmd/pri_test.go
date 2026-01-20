package cmd

import (
	"testing"
)

func TestPriCmd_SingleTask(t *testing.T) {
	repo, buf := setupTestRepository(t)

	// Test setting priority for task 2 (which has priority B)
	indices, priority, err := parsePriorityArgs([]string{"2", "A"})
	assertNoError(t, err)

	todos := setPriorities(t, repo, indices, priority)

	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}

	output := getRepositoryString(t, repo, buf)

	expectedOutput := "(A) test todo 1 +project2 @context1\n" +
		"(A) test todo 2 +project1 @context2\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestPriCmd_MultipleTasks(t *testing.T) {
	repo, buf := setupTestRepository(t)

	// Test setting priority for multiple tasks (tasks 1 and 2)
	indices, priority, err := parsePriorityArgs([]string{"1", "2", "C"})
	assertNoError(t, err)

	todos := setPriorities(t, repo, indices, priority)

	if len(todos) != 2 {
		t.Fatalf("Expected 2 todos, got %d", len(todos))
	}

	output := getRepositoryString(t, repo, buf)

	expectedOutput := "(C) test todo 1 +project2 @context1\n" +
		"(C) test todo 2 +project1 @context2\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestPriCmd_RemovePriority(t *testing.T) {
	repo, buf := setupTestRepository(t)

	// Test removing priority by setting empty string
	indices, priority, err := parsePriorityArgs([]string{"1", ""})
	assertNoError(t, err)

	todos := setPriorities(t, repo, indices, priority)

	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}

	output := getRepositoryString(t, repo, buf)

	expectedOutput := "test todo 1 +project2 @context1\n" +
		"(B) test todo 2 +project1 @context2\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestPriCmd_InvalidLineNumber(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with invalid line number (too high)
	indices, priority, err := parsePriorityArgs([]string{"10", "A"})
	assertNoError(t, err)

	_, err = setPrioritiesWithError(repo, indices, priority)
	assertError(t, err)
	assertContains(t, err.Error(), "index out of bounds")
}

func TestPriCmd_InvalidLineNumberFormat(t *testing.T) {
	// Test with invalid line number format
	_, _, err := parsePriorityArgs([]string{"abc", "A"})
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")
}

func TestPriCmd_ZeroLineNumber(t *testing.T) {
	// Test with line number 0 (should fail - line numbers start at 1)
	_, _, err := parsePriorityArgs([]string{"0", "A"})
	assertError(t, err)
	assertContains(t, err.Error(), "line number must be positive")
}

func TestPriCmd_NegativeLineNumber(t *testing.T) {
	// Test with negative line number
	_, _, err := parsePriorityArgs([]string{"-1", "A"})
	assertError(t, err)
	assertContains(t, err.Error(), "line number must be positive")
}

func TestPriCmd_EmptyRepository(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test setting priority on empty repository
	indices, priority, err := parsePriorityArgs([]string{"1", "A"})
	assertNoError(t, err)

	_, err = setPrioritiesWithError(repo, indices, priority)
	assertError(t, err)
	assertContains(t, err.Error(), "index out of bounds")
}

func TestPriCmd_MixedValidInvalidNumbers(t *testing.T) {
	// Test with mix of valid and invalid line numbers
	// parsePriorityArgs should fail early before any setting happens
	_, _, err := parsePriorityArgs([]string{"2", "abc", "3", "A"})
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")
}

func TestPriCmd_DoneTaskPriority(t *testing.T) {
	repo, buf := setupTestRepository(t)

	// Test setting priority on a done task (line 3 is done)
	indices, priority, err := parsePriorityArgs([]string{"3", "A"})
	assertNoError(t, err)

	todos := setPriorities(t, repo, indices, priority)

	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}

	// Verify the done task can have its priority changed
	output := getRepositoryString(t, repo, buf)

	expectedOutput := "(A) test todo 1 +project2 @context1\n" +
		"(B) test todo 2 +project1 @context2\n" +
		"x (A) test todo 3 +project1 @context1\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
