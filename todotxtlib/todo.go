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

// SetContexts sets the contexts of a todo
func (t *Todo) SetContexts(contexts []string) {
	t.Contexts = contexts
}

// SetProjects sets the projects of a todo
func (t *Todo) SetProjects(projects []string) {
	t.Projects = projects
}

// AddContext adds a context to a todo if it doesn't already exist
func (t *Todo) AddContext(context string) {
	// Check if context already exists
	for _, ctx := range t.Contexts {
		if ctx == context {
			return
		}
	}
	t.Contexts = append(t.Contexts, context)
}

// AddProject adds a project to a todo if it doesn't already exist
func (t *Todo) AddProject(project string) {
	// Check if project already exists
	for _, proj := range t.Projects {
		if proj == project {
			return
		}
	}
	t.Projects = append(t.Projects, project)
}

// RemoveContext removes a context from a todo
func (t *Todo) RemoveContext(context string) {
	// Find and remove the context
	for i, ctx := range t.Contexts {
		if ctx == context {
			t.Contexts = append(t.Contexts[:i], t.Contexts[i+1:]...)
			break
		}
	}
}

// RemoveProject removes a project from a todo
func (t *Todo) RemoveProject(project string) {
	// Find and remove the project
	for i, proj := range t.Projects {
		if proj == project {
			t.Projects = append(t.Projects[:i], t.Projects[i+1:]...)
			break
		}
	}
}
