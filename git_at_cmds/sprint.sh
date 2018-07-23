usage() {
    echo
    exit 1
}

cmd_sprint() {
    if [ "$#" -lt 1 ]; then
        show_sprint; exit 0
    elif [ "$#" -eq 1 ]; then
        if [ $1 == "help" ]; then
            usage; exit 0
        fi

        set_sprint $1; exit 0
    fi

    usage; exit 1
}

set_sprint() {
    from=`git @ sprint`
    `git config --replace-all at.sprint $1`
    echo 'Sprint updated to: '`git @ sprint`" from $from"
    exit 1
}

show_sprint() {
    echo `git config at.sprint`
    exit 1
}