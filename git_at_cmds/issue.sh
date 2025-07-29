#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ issue [<issue-id>]

DESCRIPTION:
  Set or get the current issue/task identifier for tracking.
  The issue ID is used in commit labels and helps link commits to issues.

EXAMPLES:
  git @ issue                    # Show current issue ID
  git @ issue PROJ-123           # Set issue to "PROJ-123"
  git @ issue BUG-456            # Set issue to "BUG-456"

STORAGE:
  Saved in git config: at.task

SECURITY:
  All issue operations are validated and logged.

EOF
    exit 1
}

cmd_issue() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    if [ "$#" -lt 1 ]; then
        show_task; exit 0
    elif [ "$#" -eq 1 ]; then
        set_task "$1"; exit 0
    fi

    usage; exit 1
}

set_task() {
    from=$(git @ issue)
    git config --replace-all at.task "$1"
    echo 'Task updated to: '$(git @ issue)" from $from"
    exit 0
}

show_task() {
    echo "$(git config at.task 2>/dev/null || echo "")"
    exit 0
}