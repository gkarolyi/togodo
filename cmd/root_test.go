package cmd

import (
	"testing"

	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestNewRootCmd(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := todotxtlib.NewTodoService(repo)
	presenter := cli.NewPresenter()

	rootCmd := NewRootCmd(service, repo, presenter)

	if rootCmd == nil {
		t.Fatal("NewRootCmd() returned nil")
	}

	if rootCmd.Use != "togodo" {
		t.Errorf("Expected Use to be 'togodo', got '%s'", rootCmd.Use)
	}

	if rootCmd.Run == nil {
		t.Error("Expected Run function to be set")
	}
}
