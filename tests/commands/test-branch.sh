#!/bin/bash

echo "Testing branch structure..."

# Show current branch
echo "Current branch: $(git rev-parse --abbrev-ref HEAD)"

# Show all branches
echo ""
echo "All branches:"
git branch -a

# Show branch history
echo ""
echo "Branch history:"
git log --oneline --graph --all -10

# Show which branch contains the current HEAD
echo ""
echo "Branches containing current HEAD:"
git branch --contains HEAD

# Show merge base with master
echo ""
echo "Merge base with master:"
if git rev-parse --verify master >/dev/null 2>&1; then
    merge_base=$(git merge-base master HEAD 2>/dev/null || echo "")
    echo "Merge base: $merge_base"
    commits_ahead=$(git rev-list --count "$merge_base..HEAD" 2>/dev/null || echo "0")
    echo "Commits ahead of master: $commits_ahead"
else
    echo "Master branch does not exist"
fi

# Show merge base with main
echo ""
echo "Merge base with main:"
if git rev-parse --verify main >/dev/null 2>&1; then
    merge_base=$(git merge-base main HEAD 2>/dev/null || echo "")
    echo "Merge base: $merge_base"
    commits_ahead=$(git rev-list --count "$merge_base..HEAD" 2>/dev/null || echo "0")
    echo "Commits ahead of main: $commits_ahead"
else
    echo "Main branch does not exist"
fi 