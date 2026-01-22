package tests

import (
	"testing"
)

// TestListconNoContexts tests listcon with no contexts in file
// Ported from: t1310-listcon.sh "listcon no contexts"
func TestListconNoContexts(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(A) email test@example.com today
(B) no context here`)

	t.Run("listcon with no contexts", func(t *testing.T) {
		output, code := env.RunCommand("listcon")
		expectedOutput := ""
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected empty output\nGot:\n%s", output)
		}
	})
}

// TestListconSingle tests listing contexts from tasks with single contexts
// Ported from: t1310-listcon.sh "Single context per line"
func TestListconSingle(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(A) @1 -- Some context 1 task, whitespace, one char
(A) @c2 -- Some context 2 task, whitespace, two char
@con03 -- Some context 3 task, no whitespace
@con04 -- Some context 4 task, no whitespace
@con05@con06 -- weird context`)

	t.Run("listcon shows all contexts", func(t *testing.T) {
		output, code := env.RunCommand("listcon")
		expectedOutput := `@1
@c2
@con03
@con04
@con05@con06`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconMultiple tests listing contexts when tasks have multiple contexts
// Ported from: t1310-listcon.sh "Multi-context per line"
func TestListconMultiple(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`@con01 -- Some context 1 task
@con02 -- Some context 2 task
@con02 @con03 -- Multi-context task`)

	t.Run("listcon with duplicates", func(t *testing.T) {
		output, code := env.RunCommand("listcon")
		// Should list unique contexts, sorted
		// @con02 appears twice but should only be listed once
		expectedOutput := `@con01
@con02
@con03`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconEmailAddress tests that email addresses don't trigger context detection
// Ported from: t1310-listcon.sh "listcon e-mail address test"
func TestListconEmailAddress(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`@con01 -- Some context 1 task
@con02 -- Some context 2 task
email test@example.com today`)

	t.Run("listcon ignores email addresses", func(t *testing.T) {
		output, code := env.RunCommand("listcon")
		// Email address "test@example.com" should NOT be detected as context
		expectedOutput := `@con01
@con02`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconWithProject tests listing contexts filtered by project
// Ported from: t1310-listcon.sh "listcon with project"
func TestListconWithProject(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`x (A) @outdoors shopping with @alice +project02
work on business plan @garage +project01
water the plants @garden +landscape`)

	t.Run("listcon filtered by project", func(t *testing.T) {
		output, code := env.RunCommand("listcon", "+landscape")
		expectedOutput := "@garden"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconDefault tests default context detection patterns
// Ported from: t1310-listcon.sh "listcon with default configuration"
func TestListconDefault(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`Review todo.txt-cli's code with @GinaTrapani
call Mom @home) and arrange schedule
do something @x
do something else @y`)

	t.Run("listcon with default config", func(t *testing.T) {
		output, code := env.RunCommand("listcon")
		expectedOutput := `@GinaTrapani
@home)
@x
@y`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconLimitingMultiChar tests limiting contexts to multi-character sequences
// Ported from: t1310-listcon.sh "listcon limiting to multi-character sequences"
func TestListconLimitingMultiChar(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SIGIL_VALID_PATTERN configuration")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`Review todo.txt-cli's code with @GinaTrapani
call Mom @home) and arrange schedule
do something @x
do something else @y`)

	t.Run("listcon with TODOTXT_SIGIL_VALID_PATTERN", func(t *testing.T) {
		// TODO: Set TODOTXT_SIGIL_VALID_PATTERN='.\{2,\}' environment variable
		output, code := env.RunCommand("listcon")
		expectedOutput := `@GinaTrapani
@home)`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconAllowingMarkerBefore tests allowing w: prefix before contexts
// Ported from: t1310-listcon.sh "listcon allowing w: marker before contexts"
func TestListconAllowingMarkerBefore(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SIGIL_BEFORE_PATTERN configuration")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`Review todo.txt-cli's code w:@GinaTrapani and w:@OtherContributors
call Mom @home) and arrange schedule
do something @x
do something else @y`)

	t.Run("listcon with TODOTXT_SIGIL_BEFORE_PATTERN", func(t *testing.T) {
		// TODO: Set TODOTXT_SIGIL_BEFORE_PATTERN='\(w:\)\{0,1\}' environment variable
		output, code := env.RunCommand("listcon")
		expectedOutput := `@GinaTrapani
@OtherContributors
@home)
@x
@y`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconAllowingParentheses tests allowing parentheses around contexts
// Ported from: t1310-listcon.sh "listcon allowing parentheses around contexts"
func TestListconAllowingParentheses(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SIGIL_BEFORE_PATTERN and TODOTXT_SIGIL_AFTER_PATTERN configuration")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`Review todo.txt-cli's code with @GinaTrapani
call Mom (@home) and then visit the (@school)
do something @x
do something else @y`)

	t.Run("listcon with parentheses patterns", func(t *testing.T) {
		// TODO: Set TODOTXT_SIGIL_BEFORE_PATTERN='(\{0,1\}' TODOTXT_SIGIL_AFTER_PATTERN=')\{0,1\}'
		output, code := env.RunCommand("listcon")
		expectedOutput := `@GinaTrapani
@home
@school
@x
@y`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconAllCustomizations tests all customizations combined
// Ported from: t1310-listcon.sh "listcon with all customizations combined"
func TestListconAllCustomizations(t *testing.T) {
	t.Skip("TODO: Implement all TODOTXT_SIGIL_* configuration variables")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`Review todo.txt-cli's code w:@GinaTrapani and w:@OtherContributors
call Mom (@home) and then visit the (@school)
do something @x
do something else @y`)

	t.Run("listcon with all custom patterns", func(t *testing.T) {
		// TODO: Set multiple TODOTXT_SIGIL_* environment variables
		output, code := env.RunCommand("listcon")
		expectedOutput := `@GinaTrapani
@OtherContributors
@home
@school`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconFromDoneTasks tests listing contexts from done.txt
// Ported from: t1310-listcon.sh "listcon from done tasks"
func TestListconFromDoneTasks(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SOURCEVAR configuration for done.txt")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`@con01 -- Some context 1 task`)
	// TODO: Write to done.txt file
	// env.WriteDoneFile(`x 2009-02-13 @done01 -- completed task 1
	// x 2009-02-14 @done02 -- completed task 2`)

	t.Run("listcon from done file", func(t *testing.T) {
		// TODO: Set TODOTXT_SOURCEVAR=$DONE_FILE
		output, code := env.RunCommand("listcon")
		expectedOutput := `@done01
@done02`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconCombinedSources tests listing contexts from both todo.txt and done.txt
// Ported from: t1310-listcon.sh "listcon from combined open + done tasks"
func TestListconCombinedSources(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SOURCEVAR configuration for multiple sources")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`@con01 -- Some context 1 task`)
	// TODO: Write to done.txt file
	// env.WriteDoneFile(`x 2009-02-13 @done01 -- completed task 1
	// x 2009-02-14 @done02 -- completed task 2`)

	t.Run("listcon from combined sources", func(t *testing.T) {
		// TODO: Set TODOTXT_SOURCEVAR='("$TODO_FILE" "$DONE_FILE")'
		output, code := env.RunCommand("listcon")
		expectedOutput := `@con01
@done01
@done02`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
