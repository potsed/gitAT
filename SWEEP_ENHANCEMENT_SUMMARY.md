# GitAT Sweep Command Enhancement

## 🎉 Major Enhancement: Remote Tracking as Default

### **What Changed:**

The `git @ sweep` command now **by default** checks for both:

- ✅ **Locally merged branches** (existing functionality)
- ✅ **Remote-deleted branches** (new default behavior)

### **New Default Behavior:**

```bash
# Default: checks both local merges AND remote deletions
git @ sweep

# Explicitly check both (same as default)
git @ sweep --remote

# Only check locally merged branches
git @ sweep --local-only
```

### **Why This Makes Sense:**

1. **Common Workflow**: Most teams use GitHub/GitLab with "Delete branch after merge"
2. **Team Collaboration**: Branches get deleted on remote after PRs are merged
3. **Cleaner Workspace**: Keeps local branches in sync with remote automatically
4. **Less Manual Work**: No need to remember to add `--remote` flag

### **Safety Features:**

- 🛡️ **Only affects tracking branches** - won't touch local-only branches
- 🛡️ **Preserves important branches** - master, main, dev, etc.
- 🛡️ **Recent activity detection** - flags branches with commits in last 30 days
- 🛡️ **Interactive confirmation** - asks before deleting recent branches
- 🛡️ **Dry run mode** - preview what would be deleted

### **Complete Usage Examples:**

```bash
# Default: clean up merged + remote-deleted branches
git @ sweep

# Only locally merged branches (more conservative)
git @ sweep --local-only

# Force cleanup (including squash-merged)
git @ sweep --force

# Force cleanup with local-only
git @ sweep --force --local-only

# Preview what would be deleted
git @ sweep --dry-run

# Preview local-only cleanup
git @ sweep --local-only --dry-run
```

### **Perfect For:**

- ✅ **Regular cleanup** after team collaboration
- ✅ **GitHub/GitLab web merges** with auto-delete enabled
- ✅ **Post-PR cleanup** when branches are deleted on remote
- ✅ **Keeping local workspace clean** and in sync

### **When to Use --local-only:**

- 🔧 **Conservative cleanup** - when you want to be extra careful
- 🔧 **Offline work** - when you can't check remote
- 🔧 **Specific scenarios** - when you only want local merge detection

### **Migration:**

- **No breaking changes** - existing scripts continue to work
- **More comprehensive** - catches more cleanup opportunities by default
- **Safer** - better detection and confirmation for active branches

This enhancement makes `git @ sweep` much more powerful and user-friendly for modern Git workflows! 🚀
