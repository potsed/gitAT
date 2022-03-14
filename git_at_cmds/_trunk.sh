usage() {
    echo
    exit 1
}

cmd__trunk() {
    if [ "$#" -lt 1 ]; then
        show_trunk; exit 0
    elif [ "$#" -eq 1 ]; then
        if [ $1 == "help" ]; then
            usage; exit 0
        fi

        set_trunk $1; exit 0
    fi

    usage; exit 1
}

set_trunk() {
    from=`git @ _trunk`
    `git config --replace-all at.trunk $1`
    echo 'Base branch updated to: '`git @ _trunk`" from $from"
    exit 1
}

show_trunk() {
    current=`git config at.trunk`
    if [ $current == "" ]; then
        current=`git branch -rl "*/HEAD" | rev | cut -d/ -f1 | rev`
        set_trunk $current
    fi
    echo $current
    exit 1
}