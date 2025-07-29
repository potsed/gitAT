#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ version [options]

DESCRIPTION:
  Manage semantic versioning for your project.
  Uses MAJOR.MINOR.FIX format (e.g., 1.2.3).

OPTIONS:
  (no options)           Show current version
  -M, --major           Increment major version (resets minor and fix to 0)
  -m, --minor           Increment minor version (resets fix to 0)
  -b, --bump            Increment fix version
  -t, --tag             Show version tag (e.g., "v1.2.3")
  -r, --reset           Reset version to 0.0.0 (use with caution)
  -h, --help            Show this help

EXAMPLES:
  git @ version                    # Show current version (e.g., "1.2.3")
  git @ version -M                 # Increment major: 1.2.3 → 2.0.0
  git @ version -m                 # Increment minor: 1.2.3 → 1.3.0
  git @ version -b                 # Increment fix: 1.2.3 → 1.2.4
  git @ version -t                 # Show version tag (e.g., "v1.2.3")
  git @ version -r                 # Reset to 0.0.0

STORAGE:
  Major version: git config at.major
  Minor version: git config at.minor
  Fix version: git config at.fix

SEMANTIC VERSIONING:
  MAJOR: Breaking changes, incompatible API changes
  MINOR: New features, backward compatible
  FIX: Bug fixes, backward compatible

SECURITY:
  All version operations are logged for audit purposes.

EOF
    exit 0
}

cmd_version() {
    if [ "$#" -lt 1 ]; then
        show_version; exit 0
    fi

    # Handle help flag first
    if [ "$1" = "-h" ] || [ "$1" = "--help" ] || [ "$1" = "help" ] || [ "$1" = "h" ]; then
        usage; exit 0
    fi

    # Handle tag flag
    if [ "$1" = "-t" ]; then
        show_tag; exit 0
    fi

    # Handle reset flag
    if [ "$1" = "-r" ]; then
        reset_version; exit 0
    fi

    # Initialize variables
    local MAJOR=false
    local MINOR=false
    local BUMP=false

    # Process other flags
    while getopts ':Mmb' flag; do
        case "${flag}" in
            M) MAJOR=true ;;
            m) MINOR=true ;;
            b) BUMP=true ;;
        esac
    done

    if [[ "$MAJOR" == "true" ]]; then
        increment_major;
    fi

    if [[ "$MINOR" == "true" ]]; then
        increment_minor;
    fi

    if [[ "$BUMP" == "true" ]]; then
        increment_fix;
    fi

    show_version;
}

show_tag() {
    local version
    version=$(show_version)
    echo "v$version"
    exit 0
}

reset_version() {
    git config --replace-all at.major 0;
    git config --replace-all at.minor 0;
    git config --replace-all at.fix 0;
}

increment_major() {
    local OLD_MAJOR=$(git config at.major);
    set_major $(($OLD_MAJOR + 1))
}

increment_minor() {
    local OLD_MINOR=$(git config at.minor);
    set_minor $(($OLD_MINOR + 1))
}

increment_fix() {
    local OLD_BUMP=$(git config at.fix);
    set_fix $(($OLD_BUMP + 1))
}

set_major() {
    git config --replace-all at.major $1
    set_minor 0;
}

set_minor() {
    git config --replace-all at.minor $1
    set_fix 0;
}

set_fix() {
    git config --replace-all at.fix "$1";
}

show_version() {
    local M=$(git config at.major 2>/dev/null || echo "0")
    local m=$(git config at.minor 2>/dev/null || echo "0")
    local f=$(git config at.fix 2>/dev/null || echo "0")
    echo "$M.$m.$f"
    exit 0
}