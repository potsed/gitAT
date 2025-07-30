#!/bin/bash

echo "Testing committer extraction functionality..."

# Test 1: Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

# Test 2: Get the latest commit hash
LATEST_HASH=$(git log -1 --pretty=format:"%H" 2>/dev/null || echo "")
if [ -z "$LATEST_HASH" ]; then
    echo "❌ Could not get latest commit hash"
    exit 1
fi

echo "✅ Latest commit hash: ${LATEST_HASH:0:7}"

# Test 3: Extract committer name
COMMITTER=$(git log -1 --pretty=format:"%an" "$LATEST_HASH" 2>/dev/null || echo "Unknown")
echo "✅ Committer: $COMMITTER"

# Test 4: Test the exact command used in hash.sh
echo "Test 4: Testing committer extraction command"
TEST_HASH=$(git log -1 --oneline | cut -d' ' -f1)
TEST_COMMITTER=$(git log -1 --pretty=format:"%an" "$TEST_HASH" 2>/dev/null || echo "Unknown")
echo "✅ Test hash: $TEST_HASH"
echo "✅ Test committer: $TEST_COMMITTER"

# Test 5: Test with short hash
SHORT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "")
SHORT_COMMITTER=$(git log -1 --pretty=format:"%an" "$SHORT_HASH" 2>/dev/null || echo "Unknown")
echo "✅ Short hash: $SHORT_HASH"
echo "✅ Short hash committer: $SHORT_COMMITTER"

# Test 6: Test error handling
echo "Test 6: Testing error handling"
INVALID_COMMITTER=$(git log -1 --pretty=format:"%an" "invalid-hash" 2>/dev/null || echo "Unknown")
echo "✅ Invalid hash committer: $INVALID_COMMITTER"

echo "Committer extraction tests completed!" 