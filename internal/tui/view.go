package tui

import (
	"fmt"
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

func (m model) View() string {
	// First build the main view
	var mainView string
	if m.filtering {
		mainView += fmt.Sprintf("\nFilter: %s", m.filter)
	}
	mainView += "\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		mainView += fmt.Sprintf("%s ", cursor)
		mainView += formatTodo(choice) + "\n"
	}

	mainView += "\nx: toggle | p: set priority | /: filter | a: add | q: quit\n"

	// If we're setting priority, show the priority overlay
	if m.setting {
		width := 40
		height := 3

		popup := stylePrimaryBold.Render("Set Priority") + "\n"
		popup += "Press A-D to set priority, ESC to cancel" + "\n"
		popup += styleHelp.Render("(A is highest, D is lowest)")

		overlay := stylePrimary.
			Width(width).
			Height(height).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			Render(popup)

		return lipgloss.Place(
			lipgloss.Width(mainView),
			lipgloss.Height(mainView),
			lipgloss.Center,
			lipgloss.Center,
			overlay,
		)
	}

	// If we're adding, overlay the popup on top
	if m.adding {
		mainWidth := lipgloss.Width(mainView)
		width := int(float64(mainWidth) * 0.75)
		height := 3

		popup := stylePrimaryBold.Render("Add New Todo") + "\n"

		// Create a Todo with proper priority parsing
		text := m.input.Value()
		priority := ""
		if len(text) >= 3 && text[0] == '(' && text[2] == ')' {
			priority = string(text[1])
		}
		popup += formatTodo(todotxtlib.Todo{
			Text:     text,
			Priority: priority,
		}) + "\n"
		popup += styleHelp.Render("(esc to cancel, enter to save)")

		overlay := stylePrimary.
			Width(width).
			Height(height).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			Render(popup)

		return lipgloss.Place(
			lipgloss.Width(mainView),
			lipgloss.Height(mainView),
			lipgloss.Center,
			lipgloss.Center,
			overlay,
		)
	}

	return mainView
}

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
