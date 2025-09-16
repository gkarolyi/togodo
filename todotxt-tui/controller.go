package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gkarolyi/togodo/todotxtlib"
)

// Controller handles TUI initialization and execution
type Controller struct {
	repo *todotxtlib.Repository
}

// NewController creates a new TUI controller
func NewController(repo *todotxtlib.Repository) *Controller {
	return &Controller{
		repo: repo,
	}
}

// Run starts the TUI interface
func (c *Controller) Run() error {
	model := initialModel(c.repo)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
