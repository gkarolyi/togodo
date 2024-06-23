package todolib

import "regexp"

var projectRe = regexp.MustCompile(`\+(\w+)`)
var contextRe = regexp.MustCompile(`@\w+`)
var priorityRe = regexp.MustCompile(`^\(([A-Z])\)`)
var doneRe = regexp.MustCompile(`^x `)
var tagRe = regexp.MustCompile(`\w+:\S+`)

// var spacingRe = regexp.MustCompile(`\s{2,}`)

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
