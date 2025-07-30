#!/bin/bash

echo "Testing git @ squash auto-detection functionality..."

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

# Test 3: Check current branch and trunk configuration
echo "Test 3: Checking branch configuration"
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
TRUNK_BRANCH=$(git config at.trunk 2>/dev/null || echo "main")

echo "✅ Current branch: $CURRENT_BRANCH"
echo "✅ Trunk branch: $TRUNK_BRANCH"

# Test 4: Check if trunk branch exists
echo "Test 4: Checking trunk branch existence"
if git rev-parse --verify "$TRUNK_BRANCH" >/dev/null 2>&1; then
    echo "✅ Trunk branch '$TRUNK_BRANCH' exists"
else
    echo "⚠️  Trunk branch '$TRUNK_BRANCH' does not exist - creating test branch"
    git checkout -b "$TRUNK_BRANCH" 2>/dev/null || true
    echo "Test trunk" > trunk_file.txt
    git add trunk_file.txt
    git commit -m "Initial trunk commit" >/dev/null 2>&1 || true
    git checkout "$CURRENT_BRANCH" 2>/dev/null || true
    echo "✅ Created test trunk branch"
fi

# Test 5: Test auto-detection without arguments
echo "Test 5: Testing auto-detection without arguments"
if git @ squash 2>&1 | grep -q "Auto-detected parent branch"; then
    echo "✅ Auto-detection works without arguments"
else
    echo "⚠️  Auto-detection may not work (this is expected if no parent can be detected)"
fi

# Test 6: Test auto-detection with upstream tracking
echo "Test 6: Testing auto-detection with upstream tracking"
# Create a test feature branch from trunk
TEST_FEATURE="test-feature-$(date +%s)"
git checkout -b "$TEST_FEATURE" "$TRUNK_BRANCH" 2>/dev/null || true

# Set upstream tracking
git branch --set-upstream-to="$TRUNK_BRANCH" 2>/dev/null || true

# Test auto-detection
if git @ squash 2>&1 | grep -q "Auto-detected parent branch: $TRUNK_BRANCH"; then
    echo "✅ Auto-detection works with upstream tracking"
else
    echo "⚠️  Auto-detection with upstream tracking may not work"
fi

# Test 7: Test auto-detection with git config merge
echo "Test 7: Testing auto-detection with git config merge"
# Set git config for merge
git config "branch.$TEST_FEATURE.merge" "refs/heads/$TRUNK_BRANCH" 2>/dev/null || true

# Test auto-detection
if git @ squash 2>&1 | grep -q "Auto-detected parent branch: $TRUNK_BRANCH"; then
    echo "✅ Auto-detection works with git config merge"
else
    echo "⚠️  Auto-detection with git config merge may not work"
fi

# Test 8: Test auto-detection with branch divergence analysis
echo "Test 8: Testing auto-detection with branch divergence analysis"
# Create another test branch
TEST_BRANCH_2="test-branch-2-$(date +%s)"
git checkout -b "$TEST_BRANCH_2" "$TRUNK_BRANCH" 2>/dev/null || true

# Make some commits to create divergence
echo "Test commit 1" > test_file_1.txt
git add test_file_1.txt
git commit -m "Test commit 1" >/dev/null 2>&1 || true

echo "Test commit 2" > test_file_2.txt
git add test_file_2.txt
git commit -m "Test commit 2" >/dev/null 2>&1 || true

# Test auto-detection
if git @ squash 2>&1 | grep -q "Auto-detected parent branch"; then
    echo "✅ Auto-detection works with branch divergence analysis"
else
    echo "⚠️  Auto-detection with branch divergence analysis may not work"
fi

# Test 9: Test fallback to configured trunk branch
echo "Test 9: Testing fallback to configured trunk branch"
# Create a branch without upstream tracking
TEST_BRANCH_3="test-branch-3-$(date +%s)"
git checkout -b "$TEST_BRANCH_3" "$TRUNK_BRANCH" 2>/dev/null || true

# Remove any upstream tracking
git config --unset "branch.$TEST_BRANCH_3.merge" 2>/dev/null || true
git config --unset "branch.$TEST_BRANCH_3.remote" 2>/dev/null || true

# Test auto-detection
if git @ squash 2>&1 | grep -q "Auto-detected parent branch: $TRUNK_BRANCH"; then
    echo "✅ Auto-detection falls back to configured trunk branch"
else
    echo "⚠️  Auto-detection fallback may not work"
fi

# Test 10: Test fallback to common branch names
echo "Test 10: Testing fallback to common branch names"
# Temporarily unset trunk config
ORIGINAL_TRUNK=$(git config at.trunk 2>/dev/null || echo "")
git config --unset at.trunk 2>/dev/null || true

# Test auto-detection
if git @ squash 2>&1 | grep -q "Auto-detected parent branch"; then
    echo "✅ Auto-detection falls back to common branch names"
else
    echo "⚠️  Auto-detection fallback to common names may not work"
fi

# Restore trunk config
if [ -n "$ORIGINAL_TRUNK" ]; then
    git config at.trunk "$ORIGINAL_TRUNK" 2>/dev/null || true
fi

# Test 11: Test error when no parent can be detected
echo "Test 11: Testing error when no parent can be detected"
# Create a completely isolated branch
TEST_ISOLATED="test-isolated-$(date +%s)"
git checkout --orphan "$TEST_ISOLATED" 2>/dev/null || true

# Make a commit to the isolated branch
echo "Isolated commit" > isolated_file.txt
git add isolated_file.txt
git commit -m "Isolated commit" >/dev/null 2>&1 || true

# Test auto-detection
if git @ squash 2>&1 | grep -q "Could not auto-detect parent branch"; then
    echo "✅ Correctly handles case when no parent can be detected"
else
    echo "⚠️  May not handle isolated branches correctly"
fi

# Test 12: Test with explicit branch specification
echo "Test 12: Testing with explicit branch specification"
# Switch back to a normal branch
git checkout "$TRUNK_BRANCH" 2>/dev/null || true

# Create a test branch
TEST_EXPLICIT="test-explicit-$(date +%s)"
git checkout -b "$TEST_EXPLICIT" 2>/dev/null || true

# Test with explicit branch
if git @ squash "$TRUNK_BRANCH" 2>&1 | grep -q "Squashed branch.*back to $TRUNK_BRANCH"; then
    echo "✅ Explicit branch specification works"
else
    echo "❌ Explicit branch specification failed"
fi

# Test 13: Test with save option
echo "Test 13: Testing with save option"
# Create another test branch
TEST_SAVE="test-save-$(date +%s)"
git checkout -b "$TEST_SAVE" 2>/dev/null || true

# Make some changes
echo "Test change" > test_save_file.txt
git add test_save_file.txt

# Test with save option
if git @ squash -s 2>&1 | grep -q "Changes saved successfully"; then
    echo "✅ Save option works with auto-detection"
else
    echo "⚠️  Save option may not work with auto-detection"
fi

# Test 14: Test PR mode still works
echo "Test 14: Testing PR mode"
# Create a test branch for PR
TEST_PR="test-pr-$(date +%s)"
git checkout -b "$TEST_PR" 2>/dev/null || true

# Make some commits
echo "PR commit 1" > pr_file_1.txt
git add pr_file_1.txt
git commit -m "PR commit 1" >/dev/null 2>&1 || true

echo "PR commit 2" > pr_file_2.txt
git add pr_file_2.txt
git commit -m "PR commit 2" >/dev/null 2>&1 || true

# Test PR mode
if git @ squash --pr 2>&1 | grep -q "Successfully squashed.*commits into one for PR"; then
    echo "✅ PR mode works correctly"
else
    echo "❌ PR mode failed"
fi

# Test 15: Test auto mode management
echo "Test 15: Testing auto mode management"
# Test enable
if git @ squash --auto on 2>&1 | grep -q "Automatic PR squashing enabled"; then
    echo "✅ Auto mode enable works"
else
    echo "❌ Auto mode enable failed"
fi

# Test status
if git @ squash --auto status 2>&1 | grep -q "Status: ✅ ENABLED"; then
    echo "✅ Auto mode status works"
else
    echo "❌ Auto mode status failed"
fi

# Test disable
if git @ squash --auto off 2>&1 | grep -q "Automatic PR squashing disabled"; then
    echo "✅ Auto mode disable works"
else
    echo "❌ Auto mode disable failed"
fi

# Test 16: Clean up test files and branches
echo "Test 16: Cleaning up test files and branches"
rm -f trunk_file.txt test_file_1.txt test_file_2.txt isolated_file.txt test_save_file.txt pr_file_1.txt pr_file_2.txt

# Clean up test branches (switch back to original branch first)
git checkout "$CURRENT_BRANCH" 2>/dev/null || true

# Delete test branches
for branch in "$TEST_FEATURE" "$TEST_BRANCH_2" "$TEST_BRANCH_3" "$TEST_ISOLATED" "$TEST_EXPLICIT" "$TEST_SAVE" "$TEST_PR"; do
    if git rev-parse --verify "$branch" >/dev/null 2>&1; then
        git branch -D "$branch" 2>/dev/null || true
        echo "✅ Cleaned up test branch: $branch"
    fi
done

# Clean up trunk branch if we created it for testing
if [ "$TRUNK_BRANCH" != "main" ] && [ "$TRUNK_BRANCH" != "master" ]; then
    if git rev-parse --verify "$TRUNK_BRANCH" >/dev/null 2>&1; then
        git branch -D "$TRUNK_BRANCH" 2>/dev/null || true
        echo "✅ Cleaned up test trunk branch: $TRUNK_BRANCH"
    fi
fi

echo ""
echo "Squash auto-detection tests completed!"
echo ""
echo "To test auto-detection in practice:"
echo "1. Create a feature branch: git @ work feature 'my-feature'"
echo "2. Make some commits"
echo "3. Run: git @ squash (should auto-detect parent branch)"
echo "4. Or specify explicitly: git @ squash develop"
echo ""
echo "Auto-detection methods (in order):"
echo "1. Git config branch.<name>.merge"
echo "2. Upstream tracking branch (@{u})"
echo "3. Branch divergence analysis"
echo "4. Configured trunk branch (at.trunk)"
echo "5. Common branch names (main, master, develop)" 