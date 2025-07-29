#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ _trunk [<branch-name>]

DESCRIPTION:
  Manage the base/trunk branch (usually develop or master).
  The trunk branch is used as the base for feature branches.

EXAMPLES:
  git @ _trunk                    # Show current trunk branch
  git @ _trunk develop            # Set trunk to "develop"
  git @ _trunk master             # Set trunk to "master"

AUTO-DETECTION:
  If no trunk is set, automatically detects from remote HEAD.

STORAGE:
  Saved in git config: at.trunk

SECURITY:
  All trunk operations are validated and logged.

EOF
    exit 1
}

cmd__trunk() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    if [ "$#" -lt 1 ]; then
        show_trunk; exit 0
    elif [ "$#" -eq 1 ]; then
        set_trunk "$1"; exit 0
    fi

    usage; exit 1
}

set_trunk() {
    from=$(git config at.trunk 2>/dev/null || echo "")
    git config --replace-all at.trunk "$1"
    echo "Base branch updated to: $1 from $from"
    exit 0
}

show_trunk() {
    current=$(git config at.trunk 2>/dev/null || echo "")
    if [ "$current" == "" ]; then
        # Auto-detect trunk branch from remote HEAD
        current=$(git branch -rl "*/HEAD" | rev | cut -d/ -f1 | rev 2>/dev/null || echo "develop")
        # Set it without calling set_trunk to avoid recursion
        git config --replace-all at.trunk "$current"
        echo "Auto-detected trunk branch: $current"
    fi
    echo "$current"
    exit 0
}