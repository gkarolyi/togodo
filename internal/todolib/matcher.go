package todolib

import "regexp"

var projectRe = regexp.MustCompile(`\+(\w+)`)
var contextRe = regexp.MustCompile(`@\w+`)
var priorityRe = regexp.MustCompile(`^\(([A-Z])\)`)
var doneRe = regexp.MustCompile(`^x `)
var tagRe = regexp.MustCompile(`\w+:\S+`)

func findProjects(text string) []string {
	return projectRe.FindAllString(text, -1)
}

func findContexts(text string) []string {
	return contextRe.FindAllString(text, -1)
}

func findPriority(text string) string {
	if priorityRe.MatchString(text) {
		return priorityRe.FindStringSubmatch(text)[1]
	}
	return ""
}

func findDone(text string) bool {
	return doneRe.MatchString(text)
}

func isDone(text string) bool {
	return doneRe.MatchString(text)
}

func isProject(word string) bool {
	return projectRe.MatchString(word)
}

func isContext(word string) bool {
	return contextRe.MatchString(word)
}

func isPriority(word string) bool {
	return priorityRe.MatchString(word)
}

func isTag(word string) bool {
	return tagRe.MatchString(word)
}
