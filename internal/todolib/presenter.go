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
var priorityStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#fc1c03"))

func Render(todo Todo) {
	fmt.Printf("%d ", todo.Number)
	words := strings.Fields(todo.Text)

	if todo.Done {
		renderStyle(todo.Text, doneStyle)
	} else if todo.Prioritised() {
		renderStyle(todo.Text, priorityStyle)
	} else {
		renderWords(words)
	}

	fmt.Println()
}

func renderWords(words []string) {
	for _, word := range words {
		if IsProject(word) {
			renderStyle(word, projectStyle)
		} else if IsContext(word) {
			renderStyle(word, contextStyle)
		} else {
			fmt.Print(word)
		}
		fmt.Print(" ")
	}
}

func renderStyle(word string, style lipgloss.Style) {
	fmt.Print(style.Render(word))
}
