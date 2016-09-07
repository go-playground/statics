package static

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	pathSep = "/"
)

// DirFile contains the static directory and file content info
type DirFile struct {
	Path       string
	Name       string
	Size       int64
	Mode       os.FileMode
	ModTime    int64
	IsDir      bool
	Compressed string
	Files      []*DirFile
}

// Files contains a full instance of a static file collection
type Files struct {
	dir dir
}

// Config contains information about how extracting the data should behave
// NOTE: FallbackToDisk falls back to disk when file not found in static assets
// usefull when you have a mixture of static assets and some that need to remain
// on disk i.e. a users avatar image
type Config struct {
	UseStaticFiles bool
	FallbackToDisk bool   // falls back to disk when file not found in static assets
	AbsPkgPath     string // the Absolute package path used for local file reading when UseStaticFiles is false
}

// New create a new static file instance.
func New(config *Config, dirFile *DirFile) (*Files, error) {

	files := map[string]*file{}

	if config.UseStaticFiles {
		processFiles(files, dirFile)
	} else {

		if config.AbsPkgPath[:7] == "$GOPATH" {

			gopath := os.Getenv("GOPATH")

			if len(gopath) == 0 {
				return nil, errors.New("$GOPATH could not be found; you're setup is not correct")
			}

			config.AbsPkgPath = gopath + config.AbsPkgPath[7:]
		}

		if !filepath.IsAbs(config.AbsPkgPath) {
			return nil, errors.New("AbsPkgPath is required when not using static files otherwise the static package has no idea where to grab local files from when your package is used from within another package.")
		}
	}

	return &Files{
		dir: dir{
			useStaticFiles: config.UseStaticFiles,
			fallbackToDisk: config.FallbackToDisk,
			files:          files,
			absPkgPath:     filepath.Clean(config.AbsPkgPath),
		},
	}, nil
}

func processFiles(files map[string]*file, dirFile *DirFile) *file {

	f := &file{
		path:    dirFile.Path,
		name:    dirFile.Name,
		size:    dirFile.Size,
		mode:    dirFile.Mode,
		modTime: dirFile.ModTime,
		isDir:   dirFile.IsDir,
		files:   []*file{},
	}

	files[filepath.ToSlash(f.path)] = f

	if dirFile.IsDir {
		for _, nestedFile := range dirFile.Files {
			resultFile := processFiles(files, nestedFile)
			f.files = append(f.files, resultFile)
		}

		return f
	}

	// decompress file contents
	b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(dirFile.Compressed))
	reader, err := gzip.NewReader(b64)
	if err != nil {
		log.Fatal(err)
	}

	f.data, err = ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

// FS returns an http.FileSystem object for serving files over http
func (f *Files) FS() http.FileSystem {
	return f.dir
}

func (f *Files) determinePath(name string) string {

	if f.dir.useStaticFiles {
		return name
	}

	return f.dir.absPkgPath + name
}

// GetHTTPFile returns an http.File object
func (f *Files) GetHTTPFile(filename string) (http.File, error) {
	return f.dir.Open(f.determinePath(filename))
}

// ReadFile returns a files contents as []byte from the filesystem, static or local
func (f *Files) ReadFile(filename string) ([]byte, error) {

	file, err := f.dir.Open(f.determinePath(filename))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(file)
}

// ReadDir reads the directory named by dirname and returns
// a list of sorted directory entries.
func (f *Files) ReadDir(dirname string) ([]os.FileInfo, error) {

	file, err := f.dir.Open(f.determinePath(dirname))
	if err != nil {
		return nil, err
	}

	results, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}

	sort.Sort(byName(results))

	return results, nil
}

// ReadFiles returns a directories file contents as a map[string][]byte from the filesystem, static or local
func (f *Files) ReadFiles(dirname string, recursive bool) (map[string][]byte, error) {

	dirname = f.determinePath(dirname)

	results := map[string][]byte{}

	file, err := f.dir.Open(dirname)
	if err != nil {
		return nil, err
	}

	if err = f.readFilesRecursive(dirname+pathSep, file, results, recursive); err != nil {
		return nil, err
	}

	return results, nil
}

func (f *Files) readFilesRecursive(dirname string, file http.File, results map[string][]byte, recursive bool) error {

	files, err := file.Readdir(-1)
	if err != nil {
		return err
	}

	var fpath string

	for _, fi := range files {

		fpath = dirname + fi.Name()

		newFile, err := f.dir.Open(fpath)
		if err != nil {
			return err
		}

		if fi.IsDir() {

			if !recursive {
				continue
			}

			err := f.readFilesRecursive(fpath+pathSep, newFile, results, recursive)
			if err != nil {
				return err
			}

			continue
		}

		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {

			link, err := filepath.EvalSymlinks(fpath)
			if err != nil {
				log.Panic("Error Resolving Symlink", err)
			}

			fi, err := os.Stat(link)
			if err != nil {
				log.Panic(err)
			}

			if fi.IsDir() {

				if !recursive {
					continue
				}

				err := f.readFilesRecursive(fpath+pathSep, newFile, results, recursive)
				if err != nil {
					return err
				}

				continue
			}
		}

		if !f.dir.useStaticFiles {
			fpath = strings.Replace(fpath, f.dir.absPkgPath, "", 1)
		}

		results[fpath], err = ioutil.ReadAll(newFile)
		if err != nil {
			return err
		}
	}

	return nil
}
