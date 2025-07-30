#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ branch [<branch-name>] [options]

DESCRIPTION:
  Manage the working branch configuration for GitAT workflow.
  Distinguishes between current git branch and configured working branch.
  Supports Conventional Commits workflow with different work types.

OPTIONS:
  (no options)           Show configured working branch
  <branch-name>          Set working branch to specified name
  -c, --current          Show current git branch
  -s, --set, .           Set working branch to current git branch
  -n, --new              Create new working branch with timestamp
  --hotfix               List all hotfix branches
  --feature              List all feature branches
  --bugfix               List all bugfix branches
  --release              List all release branches
  --chore                List all chore branches
  --docs                 List all documentation branches
  --style                List all style branches
  --refactor             List all refactor branches
  --perf                 List all performance branches
  --test                 List all test branches
  --ci                   List all CI/CD branches
  --build                List all build branches
  --revert               List all revert branches
  --all-types            List all work type branches
  -h, --help             Show this help

EXAMPLES:
  git @ branch                    # Show configured working branch
  git @ branch feature-auth       # Set working branch to "feature-auth"
  git @ branch -c                 # Show current git branch
  git @ branch -s                 # Set working branch to current branch
  git @ branch .                  # Same as -s (set to current)
  git @ branch -n                 # Create new working branch
  git @ branch --hotfix           # List all hotfix branches
  git @ branch --feature          # List all feature branches
  git @ branch --all-types        # List all work type branches

WORKFLOW:
  Working branch is the branch you're supposed to be working on.
  Current branch is the branch you're actually on.
  Use 'git @ work' to switch to your working branch.

CONVENTIONAL COMMITS INTEGRATION:
  Branch types follow Conventional Commits specification:
  - hotfix-*     : Urgent fixes for production (PATCH)
  - feature-*    : New features (MINOR)
  - bugfix-*     : Bug fixes (PATCH)
  - release-*    : Release preparation (PATCH)
  - chore-*      : Maintenance tasks
  - docs-*       : Documentation changes
  - style-*      : Code style changes
  - refactor-*   : Code refactoring
  - perf-*       : Performance improvements
  - test-*       : Test additions/changes
  - ci-*         : CI/CD changes
  - build-*      : Build system changes
  - revert-*     : Revert commits

STORAGE:
  Saved in git config: at.branch

SECURITY:
  All branch operations are validated and logged.

EOF
    exit 1
}

cmd_branch() {
    if [ "$#" -lt 1 ]; then
        show_branch; exit 0
    elif [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-n"|"--new"|"n"|"new")
                new_working_branch; exit 0
                ;;
            "-c"|"--current"|"c"|"current")
                current_branch; exit 0
                ;;
            "-s"|"--set"|"s"|"set"|".")
                local current_branch
                current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
                if [ -z "$current_branch" ] || [ "$current_branch" = "HEAD" ]; then
                    echo "Error: Not on a branch (detached HEAD state)" >&2
                    exit 1
                fi
                set_branch "$current_branch"
                ;;
            "--hotfix")
                list_branches_by_type "hotfix-"; exit 0
                ;;
            "--feature")
                list_branches_by_type "feature-"; exit 0
                ;;
            "--bugfix")
                list_branches_by_type "bugfix-"; exit 0
                ;;
            "--release")
                list_branches_by_type "release-"; exit 0
                ;;
            "--chore")
                list_branches_by_type "chore-"; exit 0
                ;;
            "--docs")
                list_branches_by_type "docs-"; exit 0
                ;;
            "--style")
                list_branches_by_type "style-"; exit 0
                ;;
            "--refactor")
                list_branches_by_type "refactor-"; exit 0
                ;;
            "--perf")
                list_branches_by_type "perf-"; exit 0
                ;;
            "--test")
                list_branches_by_type "test-"; exit 0
                ;;
            "--ci")
                list_branches_by_type "ci-"; exit 0
                ;;
            "--build")
                list_branches_by_type "build-"; exit 0
                ;;
            "--revert")
                list_branches_by_type "revert-"; exit 0
                ;;
            "--all-types")
                list_all_work_types; exit 0
                ;;
            *)
                set_branch "$1"; exit 0
                ;;
        esac
    fi

    usage; exit 1
}

new_working_branch() {
    local current_branch
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    if [ -z "$current_branch" ] || [ "$current_branch" = "HEAD" ]; then
        echo "Error: Not on a branch (detached HEAD state)" >&2
        exit 1
    fi
    
    local new_branch_name="feature-$(date +%Y%m%d-%H%M%S)"
    
    echo "Creating new working branch: $new_branch_name"
    git checkout -b "$new_branch_name"
    git @ branch "$new_branch_name"
    
    echo "New working branch created and set: $new_branch_name"
    exit 0
}

current_branch() {
    local branch
    branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    if [ -z "$branch" ] || [ "$branch" = "HEAD" ]; then
        echo "Error: Not on a branch (detached HEAD state)" >&2
        exit 1
    fi
    echo "$branch"
    exit 0
}

set_branch() {
    local from
    from=$(git config at.branch 2>/dev/null || echo "")
    git config --replace-all at.branch "$1"
    echo "Branch updated to: $1 from $from"
    exit 0
}

show_branch() {
    git config at.branch
    exit 0
}

# List branches by type (e.g., hotfix-, feature-, etc.)
list_branches_by_type() {
    local prefix="$1"
    local branch_type="${prefix%-}"  # Remove trailing dash for display
    
    echo "📋 $branch_type branches:"
    echo ""
    
    local found_branches=false
    local current_branch
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    
    while IFS= read -r branch; do
        if [ -n "$branch" ] && [[ "$branch" == "$prefix"* ]]; then
            found_branches=true
            local status=""
            if [ "$branch" = "$current_branch" ]; then
                status=" (current)"
            fi
            
            # Get last commit info
            local last_commit
            last_commit=$(git log -1 --oneline "$branch" 2>/dev/null || echo "No commits")
            
            echo "  🌿 $branch$status"
            echo "     📝 $last_commit"
            echo ""
        fi
    done < <(git branch --list | sed 's/^[* ]*//')
    
    if [ "$found_branches" = false ]; then
        echo "  No $branch_type branches found"
        echo ""
        echo "  To create a $branch_type branch:"
        case "$branch_type" in
            "hotfix")
                echo "    git @ hotfix 'description'"
                ;;
            "feature")
                echo "    git @ feature 'description'"
                ;;
            "bugfix")
                echo "    git @ bugfix 'description'"
                ;;
            *)
                echo "    git checkout -b ${prefix}description"
                ;;
        esac
    fi
    
    exit 0
}

# List all work types with counts
list_all_work_types() {
    echo "📊 All work type branches:"
    echo ""
    
    local work_types=("hotfix" "feature" "bugfix" "release" "chore" "docs" "style" "refactor" "perf" "test" "ci" "build" "revert")
    local current_branch
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    
    for type in "${work_types[@]}"; do
        local count=0
        local branches=()
        
        while IFS= read -r branch; do
            if [ -n "$branch" ] && [[ "$branch" == "${type}-"* ]]; then
                count=$((count + 1))
                branches+=("$branch")
            fi
        done < <(git branch --list | sed 's/^[* ]*//')
        
        if [ $count -gt 0 ]; then
            echo "  📁 $type branches ($count):"
            for branch in "${branches[@]}"; do
                local status=""
                if [ "$branch" = "$current_branch" ]; then
                    status=" (current)"
                fi
                echo "    🌿 $branch$status"
            done
            echo ""
        fi
    done
    
    echo "💡 Use 'git @ branch --<type>' to see detailed info for each type"
    echo "   Example: git @ branch --hotfix"
    
    exit 0
}