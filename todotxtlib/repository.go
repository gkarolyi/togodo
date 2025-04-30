package todotxtlib

import (
	"fmt"
	"sort"
)

// Repository handles storing and manipulating Todos
type Repository struct {
	todos  []Todo
	reader Reader
	writer Writer
}

// NewRepository creates a new repository with custom reader and writer
func NewRepository(reader Reader, writer Writer) (*Repository, error) {
	todos, err := reader.Read()
	if err != nil {
		return nil, err
	}

	return &Repository{
		todos:  todos,
		reader: reader,
		writer: writer,
	}, nil
}

// Add adds a todo to the repository
func (r *Repository) Add(todoText string) (Todo, error) {
	newTodo := NewTodo(todoText)
	r.todos = append(r.todos, newTodo)
	return newTodo, nil
}

// Remove removes a todo from the repository
func (r *Repository) Remove(index int) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	todo := r.todos[index]
	r.todos = append(r.todos[:index], r.todos[index+1:]...)
	return todo, nil
}

// ListTodos returns all todos
func (r *Repository) ListTodos() ([]Todo, error) {
	return r.todos, nil
}

// ListDone returns all done todos
func (r *Repository) ListDone() ([]Todo, error) {
	done := []Todo{}
	for _, todo := range r.todos {
		if todo.Done {
			done = append(done, todo)
		}
	}
	return done, nil
}

// ListProjects returns all projects
func (r *Repository) ListProjects() ([]string, error) {
	projects := []string{}
	for _, todo := range r.todos {
		projects = append(projects, todo.Projects...)
	}
	sort.Strings(projects)
	return projects, nil
}

// ListContexts returns all contexts
func (r *Repository) ListContexts() ([]string, error) {
	contexts := []string{}
	for _, todo := range r.todos {
		contexts = append(contexts, todo.Contexts...)
	}
	sort.Strings(contexts)
	return contexts, nil
}

// Save saves the todos using the configured writer
func (r *Repository) Save() error {
	return r.writer.Write(r.todos)
}
