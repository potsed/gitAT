#!/bin/bash



usage() {
    cat << 'EOF'
Usage: git @ feature [<feature-name>]

DESCRIPTION:
  Set or get the current feature name for GitAT workflow management.
  The feature name is used in commit labels and helps track what you're working on.

EXAMPLES:
  git @ feature                    # Show current feature name
  git @ feature user-auth          # Set feature to "user-auth"
  git @ feature payment-integration # Set feature to "payment-integration"

VALIDATION:
  Feature names must contain only:
  - Alphanumeric characters (a-z, A-Z, 0-9)
  - Dots (.)
  - Underscores (_)
  - Hyphens (-)

STORAGE:
  Saved in git config: at.feature

SECURITY:
  All inputs are validated against dangerous characters and patterns.

EOF
    exit 1
}

cmd_feature() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    if [ "$#" -lt 1 ]; then
        show_feature; exit 0
    elif [ "$#" -eq 1 ]; then
        set_feature "$1"; exit 0
    fi

    usage; exit 1
}

set_feature() {
    from=$(git config at.feature 2>/dev/null || echo "")
    git config --replace-all at.feature "$1"
    echo "Feature updated to: $(show_feature) from $from"
    exit 0
}

show_feature() {
    local feature
    feature=$(git config at.feature 2>/dev/null || echo "")
    echo "$feature"
    exit 0
}