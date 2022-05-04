package cmd

import (
	"github.com/spf13/cobra"
)

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	awsCmd.PersistentFlags().StringP("profile", "p", "default", `your aws-cli profile. 
This will overwrite AWS_PROFILE environment variable if it's set explicitly. Otherwise, AWS_PROFILE will be used and if
neither values are present, then "default" will be used`)

	rootCmd.AddCommand(awsCmd)
}
