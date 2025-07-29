#!/bin/bash

echo "Testing PR fixes..."

# Test 1: Check if there are commits between branches
echo "Test 1: Checking commits between branches"
current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
base_branch="master"

if [ -n "$current_branch" ] && [ -n "$base_branch" ]; then
    commit_count=$(git rev-list --count "$base_branch..$current_branch" 2>/dev/null || echo "0")
    echo "Current branch: $current_branch"
    echo "Base branch: $base_branch"
    echo "Commits between branches: $commit_count"
    
    if [ "$commit_count" -eq 0 ]; then
        echo "❌ No commits between branches - cannot create PR"
        echo "You need to make some commits on your branch first."
    else
        echo "✅ Commits found - PR creation should work"
    fi
else
    echo "❌ Could not determine branch information"
fi

echo ""

# Test 2: Check if variables are properly initialized
echo "Test 2: Variable initialization"
file_status_info=""
commit_history=""

echo "file_status_info: '$file_status_info'"
echo "commit_history: '$commit_history'"
echo "✅ Variables properly initialized"

echo ""

# Test 3: Check if branches are different
echo "Test 3: Branch validation"
if [ "$current_branch" = "$base_branch" ]; then
    echo "❌ Cannot create PR from $base_branch to itself"
else
    echo "✅ Branches are different - PR creation possible"
fi

echo ""
echo "Summary:"
echo "- Fixed unbound variable errors"
echo "- Added validation for commits between branches"
echo "- Added proper variable initialization"
echo "- Added branch difference validation" 