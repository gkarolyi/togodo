package todotxtlib

import (
	"bytes"
	"io"
	"os"
)

// Writer handles writing Todo structs to an output destination
type Writer interface {
	Write(todos []Todo) error
}

// NewFileWriter returns a new FileWriter that writes to the specified file
func NewFileWriter(path string) Writer {
	return &fileWriter{
		path: path,
	}
}

// fileWriter is a Writer that writes Todo structs to a file
type fileWriter struct {
	path string
}

// Write writes the given todos to the file
func (w *fileWriter) Write(todos []Todo) error {
	file, err := os.OpenFile(w.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, todo := range todos {
		if _, err := file.WriteString(todo.Text + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// NewBufferWriter returns a new Writer that writes to an io.Writer
func NewBufferWriter(w io.Writer) Writer {
	return &bufferWriter{
		writer: w,
	}
}

// bufferWriter is a Writer that writes Todo structs to an io.Writer
type bufferWriter struct {
	writer io.Writer
}

// Write writes the given todos to the io.Writer
func (w *bufferWriter) Write(todos []Todo) error {
	// If the writer is a *bytes.Buffer, reset it to avoid appending
	if buf, ok := w.writer.(*bytes.Buffer); ok {
		buf.Reset()
	}

	for _, todo := range todos {
		if _, err := w.writer.Write([]byte(todo.Text + "\n")); err != nil {
			return err
		}
	}
	return nil
}
