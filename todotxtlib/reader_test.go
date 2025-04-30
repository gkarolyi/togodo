package todotxtlib

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFileReader_Read(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		content string
		want    []Todo
		wantErr bool
	}{
		{
			name:    "empty file",
			content: "",
			want:    []Todo{},
			wantErr: false,
		},
		{
			name:    "single todo",
			content: "Buy groceries",
			want: []Todo{
				{Text: "Buy groceries"},
			},
			wantErr: false,
		},
		{
			name:    "multiple todos",
			content: "Buy groceries\nCall mom\nPay bills",
			want: []Todo{
				{Text: "Buy groceries"},
				{Text: "Call mom"},
				{Text: "Pay bills"},
			},
			wantErr: false,
		},
		{
			name:    "with empty lines",
			content: "Buy groceries\n\nCall mom\n\n",
			want: []Todo{
				{Text: "Buy groceries"},
				{Text: "Call mom"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tempFile := filepath.Join(tempDir, "test.todo.txt")
			if err := os.WriteFile(tempFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test reading
			reader := NewFileReader(tempFile)
			got, err := reader.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("FileReader.Read() got %d todos, want %d", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i].Text != tt.want[i].Text {
					t.Errorf("FileReader.Read() todo[%d] = %v, want %v", i, got[i].Text, tt.want[i].Text)
				}
			}
		})
	}

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		reader := NewFileReader(filepath.Join(tempDir, "nonexistent.todo.txt"))
		got, err := reader.Read()
		if err != nil {
			t.Errorf("FileReader.Read() error = %v, want nil", err)
		}
		if len(got) != 0 {
			t.Errorf("FileReader.Read() expected empty slice for non-existent file, got %v", got)
		}
	})
}

func TestBufferReader_Read(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []Todo
		wantErr bool
	}{
		{
			name:    "empty content",
			content: "",
			want:    []Todo{},
			wantErr: false,
		},
		{
			name:    "single todo",
			content: "Buy groceries",
			want: []Todo{
				{Text: "Buy groceries"},
			},
			wantErr: false,
		},
		{
			name:    "multiple todos",
			content: "Buy groceries\nCall mom\nPay bills",
			want: []Todo{
				{Text: "Buy groceries"},
				{Text: "Call mom"},
				{Text: "Pay bills"},
			},
			wantErr: false,
		},
		{
			name:    "with empty lines",
			content: "Buy groceries\n\nCall mom\n\n",
			want: []Todo{
				{Text: "Buy groceries"},
				{Text: "Call mom"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewBufferReader(bytes.NewBufferString(tt.content))
			got, err := reader.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("BufferReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("BufferReader.Read() got %d todos, want %d", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i].Text != tt.want[i].Text {
					t.Errorf("BufferReader.Read() todo[%d] = %v, want %v", i, got[i].Text, tt.want[i].Text)
				}
			}
		})
	}
}
