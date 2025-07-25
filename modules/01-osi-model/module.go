package osimodel

import (
	"fmt"
	"strings"

	"netlab/pkg/styles"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// listItem represents an OSI layer in the list
type listItem struct {
	layer OSILayer
}

func (i listItem) FilterValue() string { return i.layer.Name }
func (i listItem) Title() string       { return fmt.Sprintf("Layer %d: %s", i.layer.Number, i.layer.Name) }
func (i listItem) Description() string { return i.layer.Function }

// Model represents the main TUI model for the OSI module
type Model struct {
	list           list.Model
	layers         []OSILayer
	selectedLayer  *OSILayer
	showMnemonic   bool
	ready          bool
	width          int
	height         int
	listWidth      int
	listHeight     int
	viewportWidth  int
	viewportHeight int
	viewport       viewport.Model
	quitting       bool
	advanceToLab   bool
}

// NewModel creates a new OSI module model
func NewModel() Model {
	layers := GetOSILayers()

	// Create list items from layers (reverse order to show 7 at top)
	items := make([]list.Item, len(layers))
	for i, layer := range layers {
		items[len(layers)-1-i] = listItem{layer: layer}
	}

	// Configure list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "OSI Model - Seven Layers"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styles.H2
	l.Styles.PaginationStyle = styles.BodyMuted
	l.Styles.HelpStyle = styles.Help

	// Set selected layer to the first item (Layer 7)
	var selectedLayer *OSILayer
	if len(layers) > 0 {
		selectedLayer = &layers[0] // Layer 7 (Application)
	}

	return Model{
		list:          l,
		layers:        layers,
		selectedLayer: selectedLayer,
		showMnemonic:  false,
		ready:         false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnableMouseCellMotion
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Check for minimum terminal size
		if m.width < 60 || m.height < 20 {
			return m, nil // Will show error message in View()
		}

		// Calculate dynamic dimensions
		headerHeight := 3 // Minimal header
		footerHeight := 2 // Minimal footer
		availableHeight := m.height - headerHeight - footerHeight

		// Left pane (OSI list): 1/4 of width, minimum 25 chars
		m.listWidth = m.width / 4
		if m.listWidth < 25 {
			m.listWidth = 25
		}
		m.listHeight = availableHeight

		// Right pane (detail viewport): remaining width minus separator
		m.viewportWidth = m.width - m.listWidth - 2
		m.viewportHeight = availableHeight

		// Update list size
		m.list.SetSize(m.listWidth, m.listHeight)

		// Initialize or update viewport
		if !m.ready {
			m.viewport = viewport.New(m.viewportWidth, m.viewportHeight)
			m.viewport.SetContent(m.getDetailContent())
			m.ready = true
		} else {
			m.viewport.Width = m.viewportWidth
			m.viewport.Height = m.viewportHeight
			m.viewport.SetContent(m.getDetailContent())
		}

	case tea.MouseMsg:
		// Forward mouse events to viewport for text selection
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			if m.showMnemonic {
				m.showMnemonic = false
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit

		case "m":
			m.showMnemonic = !m.showMnemonic
			return m, nil

		case "v":
			m.advanceToLab = true
			return m, tea.Quit

		case "j", "down":
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			m.updateSelectedLayer()
			if m.ready {
				m.viewport.SetContent(m.getDetailContent())
			}
			return m, cmd

		case "k", "up":
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			m.updateSelectedLayer()
			if m.ready {
				m.viewport.SetContent(m.getDetailContent())
			}
			return m, cmd
		}
	}

	// Handle viewport scrolling and list updates
	var cmd tea.Cmd
	var vpCmd tea.Cmd

	m.list, cmd = m.list.Update(msg)
	m.updateSelectedLayer()

	if m.ready {
		m.viewport.SetContent(m.getDetailContent())
		m.viewport, vpCmd = m.viewport.Update(msg)
	}

	return m, tea.Batch(cmd, vpCmd)
}

func (m *Model) updateSelectedLayer() {
	if item, ok := m.list.SelectedItem().(listItem); ok {
		m.selectedLayer = &item.layer
	}
}

// getDetailContent returns the content for the detail viewport
func (m Model) getDetailContent() string {
	if m.selectedLayer == nil {
		return "Select a layer to view details"
	}

	layer := *m.selectedLayer
	var content strings.Builder

	// Layer title and number
	content.WriteString(styles.H2.Render(fmt.Sprintf("Layer %d: %s", layer.Number, layer.Name)))
	content.WriteString("\n\n")

	// Description
	content.WriteString(styles.ModuleSection.Render(layer.Description))
	content.WriteString("\n\n")

	// Function
	content.WriteString(styles.H3.Render("Function"))
	content.WriteString("\n")
	content.WriteString(styles.Body.Render(layer.Function))
	content.WriteString("\n\n")

	// Protocols
	content.WriteString(styles.H3.Render("Common Protocols"))
	content.WriteString("\n")
	protocolList := "• " + strings.Join(layer.Protocols, "\n• ")
	content.WriteString(styles.ModuleExample.Render(protocolList))
	content.WriteString("\n\n")

	// Real-world analogy
	content.WriteString(styles.H3.Render("Real-World Analogy"))
	content.WriteString("\n")
	content.WriteString(styles.Highlight.Render(layer.Analogy))
	content.WriteString("\n\n")

	// Header type
	content.WriteString(styles.H3.Render("Header Information"))
	content.WriteString("\n")
	content.WriteString(styles.Body.Render(layer.HeaderType))
	content.WriteString("\n\n")

	// CLI Tools
	content.WriteString(styles.H3.Render("Useful CLI Tools"))
	content.WriteString("\n")
	toolList := strings.Join(layer.CLITools, ", ")
	content.WriteString(styles.Code.Render(toolList))
	content.WriteString("\n\n")

	// Examples
	content.WriteString(styles.H3.Render("Examples"))
	content.WriteString("\n")
	exampleList := "• " + strings.Join(layer.Examples, "\n• ")
	content.WriteString(styles.Body.Render(exampleList))
	content.WriteString("\n\n")

	// External documentation
	content.WriteString(styles.H3.Render("Learn More"))
	content.WriteString("\n")
	content.WriteString(styles.BodyMuted.Render("📖 " + layer.ExternalDoc))
	content.WriteString("\n\n")

	// Kubernetes context
	k8sContext := GetKubernetesContext()
	if context, exists := k8sContext[layer.Number]; exists {
		content.WriteString(styles.H3.Render("☸️ Kubernetes Context"))
		content.WriteString("\n")
		content.WriteString(styles.ModuleSection.Render(context))
	}

	return content.String()
}

func (m Model) View() string {
	// Check for minimum terminal size first
	if m.width < 60 || m.height < 20 {
		return fmt.Sprintf("⚠️  Terminal too small (%dx%d), please resize to at least 60x20", m.width, m.height)
	}

	if !m.ready {
		return styles.Body.Render("\n  Initializing OSI Model module...")
	}

	if m.showMnemonic {
		return m.mnemonicView()
	}

	header := m.headerView()
	footer := m.footerView()

	// Create styled left and right panes
	leftStyle := lipgloss.NewStyle().
		Width(m.listWidth).
		Height(m.listHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border)

	rightStyle := lipgloss.NewStyle().
		Width(m.viewportWidth).
		Height(m.viewportHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border)

	left := leftStyle.Render(m.list.View())
	right := rightStyle.Render(m.viewport.View())

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

func (m Model) headerView() string {
	// Minimal header to prevent cutoff
	title := styles.H2.
		Width(m.width-4).
		Align(lipgloss.Center).
		Margin(0, 0).
		Render("NetLab OSI Model - Interactive Layer Explorer")

	breadcrumb := styles.BodyMuted.
		Margin(0, 0).
		Render("NetLab > Fundamentals > OSI Model")

	return lipgloss.JoinVertical(lipgloss.Left, title, breadcrumb)
}

func (m Model) footerView() string {
	helpKeys := []string{
		styles.KeyBinding.Render("↑/↓") + " navigate",
		styles.KeyBinding.Render("m") + " mnemonic",
		styles.KeyBinding.Render("v") + " packet lab",
		styles.KeyBinding.Render("q") + " quit",
	}
	helpText := styles.Help.Render(strings.Join(helpKeys, " • "))

	// Create separator line
	line := strings.Repeat("─", m.width-2)
	separator := styles.BodyDim.Render(line)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		separator,
		helpText,
	)
}

func (m Model) mnemonicView() string {
	// Check for minimum terminal size first
	if m.width < 60 || m.height < 20 {
		return fmt.Sprintf("⚠️  Terminal too small (%dx%d), please resize to at least 60x20", m.width, m.height)
	}

	var content strings.Builder

	content.WriteString(styles.H1.Render("🧠 OSI Layer Mnemonics"))
	content.WriteString("\n\n")

	content.WriteString(styles.Body.Render("Popular phrases to remember the OSI layers (Physical → Application):"))
	content.WriteString("\n\n")

	mnemonics := GetMnemonics()
	for i, mnemonic := range mnemonics {
		style := styles.ModuleExample
		if i == 0 {
			style = styles.Highlight
		}
		content.WriteString(style.Render(fmt.Sprintf("%d. %s", i+1, mnemonic)))
		content.WriteString("\n\n")
	}

	// Show the mapping
	content.WriteString(styles.H2.Render("Layer Mapping"))
	content.WriteString("\n")

	mapping := `P - Physical      (Layer 1)
D - Data Link     (Layer 2)  
N - Network       (Layer 3)
T - Transport     (Layer 4)
S - Session       (Layer 5)
P - Presentation  (Layer 6)
A - Application   (Layer 7)`

	content.WriteString(styles.ModuleExample.Render(mapping))
	content.WriteString("\n\n")

	// Instructions
	content.WriteString(styles.Help.Render("Press 'm' to return to layer explorer • Press 'q' to quit"))

	// Calculate proper dimensions for the mnemonic view
	headerHeight := 3
	footerHeight := 2
	availableHeight := m.height - headerHeight - footerHeight
	availableWidth := m.width - 4

	// Create a viewport-style container that respects terminal size
	containerStyle := lipgloss.NewStyle().
		Width(availableWidth).
		Height(availableHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Padding(1).
		Align(lipgloss.Center, lipgloss.Top)

	header := m.headerView()
	footer := m.footerView()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		containerStyle.Render(content.String()),
		footer,
	)
}

// Run starts the interactive OSI model TUI
func Run() error {
	m := NewModel()

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// Check if user wants to advance to packet lab
	if model, ok := finalModel.(Model); ok && model.advanceToLab {
		return RunWalkthrough()
	}

	return nil
}
