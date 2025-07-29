#!/bin/bash

cmd_ignore() {
     if [ "$#" -lt 1 ]; then
        usage; exit 0
    elif [ "$#" -eq 1 ]; then
        local REPOPATH=$(git @ _path)"/.gitignore";
        if grep -Fxq "$1" "$REPOPATH"; then
            echo "String $1 exists in $REPOPATH"
        else
            echo "$1" >> "$REPOPATH"
            echo "String $1 appended to $REPOPATH"
        fi
    fi
    exit 0;
}

usage() {
    echo 'Usage: git @ ignore ignore/pathto/file';
    exit 1;
}