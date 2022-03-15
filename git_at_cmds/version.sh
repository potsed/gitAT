usage() {
    cat << EOF

View current version: ................................ git @ version
View current tagged version: ......................... git @ version -t
Show this help: ...................................... git @ version -h

Increment Version
MAJOR version number: ................................ git @ version -M
MINOR version mnumber: ............................... git @ version -m
BUMP version number: ................................. git @ version -b

BE CAREFUL WITH THIS ONE:
Reset the version to 0.0.0: .......................... git @ version -r

EOF
    exit 0
}

cmd_version() {
    if [ "$#" -lt 1 ]; then
        show_version; exit 0;
    fi

    while getopts ':htrMmb' flag; do
        case "${flag}" in
            h) usage; exit 0 ;;
            t) show_tag; exit 0;;
            r) reset_version; exit 0 ;;
            M) local MAJOR=true ;;
            m) local MINOR=true ;;
            b) local BUMP=true ;;
        esac
    done

    if [[ $MAJOR ]]; then
        increment_major;
    fi

    if [[ $MINOR ]]; then
        increment_minor;
    fi

    if [[ $BUMP ]]; then
        increment_fix;
    fi

    show_version;
}

show_tag() {
    echo "v"`git @ version`
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
    local OLD_BUMP=$(git config at.fix);
    set_fix $(($OLD_BUMP + 1))
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