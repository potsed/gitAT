usage() {
    echo
    exit 1
}

cmd_project() {
    if [ "$#" -lt 1 ]; then
        show_project; exit 0
    elif [ "$#" -eq 1 ]; then
        if [ $1 == "help" ]; then
            usage; exit 0
        fi

        set_project $1; exit 0
    fi

    usage; exit 1
}

set_project() {
    from=`git @ project`
    `git config --replace-all at.project $1`
    echo 'project updated to: '`git @ project`" from $from"
    exit 1
}

show_project() {
    echo `git config at.project`
    exit 1
}