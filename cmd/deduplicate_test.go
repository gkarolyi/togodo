package cmd

import (
	"testing"
)

func TestDeduplicateCmd_WithDuplicates(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	// Add tasks with duplicates
	repo.Add("task one")
	repo.Add("task two")
	repo.Add("task one")
	repo.Add("task three")
	repo.Add("task two")

	// Execute deduplicate
	result, err := Deduplicate(repo)
	assertNoError(t, err)

	// Verify 2 duplicates were removed
	if result.RemovedCount != 2 {
		t.Errorf("Expected 2 removed duplicates, got %d", result.RemovedCount)
	}

	output := getRepositoryString(t, repo, buf)

	// Should only have 3 unique tasks now (alphabetically sorted)
	expectedOutput := "task one\n" +
		"task three\n" +
		"task two\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestDeduplicateCmd_NoDuplicates(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	// Add only unique tasks
	repo.Add("task 1")
	repo.Add("task 2")
	repo.Add("task 3")

	// Execute deduplicate
	result, err := Deduplicate(repo)
	assertNoError(t, err)

	// Verify no tasks were removed
	if result.RemovedCount != 0 {
		t.Errorf("Expected 0 removed duplicates, got %d", result.RemovedCount)
	}

	output := getRepositoryString(t, repo, buf)

	expectedOutput := "task 1\n" +
		"task 2\n" +
		"task 3\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestDeduplicateCmd_EmptyRepository(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Execute deduplicate on empty repository
	result, err := Deduplicate(repo)
	assertNoError(t, err)

	// Verify no tasks were removed
	if result.RemovedCount != 0 {
		t.Errorf("Expected 0 removed duplicates, got %d", result.RemovedCount)
	}
}

func TestDeduplicateCmd_KeepsFirstOccurrence(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	// Add same task multiple times
	repo.Add("repeated task")
	repo.Add("other task")
	repo.Add("repeated task")
	repo.Add("repeated task")

	// Execute deduplicate
	result, err := Deduplicate(repo)
	assertNoError(t, err)

	// Verify 2 duplicates were removed
	if result.RemovedCount != 2 {
		t.Errorf("Expected 2 removed duplicates, got %d", result.RemovedCount)
	}

	output := getRepositoryString(t, repo, buf)

	// Should have kept first occurrence of each (alphabetically sorted)
	expectedOutput := "other task\n" +
		"repeated task\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestDeduplicateCmd_CaseSensitive(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	// Add tasks with different cases
	repo.Add("Task One")
	repo.Add("task one")
	repo.Add("TASK ONE")

	// Execute deduplicate
	result, err := Deduplicate(repo)
	assertNoError(t, err)

	// Verify no tasks were removed (case-sensitive comparison)
	if result.RemovedCount != 0 {
		t.Errorf("Expected 0 removed duplicates, got %d", result.RemovedCount)
	}

	output := getRepositoryString(t, repo, buf)

	// All three should remain as they differ in case
	expectedOutput := "Task One\n" +
		"task one\n" +
		"TASK ONE\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestDeduplicateCmd_WithPriorities(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)

	// Add tasks with same text but different priorities
	repo.Add("(A) important task")
	repo.Add("(B) important task")
	repo.Add("important task")

	// Execute deduplicate
	result, err := Deduplicate(repo)
	assertNoError(t, err)

	// Verify no tasks were removed (text is different due to priority)
	if result.RemovedCount != 0 {
		t.Errorf("Expected 0 removed duplicates, got %d", result.RemovedCount)
	}

	output := getRepositoryString(t, repo, buf)

	// All three should remain as the text representation differs
	expectedOutput := "(A) important task\n" +
		"(B) important task\n" +
		"important task\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
