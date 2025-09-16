# Gemini Project Context: togodo

## Project Overview

**Name:** togodo
**Type:** Terminal User Interface (TUI) task management application
**Purpose:** A CLI/TUI client for managing tasks in the todo.txt format
**Language:** Go
**License:** MIT License

This project is a terminal-based task management application that uses the `todo.txt` format. It provides both a command-line interface (CLI) for quick actions and a full-screen terminal user interface (TUI) for interactive task management.

## Technology Stack

- **Language:** Go
- **CLI Framework:** Cobra
- **TUI Framework:** Bubble Tea
- **UI Components:** Charm Bracelet ecosystem (bubbles, lipgloss)

## Project Architecture

The project is structured into several packages:

- `main.go`: The application entry point.
- `cmd/`: Contains the implementation of the CLI commands (e.g., `add`, `list`, `do`).
- `todotxtlib/`: The core logic for parsing and managing `todo.txt` files.
- `tui/`: The implementation of the interactive TUI.
- `cli/`: Formatting and presentation logic for the CLI output.

## Building and Running

### Prerequisites

- Go 1.24.4 or later

### Building the Application

To build the application, run the following command in the project root:

```bash
go build
```

This will create a `togodo` executable in the project directory.

### Running the Application

To run the application in TUI mode, execute the following command:

```bash
./togodo
```

You can also use the CLI commands:

```bash
./togodo list
./togodo add "My new task"
./togodo do 1
```

### Running Tests

To run the test suite, use the following command:

```bash
go test -v ./...
```

## Development Conventions

- **Testing:** The project has a strong emphasis on testing. New features should include both unit and integration tests.
- **Atomic Commits:** Commits should be small and atomic, with each commit representing a single, complete feature.
- **Code Style:** Follow the existing code style and architectural patterns.
- **Error Handling:** Use the custom error types defined in `todotxtlib/errors.go`.
