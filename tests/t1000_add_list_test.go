package tests

import (
	"testing"
)

// TestBasicAddList tests basic add and list functionality
// Ported from: t1000-addlist.sh
func TestBasicAddList(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("basic add/list 1", func(t *testing.T) {
		output, code := env.RunCommand("add", "notice the daisies")
		expectedOutput := "1 notice the daisies\nTODO: 1 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("basic add/list 2", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := "1 notice the daisies\n--\nTODO: 1 of 1 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("basic add/list 3", func(t *testing.T) {
		output, code := env.RunCommand("add", "smell the roses")
		expectedOutput := "2 smell the roses\nTODO: 2 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("basic add/list 4", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := "1 notice the daisies\n2 smell the roses\n--\nTODO: 2 of 2 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListFiltering tests basic list filtering
// Ported from: t1000-addlist.sh
func TestListFiltering(t *testing.T) {
	env := SetupTestEnv(t)

	// Setup: add tasks
	env.RunCommand("add", "notice the daisies")
	env.RunCommand("add", "smell the roses")

	t.Run("filter by daisies", func(t *testing.T) {
		output, code := env.RunCommand("list", "daisies")
		expectedOutput := "1 notice the daisies\n--\nTODO: 1 of 2 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("filter by smell", func(t *testing.T) {
		output, code := env.RunCommand("list", "smell")
		expectedOutput := "2 smell the roses\n--\nTODO: 1 of 2 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestCaseInsensitiveFiltering tests case-insensitive filtering
// Ported from: t1000-addlist.sh
func TestCaseInsensitiveFiltering(t *testing.T) {
	env := SetupTestEnv(t)

	// Setup: add tasks
	env.RunCommand("add", "notice the daisies")
	env.RunCommand("add", "smell the roses")

	t.Run("add uppercase task", func(t *testing.T) {
		output, code := env.RunCommand("add", "smell the uppercase Roses")
		expectedOutput := "3 smell the uppercase Roses\nTODO: 3 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("filter lowercase roses", func(t *testing.T) {
		output, code := env.RunCommand("list", "roses")
		expectedOutput := "2 smell the roses\n3 smell the uppercase Roses\n--\nTODO: 2 of 3 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestAddWithSymbols tests adding tasks with special characters
// Ported from: t1000-addlist.sh
func TestAddWithSymbols(t *testing.T) {
	env := SetupTestEnv(t)

	// Setup: add previous tasks
	env.RunCommand("add", "notice the daisies")
	env.RunCommand("add", "smell the roses")
	env.RunCommand("add", "smell the uppercase Roses")

	t.Run("add symbols", func(t *testing.T) {
		output, code := env.RunCommand("add", "~@#$%^&*()-_=+[{]}|;:',<.>/?")
		expectedOutput := "4 ~@#$%^&*()-_=+[{]}|;:',<.>/?\nTODO: 4 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("add backtick and quotes", func(t *testing.T) {
		output, code := env.RunCommand("add", "`!\\\"")
		expectedOutput := "5 `!\\\"\nTODO: 5 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list all with symbols", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := "1 notice the daisies\n2 smell the roses\n3 smell the uppercase Roses\n5 `!\\\"\n4 ~@#$%^&*()-_=+[{]}|;:',<.>/?\n--\nTODO: 5 of 5 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestAddWithSpaces tests adding tasks with multiple spaces
// Ported from: t1000-addlist.sh
func TestAddWithSpaces(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("add with quoted spaces", func(t *testing.T) {
		output, code := env.RunCommand("add", "notice the   three   spaces")
		expectedOutput := "1 notice the   three   spaces\nTODO: 1 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("add with unquoted spaces", func(t *testing.T) {
		// Note: In bash "notice how the   spaces    get lost" becomes separate args
		// In Go test, we pass as single string which should preserve spaces
		// TODO: May need to adjust based on actual CLI behavior
		output, code := env.RunCommand("add", "notice", "how", "the", "spaces", "get", "lost")
		expectedOutput := "2 notice how the spaces get lost\nTODO: 2 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list with spaces", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := "2 notice how the spaces get lost\n1 notice the   three   spaces\n--\nTODO: 2 of 2 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
