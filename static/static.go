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
	dir Dir
}

// File contains the static FileInfo
type File struct {
	data         []byte
	path         string
	name         string
	size         int64
	mode         os.FileMode
	modTime      int64
	isDir        bool
	files        []*File
	lastDirIndex int
}

// Dir implements the FileSystem interface
type Dir struct {
	name             string
	isProductionMode bool
	files            map[string]*File
}

type httpFile struct {
	*bytes.Reader
	*File
}

// Config contains information about how extracting the data should behave
type Config struct {
	IsProductionMode bool
	Name             string
}

// Open returns the FileSystem DIR
func (dir Dir) Open(name string) (http.File, error) {

	if dir.isProductionMode {
		f, found := dir.files[path.Clean(name)]
		if !found {
			return nil, os.ErrNotExist
		}

		return f.File()
	}

	return os.Open(name)
}

// File returns an http.File or error
func (f File) File() (http.File, error) {

	// if production read filesystem file
	return &httpFile{
		Reader: bytes.NewReader(f.data),
		File:   &f,
	}, nil
}

// Close closes the File, rendering it unusable for I/O. It returns an error, if any.
func (f File) Close() error {
	return nil
}

// Readdir returns nil fileinfo and an error because the static FileSystem does not store directories
func (f File) Readdir(count int) ([]os.FileInfo, error) {

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
func (f File) Stat() (os.FileInfo, error) {
	return f, nil
}

// Name returns the name of the file as presented to Open.
func (f File) Name() string {
	return f.name
}

// Size length in bytes for regular files; system-dependent for others
func (f File) Size() int64 {
	return f.size
}

// Mode returns file mode bits
func (f File) Mode() os.FileMode {
	mode := os.FileMode(0644)
	if f.IsDir() {
		return mode | os.ModeDir
	}
	return mode
}

// ModTime returns the files modification time
func (f File) ModTime() time.Time {
	return time.Unix(f.modTime, 0)
}

// IsDir reports whether f describes a directory.
func (f File) IsDir() bool {
	return f.isDir
}

// Sys returns the underlying data source (can return nil)
func (f File) Sys() interface{} {
	return f
}

// New create a new static file instance.
func New(config *Config, file *DirFile) (*Files, error) {
	files := map[string]*File{}

	if config.IsProductionMode {
		processFiles(files, file)
	}

	return &Files{
		dir: Dir{
			name:             config.Name,
			isProductionMode: config.IsProductionMode,
			files:            files,
		},
	}, nil
}

func processFiles(files map[string]*File, file *DirFile) *File {

	f := &File{
		path:    file.Path,
		name:    file.Name,
		size:    file.Size,
		mode:    file.Mode,
		modTime: file.ModTime,
		isDir:   file.IsDir,
		files:   []*File{},
	}

	files[f.path] = f

	if file.IsDir {
		for _, dirFile := range file.Files {
			resultFile := processFiles(files, dirFile)
			f.files = append(f.files, resultFile)
		}

		return f
	}

	// decompress file contents
	b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(file.Compressed))
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
func (s *Files) FS() http.FileSystem {
	return s.dir
}

// GetFile returns an http.File object
func (s *Files) GetFile(name string) (http.File, error) {
	return s.dir.Open(name)
}

// GetFileBytes return a files contents as []byte from the filesystem, static or local
func (s *Files) GetFileBytes(name string) ([]byte, error) {
	f, err := s.dir.Open(name)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

// GetFileString return a files contents as string from the filesystem, static or local
func (s *Files) GetFileString(name string) (string, error) {
	f, err := s.dir.Open(name)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
