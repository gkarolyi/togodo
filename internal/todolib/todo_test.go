package todolib

import (
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
		is.Equal(todo.Done, false)
	})

	t.Run("project", func(t *testing.T) {
		is.Equal(todo.Projects, []string{"projectName"})
	})

	t.Run("context", func(t *testing.T) {
		is.Equal(todo.Contexts, []string{"contextName"})
	})

	t.Run("priority", func(t *testing.T) {
		is.Equal(todo.Priority, "B")
	})

	t.Run("text", func(t *testing.T) {
		is.Equal(todo.Text, "(B) random fake task with a +projectName and @contextName")
	})

	t.Run("done todo", func(t *testing.T) {
		todo := Todo{
			Text: "x this todo is done",
			Done: true,
		}

		is.True(todo.Done == true)
	})

	t.Run("line number", func(t *testing.T) {
		todo := Todo{
			Text:   "here's another todo item",
			Number: 3,
		}

		is.Equal(todo.Number, 3)
	})
}

func TestPrioritised(t *testing.T) {
	is := is.New(t)

	t.Run("when todo has priority", func(t *testing.T) {
		todo := Todo{Priority: "(A)"}
		is.True(todo.Prioritised())
	})

	t.Run("when todo does not have priority", func(t *testing.T) {
		todo := Todo{}
		is.Equal(todo.Prioritised(), false)
	})
}
