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

func Render(number int, todo Todo) {
	fmt.Printf("%d ", number+1)
	words := strings.Fields(todo.Text)

	if todo.Done {
		renderStyle(todo.Text, doneStyle)
	} else {
		for _, word := range words {
			if IsProject(word) {
				renderStyle(word, projectStyle)
			} else if IsContext(word) {
				renderStyle(word, contextStyle)
			} else {
				if todo.Prioritised() {
					renderStyle(word, priorityStyle)
				} else {
					fmt.Print(word)
				}
			}
			fmt.Print(" ")
		}
	}

	fmt.Println()
}

func renderStyle(word string, style lipgloss.Style) {
	fmt.Print(style.Render(word))
}
