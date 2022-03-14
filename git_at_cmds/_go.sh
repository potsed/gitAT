usage() {
    echo 'Runs at first use in a repo'
    echo 'Initialises all the main settings for general use of the tool'
    exit 1
}

cmd__go() {
    # Set the base branch based on the origin
    local ORIGIN_BRANCH=`git branch -rl "*/HEAD" | rev | cut -d/ -f1 | rev`
    git @ base $ORIGIN_BRANCH

    # Reset the version
    git @ version -r

    # set the current working branch
    git @ branch .



    git config --replace-all at.initialised true
}