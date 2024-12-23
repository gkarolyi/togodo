package main

import (
	demo "github.com/saschagrunert/demo"
	"github.com/urfave/cli/v2"
)

func main() {
	// Create a new demo CLI application
	d := demo.New()
	d.Setup(setup)

	// introduce the project
	d.Add(intro(), "intro", "introduce the project")

	// demo features
	d.Add(add(), "add", "togodo add")
	d.Add(list(), "list", "togodo list")
	d.Add(do(), "do", "togodo do")
	d.Add(tidy(), "tidy", "togodo tidy")

	// Run the application, which registers all signal handlers and waits for
	// the app to exit
	d.Run()
}

// setup is a function that is run before each demo run
func setup(ctx *cli.Context) error {
	// Ensure can be used for easy sequential command execution
	return demo.Ensure(
		"echo '(B) an important task @home' > todo.txt",
		"echo 'buy milk @shop' >> todo.txt",
		"echo 'buy eggs @shop' >> todo.txt",
		"echo 'buy bread @shop' >> todo.txt",
	)
}

// intro is a demo run to introduce the project
func intro() *demo.Run {
	r := demo.NewRun(
		"CS50 Final Project: togodo",
		"Gergely Karolyi",
		"github: @gkarolyi",
		"edx: @gkarolyi",
		"Melbourne, Australia",
	)

	r.Step(demo.S(
		"togodo is a todo.txt manager for the command line written in Go as a final project for CS50.",
		"todo.txt is a text file format for managing todo lists. It's easy to read and write, and it's easy to parse.",
		"The format is simple: each line is a task, and each task can have a priority, a context, and a project.",
		"More about the todo.txt specification here: http://todotxt.org/",
		"togodo will look for a todo.txt file in the current directory, or at the TODO_TXT_PATH environment variable.",
		"Let's see how it works!",
		"We're going to start in a directory that already has a todo.txt file:",
	), demo.S(
		"ls",
	))

	r.Step(demo.S(
		"let's see what's in the todo.txt file, which already contains a few sample tasks:",
	), demo.S(
		"cat todo.txt",
	))

	r.Step(demo.S(
		"Let's see the main features of togodo.",
	), demo.S(
		"clear",
	))

	return r
}

func add() *demo.Run {
	r := demo.NewRun(
		"togodo add [TASK]",
	)

	r.Step(demo.S(
		"Add a new task to the todo.txt file:",
	), demo.S(
		"togodo add 'write a demo for togodo'",
	))

	r.Step(demo.S(
		"Tasks can have a priority, a context, and a project.",
		"Priority is a single uppercase letter in parentheses at the beginning of the task.",
		"This new task has top priority, so it'll be added to the top of the list and highlighted:",
	), demo.S(
		"togodo add '(A) upload the demo video somewhere @home +cs50'",
	))

	r.Step(demo.S(
		"Let's see the updated todo.txt file:",
	), demo.S(
		"togodo list",
	))

	r.Step(demo.S(
		"Apart from priority, context and project, you can add arbitrary metatada such as due dates, tags etc,",
		"using any notation you like. For example: due:today, #learning, or \\go.",
	), demo.S(
		"clear",
	))

	return r
}

func list() *demo.Run {
	r := demo.NewRun(
		"togodo list [OPTIONAL_FILTER]",
	)

	r.Step(demo.S(
		"List all tasks in the todo.txt file:",
	), demo.S(
		"togodo list",
	))

	r.Step(demo.S(
		"List all the things I need to buy:",
	), demo.S(
		"togodo list buy",
	))

	r.Step(demo.S(
		"List everything I need to do at home:",
	), demo.S(
		"togodo list '@home'",
	))

	r.Step(demo.S(
		"",
	), demo.S(
		"clear",
	))

	return r
}

// do is a demo run to show how to mark tasks as done
func do() *demo.Run {
	r := demo.NewRun(
		"togodo do [LINE_NUMBER]",
	)

	r.Step(demo.S(
		"Mark a task as done:",
	), demo.S(
		"togodo do 1",
	))

	r.Step(demo.S(
		"Mark multiple tasks as done:",
	), demo.S(
		"togodo do 1 2",
	))

	r.Step(demo.S(
		"The todo.txt file is always sorted by priority, with done tasks at the bottom.",
	), demo.S(
		"togodo list",
	))

	r.Step(demo.S(
		"Mark some tasks as done and some as not done in one command:",
	), demo.S(
		"togodo do 1 3 4",
	))

	r.Step(demo.S(
		"Let's see what's left",
	), demo.S(
		"togodo list",
	))

	r.Step(demo.S(
		"In future, done tasks will be moved to a done.txt file instead of deleted.",
	), demo.S(
		"clear",
	))

	return r
}

// tidy is a demo run to show how to tidy up the todo.txt file
func tidy() *demo.Run {
	r := demo.NewRun(
		"togodo tidy",
	)

	r.Step(demo.S(
		"Let's add some completed tasks to the todo.txt file:",
	), demo.S(
		"togodo add 'x this is a completed task\nx this is another completed task\nx and another task I already completed'",
	))

	r.Step(demo.S(
		"This is what the todo.txt file looks like now:",
	), demo.S(
		"togodo list",
	))

	r.Step(demo.S(
		"Let's clean up the todo.txt file a bit:",
	), demo.S(
		"togodo tidy",
	))

	r.Step(demo.S(
		"This removes deleted tasks and reorders the remaining ones by priority.",
	), demo.S(
		"togodo list",
	))

	r.Step(demo.S(
		"That's it for now! These are the currently implemented features of togodo, but there's more to come.",
		"Thank you for watching!",
	), demo.S(
		"clear",
	))

	return r
}
