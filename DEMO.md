# Task CLI Demo Guide

This guide walks you through the basic functionality of the Task CLI application.

## 🎬 Quick Start Demo

### Step 1: Launch the Application
```bash
./task-cli
```

You'll see the main interface with:
- Empty task list (initially)
- Help text at the bottom: "Keys: n=New, e=Edit, d=Delete, t=Toggle, q=Quit, /=Search"

### Step 2: Create Your First Task
1. Press `n` (New task)
2. You'll switch to the form view
3. Fill in the fields:
   - **Title**: "Learn Go programming"
   - **Description**: "Complete Go tutorial and build a project"
   - **Priority**: Select "High" using arrow keys
   - **Tags**: "learning,programming,go"
4. Press `Ctrl+S` to submit or use Tab to navigate to Submit button

### Step 3: Create More Tasks
Repeat the process to create several tasks:

**Task 2:**
- Title: "Buy groceries"
- Description: "Milk, bread, eggs, and vegetables"
- Priority: Medium
- Tags: "shopping,daily"

**Task 3:**
- Title: "Exercise routine"
- Description: "30 minutes cardio and strength training"
- Priority: High
- Tags: "health,fitness"

**Task 4:**
- Title: "Read book"
- Description: "Finish reading 'Clean Code'"
- Priority: Low
- Tags: "reading,learning"

### Step 4: Navigate and Manage Tasks
- Use `↑` and `↓` arrow keys to navigate between tasks
- Press `t` to toggle task status (Todo → In Progress → Completed)
- Notice how the status symbols change: ◯ → ◐ → ●

### Step 5: Edit a Task
1. Select a task with arrow keys
2. Press `e` (Edit)
3. Modify any field
4. Press `Ctrl+S` to save changes

### Step 6: Try Different Themes
1. Press `q` to quit
2. Launch with different themes:
   ```bash
   ./task-cli --theme dark
   ./task-cli --theme light
   ```

## 🎯 Advanced Usage Examples

### Custom Data Directory
```bash
# Use a specific directory for your tasks
./task-cli --data-dir ~/work-tasks

# This creates: ~/work-tasks/tasks.json
```

### Theme Comparison
```bash
# Default theme (balanced)
./task-cli

# High contrast dark theme
./task-cli --theme dark

# Clean light theme  
./task-cli --theme light
```

## 📊 Understanding the Interface

### List View Layout
```
┌─Status─┬─Priority─┬─Title────────────┬─Description─────────┐
│   ◯    │   !!!    │ Learn Go         │ Complete Go tutorial│
│   ◐    │    !!    │ Buy groceries    │ Milk, bread, eggs   │
│   ●    │    !     │ Exercise routine │ 30 minutes cardio   │
└────────┴──────────┴──────────────────┴─────────────────────┘
Keys: n=New, e=Edit, d=Delete, t=Toggle, q=Quit, /=Search
```

### Form View Layout
```
┌─────────────────────────────────────────┐
│ Title: [Learn Go programming___________] │
│                                         │
│ Description:                            │
│ ┌─────────────────────────────────────┐ │
│ │Complete Go tutorial and build a    │ │  
│ │project                             │ │
│ └─────────────────────────────────────┘ │
│                                         │
│ Priority: [High          ▼]            │
│                                         │
│ Tags: [learning,programming,go_______]  │
│                                         │
│       [Submit]    [Cancel]              │
└─────────────────────────────────────────┘
Keys: Ctrl+S=Submit, Escape=Cancel
```

## 🎨 Visual Elements

### Priority Indicators
- `!!!` = High Priority (Red)
- ` !! ` = Medium Priority (Yellow)  
- ` !  ` = Low Priority (Green)

### Status Symbols
- `◯` = Todo (White/Default)
- `◐` = In Progress (Blue)
- `●` = Completed (Green)

## 💡 Pro Tips

### Keyboard Efficiency
1. **Quick Toggle**: Navigate with arrows, press `t` to mark complete
2. **Rapid Entry**: Press `n`, type title, Tab to description, Tab to priority
3. **Mass Management**: Use `t` to quickly complete multiple tasks

### Organization Strategy
1. **Use Meaningful Tags**: `work`, `personal`, `urgent`, `project-name`
2. **Priority Levels**: 
   - High: Deadlines, urgent issues
   - Medium: Regular tasks, planned work  
   - Low: Ideas, someday/maybe items
3. **Status Workflow**: Todo → In Progress → Completed

### Data Management
- Tasks are automatically saved to `~/.task-cli/tasks.json`
- Backups are created automatically in `~/.task-cli/backups/`
- Data persists between sessions

## 🔄 Workflow Examples

### Daily Task Management
```bash
# Morning: Check tasks
./task-cli

# Add new tasks as they come up (press 'n')
# Mark tasks in progress as you work (press 't')  
# Complete tasks throughout the day (press 't' again)
```

### Project-Based Organization
```bash
# Work tasks
./task-cli --data-dir ~/work-tasks

# Personal tasks  
./task-cli --data-dir ~/personal-tasks

# Project-specific tasks
./task-cli --data-dir ~/project-alpha-tasks
```

## 🐛 Troubleshooting

### Common Issues

**Application won't start:**
```bash
# Check if binary exists
ls -la task-cli

# Ensure it's executable
chmod +x task-cli

# Try running with full path
./task-cli
```

**Data not persisting:**
- Check permissions on data directory
- Ensure sufficient disk space
- Verify data directory path exists

**Theme not working:**
```bash
# Check available themes
./task-cli --help

# Valid themes: default, dark, light
./task-cli --theme dark
```

---

**Enjoy using Task CLI! 🚀**

For more advanced usage and development information, see the main [README.md](README.md).