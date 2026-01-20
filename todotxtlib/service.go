package todotxtlib

import "fmt"

// TodoService provides high-level operations for managing todos
type TodoService interface {
	AddTodos(texts []string) ([]Todo, error)
	ToggleTodos(indices []int) ([]Todo, error)
	SetPriorities(indices []int, priority string) ([]Todo, error)
	RemoveDoneTodos() ([]Todo, error)
	SearchTodos(query string) ([]Todo, error)
}

// DefaultTodoService implements TodoService using a TodoRepository
type DefaultTodoService struct {
	repo TodoRepository
}

// NewTodoService creates a new TodoService with the given repository
func NewTodoService(repo TodoRepository) TodoService {
	return &DefaultTodoService{repo: repo}
}

// AddTodos adds multiple todos, sorts the list, and saves
// Returns the added todos
func (s *DefaultTodoService) AddTodos(texts []string) ([]Todo, error) {
	addedTodos := make([]Todo, 0, len(texts))

	for _, text := range texts {
		todo, err := s.repo.Add(text)
		if err != nil {
			return nil, fmt.Errorf("failed to add todo: %w", err)
		}
		addedTodos = append(addedTodos, todo)
	}

	s.repo.Sort(nil)
	if err := s.repo.Save(); err != nil {
		return nil, fmt.Errorf("failed to save todos: %w", err)
	}

	return addedTodos, nil
}

// ToggleTodos toggles the done status of todos at the given indices (0-based)
// Returns the toggled todos
func (s *DefaultTodoService) ToggleTodos(indices []int) ([]Todo, error) {
	toggledTodos := make([]Todo, 0, len(indices))

	for _, index := range indices {
		todo, err := s.repo.ToggleDone(index)
		if err != nil {
			return nil, fmt.Errorf("failed to toggle todo at index %d: %w", index, err)
		}
		toggledTodos = append(toggledTodos, todo)
	}

	s.repo.Sort(nil)
	if err := s.repo.Save(); err != nil {
		return nil, fmt.Errorf("failed to save todos: %w", err)
	}

	return toggledTodos, nil
}

// SetPriorities sets the priority for todos at the given indices (0-based)
// Returns the updated todos
// Note: Does not sort after setting priorities to preserve user's intended order
func (s *DefaultTodoService) SetPriorities(indices []int, priority string) ([]Todo, error) {
	updatedTodos := make([]Todo, 0, len(indices))

	for _, index := range indices {
		todo, err := s.repo.SetPriority(index, priority)
		if err != nil {
			return nil, fmt.Errorf("failed to set priority for todo at index %d: %w", index, err)
		}
		updatedTodos = append(updatedTodos, todo)
	}

	// Note: Pri command doesn't sort - preserves user's order
	if err := s.repo.Save(); err != nil {
		return nil, fmt.Errorf("failed to save todos: %w", err)
	}

	return updatedTodos, nil
}

// RemoveDoneTodos removes all completed todos
// Returns the removed todos
func (s *DefaultTodoService) RemoveDoneTodos() ([]Todo, error) {
	// Get done todos before removing
	doneTodos, err := s.repo.ListDone()
	if err != nil {
		return nil, fmt.Errorf("failed to list done todos: %w", err)
	}

	// Get all todos to iterate
	allTodos, err := s.repo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list all todos: %w", err)
	}

	// Remove backwards to avoid index shifting
	for i := len(allTodos) - 1; i >= 0; i-- {
		if allTodos[i].Done {
			if _, err := s.repo.Remove(i); err != nil {
				return nil, fmt.Errorf("failed to remove todo at index %d: %w", i, err)
			}
		}
	}

	s.repo.Sort(nil)
	if err := s.repo.Save(); err != nil {
		return nil, fmt.Errorf("failed to save todos: %w", err)
	}

	return doneTodos, nil
}

// SearchTodos searches for todos matching the given query
// Returns matching todos
func (s *DefaultTodoService) SearchTodos(query string) ([]Todo, error) {
	return s.repo.Search(query)
}
