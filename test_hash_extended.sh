#!/bin/bash

echo "Testing extended git @ hash command with committer information..."

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

# Test 3: Check if we have commits
echo "Test 3: Checking commit history"
COMMIT_COUNT=$(git log --oneline | wc -l | tr -d ' ')
if [ "$COMMIT_COUNT" -gt 0 ]; then
    echo "✅ Found $COMMIT_COUNT commits"
else
    echo "⚠️  No commits found"
fi

# Test 4: Test git @ hash with recent commits and committers
echo "Test 4: Testing git @ hash with recent commits and committers"
echo "Output:"
echo "────────────────────────────────────────────────────────────────────────────────"
if timeout 10s git @ hash; then
    echo "────────────────────────────────────────────────────────────────────────────────"
    echo "✅ git @ hash works with recent commits and committers"
else
    echo "❌ git @ hash failed or timed out"
fi

# Test 5: Verify recent commits section
echo "Test 5: Verifying recent commits section"
if git @ hash 2>/dev/null | grep -q "RECENT COMMITS (last 5):"; then
    echo "✅ Recent commits section found"
else
    echo "❌ Recent commits section not found"
fi

# Test 6: Check commit format with committer
echo "Test 6: Checking commit format with committer"
if git @ hash 2>/dev/null | grep -E "^[a-f0-9]{7} │" > /dev/null; then
    echo "✅ Commit format is correct (hash │ committer │ message)"
else
    echo "⚠️  Commit format may be incorrect or no commits to display"
fi

# Test 7: Verify committer information is present
echo "Test 7: Verifying committer information"
if git @ hash 2>/dev/null | grep -E "│ [^│]+ │" > /dev/null; then
    echo "✅ Committer information is present"
else
    echo "⚠️  Committer information may be missing or no commits to display"
fi

# Test 8: Check for three-column format
echo "Test 8: Checking three-column format"
if git @ hash 2>/dev/null | grep -E "^[a-f0-9]{7} │ [^│]+ │ [^│]+" > /dev/null; then
    echo "✅ Three-column format is correct (hash │ committer │ message)"
else
    echo "⚠️  Three-column format may be incorrect or no commits to display"
fi

echo "Tests completed!" 