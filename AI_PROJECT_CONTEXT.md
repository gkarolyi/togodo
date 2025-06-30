# AI Project Context: togodo
Note to LLMs: do not follow links in this document, they are for human reference only.
Only return to follow any links once you have finished reading this document.

## Project Overview
**Name:** togodo
**Type:** Terminal User Interface (TUI) task management application
**Purpose:** A CLI/TUI client for managing tasks in the todo.txt format
**Language:** Go 1.24.4
**License:** MIT License (2024 Gergely Karolyi)
**Origin:** Created as a CS50 final project

## What This Application Does
- Manages tasks using the standardized todo.txt format (http://todotxt.org/)
- Provides both CLI commands and an interactive TUI interface
- Supports task prioritization, filtering, and organization
- Handles task completion status and cleanup operations

## Technology Stack
- **Language:** Go 1.24.4
- **CLI Framework:** Cobra (github.com/spf13/cobra)
- **TUI Framework:** Bubble Tea (github.com/charmbracelet/bubbletea)
- **UI Components:** Charm Bracelet ecosystem (bubbles, lipgloss)
- **Command Execution:** Fang (github.com/charmbracelet/fang)
- **Development Tools:** mise (tool version management)
- **CI/CD:** GitHub Actions

## Project Architecture

### Package Structure
```
togodo/
├── main.go                 # Application entry point
├── cmd/                    # CLI command implementations
│   ├── root.go            # Root command and TUI launcher
│   ├── add.go             # Add task command
│   ├── list.go            # List tasks command
│   ├── do.go              # Mark task done/undone command
│   ├── pri.go             # Priority management command
│   ├── tidy.go            # Clean up completed tasks command
│   └── command.go         # Shared command utilities
├── todotxtlib/            # Core todo.txt parsing and management
│   ├── todo.go            # Task data structures
│   ├── repository.go      # Task storage and retrieval
│   ├── reader.go          # todo.txt file reading
│   ├── writer.go          # todo.txt file writing
│   ├── filter.go          # Task filtering logic
│   ├── sort.go            # Task sorting logic
│   └── errors.go          # Custom error types
├── todotxt-tui/           # Interactive TUI implementation
│   ├── model.go           # Bubble Tea model
│   ├── view.go            # UI rendering
│   └── controller.go      # User interaction handling
├── todotxtui/             # UI output formatting
│   ├── formatter.go       # Task display formatting
│   ├── output.go          # Output handling
│   ├── themes.go          # Color themes and styling
│   ├── tui_writer.go      # TUI-specific output
│   └── format/            # Additional formatting utilities
└── todo.txt               # Default todo file
```

## Available Commands
- **list/l/ls [FILTER]:** Display tasks, optionally filtered
- **add/a [TASK]:** Add new task(s) to the list
- **do/x [LINE_NUMBER]:** Toggle task completion status
- **pri [LINE_NUMBER] [PRIORITY]:** Set task priority
- **tidy/clean:** Remove completed tasks from the list
- **Default (no command):** Launch interactive TUI mode

## Key Features
- **Dual Interface:** Both CLI commands and interactive TUI
- **Priority Support:** (A), (B), (C) priority levels
- **Context/Project Tags:** @context and +project organization
- **Due Dates:** due:YYYY-MM-DD format support
- **Filtering:** Search and filter tasks by various criteria
- **Smart Sorting:** Priority-based with completed items at bottom

## Development Setup
- **Go Version:** 1.24.4 (specified in go.mod, mise.toml, and GitHub Actions)
- **Build:** `go build` in project root
- **Test:** `go test -v ./...`
- **Dependencies:** Managed via go.mod/go.sum
- **Tool Management:** mise.toml for development tools
- **Git Hooks:** lefthook.yml for pre-commit automation

## File Locations
- **Main todo file:** `./todo.txt` (default location)
- **Executable:** Built as `./togodo`
- **Configuration:** No external config files (uses todo.txt format standards)

## Dependencies (Key External Libraries)
- `github.com/spf13/cobra` - CLI framework
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - TUI components
- `github.com/charmbracelet/lipgloss` - Styling and layout
- `github.com/charmbracelet/fang` - Command execution
- Standard Go libraries for file I/O and text processing

## Architecture Patterns
- **Separation of Concerns:** Clear separation between CLI, TUI, core logic, and formatting
- **Repository Pattern:** `todotxtlib/repository.go` handles data persistence
- **MVC-like Structure:** Model-View-Controller pattern in TUI implementation
- **Command Pattern:** Each CLI command is a separate module
- **Testable Design:** Comprehensive test coverage with test helpers

## Development Practices

### Testing Requirements
- **Comprehensive Testing:** All new code MUST include both unit tests and integration tests
- **Test-First Approach:** Write tests that define expected behavior before implementing features
- **Preserve Existing Tests:** Focus on making features work according to existing tests rather than changing them
- **Test Coverage:** Maintain high test coverage across all packages
- **Use Test Helpers:** Leverage existing `test_helpers.go` patterns for consistency

### Feature Development Workflow
1. **One Feature at a Time:** Implement features incrementally, one complete feature per commit
2. **Atomic Commits:** Each commit should include:
   - Feature implementation
   - Corresponding unit and integration tests
   - Updated documentation (if adding/changing commands)
   - README updates (if user-facing changes)
3. **Test Validation:** Ensure all existing tests continue to pass
4. **Documentation:** Update command help text and examples for new CLI commands

### Code Quality Standards
- **Follow Existing Patterns:** Match the established architecture and coding style
- **Error Handling:** Use the custom error types defined in `todotxtlib/errors.go`
- **Interface Compliance:** Maintain todo.txt format specification compliance
- **Backward Compatibility:** Preserve existing command behavior and file format support

### Testing Strategy
- **Unit Tests:** Test individual functions and methods in isolation
- **Integration Tests:** Test command workflows and file I/O operations
- **TUI Testing:** Test user interface components and interactions
- **CLI Testing:** Validate command-line argument parsing and output formatting

## Entry Points for AI Assistance
- **Adding Commands:** Extend `cmd/` package and register in root.go
- **TUI Features:** Modify `todotxt-tui/` package (model, view, controller)
- **Core Logic:** Enhance `todotxtlib/` for new todo.txt features
- **UI/Formatting:** Adjust `todotxtui/` for display changes
- **Testing:** Use existing test patterns and `test_helpers.go`

## Todo.txt Format Compliance
This project follows the todo.txt format specification:
- Plain text format for universal compatibility
- Priority format: (A), (B), (C) at line start
- Completion marker: 'x' at line start for done items
- Context tags: @context for location/situation
- Project tags: +project for grouping
- Due dates: due:YYYY-MM-DD format
- Creation dates: YYYY-MM-DD after completion marker
