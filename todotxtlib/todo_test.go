package todotxtlib

import (
	"testing"
)

func TestNewTodo(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    Todo
		wantErr bool
	}{
		{
			name: "simple todo",
			text: "Buy groceries",
			want: Todo{
				Text:     "Buy groceries",
				Done:     false,
				Priority: "",
				Projects: []string{},
				Contexts: []string{},
			},
		},
		{
			name: "done todo",
			text: "x Buy groceries",
			want: Todo{
				Text:     "x Buy groceries",
				Done:     true,
				Priority: "",
				Projects: []string{},
				Contexts: []string{},
			},
		},
		{
			name: "with priority",
			text: "(A) Buy groceries",
			want: Todo{
				Text:     "(A) Buy groceries",
				Done:     false,
				Priority: "A",
				Projects: []string{},
				Contexts: []string{},
			},
		},
		{
			name: "with projects",
			text: "Buy groceries +shopping +food",
			want: Todo{
				Text:     "Buy groceries +shopping +food",
				Done:     false,
				Priority: "",
				Projects: []string{"+shopping", "+food"},
				Contexts: []string{},
			},
		},
		{
			name: "with contexts",
			text: "Buy groceries @home @store",
			want: Todo{
				Text:     "Buy groceries @home @store",
				Done:     false,
				Priority: "",
				Projects: []string{},
				Contexts: []string{"@home", "@store"},
			},
		},
		{
			name: "complete todo",
			text: "x (A) Buy groceries +shopping @store",
			want: Todo{
				Text:     "x (A) Buy groceries +shopping @store",
				Done:     true,
				Priority: "",
				Projects: []string{"+shopping"},
				Contexts: []string{"@store"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTodo(tt.text)
			if got.Text != tt.want.Text {
				t.Errorf("NewTodo() Text = %v, want %v", got.Text, tt.want.Text)
			}
			if got.Done != tt.want.Done {
				t.Errorf("NewTodo() Done = %v, want %v", got.Done, tt.want.Done)
			}
			if got.Priority != tt.want.Priority {
				t.Errorf("NewTodo() Priority = %v, want %v", got.Priority, tt.want.Priority)
			}
			if len(got.Projects) != len(tt.want.Projects) {
				t.Errorf("NewTodo() Projects length = %v, want %v", len(got.Projects), len(tt.want.Projects))
			} else {
				for i := range got.Projects {
					if got.Projects[i] != tt.want.Projects[i] {
						t.Errorf("NewTodo() Projects[%d] = %v, want %v", i, got.Projects[i], tt.want.Projects[i])
					}
				}
			}
			if len(got.Contexts) != len(tt.want.Contexts) {
				t.Errorf("NewTodo() Contexts length = %v, want %v", len(got.Contexts), len(tt.want.Contexts))
			} else {
				for i := range got.Contexts {
					if got.Contexts[i] != tt.want.Contexts[i] {
						t.Errorf("NewTodo() Contexts[%d] = %v, want %v", i, got.Contexts[i], tt.want.Contexts[i])
					}
				}
			}
		})
	}
}

func TestTodo_Prioritised(t *testing.T) {
	tests := []struct {
		name string
		todo Todo
		want bool
	}{
		{
			name: "no priority",
			todo: Todo{Text: "Buy groceries"},
			want: false,
		},
		{
			name: "with priority",
			todo: Todo{Text: "(A) Buy groceries", Priority: "A"},
			want: true,
		},
		{
			name: "completed task",
			todo: Todo{Text: "x (A) Buy groceries", Done: true},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.todo.Prioritised(); got != tt.want {
				t.Errorf("Prioritised() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodo_ToggleDone(t *testing.T) {
	tests := []struct {
		name string
		todo Todo
		want string
	}{
		{
			name: "not done to done",
			todo: Todo{Text: "Buy groceries"},
			want: "x Buy groceries",
		},
		{
			name: "done to not done",
			todo: Todo{Text: "x Buy groceries", Done: true},
			want: "Buy groceries",
		},
		{
			name: "with priority",
			todo: Todo{Text: "(A) Buy groceries", Priority: "A"},
			want: "x (A) Buy groceries",
		},
		{
			name: "completed task",
			todo: Todo{Text: "x (A) Buy groceries", Done: true},
			want: "(A) Buy groceries",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.todo.ToggleDone()
			if tt.todo.Text != tt.want {
				t.Errorf("ToggleDone() Text = %v, want %v", tt.todo.Text, tt.want)
			}
		})
	}
}

func TestTodo_SetPriority(t *testing.T) {
	tests := []struct {
		name     string
		todo     Todo
		priority string
		want     string
	}{
		{
			name:     "set priority",
			todo:     Todo{Text: "Buy groceries"},
			priority: "A",
			want:     "(A) Buy groceries",
		},
		{
			name:     "change priority",
			todo:     Todo{Text: "(B) Buy groceries", Priority: "B"},
			priority: "A",
			want:     "(A) Buy groceries",
		},
		{
			name:     "remove priority",
			todo:     Todo{Text: "(A) Buy groceries", Priority: "A"},
			priority: "",
			want:     "Buy groceries",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.todo.SetPriority(tt.priority)
			if tt.todo.Text != tt.want {
				t.Errorf("SetPriority() Text = %v, want %v", tt.todo.Text, tt.want)
			}
			if tt.todo.Priority != tt.priority {
				t.Errorf("SetPriority() Priority = %v, want %v", tt.todo.Priority, tt.priority)
			}
		})
	}
}

func TestTodo_Equals(t *testing.T) {
	tests := []struct {
		name  string
		todo1 Todo
		todo2 Todo
		want  bool
	}{
		{
			name:  "equal todos",
			todo1: Todo{Text: "Buy groceries"},
			todo2: Todo{Text: "Buy groceries"},
			want:  true,
		},
		{
			name:  "different text",
			todo1: Todo{Text: "Buy groceries"},
			todo2: Todo{Text: "Buy milk"},
			want:  false,
		},
		{
			name:  "different done status",
			todo1: Todo{Text: "Buy groceries", Done: true},
			todo2: Todo{Text: "Buy groceries", Done: false},
			want:  false,
		},
		{
			name:  "different priority",
			todo1: Todo{Text: "(A) Buy groceries", Priority: "A"},
			todo2: Todo{Text: "(B) Buy groceries", Priority: "B"},
			want:  false,
		},
		{
			name:  "different projects",
			todo1: Todo{Text: "Buy groceries +shopping", Projects: []string{"+shopping"}},
			todo2: Todo{Text: "Buy groceries +food", Projects: []string{"+food"}},
			want:  false,
		},
		{
			name:  "different contexts",
			todo1: Todo{Text: "Buy groceries @home", Contexts: []string{"@home"}},
			todo2: Todo{Text: "Buy groceries @work", Contexts: []string{"@work"}},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.todo1.Equals(tt.todo2); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
