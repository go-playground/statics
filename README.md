Package statics
===============

[![Build Status](https://semaphoreci.com/api/v1/projects/1b97afa9-77f3-43ff-ad26-749958500745/601363/badge.svg)](https://semaphoreci.com/joeybloggs/statics)
[![Go Report Card](http://goreportcard.com/badge/go-playground/statics)](http://goreportcard.com/report/go-playground/statics)
[![GoDoc](https://godoc.org/github.com/go-playground/statics/static?status.svg)](https://godoc.org/github.com/go-playground/statics/static)

Package statics embeds static files into your go applications. It provides helper methods and objects to retrieve embeded files and serve via http.

It has the following **unique** features:

-   ```Follows Symlinks!``` uses the symlinked file/directory name and path but the contents and other fileinfo of the original file.
-   Embeds static generation command for using ```go:generate``` in each static file for easy generation in production mode.
-   Handles multiple static go files, even inside of the same package.
-   Handles development (aka local files) vs production (aka static files) across packages.

Installation
------------
Use go get.

	go get github.com/go-playground/statics

or to update

	go get -u github.com/go-playground/statics

Then import the validator package into your own code.

	import "github.com/go-playground/statics"

Usage and documentation
------

Please see https://godoc.org/github.com/go-playground/statics/static for detailed usage docs.

NOTE: when specifying path or directory name in code always use "/", even for you windows users,
     the package handles any conversion to you local filesystem paths; Except for the AbsPkgPath
     variable in the config.

run statics -h to see the options/arguments

##### Examples:

Embedding in Source Control

	statics -i=assets -o=assets.go -pkg=main -group=Assets -ignore=//.gitignore -init=true

	when using arg init=true statics package generates a minimal configuration with no 
	files embeded; you can then add it to source control, ignore the file locally using
	git update-index --assume-unchanged [filename(s)] and then when ready for generation 
	just run go generate from the project root and the files will get embedded ready for 
	compilation.

	Be sure to check out this packages best buddy https://github.com/go-playground/generate
	to help get everything generated and ready for compilation.

Example Usage
```go
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

	// NOTE: Assets in the function name below is the group in the generation command
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
