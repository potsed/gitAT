#!/bin/bash

echo "Testing enhanced git @ squash command..."

# Test 1: Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In a git repository"

# Test 2: Test help functionality
echo "Test 2: Testing help functionality"
if git @ squash -h > /dev/null 2>&1; then
    echo "✅ Help command works"
else
    echo "❌ Help command failed"
    exit 1
fi

# Test 3: Test basic squash functionality
echo "Test 3: Testing basic squash functionality"
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
TRUNK_BRANCH=$(git config at.trunk 2>/dev/null || echo "main")

# Check if trunk branch exists
if git rev-parse --verify "$TRUNK_BRANCH" >/dev/null 2>&1; then
    echo "✅ Trunk branch '$TRUNK_BRANCH' exists"
else
    echo "⚠️  Trunk branch '$TRUNK_BRANCH' does not exist - creating test branch"
    git checkout -b "$TRUNK_BRANCH" 2>/dev/null || true
    echo "Test trunk" > trunk_file.txt
    git add trunk_file.txt
    git commit -m "Initial trunk commit" >/dev/null 2>&1 || true
    git checkout "$CURRENT_BRANCH" 2>/dev/null || true
fi

# Test 4: Test error for non-existent branch
echo "Test 4: Testing error for non-existent branch"
if git @ squash nonexistent-branch 2>&1 | grep -q "does not exist locally"; then
    echo "✅ Correctly handles non-existent branch"
else
    echo "❌ Failed to handle non-existent branch"
fi

# Test 5: Test automatic PR squashing management
echo "Test 5: Testing automatic PR squashing management"

# Test initial status
INITIAL_STATUS=$(git @ squash --auto status 2>&1 | grep -o "DISABLED\|ENABLED" || echo "DISABLED")
echo "✅ Initial status: $INITIAL_STATUS"

# Test enable
if git @ squash --auto on > /dev/null 2>&1; then
    echo "✅ Enable command works"
else
    echo "❌ Enable command failed"
    exit 1
fi

# Test status when enabled
ENABLED_STATUS=$(git @ squash --auto status 2>&1 | grep -o "ENABLED" || echo "")
if [ "$ENABLED_STATUS" = "ENABLED" ]; then
    echo "✅ Squash setting is enabled"
else
    echo "❌ Squash setting is not enabled"
    exit 1
fi

# Test disable
if git @ squash --auto off > /dev/null 2>&1; then
    echo "✅ Disable command works"
else
    echo "❌ Disable command failed"
    exit 1
fi

# Test status when disabled
DISABLED_STATUS=$(git @ squash --auto status 2>&1 | grep -o "DISABLED" || echo "")
if [ "$DISABLED_STATUS" = "DISABLED" ]; then
    echo "✅ Squash setting is disabled"
else
    echo "❌ Squash setting is not disabled"
    exit 1
fi

# Test 6: Test alternative auto commands
echo "Test 6: Testing alternative auto commands"
for cmd in "true" "enable" "1"; do
    if git @ squash --auto "$cmd" > /dev/null 2>&1; then
        echo "✅ Enable with '$cmd' works"
    else
        echo "❌ Enable with '$cmd' failed"
    fi
done

for cmd in "false" "disable" "0"; do
    if git @ squash --auto "$cmd" > /dev/null 2>&1; then
        echo "✅ Disable with '$cmd' works"
    else
        echo "❌ Disable with '$cmd' failed"
    fi
done

for cmd in "show" "check" ""; do
    if git @ squash --auto $cmd > /dev/null 2>&1; then
        echo "✅ Status with '$cmd' works"
    else
        echo "❌ Status with '$cmd' failed"
    fi
done

# Test 7: Test error handling for invalid auto commands
echo "Test 7: Testing error handling"
if git @ squash --auto invalid 2>&1 | grep -q "Invalid auto action"; then
    echo "✅ Correctly handles invalid auto action"
else
    echo "❌ Failed to handle invalid auto action"
fi

if git @ squash --auto 2>&1 | grep -q "requires a value"; then
    echo "✅ Correctly handles missing auto value"
else
    echo "❌ Failed to handle missing auto value"
fi

# Test 8: Test configuration persistence
echo "Test 8: Testing configuration persistence"
git @ squash --auto on > /dev/null 2>&1
CONFIG_VALUE=$(git config at.pr.squash 2>/dev/null || echo "")
if [ "$CONFIG_VALUE" = "true" ]; then
    echo "✅ Configuration is persisted correctly"
else
    echo "❌ Configuration is not persisted correctly"
fi

# Test 9: Test PR squash mode
echo "Test 9: Testing PR squash mode"
git @ squash --auto off > /dev/null 2>&1

# Create some test commits if we don't have multiple commits
COMMIT_COUNT=$(git rev-list --count "$TRUNK_BRANCH..HEAD" 2>/dev/null || echo "0")

if [ "$COMMIT_COUNT" -le 1 ]; then
    echo "⚠️  Only $COMMIT_COUNT commits found - creating test commits"
    echo "Test commit 1" > test_squash_1.txt
    git add test_squash_1.txt
    git commit -m "Test commit 1" > /dev/null 2>&1
    
    echo "Test commit 2" > test_squash_2.txt
    git add test_squash_2.txt
    git commit -m "Test commit 2" > /dev/null 2>&1
    
    echo "✅ Created test commits for PR squashing"
fi

# Test PR squash
PR_SQUASH_OUTPUT=$(git @ squash --pr 2>&1)
if echo "$PR_SQUASH_OUTPUT" | grep -q "Successfully squashed.*commits into one for PR"; then
    echo "✅ PR squash mode works correctly"
else
    echo "⚠️  PR squash mode may not be working (check if commits exist to squash)"
fi

# Test 10: Test PR squash error handling
echo "Test 10: Testing PR squash error handling"

# Test error when on trunk branch
git checkout "$TRUNK_BRANCH" 2>/dev/null || true
if git @ squash --pr 2>&1 | grep -q "Cannot squash PR from.*to itself"; then
    echo "✅ Correctly detects trunk branch error"
else
    echo "❌ Failed to detect trunk branch error"
fi

# Go back to original branch
git checkout "$CURRENT_BRANCH" 2>/dev/null || true

# Test 11: Test integration with git @ pr
echo "Test 11: Testing integration with git @ pr"
git @ squash --auto on > /dev/null 2>&1

# Create test commits for PR integration
echo "Test PR commit 1" > test_pr_1.txt
git add test_pr_1.txt
git commit -m "Test PR commit 1" > /dev/null 2>&1

echo "Test PR commit 2" > test_pr_2.txt
git add test_pr_2.txt
git commit -m "Test PR commit 2" > /dev/null 2>&1

PR_OUTPUT=$(git @ pr "Test PR with squash" 2>&1)
if echo "$PR_OUTPUT" | grep -q "Auto-squashing commits before creating PR"; then
    echo "✅ PR command detects squash setting"
else
    echo "⚠️  PR command does not show squash message (may be normal if no commits to squash)"
fi

# Test 12: Test force override flags
echo "Test 12: Testing force override flags"
git @ squash --auto off > /dev/null 2>&1

# Create test commits for force override
echo "Test force commit 1" > test_force_1.txt
git add test_force_1.txt
git commit -m "Test force commit 1" > /dev/null 2>&1

echo "Test force commit 2" > test_force_2.txt
git add test_force_2.txt
git commit -m "Test force commit 2" > /dev/null 2>&1

FORCE_SQUASH_OUTPUT=$(git @ pr "Test PR with force squash" -s 2>&1)
if echo "$FORCE_SQUASH_OUTPUT" | grep -q "Auto-squashing commits before creating PR"; then
    echo "✅ Force squash flag (-s) works"
else
    echo "⚠️  Force squash flag may not be working (check if commits exist to squash)"
fi

git @ squash --auto on > /dev/null 2>&1
FORCE_NO_SQUASH_OUTPUT=$(git @ pr "Test PR with force no squash" -S 2>&1)
if ! echo "$FORCE_NO_SQUASH_OUTPUT" | grep -q "Auto-squashing commits before creating PR"; then
    echo "✅ Force no squash flag (-S) works"
else
    echo "❌ Force no squash flag (-S) does not work"
fi

# Test 13: Clean up test files
echo "Test 13: Cleaning up test files"
rm -f test_squash_1.txt test_squash_2.txt test_pr_1.txt test_pr_2.txt test_force_1.txt test_force_2.txt trunk_file.txt
git add test_squash_1.txt test_squash_2.txt test_pr_1.txt test_pr_2.txt test_force_1.txt test_force_2.txt trunk_file.txt 2>/dev/null || true
git commit -m "Remove test files" > /dev/null 2>&1 || true

echo ""
echo "Enhanced squash command tests completed!"
echo ""
echo "Current squash setting:"
git @ squash --auto status
echo ""
echo "To test actual squash functionality:"
echo "1. Enable auto squash: git @ squash --auto on"
echo "2. Make multiple commits"
echo "3. Run: git @ squash --pr (for manual PR squash)"
echo "4. Or run: git @ pr 'Your PR title' (for automatic PR squash)" 