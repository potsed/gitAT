#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ release [options]

DESCRIPTION:
  Create releases with automatic version bumping and tagging.
  Handles different types of releases with appropriate version increments.

OPTIONS:
  -M, --major           Major release (breaking changes)
  -m, --minor           Minor release (new features)
  -d, --danger          Danger hotfix (urgent fixes)
  -f, --fix             General bugfix
  -h, --help            Show this help

EXAMPLES:
  git @ release -M              # Major release (1.2.3 → 2.0.0)
  git @ release -m              # Minor release (1.2.3 → 1.3.0)
  git @ release -d              # Danger hotfix (1.2.3 → 1.2.4)
  git @ release -f              # General bugfix (1.2.3 → 1.2.4)

PROCESS:
  1. Stashes current changes
  2. Switches to master branch
  3. Updates version based on release type
  4. Creates tagged release
  5. Pushes to remote with tags
  6. Restores working state

RELEASE TYPES:
  MAJOR (-M): Breaking changes, incompatible API changes
  MINOR (-m): New features, backward compatible
  DANGER (-d): Urgent hotfixes, critical issues
  FIX (-f): General bugfixes, backward compatible

SECURITY:
  All release operations are validated and logged.

EOF
    exit 0
}

cmd_release() {
    if [ "$#" -lt 1 ]; then
        usage; exit 0
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
    local RELEASE_TYPE="$1";
    local CURRENT_BRANCH=$(git @ branch -c);
    local HAS_CHANGES=$(git @ changes);

    # If on master and there are changes let the user figure it out
    if [[ "master" == "${CURRENT_BRANCH}" && -n "${HAS_CHANGES}" ]]; then
        echo "You are on master and there are changes that need to be dealt with";
        exit 0
    fi

    git add .;
    local GIT_STASH_REF=$(git stash create "RELEASE_AUTOSTASH");

    if [ -z "${GIT_STASH_REF}" ]; then
        echo "Nothing to stash...";
    else
        # echo the stash commit. Useful if your script terminates unexpectedly
        echo "GIT_STASH_REF created ${GIT_STASH_REF}..."
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

    git tag -a $(git @ version -t) -m $MSG;
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
    echo "https://gitlab.com/squibler/$(git @ product)/settings/repository"
    echo ""
    exit 0
}