package todotxtlib

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFileWriter_Write(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		todos   []Todo
		want    string
		wantErr bool
	}{
		{
			name:    "empty todos",
			todos:   []Todo{},
			want:    "",
			wantErr: false,
		},
		{
			name: "single todo",
			todos: []Todo{
				{Text: "Buy groceries"},
			},
			want:    "Buy groceries\n",
			wantErr: false,
		},
		{
			name: "multiple todos",
			todos: []Todo{
				{Text: "Buy groceries"},
				{Text: "Call mom"},
				{Text: "Pay bills"},
			},
			want:    "Buy groceries\nCall mom\nPay bills\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tempFile := filepath.Join(tempDir, "test.todo.txt")
			writer := NewFileWriter(tempFile)

			// Test writing
			err := writer.Write(tt.todos)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Read back the file to verify contents
			content, err := os.ReadFile(tempFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			if string(content) != tt.want {
				t.Errorf("FileWriter.Write() content = %q, want %q", string(content), tt.want)
			}
		})
	}

	// Test writing to read-only directory
	t.Run("read-only directory", func(t *testing.T) {
		readOnlyDir := filepath.Join(tempDir, "readonly")
		if err := os.Mkdir(readOnlyDir, 0444); err != nil {
			t.Fatalf("Failed to create read-only directory: %v", err)
		}

		writer := NewFileWriter(filepath.Join(readOnlyDir, "test.todo.txt"))
		err := writer.Write([]Todo{{Text: "test"}})
		if err == nil {
			t.Error("FileWriter.Write() expected error for read-only directory, got nil")
		}
	})
}

func TestBufferWriter_Write(t *testing.T) {
	tests := []struct {
		name    string
		todos   []Todo
		want    string
		wantErr bool
	}{
		{
			name:    "empty todos",
			todos:   []Todo{},
			want:    "",
			wantErr: false,
		},
		{
			name: "single todo",
			todos: []Todo{
				{Text: "Buy groceries"},
			},
			want:    "Buy groceries\n",
			wantErr: false,
		},
		{
			name: "multiple todos",
			todos: []Todo{
				{Text: "Buy groceries"},
				{Text: "Call mom"},
				{Text: "Pay bills"},
			},
			want:    "Buy groceries\nCall mom\nPay bills\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewBufferWriter(&buf)

			err := writer.Write(tt.todos)
			if (err != nil) != tt.wantErr {
				t.Errorf("BufferWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if buf.String() != tt.want {
				t.Errorf("BufferWriter.Write() content = %q, want %q", buf.String(), tt.want)
			}
		})
	}
}
