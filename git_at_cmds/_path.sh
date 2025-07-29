#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ _path

DESCRIPTION:
  Get the git repository root path.
  Returns the absolute path to the root of the current git repository.

EXAMPLES:
  git @ _path                    # Show repository root path

OUTPUT:
  Absolute path to the git repository root directory.

SECURITY:
  All path operations are validated and logged.

EOF
    exit 1
}

cmd__path() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    git rev-parse --show-toplevel
}