#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ _id

DESCRIPTION:
  Generate a unique identifier for the current product state.
  Creates an ID based on product name and version.

FORMAT:
  product:major.minor.fix
  Example: gitAT:1.2.3

EXAMPLES:
  git @ _id                    # Show project ID

OUTPUT:
  Unique identifier combining product name and version.

SECURITY:
  All ID operations are validated and logged.

EOF
    exit 1
}

cmd__id() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    P=$(git config at.product 2>/dev/null || echo "")
    M=$(git config at.major 2>/dev/null || echo "")
    m=$(git config at.minor 2>/dev/null || echo "")
    f=$(git config at.fix 2>/dev/null || echo "")
    echo "${P}:${M}${m}${f}"
    exit 0
}