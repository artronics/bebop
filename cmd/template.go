package cmd

import (
	"github.com/nhsdigital/bebop-cli/internal"
	"github.com/nhsdigital/bebop-cli/pkg"
	"github.com/spf13/cobra"
	"log"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "template files in a directory given a data file in either json or yaml format",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		sd, err := setTemplateSourceData(cmd)
		sd.Url = "file://" + sd.Url

		err = pkg.Template(sd)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	// FIXME: there are duplicated logic between this and "project init". Check if we can get cmd and copy and modify it
	templateCmd.Flags().String("template", "", "input directory containing template files")
	_ = templateCmd.MarkFlagRequired("template")
	_ = templateCmd.MarkFlagDirname("template")

	templateCmd.Flags().String("out", ".", "generated project output path. Default is current directory")
	_ = templateCmd.MarkFlagDirname("out")

	templateCmd.Flags().String("var-file", "", "path to either json or yaml file containing key value pairs")
	_ = templateCmd.MarkFlagRequired("var-file")
	_ = templateCmd.MarkFlagFilename("var-file")

	templateCmd.Flags().StringSlice("exc-dir", []string{".git", "node_modules", ".idea", ".vscode"}, "comma separated list of directories to exclude")
	templateCmd.Flags().StringSlice("exc-file-ext", []string{".zip", ".exe", ".tar", ".tar.gz", ".jar"}, "comma separated list of file extensions to exclude")

	rootCmd.AddCommand(templateCmd)
}

func setTemplateSourceData(cmd *cobra.Command) (pkg.SourceData, error) {
	sd := pkg.SourceData{}

	sd.Url = cmd.Flags().Lookup("template").Value.String()
	sd.OutputDir = cmd.Flags().Lookup("out").Value.String()

	excDir, err := cmd.Flags().GetStringSlice("exc-dir")
	if err != nil {
		log.Fatal(err.Error())
	}
	sd.ExcludedDirs = excDir

	excFilesExt, err := cmd.Flags().GetStringSlice("exc-file-ext")
	if err != nil {
		log.Fatal(err.Error())
	}
	sd.ExcludedFiles = excFilesExt

	templateData, err := internal.ParseDataFile(cmd.Flags().Lookup("var-file").Value.String())
	if err != nil {
		log.Fatal(err.Error())
	}
	sd.TemplateData = templateData

	return sd, nil
}
