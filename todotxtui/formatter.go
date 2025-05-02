package todotxtui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gkarolyi/togodo/todotxtlib"
)

// TodoFormatter defines the interface for formatting todos
type TodoFormatter interface {
	Format(todo todotxtlib.Todo) string
	FormatList(todos []todotxtlib.Todo) []string
}

// LipglossFormatter implements TodoFormatter using lipgloss for styling
type LipglossFormatter struct {
	projectStyle    lipgloss.Style
	contextStyle    lipgloss.Style
	doneStyle       lipgloss.Style
	priorityStyle   map[string]lipgloss.Style
	tagStyle        lipgloss.Style
	lineNumberStyle lipgloss.Style
}

// NewLipglossFormatter creates a new LipglossFormatter with default styles
func NewLipglossFormatter() *LipglossFormatter {
	return &LipglossFormatter{
		projectStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")),
		contextStyle: lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#04B575")),
		doneStyle: lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.Color("#7F98AF")),
		priorityStyle: map[string]lipgloss.Style{
			"A": lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#D40B23")),
			"B": lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF6700")),
			"C": lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#0FFF95")),
		},
		tagStyle: lipgloss.NewStyle().
			Underline(true).
			Foreground(lipgloss.Color("#96C5B0")),
		lineNumberStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7F98AF")),
	}
}

// Format implements TodoFormatter for LipglossFormatter
func (f *LipglossFormatter) Format(todo todotxtlib.Todo) string {
	var builder strings.Builder
	words := strings.Fields(todo.Text)
	stdStyle := f.priorityStyle[todo.Priority]

	if todo.Done {
		builder.WriteString(f.doneStyle.Render(todo.Text))
	} else {
		for i, word := range words {
			if isProject(word) {
				builder.WriteString(f.projectStyle.Render(word))
			} else if isContext(word) {
				builder.WriteString(f.contextStyle.Render(word))
			} else if isTag(word) {
				builder.WriteString(f.tagStyle.Render(word))
			} else {
				builder.WriteString(stdStyle.Render(word))
			}
			if i < len(words)-1 {
				builder.WriteString(" ")
			}
		}
	}

	return builder.String()
}

// FormatList implements TodoFormatter for LipglossFormatter
func (f *LipglossFormatter) FormatList(todos []todotxtlib.Todo) []string {
	formatted := make([]string, len(todos))
	for i, todo := range todos {
		// Add line number (index + 1) before the formatted todo
		lineNumber := fmt.Sprintf("%3d ", i+1)
		formatted[i] = f.lineNumberStyle.Render(lineNumber) + f.Format(todo)
	}
	return formatted
}

// PlainFormatter implements TodoFormatter for simple text output
type PlainFormatter struct{}

// NewPlainFormatter creates a new PlainFormatter
func NewPlainFormatter() *PlainFormatter {
	return &PlainFormatter{}
}

// Format implements TodoFormatter for PlainFormatter
func (f *PlainFormatter) Format(todo todotxtlib.Todo) string {
	return todo.Text
}

// FormatList implements TodoFormatter for PlainFormatter
func (f *PlainFormatter) FormatList(todos []todotxtlib.Todo) []string {
	formatted := make([]string, len(todos))
	for i, todo := range todos {
		// Add line number (index + 1) before the formatted todo
		formatted[i] = fmt.Sprintf("%3d %s", i+1, f.Format(todo))
	}
	return formatted
}

// Helper functions
func isProject(word string) bool {
	return strings.HasPrefix(word, "+")
}

func isContext(word string) bool {
	return strings.HasPrefix(word, "@")
}

func isTag(word string) bool {
	return strings.Contains(word, ":")
}
