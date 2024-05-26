package todolib

import (
	"bufio"
	"os"
)

type TodoRepository struct {
	Todos []Todo
	Done  []Todo
}

func (t *TodoRepository) Add(line string) Todo {
	todo := Todo{line}
	t.Todos = append(t.Todos, todo)
	return todo
}

func (t *TodoRepository) ReadFile(path string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		t.Add(line)
	}

	return nil
}

func (t *TodoRepository) Find(lineNumber int) Todo {
	todo := t.Todos[lineNumber-1]
	return todo
}

func (t *TodoRepository) Do(lineNumber int) {
	t.Done = append(t.Done, t.Find(lineNumber))
}
