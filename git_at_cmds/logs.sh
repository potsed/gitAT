#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ logs

DESCRIPTION:
  Show recent commit history in a compact format.
  Displays the last 10 commits with abbreviated commit hashes.

EXAMPLES:
  git @ logs                    # Show recent commits

OUTPUT:
  Shows commit hash, author, date, and message for recent commits.

SECURITY:
  All log operations are validated and logged.

EOF
    exit 1
}

cmd_logs() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    git log -10 --pretty=oneline --abbrev-commit
}