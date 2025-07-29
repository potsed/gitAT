#!/bin/bash

echo "🔍 Checking Branch Status"
echo "========================"
echo ""

# Check current branch
current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "Current branch: $current_branch"
echo ""

# Check all local branches
echo "📋 All Local Branches:"
git branch -v 2>/dev/null || echo "Could not list branches"
echo ""

# Check which branches are merged into master
echo "✅ Branches merged into master:"
git branch --merged master 2>/dev/null || echo "Could not check merged branches"
echo ""

# Check which branches are NOT merged into master
echo "❌ Branches NOT merged into master:"
git branch --no-merged master 2>/dev/null || echo "Could not check unmerged branches"
echo ""

# Check specific branches
echo "🔍 Checking specific branches:"
echo ""

# Check bugfix-incorrect-help-text
if git show-ref --verify --quiet refs/heads/bugfix-incorrect-help-text 2>/dev/null; then
    echo "bugfix-incorrect-help-text:"
    echo "  - Exists: ✅"
    echo "  - Commits ahead of master:"
    git log --oneline master..bugfix-incorrect-help-text 2>/dev/null || echo "    (Could not check)"
    echo "  - Would be deleted by sweep:"
    if git branch --merged master | grep -q "bugfix-incorrect-help-text"; then
        echo "    ✅ Yes (fully merged)"
    else
        echo "    ❌ No (has unmerged changes)"
    fi
else
    echo "bugfix-incorrect-help-text: ❌ Branch does not exist"
fi
echo ""

# Check test-suite
if git show-ref --verify --quiet refs/heads/test-suite 2>/dev/null; then
    echo "test-suite:"
    echo "  - Exists: ✅"
    echo "  - Commits ahead of master:"
    git log --oneline master..test-suite 2>/dev/null || echo "    (Could not check)"
    echo "  - Would be deleted by sweep:"
    if git branch --merged master | grep -q "test-suite"; then
        echo "    ✅ Yes (fully merged)"
    else
        echo "    ❌ No (has unmerged changes)"
    fi
else
    echo "test-suite: ❌ Branch does not exist"
fi
echo ""

echo "💡 Explanation:"
echo "git @ sweep only deletes branches that are FULLY merged into the trunk branch."
echo "Branches with unmerged changes are preserved to prevent data loss."
echo ""
echo "To force delete a branch (even with unmerged changes):"
echo "  git branch -D <branch-name>"
echo ""
echo "To merge a branch before sweeping:"
echo "  git checkout master"
echo "  git merge <branch-name>"
echo "  git @ sweep" 