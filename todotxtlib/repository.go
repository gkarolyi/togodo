package todotxtlib

import (
	"bytes"
	"fmt"
	"sort"
)

// TodoRepository defines the interface for storing and manipulating Todos.
type TodoRepository interface {
	Add(todoText string) (Todo, error)
	Get(index int) (Todo, error)
	Remove(index int) (Todo, error)
	Update(index int, todo Todo) (Todo, error)
	ToggleDone(index int) (Todo, error)
	SetPriority(index int, priority string) (Todo, error)
	SetContexts(index int, contexts []string) (Todo, error)
	SetProjects(index int, projects []string) (Todo, error)
	AddContext(index int, context string) (Todo, error)
	AddProject(index int, project string) (Todo, error)
	RemoveContext(index int, context string) (Todo, error)
	RemoveProject(index int, project string) (Todo, error)
	Filter(filter Filter) ([]Todo, error)
	Sort(sort *Sort)
	ListAll() ([]Todo, error)
	ListTodos() ([]Todo, error)
	ListDone() ([]Todo, error)
	ListProjects() ([]string, error)
	ListContexts() ([]string, error)
	Save() error
	WriteToString() (string, error)
}

// FileRepository handles storing and manipulating Todos in a file.
type FileRepository struct {
	todos  []Todo
	reader Reader
	writer Writer
}

// NewFileRepository creates a new repository with custom reader and writer
func NewFileRepository(reader Reader, writer Writer) (TodoRepository, error) {
	todos, err := reader.Read()
	if err != nil {
		return nil, err
	}

	return &FileRepository{
		todos:  todos,
		reader: reader,
		writer: writer,
	}, nil
}

// Add adds a todo to the repository
func (r *FileRepository) Add(todoText string) (Todo, error) {
	newTodo := NewTodo(todoText)
	r.todos = append(r.todos, newTodo)
	return newTodo, nil
}

// Get returns a todo at the given index
func (r *FileRepository) Get(index int) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	return r.todos[index], nil
}

// Remove removes a todo from the repository
func (r *FileRepository) Remove(index int) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	todo := r.todos[index]
	r.todos = append(r.todos[:index], r.todos[index+1:]...)
	return todo, nil
}

// Update updates a todo in the repository
func (r *FileRepository) Update(index int, todo Todo) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	r.todos[index] = todo
	return todo, nil
}

// ToggleDone toggles the done status of a todo
func (r *FileRepository) ToggleDone(index int) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	r.todos[index].ToggleDone()
	return r.todos[index], nil
}

// SetPriority sets the priority of a todo
func (r *FileRepository) SetPriority(index int, priority string) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	r.todos[index].SetPriority(priority)
	return r.todos[index], nil
}

// SetContexts sets the contexts of a todo
func (r *FileRepository) SetContexts(index int, contexts []string) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	r.todos[index].SetContexts(contexts)
	return r.todos[index], nil
}

// SetProjects sets the projects of a todo
func (r *FileRepository) SetProjects(index int, projects []string) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	r.todos[index].SetProjects(projects)
	return r.todos[index], nil
}

// AddContext adds a context to a todo if it doesn't already exist
func (r *FileRepository) AddContext(index int, context string) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	r.todos[index].AddContext(context)
	return r.todos[index], nil
}

// AddProject adds a project to a todo if it doesn't already exist
func (r *FileRepository) AddProject(index int, project string) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	r.todos[index].AddProject(project)
	return r.todos[index], nil
}

// RemoveContext removes a context from a todo
func (r *FileRepository) RemoveContext(index int, context string) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	r.todos[index].RemoveContext(context)
	return r.todos[index], nil
}

// RemoveProject removes a project from a todo
func (r *FileRepository) RemoveProject(index int, project string) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	r.todos[index].RemoveProject(project)
	return r.todos[index], nil
}

// Filter returns todos that match all the specified criteria
func (r FileRepository) Filter(filter Filter) ([]Todo, error) {
	return filter.Apply(r.todos), nil
}

// Sort sorts the todos in the repository according to the specified criteria
// Pass nil to use default sort
func (r *FileRepository) Sort(sort *Sort) {
	if sort == nil {
		defaultSort := NewDefaultSort()
		sort = &defaultSort
	}
	sort.Apply(r.todos)
}

// ListAll returns all todos
func (r FileRepository) ListAll() ([]Todo, error) {
	return r.todos, nil
}

// ListTodos returns all todos that are not done
func (r FileRepository) ListTodos() ([]Todo, error) {
	notDone := []Todo{}
	for _, todo := range r.todos {
		if !todo.Done {
			notDone = append(notDone, todo)
		}
	}
	return notDone, nil
}

// ListDone returns all done todos
func (r FileRepository) ListDone() ([]Todo, error) {
	done := []Todo{}
	for _, todo := range r.todos {
		if todo.Done {
			done = append(done, todo)
		}
	}
	return done, nil
}

// ListProjects returns all unique projects sorted alphabetically
func (r FileRepository) ListProjects() ([]string, error) {
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
func (r FileRepository) ListContexts() ([]string, error) {
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
func (r *FileRepository) Save() error {
	return r.writer.Write(r.todos)
}

// WriteToString returns the todos as a string representation
func (r *FileRepository) WriteToString() (string, error) {
	var buffer bytes.Buffer
	writer := NewBufferWriter(&buffer)
	err := writer.Write(r.todos)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
