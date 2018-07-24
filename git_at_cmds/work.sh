usage() {
    echo 'Usage: git @ work'
    exit 1
}

cmd_work() {
    git fetch
    # git @ stop
    local current=`git @ branch -c`
    local branch=`git @ branch`
    if [ "$current" == "$branch" ]; then
        # git @ start
        echo "You're already in the working branch"
        echo
        exit 1;
    fi

    if [ ! `git branch --list $branch` ]; then
        git branch $branch
        exit 1;
    fi;

    git checkout $branch
    # git @ start
    exit 1;
}