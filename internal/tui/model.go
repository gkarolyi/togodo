package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gkarolyi/togodo/todotxtlib"
)

var (
	stylePrimary     = lipgloss.NewStyle().Padding(1, 2)
	stylePrimaryBold = stylePrimary.Bold(true)
	styleHelp        = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Italic(true)
)

type model struct {
	choices    []todotxtlib.Todo // items on the to-do list
	cursor     int               // which to-do list item our cursor is pointing at
	selected   map[int]struct{}  // which to-do items are selected
	repository *todotxtlib.Repository
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
