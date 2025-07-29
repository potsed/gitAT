#!/bin/bash

usage() {
    echo "Usage: git @ product [<project-name>]"
    echo "  Set or get the current project name"
    echo "  Project name must be alphanumeric with dots, underscores, and hyphens only"
    exit 1
}

cmd_product() {
    if [ "$#" -lt 1 ]; then
        show_project; exit 0
    elif [ "$#" -gt 0 ]; then
        if [ "$1" == "help" ] || [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
            usage; exit 0
        fi

        # Validate input before setting
        if ! validate_input "$*"; then
            echo "Error: Invalid project name. Use only alphanumeric characters, dots, underscores, and hyphens." >&2
            exit 1
        fi

        set_project "$@"; exit 0
    fi

    usage; exit 0
}

set_project() {
    # Use secure configuration function
    if secure_config "at.project" "$*"; then
        echo "Project updated to: $(show_project)"
        exit 0
    else
        echo "Error: Failed to update project configuration" >&2
        exit 1
    fi
}

show_project() {
    local project
    project=$(git config at.project 2>/dev/null || echo "")
    echo "$project"
    exit 0
}