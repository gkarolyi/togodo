package todolib

import (
	"bufio"
	"os"
	"sort"
	"strings"
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
		t.Done = append(t.Done, todo)

		return todo

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

func (t TodoRepository) Find(lineNumber int) Todo {
	todo := t.Todos[lineNumber-1]
	return todo
}

func (t *TodoRepository) Do(lineNumber int) {
	todo := t.Find(lineNumber)
	if !todo.Done {
		t.Done = append(t.Done, todo)
		t.Todos = removeIndex(t.Todos, lineNumber-1)
	}
}

func (t TodoRepository) All() (todos []Todo) {
	sort.SliceStable(t.Todos, func(i, j int) bool {
		iPrioritised := t.Todos[i].Prioritised()
		jPrioritised := t.Todos[j].Prioritised()

		if jPrioritised && !iPrioritised {
			return false
		} else if iPrioritised && !jPrioritised {
			return true
		} else {
			return t.Todos[i].Priority < t.Todos[j].Priority
		}
	})

	return append(t.Todos, t.Done...)
}

func (t TodoRepository) Filter(query string) (matched []Todo) {
	for _, todo := range t.All() {
		if strings.Contains(todo.Text, query) {
			matched = append(matched, todo)
		}
	}

	return matched
}

func removeIndex(s []Todo, index int) []Todo {
	return append(s[:index], s[index+1:]...)
}
