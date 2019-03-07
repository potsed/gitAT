usage() {
    cat << 'EOF'
       _ _         ____                                _
  __ _(_) |_      / __ \     ___  __ _ _   _  __ _ ___| |__
 / _` | | __|    / / _` |   / __|/ _` | | | |/ _` / __| '_ \
| (_| | | |_    | | (_| |   \__ \ (_| | |_| | (_| \__ \ | | |
 \__, |_|\__|    \ \__,_|   |___/\__, |\__,_|\__,_|___/_| |_|
 |___/            \____/            |_|


Usage:
--------------------------------------------------------------------------------
  git @ squash [options] <head_branch>

Description:
--------------------------------------------------------------------------------
  Reset the branch you're currently working in to the head of another branch
  (usually develop or master) to make a clean commit without all of the
  working history.

  NB. You may have to force push your branch to the repo after squashing the
  commits, if on a shared branch this may not be best option.

  Why is it useful?
  ----------------------
  When you have been commiting locally for a period of time working on one
  feature, then you may wish to clean up the repo before pushing it for PR

Options:
--------------------------------------------------------------------------------
  --help|-h           Show this help screen
  --save|-s           Run `git @ save` after squashing

GIT Commands Used
--------------------------------------------------------------------------------
  For transparency and education purposes I like to disclose the underlying
  git commands being run during this operation.

  01. Retrieve the HEAD SHA of the given branch
      CMD: `git rev-parse --verify --quiet --long ${BRANCH}`

  02. Resets the current branch to the SHA of the given branch without losing
      any of the work you have done since that SHA
      CMD: `git reset --soft ${SHA}`

EOF
    exit 1
}

cmd_squash() {
    while getopts ':hs' flag; do
        case "${flag}" in
            h) usage; exit 0 ;;
            s) local DOSAVE=1 ;;
        esac
    done
    shift $(expr ${OPTIND} - 1)

    if [ "${1}" == "" ]; then
        usage; exit 0;
    fi

    local HEAD="$(head "$1")"
    if [ "${HEAD}" == "0" ]; then
        echo
        echo "ERROR: Branch \"${1}\" does not exist locally"
        echo
        exit 0;
    fi

    squash $HEAD;

    echo
    echo 'Squashed branch '"$(git @ branch -c)"' back to '${1}

    if [ "${DOSAVE}" == "1" ]; then
        save
    fi
}

head() {
    local HEAD=$(git rev-parse --verify --quiet --long "$1");
    if [ "${HEAD}" != "" ]; then
        echo ${HEAD};
    else
        echo 0;
    fi
}

# exists() {

# }

squash() {
    git reset --soft ${1}
}

save() {
    git @ save
}