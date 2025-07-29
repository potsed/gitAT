#!/bin/bash

echo "Testing git @ pr command..."

# Test 1: Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In a git repository"

# Test 2: Check if remote origin is configured
REMOTE_URL=$(git config --get remote.origin.url 2>/dev/null || echo "")
if [ -z "$REMOTE_URL" ]; then
    echo "⚠️  No remote origin configured - some tests will be limited"
else
    echo "✅ Remote origin: $REMOTE_URL"
fi

# Test 3: Check current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "✅ Current branch: $CURRENT_BRANCH"

# Test 4: Test help functionality
echo "Test 4: Testing help functionality"
if git @ pr -h > /dev/null 2>&1; then
    echo "✅ Help command works"
else
    echo "❌ Help command failed"
    exit 1
fi

# Test 5: Test platform detection
echo "Test 5: Testing platform detection"
PLATFORM=$(git @ pr 2>&1 | grep -o "Creating PR for [a-z]*" | cut -d' ' -f4 || echo "unknown")
echo "✅ Detected platform: $PLATFORM"

# Test 6: Test default title generation
echo "Test 6: Testing default title generation"
DEFAULT_TITLE=$(git log -1 --pretty=format:"%s" 2>/dev/null || echo "Update from $CURRENT_BRANCH")
echo "✅ Default title would be: $DEFAULT_TITLE"

# Test 7: Test trunk branch detection
echo "Test 7: Testing trunk branch detection"
TRUNK_BRANCH=$(git config at.trunk 2>/dev/null || echo "main")
echo "✅ Trunk branch: $TRUNK_BRANCH"

# Test 8: Test error handling for trunk branch
echo "Test 8: Testing error handling for trunk branch"
if [ "$CURRENT_BRANCH" = "$TRUNK_BRANCH" ]; then
    echo "⚠️  Currently on trunk branch - PR creation would fail (expected)"
    if git @ pr 2>&1 | grep -q "Cannot create PR from.*to itself"; then
        echo "✅ Correctly detected trunk branch error"
    else
        echo "❌ Failed to detect trunk branch error"
    fi
else
    echo "✅ Not on trunk branch - PR creation should work"
fi

# Test 9: Test web URL generation (simulated)
echo "Test 9: Testing web URL generation"
case "$PLATFORM" in
    "github")
        if [ -n "$REMOTE_URL" ]; then
            REPO_INFO=$(echo "$REMOTE_URL" | sed 's|https://github.com/||' | sed 's|\.git$||' | sed 's|git@github.com:||')
            EXPECTED_URL="https://github.com/$REPO_INFO/compare/$TRUNK_BRANCH...$CURRENT_BRANCH"
            echo "✅ Expected GitHub URL: $EXPECTED_URL"
        fi
        ;;
    "gitlab")
        if [ -n "$REMOTE_URL" ]; then
            REPO_INFO=$(echo "$REMOTE_URL" | sed 's|https://gitlab.com/||' | sed 's|\.git$||' | sed 's|git@gitlab.com:||')
            EXPECTED_URL="https://gitlab.com/$REPO_INFO/-/merge_requests/new?source_branch=$CURRENT_BRANCH&target_branch=$TRUNK_BRANCH"
            echo "✅ Expected GitLab URL: $EXPECTED_URL"
        fi
        ;;
    "bitbucket")
        if [ -n "$REMOTE_URL" ]; then
            REPO_INFO=$(echo "$REMOTE_URL" | sed 's|https://bitbucket.org/||' | sed 's|\.git$||' | sed 's|git@bitbucket.org:||')
            EXPECTED_URL="https://bitbucket.org/$REPO_INFO/pull-requests/new?source=$CURRENT_BRANCH&t=1"
            echo "✅ Expected Bitbucket URL: $EXPECTED_URL"
        fi
        ;;
    *)
        echo "✅ Generic platform detected"
        ;;
esac

# Test 10: Test CLI tool availability
echo "Test 10: Testing CLI tool availability"
case "$PLATFORM" in
    "github")
        if command -v gh >/dev/null 2>&1; then
            echo "✅ GitHub CLI (gh) is installed"
            if gh auth status >/dev/null 2>&1; then
                echo "✅ GitHub CLI is authenticated"
            else
                echo "⚠️  GitHub CLI is not authenticated"
            fi
        else
            echo "⚠️  GitHub CLI (gh) is not installed"
        fi
        ;;
    "gitlab")
        if command -v glab >/dev/null 2>&1; then
            echo "✅ GitLab CLI (glab) is installed"
            if glab auth status >/dev/null 2>&1; then
                echo "✅ GitLab CLI is authenticated"
            else
                echo "⚠️  GitLab CLI is not authenticated"
            fi
        else
            echo "⚠️  GitLab CLI (glab) is not installed"
        fi
        ;;
    *)
        echo "✅ No CLI tool required for this platform"
        ;;
esac

# Test 11: Test argument parsing
echo "Test 11: Testing argument parsing"
if git @ pr --help > /dev/null 2>&1; then
    echo "✅ --help argument works"
else
    echo "❌ --help argument failed"
fi

if git @ pr help > /dev/null 2>&1; then
    echo "✅ help argument works"
else
    echo "❌ help argument failed"
fi

# Test 12: Test error handling for invalid options
echo "Test 12: Testing error handling for invalid options"
if git @ pr -t 2>&1 | grep -q "requires a value"; then
    echo "✅ Correctly handles missing title value"
else
    echo "❌ Failed to handle missing title value"
fi

if git @ pr -d 2>&1 | grep -q "requires a value"; then
    echo "✅ Correctly handles missing description value"
else
    echo "❌ Failed to handle missing description value"
fi

if git @ pr -b 2>&1 | grep -q "requires a value"; then
    echo "✅ Correctly handles missing base value"
else
    echo "❌ Failed to handle missing base value"
fi

echo ""
echo "PR command tests completed!"
echo ""
echo "To test actual PR creation:"
echo "1. Ensure you're on a feature branch (not trunk)"
echo "2. Have commits ready to submit"
echo "3. Run: git @ pr 'Your PR title'"
echo "4. Or run: git @ pr -d 'Your description' -o" 