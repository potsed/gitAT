#!/bin/bash

echo "Testing git state..."

# Check current branch
echo "Current branch: $(git rev-parse --abbrev-ref HEAD)"

# Check if working directory is clean
echo ""
echo "Working directory status:"
if git diff --quiet && git diff --cached --quiet; then
    echo "✅ Working directory is clean"
else
    echo "❌ Working directory has uncommitted changes"
    echo "Staged changes:"
    git diff --cached --name-only
    echo "Unstaged changes:"
    git diff --name-only
fi

# Check git status
echo ""
echo "Git status:"
git status --porcelain

# Test branch creation
echo ""
echo "Testing branch creation..."
TEST_BRANCH="test-branch-$(date +%s)"
echo "Creating test branch: $TEST_BRANCH"

if git checkout -b "$TEST_BRANCH" 2>&1; then
    echo "✅ Test branch creation successful"
    git checkout - 2>/dev/null  # Switch back
    git branch -D "$TEST_BRANCH" 2>/dev/null  # Clean up
else
    echo "❌ Test branch creation failed"
fi

# Check if we can checkout master
echo ""
echo "Testing checkout to master..."
if git checkout master 2>&1; then
    echo "✅ Can checkout to master"
    git checkout - 2>/dev/null  # Switch back
else
    echo "❌ Cannot checkout to master"
fi 