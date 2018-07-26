usage() {
    echo 'Usage: git @ save'
    exit 1
}

cmd_save() {
    if [ "$#" -eq 1 ]; then
        case $1 in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    save_work $@; exit 0
}

save_work() {
    local current=`git @ branch -c`
    local branch=`git @ branch`
    if [ "$current" != "$branch" ]; then
        echo 'Oops, cannot save changes.'
        echo "You don't appear to be on branch $branch"
        echo
        exit 1;
    fi

    local thePWD=`pwd`
    local root=`git @ root`

    if [ "$#" -eq 1 ]; then
        `git @ label` "$1"
    fi

    local message=`git @ label`



    cd $root
    git add . && git cz -m \""$message\""
    cd $thePWD
    exit 1;
}

# set_branch() {
#     `git config --replace-all at.branch $1`
#     echo 'Branch updated'
#     show_branch; exit 1
# }

# show_branch() {
#     echo "Current WOT Branch: "`git config at.branch`
#     echo
#     exit 1
# }