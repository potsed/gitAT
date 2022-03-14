cmd_info() {

    local PROJECT=`git @ project`
    local VERSION=`git @ version`
    local TAG=`git @ tag`
    local FEATURE=`git @ feature`
    local ISSUE=`git @ issue`
    local BRANCH=`git @ branch`
    local GITAT_PATH=`git @ _path`
    local TRUNK=`git @ _trunk`
    local WIP=`git @ wip`
    local GITAT_ID=`git @ id`
    local LABEL=`git @ label`

    cat << EOF

| Command       | Current Value
---------------------------------------
| git @ project | ${PROJECT}
| git @ feature | ${FEATURE}
| git @ version | ${VERSION}
| git @ tag     | ${TAG}
| git @ branch  | ${BRANCH}
| git @ issue   | ${ISSUE}
| git @ wip     | ${WIP}
| git @ id      | ${GITAT_ID}
| git @ label   | ${LABEL}
| git @ _path   | ${GITAT_PATH}
| git @ _trunk  | ${TRUNK}

EOF
    exit 1
}