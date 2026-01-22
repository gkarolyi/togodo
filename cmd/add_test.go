package cmd

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
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

func TestAddCreationDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantDate bool // Whether output should contain a date
	}{
		{
			name:     "simple task gets date",
			input:    "Buy groceries",
			wantDate: true,
		},
		{
			name:     "task with priority gets date after priority",
			input:    "(A) Important task",
			wantDate: true,
		},
		{
			name:     "task with existing date unchanged",
			input:    "2024-01-15 Task with date",
			wantDate: true,
		},
		{
			name:     "priority task with existing date unchanged",
			input:    "(A) 2024-01-15 Task with priority and date",
			wantDate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addCreationDate(tt.input)

			if tt.wantDate {
				// Check if result contains a date in YYYY-MM-DD format anywhere in the string
				datePattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
				if !datePattern.MatchString(result) {
					t.Errorf("Expected date in result, got: %s", result)
				}
			}

			// Verify priority is preserved if present
			if strings.HasPrefix(tt.input, "(") {
				if !strings.HasPrefix(result, "(") {
					t.Errorf("Expected priority to be preserved, got: %s", result)
				}
			}
		})
	}
}

func TestAdd(t *testing.T) {
	t.Run("adds single task", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute
		result, err := Add(repo, []string{"test", "task"}, false)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Todo.Text != "test task" {
			t.Errorf("expected 'test task', got '%s'", result.Todo.Text)
		}
		if result.LineNumber != 1 {
			t.Errorf("expected line number 1, got %d", result.LineNumber)
		}
	})
}

func TestAddMultiple(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add multiple tasks from text with newlines
	text := "task one\ntask two\ntask three"
	result, err := AddMultiple(repo, text, false)
	if err != nil {
		t.Fatalf("AddMultiple failed: %v", err)
	}

	// Should have added 3 tasks
	if len(result.Todos) != 3 {
		t.Errorf("Expected 3 todos, got %d", len(result.Todos))
	}

	// Verify tasks in repository
	todos, _ := repo.ListAll()
	if len(todos) != 3 {
		t.Errorf("Expected 3 todos in repository, got %d", len(todos))
		for i, todo := range todos {
			t.Logf("Todo %d: %s", i, todo.Text)
		}
	}

	// Check task texts
	texts := make(map[string]bool)
	for _, todo := range todos {
		texts[todo.Text] = true
	}

	if !texts["task one"] {
		t.Error("Expected 'task one' in repository")
	}
	if !texts["task two"] {
		t.Error("Expected 'task two' in repository")
	}
	if !texts["task three"] {
		t.Error("Expected 'task three' in repository")
	}
}
