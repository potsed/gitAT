#!/bin/bash

# Master Test Runner for GitAT
# Runs all test suites: security, unit, integration, and shellcheck

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TOTAL_SUITES=0
PASSED_SUITES=0
FAILED_SUITES=0

# Helper functions
log_suite() {
    local suite_name="$1"
    local result="$2"
    TOTAL_SUITES=$((TOTAL_SUITES + 1))
    
    if [ "$result" = "PASS" ]; then
        echo -e "${GREEN}✓${NC} $suite_name"
        PASSED_SUITES=$((PASSED_SUITES + 1))
    else
        echo -e "${RED}✗${NC} $suite_name"
        FAILED_SUITES=$((FAILED_SUITES + 1))
    fi
}

run_test_suite() {
    local suite_name="$1"
    local script_path="$2"
    
    echo -e "\n${BLUE}Running $suite_name...${NC}"
    
    if [ -f "$script_path" ] && [ -x "$script_path" ]; then
        if "$script_path"; then
            log_suite "$suite_name" "PASS"
        else
            log_suite "$suite_name" "FAIL"
        fi
    else
        echo -e "${YELLOW}Warning: $script_path not found or not executable${NC}"
        log_suite "$suite_name" "FAIL"
    fi
}

check_shellcheck() {
    echo -e "\n${BLUE}Running ShellCheck...${NC}"
    
    if command -v shellcheck >/dev/null 2>&1; then
        if shellcheck git-@ git_at_cmds/*.sh 2>/dev/null; then
            log_suite "ShellCheck" "PASS"
        else
            log_suite "ShellCheck" "FAIL"
        fi
    else
        echo -e "${YELLOW}Warning: ShellCheck not installed. Install with: brew install shellcheck${NC}"
        log_suite "ShellCheck" "FAIL"
    fi
}

check_syntax() {
    echo -e "\n${BLUE}Checking shell script syntax...${NC}"
    
    local failed=false
    
    # Check main script
    if ! bash -n git-@; then
        echo -e "${RED}Syntax error in git-@${NC}"
        failed=true
    fi
    
    # Check all command scripts
    for script in git_at_cmds/*.sh; do
        if [ -f "$script" ]; then
            if ! bash -n "$script"; then
                echo -e "${RED}Syntax error in $script${NC}"
                failed=true
            fi
        fi
    done
    
    if [ "$failed" = true ]; then
        log_suite "Syntax Check" "FAIL"
    else
        log_suite "Syntax Check" "PASS"
    fi
}

check_permissions() {
    echo -e "\n${BLUE}Checking file permissions...${NC}"
    
    local failed=false
    
    # Check if main script is executable
    if [ ! -x git-@ ]; then
        echo -e "${RED}git-@ is not executable${NC}"
        failed=true
    fi
    
    # Check if command scripts are executable
    for script in git_at_cmds/*.sh; do
        if [ -f "$script" ] && [ ! -x "$script" ]; then
            echo -e "${RED}$script is not executable${NC}"
            failed=true
        fi
    done
    
    if [ "$failed" = true ]; then
        log_suite "Permission Check" "FAIL"
    else
        log_suite "Permission Check" "PASS"
    fi
}

check_dependencies() {
    echo -e "\n${BLUE}Checking dependencies...${NC}"
    
    local failed=false
    local dependencies=("git" "bash")
    
    for dep in "${dependencies[@]}"; do
        if ! command -v "$dep" >/dev/null 2>&1; then
            echo -e "${RED}Required dependency not found: $dep${NC}"
            failed=true
        fi
    done
    
    if [ "$failed" = true ]; then
        log_suite "Dependency Check" "FAIL"
    else
        log_suite "Dependency Check" "PASS"
    fi
}

# Main test runner
main() {
    echo "=== GitAT Test Suite Runner ==="
    echo
    
    # Get script directory
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
    
    # Change to project directory
    cd "$PROJECT_DIR"
    
    # Run all test suites
    check_dependencies
    check_syntax
    check_permissions
    run_test_suite "Unit Tests" "tests/unit.test.sh"
    run_test_suite "Security Tests" "tests/security.test.sh"
    run_test_suite "Integration Tests" "tests/integration.test.sh"
    check_shellcheck
    
    echo
    echo "=== Test Results Summary ==="
    echo "Total Test Suites: $TOTAL_SUITES"
    echo -e "Passed: ${GREEN}$PASSED_SUITES${NC}"
    echo -e "Failed: ${RED}$FAILED_SUITES${NC}"
    
    if [ $FAILED_SUITES -gt 0 ]; then
        echo -e "\n${RED}Some test suites failed!${NC}"
        echo -e "${YELLOW}Please review the output above and fix any issues.${NC}"
        exit 1
    else
        echo -e "\n${GREEN}All test suites passed!${NC}"
        echo -e "${GREEN}GitAT is ready for use.${NC}"
        exit 0
    fi
}

# Run main function
main "$@" 