usage() {
    echo 'GitAT new repo onboarding'
    echo 'Initialises all the main settings for general use of the tool'
    exit 1
}

cmd__go() {
    # Set the base branch based on the origin
    local ORIGIN_BRANCH=`git branch -rl "*/HEAD" | rev | cut -d/ -f1 | rev`
    git @ _trunk $ORIGIN_BRANCH

    # Reset the version
    git @ version -r

    # set the current working branch
    git @ branch .

    # set the current WIP branch
    git @ wip .

    git config --replace-all at.initialised true
}