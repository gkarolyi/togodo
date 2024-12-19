package todolib

import (
	"os"
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

func TestNew(t *testing.T) {
	repo, err := New(TestTodoTxtPath)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("setting the corrext todoTxtPath", func(t *testing.T) {
		if repo.todoTxtPath != TestTodoTxtPath {
			t.Errorf("expected todoTxtPath to be %s, got %s", TestTodoTxtPath, repo.todoTxtPath)
		}
	})

	t.Run("initialising Todos", func(t *testing.T) {
		if len(repo.Todos()) != 5 {
			t.Errorf("expected 5 todos, got %d", len(repo.Todos()))
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		_, err := New("nonexistent.txt")
		if err == nil {
			t.Errorf("expected an error, got nil")
		}
	})
}

func TestAdd(t *testing.T) {
	t.Run("when todo is not done", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test_todo.txt")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		repo, err := New(tmpfile.Name())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		todo, err := repo.Add("(B) random fake task with a +projectName and @contextName @home")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		t.Run("adding to Todos", func(t *testing.T) {
			length := len(repo.Todos())
			if length != 1 {
				t.Errorf("expected 1 todo, got %d", length)
			}
		})

		t.Run("adding to Done", func(t *testing.T) {
			doneLength := len(repo.Done())
			if doneLength != 0 {
				t.Errorf("expected 0 done todos, got %d", doneLength)
			}
		})

		t.Run("last Todo matches", func(t *testing.T) {
			got := repo.Todos()[len(repo.Todos())-1]
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
			todo, err = repo.Add("another todo item to increment number")
			if todo.Number != 2 {
				t.Errorf("expected line number to be 2, got %d", todo.Number)
			}
		})

		t.Run("adding to file", func(t *testing.T) {
			content, err := os.ReadFile(tmpfile.Name())
			if err != nil {
				t.Fatalf("failed to read temp file: %v", err)
			}

			expectedContent := "(B) random fake task with a +projectName and @contextName @home\nanother todo item to increment number\n"
			if string(content) != expectedContent {
				t.Errorf("expected content to be %s, got %s", expectedContent, string(content))
			}
		})
	})

	t.Run("when todo is done", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test_todo.txt")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		repo, err := New(tmpfile.Name())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		todo, err := repo.Add("x (B) random fake task with a +projectName and @contextName @home")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		t.Run("adding to Todos", func(t *testing.T) {
			todosLength := len(repo.Todos())
			if todosLength != 0 {
				t.Errorf("expected 0 todos, got %d", todosLength)
			}
		})

		t.Run("adding to Done", func(t *testing.T) {
			doneLength := len(repo.Done())
			if doneLength != 1 {
				t.Errorf("expected 1 done todo, got %d", doneLength)
			}
		})

		t.Run("last Todo matches", func(t *testing.T) {
			got := repo.Done()[len(repo.Done())-1]
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
			todo, err = repo.Add("x another done todo item to increment number")
			if todo.Number != 2 {
				t.Errorf("expected line number to be 2, got %d", todo.Number)
			}
		})
	})
}

func TestRead(t *testing.T) {
	repo := TodoRepository{todoTxtPath: TestTodoTxtPath}
	err := repo.Read(TestTodoTxtPath)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("Todos length", func(t *testing.T) {
		length := len(repo.Todos())

		if length != 5 {
			t.Errorf("expected 5 todos, got %d", length)
		}
	})

	t.Run("Done length", func(t *testing.T) {
		length := len(repo.Done())

		if length != 2 {
			t.Errorf("expected 2 done, got %d", length)
		}
	})
}

func TestSave(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test_todo.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := `do some super cool hacker stuff +hacking @basement
x and here is a done todo
(B) here's an existing task +project @shop`

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	TmpTodoTxtPath := tmpfile.Name()

	repo, err := New(TmpTodoTxtPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	repo.Add("(A) here's a new task to be saved @home +testing")

	err = repo.Save()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	savedContent, err := os.ReadFile(TmpTodoTxtPath)
	if err != nil {
		t.Fatalf("failed to read saved todo.txt file: %v", err)
	}

	expectedContent := `do some super cool hacker stuff +hacking @basement
(B) here's an existing task +project @shop
(A) here's a new task to be saved @home +testing
x and here is a done todo
`

	if string(savedContent) != expectedContent {
		t.Errorf("Expected saved content to be: \n%v \n ----------\nGot:\n%v", expectedContent, string(savedContent))
	}
}

func TestFind(t *testing.T) {
	repo, err := New(TestTodoTxtPath)

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
		repo, err := New(TestTodoTxtPath)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		initialDoneLength := repo.DoneCount()
		initialTodosLength := repo.TodoCount()
		repo.Do(1)

		t.Run("item in done", func(t *testing.T) {
			doneLength := repo.DoneCount()
			if doneLength != initialDoneLength+1 {
				t.Errorf("expected %d done items, got %d", initialDoneLength+1, doneLength)
			}
		})

		t.Run("item not in todos", func(t *testing.T) {
			todosLength := repo.TodoCount()
			if todosLength != initialTodosLength-1 {
				t.Errorf("expected %d todos, got %d", initialTodosLength-1, todosLength)
			}
		})
	})

	// t.Run("toggling done item", func(t *testing.T) {
	// 	repo := TodoRepository{todoTxtPath: TestTodoTxtPath}}
	// 	err := repo.Read(TestTodoTxtPath)

	// 	if err != nil {
	// 		t.Fatalf("expected no error, got %v", err)
	// 	}

	// 	initialDoneLength := len(repo.Done())
	// 	if initialDoneLength != 2 {
	// 		t.Errorf("expected 2 done items, got %d", initialDoneLength)
	// 	}
	// 	initialTodosLength := len(repo.Todos())
	// 	if initialTodosLength != 5 {
	// 		t.Errorf("expected 5 todos, got %d", initialTodosLength)
	// 	}
	// 	repo.Toggle(4)

	// 	t.Run("item not in done", func(t *testing.T) {
	// 		doneLength := len(repo.Done())
	// 		if doneLength != initialDoneLength-1 {
	// 			t.Errorf("expected %d done items, got %d", initialDoneLength-1, doneLength)
	// 		}
	// 	})

	// t.Run("item back in todos", func(t *testing.T) {
	// 	todosLength := len(repo.Todos())
	// 	if todosLength != initialTodosLength+1 {
	// 		t.Errorf("expected %d todos, got %d", initialTodosLength+1, todosLength)
	// 	}
	// })
	// })
}

func TestAll(t *testing.T) {
	repo, err := New(TestTodoTxtPath)
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
		if todos[0].Priority != "A" {
			t.Errorf("expected top todo priority to be 'A', got %v", todos[0].Priority)
		}
		if todos[1].Priority != "B" {
			t.Errorf("expected second todo priority to be 'B', got %v", todos[1].Priority)
		}
	})

	t.Run("line numbers are in order", func(t *testing.T) {
		for i, todo := range todos {
			if todo.Number != i+1 {
				t.Errorf("expected todo number to be %d, got %d", i+1, todo.Number)
			}
		}
	})
}

func TestFilter(t *testing.T) {
	repo, err := New(TestTodoTxtPath)
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
