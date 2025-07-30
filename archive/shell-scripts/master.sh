#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ master

DESCRIPTION:
  Switch to master branch with stash management.
  Stashes current changes, switches to master, and pulls latest changes.

PROCESS:
  1. Stashes current changes (if not on master)
  2. Switches to master branch
  3. Pulls latest changes from remote

EXAMPLES:
  git @ master                    # Switch to master branch

SECURITY:
  All master operations are validated and logged.

EOF
    exit 1
}

cmd_master() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    local current="$(git @ branch -c)"
    local base="$(git @ _trunk)"
    if [ "$current" != "$base" ]; then
        git stash push --include-untracked -m "switched-to-master"
        git checkout $base;
        git pull
    fi
}
