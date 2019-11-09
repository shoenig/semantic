// Package semantic provides utilities for parsing and creating semver2.0 tags.
//
// For more information, see https://semver.org/
package semantic // import "gophers.dev/pkgs/semantic"

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gophers.dev/pkgs/regexplus"
)

var (
	// Unlike the official example regexp given, this one enforces the
	// 'v' prefix, which is always required in Go module version strings.
	//
	// The old regexp used was lenient around having sensible numbers. It would
	// allow things like "v00.00.00", whereas the new expression does not.
	//
	// The example regexp is available at
	// https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
	semverRe = regexp.MustCompile(`^v(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<pr>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<bm>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
)

// New creates a new Tag with the most basic amount of information, which includes
// the major, minor, and patch version levels.
//
// Examples of a basic tag:
//   v1.0.0
//   v3.20.100
func New(major, minor, patch int) Tag {
	return New3(major, minor, patch, "", "")
}

// New2 creates a new Tag with the basic version information, plus an additional
// associated pre-release suffix.
//
// The pre-release information is appended to the end of a version string,
// denoted with a '-' prefix.
//
// Examples of a tag with "pre-release" information:
//   v1.0.0-alpha
//   v1.0.0-rc1
func New2(major, minor, patch int, preRelease string) Tag {
	return New3(major, minor, patch, preRelease, "")
}

// New3 creates a new Tag with the basic version information, plus an additional
// associated pre-release suffix plus an additional build-metadata suffix.
//
// The pre-release information is appended to the end of a version string, but
// before the build-metadata, denoted with a '-' prefix.
//
// The build-metadata information is appended to the end of a version string,
// denoted with a '+' prefix.
//
// Examples of a tag with "pre-release" and "build-metadata" information:
//   v1.0.0-beta+exp.sha.5114f85
//   v1.0.0-rc1+20130313144700
//   1.0.0-alpha+001
func New3(major, minor, patch int, preRelease, buildMetadata string) Tag {
	return Tag{
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		PreRelease:    normalize(preRelease),
		BuildMetadata: normalize(buildMetadata),
	}
}

// New4 creates a new Tag with basic version information, plus an additional
// associated build-metadata suffix.
//
// The build-metadata information is appended to the end of a version string,
// denoted with a '+' prefix.
//
// Examples of a tag with "build-metadata" information:
//   v1.0.0+exp.sha.5114f85
//   v1.0.0+20130313144700
func New4(major, minor, patch int, buildMetadata string) Tag {
	return New3(major, minor, patch, "", buildMetadata)
}

func normalize(s string) string {
	noDash := strings.TrimPrefix(s, "-")
	noPlus := strings.TrimPrefix(noDash, "+")
	return noPlus
}

func Parse(s string) (Tag, bool) {
	matches := regexplus.FindNamedSubmatches(semverRe, s)

	major, exists := matches["major"]
	if !exists {
		return Tag{}, false
	}

	minor, exists := matches["minor"]
	if !exists {
		return Tag{}, false
	}

	patch, exists := matches["patch"]
	if !exists {
		return Tag{}, false
	}

	preRelease := matches["pr"]
	buildMetadata := matches["bm"]

	return Tag{
		Major:         number(major),
		Minor:         number(minor),
		Patch:         number(patch),
		PreRelease:    preRelease,
		BuildMetadata: buildMetadata,
	}, true
}

func number(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic("bug in our tag regexp")
	}
	return n
}

type Tag struct {
	Major         int
	Minor         int
	Patch         int
	PreRelease    string
	BuildMetadata string
}

func (t Tag) String() string {
	base := fmt.Sprintf(
		"v%d.%d.%d",
		t.Major,
		t.Minor,
		t.Patch,
	)

	if t.PreRelease == "" && t.BuildMetadata == "" {
		return base
	}

	if t.BuildMetadata == "" {
		return base + "-" + t.PreRelease
	}

	if t.PreRelease == "" {
		return base + "+" + t.BuildMetadata
	}

	return base + "-" + t.PreRelease + "+" + t.BuildMetadata
}

func (t Tag) Base() Tag {
	return New(t.Major, t.Minor, t.Patch)
}

func (t Tag) IsBase() bool {
	return t.PreRelease == ""
}

func (t Tag) Less(o Tag) bool {
	// build-metadata should be explicitly ignored for comparisons ; see https://semver.org/#spec-item-10
	// pre-release is NOT ignored ;

	if t.Major < o.Major {
		return true
	} else if t.Major > o.Major {
		return false
	} else if prAlessB(t.PreRelease, o.PreRelease) {
		return true
	}

	if t.Minor < o.Minor {
		return true
	} else if t.Minor > o.Minor {
		return false
	} else if prAlessB(t.PreRelease, o.PreRelease) {
		return true
	}

	if t.Patch < o.Patch {
		return true
	} else if t.Patch > o.Patch {
		return false
	} else if prAlessB(t.PreRelease, o.PreRelease) {
		return true
	}

	return false
}

// return true if a's extension precedes b's extension
// normally this is ascibetical, however the empty string
// is a special case that is higher priority than all else
func prAlessB(a, b string) bool {
	if a == "" {
		return false
	}

	if b == "" {
		return true
	}

	return a < b
}

type BySemver []Tag

func (tags BySemver) Len() int      { return len(tags) }
func (tags BySemver) Swap(x, y int) { tags[x], tags[y] = tags[y], tags[x] }
func (tags BySemver) Less(x, y int) bool {
	X := tags[x]
	Y := tags[y]
	return X.Less(Y)
}
