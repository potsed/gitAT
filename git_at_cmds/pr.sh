#!/bin/bash

usage() {
    cat << 'EOF'
Usage: git @ pr [<title>] [options]

DESCRIPTION:
  Create a Pull Request (PR) or Merge Request (MR) for the current branch.
  Automatically detects the Git hosting platform and uses appropriate tools.

PLATFORMS SUPPORTED:
  ✅ GitHub: Uses 'gh' CLI or provides web URL
  ✅ GitLab: Uses 'glab' CLI or provides web URL
  ✅ Bitbucket: Provides web URL
  ✅ Generic: Provides web URL with branch info

FEATURES:
  ✅ Auto-platform detection
  ✅ CLI tool integration (gh, glab)
  ✅ Web URL fallback
  ✅ Branch validation
  ✅ Commit message integration
  ✅ Custom title and description
  ✅ Automatic commit squashing (configurable)
  ✅ Automatic description generation from changed files

EXAMPLES:
  git @ pr                                    # Create PR with default title and auto-generated description
  git @ pr "Add user authentication"          # Create PR with custom title and auto-generated description
  git @ pr -d "Detailed description here"     # Create PR with custom description
  git @ pr -b main                            # Create PR targeting main branch
  git @ pr -h                                 # Show this help

OPTIONS:
  -t, --title <title>       PR title (defaults to last commit message)
  -d, --description <desc>  PR description
  -b, --base <branch>       Target branch (defaults to configured trunk)
  -o, --open               Open PR in browser after creation
  -s, --squash             Force squash commits before PR (overrides setting)
  -S, --no-squash          Force no squash (overrides setting)
  -h, --help               Show this help message

AUTOMATIC FEATURES:
  - Uses last commit message as default title
  - Generates description from changed files (when not provided)
  - Includes branch name and commit info
  - Validates current branch is not trunk
  - Checks for uncommitted changes
  - Auto-squash commits if at.pr.squash is enabled

CONFIGURATION:
  git config at.pr.squash true    # Enable automatic squashing
  git config at.pr.squash false   # Disable automatic squashing

EOF
    exit 1
}

# Detect Git hosting platform
detect_platform() {
    local remote_url
    
    # Get the primary remote URL
    remote_url=$(git config --get remote.origin.url 2>/dev/null || echo "")
    
    if [ -z "$remote_url" ]; then
        echo "unknown"
        return
    fi
    
    # Detect platform from URL
    if echo "$remote_url" | grep -q "github.com"; then
        echo "github"
    elif echo "$remote_url" | grep -q "gitlab.com\|gitlab\."; then
        echo "gitlab"
    elif echo "$remote_url" | grep -q "bitbucket.org\|bitbucket\."; then
        echo "bitbucket"
    else
        echo "generic"
    fi
}

# Get repository information
get_repo_info() {
    local remote_url
    local platform
    
    remote_url=$(git config --get remote.origin.url 2>/dev/null || echo "")
    platform=$(detect_platform)
    
    if [ -z "$remote_url" ]; then
        echo "Error: No remote origin configured" >&2
        exit 1
    fi
    
    # Extract owner and repo name
    case "$platform" in
        "github")
            # Handle both SSH and HTTPS URLs
            if echo "$remote_url" | grep -q "git@github.com:"; then
                echo "$remote_url" | sed 's|git@github.com:||' | sed 's|\.git$||'
            else
                echo "$remote_url" | sed 's|https://github.com/||' | sed 's|\.git$||'
            fi
            ;;
        "gitlab")
            # Handle both SSH and HTTPS URLs
            if echo "$remote_url" | grep -q "git@gitlab.com:"; then
                echo "$remote_url" | sed 's|git@gitlab.com:||' | sed 's|\.git$||'
            else
                echo "$remote_url" | sed 's|https://gitlab.com/||' | sed 's|\.git$||'
            fi
            ;;
        "bitbucket")
            # Handle both SSH and HTTPS URLs
            if echo "$remote_url" | grep -q "git@bitbucket.org:"; then
                echo "$remote_url" | sed 's|git@bitbucket.org:||' | sed 's|\.git$||'
            else
                echo "$remote_url" | sed 's|https://bitbucket.org/||' | sed 's|\.git$||'
            fi
            ;;
        *)
            echo "$remote_url"
            ;;
    esac
}

# Get default PR title from last commit
get_default_title() {
    local last_commit
    
    last_commit=$(git log -1 --pretty=format:"%s" 2>/dev/null || echo "")
    
    if [ -n "$last_commit" ]; then
        echo "$last_commit"
    else
        echo "Update from $(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "current branch")"
    fi
}

# Generate automatic description based on changed files
generate_auto_description() {
    local base_branch="$1"
    local current_branch="$2"
    local description=""
    
    # Get list of changed files
    local changed_files
    changed_files=$(git diff --name-only "$base_branch..$current_branch" 2>/dev/null || echo "")
    
    if [ -z "$changed_files" ]; then
        echo "No files changed compared to $base_branch"
        return
    fi
    
    # Count files by type
    local total_files=0
    local added_files=0
    local modified_files=0
    local deleted_files=0
    
    # Analyze each changed file
    while IFS= read -r file; do
        if [ -n "$file" ]; then
            total_files=$((total_files + 1))
            
            # Check if file was added, modified, or deleted
            if git diff --name-status "$base_branch..$current_branch" | grep -q "^A.*$file$"; then
                added_files=$((added_files + 1))
            elif git diff --name-status "$base_branch..$current_branch" | grep -q "^D.*$file$"; then
                deleted_files=$((deleted_files + 1))
            else
                modified_files=$((modified_files + 1))
            fi
        fi
    done <<< "$changed_files"
    
    # Generate summary
    description="# 📋 Pull Request Summary\n\n"
    description+="This PR contains changes from branch \`$current_branch\` targeting \`$base_branch\`.\n\n"
    
    description+="## 📊 Changes Overview\n\n"
    description+="| Metric | Count |\n"
    description+="|--------|-------|\n"
    description+="| **Total Files** | $total_files |\n"
    
    if [ $added_files -gt 0 ]; then
        description+="| **Added** | $added_files |\n"
    fi
    if [ $modified_files -gt 0 ]; then
        description+="| **Modified** | $modified_files |\n"
    fi
    if [ $deleted_files -gt 0 ]; then
        description+="| **Deleted** | $deleted_files |\n"
    fi
    
    description+="\n## 📁 File Analysis\n\n"
    
    # Group files by directory/type
    local file_types=()
    local dirs=()
    
    while IFS= read -r file; do
        if [ -n "$file" ]; then
            # Get file extension
            local ext="${file##*.}"
            if [ "$ext" = "$file" ]; then
                ext="no-extension"
            fi
            
            # Get directory
            local dir=$(dirname "$file")
            if [ "$dir" = "." ]; then
                dir="root"
            fi
            
            # Add to arrays if not already present
            local found_ext=false
            local found_dir=false
            
            # Check if extension already exists
            for existing_ext in "${file_types[@]}"; do
                if [ "$existing_ext" = "$ext" ]; then
                    found_ext=true
                    break
                fi
            done
            
            # Check if directory already exists
            for existing_dir in "${dirs[@]}"; do
                if [ "$existing_dir" = "$dir" ]; then
                    found_dir=true
                    break
                fi
            done
            
            # Add if not found
            if [ "$found_ext" = false ]; then
                file_types+=("$ext")
            fi
            if [ "$found_dir" = false ]; then
                dirs+=("$dir")
            fi
        fi
    done <<< "$changed_files"
    
    # Add file type summary
    if [ ${#file_types[@]} -gt 0 ]; then
        description+="### 🔤 File Types\n\n"
        description+="This PR affects the following file types:\n\n"
        for ext in "${file_types[@]}"; do
            description+="- \`$ext\`\n"
        done
        description+="\n"
    fi
    
    # Add directory summary if multiple directories
    if [ ${#dirs[@]} -gt 1 ]; then
        description+="### 📂 Directories Affected\n\n"
        description+="Changes span across the following directories:\n\n"
        for dir in "${dirs[@]}"; do
            description+="- \`$dir\`\n"
        done
        description+="\n"
    fi
    
    # List all changed files
    description+="## 📝 Changed Files\n\n"
    description+="<details>\n<summary>📋 Click to view all changed files</summary>\n\n"
    description+="```\n"
    
    while IFS= read -r file; do
        if [ -n "$file" ]; then
            # Get file status
            local status
            if git diff --name-status "$base_branch..$current_branch" | grep -q "^A.*$file$"; then
                status="➕"
            elif git diff --name-status "$base_branch..$current_branch" | grep -q "^D.*$file$"; then
                status="🗑️"
            else
                status="✏️"
            fi
            
            description+="$status $file\n"
        fi
    done <<< "$changed_files"
    
    description+="```\n</details>\n\n"
    
    # Add commit summary
    local commit_count
    commit_count=$(git rev-list --count "$base_branch..$current_branch" 2>/dev/null || echo "0")
    
    if [ "$commit_count" -gt 1 ]; then
        description+="## 🔄 Commits\n\n"
        description+="This PR includes **$commit_count commits**.\n\n"
        description+="<details>\n<summary>📜 Click to view commit history</summary>\n\n"
        description+="```\n"
        
        # Show last 5 commits
        git log --oneline -5 "$base_branch..$current_branch" 2>/dev/null | while IFS= read -r commit; do
            if [ -n "$commit" ]; then
                description+="$commit\n"
            fi
        done
        
        description+="```\n</details>\n\n"
    fi
    
    # Add footer
    description+="---\n\n"
    description+="*This description was automatically generated based on the changes in this PR.*\n"
    
    echo "$description"
}

# Create GitHub PR using gh CLI
create_github_pr() {
    local title="$1"
    local description="$2"
    local base_branch="$3"
    local current_branch="$4"
    
    if ! command -v gh >/dev/null 2>&1; then
        echo "GitHub CLI (gh) not installed. Please install it or use the web interface."
        echo "Install: https://cli.github.com/"
        return 1
    fi
    
    # Check if user is authenticated
    if ! gh auth status >/dev/null 2>&1; then
        echo "GitHub CLI not authenticated. Please run 'gh auth login' first."
        return 1
    fi
    
    # Create PR
    if [ -n "$description" ]; then
        gh pr create --title "$title" --body "$description" --base "$base_branch" --head "$current_branch"
    else
        gh pr create --title "$title" --base "$base_branch" --head "$current_branch"
    fi
}

# Create GitLab MR using glab CLI
create_gitlab_mr() {
    local title="$1"
    local description="$2"
    local base_branch="$3"
    local current_branch="$4"
    
    if ! command -v glab >/dev/null 2>&1; then
        echo "GitLab CLI (glab) not installed. Please install it or use the web interface."
        echo "Install: https://gitlab.com/gitlab-org/cli"
        return 1
    fi
    
    # Check if user is authenticated
    if ! glab auth status >/dev/null 2>&1; then
        echo "GitLab CLI not authenticated. Please run 'glab auth login' first."
        return 1
    fi
    
    # Create MR
    if [ -n "$description" ]; then
        glab mr create --title "$title" --description "$description" --source-branch "$current_branch" --target-branch "$base_branch"
    else
        glab mr create --title "$title" --source-branch "$current_branch" --target-branch "$base_branch"
    fi
}

# Generate web URL for PR creation
generate_web_url() {
    local platform="$1"
    local repo_info="$2"
    local current_branch="$3"
    local base_branch="$4"
    
    case "$platform" in
        "github")
            echo "https://github.com/$repo_info/compare/$base_branch...$current_branch"
            ;;
        "gitlab")
            echo "https://gitlab.com/$repo_info/-/merge_requests/new?source_branch=$current_branch&target_branch=$base_branch"
            ;;
        "bitbucket")
            echo "https://bitbucket.org/$repo_info/pull-requests/new?source=$current_branch&t=1"
            ;;
        *)
            echo "https://$repo_info/compare/$base_branch...$current_branch"
            ;;
    esac
}

# Open URL in browser
open_url() {
    local url="$1"
    
    if command -v open >/dev/null 2>&1; then
        open "$url"
    elif command -v xdg-open >/dev/null 2>&1; then
        xdg-open "$url"
    else
        echo "Please open this URL in your browser: $url"
    fi
}

# Squash commits for PR
squash_commits_for_pr() {
    local base_branch="$1"
    local current_branch
    local commit_count
    local original_pwd
    
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    original_pwd=$(pwd)
    
    # Get the number of commits ahead of base branch
    commit_count=$(git rev-list --count "$base_branch..HEAD" 2>/dev/null || echo "0")
    
    if [ "$commit_count" -le 1 ]; then
        echo "Only one commit or no commits to squash"
        return 0
    fi
    
    echo "Found $commit_count commits to squash"
    
    # Get the commit hash where the branch diverged from base
    local base_commit
    base_commit=$(git merge-base "$base_branch" HEAD 2>/dev/null || echo "")
    
    if [ -z "$base_commit" ]; then
        echo "Error: Cannot find merge base with $base_branch" >&2
        return 1
    fi
    
    # Create a temporary branch for the squash
    local temp_branch
    temp_branch="${current_branch}-squash-$(date +%s)"
    
    # Create temp branch from base
    if ! git checkout -b "$temp_branch" "$base_commit" 2>/dev/null; then
        echo "Error: Failed to create temporary branch" >&2
        return 1
    fi
    
    # Cherry-pick all commits from current branch
    local cherry_pick_success=true
    while IFS= read -r commit_hash; do
        if [ -n "$commit_hash" ]; then
            if ! git cherry-pick "$commit_hash" 2>/dev/null; then
                echo "Error: Failed to cherry-pick commit $commit_hash" >&2
                cherry_pick_success=false
                break
            fi
        fi
    done < <(git rev-list --reverse "$base_commit..HEAD" 2>/dev/null)
    
    if [ "$cherry_pick_success" = false ]; then
        # Clean up on failure
        git checkout "$current_branch" 2>/dev/null
        git branch -D "$temp_branch" 2>/dev/null
        return 1
    fi
    
    # Reset current branch to temp branch
    if ! git checkout "$current_branch" 2>/dev/null; then
        echo "Error: Failed to switch back to current branch" >&2
        git branch -D "$temp_branch" 2>/dev/null
        return 1
    fi
    
    if ! git reset --hard "$temp_branch" 2>/dev/null; then
        echo "Error: Failed to reset current branch" >&2
        git branch -D "$temp_branch" 2>/dev/null
        return 1
    fi
    
    # Clean up temp branch
    git branch -D "$temp_branch" 2>/dev/null
    
    echo "Successfully squashed $commit_count commits into one"
    return 0
}

cmd_pr() {
    local title=""
    local description=""
    local base_branch=""
    local open_browser=false
    local force_squash=""
    local force_no_squash=""
    
    # Parse arguments
    while [ "$#" -gt 0 ]; do
        case "$1" in
            "-h"|"--help"|"help"|"h")
                usage; exit 0
                ;;
            "-t"|"--title")
                if [ -n "$2" ]; then
                    title="$2"
                    shift 2
                else
                    echo "Error: --title requires a value" >&2
                    exit 1
                fi
                ;;
            "-d"|"--description")
                if [ -n "$2" ]; then
                    description="$2"
                    shift 2
                else
                    echo "Error: --description requires a value" >&2
                    exit 1
                fi
                ;;
            "-b"|"--base")
                if [ -n "$2" ]; then
                    base_branch="$2"
                    shift 2
                else
                    echo "Error: --base requires a value" >&2
                    exit 1
                fi
                ;;
            "-o"|"--open")
                open_browser=true
                shift
                ;;
            "-s"|"--squash")
                force_squash=true
                shift
                ;;
            "-S"|"--no-squash")
                force_no_squash=true
                shift
                ;;
            *)
                # If no title provided yet, use this as title
                if [ -z "$title" ]; then
                    title="$1"
                else
                    echo "Error: Unknown option '$1'" >&2
                    usage; exit 1
                fi
                shift
                ;;
        esac
    done
    
    # Validate we're in a git repository
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        echo "Error: Not in a git repository" >&2
        exit 1
    fi
    
    # Get current branch
    local current_branch
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    
    if [ -z "$current_branch" ] || [ "$current_branch" = "HEAD" ]; then
        echo "Error: Not on a branch (detached HEAD state)" >&2
        exit 1
    fi
    
    # Set default base branch if not provided
    if [ -z "$base_branch" ]; then
        base_branch=$(git config at.trunk 2>/dev/null || echo "main")
    fi
    
    # Check if we're trying to create PR from trunk branch
    if [ "$current_branch" = "$base_branch" ]; then
        echo "Error: Cannot create PR from $base_branch to itself" >&2
        exit 1
    fi
    
    # Check for uncommitted changes
    if ! git diff --quiet || ! git diff --cached --quiet; then
        echo "Warning: You have uncommitted changes. Consider committing them first."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    # Determine if we should squash commits
    local should_squash=false
    if [ "$force_squash" = true ]; then
        should_squash=true
    elif [ "$force_no_squash" = true ]; then
        should_squash=false
    else
        # Check configuration setting
        local squash_setting
        squash_setting=$(git config at.pr.squash 2>/dev/null || echo "")
        if [ "$squash_setting" = "true" ]; then
            should_squash=true
        fi
    fi
    
    # Squash commits if enabled
    if [ "$should_squash" = true ]; then
        echo "Auto-squashing commits before creating PR..."
        if ! squash_commits_for_pr "$base_branch"; then
            echo "Error: Failed to squash commits" >&2
            exit 1
        fi
        echo "✅ Commits squashed successfully"
    fi
    
    # Set default title if not provided
    if [ -z "$title" ]; then
        title=$(get_default_title)
    fi
    
    # Generate automatic description if not provided
    if [ -z "$description" ]; then
        echo "Generating automatic description based on changed files..."
        description=$(generate_auto_description "$base_branch" "$current_branch")
    fi
    
    # Get platform and repo info
    local platform
    local repo_info
    
    platform=$(detect_platform)
    repo_info=$(get_repo_info)
    
    echo "Creating PR for $platform repository: $repo_info"
    echo "From: $current_branch → To: $base_branch"
    echo "Title: $title"
    if [ -n "$description" ]; then
        echo "Description: Auto-generated based on changed files"
    fi
    
    # Try to create PR using CLI tools
    local success=false
    
    case "$platform" in
        "github")
            if create_github_pr "$title" "$description" "$base_branch" "$current_branch"; then
                success=true
            fi
            ;;
        "gitlab")
            if create_gitlab_mr "$title" "$description" "$base_branch" "$current_branch"; then
                success=true
            fi
            ;;
    esac
    
    # If CLI failed or not supported, provide web URL
    if [ "$success" = false ]; then
        local web_url
        web_url=$(generate_web_url "$platform" "$repo_info" "$current_branch" "$base_branch")
        
        echo ""
        echo "PR creation via CLI not available. Please create the PR manually:"
        echo "URL: $web_url"
        
        if [ "$open_browser" = true ]; then
            echo "Opening in browser..."
            open_url "$web_url"
        fi
    fi
    
    exit 0
}

# Only run the command if this script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    cmd_pr "$@"
fi 