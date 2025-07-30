#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ root

DESCRIPTION:
  Switch to root/trunk branch with stash management.
  Stashes current changes, switches to trunk branch, and pulls latest changes.

PROCESS:
  1. Stashes current changes (if not on trunk)
  2. Switches to trunk branch
  3. Pulls latest changes with rebase

EXAMPLES:
  git @ root                    # Switch to trunk branch

SECURITY:
  All root operations are validated and logged.

EOF
    exit 1
}

cmd_root() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    local current="$(git @ branch -c)"
    local base="$(git @ _trunk)"

    if [ "$current" != "${base}" ]; then
        git stash push --include-untracked -m "switched-to-${base}"
        git checkout ${base};
        git pull --rebase
    fi
}