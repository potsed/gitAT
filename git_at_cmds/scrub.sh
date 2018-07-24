usage() {
    echo 'Usage: git @ scrub /path/to/scrub/relative/to/app/root'
    echo 'Set the scrub tool: git @ scrub -t [/path/to/tool in relation to app root eg: /vendor/squizlabs/php_codesniffer/bin/phpcbf]'
    exit 1
}

cmd_scrub() {
    if [ "$#" -lt 1 ]; then
        usage; exit 1;
    elif [ "$#" -lt 2 ]; then
        case $1 in
            "-t"|"--tool"|"tool"|"t")
                show_tool; exit 0
                ;;
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            *)
                run_scrubber $1; exit 0
                ;;
        esac
    else
        case $1 in
            "-t"|"--tool"|"tool"|"t")
                set_tool $2; exit 0
                ;;
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    usage; exit 1
}

run_scrubber() {
    local theTOOL=`git config at.scrubber`
    local thePWD=`pwd`
    local root=`git @ root`
    # local message="$@ - "`git @ label`

    # cd $root;

    thePATH="$root$1";
    theCMD="$root$theTOOL"
    echo "Scrubbing: $thePATH";
    $theCMD $thePATH;
    exit 1;
}

set_tool() {
    `git config --replace-all at.scrubber $1`
    echo 'Scrub tool updated'
    show_tool; exit 1
}

show_tool() {
    echo "Current Scrub Tool "`git config at.scrubber`
    echo
    exit 1
}