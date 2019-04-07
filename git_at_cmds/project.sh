usage() {
    echo
    exit 1
}

cmd_project() {
    if [ "$#" -lt 1 ]; then
        show_project; exit 0
    elif [ "$#" -gt 0 ]; then
        if [ "$1" == "help" ]; then
            usage; exit 0
        fi

        set_project "$@"; exit 0
    fi

    usage; exit 0
}

set_project() {
    git config --replace-all at.project "$@"
    echo 'project updated to: '$(git @ project);
    exit 0
}

show_project() {
    echo $(git config at.project);
    exit 0
}