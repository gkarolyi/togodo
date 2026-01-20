package cmd

import (
	"testing"
)

func TestDoCmd_SingleTask(t *testing.T) {
	repo, buf := setupTestRepository(t)

	// Test toggling a task that's not done (line 1)
	indices, err := parseLineNumbers([]string{"1"})
	assertNoError(t, err)

	todos := toggleTodos(t, repo, indices)

	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}

	output := getRepositoryString(t, repo, buf)

	expectedOutput := "(B) test todo 2 +project1 @context2\n" +
		"x (A) test todo 1 +project2 @context1\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestDoCmd_MultipleTask(t *testing.T) {
	repo, buf := setupTestRepository(t)

	// Test toggling multiple tasks
	indices, err := parseLineNumbers([]string{"1", "2"})
	assertNoError(t, err)

	todos := toggleTodos(t, repo, indices)

	if len(todos) != 2 {
		t.Fatalf("Expected 2 todos, got %d", len(todos))
	}

	output := getRepositoryString(t, repo, buf)

	expectedOutput := "x (A) test todo 1 +project2 @context1\n" +
		"x (B) test todo 2 +project1 @context2\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestDoCmd_ToggleAlreadyDone(t *testing.T) {
	repo, buf := setupTestRepository(t)

	// Test toggling a task to not done
	indices, err := parseLineNumbers([]string{"3"})
	assertNoError(t, err)

	todos := toggleTodos(t, repo, indices)

	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}

	output := getRepositoryString(t, repo, buf)

	expectedOutput := "(A) test todo 1 +project2 @context1\n" +
		"(B) test todo 2 +project1 @context2\n" +
		"(C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestDoCmd_InvalidLineNumber(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with invalid line number (too high)
	indices, err := parseLineNumbers([]string{"10"})
	assertNoError(t, err)

	_, err = toggleTodosWithError(repo, indices)
	assertError(t, err)
	assertContains(t, err.Error(), "index out of bounds")
}

func TestDoCmd_InvalidLineNumberFormat(t *testing.T) {
	// Test with invalid line number format
	_, err := parseLineNumbers([]string{"abc"})
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")
}

func TestDoCmd_ZeroLineNumber(t *testing.T) {
	// Test with line number 0 (should fail - line numbers start at 1)
	_, err := parseLineNumbers([]string{"0"})
	assertError(t, err)
	assertContains(t, err.Error(), "line number must be positive")
}

func TestDoCmd_NegativeLineNumber(t *testing.T) {
	// Test with negative line number
	_, err := parseLineNumbers([]string{"-1"})
	assertError(t, err)
	assertContains(t, err.Error(), "line number must be positive")
}

func TestDoCmd_EmptyRepository(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test toggling on empty repository
	indices, err := parseLineNumbers([]string{"1"})
	assertNoError(t, err)

	_, err = toggleTodosWithError(repo, indices)
	assertError(t, err)
	assertContains(t, err.Error(), "index out of bounds")
}

func TestDoCmd_MixedValidInvalidNumbers(t *testing.T) {
	// Test with mix of valid and invalid line numbers
	// parseLineNumbers should fail early before any toggling happens
	_, err := parseLineNumbers([]string{"1", "abc", "2"})
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")
}
