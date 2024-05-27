package todolib

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestTodo(t *testing.T) {
	is := is.New(t)
	todo := Todo{
		Text:     "(B) random fake task with a +projectName and @contextName",
		Done:     false,
		Priority: "B",
		Projects: []string{"projectName"},
		Contexts: []string{"contextName"},
	}

	t.Run("Done", func(t *testing.T) {
		is.Equal(todo.Done, false) // should not be done
	})

	t.Run("project", func(t *testing.T) {
		is.Equal(todo.Projects, []string{"projectName"}) // Projects should be ["projectName"]
	})

	t.Run("context", func(t *testing.T) {
		is.Equal(todo.Contexts, []string{"contextName"}) // Contexts should be ["contextName"]
	})

	t.Run("priority", func(t *testing.T) {
		is.Equal(todo.Priority, "B") // priority should be "B"
	})

	t.Run("text", func(t *testing.T) {
		is.Equal(todo.Text, "(B) random fake task with a +projectName and @contextName") // text should match
	})

	t.Run("done todo", func(t *testing.T) {
		todo := Todo{
			Text: "x this todo is done",
			Done: true,
		}
		is.True(todo.Done == true) // should be done
	})

	t.Run("line number", func(t *testing.T) {
		todo := Todo{Text: "here's another todo item", Number: 3}
		is.Equal(todo.Number, 3) // index should be 3
	})
}

func TestPrioritised(t *testing.T) {
	is := is.New(t)

	t.Run("when todo has priority", func(t *testing.T) {
		todo := Todo{Priority: "(A)"}
		is.True(todo.Prioritised()) // should be prioritised
	})

	t.Run("when todo does not have priority", func(t *testing.T) {
		todo := Todo{}
		is.Equal(todo.Prioritised(), false) // should not be prioritised
	})
}

func TestToggleDone(t *testing.T) {
	is := is.New(t)

	t.Run("when todo is not done", func(t *testing.T) {
		todo := Todo{Text: "this todo is not done"}
		todo.ToggleDone()

		t.Run("state change to done", func(t *testing.T) {
			is.Equal(todo.Done, true) // state should be set to done
		})

		t.Run("todo text is prepended with x", func(t *testing.T) {
			is.True(strings.HasPrefix(todo.Text, "x ")) // line should begin with x
		})
	})

	t.Run("when todo is done", func(t *testing.T) {
		todo := Todo{Text: "x this todo is done", Done: true}
		todo.ToggleDone()

		t.Run("state change to not done", func(t *testing.T) {
			is.Equal(todo.Done, false) // state should be set to not done
		})

		t.Run("x at beginning of line is removed", func(t *testing.T) {
			is.True(!strings.HasPrefix(todo.Text, "x ")) // line should not begin with x
		})
	})
}
