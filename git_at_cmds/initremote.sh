cmd_initremote() {
    echo $1;
     if [ "$#" -lt 1 ]; then
        usage; exit 0
    elif [ "$#" -eq 1 ]; then
        git init
        git remote add origin $1;
        git add .;
        git commit -m "Initial commit";
        git push -u origin master;
    fi
    exit 0;
}

usage() {
    echo 'Usage: git @ initremote <REPO>';
    exit 1;
}