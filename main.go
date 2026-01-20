/*
Copyright © 2024 Gergely Karolyi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files ("the Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/internal/config"
	"github.com/gkarolyi/togodo/todotxtlib"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	todoTxtPath := config.GetTodoTxtPath()
	reader := todotxtlib.NewFileReader(todoTxtPath)
	writer := todotxtlib.NewFileWriter(todoTxtPath)

	repo, err := todotxtlib.NewFileRepository(reader, writer)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Create service layer
	service := todotxtlib.NewTodoService(repo)

	presenter := cli.NewPresenter()

	rootCmd := cmd.NewRootCmd(service, repo, presenter)

	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}
