cmd_id() {
    V=`git config at.version`
    v=`git config at.minor`
    T=`git config at.tag`
    echo "${V}-${v}-${T}"
    exit 1
}