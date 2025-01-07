package todolib

import "regexp"

var projectRe = regexp.MustCompile(`\+(\w+)`)
var contextRe = regexp.MustCompile(`@\w+`)
var priorityRe = regexp.MustCompile(`^(?:x )?(\(([A-Z])\))`)
var doneRe = regexp.MustCompile(`^x `)
var tagRe = regexp.MustCompile(`\w+:\S+`)

func FindProjects(text string) []string {
	return projectRe.FindAllString(text, -1)
}

func FindContexts(text string) []string {
	return contextRe.FindAllString(text, -1)
}

func FindPriority(text string) string {
	if priorityRe.MatchString(text) {
		return priorityRe.FindStringSubmatch(text)[2]
	}
	return ""
}

func FindDone(text string) bool {
	return doneRe.MatchString(text)
}

func IsDone(text string) bool {
	return doneRe.MatchString(text)
}

func IsProject(word string) bool {
	return projectRe.MatchString(word)
}

func IsContext(word string) bool {
	return contextRe.MatchString(word)
}

func IsPriority(word string) bool {
	return priorityRe.MatchString(word)
}

func IsTag(word string) bool {
	return tagRe.MatchString(word)
}
