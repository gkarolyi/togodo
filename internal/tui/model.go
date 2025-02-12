package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	todolib "github.com/gkarolyi/togodo/internal/todolib"
)

func NewProgram(repository todolib.TodoRepository) *tea.Program {
	model := initialModel(repository)
	return tea.NewProgram(model)
}

type model struct {
	choices    []todolib.Todo   // items on the to-do list
	cursor     int              // which to-do list item our cursor is pointing at
	selected   map[int]struct{} // which to-do items are selected
	repository todolib.TodoRepository
	filtering  bool   // whether we're currently filtering
	filter     string // the current filter string
}

func initialModel(repository todolib.TodoRepository) model {
	return model{
		// Our to-do list is a grocery list
		repository: repository,
		choices:    repository.Items(),

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected:  make(map[int]struct{}),
		filtering: false,
		filter:    "",
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		// The "x" key toggles the done status for all selected items.
		case "x":
			var lineNumbers []int
			if len(m.selected) == 0 {
				lineNumbers = []int{m.cursor + 1}
			} else {
				for i := range m.selected {
					lineNumbers = append(lineNumbers, i+1)
				}
			}
			m.repository.Toggle(lineNumbers)

		case "/":
			if !m.filtering {
				m.filtering = true
				m.filter = ""
				return m, nil
			}

		case "esc":
			if m.filtering {
				m.filtering = false
				m.filter = ""
				m.choices = m.repository.Items()
				return m, nil
			}

		default:
			if m.filtering {
				switch msg.Type {
				case tea.KeyRunes:
					m.filter += msg.String()
					m.choices = m.repository.Filter(m.filter)
					if m.cursor >= len(m.choices) {
						m.cursor = len(m.choices) - 1
					}
					if m.cursor < 0 {
						m.cursor = 0
					}
					return m, nil
				case tea.KeyBackspace:
					if len(m.filter) > 0 {
						m.filter = m.filter[:len(m.filter)-1]
						m.choices = m.repository.Filter(m.filter)
						if m.cursor >= len(m.choices) {
							m.cursor = len(m.choices) - 1
						}
						if m.cursor < 0 {
							m.cursor = 0
						}
					}
					return m, nil
				case tea.KeyEnter:
					m.filtering = false
					return m, nil
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := m.repository.Path()
	if m.filtering {
		s += fmt.Sprintf("\nFilter: %s", m.filter)
	}
	s += "\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s ", cursor)
		s += todolib.RenderToString(choice)
		s += "\n"
	}

	// The footer
	s += "\nx: toggle | p: set priority | /: filter | q: quit\n"

	// Send the UI for rendering
	return s
}
