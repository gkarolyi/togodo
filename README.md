# togodo
![GitHub branch check runs](https://img.shields.io/github/check-runs/gkarolyi/togodo/master)


A simple TUI task management client based on the [todo.txt format](http://todotxt.org/), created as a final project
for [CS50](https://www.edx.org/learn/computer-science/harvard-university-cs50-s-introduction-to-computer-science).

[Watch demo video](https://github.com/user-attachments/assets/2e6e5453-d7b1-40b8-b696-a2dc375f1d7a)

## Description

The idea is simple: tasks are written in plain text (in a file named `todo.txt` by default), and each line represents one task.
Tasks can have various attributes, all written in plain text.
A done task is simply denoted with an "x" at the beginning of the line.
Tasks can also have priorities and contexts, belong to projects, have due dates, and any other attribute you want,
as long as they can be represented in plain text:

Example todo.txt for the rest of the usage examples:
```
(A) this is the most urgent task +importantProject @work due:2024-12-31
(B) this is less important but needs to be done @home
this is a task without an assigned priority @work
x this is a finished task
```

The project is structured according to the [Go Docs on module organisation](https://go.dev/doc/modules/layout#packages-and-commands-in-the-same-repository). The installable command line application is located in `cmd`, while the `todolib` library is inside the `internal` directory. Each CLI subcommand is in its own file, where they are registered with `cobra`. The role of the subcommand files is to describe the command and permissible arguments, parse the list of arguments actually passed to the command, and call the appropriate `todolib` functions to make changes to the todo.txt file and display the results.

Inside the `internal/todolib` directory, the `todolib` module is organised in a few different files:
- `todo.go` describes the `Todo` struct which is the datatype I have chosen to represent a single line in the todo.txt file. This struct only has a few simple methods to check whether a todo is prioritised, toggle its done status on/off, and compare two `Todo` structs.
- `matcher.go` contains regular expressions and matchers for parsing the todo.txt format, which are used to determine the priority, context, project, done status, and tags in a particular line in the todo.txt file.
- `presenter.go` contains style definitions and exports a single `Render` function which is responsible for formatting each line (and each word in the line) according to its priority, context, project, and done status.
- `repository.go` is where most of the real 'work' is done. This file defines the `TodoRepository` struct, which is responsible for reading a todo.txt file and parsing it as a list of `Todo` structs. It exports a `New` function which builds a `TodoRepository` from a file path, and several functions for querying and manipulating `Todos` as well as persisting them to the appropriate file.
- `todo_test.go` and `repository_test.go` contain unit tests for most of the exported functions of the `Todo` and `TodoRepository` structs, covering the main functionality of the app to ensure that a change doesn't break existing features.

I chose to experiment with the repository software design pattern for this project in order to try and create a Go library focused on the todo.txt format that could then be used in other projects or slightly different applications. For example, the next feature I would like to start implementing is a TUI powered by the `bubbletea` library accessible at the main `togodo` command. The repository pattern provides an easily graspable layer of abstraction so that the calling application does not need to know the details of the todo.txt specification or deal with saving and parsing correctly. The CLI can simply pass the user's arguments to the repository and trust it to handle the I/O, and a new TUI component could similarly rely on the library and just deal with user input and rendering the interface.

This separation of concerns makes it easier to test the different parts of the app separately, which in turn makes it easier to modify the `todolib` library or the calling command line application in isolation, making the project cleaner and more maintainable. The repository contains business logic such as ordering items by done status and priority by default and handling line number management. Arguably, the presenter portion of `todolib` actually belongs in `cmd`, since not all applications will want to format items the same way - but the separation of concerns in the library means that importing modules are free to only use the `Todo` or the `TodoRepository` and ignore the `Render` function altogether.

Inside the repository, I have chosen to keep `Todos` an array rather than a hash table because line numbers are the natural way to refer to them anyway, so using an ordered list and accessing items by index skips a lot of the complexity of organising them in a hash table by priority or context. Keeping all the items in memory and iterating over them for eg. filtering is very fast, and it is also easy to sort them and write them to a file. At the moment, making any changes to the repository overwrites the whole todo.txt file, which is neat for relatively small files but might become a problem with very large files, or if I add a feature to recursively enumerate all todos in subdirectories. In those cases, it might make sense to keep items in hash tables and make smaller edits to existing files instead.

## Usage/Examples

### `list`

Lists tasks sorted in order of priority, with done items at the bottom of the list. Tasks can optionally be filtered
by passing an optional `[FILTER]` argument. If no filter is passed, `list` shows all items in the list. Tasks are shown
with a line number to allow you to easily refer to them.

```bash
# usage: togodo list [FILTER]
# alias: l, ls
> togodo list "@work"
```
```
1 (A) this is the most urgent task +importantProject @work due:2024-12-31
2 this is a task without an assigned priority @work

```

### `add`

Adds a new task to the list and prints the newly added task. If `[TASK]` contains multiple lines, each line is added as
a separate task to the list.

```bash
# usage: togodo add [TASK]
# alias: a
> togodo add "not a very important task @home
x something I've already done and want to note
(A) a more important task @work I just remembered"
```
```
1 (A) a more important task @work I just remembered
2 not a very important task @home
3 x something I've already done and want to note
```

### `do`

Marks a task as done or not done depending on its current status, and prints the toggled task. If `[LINE_NUMBER]`
contains multiple line numbers, each todo will be toggled.

```bash
# usage: togodo do [LINE_NUMBER]
# alias: x
> togodo do 1 2 3 4
```
```
1 this is a finished task
2 x (A) this is the most urgent task +importantProject @work due:2024-12-31
3 x (B) this is less important but needs to be done @home
4 x this is a task without an assigned priority @work
```

### `tidy`

Cleans up your todo.txt by removing done tasks, and prints the tasks that were removed.

```bash
# usage: togodo tidy
# alias: clean
> togodo tidy
```
```
x this is a finished task
```

## Installation

From the GitHub repo:
```bash
go install github.com/gkarolyi/togodo@latest
```
Alternatively, you can pull the repo and build:
```bash
  gh repo clone gkarolyi/togodo
  cd togodo
  go build
```

## Acknowledgements

This project relies on the [cobra](https://github.com/spf13/cobra) library for handling CLI interactions, and the
[bubbletea](https://github.com/charmbracelet/bubbletea) TUI framework and components for rendering.
