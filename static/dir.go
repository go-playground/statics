package static

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// dir implements the FileSystem interface
type dir struct {
	useStaticFiles bool
	absPkgPath     string
	files          map[string]*file
}

// Open returns the FileSystem DIR
func (d dir) Open(name string) (http.File, error) {

	if d.useStaticFiles {
		f, found := d.files[name]
		if !found {
			return nil, os.ErrNotExist
		}

		return f.File()
	}

	if !strings.HasPrefix(name, d.absPkgPath) {
		name = filepath.FromSlash(d.absPkgPath + name)
	} else {
		name = filepath.FromSlash(name)
	}

	return os.Open(name)
}
