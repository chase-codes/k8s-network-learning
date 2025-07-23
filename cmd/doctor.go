package cmd

import (
	"netlab/internal/utils"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run environment diagnostics",
	Long:  "Check system requirements and validate that all necessary tools are installed and configured correctly.",
	Run: func(cmd *cobra.Command, args []string) {
		utils.RunDiagnostics()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
