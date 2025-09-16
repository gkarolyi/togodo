package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gkarolyi/togodo/todotxtlib"
)

// Styling for todo formatting
var (
	projectStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	contextStyle   = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#04B575"))
	doneStyle      = lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("#7F98AF"))
	tagStyle       = lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("#96C5B0"))
	priorityAStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#D40B23"))
	priorityBStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6700"))
	priorityCStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0FFF95"))
)

// formatTodo formats a single todo item for display in the TUI
func formatTodo(todo todotxtlib.Todo) string {
	var builder strings.Builder
	words := strings.Fields(todo.Text)

	// Choose priority style
	var stdStyle lipgloss.Style
	switch todo.Priority {
	case "A":
		stdStyle = priorityAStyle
	case "B":
		stdStyle = priorityBStyle
	case "C":
		stdStyle = priorityCStyle
	default:
		stdStyle = lipgloss.NewStyle() // Default unstyled
	}

	if todo.Done {
		builder.WriteString(doneStyle.Render(todo.Text))
	} else {
		for i, word := range words {
			if strings.HasPrefix(word, "+") {
				builder.WriteString(projectStyle.Render(word))
			} else if strings.HasPrefix(word, "@") {
				builder.WriteString(contextStyle.Render(word))
			} else if strings.Contains(word, ":") {
				builder.WriteString(tagStyle.Render(word))
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
