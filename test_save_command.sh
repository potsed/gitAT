#!/bin/bash

echo "Testing git @ save command..."

# Test 1: Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In a git repository"

# Test 2: Check current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "✅ Current branch: $CURRENT_BRANCH"

# Test 3: Check if we have uncommitted changes
if git diff --quiet && git diff --cached --quiet; then
    echo "⚠️  No uncommitted changes to save"
    echo "Creating a test file to test save functionality..."
    echo "Test file for save command" > test_save_file.txt
fi

# Test 4: Test save command with message
echo "Test 4: Testing save command with message"
if git @ save "Test commit message"; then
    echo "✅ Save command worked successfully"
else
    echo "❌ Save command failed"
    exit 1
fi

# Test 5: Verify commit was created
echo "Test 5: Verifying commit was created"
LATEST_COMMIT=$(git log -1 --oneline 2>/dev/null || echo "")
if [[ "$LATEST_COMMIT" == *"Test commit message"* ]]; then
    echo "✅ Commit created successfully: $LATEST_COMMIT"
else
    echo "❌ Commit not found or message doesn't match"
    exit 1
fi

# Test 6: Clean up test file
echo "Test 6: Cleaning up"
rm -f test_save_file.txt
git add test_save_file.txt 2>/dev/null || true
git commit -m "Remove test file" 2>/dev/null || true

echo "All save command tests passed!" 