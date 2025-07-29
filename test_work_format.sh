#!/bin/bash

# Source the work.sh file to get the function
source git_at_cmds/work.sh

echo "Testing format_branch_name function:"
echo "Input: 'sweep command fixes'"
result=$(format_branch_name "sweep command fixes")
echo "Output: '$result'"
echo "Expected: 'sweep-command-fixes'"

if [ "$result" = "sweep-command-fixes" ]; then
    echo "✅ SUCCESS: Function works correctly!"
else
    echo "❌ FAILED: Expected 'sweep-command-fixes', got '$result'"
fi 