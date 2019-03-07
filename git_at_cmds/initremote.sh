cmd_initremote() {
    echo $1;
     if [ "$#" -lt 1 ]; then
        usage; exit 0
    elif [ "$#" -eq 1 ]; then
        git init
        git remote add origin $1
        git add .
        git commit -m "Initial commit"
        git push --set-upstream origin master
        git checkout -b develop
        touch CHANGELOG
        echo "[$(date)]\r- CHANGELOG CREATED\r- INITIAL COMMIT\r\r" >> CHANGELOG
        git add .
        git commit -m "CHANGELOG CREATED"
        git push --set-upstream origin develop;
    fi
    exit 0;
}

usage() {
    echo 'Usage: git @ initremote <REPO>';
    exit 1;
}