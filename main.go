package semver

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	reSemVer = regexp.MustCompile(`(?:v|V|)((?:\d+\.){0,2}\d+)-{0,1}(.*)`)
)

type Version struct {
	major  int
	minor  int
	patch  int
	suffix string
}

func (v *Version) SetMajor(version int) *Version {
	v.major = version
	return v
}

func (v *Version) SetMinor(version int) *Version {
	v.minor = version
	return v
}

func (v *Version) SetPatch(version int) *Version {
	v.patch = version
	return v
}

func (v *Version) SetSuffix(elements ...string) *Version {
	v.suffix = ""
	if len(elements) > 0 {
		v.suffix = strings.Join(elements, ".")
	}
	return v
}

func (v *Version) Set(major, minor, patch int, suffixes ...string) *Version {
	return v.SetMajor(major).SetMinor(minor).SetPatch(patch).SetSuffix(suffixes...)
}

func (v *Version) SetFromString(str string) *Version {
	version, err := NewVersionFromString(str)
	if err != nil {
		return v
	}
	return version
}

func (v *Version) String() string {
	var s string
	if v.minor == 0 && v.patch == 0 { // only major set
		s = fmt.Sprintf("v%d", v.major)
	} else if v.patch == 0 { // major and minor set
		s = fmt.Sprintf("v%d.%d", v.major, v.minor)
	} else { // all components set
		s = fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
	}
	if v.suffix != "" {
		s += "-" + v.suffix
	}
	return s
}

func NewVersion() *Version {
	v := &Version{
		major:  0,
		minor:  0,
		patch:  0,
		suffix: "",
	}
	return v
}

func NewVersionFromString(str string) (*Version, error) {
	version := &Version{
		major:  0,
		minor:  0,
		patch:  0,
		suffix: "",
	}

	matches := reSemVer.FindStringSubmatch(str)

	if len(matches) != 3 {
		return version, fmt.Errorf("invalid version format: %s", str)
	}

	// Parse version numbers
	versionNumbers := strings.Split(matches[1], ".")
	if len(versionNumbers) < 1 || len(versionNumbers) > 3 {
		return version, fmt.Errorf("invalid version format: %s", str)
	}
	major, err := strconv.Atoi(versionNumbers[0])
	if err != nil {
		return version, fmt.Errorf("invalid major version: %s", versionNumbers[0])
	}
	version.SetMajor(major)
	if len(versionNumbers) >= 2 {
		minor, err := strconv.Atoi(versionNumbers[1])
		if err != nil {
			return version, fmt.Errorf("invalid minor version: %s", versionNumbers[1])
		}
		version.SetMinor(minor)
	}
	if len(versionNumbers) == 3 {
		patch, err := strconv.Atoi(versionNumbers[2])
		if err != nil {
			return version, fmt.Errorf("invalid patch version: %s", versionNumbers[2])
		}
		version.SetPatch(patch)
	}

	version.SetSuffix(matches[2])

	return version, nil
}

// SortVersions sorts a slice of parsed semantic versions.
func SortVersions(versions []*Version) {
	sort.Slice(versions, func(i, j int) bool {
		if versions[i].major != versions[j].major {
			return versions[i].major < versions[j].major
		}
		if versions[i].minor != versions[j].minor {
			return versions[i].minor < versions[j].minor
		}
		if versions[i].patch != versions[j].patch {
			return versions[i].patch < versions[j].patch
		}
		return versions[i].suffix < versions[j].suffix
	})
}

func IsValidVersionString(version string) bool {
	_, err := NewVersionFromString(version)
	return err == nil
}
