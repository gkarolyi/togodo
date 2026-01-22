package tests

import (
	"testing"
)

// TestConfigFilePriority tests TODOTXT_PRIORITY_ON_ADD config variable
// Ported from: t1040-add-priority.sh "config file priority"
func TestConfigFilePriority(t *testing.T) {
	env := SetupTestEnv(t)

	// Note: Upstream sets: export TODOTXT_PRIORITY_ON_ADD=A in todo.cfg
	// This feature auto-adds priority (A) to every new task
	// If togodo doesn't support TODOTXT_PRIORITY_ON_ADD, this test will fail
	// TODO: Implement TODOTXT_PRIORITY_ON_ADD config variable support

	t.Run("set TODOTXT_PRIORITY_ON_ADD=A", func(t *testing.T) {
		// TODO: Implement setting TODOTXT_PRIORITY_ON_ADD
		// This may need environment variable or config file support
		t.Skip("TODO: Implement TODOTXT_PRIORITY_ON_ADD config variable")
	})

	t.Run("add task with priority auto-added", func(t *testing.T) {
		output, code := env.RunCommand("add", "take out the trash")
		// Upstream expects: "1 (A) take out the trash\nTODO: 1 added."
		expectedOutput := "1 (A) take out the trash\nTODO: 1 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Expected output:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list with -p flag shows priority", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		// Upstream expects: "1 (A) take out the trash\n--\nTODO: 1 of 1 tasks shown"
		expectedOutput := "1 (A) take out the trash\n--\nTODO: 1 of 1 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Expected output:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestConfigFileWrongPriority tests error handling for invalid TODOTXT_PRIORITY_ON_ADD values
// Ported from: t1040-add-priority.sh "config file wrong priority"
func TestConfigFileWrongPriority(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("set TODOTXT_PRIORITY_ON_ADD to invalid value", func(t *testing.T) {
		// TODO: Set TODOTXT_PRIORITY_ON_ADD=1 (invalid - must be A-Z)
		t.Skip("TODO: Implement TODOTXT_PRIORITY_ON_ADD config variable")
	})

	t.Run("add task should fail with error", func(t *testing.T) {
		output, code := env.RunCommand("add", "fail to take out the trash")
		// Upstream expects exit code 1 and error message
		expectedCode := 1
		expectedError := `TODOTXT_PRIORITY_ON_ADD should be a capital letter from A to Z (it is now "1").`

		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedError {
			t.Errorf("Expected error:\n%s\n\nGot:\n%s", expectedError, output)
		}
	})

	t.Run("list should also fail with error", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		// Upstream expects same error on list command too
		expectedCode := 1
		expectedError := `TODOTXT_PRIORITY_ON_ADD should be a capital letter from A to Z (it is now "1").`

		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedError {
			t.Errorf("Expected error:\n%s\n\nGot:\n%s", expectedError, output)
		}
	})
}
