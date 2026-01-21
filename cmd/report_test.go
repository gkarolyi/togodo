package cmd

import (
	"bytes"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestReport(t *testing.T) {
	t.Run("generates statistics for mixed tasks", func(t *testing.T) {
		// Setup repository with mixed tasks
		buffer := &bytes.Buffer{}
		buffer.WriteString(`(A) task one
(B) task two
task three
x done task`)

		reader := todotxtlib.NewBufferReader(buffer)
		writer := todotxtlib.NewBufferWriter(buffer)
		repo, err := todotxtlib.NewFileRepository(reader, writer)
		if err != nil {
			t.Fatalf("Failed to create repository: %v", err)
		}

		// Setup done.txt with archived tasks
		doneBuffer := &bytes.Buffer{}
		doneBuffer.WriteString(`x 2024-01-01 archived task one
x 2024-01-02 archived task two`)

		doneReader := todotxtlib.NewBufferReader(doneBuffer)

		// Call Report
		result, err := Report(repo, doneReader)
		if err != nil {
			t.Fatalf("Report failed: %v", err)
		}

		// Verify counts
		if result.Total != 4 {
			t.Errorf("Expected total 4, got %d", result.Total)
		}
		if result.Done != 2 {
			t.Errorf("Expected done 2, got %d", result.Done)
		}
		if result.Pending != 3 {
			t.Errorf("Expected pending 3, got %d", result.Pending)
		}
	})

	t.Run("handles empty done.txt", func(t *testing.T) {
		// Setup repository
		buffer := &bytes.Buffer{}
		buffer.WriteString(`task one
task two`)

		reader := todotxtlib.NewBufferReader(buffer)
		writer := todotxtlib.NewBufferWriter(buffer)
		repo, err := todotxtlib.NewFileRepository(reader, writer)
		if err != nil {
			t.Fatalf("Failed to create repository: %v", err)
		}

		// Setup empty done.txt
		doneBuffer := &bytes.Buffer{}
		doneReader := todotxtlib.NewBufferReader(doneBuffer)

		// Call Report
		result, err := Report(repo, doneReader)
		if err != nil {
			t.Fatalf("Report failed: %v", err)
		}

		// Verify counts
		if result.Total != 2 {
			t.Errorf("Expected total 2, got %d", result.Total)
		}
		if result.Done != 0 {
			t.Errorf("Expected done 0, got %d", result.Done)
		}
		if result.Pending != 2 {
			t.Errorf("Expected pending 2, got %d", result.Pending)
		}
	})

	t.Run("handles all completed tasks", func(t *testing.T) {
		// Setup repository with all completed tasks
		buffer := &bytes.Buffer{}
		buffer.WriteString(`x done task one
x done task two`)

		reader := todotxtlib.NewBufferReader(buffer)
		writer := todotxtlib.NewBufferWriter(buffer)
		repo, err := todotxtlib.NewFileRepository(reader, writer)
		if err != nil {
			t.Fatalf("Failed to create repository: %v", err)
		}

		// Setup done.txt
		doneBuffer := &bytes.Buffer{}
		doneBuffer.WriteString(`x 2024-01-01 archived task`)

		doneReader := todotxtlib.NewBufferReader(doneBuffer)

		// Call Report
		result, err := Report(repo, doneReader)
		if err != nil {
			t.Fatalf("Report failed: %v", err)
		}

		// Verify counts
		if result.Total != 2 {
			t.Errorf("Expected total 2, got %d", result.Total)
		}
		if result.Done != 1 {
			t.Errorf("Expected done 1, got %d", result.Done)
		}
		if result.Pending != 0 {
			t.Errorf("Expected pending 0, got %d", result.Pending)
		}
	})

	t.Run("handles empty todo.txt", func(t *testing.T) {
		// Setup empty repository
		buffer := &bytes.Buffer{}
		reader := todotxtlib.NewBufferReader(buffer)
		writer := todotxtlib.NewBufferWriter(buffer)
		repo, err := todotxtlib.NewFileRepository(reader, writer)
		if err != nil {
			t.Fatalf("Failed to create repository: %v", err)
		}

		// Setup done.txt
		doneBuffer := &bytes.Buffer{}
		doneBuffer.WriteString(`x 2024-01-01 archived task`)

		doneReader := todotxtlib.NewBufferReader(doneBuffer)

		// Call Report
		result, err := Report(repo, doneReader)
		if err != nil {
			t.Fatalf("Report failed: %v", err)
		}

		// Verify counts
		if result.Total != 0 {
			t.Errorf("Expected total 0, got %d", result.Total)
		}
		if result.Done != 1 {
			t.Errorf("Expected done 1, got %d", result.Done)
		}
		if result.Pending != 0 {
			t.Errorf("Expected pending 0, got %d", result.Pending)
		}
	})
}
