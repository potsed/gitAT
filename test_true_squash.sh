#!/bin/bash

echo "Testing true Git squashing functionality..."

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In a git repository"

# Get current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "Current branch: $CURRENT_BRANCH"

# Check if we have multiple commits to squash
COMMIT_COUNT=$(git rev-list --count HEAD 2>/dev/null || echo "0")
echo "Total commits in current branch: $COMMIT_COUNT"

# Test the squash command
echo ""
echo "Testing git @ squash with auto-detection:"
git @ squash

# Check the result
echo ""
echo "Checking result:"
NEW_COMMIT_COUNT=$(git rev-list --count HEAD 2>/dev/null || echo "0")
echo "Commits after squash: $NEW_COMMIT_COUNT"

if [ "$NEW_COMMIT_COUNT" -lt "$COMMIT_COUNT" ]; then
    echo "✅ Squashing successful - reduced from $COMMIT_COUNT to $NEW_COMMIT_COUNT commits"
else
    echo "⚠️  No squashing occurred (may be expected if only one commit)"
fi

# Show the commit history
echo ""
echo "Recent commit history:"
git log --oneline -5

echo ""
echo "True squashing test completed!" 