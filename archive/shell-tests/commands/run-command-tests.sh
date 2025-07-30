#!/bin/bash

# Command Test Runner for GitAT
# Runs all command tests in the commands directory

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMMANDS_DIR="$SCRIPT_DIR"
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

run_test_file() {
    local test_file="$1"
    local test_name
    
    # Extract test name from filename
    test_name=$(basename "$test_file" .sh | sed 's/test-//')
    
    echo -e "${BLUE}Running: $test_name${NC}"
    
    if bash "$test_file" >/dev/null 2>&1; then
        log_test "$test_name" "PASS"
    else
        log_test "$test_name" "FAIL"
    fi
}

# Main function
main() {
    echo "=== GitAT Command Test Runner ==="
    echo "Commands tests directory: $COMMANDS_DIR"
    echo
    
    # Find all test files
    local test_files=()
    while IFS= read -r -d '' file; do
        test_files+=("$file")
    done < <(find "$COMMANDS_DIR" -name "test-*.sh" -type f -print0)
    
    if [ ${#test_files[@]} -eq 0 ]; then
        echo -e "${YELLOW}No command test files found in $COMMANDS_DIR${NC}"
        exit 0
    fi
    
    echo "Found ${#test_files[@]} command test file(s):"
    for file in "${test_files[@]}"; do
        echo "  - $(basename "$file")"
    done
    echo
    
    # Run each test file
    for test_file in "${test_files[@]}"; do
        run_test_file "$test_file"
    done
    
    echo
    echo "=== Command Test Results ==="
    echo "Total Test Files: ${#test_files[@]}"
    echo "Total Tests: $TOTAL_TESTS"
    echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
    
    if [ $FAILED_TESTS -gt 0 ]; then
        echo -e "${RED}Command tests failed!${NC}"
        exit 1
    else
        echo -e "${GREEN}All command tests passed!${NC}"
        exit 0
    fi
}

# Run main function
main "$@" 