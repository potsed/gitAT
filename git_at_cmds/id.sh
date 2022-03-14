cmd_id() {
    P=`git config at.project`
    M=`git config at.major`
    m=`git config at.minor`
    f=`git config at.fix`
    echo "${P}:${M}${m}${f}"
    exit 1
}