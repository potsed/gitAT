cmd_info() {

    local PROJECT=`git @ product`
    local VERSION=`git @ version`
    local TAG=`git @ version -t`
    local FEATURE=`git @ feature`
    local ISSUE=`git @ issue`
    local BRANCH=`git @ branch`
    local GITAT_PATH=`git @ _path`
    local TRUNK=`git @ _trunk`
    local WIP=`git @ wip`
    local GITAT_ID=`git @ _id`
    local LABEL=`git @ _label`

    cat << EOF

| Command              | Current Value
----------------------------------------------------
| git @ product        | ${PROJECT}
| git @ feature        | ${FEATURE}
| git @ version        | ${VERSION}
| git @ version -t     | ${TAG}
| git @ branch         | ${BRANCH}
| git @ issue          | ${ISSUE}
| git @ wip            | ${WIP}
| git @ _id            | ${GITAT_ID}
| git @ _label         | ${LABEL}
| git @ _path          | ${GITAT_PATH}
| git @ _trunk         | ${TRUNK}

EOF
    exit 1
}