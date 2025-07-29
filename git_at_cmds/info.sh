#!/bin/bash

cmd_info() {
    # Source all the required command files to get access to their functions
    local SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    
    # Source the command files to get access to their functions
    source "$SCRIPT_DIR/product.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/version.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/feature.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/issue.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/branch.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/_path.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/_trunk.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/wip.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/_id.sh" 2>/dev/null || true
    source "$SCRIPT_DIR/_label.sh" 2>/dev/null || true

    # Call functions directly instead of using git @ commands
    local PROJECT=$(cmd_product 2>/dev/null || echo "")
    local VERSION=$(cmd_version 2>/dev/null || echo "")
    local TAG=$(cmd_version -t 2>/dev/null || echo "")
    local FEATURE=$(cmd_feature 2>/dev/null || echo "")
    local ISSUE=$(cmd_issue 2>/dev/null || echo "")
    local BRANCH=$(cmd_branch 2>/dev/null || echo "")
    local GITAT_PATH=$(cmd_path 2>/dev/null || echo "")
    local TRUNK=$(cmd_trunk 2>/dev/null || echo "")
    local WIP=$(cmd_wip 2>/dev/null || echo "")
    local GITAT_ID=$(cmd_id 2>/dev/null || echo "")
    local LABEL=$(cmd_label 2>/dev/null || echo "")

    cat << EOF

| Command              | Current Value
----------------------------------------------------
| git @ product        | ${PROJECT}
| git @ feature        | ${FEATURE}
| git @ version        | ${VERSION}
| git @ version -t     | ${TAG}
| git @ branch         | ${BRANCH}
| git @ issue          | ${ISSUE}
| git @ wip            | ${WIP}
| git @ _id            | ${GITAT_ID}
| git @ _label         | ${LABEL}
| git @ _path          | ${GITAT_PATH}
| git @ _trunk         | ${TRUNK}

EOF
    exit 0
}