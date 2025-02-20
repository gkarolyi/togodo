package todolib

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func List(todoTxtPath string, args []string) {
	repo, err := New(todoTxtPath)
	if err != nil {
		log.Fatal(err)
	}

	var todos []Todo
	if len(args) == 0 {
		todos = repo.All()
	} else {
		query := strings.Join(args, " ")
		todos = repo.Filter(query)
	}

	for _, todo := range todos {
		render(todo)
	}
}

func Add(todoTxtPath string, args []string) {
	repo, err := New(todoTxtPath)
	if err != nil {
		fmt.Println(err)
	}

	todos, err := repo.Add(args[0])
	if err != nil {
		fmt.Println(err)
	}
	for _, todo := range todos {
		render(todo)
	}
}

func Do(todoTxtPath string, args []string) {
	repo, err := New(todoTxtPath)
	if err != nil {
		fmt.Println(err)
	}
	var lineNumbers []int
	for _, arg := range args {
		lineNumber, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Invalid argument:", arg)
			continue
		}
		lineNumbers = append(lineNumbers, lineNumber)
	}
	todos, err := repo.Toggle(lineNumbers)
	if err != nil {
		fmt.Println(err)
	}
	for _, todo := range todos {
		render(todo)
	}
}
func Pri(todoTxtPath string, args []string) {
	repo, err := New(todoTxtPath)
	if err != nil {
		fmt.Println(err)
	}

	var lineNumbers []int
	var priority string
	for i, arg := range args {
		lineNumber, err := strconv.Atoi(arg)
		if err != nil {
			if i == len(args)-1 {
				priority = arg
			} else {
				fmt.Println("Invalid argument:", arg)
			}
			continue
		}
		lineNumbers = append(lineNumbers, lineNumber)
	}

	todos, err := repo.SetPriority(lineNumbers, priority)
	if err != nil {
		fmt.Println(err)
	}
	for _, todo := range todos {
		render(todo)
	}
}
func Tidy(todoTxtPath string, args []string) {
	repo, err := New(todoTxtPath)
	if err != nil {
		fmt.Println(err)
	}
	todos, err := repo.Tidy()

	if err != nil {
		fmt.Println(err)
	}
	for _, todo := range todos {
		fmt.Println(todo.Text)
	}
}
