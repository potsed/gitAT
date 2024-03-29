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

    if [ "$current" == "master" ] || [ "$current" == "develop" ]; then
        echo "Oops, cannot save changes on ${current} create a new branch instead!"
        echo
        exit 1;
    fi

    if [ "$current" == "prod" ]; then
        read -p "Are you sure you want to commit this? (Y/N): " CONFIRMATION && [[ $CONFIRMATION == [yY] || $CONFIRMATION == [yY][eE][sS] ]] || exit 1;
        git @ version -t +
        git tag `git @ version -t`

    elif [ "$current" != "$branch" ]; then
        echo "Oops cannot save the changes, you're not on the correct working branch $branch"
        echo
        exit 1;
    fi

    if [ "$#" -eq 1 ]; then
        git @ _label "$1";
    fi

    local thePWD=`pwd`
    local path=`git @ _path`
    local message=`git @ _label`

    cd $path
    git add -p && git cz -m \""$message\""
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