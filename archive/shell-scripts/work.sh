#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ work <type> [<description>] [options]

DESCRIPTION:
  Create work branches following Conventional Commits specification.
  Supports all standard commit types for organized development workflow.

WORK TYPES (Conventional Commits):
  hotfix    : Urgent fixes for production (PATCH version)
  feature   : New features (MINOR version)
  bugfix    : Bug fixes (PATCH version)
  release   : Release preparation (PATCH version)
  chore     : Maintenance tasks
  docs      : Documentation changes
  style     : Code style changes
  refactor  : Code refactoring
  perf      : Performance improvements
  test      : Test additions/changes
  ci        : CI/CD changes
  build     : Build system changes
  revert    : Revert commits

OPTIONS:
  -n, --name <name>     Specify full branch name
  -h, --help           Show this help message

EXAMPLES:
  git @ work hotfix "fix-login-bug"           # Creates hotfix-fix-login-bug
  git @ work feature "add-user-auth"          # Creates feature-add-user-auth
  git @ work bugfix "fix-crash-on-startup"    # Creates bugfix-fix-crash-on-startup
  git @ work docs "update-api-documentation"  # Creates docs-update-api-documentation
  git @ work chore "update-dependencies"      # Creates chore-update-dependencies
  git @ work feature "Incorrect Branch Name"  # Creates feature-incorrect-branch-name
  git @ work hotfix "Fix Login Bug!"          # Creates hotfix-fix-login-bug

BRANCH NAMING:
  Format: <type>-<description>
  Automatic formatting: Descriptions are automatically converted to kebab-case
  Examples:
    hotfix-fix-login-bug
    feature-add-user-auth
    bugfix-fix-crash-on-startup
    docs-update-api-documentation
    chore-update-dependencies
    feature-incorrect-branch-name (from "Incorrect Branch Name")
    hotfix-fix-login-bug (from "Fix Login Bug!")

CONVENTIONAL COMMITS INTEGRATION:
  - Branch types follow Conventional Commits specification
  - Commit messages will include [TYPE] prefix
  - Supports semantic versioning correlation
  - Integrates with git @ branch --<type> listing

WORKFLOW:
  1. Creates branch from current branch or trunk
  2. Switches to new work branch
  3. Sets working branch to new branch
  4. Provides next steps guidance

EOF
    exit 1
}

# Format branch name to kebab-case
format_branch_name() {
    local input="$1"
    
    # Convert to lowercase
    input=$(echo "$input" | tr '[:upper:]' '[:lower:]')
    
    # Replace spaces, underscores, and other separators with hyphens
    input=$(echo "$input" | sed 's/[[:space:]]/-/g')
    input=$(echo "$input" | sed 's/_/-/g')
    input=$(echo "$input" | sed 's/\./-/g')
    input=$(echo "$input" | sed 's/\\/-/g')
    
    # Remove any non-alphanumeric characters except hyphens
    input=$(echo "$input" | sed 's/[^a-z0-9-]//g')
    
    # Remove multiple consecutive hyphens
    input=$(echo "$input" | sed 's/--\+/-/g')
    
    # Remove leading and trailing hyphens
    input=$(echo "$input" | sed 's/^-\+//; s/-\+$//')
    
    # If empty after formatting, use "update"
    if [ -z "$input" ]; then
        input="update"
    fi
    
    echo "$input"
}

cmd_work() {
    local work_type=""
    local description=""
    local full_name=""
    
    # Parse arguments
    while [ "$#" -gt 0 ]; do
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-n"|"--name")
                if [ -n "$2" ] && [[ "$2" != -* ]]; then
                    full_name="$2"
                    shift 2
                else
                    echo "Error: --name requires a value" >&2
                    usage; exit 1
                fi
                ;;
            *)
                # First argument is work type, second is description
                if [ -z "$work_type" ]; then
                    work_type="$1"
                elif [ -z "$description" ]; then
                    description="$1"
                else
                    echo "Error: Too many arguments" >&2
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
    
    # If full name provided, use it directly
    if [ -n "$full_name" ]; then
        create_work_branch "$full_name" "$current_branch"
        exit 0
    fi
    
    # Validate work type
    if [ -z "$work_type" ]; then
        echo "Error: Work type is required" >&2
        echo "Available types: hotfix, feature, bugfix, release, chore, docs, style, refactor, perf, test, ci, build, revert" >&2
        usage; exit 1
    fi
    
    # Validate work type against allowed types
    local allowed_types=("hotfix" "feature" "bugfix" "release" "chore" "docs" "style" "refactor" "perf" "test" "ci" "build" "revert")
    local valid_type=false
    
    for type in "${allowed_types[@]}"; do
        if [ "$work_type" = "$type" ]; then
            valid_type=true
            break
        fi
    done
    
    if [ "$valid_type" = false ]; then
        echo "Error: Invalid work type '$work_type'" >&2
        echo "Available types: ${allowed_types[*]}" >&2
        usage; exit 1
    fi
    
    # Prompt for description if not provided
    if [ -z "$description" ]; then
        echo "Creating $work_type branch"
        echo "Please enter a description for the $work_type:"
        read -p "$work_type description: " description
        
        if [ -z "$description" ]; then
            echo "Error: Description cannot be empty" >&2
            exit 1
        fi
    fi
    
    # Format description into kebab-case
    local formatted_description
    formatted_description=$(format_branch_name "$description")
    
    # Create branch name
    local branch_name="${work_type}-${formatted_description}"
    
    # Special handling for hotfix - should come from trunk
    if [ "$work_type" = "hotfix" ]; then
        local trunk_branch
        trunk_branch=$(git config at.trunk 2>/dev/null || echo "main")
        
        # Validate trunk branch exists
        if ! git rev-parse --verify "$trunk_branch" >/dev/null 2>&1; then
            echo "Error: Trunk branch '$trunk_branch' does not exist" >&2
            echo "Please ensure the trunk branch exists or configure it with: git @ _trunk <branch>" >&2
            exit 1
        fi
        
        # Switch to trunk branch first
        echo "Switching to trunk branch: $trunk_branch"
        if ! git checkout "$trunk_branch" 2>/dev/null; then
            echo "Error: Failed to switch to trunk branch '$trunk_branch'" >&2
            exit 1
        fi
        
        # Update trunk branch
        echo "Updating trunk branch..."
        if git remote get-url origin >/dev/null 2>&1; then
            if ! git pull origin "$trunk_branch" 2>/dev/null; then
                echo "Warning: Failed to pull latest changes from remote, but continuing..."
            fi
        fi
        
        current_branch="$trunk_branch"
    fi
    
    # Create the work branch
    create_work_branch "$branch_name" "$current_branch"
    
    exit 0
}

create_work_branch() {
    local branch_name="$1"
    local base_branch="$2"
    local work_type
    work_type=$(echo "$branch_name" | cut -d'-' -f1)
    
    # Validate branch name
    if ! validate_branch_name "$branch_name"; then
        echo "Error: Invalid branch name '$branch_name'" >&2
        echo "Branch names must contain only alphanumeric characters, hyphens, underscores, and slashes" >&2
        exit 1
    fi
    
    # Check if branch already exists
    if git rev-parse --verify "$branch_name" >/dev/null 2>&1; then
        echo "Error: Branch '$branch_name' already exists" >&2
        exit 1
    fi
    
    # Check for uncommitted changes
    if ! git diff --quiet || ! git diff --cached --quiet; then
        echo "Warning: You have uncommitted changes"
        echo "These will be saved to WIP before creating work branch"
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
    
    # Create and switch to work branch
    echo "Creating $work_type branch: $branch_name"
    if ! git checkout -b "$branch_name" 2>/dev/null; then
        echo "Error: Failed to create $work_type branch '$branch_name'" >&2
        exit 1
    fi
    
    # Set working branch to new branch
    echo "Setting working branch to $work_type branch..."
    if ! git @ branch "$branch_name" >/dev/null 2>&1; then
        echo "Warning: Failed to set working branch, but $work_type branch is ready"
    fi
    
    echo ""
    echo "✅ $work_type branch '$branch_name' created successfully!"
    echo ""
    echo "Current status:"
    echo "  Branch: $branch_name"
    echo "  Base: $base_branch"
    echo "  Working branch: $(git config at.branch 2>/dev/null || echo 'not set')"
    echo ""
    # Format work type for display
    local display_type
    case "$work_type" in
        "hotfix") display_type="HOTFIX" ;;
        "feature") display_type="FEATURE" ;;
        "bugfix") display_type="BUGFIX" ;;
        "release") display_type="RELEASE" ;;
        "chore") display_type="CHORE" ;;
        "docs") display_type="DOCS" ;;
        "style") display_type="STYLE" ;;
        "refactor") display_type="REFACTOR" ;;
        "perf") display_type="PERF" ;;
        "test") display_type="TEST" ;;
        "ci") display_type="CI" ;;
        "build") display_type="BUILD" ;;
        "revert") display_type="REVERT" ;;
        *) display_type="$(echo "$work_type" | tr '[:lower:]' '[:upper:]')" ;;
    esac
    
    local title_type
    case "$work_type" in
        "hotfix") title_type="Hotfix" ;;
        "feature") title_type="Feature" ;;
        "bugfix") title_type="Bugfix" ;;
        "release") title_type="Release" ;;
        "chore") title_type="Chore" ;;
        "docs") title_type="Docs" ;;
        "style") title_type="Style" ;;
        "refactor") title_type="Refactor" ;;
        "perf") title_type="Perf" ;;
        "test") title_type="Test" ;;
        "ci") title_type="CI" ;;
        "build") title_type="Build" ;;
        "revert") title_type="Revert" ;;
        *) title_type="$(echo "$work_type" | sed 's/^./\U&/')" ;;
    esac
    
    echo "Next steps:"
    echo "  1. Make your changes"
    echo "  2. git @ save '[$display_type] Description of changes'"
    echo "  3. git @ pr '$title_type: Description of changes'"
    
    # Special guidance for different types
    case "$work_type" in
        "hotfix")
            echo "  4. After merge, consider: git @ release -p (patch release)"
            ;;
        "feature")
            echo "  4. After merge, consider: git @ release -m (minor release)"
            ;;
        "release")
            echo "  4. After merge, consider: git @ release -M (major release)"
            ;;
    esac
    echo ""
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

# Only run the command if this script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    cmd_work "$@"
fi
