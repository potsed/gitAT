#!/bin/bash

echo "🧪 Testing Branch Deletion"
echo "=========================="
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

# Try to delete the branch and capture the error
echo "🔍 Testing git branch -D bugfix-incorrect-help-text:"
echo "Command: git branch -D bugfix-incorrect-help-text"
echo "Output:"
git branch -D bugfix-incorrect-help-text
exit_code=$?

echo ""
echo "Exit code: $exit_code"

if [ $exit_code -eq 0 ]; then
    echo "✅ Branch deleted successfully"
else
    echo "❌ Branch deletion failed"
    echo ""
    echo "💡 Possible reasons:"
    echo "1. Branch is the current branch"
    echo "2. Branch has uncommitted changes"
    echo "3. Branch is protected by Git configuration"
    echo "4. Branch is being used by another process"
fi 