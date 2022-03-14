cmd_master() {
    local current="$(git @ branch -c)"
    local base="$(git @ _trunk)"
    if [ "$current" != "$base" ]; then
        git stash push --include-untracked -m "switched-to-master"
        git checkout $base;
        git pull
    fi
}
