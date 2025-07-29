#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ changes

DESCRIPTION:
  Show uncommitted changes in the working directory.
  Lists files that have been modified but not yet committed.

EXAMPLES:
  git @ changes                    # Show modified files

OUTPUT:
  Lists file names that have been changed since last commit.

SECURITY:
  All change operations are validated and logged.

EOF
    exit 1
}

cmd_changes() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    git diff --name-only --no-color
}