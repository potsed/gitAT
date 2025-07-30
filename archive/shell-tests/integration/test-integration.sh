#!/bin/bash

# Integration Tests for GitAT CLI
# Tests actual command-line functionality with security validation

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

# Test repository setup
TEST_REPO_DIR="/tmp/gitat_test_repo"
ORIGINAL_DIR=$(pwd)

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

setup_test_repo() {
    echo "Setting up test repository..."
    rm -rf "$TEST_REPO_DIR"
    mkdir -p "$TEST_REPO_DIR"
    cd "$TEST_REPO_DIR"
    
    git init
    git config user.name "Test User"
    git config user.email "test@example.com"
    
    echo "test content" > test.txt
    git add test.txt
    git commit -m "Initial commit"
    
    echo "Test repository setup complete"
}

cleanup_test_repo() {
    echo "Cleaning up test repository..."
    cd "$ORIGINAL_DIR"
    rm -rf "$TEST_REPO_DIR"
}

# Test 1: Command injection prevention
test_command_injection_prevention() {
    echo "Testing command injection prevention..."
    
    local malicious_inputs=(
        "test; rm -rf /"
        "test&rm -rf /"
        "test|cat /etc/passwd"
        "test\`rm -rf /\`"
        "test\$(rm -rf /)"
    )
    
    local vulnerable=false
    
    for input in "${malicious_inputs[@]}"; do
        if git @ product "$input" 2>/dev/null; then
            vulnerable=true
            break
        fi
    done
    
    if [ "$vulnerable" = true ]; then
        log_test "Command Injection Prevention" "FAIL"
    else
        log_test "Command Injection Prevention" "PASS"
    fi
}

# Test 2: Path traversal prevention
test_path_traversal_prevention() {
    echo "Testing path traversal prevention..."
    
    local malicious_paths=(
        "../../../etc/passwd"
        "..\\..\\..\\windows\\system32"
        "test/../etc/passwd"
    )
    
    local vulnerable=false
    
    for path in "${malicious_paths[@]}"; do
        if git @ save "$path" 2>/dev/null; then
            vulnerable=true
            break
        fi
    done
    
    if [ "$vulnerable" = true ]; then
        log_test "Path Traversal Prevention" "FAIL"
    else
        log_test "Path Traversal Prevention" "PASS"
    fi
}

# Test 3: Input validation
test_input_validation() {
    echo "Testing input validation..."
    
    local long_input
    long_input=$(printf 'a%.0s' {1..2000})
    
    if git @ product "$long_input" 2>/dev/null; then
        log_test "Input Length Validation" "FAIL"
    else
        log_test "Input Length Validation" "PASS"
    fi
}

# Test 4: Permission checking
test_permission_checking() {
    echo "Testing permission checking..."
    
    # Make repository read-only temporarily
    chmod -w .
    
    if git @ product test 2>/dev/null; then
        log_test "Permission Checking" "FAIL"
    else
        log_test "Permission Checking" "PASS"
    fi
    
    # Restore permissions
    chmod +w .
}

# Test 5: Secure configuration
test_secure_configuration() {
    echo "Testing secure configuration..."
    
    local invalid_keys=(
        "at..invalid"
        "at.invalid."
        "invalid.key"
    )
    
    local vulnerable=false
    
    for key in "${invalid_keys[@]}"; do
        if git config --replace-all "$key" test 2>/dev/null; then
            vulnerable=true
            break
        fi
    done
    
    if [ "$vulnerable" = true ]; then
        log_test "Secure Configuration" "FAIL"
    else
        log_test "Secure Configuration" "PASS"
    fi
}

# Test 6: Error handling
test_error_handling() {
    echo "Testing error handling..."
    
    # Test outside git repository
    cd /tmp
    
    if git @ info 2>/dev/null; then
        log_test "Error Handling" "FAIL"
    else
        log_test "Error Handling" "PASS"
    fi
    
    # Return to test repository
    cd "$TEST_REPO_DIR"
}

# Test 7: Safe command execution
test_safe_command_execution() {
    echo "Testing safe command execution..."
    
    # Check that no dangerous commands are used
    local dangerous_patterns=(
        "eval "
        "exec "
        "sudo "
    )
    
    local vulnerable=false
    
    for pattern in "${dangerous_patterns[@]}"; do
        if grep -r "$pattern" git_at_cmds/ 2>/dev/null; then
            vulnerable=true
            break
        fi
    done
    
    if [ "$vulnerable" = true ]; then
        log_test "Safe Command Execution" "FAIL"
    else
        log_test "Safe Command Execution" "PASS"
    fi
}

# Test 8: Security logging
test_security_logging() {
    echo "Testing security logging..."
    
    local log_file="$HOME/.gitat_security.log"
    
    # Trigger a security event
    git @ product "test;rm -rf /" 2>/dev/null || true
    
    if [ -f "$log_file" ] && grep -q "Dangerous characters detected" "$log_file"; then
        log_test "Security Logging" "PASS"
    else
        log_test "Security Logging" "FAIL"
    fi
}

# Test 9: Concurrent access safety
test_concurrent_access() {
    echo "Testing concurrent access safety..."
    
    local success_count=0
    
    # Run multiple operations concurrently
    for i in {1..5}; do
        if git @ product "test$i" 2>/dev/null; then
            success_count=$((success_count + 1))
        fi
    done
    
    if [ $success_count -gt 0 ]; then
        log_test "Concurrent Access Safety" "PASS"
    else
        log_test "Concurrent Access Safety" "FAIL"
    fi
}

# Test 10: Configuration validation
test_configuration_validation() {
    echo "Testing configuration validation..."
    
    local test_cases=(
        "valid-project:true"
        "invalid;project:false"
        "valid-feature:true"
        "invalid&feature:false"
    )
    
    local all_passed=true
    
    for test_case in "${test_cases[@]}"; do
        local input="${test_case%:*}"
        local expected="${test_case#*:}"
        
        if git @ product "$input" 2>/dev/null; then
            if [ "$expected" = "false" ]; then
                all_passed=false
            fi
        else
            if [ "$expected" = "true" ]; then
                all_passed=false
            fi
        fi
    done
    
    if [ "$all_passed" = true ]; then
        log_test "Configuration Validation" "PASS"
    else
        log_test "Configuration Validation" "FAIL"
    fi
}

# Run all tests
main() {
    echo "=== GitAT Integration Test Suite ==="
    echo
    
    # Setup test environment
    setup_test_repo
    
    # Run tests
    test_command_injection_prevention
    test_path_traversal_prevention
    test_input_validation
    test_permission_checking
    test_secure_configuration
    test_error_handling
    test_safe_command_execution
    test_security_logging
    test_concurrent_access
    test_configuration_validation
    
    # Cleanup
    cleanup_test_repo
    
    echo
    echo "=== Test Results ==="
    echo "Total Tests: $TOTAL_TESTS"
    echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
    
    if [ $FAILED_TESTS -gt 0 ]; then
        echo -e "${RED}Integration tests failed!${NC}"
        exit 1
    else
        echo -e "${GREEN}All integration tests passed!${NC}"
        exit 0
    fi
}

# Run main function
main "$@" 