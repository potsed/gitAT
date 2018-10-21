cmd_master() {
    local current="$(git @ branch -c)"
    if [ "$current" != "master" ]; then
        git stash push --all -m "switched-to-master"
        git checkout master;
        git pull
    fi
}
