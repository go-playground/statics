package static

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

type httpFile struct {
	*bytes.Reader
	*file
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

// File returns an http.File or error
func (f *file) File() (http.File, error) {

	// if production read filesystem file
	return &httpFile{
		bytes.NewReader(f.data),
		f,
	}, nil
}

// Close closes the File, rendering it unusable for I/O. It returns an error, if any.
func (f *file) Close() error {
	return nil
}

// Readdir returns nil fileinfo and an error because the static FileSystem does not store directories
func (f *file) Readdir(count int) ([]os.FileInfo, error) {

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
		count += f.lastDirIndex
	}

	if f.lastDirIndex >= len(f.files) {
		f.lastDirIndex = 0
		return nil, io.EOF
	}

	if count+f.lastDirIndex >= len(f.files) {
		count = len(f.files)
	}

	i := f.lastDirIndex

	var j int

	for i = f.lastDirIndex; i < count; i++ {
		files[j] = f.files[i]
		j++
	}

	if count > 0 {
		f.lastDirIndex += j
	}

	return files[:j], nil
}

// Stat returns the FileInfo structure describing file. If there is an error, it will be of type *PathError.
func (f *file) Stat() (os.FileInfo, error) {
	return f, nil
}

// Name returns the name of the file as presented to Open.
func (f *file) Name() string {
	return f.name
}

// Size length in bytes for regular files; system-dependent for others
func (f *file) Size() int64 {
	return f.size
}

// Mode returns file mode bits
func (f *file) Mode() os.FileMode {
	return os.FileMode(f.mode)
}

// ModTime returns the files modification time
func (f *file) ModTime() time.Time {
	return time.Unix(f.modTime, 0)
}

// IsDir reports whether f describes a directory.
func (f *file) IsDir() bool {
	return f.isDir
}

// Sys returns the underlying data source (can return nil)
func (f *file) Sys() interface{} {
	return f
}
