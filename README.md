# togodo
![GitHub branch check runs](https://img.shields.io/github/check-runs/gkarolyi/togodo/master)


A simple TUI task management client based on the [todo.txt format](http://todotxt.org/), created as a final project
for [CS50](https://www.edx.org/learn/computer-science/harvard-university-cs50-s-introduction-to-computer-science).

[Watch demo video](https://github.com/user-attachments/assets/2e6e5453-d7b1-40b8-b696-a2dc375f1d7a)

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
