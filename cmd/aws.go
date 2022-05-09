package cmd

import (
	"github.com/spf13/cobra"
)

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Managing commands related to aws. It requires aws-cli as a dependency",
	Long:  `You can use this command to login and TODO:`,
}

func init() {
	awsCmd.PersistentFlags().StringP("profile", "p", "default", `your aws-cli profile. 
This will overwrite AWS_PROFILE environment variable if it's set explicitly. Otherwise, AWS_PROFILE will be used and if
neither values are present, then "default" will be used`)

	rootCmd.AddCommand(awsCmd)
}
