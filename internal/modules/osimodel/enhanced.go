package osimodel

import (
	"fmt"
	"strings"

	"netlab/pkg/components"
	"netlab/pkg/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type enhancedOSIModel struct {
	content  string
	ready    bool
	viewport viewport.Model
	width    int
	height   int
}

func (m enhancedOSIModel) Init() tea.Cmd {
	return nil
}

func (m enhancedOSIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m enhancedOSIModel) View() string {
	if !m.ready {
		return styles.Body.Render("\n  Initializing OSI Model module...")
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m enhancedOSIModel) headerView() string {
	// Compact logo and module title
	logo := components.RenderCompactLogo()
	title := styles.ModuleTitle.
		Width(m.width - 4).
		Align(lipgloss.Center).
		Render("OSI Model - Seven Layers of Network Communication")

	breadcrumb := styles.Breadcrumb.Render("NetLab > Fundamentals > OSI Model")

	return lipgloss.JoinVertical(lipgloss.Left, logo, title, breadcrumb)
}

func (m enhancedOSIModel) footerView() string {
	// Progress indicator
	progress := fmt.Sprintf("%.0f%%", m.viewport.ScrollPercent()*100)
	progressText := styles.BodyMuted.Render(fmt.Sprintf("Progress: %s", progress))

	// Help keys
	helpKeys := []string{
		styles.KeyBinding.Render("↑/↓") + " scroll",
		styles.KeyBinding.Render("PgUp/PgDn") + " page",
		styles.KeyBinding.Render("q") + " back to menu",
	}
	helpText := styles.Help.Render(strings.Join(helpKeys, " • "))

	// Footer line
	line := strings.Repeat("─", maxInt(0, m.width-2))
	separator := styles.BodyDim.Render(line)

	// Create footer with progress on left, help on right
	footerContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		progressText,
		strings.Repeat(" ", maxInt(0, m.width-lipgloss.Width(progressText)-lipgloss.Width(helpText))),
		helpText,
	)

	footer := lipgloss.JoinVertical(
		lipgloss.Left,
		separator,
		footerContent,
	)

	return lipgloss.Place(m.width, lipgloss.Height(footer), lipgloss.Center, lipgloss.Top, footer)
}

// RunEnhanced starts the enhanced OSI Model learning module
func RunEnhanced() error {
	content := getEnhancedOSIContent()

	m := enhancedOSIModel{
		content: content,
		ready:   false,
	}

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func getEnhancedOSIContent() string {
	var content strings.Builder

	// Introduction
	content.WriteString(styles.H1.Render("🌐 The OSI Model"))
	content.WriteString("\n\n")
	content.WriteString(styles.Body.Render("The Open Systems Interconnection (OSI) model is a conceptual framework that standardizes the functions of a telecommunication or computing system into seven abstraction layers."))
	content.WriteString("\n\n")

	// Layer 7
	content.WriteString(styles.H2.Render("Layer 7: Application Layer"))
	content.WriteString("\n")
	layerBox := `┌─────────────────────────────────┐
│        APPLICATION LAYER        │  ← HTTP, FTP, SMTP, DNS
│         (Layer 7)               │    User Interface
└─────────────────────────────────┘`
	content.WriteString(styles.ModuleExample.Render(layerBox))
	content.WriteString("\n")
	content.WriteString(styles.ModuleSection.Render("The application layer provides network services directly to end-users. This is where human-computer interaction happens through network-aware applications.\n\nExamples: Web browsers (HTTP/HTTPS), Email clients (SMTP/POP3/IMAP), File transfer (FTP), Domain Name System (DNS)"))
	content.WriteString("\n\n")

	// Layer 6
	content.WriteString(styles.H2.Render("Layer 6: Presentation Layer"))
	content.WriteString("\n")
	layerBox = `┌─────────────────────────────────┐
│       PRESENTATION LAYER        │  ← SSL/TLS, JPEG, MPEG
│         (Layer 6)               │    Encryption & Compression
└─────────────────────────────────┘`
	content.WriteString(styles.ModuleExample.Render(layerBox))
	content.WriteString("\n")
	content.WriteString(styles.ModuleSection.Render("The presentation layer transforms data between the application layer and the network. It handles encryption, compression, and data formatting.\n\nExamples: SSL/TLS encryption, JPEG/PNG image formats, MPEG video compression"))
	content.WriteString("\n\n")

	// Layer 5
	content.WriteString(styles.H2.Render("Layer 5: Session Layer"))
	content.WriteString("\n")
	layerBox = `┌─────────────────────────────────┐
│         SESSION LAYER           │  ← NetBIOS, RPC, SQL
│         (Layer 5)               │    Session Management
└─────────────────────────────────┘`
	content.WriteString(styles.ModuleExample.Render(layerBox))
	content.WriteString("\n")
	content.WriteString(styles.ModuleSection.Render("The session layer manages connections (sessions) between applications. It establishes, maintains, and terminates connections.\n\nExamples: NetBIOS, Remote Procedure Calls (RPC), SQL sessions"))
	content.WriteString("\n\n")

	// Layer 4
	content.WriteString(styles.H2.Render("Layer 4: Transport Layer"))
	content.WriteString("\n")
	layerBox = `┌─────────────────────────────────┐
│        TRANSPORT LAYER          │  ← TCP, UDP
│         (Layer 4)               │    End-to-End Delivery
└─────────────────────────────────┘`
	content.WriteString(styles.ModuleExample.Render(layerBox))
	content.WriteString("\n")
	content.WriteString(styles.ModuleSection.Render("The transport layer provides reliable (TCP) or unreliable (UDP) delivery of data between end systems. It handles error detection, flow control, and segmentation.\n\nExamples: TCP (reliable), UDP (fast but unreliable)"))
	content.WriteString("\n")
	concepts := `Key Concepts:
• Port numbers (80 for HTTP, 443 for HTTPS, 22 for SSH)
• Segmentation and reassembly
• Flow control and error recovery`
	content.WriteString(styles.ModuleExample.Render(concepts))
	content.WriteString("\n\n")

	// Layer 3
	content.WriteString(styles.H2.Render("Layer 3: Network Layer"))
	content.WriteString("\n")
	layerBox = `┌─────────────────────────────────┐
│         NETWORK LAYER           │  ← IP, ICMP, OSPF
│         (Layer 3)               │    Routing & Addressing
└─────────────────────────────────┘`
	content.WriteString(styles.ModuleExample.Render(layerBox))
	content.WriteString("\n")
	content.WriteString(styles.ModuleSection.Render("The network layer handles routing of data between different networks. It determines the best path for data to travel across multiple networks.\n\nExamples: IP (Internet Protocol), ICMP (ping), routing protocols (OSPF, BGP)"))
	content.WriteString("\n")
	concepts = `Key Concepts:
• IP addresses (IPv4: 192.168.1.1, IPv6: 2001:db8::1)
• Routing tables
• Subnets and VLANs`
	content.WriteString(styles.ModuleExample.Render(concepts))
	content.WriteString("\n\n")

	// Layer 2
	content.WriteString(styles.H2.Render("Layer 2: Data Link Layer"))
	content.WriteString("\n")
	layerBox = `┌─────────────────────────────────┐
│        DATA LINK LAYER          │  ← Ethernet, WiFi, PPP
│         (Layer 2)               │    Node-to-Node Delivery
└─────────────────────────────────┘`
	content.WriteString(styles.ModuleExample.Render(layerBox))
	content.WriteString("\n")
	content.WriteString(styles.ModuleSection.Render("The data link layer provides node-to-node delivery of data over a physical link. It handles error detection and correction for the physical layer.\n\nExamples: Ethernet, Wi-Fi (802.11), PPP (Point-to-Point Protocol)"))
	content.WriteString("\n")
	concepts = `Key Concepts:
• MAC addresses (48-bit hardware addresses)
• Frame formatting
• Error detection (CRC)
• Switches operate at this layer`
	content.WriteString(styles.ModuleExample.Render(concepts))
	content.WriteString("\n\n")

	// Layer 1
	content.WriteString(styles.H2.Render("Layer 1: Physical Layer"))
	content.WriteString("\n")
	layerBox = `┌─────────────────────────────────┐
│         PHYSICAL LAYER          │  ← Cables, Radio, Fiber
│         (Layer 1)               │    Bits over Physical Media
└─────────────────────────────────┘`
	content.WriteString(styles.ModuleExample.Render(layerBox))
	content.WriteString("\n")
	content.WriteString(styles.ModuleSection.Render("The physical layer defines the electrical, mechanical, and procedural interface to the physical transmission medium.\n\nExamples: Ethernet cables (Cat5e, Cat6), Fiber optic cables, Radio frequencies (Wi-Fi, Bluetooth)"))
	content.WriteString("\n")
	concepts = `Key Concepts:
• Electrical signals
• Cable specifications
• Connectors and ports
• Hubs and repeaters`
	content.WriteString(styles.ModuleExample.Render(concepts))
	content.WriteString("\n\n")

	// Memory Aid
	content.WriteString(styles.H2.Render("🧠 Memory Aid"))
	content.WriteString("\n")
	memoryAid := `"Please Do Not Throw Sausage Pizza Away"

Physical      (Layer 1)
Data Link     (Layer 2)
Network       (Layer 3)
Transport     (Layer 4)
Session       (Layer 5)
Presentation  (Layer 6)
Application   (Layer 7)`
	content.WriteString(styles.Highlight.Render(memoryAid))
	content.WriteString("\n\n")

	// Real-world example
	content.WriteString(styles.H2.Render("🌍 Real-World Example: Web Browsing"))
	content.WriteString("\n")
	example := `When you visit https://example.com:

7. Application:   Browser sends HTTP request
6. Presentation:  SSL/TLS encrypts the data
5. Session:       Session established between browser and server
4. Transport:     TCP ensures reliable delivery (port 443)
3. Network:       IP routes packets across the internet
2. Data Link:     Ethernet frames carry data between local devices
1. Physical:      Electrical signals travel over cables/wireless`
	content.WriteString(styles.ModuleExample.Render(example))
	content.WriteString("\n\n")

	// Quiz section
	content.WriteString(styles.H2.Render("🧪 Knowledge Check"))
	content.WriteString("\n")
	quiz := `1. Which layer handles IP addresses?
   Answer: Layer 3 (Network Layer)

2. What layer does a switch operate at?
   Answer: Layer 2 (Data Link Layer)

3. Which layer is responsible for encryption?
   Answer: Layer 6 (Presentation Layer)

4. What layer uses port numbers?
   Answer: Layer 4 (Transport Layer)`
	content.WriteString(styles.ModuleQuiz.Render(quiz))
	content.WriteString("\n\n")

	// Kubernetes context
	content.WriteString(styles.H2.Render("☸️ Kubernetes Networking Context"))
	content.WriteString("\n")
	k8sContext := `Understanding the OSI model is crucial for Kubernetes networking:

• Container networking involves layers 2-4
• Service discovery operates at layers 3-7
• Ingress controllers work primarily at layer 7
• Network policies function at layers 3-4

This foundation will help you understand how Kubernetes implements networking concepts across these layers.`
	content.WriteString(styles.ModuleSection.Render(k8sContext))
	content.WriteString("\n\n")

	// Next steps
	content.WriteString(styles.H2.Render("➡️ Next Steps"))
	content.WriteString("\n")
	nextSteps := `Continue your learning journey:

🔍 netlab module 02-tcp-ip
   Deep dive into the TCP/IP protocol stack

🏠 netlab start
   Return to the main module menu`
	content.WriteString(styles.Body.Render(nextSteps))
	content.WriteString("\n\n")

	return content.String()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
