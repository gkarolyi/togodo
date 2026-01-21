package todotxtlib

import (
	"regexp"
	"slices"
	"strings"
)

var projectRe = regexp.MustCompile(`\+\S+`)
var contextRe = regexp.MustCompile(`@\S+`)
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
		Done:     parseDone(text),
		Priority: parsePriority(text),
		Projects: parseProjects(text),
		Contexts: parseContexts(text),
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
	// Remove existing priority from the text
	if t.Priority != "" {
		if t.Done {
			// For done tasks, remove priority after "x "
			oldPattern := "x (" + t.Priority + ") "
			t.Text = strings.Replace(t.Text, oldPattern, "x ", 1)
		} else {
			// For regular tasks, remove priority from beginning
			oldPattern := "(" + t.Priority + ") "
			t.Text = strings.TrimPrefix(t.Text, oldPattern)
		}
	}

	t.Priority = priority

	if t.Prioritised() {
		// Insert priority in the correct position
		if t.Done {
			// For done tasks: "x " + "(A) " + rest
			t.Text = strings.Replace(t.Text, "x ", "x ("+t.Priority+") ", 1)
		} else {
			// For regular tasks: "(A) " + text
			t.Text = strings.Join([]string{"(", t.Priority, ") ", t.Text}, "")
		}
	}
}

func (t Todo) Equals(other Todo) bool {
	if t.Text != other.Text || t.Done != other.Done || t.Priority != other.Priority {
		return false
	}

	if !slices.Equal(t.Projects, other.Projects) || !slices.Equal(t.Contexts, other.Contexts) {
		return false
	}

	return true
}

func parseProjects(text string) []string {
	projects := projectRe.FindAllString(text, -1)
	slices.Sort(projects)
	return slices.Compact(projects)
}

func parseContexts(text string) []string {
	contexts := contextRe.FindAllString(text, -1)
	slices.Sort(contexts)
	return slices.Compact(contexts)
}

func parsePriority(text string) string {
	// Check for priority at start of text
	if priorityRe.MatchString(text) {
		return priorityRe.FindStringSubmatch(text)[1]
	}
	// Check for priority after "x " marker for done tasks
	if strings.HasPrefix(text, "x ") {
		remaining := text[2:] // Remove "x "
		if priorityRe.MatchString(remaining) {
			return priorityRe.FindStringSubmatch(remaining)[1]
		}
	}
	return ""
}

func parseDone(text string) bool {
	return doneRe.MatchString(text)
}

// SetContexts sets the contexts of a todo
func (t *Todo) SetContexts(contexts []string) {
	// Remove all existing contexts from text
	for _, ctx := range t.Contexts {
		t.removeFromText(ctx)
	}

	// Set new contexts
	t.Contexts = slices.Clone(contexts)
	slices.Sort(t.Contexts)
	t.Contexts = slices.Compact(t.Contexts)

	// Add new contexts to text
	for _, ctx := range t.Contexts {
		t.addToText(ctx)
	}
}

// SetProjects sets the projects of a todo
func (t *Todo) SetProjects(projects []string) {
	// Remove all existing projects from text
	for _, proj := range t.Projects {
		t.removeFromText(proj)
	}

	// Set new projects
	t.Projects = slices.Clone(projects)
	slices.Sort(t.Projects)
	t.Projects = slices.Compact(t.Projects)

	// Add new projects to text
	for _, proj := range t.Projects {
		t.addToText(proj)
	}
}

// AddContext adds a context to a todo if it doesn't already exist
func (t *Todo) AddContext(context string) {
	if slices.Contains(t.Contexts, context) {
		return
	}

	t.Contexts = append(t.Contexts, context)
	slices.Sort(t.Contexts)
	t.addToText(context)
}

// AddProject adds a project to a todo if it doesn't already exist
func (t *Todo) AddProject(project string) {
	if slices.Contains(t.Projects, project) {
		return
	}

	t.Projects = append(t.Projects, project)
	slices.Sort(t.Projects)
	t.addToText(project)
}

// RemoveContext removes a context from a todo
func (t *Todo) RemoveContext(context string) {
	// Find and remove the context
	for i, ctx := range t.Contexts {
		if ctx == context {
			t.Contexts = append(t.Contexts[:i], t.Contexts[i+1:]...)
			t.removeFromText(context)
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
			t.removeFromText(project)
			break
		}
	}
}

// addToText adds a project or context to the end of the todo text
func (t *Todo) addToText(item string) {
	if !strings.Contains(t.Text, item) {
		t.Text = strings.TrimSpace(t.Text + " " + item)
	}
}

// removeFromText removes a project or context from the todo text
func (t *Todo) removeFromText(item string) {
	t.Text = strings.ReplaceAll(t.Text, item, "")
}
