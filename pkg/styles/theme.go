package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// NetLab Color Palette
var (
	// Primary colors
	Primary   = lipgloss.Color("#00D4FF") // Bright cyan - tech/network theme
	Secondary = lipgloss.Color("#7C3AED") // Purple - learning theme
	Accent    = lipgloss.Color("#F59E0B") // Amber - highlights

	// Status colors
	Success = lipgloss.Color("#10B981") // Green
	Warning = lipgloss.Color("#F59E0B") // Amber
	Error   = lipgloss.Color("#EF4444") // Red
	Info    = lipgloss.Color("#3B82F6") // Blue

	// Neutral colors
	Text      = lipgloss.Color("#F8FAFC") // Near white
	TextMuted = lipgloss.Color("#94A3B8") // Light gray
	TextDim   = lipgloss.Color("#64748B") // Medium gray

	// Background colors
	Background     = lipgloss.Color("#0F172A") // Dark slate
	BackgroundCard = lipgloss.Color("#1E293B") // Lighter slate
	Border         = lipgloss.Color("#334155") // Border gray
)

// Typography Styles
var (
	// Headers
	H1 = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true).
		Margin(1, 0).
		Padding(0, 1)

	H2 = lipgloss.NewStyle().
		Foreground(Secondary).
		Bold(true).
		Margin(1, 0, 0, 0)

	H3 = lipgloss.NewStyle().
		Foreground(Accent).
		Bold(true)

	// Body text
	Body = lipgloss.NewStyle().
		Foreground(Text).
		Margin(0, 0, 1, 0)

	BodyMuted = lipgloss.NewStyle().
			Foreground(TextMuted)

	BodyDim = lipgloss.NewStyle().
		Foreground(TextDim)

	// Special text
	Code = lipgloss.NewStyle().
		Foreground(Primary).
		Background(BackgroundCard).
		Padding(0, 1)

	Highlight = lipgloss.NewStyle().
			Foreground(Background).
			Background(Accent).
			Bold(true).
			Padding(0, 1)
)

// Layout Styles
var (
	// Containers
	Container = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 0)

	Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Border).
		Background(BackgroundCard).
		Padding(1, 2).
		Margin(1, 0)

	Panel = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(Primary).
		Padding(1, 2).
		Margin(1, 0)

	// Interactive elements
	Button = lipgloss.NewStyle().
		Foreground(Background).
		Background(Primary).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1)

	ButtonSecondary = lipgloss.NewStyle().
			Foreground(Text).
			Background(Secondary).
			Bold(true).
			Padding(0, 2).
			Margin(0, 1)

	// Status indicators
	StatusSuccess = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	StatusWarning = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)

	StatusError = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	StatusInfo = lipgloss.NewStyle().
			Foreground(Info).
			Bold(true)
)

// Component Styles
var (
	// Lists
	ListItem = lipgloss.NewStyle().
			PaddingLeft(2).
			Margin(0, 0, 0, 2)

	ListItemActive = lipgloss.NewStyle().
			Foreground(Background).
			Background(Primary).
			Bold(true).
			PaddingLeft(1).
			Margin(0, 0, 0, 1)

	ListItemHover = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			PaddingLeft(1).
			Margin(0, 0, 0, 1)

	// Navigation
	Breadcrumb = lipgloss.NewStyle().
			Foreground(TextMuted).
			Margin(0, 0, 1, 0)

	NavItem = lipgloss.NewStyle().
		Foreground(Text).
		Padding(0, 1)

	NavItemActive = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Padding(0, 1)

	// Help and instructions
	Help = lipgloss.NewStyle().
		Foreground(TextDim).
		Margin(1, 0, 0, 0).
		Italic(true)

	KeyBinding = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	Instruction = lipgloss.NewStyle().
			Foreground(TextMuted).
			Margin(0, 0, 1, 0)
)

// Utility functions
func WithWidth(style lipgloss.Style, width int) lipgloss.Style {
	return style.Width(width)
}

func WithHeight(style lipgloss.Style, height int) lipgloss.Style {
	return style.Height(height)
}

func WithBorder(style lipgloss.Style, border lipgloss.Border, color lipgloss.Color) lipgloss.Style {
	return style.Border(border).BorderForeground(color)
}

// Logo styling
var Logo = lipgloss.NewStyle().
	Foreground(Primary).
	Bold(true).
	Align(lipgloss.Center).
	Margin(1, 0, 2, 0)

// Progress indicators
var (
	ProgressBar = lipgloss.NewStyle().
			Background(Border).
			Height(1)

	ProgressFill = lipgloss.NewStyle().
			Background(Primary).
			Height(1)

	Spinner = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)
)

// Module-specific styles
var (
	ModuleTitle = lipgloss.NewStyle().
			Foreground(Primary).
			Background(BackgroundCard).
			Bold(true).
			Padding(1, 2).
			Margin(0, 0, 1, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary)

	ModuleSection = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(Secondary).
			PaddingLeft(2).
			Margin(1, 0)

	ModuleExample = lipgloss.NewStyle().
			Background(BackgroundCard).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Accent).
			Padding(1, 2).
			Margin(1, 0)

	ModuleQuiz = lipgloss.NewStyle().
			Background(lipgloss.Color("#1F2937")).
			Border(lipgloss.ThickBorder()).
			BorderForeground(Warning).
			Padding(1, 2).
			Margin(1, 0)
)
