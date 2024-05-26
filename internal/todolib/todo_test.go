package todolib

import "testing"

func TestTodo(t *testing.T) {
	todo := Todo{"(B) random fake task with a +project and @context"}

	t.Run("project", func(t *testing.T) {
		got := todo.Project()
		want := "+project"

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("context", func(t *testing.T) {
		got := todo.Context()
		want := "@context"

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("priority", func(t *testing.T) {
		got := todo.Priority()
		want := "(B)"

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("text", func(t *testing.T) {
		got := todo.Text
		want := "(B) random fake task with a +project and @context"

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
