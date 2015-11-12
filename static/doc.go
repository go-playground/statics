/*
Package static reads package statics generated go files and provides helper methods and objects to retrieve the embeded files and even serve then via http.FileSystem.

Embedding in Source Control

	statics -i=assets -o=assets.go -pkg=main -group=Assets -ignore=//.gitignore -init=true

	when using arg init=true statics package generates a minimal configuration with no
	files embeded; you can then add it to source control, add the file to .gitignore and
	then when ready for generation just run go generate from the project root and the
	files will get embedded ready for compilation.


NOTE: when specifying paths or directory name in code always use "/", even for you windows users,
      the package handles any conversion to you local filesystem paths; Except for the AbsPkgPath
      variable in the config.

run statics -h to see the options/arguments

Example Usages

	// generated via command:
	// statics -i=assets -o=assets.go -pkg=main -group=Assets -ignore=//.gitignore

	gopath := getGopath() // retrieved from environment variable
	pkgPath := "/src/github.com/username/project"

	// get absolute directory path of the -i arguments parent directory + any prefix
	// removed, used when UseStaticFiles=false this is so even when referencing this
	// package from another project and your PWD is not for this package anymore the
	// file paths will still work.
	pkg := goapth + pkgPath

	config := &static.Config{
		UseStaticFiles: true,
		AbsPkgPath:     pkg,
	}

	assets, err := newStaticAssets(config)
	if err != nil {
		log.Println(err)
	}

	// when using http
	http.Handle("/assets", http.FileServer(assets.FS()))

	// other methods for direct access
	assets.GetHTTPFile
	assets.ReadFile
	assets.ReadDir
	assets.ReadFiles
*/
package static
