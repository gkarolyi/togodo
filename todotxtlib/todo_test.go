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

func TestTodo_SetContexts(t *testing.T) {
	tests := []struct {
		name     string
		todo     Todo
		contexts []string
		want     []string
	}{
		{
			name:     "set empty contexts",
			todo:     Todo{Text: "Buy groceries", Contexts: []string{"@home", "@store"}},
			contexts: []string{},
			want:     []string{},
		},
		{
			name:     "set new contexts",
			todo:     Todo{Text: "Buy groceries"},
			contexts: []string{"@home", "@store"},
			want:     []string{"@home", "@store"},
		},
		{
			name:     "replace existing contexts",
			todo:     Todo{Text: "Buy groceries", Contexts: []string{"@old"}},
			contexts: []string{"@new"},
			want:     []string{"@new"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.todo.SetContexts(tt.contexts)
			if len(tt.todo.Contexts) != len(tt.want) {
				t.Errorf("SetContexts() Contexts length = %v, want %v", len(tt.todo.Contexts), len(tt.want))
			} else {
				for i := range tt.todo.Contexts {
					if tt.todo.Contexts[i] != tt.want[i] {
						t.Errorf("SetContexts() Contexts[%d] = %v, want %v", i, tt.todo.Contexts[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestTodo_SetProjects(t *testing.T) {
	tests := []struct {
		name     string
		todo     Todo
		projects []string
		want     []string
	}{
		{
			name:     "set empty projects",
			todo:     Todo{Text: "Buy groceries", Projects: []string{"+shopping", "+food"}},
			projects: []string{},
			want:     []string{},
		},
		{
			name:     "set new projects",
			todo:     Todo{Text: "Buy groceries"},
			projects: []string{"+shopping", "+food"},
			want:     []string{"+shopping", "+food"},
		},
		{
			name:     "replace existing projects",
			todo:     Todo{Text: "Buy groceries", Projects: []string{"+old"}},
			projects: []string{"+new"},
			want:     []string{"+new"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.todo.SetProjects(tt.projects)
			if len(tt.todo.Projects) != len(tt.want) {
				t.Errorf("SetProjects() Projects length = %v, want %v", len(tt.todo.Projects), len(tt.want))
			} else {
				for i := range tt.todo.Projects {
					if tt.todo.Projects[i] != tt.want[i] {
						t.Errorf("SetProjects() Projects[%d] = %v, want %v", i, tt.todo.Projects[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestTodo_AddContext(t *testing.T) {
	tests := []struct {
		name    string
		todo    Todo
		context string
		want    []string
	}{
		{
			name:    "add to empty contexts",
			todo:    Todo{Text: "Buy groceries"},
			context: "@home",
			want:    []string{"@home"},
		},
		{
			name:    "add to existing contexts",
			todo:    Todo{Text: "Buy groceries", Contexts: []string{"@store"}},
			context: "@home",
			want:    []string{"@store", "@home"},
		},
		{
			name:    "do not add duplicate context",
			todo:    Todo{Text: "Buy groceries", Contexts: []string{"@home"}},
			context: "@home",
			want:    []string{"@home"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.todo.AddContext(tt.context)
			if len(tt.todo.Contexts) != len(tt.want) {
				t.Errorf("AddContext() Contexts length = %v, want %v", len(tt.todo.Contexts), len(tt.want))
			} else {
				for i := range tt.todo.Contexts {
					if tt.todo.Contexts[i] != tt.want[i] {
						t.Errorf("AddContext() Contexts[%d] = %v, want %v", i, tt.todo.Contexts[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestTodo_AddProject(t *testing.T) {
	tests := []struct {
		name    string
		todo    Todo
		project string
		want    []string
	}{
		{
			name:    "add to empty projects",
			todo:    Todo{Text: "Buy groceries"},
			project: "+shopping",
			want:    []string{"+shopping"},
		},
		{
			name:    "add to existing projects",
			todo:    Todo{Text: "Buy groceries", Projects: []string{"+food"}},
			project: "+shopping",
			want:    []string{"+food", "+shopping"},
		},
		{
			name:    "do not add duplicate project",
			todo:    Todo{Text: "Buy groceries", Projects: []string{"+shopping"}},
			project: "+shopping",
			want:    []string{"+shopping"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.todo.AddProject(tt.project)
			if len(tt.todo.Projects) != len(tt.want) {
				t.Errorf("AddProject() Projects length = %v, want %v", len(tt.todo.Projects), len(tt.want))
			} else {
				for i := range tt.todo.Projects {
					if tt.todo.Projects[i] != tt.want[i] {
						t.Errorf("AddProject() Projects[%d] = %v, want %v", i, tt.todo.Projects[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestTodo_RemoveContext(t *testing.T) {
	tests := []struct {
		name    string
		todo    Todo
		context string
		want    []string
	}{
		{
			name:    "remove existing context",
			todo:    Todo{Text: "Buy groceries", Contexts: []string{"@home", "@store"}},
			context: "@home",
			want:    []string{"@store"},
		},
		{
			name:    "remove non-existent context",
			todo:    Todo{Text: "Buy groceries", Contexts: []string{"@store"}},
			context: "@home",
			want:    []string{"@store"},
		},
		{
			name:    "remove from empty contexts",
			todo:    Todo{Text: "Buy groceries"},
			context: "@home",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.todo.RemoveContext(tt.context)
			if len(tt.todo.Contexts) != len(tt.want) {
				t.Errorf("RemoveContext() Contexts length = %v, want %v", len(tt.todo.Contexts), len(tt.want))
			} else {
				for i := range tt.todo.Contexts {
					if tt.todo.Contexts[i] != tt.want[i] {
						t.Errorf("RemoveContext() Contexts[%d] = %v, want %v", i, tt.todo.Contexts[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestTodo_RemoveProject(t *testing.T) {
	tests := []struct {
		name    string
		todo    Todo
		project string
		want    []string
	}{
		{
			name:    "remove existing project",
			todo:    Todo{Text: "Buy groceries", Projects: []string{"+shopping", "+food"}},
			project: "+shopping",
			want:    []string{"+food"},
		},
		{
			name:    "remove non-existent project",
			todo:    Todo{Text: "Buy groceries", Projects: []string{"+food"}},
			project: "+shopping",
			want:    []string{"+food"},
		},
		{
			name:    "remove from empty projects",
			todo:    Todo{Text: "Buy groceries"},
			project: "+shopping",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.todo.RemoveProject(tt.project)
			if len(tt.todo.Projects) != len(tt.want) {
				t.Errorf("RemoveProject() Projects length = %v, want %v", len(tt.todo.Projects), len(tt.want))
			} else {
				for i := range tt.todo.Projects {
					if tt.todo.Projects[i] != tt.want[i] {
						t.Errorf("RemoveProject() Projects[%d] = %v, want %v", i, tt.todo.Projects[i], tt.want[i])
					}
				}
			}
		})
	}
}
