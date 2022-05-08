package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nhsdigital/bebop-cli/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

		templateData, err := readTemplateData(cmd.Flags().Lookup("var-file").Value.String())
		if err != nil {
			log.Fatal(err.Error())
		}
		sd.TemplateData = templateData

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

func readTemplateData(path string) (map[string]interface{}, error) {
	var data map[string]interface{}

	file, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer func() {
		err = file.Close()
	}()
	bb, err := ioutil.ReadAll(file)
	if err != nil {
		return data, err
	}

	ext := filepath.Ext(file.Name())
	fmt.Println(ext)
	if ext == ".json" {
		err = json.Unmarshal(bb, &data)
		if err != nil {
			return data, err
		}

	} else if ext == ".yml" || ext == ".yaml" {
		err = yaml.Unmarshal(bb, &data)
		if err != nil {
			return data, err
		}

	} else {
		return data, errors.New("this file format is not supported")
	}

	return data, err
}
