cmd_dev() {
    local current="$(git @ branch -c)"

    if [ "$current" != "develop" ]; then
        git stash
        git checkout develop;
        git pull
        git stash pop
    fi
}