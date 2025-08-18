# Binary Naming Guidelines

## Overview

GitAT uses `git-@` as its binary name to function as a Git extension. This allows users to run commands like `git @ status` or `git @ version`.

## ⚠️ Important: Deprecated Binary Name

The binary name `gitat` is **deprecated** and should never be used. Always use `git-@`.

## Safeguards in Place

To prevent accidental use of the deprecated binary name, several safeguards are implemented:

### 1. Makefile Safeguards

- **Binary Name Validation**: The Makefile validates that `BINARY_NAME=git-@` before building
- **Deprecated Binary Prevention**: A specific `gitat` target that shows an error message
- **Automatic Cleanup**: The `clean` target removes both `git-@` and any deprecated `gitat` binaries
- **Build Validation**: All build targets run `check-binary-name` first

```bash
# ✅ Correct usage
make build-local    # Builds git-@
make build          # Builds git-@ in build/ directory

# ❌ This will show an error
make gitat          # Shows deprecation error
```

### 2. Git Hooks

Pre-commit hooks prevent committing deprecated binaries:

```bash
# Setup hooks (run once)
./scripts/setup-hooks.sh

# Or as part of development setup
make dev-setup
```

The pre-commit hook will:
- Prevent committing any `gitat` binary
- Prevent staging `gitat` in git
- Check for build references to `gitat`

### 3. CI/CD Validation

The `scripts/validate-binary-name.sh` script provides comprehensive validation:

```bash
# Run validation
./scripts/validate-binary-name.sh

# Or through make
make validate
```

This checks:
- No deprecated `gitat` binary exists
- Makefile uses correct binary name
- `.gitignore` includes deprecated binary
- No problematic references in build scripts

### 4. .gitignore Protection

The `.gitignore` file includes:
```
# Build artifacts
gitat
gitat.exe
gitat-*

# Correct binary name (should be tracked)
# git-@
```

## Development Workflow

### Initial Setup

```bash
# Clone repository
git clone <repository-url>
cd gitAT

# Setup development environment (includes hooks)
make dev-setup
```

### Building

```bash
# Build for local development
make build-local

# Build for distribution
make build

# Clean all artifacts
make clean
```

### Validation

```bash
# Validate binary naming
make validate

# Run all checks
make test validate lint
```

## Troubleshooting

### If you accidentally create `gitat`

```bash
# Clean and rebuild
make clean
make build-local
```

### If git hooks are not working

```bash
# Reinstall hooks
./scripts/setup-hooks.sh

# Or bypass temporarily (not recommended)
git commit --no-verify
```

### If CI/CD fails on binary naming

1. Check that you're building `git-@` not `gitat`
2. Run `make validate` locally
3. Ensure no `gitat` binary is committed
4. Update any scripts that reference `gitat`

## Why `git-@` Instead of `gitat`?

The `git-@` binary name allows GitAT to function as a Git extension:

- **Git Extension**: Git automatically finds `git-@` and allows `git @` commands
- **Namespace Separation**: Clear distinction between Git and GitAT commands
- **User Experience**: Natural syntax like `git @ version` or `git @ status`
- **Tool Integration**: Works seamlessly with existing Git workflows

## Migration from `gitat`

If you have existing references to `gitat`:

1. **Update build scripts**: Change `gitat` to `git-@`
2. **Update documentation**: Reference `git-@` instead of `gitat`
3. **Update installation**: Install as `git-@` in PATH
4. **Update aliases**: Change shell aliases to use `git-@`

Example migration:
```bash
# Old (deprecated)
./gitat version
alias gat="gitat"

# New (correct)
./git-@ version
alias gat="git-@"
# Or better yet
alias gat="git @"
```