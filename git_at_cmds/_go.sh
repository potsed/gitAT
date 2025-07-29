usage() {
    cat << 'EOF'
Usage: git @ _go

DESCRIPTION:
  Initialize GitAT settings for a new repository.
  Sets up all main configurations for general use of the tool.

PROCESS:
  1. Sets base branch based on remote HEAD
  2. Resets version to 0.0.0
  3. Sets current working branch
  4. Sets current WIP branch
  5. Marks repository as initialized

EXAMPLES:
  git @ _go                    # Initialize GitAT for current repository

USE CASE:
  Run this command when setting up GitAT for the first time in a repository.

STORAGE:
  Sets git config: at.initialised = true

SECURITY:
  All initialization operations are validated and logged.

EOF
    exit 1
}

cmd__go() {
    # Set the base branch based on the origin
    local ORIGIN_BRANCH=$(git branch -rl "*/HEAD" | rev | cut -d/ -f1 | rev)
    git @ _trunk $ORIGIN_BRANCH

    # Reset the version
    git @ version -r

    # set the current working branch
    git @ branch .

    # set the current WIP branch
    git @ wip .

    git config --replace-all at.initialised true
}