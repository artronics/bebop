package cmd

import (
	"fmt"
	"github.com/nhsdigital/bebop-cli/internal"
	"github.com/nhsdigital/bebop-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/valyala/fasttemplate"
	"log"
	"os"
	"strings"
	"time"
)

const releaseTemplate = ` #  --- DO NOT EDIT --- Auto-generated at: {{ time }}
version: {{ version }}
releaseId: {{ releaseId }}
commitId: {{ commitId }}
`

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Manage project version during release.",
	Long: `This command is used to store project version information in a release.yml file. The information contains
not only the version but also commit sha which can be used as template for some of the healthcheck responses. This
command can also be used to bump versions which is useful during release process.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flags().Lookup("release-file").Value.String()
		data, err := internal.ParseDataFile(path)
		if err != nil {
			log.Fatalln(err.Error())
		}

		relData := pkg.ReleaseData{}
		if version, ok := data["version"]; ok {
			relData.Version = version.(string)
		}
		relData.CommitId = cmd.Flags().Lookup("commitId").Value.String()
		relData.ReleaseId = cmd.Flags().Lookup("releaseId").Value.String()

		bump := validateOperation(cmd.Flags().Lookup("bump").Value.String())
		if err != nil {
			log.Fatalln(err.Error())
		}

		updated, err := pkg.Version(relData, bump)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if updated.ReleaseId == "" && updated.CommitId == "" {
			fmt.Println(updated.Version)
			return
		}

		content := renderTemplate(updated)
		f, err := os.Create(path)
		defer f.Close()

		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = fmt.Fprintf(f, "%s", content)
		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Println(updated.Version)
	},
}

func init() {
	versionCmd.Flags().String("release-file", "release.yml", "the release file containing version information")
	_ = versionCmd.MarkFlagFilename("release-file")

	versionCmd.Flags().String("bump", "minor", "bump version. Valid values are major, minor, patch")
	versionCmd.Flags().String("releaseId", "", "The pipeline release number")
	versionCmd.Flags().String("commitId", "", "The git sha code of this version")

	projectCmd.AddCommand(versionCmd)
}

func validateOperation(bump string) pkg.BumpVersion {
	switch strings.ToLower(bump) {
	case "major":
		return pkg.Major
	case "minor":
		return pkg.Minor
	case "patch":
		return pkg.Patch
	default:
		log.Fatalln("wrong bump version value. Only major minor and patch are valid.")
	}

	return pkg.Patch
}

func renderTemplate(rd pkg.ReleaseData) string {
	template := fasttemplate.New(releaseTemplate, "{{ ", " }}")

	data := map[string]interface{}{
		"time":      time.Now().UTC().Format("2006-01-02 15:04:05"),
		"commitId":  rd.CommitId,
		"releaseId": rd.ReleaseId,
		"version":   rd.Version}

	return template.ExecuteString(data)
}
