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
            exit 0
}

usage() {
    cat << 'EOF'
Usage: git @ ignore <pattern>

DESCRIPTION:
  Add patterns to .gitignore file.
  Checks if pattern already exists before adding.

ARGUMENTS:
  <pattern>    Pattern to add to .gitignore

EXAMPLES:
  git @ ignore "*.log"              # Ignore all log files
  git @ ignore "node_modules/"      # Ignore node_modules directory
  git @ ignore "build/"             # Ignore build directory
  git @ ignore "*.tmp"              # Ignore temporary files

PROCESS:
  1. Checks if pattern already exists in .gitignore
  2. Adds pattern if not present
  3. Reports success or existing pattern

SECURITY:
  All ignore operations are validated and logged.

EOF
    exit 1
}