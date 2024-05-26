package todolib

import "regexp"

type Todo struct {
	Text string
}

var projectRe = regexp.MustCompile(`\+\w+`)
var contextRe = regexp.MustCompile(`@\w+`)
var priorityRe = regexp.MustCompile(`\([A-Z]\)`)

func (t Todo) Project() string {
	return projectRe.FindString(t.Text)
}

func (t Todo) Context() string {
	return contextRe.FindString(t.Text)
}

func (t Todo) Priority() string {
	return priorityRe.FindString(t.Text)
}
