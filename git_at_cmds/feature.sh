#!/bin/bash

usage() {
    echo "Usage: git @ feature [<feature-name>]"
    echo "  Set or get the current feature name"
    echo "  Feature name must be alphanumeric with dots, underscores, and hyphens only"
    exit 1
}

cmd_feature() {
    if [ "$#" -lt 1 ]; then
        show_feature; exit 0
    elif [ "$#" -eq 1 ]; then
        if [ "$1" == "help" ] || [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
            usage; exit 0
        fi

        # Validate input before setting
        if ! validate_input "$1"; then
            echo "Error: Invalid feature name. Use only alphanumeric characters, dots, underscores, and hyphens." >&2
            exit 1
        fi

        set_feature "$1"; exit 0
    fi

    usage; exit 1
}

set_feature() {
    # Use secure configuration function
    if secure_config "at.feature" "$1"; then
        echo "Feature updated to: $(show_feature)"
        exit 0
    else
        echo "Error: Failed to update feature configuration" >&2
        exit 1
    fi
}

show_feature() {
    local feature
    feature=$(git config at.feature 2>/dev/null || echo "")
    echo "$feature"
    exit 0
}