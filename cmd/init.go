package cmd

import (
	"github.com/nhsdigital/bebop-cli/pkg"
	"github.com/spf13/cobra"
	"log"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new project",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sd, err := setTemplateSourceData(cmd)

		err = pkg.Template(sd)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	initCmd.Flags().String("template", "https://github.com/artronics/bebop-proto", "git repo url")

	initCmd.Flags().String("out", ".", "generated project output path. Default is current directory")
	err := initCmd.MarkFlagDirname("out")
	if err != nil {
		log.Fatal(err.Error())
	}

	initCmd.Flags().String("var-file", "", "path to either json or yaml file containing key value pairs")
	_ = initCmd.MarkFlagRequired("var-file")
	err = initCmd.MarkFlagFilename("var-file")
	if err != nil {
		log.Fatal(err.Error())
	}

	initCmd.Flags().StringSlice("exc-dir", []string{".git", "node_modules", ".idea", ".vscode"}, "comma separated list of directories to exclude")
	initCmd.Flags().StringSlice("exc-file-ext", []string{".zip", ".exe", ".tar", ".tar.gz", ".jar"}, "comma separated list of file extensions to exclude")

	projectCmd.AddCommand(initCmd)
}
