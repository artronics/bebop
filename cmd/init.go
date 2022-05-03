package cmd

import (
	"github.com/nhsdigital/bebop-cli/template"
	"github.com/spf13/cobra"
	"log"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new project",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sd := template.SourceData{}
		sd.Url = cmd.Flags().Lookup("template").Value.String()

		err := template.Template(sd)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	initCmd.Flags().String("template", "https://github.com/artronics/bebop-proto", "git repo url")
	projectCmd.AddCommand(initCmd)
}
