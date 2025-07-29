#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ hotfix [options] [<name>]

DESCRIPTION:
  Create a hotfix branch for urgent fixes that need to be deployed immediately.
  Creates a new branch from the trunk branch (master/main) and switches to it.

FEATURES:
  ✅ Creates hotfix branch from trunk (master/main)
  ✅ Saves current WIP state before switching
  ✅ Integrates with existing GitAT workflow
  ✅ Interactive name prompt if not provided
  ✅ Validates branch name and repository state

WORKFLOW:
  1. Saves current WIP state (if any)
  2. Switches to trunk branch
  3. Creates new hotfix branch from trunk
  4. Switches to hotfix branch
  5. Sets working branch to hotfix branch

OPTIONS:
  -n, --name <name>     Specify hotfix branch name
  -h, --help           Show this help message

EXAMPLES:
  git @ hotfix                    # Interactive name prompt
  git @ hotfix "fix-login-bug"    # Create hotfix with specific name
  git @ hotfix -n "security-patch" # Create hotfix with name option
  git @ hotfix --name "urgent-fix" # Create hotfix with name option

BRANCH NAMING:
  Format: hotfix-description (single hotfix- prefix)
  Examples:
    hotfix-fix-login-bug
    hotfix-security-patch
    hotfix-urgent-database-fix

INTEGRATION:
  - Uses git @ wip to save/restore work state
  - Uses git @ branch to set working branch
  - Uses git @ _trunk to get trunk branch name
  - Follows GitAT workflow patterns

EOF
    exit 1
}

cmd_hotfix() {
    local hotfix_name=""
    
    # Parse arguments
    while [ "$#" -gt 0 ]; do
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-n"|"--name")
                if [ -n "$2" ] && [[ "$2" != -* ]]; then
                    hotfix_name="$2"
                    shift 2
                else
                    echo "Error: --name requires a value" >&2
                    usage; exit 1
                fi
                ;;
            *)
                # If no name provided yet, use this as name
                if [ -z "$hotfix_name" ]; then
                    hotfix_name="$1"
                else
                    echo "Error: Unknown option '$1'" >&2
                    usage; exit 1
                fi
                shift
                ;;
        esac
    done
    
    # Validate we're in a git repository
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        echo "Error: Not in a git repository" >&2
        exit 1
    fi
    
    # Get current branch
    local current_branch
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    
    if [ -z "$current_branch" ] || [ "$current_branch" = "HEAD" ]; then
        echo "Error: Not on a branch (detached HEAD state)" >&2
        exit 1
    fi
    
    # Get trunk branch
    local trunk_branch
    trunk_branch=$(git config at.trunk 2>/dev/null || echo "main")
    
    # Validate trunk branch exists
    if ! git rev-parse --verify "$trunk_branch" >/dev/null 2>&1; then
        echo "Error: Trunk branch '$trunk_branch' does not exist" >&2
        echo "Please ensure the trunk branch exists or configure it with: git @ _trunk <branch>" >&2
        exit 1
    fi
    
    # Prompt for hotfix name if not provided
    if [ -z "$hotfix_name" ]; then
        echo "Creating hotfix branch from $trunk_branch"
        echo "Please enter a name for the hotfix branch:"
        echo "Format: description (e.g., fix-login-bug)"
        echo "Branch will be created as: hotfix-description"
        read -p "Hotfix description: " hotfix_name
        
        if [ -z "$hotfix_name" ]; then
            echo "Error: Hotfix description cannot be empty" >&2
            exit 1
        fi
    fi
    
    # Ensure hotfix name has the correct prefix
    if [[ "$hotfix_name" != "hotfix-"* ]]; then
        hotfix_name="hotfix-$hotfix_name"
    fi
    
    # Validate hotfix name
    if ! validate_branch_name "$hotfix_name"; then
        echo "Error: Invalid branch name '$hotfix_name'" >&2
        echo "Branch names must contain only alphanumeric characters, hyphens, underscores, and slashes" >&2
        exit 1
    fi
    
    # Check if hotfix branch already exists
    if git rev-parse --verify "$hotfix_name" >/dev/null 2>&1; then
        echo "Error: Branch '$hotfix_name' already exists" >&2
        exit 1
    fi
    
    # Check for uncommitted changes
    if ! git diff --quiet || ! git diff --cached --quiet; then
        echo "Warning: You have uncommitted changes"
        echo "These will be saved to WIP before creating hotfix branch"
        read -p "Continue? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Operation cancelled"
            exit 1
        fi
    fi
    
    # Save current WIP state
    echo "Saving current work state..."
    if ! git @ wip -s >/dev/null 2>&1; then
        echo "Warning: Failed to save WIP state, but continuing..."
    fi
    
    # Switch to trunk branch
    echo "Switching to trunk branch: $trunk_branch"
    if ! git checkout "$trunk_branch" 2>/dev/null; then
        echo "Error: Failed to switch to trunk branch '$trunk_branch'" >&2
        exit 1
    fi
    
    # Ensure trunk branch is up to date
    echo "Updating trunk branch..."
    if git remote get-url origin >/dev/null 2>&1; then
        if ! git pull origin "$trunk_branch" 2>/dev/null; then
            echo "Warning: Failed to pull latest changes from remote, but continuing..."
        fi
    fi
    
    # Create and switch to hotfix branch
    echo "Creating hotfix branch: $hotfix_name"
    if ! git checkout -b "$hotfix_name" 2>/dev/null; then
        echo "Error: Failed to create hotfix branch '$hotfix_name'" >&2
        exit 1
    fi
    
    # Set working branch to hotfix branch
    echo "Setting working branch to hotfix branch..."
    if ! git @ branch "$hotfix_name" >/dev/null 2>&1; then
        echo "Warning: Failed to set working branch, but hotfix branch is ready"
    fi
    
    echo ""
    echo "✅ Hotfix branch '$hotfix_name' created successfully!"
    echo ""
    echo "Current status:"
    echo "  Branch: $hotfix_name"
    echo "  Base: $trunk_branch"
    echo "  Working branch: $(git config at.branch 2>/dev/null || echo 'not set')"
    echo ""
    echo "Next steps:"
    echo "  1. Make your urgent fixes"
    echo "  2. git @ save 'Fix description'"
    echo "  3. git @ pr 'Hotfix: description'"
    echo "  4. After merge, consider: git @ release -p (patch release)"
    echo ""
    
    exit 0
}

# Validate branch name
validate_branch_name() {
    local name="$1"
    
    # Check if name is empty
    if [ -z "$name" ]; then
        return 1
    fi
    
    # Check for dangerous characters
    if echo "$name" | grep -q '[;&|`$(){}]'; then
        return 1
    fi
    
    # Check for path traversal attempts
    if echo "$name" | grep -q '\.\./'; then
        return 1
    fi
    
    # Check for valid characters (alphanumeric, hyphens, underscores, slashes)
    if ! echo "$name" | grep -qE '^[a-zA-Z0-9._/-]+$'; then
        return 1
    fi
    
    # Check for reserved names
    case "$name" in
        "HEAD"|"head"|"master"|"main"|"develop"|"development")
            return 1
            ;;
    esac
    
    return 0
}

# Run the command
cmd_hotfix "$@" 