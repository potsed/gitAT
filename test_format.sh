#!/bin/bash

# Format branch name to kebab-case
format_branch_name() {
    local input="$1"
    
    # Convert to lowercase
    input=$(echo "$input" | tr '[:upper:]' '[:lower:]')
    
    # Replace spaces, underscores, and other separators with hyphens
    input=$(echo "$input" | sed 's/[[:space:]_\.\/\\]+/-/g')
    
    # Remove any non-alphanumeric characters except hyphens
    input=$(echo "$input" | sed 's/[^a-z0-9-]//g')
    
    # Remove multiple consecutive hyphens
    input=$(echo "$input" | sed 's/--\+/-/g')
    
    # Remove leading and trailing hyphens
    input=$(echo "$input" | sed 's/^-\+//; s/-\+$//')
    
    # If empty after formatting, use "update"
    if [ -z "$input" ]; then
        input="update"
    fi
    
    echo "$input"
}

echo "Testing branch name formatting:"
echo "================================"

# Test cases
echo "Input: 'Incorrect Branch Name'"
echo "Output: '$(format_branch_name "Incorrect Branch Name")'"
echo "Expected: 'incorrect-branch-name'"
echo ""

echo "Input: 'Fix Login Bug!'"
echo "Output: '$(format_branch_name "Fix Login Bug!")'"
echo "Expected: 'fix-login-bug'"
echo ""

echo "Input: 'Add User Authentication'"
echo "Output: '$(format_branch_name "Add User Authentication")'"
echo "Expected: 'add-user-authentication'"
echo ""

echo "Input: 'Update API Documentation'"
echo "Output: '$(format_branch_name "Update API Documentation")'"
echo "Expected: 'update-api-documentation'"
echo ""

echo "Input: 'Fix_Crash_On_Startup'"
echo "Output: '$(format_branch_name "Fix_Crash_On_Startup")'"
echo "Expected: 'fix-crash-on-startup'"
echo ""

echo "Input: 'Remove Old Code!!!'"
echo "Output: '$(format_branch_name "Remove Old Code!!!")'"
echo "Expected: 'remove-old-code'"
echo ""

echo "✅ Formatting function test complete!" 