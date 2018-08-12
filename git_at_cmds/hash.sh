
cmd_hash() {
    git fetch
    BRANCH=`git @ branch -c`
    REM_HAS_MASTER=`git ls-remote --heads origin master | wc -l | tr -d ' '`
    REM_HAS_DEVELOP=`git ls-remote --heads origin develop | wc -l | tr -d ' '`
    REM_HAS_CURRENT=`git ls-remote --heads origin ${BRANCH} | wc -l | tr -d ' '`

    REM_CURRENT_HASH="UNKNOWN"
    CURRENT_BEHIND_REM="0"
    CURRENT_AHEAD_REM="0"

    REM_DEVELOP_BR_HASH="UNKNOWN"
    LOC_BR_BEHIND_REM_DEVELOP="0"
    LOC_BR_AHEAD_REM_DEVELOP="0"

    REM_MASTER_BR_HASH="UNKNOWN"
    LOC_BR_BEHIND_REM_MASTER="0"
    LOC_BR_AHEAD_REM_MASTER="0"

    LOC_BR_HASH=`git rev-parse --short HEAD`
    LOC_DEVELOP_HASH=`git rev-parse --short develop`
    LOC_MASTER_HASH=`git rev-parse --short master`

    if [ "$REM_HAS_CURRENT" == "1" ]; then
        REMBR="origin/${BRANCH}"
        REM_CURRENT_HASH=`git rev-parse --short ${REMBR}`
        CURRENT_BEHIND_REM=`git rev-list --left-right --count ${REMBR}...@ | cut -f1`
        CURRENT_AHEAD_REM=`git rev-list --left-right --count ${REMBR}...@ | cut -f2`
    fi

    echo
    echo "BRANCH ${BRANCH}"
    echo "- CURRENT HASH: ${LOC_BR_HASH}"
    echo "- ORIGIN HASH: ${REM_CURRENT_HASH}"
    echo "CURRENT is ${CURRENT_BEHIND_REM} behind and ${CURRENT_AHEAD_REM} ahead of origin/${BRANCH}"

    if [ "$REM_HAS_DEVELOP" == "1" ]; then
        LOC_BR_BEHIND_REM_DEVELOP=`git rev-list --left-right --count origin/develop...@ | cut -f1`
        LOC_BR_BEHIND_REM_DEVELOP=`git rev-list --left-right --count origin/develop...@ | cut -f2`
        echo "CURRENT is ${LOC_BR_BEHIND_REM_DEVELOP} behind and ${LOC_BR_BEHIND_REM_DEVELOP} ahead of origin/develop"
    fi

    if [ "$REM_HAS_MASTER" == "1" ]; then
        LOC_BR_BEHIND_REM_MASTER=`git rev-list --left-right --count origin/master...@ | cut -f1`
        LOC_BR_AHEAD_REM_MASTER=`git rev-list --left-right --count origin/master...@ | cut -f2`
        echo "CURRENT is ${LOC_BR_BEHIND_REM_DEVELOP} behind and ${LOC_BR_BEHIND_REM_DEVELOP} ahead of origin/master"
    fi
    echo

    exit 1
}