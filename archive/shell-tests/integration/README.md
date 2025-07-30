# Integration Tests

This directory contains integration tests that test the interaction between multiple GitAT commands and real Git operations.

## Test Files

- `test-sweep-remote.sh` - Tests the sweep command with remote branch cleanup
- `test-branch-delete.sh` - Tests branch deletion functionality
- `test-sweep-fix.sh` - Tests sweep command fixes
- `test-hash-fix.sh` - Tests hash command fixes
- `test-squash-fix.sh` - Tests squash command fixes
- `test-individual-commands.sh` - Tests individual command functionality

## Running Integration Tests

```bash
# Run all integration tests
./test/integration/run-integration-tests.sh

# Run a specific test
./test/integration/test-sweep-remote.sh
```

## Requirements

- Must be run in a Git repository
- Some tests may modify the repository state
- Tests should be run in a clean repository or test environment
