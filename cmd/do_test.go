package cmd

import (
	"testing"
)

func TestExecuteDo_SingleTask(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test toggling a task that's not done (line 1)
	args := []string{"1"}
	err := executeDo(repo, args)
	assertNoError(t, err)

	output, err := repo.WriteToString()

	expectedOutput := "(B) test todo 2 +project1 @context2\n" +
		"x (A) test todo 1 +project2 @context1\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecuteDo_MultipleTask(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test toggling multiple tasks
	args := []string{"1", "2"}
	err := executeDo(repo, args)
	assertNoError(t, err)

	output, err := repo.WriteToString()

	expectedOutput := "x (A) test todo 1 +project2 @context1\n" +
		"x (B) test todo 2 +project1 @context2\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecuteDo_ToggleAlreadyDone(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test toggling a task to not done
	args := []string{"3"}
	err := executeDo(repo, args)
	assertNoError(t, err)

	output, err := repo.WriteToString()

	expectedOutput := "(A) test todo 1 +project2 @context1\n" +
		"(B) test todo 2 +project1 @context2\n" +
		"(C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecuteDo_InvalidLineNumber(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with invalid line number (too high)
	args := []string{"10"}
	err := executeDo(repo, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to toggle todo at line 10")
}

func TestExecuteDo_InvalidLineNumberFormat(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with invalid line number format
	args := []string{"abc"}
	err := executeDo(repo, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")
}

func TestExecuteDo_ZeroLineNumber(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with line number 0 (should fail - line numbers start at 1)
	args := []string{"0"}
	err := executeDo(repo, args)
	assertError(t, err)
}

func TestExecuteDo_NegativeLineNumber(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with negative line number
	args := []string{"-1"}
	err := executeDo(repo, args)
	assertError(t, err)
}

func TestExecuteDo_EmptyRepository(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test toggling on empty repository
	args := []string{"1"}
	err := executeDo(repo, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to toggle todo at line 1")
}

func TestExecuteDo_MixedValidInvalidNumbers(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with mix of valid and invalid line numbers
	args := []string{"1", "abc", "2"}
	err := executeDo(repo, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")

	output, err := repo.WriteToString()

	expectedOutput := "x (A) test todo 1 +project2 @context1\n" +
		"(B) test todo 2 +project1 @context2\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
