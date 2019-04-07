usage() {
    echo "View current tag: git @ tag"
    echo "Show this help: git @ tag -h"
    exit 0
}

cmd_tag() {
    if [ "$#" -lt 1 ]; then
        show_tag; exit 0;
    fi

    usage; exit 1
}

show_tag() {
    local V=`git @ version`
    echo "v$V"
    exit 0
}