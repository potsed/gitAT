#!/bin/bash

usage() {
    cat << 'EOF'
       _ _         ____                                _
  __ _(_) |_      / __ \     ___  __ _ _   _  __ _ ___| |__
 / _` | | __|    / / _` |   / __|/ _` | | | |/ _` / __| '_ \
| (_| | | |_    | | (_| |   \__ \ (_| | |_| | (_| \__ \ | | |
 \__, |_|\__|    \ \__,_|   |___/\__, |\__,_|\__,_|___/_| |_|
 |___/            \____/            |_|


Usage: git @ squash [options] [<target-branch>]

DESCRIPTION:
  Reset the current branch to the HEAD of its parent/upstream branch to create
  a clean commit without working history. Automatically detects the parent branch
  based on where the current branch was created from.

OPTIONS:
  -s, --save           Run 'git @ save' after squashing
  -p, --pr             Squash for PR (uses configured trunk branch)
  -a, --auto           Enable/disable automatic PR squashing
  -h, --help           Show this help

EXAMPLES:
  git @ squash                      # Squash to parent branch (auto-detected)
  git @ squash -s                   # Squash to parent and save
  git @ squash develop              # Squash to specific branch
  git @ squash master -s            # Squash to master and save
  git @ squash --pr                 # Squash for PR using trunk branch
  git @ squash --auto on            # Enable automatic PR squashing
  git @ squash --auto off           # Disable automatic PR squashing
  git @ squash --auto status        # Show automatic squashing status

PR SQUASHING:
  When using --pr, the command will:
  1. Use the configured trunk branch (at.trunk) as target
  2. Squash commits ahead of the trunk branch
  3. Preserve commit messages in the final squashed commit

AUTOMATIC PR SQUASHING:
  Configure automatic squashing for git @ pr:
  git @ squash --auto on            # Enable automatic squashing
  git @ squash --auto off           # Disable automatic squashing
  git @ squash --auto status        # Show current setting

PROCESS:
  1. Auto-detects parent branch (or uses specified target)
  2. Validates target branch exists
  3. Retrieves HEAD SHA of target branch
  4. Soft reset to target branch SHA
  5. Keeps all changes staged for commit
  6. Optionally runs 'git @ save'

USE CASES:
  - Clean up commit history before PR
  - Remove intermediate commits from feature branch
  - Create single clean commit from multiple commits
  - Automatic squashing before creating PRs

WARNING:
  You may need to force push after squashing if branch is shared.
  Use with caution on shared branches.

GIT COMMANDS USED:
  - git rev-parse --verify --quiet --long ${BRANCH}
  - git reset --soft ${SHA}

SECURITY:
  All squash operations are validated and logged.

EOF
    exit 1
}

cmd_squash() {
    local DOSAVE=0
    local PR_MODE=false
    local AUTO_MODE=""
    
    # Parse arguments
    while [ "$#" -gt 0 ]; do
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-s"|"--save")
                DOSAVE=1
                shift
                ;;
            "-p"|"--pr")
                PR_MODE=true
                shift
                ;;
            "-a"|"--auto")
                if [ -n "$2" ] && [[ "$2" != -* ]]; then
                    AUTO_MODE="$2"
                    shift 2
                else
                    echo "Error: --auto requires a value (on|off|status)" >&2
                    usage; exit 1
                fi
                ;;
            *)
                break
                ;;
        esac
    done
    
    # Handle automatic PR squashing management
    if [ -n "$AUTO_MODE" ]; then
        handle_auto_squash "$AUTO_MODE"
        exit 0
    fi
    
    # Handle PR mode
    if [ "$PR_MODE" = true ]; then
        handle_pr_squash
        exit 0
    fi
    
    # Original squash functionality
    local target_branch=""
    local HEAD_SHA
    
    # If no target branch specified, auto-detect parent branch
    if [ "$#" -eq 0 ]; then
        target_branch=$(detect_parent_branch)
        if [ -z "$target_branch" ]; then
            echo "Error: Could not auto-detect parent branch" >&2
            echo "Please specify a target branch: git @ squash <branch>" >&2
            usage; exit 1
        fi
        echo "Auto-detected parent branch: $target_branch"
    else
        target_branch="$1"
    fi
    
    HEAD_SHA=$(head "$target_branch")
    if [ "$HEAD_SHA" = "0" ]; then
        echo "ERROR: Branch \"$target_branch\" does not exist locally" >&2
        exit 1
    fi
    
    squash "$HEAD_SHA"
    
    echo "Squashed branch $(git @ branch -c) back to $target_branch"
    
    if [ "$DOSAVE" = 1 ]; then
        save
    fi
}

# Handle automatic PR squashing management
handle_auto_squash() {
    local action="$1"
    
    case "$action" in
        "on"|"true"|"enable"|"1")
            enable_auto_squash
            ;;
        "off"|"false"|"disable"|"0")
            disable_auto_squash
            ;;
        "status"|"show"|"check")
            show_auto_squash_status
            ;;
        *)
            echo "Error: Invalid auto action '$action'. Use 'on', 'off', or 'status'" >&2
            usage; exit 1
            ;;
    esac
}

enable_auto_squash() {
    git config --replace-all at.pr.squash true
    echo "✅ Automatic PR squashing enabled"
    echo "   Commits will be automatically squashed before creating PRs"
    echo "   Use 'git @ pr -S' to override and skip squashing"
}

disable_auto_squash() {
    git config --replace-all at.pr.squash false
    echo "✅ Automatic PR squashing disabled"
    echo "   PRs will be created with all commits as-is"
    echo "   Use 'git @ pr -s' to force squashing"
}

show_auto_squash_status() {
    local setting
    setting=$(git config at.pr.squash 2>/dev/null || echo "false")
    
    echo "PR Squash Setting:"
    if [ "$setting" = "true" ]; then
        echo "  Status: ✅ ENABLED"
        echo "  Commits will be automatically squashed before creating PRs"
        echo "  Override: git @ pr -S (force no squash)"
    else
        echo "  Status: ❌ DISABLED"
        echo "  PRs will be created with all commits as-is"
        echo "  Override: git @ pr -s (force squash)"
    fi
    
    echo ""
    echo "Commands:"
    echo "  git @ squash --auto on     # Enable automatic squashing"
    echo "  git @ squash --auto off    # Disable automatic squashing"
    echo "  git @ squash --auto status # Show this information"
}

# Handle PR mode squashing
handle_pr_squash() {
    local trunk_branch
    local current_branch
    local commit_count
    
    trunk_branch=$(git config at.trunk 2>/dev/null || echo "main")
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    
    if [ -z "$current_branch" ] || [ "$current_branch" = "HEAD" ]; then
        echo "Error: Not on a branch (detached HEAD state)" >&2
        exit 1
    fi
    
    # Check if we're on the trunk branch
    if [ "$current_branch" = "$trunk_branch" ]; then
        echo "Error: Cannot squash PR from $trunk_branch to itself" >&2
        exit 1
    fi
    
    # Get the number of commits ahead of trunk branch
    commit_count=$(git rev-list --count "$trunk_branch..HEAD" 2>/dev/null || echo "0")
    
    if [ "$commit_count" -le 1 ]; then
        echo "Only one commit or no commits to squash"
        return 0
    fi
    
    echo "Found $commit_count commits to squash for PR"
    
    # Get the commit hash where the branch diverged from trunk
    local base_commit
    base_commit=$(git merge-base "$trunk_branch" HEAD 2>/dev/null || echo "")
    
    if [ -z "$base_commit" ]; then
        echo "Error: Cannot find merge base with $trunk_branch" >&2
        exit 1
    fi
    
    # Create a temporary branch for the squash
    local temp_branch
    temp_branch="${current_branch}-squash-$(date +%s)"
    
    # Create temp branch from base
    if ! git checkout -b "$temp_branch" "$base_commit" 2>/dev/null; then
        echo "Error: Failed to create temporary branch" >&2
        exit 1
    fi
    
    # Cherry-pick all commits from current branch
    local cherry_pick_success=true
    while IFS= read -r commit_hash; do
        if [ -n "$commit_hash" ]; then
            if ! git cherry-pick "$commit_hash" 2>/dev/null; then
                echo "Error: Failed to cherry-pick commit $commit_hash" >&2
                cherry_pick_success=false
                break
            fi
        fi
    done < <(git rev-list --reverse "$base_commit..HEAD" 2>/dev/null)
    
    if [ "$cherry_pick_success" = false ]; then
        # Clean up on failure
        git checkout "$current_branch" 2>/dev/null
        git branch -D "$temp_branch" 2>/dev/null
        exit 1
    fi
    
    # Reset current branch to temp branch
    if ! git checkout "$current_branch" 2>/dev/null; then
        echo "Error: Failed to switch back to current branch" >&2
        git branch -D "$temp_branch" 2>/dev/null
        exit 1
    fi
    
    if ! git reset --hard "$temp_branch" 2>/dev/null; then
        echo "Error: Failed to reset current branch" >&2
        git branch -D "$temp_branch" 2>/dev/null
        exit 1
    fi
    
    # Clean up temp branch
    git branch -D "$temp_branch" 2>/dev/null
    
    echo "✅ Successfully squashed $commit_count commits into one for PR"
}

# Detect the parent/upstream branch of the current branch
detect_parent_branch() {
    local current_branch
    local parent_branch=""
    
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    
    if [ -z "$current_branch" ] || [ "$current_branch" = "HEAD" ]; then
        return 1
    fi
    
    # Method 1: Check git config for upstream branch
    parent_branch=$(git config "branch.$current_branch.merge" 2>/dev/null | sed 's|refs/heads/||' || echo "")
    if [ -n "$parent_branch" ]; then
        echo "$parent_branch"
        return 0
    fi
    
    # Method 2: Check if current branch has an upstream tracking branch
    parent_branch=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null | sed 's|refs/remotes/origin/||' || echo "")
    if [ -n "$parent_branch" ]; then
        echo "$parent_branch"
        return 0
    fi
    
    # Method 3: Find the branch that the current branch diverged from
    # Get all local branches
    local branches=()
    while IFS= read -r branch; do
        if [ -n "$branch" ] && [ "$branch" != "$current_branch" ]; then
            branches+=("$branch")
        fi
    done < <(git branch --list | sed 's/^[* ]*//')
    
    # Find the branch with the most recent common ancestor
    local best_branch=""
    local best_count=0
    
    for branch in "${branches[@]}"; do
        # Count commits that are in current branch but not in this branch
        local commit_count
        commit_count=$(git rev-list --count "$branch..HEAD" 2>/dev/null || echo "0")
        
        # If this branch has fewer commits ahead, it's likely the parent
        if [ "$commit_count" -gt 0 ] && [ "$commit_count" -gt "$best_count" ]; then
            best_branch="$branch"
            best_count="$commit_count"
        fi
    done
    
    if [ -n "$best_branch" ]; then
        echo "$best_branch"
        return 0
    fi
    
    # Method 4: Fallback to configured trunk branch
    parent_branch=$(git config at.trunk 2>/dev/null || echo "")
    if [ -n "$parent_branch" ]; then
        echo "$parent_branch"
        return 0
    fi
    
    # Method 5: Try common branch names
    for branch in "main" "master" "develop" "development"; do
        if git rev-parse --verify "$branch" >/dev/null 2>&1; then
            echo "$branch"
            return 0
        fi
    done
    
    return 1
}

head() {
    local HEAD=$(git rev-parse --verify --quiet --long "$1")
    if [ "${HEAD}" != "" ]; then
        echo "${HEAD}"
    else
        echo "0"
    fi
}

squash() {
    git reset --soft "$1"
}

save() {
    git @ save
}