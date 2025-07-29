#!/bin/bash

usage() {
    echo
    exit 1
}

cmd_issue() {
    if [ "$#" -lt 1 ]; then
        show_task; exit 0
    elif [ "$#" -eq 1 ]; then
        if [ "$1" == "help" ]; then
            usage; exit 0
        fi

        set_task "$1"; exit 0
    fi

    usage; exit 1
}

set_task() {
    from=`git @ issue`
    git config --replace-all at.task "$1"
    echo 'Task updated to: '`git @ issue`" from $from"
    exit 1
}

show_task() {
    echo `git config at.task`
    exit 1
}