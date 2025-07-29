Feature: Pull Request Creation
  As a developer using GitAT
  I want to create pull requests or merge requests
  So that I can easily submit my changes for review

  Background:
    Given I am in a git repository
    And I have a remote origin configured
    And I am on a feature branch

  Scenario: Create PR with default title
    When I run "git @ pr"
    Then the command should detect the platform automatically
    And it should use the last commit message as the PR title
    And it should target the configured trunk branch
    And it should create a PR using the appropriate CLI tool
    And it should show success message

  Scenario: Create PR with custom title
    When I run "git @ pr 'Add user authentication feature'"
    Then the command should use "Add user authentication feature" as the title
    And it should create a PR with the custom title

  Scenario: Create PR with description
    When I run "git @ pr -d 'This PR adds user authentication with JWT tokens'"
    Then the command should include the description in the PR
    And it should create a PR with both title and description

  Scenario: Create PR targeting specific base branch
    When I run "git @ pr -b develop"
    Then the command should target the develop branch
    And it should create a PR from current branch to develop

  Scenario: Create PR and open in browser
    When I run "git @ pr -o"
    Then the command should create the PR
    And it should open the PR URL in the browser

  Scenario: Create PR with all options
    When I run "git @ pr 'Feature: Add OAuth support' -d 'Implements OAuth 2.0 authentication' -b main -o"
    Then the command should use the custom title
    And it should include the description
    And it should target the main branch
    And it should open in browser

  Scenario: Show help for PR command
    When I run "git @ pr -h"
    Then it should display the usage information
    And it should show all available options
    And it should show examples

  Scenario: Show help with --help
    When I run "git @ pr --help"
    Then it should display the same usage information as -h

  Scenario: Show help with help
    When I run "git @ pr help"
    Then it should display the same usage information as -h

  Scenario: Error when not in git repository
    Given I am not in a git repository
    When I run "git @ pr"
    Then it should show "Error: Not in a git repository"

  Scenario: Error when no remote origin
    Given I have no remote origin configured
    When I run "git @ pr"
    Then it should show "Error: No remote origin configured"

  Scenario: Error when on trunk branch
    Given I am on the trunk branch
    When I run "git @ pr"
    Then it should show "Error: Cannot create PR from trunk to itself"

  Scenario: Error when in detached HEAD state
    Given I am in a detached HEAD state
    When I run "git @ pr"
    Then it should show "Error: Not on a branch (detached HEAD state)"

  Scenario: Warning for uncommitted changes
    Given I have uncommitted changes
    When I run "git @ pr"
    Then it should show a warning about uncommitted changes
    And it should ask for confirmation to continue

  Scenario: Continue with uncommitted changes after confirmation
    Given I have uncommitted changes
    When I run "git @ pr" and confirm with "y"
    Then it should proceed to create the PR

  Scenario: Abort with uncommitted changes after rejection
    Given I have uncommitted changes
    When I run "git @ pr" and reject with "n"
    Then it should abort the operation

  Scenario: GitHub platform detection
    Given my remote origin is "https://github.com/user/repo.git"
    When I run "git @ pr"
    Then it should detect GitHub as the platform
    And it should try to use GitHub CLI (gh)

  Scenario: GitHub CLI not installed
    Given my remote origin is "https://github.com/user/repo.git"
    And GitHub CLI is not installed
    When I run "git @ pr"
    Then it should show GitHub CLI installation instructions
    And it should provide a web URL for manual PR creation

  Scenario: GitHub CLI not authenticated
    Given my remote origin is "https://github.com/user/repo.git"
    And GitHub CLI is installed but not authenticated
    When I run "git @ pr"
    Then it should show authentication instructions
    And it should provide a web URL for manual PR creation

  Scenario: GitLab platform detection
    Given my remote origin is "https://gitlab.com/user/repo.git"
    When I run "git @ pr"
    Then it should detect GitLab as the platform
    And it should try to use GitLab CLI (glab)

  Scenario: GitLab CLI not installed
    Given my remote origin is "https://gitlab.com/user/repo.git"
    And GitLab CLI is not installed
    When I run "git @ pr"
    Then it should show GitLab CLI installation instructions
    And it should provide a web URL for manual PR creation

  Scenario: GitLab CLI not authenticated
    Given my remote origin is "https://gitlab.com/user/repo.git"
    And GitLab CLI is installed but not authenticated
    When I run "git @ pr"
    Then it should show authentication instructions
    And it should provide a web URL for manual PR creation

  Scenario: Bitbucket platform detection
    Given my remote origin is "https://bitbucket.org/user/repo.git"
    When I run "git @ pr"
    Then it should detect Bitbucket as the platform
    And it should provide a web URL for manual PR creation

  Scenario: Generic platform detection
    Given my remote origin is "https://custom-git-server.com/user/repo.git"
    When I run "git @ pr"
    Then it should detect generic as the platform
    And it should provide a generic web URL

  Scenario: SSH URL handling for GitHub
    Given my remote origin is "git@github.com:user/repo.git"
    When I run "git @ pr"
    Then it should correctly extract "user/repo" from the SSH URL
    And it should detect GitHub as the platform

  Scenario: SSH URL handling for GitLab
    Given my remote origin is "git@gitlab.com:user/repo.git"
    When I run "git @ pr"
    Then it should correctly extract "user/repo" from the SSH URL
    And it should detect GitLab as the platform

  Scenario: SSH URL handling for Bitbucket
    Given my remote origin is "git@bitbucket.org:user/repo.git"
    When I run "git @ pr"
    Then it should correctly extract "user/repo" from the SSH URL
    And it should detect Bitbucket as the platform

  Scenario: Default title from last commit
    Given my last commit message is "Add user authentication"
    When I run "git @ pr"
    Then it should use "Add user authentication" as the default title

  Scenario: Default title when no commits
    Given I have no commits
    When I run "git @ pr"
    Then it should use "Update from current-branch" as the default title

  Scenario: Error for invalid title option
    When I run "git @ pr -t"
    Then it should show "Error: --title requires a value"

  Scenario: Error for invalid description option
    When I run "git @ pr -d"
    Then it should show "Error: --description requires a value"

  Scenario: Error for invalid base option
    When I run "git @ pr -b"
    Then it should show "Error: --base requires a value"

  Scenario: Error for unknown option
    When I run "git @ pr --unknown"
    Then it should show "Error: Unknown option '--unknown'"

  Scenario: Web URL generation for GitHub
    Given my remote origin is "https://github.com/user/repo.git"
    And I am on branch "feature/auth"
    And the trunk branch is "main"
    When I run "git @ pr"
    Then it should generate URL "https://github.com/user/repo/compare/main...feature/auth"

  Scenario: Web URL generation for GitLab
    Given my remote origin is "https://gitlab.com/user/repo.git"
    And I am on branch "feature/auth"
    And the trunk branch is "main"
    When I run "git @ pr"
    Then it should generate URL "https://gitlab.com/user/repo/-/merge_requests/new?source_branch=feature/auth&target_branch=main"

  Scenario: Web URL generation for Bitbucket
    Given my remote origin is "https://bitbucket.org/user/repo.git"
    And I am on branch "feature/auth"
    When I run "git @ pr"
    Then it should generate URL "https://bitbucket.org/user/repo/pull-requests/new?source=feature/auth&t=1"

  Scenario: Browser opening on macOS
    Given I am on macOS
    And the platform is GitHub
    When I run "git @ pr -o"
    Then it should use "open" command to open the URL

  Scenario: Browser opening on Linux
    Given I am on Linux
    And the platform is GitHub
    When I run "git @ pr -o"
    Then it should use "xdg-open" command to open the URL

  Scenario: Browser opening fallback
    Given I am on an unsupported system
    And the platform is GitHub
    When I run "git @ pr -o"
    Then it should show "Please open this URL in your browser: [URL]"

  Scenario: Use configured trunk branch as default
    Given the configured trunk branch is "develop"
    When I run "git @ pr"
    Then it should use "develop" as the target branch

  Scenario: Fallback to main when no trunk configured
    Given no trunk branch is configured
    When I run "git @ pr"
    Then it should use "main" as the target branch

  Scenario: Successful GitHub PR creation
    Given my remote origin is "https://github.com/user/repo.git"
    And GitHub CLI is installed and authenticated
    When I run "git @ pr 'Test PR'"
    Then it should create a PR successfully
    And it should show success message

  Scenario: Successful GitLab MR creation
    Given my remote origin is "https://gitlab.com/user/repo.git"
    And GitLab CLI is installed and authenticated
    When I run "git @ pr 'Test MR'"
    Then it should create an MR successfully
    And it should show success message

  Scenario: PR with automatic squashing enabled
    Given automatic PR squashing is enabled
    And I have multiple commits on my feature branch
    When I run "git @ pr 'Feature: Add authentication'"
    Then it should show "Auto-squashing commits before creating PR"
    And it should squash the commits
    And it should show "Commits squashed successfully"
    And it should create the PR with squashed commits

  Scenario: PR with automatic squashing disabled
    Given automatic PR squashing is disabled
    And I have multiple commits on my feature branch
    When I run "git @ pr 'Feature: Add authentication'"
    Then it should not show squashing messages
    And it should create the PR with all commits as-is

  Scenario: PR with force squash override
    Given automatic PR squashing is disabled
    And I have multiple commits on my feature branch
    When I run "git @ pr 'Feature: Add authentication' -s"
    Then it should show "Auto-squashing commits before creating PR"
    And it should squash the commits despite the setting being disabled

  Scenario: PR with force no squash override
    Given automatic PR squashing is enabled
    And I have multiple commits on my feature branch
    When I run "git @ pr 'Feature: Add authentication' -S"
    Then it should not show squashing messages
    And it should create the PR with all commits as-is despite the setting being enabled

  Scenario: No squashing when only one commit
    Given automatic PR squashing is enabled
    And I have only one commit on my feature branch
    When I run "git @ pr 'Feature: Add authentication'"
    Then it should show "Only one commit or no commits to squash"
    And it should not perform any squashing

  Scenario: Squash preserves meaningful commit message
    Given automatic PR squashing is enabled
    And I have multiple commits with meaningful messages
    When I run "git @ pr 'Feature: Add authentication'"
    Then it should squash the commits
    And the final commit should have a meaningful message from the commits 