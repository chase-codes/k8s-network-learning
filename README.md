# NetLab 🚀

**Professional Interactive Terminal-Based Learning Environment for Networking Fundamentals**

NetLab is a modern, TUI-based learning platform that teaches networking concepts and Kubernetes networking through interactive terminal modules. Built with Go and the Charm ecosystem for a fast, responsive, and beautiful CLI experience.

## ✨ Features

- **🎯 Interactive Learning**: Hands-on modules with immediate feedback
- **🎨 Professional TUI Design**: Beautiful, consistent interface with custom styling system
- **📚 Progressive Curriculum**: From basic networking to advanced Kubernetes concepts
- **🔧 Environment Validation**: Built-in diagnostics for required tools
- **⚡ Fast & Efficient**: Native Go performance with minimal resource usage
- **🔄 Modular Design**: Self-contained lessons that build upon each other
- **📱 Responsive**: Adapts to different terminal sizes gracefully
- **🎭 Themed Interface**: Consistent color palette and typography across all modules

## 🎨 Design System

NetLab features a comprehensive design system with:

- **Professional Color Palette**: Cyan primary, purple secondary, amber accents
- **Consistent Typography**: H1/H2/H3 headers, body text, code formatting
- **Status Indicators**: Ready ✅, Work-in-Progress 🚧, Planned 📋
- **Interactive Elements**: Styled buttons, lists, and navigation
- **Visual Hierarchy**: Cards, panels, sections with proper spacing
- **Accessibility**: Color contrast, keyboard navigation, clear feedback

See the complete [Style Guide](docs/style-guide.md) for development standards.

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** (required)
- **Docker** (optional, for advanced modules)
- **kubectl** (optional, for Kubernetes modules)

### Installation

```bash
# Clone the repository
git clone https://github.com/your-username/netlab.git
cd netlab

# Run setup validation
make setup

# Build and start NetLab
make start
```

### Alternative: Direct Run

```bash
# Run without building
go run . start

# Or run a specific module
go run . module 01-osi-model
```

## 📋 Available Commands

```bash
# Main commands
netlab start              # Launch interactive module menu
netlab module <id>        # Jump to specific module
netlab doctor             # Run environment diagnostics
netlab --help             # Show help and options

# Development commands (via Makefile)
make build               # Build binary to bin/netlab
make run                 # Run in development mode
make test                # Run tests
make fmt                 # Format code
make doctor              # Run diagnostics
make clean               # Clean build artifacts
make install             # Install to /usr/local/bin
```

## 📚 Learning Modules

### Current Modules

| Module | Topic | Status | Prerequisites |
|--------|-------|---------|---------------|
| `01-osi-model` | OSI Model Fundamentals | ✅ **ENHANCED** | Basic networking knowledge |
| `02-tcp-ip` | TCP/IP Stack Deep Dive | 📋 Planned | OSI Model |
| `03-subnetting` | Subnetting and CIDR | 📋 Planned | TCP/IP basics |
| `04-routing` | Routing Protocols | 📋 Planned | Subnetting |
| `05-k8s-networking` | Kubernetes Networking | 📋 Planned | Basic K8s knowledge |
| `06-cni` | Container Network Interface | 📋 Planned | K8s networking |
| `07-service-mesh` | Service Mesh Concepts | 📋 Planned | Advanced K8s |

### Enhanced Module Features

The OSI Model module now includes:
- **Professional styling** with consistent theming
- **Rich visual elements** including ASCII diagrams and status indicators
- **Interactive content** with progress tracking and smooth navigation
- **Comprehensive coverage** from basic concepts to Kubernetes networking context
- **Knowledge checks** and real-world examples
- **Responsive design** that adapts to your terminal size

### Module Structure

Each module includes:
- **Interactive TUI**: Navigate with keyboard controls
- **Visual Diagrams**: ASCII art and charts with consistent styling
- **Practical Examples**: Real-world scenarios in styled containers
- **Knowledge Checks**: Practice questions in highlighted quiz sections
- **Reference Materials**: Quick lookup guides with proper typography

## 🛠️ Development

### Project Structure

```
netlab/
├── cmd/                # CLI commands (cobra)
│   ├── root.go        # Root command
│   ├── start.go       # Start TUI
│   ├── module.go      # Module runner
│   └── doctor.go      # Diagnostics
├── internal/
│   ├── tui/           # TUI components
│   │   ├── welcome.go           # Basic welcome screen
│   │   └── welcome_enhanced.go  # Enhanced welcome screen
│   ├── modules/       # Module management
│   │   ├── runner.go  # Module dispatcher
│   │   └── osimodel/  # OSI Model learning module
│   │       ├── module.go    # Original implementation
│   │       └── enhanced.go  # Enhanced styled version
│   └── utils/         # Utilities & diagnostics
├── pkg/               # Shared components
│   ├── styles/        # Design system and theming
│   │   └── theme.go   # Complete style definitions
│   └── components/    # Reusable UI components
│       └── logo.go    # Logo and header components
├── modules/           # Learning content
│   └── 01-osi-model/ # Example module with README
├── docs/              # Documentation
│   └── style-guide.md # Complete style guide for developers
├── assets/            # Static assets
├── scripts/           # Setup & validation tools
├── go.mod            # Go dependencies
├── Makefile          # Build automation
└── README.md         # This file
```

### Tech Stack

- **Go 1.21+**: Core language
- **Bubble Tea**: TUI framework
- **Lip Gloss**: Styling and layout with custom design system
- **Bubbles**: Pre-built UI components
- **Cobra**: CLI framework
- **Custom Styling**: Professional design system with consistent theming

### Style System

NetLab includes a comprehensive styling system:

```go
// Use the design system
import "netlab/pkg/styles"

// Headers
content := styles.H1.Render("Module Title")
content += styles.H2.Render("Section Header")

// Content containers
example := styles.ModuleExample.Render(codeBlock)
section := styles.ModuleSection.Render(explanation)
quiz := styles.ModuleQuiz.Render(questions)

// Status indicators
status := styles.StatusSuccess.Render("✅ READY")
```

### Adding New Modules

1. Create module directory: `modules/XX-topic-name/`
2. Add module implementation in `internal/modules/`
3. **Follow the style guide**: Use consistent colors, typography, and layouts
4. Update the module runner in `internal/modules/runner.go`
5. Add module to the welcome screen list
6. Create README.md with learning objectives

See [`docs/style-guide.md`](docs/style-guide.md) for complete development standards.

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-module`
3. **Follow the style guide** for consistent UI/UX
4. Make your changes and test: `make test`
5. Format code: `make fmt`
6. Commit and push your changes
7. Create a Pull Request

## 🔧 System Requirements

### Required
- **Go 1.21+**: For building and running NetLab
- **Terminal**: Modern terminal with Unicode support (recommended: iTerm2, Windows Terminal)

### Optional (for advanced modules)
- **Docker**: Container networking experiments
- **kubectl**: Kubernetes cluster interaction
- **kind**: Local Kubernetes clusters
- **tcpdump**: Network packet analysis
- **ip/iptables**: Linux networking tools (Linux only)

Run `netlab doctor` or `make setup` to validate your environment.

## 📖 Usage Examples

### Basic Usage

```bash
# Start the interactive menu
netlab start

# Run a specific module
netlab module 01-osi-model

# Check your environment
netlab doctor
```

### Navigation

Within the application:
- **↑/↓ arrows**: Navigate lists and scroll content
- **Page Up/Down**: Fast scroll through module content
- **Enter**: Select items and activate modules
- **q/Ctrl+C/Esc**: Exit to previous screen or quit
- **Mouse support**: Scroll with mouse wheel (where supported)

### Development Workflow

```bash
# Full development cycle
make clean          # Clean previous builds
make deps           # Download dependencies  
make fmt            # Format code
make test           # Run tests
make build          # Build binary
make start          # Launch NetLab

# Quick iteration
make run            # Run without building
```

## 🎯 Learning Path

### Beginner Path
1. **OSI Model** (`01-osi-model`) - ✅ **Enhanced** - Fundamental network layers
2. **TCP/IP** (`02-tcp-ip`) - Internet protocol deep dive
3. **Subnetting** (`03-subnetting`) - Network segmentation

### Intermediate Path
4. **Routing** (`04-routing`) - How packets find their way
5. **Kubernetes Networking** (`05-k8s-networking`) - Container networking basics

### Advanced Path
6. **CNI** (`06-cni`) - Container Network Interface details
7. **Service Mesh** (`07-service-mesh`) - Advanced traffic management

## ❓ Troubleshooting

### Common Issues

**Build fails with missing dependencies:**
```bash
go mod download
go mod tidy
```

**TUI display issues:**
- Ensure terminal supports Unicode and colors
- Try different terminal (iTerm2, Windows Terminal, etc.)
- Check terminal size (minimum 80x24 recommended)
- Verify color support: `echo $COLORTERM`

**Permission denied on scripts:**
```bash
chmod +x scripts/setup.sh
```

### Getting Help

- Run `netlab --help` for command information
- Run `netlab doctor` for environment diagnostics
- Check module READMEs for specific guidance
- Review the [Style Guide](docs/style-guide.md) for development standards
- Open an issue on GitHub for bugs or questions

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) team for the amazing TUI toolkit
- [Cobra](https://cobra.dev/) for CLI framework
- The Go community for excellent tooling and libraries

---

**Happy Learning!** 🎓 Start your networking journey with `netlab start`

**Professional. Beautiful. Educational.** ✨
