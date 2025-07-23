package cmd

import (
	"fmt"
	"log"

	"netlab/internal/modules"
	"netlab/internal/tui"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch the NetLab TUI welcome screen",
	Long:  "Start the interactive NetLab terminal interface with module selection menu.",
	Run: func(cmd *cobra.Command, args []string) {
		moduleID, err := tui.StartEnhancedWelcome()
		if err != nil {
			log.Fatal(fmt.Errorf("failed to start NetLab TUI: %w", err))
		}

		// If a module was selected, launch it with dependency checking
		if moduleID != "" {
			if err := modules.RunModuleWithDependencyCheck(moduleID); err != nil {
				log.Fatal(fmt.Errorf("failed to run module %s: %w", moduleID, err))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
