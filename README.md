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

Generation + Usage

License
------
Distributed under MIT License, please see license file in code for more details.