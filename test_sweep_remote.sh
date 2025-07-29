#!/bin/bash

echo "🧪 Testing Enhanced Sweep with Remote Tracking (Default)"
echo "======================================================"
echo ""

# Check if we're in a git repo
if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In git repository"
echo ""

# Check current branch
current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
echo "Current branch: $current_branch"
echo ""

# Check if remote is configured
if git remote get-url origin >/dev/null 2>&1; then
    echo "✅ Remote origin configured"
    remote_url=$(git remote get-url origin)
    echo "   URL: $remote_url"
else
    echo "⚠️  No remote origin configured"
    echo "   Remote tracking features won't work without a remote"
fi
echo ""

# Show all branches with their tracking info
echo "📋 All branches with tracking information:"
git branch -vv 2>/dev/null || echo "Could not list branches"
echo ""

# Test the different sweep modes
echo "🧹 Testing Sweep Modes:"
echo "======================="
echo ""

echo "1️⃣ Default sweep (merged + remote-deleted branches):"
echo "   git @ sweep"
echo ""

echo "2️⃣ Local-only sweep (merged branches only):"
echo "   git @ sweep --local-only"
echo ""

echo "3️⃣ Force sweep (including squash-merged):"
echo "   git @ sweep --force"
echo ""

echo "4️⃣ Force sweep with local-only:"
echo "   git @ sweep --force --local-only"
echo ""

echo "5️⃣ Dry run modes (preview only):"
echo "   git @ sweep --dry-run"
echo "   git @ sweep --local-only --dry-run"
echo "   git @ sweep --force --dry-run"
echo ""

# Check for branches that might be affected
echo "🔍 Analyzing branches:"
echo "======================"
echo ""

# Check merged branches
echo "✅ Branches merged into master:"
merged_branches=$(git branch --merged master 2>/dev/null || echo "")
if [ -n "$merged_branches" ]; then
    echo "$merged_branches"
else
    echo "(none found)"
fi
echo ""

# Check branches with remote tracking
echo "🌐 Branches with remote tracking:"
tracking_branches=$(git branch -vv | grep -E '\[origin/' | awk '{print $1}' | sed 's/^[[:space:]]*//' || echo "")
if [ -n "$tracking_branches" ]; then
    echo "$tracking_branches"
else
    echo "(none found)"
fi
echo ""

# Check for potentially deleted remote branches
echo "🗑️  Checking for branches that might be deleted on remote:"
if git remote get-url origin >/dev/null 2>&1; then
    # Prune remote tracking branches
    git remote prune origin 2>/dev/null || true
    
    # Check for local branches tracking non-existent remote branches
    while IFS= read -r branch; do
        if [ -n "$branch" ]; then
            # Check if remote branch exists
            remote_exists=$(git ls-remote --heads origin "$branch" 2>/dev/null | wc -l)
            if [ "$remote_exists" -eq 0 ]; then
                echo "   - $branch (remote branch deleted)"
            fi
        fi
    done < <(git branch -vv | grep -E '\[origin/' | awk '{print $1}' | sed 's/^[[:space:]]*//')
else
    echo "   (no remote configured)"
fi
echo ""

echo "💡 Usage Examples:"
echo "=================="
echo ""
echo "• Clean up merged + remote-deleted branches (default):"
echo "  git @ sweep"
echo ""
echo "• Clean up only locally merged branches:"
echo "  git @ sweep --local-only"
echo ""
echo "• Force clean up (including squash-merged):"
echo "  git @ sweep --force"
echo ""
echo "• Preview what would be deleted:"
echo "  git @ sweep --dry-run"
echo ""
echo "🎯 Default behavior is perfect for:"
echo "   • Regular cleanup after team collaboration"
echo "   • Branches merged via GitHub/GitLab web interface"
echo "   • Branches deleted after merge (auto-delete enabled)"
echo "   • Keeping local branches in sync with remote"
echo ""
echo "🔧 Use --local-only when:"
echo "   • You want to be more conservative"
echo "   • You're working offline"
echo "   • You only want to clean up locally merged branches" 