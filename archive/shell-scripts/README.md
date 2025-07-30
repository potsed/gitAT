# Archived Shell Scripts

This directory contains the original shell script implementations of GitAT commands that have been migrated to Go.

## Files Archived

- **All command scripts** (`*.sh`) - Migrated to Go functions in `internal/commands/manager.go`
- **git-@** - The original bash wrapper script, now replaced by the Go binary

## Migration Status

All 25 shell scripts have been successfully migrated to Go and are now available in the main Go implementation.

### Commands Migrated

| Shell Script | Go Implementation | Status |
|--------------|-------------------|--------|
| `work.sh` | `Work()` | ✅ Complete |
| `hotfix.sh` | `Hotfix()` | ✅ Complete |
| `save.sh` | `Save()` | ✅ Complete |
| `squash.sh` | `Squash()` | ✅ Complete |
| `pr.sh` | `PullRequest()` | ✅ Complete |
| `branch.sh` | `Branch()` | ✅ Complete |
| `sweep.sh` | `Sweep()` | ✅ Complete |
| `info.sh` | `Info()` | ✅ Complete |
| `hash.sh` | `Hash()` | ✅ Complete |
| `product.sh` | `Product()` | ✅ Complete |
| `feature.sh` | `Feature()` | ✅ Complete |
| `issue.sh` | `Issue()` | ✅ Complete |
| `version.sh` | `Version()` | ✅ Complete |
| `release.sh` | `Release()` | ✅ Complete |
| `master.sh` | `Master()` | ✅ Complete |
| `root.sh` | `Master()` | ✅ Complete |
| `wip.sh` | `WIP()` | ✅ Complete |
| `changes.sh` | `Changes()` | ✅ Complete |
| `logs.sh` | `Logs()` | ✅ Complete |
| `_label.sh` | `Label()` | ✅ Complete |
| `_id.sh` | `ID()` | ✅ Complete |
| `_path.sh` | `Path()` | ✅ Complete |
| `_trunk.sh` | `Trunk()` | ✅ Complete |
| `ignore.sh` | `Ignore()` | ✅ Complete |
| `initlocal.sh` | `InitLocal()` | ✅ Complete |
| `initremote.sh` | `InitRemote()` | ✅ Complete |
| `_security.sh` | `Security()` | ✅ Complete |
| `_go.sh` | `Go()` | ✅ Complete |

## Go Implementation

The Go implementation provides:

- Better performance and startup time
- Type safety and error handling
- Cross-platform compatibility
- Consistent command-line interface
- Comprehensive help documentation
- Unit and integration tests

## Usage

To use the new Go implementation:

```bash
# Build the Go binary
go build -o gitat cmd/gitat/main.go

# Use any command
./gitat work feature my-feature
./gitat save "My commit message"
./gitat pr
```

## Archive Date

Shell scripts archived on: $(date)

## Notes

- All functionality from the shell scripts has been preserved
- Help text and documentation have been maintained
- Error handling has been improved
- Performance has been enhanced
