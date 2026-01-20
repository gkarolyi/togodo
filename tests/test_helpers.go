package tests

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// TestEnvironment holds the test setup for integration tests
type TestEnvironment struct {
	buffer      *bytes.Buffer
	output      *bytes.Buffer
	repo        todotxtlib.TodoRepository
	presenter   *cli.Presenter
	rootCmd     *cobra.Command
	t           *testing.T
}

// BufferOutputWriter captures output to a buffer for testing
type BufferOutputWriter struct {
	buffer *bytes.Buffer
}

func NewBufferOutputWriter(buffer *bytes.Buffer) *BufferOutputWriter {
	return &BufferOutputWriter{buffer: buffer}
}

func (w *BufferOutputWriter) WriteLine(line string) {
	w.buffer.WriteString(line + "\n")
}

func (w *BufferOutputWriter) WriteLines(lines []string) {
	for _, line := range lines {
		w.WriteLine(line)
	}
}

func (w *BufferOutputWriter) WriteError(err error) {
	w.buffer.WriteString("Error: " + err.Error() + "\n")
}

func (w *BufferOutputWriter) Run() error {
	return nil
}

// SetupTestEnv creates a new test environment with buffer-based repository
func SetupTestEnv(t *testing.T) *TestEnvironment {
	t.Helper()

	// Create buffer for todo.txt content
	buffer := &bytes.Buffer{}

	// Create output buffer to capture command output
	output := &bytes.Buffer{}

	// Create buffer-based reader and writer
	reader := todotxtlib.NewBufferReader(buffer)
	writer := todotxtlib.NewBufferWriter(buffer)

	// Create repository
	repo, err := todotxtlib.NewFileRepository(reader, writer)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create presenter with buffer output
	outputWriter := NewBufferOutputWriter(output)
	formatter := cli.NewPlainFormatter()
	presenter := cli.NewPresenterWithDeps(formatter, outputWriter)

	env := &TestEnvironment{
		buffer:    buffer,
		output:    output,
		repo:      repo,
		presenter: presenter,
		t:         t,
	}

	// Create root command with injected dependencies
	env.rootCmd = cmd.NewRootCmd(repo, env.presenter)

	return env
}

// RunCommand executes a command with args and returns output and exit code
func (env *TestEnvironment) RunCommand(args ...string) (string, int) {
	env.t.Helper()

	// Clear output buffer
	env.output.Reset()

	// Reset command state
	env.rootCmd = cmd.NewRootCmd(env.repo, env.presenter)
	env.rootCmd.SetArgs(args)
	env.rootCmd.SetOut(env.output)
	env.rootCmd.SetErr(env.output)

	// Execute command
	err := env.rootCmd.Execute()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}

	// Get output
	output := env.output.String()
	output = strings.TrimRight(output, "\n")

	return output, exitCode
}

// WriteTodoFile writes content to the todo.txt buffer
func (env *TestEnvironment) WriteTodoFile(content string) {
	env.t.Helper()

	// Clear buffer and write new content
	env.buffer.Reset()
	env.buffer.WriteString(content)

	// Recreate repository with new content
	reader := todotxtlib.NewBufferReader(env.buffer)
	writer := todotxtlib.NewBufferWriter(env.buffer)

	repo, err := todotxtlib.NewFileRepository(reader, writer)
	if err != nil {
		env.t.Fatalf("Failed to recreate repository: %v", err)
	}

	env.repo = repo

	// Recreate root command with new repository
	env.rootCmd = cmd.NewRootCmd(env.repo, env.presenter)
}

// ReadTodoFile returns the contents of the todo.txt buffer
func (env *TestEnvironment) ReadTodoFile() string {
	env.t.Helper()

	// Save to buffer to get current state
	if err := env.repo.Save(); err != nil {
		env.t.Fatalf("Failed to save repository: %v", err)
	}

	content := env.buffer.String()
	return strings.TrimRight(content, "\n")
}

// ReadDoneFile returns the contents of done.txt (not yet implemented)
func (env *TestEnvironment) ReadDoneFile() string {
	env.t.Helper()
	// TODO: Implement done.txt support
	return ""
}

// AssertOutput runs command and asserts the output matches expected
func (env *TestEnvironment) AssertOutput(expected string, args ...string) {
	env.t.Helper()
	output, exitCode := env.RunCommand(args...)
	if exitCode != 0 {
		env.t.Errorf("Command failed with exit code %d\nArgs: %v\nOutput: %s", exitCode, args, output)
	}
	if output != expected {
		env.t.Errorf("Output mismatch\nArgs: %v\nExpected:\n%s\n\nGot:\n%s", args, expected, output)
	}
}

// AssertOutputAndCode runs command and asserts both output and exit code
func (env *TestEnvironment) AssertOutputAndCode(expectedCode int, expected string, args ...string) {
	env.t.Helper()
	output, exitCode := env.RunCommand(args...)
	if exitCode != expectedCode {
		env.t.Errorf("Exit code mismatch\nArgs: %v\nExpected: %d\nGot: %d\nOutput: %s", args, expectedCode, exitCode, output)
	}
	if output != expected {
		env.t.Errorf("Output mismatch\nArgs: %v\nExpected:\n%s\n\nGot:\n%s", args, expected, output)
	}
}

// AssertExitCode runs command and asserts only the exit code
func (env *TestEnvironment) AssertExitCode(expectedCode int, args ...string) {
	env.t.Helper()
	_, exitCode := env.RunCommand(args...)
	if exitCode != expectedCode {
		env.t.Errorf("Exit code mismatch\nArgs: %v\nExpected: %d\nGot: %d", args, expectedCode, exitCode)
	}
}

// AssertFileContent checks that todo.txt contains expected content
func (env *TestEnvironment) AssertFileContent(expected string) {
	env.t.Helper()
	content := env.ReadTodoFile()
	if content != expected {
		env.t.Errorf("File content mismatch\nExpected:\n%s\n\nGot:\n%s", expected, content)
	}
}

// AssertContains checks that output contains expected substring
func (env *TestEnvironment) AssertContains(output, expected string) {
	env.t.Helper()
	if !strings.Contains(output, expected) {
		env.t.Errorf("Output does not contain expected substring\nExpected substring: %s\nGot: %s", expected, output)
	}
}
