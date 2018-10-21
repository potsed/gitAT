usage() {
    echo "View current tag: git @ tag"
    echo "Increase tag number: git @ tag +"
    echo "Decrease tag number: git @ tag -"
    echo "Show this help: git @ tag -h"
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
            "-a"|"--add"|"+"|"add"|"up")
                add_tag; exit 0
                ;;
            "-s"|"--sub"|"-"|"sub"|"down")
                sub_tag; exit 0
                ;;
        esac
    fi

    usage; exit 1
}

add_tag() {
    OLDVER=`git @ tag`
    OLDTAG=`git config at.tag`
    NEWTAG=$(($OLDTAG + 1))
    git config --replace-all at.tag ${NEWTAG}
    echo ${OLDVER} 'updated to '`git @ tag`
    exit 1
}

sub_tag() {
    OLDVER=`git @ tag`;
    OLDTAG=`git config at.tag`;
    if [ $OLDTAG -gt 0 ]; then
        NEWTAG=$(($OLDTAG - 1));
    else
        NEWTAG=0;
    fi
    `git config --replace-all at.tag $NEWTAG`
    echo $OLDVER 'updated to '`git @ tag`
    exit 1
}

show_tag() {
    V=`git @ version`
    echo "v$V"
    exit 1
}