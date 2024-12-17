package todolib

import (
	"strings"
	"testing"
)

func TestTodo(t *testing.T) {
	todo := Todo{
		Text:     "(B) random fake task with a +projectName and @contextName",
		Done:     false,
		Priority: "B",
		Projects: []string{"projectName"},
		Contexts: []string{"contextName"},
	}

	t.Run("Done", func(t *testing.T) {
		if todo.Done != false {
			t.Errorf("expected Done to be false, got %v", todo.Done)
		}
	})

	t.Run("project", func(t *testing.T) {
		if !equalStringSlices(todo.Projects, []string{"projectName"}) {
			t.Errorf("expected Projects to be [projectName], got %v", todo.Projects)
		}
	})

	t.Run("context", func(t *testing.T) {
		if !equalStringSlices(todo.Contexts, []string{"contextName"}) {
			t.Errorf("expected Contexts to be [contextName], got %v", todo.Contexts)
		}
	})

	t.Run("priority", func(t *testing.T) {
		if todo.Priority != "B" {
			t.Errorf("expected Priority to be B, got %v", todo.Priority)
		}
	})

	t.Run("text", func(t *testing.T) {
		if todo.Text != "(B) random fake task with a +projectName and @contextName" {
			t.Errorf("expected Text to be '(B) random fake task with a +projectName and @contextName', got %v", todo.Text)
		}
	})

	t.Run("done todo", func(t *testing.T) {
		todo := Todo{
			Text: "x this todo is done",
			Done: true,
		}
		if todo.Done != true {
			t.Errorf("expected Done to be true, got %v", todo.Done)
		}
	})

	t.Run("line number", func(t *testing.T) {
		todo := Todo{Text: "here's another todo item", Number: 3}
		if todo.Number != 3 {
			t.Errorf("expected Number to be 3, got %v", todo.Number)
		}
	})
}

func TestPrioritised(t *testing.T) {
	t.Run("when todo has priority", func(t *testing.T) {
		todo := Todo{Priority: "(A)"}
		if !todo.Prioritised() {
			t.Errorf("expected todo to be prioritised")
		}
	})

	t.Run("when todo does not have priority", func(t *testing.T) {
		todo := Todo{}
		if todo.Prioritised() {
			t.Errorf("expected todo to not be prioritised")
		}
	})
}

func TestToggleDone(t *testing.T) {
	t.Run("when todo is not done", func(t *testing.T) {
		todo := Todo{Text: "this todo is not done"}
		todo.ToggleDone()

		t.Run("state change to done", func(t *testing.T) {
			if todo.Done != true {
				t.Errorf("expected Done to be true, got %v", todo.Done)
			}
		})

		t.Run("todo text is prepended with x", func(t *testing.T) {
			if !strings.HasPrefix(todo.Text, "x ") {
				t.Errorf("expected Text to start with 'x ', got %v", todo.Text)
			}
		})
	})

	t.Run("when todo is done", func(t *testing.T) {
		todo := Todo{Text: "x this todo is done", Done: true}
		todo.ToggleDone()

		t.Run("state change to not done", func(t *testing.T) {
			if todo.Done != false {
				t.Errorf("expected Done to be false, got %v", todo.Done)
			}
		})

		t.Run("x at beginning of line is removed", func(t *testing.T) {
			if strings.HasPrefix(todo.Text, "x ") {
				t.Errorf("expected Text to not start with 'x ', got %v", todo.Text)
			}
		})
	})
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
