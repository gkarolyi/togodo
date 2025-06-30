package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/gkarolyi/togodo/todotxtui/format"
)

var (
	stylePrimary     = lipgloss.NewStyle().Padding(1, 2)
	stylePrimaryBold = stylePrimary.Bold(true)
	styleHelp        = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Italic(true)
)

func NewProgram(repository *todotxtlib.Repository) *tea.Program {
	model := initialModel(repository)
	return tea.NewProgram(model)
}

type model struct {
	choices    []todotxtlib.Todo // items on the to-do list
	cursor     int               // which to-do list item our cursor is pointing at
	selected   map[int]struct{}  // which to-do items are selected
	repository *todotxtlib.Repository
	formatter  format.TodoFormatter
	filtering  bool            // whether we're currently filtering
	filter     string          // the current filter string
	adding     bool            // whether we're currently adding a new item
	input      textinput.Model // text input for new items
	setting    bool            // whether we're currently setting priority
}

func initialModel(repository *todotxtlib.Repository) model {
	ti := textinput.New()
	ti.Placeholder = "Enter new todo item..."
	ti.CharLimit = 150
	ti.Width = 50

	allTodos, err := repository.ListAll()
	if err != nil {
		return model{
			repository: repository,
			formatter:  format.NewLipglossFormatter(),
			choices:    []todotxtlib.Todo{},
			selected:   make(map[int]struct{}),
			filtering:  false,
			filter:     "",
			adding:     false,
			setting:    false,
			input:      ti,
		}
	} else {
		return model{
			repository: repository,
			formatter:  format.NewLipglossFormatter(),
			choices:    allTodos,
			selected:   make(map[int]struct{}),
			filtering:  false,
			filter:     "",
			adding:     false,
			setting:    false,
			input:      ti,
		}
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If we're setting priority, handle priority keys
		if m.setting {
			switch msg.String() {
			case "esc":
				m.setting = false
				return m, nil
			case "a", "b", "c", "d", "A", "B", "C", "D":
				priority := strings.ToUpper(msg.String())
				for i := range m.selected {
					m.repository.SetPriority(i, priority)
				}
				allTodos, _ := m.repository.ListAll()
				m.choices = allTodos
				m.setting = false
				return m, nil
			}
			return m, nil
		}

		// If we're adding, handle all keys through the text input except ESC and Enter
		if m.adding {
			switch msg.Type {
			case tea.KeyEsc:
				m.adding = false
				m.input.Reset()
				m.input.Blur()
				return m, nil
			case tea.KeyEnter:
				if m.input.Value() != "" {
					m.repository.Add(m.input.Value())
					allTodos, _ := m.repository.ListAll()
					m.choices = allTodos
					m.adding = false
					m.input.Reset()
					m.input.Blur()
				}
				return m, nil
			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
		}

		// If we're filtering, handle all keys except ESC and Enter
		if m.filtering {
			switch msg.Type {
			case tea.KeyEsc:
				m.filtering = false
				m.filter = ""
				allTodos, _ := m.repository.ListAll()
				m.choices = allTodos
				return m, nil
			case tea.KeyEnter:
				m.filtering = false
				return m, nil
			case tea.KeyBackspace:
				if len(m.filter) > 0 {
					m.filter = m.filter[:len(m.filter)-1]
					filteredTodos, _ := m.repository.Search(m.filter)
					m.choices = filteredTodos
					if m.cursor >= len(m.choices) {
						m.cursor = len(m.choices) - 1
					}
					if m.cursor < 0 {
						m.cursor = 0
					}
				}
				return m, nil
			default:
				m.filter += msg.String()
				filteredTodos, _ := m.repository.Search(m.filter)
				m.choices = filteredTodos
				if m.cursor >= len(m.choices) {
					m.cursor = len(m.choices) - 1
				}
				if m.cursor < 0 {
					m.cursor = 0
				}
				return m, nil
			}
		}

		// Regular key handling when not adding or filtering
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		case "x":
			for i := range m.selected {
				m.repository.ToggleDone(i)
			}
			allTodos, _ := m.repository.ListAll()
			m.choices = allTodos

		case "/":
			if !m.adding {
				m.filtering = true
				m.filter = ""
				return m, nil
			}

		case "a":
			if !m.filtering {
				m.adding = true
				m.input.Reset()
				m.input.Focus()
				return m, textinput.Blink
			}

		case "p":
			if !m.filtering && !m.adding {
				m.setting = true
				return m, nil
			}
		}
	}

	return m, nil
}

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
		mainView += m.formatter.Format(choice) + "\n"
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
		popup += m.formatter.Format(todotxtlib.Todo{
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
