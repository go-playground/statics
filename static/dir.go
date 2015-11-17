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
	fallbackToDisk bool
	absPkgPath     string
	files          map[string]*file
}

// Open returns the FileSystem DIR
func (d dir) Open(name string) (http.File, error) {

	if d.useStaticFiles {
		f, found := d.files[name]

		if found {
			return f.File()
		}

		if !d.fallbackToDisk {
			return nil, os.ErrNotExist
		}
	}

	if !strings.HasPrefix(name, d.absPkgPath) {
		name = filepath.FromSlash(d.absPkgPath + name)
	} else {
		name = filepath.FromSlash(name)
	}

	return os.Open(name)
}
