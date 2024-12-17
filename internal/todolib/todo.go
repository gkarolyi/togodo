package todolib

import "strings"

type Todo struct {
	Text     string
	Done     bool
	Priority string
	Projects []string
	Contexts []string
	Number   int
	Index    int
}

func (t Todo) Prioritised() bool {
	return t.Priority != ""
}

func (t *Todo) ToggleDone() {
	if t.Done {
		t.Done = false
		t.Text = strings.TrimPrefix(t.Text, "x ")
	} else {
		t.Done = true
		t.Text = strings.Join([]string{"x ", t.Text}, "")
	}
}

func (t Todo) Equals(other Todo) bool {
	return t.Text == other.Text
}

// func (t Todo) Project() string {
// 	return projectRe.FindString(t.Text)
// }

// func (t Todo) Context() string {
// 	return contextRe.FindString(t.Text)
// }

// func (t Todo) Priority() string {
// 	return priorityRe.FindString(t.Text)
// }

// func (t Todo) Description() string {
// 	noProj := projectRe.ReplaceAllString(t.Text, "")
// 	noCon := contextRe.ReplaceAllString(noProj, "")
// 	noPri := priorityRe.ReplaceAllString(noCon, "")
// 	return spacingRe.ReplaceAllString(noPri, " ")
// }
