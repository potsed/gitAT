# Command Tests

Tests for specific git @ commands and their functionality.

## Test Files

- `test-work.sh` - Tests for `git @ work` command
- `test-sweep.sh` - Tests for `git @ sweep` command
- `test-pr.sh` - Tests for `git @ pr` command
- `test-squash.sh` - Tests for `git @ squash` command
- `test-hotfix.sh` - Tests for `git @ hotfix` command
- `test-save.sh` - Tests for `git @ save` command
- `test-info.sh` - Tests for `git @ info` command
- `test-branch.sh` - Tests for `git @ branch` command
- `test-version.sh` - Tests for `git @ version` command

## Running Command Tests

```bash
./test/commands/run-command-tests.sh
```

## Test Structure

Each command test should:

1. Test the command's basic functionality
2. Test error handling and edge cases
3. Test command options and flags
4. Verify output formatting
5. Test integration with other commands
