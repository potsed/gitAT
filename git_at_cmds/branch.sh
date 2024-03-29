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
            "-n"|"--new"|"n"|"new")
                new_working_branch; exit 0
                ;;
            "-c"|"--current"|"c"|"current")
                current_branch; exit 0
                ;;
            "-s"|"--set"|"s"|"set"|".")
                set_branch `git branch --show-current`
                ;;
        esac
        set_branch $1; exit 0
    fi

    usage; exit 1
}

new_working_branch() {
    echo "Hello"; exit 0
}

current_branch() {
    echo `git rev-parse --abbrev-ref HEAD`; exit 0
}

set_branch() {
    from=`git @ branch`
    git config --replace-all at.branch $1
    echo 'Branch updated to '`git @ branch`" from $from"
    exit 1
}

show_branch() {
    echo `git config at.branch`
    exit 1
}