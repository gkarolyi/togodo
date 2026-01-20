package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gkarolyi/togodo/business"
	"github.com/gkarolyi/togodo/todotxtlib"
)

// Run starts the TUI interface
func Run(repo todotxtlib.TodoRepository) error {
	model := initialModel(repo)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
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
				// Convert selected map to slice of indices
				indices := make([]int, 0, len(m.selected))
				for i := range m.selected {
					indices = append(indices, i)
				}

				if len(indices) > 0 {
					_, err := business.SetPriority(m.repository, indices, priority)
					if err != nil {
						// TODO: Show error in UI
					}

					// Refresh display
					allTodos, _ := m.repository.ListAll()
					m.choices = allTodos
					m.selected = make(map[int]struct{}) // Clear selection
				}
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
					_, err := business.Add(m.repository, []string{m.input.Value()})
					if err != nil {
						// TODO: Show error in UI
					}
					// Refresh display
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
					filter := todotxtlib.Filter{Text: m.filter}
					filteredTodos, _ := m.repository.Filter(filter)
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
				filter := todotxtlib.Filter{Text: m.filter}
				filteredTodos, _ := m.repository.Filter(filter)
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
			// Convert selected map to slice of indices
			indices := make([]int, 0, len(m.selected))
			for i := range m.selected {
				indices = append(indices, i)
			}

			if len(indices) > 0 {
				_, err := business.Do(m.repository, indices)
				if err != nil {
					// TODO: Show error in UI
				}

				// Refresh display
				allTodos, _ := m.repository.ListAll()
				m.choices = allTodos
				m.selected = make(map[int]struct{}) // Clear selection
			}

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
