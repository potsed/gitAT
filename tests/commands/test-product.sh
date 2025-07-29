#!/bin/bash

echo "Testing product configuration key change..."

# Test 1: Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In a git repository"

# Test 2: Set product name using new configuration key
echo "Test 2: Setting product name"
if git @ product "TestProduct"; then
    echo "✅ Product name set successfully"
else
    echo "❌ Failed to set product name"
    exit 1
fi

# Test 3: Verify the configuration key is at.product
echo "Test 3: Verifying configuration key"
PRODUCT_VALUE=$(git config at.product 2>/dev/null || echo "")
if [ "$PRODUCT_VALUE" = "TestProduct" ]; then
    echo "✅ Configuration key is at.product with value: $PRODUCT_VALUE"
else
    echo "❌ Configuration key issue - expected 'TestProduct', got '$PRODUCT_VALUE'"
    exit 1
fi

# Test 4: Verify label generation uses at.product
echo "Test 4: Testing label generation"
LABEL=$(git @ _label 2>/dev/null || echo "")
if [[ "$LABEL" == *"TestProduct"* ]]; then
    echo "✅ Label generation uses at.product: $LABEL"
else
    echo "❌ Label generation issue - expected to contain 'TestProduct', got: $LABEL"
    exit 1
fi

# Test 5: Verify project ID generation uses at.product
echo "Test 5: Testing project ID generation"
PROJECT_ID=$(git @ _id 2>/dev/null || echo "")
if [[ "$PROJECT_ID" == *"TestProduct"* ]]; then
    echo "✅ Project ID generation uses at.product: $PROJECT_ID"
else
    echo "❌ Project ID generation issue - expected to contain 'TestProduct', got: $PROJECT_ID"
    exit 1
fi

# Test 6: Clean up
echo "Test 6: Cleaning up"
git config --unset at.product 2>/dev/null || true
echo "✅ Cleanup completed"

echo "All tests passed! Product configuration key change is working correctly." 