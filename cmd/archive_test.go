package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestArchive(t *testing.T) {
	t.Run("archive completed todos", func(t *testing.T) {
		// Create temporary directory for test files
		tmpDir := t.TempDir()
		todoFile := filepath.Join(tmpDir, "todo.txt")
		doneFile := filepath.Join(tmpDir, "done.txt")

		// Create todo.txt with mixed completed and active todos
		initialContent := `task one
x completed task
task two
x another completed`

		err := os.WriteFile(todoFile, []byte(initialContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Create repository
		reader := todotxtlib.NewFileReader(todoFile)
		writer := todotxtlib.NewFileWriter(todoFile)
		repo, err := todotxtlib.NewFileRepository(reader, writer)
		if err != nil {
			t.Fatalf("Failed to create repository: %v", err)
		}

		// Create reader and writer for done.txt
		doneReader := todotxtlib.NewFileReader(doneFile)
		doneWriter := todotxtlib.NewFileWriter(doneFile)

		// Execute archive
		result, err := Archive(repo, doneReader, doneWriter)
		if err != nil {
			t.Fatalf("Archive failed: %v", err)
		}

		// Verify archived todos count
		if len(result.ArchivedTodos) != 2 {
			t.Errorf("Expected 2 archived todos, got %d", len(result.ArchivedTodos))
		}

		// Verify archived todos are done
		for _, todo := range result.ArchivedTodos {
			if !todo.Done {
				t.Errorf("Archived todo should be done: %s", todo.Text)
			}
		}

		// Verify todo.txt only contains active todos
		content, err := os.ReadFile(todoFile)
		if err != nil {
			t.Fatalf("Failed to read todo.txt: %v", err)
		}

		expectedTodoContent := "task one\ntask two\n"
		if string(content) != expectedTodoContent {
			t.Errorf("todo.txt content mismatch\nExpected:\n%s\nGot:\n%s", expectedTodoContent, string(content))
		}

		// Verify done.txt contains completed todos
		doneContent, err := os.ReadFile(doneFile)
		if err != nil {
			t.Fatalf("Failed to read done.txt: %v", err)
		}

		expectedDoneContent := "x completed task\nx another completed\n"
		if string(doneContent) != expectedDoneContent {
			t.Errorf("done.txt content mismatch\nExpected:\n%s\nGot:\n%s", expectedDoneContent, string(doneContent))
		}
	})

	t.Run("archive with no completed todos", func(t *testing.T) {
		// Setup buffer-based repository
		buffer := &bytes.Buffer{}
		buffer.WriteString("task one\ntask two\n")

		reader := todotxtlib.NewBufferReader(buffer)
		writer := todotxtlib.NewBufferWriter(buffer)
		repo, err := todotxtlib.NewFileRepository(reader, writer)
		if err != nil {
			t.Fatalf("Failed to create repository: %v", err)
		}

		// Create temporary directory for done.txt
		tmpDir := t.TempDir()
		doneFile := filepath.Join(tmpDir, "done.txt")

		// Create reader and writer for done.txt
		doneReader := todotxtlib.NewFileReader(doneFile)
		doneWriter := todotxtlib.NewFileWriter(doneFile)

		// Execute archive
		result, err := Archive(repo, doneReader, doneWriter)
		if err != nil {
			t.Fatalf("Archive failed: %v", err)
		}

		// Verify no todos were archived
		if len(result.ArchivedTodos) != 0 {
			t.Errorf("Expected 0 archived todos, got %d", len(result.ArchivedTodos))
		}
	})

	t.Run("archive appends to existing done.txt", func(t *testing.T) {
		// Create temporary directory for test files
		tmpDir := t.TempDir()
		todoFile := filepath.Join(tmpDir, "todo.txt")
		doneFile := filepath.Join(tmpDir, "done.txt")

		// Create todo.txt with completed todos
		initialContent := `x new completed task`
		err := os.WriteFile(todoFile, []byte(initialContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Create done.txt with existing content
		existingDoneContent := `x old completed task`
		err = os.WriteFile(doneFile, []byte(existingDoneContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write done.txt: %v", err)
		}

		// Create repository
		reader := todotxtlib.NewFileReader(todoFile)
		writer := todotxtlib.NewFileWriter(todoFile)
		repo, err := todotxtlib.NewFileRepository(reader, writer)
		if err != nil {
			t.Fatalf("Failed to create repository: %v", err)
		}

		// Create reader and writer for done.txt
		doneReader := todotxtlib.NewFileReader(doneFile)
		doneWriter := todotxtlib.NewFileWriter(doneFile)

		// Execute archive
		result, err := Archive(repo, doneReader, doneWriter)
		if err != nil {
			t.Fatalf("Archive failed: %v", err)
		}

		// Verify archived todos count
		if len(result.ArchivedTodos) != 1 {
			t.Errorf("Expected 1 archived todo, got %d", len(result.ArchivedTodos))
		}

		// Verify done.txt contains both old and new completed todos
		doneContent, err := os.ReadFile(doneFile)
		if err != nil {
			t.Fatalf("Failed to read done.txt: %v", err)
		}

		expectedDoneContent := "x old completed task\nx new completed task\n"
		if string(doneContent) != expectedDoneContent {
			t.Errorf("done.txt content mismatch\nExpected:\n%s\nGot:\n%s", expectedDoneContent, string(doneContent))
		}
	})
}
