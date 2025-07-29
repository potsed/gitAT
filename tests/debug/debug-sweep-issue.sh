#!/bin/bash

echo "🔍 Debugging Sweep Issue"
echo "========================"
echo ""

# Check if we're in a git repo
if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In git repository"
echo ""

# Check current branch
current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "Current branch: $current_branch"
echo ""

# Check if remote is configured
if git remote get-url origin >/dev/null 2>&1; then
    echo "✅ Remote origin configured"
else
    echo "❌ No remote origin configured"
    exit 1
fi

echo ""

# Check specific branches
echo "🔍 Checking specific branches:"
echo ""

# Check bugfix-incorrect-help-text
if git show-ref --verify --quiet refs/heads/bugfix-incorrect-help-text 2>/dev/null; then
    echo "bugfix-incorrect-help-text:"
    echo "  - Exists locally: ✅"
    
    # Check if it's merged
    if git branch --merged master | grep -q "bugfix-incorrect-help-text"; then
        echo "  - Merged into master: ✅"
    else
        echo "  - Merged into master: ❌"
    fi
    
    # Check if it exists remotely
    remote_exists=$(git ls-remote --heads origin "bugfix-incorrect-help-text" 2>/dev/null | wc -l)
    if [ "$remote_exists" -eq 0 ]; then
        echo "  - Exists on remote: ❌ (deleted on remote)"
    else
        echo "  - Exists on remote: ✅"
    fi
    
    # Check if we can delete it safely
    echo "  - Testing git branch -d:"
    if git branch -d bugfix-incorrect-help-text 2>&1; then
        echo "    ✅ Safe delete works"
    else
        echo "    ❌ Safe delete fails"
    fi
    
    # Check if we can force delete it
    echo "  - Testing git branch -D:"
    if git branch -D bugfix-incorrect-help-text 2>&1; then
        echo "    ✅ Force delete works"
        # Recreate the branch for testing
        git checkout -b bugfix-incorrect-help-text HEAD~1 2>/dev/null || true
    else
        echo "    ❌ Force delete fails"
    fi
    
else
    echo "bugfix-incorrect-help-text: ❌ Branch does not exist"
fi
echo ""

# Check test-suite
if git show-ref --verify --quiet refs/heads/test-suite 2>/dev/null; then
    echo "test-suite:"
    echo "  - Exists locally: ✅"
    
    # Check if it's merged
    if git branch --merged master | grep -q "test-suite"; then
        echo "  - Merged into master: ✅"
    else
        echo "  - Merged into master: ❌"
    fi
    
    # Check if it exists remotely
    remote_exists=$(git ls-remote --heads origin "test-suite" 2>/dev/null | wc -l)
    if [ "$remote_exists" -eq 0 ]; then
        echo "  - Exists on remote: ❌ (deleted on remote)"
    else
        echo "  - Exists on remote: ✅"
    fi
    
    # Check if we can delete it safely
    echo "  - Testing git branch -d:"
    if git branch -d test-suite 2>&1; then
        echo "    ✅ Safe delete works"
    else
        echo "    ❌ Safe delete fails"
    fi
    
    # Check if we can force delete it
    echo "  - Testing git branch -D:"
    if git branch -D test-suite 2>&1; then
        echo "    ✅ Force delete works"
        # Recreate the branch for testing
        git checkout -b test-suite HEAD~1 2>/dev/null || true
    else
        echo "    ❌ Force delete fails"
    fi
    
else
    echo "test-suite: ❌ Branch does not exist"
fi
echo ""

echo "💡 Analysis:"
echo "If branches are detected as 'fully merged' but force delete fails,"
echo "it might be because they're the current branch or have special protection."
echo ""
echo "Try switching to a different branch first:"
echo "  git checkout master"
echo "  git @ sweep" 