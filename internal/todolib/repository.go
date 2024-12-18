package todolib

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type TodoRepository struct {
	Todos       []Todo
	Done        []Todo
	todoCount   int
	todoTxtPath string
}

func New(todoTxtPath string) (TodoRepository, error) {
	repo := TodoRepository{todoTxtPath: todoTxtPath}
	err := repo.Read(todoTxtPath)
	if err != nil {
		return TodoRepository{}, err
	}
	return repo, nil
}

func (t *TodoRepository) Save() error {
	file, err := os.Create(t.todoTxtPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, todo := range append(t.Todos, t.Done...) {
		_, err := writer.WriteString(todo.Text + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

func (t *TodoRepository) TodoCount() int {
	t.todoCount++
	return t.todoCount
}

func (t *TodoRepository) Add(line string) (Todo, error) {
	todo := Todo{Number: t.TodoCount(), Text: line}

	if doneRe.MatchString(line) {
		todo.Done = true
		t.Done = append(t.Done, todo)
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
		t.Todos = append(t.Todos, todo)
	}

	// err := t.Save()
	// if err != nil {
	// 	return Todo{}, err
	// }

	return todo, nil
}

func (t *TodoRepository) Read(path string) error {
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

func (t TodoRepository) Find(index int) Todo {
	todo := t.Todos[index-1]
	return todo
}

func (t *TodoRepository) Do(index int) {
	todo := t.Find(index)
	if !todo.Done {
		t.Done = append(t.Done, todo)
		t.Todos = removeIndex(t.Todos, index-1)
	}
}

func (t TodoRepository) All() (todos []Todo) {
	sort.SliceStable(t.Todos, func(i, j int) bool {
		iPrioritised := t.Todos[i].Prioritised()
		jPrioritised := t.Todos[j].Prioritised()

		if iPrioritised && jPrioritised {
			return t.Todos[i].Priority < t.Todos[j].Priority
		} else if iPrioritised {
			return true
		} else if jPrioritised {
			return false
		} else {
			return false
		}
	})

	sort.SliceStable(t.Done, func(i, j int) bool {
		return t.Done[i].Priority < t.Done[j].Priority
	})

	return append(t.Todos, t.Done...)
}

func (t TodoRepository) Filter(query string) (matched []Todo) {
	for _, todo := range t.All() {
		if strings.Contains(strings.ToLower(todo.Text), strings.ToLower(query)) {
			matched = append(matched, todo)
		}
	}

	return matched
}

func removeIndex(s []Todo, index int) []Todo {
	return append(s[:index], s[index+1:]...)
}
