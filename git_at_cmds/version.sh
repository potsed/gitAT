usage() {
    echo "View current version: git @ version"
    echo "Set major version Number: git @ version -M 1"
    echo "Set minor version Number: git @ version -m 1"
    echo "Show this help: git @ version -h"
    exit 1
}

cmd_version() {
    if [ "$#" -lt 1 ]; then
        show_version; exit 1;
    elif [ "$#" -lt 2 ]; then
        usage; exit 1;
    else
        case $1 in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-m"|"--minor"|"m"|"minor")
                set_minor $2; exit 0
                ;;
            "-M"|"--major"|"M"|"major")
                set_major $2; exit 0
                ;;
        esac
    fi

    usage; exit 1
}

set_major() {
    `git config --replace-all at.version $1`
    set_minor 0;
    echo 'Version updated to '`git @ tag`
    exit 1
}

set_minor() {
    `git config --replace-all at.minor $1`
    `git config --replace-all at.tag 0`
    echo 'Version updated to '`git @ tag`
    exit 1
}

show_version() {
    V=`git config at.version`
    v=`git config at.minor`
    t=`git config at.tag`
    echo $V'.'$v'.'$t
    exit 1
}