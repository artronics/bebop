package cmd

import (
	"github.com/nhsdigital/bebop-cli/pkg"
	"github.com/spf13/cobra"
	"log"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new project",
	Long: `This command will create a new project based on a default template hosted on ???. You can provide 
a file containing ...TODO:`,
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
