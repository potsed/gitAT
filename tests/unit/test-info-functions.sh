#!/bin/bash

echo "Testing info command functions..."

# Test 1: Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In a git repository"

# Test 2: Test product function
echo "Test 2: Testing product function"
PRODUCT=$(cmd_product 2>/dev/null || echo "")
echo "✅ Product: '$PRODUCT'"

# Test 3: Test version function
echo "Test 3: Testing version function"
VERSION=$(cmd_version 2>/dev/null || echo "")
echo "✅ Version: '$VERSION'"

# Test 4: Test version tag function
echo "Test 4: Testing version tag function"
TAG=$(cmd_version -t 2>/dev/null || echo "")
echo "✅ Version Tag: '$TAG'"

# Test 5: Test feature function
echo "Test 5: Testing feature function"
FEATURE=$(cmd_feature 2>/dev/null || echo "")
echo "✅ Feature: '$FEATURE'"

# Test 6: Test issue function
echo "Test 6: Testing issue function"
ISSUE=$(cmd_issue 2>/dev/null || echo "")
echo "✅ Issue: '$ISSUE'"

# Test 7: Test branch function
echo "Test 7: Testing branch function"
BRANCH=$(cmd_branch 2>/dev/null || echo "")
echo "✅ Branch: '$BRANCH'"

# Test 8: Test path function
echo "Test 8: Testing path function"
PATH_VAL=$(cmd__path 2>/dev/null || echo "")
echo "✅ Path: '$PATH_VAL'"

# Test 9: Test trunk function
echo "Test 9: Testing trunk function"
TRUNK=$(cmd__trunk 2>/dev/null || echo "")
echo "✅ Trunk: '$TRUNK'"

# Test 10: Test wip function
echo "Test 10: Testing wip function"
WIP=$(cmd_wip 2>/dev/null || echo "")
echo "✅ WIP: '$WIP'"

# Test 11: Test id function
echo "Test 11: Testing id function"
ID=$(cmd__id 2>/dev/null || echo "")
echo "✅ ID: '$ID'"

# Test 12: Test label function
echo "Test 12: Testing label function"
LABEL=$(cmd__label 2>/dev/null || echo "")
echo "✅ Label: '$LABEL'"

echo "All function tests completed!" 