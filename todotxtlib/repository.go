package todotxtlib

// Repository handles storing and manipulating Todos
type Repository interface {
	// CRUD operations
	// Add(todo Todo) error
	// Remove(todo Todo) error
	// Update(todo Todo) error

	// Query operations
	// Find(filter Filter) ([]Todo, error)
	// All() ([]Todo, error)
	// GetByNumber(number int) (Todo, error)

	// Update operations
	// ToggleDone(todo Todo) error
	// SetPriority(todo Todo, priority string) error
	// SetContext(todo Todo, context string) error
	// SetProject(todo Todo, project string) error
	// SetTags(todo Todo, tags []string) error

	// List operations
	// ListContexts() ([]string, error)
	// ListProjects() ([]string, error)

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

func (r *repository) Save() error {
	return r.writer.Write(r.path, r.todos)
}
