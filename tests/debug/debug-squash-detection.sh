#!/bin/bash

echo "Debugging git @ squash auto-detection..."

# Check current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "Current branch: $CURRENT_BRANCH"

# Check all local branches
echo ""
echo "All local branches:"
git branch --list

# Check upstream tracking
echo ""
echo "Upstream tracking:"
git branch -vv

# Check git config for merge
echo ""
echo "Git config for merge:"
git config "branch.$CURRENT_BRANCH.merge" 2>/dev/null || echo "Not set"

# Check upstream symbolic reference
echo ""
echo "Upstream symbolic reference:"
git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null || echo "Not set"

# Test merge base with each branch
echo ""
echo "Merge base analysis:"
while IFS= read -r branch; do
    if [ -n "$branch" ] && [ "$branch" != "$CURRENT_BRANCH" ]; then
        echo "Branch: $branch"
        
        # Get merge base
        merge_base=$(git merge-base "$branch" HEAD 2>/dev/null || echo "")
        echo "  Merge base: $merge_base"
        
        if [ -n "$merge_base" ]; then
            # Count commits from merge base to HEAD
            commit_count=$(git rev-list --count "$merge_base..HEAD" 2>/dev/null || echo "0")
            echo "  Commits ahead: $commit_count"
        fi
        
        echo ""
    fi
done < <(git branch --list | sed 's/^[* ]*//')

# Test the actual detect_parent_branch function
echo ""
echo "Testing detect_parent_branch function:"
source git_at_cmds/squash.sh
detect_parent_branch 