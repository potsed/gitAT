# Git Extension Mechanism

## Overview

GitAT uses Git's built-in extension mechanism to provide the `git @` command. This approach is cleaner and more standard than using Git aliases.

## How Git Extensions Work

Git automatically recognizes any executable named `git-<command>` as a Git extension. When you run `git <command>`, Git will:

1. Look for an executable named `git-<command>` in your PATH
2. Execute it with all the arguments you provided
3. Pass through the exit code

## GitAT Implementation

GitAT provides a `git-@` executable that:

- Is a bash script that handles command routing
- Sources individual command scripts from `git_at_cmds/`
- Provides comprehensive help and error handling
- Integrates with the Go binary for TUI functionality

## Installation

The installation process copies the necessary files to `/usr/local/bin/`:

```bash
# Install the Go binary (for TUI functionality)
sudo cp gitat /usr/local/bin/

# Install the Git extension
sudo cp git-@ /usr/local/bin/
sudo chmod +x /usr/local/bin/git-@

# Install command scripts
sudo cp -r git_at_cmds /usr/local/bin/
```

## Usage

Once installed, you can use GitAT as a native Git command:

```bash
# Show help
git @ --help

# Show status
git @ info

# Create a feature branch
git @ work feature add-authentication

# Save changes
git @ save "Add user authentication"

# Launch TUI
git @ info
```

## Advantages Over Aliases

1. **No Configuration Required**: No need to set up Git aliases
2. **Standard Git Mechanism**: Uses Git's official extension system
3. **Better Integration**: Works seamlessly with Git's command completion
4. **Cleaner**: No wrapper scripts or complex alias configurations
5. **Portable**: Works on any system where Git is installed

## Troubleshooting

### "No manual entry for git-@"

This is normal! Git extensions don't automatically have man pages. The functionality still works perfectly.

### Command Not Found

Ensure the `git-@` file is:

- In your PATH
- Executable (`chmod +x git-@`)
- Has the correct shebang line

### Missing Command Scripts

Ensure the `git_at_cmds/` directory is copied to the same location as `git-@` and contains all the command scripts.

## Development

When developing GitAT:

1. Make `git-@` executable: `chmod +x git-@`
2. Add the project directory to PATH for testing: `export PATH="/path/to/gitAT:$PATH"`
3. Test commands: `git @ info`, `git @ --help`, etc.
