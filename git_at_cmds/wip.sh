usage() {
    cat << 'EOF'
          __
       __/\ \__          __                    __
   __ /\_\ \ ,_\        /'_`\_      __  __  __/\_\  _____
 /'_ `\/\ \ \ \/       /'/'_` \    /\ \/\ \/\ \/\ \/\ '__`\
/\ \L\ \ \ \ \ \_     /\ \ \L\ \   \ \ \_/ \_/ \ \ \ \ \L\ \
\ \____ \ \_\ \__\    \ \ `\__,_\   \ \___x___/'\ \_\ \ ,__/
 \/___L\ \/_/\/__/     \ `\_____\    \/__//__/   \/_/\ \ \/
   /\____/              `\/_____/                     \ \_\
   \_/__/

Usage:
---------
git @ wip [<options>]

Without options, echo out the current branch stored as a work in progress.

Options:
---------
    --help|-h           Show this help screen
    --set|-s            Sets the current `git @ branch` to the wip
    --checkout|-c       Checks out the wip branch
    --restore|-r        Restores the wip branch to the primary `git @ branch`

EOF
    exit 1
}

cmd_wip() {
    if [ "$#" -lt 1 ]; then
        show_wip; exit 0
    elif [ "$#" -eq 1 ]; then
        case $1 in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-s"|"--set"|"s"|"set"|".")
                set_wip; exit 0
                ;;
            "-c"|"--checkout"|"c"|"checkout")
                checkout_wip; exit 0
                ;;
        esac
    fi

    usage; exit 1
}

restore_wip() {
    # echo `git rev-parse --abbrev-ref HEAD`; exit 0
    from=`git @ wip`
    git @ branch $from
    git @ work
}

set_wip() {
    from=`git config at.wip`
    branch=`git @ branch -c`
    `git config --replace-all at.wip $branch`
    echo "wip updated to $branch from $from"
}

show_wip() {
    echo "Current WIP is: "`git config at.wip`
}

checkout_wip() {
    from=`git config at.wip`
    git checkout $from
}

