#!/bin/bash

# Security Test Suite for GitAT
# Tests for OWASP Top 10 and common shell script vulnerabilities

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

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

# Test 1: Command Injection via Unquoted Variables
test_command_injection() {
    echo "Testing for command injection vulnerabilities..."
    
    # Test for unquoted variables in various scripts
    local vulnerable=false
    
    # Check for actual command injection vulnerabilities by looking for specific dangerous patterns
    # Only check for patterns that are not in comments and not in safe contexts
    local dangerous_patterns=(
        "git config.*\$1[^\"']"
        "git.*\$1[^\"']"
        "cd.*\$1[^\"']"
    )
    
    for pattern in "${dangerous_patterns[@]}"; do
        # Use awk to exclude comment lines and safe contexts
        if awk '!/^\s*#/ && !/^\s*\/\// && /'"$pattern"'/ && !/validate_input/ && !/^\s*#/' git_at_cmds/*.sh 2>/dev/null | grep -q .; then
            vulnerable=true
            break
        fi
    done
    
    if [ "$vulnerable" = true ]; then
        log_test "Command Injection Protection" "FAIL"
    else
        log_test "Command Injection Protection" "PASS"
    fi
}

# Test 2: Path Traversal Protection
test_path_traversal() {
    echo "Testing for path traversal vulnerabilities..."
    
    # Check if scripts validate paths before operations
    local vulnerable=false
    
    if grep -q "cd \$path" git_at_cmds/save.sh; then
        vulnerable=true
    fi
    
    if [ "$vulnerable" = true ]; then
        log_test "Path Traversal Protection" "FAIL"
    else
        log_test "Path Traversal Protection" "PASS"
    fi
}

# Test 3: Input Validation
test_input_validation() {
    echo "Testing input validation..."
    
    local vulnerable=false
    
    # Check for proper input validation in critical functions
    if ! grep -r "validate_input" git_at_cmds/*.sh 2>/dev/null; then
        vulnerable=true
    fi
    
    if [ "$vulnerable" = true ]; then
        log_test "Input Validation" "FAIL"
    else
        log_test "Input Validation" "PASS"
    fi
}

# Test 4: Privilege Escalation Protection
test_privilege_escalation() {
    echo "Testing privilege escalation protection..."
    
    local vulnerable=false
    
    # Check if scripts run with elevated privileges unnecessarily
    if grep -q "sudo" git_at_cmds/*.sh; then
        vulnerable=true
    fi
    
    if [ "$vulnerable" = true ]; then
        log_test "Privilege Escalation Protection" "FAIL"
    else
        log_test "Privilege Escalation Protection" "PASS"
    fi
}

# Test 5: Data Exposure Protection
test_data_exposure() {
    echo "Testing data exposure protection..."
    
    local vulnerable=false
    
    # Check if sensitive data is properly handled
    if grep -q "echo.*password" git_at_cmds/*.sh; then
        vulnerable=true
    fi
    
    if [ "$vulnerable" = true ]; then
        log_test "Data Exposure Protection" "FAIL"
    else
        log_test "Data Exposure Protection" "PASS"
    fi
}

# Test 6: Error Handling
test_error_handling() {
    echo "Testing error handling..."
    
    local vulnerable=false
    
    # Check if scripts have proper error handling
    if ! grep -q "set -e" git_at_cmds/*.sh; then
        vulnerable=true
    fi
    
    if [ "$vulnerable" = true ]; then
        log_test "Error Handling" "FAIL"
    else
        log_test "Error Handling" "PASS"
    fi
}

# Test 7: Safe Command Execution
test_safe_command_execution() {
    echo "Testing safe command execution..."
    
    local vulnerable=false
    
    # Check for unsafe command execution patterns
    if grep -r "eval " git_at_cmds/*.sh 2>/dev/null; then
        vulnerable=true
    fi
    
    if grep -r "exec " git_at_cmds/*.sh 2>/dev/null; then
        vulnerable=true
    fi
    
    if [ "$vulnerable" = true ]; then
        log_test "Safe Command Execution" "FAIL"
    else
        log_test "Safe Command Execution" "PASS"
    fi
}

# Test 8: Authentication and Authorization
test_auth_authorization() {
    echo "Testing authentication and authorization..."
    
    local vulnerable=false
    
    # Check if scripts verify user permissions
    if ! grep -r "check_permissions" git_at_cmds/*.sh 2>/dev/null; then
        vulnerable=true
    fi
    
    if [ "$vulnerable" = true ]; then
        log_test "Authentication and Authorization" "FAIL"
    else
        log_test "Authentication and Authorization" "PASS"
    fi
}

# Test 9: Logging and Monitoring
test_logging_monitoring() {
    echo "Testing logging and monitoring..."
    
    local vulnerable=false
    
    # Check if security logging is implemented
    if ! grep -r "log_security_event" git_at_cmds/*.sh 2>/dev/null; then
        vulnerable=true
    fi
    
    if [ "$vulnerable" = true ]; then
        log_test "Logging and Monitoring" "FAIL"
    else
        log_test "Logging and Monitoring" "PASS"
    fi
}

# Test 10: Configuration Security
test_configuration_security() {
    echo "Testing configuration security..."
    
    local vulnerable=false
    
    # Check if configuration files are properly secured (macOS compatible)
    if [ -f ".git/config" ]; then
        local perms
        perms=$(stat -f %Lp .git/config 2>/dev/null || stat -c %a .git/config 2>/dev/null || echo "644")
        if [ "$perms" != "600" ]; then
            vulnerable=true
        fi
    fi
    
    if [ "$vulnerable" = true ]; then
        log_test "Configuration Security" "FAIL"
    else
        log_test "Configuration Security" "PASS"
    fi
}

# Run all tests
main() {
    echo "=== GitAT Security Test Suite ==="
    echo
    
    test_command_injection
    test_path_traversal
    test_input_validation
    test_privilege_escalation
    test_data_exposure
    test_error_handling
    test_safe_command_execution
    test_auth_authorization
    test_logging_monitoring
    test_configuration_security
    
    echo
    echo "=== Test Results ==="
    echo "Total Tests: $TOTAL_TESTS"
    echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
    
    if [ $FAILED_TESTS -gt 0 ]; then
        echo -e "${RED}Security vulnerabilities detected!${NC}"
        exit 1
    else
        echo -e "${GREEN}All security tests passed!${NC}"
        exit 0
    fi
}

main "$@" 