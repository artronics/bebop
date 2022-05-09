package cmd

import (
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Creating, managing and maintaining a project.",
	Long: `A proxy project is following a standard template. This command provides 
all the tools necessary to initialise and maintain certain aspects of a project.`,
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
