#!/bin/bash

echo "Testing git @ work branch name formatting..."

# Source the work.sh file to test the format_branch_name function
source git_at_cmds/work.sh

# Test cases
test_cases=(
    "Incorrect Branch Name"
    "Fix Login Bug!"
    "Add User Authentication"
    "Update API Documentation"
    "Fix_Crash_On_Startup"
    "Update.Dependencies"
    "Remove Old Code!!!"
    "Add New Feature (v2.0)"
    "Fix Bug #123"
    "Update README.md"
    "Add Tests for Login"
    "Refactor User Model"
)

echo "Testing format_branch_name function:"
echo "====================================="

for test_input in "${test_cases[@]}"; do
    expected_output=$(format_branch_name "$test_input")
    echo "Input:  '$test_input'"
    echo "Output: '$expected_output'"
    echo "---"
done

echo ""
echo "Testing complete! The format_branch_name function should convert:"
echo "- Uppercase to lowercase"
echo "- Spaces, underscores, dots to hyphens"
echo "- Remove special characters"
echo "- Remove multiple consecutive hyphens"
echo "- Remove leading/trailing hyphens" 