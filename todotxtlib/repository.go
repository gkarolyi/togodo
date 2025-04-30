package todotxtlib

import (
	"fmt"
	"sort"
)

// Repository handles storing and manipulating Todos
type Repository interface {
	// CRUD operations
	Add(todoText string) (Todo, error)
	Remove(index int) (Todo, error)
	// Update(todo Todo) error

	// Query operations
	// Find(filter Filter) ([]Todo, error)
	// GetByNumber(number int) (Todo, error)

	// Update operations
	// ToggleDone(todo Todo) error
	// SetPriority(todo Todo, priority string) error
	// SetContext(todo Todo, context string) error
	// SetProject(todo Todo, project string) error
	// SetTags(todo Todo, tags []string) error

	// List operations
	ListTodos() ([]Todo, error)
	ListDone() ([]Todo, error)
	ListProjects() ([]string, error)
	ListContexts() ([]string, error)

	Save() error
}

type repository struct {
	todos  []Todo
	reader Reader
	writer Writer
	path   string
}

func NewRepository(path string) (Repository, error) {
	reader := NewFileReader()
	writer := NewFileWriter()

	todos, err := reader.Read(path)
	if err != nil {
		return nil, err
	}

	return &repository{
		todos:  todos,
		reader: reader,
		writer: writer,
		path:   path,
	}, nil
}

// Add adds a todo to the repository
func (r *repository) Add(todoText string) (Todo, error) {
	newTodo := NewTodo(todoText)
	r.todos = append(r.todos, newTodo)
	return newTodo, nil
}

// Remove removes a todo from the repository
func (r *repository) Remove(index int) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	todo := r.todos[index]
	r.todos = append(r.todos[:index], r.todos[index+1:]...)
	return todo, nil
}

// ListTodos returns all todos
func (r *repository) ListTodos() ([]Todo, error) {
	return r.todos, nil
}

// ListDone returns all done todos
func (r *repository) ListDone() ([]Todo, error) {
	done := []Todo{}
	for _, todo := range r.todos {
		if todo.Done {
			done = append(done, todo)
		}
	}
	return done, nil
}

// ListProjects returns all projects
func (r *repository) ListProjects() ([]string, error) {
	projects := []string{}
	for _, todo := range r.todos {
		projects = append(projects, todo.Projects...)
	}
	sort.Strings(projects)
	return projects, nil
}

// ListContexts returns all contexts
func (r *repository) ListContexts() ([]string, error) {
	contexts := []string{}
	for _, todo := range r.todos {
		contexts = append(contexts, todo.Contexts...)
	}
	sort.Strings(contexts)
	return contexts, nil
}

// Save saves the todos to the file at the given path
func (r *repository) Save() error {
	return r.writer.Write(r.path, r.todos)
}
