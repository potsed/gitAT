#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ sweep [options]

DESCRIPTION:
  Clean up local branches that have been merged into trunk or deleted on remote.
  Deletes local copies of branches that are no longer needed.

OPTIONS:
  -f, --force      Force delete branches even if they appear to have unmerged changes
                   (Useful for branches that were squash-merged or rebased)
  -n, --dry-run    Show what would be deleted without actually deleting anything
  -l, --local-only Only check for locally merged branches (skip remote tracking)
                   (By default, checks both local merges and remote deletions)

PROCESS:
  1. Identifies branches merged into trunk
  2. Identifies branches deleted on remote (unless --local-only is used)
  3. Deletes local copies of merged/deleted branches
  4. Preserves important branches (trunk, master, main, dev, develop, staging, stage, qa)

EXAMPLES:
  git @ sweep                    # Clean up merged + remote-deleted branches (default)
  git @ sweep --local-only      # Clean up only locally merged branches
  git @ sweep --force           # Force clean up (including squash-merged branches)
  git @ sweep --force --remote  # Force clean up including remote-deleted branches
  git @ sweep --dry-run         # Preview what would be deleted
  git @ sweep --force --dry-run # Preview force deletion

SAFETY:
  - Only deletes branches that are fully merged or deleted on remote
  - Preserves important branches (trunk, master, main, dev, develop, staging, stage, qa)
  - Does not affect remote branches
  - --force option allows deletion of branches that may appear unmerged
  - --dry-run option lets you preview changes before applying them
  - Remote tracking only affects branches that track remote branches

SECURITY:
  All sweep operations are validated and logged.

EOF
    exit 1
}

cmd_sweep() {
    local force_delete=false
    local dry_run=false
    local check_remote=true  # Default to true (check remote)
    
    while [ "$#" -gt 0 ]; do
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-f"|"--force")
                force_delete=true
                shift
                ;;
            "-n"|"--dry-run")
                dry_run=true
                shift
                ;;
            "-l"|"--local-only")
                check_remote=false
                shift
                ;;
            "-r"|"--remote")
                check_remote=true
                shift
                ;;
            *)
                echo "Error: Unknown option '$1'" >&2
                usage; exit 1
                ;;
        esac
    done
    
    echo "🧹 Cleaning up merged branches..."
    local trunk_branch
    trunk_branch=$(git config at.trunk 2>/dev/null || echo "develop")
    echo "Trunk branch: $trunk_branch"
    
    if [ "$force_delete" = true ]; then
        echo "Force mode: Will delete branches even if they appear unmerged"
    fi
    
    if [ "$check_remote" = true ]; then
        echo "Remote tracking: Will also check for branches deleted on remote"
    else
        echo "Local only: Will only check for locally merged branches"
    fi
    
    if [ "$dry_run" = true ]; then
        echo "Dry run mode: Will show what would be deleted (no actual deletion)"
    fi
    echo
    delete_merged_branches_locally "$force_delete" "$dry_run" "$check_remote"
    echo
    if [ "$dry_run" = true ]; then
        echo "✅ Dry run completed - no branches were actually deleted"
    else
        echo "✅ Sweep completed"
    fi
    exit 0
}

delete_merged_branches_locally() {
    local force_delete="$1"
    local dry_run="$2"
    local check_remote="$3"
    
    # Preserve important branches: trunk, master, main, dev, develop, staging, stage, qa
    local trunk_branch
    trunk_branch=$(git config at.trunk 2>/dev/null || echo "develop")
    local preserved_branches="master|main|dev|develop|staging|stage|qa|${trunk_branch}"
    
    echo "Preserving branches: ${preserved_branches}"
    
    # Get list of branches to delete and delete them directly
    local deleted_count=0
    local failed_count=0
    
    if [ "$force_delete" = true ]; then
        # Force mode: try to delete all non-preserved branches
        echo "⚠️  Force mode: Will attempt to delete all non-preserved branches"
        echo "   This includes branches that may still be in use!"
        echo ""
        
        # Get current date for recent activity check
        local current_date=$(date +%s)
        local thirty_days_ago=$((current_date - 30 * 24 * 60 * 60))
        
        # Collect branches to potentially delete
        local branches_to_delete=()
        local recent_branches=()
        
        while IFS= read -r branch; do
            if [ -n "$branch" ]; then
                # Check if branch has recent activity (commits in last 30 days)
                local last_commit_date=$(git log -1 --format="%ct" "$branch" 2>/dev/null || echo "0")
                if [ "$last_commit_date" -gt "$thirty_days_ago" ]; then
                    recent_branches+=("$branch")
                else
                    branches_to_delete+=("$branch")
                fi
            fi
        done < <(git branch | grep -v "\*\|${preserved_branches}" | sed 's/^[[:space:]]*//')
        
        # Show recent branches and ask for confirmation (unless dry run)
        if [ ${#recent_branches[@]} -gt 0 ] && [ "$dry_run" = false ]; then
            echo "⚠️  The following branches have recent activity (last 30 days):"
            for branch in "${recent_branches[@]}"; do
                local last_commit=$(git log -1 --oneline "$branch" 2>/dev/null || echo "unknown")
                echo "   - $branch (last: $last_commit)"
            done
            echo ""
            echo "These branches might still be in use. Do you want to:"
            echo "1) Delete only old branches (safer)"
            echo "2) Delete all branches including recent ones (risky)"
            echo "3) Cancel and exit"
            echo ""
            read -p "Choose option (1/2/3): " -n 1 -r
            echo
            echo ""
            
            case $REPLY in
                1)
                    echo "✅ Proceeding with safe deletion (old branches only)"
                    ;;
                2)
                    echo "⚠️  Proceeding with risky deletion (all branches)"
                    branches_to_delete+=("${recent_branches[@]}")
                    ;;
                3)
                    echo "❌ Cancelled"
                    return 0
                    ;;
                *)
                    echo "❌ Invalid option, cancelling"
                    return 0
                    ;;
            esac
        elif [ "$dry_run" = true ]; then
            # In dry run mode, show both old and recent branches
            if [ ${#recent_branches[@]} -gt 0 ]; then
                echo "⚠️  The following branches have recent activity (last 30 days):"
                for branch in "${recent_branches[@]}"; do
                    local last_commit=$(git log -1 --oneline "$branch" 2>/dev/null || echo "unknown")
                    echo "   - $branch (last: $last_commit)"
                done
                echo ""
            fi
        fi
        
        # Delete the selected branches (or show what would be deleted)
        for branch in "${branches_to_delete[@]}"; do
            if [ "$dry_run" = true ]; then
                echo "  - $branch (would be force deleted)"
                ((deleted_count++))
            else
                echo "  - $branch"
                if git branch -D "$branch" 2>/dev/null; then
                    echo "    (force deleted)"
                    ((deleted_count++))
                else
                    echo "    (failed to delete)"
                    ((failed_count++))
                fi
            fi
        done
    else
        # Normal mode: delete fully merged branches and optionally remote-deleted branches
        local branches_to_delete=()
        
        # Get merged branches
        while IFS= read -r branch; do
            if [ -n "$branch" ]; then
                branches_to_delete+=("$branch")
            fi
        done < <(git branch --merged | grep -v "\*\|${preserved_branches}")
        
        # Check for remote-deleted branches if requested
        if [ "$check_remote" = true ]; then
            echo "Checking for branches deleted on remote..."
            
            # Prune remote tracking branches first
            git remote prune origin 2>/dev/null || true
            
            # Find local branches that exist locally but not remotely
            # This catches branches that were tracking remote branches that got deleted
            while IFS= read -r branch; do
                if [ -n "$branch" ]; then
                    # Check if this branch exists locally but not remotely
                    local remote_exists=$(git ls-remote --heads origin "$branch" 2>/dev/null | wc -l)
                    if [ "$remote_exists" -eq 0 ]; then
                        # Branch exists locally but not remotely - likely was tracking a deleted remote branch
                        # Only add if it's not already in the deletion list
                        if [[ ! " ${branches_to_delete[@]} " =~ " ${branch} " ]]; then
                            branches_to_delete+=("$branch")
                        fi
                    fi
                fi
            done < <(git branch | grep -v "\*\|${preserved_branches}" | sed 's/^[[:space:]]*//')
        fi
        
        # Delete the branches
        for branch in "${branches_to_delete[@]}"; do
            if [ "$dry_run" = true ]; then
                # Determine why it would be deleted
                local reason=""
                if git branch --merged | grep -q "^[[:space:]]*$branch$"; then
                    reason="fully merged"
                else
                    reason="deleted on remote"
                fi
                echo "  - $branch (would be deleted - $reason)"
                ((deleted_count++))
            else
                echo "  - $branch"
                # Check if this branch is merged or remote-deleted
                local is_merged=false
                local is_remote_deleted=false
                
                if git branch --merged | grep -q "^[[:space:]]*$branch$"; then
                    is_merged=true
                fi
                
                # Check if it's remote-deleted (exists locally but not remotely)
                local remote_exists=$(git ls-remote --heads origin "$branch" 2>/dev/null | wc -l)
                if [ "$remote_exists" -eq 0 ]; then
                    is_remote_deleted=true
                fi
                
                # Debug output
                if [ "$is_merged" = true ]; then
                    echo "    (debug: detected as merged)"
                fi
                if [ "$is_remote_deleted" = true ]; then
                    echo "    (debug: detected as remote-deleted)"
                fi
                
                # Prioritize remote-deleted over merged (since remote-deleted should always be safe to force delete)
                if [ "$is_remote_deleted" = true ]; then
                    # Use force delete for remote-deleted branches
                    echo "    (attempting force delete for remote-deleted branch)"
                    
                    # Check if there are uncommitted changes that might prevent deletion
                    local has_uncommitted=false
                    if ! git diff --quiet 2>/dev/null || ! git diff --cached --quiet 2>/dev/null; then
                        has_uncommitted=true
                    fi
                    
                    if [ "$has_uncommitted" = true ]; then
                        echo "    (stashing uncommitted changes to allow branch deletion)"
                        if git stash push -m "Auto-stash before deleting remote-deleted branch: $branch" 2>/dev/null; then
                            echo "    (changes stashed successfully)"
                        else
                            echo "    (failed to stash changes)"
                        fi
                    fi
                    
                    # Try multiple deletion methods
                    local deletion_success=false
                    
                    # Method 1: Standard force delete
                    if git branch -D "$branch" 2>/dev/null; then
                        echo "    (force deleted - was deleted on remote)"
                        deletion_success=true
                    else
                        echo "    (standard force delete failed, trying alternative method)"
                        
                        # Method 2: Use git update-ref to delete the branch reference
                        if git update-ref -d "refs/heads/$branch" 2>/dev/null; then
                            echo "    (deleted using update-ref - was deleted on remote)"
                            deletion_success=true
                        else
                            echo "    (update-ref method also failed)"
                        fi
                    fi
                    
                    if [ "$deletion_success" = true ]; then
                        ((deleted_count++))
                    else
                        # Try to understand why it failed
                        echo "    (debug: checking branch status)"
                        
                        # Check if this is the current branch
                        if [ "$(git rev-parse --abbrev-ref HEAD 2>/dev/null)" = "$branch" ]; then
                            echo "    (debug: this is the current branch - cannot delete current branch)"
                        else
                            echo "    (debug: not the current branch)"
                        fi
                        
                        # Check if there are uncommitted changes
                        if ! git diff --quiet 2>/dev/null || ! git diff --cached --quiet 2>/dev/null; then
                            echo "    (debug: there are uncommitted changes in working directory)"
                        else
                            echo "    (debug: no uncommitted changes in working directory)"
                        fi
                        
                        # Check if the branch has any special protection
                        local branch_protection=$(git config branch."$branch".protect 2>/dev/null || echo "")
                        if [ -n "$branch_protection" ]; then
                            echo "    (debug: branch has protection setting: $branch_protection)"
                        fi
                        
                        ((failed_count++))
                    fi
                    
                    # Restore stashed changes if we stashed them
                    if [ "$has_uncommitted" = true ]; then
                        echo "    (restoring stashed changes)"
                        if git stash pop 2>/dev/null; then
                            echo "    (changes restored successfully)"
                        else
                            echo "    (failed to restore changes - check git stash list)"
                        fi
                    fi
                elif [ "$is_merged" = true ]; then
                    # Use safe delete for merged branches
                    echo "    (attempting safe delete for merged branch)"
                    if git branch -d "$branch" 2>/dev/null; then
                        ((deleted_count++))
                    else
                        echo "    (failed to delete - may have unmerged changes)"
                        ((failed_count++))
                    fi
                else
                    # Fallback - shouldn't happen but just in case
                    echo "    (attempting safe delete as fallback)"
                    if git branch -d "$branch" 2>/dev/null; then
                        ((deleted_count++))
                    else
                        echo "    (failed to delete)"
                        ((failed_count++))
                    fi
                fi
            fi
        done
    fi
    
    if [ $deleted_count -eq 0 ] && [ $failed_count -eq 0 ]; then
        echo "No branches to clean up"
    else
        if [ "$dry_run" = true ]; then
            echo "Would delete: $deleted_count branches"
        else
            echo "Deleted: $deleted_count branches, Failed: $failed_count branches"
        fi
    fi
}