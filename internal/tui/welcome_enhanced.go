package tui

import (
	"fmt"
	"io"
	"strings"

	"netlab/pkg/components"
	"netlab/pkg/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	enhancedListHeight = 12
	minWidth           = 80
)

type enhancedItem struct {
	title       string
	description string
	moduleID    string
	status      string // "ready", "planned", "wip"
}

func (i enhancedItem) FilterValue() string { return "" }
func (i enhancedItem) Title() string       { return i.title }
func (i enhancedItem) Description() string { return i.description }

type enhancedItemDelegate struct{}

func (d enhancedItemDelegate) Height() int                             { return 3 }
func (d enhancedItemDelegate) Spacing() int                            { return 1 }
func (d enhancedItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d enhancedItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(enhancedItem)
	if !ok {
		return
	}

	// Status indicator
	var statusStyle lipgloss.Style
	var statusText string
	switch i.status {
	case "ready":
		statusStyle = styles.StatusSuccess
		statusText = "✅ READY"
	case "wip":
		statusStyle = styles.StatusWarning
		statusText = "🚧 WIP"
	default: // planned
		statusStyle = styles.StatusInfo
		statusText = "📋 PLANNED"
	}

	// Module number and title
	moduleNum := fmt.Sprintf("%d.", index+1)
	title := fmt.Sprintf("%s %s", moduleNum, i.title)

	var titleStyle lipgloss.Style
	var itemStyle lipgloss.Style

	if index == m.Index() {
		// Selected item
		titleStyle = styles.ListItemActive.Copy().Width(m.Width() - 4)
		itemStyle = styles.ListItemActive.Copy().
			Foreground(styles.Background).
			Width(m.Width() - 4)
	} else {
		// Normal item
		titleStyle = styles.ListItem.Copy().
			Foreground(styles.Text).
			Bold(true)
		itemStyle = styles.ListItem.Copy().
			Foreground(styles.TextMuted)
	}

	// Render the item
	renderedTitle := titleStyle.Render(title)
	renderedDesc := itemStyle.Render(i.description)
	renderedStatus := statusStyle.Render(statusText)

	// Combine with status
	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		renderedTitle,
		strings.Repeat(" ", max(0, m.Width()-lipgloss.Width(renderedTitle)-lipgloss.Width(renderedStatus)-4)),
		renderedStatus)

	output := lipgloss.JoinVertical(lipgloss.Left, line1, renderedDesc)
	fmt.Fprint(w, output)
}

type enhancedModel struct {
	list     list.Model
	choice   string
	quitting bool
	width    int
	height   int
}

func (m enhancedModel) Init() tea.Cmd {
	return nil
}

func (m enhancedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 15) // Account for header and footer
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(enhancedItem)
			if ok {
				m.choice = i.moduleID
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m enhancedModel) View() string {
	if m.quitting {
		return styles.Help.Render("Thanks for using NetLab! 🚀")
	}

	// Header with logo
	header := components.RenderWelcomeHeader(max(m.width, minWidth))

	// Instructions
	instructions := styles.Instruction.
		Align(lipgloss.Center).
		Margin(1, 0).
		Render("Select a learning module to begin your networking journey")

	// List of modules
	listView := m.list.View()

	// Footer with help
	helpKeys := []string{
		styles.KeyBinding.Render("↑/↓") + " navigate",
		styles.KeyBinding.Render("Enter") + " select",
		styles.KeyBinding.Render("q") + " quit",
	}
	helpText := styles.Help.
		Align(lipgloss.Center).
		Margin(1, 0).
		Render(strings.Join(helpKeys, " • "))

	// Progress indicator
	totalModules := len(m.list.Items())
	readyModules := 0
	for _, item := range m.list.Items() {
		if item.(enhancedItem).status == "ready" {
			readyModules++
		}
	}
	progress := fmt.Sprintf("Progress: %d/%d modules ready", readyModules, totalModules)
	progressText := styles.BodyMuted.
		Align(lipgloss.Center).
		Render(progress)

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		instructions,
		listView,
		progressText,
		helpText,
	)

	// Center everything and add some padding
	return lipgloss.Place(
		max(m.width, minWidth),
		max(m.height, 30),
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// StartEnhancedWelcome launches the enhanced welcome screen
func StartEnhancedWelcome() (string, error) {
	items := []list.Item{
		enhancedItem{
			title:       "OSI Model Fundamentals",
			description: "Learn the seven layers of network communication",
			moduleID:    "01-osi-model",
			status:      "ready",
		},
		enhancedItem{
			title:       "TCP/IP Stack Deep Dive",
			description: "Explore the Internet Protocol suite in detail",
			moduleID:    "02-tcp-ip",
			status:      "planned",
		},
		enhancedItem{
			title:       "Subnetting and CIDR",
			description: "Master network segmentation and addressing",
			moduleID:    "03-subnetting",
			status:      "planned",
		},
		enhancedItem{
			title:       "Routing Protocols",
			description: "Understand how packets find their destination",
			moduleID:    "04-routing",
			status:      "planned",
		},
		enhancedItem{
			title:       "Kubernetes Networking",
			description: "Container networking in orchestrated environments",
			moduleID:    "05-k8s-networking",
			status:      "planned",
		},
		enhancedItem{
			title:       "Container Network Interface",
			description: "CNI specifications and implementations",
			moduleID:    "06-cni",
			status:      "planned",
		},
		enhancedItem{
			title:       "Service Mesh Concepts",
			description: "Advanced traffic management and observability",
			moduleID:    "07-service-mesh",
			status:      "planned",
		},
	}

	l := list.New(items, enhancedItemDelegate{}, minWidth, enhancedListHeight)
	l.Title = "" // We'll handle the title in our custom view
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false) // We'll show custom help

	// Custom list styles
	l.Styles.NoItems = styles.BodyMuted
	l.Styles.PaginationStyle = styles.BodyDim.Align(lipgloss.Center)

	m := enhancedModel{
		list:   l,
		width:  minWidth,
		height: 30,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	// Extract the selected module ID
	if finalModel.(enhancedModel).choice != "" {
		return finalModel.(enhancedModel).choice, nil
	}

	return "", nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
