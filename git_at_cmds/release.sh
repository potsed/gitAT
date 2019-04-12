usage() {
    cat << 'EOF'

Usage:
------------------------------------------------------------------
git @ initlocal <origin-user-url> <project-name>

Once completed you should visit setting
(Gitlab EG <origin-user-url>/<project-name>/settings/repository)
And set the default branch to `develop`

EOF
    exit 0
}

cmd_release() {
    if [ "$#" -lt 1 ]; then
        usage; exit 0;
    fi

    while getopts ':Mmfdh' flag; do
        case "${flag}" in
            h) usage; exit 0 ;;
            M) local RELEASE_TYPE=1 ;; # MAJOR
            m) local RELEASE_TYPE=2 ;; # MINOR
            d) local RELEASE_TYPE=3 ;; # DANGER HOTFIX
            f) local RELEASE_TYPE=4 ;; # GENERAL BUGFIX
        esac
    done

    if [[ $RELEASE_TYPE == 1 ]]; then
        release 1
    elif [[ $RELEASE_TYPE == 2 ]]; then
        release 2
    elif [[ $RELEASE_TYPE == 3 ]]; then
        release 3
    elif [[ $RELEASE_TYPE == 4 ]]; then
        release 4
    fi
}

release() {
    # References https://stackoverflow.com/a/28804778
    local RELEASE_TYPE=$1;
    local CURRENT_BRANCH=$(git @ branch -c);
    local HAS_CHANGES=$(git @ changes);

    # If on master and there are changes let the user figure it out
    if [[ "master" == "${CURRENT_BRANCH}" && -n "${HAS_CHANGES}" ]]; then
        echo "You are on master and there are changes that need to be dealt with";
        exit 0;
    fi

    git add .;
    local GIT_STASH_REF=$(git stash create "RELEASE_AUTOSTASH");

    if [ -z "${GIT_STASH_REF}" ]; then
        echo "Nothing to stash...";
    else
        # echo the stash commit. Useful if your script terminates unexpectedly
        echo "GIT_STASH_REF created ${GIT_STASH_REF}...";
    fi

    git reset --hard;

    if [[ "master" != "${CURRENT_BRANCH}" ]]; then
        git checkout master;
    fi

    git pull;

    if [[ $RELEASE_TYPE == 1 ]]; then
        git @ version -M;
        local MSG="Major Release - Breaking Changes";
    elif [[ $RELEASE_TYPE == 2 ]]; then
        git @ version -m;
        local MSG="Minor Release";
    elif [[ $RELEASE_TYPE == 3 ]]; then
        git @ version -f;
        local MSG="Hotfix";
    elif [[ $RELEASE_TYPE == 4 ]]; then
        git @ version -f;
        local MSG="Bugfix";
    fi

    git tag -a $(git @ tag) -m $MSG;
    git push origin master --tags;

    if [[ "master" != "${CURRENT_BRANCH}" ]]; then
        git checkout $CURRENT_BRANCH;
    fi

    if [ -n "${GIT_STASH_REF}" ] ; then
        git stash apply ${GIT_STASH_REF}
    fi
}

next_step() {
    echo ""
    echo "You should now visit the settings for this project and set"
    echo "the default branch to 'develop'"
    echo ""
    echo "https://gitlab.com/squibler/$(git @ project)/settings/repository"
    echo ""
    exit 0
}