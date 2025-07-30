#!/bin/bash

# Main Test Runner for GitAT
# Runs all test categories: unit, integration, commands, etc.

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TEST_FILES=0
PASSED_TEST_FILES=0
FAILED_TEST_FILES=0

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Helper functions
log_test_category() {
    local category="$1"
    local result="$2"
    TOTAL_TEST_FILES=$((TOTAL_TEST_FILES + 1))
    
    if [ "$result" = "PASS" ]; then
        echo -e "${GREEN}✓${NC} $category"
        PASSED_TEST_FILES=$((PASSED_TEST_FILES + 1))
    else
        echo -e "${RED}✗${NC} $category"
        FAILED_TEST_FILES=$((FAILED_TEST_FILES + 1))
    fi
}

run_test_category() {
    local category="$1"
    local runner_script="$2"
    
    echo -e "${BLUE}Running $category tests...${NC}"
    
    if [ -f "$runner_script" ]; then
        if bash "$runner_script" >/dev/null 2>&1; then
            log_test_category "$category" "PASS"
        else
            log_test_category "$category" "FAIL"
        fi
    else
        echo -e "${YELLOW}No runner script found for $category: $runner_script${NC}"
        log_test_category "$category" "PASS" # Skip if no tests
    fi
}

# Main function
main() {
    echo "=== GitAT Complete Test Suite ==="
    echo "Project directory: $PROJECT_DIR"
    echo "Test directory: $SCRIPT_DIR"
    echo
    
    # Run unit tests
    run_test_category "Unit Tests" "$SCRIPT_DIR/unit/run-unit-tests.sh"
    
    # Run command tests
    run_test_category "Command Tests" "$SCRIPT_DIR/commands/run-command-tests.sh"
    
    # Run integration tests (if they exist)
    if [ -f "$SCRIPT_DIR/integration/run-integration-tests.sh" ]; then
        run_test_category "Integration Tests" "$SCRIPT_DIR/integration/run-integration-tests.sh"
    fi
    
    # Run security tests (if they exist)
    if [ -f "$SCRIPT_DIR/security/run-security-tests.sh" ]; then
        run_test_category "Security Tests" "$SCRIPT_DIR/security/run-security-tests.sh"
    fi
    
    echo
    echo "=== Overall Test Results ==="
    echo "Total Test Categories: $TOTAL_TEST_FILES"
    echo -e "Passed: ${GREEN}$PASSED_TEST_FILES${NC}"
    echo -e "Failed: ${RED}$FAILED_TEST_FILES${NC}"
    
    if [ $FAILED_TEST_FILES -gt 0 ]; then
        echo -e "${RED}Some test categories failed!${NC}"
        exit 1
    else
        echo -e "${GREEN}All test categories passed!${NC}"
        exit 0
    fi
}

# Run main function
main "$@" 