package todolib

import (
	"bufio"
	"errors"
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
	for _, todo := range t.All() {
		_, err := writer.WriteString(todo.Text + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

// Add creates a new item from the given line and appends it to the repository.
func (t *TodoRepository) Add(line string) ([]Todo, error) {
	var newTodos []Todo
	// Handle multi-line input
	for _, todoText := range strings.Split(line, "\n") {
		if strings.TrimSpace(todoText) == "" {
			continue
		}
		// Create new todo with temporary index/number
		todo := NewTodo(todoText, len(t.items))
		t.items = append(t.items, todo)
		newTodos = append(newTodos, todo)
	}

	// Sort all items and update numbers
	t.items = sortByPriority(t.items)
	t.reassignNumbers()

	// Update the numbers in our newTodos slice to match their sorted position
	for i := range newTodos {
		for _, todo := range t.items {
			if newTodos[i].Text == todo.Text {
				newTodos[i].Number = todo.Number
				break
			}
		}
	}

	err := t.Save()
	if err != nil {
		return nil, err
	}

	return newTodos, nil
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

func (t TodoRepository) Find(lineNumber int) Todo {
	todo := t.Items()[lineNumber-1]
	return todo
}

func (t *TodoRepository) Toggle(lineNumbers []int) ([]Todo, error) {
	var toggledTodos []Todo

	// First, toggle the specified todos
	for _, lineNumber := range lineNumbers {
		todo, err := t.get(lineNumber)
		if err != nil {
			return nil, err
		}
		todo.ToggleDone()
		// Store the original text for matching later
		toggledTodos = append(toggledTodos, *todo)
	}

	// Sort all items and update numbers
	t.items = sortByPriority(t.items)
	t.reassignNumbers()

	// Update the numbers in our toggledTodos slice to match their new sorted position
	for i := range toggledTodos {
		for _, todo := range t.items {
			if toggledTodos[i].Text == todo.Text {
				toggledTodos[i].Number = todo.Number
				break
			}
		}
	}

	err := t.Save()
	if err != nil {
		return nil, err
	}

	return toggledTodos, nil
}

func (t *TodoRepository) All() (todos []Todo) {
	t.items = sortByPriority(t.Items())
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

// Tidy removes all done items from the repository.
func (t *TodoRepository) Tidy() ([]Todo, error) {
	done := t.Done()
	// t.items = sortByPriority(t.Todos())
	t.items = t.Todos()
	err := t.Save()
	if err != nil {
		return nil, err
	}

	return done, nil
}

func (t *TodoRepository) get(lineNumber int) (*Todo, error) {
	if lineNumber < 1 || lineNumber > len(t.items) {
		return nil, errorInvalidLineNumber
	}
	return &t.items[lineNumber-1], nil
}

func (t *TodoRepository) reassignNumbers() {
	for i := range t.items {
		t.items[i].Number = i + 1
	}
}

// sortByPriority sorts the todos by priority, with done items last.
func sortByPriority(todos []Todo) []Todo {
	sort.SliceStable(todos, func(i, j int) bool {
		if todos[i].Done != todos[j].Done {
			return !todos[i].Done
		}
		if todos[i].Priority != todos[j].Priority {
			if todos[i].Priority == "" {
				return false
			}
			if todos[j].Priority == "" {
				return true
			}
			return todos[i].Priority < todos[j].Priority
		}
		return false
	})

	return todos
}

var errorInvalidLineNumber = errors.New("invalid line number")
