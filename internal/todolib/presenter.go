package todolib

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var projectStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4"))
var contextStyle = lipgloss.NewStyle().
	Italic(true).
	Foreground(lipgloss.Color("#04B575"))
var doneStyle = lipgloss.NewStyle().
	Strikethrough(true).
	Foreground(lipgloss.Color("#3C3C3C"))
var priorityAStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#D40B23"))
var priorityBStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FF6700"))
var priorityCStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#0FFF95"))
var tagStyle = lipgloss.NewStyle().
	Underline(true).
	Foreground(lipgloss.Color("#96C5B0"))

func Render(number int, todo Todo) {
	fmt.Printf("%d ", number+1)
	words := strings.Fields(todo.Text)
	stdStyle := priorityStyle(todo.Priority)

	if todo.Done {
		renderStyle(todo.Text, doneStyle)
	} else {
		for _, word := range words {
			if IsProject(word) {
				renderStyle(word, projectStyle)
			} else if IsContext(word) {
				renderStyle(word, contextStyle)
			} else if IsTag(word) {
				renderStyle(word, tagStyle)
			} else {
				renderStyle(word, stdStyle)
			}
			fmt.Print(" ")
		}
	}

	fmt.Println()
}

func renderStyle(word string, style lipgloss.Style) {
	fmt.Print(style.Render(word))
}

func priorityStyle(priority string) lipgloss.Style {
	switch priority {
	case "A":
		return priorityAStyle
	case "B":
		return priorityBStyle
	case "C":
		return priorityCStyle
	default:
		return lipgloss.NewStyle()
	}
}
