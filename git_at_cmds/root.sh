#!/bin/bash

cmd_root() {
    local current="$(git @ branch -c)"
    local base="$(git @ _trunk)"

    if [ "$current" != "${base}" ]; then
        git stash push --include-untracked -m "switched-to-${base}"
        git checkout ${base};
        git pull --rebase
    fi
}