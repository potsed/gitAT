usage() {
    echo
    exit 1
}

cmd_branch() {
    if [ "$#" -lt 1 ]; then
        show_branch; exit 0
    elif [ "$#" -eq 1 ]; then
        case $1 in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-c"|"--current"|"c"|"current")
                current_branch; exit 0
                ;;
            "-s"|"--sweep"|"s"|"sweep")
                delete_merged_branches_locally; exit 0
                ;;
            ".")
                set_branch `git @ branch --current`
                ;;
        esac
        set_branch $1; exit 0
    fi

    usage; exit 1
}

current_branch() {
    echo `git rev-parse --abbrev-ref HEAD`; exit 0
}

set_branch() {
    from=`git @ branch`
    `git config --replace-all at.branch $1`
    echo 'Branch updated to '`git @ branch`" from $from"
    exit 1
}

show_branch() {
    echo `git config at.branch`
    exit 1
}

delete_merged_branches_locally() {
    git branch --merged | grep -v '\*\|master\|dev|develop' | xargs -n 1 git branch -d
}