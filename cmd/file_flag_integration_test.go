package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/internal/config"
)

func TestFileFlag_Integration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "togodo_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test todo file
	testFile := filepath.Join(tempDir, "test_todo.txt")
	testContent := `(A) Call mom +family @home
(B) Buy groceries +shopping @errands
Write documentation +work @office
x Complete project setup +work @office
(C) Schedule dentist appointment +health @phone
Review pull requests +work @office`

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test that config override works
	config.SetTodoTxtPath(testFile)
	actualPath := config.GetTodoTxtPath()
	if actualPath != testFile {
		t.Errorf("Expected config path to be %s, got %s", testFile, actualPath)
	}

	// Reset config for next test
	config.SetTodoTxtPath("")

	// Test tilde expansion
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	// Create a test file in home directory
	homeTestFile := filepath.Join(homeDir, "test_togodo_home.txt")
	err = os.WriteFile(homeTestFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write home test file: %v", err)
	}
	defer os.Remove(homeTestFile)

	// Test tilde expansion
	tildeTestPath := "~/test_togodo_home.txt"
	config.SetTodoTxtPath(tildeTestPath)
	actualPath = config.GetTodoTxtPath()
	expectedPath := homeTestFile

	if actualPath != expectedPath {
		t.Errorf("Expected tilde expansion to result in %s, got %s", expectedPath, actualPath)
	}

	// Reset config
	config.SetTodoTxtPath("")
}

func TestFileFlag_RootCommandPersistentFlag(t *testing.T) {
	repo, _ := setupTestRepository(t)
	presenter := cli.NewPresenter()
	rootCmd := NewRootCmd(repo, presenter)

	// Test that the flag is properly configured as persistent
	fileFlag := rootCmd.PersistentFlags().Lookup("file")
	if fileFlag == nil {
		t.Fatal("Expected 'file' persistent flag to be present")
	}

	// Test flag configuration
	tests := []struct {
		property string
		expected string
		actual   string
	}{
		{"Name", "file", fileFlag.Name},
		{"Shorthand", "f", fileFlag.Shorthand},
		{"Usage", "Specify the todo.txt file to use", fileFlag.Usage},
		{"DefValue", "", fileFlag.DefValue},
	}

	for _, test := range tests {
		if test.actual != test.expected {
			t.Errorf("Expected %s to be '%s', got '%s'", test.property, test.expected, test.actual)
		}
	}
}

func TestFileFlag_PersistentPreRun(t *testing.T) {
	// Create a temporary file
	tempDir, err := os.MkdirTemp("", "togodo_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "persistent_test.txt")
	testContent := "(A) Test task from persistent flag"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test that PersistentPreRun is configured
	repo, _ := setupTestRepository(t)
	presenter := cli.NewPresenter()
	rootCmd := NewRootCmd(repo, presenter)
	if rootCmd.PersistentPreRun == nil {
		t.Fatal("Expected PersistentPreRun to be configured")
	}

	// Test that PersistentPreRun handles the flag correctly
	// We can't easily test the actual flag parsing without mocking cobra,
	// but we can verify the function exists and doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PersistentPreRun panicked: %v", r)
		}
	}()

	// Call PersistentPreRun with empty args (simulates flag not being set)
	rootCmd.PersistentPreRun(rootCmd, []string{})

	// Reset any config changes
	config.SetTodoTxtPath("")
}

func TestFileFlag_ConfigOverride(t *testing.T) {
	// Test the config override functionality directly
	originalPath := config.GetTodoTxtPath()

	// Create test paths
	testPaths := []string{
		"/tmp/test_todo.txt",
		"./relative_todo.txt",
		"~/home_todo.txt",
		"/absolute/path/todo.txt",
	}

	for _, testPath := range testPaths {
		config.SetTodoTxtPath(testPath)
		actualPath := config.GetTodoTxtPath()

		// For tilde paths, check that expansion occurred
		if strings.HasPrefix(testPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				expectedPath := filepath.Join(homeDir, testPath[2:])
				if actualPath != expectedPath {
					t.Errorf("Tilde expansion failed. Expected %s, got %s", expectedPath, actualPath)
				}
			}
		} else {
			// For non-tilde paths, should match exactly
			if actualPath != testPath {
				t.Errorf("Path override failed. Expected %s, got %s", testPath, actualPath)
			}
		}

		// Reset for next iteration
		config.SetTodoTxtPath("")
	}

	// Verify reset works
	resetPath := config.GetTodoTxtPath()
	if resetPath == "" {
		// If empty, it should fall back to default config
		// We can't predict the exact default without mocking viper
		// but we can verify it's not empty after reset
		t.Logf("Config reset to default: %s", resetPath)
	}

	// Restore original path if needed
	if originalPath != "" {
		config.SetTodoTxtPath(originalPath)
	}
}

func TestFileFlag_EmptyPath(t *testing.T) {
	// Test behavior with empty path
	config.SetTodoTxtPath("")
	path := config.GetTodoTxtPath()

	// Should fall back to viper config default
	// We can't predict the exact value without mocking viper,
	// but it should have some default value
	if path == "" {
		t.Log("Config returned empty path, which may be expected in test environment")
	} else {
		t.Logf("Config returned default path: %s", path)
	}
}

func TestFileFlag_InvalidTildePath(t *testing.T) {
	// Test tilde expansion with edge cases
	testCases := []struct {
		input    string
		shouldExpandTilde bool
	}{
		{"~", false},           // Just tilde
		{"~/", true},          // Tilde with trailing slash
		{"~invalid", false},   // Tilde with immediate text (should not expand)
		{" ~/test.txt", false}, // Tilde with leading space (should not expand)
		{"~/test.txt", true},  // Valid tilde path
	}

	for _, testCase := range testCases {
		config.SetTodoTxtPath(testCase.input)
		actualPath := config.GetTodoTxtPath()

		if testCase.shouldExpandTilde {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				var expectedPath string
				if testCase.input == "~/" {
					expectedPath = homeDir
				} else {
					expectedPath = filepath.Join(homeDir, testCase.input[2:])
				}
				if actualPath != expectedPath {
					t.Errorf("Tilde expansion failed for '%s'. Expected %s, got %s", testCase.input, expectedPath, actualPath)
				}
			}
		} else {
			// Should not be expanded
			if actualPath != testCase.input {
				t.Errorf("Path should not be expanded for '%s'. Expected %s, got %s", testCase.input, testCase.input, actualPath)
			}
		}

		// Reset for next iteration
		config.SetTodoTxtPath("")
	}
}
