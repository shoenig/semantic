package semantic

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/shoenig/ignore"
	"github.com/shoenig/test/must"
)

func Test_New(t *testing.T) {
	tag := New(1, 2, 3)
	must.Eq(t, Tag{
		Major: 1,
		Minor: 2,
		Patch: 3,
	}, tag)
}

func Test_New2(t *testing.T) {
	tag := New2(1, 2, 3, "rc1")
	must.Eq(t, Tag{
		Major:      1,
		Minor:      2,
		Patch:      3,
		PreRelease: "rc1",
	}, tag)
}

func Test_New3(t *testing.T) {
	tag := New3(1, 2, 3, "alpha", "linux")
	must.Eq(t, Tag{
		Major:         1,
		Minor:         2,
		Patch:         3,
		PreRelease:    "alpha",
		BuildMetadata: "linux",
	}, tag)
}

func Test_New4(t *testing.T) {
	tag := New4(1, 2, 3, "abc123")
	must.Eq(t, Tag{
		Major:         1,
		Minor:         2,
		Patch:         3,
		BuildMetadata: "abc123",
	}, tag)
}

func Test_Parse(t *testing.T) {
	try := func(s string, exp Tag, expOK bool) {
		result, ok := Parse(s)
		must.Eq(t, expOK, ok)
		must.Eq(t, exp, result)
	}

	try("v1.3.5", New(1, 3, 5), true)
	try("v111.222.333", New(111, 222, 333), true)
	try("v1.2.3-alpha", New2(1, 2, 3, "alpha"), true)
	try("v1.2.3-alpha2", New2(1, 2, 3, "alpha2"), true)
	try("1.2.3", empty, false)       // missing v
	try("v1.2.3_beta", empty, false) // dash required for extension
	try("v0.8.2-0.20190227000051-27936f6d90f9", New2(0, 8, 2, "0.20190227000051-27936f6d90f9"), true)
	try("v2.0.0+incompatible", New4(2, 0, 0, "incompatible"), true)
	try("v2.0.0-pre+incompatible", New3(2, 0, 0, "pre", "incompatible"), true)
}

func Test_String(t *testing.T) {
	must.Eq(t, "v1.2.3", New(1, 2, 3).String())
	must.Eq(t, "v0.8.2-0.20190227000051-27936f6d90f9", New2(0, 8, 2, "0.20190227000051-27936f6d90f9").String())
	must.Eq(t, "v2.0.0+incompatible", New4(2, 0, 0, "incompatible").String())
	must.Eq(t, "v2.0.0-pre+incompatible", New3(2, 0, 0, "pre", "incompatible").String())
}

func Test_Sort_BySemver(t *testing.T) {
	list := []Tag{
		New(3, 1, 2),
		New(3, 3, 1),
		New(1, 3, 2),
		New(2, 1, 1),
		New(1, 6, 2),
		New(3, 3, 3),
		New(2, 4, 1),
		New(1, 8, 2),
		New(1, 7, 0),
	}
	sort.Sort(sort.Reverse(BySemver(list)))
	must.Eq(t, []Tag{
		New(3, 3, 3),
		New(3, 3, 1),
		New(3, 1, 2),
		New(2, 4, 1),
		New(2, 1, 1),
		New(1, 8, 2),
		New(1, 7, 0),
		New(1, 6, 2),
		New(1, 3, 2),
	}, list)
}

func Test_cmpChunk(t *testing.T) {
	try := func(a, b string, exp int) {
		result := cmpChunk(a, b)
		must.Eq(t, exp, result)
	}

	try("alpha", "alpha", 0)
	try("alpha", "beta", -1)
	try("beta", "alpha", 1)

	try("alpha", "alpha-1", -1)
	try("alpha-1", "alpha", 1)

	try("1", "1", 0)
	try("10", "1", 1)
	try("1", "10", -1)

	try("10", "2", 1)
	try("2", "10", -1)

	try("alpha", "999999999", 1)
	try("999999999", "alpha", -1)
}

func Test_cmpPre(t *testing.T) {
	try := func(a, b string, exp bool) {
		result := cmpPreRelease(a, b)
		must.Eq(t, exp, result)
	}

	try("", "alpha", false)
	try("alpha", "", true)

	try("alpha", "alpha.1", true)
	try("alpha.1", "alpha", false)

	try("alpha.1", "alpha.beta", true)
	try("alpha.beta", "alpha.1", false)

	try("alpha.beta", "beta", true)
	try("beta", "alpha.beta", false)

	try("beta", "beta.2", true)
	try("beta.2", "beta", false)

	try("beta.2", "beta.11", true)
	try("beta.11", "beta.2", false)

	try("beta.11", "rc.1", true)
	try("rc.1", "beta.11", false)
}

// Example: 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0
func Test_Sort_BySemver_preReleases(t *testing.T) {
	expected := []Tag{
		New2(1, 1, 1, "alpha"),
		New2(1, 1, 1, "alpha.1"),
		New2(1, 1, 1, "alpha.beta"),
		New2(1, 1, 1, "beta"),
		New2(1, 1, 1, "beta.2"),
		New2(1, 1, 1, "beta.11"),
		New2(1, 1, 1, "rc.1"),
		New(1, 1, 1), // non-prerelease tags are always highest
	}

	tags := []Tag{
		New2(1, 1, 1, "alpha.1"),
		New(1, 1, 1),
		New2(1, 1, 1, "alpha"),
		New2(1, 1, 1, "beta.11"),
		New2(1, 1, 1, "beta.2"),
		New2(1, 1, 1, "rc.1"),
		New2(1, 1, 1, "beta"),
		New2(1, 1, 1, "alpha.beta"),
	}

	sort.Sort(BySemver(tags))
	must.Eq(t, expected, tags)
}

func load(t *testing.T, filename string) []Tag {
	f, err := os.Open(filename)
	must.NoError(t, err)
	defer ignore.Close(f)

	return parse(t, f)
}

func parse(t *testing.T, r io.Reader) []Tag {
	var tags []Tag
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if tag, ok := Parse(line); ok {
			tags = append(tags, tag)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Error("scan lines:", err)
	}
	return tags
}

func testFile(t *testing.T, input, expected string) {
	orig := load(t, filepath.Join("hack", input))
	exp := load(t, filepath.Join("hack", expected))

	result := make([]Tag, 0, len(orig))
	for _, tag := range orig {
		result = append(result, tag)
	}
	sort.Sort(BySemver(result))

	for i := 0; i < len(orig); i++ {
		a, b, c := orig[i], exp[i], result[i]
		triple := fmt.Sprintf("(%s, %s | %s)", b, c, a)
		t.Logf("triple[%d]: %s", i, triple)
		if !b.Equal(c) {
			t.Errorf("  out of order tag %s", triple)
		}
	}
}

func Test_Sort_nomad(t *testing.T) {
	testFile(t, "nomad.orig.txt", "nomad.exp.txt")
}

func Test_Sort_consul(t *testing.T) {
	testFile(t, "consul.orig.txt", "consul.exp.txt")
}
