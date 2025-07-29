#!/bin/bash

# Security utilities for GitAT
# Implements defensive coding practices and security checks

set -euo pipefail

# Security configuration
readonly SECURITY_LOG_FILE="${HOME}/.gitat_security.log"
readonly MAX_INPUT_LENGTH=1000
readonly ALLOWED_CHARS='[a-zA-Z0-9._-]'

# Logging function for security events
log_security_event() {
    local level="$1"
    local message="$2"
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo "[${timestamp}] [${level}] ${message}" >> "$SECURITY_LOG_FILE"
}

# Input validation function
validate_input() {
    local input="$1"
    local max_length="${2:-$MAX_INPUT_LENGTH}"
    
    # Check length
    if [ ${#input} -gt "$max_length" ]; then
        log_security_event "WARNING" "Input too long: ${#input} characters"
        return 1
    fi
    
    # Check for dangerous characters
    if echo "$input" | grep -q '[;&|`$(){}]'; then
        log_security_event "WARNING" "Dangerous characters detected in input"
        return 1
    fi
    
    # Check for path traversal attempts
    if echo "$input" | grep -q '\.\./'; then
        log_security_event "WARNING" "Path traversal attempt detected"
        return 1
    fi
    
    return 0
}

# Path validation function
validate_path() {
    local path="$1"
    local repo_root
    
    # Get repository root
    repo_root=$(git rev-parse --show-toplevel 2>/dev/null || echo "")
    
    if [ -z "$repo_root" ]; then
        log_security_event "ERROR" "Not in a git repository"
        return 1
    fi
    
    # Resolve path and check if it's within repository
    local resolved_path
    resolved_path=$(realpath "$path" 2>/dev/null || echo "")
    
    if [ -z "$resolved_path" ]; then
        log_security_event "WARNING" "Invalid path: $path"
        return 1
    fi
    
    # Check if path is within repository
    if [[ "$resolved_path" != "$repo_root"* ]]; then
        log_security_event "WARNING" "Path traversal attempt: $path"
        return 1
    fi
    
    return 0
}

# Command injection protection
sanitize_command() {
    local command="$1"
    
    # Remove dangerous characters
    command=$(echo "$command" | sed 's/[;&|`$(){}]//g')
    
    # Escape remaining special characters
    command=$(printf '%s' "$command" | sed 's/[][\\^$.*+?{}()|]/\\&/g')
    
    echo "$command"
}

# Permission checking
check_permissions() {
    local operation="$1"
    local target="$2"
    
    # Check if user has write permissions to repository
    if ! [ -w "$(git rev-parse --show-toplevel)" ]; then
        log_security_event "ERROR" "Insufficient permissions for operation: $operation"
        return 1
    fi
    
    # Check if target is writable
    if [ -n "$target" ] && [ -e "$target" ] && ! [ -w "$target" ]; then
        log_security_event "ERROR" "Target not writable: $target"
        return 1
    fi
    
    return 0
}

# Safe command execution
safe_execute() {
    local command="$1"
    local args="$2"
    
    # Validate command
    if ! validate_input "$command" 100; then
        log_security_event "ERROR" "Invalid command: $command"
        return 1
    fi
    
    # Validate arguments
    if [ -n "$args" ] && ! validate_input "$args" 500; then
        log_security_event "ERROR" "Invalid arguments: $args"
        return 1
    fi
    
    # Log command execution
    log_security_event "INFO" "Executing command: $command $args"
    
    # Execute command with proper error handling
    if "$command" $args; then
        log_security_event "INFO" "Command executed successfully: $command"
        return 0
    else
        log_security_event "ERROR" "Command failed: $command"
        return 1
    fi
}

# Configuration security
secure_config() {
    local key="$1"
    local value="$2"
    
    # Validate key
    if ! echo "$key" | grep -qE '^at\.[a-zA-Z0-9._-]+$'; then
        log_security_event "ERROR" "Invalid config key: $key"
        return 1
    fi
    
    # Validate value
    if ! validate_input "$value"; then
        log_security_event "ERROR" "Invalid config value for key: $key"
        return 1
    fi
    
    # Set configuration
    if git config --replace-all "$key" "$value"; then
        log_security_event "INFO" "Configuration updated: $key"
        return 0
    else
        log_security_event "ERROR" "Failed to update configuration: $key"
        return 1
    fi
}

# Error handling wrapper
handle_errors() {
    local exit_code=$?
    
    if [ $exit_code -ne 0 ]; then
        log_security_event "ERROR" "Command failed with exit code: $exit_code"
        echo "Error: Operation failed. Check logs for details." >&2
        exit $exit_code
    fi
}

# Set up error handling
trap handle_errors ERR

# Initialize security logging
init_security_logging() {
    # Create log directory if it doesn't exist
    local log_dir=$(dirname "$SECURITY_LOG_FILE")
    if [ ! -d "$log_dir" ]; then
        mkdir -p "$log_dir"
        chmod 700 "$log_dir"
    fi
    
    # Create log file if it doesn't exist
    if [ ! -f "$SECURITY_LOG_FILE" ]; then
        touch "$SECURITY_LOG_FILE"
        chmod 600 "$SECURITY_LOG_FILE"
    fi
    
    log_security_event "INFO" "Security logging initialized"
}

# Initialize security logging when this file is sourced
init_security_logging 