#!/bin/bash

input="sweep command fixes"
echo "Original: '$input'"

# Convert to lowercase
input=$(echo "$input" | tr '[:upper:]' '[:lower:]')
echo "After lowercase: '$input'"

# Replace spaces with hyphens
input=$(echo "$input" | sed 's/[[:space:]]/-/g')
echo "After space replacement: '$input'"

# Replace underscores with hyphens
input=$(echo "$input" | sed 's/_/-/g')
echo "After underscore replacement: '$input'"

# Replace dots with hyphens
input=$(echo "$input" | sed 's/\./-/g')
echo "After dot replacement: '$input'"

# Replace backslashes with hyphens
input=$(echo "$input" | sed 's/\\/-/g')
echo "After backslash replacement: '$input'"

# Remove any non-alphanumeric characters except hyphens
input=$(echo "$input" | sed 's/[^a-z0-9-]//g')
echo "After non-alphanumeric removal: '$input'"

# Remove multiple consecutive hyphens
input=$(echo "$input" | sed 's/--\+/-/g')
echo "After multiple hyphen removal: '$input'"

# Remove leading and trailing hyphens
input=$(echo "$input" | sed 's/^-\+//; s/-\+$//')
echo "Final result: '$input'" 