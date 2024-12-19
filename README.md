# togodo
![GitHub branch status](https://img.shields.io/github/checks-status/gkarolyi/togodo/master)


A simple TUI task management client based on the [todo.txt format](http://todotxt.org/).

The idea is simple: tasks are written in plain text (in a file named `todo.txt` by default), and each line represents one task.
Tasks can have various attributes, all written in plain text.

## Usage/Examples

```bash
❯ ./togodo list
1 (A) test the saving setup to make sure content isn't lost +test
10 (A) change TodoRepository API to only export New() and automatically read the right file +refactor
3 (B) make saving work and save on add, complete etc +feature
4 (B) add Add command to CLI +feature
5 (C) add Do command to CLI +feature
6 add tui for browsing and acting on todos +feature
7 add quick entry interface? +idea
8 support setting priorities with !! shortcuts +idea
9 create TodoRepository with a todo.txt path +refactor
11 implement recursive flag to list todo.txts from subdirectories +idea
12 might be faster to just use text operations rather than classes +idea
2 x (A) list todos from current directory by default
```
## Installation

Clone the project

```bash
  gh repo clone gkarolyi/togodo
```

Go to the project directory

```bash
  cd togodo
```

Build the executable

```bash
  go build
```

Alternatively, using go:

```bash
go install github.com/gkarolyi/togodo@latest
```
