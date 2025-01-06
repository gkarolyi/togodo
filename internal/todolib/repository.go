package todolib

import (
	"bufio"
	"errors"
	"os"
	"sort"
	"strings"
)

type TodoRepository struct {
	itemsTable  map[string]Todo
	order       []string
	todoTxtPath string
}

// Items returns all items in the repository.
func (t TodoRepository) Items() []Todo {
	var items []Todo
	for _, text := range t.order {
		items = append(items, t.itemsTable[text])
	}
	return items
}

// Todos returns all items that are not done.
func (t TodoRepository) Todos() []Todo {
	var todos []Todo
	for _, text := range t.order {
		todo := t.itemsTable[text]
		if !todo.Done {
			todos = append(todos, todo)
		}
	}
	return todos
}

// Done returns all items that are done.
func (t TodoRepository) Done() []Todo {
	var done []Todo
	for _, text := range t.order {
		todo := t.itemsTable[text]
		if todo.Done {
			done = append(done, todo)
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
	repo := TodoRepository{
		todoTxtPath: todoTxtPath,
		itemsTable:  make(map[string]Todo),
		order:       make([]string, 0),
	}
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
func (t *TodoRepository) Add(line string) (Todo, error) {
	todo := NewTodo(line, len(t.itemsTable))
	t.itemsTable[todo.hash()] = todo
	t.order = append(t.order, todo.hash())

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

func (t TodoRepository) Find(lineNumber int) Todo {
	return t.itemsTable[t.order[lineNumber-1]]
}

func (t *TodoRepository) Toggle(lineNumbers []int) ([]Todo, error) {
	var modifiedTodos []Todo

	// Get all texts before modification
	todoHashes := make([]string, 0, len(lineNumbers))
	for _, ln := range lineNumbers {
		if ln < 1 || ln > len(t.order) {
			return nil, errorInvalidLineNumber
		}
		todoHashes = append(todoHashes, t.order[ln-1])
	}

	for _, hash := range todoHashes {
		todo := t.itemsTable[hash]
		oldHash := todo.hash()
		todo.ToggleDone()

		// Update itemsTable with new text
		t.itemsTable[todo.hash()] = todo
		delete(t.itemsTable, oldHash)

		// Update order slice with new text
		for i, orderText := range t.order {
			if orderText == oldHash {
				t.order[i] = todo.hash()
				break
			}
		}

		modifiedTodos = append(modifiedTodos, todo)
	}

	// Update order and numbers
	t.updateOrder()

	// Get final numbers for modified todos
	result := make([]Todo, len(modifiedTodos))
	for i, modifiedTodo := range modifiedTodos {
		for j, hash := range t.order {
			if hash == modifiedTodo.hash() { // Match by new text
				todo := t.itemsTable[hash]
				todo.Number = j + 1
				result[i] = todo
				break
			}
		}
	}

	err := t.Save()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *TodoRepository) All() (todos []Todo) {
	t.updateOrder()
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
	for _, todo := range done {
		delete(t.itemsTable, todo.Text)
	}

	t.updateOrder()

	err := t.Save()
	if err != nil {
		return nil, err
	}

	return done, nil
}

// updateOrder sorts the items by priority and updates the order slice.
func (t *TodoRepository) updateOrder() {
	// First, build a slice maintaining original order from t.order
	var todos []Todo
	for _, text := range t.order {
		if todo, exists := t.itemsTable[text]; exists {
			todos = append(todos, todo)
		}
	}

	// Sort todos while maintaining relative order within same priority groups
	todos = sortByPriority(todos)

	// Rebuild order slice and update numbers
	t.order = make([]string, len(todos))
	for i, todo := range todos {
		todo.Number = i + 1
		t.order[i] = todo.Text
		t.itemsTable[todo.Text] = todo
	}
}

func sortByPriority(todos []Todo) []Todo {
	sort.SliceStable(todos, func(i, j int) bool {
		// First separate done and not done
		if todos[i].Done != todos[j].Done {
			return !todos[i].Done
		}

		// Then sort by priority within each group
		if todos[i].Priority != todos[j].Priority {
			// If one has priority and other doesn't
			if todos[i].Priority == "" {
				return false
			}
			if todos[j].Priority == "" {
				return true
			}
			// Both have priorities, sort by priority
			return todos[i].Priority < todos[j].Priority
		}

		// If same done status and same priority, maintain original order
		return false
	})

	return todos
}

var errorInvalidLineNumber = errors.New("invalid line number")
