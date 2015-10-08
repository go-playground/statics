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
	helpString = `
nothing yet
`
	startFile = `package %s

import (
	"github.com/joeybloggs/statics/statics"
)

func NewStatic%s(config *statics.Config) *statics.Files {

	return statics.New(config, map[string]*statics.FileInfo{
`
	endfile = `},
	)
}`
	mapEntry = `"%s" : {
		Contents: "%s",
	},
`
)

var (
	flagStaticDir = flag.String("i", "static", "Static File Directory to compile")
	flagOuputFile = flag.String("o", "", "Output File Path to write to")
	flagPkg       = flag.String("pkg", "main", "Package name of the generated static file")
	flagGroup     = flag.String("group", "assets", "The group name of the static files i.e. CSS, JS, Assets, HTML")
	flagHelp      = flag.Bool("help", false, "-help")
	writer        *bufio.Writer
)

func help() {
	fmt.Printf(helpString)
}

func main() {
	parseFlags()

	os.Remove(*flagOuputFile)
	f, err := os.Create(*flagOuputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer = bufio.NewWriter(f)

	writer.WriteString(fmt.Sprintf(startFile, *flagPkg, strings.ToUpper((*flagGroup)[0:1])+(*flagGroup)[1:]))

	processFiles(*flagStaticDir, false, "")

	writer.WriteString(endfile)
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

	if *flagHelp {
		help()
		os.Exit(0)
	}

	s := filepath.Clean(*flagStaticDir)
	flagStaticDir = &s

	if len(*flagStaticDir) == 0 || *flagStaticDir == "." {
		fmt.Printf("\n**invalid Static File Directoy '%s'\n", *flagStaticDir)
		help()
		os.Exit(1)
	}

	if len(*flagOuputFile) == 0 {
		fmt.Printf("\n**invalid Output Directory '%s'\n", *flagOuputFile)
		help()
		os.Exit(1)
	}

	if len(*flagPkg) == 0 {
		fmt.Printf("\n**invalid Package Name '%s'\n", *flagPkg)
		help()
		os.Exit(1)
	}
}

// need isSymlinkDir variable as it is valid for symlinkDir to be blank
func processFiles(dir string, isSymlinkDir bool, symlinkDir string) {

	walker := func(path string, info os.FileInfo, err error) error {

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

		// read file
		b, err := ioutil.ReadFile(path)
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

		gz.Flush()

		// turn into base64
		b64File := base64.StdEncoding.EncodeToString(gzBuff.Bytes())

		// read from realpath but key should be symlink one
		if isSymlinkDir {
			path = strings.Replace(path, dir, symlinkDir, 1)
		}

		path = strings.Replace(path, *flagStaticDir, filepath.Base(*flagStaticDir), 1)

		fmt.Println("Processing:", path)

		writer.WriteString(fmt.Sprintf(mapEntry, path, b64File))

		return nil
	}

	if err := filepath.Walk(dir, walker); err != nil {
		fmt.Printf("\n**could not walk project path '%s'\n%s\n", *flagStaticDir, err)
		os.Exit(1)
	}
}
