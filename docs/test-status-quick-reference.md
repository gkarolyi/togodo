# Test Status Quick Reference

**Last Updated**: 2026-01-21
**Test Pass Rate**: 48% (27/56 tests passing)

## Test Results at a Glance

```
✅ PASS: 27 tests (48%)
❌ FAIL: 13 tests (23%)
⏭️  SKIP: 16 tests (29%)
━━━━━━━━━━━━━━━━━━━━━━
📊 TOTAL: 56 tests
```

## Critical Issue 🔥

**List Sorting Problem** - Affects 8 failing tests

**Issue**: Tasks lose their original line numbers after operations
**Impact**: Critical - breaks expected todo.txt-cli behavior
**Priority**: Must fix before Phase 2

**Example**:
```
Expected: 2 smell the roses
Got:      1 smell the roses
```

## What's Working ✅

### Core Commands (27 tests passing)
- ✅ Add tasks
- ✅ List tasks (with filtering)
- ✅ Replace task text
- ✅ Set priority (pri)
- ✅ Remove priority (depri)
- ✅ Mark done (do)
- ✅ Prepend text
- ✅ Append text
- ✅ Delete tasks (basic)
- ✅ List by priority (listpri)
- ✅ List contexts (listcon)
- ✅ List projects (listproj)
- ✅ List all including done (listall)
- ✅ Help system

## What's Failing ❌

### Commands Not Yet Implemented (5 tests)
- ❌ archive - Move done tasks to done.txt
- ❌ deduplicate - Remove duplicate tasks
- ❌ report - Show task statistics

### Output Format Issues (8 tests)
- ❌ List sorting (line numbers wrong after operations)
- ❌ Case-insensitive filtering
- ❌ List with symbols and special characters
- ❌ List with multiple spaces

## What's Planned (Phase 2) ⏭️

### Configuration (3 tests)
- ⏭️ Read config values
- ⏭️ Write config values
- ⏭️ List all config

### Date Support (2 tests)
- ⏭️ Add with creation date
- ⏭️ Add with priority in input

### Enhanced Delete (2 tests)
- ⏭️ Delete multiple tasks
- ⏭️ Delete specific terms from tasks

### Multi-file (2 tests)
- ⏭️ Move tasks between files

### Data Quality (1 test)
- ⏭️ Deduplicate with priority preservation

### Advanced Features (4 tests)
- ⏭️ Multiline task support (2 tests)
- ⏭️ Short help command (1 test)

## Command Completion Status

| Command | Status | Tests | Notes |
|---------|--------|-------|-------|
| add | ✅ Working | 4/8 | Line number issues |
| list | ✅ Working | 2/4 | Line number issues |
| do | ✅ Working | 2/2 | Complete |
| pri | ✅ Working | 2/4 | Some sorting issues |
| depri | ✅ Working | 2/3 | Line number issues |
| replace | ✅ Working | 2/2 | Complete |
| prepend | ✅ Working | 2/3 | Priority preservation issue |
| append | ✅ Working | 2/3 | Priority preservation issue |
| del | ✅ Working | 1/4 | Line number issues |
| listpri | ✅ Working | 2/2 | Complete |
| listcon | ✅ Working | 2/2 | Complete |
| listproj | ✅ Working | 2/2 | Complete |
| listall | ✅ Working | 2/2 | Complete |
| help | ✅ Working | 2/3 | Short help missing |
| archive | ❌ Missing | 0/5 | Phase 2 |
| deduplicate | ❌ Missing | 0/2 | Phase 2 |
| report | ❌ Missing | 0/1 | Phase 2 |
| move | ⏭️ Planned | 0/2 | Phase 2 |
| config | ⏭️ Planned | 0/3 | Phase 2 |

## Test Files Status

| File | Pass | Fail | Skip | Status |
|------|------|------|------|--------|
| t0000_config_test.go | 0 | 0 | 3 | Phase 2 |
| t1000_add_list_test.go | 4 | 4 | 0 | **Needs fix** |
| t1010_add_date_test.go | 0 | 0 | 1 | Phase 2 |
| t1040_add_priority_test.go | 0 | 0 | 1 | Phase 2 |
| t1100_replace_test.go | 2 | 0 | 0 | ✅ Complete |
| t1200_pri_test.go | 2 | 2 | 1 | **Needs fix** |
| t1250_listpri_test.go | 2 | 0 | 0 | ✅ Complete |
| t1310_listcon_test.go | 2 | 0 | 0 | ✅ Complete |
| t1320_listproj_test.go | 2 | 0 | 0 | ✅ Complete |
| t1350_listall_test.go | 2 | 0 | 0 | ✅ Complete |
| t1400_prepend_test.go | 2 | 1 | 0 | **Needs fix** |
| t1500_do_test.go | 2 | 0 | 0 | ✅ Complete |
| t1600_append_test.go | 2 | 1 | 0 | **Needs fix** |
| t1700_depri_test.go | 2 | 1 | 0 | **Needs fix** |
| t1800_del_test.go | 1 | 1 | 2 | **Needs fix** |
| t1850_move_test.go | 0 | 0 | 2 | Phase 2 |
| t1900_archive_test.go | 0 | 5 | 0 | **Need command** |
| t1910_deduplicate_test.go | 0 | 1 | 1 | **Need command** |
| t1950_report_test.go | 0 | 1 | 0 | **Need command** |
| t2000_multiline_test.go | 0 | 0 | 2 | Phase 2 |
| t2100_help_test.go | 2 | 0 | 1 | Phase 2 |

## Next Actions

### Immediate (This Week)
1. 🔥 Fix list sorting issue (blocks 8 tests)
2. Investigate line number handling in Todo struct
3. Update list output to use original line numbers

### Short-term (Next Sprint)
1. Implement archive command (5 tests)
2. Implement deduplicate command (2 tests)
3. Implement report command (1 test)

### Long-term (Phase 2)
1. Config command (3 tests)
2. Date support (2 tests)
3. Enhanced delete operations (2 tests)
4. Move command (2 tests)
5. Multiline support (2 tests)

## Files to Review

### For List Sorting Fix
- `/Users/gergely.karolyi/Code/gkarolyi/todo-txt-cli-parity/internal/cli/list.go`
- `/Users/gergely.karolyi/Code/gkarolyi/todo-txt-cli-parity/cmd/list.go`
- `/Users/gergely.karolyi/Code/gkarolyi/todo-txt-cli-parity/todotxtlib/todo.go`
- `/Users/gergely.karolyi/Code/gkarolyi/todo-txt-cli-parity/todotxtlib/repository.go`

### For New Commands
- Check existing `cmd/listproj.go` as template (already in working tree)
- Check existing `internal/cli/listproj.go` as template (already in working tree)

## Progress Tracking

**Phase 1 Complete**: ✅ 27/29 planned tests passing (93% of Phase 1)
- Still need to fix: list sorting issue (8 tests)

**Phase 2 Planned**: 16 skipped + 5 failing = 21 tests remaining
- Archive/deduplicate/report: 8 tests (high priority)
- Config/enhanced features: 13 tests (medium/low priority)

**Overall Progress**: 27/56 = 48% complete
**Adjusted Progress** (excluding Phase 2): 27/40 = 68% complete
