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
- **TUI Interface**: `/todotxt-tui/` implements Bubble Tea model/view/controller
- **Formatting**: `/todotxtui/` handles output formatting and theming with Lipgloss

### Key Patterns
1. **Repository Pattern**: `todotxtlib/repository.go` with pluggable readers/writers
2. **BaseCommand**: Shared functionality between CLI commands via `cmd/command.go`
3. **Dual Command Types**:
   - `NewDefaultBaseCommand()`: CLI commands (stdout output)
   - `NewTUIBaseCommand()`: Interactive mode (TUI output)

### Data Flow
- File operations abstracted through `todotxtlib/reader.go` and `todotxtlib/writer.go`
- Core `Todo` struct in `todotxtlib/todo.go` with priority, projects (+), contexts (@)
- Filtering and sorting logic separated in `todotxtlib/filter.go` and `todotxtlib/sort.go`

## File Location Strategy
The application searches for `todo.txt` in this order:
1. `TODO_TXT_PATH` environment variable
2. Current working directory  
3. User's home directory

## Testing
- Comprehensive test coverage with `*_test.go` files
- Test helpers in `test_helpers.go` files provide shared utilities
- Both unit and integration tests covering CLI and TUI functionality
- Always run the full test suite before committing changes

## Key Technologies
- **CLI Framework**: Cobra (`github.com/spf13/cobra`)
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