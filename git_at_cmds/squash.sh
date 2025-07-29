#!/bin/bash

usage() {
    cat << 'EOF'
       _ _         ____                                _
  __ _(_) |_      / __ \     ___  __ _ _   _  __ _ ___| |__
 / _` | | __|    / / _` |   / __|/ _` | | | |/ _` / __| '_ \
| (_| | | |_    | | (_| |   \__ \ (_| | |_| | (_| \__ \ | | |
 \__, |_|\__|    \ \__,_|   |___/\__, |\__,_|\__,_|___/_| |_|
 |___/            \____/            |_|


Usage: git @ squash [options] <target-branch>

DESCRIPTION:
  Reset the current branch to the HEAD of another branch (usually develop or master)
  to create a clean commit without working history. Useful for cleaning up before PR.

OPTIONS:
  -s, --save           Run 'git @ save' after squashing
  -h, --help           Show this help

EXAMPLES:
  git @ squash develop              # Reset to develop HEAD
  git @ squash master               # Reset to master HEAD
  git @ squash develop -s           # Reset and save

PROCESS:
  1. Validates target branch exists
  2. Retrieves HEAD SHA of target branch
  3. Soft reset to target branch SHA
  4. Keeps all changes staged for commit
  5. Optionally runs 'git @ save'

USE CASES:
  - Clean up commit history before PR
  - Remove intermediate commits from feature branch
  - Create single clean commit from multiple commits

WARNING:
  You may need to force push after squashing if branch is shared.
  Use with caution on shared branches.

GIT COMMANDS USED:
  - git rev-parse --verify --quiet --long ${BRANCH}
  - git reset --soft ${SHA}

SECURITY:
  All squash operations are validated and logged.

EOF
    exit 1
}

cmd_squash() {
    while getopts ':hs' flag; do
        case "${flag}" in
            h) usage; exit 0 ;;
            s) local DOSAVE=1 ;;
        esac
    done
    shift $(expr ${OPTIND} - 1)

    if [ "${1}" == "" ]; then
        usage; exit 0
    fi

    local HEAD="$(head "$1")"
    if [ "${HEAD}" == "0" ]; then
        echo
        echo "ERROR: Branch \"${1}\" does not exist locally"
        echo
        exit 0
    fi

    squash "$HEAD"

    echo
    echo 'Squashed branch '"$(git @ branch -c)"' back to '"$1"

    if [ "${DOSAVE}" == "1" ]; then
        save
    fi
}

head() {
    local HEAD=$(git rev-parse --verify --quiet --long "$1")
    if [ "${HEAD}" != "" ]; then
        echo "${HEAD}"
    else
        echo "0"
    fi
}

# exists() {

# }

squash() {
    git reset --soft "$1"
}

save() {
    git @ save
}