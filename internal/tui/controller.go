package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gkarolyi/togodo/todotxtlib"
)

// controller handles TUI initialization and execution
type controller struct {
	repo *todotxtlib.Repository
}

// NewController creates a new TUI controller
func NewController(repo *todotxtlib.Repository) *controller {
	return &controller{
		repo: repo,
	}
}

// Run starts the TUI interface
func (c *controller) Run() error {
	model := initialModel(c.repo)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
