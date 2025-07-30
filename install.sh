#!/bin/bash

# GitAT Installation Script
# This script installs GitAT as a proper Git extension

set -e

echo "🚀 Installing GitAT..."

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Build the Go binary
echo "📦 Building GitAT binary..."
cd "$SCRIPT_DIR"
go build -o git-@ ./cmd/gitat

# Install the Go binary
echo "📥 Installing binary to /usr/local/bin..."
sudo cp git-@ /usr/local/bin/
sudo chmod +x /usr/local/bin/git-@



echo "✅ GitAT installed successfully!"
echo ""
echo "Usage:"
echo "  git @ --help          # Show help"
echo "  git @ work feature    # Create work branch"
echo "  git @ save \"message\"  # Save changes"
echo ""
echo "Git extension 'git @' is now available from anywhere!"
echo "For more information, visit: https://github.com/potsed/gitAT" 