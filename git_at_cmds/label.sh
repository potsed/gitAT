usage() {
    echo
    exit 1
}

cmd_label() {
    if [ "$#" -lt 1 ]; then
        show_label; exit 0
    elif [ "$#" -eq 1 ]; then
        if [ "$1" == "help" ]; then
            usage; return;
        fi

        set_label "$@"; return;
    fi

    usage;
}

set_label() {
    `git config --replace-all at.label "$@ "`

    echo 'Label updated to: '`git @ label`
    return;
}

show_label() {
    userLabel=`git config at.label`



    echo $userLabel"(S"`git @ sprint`".F"`git @ feature`".T"`git @ task`")"
    return;
}