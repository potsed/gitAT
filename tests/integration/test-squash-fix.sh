#!/bin/bash

echo "Testing squash auto-detection fix..."

# Check current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "Current branch: $CURRENT_BRANCH"

# Test the detect_parent_branch function
echo ""
echo "Testing detect_parent_branch function:"
source git_at_cmds/squash.sh
DETECTED_PARENT=$(detect_parent_branch)
echo "Detected parent: $DETECTED_PARENT"

# Show all branches for context
echo ""
echo "All local branches:"
git branch --list

# Test the actual squash command
echo ""
echo "Testing git @ squash (should auto-detect parent):"
git @ squash 