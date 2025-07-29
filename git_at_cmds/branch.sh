#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ branch [<branch-name>] [options]

DESCRIPTION:
  Manage the working branch configuration for GitAT workflow.
  Distinguishes between current git branch and configured working branch.

OPTIONS:
  (no options)           Show configured working branch
  <branch-name>          Set working branch to specified name
  -c, --current          Show current git branch
  -s, --set, .           Set working branch to current git branch
  -n, --new              Create new working branch with timestamp
  -h, --help             Show this help

EXAMPLES:
  git @ branch                    # Show configured working branch
  git @ branch feature-auth       # Set working branch to "feature-auth"
  git @ branch -c                 # Show current git branch
  git @ branch -s                 # Set working branch to current branch
  git @ branch .                  # Same as -s (set to current)
  git @ branch -n                 # Create new working branch

WORKFLOW:
  Working branch is the branch you're supposed to be working on.
  Current branch is the branch you're actually on.
  Use 'git @ work' to switch to your working branch.

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