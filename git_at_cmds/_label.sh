usage() {
    echo
    exit 1
}

cmd__label() {
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

    echo 'Label updated to: '`git @ _label`
    return;
}

show_label() {
    echo "[\
PRODUCT: `git @ product`; \
FEATURE: `git @ feature`; \
ISSUE: `git @ issue`; \
VERSION: `git @ version -t` \
]"
    return;

}