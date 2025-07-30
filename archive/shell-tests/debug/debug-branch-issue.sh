#!/bin/bash

echo "🔍 Debugging Branch Deletion Issue"
echo "=================================="
echo ""

# Check current branch
current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "Current branch: $current_branch"
echo ""

# Check if the branch exists
if git show-ref --verify --quiet refs/heads/bugfix-incorrect-help-text 2>/dev/null; then
    echo "✅ bugfix-incorrect-help-text exists"
else
    echo "❌ bugfix-incorrect-help-text does not exist"
    exit 1
fi

echo ""

# Check working directory status
echo "📋 Working directory status:"
git status --porcelain
echo ""

# Check if there are any Git hooks that might prevent deletion
echo "🔧 Checking for Git hooks:"
if [ -d ".git/hooks" ]; then
    echo "Git hooks directory exists"
    ls -la .git/hooks/ 2>/dev/null || echo "Could not list hooks"
else
    echo "No Git hooks directory"
fi
echo ""

# Check Git configuration for any protection settings
echo "⚙️  Checking Git configuration:"
echo "core.protectNTFS: $(git config core.protectNTFS 2>/dev/null || echo 'not set')"
echo "core.protectHFS: $(git config core.protectHFS 2>/dev/null || echo 'not set')"
echo "branch.bugfix-incorrect-help-text.protect: $(git config branch.bugfix-incorrect-help-text.protect 2>/dev/null || echo 'not set')"
echo ""

# Try to delete the branch with verbose output
echo "🗑️  Attempting to delete branch:"
echo "Command: git branch -D bugfix-incorrect-help-text"
echo "Output:"
git branch -D bugfix-incorrect-help-text 2>&1
exit_code=$?

echo ""
echo "Exit code: $exit_code"

if [ $exit_code -eq 0 ]; then
    echo "✅ Branch deleted successfully"
else
    echo "❌ Branch deletion failed"
    echo ""
    echo "💡 Let's try a different approach..."
    echo ""
    
    # Try to understand what's preventing deletion
    echo "🔍 Additional diagnostics:"
    
    # Check if the branch is referenced anywhere
    echo "Checking branch references:"
    git for-each-ref --format='%(refname)' refs/heads/ | grep -i bugfix || echo "No bugfix branches found"
    
    # Check if there are any reflog entries
    echo "Checking reflog:"
    git reflog --oneline | head -5
    
    # Try to understand the branch structure
    echo "Branch structure:"
    git log --oneline --graph bugfix-incorrect-help-text | head -5
fi 