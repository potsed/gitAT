#!/bin/bash

echo "Testing git @ work branch name formatting..."

# Test the format_branch_name function
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

echo "Testing format_branch_name function:"
echo "===================================="

# Test cases
test_cases=(
    "Incorrect Branch Name"
    "Fix Login Bug!"
    "Add User Authentication"
    "Update API Documentation"
    "Fix_Crash_On_Startup"
    "Update.Dependencies"
    "Remove Old Code!!!"
    "Add New Feature (v2.0)"
    "Fix Bug #123"
    "Update README.md"
    "Add Tests for Login"
    "Refactor User Model"
)

for test_input in "${test_cases[@]}"; do
    expected_output=$(format_branch_name "$test_input")
    echo "Input:  '$test_input'"
    echo "Output: '$expected_output'"
    echo "---"
done

echo ""
echo "✅ Branch name formatting test complete!"
echo ""
echo "The format_branch_name function should convert:"
echo "- Uppercase to lowercase"
echo "- Spaces, underscores, dots to hyphens"
echo "- Remove special characters"
echo "- Remove multiple consecutive hyphens"
echo "- Remove leading/trailing hyphens" 