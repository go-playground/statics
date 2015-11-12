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
	"regexp"
	"strings"
)

var (
	flagStaticDir = flag.String("i", "static", "Static directory to embed")
	flagOuputFile = flag.String("o", "", "Output file to write to")
	flagPkg       = flag.String("pkg", "main", "Package name of the generated static file")
	flagGroup     = flag.String("group", "Assets", "The group name of the static files i.e. CSS, JS, Assets, HTML. It will be added to the generated static function name.")
	flagIgnore    = flag.String("ignore", "", "Regexp for files/dirs we should ignore i.e. \\.gitignore")
	flagPrefix    = flag.String("prefix", "", "Prefix to strip from file paths")
	flagInit      = flag.Bool("init", false, " determines if only initializing the static file without contents")

	ignoreRegexp *regexp.Regexp
	writer       *bufio.Writer
)

func main() {
	parseFlags()

	os.Remove(*flagOuputFile)

	f, err := os.Create(*flagOuputFile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	funcName := strings.ToUpper((*flagGroup)[0:1]) + (*flagGroup)[1:]

	writer = bufio.NewWriter(f)

	ignoreFlag := ""
	prefixFlag := ""

	if len(*flagIgnore) > 0 {
		ignoreFlag = "-ignore=" + *flagIgnore
	}

	if len(*flagPrefix) > 0 {
		prefixFlag = "-prefix=" + *flagPrefix
	}

	if *flagInit {
		writer.WriteString(fmt.Sprintf(initStartFile, *flagStaticDir, *flagOuputFile, *flagPkg, *flagGroup, ignoreFlag, prefixFlag, *flagPkg, funcName, funcName))
		writer.WriteString(initEndfile)
	} else {
		writer.WriteString(fmt.Sprintf(startFile, *flagStaticDir, *flagOuputFile, *flagPkg, *flagGroup, ignoreFlag, prefixFlag, *flagPkg, funcName, funcName))
		processFiles(filepath.Clean(*flagStaticDir))
		writer.WriteString(endfile)
	}

	writer.Flush()

	f.Close()

	// after file written run gofmt on file
	cmd := exec.Command("gofmt", "-s", "-w", *flagOuputFile)
	if err = cmd.Run(); err != nil {
		log.Panic(err)
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

	if len(*flagIgnore) > 0 {

		var err error

		ignoreRegexp, err = regexp.Compile(*flagIgnore)
		if err != nil {
			panic("**Error Compiling Regex:" + err.Error())
		}
	}
}

func processFiles(dir string) {

	fi, err := os.Stat(dir)
	if err != nil {
		log.Panic(err)
	}

	dPath := applyPathOptions(dir)

	writer.WriteString(fmt.Sprintf(dirFileStart, dPath, fi.Name(), fi.Size(), fi.Mode(), fi.ModTime().Unix(), true, ""))
	processFilesRecursive(dir, "", false, "")
	writer.WriteString(dirFileEnd)
}

// need isSymlinkDir variable as it is valid for symlinkDir to be blank
func processFilesRecursive(path string, dir string, isSymlinkDir bool, symlinkDir string) {

	var b64File string
	var p string
	var tmpPath string

	f, err := os.Open(path)
	if err != nil {
		log.Panic(err)
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

		if ignoreRegexp != nil && ignoreRegexp.MatchString(fPath) {
			continue
		}

		if file.IsDir() {

			tmpPath = applyPathOptions(fPath)

			// write out here
			writer.WriteString(fmt.Sprintf(dirFileStart, tmpPath, info.Name(), info.Size(), info.Mode(), info.ModTime().Unix(), true, ""))
			processFilesRecursive(p, p, isSymlinkDir, symlinkDir+string(os.PathSeparator)+info.Name())
			writer.WriteString(dirFileEndArray)
			continue
		}

		if file.Mode()&os.ModeSymlink == os.ModeSymlink {

			link, err := filepath.EvalSymlinks(p)
			if err != nil {
				log.Panic("Error Resolving Symlink", err)
			}

			fi, err := os.Stat(link)
			if err != nil {
				log.Panic(err)
			}

			info = fi

			if fi.IsDir() {

				tmpPath = applyPathOptions(fPath)

				// write out here
				writer.WriteString(fmt.Sprintf(dirFileStart, tmpPath, file.Name(), info.Size(), info.Mode(), info.ModTime().Unix(), true, ""))
				processFilesRecursive(link, link, true, fPath)
				writer.WriteString(dirFileEndArray)
				continue
			}
		}

		// if we get here it's a file

		// read file
		b, err := ioutil.ReadFile(p)
		if err != nil {
			log.Panic(err)
		}

		// gzip
		var gzBuff bytes.Buffer
		gz := gzip.NewWriter(&gzBuff)
		defer gz.Close()

		_, err = gz.Write(b)
		if err != nil {
			log.Panic(err)
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

		fPath = applyPathOptions(fPath)

		fmt.Println("Processing:", fPath)

		// write out here
		writer.WriteString(fmt.Sprintf(dirFileStart, fPath, file.Name(), info.Size(), info.Mode(), info.ModTime().Unix(), false, b64File))
		writer.WriteString(dirFileEndArray)
	}
}

func applyPathOptions(path string) string {
	path = strings.TrimPrefix(path, *flagPrefix)
	path = strings.TrimLeft(path, string(os.PathSeparator))
	return string(os.PathSeparator) + path
}
