#!/bin/bash

echo "Testing improved markdown formatting for PR descriptions..."

# Create a sample markdown description to show the new format
cat << 'EOF'

# 📋 Pull Request Summary

This PR contains changes from branch `feature-auto-description` targeting `main`.

## 📊 Changes Overview

| Metric | Count |
|--------|-------|
| **Total Files** | 4 |
| **Added** | 1 |
| **Modified** | 2 |
| **Deleted** | 1 |

## 📁 File Analysis

### 🔤 File Types

This PR affects the following file types:

- `txt`
- `py`
- `md`
- `js`

### 📂 Directories Affected

Changes span across the following directories:

- `test-dir`
- `src`
- `root`

## 📝 Changed Files

<details>
<summary>📋 Click to view all changed files</summary>

```
➕ test-dir/newfile.py
✏️ test-dir/file1.txt
🗑️ test-dir/file2.sh
✏️ README.md
```
</details>

## 🔄 Commits

This PR includes **2 commits**.

<details>
<summary>📜 Click to view commit history</summary>

```
abc1234 Update files for testing auto-description
def5678 Add test files for auto-description testing
```
</details>

---

*This description was automatically generated based on the changes in this PR.*

EOF

echo ""
echo "✅ The improved markdown formatting includes:"
echo "  📋 Professional header with emojis"
echo "  📊 Clean table format for statistics"
echo "  📁 Organized sections with clear hierarchy"
echo "  📝 Collapsible file lists"
echo "  🔄 Commit history in collapsible sections"
echo "  --- Footer with attribution"
echo ""
echo "This format will render beautifully in GitHub, GitLab, and other platforms!" 