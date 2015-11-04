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

	return static.New(config, &static.DirFile{
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

var (
	flagStaticDir = flag.String("i", "static", "Static File Directory to compile")
	flagOuputFile = flag.String("o", "", "Output File Path to write to")
	flagPkg       = flag.String("pkg", "main", "Package name of the generated static file")
	flagGroup     = flag.String("group", "Assets", "The group name of the static files i.e. CSS, JS, Assets, HTML")
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
		processFiles(filepath.Clean(*flagStaticDir))
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

func processFiles(dir string) {

	fi, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err)
	}

	writer.WriteString(fmt.Sprintf(dirFileStart, dir, fi.Name(), fi.Size(), fi.Mode(), fi.ModTime().Unix(), true, ""))
	processFilesRecursive(dir, "", false, "")
	writer.WriteString(dirFileEnd)
}

// need isSymlinkDir variable as it is valid for symlinkDir to be blank
func processFilesRecursive(path string, dir string, isSymlinkDir bool, symlinkDir string) {

	var b64File string
	var p string

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(0)

	for _, file := range files {

		info := file
		b64File = ""
		p = path + string(os.PathSeparator) + file.Name()
		fPath := p

		if isSymlinkDir {
			fPath = strings.Replace(p, dir, symlinkDir, 1)
		}

		fmt.Println("Processing:", fPath)

		if file.IsDir() {

			// write out here
			writer.WriteString(fmt.Sprintf(dirFileStart, fPath, info.Name(), info.Size(), info.Mode(), info.ModTime().Unix(), true, ""))
			processFilesRecursive(p, p, isSymlinkDir, symlinkDir+string(os.PathSeparator)+info.Name())
			writer.WriteString(dirFileEndArray)
			continue
		}

		if file.Mode()&os.ModeSymlink == os.ModeSymlink {

			link, err := filepath.EvalSymlinks(p)
			if err != nil {
				log.Fatal("Error Resolving Symlink", err)
			}

			fi, err := os.Stat(link)
			if err != nil {
				log.Fatal(err)
			}

			info = fi

			if fi.IsDir() {
				// write out here
				writer.WriteString(fmt.Sprintf(dirFileStart, fPath, file.Name(), info.Size(), info.Mode(), info.ModTime().Unix(), true, ""))
				processFilesRecursive(link, link, true, fPath)
				writer.WriteString(dirFileEndArray)
				continue
			}
		}

		// if we get here it's a file

		// read file
		b, err := ioutil.ReadFile(p)
		if err != nil {
			log.Fatal(err)
		}

		// gzip
		var gzBuff bytes.Buffer
		gz := gzip.NewWriter(&gzBuff)
		defer gz.Close()

		_, err = gz.Write(b)
		if err != nil {
			log.Fatal(err)
		}

		// Flush not quaranteed to flush, must close
		// gz.Flush()
		gz.Close()

		// turn into chunked base64 string
		var bb bytes.Buffer
		b64 := base64.NewEncoder(base64.StdEncoding, &bb)
		b64.Write(gzBuff.Bytes())
		b64.Close()
		// b64File += "\n"
		chunk := make([]byte, 80)

		for n, _ := bb.Read(chunk); n > 0; n, _ = bb.Read(chunk) {
			b64File += string(chunk[0:n]) + "\n"
		}

		// write out here
		writer.WriteString(fmt.Sprintf(dirFileStart, fPath, file.Name(), info.Size(), info.Mode(), info.ModTime().Unix(), false, b64File))
		writer.WriteString(dirFileEndArray)
	}
}
