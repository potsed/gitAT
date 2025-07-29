#!/bin/bash

echo "Debugging squash issue..."

# Check current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "Current branch: $CURRENT_BRANCH"

# Check if master branch exists
echo ""
echo "Checking master branch:"
if git rev-parse --verify master >/dev/null 2>&1; then
    echo "✅ Master branch exists"
    MASTER_SHA=$(git rev-parse master 2>/dev/null || echo "")
    echo "Master SHA: $MASTER_SHA"
else
    echo "❌ Master branch does not exist"
fi

# Check all local branches
echo ""
echo "All local branches:"
git branch --list

# Test the head function
echo ""
echo "Testing head function:"
source git_at_cmds/squash.sh
MASTER_HEAD=$(head "master")
echo "head('master') = $MASTER_HEAD"

# Test squash function directly
echo ""
echo "Testing squash function:"
if [ "$MASTER_HEAD" != "0" ]; then
    echo "Calling squash with SHA: $MASTER_HEAD"
    squash "$MASTER_HEAD"
else
    echo "Cannot test squash - master branch not found"
fi

# Check git status
echo ""
echo "Git status:"
git status --porcelain

# Check if we're in a clean state
echo ""
echo "Checking if working directory is clean:"
if git diff --quiet && git diff --cached --quiet; then
    echo "✅ Working directory is clean"
else
    echo "⚠️  Working directory has uncommitted changes"
fi 