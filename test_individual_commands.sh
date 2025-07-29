#!/bin/bash

echo "Testing individual commands..."

# Test 1: Test git @ _label
echo "Test 1: git @ _label"
LABEL=$(git @ _label 2>/dev/null || echo "FAILED")
echo "Label: '$LABEL'"

# Test 2: Test git @ _id
echo "Test 2: git @ _id"
ID=$(git @ _id 2>/dev/null || echo "FAILED")
echo "ID: '$ID'"

# Test 3: Test git @ _path
echo "Test 3: git @ _path"
PATH_VAL=$(git @ _path 2>/dev/null || echo "FAILED")
echo "Path: '$PATH_VAL'"

# Test 4: Test git @ _trunk
echo "Test 4: git @ _trunk"
TRUNK=$(git @ _trunk 2>/dev/null || echo "FAILED")
echo "Trunk: '$TRUNK'"

# Test 5: Test git @ version -t
echo "Test 5: git @ version -t"
TAG=$(git @ version -t 2>/dev/null || echo "FAILED")
echo "Tag: '$TAG'"

# Test 6: Test git @ wip
echo "Test 6: git @ wip"
WIP=$(git @ wip 2>/dev/null || echo "FAILED")
echo "WIP: '$WIP'"

echo "Individual command tests completed!" 