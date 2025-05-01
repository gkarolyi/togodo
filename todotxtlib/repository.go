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

// Update updates a todo in the repository
func (r *Repository) Update(index int, todo Todo) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	r.todos[index] = todo
	return todo, nil
}

// ListAll returns all todos
func (r Repository) ListAll() ([]Todo, error) {
	return r.todos, nil
}

// ListTodos returns all todos that are not done
func (r Repository) ListTodos() ([]Todo, error) {
	notDone := []Todo{}
	for _, todo := range r.todos {
		if !todo.Done {
			notDone = append(notDone, todo)
		}
	}
	return notDone, nil
}

// ListDone returns all done todos
func (r Repository) ListDone() ([]Todo, error) {
	done := []Todo{}
	for _, todo := range r.todos {
		if todo.Done {
			done = append(done, todo)
		}
	}
	return done, nil
}

// ListProjects returns all unique projects sorted alphabetically
func (r Repository) ListProjects() ([]string, error) {
	projectMap := make(map[string]struct{})
	for _, todo := range r.todos {
		for _, project := range todo.Projects {
			projectMap[project] = struct{}{}
		}
	}

	projects := make([]string, 0, len(projectMap))
	for project := range projectMap {
		projects = append(projects, project)
	}

	sort.Strings(projects)
	return projects, nil
}

// ListContexts returns all unique contexts sorted alphabetically
func (r Repository) ListContexts() ([]string, error) {
	contextMap := make(map[string]struct{})
	for _, todo := range r.todos {
		for _, context := range todo.Contexts {
			contextMap[context] = struct{}{}
		}
	}

	contexts := make([]string, 0, len(contextMap))
	for context := range contextMap {
		contexts = append(contexts, context)
	}

	sort.Strings(contexts)
	return contexts, nil
}

// Save saves the todos using the configured writer
func (r *Repository) Save() error {
	return r.writer.Write(r.todos)
}
