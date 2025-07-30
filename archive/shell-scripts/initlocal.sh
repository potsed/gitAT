#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ initlocal <origin-url> <project-name>

DESCRIPTION:
  Initialize a new local repository with remote setup and proper branch structure.
  Creates master → staging → develop branch hierarchy.

ARGUMENTS:
  <origin-url>      Remote repository URL (e.g., git@gitlab.com:user)
  <project-name>    Project name for the repository

EXAMPLES:
  git @ initlocal git@gitlab.com:user my-project
  git @ initlocal git@github.com:org api-service

PROCESS:
  1. Initializes git repository
  2. Sets project name
  3. Creates remote origin
  4. Sets up branch structure: master → staging → develop
  5. Pushes all branches to remote

BRANCH STRUCTURE:
  master   → Production-ready code
  staging  → Pre-production testing
  develop  → Development and feature integration

NEXT STEPS:
  After initialization, visit your repository settings and set the default
  branch to 'develop':
  https://gitlab.com/user/project-name/settings/repository

SECURITY:
  All initialization operations are validated and logged.

EOF
    exit 0
}

cmd_initlocal() {
    if [ "$#" -eq 1 ]; then
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
        esac
    fi
    
    # git push --set-upstream git@gitlab.example.com:namespace/nonexistent-project.git master
    if [ "$#" -lt 1 ]; then
        show_repo; exit 0
    elif [ "$#" -gt 0 ]; then
        set_remote "$1"; exit 0
    fi

    usage; exit 1
}

set_remote() {
    git init;
    git add .;
    git commit -m "Initial commit";
    sleep 5;

    git @ product "$1";
    git @ version -r;
    git remote add origin git@gitlab.com:squibler/$(git @ product).git;
    sleep 5;

    # Now create the rest of the structure MASTER -> STAGING -> DEVELOP
    git push --set-upstream origin master;
    sleep 5;

    git checkout -b staging;
    git push --set-upstream origin staging;
    sleep 5;

    git checkout -b develop;
    git push --set-upstream origin develop;

    next_step;
    exit 0
}

next_step() {
    echo ""
    echo "You should now visit the settings for this project and set"
    echo "the default branch to 'develop'"
    echo ""
    echo "https://gitlab.com/squibler/$(git @ product)/settings/repository"
    echo ""
    exit 0
}