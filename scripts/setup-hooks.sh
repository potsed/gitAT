#!/bin/bash

# Setup script to install git hooks for binary name validation

echo "🔧 Setting up Git hooks for GitAT development..."

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Install pre-commit hook
if [ -f ".githooks/pre-commit" ]; then
    cp .githooks/pre-commit .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
    echo "✅ Installed pre-commit hook"
else
    echo "❌ Error: .githooks/pre-commit not found"
    exit 1
fi

# Configure git to use our hooks directory
git config core.hooksPath .githooks

echo "✅ Git hooks setup complete"
echo ""
echo "The following hooks are now active:"
echo "  - pre-commit: Prevents committing deprecated 'gitat' binary"
echo ""
echo "To disable hooks temporarily, use: git commit --no-verify"