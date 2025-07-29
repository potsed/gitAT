#!/bin/bash

echo "Testing git @ work command (Conventional Commits integration)..."

# Test 1: Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In a git repository"

# Test 2: Test help functionality
echo "Test 2: Testing help functionality"
if git @ work -h > /dev/null 2>&1; then
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

# Test 5: Test error handling for missing work type
echo "Test 5: Testing missing work type"
if git @ work 2>&1 | grep -q "Work type is required"; then
    echo "✅ Correctly handles missing work type"
else
    echo "❌ Failed to handle missing work type"
fi

# Test 6: Test error handling for invalid work type
echo "Test 6: Testing invalid work type"
if git @ work invalid-type test 2>&1 | grep -q "Invalid work type"; then
    echo "✅ Correctly handles invalid work type"
else
    echo "❌ Failed to handle invalid work type"
fi

# Test 7: Test error handling for missing name option value
echo "Test 7: Testing missing name option value"
if git @ work -n 2>&1 | grep -q "requires a value"; then
    echo "✅ Correctly handles missing -n value"
else
    echo "❌ Failed to handle missing -n value"
fi

if git @ work --name 2>&1 | grep -q "requires a value"; then
    echo "✅ Correctly handles missing --name value"
else
    echo "❌ Failed to handle missing --name value"
fi

# Test 8: Test error handling for too many arguments
echo "Test 8: Testing too many arguments"
if git @ work feature add-auth extra-arg 2>&1 | grep -q "Too many arguments"; then
    echo "✅ Correctly handles too many arguments"
else
    echo "❌ Failed to handle too many arguments"
fi

# Test 9: Test feature branch creation
echo "Test 9: Testing feature branch creation"
FEATURE_NAME="test-feature-$(date +%s)"
if git @ work feature "$FEATURE_NAME" 2>&1 | grep -q "feature branch.*created successfully"; then
    echo "✅ Feature branch creation works"
    
    # Check if we're on the feature branch
    NEW_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    if [ "$NEW_BRANCH" = "feature-$FEATURE_NAME" ]; then
        echo "✅ Successfully switched to feature branch"
    else
        echo "❌ Failed to switch to feature branch"
    fi
    
    # Check if working branch is set
    WORKING_BRANCH=$(git config at.branch 2>/dev/null || echo "")
    if [ "$WORKING_BRANCH" = "feature-$FEATURE_NAME" ]; then
        echo "✅ Working branch set to feature branch"
    else
        echo "❌ Working branch not set to feature branch"
    fi
else
    echo "❌ Feature branch creation failed"
fi

# Test 10: Test hotfix branch creation
echo "Test 10: Testing hotfix branch creation"
HOTFIX_NAME="test-hotfix-$(date +%s)"
if git @ work hotfix "$HOTFIX_NAME" 2>&1 | grep -q "hotfix branch.*created successfully"; then
    echo "✅ Hotfix branch creation works"
    
    # Check if we're on the hotfix branch
    NEW_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    if [ "$NEW_BRANCH" = "hotfix-$HOTFIX_NAME" ]; then
        echo "✅ Successfully switched to hotfix branch"
    else
        echo "❌ Failed to switch to hotfix branch"
    fi
else
    echo "❌ Hotfix branch creation failed"
fi

# Test 11: Test bugfix branch creation
echo "Test 11: Testing bugfix branch creation"
BUGFIX_NAME="test-bugfix-$(date +%s)"
if git @ work bugfix "$BUGFIX_NAME" 2>&1 | grep -q "bugfix branch.*created successfully"; then
    echo "✅ Bugfix branch creation works"
else
    echo "❌ Bugfix branch creation failed"
fi

# Test 12: Test docs branch creation
echo "Test 12: Testing docs branch creation"
DOCS_NAME="test-docs-$(date +%s)"
if git @ work docs "$DOCS_NAME" 2>&1 | grep -q "docs branch.*created successfully"; then
    echo "✅ Docs branch creation works"
else
    echo "❌ Docs branch creation failed"
fi

# Test 13: Test chore branch creation
echo "Test 13: Testing chore branch creation"
CHORE_NAME="test-chore-$(date +%s)"
if git @ work chore "$CHORE_NAME" 2>&1 | grep -q "chore branch.*created successfully"; then
    echo "✅ Chore branch creation works"
else
    echo "❌ Chore branch creation failed"
fi

# Test 14: Test work branch with full name option
echo "Test 14: Testing work branch with full name option"
CUSTOM_NAME="custom-test-branch-$(date +%s)"
if git @ work -n "$CUSTOM_NAME" 2>&1 | grep -q "branch.*created successfully"; then
    echo "✅ Custom branch creation works with -n option"
else
    echo "❌ Custom branch creation failed with -n option"
fi

# Test 15: Test error for existing branch
echo "Test 15: Testing error for existing branch"
if git @ work feature "$FEATURE_NAME" 2>&1 | grep -q "already exists"; then
    echo "✅ Correctly handles existing branch error"
else
    echo "❌ Failed to handle existing branch error"
fi

# Test 16: Test warning for uncommitted changes
echo "Test 16: Testing warning for uncommitted changes"
echo "Test uncommitted change" > test_uncommitted.txt
git add test_uncommitted.txt

# Note: This test would require interactive input, so we'll just check the warning
WORK_WITH_CHANGES="test-work-changes-$(date +%s)"
WORK_OUTPUT=$(git @ work feature "$WORK_WITH_CHANGES" 2>&1)
if echo "$WORK_OUTPUT" | grep -q "uncommitted changes"; then
    echo "✅ Correctly warns about uncommitted changes"
else
    echo "⚠️  No uncommitted changes warning (may be normal if changes were committed)"
fi

# Test 17: Test integration with existing GitAT commands
echo "Test 17: Testing integration with existing GitAT commands"
WORK_INTEGRATION="test-work-integration-$(date +%s)"

# Check if wip command is used
if git @ work feature "$WORK_INTEGRATION" 2>&1 | grep -q "Saving current work state"; then
    echo "✅ Integrates with git @ wip command"
else
    echo "⚠️  WIP integration not visible (may be normal if no changes)"
fi

# Check if branch command is used
if git @ work feature "$WORK_INTEGRATION" 2>&1 | grep -q "Setting working branch"; then
    echo "✅ Integrates with git @ branch command"
else
    echo "❌ Branch integration not working"
fi

# Test 18: Test next steps guidance
echo "Test 18: Testing next steps guidance"
WORK_GUIDANCE="test-work-guidance-$(date +%s)"
GUIDANCE_OUTPUT=$(git @ work feature "$WORK_GUIDANCE" 2>&1)
if echo "$GUIDANCE_OUTPUT" | grep -q "Next steps:"; then
    echo "✅ Shows next steps guidance"
    if echo "$GUIDANCE_OUTPUT" | grep -q "git @ save"; then
        echo "✅ Mentions git @ save in guidance"
    fi
    if echo "$GUIDANCE_OUTPUT" | grep -q "git @ pr"; then
        echo "✅ Mentions git @ pr in guidance"
    fi
    if echo "$GUIDANCE_OUTPUT" | grep -q "git @ release"; then
        echo "✅ Mentions release guidance"
    fi
else
    echo "❌ Does not show next steps guidance"
fi

# Test 19: Test Conventional Commits integration with save command
echo "Test 19: Testing Conventional Commits integration with save command"
# Switch to a feature branch
git checkout "feature-$FEATURE_NAME" 2>/dev/null || true

# Create a test file and save it
echo "Test feature change" > test_feature_change.txt
git add test_feature_change.txt

# Test save command with Conventional Commits prefix
SAVE_OUTPUT=$(git @ save "Add new feature functionality" 2>&1)
if echo "$SAVE_OUTPUT" | grep -q "Changes saved successfully"; then
    echo "✅ Save command works with feature branch"
    
    # Check if commit message has [FEATURE] prefix
    LATEST_COMMIT=$(git log -1 --oneline 2>/dev/null || echo "")
    if echo "$LATEST_COMMIT" | grep -q "\[FEATURE\]"; then
        echo "✅ Commit message includes [FEATURE] prefix"
    else
        echo "❌ Commit message missing [FEATURE] prefix"
    fi
else
    echo "❌ Save command failed with feature branch"
fi

# Test 20: Test all work types
echo "Test 20: Testing all work types"
WORK_TYPES=("hotfix" "feature" "bugfix" "release" "chore" "docs" "style" "refactor" "perf" "test" "ci" "build" "revert")

for work_type in "${WORK_TYPES[@]}"; do
    TEST_NAME="test-${work_type}-$(date +%s)"
    if git @ work "$work_type" "$TEST_NAME" 2>&1 | grep -q "${work_type} branch.*created successfully"; then
        echo "✅ $work_type branch creation works"
    else
        echo "❌ $work_type branch creation failed"
    fi
done

# Test 21: Clean up test files and branches
echo "Test 21: Cleaning up test files and branches"
rm -f test_uncommitted.txt trunk_file.txt test_feature_change.txt

# Clean up test branches (switch back to original branch first)
git checkout "$CURRENT_BRANCH" 2>/dev/null || true

# Delete test branches
for branch in "feature-$FEATURE_NAME" "hotfix-$HOTFIX_NAME" "bugfix-$BUGFIX_NAME" "docs-$DOCS_NAME" "chore-$CHORE_NAME" "$CUSTOM_NAME" "$WORK_WITH_CHANGES" "$WORK_INTEGRATION" "$WORK_GUIDANCE"; do
    if git rev-parse --verify "$branch" >/dev/null 2>&1; then
        git branch -D "$branch" 2>/dev/null || true
        echo "✅ Cleaned up test branch: $branch"
    fi
done

# Clean up work type test branches
for work_type in "${WORK_TYPES[@]}"; do
    for branch in $(git branch --list | grep "test-${work_type}-" | sed 's/^[* ]*//'); do
        if [ -n "$branch" ]; then
            git branch -D "$branch" 2>/dev/null || true
            echo "✅ Cleaned up test branch: $branch"
        fi
    done
done

# Clean up trunk branch if we created it for testing
if [ "$TRUNK_BRANCH" != "main" ] && [ "$TRUNK_BRANCH" != "master" ]; then
    if git rev-parse --verify "$TRUNK_BRANCH" >/dev/null 2>&1; then
        git branch -D "$TRUNK_BRANCH" 2>/dev/null || true
        echo "✅ Cleaned up test trunk branch: $TRUNK_BRANCH"
    fi
fi

echo ""
echo "Work command tests completed!"
echo ""
echo "To test interactive work branch creation:"
echo "1. Run: git @ work feature"
echo "2. Enter a description when prompted"
echo "3. Follow the next steps guidance"
echo ""
echo "To test Conventional Commits workflow:"
echo "1. git @ work feature 'add-user-auth'"
echo "2. Make your changes"
echo "3. git @ save 'Add user authentication system'"
echo "4. git @ pr 'Feature: Add user authentication system'"
echo ""
echo "To test branch listing:"
echo "1. git @ branch --hotfix"
echo "2. git @ branch --feature"
echo "3. git @ branch --all-types" 