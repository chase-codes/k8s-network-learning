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
}

var diagnostics = []diagnostic{
	{
		name:        "Go",
		command:     "go",
		args:        []string{"version"},
		required:    true,
		description: "Go programming language (required for building NetLab)",
	},
	{
		name:        "Docker",
		command:     "docker",
		args:        []string{"--version"},
		required:    false,
		description: "Docker for containerized network experiments",
	},
	{
		name:        "kubectl",
		command:     "kubectl",
		args:        []string{"version", "--client"},
		required:    false,
		description: "Kubernetes CLI for cluster networking modules",
	},
	{
		name:        "kind",
		command:     "kind",
		args:        []string{"version"},
		required:    false,
		description: "Kubernetes in Docker for local cluster setup",
	},
	{
		name:        "tcpdump",
		command:     "tcpdump",
		args:        []string{"--version"},
		required:    false,
		description: "Network packet analyzer for traffic inspection",
	},
	{
		name:        "ip",
		command:     "ip",
		args:        []string{"--version"},
		required:    false,
		description: "Network configuration tool (Linux)",
	},
	{
		name:        "iptables",
		command:     "iptables",
		args:        []string{"--version"},
		required:    false,
		description: "Firewall administration tool (Linux)",
	},
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
