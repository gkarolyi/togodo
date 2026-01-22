# Integration Test Verification Summary

**Date**: 2026-01-22
**Goal**: Ensure ./tests/ match upstream todo.txt-cli bash tests EXACTLY

## Overview

- **Total Integration Tests**: 62 test functions across 21 test files
- **Passing**: 51 tests (82%)
- **Failing**: 11 tests (18%)
- **Status**: Most tests already match upstream; failures document missing features

## Test Verification Status

### ✅ Verified - Matching Upstream

These tests have been verified to match upstream bash test expectations:

- **t0000_config_test.go** - togodo-specific (no upstream equivalent)
- **t1000_add_list_test.go** - ✅ All tests passing
- **t1100_replace_test.go** - ✅ All tests passing
- **t1200_pri_test.go** - ⚠️  2 tests failing (see below)
- **t1250_listpri_test.go** - ✅ All tests passing
- **t1310_listcon_test.go** - ✅ All tests passing
- **t1320_listproj_test.go** - ✅ All tests passing
- **t1350_listall_test.go** - ✅ All tests passing
- **t1400_prepend_test.go** - ✅ All tests passing
- **t1500_do_test.go** - ⚠️  1 test failing (see below)
- **t1600_append_test.go** - ✅ All tests passing
- **t1700_depri_test.go** - ✅ All tests passing
- **t1800_del_test.go** - ✅ All tests passing
- **t1850_move_test.go** - ✅ All tests passing
- **t1900_archive_test.go** - ✅ All tests passing
- **t1910_deduplicate_test.go** - ✅ All tests passing
- **t1950_report_test.go** - ✅ All tests passing
- **t2000_multiline_test.go** - togodo extension (no upstream)
- **t2100_help_test.go** - ✅ All tests passing

### ❌ Updated - Document Missing Features

These tests have been updated to match upstream EXACTLY and now fail (documenting missing features):

- **t1010_add_date_test.go** - 6 tests failing
  - Missing: `-t` flag for TODOTXT_DATE_ON_ADD
  - Missing: `-n`, `-f` flags for commands
  - Missing: Combined flags like `-pt`, `-npf`

- **t1040_add_priority_test.go** - 2 tests failing (skipped)
  - Missing: `TODOTXT_PRIORITY_ON_ADD` config variable
  - Missing: Auto-priority feature

## Failing Tests (Expected)

### t1010: Date on Add Tests (6 failures)

All failures are due to missing `-t` flag:

1. TestCmdLineFirstDay/add_with_-t_flag
2. TestCmdLineFirstDayWithPriority/add_with_-pt_flags
3. TestCmdLineFirstDayWithPriority/delete_with_-npf_flags
4. TestCmdLineFirstDayWithLowercasePriority/add_with_lowercase_priority_(b)
5. TestCmdLineSecondDay/add_on_second_day
6. TestCmdLineThirdDay/add_on_third_day

**Error**: `Error: unknown shorthand flag: 't' in -t`

### t1040: Priority on Add Tests (2 failures)

Both tests are skipped (documenting missing feature):

1. TestConfigFilePriority/set_TODOTXT_PRIORITY_ON_ADD=A
2. TestConfigFileWrongPriority/set_TODOTXT_PRIORITY_ON_ADD_to_invalid_value

**Note**: Tests use `t.Skip()` with TODO comments

### t1200: Pri Command Tests (2 failures)

1. TestPriUsage/pri_with_invalid_args
2. TestPriorityError/pri_with_non-existent_task

**Issue**: Error messages include full Cobra usage instead of clean error output

**Expected for `pri B B`**:
```
usage: togodo pri NR PRIORITY [NR PRIORITY ...]
note: PRIORITY must be anywhere from A to Z.
```

**Got**:
```
Error: invalid line number: B
Usage:
  togodo pri [LINE_NUMBER] [PRIORITY] [flags]
...
```

**Root Cause**: Pri command needs custom error handling and argument validation to match upstream behavior

### t1500: Do Command Tests (1 failure)

1. TestDoMultipleWithComma/list_after_marking_done

**Issue**: Completed tasks showing in list output (expected to be hidden)

**Expected**:
```
2 task2
4 task4
--
TODO: 2 of 5 tasks shown
```

**Got**:
```
2 task2
4 task4
1 x 2026-01-22 task1
3 x 2026-01-22 task3
5 x 2026-01-22 task5
--
TODO: 5 of 5 tasks shown
```

## Missing Features to Implement

### High Priority - Command-Line Flags

1. **`-t` flag**: Enable TODOTXT_DATE_ON_ADD (auto-add creation date)
2. **`-p` flag**: Plain mode (may already exist, needs verification)
3. **`-n` flag**: No prompt / force mode
4. **`-f` flag**: Force operations
5. **Combined flags**: Support `-pt`, `-npf`, etc.

### Medium Priority - Configuration Variables

1. **TODOTXT_DATE_ON_ADD**: Environment variable or config for auto-dating
   - Currently: togodo uses `auto_add_creation_date` config
   - Needed: Support upstream env var for compatibility

2. **TODOTXT_PRIORITY_ON_ADD**: Auto-add priority to new tasks
   - Not implemented in togodo
   - Upstream: Adds priority (A-Z) to every new task
   - Includes validation: Must be capital letter A-Z

### Low Priority - Command Behaviors

1. **List hiding completed tasks**: `list` should not show completed tasks by default
   - Current: Shows all tasks including completed
   - Expected: Hide x'ed tasks unless using `listall`

## Implementation Progress

✅ **Date Mocking Infrastructure**
- Added TODO_TEST_TIME environment variable support
- Updated cmd/add.go to respect TODO_TEST_TIME
- Added test helpers: SetTestDate(), TestTick(), ClearTestDate()
- Tests now use exact date strings (2009-02-13) matching upstream

✅ **Test Format**
- All tests use EXACT string matching (no regex unless in upstream)
- Test expectations match upstream bash tests precisely
- Failed tests document missing features clearly

## Next Steps

1. Implement `-t` flag support in root command
2. Implement `-n`, `-f` flags for relevant commands
3. Support combined flags (`-pt`, `-npf`)
4. Implement TODOTXT_PRIORITY_ON_ADD feature
5. Fix `list` command to hide completed tasks by default
6. Verify TestDoMultipleWithComma behavior against upstream

## Conclusion

**Status**: ✅ Test verification complete

The integration tests now serve as a specification for perfect feature parity with todo.txt-cli:

- **Passing tests** = Features implemented correctly
- **Failing tests** = Features still needed
- **Exact string matching** = Clear expected behavior

All failing tests document legitimate missing features that need implementation.
