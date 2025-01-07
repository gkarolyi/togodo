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
		Text:   text,
		Number: index + 1,
	}
	todo.reloadProperties()

	return todo
}

func (t Todo) Prioritised() bool {
	return t.Priority != ""
}

func (t *Todo) ToggleDone() {
	if t.Done {
		t.Text = strings.TrimPrefix(t.Text, "x ")
	} else {
		t.Text = strings.Join([]string{"x ", t.Text}, "")
	}
	t.reloadProperties()
}

func (t *Todo) reloadProperties() {
	t.Done = IsDone(t.Text)
	t.Priority = FindPriority(t.Text)
	t.Projects = FindProjects(t.Text)
	t.Contexts = FindContexts(t.Text)
}

func (t Todo) Equals(other Todo) bool {
	return t.Text == other.Text
}

func (t Todo) hash() string {
	return t.Text
}
