package todolib

import (
	"os"
	"slices"
	"sort"
	"strconv"
	"testing"

	"golang.org/x/exp/rand"
)

func tempTodoTxtFile(t *testing.T) string {
	tmpfile, err := os.CreateTemp("", "test_todo.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	content := `do some super cool hacker stuff +hacking @basement
buy some carrots and potatoes for +dinner @shop
other random fake task with a +project and @context
(A) don't forget about priorities!!! +GTD @everywhere
(B) finish amazing todo.txt project due:someday
x this is a done todo
x and here is another done todo
`
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(tmpfile.Name())
	})

	return tmpfile.Name()
}

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
	todoTxtPath := tempTodoTxtFile(t)
	repo, err := New(todoTxtPath)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("setting the corrext todoTxtPath", func(t *testing.T) {
		if repo.todoTxtPath != todoTxtPath {
			t.Errorf("expected todoTxtPath to be %s, got %s", todoTxtPath, repo.todoTxtPath)
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
		todos, err := repo.Add("(B) random fake task with a +projectName and @contextName @home")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(todos) != 1 {
			t.Errorf("expected 1 todo, got %d", len(todos))
		}
		todo := todos[0]

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
			todos, err = repo.Add("another todo item to increment number")
			todo = todos[0]
			if todo.Number != 2 {
				t.Errorf("expected line number to be 2, got %d", todo.Number)
			}
			todos, err = repo.Add("(A) this item should be added to the top of the list")
			todo = todos[0]
			if todo.Number != 1 {
				t.Errorf("expected line number to be 1, got %d", todo.Number)
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
		todos, err := repo.Add("x (B) random fake task with a +projectName and @contextName @home")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(todos) != 1 {
			t.Errorf("expected 1 todo, got %d", len(todos))
		}
		todo := todos[0]

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
			if todo.Done != true {
				t.Errorf("expected done to be true, got %v", todo.Done)
			}
		})

		t.Run("line number", func(t *testing.T) {
			if todo.Number != 1 {
				t.Errorf("expected line number to be 1, got %d", todo.Number)
			}
			todos, err = repo.Add("x another done todo item to increment number")
			todo = todos[0]
			if todo.Number != 2 {
				t.Errorf("expected line number to be 2, got %d", todo.Number)
			}
		})
	})
}

func TestRead(t *testing.T) {
	todoTxtPath := tempTodoTxtFile(t)
	repo := TodoRepository{todoTxtPath: todoTxtPath}
	err := repo.Read(todoTxtPath)

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

	expectedContent := `(A) here's a new task to be saved @home +testing
(B) here's an existing task +project @shop
do some super cool hacker stuff +hacking @basement
x and here is a done todo
`

	if string(savedContent) != expectedContent {
		t.Errorf("Expected saved content to be: \n%v \n ----------\nGot:\n%v", expectedContent, string(savedContent))
	}
}

func TestFind(t *testing.T) {
	todoTxtPath := tempTodoTxtFile(t)
	repo, err := New(todoTxtPath)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("find correct item", func(t *testing.T) {
		got := repo.Find(3)
		want := Todo{
			Text: "do some super cool hacker stuff +hacking @basement",
		}

		if got.Text != want.Text {
			t.Errorf("expected %v, got %v", want.Text, got.Text)
		}
	})
}

func TestToggle(t *testing.T) {
	todoTxtPath := tempTodoTxtFile(t)
	repo, err := New(todoTxtPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("toggling not done item", func(t *testing.T) {
		initialDoneLength := repo.DoneCount()
		initialTodosLength := repo.TodoCount()
		todos, err := repo.Toggle([]int{1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		t.Run("item added to done", func(t *testing.T) {
			doneLength := repo.DoneCount()
			if doneLength != initialDoneLength+1 {
				t.Errorf("expected %d done items, got %d", initialDoneLength+1, doneLength)
			}
		})

		t.Run("item removed from todos", func(t *testing.T) {
			todosLength := repo.TodoCount()
			if todosLength != initialTodosLength-1 {
				t.Errorf("expected %d todos, got %d", initialTodosLength-1, todosLength)
			}
		})

		t.Run("item marked as done", func(t *testing.T) {
			if !todos[0].Done {
				t.Errorf("expected todo to be done, got %v", todos[0].Done)
			}
		})

		t.Run("saving change to file", func(t *testing.T) {
			content, err := os.ReadFile(todoTxtPath)
			if err != nil {
				t.Fatalf("failed to read todo.txt file: %v", err)
			}
			expectedContent := `(B) finish amazing todo.txt project due:someday
do some super cool hacker stuff +hacking @basement
buy some carrots and potatoes for +dinner @shop
other random fake task with a +project and @context
x (A) don't forget about priorities!!! +GTD @everywhere
x this is a done todo
x and here is another done todo
`

			if string(content) != expectedContent {
				t.Errorf("expected content to be: \n%s\n -------\n Got: \n%s", expectedContent, string(content))
			}
		})
	})

	t.Run("toggling done item", func(t *testing.T) {
		initialDoneLength := repo.DoneCount()
		initialTodosLength := repo.TodoCount()

		todos, err := repo.Toggle([]int{6})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		t.Run("item not in done", func(t *testing.T) {
			doneLength := repo.DoneCount()
			if doneLength != initialDoneLength-1 {
				t.Errorf("expected %d done items, got %d", initialDoneLength-1, doneLength)
			}
		})

		t.Run("item back in todos", func(t *testing.T) {
			todosLength := repo.TodoCount()
			if todosLength != initialTodosLength+1 {
				t.Errorf("expected %d todos, got %d", initialTodosLength+1, todosLength)
			}
		})

		t.Run("item marked as not done", func(t *testing.T) {
			if todos[0].Done {
				t.Errorf("expected todo to be not done, got %v", todos[0].Done)
			}
		})

		t.Run("saving change to file", func(t *testing.T) {
			content, err := os.ReadFile(todoTxtPath)
			if err != nil {
				t.Fatalf("failed to read todo.txt file: %v", err)
			}
			expectedContent := `(B) finish amazing todo.txt project due:someday
do some super cool hacker stuff +hacking @basement
buy some carrots and potatoes for +dinner @shop
other random fake task with a +project and @context
this is a done todo
x (A) don't forget about priorities!!! +GTD @everywhere
x and here is another done todo
`

			if string(content) != expectedContent {
				t.Errorf("expected content to be: \n%s\n -------\n Got: \n%s", expectedContent, string(content))
			}
		})
	})

	t.Run("toggling non-existent item", func(t *testing.T) {
		_, err := repo.Toggle([]int{100})

		if err == nil {
			t.Errorf("expected an error, got nil")
		}
	})

	t.Run("toggling multiple items", func(t *testing.T) {
		initialDoneLength := repo.DoneCount()
		initialTodosLength := repo.TodoCount()

		todos, err := repo.Toggle([]int{1, 2, 3})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		t.Run("items added to done", func(t *testing.T) {
			doneLength := repo.DoneCount()
			if doneLength != initialDoneLength+3 {
				t.Errorf("expected %d done items, got %d", initialDoneLength+3, doneLength)
			}
		})

		t.Run("items removed from todos", func(t *testing.T) {
			todosLength := repo.TodoCount()
			if todosLength != initialTodosLength-3 {
				t.Errorf("expected %d todos, got %d", initialTodosLength-3, todosLength)
			}
		})

		t.Run("items marked as done", func(t *testing.T) {
			for _, todo := range todos {
				if !todo.Done {
					t.Errorf("expected todo to be done, got %v", todo.Done)
				}
			}
		})
	})
}

func TestAll(t *testing.T) {
	todoTxtPath := tempTodoTxtFile(t)
	repo, err := New(todoTxtPath)
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
	todoTxtPath := tempTodoTxtFile(t)
	repo, err := New(todoTxtPath)
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

func TestTidy(t *testing.T) {
	todoTxtPath := tempTodoTxtFile(t)
	repo, err := New(todoTxtPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("removes done items from the list", func(t *testing.T) {
		removedTodos, err := repo.Tidy()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if repo.DoneCount() != 0 {
			t.Errorf("expected 0 done items, got %d", repo.DoneCount())
		}

		if repo.TodoCount() != 5 {
			t.Errorf("expected 5 todos, got %d", repo.TodoCount())
		}

		for _, todo := range removedTodos {
			if !todo.Done {
				t.Errorf("expected all removed todos to be done, got %v", todo.Text)
			}
		}

		for _, todo := range repo.Items() {
			if todo.Done {
				t.Errorf("expected all done items to be removed, got %v", todo.Text)
			}
		}
	})
}

// Helper function for the old implementation
func sortByPriorityOld(todos []Todo) []Todo {
	sort.SliceStable(todos, func(i, j int) bool {
		if todos[i].Done != todos[j].Done {
			return !todos[i].Done
		}
		if todos[i].Priority != todos[j].Priority {
			if todos[i].Priority == "" {
				return false
			}
			if todos[j].Priority == "" {
				return true
			}
			return todos[i].Priority < todos[j].Priority
		}
		return false
	})
	return todos
}

// Helper function for the new implementation
func sortByPriorityNew(todos []Todo) []Todo {
	slices.SortStableFunc(todos, func(a, b Todo) int {
		if a.Done != b.Done {
			if a.Done {
				return 1
			}
			return -1
		}
		if a.Priority != b.Priority {
			if a.Priority == "" {
				return 1
			}
			if b.Priority == "" {
				return -1
			}
			if a.Priority < b.Priority {
				return -1
			}
			return 1
		}
		return 0
	})
	return todos
}

// Helper function to generate test data
func generateTestTodos(n int) []Todo {
	priorities := []string{"A", "B", "C", "D", "", ""}
	todos := make([]Todo, n)
	for i := 0; i < n; i++ {
		todos[i] = Todo{
			Text:     "Test todo",
			Done:     rand.Float32() < 0.3, // 30% chance of being done
			Priority: priorities[rand.Intn(len(priorities))],
			Number:   i + 1,
		}
	}
	return todos
}

func BenchmarkSortByPriority(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		testData := generateTestTodos(size)

		b.Run("Old-"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Create a fresh copy for each iteration
				todos := make([]Todo, len(testData))
				copy(todos, testData)
				sortByPriorityOld(todos)
			}
		})

		b.Run("New-"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Create a fresh copy for each iteration
				todos := make([]Todo, len(testData))
				copy(todos, testData)
				sortByPriorityNew(todos)
			}
		})
	}
}

// Additional test to verify both implementations sort identically
func TestSortByPriorityEquivalence(t *testing.T) {
	testData := generateTestTodos(1000)

	oldResult := make([]Todo, len(testData))
	copy(oldResult, testData)
	sortByPriorityOld(oldResult)

	newResult := make([]Todo, len(testData))
	copy(newResult, testData)
	sortByPriorityNew(newResult)

	if len(oldResult) != len(newResult) {
		t.Errorf("Different lengths: old=%d, new=%d", len(oldResult), len(newResult))
		return
	}

	for i := range oldResult {
		if !oldResult[i].Equals(newResult[i]) {
			t.Errorf("Mismatch at position %d:\nold: %+v\nnew: %+v", i, oldResult[i], newResult[i])
		}
	}
}
