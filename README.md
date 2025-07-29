# git @

A secure Git plugin to help with git workflows from the command line and some utilities.

## Security Features

This plugin implements comprehensive security measures to protect against common vulnerabilities:

- **Command Injection Protection**: All user inputs are validated and sanitized
- **Path Traversal Protection**: Paths are validated to prevent directory traversal attacks
- **Input Validation**: Strict validation of all user inputs with length limits
- **Permission Checking**: Verification of user permissions before operations
- **Secure Logging**: All security events are logged for audit purposes
- **Error Handling**: Secure error handling that doesn't expose sensitive information

## Usage

Once installed, the command `git @` is available from the command line.

### Basic Commands

Running `git @` from within a repository will display the list of sub-commands available:

```bash
git @
```

### Available Commands

- **`git @ product`** - Set or get the current project name
- **`git @ feature`** - Set or get the current feature name  
- **`git @ issue`** - Set or get the current issue/task number
- **`git @ branch`** - Manage working branches
- **`git @ save`** - Save current changes with validation
- **`git @ version`** - Manage semantic versioning
- **`git @ release`** - Create releases with proper tagging
- **`git @ work`** - Switch to work branch and stash changes
- **`git @ wip`** - Manage work-in-progress branches
- **`git @ master`** - Switch to master/trunk branch
- **`git @ root`** - Switch to root/trunk branch
- **`git @ logs`** - Show recent commit logs
- **`git @ changes`** - Show changed files
- **`git @ ignore`** - Manage .gitignore entries
- **`git @ info`** - Display project information
- **`git @ hash`** - Show commit hashes and branch info
- **`git @ sweep`** - Clean up merged branches
- **`git @ squash`** - Squash commits on a branch

### Example Workflow

```bash
# Set up project
git @ product my-awesome-project
git @ feature user-authentication
git @ issue 123

# Create and switch to work branch
git @ work

# Make changes and save
git @ save "Add login form validation"

# Create release
git @ release patch

# Switch back to master
git @ master
```

### Getting Help

For help on any command:

```bash
git @ <command> help
# or
git @ <command> -h
```

## Installing

### Automated Installation (Recommended)

The easiest way to install GitAT is using the automated install script:

```bash
# Clone the repository
git clone https://github.com/potsed/gitAT.git
cd gitAT

# Run the install script
./install.sh
```

The install script will:

- ✅ Check dependencies (Git)
- ✅ Set proper permissions
- ✅ Auto-detect your shell and installation method
- ✅ Install GitAT to your system
- ✅ Verify the installation
- ✅ Run basic tests

**Installation Options:**

```bash
# Auto-detect best method (recommended)
./install.sh

# Specify installation method
./install.sh --method profile    # Add to shell profile
./install.sh --method link       # Create symbolic link
./install.sh --method copy       # Copy to directory

# Copy to specific directory
./install.sh --method copy --directory ~/.local/bin

# Backup existing installation
./install.sh --backup
```

### Manual Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/potsed/gitAT.git
   cd gitAT
   ```

2. **Make the script executable:**

   ```bash
   chmod +x git-@
   ```

3. **Add to your PATH:**

   **Option A: Add to your shell profile (recommended)**

   For bash/zsh (add to `~/.bashrc`, `~/.zshrc`, or `~/.bash_profile`):

   ```bash
   export PATH="$PATH:$(pwd)"
   ```

   For fish shell (add to `~/.config/fish/config.fish`):

   ```bash
   set -gx PATH $PATH (pwd)
   ```

   **Option B: Create a symbolic link**

   ```bash
   sudo ln -s "$(pwd)/git-@" /usr/local/bin/git-@
   ```

   **Option C: Copy to a directory in your PATH**

   ```bash
   cp git-@ ~/.local/bin/
   # or
   sudo cp git-@ /usr/local/bin/
   ```

4. **Reload your shell or restart your terminal**

5. **Verify installation:**

   ```bash
   git @
   ```

### Installation Verification

After installation, you should be able to run:

```bash
git @
```

This should display the available sub-commands and version information.

### Troubleshooting

- **"Command not found"**: Make sure the directory containing `git-@` is in your PATH
- **"Permission denied"**: Run `chmod +x git-@` to make the script executable
- **"Not a git repository"**: Run `git @` from within a git repository

## Testing

Run the comprehensive test suite to verify security features:

```bash
# Run all tests
bash tests/run-all-tests.sh

# Run security tests only
bash tests/security.test.sh

# Run integration tests
bash tests/integration.test.sh

# Run unit tests
bash tests/unit.test.sh

# Run shellcheck (requires shellcheck to be installed)
shellcheck git-@ git_at_cmds/*.sh
```

### Installing ShellCheck

For macOS:

```bash
brew install shellcheck
```

For Ubuntu/Debian:

```bash
sudo apt-get install shellcheck
```

For CentOS/RHEL:

```bash
sudo yum install epel-release
sudo yum install shellcheck
```

For Windows (with Chocolatey):

```bash
choco install shellcheck
```

## Security Considerations

- All inputs are validated against dangerous characters and patterns
- Path operations are restricted to the repository root
- Configuration values are validated before storage
- Security events are logged to `~/.gitat_security.log`
- No operations require elevated privileges
- Error messages don't expose sensitive information

## Development

### Running Tests

The project includes a comprehensive test suite to ensure security and functionality:

```bash
# Run all tests
bash tests/run-all-tests.sh

# Run specific test suites
bash tests/unit.test.sh      # Unit tests for security functions
bash tests/security.test.sh  # Security vulnerability tests
bash tests/integration.test.sh # End-to-end integration tests
shellcheck git-@ git_at_cmds/*.sh # Static analysis with ShellCheck
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git @ work`
3. Make your changes
4. Run tests: `bash tests/run-all-tests.sh`
5. Commit your changes: `git @ save "Description of changes"`
6. Push to your fork and submit a pull request

### Security Reporting

If you discover a security vulnerability, please report it privately to the maintainers. Do not create public issues for security concerns.
