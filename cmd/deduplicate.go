package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// DeduplicateResult contains the result of a Deduplicate operation
type DeduplicateResult struct {
	RemovedCount int
}

// Deduplicate removes duplicate tasks based on their text content (ignoring priority).
// When duplicates exist, keeps the task with the highest priority (A > B > C > ... > no priority)
func Deduplicate(repo todotxtlib.TodoRepository) (DeduplicateResult, error) {
	// Get all todos
	allTodos, err := repo.ListAll()
	if err != nil {
		return DeduplicateResult{}, fmt.Errorf("failed to list all todos: %w", err)
	}

	// Group todos by their base text (without priority)
	// Map: base text -> list of indices with that text
	groups := make(map[string][]int)

	for i, todo := range allTodos {
		baseText := getBaseText(todo.Text)
		groups[baseText] = append(groups[baseText], i)
	}

	// Find which todos to keep (those with highest priority in each group)
	toKeep := make(map[int]bool)

	for _, indices := range groups {
		if len(indices) == 1 {
			// No duplicates for this text
			toKeep[indices[0]] = true
			continue
		}

		// Find the todo with highest priority
		bestIdx := indices[0]
		bestPriority := allTodos[indices[0]].Priority

		for _, idx := range indices[1:] {
			priority := allTodos[idx].Priority
			if comparePriority(priority, bestPriority) < 0 {
				// This priority is better (higher)
				bestIdx = idx
				bestPriority = priority
			}
		}

		toKeep[bestIdx] = true
	}

	// Remove todos that are not in toKeep (process backwards to avoid index shifting)
	removedCount := 0
	for i := len(allTodos) - 1; i >= 0; i-- {
		if !toKeep[i] {
			if _, err := repo.Remove(i); err != nil {
				return DeduplicateResult{}, fmt.Errorf("failed to remove todo at index %d: %w", i, err)
			}
			removedCount++
		}
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return DeduplicateResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DeduplicateResult{RemovedCount: removedCount}, nil
}

// getBaseText extracts the text content without priority prefix or completion marker
func getBaseText(text string) string {
	// Remove "x " prefix if present (completed tasks)
	text = strings.TrimPrefix(text, "x ")
	text = strings.TrimSpace(text)

	// Remove priority prefix if present
	priorityRe := regexp.MustCompile(`^\([A-Z]\) `)
	text = priorityRe.ReplaceAllString(text, "")

	return text
}

// comparePriority compares two priorities, returning:
// - negative if p1 has higher priority than p2
// - positive if p2 has higher priority than p1
// - zero if they are equal
// Priority order: A > B > C > ... > Z > (no priority)
func comparePriority(p1, p2 string) int {
	// Empty priority is lowest
	if p1 == "" && p2 == "" {
		return 0
	}
	if p1 == "" {
		return 1 // p2 is higher
	}
	if p2 == "" {
		return -1 // p1 is higher
	}

	// Compare alphabetically (A < B < C, etc.)
	// Lower letter = higher priority, so we return negative if p1 < p2
	if p1 < p2 {
		return -1
	}
	if p1 > p2 {
		return 1
	}
	return 0
}
