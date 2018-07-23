usage() {
    echo
    exit 1
}

cmd_tag() {
    if [ "$#" -lt 1 ]; then
        show_tag; exit 1;
    else
        case $1 in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-i"|"--inc"|"+"|"inc")
                set_tag; exit 0
                ;;
        esac
    fi

    usage; exit 1
}

set_tag() {
    OLDVER=`git @ tag`
    OLDTAG=`git config at.tag`
    NEWTAG=$(($OLDTAG + 1))
    `git config --replace-all at.tag $NEWTAG`
    echo $OLDVER 'updated to '`git @ tag`
    exit 1
}

show_tag() {
    V=`git @ version`
    echo "v$V"
    exit 1
}