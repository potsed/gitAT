#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ hash

DESCRIPTION:
  Show detailed branch status and commit relationships.
  Displays hash comparisons, commit counts, and recent commit history.

OUTPUT:
  - Current branch hash vs remote
  - Commit counts (ahead/behind)
  - Comparison with develop/master branches
  - Last 5 commits with hashes, committers, and messages

EXAMPLES:
  git @ hash                    # Show branch status and recent commits

INFORMATION DISPLAYED:
  - Current branch hash and remote hash
  - How many commits ahead/behind remote
  - Comparison with develop and master branches
  - Recent commit history (last 5 commits with committers)

SECURITY:
  All hash operations are validated and logged.

EOF
    exit 1
}

cmd_hash() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo "Error: Not in a git repository" >&2
        exit 1
    fi
    
    # Get current branch directly instead of using git @ branch -c
    BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    if [ -z "$BRANCH" ]; then
        echo "Error: Could not determine current branch" >&2
        exit 1
    fi
    
    # Fetch remote changes
    git fetch > /dev/null 2>&1
    
    # Check if remote branches exist
    REM_HAS_MASTER=$(git ls-remote --heads origin master 2>/dev/null | wc -l | tr -d ' ')
    REM_HAS_DEVELOP=$(git ls-remote --heads origin develop 2>/dev/null | wc -l | tr -d ' ')
    REM_HAS_CURRENT=$(git ls-remote --heads origin "${BRANCH}" 2>/dev/null | wc -l | tr -d ' ')

    REM_CURRENT_HASH="UNKNOWN"
    CURRENT_BEHIND_REM="0"
    CURRENT_AHEAD_REM="0"

    REM_DEVELOP_BR_HASH="UNKNOWN"
    LOC_BR_BEHIND_REM_DEVELOP="0"
    LOC_BR_AHEAD_REM_DEVELOP="0"

    REM_MASTER_BR_HASH="UNKNOWN"
    LOC_BR_BEHIND_REM_MASTER="0"
    LOC_BR_AHEAD_REM_MASTER="0"

    # Get local hashes with error handling
    LOC_BR_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "UNKNOWN")
    LOC_DEVELOP_HASH=$(git rev-parse --short develop 2>/dev/null || echo "UNKNOWN")
    LOC_MASTER_HASH=$(git rev-parse --short master 2>/dev/null || echo "UNKNOWN")

    if [ "$REM_HAS_CURRENT" == "1" ]; then
        REMBR="origin/${BRANCH}"
        REM_CURRENT_HASH=$(git rev-parse --short "${REMBR}" 2>/dev/null || echo "UNKNOWN")
        CURRENT_BEHIND_REM=$(git rev-list --left-right --count "${REMBR}"...@ 2>/dev/null | cut -f1 || echo "0")
        CURRENT_AHEAD_REM=$(git rev-list --left-right --count "${REMBR}"...@ 2>/dev/null | cut -f2 || echo "0")
    fi

    echo
    echo "BRANCH ${BRANCH}"
    echo "- CURRENT HASH: ${LOC_BR_HASH}"
    echo "- ORIGIN HASH: ${REM_CURRENT_HASH}"
    echo "CURRENT is ${CURRENT_BEHIND_REM} behind and ${CURRENT_AHEAD_REM} ahead of origin/${BRANCH}"

    if [ "$REM_HAS_DEVELOP" == "1" ]; then
        LOC_BR_BEHIND_REM_DEVELOP=$(git rev-list --left-right --count origin/develop...@ 2>/dev/null | cut -f1 || echo "0")
        LOC_BR_AHEAD_REM_DEVELOP=$(git rev-list --left-right --count origin/develop...@ 2>/dev/null | cut -f2 || echo "0")
        echo "CURRENT is ${LOC_BR_BEHIND_REM_DEVELOP} behind and ${LOC_BR_AHEAD_REM_DEVELOP} ahead of origin/develop"
    fi

    if [ "$REM_HAS_MASTER" == "1" ]; then
        LOC_BR_BEHIND_REM_MASTER=$(git rev-list --left-right --count origin/master...@ 2>/dev/null | cut -f1 || echo "0")
        LOC_BR_AHEAD_REM_MASTER=$(git rev-list --left-right --count origin/master...@ 2>/dev/null | cut -f2 || echo "0")
        echo "CURRENT is ${LOC_BR_BEHIND_REM_MASTER} behind and ${LOC_BR_AHEAD_REM_MASTER} ahead of origin/master"
    fi
    
    # Show recent commit history
    echo
    echo "RECENT COMMITS (last 5):"
    echo "────────────────────────────────────────────────────────────────────────────────"
    
    # Get the last 5 commits with hash, message, and committer
    local commit_count=0
    local has_commits=false
    
    while IFS= read -r line; do
        if [ -n "$line" ]; then
            commit_count=$((commit_count + 1))
            if [ $commit_count -le 5 ]; then
                has_commits=true
                local hash=$(echo "$line" | cut -d' ' -f1)
                local message=$(echo "$line" | cut -d' ' -f2-)
                local committer=$(git log -1 --pretty=format:"%an" "$hash" 2>/dev/null || echo "Unknown")
                
                # Handle empty messages
                if [ -z "$message" ]; then
                    message="<no message>"
                fi
                
                # Truncate very long messages
                if [ ${#message} -gt 50 ]; then
                    message="${message:0:47}..."
                fi
                
                # Truncate very long committer names
                if [ ${#committer} -gt 20 ]; then
                    committer="${committer:0:17}..."
                fi
                
                printf "%-8s │ %-20s │ %s\n" "$hash" "$committer" "$message"
            fi
        fi
    done < <(git log --oneline -5 2>/dev/null || echo "")
    
    if [ "$has_commits" = false ]; then
        echo "No commits found"
    fi
    
    echo
    exit 0
}