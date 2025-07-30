# Unit Tests

Unit tests for individual functions and components of GitAT.

## Test Files

- `test-array-handling.sh` - Tests for array operations
- `test-variable-initialization.sh` - Tests for variable handling
- `test-branch-validation.sh` - Tests for branch name validation
- `test-markdown-formatting.sh` - Tests for markdown generation
- `test-git-commands.sh` - Tests for Git command execution

## Running Unit Tests

```bash
./test/unit/run-unit-tests.sh
```

## Test Structure

Each unit test should:

1. Test a single function or component
2. Have clear test cases
3. Provide meaningful error messages
4. Exit with status 0 for success, non-zero for failure
