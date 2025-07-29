#!/bin/bash

echo "Testing git @ info..."

# Test if the info.sh file exists
if [ -f "git_at_cmds/info.sh" ]; then
    echo "✅ info.sh file exists"
else
    echo "❌ info.sh file not found"
    exit 1
fi

# Test if the file is executable
if [ -x "git_at_cmds/info.sh" ]; then
    echo "✅ info.sh is executable"
else
    echo "❌ info.sh is not executable"
fi

# Test sourcing the file
echo "Testing source of info.sh..."
if source git_at_cmds/info.sh 2>&1; then
    echo "✅ info.sh sourced successfully"
else
    echo "❌ Failed to source info.sh"
fi

# Test if cmd_info function exists
if type cmd_info >/dev/null 2>&1; then
    echo "✅ cmd_info function exists"
else
    echo "❌ cmd_info function not found"
fi

# Test calling cmd_info directly
echo "Testing cmd_info function..."
if cmd_info 2>&1; then
    echo "✅ cmd_info executed successfully"
else
    echo "❌ cmd_info failed"
fi 