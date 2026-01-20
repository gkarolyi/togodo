package cmd

import (
	"strings"
	"testing"
)

func TestExecuteAdd_SingleTask(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test adding a single task
	args := []string{"(A) new task +project @context"}
	_, err := executeAdd(repo, args)
	assertNoError(t, err)

	// Verify the task was added
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "(A) new task +project @context\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecuteAdd_MultipleTasks(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	args := []string{
		"(A) first task +project1 @context1",
		"(B) second task +project2 @context2",
		"third task without priority",
	}
	_, err := executeAdd(repo, args)

	assertNoError(t, err)

	// Check that tasks are in correct sorted order
	output, err := repo.WriteToString()
	assertNoError(t, err)

	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedOrder := []string{
		"(A) first task +project1 @context1",
		"(B) second task +project2 @context2",
		"third task without priority",
	}

	if len(lines) != len(expectedOrder) {
		t.Fatalf("Expected %d lines, got %d", len(expectedOrder), len(lines))
	}

	for i, expected := range expectedOrder {
		if lines[i] != expected {
			t.Errorf("Line %d: expected '%s', got '%s'", i+1, expected, lines[i])
		}
	}
}

func TestExecuteAdd_EmptyArgs(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test adding with empty args (should not add anything)
	args := []string{}
	_, err := executeAdd(repo, args)
	assertNoError(t, err)

	// Verify no tasks were added
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := ""
	if output != expectedOutput {
		t.Errorf("Expected empty output, got:\n%s", output)
	}
}

func TestExecuteAdd_SortingBehavior(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks in mixed priority order
	args := []string{
		"no priority task",
		"(C) low priority task",
		"(A) high priority task",
		"(B) medium priority task",
	}

	for _, arg := range args {
		_, err := executeAdd(repo, []string{arg})
		assertNoError(t, err)
	}

	// Verify the actual order using output
	output, err := repo.WriteToString()
	assertNoError(t, err)

	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedOrder := []string{
		"(A) high priority task",
		"(B) medium priority task",
		"(C) low priority task",
		"no priority task",
	}

	if len(lines) != len(expectedOrder) {
		t.Fatalf("Expected %d lines, got %d", len(expectedOrder), len(lines))
	}

	for i, expected := range expectedOrder {
		if lines[i] != expected {
			t.Errorf("Line %d: expected '%s', got '%s'", i+1, expected, lines[i])
		}
	}
}

func TestExecuteAdd_MultilineString(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test adding a task that contains newlines (should be treated as one task)
	args := []string{"(A) task with\nnewline characters\nin the text"}
	_, err := executeAdd(repo, args)
	assertNoError(t, err)

	// Verify only one task was added
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "(A) task with\nnewline characters\nin the text\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecuteAdd_DuplicateTasks(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test adding duplicate tasks
	args := []string{
		"(A) duplicate task +project @context",
		"(A) duplicate task +project @context",
	}
	_, err := executeAdd(repo, args)
	assertNoError(t, err)

	// Verify both duplicate tasks were added
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "(A) duplicate task +project @context\n" +
		"(A) duplicate task +project @context\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
