# GitAT Terminal User Interface (TUI)

GitAT now includes a beautiful, interactive terminal user interface built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Bubbles](https://github.com/charmbracelet/bubbles).

## Launching the TUI

To launch the interactive TUI, run:

```bash
git @ --tui
```

## TUI Features

### 🎨 **Beautiful Interface**

- Modern, colorful terminal interface
- Responsive design that adapts to terminal size
- Smooth animations and transitions
- Professional styling with Lip Gloss

### 📊 **Four Main Tabs**

#### 1. **Status Tab**

- Real-time repository status
- Current branch information
- Working directory status
- Quick action shortcuts
- Auto-refreshing data

#### 2. **Commands Tab**

- Interactive list of all GitAT commands
- Command descriptions and shortcuts
- One-click command execution
- Keyboard navigation

#### 3. **Branches Tab**

- List of all local branches
- Current branch highlighting
- Branch switching with Enter key
- Branch deletion with 'd' key
- Real-time branch refresh

#### 4. **Info Tab**

- Comprehensive repository information
- GitAT configuration display
- Git repository details
- Commit statistics
- Remote information

## Navigation

### **Tab Navigation**

- **Tab** / **Shift+Tab**: Switch between tabs
- **1-4**: Jump directly to specific tabs
- **q**: Quit the application

### **General Navigation**

- **↑/↓**: Navigate lists and menus
- **Enter**: Execute selected action
- **Ctrl+C**: Quit the application

### **Branch Management**

- **Enter**: Switch to selected branch
- **d**: Delete selected branch
- **r**: Refresh branch list

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Next tab |
| `Shift+Tab` | Previous tab |
| `1-4` | Jump to tab 1-4 |
| `↑/↓` | Navigate |
| `Enter` | Execute/Select |
| `q` | Quit |
| `Ctrl+C` | Quit |
| `r` | Refresh (in branches tab) |
| `d` | Delete (in branches tab) |

## Features in Detail

### **Status Tab**

- **Current Branch**: Shows the active branch
- **Repository Path**: Full path to the Git repository
- **Status**: Clean/Modified working directory
- **Last Updated**: Timestamp of last refresh
- **Quick Actions**: Common command shortcuts

### **Commands Tab**

- **Work Branch**: Create new work branches
- **Hotfix Branch**: Create hotfix branches
- **Save Changes**: Commit changes with validation
- **Squash Commits**: Combine commits
- **Pull Request**: Create PRs with auto-description
- **Branch Management**: Manage branches
- **Sweep Branches**: Clean up branches
- **Info Display**: Show repository info
- **Hash Display**: Show commit relationships
- **Product Config**: Set product name
- **Version Management**: Handle versions
- **Release Creation**: Create releases

### **Branches Tab**

- **Branch List**: All local branches
- **Current Indicator**: Highlights active branch
- **Status Display**: Shows branch status
- **Switch Function**: Change branches instantly
- **Delete Function**: Remove branches safely
- **Auto-refresh**: Updates when switching tabs

### **Info Tab**

- **Configuration**: GitAT settings
- **Git Repository**: Repository details
- **Statistics**: Commit counts and history
- **Remote Info**: Remote repository details

## Technical Details

### **Built With**

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling

### **Architecture**

- **Model-View-Update**: Elm-inspired architecture
- **Component-based**: Modular UI components
- **Message-driven**: Event-based communication
- **Responsive**: Adapts to terminal size

### **Performance**

- **Efficient**: Minimal resource usage
- **Fast**: Instant response to user input
- **Smooth**: 60fps animations
- **Lightweight**: Small binary size

## Customization

The TUI uses a consistent color scheme and styling that can be easily customized:

- **Primary Color**: `#7D56F4` (Purple)
- **Background**: Terminal default
- **Text**: `#FAFAFA` (Light gray)
- **Secondary**: `#666666` (Gray)
- **Error**: `#FF0000` (Red)

## Future Enhancements

Planned features for the TUI:

1. **Interactive Forms**: Command input forms
2. **Real-time Updates**: Live status updates
3. **Custom Themes**: User-defined color schemes
4. **Keyboard Shortcuts**: Customizable shortcuts
5. **Plugin Support**: Extensible command system
6. **Search Functionality**: Find commands and branches
7. **History**: Command execution history
8. **Notifications**: Success/error notifications

## Troubleshooting

### **TUI Not Starting**

- Ensure terminal supports colors
- Check terminal size (minimum 80x24)
- Verify Go installation

### **Display Issues**

- Try resizing terminal window
- Check terminal color support
- Restart the application

### **Performance Issues**

- Close other terminal applications
- Reduce terminal window size
- Check system resources

## Examples

### **Basic Usage**

```bash
# Launch TUI
git @ --tui

# Navigate to Commands tab
# Press '2' or use Tab

# Execute a command
# Select "Create Work Branch" and press Enter

# Switch branches
# Navigate to Branches tab, select branch, press Enter
```

### **Workflow Example**

```bash
# 1. Launch TUI
git @ --tui

# 2. Create work branch
# Tab to Commands → Select "Create Work Branch" → Enter

# 3. Check status
# Tab to Status → View current branch and status

# 4. Switch branches
# Tab to Branches → Select target branch → Enter

# 5. Save changes
# Tab to Commands → Select "Save Changes" → Enter
```

The TUI provides an intuitive, visual way to interact with GitAT, making Git workflow management more accessible and enjoyable!
