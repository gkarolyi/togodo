package todotxtlib

import (
	"regexp"
	"strings"
)

var projectRe = regexp.MustCompile(`\+(\w+)`)
var contextRe = regexp.MustCompile(`@\w+`)
var priorityRe = regexp.MustCompile(`^\(([A-Z])\)`)
var doneRe = regexp.MustCompile(`^x `)
var tagRe = regexp.MustCompile(`\w+:\S+`)

type Todo struct {
	Text     string
	Done     bool
	Priority string
	Projects []string
	Contexts []string
}

func NewTodo(text string) Todo {
	todo := Todo{
		Text:     text,
		Done:     isDone(text),
		Priority: findPriority(text),
		Projects: findProjects(text),
		Contexts: findContexts(text),
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
	return t.Text == other.Text && t.Done == other.Done
}

func findProjects(text string) []string {
	return projectRe.FindAllString(text, -1)
}

func findContexts(text string) []string {
	return contextRe.FindAllString(text, -1)
}

func findPriority(text string) string {
	if priorityRe.MatchString(text) {
		return priorityRe.FindStringSubmatch(text)[1]
	}
	return ""
}

func isDone(text string) bool {
	return doneRe.MatchString(text)
}
