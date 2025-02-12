package todolib

import "strings"

type Todo struct {
	Text     string
	Done     bool
	Priority string
	Projects []string
	Contexts []string
	Number   int
}

func NewTodo(text string, index int) Todo {
	todo := Todo{
		Text:     text,
		Done:     IsDone(text),
		Priority: FindPriority(text),
		Projects: FindProjects(text),
		Contexts: FindContexts(text),
		Number:   index + 1,
	}

	return todo
}

func (t Todo) Prioritised() bool {
	return t.Priority != ""
}

func (t *Todo) ToggleDone() {
	if t.Done {
		t.Done = false
		t.Text = strings.TrimPrefix(t.Text, "x ")
	} else {
		t.Done = true
		t.Text = strings.Join([]string{"x ", t.Text}, "")
	}
}

// SetPriority sets the priority of the todo item.
func (t *Todo) SetPriority(priority string) {
	t.Text = strings.TrimPrefix(t.Text, "("+t.Priority+") ")
	t.Priority = priority

	if t.Prioritised() {
		t.Text = strings.Join([]string{"(", t.Priority, ") ", t.Text}, "")
	}
}

func (t Todo) Equals(other Todo) bool {
	return t.Text == other.Text
}
