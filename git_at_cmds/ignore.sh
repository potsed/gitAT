cmd_ignore() {
     if [ "$#" -lt 1 ]; then
        usage; exit 0
    elif [ "$#" -eq 1 ]; then
        PATH=`git @ path`"/.gitignore";
        echo $1 >> $PATH;
        echo 'String appended to .gitignore';
    fi
    exit 0;
}

usage() {
    echo 'Usage: git @ ignore ignore/pathto/file';
    exit 1;
}