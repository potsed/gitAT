cmd_lint() {
    ROOT=`git rev-parse --show-toplevel`
    CHANGES=`git @ changes`


    echo $CHANGES
}