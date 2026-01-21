# Phase 2 Work Breakdown

**Date**: 2026-01-21
**Source**: Integration test execution results
**Total Phase 2 Items**: 16 skipped tests + 5 failing tests requiring new commands

## Overview

This document groups all remaining work by feature type and provides implementation guidance for Phase 2.

## Feature Groups

### Group 1: Archive & Done.txt Management (5 failing tests)

**Status**: Critical - Required for todo.txt-cli parity
**Complexity**: Medium
**Estimated Effort**: 12-16 hours

#### Tests Requiring Implementation
- `TestArchiveWithDuplicates` (4 sub-tests)
  - archive_done_tasks
  - verify_done_task_removed_from_todo.txt
  - verify_done_task_in_done.txt
  - list_after_archive
- `TestArchiveWarning` (1 sub-test)
  - archive_with_no_done_tasks

#### Implementation Requirements

**New Files Needed**:
- `cmd/archive.go` - Business logic for archiving
- `cmd/archive_test.go` - Unit tests
- `internal/cli/archive.go` - CLI command wrapper

**Core Functionality**:
1. Read todo.txt and identify completed tasks (marked with 'x')
2. Append completed tasks to done.txt (create if doesn't exist)
3. Remove completed tasks from todo.txt
4. Save both files
5. Handle edge cases:
   - No completed tasks (show warning)
   - done.txt doesn't exist (create it)
   - Duplicate tasks in done.txt (allow duplicates per spec)

**Expected Behavior**:
```bash
# Archive completed tasks
togodo archive
# Output: TODO: todo.txt archived.

# If no completed tasks
togodo archive
# Output: TODO: todo.txt does not contain any done tasks.
```

**Dependencies**:
- `todotxtlib/repository.go` - May need `ListDone()` and `ArchiveDone()` methods
- Config system - Need to handle `done_txt_path` configuration

**Testing Approach**:
1. Unit test `cmd.Archive()` function
2. Integration test with buffer-based repository
3. Verify file operations (both todo.txt and done.txt)
4. Test with various edge cases (no done tasks, empty files, duplicates)

---

### Group 2: Configuration Management (3 skipped tests)

**Status**: Important - Needed for usability
**Complexity**: Low-Medium
**Estimated Effort**: 4-6 hours

#### Tests Requiring Implementation
- `TestConfigRead` - Read configuration values
- `TestConfigWrite` - Write configuration values
- `TestConfigList` - List all configuration settings

#### Implementation Requirements

**New Files Needed**:
- `cmd/config.go` - Business logic for config operations
- `cmd/config_test.go` - Unit tests
- `internal/cli/config.go` - CLI command wrapper (may already exist)

**Core Functionality**:
1. Read configuration value by key
2. Write configuration value by key
3. List all configuration settings
4. Validate configuration keys and values

**Expected Behavior**:
```bash
# Read config value
togodo config todo_txt_path
# Output: /path/to/todo.txt

# Write config value
togodo config todo_txt_path ~/newtodo.txt
# Output: Configuration updated

# List all config
togodo config
# Output:
# todo_txt_path: /path/to/todo.txt
# done_txt_path: /path/to/done.txt
```

**Dependencies**:
- Viper configuration system (already in use)
- Config file at `~/.config/togodo/config.toml`

**Testing Approach**:
1. Unit test config read/write/list operations
2. Integration test with temporary config files
3. Test invalid keys and values

**Notes**:
- This may already be partially implemented - check `internal/cli/config.go`
- If exists, may just need output format adjustments

---

### Group 3: Data Quality Commands (2 tests)

**Status**: Nice-to-have - Quality of life improvements
**Complexity**: Medium
**Estimated Effort**: 8-10 hours

#### Tests Requiring Implementation
- `TestDeduplicate` (1 failing test) - Remove duplicate tasks
- `TestDeduplicateWithPriority` (1 skipped test) - Preserve higher priority

#### Implementation Requirements

**New Files Needed**:
- `cmd/deduplicate.go` - Business logic for deduplication
- `cmd/deduplicate_test.go` - Unit tests
- `internal/cli/deduplicate.go` - CLI command wrapper

**Core Functionality**:
1. Scan all tasks and identify duplicates (exact text match)
2. When duplicates found:
   - If priorities differ, keep higher priority
   - If priorities same or absent, keep first occurrence
3. Remove duplicate tasks
4. Save updated todo.txt

**Expected Behavior**:
```bash
# Remove duplicates
togodo deduplicate
# Output: TODO: 3 duplicate tasks removed.

# If no duplicates
togodo deduplicate
# Output: TODO: No duplicates found.
```

**Algorithm**:
```
1. Load all tasks
2. Create map: text -> list of tasks with that text
3. For each duplicate group:
   a. If priorities differ: keep highest priority
   b. If priorities same: keep first occurrence
   c. Mark others for deletion
4. Remove marked tasks
5. Save
```

**Testing Approach**:
1. Unit test deduplication logic
2. Test with same text, no priorities
3. Test with same text, different priorities
4. Test with no duplicates
5. Test with multiple sets of duplicates

---

### Group 4: Reporting & Statistics (1 failing test)

**Status**: Nice-to-have - Analytics feature
**Complexity**: Low-Medium
**Estimated Effort**: 4-6 hours

#### Tests Requiring Implementation
- `TestReport` - Show task statistics

#### Implementation Requirements

**New Files Needed**:
- `cmd/report.go` - Business logic for reporting
- `cmd/report_test.go` - Unit tests
- `internal/cli/report.go` - CLI command wrapper

**Core Functionality**:
1. Count total tasks
2. Count completed tasks
3. Count by priority (A, B, C, etc.)
4. Count by context and project
5. Calculate completion rate

**Expected Behavior**:
```bash
# Show report
togodo report
# Output:
# Total tasks: 25
# Completed: 10 (40%)
# Priority A: 5
# Priority B: 8
# Priority C: 2
# Contexts: @home (12), @work (8), @errands (5)
# Projects: +website (15), +blog (7), +personal (3)
```

**Testing Approach**:
1. Unit test report generation
2. Test with various task combinations
3. Test with empty todo.txt
4. Verify statistics accuracy

---

### Group 5: Date Support (2 skipped tests)

**Status**: Important - Standard todo.txt format feature
**Complexity**: Medium-High
**Estimated Effort**: 12-16 hours

#### Tests Requiring Implementation
- `TestAddWithDate` - Adding tasks with creation dates
- `TestAddWithPriority` - Adding tasks with priority in input

#### Implementation Requirements

**Modifications Needed**:
- `todotxtlib/todo.go` - Parse dates from todo text
- `todotxtlib/parser.go` - Date parsing logic (may need to create)
- `cmd/add.go` - Support date flags or automatic date addition
- Tests for date parsing and formatting

**Core Functionality**:
1. Parse creation dates from todo text (format: YYYY-MM-DD)
2. Parse completion dates (format: x YYYY-MM-DD YYYY-MM-DD)
3. Optionally add creation date when adding new tasks
4. Support priority with date (format: (A) YYYY-MM-DD task text)

**Expected Behavior**:
```bash
# Add with auto-date
togodo add --with-date "new task"
# Saved: 2026-01-21 new task

# Add with priority and date
togodo add "(A) important task"
# Saved: (A) 2026-01-21 important task
```

**Todo.txt Date Format**:
```
# Creation date
2026-01-21 task text

# Priority with creation date
(A) 2026-01-21 task text

# Completed with dates
x 2026-01-21 2026-01-15 task text
  ^completion ^creation
```

**Testing Approach**:
1. Unit test date parsing
2. Test various date formats
3. Test priority + date combinations
4. Test completion date handling
5. Integration tests for add with dates

**Complexity Notes**:
- Requires careful parsing logic
- Must maintain todo.txt format compliance
- Need to handle timezones and date validation
- May affect multiple commands (add, do, list display)

---

### Group 6: Enhanced Delete Operations (2 skipped tests)

**Status**: Nice-to-have - Power user features
**Complexity**: Medium
**Estimated Effort**: 6-8 hours

#### Tests Requiring Implementation
- `TestDelMultiple` - Delete multiple tasks at once
- `TestDelWithTerm` - Remove specific terms from tasks

#### Implementation Requirements

**Modifications Needed**:
- `cmd/del.go` - Support multiple line numbers
- `cmd/del.go` - Support term removal syntax
- `internal/cli/del.go` - Parse multiple arguments

**Core Functionality**:

**Delete Multiple**:
```bash
# Delete tasks 1, 3, and 5
togodo del 1 3 5
# Output:
# 1 task one
# 3 task three
# 5 task five
# TODO: 3 tasks deleted.
```

**Delete Term**:
```bash
# Remove @work from task 2
togodo del 2 @work
# Before: 2 review code @work @urgent
# After:  2 review code @urgent
# Output: TODO: Removed '@work' from task 2
```

**Testing Approach**:
1. Unit test multiple deletion
2. Unit test term removal
3. Test error cases (invalid numbers, missing terms)
4. Integration tests for both features

---

### Group 7: Move Command (2 skipped tests)

**Status**: Nice-to-have - Multi-file management
**Complexity**: Medium
**Estimated Effort**: 8-10 hours

#### Tests Requiring Implementation
- `TestBasicMove` - Move tasks between files
- `TestMoveUsage` - Move command validation

#### Implementation Requirements

**New Files Needed**:
- `cmd/move.go` - Business logic for moving tasks
- `cmd/move_test.go` - Unit tests
- `internal/cli/move.go` - CLI command wrapper

**Core Functionality**:
1. Remove task from source file (todo.txt)
2. Append task to destination file
3. Support moving to done.txt or custom files
4. Preserve task format (priority, contexts, projects)

**Expected Behavior**:
```bash
# Move task to done.txt
togodo move 1 done.txt
# Output: TODO: Moved task 1 to done.txt

# Move to custom file
togodo move 3 someday.txt
# Output: TODO: Moved task 3 to someday.txt
```

**Dependencies**:
- File I/O for multiple todo files
- Repository pattern extension to support multiple files
- Configuration for file paths

**Testing Approach**:
1. Unit test move operation
2. Test with done.txt
3. Test with custom files
4. Test error cases (file not found, invalid task number)

---

### Group 8: Multiline Support (2 skipped tests)

**Status**: Low priority - Advanced formatting
**Complexity**: High
**Estimated Effort**: 16-20 hours

#### Tests Requiring Implementation
- `TestMultilineAdd` - Adding multiline tasks
- `TestMultilineHandling` - Handling multiline task text

#### Implementation Requirements

**Modifications Needed**:
- `todotxtlib/parser.go` - Parse multiline tasks
- `todotxtlib/writer.go` - Write multiline tasks
- `cmd/add.go` - Accept multiline input
- `internal/cli/list.go` - Display multiline tasks

**Core Functionality**:
1. Accept multiline input (via heredoc, interactive editor, or escaped newlines)
2. Store multiline tasks in todo.txt (format TBD)
3. Display multiline tasks correctly
4. Preserve multiline format through operations

**Challenges**:
- Todo.txt format is line-based (one task per line)
- Need to escape or encode newlines
- Display must handle terminal width
- Editing multiline tasks is complex

**Possible Formats**:
```
# Option 1: Escape newlines
task text line 1\ntask text line 2

# Option 2: Continuation marker
task text line 1
+ task text line 2

# Option 3: Store in extended attributes
```

**Testing Approach**:
1. Define multiline format spec
2. Unit test parsing
3. Unit test writing
4. Integration test add/list/edit workflow
5. Test edge cases (very long lines, many lines)

**Notes**:
- This is a non-standard extension
- May break compatibility with other todo.txt tools
- Consider if this feature is truly needed
- Recommend deferring until other features complete

---

### Group 9: Help System Enhancements (1 skipped test)

**Status**: Low priority - Documentation feature
**Complexity**: Low
**Estimated Effort**: 2-3 hours

#### Tests Requiring Implementation
- `TestShortHelp` - Condensed help output

#### Implementation Requirements

**New Files Needed**:
- `internal/cli/shorthelp.go` - Short help command (or modify existing help)

**Core Functionality**:
1. Display condensed list of commands
2. Show one-line description per command
3. More compact than `--help`

**Expected Behavior**:
```bash
# Short help
togodo shorthelp
# Output:
# add (a)      - Add new task
# list (ls)    - List tasks
# do (x)       - Mark task done
# pri (p)      - Set priority
# del (rm)     - Delete task
# [... etc ...]
```

**Testing Approach**:
1. Integration test output format
2. Verify all commands listed
3. Verify aliases shown

---

## Implementation Priority Recommendations

### Phase 2A: Core Features (High Priority)
**Duration**: 2-3 weeks
**Goal**: Achieve 80%+ test pass rate

1. **Fix list sorting issue** (Priority 0 - Blocks 8 tests)
2. **Archive command** (5 tests) - Most requested feature
3. **Config command** (3 tests) - Important for usability
4. **Deduplicate command** (2 tests) - Data quality
5. **Report command** (1 test) - User value

**Outcome**: 27 + 8 + 5 + 3 + 2 + 1 = 46 passing tests (82% pass rate)

### Phase 2B: Enhanced Features (Medium Priority)
**Duration**: 2-3 weeks
**Goal**: Improve power user experience

1. **Date support** (2 tests) - Standard format compliance
2. **Enhanced delete** (2 tests) - Power user features
3. **Move command** (2 tests) - Multi-file support
4. **Short help** (1 test) - Nice-to-have

**Outcome**: 46 + 7 = 53 passing tests (95% pass rate)

### Phase 2C: Advanced Features (Low Priority)
**Duration**: 3-4 weeks
**Goal**: Complete feature set

1. **Multiline support** (2 tests) - Complex, non-standard

**Outcome**: 55 passing tests (98% pass rate)

---

## Estimated Total Effort

**Phase 2A**: 30-40 hours (2-3 weeks)
**Phase 2B**: 28-36 hours (2-3 weeks)
**Phase 2C**: 16-20 hours (1-2 weeks)

**Total Phase 2**: 74-96 hours (5-8 weeks)

---

## Success Metrics

- **Code Coverage**: Maintain >80% test coverage
- **Test Pass Rate**: Achieve 95%+ passing tests
- **Compatibility**: Full todo.txt-cli command parity
- **Performance**: All commands complete in <100ms for typical files
- **Documentation**: All commands documented with examples
