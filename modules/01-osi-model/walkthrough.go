package osimodel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"netlab/pkg/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PacketLayer represents a parsed layer from a network packet
type PacketLayer struct {
	OSILayer    int
	Name        string
	Headers     map[string]string
	RawData     string
	Explanation string
}

// WalkthroughModel represents the packet walkthrough TUI
type WalkthroughModel struct {
	layers         []PacketLayer
	currentIdx     int
	viewport       viewport.Model
	ready          bool
	width          int
	height         int
	labReady       bool
	labRunning     bool
	labError       string
	labProgress    float64
	labOutput      []string
	outputViewport viewport.Model
	showLabSetup   bool
}

func NewWalkthroughModel() WalkthroughModel {
	// Initialize with sample packet data
	layers := getSamplePacketLayers()

	return WalkthroughModel{
		layers:     layers,
		currentIdx: 0,
		ready:      false,
		labReady:   false,
	}
}

func (m WalkthroughModel) Init() tea.Cmd {
	return tea.Batch(
		m.checkLabStatus(),
		tea.EnableMouseCellMotion,
	)
}

func (m WalkthroughModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Check for minimum terminal size
		if m.width < 60 || m.height < 20 {
			return m, nil // Will show error message in View()
		}

		// Conservative height estimates to prevent cutoff
		headerHeight := 5 // Title + Breadcrumb + Progress + spacing
		footerHeight := 3 // Separator + help text + spacing
		availableHeight := msg.Height - headerHeight - footerHeight

		if availableHeight < 8 {
			availableHeight = 8 // Minimum height
		}

		if !m.ready {
			m.viewport = viewport.New(msg.Width-4, availableHeight) // Add margin
			m.viewport.SetContent(m.getLayerContent())

			// Initialize output viewport for lab setup
			m.outputViewport = viewport.New(msg.Width-8, availableHeight/2)

			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4 // Add margin
			m.viewport.Height = availableHeight
			m.viewport.SetContent(m.getLayerContent())

			// Update output viewport dimensions
			m.outputViewport.Width = msg.Width - 8
			m.outputViewport.Height = availableHeight / 2
		}

	case tea.MouseMsg:
		// Forward mouse events to the appropriate viewport
		if m.showLabSetup {
			var cmd tea.Cmd
			m.outputViewport, cmd = m.outputViewport.Update(msg)
			return m, cmd
		} else {
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "n", "right", "space":
			if m.currentIdx < len(m.layers)-1 {
				m.currentIdx++
				m.viewport.SetContent(m.getLayerContent())
			}
			return m, nil

		case "p", "left", "backspace":
			if m.currentIdx > 0 {
				m.currentIdx--
				m.viewport.SetContent(m.getLayerContent())
			}
			return m, nil

		case "r":
			if !m.labRunning {
				m.labRunning = true
				m.labError = ""
				m.showLabSetup = true
				m.labProgress = 0
				m.labOutput = []string{}
				m.labReady = false // Reset lab ready state
				return m, m.runLabSetup()
			}

		case "c":
			if !m.labRunning {
				m.labRunning = true
				m.labError = ""
				m.showLabSetup = true
				m.labProgress = 0
				m.labOutput = []string{}
				return m, m.runLabCleanup()
			}
			return m, nil

		case "e":
			// Export logs to file
			if len(m.labOutput) > 0 {
				return m, m.exportLogs()
			}
			return m, nil
		}

	case labStatusMsg:
		m.labRunning = false
		m.labReady = msg.ready
		m.labError = msg.error

		if m.labReady {
			// Load actual packet data and update layers
			if layers, err := LoadPacketData(); err == nil {
				m.layers = layers
			}
			// Update viewport content
			if m.ready {
				m.viewport.SetContent(m.getLayerContent())
			}
		}
		return m, nil

	case labSetupStartMsg:
		m.labProgress = 0
		m.labOutput = append(m.labOutput, "🚀 Starting lab setup...")

		// Add debugging info
		if wd, err := os.Getwd(); err == nil {
			m.labOutput = append(m.labOutput, fmt.Sprintf("📁 Working directory: %s", wd))
		}

		// Check if script exists
		if _, err := os.Stat("./scripts/k8s_lab.sh"); err == nil {
			m.labOutput = append(m.labOutput, "✅ Found lab script: ./scripts/k8s_lab.sh")
		} else {
			m.labOutput = append(m.labOutput, "❌ Lab script not found: ./scripts/k8s_lab.sh")
		}

		m.updateOutputViewport()
		return m, m.streamLabOutput()

	case labCleanupStartMsg:
		m.labProgress = 0
		m.labOutput = append(m.labOutput, "🧹 Starting lab cleanup...")

		// Add debugging info
		if wd, err := os.Getwd(); err == nil {
			m.labOutput = append(m.labOutput, fmt.Sprintf("📁 Working directory: %s", wd))
		}

		// Check if script exists
		if _, err := os.Stat("./scripts/k8s_lab.sh"); err == nil {
			m.labOutput = append(m.labOutput, "✅ Found lab script: ./scripts/k8s_lab.sh")
		} else {
			m.labOutput = append(m.labOutput, "❌ Lab script not found: ./scripts/k8s_lab.sh")
		}

		m.updateOutputViewport()
		return m, m.streamLabCleanup()

	case labOutputMsg:
		if msg.finished {
			m.labRunning = false
			if msg.success {
				m.labReady = true
				m.labError = ""
				m.labProgress = 100
				m.labOutput = append(m.labOutput, "✅ Lab setup completed successfully!")
				m.showLabSetup = false // Only hide on success
				// Load actual packet data and update layers
				if layers, err := LoadPacketData(); err == nil {
					m.layers = layers
				}
				if m.ready {
					m.viewport.SetContent(m.getLayerContent())
				}
			} else {
				m.labError = msg.error
				m.labProgress = 0
				// Add the full output message which includes user-friendly error and troubleshooting
				lines := strings.Split(msg.output, "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						m.labOutput = append(m.labOutput, line)
					}
				}
				// Keep showLabSetup = true so user can see the error and retry
			}
			m.updateOutputViewport()
		} else {
			// Parse output for progress and add to output log
			lines := strings.Split(msg.output, "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					m.labOutput = append(m.labOutput, line)
					// Update progress based on output content
					if progress := parseProgressFromOutput(line); progress > m.labProgress {
						m.labProgress = progress
					}
				}
			}
			m.updateOutputViewport()
		}
		return m, nil

	case logExportMsg:
		if msg.success {
			m.labOutput = append(m.labOutput, fmt.Sprintf("✅ Logs exported to: %s", msg.file))
		} else {
			m.labOutput = append(m.labOutput, fmt.Sprintf("❌ Export failed: %s", msg.error))
		}
		m.updateOutputViewport()
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m WalkthroughModel) View() string {
	// Check for minimum terminal size first
	if m.width < 60 || m.height < 20 {
		return fmt.Sprintf("⚠️  Terminal too small (%dx%d), please resize to at least 60x20", m.width, m.height)
	}

	if !m.ready {
		return styles.Body.Render("\n  Initializing packet walkthrough...")
	}

	if m.showLabSetup {
		return m.labSetupView()
	}

	header := m.headerView()
	footer := m.footerView()

	// Create a styled container for the viewport
	viewportStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Padding(1).
		Margin(0, 1)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		viewportStyle.Render(m.viewport.View()),
		footer,
	)
}

func (m WalkthroughModel) headerView() string {
	// Minimal header to prevent cutoff
	title := styles.H2.
		Width(m.width-4).
		Align(lipgloss.Center).
		Margin(0, 0).
		Render("NetLab OSI Model - Packet Analysis Lab")

	breadcrumb := styles.BodyMuted.
		Margin(0, 0).
		Render("NetLab > OSI Model > Packet Walkthrough")

	// Layer progress indicator
	current := m.currentIdx + 1
	total := len(m.layers)
	progress := fmt.Sprintf("Layer %d of %d", current, total)
	progressText := styles.BodyMuted.
		Margin(0, 0).
		Render(progress)

	return lipgloss.JoinVertical(lipgloss.Left, title, breadcrumb, progressText)
}

func (m WalkthroughModel) footerView() string {
	var helpKeys []string

	if m.labRunning {
		helpKeys = []string{
			styles.KeyBinding.Render("q") + " quit",
			styles.BodyMuted.Render("(lab setup running...)"),
		}
	} else if m.showLabSetup && m.labError != "" {
		helpKeys = []string{
			styles.KeyBinding.Render("r") + " retry lab setup",
			styles.KeyBinding.Render("c") + " cleanup lab",
		}
		if len(m.labOutput) > 0 {
			helpKeys = append(helpKeys, styles.KeyBinding.Render("e")+" export logs")
		}
		helpKeys = append(helpKeys,
			styles.KeyBinding.Render("q")+" quit",
			styles.BodyMuted.Render("(setup failed)"),
		)
	} else if !m.labReady {
		helpKeys = []string{
			styles.KeyBinding.Render("r") + " run lab setup",
			styles.KeyBinding.Render("c") + " cleanup lab",
		}
		if len(m.labOutput) > 0 {
			helpKeys = append(helpKeys, styles.KeyBinding.Render("e")+" export logs")
		}
		helpKeys = append(helpKeys, styles.KeyBinding.Render("q")+" quit")
	} else {
		helpKeys = []string{
			styles.KeyBinding.Render("←/→") + " navigate",
			styles.KeyBinding.Render("n") + " next",
			styles.KeyBinding.Render("p") + " prev",
			styles.KeyBinding.Render("c") + " cleanup lab",
			styles.KeyBinding.Render("q") + " quit",
		}
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

func (m WalkthroughModel) getLayerContent() string {
	if m.currentIdx >= len(m.layers) {
		return "No more layers"
	}

	layer := m.layers[m.currentIdx]
	var content strings.Builder

	// Check if lab is ready
	if !m.labReady {
		return m.getLabSetupContent()
	}

	// Layer header
	content.WriteString(styles.H1.Render(fmt.Sprintf("OSI Layer %d: %s", layer.OSILayer, layer.Name)))
	content.WriteString("\n\n")

	// Explanation
	content.WriteString(styles.ModuleSection.Render(layer.Explanation))
	content.WriteString("\n\n")

	// Headers section
	if len(layer.Headers) > 0 {
		content.WriteString(styles.H2.Render("📋 Headers & Fields"))
		content.WriteString("\n")

		var headerContent strings.Builder
		for field, value := range layer.Headers {
			headerContent.WriteString(fmt.Sprintf("%-20s: %s\n", field, value))
		}

		content.WriteString(styles.ModuleExample.Render(headerContent.String()))
		content.WriteString("\n\n")
	}

	// Raw data section
	if layer.RawData != "" {
		content.WriteString(styles.H2.Render("🔍 Raw Data"))
		content.WriteString("\n")
		content.WriteString(styles.Code.Render(layer.RawData))
		content.WriteString("\n\n")
	}

	// OSI Layer context
	osiLayers := GetOSILayers()
	for _, osiLayer := range osiLayers {
		if osiLayer.Number == layer.OSILayer {
			content.WriteString(styles.H2.Render("📚 OSI Layer Context"))
			content.WriteString("\n")
			content.WriteString(styles.Body.Render(osiLayer.Description))
			content.WriteString("\n\n")

			content.WriteString(styles.H3.Render("Key Concepts:"))
			content.WriteString("\n")
			conceptList := "• " + strings.Join(osiLayer.KeyConcepts, "\n• ")
			content.WriteString(styles.ModuleSection.Render(conceptList))
			break
		}
	}

	return content.String()
}

func (m WalkthroughModel) getLabSetupContent() string {
	var content strings.Builder

	content.WriteString(styles.H1.Render("🚀 Kubernetes Packet Analysis Lab"))
	content.WriteString("\n\n")

	// Show lab status
	if m.labRunning {
		content.WriteString(styles.StatusInfo.Render("🔄 Lab setup in progress..."))
		content.WriteString("\n")
		content.WriteString(styles.BodyMuted.Render("Setting up kind cluster, deploying nginx, and capturing packets. This may take several minutes."))
		content.WriteString("\n\n")
		return content.String()
	}

	if m.labError != "" {
		content.WriteString(styles.StatusError.Render("❌ Lab setup failed"))
		content.WriteString("\n")
		content.WriteString(styles.BodyMuted.Render(m.labError))
		content.WriteString("\n\n")
		content.WriteString(styles.Body.Render("Please check that Docker is running and all prerequisites are installed."))
		content.WriteString("\n")
		content.WriteString(styles.Help.Render("Run 'netlab doctor' to check dependencies"))
		content.WriteString("\n\n")
		content.WriteString(styles.Highlight.Render("Press 'r' to retry lab setup"))
		content.WriteString("\n\n")
		return content.String()
	}

	content.WriteString(styles.Body.Render("This lab will demonstrate OSI layers in action using a real HTTP request from a Pod to nginx running in a kind cluster."))
	content.WriteString("\n\n")

	content.WriteString(styles.H2.Render("Lab Components:"))
	content.WriteString("\n")

	components := `• kind cluster (local Kubernetes)
• nginx Deployment and Service  
• busybox Pod for making requests
• tcpdump for packet capture
• Parsed packet data for analysis`

	content.WriteString(styles.ModuleExample.Render(components))
	content.WriteString("\n\n")

	content.WriteString(styles.H2.Render("Prerequisites:"))
	content.WriteString("\n")

	prereqs := `• Docker (for kind) - will be started automatically if needed
• kubectl (Kubernetes CLI)
• kind (Kubernetes in Docker)  
• tcpdump (packet capture)
• tshark (packet analysis)`

	content.WriteString(styles.ModuleSection.Render(prereqs))
	content.WriteString("\n\n")

	content.WriteString(styles.H2.Render("What happens during setup:"))
	content.WriteString("\n")

	setupSteps := `1. Check and start Docker if needed (may take 30-60 seconds)
2. Create a kind Kubernetes cluster
3. Deploy nginx and busybox pods
4. Capture HTTP packets between pods
5. Generate packet analysis data`

	content.WriteString(styles.ModuleSection.Render(setupSteps))
	content.WriteString("\n\n")

	content.WriteString(styles.Highlight.Render("Press 'r' to run the lab setup script"))
	content.WriteString("\n")
	content.WriteString(styles.Highlight.Render("Press 'c' to cleanup the lab environment"))
	content.WriteString("\n\n")

	content.WriteString(styles.BodyMuted.Render("Note: If Docker isn't running, the script will attempt to start it automatically. First-time Docker startup may take longer."))
	content.WriteString("\n\n")

	content.WriteString(styles.Help.Render("💡 Tip: Run './scripts/k8s_lab.sh setup' or './scripts/k8s_lab.sh cleanup' manually if you prefer to see detailed output"))

	return content.String()
}

func (m WalkthroughModel) checkLabStatus() tea.Cmd {
	return func() tea.Msg {
		// Check if packet capture file exists
		_, err := os.Stat("modules/01-osi-model/assets/https-nginx.pcap")
		return labStatusMsg{
			ready: err == nil,
			error: "",
		}
	}
}

func (m WalkthroughModel) runLabSetup() tea.Cmd {
	return func() tea.Msg {
		// Start the lab setup process
		return labSetupStartMsg{}
	}
}

func (m WalkthroughModel) runLabCleanup() tea.Cmd {
	return func() tea.Msg {
		// Start the lab cleanup process
		return labCleanupStartMsg{}
	}
}

func (m WalkthroughModel) streamLabOutput() tea.Cmd {
	return func() tea.Msg {
		// First check if script exists
		if _, err := os.Stat("./scripts/k8s_lab.sh"); os.IsNotExist(err) {
			return labOutputMsg{
				output:   "📁 Lab script not found\n\nThe setup script './scripts/k8s_lab.sh' was not found.\n\nThis usually means:\n• You're not running from the NetLab project root directory\n• The script file is missing or moved\n\nTroubleshooting:\n• Make sure you're in the k8s-network-learning directory\n• Check if the file exists: 'ls -la scripts/'\n• If missing, clone the repository again or check if it was accidentally deleted",
				progress: 0,
				finished: true,
				error:    "📁 Lab script not found in expected location",
				success:  false,
			}
		}

		// Run the k8s lab setup script
		cmd := exec.Command("./scripts/k8s_lab.sh", "setup")
		cmd.Dir = "." // Ensure we're in the right directory

		// Execute and capture output
		output, err := cmd.CombinedOutput()

		if err != nil {
			// Parse the error and provide user-friendly explanation
			userError, troubleshooting := parseLabError(string(output), err)

			return labOutputMsg{
				output:   fmt.Sprintf("%s\n\n%s\n\nRaw output:\n%s", userError, troubleshooting, string(output)),
				progress: 0,
				finished: true,
				error:    userError,
				success:  false,
			}
		}

		// Check final status
		_, fileErr := os.Stat("modules/01-osi-model/assets/https-nginx.pcap")
		success := fileErr == nil

		return labOutputMsg{
			output:   string(output),
			progress: 100,
			finished: true,
			error:    "",
			success:  success,
		}
	}
}

func (m WalkthroughModel) streamLabCleanup() tea.Cmd {
	return func() tea.Msg {
		// First check if script exists
		if _, err := os.Stat("./scripts/k8s_lab.sh"); os.IsNotExist(err) {
			return labOutputMsg{
				output:   "📁 Lab script not found\n\nThe cleanup script './scripts/k8s_lab.sh' was not found.\n\nThis usually means:\n• You're not running from the NetLab project root directory\n• The script file is missing or moved\n\nTroubleshooting:\n• Make sure you're in the k8s-network-learning directory\n• You can manually run: 'kind delete cluster --name netlab-osi'",
				progress: 0,
				finished: true,
				error:    "📁 Lab script not found in expected location",
				success:  false,
			}
		}

		// Run the k8s lab cleanup script
		cmd := exec.Command("./scripts/k8s_lab.sh", "cleanup")
		cmd.Dir = "." // Ensure we're in the right directory

		// Execute and capture output
		output, err := cmd.CombinedOutput()

		if err != nil {
			// For cleanup, most errors are non-critical
			return labOutputMsg{
				output:   fmt.Sprintf("⚠️ Cleanup completed with warnings:\n\n%s", string(output)),
				progress: 100,
				finished: true,
				error:    "",
				success:  true, // Still consider it successful
			}
		}

		return labOutputMsg{
			output:   string(output),
			progress: 100,
			finished: true,
			error:    "",
			success:  true,
		}
	}
}

func parseProgressFromOutput(line string) float64 {
	// Parse common progress indicators from our script
	lower := strings.ToLower(line)

	switch {
	case strings.Contains(lower, "checking prerequisites") || strings.Contains(lower, "checking docker"):
		return 5
	case strings.Contains(lower, "creating kind cluster"):
		return 15
	case strings.Contains(lower, "kind cluster") && strings.Contains(lower, "created successfully"):
		return 30
	case strings.Contains(lower, "deploying nginx"):
		return 40
	case strings.Contains(lower, "nginx deployed successfully"):
		return 55
	case strings.Contains(lower, "deploying busybox"):
		return 65
	case strings.Contains(lower, "busybox pod deployed successfully"):
		return 75
	case strings.Contains(lower, "setting up packet capture") || strings.Contains(lower, "installing tcpdump"):
		return 80
	case strings.Contains(lower, "starting packet capture"):
		return 85
	case strings.Contains(lower, "making http request"):
		return 90
	case strings.Contains(lower, "copying packet capture") || strings.Contains(lower, "packet capture saved"):
		return 95
	case strings.Contains(lower, "lab setup completed") || (strings.Contains(lower, "success") && strings.Contains(lower, "capture")):
		return 100
	default:
		return 0
	}
}

// parseLabError analyzes script output and provides user-friendly error messages
func parseLabError(output string, err error) (userError, troubleshooting string) {
	lower := strings.ToLower(output)

	// Docker-related errors
	if strings.Contains(lower, "docker") && (strings.Contains(lower, "not found") || strings.Contains(lower, "command not found")) {
		return "🐳 Docker is not installed",
			"Please install Docker:\n• macOS: Download Docker Desktop from docker.com\n• Linux: Run 'curl -fsSL https://get.docker.com | sh'\n• Then restart your terminal"
	}

	if strings.Contains(lower, "docker") && (strings.Contains(lower, "permission denied") || strings.Contains(lower, "cannot connect")) {
		return "🐳 Docker connection issue",
			"The lab script will attempt to start Docker automatically.\nIf this fails:\n• Start Docker Desktop manually (macOS)\n• Run 'sudo systemctl start docker' (Linux)\n• Add your user to docker group: 'sudo usermod -aG docker $USER'"
	}

	// Docker starting process messages
	if strings.Contains(lower, "docker is not running") && strings.Contains(lower, "attempting to start") {
		return "🐳 Starting Docker",
			"Docker is being started automatically.\nThis may take 30-60 seconds on first startup.\nIf startup fails, please start Docker manually and try again."
	}

	if strings.Contains(lower, "waiting for docker") {
		return "🐳 Docker is starting",
			"Please wait while Docker starts up.\nThis is normal and may take up to 60 seconds.\nDocker Desktop needs time to initialize all components."
	}

	if strings.Contains(lower, "docker failed to start") {
		return "🐳 Docker startup failed",
			"Docker could not be started automatically.\nPlease:\n• Start Docker Desktop manually (macOS)\n• Run 'sudo systemctl start docker' (Linux)\n• Wait for Docker to be fully ready\n• Then try the lab again"
	}

	// kind-related errors
	if strings.Contains(lower, "kind") && strings.Contains(lower, "not found") {
		return "☸️  kind is not installed",
			"Please install kind:\n• macOS: 'brew install kind'\n• Linux: 'curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64'\n• Then: 'chmod +x ./kind && sudo mv ./kind /usr/local/bin/kind'"
	}

	// kubectl-related errors
	if strings.Contains(lower, "kubectl") && strings.Contains(lower, "not found") {
		return "☸️  kubectl is not installed",
			"Please install kubectl:\n• macOS: 'brew install kubectl'\n• Linux: Follow instructions at kubernetes.io/docs/tasks/tools/install-kubectl-linux/"
	}

	// tcpdump/tshark errors
	if strings.Contains(lower, "tcpdump") && strings.Contains(lower, "not found") {
		return "📡 tcpdump is not installed",
			"Please install tcpdump:\n• macOS: 'brew install tcpdump'\n• Linux: 'sudo apt-get install tcpdump' or 'sudo yum install tcpdump'"
	}

	if strings.Contains(lower, "tshark") && strings.Contains(lower, "not found") {
		return "📡 tshark is not installed",
			"Please install Wireshark (includes tshark):\n• macOS: 'brew install wireshark'\n• Linux: 'sudo apt-get install tshark' or 'sudo yum install wireshark'"
	}

	// Cluster creation errors
	if strings.Contains(lower, "creating cluster") && strings.Contains(lower, "failed") {
		return "☸️  Failed to create kind cluster",
			"This usually means:\n• Docker is not running properly\n• Insufficient system resources\n• Port conflicts\n\nTry:\n• Restart Docker\n• Free up disk space\n• Run 'kind delete cluster --name netlab-osi' to clean up"
	}

	// Network/connectivity errors
	if strings.Contains(lower, "network") && (strings.Contains(lower, "unreachable") || strings.Contains(lower, "timeout")) {
		return "🌐 Network connectivity issue",
			"This usually means:\n• No internet connection for downloading images\n• Corporate firewall blocking Docker registry\n• DNS resolution issues\n\nTry:\n• Check internet connectivity\n• Configure Docker to use corporate proxy if needed"
	}

	// Permission errors
	if strings.Contains(lower, "permission denied") && !strings.Contains(lower, "docker") {
		return "🔒 Permission denied",
			"This usually means:\n• Script file is not executable\n• Insufficient privileges\n\nTry:\n• Run 'chmod +x ./scripts/k8s_lab.sh'\n• Check if you need sudo for certain operations"
	}

	// Resource errors
	if strings.Contains(lower, "no space left") || strings.Contains(lower, "disk full") {
		return "💾 Insufficient disk space",
			"Please:\n• Free up disk space (at least 2GB recommended)\n• Clean up Docker: 'docker system prune -a'\n• Remove unused kind clusters: 'kind get clusters' then 'kind delete cluster --name <cluster>'"
	}

	// Generic fallback
	if strings.Contains(output, "exit status") || strings.Contains(err.Error(), "exit status") {
		return "❌ Lab setup script failed",
			"The setup script encountered an error. Common solutions:\n• Run 'netlab doctor' to check all dependencies\n• Try running './scripts/k8s_lab.sh setup' manually for detailed output\n• Make sure Docker is running and you have sufficient resources\n• Check the raw output below for specific error details"
	}

	// Last resort
	return fmt.Sprintf("❌ Lab setup failed: %v", err),
		"Try:\n• Run 'netlab doctor' to check dependencies\n• Run './scripts/k8s_lab.sh setup' manually for more details\n• Check the raw output below for error specifics"
}

type labStatusMsg struct {
	ready bool
	error string
}

type labSetupStartMsg struct{}

type labCleanupStartMsg struct{}

type labOutputMsg struct {
	output   string
	progress float64
	finished bool
	error    string
	success  bool
}

type logExportMsg struct {
	success bool
	file    string
	error   string
}

// getSamplePacketLayers returns sample packet data for demonstration
func getSamplePacketLayers() []PacketLayer {
	return []PacketLayer{
		{
			OSILayer: 1,
			Name:     "Physical Layer",
			Headers: map[string]string{
				"Medium":       "Ethernet over copper wire",
				"Encoding":     "Manchester encoding",
				"Signal Level": "-2.5V to +2.5V",
				"Bit Rate":     "1000 Mbps",
			},
			RawData:     "10101010 10101010 10101010 10101010 10111011 ...",
			Explanation: "At the physical layer, data is transmitted as electrical signals over the network cable. This represents the raw bits being sent as voltage levels on the wire.",
		},
		{
			OSILayer: 2,
			Name:     "Data Link Layer (Ethernet)",
			Headers: map[string]string{
				"Destination MAC": "02:42:ac:12:00:02",
				"Source MAC":      "02:42:ac:12:00:01",
				"EtherType":       "0x0800 (IPv4)",
				"Frame Length":    "74 bytes",
				"FCS":             "0x12345678",
			},
			RawData:     "02:42:ac:12:00:02 02:42:ac:12:00:01 08:00 45:00...",
			Explanation: "The Ethernet frame wraps the IP packet with MAC addresses for local network delivery. The source MAC is the sending container's interface, and the destination MAC is the nginx container's interface.",
		},
		{
			OSILayer: 3,
			Name:     "Network Layer (IP)",
			Headers: map[string]string{
				"Version":        "4 (IPv4)",
				"Header Length":  "20 bytes",
				"Source IP":      "10.244.0.5",
				"Destination IP": "10.244.0.10",
				"Protocol":       "6 (TCP)",
				"TTL":            "64",
				"Packet Length":  "60 bytes",
			},
			RawData:     "45:00:00:3c:00:00:40:00:40:06:b7:c8:0a:f4:00:05:0a:f4:00:0a",
			Explanation: "The IP header contains routing information to deliver the packet from the busybox Pod IP to the nginx Pod IP within the Kubernetes cluster network.",
		},
		{
			OSILayer: 4,
			Name:     "Transport Layer (TCP)",
			Headers: map[string]string{
				"Source Port":      "38472",
				"Destination Port": "80",
				"Sequence Number":  "1234567890",
				"Ack Number":       "0",
				"Flags":            "SYN",
				"Window Size":      "65535",
				"Checksum":         "0x1234",
			},
			RawData:     "96:38:00:50:49:96:02:d2:00:00:00:00:a0:02:ff:ff:12:34:00:00",
			Explanation: "The TCP header establishes a reliable connection to nginx on port 80. This is the SYN packet that starts the TCP three-way handshake for the HTTP connection.",
		},
		{
			OSILayer: 7,
			Name:     "Application Layer (HTTP)",
			Headers: map[string]string{
				"Method":       "GET",
				"URI":          "/",
				"HTTP Version": "1.1",
				"Host":         "nginx",
				"User-Agent":   "curl/7.64.0",
				"Accept":       "*/*",
				"Connection":   "keep-alive",
			},
			RawData:     "GET / HTTP/1.1\\r\\nHost: nginx\\r\\nUser-Agent: curl/7.64.0\\r\\nAccept: */*\\r\\n\\r\\n",
			Explanation: "The HTTP GET request from the busybox Pod to fetch the nginx welcome page. This is the application-layer data that users actually care about - a web request.",
		},
	}
}

// RunWalkthrough starts the packet analysis walkthrough
func RunWalkthrough() error {
	m := NewWalkthroughModel()

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}

// LoadPacketData loads parsed packet data from the lab files
func LoadPacketData() ([]PacketLayer, error) {
	// Check if the packet capture file exists
	pcapPath := filepath.Join("modules", "01-osi-model", "assets", "https-nginx.pcap")
	if _, err := os.Stat(pcapPath); os.IsNotExist(err) {
		// Return sample data if no lab data is available
		return getSamplePacketLayers(), nil
	}

	// TODO: Parse actual packet data from pcap file
	// For now, return sample data
	return getSamplePacketLayers(), nil
}

// updateOutputViewport updates the output viewport with current lab output
func (m *WalkthroughModel) updateOutputViewport() {
	if len(m.labOutput) > 0 {
		content := strings.Join(m.labOutput, "\n")
		if m.outputViewport.Width > 0 {
			m.outputViewport.SetContent(content)
			m.outputViewport.GotoBottom()
		}
	}
}

// labSetupView renders the lab setup view with progress and output
func (m WalkthroughModel) labSetupView() string {
	// Check for minimum terminal size first
	if m.width < 60 || m.height < 20 {
		return fmt.Sprintf("⚠️  Terminal too small (%dx%d), please resize to at least 60x20", m.width, m.height)
	}

	header := m.headerView()
	footer := m.footerView()

	// Progress bar
	progressPercent := fmt.Sprintf("%.0f%%", m.labProgress)
	progressBar := strings.Repeat("█", int(m.labProgress/4)) + strings.Repeat("░", 25-int(m.labProgress/4))
	progressDisplay := fmt.Sprintf("[%s] %s", progressBar, progressPercent)

	progressStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Width(m.width - 4).
		Align(lipgloss.Center)

	// Title and status
	var title, status string
	if m.labError != "" {
		title = styles.H2.Render("❌ Lab Setup Failed")
		status = styles.StatusError.Render("Lab setup encountered an error. Check output below and try again.")
	} else if m.labRunning {
		title = styles.H2.Render("🚀 Lab Setup in Progress")
		status = styles.BodyMuted.Render("Setting up kind cluster, deploying nginx, and capturing packets...")
	} else {
		title = styles.H2.Render("✅ Lab Setup Complete")
		status = styles.StatusSuccess.Render("Lab setup completed successfully!")
	}

	// Output viewport
	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Width(m.width - 8).
		Height(m.height - 15). // Leave space for header, progress, footer
		Padding(1)

	outputTitle := styles.H3.Render("📄 Live Output")
	copyHint := styles.BodyMuted.Render("💡 Tip: Press 'e' to export logs to a file for copying")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		status,
		"",
		progressStyle.Render(progressDisplay),
		"",
		outputTitle,
		copyHint,
		outputStyle.Render(m.outputViewport.View()),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

func (m WalkthroughModel) exportLogs() tea.Cmd {
	return func() tea.Msg {
		// Create logs directory if it doesn't exist
		if err := os.MkdirAll("logs", 0755); err != nil {
			return logExportMsg{
				success: false,
				error:   fmt.Sprintf("Failed to create logs directory: %s", err.Error()),
			}
		}

		// Generate filename with timestamp
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("logs/netlab-lab-output_%s.txt", timestamp)

		// Write logs to file
		content := strings.Join(m.labOutput, "\n")
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return logExportMsg{
				success: false,
				error:   fmt.Sprintf("Failed to write log file: %s", err.Error()),
			}
		}

		return logExportMsg{
			success: true,
			file:    filename,
		}
	}
}
