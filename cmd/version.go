package cmd

import (
	"errors"
	"fmt"
	"github.com/nhsdigital/bebop-cli/internal"
	"github.com/nhsdigital/bebop-cli/pkg"
	"github.com/spf13/cobra"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const versionReg = "v(?P<major>\\d+)\\.(?P<minor>\\d+)\\.(?P<patch>\\d+)(?:-(?P<prerelease>alpha|beta))?"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Manage project version during release.",
	Long: `This command is used to store project version information in a release.yml file. The information contains
not only the version but also commit sha which can be used as template for some of the healthcheck responses. This
command can also be used to bump versions which is useful during release process.`,
	Run: func(cmd *cobra.Command, args []string) {
		d := pkg.VersionData{}

		path := cmd.Flags().Lookup("release-file").Value.String()
		data, err := internal.ParseDataFile(path)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if version, ok := data["version"]; ok {
			d.Version, err = parseVersion(version.(string))
			if err != nil {
				log.Fatalln(err.Error())
			}
		}
		if commitId, ok := data["commitId"]; ok {
			d.CommitId = commitId.(string)
		}
		if releaseId, ok := data["releaseId"]; ok {
			d.ReleaseId = releaseId.(string)
		}

		bump := validateOperation(cmd.Flags().Lookup("bump").Value.String())
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = pkg.Version(&d, bump)
		fmt.Println(d)
		if err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	versionCmd.Flags().String("release-file", "release.yml", "the release file containing version information")
	_ = versionCmd.MarkFlagFilename("release-file")

	versionCmd.Flags().String("bump", "minor", "bump version. Valid values are major, minor, patch")

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

func parseVersion(v string) (pkg.VersionParsed, error) {
	var version pkg.VersionParsed
	r := regexp.MustCompile(versionReg)
	matched := r.MatchString(v)
	if !matched {
		// FIXME: it still matches with v1.1.1-foo  prerelease part can be anything. The rest (numbers) work as expected
		return version, errors.New("provided version doesn't match pattern: v{minor}.{major}.{patch}-{[alpha|beta]}")
	}

	matches := r.FindStringSubmatch(v)
	var toInt = func(idx string) int {
		v := r.SubexpIndex(idx)
		n, _ := strconv.ParseInt(matches[v], 10, 32)

		return int(n)
	}

	version.Major = toInt("major")
	version.Minor = toInt("minor")
	version.Patch = toInt("patch")

	if pre := r.SubexpIndex("prerelease"); pre != -1 {
		version.Prerelease = matches[pre]
	}

	return version, nil
}
