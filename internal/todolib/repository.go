package todolib

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type TodoRepository struct {
	items       []Todo
	todoTxtPath string
}

// Items returns all items in the repository.
func (t TodoRepository) Items() []Todo {
	return t.items
}

// Todos returns all items that are not done.
func (t TodoRepository) Todos() []Todo {
	var todos []Todo
	for _, item := range t.Items() {
		if !item.Done {
			todos = append(todos, item)
		}
	}
	return todos
}

// Done returns all items that are done.
func (t TodoRepository) Done() []Todo {
	var done []Todo
	for _, item := range t.Items() {
		if item.Done {
			done = append(done, item)
		}
	}
	return done
}

// TodoCount returns the number of items that are not done.
func (t TodoRepository) TodoCount() int {
	return len(t.Todos())
}

// DoneCount returns the number of items that are done.
func (t TodoRepository) DoneCount() int {
	return len(t.Done())
}

// New reads the items from the todo.txt file and returns a new repository.
func New(todoTxtPath string) (TodoRepository, error) {
	repo := TodoRepository{todoTxtPath: todoTxtPath}
	err := repo.Read(todoTxtPath)
	if err != nil {
		return TodoRepository{}, err
	}
	return repo, nil
}

// Save writes the items to the todo.txt file.
func (t *TodoRepository) Save() error {
	file, err := os.Create(t.todoTxtPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, todo := range append(t.Todos(), t.Done()...) {
		_, err := writer.WriteString(todo.Text + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

// Add creates a new item from the given line and appends it to the repository.
func (t *TodoRepository) Add(line string) (Todo, error) {
	todo := Todo{Number: len(t.items) + 1, Text: line}

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
	t.items = append(t.items, todo)

	err := t.Save()
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

// Read adds the items from the given file to the repository.
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
	todo := t.Items()[index-1]
	return todo
}

func (t *TodoRepository) Do(index int) {
	t.items[index-1].ToggleDone()
}

func (t *TodoRepository) All() (todos []Todo) {
	t.items = sortByPriority(t.Todos(), t.Done())
	t.reassignNumbers()
	return t.Items()
}

func (t TodoRepository) Filter(query string) (matched []Todo) {
	for _, todo := range t.All() {
		if strings.Contains(strings.ToLower(todo.Text), strings.ToLower(query)) {
			matched = append(matched, todo)
		}
	}

	return matched
}

func (t *TodoRepository) reassignNumbers() {
	for i := range t.items {
		t.items[i].Number = i + 1
	}
}

func sortByPriority(todos, done []Todo) []Todo {
	sort.SliceStable(todos, func(i, j int) bool {
		iPrioritised := todos[i].Prioritised()
		jPrioritised := todos[j].Prioritised()

		if iPrioritised && jPrioritised {
			return todos[i].Priority < todos[j].Priority
		} else if iPrioritised {
			return true
		} else if jPrioritised {
			return false
		} else {
			return false
		}
	})

	sort.SliceStable(done, func(i, j int) bool {
		return done[i].Priority < done[j].Priority
	})

	return append(todos, done...)
}
