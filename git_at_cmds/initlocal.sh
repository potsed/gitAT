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

cmd_initlocal() {
    # git push --set-upstream git@gitlab.example.com:namespace/nonexistent-project.git master
    if [ "$#" -lt 1 ]; then
        show_repo; exit 0
    elif [ "$#" -gt 0 ]; then
        if [ $1 == "help" ]; then
            usage; exit 0
        fi

        set_remote $1; exit 0
    fi

    usage; exit 1
}

set_remote() {
    git init;
    git @ project $1

    local ORIGIN=git@gitlab.com:squibler/$(git @ project).git

    git remote add origin $ORIGIN

    git add .
    git commit -m "Initial commit"
    git @ version -r;
    git push --set-upstream origin master;

    git @ version -f;
    git tag $(git @ tag);

    # Now create the rest of the structure MASTER -> STAGING -> DEVELOP
    git checkout -b staging
    git push --set-upstream origin staging;

    git checkout -b develop
    git push --set-upstream origin develop;

    next_step;
    exit 0
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