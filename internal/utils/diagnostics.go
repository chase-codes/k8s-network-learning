package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	checkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true) // Green
	warnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true) // Yellow
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)  // Red
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true) // Blue
)

type diagnostic struct {
	name        string
	command     string
	args        []string
	required    bool
	description string
	installCmd  map[string]string // OS -> install command
}

var diagnostics = []diagnostic{
	{
		name:        "Go",
		command:     "go",
		args:        []string{"version"},
		required:    true,
		description: "Go programming language (required for building NetLab)",
		installCmd: map[string]string{
			"darwin": "brew install go",
			"linux":  "Visit https://golang.org/doc/install",
		},
	},
	{
		name:        "Docker",
		command:     "docker",
		args:        []string{"--version"},
		required:    false,
		description: "Docker for containerized network experiments",
		installCmd: map[string]string{
			"darwin": "brew install --cask docker",
			"linux":  "curl -fsSL https://get.docker.com | sh",
		},
	},
	{
		name:        "kubectl",
		command:     "kubectl",
		args:        []string{"version", "--client"},
		required:    false,
		description: "Kubernetes CLI for cluster networking modules",
		installCmd: map[string]string{
			"darwin": "brew install kubectl",
			"linux":  "curl -LO \"https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl\"",
		},
	},
	{
		name:        "kind",
		command:     "kind",
		args:        []string{"version"},
		required:    false,
		description: "Kubernetes in Docker for local cluster setup",
		installCmd: map[string]string{
			"darwin": "brew install kind",
			"linux":  "curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64 && chmod +x ./kind && sudo mv ./kind /usr/local/bin/kind",
		},
	},
	{
		name:        "tcpdump",
		command:     "tcpdump",
		args:        []string{"--version"},
		required:    false,
		description: "Network packet analyzer for traffic inspection",
		installCmd: map[string]string{
			"darwin": "brew install tcpdump",
			"linux":  "sudo apt-get install tcpdump",
		},
	},
	{
		name:        "tshark",
		command:     "tshark",
		args:        []string{"--version"},
		required:    false,
		description: "Wireshark command-line packet analyzer for OSI lab",
		installCmd: map[string]string{
			"darwin": "brew install wireshark",
			"linux":  "sudo apt-get install tshark",
		},
	},
	{
		name:        "ip",
		command:     "ip",
		args:        []string{"--version"},
		required:    false,
		description: "Network configuration tool (Linux)",
		installCmd: map[string]string{
			"linux": "sudo apt-get install iproute2",
		},
	},
	{
		name:        "iptables",
		command:     "iptables",
		args:        []string{"--version"},
		required:    false,
		description: "Firewall administration tool (Linux)",
		installCmd: map[string]string{
			"linux": "sudo apt-get install iptables",
		},
	},
}

// ModuleDependencies defines what each module actually needs
var ModuleDependencies = map[string][]string{
	"01-osi-model":      {"Docker", "kubectl", "kind", "tcpdump", "tshark"},
	"02-tcp-ip":         {"tcpdump", "tshark"},
	"03-subnetting":     {"ip"},
	"04-routing":        {"ip", "iptables"},
	"05-k8s-networking": {"Docker", "kubectl", "kind"},
	"06-cni":            {"Docker", "kubectl", "kind"},
	"07-service-mesh":   {"Docker", "kubectl", "kind"},
}

type DependencyStatus struct {
	Name       string
	Status     string // "ok", "missing", "error"
	Output     string
	Required   bool
	InstallCmd string
}

func RunDiagnostics() {
	fmt.Println(infoStyle.Render("🔍 NetLab Environment Diagnostics"))
	fmt.Println()

	allGood := true
	warnings := []string{}

	for _, diag := range diagnostics {
		status, output := checkCommand(diag.command, diag.args...)

		switch status {
		case "ok":
			fmt.Printf("%s %s: %s\n",
				checkStyle.Render("✓"),
				diag.name,
				strings.TrimSpace(output))
		case "missing":
			if diag.required {
				fmt.Printf("%s %s: %s\n",
					errorStyle.Render("✗"),
					diag.name,
					"REQUIRED - Not found")
				allGood = false
			} else {
				fmt.Printf("%s %s: %s\n",
					warnStyle.Render("⚠"),
					diag.name,
					"Optional - Not found")
				warnings = append(warnings, diag.name)
			}
		case "error":
			fmt.Printf("%s %s: %s\n",
				errorStyle.Render("✗"),
				diag.name,
				"Error running command")
			if diag.required {
				allGood = false
			}
		}
	}

	fmt.Println()

	if allGood && len(warnings) == 0 {
		fmt.Println(checkStyle.Render("🎉 All systems go! NetLab is ready to run."))
	} else if allGood {
		fmt.Println(checkStyle.Render("✅ Core requirements met!"))
		if len(warnings) > 0 {
			fmt.Println(warnStyle.Render(fmt.Sprintf("⚠️  Optional tools missing: %s", strings.Join(warnings, ", "))))
			fmt.Println("   Some advanced modules may have limited functionality.")
		}
	} else {
		fmt.Println(errorStyle.Render("❌ Missing required dependencies!"))
		fmt.Println("   Please install the missing tools and run 'netlab doctor' again.")
	}

	fmt.Println()
	fmt.Println("💡 Run 'scripts/setup.sh' for guided installation help.")
}

// CheckModuleDependencies checks dependencies specific to a module
func CheckModuleDependencies(moduleID string) ([]DependencyStatus, bool) {
	requiredDeps, exists := ModuleDependencies[moduleID]
	if !exists {
		return nil, true // No specific requirements
	}

	var results []DependencyStatus
	allGood := true

	// Create a map for quick lookup
	diagMap := make(map[string]diagnostic)
	for _, diag := range diagnostics {
		diagMap[diag.name] = diag
	}

	for _, depName := range requiredDeps {
		if diag, exists := diagMap[depName]; exists {
			status, output := checkCommand(diag.command, diag.args...)

			installCmd := ""
			if status == "missing" {
				// Determine OS and get install command
				os := getOS()
				if cmd, hasCmd := diag.installCmd[os]; hasCmd {
					installCmd = cmd
				}
			}

			result := DependencyStatus{
				Name:       diag.name,
				Status:     status,
				Output:     output,
				Required:   true, // All module deps are required for that module
				InstallCmd: installCmd,
			}
			results = append(results, result)

			if status != "ok" {
				allGood = false
			}
		}
	}

	return results, allGood
}

// GetInstallationGuide returns formatted installation instructions for missing dependencies
func GetInstallationGuide(missingDeps []DependencyStatus) string {
	if len(missingDeps) == 0 {
		return ""
	}

	var guide strings.Builder
	guide.WriteString(infoStyle.Render("📦 Installation Guide\n"))
	guide.WriteString("\n")

	os := getOS()
	guide.WriteString(fmt.Sprintf("Detected OS: %s\n\n", os))

	for _, dep := range missingDeps {
		if dep.Status == "missing" && dep.InstallCmd != "" {
			guide.WriteString(fmt.Sprintf("%s %s:\n",
				warnStyle.Render("•"), dep.Name))
			guide.WriteString(fmt.Sprintf("  %s\n\n", dep.InstallCmd))
		}
	}

	guide.WriteString("After installation, run the module again or use 'netlab doctor' to verify.\n")
	return guide.String()
}

func checkCommand(command string, args ...string) (string, string) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()

	if err != nil {
		// Check if it's a "command not found" error
		if strings.Contains(err.Error(), "executable file not found") ||
			strings.Contains(err.Error(), "command not found") {
			return "missing", ""
		}
		return "error", err.Error()
	}

	return "ok", string(output)
}

func getOS() string {
	// Simple OS detection - can be enhanced
	if cmd := exec.Command("uname", "-s"); cmd.Run() == nil {
		if output, err := cmd.Output(); err == nil {
			os := strings.ToLower(strings.TrimSpace(string(output)))
			if os == "darwin" {
				return "darwin"
			}
			if os == "linux" {
				return "linux"
			}
		}
	}
	return "unknown"
}
