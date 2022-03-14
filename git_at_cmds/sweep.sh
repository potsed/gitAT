cmd_sweep() {
    echo "Deleting local copies of branches merged into `git @ base`"
    delete_merged_branches_locally;
    echo "Done"
    exit 0
}

delete_merged_branches_locally() {
    git branch --merged | grep -v '\*\|master\|main\|dev|develop' | xargs -n 1 git branch -d
}