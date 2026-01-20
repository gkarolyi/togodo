# Plan: Achieve Feature Parity with todo.txt-cli

**Last Updated:** 2026-01-20
**Status:** Phase 1 & 2 Complete ✅ | Phase 3 In Progress 🔄

## Goal
Port all bash integration tests from todo.txt-cli to Go, then implement missing features in togodo's cmd package to achieve complete feature parity with todo.txt-cli.

## Approach: Test-Driven Feature Parity

---

### Phase 1: Set Up Test Infrastructure ✅ COMPLETE

**Goal:** Create Go test infrastructure that mirrors bash test suite structure

**Status:** ✅ Complete (2026-01-20)

**Completed Tasks:**
- ✅ Created `/tests/` directory for integration tests
- ✅ Ported bash test helpers to Go with buffer-based approach
  - Test setup/teardown using buffer repositories
  - BufferOutputWriter for capturing command output
  - Assertion helpers for exact output matching
  - Direct command execution via cmd package (no os/exec)
- ✅ Test runner executes tests in isolation with perfect isolation
- ✅ Test aggregation via standard Go test framework

**Deliverables:**
- ✅ `/tests/test_helpers.go` - Core test infrastructure (220 lines)
- ✅ `BufferOutputWriter` - Custom output capture for tests
- ✅ `internal/cli/presenter.go` - Added `NewPresenterWithDeps()` for dependency injection

**Key Innovation:**
- Buffer-based approach (no file I/O in tests)
- Direct cmd package invocation (no binary compilation)
- Tests run in ~45ms (extremely fast)

**Documentation:**
- `TEST_INFRASTRUCTURE_FIXED.md` - Complete implementation details

---

### Phase 2: Port All Integration Tests ✅ COMPLETE

**Goal:** Convert all bash tests to Go, making them compile and run

**Status:** ✅ Complete (2026-01-20)

**Completed Waves:**

**✅ Wave 1 - Core Commands (6 files, 21 test cases):**
- ✅ `t1000_add_list_test.go` - Add and list functionality
- ✅ `t1200_pri_test.go` - Priority setting
- ✅ `t1500_do_test.go` - Mark done
- ✅ `t1700_depri_test.go` - Remove priority
- ✅ `t1900_archive_test.go` - Archive to done.txt

**✅ Wave 2 - Task Modification (5 files, 13 test cases):**
- ✅ `t1100_replace_test.go` - Replace task text
- ✅ `t1400_prepend_test.go` - Prepend text
- ✅ `t1600_append_test.go` - Append text
- ✅ `t1800_del_test.go` - Delete tasks
- ✅ `t1850_move_test.go` - Move between files

**✅ Wave 3 - Listing & Filtering (4 files, 8 test cases):**
- ✅ `t1250_listpri_test.go` - List by priority
- ✅ `t1310_listcon_test.go` - List contexts
- ✅ `t1320_listproj_test.go` - List projects
- ✅ `t1350_listall_test.go` - List from todo.txt + done.txt

**✅ Wave 4 - Advanced Features (5 files, 7 test cases):**
- ✅ `t1010_add_date_test.go` - Auto-date support
- ✅ `t1040_add_priority_test.go` - Add with priority
- ✅ `t1910_deduplicate_test.go` - Remove duplicates
- ✅ `t1950_report_test.go` - Statistics
- ✅ `t2000_multiline_test.go` - Multiline handling

**✅ Wave 5 - Configuration & Help (2 files, 6 test cases):**
- ✅ `t0000_config_test.go` - Configuration management
- ✅ `t2100_help_test.go` - Help system

**✅ Deferred (appropriately skipped):**
- Completion tests (t6000-series) - shell-specific
- Actions/addon tests (t8000-series) - plugin system

**Deliverables:**
- ✅ 21 test files created
- ✅ 54 top-level test functions
- ✅ 131 total test cases (including subtests)
- ✅ All tests compile and run successfully
- ✅ Test failures are meaningful (real feature gaps, not infrastructure issues)

**Test Results:**
```
PASS:  14 tests (usage/error handling)
FAIL:  24 tests (missing features - expected)
SKIP:  16 tests (advanced features for later)
Time:  0.045s
```

**Documentation:**
- `PHASE_1_2_COMPLETION_REVIEW.md` - Comprehensive completion review

---

### Phase 3: Gap Analysis 🔄 IN PROGRESS

**Goal:** Identify all missing/different features using test results

**Status:** 🔄 Partial - Initial analysis complete, prioritization needed

**Completed:**
- ✅ Run full test suite: `go test ./tests` (all tests run successfully)
- ✅ Initial gap report generated
- ✅ Test failures categorized

**Documentation Created:**
- ✅ `FEATURE_GAP_REPORT.md` - Initial gap analysis
- ✅ `ACTUAL_ISSUES.md` - Detailed issue breakdown with 7 critical findings

**Findings Summary:**

**🔴 Critical Issues (breaks basic functionality):**
1. Add command treats each word as separate task
2. Formatter collapses multiple spaces (data loss)
3. Line numbers wrong after filtering (makes `do` unusable)

**🟡 High Priority (compatibility):**
4. Search is case-SENSITIVE (should be case-insensitive)
5. Add command missing line number and confirmation message
6. List command missing separator `--` and summary line

**📊 Medium Priority:**
7. Sorting behavior differs from todo.txt-cli

**Missing Commands (15):**
- Core: `del/rm`, `archive`, `depri`, `append`, `prepend`, `replace`
- Listing: `listpri`, `listcon`, `listproj`, `listall`, `listfile`
- Advanced: `deduplicate`, `report`, `move/mv`, `addm`

**Remaining Tasks:**
- ⏳ Create prioritized implementation roadmap
- ⏳ Estimate effort for each command
- ⏳ Group related changes for efficient implementation

---

### Phase 4: Implement Missing Features ⏳ READY TO START

**Goal:** Make tests pass by implementing features

**Status:** ⏳ Ready - Tests provide specification

**Strategy:**
1. Fix critical issues first (add command, formatter, line numbers)
2. Implement core missing commands (del, archive, depri)
3. Add listing commands (listpri, listcon, listproj)
4. Implement task modification (append, prepend, replace)
5. Add advanced features (deduplicate, report)

**Implementation Pattern:**
1. Create `cmd/[command].go` following existing patterns
2. Add service method in `todotxtlib/service.go` if needed
3. Add repository method in `todotxtlib/repository.go` if needed
4. Wire up in `cmd/root.go`
5. Run integration tests: `go test ./tests -run Test[Command]`
6. Fix issues until tests pass

**Priority Order (by impact):**

**Phase 4A - Fix Critical Issues (Week 1):**
1. Fix add command to join args into single task
2. Fix formatter to preserve multiple spaces
3. Fix line number tracking after filtering
4. Fix case-insensitive search
5. Fix add/list output format (line numbers, separators, summaries)

**Phase 4B - Core Commands (Week 2):**
1. `del/rm` - Delete tasks by line number
2. `archive` - Move done tasks to done.txt (with done.txt support)
3. `depri/dp` - Remove priority from task
4. `replace` - Replace entire task text

**Phase 4C - Task Modification (Week 3):**
1. `append/app` - Append text to task
2. `prepend/prep` - Prepend text to task
3. Verify multi-arg handling across commands

**Phase 4D - Listing Commands (Week 4):**
1. `listpri/lsp` - List tasks by priority
2. `listcon/lsc` - List all contexts
3. `listproj/lsprj` - List all projects
4. `listall/lsa` - List from todo.txt and done.txt

**Phase 4E - Advanced Features (Week 5):**
1. `deduplicate` - Remove duplicate tasks
2. `report` - Generate task statistics
3. `move/mv` - Move task between files
4. Auto-date support (optional configuration)

**Not Implementing (out of scope):**
- `addm` - Current add already handles multiple tasks
- `listfile/lf` - Low value, use list command
- Plugin/addon system - Complex, defer to future

---

### Phase 5: Verification & Documentation ⏳ PENDING

**Goal:** Confirm feature parity achieved

**Status:** ⏳ Pending Phase 4 completion

**Tasks:**
- ⏳ Run full test suite: `go test ./...` (target: 100% core tests passing)
- ⏳ Performance comparison: togodo vs todo.txt-cli
  - Command execution speed
  - Large file handling (1000+ tasks)
- ⏳ Update README with feature parity status
- ⏳ Document intentional differences (if any)
- ⏳ Verify TUI still works with new cmd implementations
- ⏳ Create migration guide from todo.txt-cli

**Success Criteria:**
- All core integration tests passing (target: 90%+)
- Performance faster than todo.txt-cli
- TUI functional
- Documentation complete

---

## Current Status Dashboard

### Overall Progress
```
✅ Phase 1: Test Infrastructure     [████████████████████] 100%
✅ Phase 2: Port Integration Tests  [████████████████████] 100%
🔄 Phase 3: Gap Analysis            [███████████░░░░░░░░░]  60%
⏳ Phase 4: Implement Features      [░░░░░░░░░░░░░░░░░░░░]   0%
⏳ Phase 5: Verification            [░░░░░░░░░░░░░░░░░░░░]   0%
```

### Test Status
```
Total Tests:    54
Passing:        14 (26%)
Failing:        24 (44%) - Features to implement
Skipped:        16 (30%) - Advanced features
```

### Commands Status
```
Implemented:     5 (add, do, list, pri, tidy)
Missing:        15 (core: 6, listing: 4, advanced: 5)
Partially:       2 (do, pri need behavior fixes)
```

---

## Critical Files Created/Modified

### New Files Created (Phase 1 & 2):
- ✅ `/tests/test_helpers.go` - Test infrastructure (220 lines)
- ✅ `/tests/t0000_config_test.go` - Config tests
- ✅ `/tests/t1000_add_list_test.go` - Add/list tests
- ✅ `/tests/t1010_add_date_test.go` - Date tests
- ✅ `/tests/t1040_add_priority_test.go` - Priority tests
- ✅ `/tests/t1100_replace_test.go` - Replace tests
- ✅ `/tests/t1200_pri_test.go` - Priority tests
- ✅ `/tests/t1250_listpri_test.go` - List priority tests
- ✅ `/tests/t1310_listcon_test.go` - List context tests
- ✅ `/tests/t1320_listproj_test.go` - List project tests
- ✅ `/tests/t1350_listall_test.go` - List all tests
- ✅ `/tests/t1400_prepend_test.go` - Prepend tests
- ✅ `/tests/t1500_do_test.go` - Do tests
- ✅ `/tests/t1600_append_test.go` - Append tests
- ✅ `/tests/t1700_depri_test.go` - Depriority tests
- ✅ `/tests/t1800_del_test.go` - Delete tests
- ✅ `/tests/t1850_move_test.go` - Move tests
- ✅ `/tests/t1900_archive_test.go` - Archive tests
- ✅ `/tests/t1910_deduplicate_test.go` - Deduplicate tests
- ✅ `/tests/t1950_report_test.go` - Report tests
- ✅ `/tests/t2000_multiline_test.go` - Multiline tests
- ✅ `/tests/t2100_help_test.go` - Help tests

### Documentation Files:
- ✅ `FEATURE_PARITY_PLAN.md` - This file
- ✅ `FEATURE_GAP_REPORT.md` - Gap analysis
- ✅ `ACTUAL_ISSUES.md` - Detailed issue breakdown
- ✅ `PHASE_1_2_COMPLETION_REVIEW.md` - Phase 1 & 2 review
- ✅ `TEST_INFRASTRUCTURE_FIXED.md` - Infrastructure details

### Modified Files (Phase 1):
- ✅ `internal/cli/presenter.go` - Added `NewPresenterWithDeps()` constructor

### Files to Create (Phase 4):
- ⏳ `/cmd/del.go` - Delete command
- ⏳ `/cmd/archive.go` - Archive command
- ⏳ `/cmd/depri.go` - Deprioritize command
- ⏳ `/cmd/replace.go` - Replace command
- ⏳ `/cmd/append.go` - Append command
- ⏳ `/cmd/prepend.go` - Prepend command
- ⏳ `/cmd/listpri.go` - List priority command
- ⏳ `/cmd/listcon.go` - List context command
- ⏳ `/cmd/listproj.go` - List project command
- ⏳ `/cmd/listall.go` - List all command
- ⏳ `/cmd/deduplicate.go` - Deduplicate command
- ⏳ `/cmd/report.go` - Report command
- ⏳ `/todotxtlib/done.go` - Done.txt support

### Files to Modify (Phase 4):
- ⏳ `/cmd/add.go` - Fix arg joining, output format
- ⏳ `/cmd/list.go` - Fix output format, summary line
- ⏳ `/cmd/root.go` - Register new commands
- ⏳ `/internal/cli/formatter.go` - Fix space preservation, line number format
- ⏳ `/todotxtlib/filter.go` - Add case-insensitive search
- ⏳ `/todotxtlib/service.go` - Add new service methods
- ⏳ `/todotxtlib/repository.go` - Track original indices

---

## Open Questions & Decisions

### 1. Done.txt Handling ✅ DECIDED
**Decision:** Implement done.txt support (archive to done.txt)
**Rationale:** Feature parity with todo.txt-cli, preserves history
**Action:** Keep `tidy` command (deletes), add `archive` command (moves to done.txt)

### 2. Config Compatibility ⏳ PENDING
**Question:** Should togodo read todo.txt-cli config files?
**Options:**
- A) Support both config formats (togodo + todo.txt-cli)
- B) Provide migration tool
- C) Togodo config only

### 3. Plugin System ⏳ DEFERRED
**Decision:** Defer to post-parity phase
**Rationale:** Complex, low priority, can add later without breaking changes

### 4. Line Number Tracking ⏳ NEEDS DESIGN
**Question:** How to track original line numbers after filtering?
**Options:**
- A) Add Index field to Todo struct
- B) Return tuples (index, todo) from filter operations
- C) Separate indices array passed to formatter
**Decision Needed By:** Phase 4A

---

## Success Criteria

### Phase Completion Criteria:
1. ✅ Phase 1: Test infrastructure working with buffer-based approach
2. ✅ Phase 2: All core tests ported, compile, and run
3. 🔄 Phase 3: Gap report complete with prioritized roadmap
4. ⏳ Phase 4: All critical issues fixed, core commands implemented
5. ⏳ Phase 5: 90%+ tests passing, performance verified

### Final Success Criteria:
1. **Test Coverage:** 90%+ of core integration tests passing
2. **Command Parity:** All core commands implemented and working
3. **Performance:** Faster than todo.txt-cli (Go vs bash)
4. **TUI:** Existing functionality preserved
5. **Documentation:** Complete migration guide

---

## Notes & Lessons Learned

### What Worked Well:
- ✅ Buffer-based testing approach (fast, isolated)
- ✅ Direct cmd package invocation (no exec overhead)
- ✅ Test-driven approach (clear specification)
- ✅ Phased approach (manageable chunks)
- ✅ Comprehensive documentation

### Challenges Encountered:
- Initial exec-based approach was slow (fixed with buffers)
- Presenter needed dependency injection (added constructor)
- Some bash test idioms don't translate directly to Go

### Best Practices Established:
- Use buffers for test I/O (not files)
- Call cmd package directly (not via exec)
- Document test source (bash file reference)
- Use t.Skip() with explanation for future features
- Keep test names descriptive

---

## Timeline

- **2026-01-20:** Phases 1 & 2 complete
- **Week 1 (target):** Phase 3 complete, Phase 4A started
- **Week 2-5 (target):** Phase 4B-E implementation
- **Week 6 (target):** Phase 5 verification

---

## Resources

- **Bash Tests Source:** `/tmp/todo.txt-cli/tests/`
- **Bash Script:** `/tmp/todo.txt-cli/todo.sh`
- **Test Command:** `go test ./tests`
- **Single Test:** `go test ./tests -run TestName`
- **Verbose:** `go test -v ./tests`
