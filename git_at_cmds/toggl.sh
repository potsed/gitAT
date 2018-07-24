usage() {
    echo "git @ toggl [options]"
    echo
    exit 0
}

start() {
    toggl start `git @ label`
}

stop() {
    toggl stop
}

current() {
    toggl current
}

cmd_toggl() {
    if [ "$#" -lt 1 ]; then
        start; exit 0
    fi

    local subcommand="$1"; shift

    case $subcommand in
        "s"|"stop"|"--stop")
            stop; exit 0
            ;;
        "c"|"-c"|"--current")
            current; exit 0
            ;;
    esac

    usage
}