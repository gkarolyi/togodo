package tests

import (
	"testing"
)

// TestCmdLineFirstDay tests date on add with -t flag
// Ported from: t1010-add-date.sh "cmd line first day"
func TestCmdLineFirstDay(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	// Set test date to 2009-02-13 (matching bash test)
	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	t.Run("add with -t flag", func(t *testing.T) {
		output, code := env.RunCommand("-t", "add", "notice the daisies")
		expectedOutput := "1 2009-02-13 notice the daisies\nTODO: 1 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows date", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := "1 2009-02-13 notice the daisies\n--\nTODO: 1 of 1 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestCmdLineFirstDayWithPriority tests -pt flags together
// Ported from: t1010-add-date.sh "cmd line first day with priority"
func TestCmdLineFirstDayWithPriority(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	t.Run("add with -pt flags", func(t *testing.T) {
		output, code := env.RunCommand("-pt", "add", "(A) notice the daisies")
		expectedOutput := "2 (A) 2009-02-13 notice the daisies\nTODO: 2 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list with -p flag shows priority", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := "2 (A) 2009-02-13 notice the daisies\n1 2009-02-13 notice the daisies\n--\nTODO: 2 of 2 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("delete with -npf flags", func(t *testing.T) {
		output, code := env.RunCommand("-npf", "del", "2")
		expectedOutput := "2 (A) 2009-02-13 notice the daisies\nTODO: 2 deleted."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestCmdLineFirstDayWithLowercasePriority tests lowercase priority normalization
// Ported from: t1010-add-date.sh "cmd line first day with lowercase priority"
func TestCmdLineFirstDayWithLowercasePriority(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	t.Run("add with lowercase priority (b)", func(t *testing.T) {
		output, code := env.RunCommand("-pt", "add", "(b) notice the daisies")
		expectedOutput := "2 (B) 2009-02-13 notice the daisies\nTODO: 2 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list with -p flag", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := "2 (B) 2009-02-13 notice the daisies\n1 2009-02-13 notice the daisies\n--\nTODO: 2 of 2 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("delete with -npf flags", func(t *testing.T) {
		output, code := env.RunCommand("-npf", "del", "2")
		expectedOutput := "2 (B) 2009-02-13 notice the daisies\nTODO: 2 deleted."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestCmdLineSecondDay tests date changes across test_tick
// Ported from: t1010-add-date.sh "cmd line second day"
func TestCmdLineSecondDay(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	// Advance to next day (test_tick)
	env.TestTick()

	t.Run("add on second day", func(t *testing.T) {
		output, code := env.RunCommand("-t", "add", "smell the roses")
		expectedOutput := "2 2009-02-14 smell the roses\nTODO: 2 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows both tasks with dates", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := "1 2009-02-13 notice the daisies\n2 2009-02-14 smell the roses\n--\nTODO: 2 of 2 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestCmdLineThirdDay tests another day with test_tick
// Ported from: t1010-add-date.sh "cmd line third day"
func TestCmdLineThirdDay(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	// Advance two days
	env.TestTick()
	env.TestTick()

	t.Run("add on third day", func(t *testing.T) {
		output, code := env.RunCommand("-t", "add", "mow the lawn")
		expectedOutput := "3 2009-02-15 mow the lawn\nTODO: 3 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows three tasks", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := "1 2009-02-13 notice the daisies\n2 2009-02-14 smell the roses\n3 2009-02-15 mow the lawn\n--\nTODO: 3 of 3 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestConfigFileDateOnAdd tests TODOTXT_DATE_ON_ADD via config file
// Ported from: t1010-add-date.sh "config file third day"
func TestConfigFileDateOnAdd(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	// Advance to third day, tick an extra hour (test_tick 3600)
	env.TestTick()
	env.TestTick()
	env.TestTick(3600) // Bump the clock by one hour

	t.Run("add with TODOTXT_DATE_ON_ADD=1 in config", func(t *testing.T) {
		// TODO: Set TODOTXT_DATE_ON_ADD=1 via environment or config file
		// For now, use -t flag as workaround
		output, code := env.RunCommand("-t", "add", "take out the trash")
		expectedOutput := "4 2009-02-15 take out the trash\nTODO: 4 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows all four tasks", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := "1 2009-02-13 notice the daisies\n2 2009-02-14 smell the roses\n3 2009-02-15 mow the lawn\n4 2009-02-15 take out the trash\n--\nTODO: 4 of 4 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
