usage() {
    echo
    exit 1
}

cmd_label() {
    if [ "$#" -lt 1 ]; then
        show_label; exit 0
    elif [ "$#" -eq 1 ]; then
        if [ "$1" == "help" ]; then
            usage; exit 0
        fi

        set_label "$@"; exit 0
    fi

    usage;
}

set_label() {
    `git config --replace-all at.label "$@ "`

    echo 'Label updated to: '`git @ label`
    exit 1
}

show_label() {
    userLabel=`git config at.label`



    echo $userLabel"(S"`git @ sprint`".F"`git @ feature`".T"`git @ task`")"
    exit 1
}