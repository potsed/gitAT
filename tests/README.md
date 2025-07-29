# GitAT Test Suite

This directory contains the comprehensive test suite for GitAT, organized by test type and functionality.

## Directory Structure

```
test/
├── README.md                    # This file
├── run-all.sh                   # Main test runner (runs all test types)
├── run-all-original.sh          # Original test runner (legacy)
├── run-all-legacy.sh            # Legacy test runner from old tests/ directory
├── commands/                    # Command-specific tests
│   ├── README.md
│   ├── run-command-tests.sh
│   ├── test-work.sh
│   ├── test-work-original.sh
│   ├── test-hotfix.sh
│   ├── test-pr.sh
│   ├── test-save.sh
│   ├── test-squash.sh
│   ├── test-squash-enhanced.sh
│   ├── test-branch.sh
│   ├── test-hash.sh
│   ├── test-info.sh
│   └── test-product.sh
├── unit/                        # Unit tests for specific functionality
│   ├── README.md
│   ├── run-unit-tests.sh
│   ├── test-security.sh
│   ├── test-security-legacy.sh
│   ├── test-work-formatting.sh
│   ├── test-pr-fixes.sh
│   ├── test-pr-auto-description.sh
│   ├── test-parent-detection.sh
│   ├── test-markdown-formatting.sh
│   ├── test-array-fix.sh
│   ├── test-array-fix-v2.sh
│   ├── test-true-squash.sh
│   ├── test-info-functions.sh
│   ├── test-committer-extraction.sh
│   ├── test-git-state.sh
│   └── test-unit.sh
├── integration/                 # Integration tests
│   ├── README.md
│   ├── run-integration-tests.sh
│   ├── test-integration.sh
│   ├── test-sweep-remote.sh
│   ├── test-branch-delete.sh
│   ├── test-sweep-fix.sh
│   ├── test-hash-fix.sh
│   ├── test-squash-fix.sh
│   └── test-individual-commands.sh
└── debug/                       # Debug scripts and utilities
    ├── README.md
    ├── debug-branch-issue.sh
    ├── debug-sweep-issue.sh
    ├── debug-squash-issue.sh
    ├── debug-squash-detection.sh
    ├── check-branch-status.sh
    └── hash-example-output.txt
```

## Running Tests

### Run All Tests

```bash
# Run the complete test suite
./test/run-all.sh
```

### Run Specific Test Types

```bash
# Run only command tests
./test/commands/run-command-tests.sh

# Run only unit tests
./test/unit/run-unit-tests.sh

# Run only integration tests
./test/integration/run-integration-tests.sh
```

### Run Individual Tests

```bash
# Run a specific command test
./test/commands/test-work.sh

# Run a specific unit test
./test/unit/test-security.sh

# Run a specific integration test
./test/integration/test-sweep-remote.sh
```

## Test Categories

### Command Tests (`commands/`)

Tests for individual GitAT commands and their functionality:

- `git @ work` - Work branch creation
- `git @ hotfix` - Hotfix branch creation
- `git @ pr` - Pull request creation
- `git @ save` - Commit saving
- `git @ squash` - Commit squashing
- `git @ branch` - Branch management
- `git @ hash` - Hash and status display
- `git @ info` - Information display
- `git @ product` - Product configuration

### Unit Tests (`unit/`)

Tests for specific functionality and edge cases:

- Security validation
- Work branch formatting
- PR auto-description generation
- Parent branch detection
- Markdown formatting
- Array handling fixes
- True squash functionality
- Info function behavior
- Committer extraction
- Git state validation

### Integration Tests (`integration/`)

Tests that verify the interaction between multiple commands and real Git operations:

- Sweep command with remote cleanup
- Branch deletion workflows
- Fix verification tests
- Individual command integration

### Debug Scripts (`debug/`)

Utilities for troubleshooting and debugging GitAT issues:

- Branch issue debugging
- Sweep issue debugging
- Squash issue debugging
- Status checking utilities
- Example outputs

## Requirements

- Must be run in a Git repository
- GitAT must be installed and available in PATH
- Some tests may modify repository state
- Integration tests should be run in a clean repository or test environment

## Legacy Tests

The following legacy test files are preserved for reference:

- `run-all-original.sh` - Original test runner
- `run-all-legacy.sh` - Legacy test runner from old tests/ directory
- `test/unit/test-security-legacy.sh` - Legacy security tests

## Contributing

When adding new tests:

1. Place command-specific tests in `commands/`
2. Place unit tests in `unit/`
3. Place integration tests in `integration/`
4. Place debug utilities in `debug/`
5. Make sure all test files are executable (`chmod +x`)
6. Update the appropriate README.md file
7. Update this main README.md if adding new categories
