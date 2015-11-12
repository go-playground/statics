package static

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	. "gopkg.in/go-playground/assert.v1"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//

var testDirFile *DirFile

func getGOPATH() string {
	gopath := os.Getenv("GOPATH")

	if len(gopath) == 0 {
		panic("$GOPATH environment is not setup correctly, ending transmission!")
	}

	return gopath
}

func TestMain(m *testing.M) {

	// setup
	testDirFile = &DirFile{
		Path:    "/static/test-files/teststart",
		Name:    "teststart",
		Size:    170,
		Mode:    os.FileMode(2147484141),
		ModTime: 1446650128,
		IsDir:   true,
		Compressed: `
`,
		Files: []*DirFile{{
			Path:    "/static/test-files/teststart/symlinkeddir",
			Name:    "symlinkeddir",
			Size:    136,
			Mode:    os.FileMode(2147484141),
			ModTime: 1446648191,
			IsDir:   true,
			Compressed: `
`,
			Files: []*DirFile{{
				Path:    "/static/test-files/teststart/symlinkeddir/realdir",
				Name:    "realdir",
				Size:    136,
				Mode:    os.FileMode(2147484141),
				ModTime: 1446648584,
				IsDir:   true,
				Compressed: `
`,
				Files: []*DirFile{{
					Path:    "/static/test-files/teststart/symlinkeddir/realdir/doublesymlinkeddir",
					Name:    "doublesymlinkeddir",
					Size:    136,
					Mode:    os.FileMode(2147484141),
					ModTime: 1447163695,
					IsDir:   true,
					Compressed: `
`,
					Files: []*DirFile{{
						Path:    "/static/test-files/teststart/symlinkeddir/realdir/doublesymlinkeddir/doublesymlinkedfile.txt",
						Name:    "doublesymlinkedfile.txt",
						Size:    5,
						Mode:    os.FileMode(420),
						ModTime: 1446648265,
						IsDir:   false,
						Compressed: `
H4sIAAAJbogA/0pJLEnkAgQAAP//gsXB5gUAAAA=
`,
						Files: []*DirFile{},
					},
						{
							Path:    "/static/test-files/teststart/symlinkeddir/realdir/doublesymlinkeddir/triplesymlinkeddir",
							Name:    "triplesymlinkeddir",
							Size:    102,
							Mode:    os.FileMode(2147484141),
							ModTime: 1447163709,
							IsDir:   true,
							Compressed: `
`,
							Files: []*DirFile{{
								Path:    "/static/test-files/teststart/symlinkeddir/realdir/doublesymlinkeddir/triplesymlinkeddir/triplefile.txt",
								Name:    "triplefile.txt",
								Size:    5,
								Mode:    os.FileMode(420),
								ModTime: 1447163511,
								IsDir:   false,
								Compressed: `
H4sIAAAJbogA/0pJLEnkAgQAAP//gsXB5gUAAAA=
`,
								Files: []*DirFile{},
							},
							},
						},
					},
				},
					{
						Path:    "/static/test-files/teststart/symlinkeddir/realdir/realdirfile.txt",
						Name:    "realdirfile.txt",
						Size:    5,
						Mode:    os.FileMode(420),
						ModTime: 1446648207,
						IsDir:   false,
						Compressed: `
H4sIAAAJbogA/0pJLEnkAgQAAP//gsXB5gUAAAA=
`,
						Files: []*DirFile{},
					},
				},
			},
				{
					Path:    "/static/test-files/teststart/symlinkeddir/symlinkeddirfile.txt",
					Name:    "symlinkeddirfile.txt",
					Size:    5,
					Mode:    os.FileMode(420),
					ModTime: 1446647769,
					IsDir:   false,
					Compressed: `
H4sIAAAJbogA/0pJLEnkAgQAAP//gsXB5gUAAAA=
`,
					Files: []*DirFile{},
				},
			},
		},
			{
				Path:    "/static/test-files/teststart/plainfile.txt",
				Name:    "plainfile.txt",
				Size:    10,
				Mode:    os.FileMode(420),
				ModTime: 1446650128,
				IsDir:   false,
				Compressed: `
H4sIAAAJbogA/ypIzMnMS0ksSeQCBAAA//9+mKzVCgAAAA==
`,
				Files: []*DirFile{},
			},
			{
				Path:    "/static/test-files/teststart/symlinkedfile.txt",
				Name:    "symlinkedfile.txt",
				Size:    5,
				Mode:    os.FileMode(420),
				ModTime: 1446647746,
				IsDir:   false,
				Compressed: `
H4sIAAAJbogA/0pJLEnkAgQAAP//gsXB5gUAAAA=
`,
				Files: []*DirFile{},
			},
		},
	}

	os.Exit(m.Run())

	// teardown
}

func TestStaticNew(t *testing.T) {

	config := &Config{
		UseStaticFiles: true,
		AbsPkgPath:     getGOPATH() + "/src/github.com/go-playground/statics",
	}

	staticFiles, err := New(config, testDirFile)
	Equal(t, err, nil)
	NotEqual(t, staticFiles, nil)

	go func(sf *Files) {

		http.Handle("/static/", http.StripPrefix("/", http.FileServer(sf.FS())))
		http.ListenAndServe("127.0.0.1:3006", nil)

	}(staticFiles)

	time.Sleep(1000)

	f, err := staticFiles.GetHTTPFile("/static/test-files/teststart/plainfile.txt")
	Equal(t, err, nil)

	fis, err := f.Readdir(-1)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "not a directory")

	fi, err := f.Stat()
	Equal(t, err, nil)
	Equal(t, fi.Name(), "plainfile.txt")
	Equal(t, fi.Size(), int64(10))
	Equal(t, fi.IsDir(), false)
	Equal(t, fi.Mode(), os.FileMode(420))
	Equal(t, fi.ModTime(), time.Unix(1446650128, 0))
	NotEqual(t, fi.Sys(), nil)

	err = f.Close()
	Equal(t, err, nil)

	f, err = staticFiles.GetHTTPFile("/static/test-files/teststart")
	Equal(t, err, nil)

	fi, err = f.Stat()
	Equal(t, err, nil)
	Equal(t, fi.Name(), "teststart")
	Equal(t, fi.Size(), int64(170))
	Equal(t, fi.IsDir(), true)
	Equal(t, fi.Mode(), os.FileMode(2147484141))
	Equal(t, fi.ModTime(), time.Unix(1446650128, 0))
	NotEqual(t, fi.Sys(), nil)

	var j int

	for err != io.EOF {

		fis, err = f.Readdir(2)

		switch j {
		case 0:
			Equal(t, len(fis), 2)
		case 1:
			Equal(t, len(fis), 1)
		case 2:
			Equal(t, len(fis), 0)
		}

		j++
	}

	err = f.Close()
	Equal(t, err, nil)

	b, err := staticFiles.ReadFile("/static/test-files/teststart/plainfile.txt")
	Equal(t, err, nil)
	Equal(t, string(b), "palindata\n")

	b, err = staticFiles.ReadFile("nonexistantfile")
	NotEqual(t, err, nil)

	bs, err := staticFiles.ReadFiles("/static/test-files/teststart", false)
	Equal(t, err, nil)
	Equal(t, len(bs), 2)
	Equal(t, string(bs["/static/test-files/teststart/plainfile.txt"]), "palindata\n")
	Equal(t, string(bs["/static/test-files/teststart/symlinkedfile.txt"]), "data\n")

	bs, err = staticFiles.ReadFiles("/static/test-files/teststart", true)
	Equal(t, err, nil)
	Equal(t, len(bs), 6)

	bs, err = staticFiles.ReadFiles("nonexistantdir", false)
	NotEqual(t, err, nil)

	fis, err = staticFiles.ReadDir("/static/test-files/teststart")
	Equal(t, err, nil)
	Equal(t, len(fis), 3)
	Equal(t, fis[0].Name(), "plainfile.txt")
	Equal(t, fis[1].Name(), "symlinkeddir")
	Equal(t, fis[2].Name(), "symlinkedfile.txt")

	fis, err = staticFiles.ReadDir("nonexistantdir")
	NotEqual(t, err, nil)

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://127.0.0.1:3006/static/test-files/teststart/plainfile.txt", nil)
	Equal(t, err, nil)

	resp, err := client.Do(req)
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, http.StatusOK)

	bytes, err := ioutil.ReadAll(resp.Body)
	Equal(t, err, nil)
	Equal(t, string(bytes), "palindata\n")

	defer resp.Body.Close()
}

func TestLocalNew(t *testing.T) {

	config := &Config{
		UseStaticFiles: false,
		AbsPkgPath:     getGOPATH() + "/src/github.com/go-playground/statics",
	}

	staticFiles, err := New(config, testDirFile)
	Equal(t, err, nil)
	NotEqual(t, staticFiles, nil)

	go func(sf *Files) {

		http.Handle("/static/test-files/", http.StripPrefix("/", http.FileServer(sf.FS())))
		http.ListenAndServe("127.0.0.1:3007", nil)

	}(staticFiles)

	time.Sleep(1000)

	f, err := staticFiles.GetHTTPFile("/static/test-files/teststart/plainfile.txt")
	Equal(t, err, nil)

	fis, err := f.Readdir(-1)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "readdirent: invalid argument")

	fi, err := f.Stat()
	Equal(t, err, nil)
	Equal(t, fi.Name(), "plainfile.txt")
	Equal(t, fi.Size(), int64(10))
	Equal(t, fi.IsDir(), false)
	Equal(t, fi.Mode(), os.FileMode(420))
	Equal(t, fi.ModTime(), time.Unix(1446650128, 0))
	NotEqual(t, fi.Sys(), nil)

	err = f.Close()
	Equal(t, err, nil)

	f, err = staticFiles.GetHTTPFile("/static/test-files/teststart")
	Equal(t, err, nil)

	fi, err = f.Stat()
	Equal(t, err, nil)
	Equal(t, fi.Name(), "teststart")
	Equal(t, fi.Size(), int64(170))
	Equal(t, fi.IsDir(), true)
	Equal(t, fi.Mode(), os.FileMode(2147484141))
	Equal(t, fi.ModTime(), time.Unix(1446650128, 0))
	NotEqual(t, fi.Sys(), nil)

	var j int

	for err != io.EOF {

		fis, err = f.Readdir(2)

		switch j {
		case 0:
			Equal(t, len(fis), 2)
		case 1:
			Equal(t, len(fis), 1)
		case 2:
			Equal(t, len(fis), 0)
		}

		j++
	}

	err = f.Close()
	Equal(t, err, nil)

	b, err := staticFiles.ReadFile("/static/test-files/teststart/plainfile.txt")
	Equal(t, err, nil)
	Equal(t, string(b), "palindata\n")

	b, err = staticFiles.ReadFile("nonexistantfile")
	NotEqual(t, err, nil)

	bs, err := staticFiles.ReadFiles("/static/test-files/teststart", false)
	Equal(t, err, nil)
	Equal(t, len(bs), 2)
	Equal(t, string(bs["/static/test-files/teststart/plainfile.txt"]), "palindata\n")
	Equal(t, string(bs["/static/test-files/teststart/symlinkedfile.txt"]), "data\n")

	bs, err = staticFiles.ReadFiles("/static/test-files/teststart", true)
	Equal(t, err, nil)
	Equal(t, len(bs), 6)

	bs, err = staticFiles.ReadFiles("nonexistantdir", false)
	NotEqual(t, err, nil)

	fis, err = staticFiles.ReadDir("/static/test-files/teststart")
	Equal(t, err, nil)
	Equal(t, len(fis), 3)
	Equal(t, fis[0].Name(), "plainfile.txt")
	Equal(t, fis[1].Name(), "symlinkeddir")
	Equal(t, fis[2].Name(), "symlinkedfile.txt")

	fis, err = staticFiles.ReadDir("nonexistantdir")
	NotEqual(t, err, nil)

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://127.0.0.1:3007/static/test-files/teststart/plainfile.txt", nil)
	Equal(t, err, nil)

	resp, err := client.Do(req)
	Equal(t, err, nil)
	Equal(t, resp.StatusCode, http.StatusOK)

	bytes, err := ioutil.ReadAll(resp.Body)
	Equal(t, err, nil)
	Equal(t, string(bytes), "palindata\n")

	defer resp.Body.Close()
}

func TestBadLocalAbsPath(t *testing.T) {

	config := &Config{
		UseStaticFiles: false,
		AbsPkgPath:     "../github.com/go-playground/statics",
	}

	staticFiles, err := New(config, testDirFile)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "AbsPkgPath is required when not using static files otherwise the static package has no idea where to grab local files from when your package is used from within another package.")
	Equal(t, staticFiles, nil)
}
