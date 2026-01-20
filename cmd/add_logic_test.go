package cmd

import (
	"testing"
)

func TestAdd(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	result := Add(repo, []string{"Buy milk", "Call dentist"})

	if len(result.Added) != 2 {
		t.Errorf("Expected 2 added todos, got %d", len(result.Added))
	}
	if result.Added[0].Text != "Buy milk" {
		t.Errorf("Expected first todo to be 'Buy milk', got '%s'", result.Added[0].Text)
	}
	if result.Added[1].Text != "Call dentist" {
		t.Errorf("Expected second todo to be 'Call dentist', got '%s'", result.Added[1].Text)
	}
	if result.Error != nil {
		t.Errorf("Expected no error, got %v", result.Error)
	}

	// Verify todos are in repository
	todos, err := repo.ListAll()
	if err != nil {
		t.Fatalf("Failed to list todos: %v", err)
	}
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos in repository, got %d", len(todos))
	}
}
