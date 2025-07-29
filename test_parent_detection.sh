#!/bin/bash

echo "Testing parent branch detection..."

# Source the squash.sh file to test the detect_parent_branch function
source git_at_cmds/squash.sh

echo "Current branch: $(git rev-parse --abbrev-ref HEAD)"
echo ""

echo "Testing detect_parent_branch function:"
echo "====================================="

# Test the function
parent_branch=$(detect_parent_branch)
echo "Result: $parent_branch"
echo ""

echo "Available branches:"
git branch --list
echo ""

echo "Git config for current branch:"
git config "branch.$(git rev-parse --abbrev-ref HEAD).merge" 2>/dev/null || echo "No merge config"
echo ""

echo "Upstream tracking:"
git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null || echo "No upstream tracking"
echo ""

echo "Configured trunk:"
git config at.trunk 2>/dev/null || echo "No trunk configured"
echo ""

echo "Merge bases with other branches:"
while IFS= read -r branch; do
    if [ -n "$branch" ] && [ "$branch" != "$(git rev-parse --abbrev-ref HEAD)" ]; then
        merge_base=$(git merge-base "$branch" HEAD 2>/dev/null || echo "")
        if [ -n "$merge_base" ]; then
            merge_date=$(git log -1 --format="%ct" "$merge_base" 2>/dev/null || echo "0")
            commits_ahead=$(git rev-list --count "$merge_base..HEAD" 2>/dev/null || echo "0")
            echo "  $branch: merge_base=$merge_base, date=$merge_date, commits_ahead=$commits_ahead"
        fi
    fi
done < <(git branch --list | sed 's/^[* ]*//') 