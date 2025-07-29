#!/bin/bash

# Source security utilities if available
if [ -f "$(dirname "${BASH_SOURCE[0]}")/_security.sh" ]; then
    source "$(dirname "${BASH_SOURCE[0]}")/_security.sh"
fi

usage() {
    echo 'Usage: git @ save [<message>]'
    echo '  Save current changes with proper validation'
    echo '  Message is optional and will be validated'
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
    
    # Validate any provided message
    if [ "$#" -gt 0 ]; then
        if ! validate_input "$*"; then
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
    current=$(git @ branch -c 2>/dev/null || echo "")
    branch=$(git @ branch 2>/dev/null || echo "")
    repo_path=$(git @ _path 2>/dev/null || echo "")

    # Validate we're in a git repository
    if [ -z "$repo_path" ]; then
        echo "Error: Not in a git repository" >&2
        exit 1
    fi

    # Check permissions
    if ! check_permissions "save" "$repo_path"; then
        echo "Error: Insufficient permissions to save changes" >&2
        exit 1
    fi

    # If no working branch is set, set it to current branch
    if [ -z "$branch" ]; then
        echo "No working branch configured. Setting current branch as working branch..."
        git @ branch "$current"
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

    if [ "$#" -eq 1 ]; then
        git @ _label "$1"
    fi

    local original_pwd
    local message
    
    original_pwd=$(pwd)
    message=$(git @ _label 2>/dev/null || echo "Update")
    
    # Validate path before changing directory
    if ! validate_path "$repo_path"; then
        echo "Error: Invalid repository path" >&2
        exit 1
    fi

    cd "$repo_path" || {
        echo "Error: Cannot change to repository directory" >&2
        exit 1
    }
    
    # Use safe command execution
    if safe_execute "git" "add -p"; then
        if safe_execute "git" "commit -m \"$message\""; then
            echo "Changes saved successfully"
        else
            echo "Error: Failed to commit changes" >&2
            exit 1
        fi
    else
        echo "Error: Failed to add changes" >&2
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