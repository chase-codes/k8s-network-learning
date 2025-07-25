package tui

import (
	"fmt"
	"netlab/internal/utils"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	depTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true).
			Margin(1, 0, 1, 0)

	okStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	missingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	depHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true).
			Margin(1, 0, 0, 0)

	buttonStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("12")).
			Foreground(lipgloss.Color("15")).
			Padding(0, 2).
			Margin(0, 1, 0, 0).
			Bold(true)

	selectedButtonStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("170")).
				Foreground(lipgloss.Color("15")).
				Padding(0, 2).
				Margin(0, 1, 0, 0).
				Bold(true)
)

type DependencyCheckResult int

const (
	CheckContinue DependencyCheckResult = iota
	CheckShowGuide
	CheckAbort
	CheckRecheck
)

type dependencyCheckModel struct {
	moduleID       string
	moduleName     string
	dependencies   []utils.DependencyStatus
	missingDeps    []utils.DependencyStatus
	allGood        bool
	showingGuide   bool
	viewport       viewport.Model
	selectedButton int
	width          int
	height         int
	result         DependencyCheckResult
	done           bool
}

func NewDependencyCheck(moduleID, moduleName string) *dependencyCheckModel {
	// Check dependencies for the module
	deps, allGood := utils.CheckModuleDependencies(moduleID)

	// Filter missing dependencies
	var missing []utils.DependencyStatus
	for _, dep := range deps {
		if dep.Status != "ok" {
			missing = append(missing, dep)
		}
	}

	vp := viewport.New(60, 10)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 2)

	return &dependencyCheckModel{
		moduleID:     moduleID,
		moduleName:   moduleName,
		dependencies: deps,
		missingDeps:  missing,
		allGood:      allGood,
		viewport:     vp,
	}
}

func (m *dependencyCheckModel) Init() tea.Cmd {
	return nil
}

func (m *dependencyCheckModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = m.width - 6
		m.viewport.Height = m.height - 12
		return m, nil

	case recheckMsg:
		m.updateDependencies(msg)
		return m, nil

	case tea.KeyMsg:
		if m.showingGuide {
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				m.showingGuide = false
				return m, nil
			default:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}

		switch msg.String() {
		case "ctrl+c", "q":
			m.result = CheckAbort
			m.done = true
			return m, tea.Quit

		case "left", "h":
			if m.selectedButton > 0 {
				m.selectedButton--
			}

		case "right", "l":
			maxButton := 2
			if m.allGood {
				maxButton = 1 // Only "Continue" and "Check Again" when all good
			} else {
				maxButton = 3 // "Continue", "Install Guide", "Check Again" when missing deps
			}
			if m.selectedButton < maxButton {
				m.selectedButton++
			}

		case "enter", " ":
			switch m.selectedButton {
			case 0: // Continue
				m.result = CheckContinue
				m.done = true
				return m, tea.Quit
			case 1: // Install Guide (only shown if missing deps)
				if !m.allGood {
					m.showingGuide = true
					m.viewport.SetContent(utils.GetInstallationGuide(m.missingDeps))
					return m, nil
				} else {
					// This is "Check Again" when all good
					m.result = CheckRecheck
					return m, m.recheckDependencies()
				}
			case 2: // Check Again
				m.result = CheckRecheck
				return m, m.recheckDependencies()
			}
		}
	}

	return m, nil
}

func (m *dependencyCheckModel) recheckDependencies() tea.Cmd {
	return func() tea.Msg {
		// Re-check dependencies
		deps, allGood := utils.CheckModuleDependencies(m.moduleID)

		// Filter missing dependencies
		var missing []utils.DependencyStatus
		for _, dep := range deps {
			if dep.Status != "ok" {
				missing = append(missing, dep)
			}
		}

		return recheckMsg{
			dependencies: deps,
			missingDeps:  missing,
			allGood:      allGood,
		}
	}
}

type recheckMsg struct {
	dependencies []utils.DependencyStatus
	missingDeps  []utils.DependencyStatus
	allGood      bool
}

func (m *dependencyCheckModel) View() string {
	if m.done {
		return ""
	}

	if m.showingGuide {
		return m.viewInstallationGuide()
	}

	var s strings.Builder

	// Title
	s.WriteString(depTitleStyle.Render(fmt.Sprintf("🔍 Dependency Check: %s", m.moduleName)))
	s.WriteString("\n\n")

	// Status summary
	if m.allGood {
		s.WriteString(okStyle.Render("✅ All dependencies satisfied!"))
		s.WriteString("\n\n")
	} else {
		s.WriteString(missingStyle.Render(fmt.Sprintf("⚠️  Missing %d dependencies", len(m.missingDeps))))
		s.WriteString("\n\n")
	}

	// Dependency list
	for _, dep := range m.dependencies {
		switch dep.Status {
		case "ok":
			s.WriteString(okStyle.Render(fmt.Sprintf("✓ %s", dep.Name)))
		case "missing":
			s.WriteString(missingStyle.Render(fmt.Sprintf("✗ %s - Not found", dep.Name)))
		case "error":
			s.WriteString(missingStyle.Render(fmt.Sprintf("✗ %s - Error", dep.Name)))
		}
		s.WriteString("\n")
	}

	s.WriteString("\n")

	// Warning if missing dependencies
	if !m.allGood {
		s.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Render("⚠️  Some features may not work without these dependencies."))
		s.WriteString("\n\n")
	}

	// Buttons
	s.WriteString(m.renderButtons())
	s.WriteString("\n\n")

	// Help
	s.WriteString(depHelpStyle.Render("Use ← → to navigate, Enter to select, q to quit"))

	return s.String()
}

func (m *dependencyCheckModel) renderButtons() string {
	buttons := []string{}

	// Continue button (always available)
	if m.selectedButton == 0 {
		buttons = append(buttons, selectedButtonStyle.Render("Continue"))
	} else {
		buttons = append(buttons, buttonStyle.Render("Continue"))
	}

	if !m.allGood {
		// Install Guide button (only when missing deps)
		if m.selectedButton == 1 {
			buttons = append(buttons, selectedButtonStyle.Render("Install Guide"))
		} else {
			buttons = append(buttons, buttonStyle.Render("Install Guide"))
		}

		// Check Again button
		if m.selectedButton == 2 {
			buttons = append(buttons, selectedButtonStyle.Render("Check Again"))
		} else {
			buttons = append(buttons, buttonStyle.Render("Check Again"))
		}
	} else {
		// Only Check Again button when all good
		if m.selectedButton == 1 {
			buttons = append(buttons, selectedButtonStyle.Render("Check Again"))
		} else {
			buttons = append(buttons, buttonStyle.Render("Check Again"))
		}
	}

	return strings.Join(buttons, "")
}

func (m *dependencyCheckModel) viewInstallationGuide() string {
	var s strings.Builder

	s.WriteString(depTitleStyle.Render("📦 Installation Guide"))
	s.WriteString("\n\n")

	s.WriteString(m.viewport.View())
	s.WriteString("\n\n")

	s.WriteString(depHelpStyle.Render("Press q or Esc to go back"))

	return s.String()
}

func (m *dependencyCheckModel) GetResult() DependencyCheckResult {
	return m.result
}

func (m *dependencyCheckModel) AllDependenciesSatisfied() bool {
	return m.allGood
}

// Handle recheck message
func (m *dependencyCheckModel) updateDependencies(msg recheckMsg) {
	m.dependencies = msg.dependencies
	m.missingDeps = msg.missingDeps
	m.allGood = msg.allGood
	m.selectedButton = 0 // Reset selection
}

// Update method to handle recheckMsg
func (m *dependencyCheckModel) handleRecheck(msg recheckMsg) (*dependencyCheckModel, tea.Cmd) {
	m.updateDependencies(msg)
	return m, nil
}

// RunDependencyCheck runs the dependency check TUI and returns the result
func RunDependencyCheck(moduleID, moduleName string) (DependencyCheckResult, error) {
	model := NewDependencyCheck(moduleID, moduleName)

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return CheckAbort, err
	}

	result := finalModel.(*dependencyCheckModel).GetResult()
	return result, nil
}
