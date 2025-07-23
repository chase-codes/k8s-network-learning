package cmd

import (
	"fmt"
	"log"

	"netlab/internal/modules"

	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
	Use:   "module [module-id]",
	Short: "Jump directly to a specific learning module",
	Long:  "Launch a specific NetLab learning module by its ID (e.g., 'netlab module 01-osi-model').",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleID := args[0]
		if err := modules.RunModuleWithDependencyCheck(moduleID); err != nil {
			log.Fatal(fmt.Errorf("failed to run module %s: %w", moduleID, err))
		}
	},
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}
