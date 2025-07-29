# GitAT - Git Workflow Management Tool

A comprehensive Git workflow management tool that provides a set of commands to streamline development processes, ensure code quality, and integrate with modern Git hosting platforms.

## Overview

GitAT (`git @`) is a powerful command-line tool that extends Git with workflow management capabilities. It provides commands for branch management, commit organization, pull request creation, and more, all following best practices and conventions.

## Features

### 🚀 **Core Workflow Commands**

- `git @ work <type> <description>` - Create work branches following Conventional Commits
- `git @ hotfix <description>` - Create hotfix branches for urgent fixes
- `git @ save "message"` - Securely save changes with validation
- `git @ squash [branch]` - Squash commits with auto-detection of parent branch
- `git @ pr [options]` - Create Pull Requests with auto-description generation

### 🌿 **Branch Management**

- `git @ branch` - Manage working branch configuration
- `git @ sweep` - Clean up local branches (merged + remote-deleted)
- `git @ master` / `git @ root` - Switch to trunk branches
- `git @ wip` - Work in progress management

### 📊 **Information & Status**

- `git @ info` - Comprehensive status report from all commands
- `git @ hash` - Detailed branch status and commit relationships
- `git @ changes` - View uncommitted changes
- `git @ logs` - View commit history

### 🏷️ **Version Management**

- `git @ version` - Semantic versioning management
- `git @ release` - Create releases with proper tagging
- `git @ product` - Product name configuration

### 🔧 **Utilities**

- `git @ _label` - Generate commit labels
- `git @ _id` - Generate unique project identifiers
- `git @ _path` - Get repository path
- `git @ _trunk` - Manage trunk branch configuration

## Installation

```bash
# Clone the repository
git clone https://github.com/potsed/gitAT.git
cd gitAT

# Run the installation script
./install.sh
```

## Quick Start

1. **Initialize in a Git repository:**

   ```bash
   git @ product "MyProject"
   git @ _trunk master
   ```

2. **Create a feature branch:**

   ```bash
   git @ work feature add-user-authentication
   ```

3. **Make changes and save:**

   ```bash
   git @ save "Add user authentication system"
   ```

4. **Create a pull request:**

   ```bash
   git @ pr
   ```

## Configuration

GitAT uses Git configuration to store project settings:

- `at.product` - Product name
- `at.feature` - Current feature name
- `at.task` - Current task/issue ID
- `at.branch` - Working branch
- `at.trunk` - Trunk branch (master/main)
- `at.version` - Current version
- `at.wip` - Work in progress branch

## Conventional Commits

GitAT follows the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - New features (MINOR version)
- `fix:` - Bug fixes (PATCH version)
- `docs:` - Documentation changes
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance tasks
- `hotfix:` - Urgent fixes

## Examples

### Creating Different Work Types

```bash
git @ work feature add-payment-gateway
git @ work bugfix fix-login-validation
git @ work docs update-api-documentation
git @ work refactor optimize-database-queries
```

### Managing Branches

```bash
git @ branch -c                    # Show current branch
git @ branch --feature             # List feature branches
git @ branch --all-types           # List all branch types
```

### Creating Pull Requests

```bash
git @ pr                           # Create PR with auto-description
git @ pr -s                        # Create PR with auto-squashing
git @ pr --title "Custom Title"    # Create PR with custom title
```

### Cleaning Up Branches

```bash
git @ sweep                        # Clean merged + remote-deleted branches
git @ sweep --local-only           # Clean only locally merged branches
git @ sweep --force                # Force clean (including squash-merged)
git @ sweep --dry-run              # Preview what would be deleted
```

## Testing

Run the comprehensive test suite:

```bash
# Run all tests
./tests/run-all.sh

# Run specific test types
./tests/unit/run-unit-tests.sh
./tests/commands/run-command-tests.sh
./tests/integration/run-integration-tests.sh
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git @ work feature your-feature-name`
3. Make your changes
4. Save your work: `git @ save "Add your feature"`
5. Create a pull request: `git @ pr`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/potsed/gitAT).
