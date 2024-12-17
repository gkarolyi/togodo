package todolib

import (
	"testing"
)

var TestTodoTxtPath = "../../testdata/todo.txt"

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// New API design:
// ONLY one export: a func .New which returns a TodoTxtRepository
// the TodoTxtRepository has these exported methods:
// .Add(text string)
// .All()
// .Filter(query string)
// .Find(lineNumber int)
// .Update(attribute, value string) or maybe a new Todo object altogether? ie. (lineNumber int, text string)
// exported func .New(path string) TodoTxtRepository

func TestAdd(t *testing.T) {
	t.Run("when todo is not done", func(t *testing.T) {
		repo := TodoRepository{}
		todo := repo.Add("(B) random fake task with a +projectName and @contextName @home")

		t.Run("adding to Todos", func(t *testing.T) {
			length := len(repo.Todos)
			if length != 1 {
				t.Errorf("expected 1 todo, got %d", length)
			}
		})

		t.Run("adding to Done", func(t *testing.T) {
			doneLength := len(repo.Done)
			if doneLength != 0 {
				t.Errorf("expected 0 done todos, got %d", doneLength)
			}
		})

		t.Run("last Todo matches", func(t *testing.T) {
			got := repo.Todos[len(repo.Todos)-1]
			if !got.Equals(todo) {
				t.Errorf("expected last todo to be %v, got %v", todo, got)
			}
		})

		t.Run("priority", func(t *testing.T) {
			if todo.Priority != "B" {
				t.Errorf("expected priority to be B, got %s", todo.Priority)
			}
		})

		t.Run("projects", func(t *testing.T) {
			expectedProjects := []string{"+projectName"}
			if !equalSlices(todo.Projects, expectedProjects) {
				t.Errorf("expected projects to be %v, got %v", expectedProjects, todo.Projects)
			}
		})

		t.Run("contexts", func(t *testing.T) {
			expectedContexts := []string{"@contextName", "@home"}
			if !equalSlices(todo.Contexts, expectedContexts) {
				t.Errorf("expected contexts to be %v, got %v", expectedContexts, todo.Contexts)
			}
		})

		t.Run("done status", func(t *testing.T) {
			if todo.Done != false {
				t.Errorf("expected done to be false, got %v", todo.Done)
			}
		})

		t.Run("line number", func(t *testing.T) {
			if todo.Number != 1 {
				t.Errorf("expected line number to be 1, got %d", todo.Number)
			}
			todo = repo.Add("another todo item to increment number")
			if todo.Number != 2 {
				t.Errorf("expected line number to be 2, got %d", todo.Number)
			}
		})

	})

	t.Run("when todo is done", func(t *testing.T) {
		repo := TodoRepository{}
		todo := repo.Add("x (B) random fake task with a +projectName and @contextName @home")

		t.Run("adding to Todos", func(t *testing.T) {
			todosLength := len(repo.Todos)
			if todosLength != 0 {
				t.Errorf("expected 0 todos, got %d", todosLength)
			}
		})

		t.Run("adding to Done", func(t *testing.T) {
			doneLength := len(repo.Done)
			if doneLength != 1 {
				t.Errorf("expected 1 done todo, got %d", doneLength)
			}
		})

		t.Run("last Todo matches", func(t *testing.T) {
			got := repo.Done[len(repo.Done)-1]
			if !got.Equals(todo) {
				t.Errorf("expected last todo to be %v, got %v", todo, got)
			}
		})

		t.Run("priority", func(t *testing.T) {
			if todo.Priority != "" {
				t.Errorf("expected priority to be empty, got %s", todo.Priority)
			}
		})

		t.Run("projects", func(t *testing.T) {
			if todo.Projects != nil {
				t.Errorf("expected projects to be nil, got %v", todo.Projects)
			}
		})

		t.Run("contexts", func(t *testing.T) {
			if todo.Contexts != nil {
				t.Errorf("expected contexts to be nil, got %v", todo.Contexts)
			}
		})

		t.Run("done status", func(t *testing.T) {
			if todo.Done != true {
				t.Errorf("expected done to be true, got %v", todo.Done)
			}
		})

		t.Run("line number", func(t *testing.T) {
			if todo.Number != 1 {
				t.Errorf("expected line number to be 1, got %d", todo.Number)
			}
			todo = repo.Add("another todo item to increment number")
			if todo.Number != 2 {
				t.Errorf("expected line number to be 2, got %d", todo.Number)
			}
		})
	})
}

func TestReadFile(t *testing.T) {
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("Todos length", func(t *testing.T) {
		length := len(repo.Todos)

		if length != 5 {
			t.Errorf("expected 5 todos, got %d", length)
		}
	})

	t.Run("index", func(t *testing.T) {
		lastNumber := repo.Todos[len(repo.Todos)-1].Number
		if lastNumber != 7 {
			t.Errorf("expected last todo number to be 7, got %d", lastNumber)
		}
	})

	// t.Run("sort done items at the bottom", func(t *testing.T) {
	// 	lastTodo := repo.Todos[len(repo.Todos)-1]
	// 	if !lastTodo.Done {
	// 		t.Errorf("expected last todo to be done")
	// 	}
	// })
}

func TestFind(t *testing.T) {
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("find correct item", func(t *testing.T) {
		got := repo.Find(1)
		want := Todo{
			Text: "do some super cool hacker stuff +hacking @basement",
		}

		if got.Text != want.Text {
			t.Errorf("expected %v, got %v", want.Text, got.Text)
		}
	})
}

func TestDo(t *testing.T) {
	t.Run("toggling not done item", func(t *testing.T) {
		repo := TodoRepository{}
		err := repo.ReadFile(TestTodoTxtPath)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		initialDoneLength := len(repo.Done)
		initialTodosLength := len(repo.Todos)
		repo.Do(1)

		t.Run("item in done", func(t *testing.T) {
			doneLength := len(repo.Done)
			if doneLength != initialDoneLength+1 {
				t.Errorf("expected %d done items, got %d", initialDoneLength+1, doneLength)
			}
		})

		t.Run("item not in todos", func(t *testing.T) {
			todosLength := len(repo.Todos)
			if todosLength != initialTodosLength-1 {
				t.Errorf("expected %d todos, got %d", initialTodosLength-1, todosLength)
			}
		})
	})

	// t.Run("toggling done item", func(t *testing.T) {
	// 	repo := TodoRepository{}
	// 	err := repo.ReadFile(TestTodoTxtPath)

	// 	if err != nil {
	// 		t.Fatalf("expected no error, got %v", err)
	// 	}

	// 	initialDoneLength := len(repo.Done)
	// 	if initialDoneLength != 2 {
	// 		t.Errorf("expected 2 done items, got %d", initialDoneLength)
	// 	}
	// 	initialTodosLength := len(repo.Todos)
	// 	if initialTodosLength != 5 {
	// 		t.Errorf("expected 5 todos, got %d", initialTodosLength)
	// 	}
	// 	repo.Toggle(4)

	// 	t.Run("item not in done", func(t *testing.T) {
	// 		doneLength := len(repo.Done)
	// 		if doneLength != initialDoneLength-1 {
	// 			t.Errorf("expected %d done items, got %d", initialDoneLength-1, doneLength)
	// 		}
	// 	})

	// t.Run("item back in todos", func(t *testing.T) {
	// 	todosLength := len(repo.Todos)
	// 	if todosLength != initialTodosLength+1 {
	// 		t.Errorf("expected %d todos, got %d", initialTodosLength+1, todosLength)
	// 	}
	// })
	// })
}

func TestAll(t *testing.T) {
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	todos := repo.All()

	t.Run("number of items correct", func(t *testing.T) {
		if len(todos) != 7 {
			t.Errorf("expected 7 todos, got %d", len(todos))
		}
	})

	t.Run("sort done to the bottom", func(t *testing.T) {
		lastTodo := todos[len(todos)-1]
		if lastTodo.Text != "x and here is another done todo" {
			t.Errorf("expected last todo text to be 'x and here is another done todo', got %v", lastTodo.Text)
		}
	})

	t.Run("sort by priority", func(t *testing.T) {
		topTodo := todos[0]
		if topTodo.Priority != "A" {
			t.Errorf("expected top todo priority to be 'A', got %v", topTodo.Priority)
		}
	})
}

func TestFilter(t *testing.T) {
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("number of items correct", func(t *testing.T) {
		todos := repo.Filter("@basement")
		if len(todos) != 1 {
			t.Errorf("expected 1 matching todo, got %d", len(todos))
		}
	})

	t.Run("number of items correct", func(t *testing.T) {
		todos := repo.Filter("@shop")
		if len(todos) != 1 {
			t.Errorf("expected 1 matching todo, got %d", len(todos))
		}
	})

	t.Run("number of items correct", func(t *testing.T) {
		todos := repo.Filter("(A)")
		if len(todos) != 1 {
			t.Errorf("expected 1 matching todo, got %d", len(todos))
		}
	})
}
