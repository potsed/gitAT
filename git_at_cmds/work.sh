#!/bin/bash

usage() {
    cat << 'EOF'
          __                                                 __
       __/\ \__          __                                 /\ \
   __ /\_\ \ ,_\        /'_`\_      __  __  __    ___   _ __\ \ \/'\
 /'_ `\/\ \ \ \/       /'/'_` \    /\ \/\ \/\ \  / __`\/\`'__\ \ , <
/\ \L\ \ \ \ \ \_     /\ \ \L\ \   \ \ \_/ \_/ \/\ \L\ \ \ \/ \ \ \\`\
\ \____ \ \_\ \__\    \ \ `\__,_\   \ \___x___/'\ \____/\ \_\  \ \_\ \_\
 \/___L\ \/_/\/__/     \ `\_____\    \/__//__/   \/___/  \/_/   \/_/\/_/
   /\____/              `\/_____/
   \_/__/

Usage: git @ work

DESCRIPTION:
  Switch to your working branch with intelligent stash management.
  This is the primary command for starting work on a feature.

PROCESS:
  1. Stashes current changes (if any)
  2. Fetches latest from remote
  3. Updates base branch (develop/master)
  4. Creates working branch if needed
  5. Restores stashed changes

EXAMPLES:
  git @ work                    # Switch to working branch

WORKFLOW:
  Use this command when you want to start working on your feature.
  It ensures you're on the correct branch with latest changes.

SECURITY:
  All operations are validated and logged for audit purposes.

EOF
    exit 0
}

cmd_work() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi

    local current=$(git @ branch -c)
    local branch=$(git @ branch)
    local STASHKEY='autostash-work-branch'
    if [ "$current" == "$branch" ]; then

        echo "You're already in the working branch"
        echo
        exit 0
    fi

    echo 'Fetching branches'
    git fetch;

    echo 'Stashing Changes'
    git stash push -m $STASHKEY;

    local HAS_STASH=$(git stash list | grep 'test' | wc -l | tr -d ' ')

    if [ ! "$(git branch --list "$branch")" ]; then

        echo 'Switching to and updating base branch'
        git checkout "$(git @ _trunk)";
        git pull;

        echo 'Creating local branch from updated dev branch'
        git branch $branch;
    fi;

    echo 'Switching to local branch'
    git checkout $branch

    if [ "$HAS_STASH" == "1" ]; then
        echo 'Reapplying stashed changes'
        git stash pop;
    fi;

    exit 0
}
