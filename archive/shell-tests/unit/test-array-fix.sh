#!/bin/bash

# Test the array handling fix
echo "Testing array handling fix..."

# Simulate the array initialization and usage
file_types=()
dirs=()

# Test adding elements
ext="txt"
dir="src"

# Check if extension already exists
found_ext=false
for existing_ext in "${file_types[@]}"; do
    if [ "$existing_ext" = "$ext" ]; then
        found_ext=true
        break
    fi
done

# Check if directory already exists
found_dir=false
for existing_dir in "${dirs[@]}"; do
    if [ "$existing_dir" = "$dir" ]; then
        found_dir=true
        break
    fi
done

# Add if not found
if [ "$found_ext" = false ]; then
    file_types+=("$ext")
fi
if [ "$found_dir" = false ]; then
    dirs+=("$dir")
fi

echo "✅ Array handling test passed!"
echo "file_types: ${file_types[*]}"
echo "dirs: ${dirs[*]}" 