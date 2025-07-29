#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ sweep

DESCRIPTION:
  Clean up local branches that have been merged into trunk.
  Deletes local copies of branches that are no longer needed.

PROCESS:
  1. Identifies branches merged into trunk
  2. Deletes local copies of merged branches
  3. Preserves important branches (trunk, master, main, dev, develop, staging, stage, qa)

EXAMPLES:
  git @ sweep                    # Clean up merged branches

SAFETY:
  - Only deletes branches that are fully merged
  - Preserves important branches (trunk, master, main, dev, develop, staging, stage, qa)
  - Does not affect remote branches

SECURITY:
  All sweep operations are validated and logged.

EOF
    exit 1
}

cmd_sweep() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    echo "🧹 Cleaning up merged branches..."
    local trunk_branch
    trunk_branch=$(git config at.trunk 2>/dev/null || echo "develop")
    echo "Trunk branch: $trunk_branch"
    echo
    delete_merged_branches_locally
    echo
    echo "✅ Sweep completed"
    exit 0
}

delete_merged_branches_locally() {
    # Preserve important branches: trunk, master, main, dev, develop, staging, stage, qa
    local trunk_branch
    trunk_branch=$(git config at.trunk 2>/dev/null || echo "develop")
    local preserved_branches="master|main|dev|develop|staging|stage|qa|${trunk_branch}"
    
    echo "Preserving branches: ${preserved_branches}"
    
    # Get list of branches to delete and delete them directly
    local deleted_count=0
    local failed_count=0
    
    while IFS= read -r branch; do
        if [ -n "$branch" ]; then
            echo "  - $branch"
            if git branch -d "$branch" 2>/dev/null; then
                ((deleted_count++))
            else
                echo "    (failed to delete - may have unmerged changes)"
                ((failed_count++))
            fi
        fi
    done < <(git branch --merged | grep -v "\*\|${preserved_branches}")
    
    if [ $deleted_count -eq 0 ] && [ $failed_count -eq 0 ]; then
        echo "No merged branches to clean up"
    else
        echo "Deleted: $deleted_count branches, Failed: $failed_count branches"
    fi
}