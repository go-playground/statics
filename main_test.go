package main

import (
	"io/ioutil"
	"os"
	"testing"

	. "gopkg.in/go-playground/assert.v1"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//

func TestMain(m *testing.M) {

	// setup

	os.Exit(m.Run())

	// teardown
}

func TestNonExistantStaticDir(t *testing.T) {

	i := "static/test-files/garbagedir"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	ignore := ""
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	PanicMatches(t, func() { main() }, "stat static/test-files/garbagedir: no such file or directory")
}

func TestBadPackage(t *testing.T) {

	i := "static/test-files/test-inner"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := ""
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	init := false
	flagInit = &init

	PanicMatches(t, func() { main() }, "**invalid Package Name")
}

func TestBadOutputDir(t *testing.T) {

	i := "static/test-files/test-inner"
	flagStaticDir = &i

	o := ""
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	init := false
	flagInit = &init

	PanicMatches(t, func() { main() }, "**invalid Output Directory")
}

func TestBadStaticDir(t *testing.T) {

	i := ""
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	init := false
	flagInit = &init

	PanicMatches(t, func() { main() }, "**invalid Static File Directoy '.'")
}

func TestGenerateInitFile(t *testing.T) {

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	init := true
	flagInit = &init

	main()

	b, err := ioutil.ReadFile("static/test-files/test.go")
	Equal(t, err, nil)

	expected := `//go:generate statics -i=static/test-files/teststart -o=static/test-files/test.go -pkg=test -group=Assets

package test

import "github.com/go-playground/statics/static"

// newStaticAssets initializes a new *static.Files instance for use
func newStaticAssets(config *static.Config) (*static.Files, error) {

	return static.New(config, &static.DirFile{})
}
`

	Equal(t, string(b), expected)
}

func TestIgnore(t *testing.T) {

	Equal(t, true, true)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	ignore := ".*.txt"
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	main()
}

func TestBadIgnore(t *testing.T) {

	Equal(t, true, true)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	ignore := "([12.gitignore"
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	PanicMatches(t, func() { main() }, "**Error Compiling Regex:error parsing regexp: missing closing ]: `[12.gitignore`")
}

func TestGenerateFilePrefix(t *testing.T) {

	Equal(t, true, true)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	ignore := ""
	flagIgnore = &ignore

	prefix := "static/"
	flagPrefix = &prefix

	init := false
	flagInit = &init

	main()
}

func TestGenerateFile(t *testing.T) {

	Equal(t, true, true)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	ignore := ""
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	main()
}
