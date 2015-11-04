package static

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
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
	absPkgPath string
	dir        Dir
}

// File contains the static FileInfo
type file struct {
	data         []byte
	path         string
	name         string
	size         int64
	mode         os.FileMode
	modTime      int64
	isDir        bool
	files        []*file
	lastDirIndex int
}

// Dir implements the FileSystem interface
type Dir struct {
	useStaticFiles bool
	files          map[string]*file
}

type httpFile struct {
	*bytes.Reader
	*file
}

// Config contains information about how extracting the data should behave
type Config struct {
	UseStaticFiles bool
	AbsPkgPath     string // the Absolute package path used for local file reading when UseStaticFiles is false
}

// Open returns the FileSystem DIR
func (dir Dir) Open(name string) (http.File, error) {

	if dir.useStaticFiles {
		f, found := dir.files[path.Clean(name)]
		if !found {
			return nil, os.ErrNotExist
		}

		return f.File()
	}

	return os.Open(name)
}

// File returns an http.File or error
func (f file) File() (http.File, error) {

	// if production read filesystem file
	return &httpFile{
		bytes.NewReader(f.data),
		&f,
	}, nil
}

// Close closes the File, rendering it unusable for I/O. It returns an error, if any.
func (f file) Close() error {
	return nil
}

// Readdir returns nil fileinfo and an error because the static FileSystem does not store directories
func (f file) Readdir(count int) ([]os.FileInfo, error) {

	if !f.IsDir() {
		return nil, errors.New("not a directory")
	}

	var files []os.FileInfo

	if count <= 0 {
		files = make([]os.FileInfo, len(f.files))
		count = len(f.files)
		f.lastDirIndex = 0
	} else {
		files = make([]os.FileInfo, count)
	}

	if f.lastDirIndex >= len(f.files) {
		return nil, io.EOF
	}

	if count+f.lastDirIndex >= len(f.files) {
		count = len(f.files)
	}

	for i := f.lastDirIndex; i < count; i++ {
		files = append(files, *f.files[i])
	}

	return files, nil
}

// Stat returns the FileInfo structure describing file. If there is an error, it will be of type *PathError.
func (f file) Stat() (os.FileInfo, error) {
	return f, nil
}

// Name returns the name of the file as presented to Open.
func (f file) Name() string {
	return f.name
}

// Size length in bytes for regular files; system-dependent for others
func (f file) Size() int64 {
	return f.size
}

// Mode returns file mode bits
func (f file) Mode() os.FileMode {
	mode := os.FileMode(0644)
	if f.IsDir() {
		return mode | os.ModeDir
	}
	return mode
}

// ModTime returns the files modification time
func (f file) ModTime() time.Time {
	return time.Unix(f.modTime, 0)
}

// IsDir reports whether f describes a directory.
func (f file) IsDir() bool {
	return f.isDir
}

// Sys returns the underlying data source (can return nil)
func (f file) Sys() interface{} {
	return f
}

// New create a new static file instance.
func New(config *Config, dirFile *DirFile) (*Files, error) {

	files := map[string]*file{}

	if config.UseStaticFiles {
		processFiles(files, dirFile)
	} else {
		if len(config.AbsPkgPath) == 0 {
			return nil, errors.New("AbsPkgPath is required when not using static files otherwise the static package has no idea where to grab local files from when your package is used from within another package.")
		}
	}

	return &Files{
		absPkgPath: config.AbsPkgPath,
		dir: Dir{
			useStaticFiles: config.UseStaticFiles,
			files:          files,
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

	files[f.path] = f

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

// GetHTTPFile returns an http.File object
func (f *Files) GetHTTPFile(name string) (http.File, error) {
	return f.dir.Open(name)
}

// ReadFile returns a files contents as []byte from the filesystem, static or local
func (f *Files) ReadFile(path string) ([]byte, error) {

	file, err := f.dir.Open(path)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(file)
}

// ADD a READ All Dir Files instead of All File

// // ReadFile returns a files contents as []byte from the filesystem, static or local
// func (f *Files) ReadAllFile() ([]byte, error) {

// 	file, err := f.dir.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ioutil.ReadAll(file)
// }
