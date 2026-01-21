package tests

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/internal/config"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TestEnvironment holds the test setup for integration tests
type TestEnvironment struct {
	buffer     *bytes.Buffer
	doneBuffer *bytes.Buffer
	output     *bytes.Buffer
	repo       todotxtlib.TodoRepository
	rootCmd    *cobra.Command
	t          *testing.T
	todoFile   string // Path to todo.txt file (for file-based tests)
	doneFile   string // Path to done.txt file (for file-based tests)
	tempDir    string // Temporary directory for file-based tests
	testTime   int64  // Mocked timestamp for TODO_TEST_TIME (Unix timestamp)
}

// SetupTestEnv creates a new test environment with buffer-based repository
func SetupTestEnv(t *testing.T) *TestEnvironment {
	t.Helper()

	// Reset config to defaults for test isolation
	// This prevents config changes from one test affecting others
	viper.Set("auto_add_creation_date", false)

	// Create buffer for todo.txt content
	buffer := &bytes.Buffer{}

	// Create buffer for done.txt content
	doneBuffer := &bytes.Buffer{}

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

	env := &TestEnvironment{
		buffer:     buffer,
		doneBuffer: doneBuffer,
		output:     output,
		repo:       repo,
		t:          t,
	}

	// Create root command with injected dependencies
	env.rootCmd = cli.NewRootCmd(repo)

	return env
}

// RunCommand executes a command with args and returns output and exit code
func (env *TestEnvironment) RunCommand(args ...string) (string, int) {
	env.t.Helper()

	// Clear output buffer
	env.output.Reset()

	// Reset command state
	env.rootCmd = cli.NewRootCmd(env.repo)
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

// SetTestDate sets a mocked date for testing (format: "2009-02-13")
// Compatible with todo.txt-cli's TODO_TEST_TIME environment variable
func (env *TestEnvironment) SetTestDate(dateStr string) error {
	env.t.Helper()

	// Parse the date string
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	// Store as Unix timestamp
	env.testTime = t.Unix()

	// Set environment variable for commands to use
	os.Setenv("TODO_TEST_TIME", fmt.Sprintf("%d", env.testTime))

	return nil
}

// TestTick advances the mocked time by the specified number of seconds
// Default is 86400 (one day), matching todo.txt-cli's test_tick
func (env *TestEnvironment) TestTick(seconds ...int64) {
	env.t.Helper()

	// Default to one day (86400 seconds)
	increment := int64(86400)
	if len(seconds) > 0 {
		increment = seconds[0]
	}

	env.testTime += increment

	// Update environment variable
	os.Setenv("TODO_TEST_TIME", fmt.Sprintf("%d", env.testTime))
}

// ClearTestDate removes the mocked date
func (env *TestEnvironment) ClearTestDate() {
	env.t.Helper()
	env.testTime = 0
	os.Unsetenv("TODO_TEST_TIME")
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
	env.rootCmd = cli.NewRootCmd(env.repo)
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

// ReadDoneFile returns the contents of done.txt
func (env *TestEnvironment) ReadDoneFile() string {
	env.t.Helper()
	content := env.doneBuffer.String()
	return strings.TrimRight(content, "\n")
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

// SetupFileBasedTestEnv creates a test environment using temporary files
// This is needed for commands like archive that work with multiple files
func SetupFileBasedTestEnv(t *testing.T) *TestEnvironment {
	t.Helper()

	// Create temporary directory
	tempDir := t.TempDir()
	todoFile := filepath.Join(tempDir, "todo.txt")
	doneFile := filepath.Join(tempDir, "done.txt")

	// Create empty todo.txt
	if err := os.WriteFile(todoFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create todo.txt: %v", err)
	}

	// Set config to use our temp file
	config.SetTodoTxtPath(todoFile)

	// Create output buffer to capture command output
	output := &bytes.Buffer{}

	// Create file-based reader and writer
	reader := todotxtlib.NewFileReader(todoFile)
	writer := todotxtlib.NewFileWriter(todoFile)

	// Create repository
	repo, err := todotxtlib.NewFileRepository(reader, writer)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	env := &TestEnvironment{
		output:   output,
		repo:     repo,
		t:        t,
		todoFile: todoFile,
		doneFile: doneFile,
		tempDir:  tempDir,
	}

	// Create root command with injected dependencies
	env.rootCmd = cli.NewRootCmd(repo)

	return env
}

// WriteTodoFileContent writes content to the todo.txt file (for file-based tests)
func (env *TestEnvironment) WriteTodoFileContent(content string) {
	env.t.Helper()

	if env.todoFile == "" {
		env.t.Fatal("This test environment is not file-based")
	}

	// Write content to file
	if err := os.WriteFile(env.todoFile, []byte(content), 0644); err != nil {
		env.t.Fatalf("Failed to write todo.txt: %v", err)
	}

	// Recreate repository with new content
	reader := todotxtlib.NewFileReader(env.todoFile)
	writer := todotxtlib.NewFileWriter(env.todoFile)

	repo, err := todotxtlib.NewFileRepository(reader, writer)
	if err != nil {
		env.t.Fatalf("Failed to recreate repository: %v", err)
	}

	env.repo = repo

	// Recreate root command with new repository
	env.rootCmd = cli.NewRootCmd(env.repo)
}

// ReadTodoFileContent returns the contents of the todo.txt file (for file-based tests)
func (env *TestEnvironment) ReadTodoFileContent() string {
	env.t.Helper()

	if env.todoFile == "" {
		env.t.Fatal("This test environment is not file-based")
	}

	content, err := os.ReadFile(env.todoFile)
	if err != nil {
		env.t.Fatalf("Failed to read todo.txt: %v", err)
	}

	return strings.TrimRight(string(content), "\n")
}

// ReadDoneFileContent returns the contents of the done.txt file (for file-based tests)
func (env *TestEnvironment) ReadDoneFileContent() string {
	env.t.Helper()

	if env.doneFile == "" {
		env.t.Fatal("This test environment is not file-based")
	}

	content, err := os.ReadFile(env.doneFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		env.t.Fatalf("Failed to read done.txt: %v", err)
	}

	return strings.TrimRight(string(content), "\n")
}
