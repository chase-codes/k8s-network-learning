package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netlab",
	Short: "NetLab - Interactive Learning Environment for Networking Fundamentals",
	Long: `NetLab is a modern, TUI-based learning environment that teaches networking 
fundamentals and Kubernetes networking through interactive terminal modules.

NetLab emphasizes conceptual understanding, visual feedback, and hands-on learning—all 
within a fast, efficient CLI application built with the Charm ecosystem.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to NetLab! Use 'netlab start' to begin or 'netlab --help' for more options.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
}
