# Test Verification Notes

## Goal: Perfect Feature Parity with todo.txt-cli

Integration tests should test EXACTLY the same inputs and outputs as upstream bash tests.
If a feature is missing from togodo, the test should FAIL, not pass with different behavior.

Failing tests = Feature roadmap

## Missing Features to Implement

### Command-Line Flags
- [ ] `-t` flag: Enable TODOTXT_DATE_ON_ADD via command line
- [ ] `-p` flag: Plain mode (already implemented?)
- [ ] `-d` flag: Specify custom config file
- [ ] `-f` flag: Force operations
- [ ] `-n` flag: Don't prompt
- [ ] `-@` flag: Hide context names
- [ ] `-+` flag: Hide project names

### Configuration Variables
- [ ] `TODOTXT_DATE_ON_ADD`: Auto-add creation date (currently: `auto_add_creation_date` config)
- [ ] `TODOTXT_PRIORITY_ON_ADD`: Auto-add priority to new tasks
- [ ] Other config variables from todo.sh

### Command Behaviors
- [ ] Carriage return handling in tasks
- [ ] Adding to files without EOL (end of line)
- [ ] Exact output format matching (line-by-line)
- [ ] Error message format matching

## Test Verification Progress

Status codes:
- ✅ Verified - matches upstream exactly
- 🔄 In Progress - being updated to match upstream
- ❌ Needs Update - currently doesn't match upstream
- 🆕 New Feature - togodo-specific, no upstream equivalent

### Tests

- 🆕 t0000_config_test.go - togodo-specific `config` command
- 🔄 t1000_add_list_test.go - missing 2 test cases (CR, EOF)
- ❌ t1010_add_date_test.go - needs `-t` flag support
- ❌ t1040_add_priority_test.go - needs `TODOTXT_PRIORITY_ON_ADD` support
- [ ] t1100_replace_test.go
- [ ] t1200_pri_test.go
- [ ] t1250_listpri_test.go
- [ ] t1310_listcon_test.go
- [ ] t1320_listproj_test.go
- [ ] t1350_listall_test.go
- [ ] t1400_prepend_test.go
- [ ] t1500_do_test.go
- [ ] t1600_append_test.go
- [ ] t1700_depri_test.go
- [ ] t1800_del_test.go
- [ ] t1850_move_test.go
- [ ] t1900_archive_test.go
- [ ] t1910_deduplicate_test.go
- [ ] t1950_report_test.go
- 🆕 t2000_multiline_test.go - togodo extension
- [ ] t2100_help_test.go
