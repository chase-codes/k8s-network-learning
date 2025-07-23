# NetLab Style Guide

This document defines the visual and interaction standards for NetLab's TUI components to ensure consistency and professional quality across all modules.

## 🎨 Color Palette

### Primary Colors
- **Primary Cyan** (`#00D4FF`) - Main brand color, used for highlights and active states
- **Secondary Purple** (`#7C3AED`) - Learning theme, used for section headers
- **Accent Amber** (`#F59E0B`) - Emphasis and call-to-action elements

### Status Colors
- **Success Green** (`#10B981`) - Completed items, success states
- **Warning Amber** (`#F59E0B`) - Caution, work-in-progress items
- **Error Red** (`#EF4444`) - Errors, failures, critical items
- **Info Blue** (`#3B82F6`) - Information, planned items

### Neutral Palette
- **Text Light** (`#F8FAFC`) - Primary text content
- **Text Muted** (`#94A3B8`) - Secondary text, descriptions
- **Text Dim** (`#64748B`) - Tertiary text, help text
- **Background Dark** (`#0F172A`) - Main background
- **Card Background** (`#1E293B`) - Component backgrounds
- **Border Gray** (`#334155`) - Borders and dividers

## 🔤 Typography Scale

### Headers
```go
// H1 - Module titles, main headings
styles.H1.Render("Module Title")

// H2 - Section headers, subsections  
styles.H2.Render("Section Header")

// H3 - Sub-section headers
styles.H3.Render("Sub-section")
```

### Body Text
```go
// Primary content
styles.Body.Render("Main content text")

// Secondary descriptions
styles.BodyMuted.Render("Supporting information")

// Help text, instructions
styles.BodyDim.Render("Help and instructions")
```

### Special Text
```go
// Code snippets, commands
styles.Code.Render("netlab start")

// Important highlights
styles.Highlight.Render("Key Concept")
```

## 📦 Layout Components

### Containers
```go
// Basic content container
styles.Container.Render(content)

// Card for grouped content
styles.Card.Render(cardContent)

// Panel for highlighted sections
styles.Panel.Render(importantContent)
```

### Interactive Elements
```go
// Primary action button
styles.Button.Render("Start Module")

// Secondary action
styles.ButtonSecondary.Render("Learn More")

// Status indicators
styles.StatusSuccess.Render("✅ READY")
styles.StatusWarning.Render("🚧 WIP")
styles.StatusError.Render("❌ ERROR")
styles.StatusInfo.Render("📋 PLANNED")
```

## 📋 List Design Patterns

### Module Lists
- Use 3-line items with title, description, and status
- Show module number, clear status indicators
- Highlight active selection with background color
- Include progress indicators

### Content Navigation
- Use breadcrumbs for deep navigation
- Consistent key binding indicators
- Clear help text at bottom

## 🎯 Module Layout Standards

### Module Header
```go
// Standard module header with logo
components.RenderCompactLogo()

// Module title with description
styles.ModuleTitle.Render("OSI Model - Layer by Layer")
```

### Content Sections
```go
// Section with left border
styles.ModuleSection.Render(sectionContent)

// Examples and code blocks
styles.ModuleExample.Render(exampleContent)

// Quiz and interactive elements
styles.ModuleQuiz.Render(quizContent)
```

### Navigation Footer
```go
// Consistent key bindings
helpKeys := []string{
    styles.KeyBinding.Render("↑/↓") + " navigate",
    styles.KeyBinding.Render("Enter") + " select", 
    styles.KeyBinding.Render("q") + " quit",
}
styles.Help.Render(strings.Join(helpKeys, " • "))
```

## 🔧 Component Guidelines

### Spacing and Margins
- Use consistent margins: `Margin(1, 0)` for vertical, `Margin(0, 1)` for horizontal
- Add padding to containers: `Padding(1, 2)` 
- Use spacing between related elements

### Borders and Visual Hierarchy
- **Rounded borders** for cards and containers
- **Thick borders** for panels and important content
- **Left borders** for content sections
- **Full borders** with colors for status indicators

### Responsive Design
- Minimum width: 80 characters
- Adapt to terminal size with `tea.WindowSizeMsg`
- Center content when space allows
- Graceful degradation for small terminals

## 📱 Interaction Patterns

### Navigation
- **Arrow keys** for list navigation
- **Enter** to select/activate
- **q/Ctrl+C/Esc** to go back/quit
- **Tab** for focus switching (future enhancement)

### Feedback
- Immediate visual feedback on selection
- Status indicators for progress
- Loading states for async operations
- Error states with clear messages

## 🎬 Animation Guidelines

### Transitions (Future Enhancement)
- Subtle fade-ins for content loading
- Smooth color transitions for state changes
- Minimal, purposeful animations
- Respect terminal capabilities

### Motion Principles
- **Functional** - Animations should serve a purpose
- **Fast** - Quick transitions (200-300ms equivalent)
- **Consistent** - Same timing across similar interactions

## 📐 Layout Patterns

### Welcome Screen
```
┌─ Header with Logo ─┐
│     Instructions   │
├─ Module List ──────┤
│ 1. Module Title    │
│    Description  ✅ │
│                    │
├─ Progress ─────────┤
│ Help & Keybinds    │
└────────────────────┘
```

### Module Content
```
┌─ Compact Logo ─────┐
├─ Module Title ─────┤
│                    │
│ Content Sections   │
│ with left borders  │
│                    │
│ ┌─ Examples ─────┐ │
│ │ Code blocks    │ │
│ └────────────────┘ │
│                    │
├─ Progress/Nav ─────┤
│ Help & Keybinds    │
└────────────────────┘
```

## 🔍 Accessibility

### Color Contrast
- Ensure sufficient contrast ratios
- Don't rely solely on color for information
- Use symbols and text alongside colors

### Keyboard Navigation
- All functions accessible via keyboard
- Clear focus indicators
- Consistent key bindings across modules

### Screen Reader Support
- Semantic content structure
- Alternative text for visual elements
- Logical reading order

## 📋 Checklist for New Modules

### Visual Design
- [ ] Uses NetLab color palette
- [ ] Consistent typography scale
- [ ] Proper spacing and margins
- [ ] Status indicators where appropriate
- [ ] Responsive to terminal size

### Interaction Design
- [ ] Standard navigation keys
- [ ] Clear help text
- [ ] Immediate feedback on actions
- [ ] Graceful error handling

### Content Structure
- [ ] Module header with title
- [ ] Logical content sections
- [ ] Examples in styled blocks
- [ ] Progress indicators
- [ ] Consistent footer

### Code Quality
- [ ] Uses `netlab/pkg/styles` package
- [ ] Reusable components when possible
- [ ] Follows TUI patterns
- [ ] Proper error handling
- [ ] Clean, readable code

## 📚 Resources

### References
- [Charm Lip Gloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Bubble Tea Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [Terminal Color Standards](https://en.wikipedia.org/wiki/ANSI_escape_code)

### Testing
- Test on different terminal sizes
- Verify colors in various terminals (iTerm2, Terminal.app, etc.)
- Check with different themes and color schemes
- Test keyboard navigation thoroughly

---

**Remember**: Consistency is key. When in doubt, reference existing patterns and components rather than creating new ones. 