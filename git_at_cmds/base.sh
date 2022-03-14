usage() {
    echo
    exit 1
}

cmd_base() {
    if [ "$#" -lt 1 ]; then
        show_base; exit 0
    elif [ "$#" -eq 1 ]; then
        if [ $1 == "help" ]; then
            usage; exit 0
        fi

        set_base $1; exit 0
    fi

    usage; exit 1
}

set_base() {
    from=`git @ base`
    `git config --replace-all at.base $1`
    echo 'Base branch updated to: '`git @ base`" from $from"
    exit 1
}

show_base() {
    current=`git config at.base`
    if [ $current == "" ]; then
        current=`git branch -rl "*/HEAD" | rev | cut -d/ -f1 | rev`
        set_base $current
    fi
    echo $current
    exit 1
}