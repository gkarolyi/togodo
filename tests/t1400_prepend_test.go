package tests

import (
	"testing"
)

// TestPrependUsage tests prepend command usage
// Ported from: t1400-prepend.sh
func TestPrependUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("prepend with wrong args", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "adf", "asdfa")
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Check error message "usage: togodo prepend NR \"TEXT TO PREPEND\""
		_ = output
	})
}

// TestBasicPrepend tests basic prepend functionality
// Ported from: t1400-prepend.sh
func TestBasicPrepend(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
notice the sunflowers
stop`)

	t.Run("prepend text to task", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "2", "really")
		expectedOutput := "2 really notice the sunflowers"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify prepend in file", func(t *testing.T) {
		content := env.ReadTodoFile()
		if !contains(content, "really notice the sunflowers") {
			t.Errorf("File should contain prepended text, got:\n%s", content)
		}
	})
}

// TestPrependPreservesPriority tests that prepend preserves priority
// Ported from: t1400-prepend.sh
func TestPrependPreservesPriority(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile("(A) task with priority")

	t.Run("prepend to prioritized task", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "1", "important")
		// Should preserve (A) priority
		expectedOutput := "1 (A) important task with priority"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || len(s) > len(substr)+1 && s[1:len(substr)+1] == substr || indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
