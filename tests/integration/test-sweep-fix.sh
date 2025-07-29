#!/bin/bash

echo "🧪 Testing Sweep Fix for Remote-Deleted Branches"
echo "================================================"
echo ""

# Check if we're in a git repo
if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In git repository"
echo ""

# Check current branch
current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "Current branch: $current_branch"
echo ""

# Check if remote is configured
if git remote get-url origin >/dev/null 2>&1; then
    echo "✅ Remote origin configured"
else
    echo "❌ No remote origin configured"
    exit 1
fi

echo ""

# Show all branches
echo "📋 All local branches:"
git branch 2>/dev/null || echo "Could not list branches"
echo ""

# Prune remote tracking branches
echo "🧹 Pruning remote tracking branches..."
git remote prune origin 2>/dev/null || true
echo ""

# Check which branches exist locally but not remotely
echo "🔍 Checking for branches that exist locally but not remotely:"
while IFS= read -r branch; do
    if [ -n "$branch" ]; then
        # Check if this branch exists remotely
        remote_exists=$(git ls-remote --heads origin "$branch" 2>/dev/null | wc -l)
        if [ "$remote_exists" -eq 0 ]; then
            echo "   - $branch (exists locally, not remotely)"
        else
            echo "   - $branch (exists both locally and remotely)"
        fi
    fi
done < <(git branch | grep -v "\*" | sed 's/^[[:space:]]*//')
echo ""

# Test the sweep command
echo "🧪 Testing sweep command..."
echo ""

echo "1️⃣ Dry run to see what would be deleted:"
echo "   git @ sweep --dry-run"
echo ""

echo "2️⃣ Actual sweep:"
echo "   git @ sweep"
echo ""

echo "💡 The sweep command should now properly detect and delete branches that:"
echo "   • Are fully merged into trunk"
echo "   • Exist locally but not remotely (remote-deleted branches)"
echo ""
echo "🎯 This should fix the issue with bugfix-incorrect-help-text not being deleted!" 