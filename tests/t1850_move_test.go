package tests

import (
	"testing"
)

// TestBasicMove tests moving tasks between files
// Ported from: t1850-move.sh
func TestBasicMove(t *testing.T) {
	t.Skip("TODO: Implement move command to move tasks between files")

	// env := SetupTestEnv(t)
	//
	// env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
	// (A) notice the sunflowers`)
	//
	// // Create done.txt with some existing entries
	// // (done.txt creation needs to be added to test helpers)
	//
	// output, code := env.RunCommand("move", "1", "done.txt")
	// // Should move task 1 from todo.txt to done.txt
	// expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
	// TODO: 1 moved from 'todo.txt' to 'done.txt'.`
	//
	// // Verify task removed from todo.txt
	// // Verify task added to done.txt
}

// TestMoveUsage tests move command usage
// Ported from: t1850-move.sh
func TestMoveUsage(t *testing.T) {
	t.Skip("TODO: Implement move command")

	// env := SetupTestEnv(t)
	//
	// output, code := env.RunCommand("move")
	// // Should show usage error
}
