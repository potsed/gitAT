#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ info

DESCRIPTION:
  Show comprehensive GitAT status and configuration.
  Displays current values for all GitAT settings and repository status.

OUTPUT:
  Displays comprehensive information including:
  - Configuration & Information (product, feature, version, etc.)
  - Git Repository Status (current branch, remote, changes)
  - Branch Information (status, protection)
  - Available Commands (categorized by function)
  - Quick Actions (common commands)

SECTIONS:
  📋 Configuration & Information - All GitAT settings
  🌿 Git Repository Status - Current git state
  📊 Branch Information - Branch status and protection
  🛠️  Available Commands - Categorized command list
  💡 Quick Actions - Common workflow commands

EXAMPLES:
  git @ info                    # Show comprehensive GitAT status

SECURITY:
  All info operations are validated and logged.

EOF
    exit 1
}

cmd_info() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    # Source all the required command files to get access to their functions
    local SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    
    # Source the command files to get access to their functions
    source "$SCRIPT_DIR/product.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/version.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/feature.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/issue.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/branch.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/_path.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/_trunk.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/wip.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/_id.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/_label.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/hash.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/changes.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/logs.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/save.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/work.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/release.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/squash.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/sweep.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/ignore.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/initlocal.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/initremote.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/_go.sh" 2>/dev/null || true

    # Get current git status
    local CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    local REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || echo "")
    local REMOTE_URL=$(git config --get remote.origin.url 2>/dev/null || echo "")
    local UNCOMMITTED_CHANGES=$(git status --porcelain 2>/dev/null | wc -l | tr -d ' ')
    local STASH_COUNT=$(git stash list 2>/dev/null | wc -l | tr -d ' ')

    # Call functions directly instead of using git @ commands
    local PROJECT=$(cmd_product 2>/dev/null || echo "")
    local VERSION=$(cmd_version 2>/dev/null || echo "")
    local TAG=$(cmd_version -t 2>/dev/null || echo "")
    local FEATURE=$(cmd_feature 2>/dev/null || echo "")
    local ISSUE=$(cmd_issue 2>/dev/null || echo "")
    local BRANCH=$(cmd_branch 2>/dev/null || echo "")
    local GITAT_PATH=$(cmd__path 2>/dev/null || echo "")
    local TRUNK=$(cmd__trunk 2>/dev/null || echo "")
    local WIP=$(cmd_wip 2>/dev/null || echo "")
    local GITAT_ID=$(cmd__id 2>/dev/null || echo "")
    local LABEL=$(cmd__label 2>/dev/null || echo "")

    cat << EOF

╔══════════════════════════════════════════════════════════════════════════════╗
║                              GitAT Status Report                             ║
╚══════════════════════════════════════════════════════════════════════════════╝

📋 Configuration & Information
─────────────────────────────────────────────────────────────────────────────────
│ Product Name        │ ${PROJECT:-<not set>}
│ Feature Name        │ ${FEATURE:-<not set>}
│ Version             │ ${VERSION:-<not set>}
│ Version Tag         │ ${TAG:-<not set>}
│ Working Branch      │ ${BRANCH:-<not set>}
│ Issue/Task ID       │ ${ISSUE:-<not set>}
│ WIP Branch          │ ${WIP:-<not set>}
│ Project ID          │ ${GITAT_ID:-<not set>}
│ Commit Label        │ ${LABEL:-<not set>}
│ Repository Path     │ ${GITAT_PATH:-<not set>}
│ Trunk Branch        │ ${TRUNK:-<not set>}

🌿 Git Repository Status
─────────────────────────────────────────────────────────────────────────────────
│ Current Branch      │ ${CURRENT_BRANCH:-<not in git repo>}
│ Repository Root     │ ${REPO_ROOT:-<not in git repo>}
│ Remote URL          │ ${REMOTE_URL:-<not set>}
│ Uncommitted Changes │ ${UNCOMMITTED_CHANGES:-0} files
│ Stash Count         │ ${STASH_COUNT:-0} stashes

📊 Branch Information
─────────────────────────────────────────────────────────────────────────────────
│ Branch Status       │ $(if [ "$CURRENT_BRANCH" != "$BRANCH" ] && [ -n "$BRANCH" ]; then echo "⚠️  On $CURRENT_BRANCH (should be on $BRANCH)"; else echo "✅ On correct branch"; fi)
│ Branch Protection   │ $(if [ "$CURRENT_BRANCH" = "master" ] || [ "$CURRENT_BRANCH" = "develop" ]; then echo "🛡️  Protected branch"; else echo "✅ Safe to work"; fi)

🛠️  Available Commands
─────────────────────────────────────────────────────────────────────────────────
│ Workflow           │ git @ save, git @ work, git @ wip
│ Version Management │ git @ version, git @ release
│ Branch Management  │ git @ branch, git @ master, git @ root, git @ sweep
│ Repository Setup   │ git @ _go, git @ initlocal, git @ initremote
│ Information        │ git @ info, git @ changes, git @ logs, git @ hash
│ Utilities          │ git @ _path, git @ _trunk, git @ _label, git @ _id

💡 Quick Actions
─────────────────────────────────────────────────────────────────────────────────
│ Switch to Work     │ git @ work
│ Save Changes       │ git @ save "message"
│ Check Changes      │ git @ changes
│ View Logs          │ git @ logs
│ Create Release     │ git @ release -m

EOF
    exit 0
}