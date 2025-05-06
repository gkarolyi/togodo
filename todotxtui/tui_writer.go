package todotxtui

import (
	tui "github.com/gkarolyi/togodo/todotxt-tui"
	"github.com/gkarolyi/togodo/todotxtlib"
)

// TUIWriter implements OutputWriter using bubbletea for interactive output
type TUIWriter struct {
	repo *todotxtlib.Repository
}

// NewTUIWriter creates a new TUIWriter
func NewTUIWriter(repo *todotxtlib.Repository) *TUIWriter {
	return &TUIWriter{
		repo: repo,
	}
}

// WriteLine is not used in TUI mode as we show everything in the interactive interface
func (w *TUIWriter) WriteLine(line string) {
	// No-op in TUI mode
}

// WriteLines is not used in TUI mode as we show everything in the interactive interface
func (w *TUIWriter) WriteLines(lines []string) {
	// No-op in TUI mode
}

// WriteError shows errors in the TUI
func (w *TUIWriter) WriteError(err error) {
	// Errors will be shown in the TUI interface
}

// Run starts the TUI interface
func (w *TUIWriter) Run() error {
	p := tui.NewProgram(w.repo)
	_, err := p.Run()
	return err
}
