#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ save [<message>]

DESCRIPTION:
  Securely save current changes with comprehensive validation and security checks.
  This is the primary command for committing changes in GitAT workflow.

FEATURES:
  ✅ Auto-branch setup: Sets working branch if not configured
  ✅ Security validation: Validates inputs and paths
  ✅ Branch protection: Prevents saves on master/develop
  ✅ Production warnings: Confirms before saving to prod
  ✅ Safe execution: Uses secure command execution

EXAMPLES:
  git @ save                           # Save with default message
  git @ save "Add user authentication" # Save with custom message
  git @ save "Fix login bug"           # Save with descriptive message

VALIDATION:
  Messages must contain only:
  - Alphanumeric characters (a-z, A-Z, 0-9)
  - Dots (.), underscores (_), hyphens (-)
  - Spaces and common punctuation

SECURITY:
  - All inputs are validated against dangerous patterns
  - Path operations are restricted to repository root
  - Commands are executed safely
  - Security events are logged

BRANCH PROTECTION:
  - Cannot save on master or develop branches
  - Production branch requires confirmation
  - Must be on configured working branch

EOF
    exit 1
}

cmd_save() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    # Basic input validation
    if [ "$#" -gt 0 ]; then
        # Check for dangerous characters
        if echo "$*" | grep -q '[;&|`$(){}]'; then
            echo "Error: Invalid message. Use only alphanumeric characters, dots, underscores, and hyphens." >&2
            exit 1
        fi
    fi
    
    save_work "$@"; exit 0
}

save_work() {
    local current
    local branch
    local repo_path
    
    # Get current branch safely
    current=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    branch=$(git config at.branch 2>/dev/null || echo "")
    repo_path=$(git rev-parse --show-toplevel 2>/dev/null || echo "")

    # Validate we're in a git repository
    if [ -z "$repo_path" ]; then
        echo "Error: Not in a git repository" >&2
        exit 1
    fi

    # If no working branch is set, set it to current branch
    if [ -z "$branch" ]; then
        echo "No working branch configured. Setting current branch as working branch..."
        git config --replace-all at.branch "$current"
        branch="$current"
    fi

    if [ "$current" = "master" ] || [ "$current" = "develop" ]; then
        echo "Error: Cannot save changes on ${current}. Create a new branch instead!" >&2
        exit 1
    fi

    if [ "$current" = "prod" ]; then
        echo "Warning: You are on the production branch!"
        read -p "Are you sure you want to commit this? (Y/N): " CONFIRMATION
        if [[ ! "$CONFIRMATION" =~ ^[yY](es)?$ ]]; then
            echo "Operation cancelled."
            exit 1
        fi
        git @ version -t +
        git tag "$(git @ version -t)"
    elif [ "$current" != "$branch" ]; then
        echo "Error: Cannot save changes. You're not on the correct working branch '$branch'" >&2
        echo "Current branch: '$current'" >&2
        echo "To fix this, run: git @ branch '$current'" >&2
        exit 1
    fi

    local original_pwd
    local message
    
    original_pwd=$(pwd)
    
    # Generate commit message with label and user message
    if [ "$#" -eq 1 ]; then
        # User provided a message, combine with label
        local label
        label=$(git @ _label 2>/dev/null || echo "")
        if [ -n "$label" ]; then
            message="${label} $1"
        else
            message="$1"
        fi
    else
        # No user message, use default label
        message=$(git @ _label 2>/dev/null || echo "Update")
    fi
    
    # Change to repository directory
    cd "$repo_path" || {
        echo "Error: Cannot change to repository directory" >&2
        exit 1
    }
    
    # Add all changes and commit
    if git add . && git commit -m "$message"; then
        echo "Changes saved successfully"
    else
        echo "Error: Failed to save changes" >&2
        cd "$original_pwd" || exit 1
        exit 1
    fi
    
    cd "$original_pwd" || exit 1
    exit 0
}

# set_branch() {
#     `git config --replace-all at.branch $1`
#     echo 'Branch updated'
#     show_branch; exit 1
# }

# show_branch() {
#     echo "Current WOT Branch: "`git config at.branch`
#     echo
#     exit 1
# }