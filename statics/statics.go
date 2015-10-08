package statics

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
)

// FileInfo contains all necessary information about the sttically compile file.
type FileInfo struct {
	Contents string
}

// Config contains information about how extracting the data should behave
type Config struct {
	IsProductionMode bool
}

// Files contains an instance of static files & functions
type Files struct {
	isProductionMode bool
	files            map[string]*FileInfo
}

// New create a new static file instance.
func New(config *Config, m map[string]*FileInfo) *Files {

	for _, v := range m {

		b, err := base64.StdEncoding.DecodeString(v.Contents)
		if err != nil {
			log.Fatal(err)
		}

		in := bytes.NewBuffer(b)
		var buff bytes.Buffer
		r, err := gzip.NewReader(in)
		if err != nil {
			log.Fatal(err)
		}

		defer r.Close()
		_, err = io.Copy(&buff, r)
		if err != nil {
			log.Fatal(err)
		}

		v.Contents = buff.String()
	}

	return &Files{
		isProductionMode: config.IsProductionMode,
		files:            m,
		// decompressOnRetirieve: config.DecompressOnRetirieve,
	}
}

// Get returns the file as a byte array either from static contents or from disk depending on the Files object settigns.
func (f *Files) Get(file string) []byte {

	if !f.isProductionMode {
		return getFileFromDisk(file)
	}

	b, ok := f.files[file]
	if !ok {
		log.Fatalf("File %s Does Not Exists", file)
	}

	return []byte(b.Contents)
}

// GetString returns the file as a string either from static contents or from disk depending on the Files object settigns.
func (f *Files) GetString(file string) string {

	if !f.isProductionMode {
		return string(getFileFromDisk(file))
	}

	b, ok := f.files[file]
	if !ok {
		log.Fatalf("File %s Does Not Exists", file)
	}

	return b.Contents
}

func getFileFromDisk(file string) []byte {

	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return b
}
