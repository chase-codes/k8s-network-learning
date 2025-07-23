package osimodel

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Initialize viewport with window size
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
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

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := titleStyle.Render("OSI Model - Layer by Layer")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Run starts the OSI Model learning module
func Run() error {
	content := getOSIContent()

	// Initialize the model with content
	m := model{
		content: content,
		ready:   false,
	}

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func getOSIContent() string {
	return `# OSI Model - The Seven Layers of Network Communication

The Open Systems Interconnection (OSI) model is a conceptual framework that standardizes the functions of a telecommunication or computing system into seven abstraction layers.

## Layer 7: Application Layer
┌─────────────────────────────────┐
│        APPLICATION LAYER        │  ← HTTP, FTP, SMTP, DNS
│         (Layer 7)               │    User Interface
└─────────────────────────────────┘

The application layer provides network services directly to end-users. This is where human-computer interaction happens through network-aware applications.

Examples: Web browsers (HTTP/HTTPS), Email clients (SMTP/POP3/IMAP), File transfer (FTP), Domain Name System (DNS)

## Layer 6: Presentation Layer
┌─────────────────────────────────┐
│       PRESENTATION LAYER        │  ← SSL/TLS, JPEG, MPEG
│         (Layer 6)               │    Encryption & Compression
└─────────────────────────────────┘

The presentation layer transforms data between the application layer and the network. It handles encryption, compression, and data formatting.

Examples: SSL/TLS encryption, JPEG/PNG image formats, MPEG video compression

## Layer 5: Session Layer
┌─────────────────────────────────┐
│         SESSION LAYER           │  ← NetBIOS, RPC, SQL
│         (Layer 5)               │    Session Management
└─────────────────────────────────┘

The session layer manages connections (sessions) between applications. It establishes, maintains, and terminates connections.

Examples: NetBIOS, Remote Procedure Calls (RPC), SQL sessions

## Layer 4: Transport Layer
┌─────────────────────────────────┐
│        TRANSPORT LAYER          │  ← TCP, UDP
│         (Layer 4)               │    End-to-End Delivery
└─────────────────────────────────┘

The transport layer provides reliable (TCP) or unreliable (UDP) delivery of data between end systems. It handles error detection, flow control, and segmentation.

Examples: TCP (reliable), UDP (fast but unreliable)

Key Concepts:
- Port numbers (80 for HTTP, 443 for HTTPS, 22 for SSH)
- Segmentation and reassembly
- Flow control and error recovery

## Layer 3: Network Layer
┌─────────────────────────────────┐
│         NETWORK LAYER           │  ← IP, ICMP, OSPF
│         (Layer 3)               │    Routing & Addressing
└─────────────────────────────────┘

The network layer handles routing of data between different networks. It determines the best path for data to travel across multiple networks.

Examples: IP (Internet Protocol), ICMP (ping), routing protocols (OSPF, BGP)

Key Concepts:
- IP addresses (IPv4: 192.168.1.1, IPv6: 2001:db8::1)
- Routing tables
- Subnets and VLANs

## Layer 2: Data Link Layer
┌─────────────────────────────────┐
│        DATA LINK LAYER          │  ← Ethernet, WiFi, PPP
│         (Layer 2)               │    Node-to-Node Delivery
└─────────────────────────────────┘

The data link layer provides node-to-node delivery of data over a physical link. It handles error detection and correction for the physical layer.

Examples: Ethernet, Wi-Fi (802.11), PPP (Point-to-Point Protocol)

Key Concepts:
- MAC addresses (48-bit hardware addresses)
- Frame formatting
- Error detection (CRC)
- Switches operate at this layer

## Layer 1: Physical Layer
┌─────────────────────────────────┐
│         PHYSICAL LAYER          │  ← Cables, Radio, Fiber
│         (Layer 1)               │    Bits over Physical Media
└─────────────────────────────────┘

The physical layer defines the electrical, mechanical, and procedural interface to the physical transmission medium.

Examples: Ethernet cables (Cat5e, Cat6), Fiber optic cables, Radio frequencies (Wi-Fi, Bluetooth)

Key Concepts:
- Electrical signals
- Cable specifications
- Connectors and ports
- Hubs and repeaters

## Memory Aid: "Please Do Not Throw Sausage Pizza Away"
- Physical
- Data Link  
- Network
- Transport
- Session
- Presentation
- Application

## Real-World Example: Web Browsing

When you visit https://example.com:

1. **Application (7)**: Browser sends HTTP request
2. **Presentation (6)**: SSL/TLS encrypts the data
3. **Session (5)**: Session established between browser and server
4. **Transport (4)**: TCP ensures reliable delivery (port 443)
5. **Network (3)**: IP routes packets across the internet
6. **Data Link (2)**: Ethernet frames carry data between local devices
7. **Physical (1)**: Electrical signals travel over cables/wireless

## Practice Questions

1. Which layer handles IP addresses?
   Answer: Layer 3 (Network Layer)

2. What layer does a switch operate at?
   Answer: Layer 2 (Data Link Layer)

3. Which layer is responsible for encryption?
   Answer: Layer 6 (Presentation Layer)

4. What layer uses port numbers?
   Answer: Layer 4 (Transport Layer)

## Kubernetes Networking Context

Understanding the OSI model is crucial for Kubernetes networking:

- **Container networking** involves layers 2-4
- **Service discovery** operates at layers 3-7  
- **Ingress controllers** work primarily at layer 7
- **Network policies** function at layers 3-4

Next Module: TCP/IP Stack Deep Dive → 'netlab module 02-tcp-ip'

---
Press 'q' or 'Ctrl+C' to return to the main menu.
`
}
