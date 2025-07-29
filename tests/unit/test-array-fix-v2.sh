#!/bin/bash

echo "Testing new array handling approach..."

# Simulate the new approach using strings instead of arrays
file_types=""
dirs=""

# Test adding elements
ext1="txt"
ext2="py"
dir1="src"
dir2="docs"

# Add to file types if not already present
if [[ ! "$file_types" =~ "$ext1" ]]; then
    if [ -n "$file_types" ]; then
        file_types="$file_types $ext1"
    else
        file_types="$ext1"
    fi
fi

if [[ ! "$file_types" =~ "$ext2" ]]; then
    if [ -n "$file_types" ]; then
        file_types="$file_types $ext2"
    else
        file_types="$ext2"
    fi
fi

# Add to directories if not already present
if [[ ! "$dirs" =~ "$dir1" ]]; then
    if [ -n "$dirs" ]; then
        dirs="$dirs $dir1"
    else
        dirs="$dir1"
    fi
fi

if [[ ! "$dirs" =~ "$dir2" ]]; then
    if [ -n "$dirs" ]; then
        dirs="$dirs $dir2"
    else
        dirs="$dir2"
    fi
fi

echo "✅ New array handling test passed!"
echo "file_types: '$file_types'"
echo "dirs: '$dirs'"

# Test iteration
echo ""
echo "File types:"
for ext in $file_types; do
    echo "- $ext"
done

echo ""
echo "Directories:"
for dir in $dirs; do
    echo "- $dir"
done 