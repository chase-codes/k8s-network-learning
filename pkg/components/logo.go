package components

import (
	"strings"

	"netlab/pkg/styles"

	"github.com/charmbracelet/lipgloss"
)

// Logo asset embedded at build time
// Note: For now we'll use inline logo, can be enhanced later

// RenderLogo returns the NetLab ASCII logo with styling
func RenderLogo() string {
	// Try to read from embedded assets first, fallback to inline
	logoText := `███╗   ██╗███████╗████████╗██╗      █████╗ ██████╗ 
████╗  ██║██╔════╝╚══██╔══╝██║     ██╔══██╗██╔══██╗
██╔██╗ ██║█████╗     ██║   ██║     ███████║██████╔╝
██║╚██╗██║██╔══╝     ██║   ██║     ██╔══██║██╔══██╗
██║ ╚████║███████╗   ██║   ███████╗██║  ██║██████╔╝
╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═════╝`

	return styles.Logo.Render(logoText)
}

// RenderLogoWithTagline returns the logo with tagline
func RenderLogoWithTagline() string {
	logo := RenderLogo()
	tagline := styles.BodyMuted.
		Align(lipgloss.Center).
		Margin(0, 0, 2, 0).
		Render("Interactive Networking Learning Environment")

	return lipgloss.JoinVertical(lipgloss.Center, logo, tagline)
}

// RenderCompactLogo returns a smaller version for headers
func RenderCompactLogo() string {
	compactText := "╔═══ NETLAB ═══╗"
	return styles.H2.
		Align(lipgloss.Center).
		Margin(0, 0, 1, 0).
		Render(compactText)
}

// RenderWelcomeHeader creates the complete welcome header
func RenderWelcomeHeader(width int) string {
	logo := RenderLogoWithTagline()

	// Create a decorative border
	border := strings.Repeat("═", width-4)
	topBorder := "╔═" + border + "═╗"
	bottomBorder := "╚═" + border + "═╝"

	primaryStyle := lipgloss.NewStyle().Foreground(styles.Primary)
	header := lipgloss.JoinVertical(
		lipgloss.Center,
		primaryStyle.Render(topBorder),
		logo,
		primaryStyle.Render(bottomBorder),
	)

	return lipgloss.Place(width, lipgloss.Height(header), lipgloss.Center, lipgloss.Center, header)
}
