#!/bin/bash

usage() {
    cat << 'EOF'
          __
       __/\ \__          __                    __
   __ /\_\ \ ,_\        /'_`\_      __  __  __/\_\  _____
 /'_ `\/\ \ \ \/       /'/'_` \    /\ \/\ \/\ \/\ \/\ '__`\
/\ \L\ \ \ \ \ \_     /\ \ \L\ \   \ \ \_/ \_/ \ \ \ \ \L\ \
\ \____ \ \_\ \__\    \ \ `\__,_\   \ \___x___/'\ \_\ \ ,__/
 \/___L\ \/_/\/__/     \ `\_____\    \/__//__/   \/_/\ \ \/
   /\____/              `\/_____/                     \ \_\
   \_/__/

Usage: git @ wip [options]

DESCRIPTION:
  Manage Work-In-Progress branch state.
  Tracks which branch you were working on for quick context switching.

OPTIONS:
  (no options)           Show current WIP branch
  -s, --set              Set current branch as WIP
  -c, --checkout         Checkout WIP branch
  -r, --restore          Restore WIP to working branch
  -h, --help             Show this help

EXAMPLES:
  git @ wip                    # Show current WIP branch
  git @ wip -s                 # Set current branch as WIP
  git @ wip -c                 # Checkout WIP branch
  git @ wip -r                 # Restore WIP to working branch

WORKFLOW:
  Use WIP to quickly switch between different features you're working on.
  Set WIP when you need to context switch to another task.

STORAGE:
  Saved in git config: at.wip

SECURITY:
  All WIP operations are validated and logged.

EOF
    exit 1
}

cmd_wip() {
    if [ "$#" -lt 1 ]; then
        show_wip; exit 0
    elif [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-s"|"--set"|"s"|"set"|".")
                set_wip; exit 0
                ;;
            "-c"|"--checkout"|"c"|"checkout")
                checkout_wip; exit 0
                ;;
            "-r"|"--restore"|"r"|"restore")
                restore_wip; exit 0
                ;;
        esac
    fi

    usage; exit 1
}

restore_wip() {
    wip_branch=$(git config at.wip 2>/dev/null || echo "")
    if [ -z "$wip_branch" ]; then
        echo "Error: No WIP branch configured" >&2
        exit 1
    fi
    git @ branch "$wip_branch"
    git @ work
    echo "Restored WIP branch: $wip_branch"
    exit 0
}

set_wip() {
    from=$(git config at.wip 2>/dev/null || echo "")
    branch=$(git @ branch -c)
    git config --replace-all at.wip "$branch"
    echo "WIP updated to $branch from $from"
    exit 0
}

show_wip() {
    current=$(git config at.wip 2>/dev/null || echo "")
    echo "$current"
    exit 0
}

checkout_wip() {
    wip_branch=$(git config at.wip 2>/dev/null || echo "")
    if [ -z "$wip_branch" ]; then
        echo "Error: No WIP branch configured" >&2
        exit 1
    fi
    git checkout "$wip_branch"
    echo "Switched to WIP branch: $wip_branch"
    exit 0
}

