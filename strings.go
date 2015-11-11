package main

const (
	functionComments = "// NewStatic%s initializes a new static.Files instance for use"
	initStartFile    = `//go:generate statics -i=%s -o=%s -pkg=%s -group=%s %s %s

	package %s

import "github.com/joeybloggs/statics/static"

// newStatic%s initializes a new *static.Files instance for use
func newStatic%s(config *static.Config) (*static.Files, error) {

	return static.New(config, &static.DirFile{
`
	initEndfile = `})
}`
	startFile = `//go:generate statics -i=%s -o=%s -pkg=%s -group=%s %s %s

	package %s

import (
	"os"

	"github.com/joeybloggs/statics/static"
)

// newStatic%s initializes a new *static.Files instance for use
func newStatic%s(config *static.Config) (*static.Files, error) {

	return static.New(config, `
	endfile = `)
}`

	dirFileEnd = `},
}`

	dirFileEndArray = `},
},
`
)

var (
	dirFileStart = `&static.DirFile{
		Path: %q,
		Name: "%s",
		Size: %d,
		Mode: os.FileMode(%d),
		ModTime: %v,
		IsDir: %t,
		Compressed: ` + "`\n%s`" + `,
		Files: []*static.DirFile{`
)
