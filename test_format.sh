#!/bin/bash

# Format branch name to kebab-case
format_branch_name() {
    local input="$1"
    
    echo "Original input: '$input'"
    
    # Convert to lowercase
    input=$(echo "$input" | tr '[:upper:]' '[:lower:]')
    echo "After lowercase: '$input'"
    
    # Replace spaces, underscores, and other separators with hyphens
    input=$(echo "$input" | sed 's/[[:space:]_\.\/\\]+/-/g')
    echo "After space replacement: '$input'"
    
    # Remove any non-alphanumeric characters except hyphens
    input=$(echo "$input" | sed 's/[^a-z0-9-]//g')
    echo "After non-alphanumeric removal: '$input'"
    
    # Remove multiple consecutive hyphens
    input=$(echo "$input" | sed 's/--\+/-/g')
    echo "After multiple hyphen removal: '$input'"
    
    # Remove leading and trailing hyphens
    input=$(echo "$input" | sed 's/^-\+//; s/-\+$//')
    echo "After leading/trailing hyphen removal: '$input'"
    
    # If empty after formatting, use "update"
    if [ -z "$input" ]; then
        input="update"
    fi
    
    echo "Final result: '$input'"
    echo "$input"
}

echo "Testing format_branch_name function:"
format_branch_name "sweep command fixes" 