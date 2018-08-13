usage() {
    cat <<'EOF'

          __                                                 __
       __/\ \__          __                                 /\ \
   __ /\_\ \ ,_\        /'_`\_      __  __  __    ___   _ __\ \ \/'\
 /'_ `\/\ \ \ \/       /'/'_` \    /\ \/\ \/\ \  / __`\/\`'__\ \ , <
/\ \L\ \ \ \ \ \_     /\ \ \L\ \   \ \ \_/ \_/ \/\ \L\ \ \ \/ \ \ \\`\
\ \____ \ \_\ \__\    \ \ `\__,_\   \ \___x___/'\ \____/\ \_\  \ \_\ \_\
 \/___L\ \/_/\/__/     \ `\_____\    \/__//__/   \/___/  \/_/   \/_/\/_/
   /\____/              `\/_____/
   \_/__/

EOF
    echo 'Usage: git @ work'
    exit 1
}

cmd_work() {
    git fetch
    # git @ stop
    local current=`git @ branch -c`
    local branch=`git @ branch`
    if [ "$current" == "$branch" ]; then
        # git @ start
        echo "You're already in the working branch"
        echo
        exit 1;
    fi

    if [ ! `git branch --list $branch` ]; then
        AUTOSTASH=1
        git stash;
        git checkout develop;
        git pull;
        git branch $branch;
    fi;

    git checkout $branch

    if [ "$AUTOSTASH" == 1 ]; then
        git stash pop
    fi;
    # git @ start
    exit 1;
}