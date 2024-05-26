package todolib

import "testing"

var TestTodoTxtPath = "../../testdata/todo.txt"

func TestAdd(t *testing.T) {
	repo := TodoRepository{}
	todo := repo.Add("(B) random fake task with a +project and @context")

	t.Run("adding to Todos", func(t *testing.T) {
		length := len(repo.Todos)
		want := 1

		if length != want {
			t.Errorf("length %d want %d", length, want)
		}
	})

	t.Run("last Todo matches", func(t *testing.T) {
		got := repo.Todos[len(repo.Todos)-1]
		want := todo

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}

func TestReadFile(t *testing.T) {
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	if err != nil {
		t.Fatal(err)
	}

	t.Run("Todos length", func(t *testing.T) {
		length := len(repo.Todos)
		want := 5

		if length != want {
			t.Errorf("length %d want %d", length, want)
		}
	})
}

func TestFind(t *testing.T) {
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	if err != nil {
		t.Fatal(err)
	}

	t.Run("find correct item", func(t *testing.T) {
		got := repo.Find(1)
		want := Todo{"do some super cool hacker stuff +hacking @basement"}

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}

func TestDo(t *testing.T) {
	repo := TodoRepository{}
	err := repo.ReadFile(TestTodoTxtPath)

	if err != nil {
		t.Fatal(err)
	}

	t.Run("first item in done", func(t *testing.T) {
		repo.Do(1)
		length := len(repo.Done)
		want := 1

		if length != want {
			t.Errorf("length %d want %d", length, want)
		}
	})
}
