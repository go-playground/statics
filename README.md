Package statics
===============

[![GoDoc](https://godoc.org/github.com/joeybloggs/statics/static?status.svg)](https://godoc.org/github.com/joeybloggs/statics/static)

Package statics embeds static files into go and provides helper methods and objects to retrieve the embeded files and even serve then via http.FileSystem.

It has the following **unique** features:

-   Follows Symlinks! and uses the symlinked file name and path but the contents + fileinfo of the original file.
-   Embeds original static generation command using go:generate in each static file for easy generation in production modes.
-   Handles multiple static go files, even inside of the same package.
-   Handles development (aka local files) vs production (aka static files) across packages.

Installation
------------

Use go get.

	go get github.com/joeybloggs/statics

or to update

	go get -u github.com/joeybloggs/statics

Then import the validator package into your own code.

	import "github.com/joeybloggs/statics"

Usage and documentation
------

Please see https://godoc.org/github.com/joeybloggs/statics/static for detailed usage docs.

##### Examples:

Embedding in Source Control

	statics -i=assets -o=assets.go -pkg=main -group=Assets -ignore=//.gitignore -init=true

	when using arg init=true statics package generates a minimal configuration with no files embeded;
	you can then add it to source control, add the file to .gitignore and then when ready for generation
	just run go generate from the project root and the files will get embedded ready for compilation.

Example Usage
```go
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
```

License
------
Distributed under MIT License, please see license file in code for more details.