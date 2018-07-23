usage() {
    echo 'Usage: git @ save "Your commit message"'
    exit 1
}

cmd_save() {
    if [ "$#" -eq 1 ]; then
        case $1 in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
        save_work $@; exit 0
    fi
    usage; exit 1
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
    local message="$@ - "`git @ label`

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