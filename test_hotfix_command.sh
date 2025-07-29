#!/bin/bash

echo "Testing git @ hotfix command..."

# Test 1: Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In a git repository"

# Test 2: Test help functionality
echo "Test 2: Testing help functionality"
if git @ hotfix -h > /dev/null 2>&1; then
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

# Test 5: Test error handling for invalid branch names
echo "Test 5: Testing invalid branch name handling"
if git @ hotfix invalid@name 2>&1 | grep -q "Invalid branch name"; then
    echo "✅ Correctly handles invalid branch name"
else
    echo "❌ Failed to handle invalid branch name"
fi

# Test 6: Test error handling for reserved branch names
echo "Test 6: Testing reserved branch name handling"
for reserved_name in "HEAD" "master" "main" "develop"; do
    if git @ hotfix "$reserved_name" 2>&1 | grep -q "Invalid branch name"; then
        echo "✅ Correctly handles reserved name: $reserved_name"
    else
        echo "❌ Failed to handle reserved name: $reserved_name"
    fi
done

# Test 7: Test error handling for missing name option value
echo "Test 7: Testing missing name option value"
if git @ hotfix -n 2>&1 | grep -q "requires a value"; then
    echo "✅ Correctly handles missing -n value"
else
    echo "❌ Failed to handle missing -n value"
fi

if git @ hotfix --name 2>&1 | grep -q "requires a value"; then
    echo "✅ Correctly handles missing --name value"
else
    echo "❌ Failed to handle missing --name value"
fi

# Test 8: Test error handling for unknown options
echo "Test 8: Testing unknown option handling"
if git @ hotfix --unknown 2>&1 | grep -q "Unknown option"; then
    echo "✅ Correctly handles unknown option"
else
    echo "❌ Failed to handle unknown option"
fi

# Test 9: Test hotfix creation with specific name
echo "Test 9: Testing hotfix creation with specific name"
HOTFIX_NAME="test-hotfix-$(date +%s)"
if git @ hotfix "$HOTFIX_NAME" 2>&1 | grep -q "Hotfix branch.*created successfully"; then
    echo "✅ Hotfix creation works with specific name"
    
    # Check if we're on the hotfix branch
    NEW_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    if [ "$NEW_BRANCH" = "$HOTFIX_NAME" ]; then
        echo "✅ Successfully switched to hotfix branch"
    else
        echo "❌ Failed to switch to hotfix branch"
    fi
    
    # Check if working branch is set
    WORKING_BRANCH=$(git config at.branch 2>/dev/null || echo "")
    if [ "$WORKING_BRANCH" = "$HOTFIX_NAME" ]; then
        echo "✅ Working branch set to hotfix branch"
    else
        echo "❌ Working branch not set to hotfix branch"
    fi
else
    echo "❌ Hotfix creation failed"
fi

# Test 10: Test hotfix creation with name option
echo "Test 10: Testing hotfix creation with name option"
HOTFIX_NAME_OPTION="test-hotfix-option-$(date +%s)"
if git @ hotfix -n "$HOTFIX_NAME_OPTION" 2>&1 | grep -q "Hotfix branch.*created successfully"; then
    echo "✅ Hotfix creation works with -n option"
else
    echo "❌ Hotfix creation failed with -n option"
fi

# Test 11: Test hotfix creation with --name option
echo "Test 11: Testing hotfix creation with --name option"
HOTFIX_NAME_LONG="test-hotfix-long-$(date +%s)"
if git @ hotfix --name "$HOTFIX_NAME_LONG" 2>&1 | grep -q "Hotfix branch.*created successfully"; then
    echo "✅ Hotfix creation works with --name option"
else
    echo "❌ Hotfix creation failed with --name option"
fi

# Test 12: Test error for existing branch
echo "Test 12: Testing error for existing branch"
if git @ hotfix "$HOTFIX_NAME" 2>&1 | grep -q "already exists"; then
    echo "✅ Correctly handles existing branch error"
else
    echo "❌ Failed to handle existing branch error"
fi

# Test 13: Test warning for uncommitted changes
echo "Test 13: Testing warning for uncommitted changes"
echo "Test uncommitted change" > test_uncommitted.txt
git add test_uncommitted.txt

# Note: This test would require interactive input, so we'll just check the warning
HOTFIX_WITH_CHANGES="test-hotfix-changes-$(date +%s)"
HOTFIX_OUTPUT=$(git @ hotfix "$HOTFIX_WITH_CHANGES" 2>&1)
if echo "$HOTFIX_OUTPUT" | grep -q "uncommitted changes"; then
    echo "✅ Correctly warns about uncommitted changes"
else
    echo "⚠️  No uncommitted changes warning (may be normal if changes were committed)"
fi

# Test 14: Test integration with existing GitAT commands
echo "Test 14: Testing integration with existing GitAT commands"
HOTFIX_INTEGRATION="test-hotfix-integration-$(date +%s)"

# Check if wip command is used
if git @ hotfix "$HOTFIX_INTEGRATION" 2>&1 | grep -q "Saving current work state"; then
    echo "✅ Integrates with git @ wip command"
else
    echo "⚠️  WIP integration not visible (may be normal if no changes)"
fi

# Check if branch command is used
if git @ hotfix "$HOTFIX_INTEGRATION" 2>&1 | grep -q "Setting working branch"; then
    echo "✅ Integrates with git @ branch command"
else
    echo "❌ Branch integration not working"
fi

# Test 15: Test next steps guidance
echo "Test 15: Testing next steps guidance"
HOTFIX_GUIDANCE="test-hotfix-guidance-$(date +%s)"
GUIDANCE_OUTPUT=$(git @ hotfix "$HOTFIX_GUIDANCE" 2>&1)
if echo "$GUIDANCE_OUTPUT" | grep -q "Next steps:"; then
    echo "✅ Shows next steps guidance"
    if echo "$GUIDANCE_OUTPUT" | grep -q "git @ save"; then
        echo "✅ Mentions git @ save in guidance"
    fi
    if echo "$GUIDANCE_OUTPUT" | grep -q "git @ pr"; then
        echo "✅ Mentions git @ pr in guidance"
    fi
    if echo "$GUIDANCE_OUTPUT" | grep -q "git @ release -p"; then
        echo "✅ Mentions patch release in guidance"
    fi
else
    echo "❌ Does not show next steps guidance"
fi

# Test 16: Clean up test files and branches
echo "Test 16: Cleaning up test files and branches"
rm -f test_uncommitted.txt trunk_file.txt

# Clean up test branches (switch back to original branch first)
git checkout "$CURRENT_BRANCH" 2>/dev/null || true

# Delete test branches
for branch in "$HOTFIX_NAME" "$HOTFIX_NAME_OPTION" "$HOTFIX_NAME_LONG" "$HOTFIX_WITH_CHANGES" "$HOTFIX_INTEGRATION" "$HOTFIX_GUIDANCE"; do
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
echo "Hotfix command tests completed!"
echo ""
echo "To test interactive hotfix creation:"
echo "1. Run: git @ hotfix"
echo "2. Enter a hotfix branch name when prompted"
echo "3. Follow the next steps guidance"
echo ""
echo "To test hotfix workflow:"
echo "1. git @ hotfix 'fix-critical-bug'"
echo "2. Make your fixes"
echo "3. git @ save 'Fix critical bug'"
echo "4. git @ pr 'Hotfix: Fix critical bug'" 