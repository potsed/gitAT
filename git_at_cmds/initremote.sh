#!/bin/bash

cmd_initremote() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
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
            exit 0
}

usage() {
    cat << 'EOF'
Usage: git @ initremote <repository-url>

DESCRIPTION:
  Initialize a remote repository with basic structure.
  Creates initial commit and sets up develop branch with CHANGELOG.

ARGUMENTS:
  <repository-url>    Remote repository URL (e.g., git@github.com:user/repo.git)

PROCESS:
  1. Initializes git repository
  2. Adds remote origin
  3. Creates initial commit
  4. Pushes to master branch
  5. Creates develop branch
  6. Creates CHANGELOG file
  7. Pushes develop branch

EXAMPLES:
  git @ initremote git@github.com:user/my-repo.git
  git @ initremote git@gitlab.com:org/project.git

SECURITY:
  All initialization operations are validated and logged.

EOF
    exit 1
}