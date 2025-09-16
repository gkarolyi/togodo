package cmd

import (
	"github.com/gkarolyi/togodo/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/gkarolyi/togodo/tui"
)

// Factory functions for creating dependencies
func createRepository() (*todotxtlib.Repository, error) {
	todoTxtPath := GetTodoTxtPath()
	reader := todotxtlib.NewFileReader(todoTxtPath)
	writer := todotxtlib.NewFileWriter(todoTxtPath)
	return todotxtlib.NewRepository(reader, writer)
}

func createCLIPresenter() *cli.Presenter {
	return cli.NewPresenter()
}

func createTUIController(repo *todotxtlib.Repository) interface{ Run() error } {
	return tui.NewController(repo)
}
