usage() {
    echo
    exit 1
}

cmd_feature() {
    if [ "$#" -lt 1 ]; then
        show_feature; exit 0
    elif [ "$#" -eq 1 ]; then
        if [ $1 == "help" ]; then
            usage; exit 0
        fi

        set_feature $1; exit 0
    fi

    usage; exit 1
}

set_feature() {
    `git config --replace-all at.feature $1`
    echo 'Feature updated to '`git @ feature`
    exit 1
}

show_feature() {
    echo `git config at.feature`
    exit 1
}