package cmd

import (
	"testing"
)

func TestExecutePri_SingleTask(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test setting priority for task 2 (which has priority B)
	args := []string{"2", "A"}
	err := executePri(repo, args)
	assertNoError(t, err)

	output, err := repo.WriteToString()

	expectedOutput := "(A) test todo 1 +project2 @context1\n" +
		"(A) test todo 2 +project1 @context2\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecutePri_MultipleTasks(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test setting priority for multiple tasks (tasks 1 and 2)
	args := []string{"1", "2", "C"}
	err := executePri(repo, args)
	assertNoError(t, err)

	output, err := repo.WriteToString()

	expectedOutput := "(C) test todo 1 +project2 @context1\n" +
		"(C) test todo 2 +project1 @context2\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecutePri_RemovePriority(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test removing priority by setting empty string
	args := []string{"1", ""}
	err := executePri(repo, args)
	assertNoError(t, err)

	output, err := repo.WriteToString()

	expectedOutput := "test todo 1 +project2 @context1\n" +
		"(B) test todo 2 +project1 @context2\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecutePri_InvalidLineNumber(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with invalid line number (too high)
	args := []string{"10", "A"}
	err := executePri(repo, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to set priority for todo at line 10")
}

func TestExecutePri_InvalidLineNumberFormat(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with invalid line number format
	args := []string{"abc", "A"}
	err := executePri(repo, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")
}

func TestExecutePri_ZeroLineNumber(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with line number 0 (should fail - line numbers start at 1)
	args := []string{"0", "A"}
	err := executePri(repo, args)
	assertError(t, err)
}

func TestExecutePri_NegativeLineNumber(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with negative line number
	args := []string{"-1", "A"}
	err := executePri(repo, args)
	assertError(t, err)
}

func TestExecutePri_EmptyRepository(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test setting priority on empty repository
	args := []string{"1", "A"}
	err := executePri(repo, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to set priority for todo at line 1")
}

func TestExecutePri_MixedValidInvalidNumbers(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with mix of valid and invalid line numbers
	args := []string{"2", "abc", "3", "A"}
	err := executePri(repo, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")

	output, err := repo.WriteToString()

	expectedOutput := "(A) test todo 1 +project2 @context1\n" +
		"(A) test todo 2 +project1 @context2\n" +
		"x (C) test todo 3 +project1 @context1\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecutePri_DoneTaskPriority(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test setting priority on a done task (line 3 is done)
	args := []string{"3", "A"}
	err := executePri(repo, args)
	assertNoError(t, err)

	// Verify the done task can have its priority changed
	output, err := repo.WriteToString()

	expectedOutput := "(A) test todo 1 +project2 @context1\n" +
		"(B) test todo 2 +project1 @context2\n" +
		"x (A) test todo 3 +project1 @context1\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
