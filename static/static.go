package static

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

// Files contains a full instance of a static file collection
type Files struct {
	dir Dir
}

// File contains the static FileInfo
type File struct {
	data       []byte
	Path       string
	Filename   string
	Filesize   int64
	FMode      os.FileMode
	Modtime    int64
	Compressed string
	isDir      bool
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
func (f *File) File() (http.File, error) {

	// if production read filesystem file
	return &httpFile{
		Reader: bytes.NewReader(f.data),
		File:   f,
	}, nil
}

// Close closes the File, rendering it unusable for I/O. It returns an error, if any.
func (f *File) Close() error {
	return nil
}

// Readdir returns nil fileinfo and an error because the static FileSystem does not store directories
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("not a directory")
}

// Stat returns the FileInfo structure describing file. If there is an error, it will be of type *PathError.
func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}

// Name returns the name of the file as presented to Open.
func (f *File) Name() string {
	return f.Filename
}

// Size length in bytes for regular files; system-dependent for others
func (f *File) Size() int64 {
	return f.Filesize
}

// Mode returns file mode bits
func (f *File) Mode() os.FileMode {
	mode := os.FileMode(0644)
	if f.IsDir() {
		return mode | os.ModeDir
	}
	return mode
	// return 0
}

// ModTime returns the files modification time
func (f *File) ModTime() time.Time {
	return time.Unix(f.Modtime, 0)
}

// IsDir reports whether f describes a directory.
func (f *File) IsDir() bool {
	return f.isDir
}

// Sys returns the underlying data source (can return nil)
func (f *File) Sys() interface{} {
	return f
}

// New create a new static file instance.
func New(config *Config, files map[string]*File) (*Files, error) {

	if config.IsProductionMode {
		var err error
		var reader *gzip.Reader
		var b64 io.Reader

		for _, f := range files {

			if f.Filesize == 0 {
				continue
			}

			b64 = base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.Compressed))
			reader, err = gzip.NewReader(b64)
			if err != nil {
				return nil, err
			}

			f.data, err = ioutil.ReadAll(reader)
			if err != nil {
				return nil, err
			}

			// original string should get garbage collected
			f.Compressed = ""
		}
	}

	return &Files{
		dir: Dir{
			name:             config.Name,
			isProductionMode: config.IsProductionMode,
			files:            files,
		},
	}, nil
}

// FS returns an http.FileSystem object for serving files over http
func (s *Files) FS() http.FileSystem {
	return s.dir
}

// GetFile return a files contents as []byte from the filesystem, static or local
func (s *Files) GetFile(name string) ([]byte, error) {
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
