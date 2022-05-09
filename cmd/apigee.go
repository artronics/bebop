package cmd

import (
	"github.com/spf13/cobra"
)

var apigeeCmd = &cobra.Command{
	Use:   "apigee",
	Short: "Managing commands related to apigee.",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(apigeeCmd)
}
