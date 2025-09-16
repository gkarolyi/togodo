package injector

import (
	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/internal/config"
	"github.com/gkarolyi/togodo/internal/tui"
	"github.com/gkarolyi/togodo/todotxtlib"
)

// CreateRepository creates a new todotxtlib.Repository
func CreateRepository() (*todotxtlib.Repository, error) {
	todoTxtPath := config.GetTodoTxtPath()
	reader := todotxtlib.NewFileReader(todoTxtPath)
	writer := todotxtlib.NewFileWriter(todoTxtPath)
	return todotxtlib.NewRepository(reader, writer)
}

// CreateCLIPresenter creates a new cli.Presenter
func CreateCLIPresenter() *cli.Presenter {
	return cli.NewPresenter()
}

// CreateTUIController creates a new tui.Controller
func CreateTUIController(repo *todotxtlib.Repository) interface{ Run() error } {
	return tui.NewController(repo)
}
