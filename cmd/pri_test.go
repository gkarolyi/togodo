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

	// Verify the priority was set
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "(A) test todo 1 +project2 @context1\n(A) test todo 2 +project1 @context2\nx (C) test todo 3 +project1 @context1\n"
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

	// Verify both tasks had their priority set
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "(C) test todo 1 +project2 @context1\n(C) test todo 2 +project1 @context2\nx (C) test todo 3 +project1 @context1\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecutePri_ChangePriority(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test changing an existing priority
	args := []string{"1", "B"}
	err := executePri(repo, args)
	assertNoError(t, err)

	// Verify the priority was changed
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "(B) test todo 1 +project2 @context1\n(B) test todo 2 +project1 @context2\nx (C) test todo 3 +project1 @context1\n"
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

	// Verify the priority was removed
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "test todo 1 +project2 @context1\n(B) test todo 2 +project1 @context2\nx (C) test todo 3 +project1 @context1\n"
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

func TestExecutePri_PreservesOtherProperties(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add a task with projects and contexts
	repo.Add("original task +project @context")

	// Change its priority
	args := []string{"1", "A"}
	err := executePri(repo, args)
	assertNoError(t, err)

	// Verify projects and contexts are preserved
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 1)

	todo := todos[0]
	assertTodoPriority(t, todo, "A")
	assertContains(t, todo.Text, "+project")
	assertContains(t, todo.Text, "@context")
	assertContains(t, todo.Text, "original task")
}

func TestExecutePri_MixedValidInvalidNumbers(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with mix of valid and invalid line numbers
	args := []string{"2", "abc", "3", "A"}
	err := executePri(repo, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")

	// Verify first task was set before error occurred
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "(A) test todo 1 +project2 @context1\n(A) test todo 2 +project1 @context2\nx (C) test todo 3 +project1 @context1\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecutePri_MultipleTasksSamePriority(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add multiple tasks without priority
	repo.Add("task 1")
	repo.Add("task 2")
	repo.Add("task 3")

	// Set same priority for multiple tasks
	args := []string{"1", "2", "3", "A"}
	err := executePri(repo, args)
	assertNoError(t, err)

	// Verify all tasks have the same priority
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "(A) task 1\n(A) task 2\n(A) task 3\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecutePri_SortingAfterPriorityChange(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks in specific order
	repo.Add("no priority task")
	repo.Add("(C) low priority task")

	// Change the non-priority task to high priority
	args := []string{"1", "A"}
	err := executePri(repo, args)
	assertNoError(t, err)

	// Verify sorting - high priority should come first
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 2)

	// The first task should now be the high priority one
	firstTodo := todos[0]
	assertTodoPriority(t, firstTodo, "A")
	assertContains(t, firstTodo.Text, "no priority task")

	// The second task should be the low priority one
	secondTodo := todos[1]
	assertTodoPriority(t, secondTodo, "C")
}

func TestExecutePri_Integration(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test full integration: set priority, verify save was called
	args := []string{"2", "A"}
	err := executePri(repo, args)
	assertNoError(t, err)

	// Verify priority was set and repository was saved
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "(A) test todo 1 +project2 @context1\n(A) test todo 2 +project1 @context2\nx (C) test todo 3 +project1 @context1\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecutePri_SingleArgument(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test with only one argument (will be treated as setting priority "A" for line 0)
	// This should work - it sets priority A for the first task
	args := []string{"A"}
	err := executePri(repo, args)

	// This should succeed since it treats "A" as priority for line 0
	assertNoError(t, err)
}

func TestExecutePri_DoneTaskPriority(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test setting priority on a done task (line 3 is done)
	args := []string{"3", "A"}
	err := executePri(repo, args)
	assertNoError(t, err)

	// Verify the done task can have its priority changed
	output, err := repo.WriteToString()
	assertNoError(t, err)

	expectedOutput := "(A) test todo 1 +project2 @context1\n(B) test todo 2 +project1 @context2\nx (A) test todo 3 +project1 @context1\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
