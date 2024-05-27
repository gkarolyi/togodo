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
			length := len(repo.Todos)
			is.Equal(length, 1) // repo should contain one todo
		})

		t.Run("last Todo matches", func(t *testing.T) {
			got := repo.Todos[len(repo.Todos)-1]
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

		is.Equal(length, 5)
	})

	t.Run("line numbers", func(t *testing.T) {
		lastNumber := repo.Todos[len(repo.Todos)-1].Number
		is.Equal(lastNumber, 5)
	})
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
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	is.NoErr(err)

	t.Run("first item in done", func(t *testing.T) {
		repo.Do(1)
		length := len(repo.Done)

		is.Equal(length, 1)
	})
}
