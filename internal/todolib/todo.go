package todolib

import "strings"

type Todo struct {
	Text     string
	Done     bool
	Priority string
	Projects []string
	Contexts []string
	Number   int
	Index    int // TODO: is this needed?
}

func NewTodo(text string, index int) Todo {
	todo := Todo{
		Text:   text,
		Done:   IsDone(text),
		Index:  index,
		Number: index + 1,
	}
	if !todo.Done {
		todo.Priority = FindPriority(text)
		todo.Projects = FindProjects(text)
		todo.Contexts = FindContexts(text)
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

func (t Todo) Equals(other Todo) bool {
	return t.Text == other.Text
}
