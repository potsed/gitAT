usage() {
    echo "View current version: git @ version"
    echo "Increment major version Number: git @ version -M"
    echo "Increment minor version Number: git @ version -m 1"
    echo "Increment fix version Number: git @ version -f 1"
    echo "Show this help: git @ version -h"
    echo "DO NOT MISUSE: Reset the version to 0.0.0: git @ version -r"
    exit 0
}

cmd_version() {
    if [ "$#" -lt 1 ]; then
        show_version; exit 0;
    fi

    while getopts ':hrMmf' flag; do
        case "${flag}" in
            h) usage; exit 0 ;;
            r) reset_version; exit 0 ;;
            M) local MAJOR=true ;;
            m) local MINOR=true ;;
            f) local FIX=true ;;
        esac
    done

    if [[ $MAJOR ]]; then
        increment_major;
    fi

    if [[ $MINOR ]]; then
        increment_minor;
    fi

    if [[ $FIX ]]; then
        increment_fix;
    fi

    show_version;
}

reset_version() {
    git config --replace-all at.major 0;
    git config --replace-all at.minor 0;
    git config --replace-all at.fix 0;
}

increment_major() {
    local OLD_MAJOR=$(git config at.major);
    set_major $(($OLD_MAJOR + 1))
}

increment_minor() {
    local OLD_MINOR=$(git config at.minor);
    set_minor $(($OLD_MINOR + 1))
}

increment_fix() {
    local OLD_FIX=$(git config at.fix);
    set_fix $(($OLD_FIX + 1))
}

set_major() {
    git config --replace-all at.major $1
    set_minor 0;
}

set_minor() {
    git config --replace-all at.minor $1
    set_fix 0;
}

set_fix() {
    git config --replace-all at.fix $1;
}

show_version() {
    local M=$(git config at.major);
    local m=$(git config at.minor);
    local f=$(git config at.fix);
    echo $M'.'$m'.'$f;
}