package todolib

import (
	"testing"

	"github.com/matryer/is"
)

var TestTodoTxtPath = "../../testdata/todo.txt"

func TestAdd(t *testing.T) {
	is := is.New(t)

	t.Run("when todo is not done", func(t *testing.T) {
		repo := TodoRepository{}
		todo := repo.Add("(B) random fake task with a +projectName and @contextName @home")

		t.Run("adding to Todos", func(t *testing.T) {
			length := len(repo.Todos)
			is.Equal(length, 1) // repo should contain one todo
		})

		t.Run("adding to Done", func(t *testing.T) {
			doneLength := len(repo.Done)
			is.Equal(doneLength, 0) // repo should contain no done todos
		})

		t.Run("last Todo matches", func(t *testing.T) {
			got := repo.Todos[len(repo.Todos)-1]
			is.Equal(got, todo) // last todo should be the same todo
		})

		t.Run("priority", func(t *testing.T) {
			is.Equal(todo.Priority, "B") // priority should be B
		})

		t.Run("projects", func(t *testing.T) {
			is.Equal(todo.Projects, []string{"+projectName"}) // projects should be ["+projectName"]
		})

		t.Run("contexts", func(t *testing.T) {
			is.Equal(todo.Contexts, []string{"@contextName", "@home"}) // contexts should be ["@contextName", "@home"]
		})

		t.Run("done status", func(t *testing.T) {
			is.Equal(todo.Done, false) // done should be false
		})

		t.Run("line number", func(t *testing.T) {
			is.Equal(todo.Number, 1) // first todo should have line number 1
			todo = repo.Add("another todo item to increment number")
			is.Equal(todo.Number, 2) // second todo should have line number 2
		})
	})

	t.Run("when todo is done", func(t *testing.T) {
		repo := TodoRepository{}
		todo := repo.Add("x (B) random fake task with a +projectName and @contextName @home")

		t.Run("adding to Todos", func(t *testing.T) {
			todosLength := len(repo.Todos)
			is.Equal(todosLength, 0) // repo should contain no todos
		})

		t.Run("adding to Done", func(t *testing.T) {
			doneLength := len(repo.Done)
			is.Equal(doneLength, 1) // repo should contain one done todo
		})

		t.Run("last Todo matches", func(t *testing.T) {
			got := repo.Done[len(repo.Done)-1]
			is.Equal(got, todo) // last todo should be the same todo
		})

		t.Run("priority", func(t *testing.T) {
			is.Equal(todo.Priority, "") // priority should be empty
		})

		t.Run("projects", func(t *testing.T) {
			is.Equal(todo.Projects, nil) // projects should be nil
		})

		t.Run("contexts", func(t *testing.T) {
			is.Equal(todo.Contexts, nil) // contexts should be nil
		})

		t.Run("done status", func(t *testing.T) {
			is.Equal(todo.Done, true) // done should be true
		})

		t.Run("line number", func(t *testing.T) {
			is.Equal(todo.Number, 1) // first todo should have line number 1
			todo = repo.Add("another todo item to increment number")
			is.Equal(todo.Number, 2) // second todo should have line number 2
		})
	})
}

func TestReadFile(t *testing.T) {
	is := is.New(t)
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	is.NoErr(err)

	t.Run("Todos length", func(t *testing.T) {
		length := len(repo.Todos)

		is.Equal(length, 5) // repo should contain 5 todos
	})

	t.Run("index", func(t *testing.T) {
		lastNumber := repo.Todos[len(repo.Todos)-1].Number
		is.Equal(lastNumber, 7) // index of last todo should be 6
	})

	// t.Run("sort done items at the bottom", func(t *testing.T) {
	// 	lastTodo := repo.Todos[len(repo.Todos)-1]
	// 	is.Equal(lastTodo.Done, true) // last todo status should be done
	// })
}

func TestFind(t *testing.T) {
	is := is.New(t)
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	is.NoErr(err)

	t.Run("find correct item", func(t *testing.T) {
		got := repo.Find(1)
		want := Todo{
			Text: "do some super cool hacker stuff +hacking @basement",
		}

		is.Equal(got.Text, want.Text)
	})
}

func TestDo(t *testing.T) {
	is := is.New(t)

	t.Run("toggling not done item", func(t *testing.T) {
		repo := TodoRepository{}
		err := repo.ReadFile(TestTodoTxtPath)

		is.NoErr(err)

		initialDoneLength := len(repo.Done)
		initialTodosLength := len(repo.Todos)
		repo.Do(1)

		t.Run("item in done", func(t *testing.T) {
			doneLength := len(repo.Done)
			is.Equal(doneLength, initialDoneLength+1) // item should be moved to done
		})

		t.Run("item not in todos", func(t *testing.T) {
			todosLength := len(repo.Todos)
			is.Equal(todosLength, initialTodosLength-1) // item should be removed from todos
		})
	})

	// t.Run("toggling done item", func(t *testing.T) {
	// 	repo := TodoRepository{}
	// 	err := repo.ReadFile(TestTodoTxtPath)

	// 	is.NoErr(err)

	// 	initialDoneLength := len(repo.Done)
	// 	is.Equal(initialDoneLength, 2)
	// 	initialTodosLength := len(repo.Todos)
	// 	is.Equal(initialTodosLength, 5)
	// 	repo.Toggle(4)

	// 	t.Run("item not in done", func(t *testing.T) {
	// 		doneLength := len(repo.Done)
	// 		is.Equal(doneLength, initialDoneLength-1) // item should be removed from done
	// 	})

	// t.Run("item back in todos", func(t *testing.T) {
	// 	todosLength := len(repo.Todos)
	// 	is.Equal(todosLength, initialTodosLength+1) // item should be moved to todos
	// })
	// })
}

func TestList(t *testing.T) {
	is := is.New(t)
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	is.NoErr(err)
	todos := repo.List()

	t.Run("number of items correct", func(t *testing.T) {
		is.Equal(len(todos), 7) // should contain 7 todos
	})

	t.Run("sort done to the bottom", func(t *testing.T) {
		lastTodo := todos[len(todos)-1]
		is.Equal(lastTodo.Text, "x and here is another done todo") // text should match last done todo
	})

	t.Run("sort by priority", func(t *testing.T) {
		topTodo := todos[0]
		is.Equal(topTodo.Priority, "A") // highest priority should be first
	})
}
