package cmd

import (
	"strings"
	"testing"
)

func TestAddCmd_SingleTask(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	// Test adding a single task
	todos := addTodos(t, repo, []string{"(A) new task +project @context"})

	// Verify the returned todos
	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}

	// Verify the task was saved
	output := getRepositoryString(t, repo, buf)

	expectedOutput := "(A) new task +project @context\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestAddCmd_MultipleTasks(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	todos := addTodos(t, repo, []string{
		"(A) first task +project1 @context1",
		"(B) second task +project2 @context2",
		"third task without priority",
	})

	// Verify the returned todos
	if len(todos) != 3 {
		t.Fatalf("Expected 3 todos, got %d", len(todos))
	}

	// Check that tasks are in correct sorted order
	output := getRepositoryString(t, repo, buf)

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

func TestAddCmd_EmptyArgs(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	// Test adding with empty args (should not add anything)
	todos := addTodos(t, repo, []string{})

	// Verify no todos were returned
	if len(todos) != 0 {
		t.Fatalf("Expected 0 todos, got %d", len(todos))
	}

	// Verify no tasks were added
	output := getRepositoryString(t, repo, buf)

	expectedOutput := ""
	if output != expectedOutput {
		t.Errorf("Expected empty output, got:\n%s", output)
	}
}

func TestAddCmd_SortingBehavior(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	// Add tasks in mixed priority order
	tasks := []string{
		"no priority task",
		"(C) low priority task",
		"(A) high priority task",
		"(B) medium priority task",
	}

	for _, task := range tasks {
		addTodos(t, repo, []string{task})
	}

	// Verify the actual order using output
	output := getRepositoryString(t, repo, buf)

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

func TestAddCmd_MultilineString(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	// Test adding a task that contains newlines (should be treated as one task)
	todos := addTodos(t, repo, []string{"(A) task with\nnewline characters\nin the text"})

	// Verify one todo was returned
	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}

	// Verify only one task was added
	output := getRepositoryString(t, repo, buf)

	expectedOutput := "(A) task with\nnewline characters\nin the text\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestAddCmd_DuplicateTasks(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	// Test adding duplicate tasks
	todos := addTodos(t, repo, []string{
		"(A) duplicate task +project @context",
		"(A) duplicate task +project @context",
	})

	// Verify two todos were returned
	if len(todos) != 2 {
		t.Fatalf("Expected 2 todos, got %d", len(todos))
	}

	// Verify both duplicate tasks were added
	output := getRepositoryString(t, repo, buf)

	expectedOutput := "(A) duplicate task +project @context\n" +
		"(A) duplicate task +project @context\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
