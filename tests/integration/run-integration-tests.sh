#!/bin/bash

# Integration Tests Runner
# Tests the interaction between multiple GitAT commands and real Git operations

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$TEST_DIR")"

echo "🧪 Running GitAT Integration Tests"
echo "=================================="

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not in a git repository"
    echo "Please run this script from within a git repository"
    exit 1
fi

# Check if GitAT is available
if ! command -v git @ > /dev/null 2>&1; then
    echo "❌ Error: GitAT not found"
    echo "Please ensure git @ is available in your PATH"
    exit 1
fi

# Function to run a test
run_test() {
    local test_file="$1"
    local test_name="$(basename "$test_file" .sh)"
    
    echo ""
    echo "🔍 Running: $test_name"
    echo "----------------------------------------"
    
    if [ -x "$test_file" ]; then
        if "$test_file"; then
            echo "✅ $test_name passed"
            return 0
        else
            echo "❌ $test_name failed"
            return 1
        fi
    else
        echo "⚠️  $test_name is not executable, skipping"
        return 0
    fi
}

# Find all test files
test_files=()
while IFS= read -r -d '' file; do
    test_files+=("$file")
done < <(find "$SCRIPT_DIR" -name "test-*.sh" -type f -print0 | sort -z)

if [ ${#test_files[@]} -eq 0 ]; then
    echo "❌ No test files found in $SCRIPT_DIR"
    exit 1
fi

echo "Found ${#test_files[@]} integration test(s):"
for test_file in "${test_files[@]}"; do
    echo "  - $(basename "$test_file")"
done

# Run tests
passed=0
failed=0

for test_file in "${test_files[@]}"; do
    if run_test "$test_file"; then
        ((passed++))
    else
        ((failed++))
    fi
done

echo ""
echo "📊 Integration Test Results"
echo "=========================="
echo "✅ Passed: $passed"
echo "❌ Failed: $failed"
echo "📋 Total: $((passed + failed))"

if [ $failed -eq 0 ]; then
    echo ""
    echo "🎉 All integration tests passed!"
    exit 0
else
    echo ""
    echo "💥 Some integration tests failed!"
    exit 1
fi 