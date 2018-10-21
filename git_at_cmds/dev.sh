cmd_dev() {
    local current="$(git @ branch -c)"

    if [ "$current" != "develop" ]; then
        git stash push --include-untracked -m "switched-to-develop"
        git checkout develop;
        git pull
    fi
}