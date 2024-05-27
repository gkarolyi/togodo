package todolib

import (
	"bufio"
	"os"
)

type TodoRepository struct {
	Todos     []Todo
	Done      []Todo
	todoCount int
}

func (t *TodoRepository) TodoCount() int {
	t.todoCount++
	return t.todoCount
}

func (t *TodoRepository) Add(line string) Todo {
	todo := Todo{Number: t.TodoCount(), Text: line}

	if doneRe.MatchString(line) {
		todo.Done = true
	} else {
		if priorityRe.MatchString(line) {
			todo.Priority = priorityRe.FindStringSubmatch(line)[1]
		}
		if projectRe.MatchString(line) {
			todo.Projects = projectRe.FindAllString(line, -1)
		}
		if contextRe.MatchString(line) {
			todo.Contexts = contextRe.FindAllString(line, -1)
		}
	}

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
