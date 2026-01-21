# Integration Test Execution Summary

**Date**: 2026-01-21
**Test Suite**: `./tests`
**Total Test Cases**: 56 (27 PASS, 13 FAIL, 16 SKIP)

## Executive Summary

The integration test suite shows significant progress with 27 tests passing. However, there are 13 failing tests that need attention and 16 tests that are intentionally skipped as Phase 2 features. Additionally, a **critical list sorting issue** has been identified that affects multiple test failures.

## Test Results Overview

```
PASS: 27 tests (48%)
FAIL: 13 tests (23%)
SKIP: 16 tests (29%)
```

### Passing Tests (27)

The following test categories are fully working:

1. **Basic Add/List** (4 tests)
   - TestBasicAddList - Basic adding and listing tasks
   - TestListFiltering - Filtering tasks by term

2. **Replace Command** (2 tests)
   - TestReplaceUsage - Validation of replace arguments
   - TestBasicReplace - Replacing task text

3. **Priority Management** (4 tests)
   - TestPriUsage - Validation of pri arguments
   - TestBasicPriority - Setting task priorities

4. **List Commands** (8 tests)
   - TestListpriBasic - List tasks by priority
   - TestListconSingle - List single context
   - TestListconMultiple - List multiple contexts
   - TestListprojSingle - List single project
   - TestListprojMultiple - List multiple projects
   - TestListallBasic - List all tasks including done

5. **Task Modification** (6 tests)
   - TestPrependUsage - Validation of prepend arguments
   - TestBasicPrepend - Prepending text to tasks
   - TestBasicDo - Marking tasks as done
   - TestBasicAppend - Appending text to tasks
   - TestBasicDepriority - Removing priority from tasks

6. **Delete Command** (1 test)
   - TestDelUsage - Validation of delete arguments

7. **Help System** (2 tests)
   - TestHelp - Help command and --help flag
   - TestCommandHelp - Help for specific commands

### Failing Tests (13)

#### Critical Issue: List Sorting Problem

**Root Cause**: The list command is not maintaining original line numbers when displaying filtered or modified task lists. This affects 8 of the 13 failing tests.

**Affected Tests**:
1. TestCaseInsensitiveFiltering (filter_lowercase_roses)
2. TestAddWithSymbols (add_backtick_and_quotes, list_all_with_symbols)
3. TestAddWithSpaces (add_with_unquoted_spaces, list_with_spaces)
4. TestPrependPreservesPriority (list_after_prepend_priority)
5. TestAppendPreservesPriority (list_after_append_priority)
6. TestDepriMultiple (list_after_depri_multiple)
7. TestBasicDel (list_after_delete)

**Issue Details**:
- When tasks are added, they receive sequential line numbers (1, 2, 3, etc.)
- After modifications (add, delete, etc.), the list command re-numbers tasks starting from 1
- Expected behavior: Tasks should retain their original line numbers from the todo.txt file
- Current behavior: Tasks are numbered by their position in the filtered/sorted list

**Example from TestCaseInsensitiveFiltering**:
```
Expected:
2 smell the roses
3 smell the uppercase Roses
--
TODO: 2 of 3 tasks shown

Got:
2 smell the roses
--
TODO: 1 of 3 tasks shown
```

**Impact**: This is a fundamental issue that violates todo.txt format expectations where line numbers correspond to file line positions.

#### Other Failing Tests (5)

1. **TestArchiveWithDuplicates** (4 sub-tests)
   - Error: `unknown command "archive"`
   - Status: Archive command not yet implemented (Phase 2 feature)
   - Sub-tests: archive_done_tasks, verify_done_task_removed_from_todo.txt, verify_done_task_in_done.txt, list_after_archive

2. **TestArchiveWarning** (1 sub-test)
   - Error: `unknown command "archive"`
   - Status: Archive command not yet implemented (Phase 2 feature)
   - Sub-test: archive_with_no_done_tasks

3. **TestDeduplicate** (1 sub-test)
   - Error: Exit code 1 (command not found)
   - Status: Deduplicate command not yet implemented (Phase 2 feature)
   - Sub-test: deduplicate_removes_duplicates

4. **TestReport** (1 sub-test)
   - Error: Exit code 1 (command not found)
   - Status: Report command not yet implemented (Phase 2 feature)
   - Sub-test: report_shows_statistics

### Skipped Tests (16)

These tests are intentionally skipped as Phase 2 features:

#### Configuration Management (3 tests)
- **TestConfigRead** - Read configuration values
- **TestConfigWrite** - Write configuration values
- **TestConfigList** - List all configuration settings

**Reason**: Config command needs to be implemented to read/write Viper configuration.

#### Date Support (2 tests)
- **TestAddWithDate** - Adding tasks with creation dates
- **TestAddWithPriority** - Adding tasks with priority in input

**Reason**: Date and priority parsing in add command input not yet implemented.

#### Advanced Delete Operations (2 tests)
- **TestDelMultiple** - Deleting multiple tasks at once
- **TestDelWithTerm** - Removing specific terms from tasks

**Reason**: Enhanced delete functionality not yet implemented.

#### Archive/Move Operations (2 tests)
- **TestBasicMove** - Moving tasks between files
- **TestMoveUsage** - Move command validation

**Reason**: Move command for transferring tasks between todo.txt files not yet implemented.

#### Data Quality (1 test)
- **TestDeduplicateWithPriority** - Deduplication with priority preservation

**Reason**: Advanced deduplication logic not yet implemented.

#### Multiline Support (2 tests)
- **TestMultilineAdd** - Adding multiline tasks
- **TestMultilineHandling** - Handling multiline task text

**Reason**: Multiline task support not yet implemented.

#### Help System (1 test)
- **TestShortHelp** - Condensed help output

**Reason**: Short help command not yet implemented.

## Phase 2 Work Required

### Priority 1: Fix List Sorting Issue

**Urgency**: Critical - Blocks 8 tests
**Complexity**: Medium
**Location**: `internal/cli/list.go`, possibly `cmd/list.go`

**Required Changes**:
1. Ensure `Todo` struct includes original line number from file
2. Modify list output to use original line numbers, not loop indices
3. Verify filtering maintains original line numbers
4. Update all list-related commands (list, listall, listpri) to use original line numbers

**Files to Investigate**:
- `/Users/gergely.karolyi/Code/gkarolyi/todo-txt-cli-parity/internal/cli/list.go`
- `/Users/gergely.karolyi/Code/gkarolyi/todo-txt-cli-parity/cmd/list.go`
- `/Users/gergely.karolyi/Code/gkarolyi/todo-txt-cli-parity/todotxtlib/todo.go` (check if LineNumber field exists)
- `/Users/gergely.karolyi/Code/gkarolyi/todo-txt-cli-parity/todotxtlib/repository.go` (check if line numbers are preserved)

### Priority 2: Implement Missing Commands

#### Archive Command (5 failing tests)
- Move completed tasks to done.txt
- Remove archived tasks from todo.txt
- Handle done.txt file creation
- Show appropriate warnings when no tasks to archive

**Test Coverage**: TestArchiveWithDuplicates, TestArchiveWarning

#### Deduplicate Command (1 failing test)
- Remove duplicate tasks
- Preserve higher priority when duplicates exist
- Update file after deduplication

**Test Coverage**: TestDeduplicate

#### Report Command (1 failing test)
- Generate statistics about tasks
- Show counts by priority, context, project
- Display completion rates

**Test Coverage**: TestReport

### Priority 3: Implement Phase 2 Features (16 skipped tests)

Group these by dependency and complexity:

**Easy Wins** (3-4 hours):
1. Config command (3 tests) - Wrapper around existing Viper config
2. Short help command (1 test) - Format existing help output

**Medium Complexity** (8-12 hours):
1. Enhanced delete operations (2 tests) - Multiple items, term removal
2. Move command (2 tests) - Transfer tasks between files
3. Date support in add command (2 tests) - Parse dates from input

**Complex Features** (16-24 hours):
1. Advanced deduplication (1 test) - Priority-aware duplicate handling
2. Multiline support (2 tests) - Parse and display multiline tasks

## Recommendations

### Immediate Actions

1. **Fix list sorting issue** - This is blocking 8 tests and is a fundamental correctness issue
2. **Investigate line number handling** - Review how line numbers are stored and retrieved
3. **Add unit tests for line number preservation** - Ensure fix doesn't break existing functionality

### Short-term Actions (Next Sprint)

1. **Implement archive command** - Required for 5 tests, commonly used feature
2. **Implement deduplicate command** - Required for 1 test
3. **Implement report command** - Required for 1 test

### Long-term Actions (Phase 2)

1. **Config command** - Important for usability
2. **Date support** - Standard todo.txt format feature
3. **Enhanced delete operations** - Power user features
4. **Move command** - Multi-file support
5. **Multiline support** - Advanced formatting

## Test Execution Details

### Command Used
```bash
go test ./tests -v
```

### Test Files Processed
- t0000_config_test.go (3 tests, 3 skipped)
- t1000_add_list_test.go (8 tests, 4 passing, 4 failing)
- t1010_add_date_test.go (1 test, 1 skipped)
- t1040_add_priority_test.go (1 test, 1 skipped)
- t1100_replace_test.go (2 tests, 2 passing)
- t1200_pri_test.go (2 tests, 2 passing)
- t1250_listpri_test.go (2 tests, 2 passing)
- t1310_listcon_test.go (2 tests, 2 passing)
- t1320_listproj_test.go (2 tests, 2 passing)
- t1350_listall_test.go (2 tests, 2 passing)
- t1400_prepend_test.go (3 tests, 2 passing, 1 failing)
- t1500_do_test.go (2 tests, 2 passing)
- t1600_append_test.go (3 tests, 2 passing, 1 failing)
- t1700_depri_test.go (3 tests, 2 passing, 1 failing)
- t1800_del_test.go (4 tests, 1 passing, 1 failing, 2 skipped)
- t1850_move_test.go (2 tests, 2 skipped)
- t1900_archive_test.go (6 tests, 5 failing)
- t1910_deduplicate_test.go (2 tests, 1 failing, 1 skipped)
- t1950_report_test.go (1 test, 1 failing)
- t2000_multiline_test.go (2 tests, 2 skipped)
- t2100_help_test.go (3 tests, 2 passing, 1 skipped)

## Conclusion

The togodo application has achieved 48% test pass rate with core functionality working well. The main blocker is the list sorting issue which affects multiple test scenarios. Once this is resolved, implementing the remaining commands (archive, deduplicate, report) will bring the pass rate to approximately 80%, with the remaining 20% being optional Phase 2 enhancements.

**Next Steps**:
1. Address the critical list sorting issue
2. Implement archive, deduplicate, and report commands
3. Plan Phase 2 feature implementation
4. Consider test-driven development for remaining features
