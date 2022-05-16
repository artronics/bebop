package pkg

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const versionReg = "v(?P<major>\\d+)\\.(?P<minor>\\d+)\\.(?P<patch>\\d+)(?:-(?P<prerelease>alpha|beta))?"

type ReleaseData struct {
	Version   string
	ReleaseId string
	CommitId  string
}

type VersionData struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string // vM.m.p-{prerelease}
}

func (v *VersionData) String() string {
	s := strings.Builder{}

	s.WriteString(fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch))
	if v.Prerelease != "" {
		s.WriteString(fmt.Sprintf("-%s", v.Prerelease))
	}

	return s.String()
}

type BumpVersion int

const (
	_ BumpVersion = iota
	Major
	Minor
	Patch
)

func Version(rd ReleaseData, bump BumpVersion) (ReleaseData, error) {
	updatedRelease := ReleaseData{}

	// No logic around commitId and releaseId i.e. user will send them
	updatedRelease.ReleaseId = rd.ReleaseId
	updatedRelease.CommitId = rd.CommitId

	ver, err := parseVersion(rd.Version)
	if err != nil {
		return updatedRelease, err
	}
	updatedVer := bumpVersion(ver, bump)
	updatedRelease.Version = updatedVer.String()

	return updatedRelease, nil
}

func bumpVersion(v VersionData, bump BumpVersion) VersionData {
	switch bump {
	case Major:
		v.Major++
	case Minor:
		v.Minor++
	case Patch:
		v.Patch++
	}

	return v
}

func parseVersion(v string) (VersionData, error) {
	var version VersionData
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
