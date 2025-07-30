#!/bin/bash

# Test for git @ work command (Conventional Commits integration)

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Helper functions
log_test() {
    local test_name="$1"
    local result="$2"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$result" = "PASS" ]; then
        echo -e "${GREEN}✓${NC} $test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗${NC} $test_name"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# Test 1: Check if we're in a git repository
test_git_repository() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_test "Git Repository Check" "FAIL"
        return 1
    fi
    log_test "Git Repository Check" "PASS"
}

# Test 2: Test help functionality
test_help_functionality() {
    if git @ work -h > /dev/null 2>&1; then
        log_test "Help Functionality" "PASS"
    else
        log_test "Help Functionality" "FAIL"
    fi
}

# Test 3: Check branch configuration
test_branch_configuration() {
    local current_branch
    local trunk_branch
    
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    trunk_branch=$(git config at.trunk 2>/dev/null || echo "main")
    
    if [ -n "$current_branch" ] && [ "$current_branch" != "HEAD" ]; then
        log_test "Branch Configuration" "PASS"
    else
        log_test "Branch Configuration" "FAIL"
    fi
}

# Test 4: Check trunk branch existence
test_trunk_branch_existence() {
    local trunk_branch
    local current_branch
    
    trunk_branch=$(git config at.trunk 2>/dev/null || echo "main")
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    
    if git rev-parse --verify "$trunk_branch" >/dev/null 2>&1; then
        log_test "Trunk Branch Existence" "PASS"
    else
        # Create test trunk branch if it doesn't exist
        git checkout -b "$trunk_branch" 2>/dev/null || true
        echo "Test trunk" > trunk_file.txt
        git add trunk_file.txt
        git commit -m "Initial trunk commit" >/dev/null 2>&1 || true
        git checkout "$current_branch" 2>/dev/null || true
        log_test "Trunk Branch Creation" "PASS"
    fi
}

# Test 5: Test error handling for missing work type
test_missing_work_type() {
    if git @ work 2>&1 | grep -q "Work type is required"; then
        log_test "Missing Work Type Error" "PASS"
    else
        log_test "Missing Work Type Error" "FAIL"
    fi
}

# Test 6: Test error handling for invalid work type
test_invalid_work_type() {
    if git @ work invalid-type test 2>&1 | grep -q "Invalid work type"; then
        log_test "Invalid Work Type Error" "PASS"
    else
        log_test "Invalid Work Type Error" "FAIL"
    fi
}

# Test 7: Test error handling for missing name option value
test_missing_name_option() {
    local failed=false
    
    if ! git @ work -n 2>&1 | grep -q "requires a value"; then
        failed=true
    fi
    
    if ! git @ work --name 2>&1 | grep -q "requires a value"; then
        failed=true
    fi
    
    if [ "$failed" = true ]; then
        log_test "Missing Name Option Error" "FAIL"
    else
        log_test "Missing Name Option Error" "PASS"
    fi
}

# Test 8: Test error handling for too many arguments
test_too_many_arguments() {
    if git @ work feature add-auth extra-arg 2>&1 | grep -q "Too many arguments"; then
        log_test "Too Many Arguments Error" "PASS"
    else
        log_test "Too Many Arguments Error" "FAIL"
    fi
}

# Test 9: Test feature branch creation
test_feature_branch_creation() {
    local feature_name
    local new_branch
    local working_branch
    
    feature_name="test-feature-$(date +%s)"
    
    if git @ work feature "$feature_name" 2>&1 | grep -q "feature branch.*created successfully"; then
        # Check if we're on the feature branch
        new_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
        if [ "$new_branch" = "feature-$feature_name" ]; then
            # Check if working branch is set
            working_branch=$(git config at.branch 2>/dev/null || echo "")
            if [ "$working_branch" = "feature-$feature_name" ]; then
                log_test "Feature Branch Creation" "PASS"
            else
                log_test "Feature Branch Creation" "FAIL"
            fi
        else
            log_test "Feature Branch Creation" "FAIL"
        fi
    else
        log_test "Feature Branch Creation" "FAIL"
    fi
}

# Test 10: Test branch name formatting
test_branch_name_formatting() {
    local test_name
    local expected_branch
    
    test_name="Test Branch Name"
    expected_branch="feature-test-branch-name"
    
    if git @ work feature "$test_name" 2>&1 | grep -q "feature branch.*created successfully"; then
        local new_branch
        new_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
        if [ "$new_branch" = "$expected_branch" ]; then
            log_test "Branch Name Formatting" "PASS"
        else
            log_test "Branch Name Formatting" "FAIL"
        fi
    else
        log_test "Branch Name Formatting" "FAIL"
    fi
}

# Run all tests
main() {
    echo "=== GitAT Work Command Test Suite ==="
    echo
    
    test_git_repository
    test_help_functionality
    test_branch_configuration
    test_trunk_branch_existence
    test_missing_work_type
    test_invalid_work_type
    test_missing_name_option
    test_too_many_arguments
    test_feature_branch_creation
    test_branch_name_formatting
    
    echo
    echo "=== Test Results ==="
    echo "Total Tests: $TOTAL_TESTS"
    echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
    
    if [ $FAILED_TESTS -gt 0 ]; then
        echo -e "${RED}Work command tests failed!${NC}"
        exit 1
    else
        echo -e "${GREEN}All work command tests passed!${NC}"
        exit 0
    fi
}

# Run main function
main "$@" 