#!/bin/bash

echo "Testing git @ hash command fixes..."

# Test 1: Check if we're in a git repository
echo "Test 1: Checking if we're in a git repository"
if git rev-parse --git-dir > /dev/null 2>&1; then
    echo "✅ In a git repository"
else
    echo "❌ Not in a git repository"
    exit 1
fi

# Test 2: Check current branch
echo "Test 2: Checking current branch"
BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
if [ -n "$BRANCH" ] && [ "$BRANCH" != "HEAD" ]; then
    echo "✅ Current branch: $BRANCH"
else
    echo "❌ Could not determine current branch or in detached HEAD"
    exit 1
fi

# Test 3: Test git @ branch -c
echo "Test 3: Testing git @ branch -c"
if git @ branch -c > /dev/null 2>&1; then
    echo "✅ git @ branch -c works"
else
    echo "❌ git @ branch -c failed"
fi

# Test 4: Test git @ hash
echo "Test 4: Testing git @ hash"
if timeout 10s git @ hash > /dev/null 2>&1; then
    echo "✅ git @ hash works"
else
    echo "❌ git @ hash failed or timed out"
fi

echo "Tests completed!" 