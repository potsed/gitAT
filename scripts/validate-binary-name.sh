#!/bin/bash

# CI/CD script to validate binary naming conventions
# This script ensures that only 'git-@' is built and no deprecated 'gitat' exists

set -e

echo "🔍 Validating binary naming conventions..."

# Check if deprecated binary exists
if [ -f "gitat" ]; then
    echo "❌ FAIL: Found deprecated binary 'gitat'"
    echo "   The correct binary name is 'git-@'"
    echo "   Please ensure your build process creates 'git-@' not 'gitat'"
    exit 1
fi

# Check if correct binary exists (when expected)
if [ "$1" = "--expect-binary" ] && [ ! -f "git-@" ]; then
    echo "❌ FAIL: Expected binary 'git-@' not found"
    echo "   Please build the project with 'make build-local'"
    exit 1
fi

# Check Makefile for correct binary name
if ! grep -q "BINARY_NAME=git-@" Makefile; then
    echo "❌ FAIL: Makefile does not set BINARY_NAME=git-@"
    exit 1
fi

# Check .gitignore includes deprecated binary
if ! grep -q "gitat" .gitignore; then
    echo "❌ FAIL: .gitignore should include 'gitat' to prevent accidental commits"
    exit 1
fi

# Check for problematic references to 'gitat' in build scripts (excluding comments and archives)
PROBLEMATIC_FILES=$(find . -name "*.sh" \
    -not -path "./scripts/validate-binary-name.sh" \
    -not -path "./.githooks/*" \
    -not -path "./archive/*" \
    -not -path "./install.sh" \
    -exec grep -l "go build.*gitat\|BINARY.*gitat\|gitat.*exe" {} \; 2>/dev/null | grep -v ".git" || true)

if [ -n "$PROBLEMATIC_FILES" ]; then
    echo "❌ FAIL: Found problematic references to 'gitat' binary in scripts:"
    echo "$PROBLEMATIC_FILES"
    echo "   Please update build commands to use 'git-@' instead"
    exit 1
fi

echo "✅ PASS: Binary naming validation successful"
echo "✅ Correct binary name: git-@"
echo "✅ No deprecated 'gitat' binary found"
echo "✅ Build configuration is correct"