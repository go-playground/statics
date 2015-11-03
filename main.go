package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	functionComments = "// NewStatic%s initializes a new static.Files instance for use"
	initStartFile    = `package %s

import "github.com/joeybloggs/statics/static"

// NewStatic%s initializes a new static.Files instance for use
func NewStatic%s(config *static.Config) (*static.Files, error) {

	return static.New(config, map[string]*static.File{
`
	initEndfile = `})
}`
	startFile = `package %s

import (
	"os"

	"github.com/joeybloggs/statics/static"
)

// NewStatic%s initializes a new static.Files instance for use
func NewStatic%s(config *static.Config) (*static.Files, error) {

	return static.New(config, map[string]*static.File{
`
	endfile = `},
	)
}`
	mapEntry = `%q : {
		Path: %q,
		Filename: "%s",
		Filesize: %d,
		FMode: os.FileMode(%d),
		Modtime: %v,
		Compressed: %s,
	},
`
)

var (
	flagStaticDir = flag.String("i", "static", "Static File Directory to compile")
	flagOuputFile = flag.String("o", "", "Output File Path to write to")
	flagPkg       = flag.String("pkg", "main", "Package name of the generated static file")
	flagGroup     = flag.String("group", "assets", "The group name of the static files i.e. CSS, JS, Assets, HTML")
	flagInit      = flag.Bool("init", false, " determines if only initializing the static file without contents")
	writer        *bufio.Writer
)

func main() {
	parseFlags()

	os.Remove(*flagOuputFile)
	f, err := os.Create(*flagOuputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	funcName := strings.ToUpper((*flagGroup)[0:1]) + (*flagGroup)[1:]

	writer = bufio.NewWriter(f)

	if *flagInit {
		writer.WriteString(fmt.Sprintf(initStartFile, *flagPkg, funcName, funcName))
		writer.WriteString(initEndfile)
	} else {
		writer.WriteString(fmt.Sprintf(startFile, *flagPkg, funcName, funcName))
		processFiles(*flagStaticDir, false, "")
		writer.WriteString(endfile)
	}

	writer.Flush()

	f.Close()

	// after file written run gofmt on file
	cmd := exec.Command("gofmt", "-s", "-w", *flagOuputFile)
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() {

	flag.Parse()

	s := filepath.Clean(*flagStaticDir)
	flagStaticDir = &s

	if len(*flagStaticDir) == 0 || *flagStaticDir == "." {
		panic("**invalid Static File Directoy '" + *flagStaticDir + "'")
	}

	if len(*flagOuputFile) == 0 {
		panic("**invalid Output Directory")
	}

	if len(*flagPkg) == 0 {
		panic("**invalid Package Name")
	}
}

// need isSymlinkDir variable as it is valid for symlinkDir to be blank
func processFiles(dir string, isSymlinkDir bool, symlinkDir string) {

	walker := func(path string, info os.FileInfo, err error) error {

		if info == nil {
			fmt.Println(path)
			fmt.Println(err)
		}

		if info.IsDir() {
			return nil
		}

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {

			link, err := filepath.EvalSymlinks(path)
			if err != nil {
				fmt.Println("Error Resolving Symlink", err)
				return err
			}

			fi, err := os.Stat(link)
			if err != nil {
				fmt.Println(err)
				return err
			}

			if fi.IsDir() {
				// call process files, otherwise just fall through and read file.
				processFiles(link, true, path)
				return nil
			}
		}

		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// read file
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// gzip
		var gzBuff bytes.Buffer
		gz := gzip.NewWriter(&gzBuff)
		defer gz.Close()

		_, err = gz.Write(b)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Flush not quaranteed to flush, must close
		// gz.Flush()
		gz.Close()

		// turn into chunked base64 string
		var bb bytes.Buffer
		b64 := base64.NewEncoder(base64.StdEncoding, &bb)
		b64.Write(gzBuff.Bytes())
		b64.Close()
		b64File := "`\n"
		chunk := make([]byte, 80)

		for n, _ := bb.Read(chunk); n > 0; n, _ = bb.Read(chunk) {
			b64File += string(chunk[0:n]) + "\n"
		}

		b64File += "`"

		// read from realpath but key should be symlink one
		if isSymlinkDir {
			path = strings.Replace(path, dir, symlinkDir, 1)
		}

		path = strings.Replace(path, *flagStaticDir, filepath.Base(*flagStaticDir), 1)

		fmt.Println("Processing:", path)

		fi, err := f.Stat()
		if err != nil {
			fmt.Println(err)
			return err
		}

		writer.WriteString(fmt.Sprintf(mapEntry, path, path, fi.Name(), fi.Size(), fi.Mode(), fi.ModTime().Unix(), b64File))

		return nil
	}

	if err := filepath.Walk(dir, walker); err != nil {
		fmt.Printf("\n**could not walk project path '%s'\n%s\n", *flagStaticDir, err)
		os.Exit(1)
	}
}
