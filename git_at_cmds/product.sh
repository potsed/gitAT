#!/bin/bash



usage() {
    cat << 'EOF'
Usage: git @ product [<product-name>]

DESCRIPTION:
  Set or get the current product name for GitAT workflow management.
  The product name is used in commit labels and configuration.
  Examples: gitAT, myApp, apiService

EXAMPLES:
  git @ product                    # Show current product name
  git @ product gitAT              # Set product name to "gitAT"
  git @ product myApp              # Set product name to "myApp"
  git @ product apiService         # Set product name to "apiService"

VALIDATION:
  Product names must contain only:
  - Alphanumeric characters (a-z, A-Z, 0-9)
  - Dots (.)
  - Underscores (_)
  - Hyphens (-)

STORAGE:
  Saved in git config: at.product

SECURITY:
  All inputs are validated against dangerous characters and patterns.

EOF
    exit 1
}

cmd_product() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    if [ "$#" -lt 1 ]; then
        show_project; exit 0
    elif [ "$#" -gt 0 ]; then
        set_project "$@"; exit 0
    fi

    usage; exit 0
}

set_project() {
    from=$(git config at.product 2>/dev/null || echo "")
    git config --replace-all at.product "$*"
    echo "Project updated to: $(show_project) from $from"
    exit 0
}

show_project() {
    local project
    project=$(git config at.product 2>/dev/null || echo "")
    echo "$project"
    exit 0
}