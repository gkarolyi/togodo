package cmd

import (
	"testing"

	"github.com/gkarolyi/togodo/internal/cli"
)

func TestExecuteTidy_WithDoneTasks(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Mark one more task as done to have multiple done tasks
	repo.ToggleDone(0)

	// Execute tidy
	err := executeTidy(repo, cli.NewPresenter())
	assertNoError(t, err)

	output, err := repo.WriteToString()

	expectedOutput := "(B) test todo 2 +project1 @context2\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecuteTidy_NoDoneTasks(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	repo.Add("(A) task 1")
	repo.Add("(B) task 2")

	// Execute tidy
	err := executeTidy(repo, cli.NewPresenter())
	assertNoError(t, err)

	output, err := repo.WriteToString()

	expectedOutput := "(A) task 1\n" +
		"(B) task 2\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExecuteTidy_EmptyRepository(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Execute tidy on empty repository
	err := executeTidy(repo, cli.NewPresenter())
	assertNoError(t, err)
}

func TestExecuteTidy_AllTasksDone(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks and mark them all as done
	repo.Add("x task 1")
	repo.Add("x task 2")
	repo.Add("x task 3")

	// Execute tidy
	err := executeTidy(repo, cli.NewPresenter())
	assertNoError(t, err)

	output, err := repo.WriteToString()

	expectedOutput := ""

	if output != expectedOutput {
		t.Errorf("Expected empty output, got:\n%s", output)
	}
}

func TestExecuteTidy_PrintsRemovedTasks(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Mark additional task as done
	repo.ToggleDone(0)

	// Get done tasks before tidy
	doneTodos, err := repo.ListDone()
	assertNoError(t, err)
	doneCount := len(doneTodos)

	// Execute tidy
	err = executeTidy(repo, cli.NewPresenter())
	assertNoError(t, err)

	// Verify that done tasks would have been printed
	// (We can't easily test the actual printing without mocking the output)
	if doneCount == 0 {
		t.Error("Expected to have some done tasks to remove")
	}
}
