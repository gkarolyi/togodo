package todotxtlib

import (
	"testing"
)

// generateTestTodos creates a slice of todos with random projects and contexts
func generateTestTodos(count int, projectsPerTodo, contextsPerTodo int) []Todo {
	todos := make([]Todo, count)
	for i := 0; i < count; i++ {
		projects := make([]string, projectsPerTodo)
		contexts := make([]string, contextsPerTodo)

		// Generate projects
		for j := 0; j < projectsPerTodo; j++ {
			projects[j] = "+project" + string(rune('a'+j%5)) // Reuse 5 different projects
		}

		// Generate contexts
		for j := 0; j < contextsPerTodo; j++ {
			contexts[j] = "@context" + string(rune('a'+j%5)) // Reuse 5 different contexts
		}

		todos[i] = Todo{
			Text:     "Test todo",
			Done:     false,
			Projects: projects,
			Contexts: contexts,
		}
	}
	return todos
}

func BenchmarkListContextsSmall(b *testing.B) {
	// Small test: 100 todos, 2 contexts each
	todos := generateTestTodos(100, 2, 2)
	repo := Repository{todos: todos}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.ListContexts()
	}
}

func BenchmarkListContextsLarge(b *testing.B) {
	// Large test: 2000 todos, 3 contexts each
	todos := generateTestTodos(2000, 3, 3)
	repo := Repository{todos: todos}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.ListContexts()
	}
}

func BenchmarkListProjectsSmall(b *testing.B) {
	// Small test: 100 todos, 2 projects each
	todos := generateTestTodos(100, 2, 2)
	repo := Repository{todos: todos}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.ListProjects()
	}
}

func BenchmarkListProjectsLarge(b *testing.B) {
	// Large test: 2000 todos, 3 projects each
	todos := generateTestTodos(2000, 3, 3)
	repo := Repository{todos: todos}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.ListProjects()
	}
}
