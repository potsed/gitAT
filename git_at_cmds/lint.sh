usage() {
    echo 'Usage: git @ lint /path/to/lint/relative/to/app/root'
    echo 'Set the lint tool: git @ lint -t [/path/to/tool in relation to app root eg: /vendor/squizlabs/php_codesniffer/bin/phpcs]'
    exit 1
}

cmd_lint() {
    while getopts ':hic:d:' flag; do
        case "${flag}" in
            h) usage; exit 0 ;;
            i) show_tool; exit 0 ;;
            c) local cmd="${OPTARG}" ;;
            d) local dir="${OPTARG}" ;;
        esac
    done

    if [ "${cmd}" != "" ]; then
        set_tool ${cmd};
    fi

    if [ "${dir}" != "" ]; then
        set_dir ${dir};
    fi

    run_linter;





#     while getopts ":a:" opt; do
#   case $opt in
#     a)
#       echo "-a was triggered, Parameter: $OPTARG" >&2
#       ;;
#     \?)
#       echo "Invalid option: -$OPTARG" >&2
#       exit 1
#       ;;
#     :)
#       echo "Option -$OPTARG requires an argument." >&2
#       exit 1
#       ;;
#   esac
# done

    # if [ "$#" -lt 1 ]; then
    #     usage; exit 0;
    # else
    #     for i in "${helpflags[@]}"
    #     do
    #         echo $@
    #     done

    #     case $1 in
    #         "-t"|"--tool"|"tool"|"t")
    #             show_tool; exit 0
    #             ;;
    #         "-h"|"--help"|"help"|"h")
    #             usage; exit 0
    #             ;;
    #         *)
    #             run_linter $1; exit 0
    #             ;;
    #     esac
    # else
    #     case $1 in
    #         "-t"|"--tool"|"tool"|"t")
    #             set_tool $2; exit 0
    #             ;;
    #         "-h"|"--help"|"help"|"h")
    #             usage; exit 0
    #             ;;
    #     esac
    # fi
    # usage; exit 1
}

run_linter() {
    local root=`git @ root`
    # local theCMD=`git config at.linter.tool`
    local thePWD=`pwd`

    php -l ${root}/backend/app
    ${root}/backend/vendor/bin/phpcbf --standard=${root}/backend/phpcs.xml
    ${root}/backend/vendor/bin/phpcs --standard=${root}/backend/phpcs.xml


    # # local thePATH="$root$1";
    # # local theCMD="$root$theTOOL"
    # echo "Linting: $thePATH";
    # $theCMD;
    # cd $thePWD;
    # exit 1;
}


set_tool() {
    `git config --replace-all at.linter.tool $1`
    echo 'Lint tool updated'
}


set_dir() {
    `git config --replace-all at.linter.dir $1`
    echo 'Lint base directory updated'
}


set_conf() {
    `git config --replace-all at.linter.conf $1`
    echo 'Lint config updated'
}

show_tool() {
    echo "CMD "`git config at.linter.tool`
    echo
    exit 1
}