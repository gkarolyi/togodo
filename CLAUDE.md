# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**togodo** is a Terminal User Interface (TUI) task management application written in Go 1.24.4 that implements the standardized todo.txt format. It provides both CLI commands and an interactive TUI interface using the Charm Bracelet ecosystem (Bubble Tea, Bubbles, Lipgloss).

## Development Commands

### Build and Test
```bash
# Build the application
go build

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./cmd
go test ./todotxtlib

# Run a specific test
go test -run TestAddCmd ./cmd

# Static analysis
go vet ./...

# Format code
gofmt -w .
```

### Pre-commit Automation
The project uses lefthook which automatically runs on commits:
- `go vet ./...` - Package consistency checks
- `gofmt -w {files}` - Code formatting
- `go test ./...` - Full test suite
- `go build` - Build verification

### Development Setup
```bash
# Install development tools (requires mise)
mise install

# Run the application
./togodo              # Launch interactive TUI mode
./togodo list         # CLI list command
./togodo add "task"   # CLI add command
```

## Architecture Overview

### Core Structure
- **Entry Point**: `main.go` uses `fang.Execute()` with root command
- **CLI Commands**: `/cmd/` directory with Cobra-based commands
- **Core Logic**: `/todotxtlib/` handles todo.txt parsing and data management
- **TUI Interface**: `/internal/tui/` implements Bubble Tea model/view/controller
- **Formatting**: `/internal/cli/` handles output formatting and theming with Lipgloss
- **Configuration**: `/internal/config/` manages configuration with Viper

### Key Patterns
1. **Dependency Injection**: Dependencies (repository, presenter) are injected through command constructors in `cmd/root.go`
2. **Repository Pattern**: `todotxtlib/repository.go` with pluggable readers/writers (FileReader, FileWriter, BufferWriter)
3. **Command Wrappers**: `cmd/root.go` uses wrapper functions (`wrapAddCmd`, `wrapDoCmd`, etc.) to:
   - Inject dependencies into command constructors
   - Execute command logic (`executeAdd`, `executeDo`, etc.)
   - Handle presentation via the Presenter
4. **Separation of Concerns**:
   - Command structs define CLI interface (flags, usage, validation)
   - Execute functions handle business logic
   - Presenter handles output formatting
   - Repository handles data persistence

### Data Flow
1. **CLI Execution**: `main.go` → config initialization → repository creation → root command construction with dependencies
2. **File Operations**: Abstracted through `todotxtlib/reader.go` and `todotxtlib/writer.go`
3. **Core Data Model**: `Todo` struct in `todotxtlib/todo.go` with priority, projects (+), contexts (@)
4. **Filtering and Sorting**: Separated in `todotxtlib/filter.go` and `todotxtlib/sort.go`
5. **Presentation**: `internal/cli/presenter.go` uses formatter and output writer for consistent display

## Configuration

Configuration is stored in `~/.config/togodo/config.toml` and managed via Viper.

### File Location Strategy
The application searches for `todo.txt` in this order:
1. `--file` flag (if provided)
2. `todo_txt_path` in config file
3. Default: `todo.txt` (in current directory)

Configuration can be managed via:
```bash
togodo config                           # Show all configuration
togodo config todo_txt_path             # Show specific value
togodo config todo_txt_path ~/todo.txt  # Set value
```

## Testing
- Comprehensive test coverage with `*_test.go` files
- Test helpers in `test_helpers.go` files provide shared utilities (e.g., `CreateTempTodoFile`, `NewTestRepository`)
- Both unit and integration tests covering CLI and TUI functionality
- Run tests before committing changes (lefthook enforces this)

## Key Technologies
- **CLI Framework**: Cobra (`github.com/spf13/cobra`)
- **Configuration**: Viper (`github.com/spf13/viper`)
- **TUI Framework**: Bubble Tea (`github.com/charmbracelet/bubbletea`)
- **UI Components**: Charm Bracelet ecosystem (bubbles, lipgloss, fang)
- **Tool Management**: mise for Go version management (Go 1.24.4)

## Available Commands
- **Default (no args)**: Interactive TUI mode
- **list/l/ls [FILTER]**: Display tasks with optional filtering
- **add/a [TASK]**: Add new task(s)
- **do/x [LINE_NUMBER]**: Toggle task completion
- **pri [LINE_NUMBER] [PRIORITY]**: Set task priority (A, B, C)
- **tidy/clean**: Remove completed tasks
- **config [KEY] [VALUE]**: View or set configuration