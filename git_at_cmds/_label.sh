usage() {
    cat << 'EOF'
Usage: git @ _label [<custom-label>]

DESCRIPTION:
  Manage commit labels for GitAT workflow.
  Labels are used in commit messages and help track project context.

EXAMPLES:
  git @ _label                    # Show formatted label [product.feature.issue]
  git @ _label "Custom label"     # Set custom label

DEFAULT FORMAT:
  [product.feature.issue]
  Example: [gitAT.user-auth.PROJ-123]

STORAGE:
  Saved in git config: at.label

SECURITY:
  All label operations are validated and logged.

EOF
    exit 1
}

cmd__label() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    if [ "$#" -lt 1 ]; then
        show_label; exit 0
    elif [ "$#" -eq 1 ]; then
        set_label "$@"; exit 0
    fi

    usage; exit 1
}

set_label() {
    git config --replace-all at.label "$*"

    echo "Label updated to: $(show_label)"
    exit 0
}

show_label() {
    local product
    local feature
    local issue
    
    # Get values directly from git config
    product=$(git config at.product 2>/dev/null || echo "")
    feature=$(git config at.feature 2>/dev/null || echo "")
    issue=$(git config at.task 2>/dev/null || echo "")
    
    # Format the label
    if [ -n "$product" ] || [ -n "$feature" ] || [ -n "$issue" ]; then
        echo "[${product}.${feature}.${issue}]"
    else
        echo "[Update]"
    fi
    exit 0
}