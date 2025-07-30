# Debug Scripts

This directory contains debug scripts and utilities for troubleshooting GitAT issues.

## Debug Files

- `debug-branch-issue.sh` - Debug script for branch-related issues
- `debug-sweep-issue.sh` - Debug script for sweep command issues
- `debug-squash-issue.sh` - Debug script for squash command issues
- `debug-squash-detection.sh` - Debug script for squash parent detection
- `check-branch-status.sh` - Utility to check branch status
- `hash-example-output.txt` - Example output from hash command

## Usage

These scripts are designed to help diagnose issues with GitAT commands. They provide detailed output and state information to help identify problems.

```bash
# Run a debug script
./test/debug/debug-branch-issue.sh

# Check branch status
./test/debug/check-branch-status.sh
```

## Notes

- Debug scripts may output sensitive information
- Some scripts may modify repository state
- Use in a test environment when possible
