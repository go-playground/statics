/*
Package static reads package statics generated go files and provides helper methods and objects to retrieve the embeded files and even serve then via http.FileSystem.

Embedding in Source Control

	statics -i=assets -o=assets.go -pkg=main -group=Assets -ignore=//.gitignore -init=true

	when using arg init=true statics package generates a minimal configuration with no files embeded;
	you can then add it to source control, add the file to .gitignore and then when ready for generation
	just run go generate from the project root and the files will get embedded ready for compilation.

Example Usages

	// generated via command: statics -i=assets -o=assets.go -pkg=main -group=Assets -ignore=//.gitignore

	pkg = // get absolute directory path of the -i arguments parent directory, used when UseStaticFiles=false

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
