#!/bin/bash

echo "Testing git @ pr automatic description generation..."

# Source the pr.sh file to test the generate_auto_description function
source git_at_cmds/pr.sh

# Create a test repository structure
echo "Setting up test environment..."

# Create test files
mkdir -p test-dir
echo "test content" > test-dir/file1.txt
echo "test content" > test-dir/file2.sh
echo "test content" > README.md
echo "test content" > src/main.js

# Add files to git
git add test-dir/file1.txt test-dir/file2.sh README.md src/main.js
git commit -m "Add test files for auto-description testing"

# Create a feature branch
git checkout -b feature-auto-description-test

# Make some changes
echo "updated content" > test-dir/file1.txt
echo "new content" > test-dir/newfile.py
rm test-dir/file2.sh
echo "updated readme" > README.md

# Add changes
git add test-dir/file1.txt test-dir/newfile.py README.md
git rm test-dir/file2.sh
git commit -m "Update files for testing auto-description"

# Test the auto-description generation
echo ""
echo "Testing generate_auto_description function:"
echo "=========================================="

# Get the current branch and base branch
current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
base_branch="main"

echo "Current branch: $current_branch"
echo "Base branch: $base_branch"
echo ""

# Generate auto description
auto_description=$(generate_auto_description "$base_branch" "$current_branch")

echo "Generated Description:"
echo "====================="
echo "$auto_description"

echo ""
echo "✅ Auto-description generation test complete!"
echo ""
echo "The description should include:"
echo "- Summary of changed files (added, modified, deleted)"
echo "- File types and directories affected"
echo "- List of all changed files with status icons"
echo "- Commit count and recent commit history" 