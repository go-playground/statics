package main

import (
	"io/ioutil"
	"os"
	"testing"

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

func TestMain(m *testing.M) {

	// setup

	os.Exit(m.Run())

	// teardown
}

func TestNonExistantStaticDir(t *testing.T) {

	i := "static/test-files/garbagedir"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	ignore := ""
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	PanicMatches(t, func() { main() }, "stat static/test-files/garbagedir: no such file or directory")
}

func TestBadPackage(t *testing.T) {

	i := "static/test-files/test-inner"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := ""
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	init := false
	flagInit = &init

	PanicMatches(t, func() { main() }, "**invalid Package Name")
}

func TestBadOutputDir(t *testing.T) {

	i := "static/test-files/test-inner"
	flagStaticDir = &i

	o := ""
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	init := false
	flagInit = &init

	PanicMatches(t, func() { main() }, "**invalid Output Directory")
}

func TestBadStaticDir(t *testing.T) {

	i := ""
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	init := false
	flagInit = &init

	PanicMatches(t, func() { main() }, "**invalid Static File Directoy '.'")
}

func TestGenerateInitFile(t *testing.T) {

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	init := true
	flagInit = &init

	main()

	b, err := ioutil.ReadFile("static/test-files/test.go")
	Equal(t, err, nil)

	expected := `//go:generate statics -i=static/test-files/teststart -o=static/test-files/test.go -pkg=test -group=Assets

package test

import "github.com/joeybloggs/statics/static"

// newStaticAssets initializes a new *static.Files instance for use
func newStaticAssets(config *static.Config) (*static.Files, error) {

	return static.New(config, &static.DirFile{})
}
`

	Equal(t, string(b), expected)
}

func TestGenerateFile(t *testing.T) {

	Equal(t, true, true)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	ignore := ""
	flagIgnore = &ignore

	prefix := "static/"
	flagPrefix = &prefix

	init := false
	flagInit = &init

	main()
}

func TestIgnore(t *testing.T) {

	Equal(t, true, true)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	ignore := ".*.txt"
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	main()
}

func TestBadIgnore(t *testing.T) {

	Equal(t, true, true)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	g := "Assets"
	flagGroup = &g

	ignore := "([12.gitignore"
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	PanicMatches(t, func() { main() }, "**Error Compiling Regex:error parsing regexp: missing closing ]: `[12.gitignore`")
}

// func TestGenerateInitFile(t *testing.T) {

// 	cmd := exec.Command("./run", "-i=static/test-files/test-inner", "-pkg=test", "-group=Assets", "-o=static/test-files/test.go", "-init=true")
// 	err := cmd.Run()

// 	Equal(t, err, nil)

// 	b, err := ioutil.ReadFile("static/test-files/test.go")

// 	expected := `package test

// import "github.com/joeybloggs/statics/static"

// // NewStaticAssets initializes a new static.Files instance for use
// func NewStaticAssets(config *static.Config) (*static.Files, error) {

// 	return static.New(config, map[string]*static.File{})
// }
// `

// 	Equal(t, string(b), expected)
// }

// func TestGenerateFile(t *testing.T) {

// 	cmd := exec.Command("./run", "-i=static/test-files/test-inner", "-pkg=test", "-group=Assets", "-o=static/test-files/test.go")
// 	err := cmd.Run()

// 	Equal(t, err, nil)

// 	b, err := ioutil.ReadFile("static/test-files/test.go")

// 	expected := `package test

// import (
// 	"os"

// 	"github.com/joeybloggs/statics/static"
// )

// // NewStaticAssets initializes a new static.Files instance for use
// func NewStaticAssets(config *static.Config) (*static.Files, error) {

// 	return static.New(config, map[string]*static.File{
// 		"test-inner/main.js": {
// 			Path:     "test-inner/main.js",
// 			Filename: "main.js",
// 			Filesize: 16500,
// 			FMode:    os.FileMode(420),
// 			Modtime:  1444221312,
// 			Compressed: ` + "`" + `
// H4sIAAAJbogA/9wbbW/bNvrz5Vdwu2ySG1tJO+yL2wTXpdcth63drQX2oSgOtETbbGVJE+k4XpH/fs9D
// UhIlkbKcbsPhCKyxJD7v75QWbAUjQpY8lsHTk5Nwuc1iyfOMhKcT8unkhMA6jZZZ9C+RZ2/pImVxnm43
// Gbkkzda8wL8CAYhZt7QkcBu2nUbsTrIsCT/dT52oooQt6TaVYkoqRE9rPCWT2zLDB/rePTDpZarGBGQb
// VtaMJjxbzUnwg/4VTOtnCZUUHnwAPP9ZcpYm1jO5Lxg8Q+W0YEReSqQ5J7LcMuu+pKVc8yRh2ZwsaSpY
// zXGfYZ/+WgoEbOINA+lkXgIAXkbCXD896epaDCvbqeYWjtN673WeScozhlRPw+BZwm+vgin5RAJU2Ezk
// ASjmOdzR14WAa+QgKuiKveG/M3I/iWiSXKdUiFCpd6Z0FvRo3nLBFUU0oLj+7kcupE3VxpOyFchGwNgz
// sZl900cmjWoRWv1uwSsOYyHUVR8W/aSCxd9XgeWHzZZf8l21q7xqsQDKp/HakNHyTC0r8yxhd1PCJdvY
// kWJhr6krtqUsQ61fDppW0BZDNRjPiq38vsy3RcdYQYxyoqXUltkK9wRomKIANb7NQ6fyXTREQbObLh28
// 6Sc0A7XnWZtcw6qLSLwwiNU2428YhIg5XrP44yK/Q5/jibqzCMhZO0TOiNHxLU23EJ/mSsGyZE6+QN1H
// VpxWDswDsxn9dgG/QiBIsxULptqv8+wHBXD93bW6b0vVVo7tEJZoKV2w1Minf58bCZd5eUgcCUE9V44T
// mWRG/HptUedL8LuO0JPWDlyfeneUO8t1BBAs7BgL1/1J71ZDyuTHrpPXeEEF1Ahfe84q3RdrHkOc1L9m
// iAjtvS7ZEvf8PejZq7ohA6Mf9BdUDqostDVW2zXl8cfGrG+AxjXesvUp1y6RCYOUPiCTjgcn6UOo+9pE
// 5dtQOvF0YO+75q43dmE7gPqmvUkZzOG6kNBBgleMJQyToyprHVxW1kNbjEx5FYFTWdrp1Jf2nO7WLjnP
// Ls1urz+XXn/G1REVi/s4Uw2n/fjGZKEbjxYaTSDZIHAIi+v8nLx8/ctPz9+Sm1f6x83rV+QfNJPrPNs7
// QVBHimwESWZD5cRNXG2FXe/0XtT9e+BEC8TS9KUCbmGa9gA8SnXoS3MWaoVEKgOQLy4vyRb0tISmIxlg
// 0yhJuUvSCrYWL10v0jawWgGLto9xjPY/i5EjdVVZ0U7iHkcyrA06uocOglkZoXQlKsc9DKt+IulCtVNG
// v9E0KEa1JEPgdSLrlUE7vL3JCJNJ6Ws/NYbAoQMIChIiAg6wF0/hzzPyE4UUHjOeNqkxSlm2gjbvvN0s
// TwhAnJ2NqZfoZVgKfwZQgg1DyOGfxxNvlWye16W2mKU8+xhYrU67JCLuXklUinGWrh67A9bRWNoqZUlX
// o22cZgJUw0+Fp0/CQnFvW98zA7UGRFzVWAD60YPKtP1Yu+OcvHvffqAHyO7dyrJz8vhi2nZFS79zq0r4
// yyPI/X3JdbLBiwn4TckyCdqLPIOVjaAuj5cNLjsrwezmMGsFaVFtZTIXNfBTuciTPVROcKaa2ESXx0ZU
// WZqKCD8GSmJFt8pk3rqpJMTMCMyGoeJ6Bi5PHpEmwCAG3LmwQvDPDDWsgBs4D1HMJpUUKrrI1aXFxddf
// k/bTZ5cVicGsbeQV63znz9y9Oy7NmJBxNGv3HW9syrvtjPnig5mhJp/cXpWwmG9oqrKdc4PYbxZ5+jMr
// Y3BV7UjOzga6mq3AgebVdrOATFGUucyxiY9k/pLfsQQiY1nmG/LTi1eR4gweLPFBf4w5f9RX2gu+4lI4
// 7PhanYHQNCJv14xkijrJlyRRAETmBJMNSEKXkuFRDKuFLnKeyacqI5EN3ZMFI1SrC37KHWMZuSAUHOrJ
// xbRPmGdxCiLfsqnawzdFyjaoJHUmoxDmhrd0T8S2KGBEAQIpLbE2lTh+IqOKoIjIzVJzAo+3iIdw0Sea
// b7iULMG2DZ5D5DEKl4QKchE5lPOc6FMv8KSiZKLiDsmCpl4vPgBJKkmSM0GyXBI8RGR3RZ7BTg4agnsa
// ACVcAxV2R2MJ4hjt6j99ug5dpzRmLRsB+zDpZtikQyBmLGZC0HKvtYmgy5LGWn8QdqWSt4CKA/t3XK77
// RH9nZS5auIjItYCgLeQesYoCOFpywKKLuFK8UQZQWCmNIu80I4/Z2ZPHU+MgDPrzpE9VoOH3JAYrC5fz
// v1EGAO9HsXRAg40ru/DMqe+oT+fR+YlneIJQqrR8uP1uAp5YgKPmIwhxKw/MyQuDCkKsue9lsbBSyKWe
// yjzzgY7AS/P3EZReR25Sqnfkpq+CUbKYzKpI1BnKKAOrTBd1J+N102/3WOmP7ghc8KqV0dXpiOKO5qg2
// cyA0N4dqwWSgbZXJO43yEltVi/AZCd63e4ShqoeY1p+N6dD5zRHMDg1XxzHrP19rXfXcpj62+is8JuEl
// 0yQ8PaTIfQ1h1DvNs4Qv2Sa/ZWYSEDMoOrMXLh5qJ63HBjFDzdZ8+XpfdRjodHLpk5UfGxgIRMsSqvYl
// jAF/flusmITnekiVSVADROy30HtcVqukvEE3UMeT2gfCEUJW0FWKzbapI/fjElBj4zVIA6of6nZjCrFY
// vVub+49ZSJPX9fkuHrN4Qq9aC6jFH4e3aPLQxh2gbdMHdQl2k5mzXc3IAU4sbg7u0ywt05yOYarH2EsE
// PJK1Y9hrWAQ/YZJv2EgubU4ztiMvAPwBfB7LKy4z7Y/n0+Z1tLt5uBy73XMa51oqz0TFVqzDTyqS5+pf
// NavNDeOuIzonWV+cdxdWfqsCQA8WPA+GQtvNNaZ/K9vRKVkcg6TihEYgJrkiC/wLCKqO7PFTx3A8HuGz
// HsLZgzDiMgg8zefQGms7e/kK8pQEzmI6SL9qkY7m4g+zcbWGTPNgy7QR950IbP5wxLgebvpqjQ5L1xrw
// hefH+kLN0MPYcZ2IayfR4/PgifeYpZoZ8wpGlnbDrci84+8j1ev0eu7P0K9MYPCV2LZ9lpVMQ9juBs3R
// tuy+qh27HmgnNdfVbrOmon0w/zkWwvbfvGto972/hReTSL10OLa0VuuBwvrPRUcSdRylDEIYL7kPP/x7
// y8r9RL2YOq1TZGhOMtAj9ADWHgpaEx1f1gP4Ur22Ue/wAcS8VLoiF3j03N3T2tA15/k5uU7zjBHzxVPr
// oRodd3+DXxhkO8CX7yb9A81TWW1RhF0bqi+rZIfzGGmHLhAYmjXSFgiMOipwashOFII4z5OEqDddU6LH
// S7Kh5Ypn6lowSXY8AV3godqupAVRej9ps9tjx548JbjtfsYyhAMZ+ntj2OYOGs3JnLhOhWEpzuYkeHxx
// 8VXgOOiPkGH9XlKLePml4QYffIkfFPU//4HgtqIa9FHChLGbgQ4dp2iwtzWT29utabFWhv28+wqv52fq
// gFRNAJprrXr1dSYJF1TweNKxQ6SOg6tP+ToyK4dwCg3Ebpbg0egqVyTWbwgFuvgVnsmyjOwAmWbHsGK+
// aVXn8RnQFNDtQpuQ3NIs7r5p0Do1XqkH++EgGyEM0EdRnM9qhiphR+lZrKEXW6NAWIkhMrOV6Ee3JoGf
// xVoRHOELeNwP+dtWduAKVL3hugpXLwaQcAD+JhOxHISvtTCA5VcM6N4ZV4MDw6TnLZYOmo8IzCdTro0g
// qiM/GECQ0pERrAwGlpivJNS/oDr+Ch0AiMuPR/tcYmFwc6qUGq3lJq2/aMVvrNTnpWftxIpsLXkp5Cxe
// 8xRRK7DJWfDsHLbDPyX+U38V23W6N1Vi7fsY5Nxf1ZPWF8/+b7S8BbbDb7dWkpAfalZsU1vlEY/NILki
// k3WF1VdDZxP3fQPVnLrt1xWkfLgI4IQtLEqCNeOrtaxFMJfDMvi/NatNii/9ehnbbTxLv0ad0tLk8ZRi
// R6y5tbDWH+qELddvnk0afnpe5Df1fT/TlKzIBUdjvWlS5xjHBjl/sd4doagxTeNtqt/DplBFsdZ4T17r
// HUANAZ9X195RBGsBvqPBFh8zX8HwS4dSvyEmVenGl4kUmqmVesm5n9E71+tjUrcSxqmgzllJt/G1AffV
// ZZn3qNdllqR0n2+lP3DwRaLRfuUnaEaQC9oHgIfblaAiLvM0fZsXA3GIElkyWCC+Kt5h5lfsJWReoDEV
// YRTHUvFWqvf0nO1Gp6WuD08G2kl75QWNudzPyWN3b2kv4HhOPJIPAg8dTw2+XavWYZ3x7C/Wl6cXt5fS
// 18XDNeN8clBfdrjgJxatkFHt8h8RL6hs9X1J7QOD4bLrRAlkhHy5hBIPYwEaFefP9p5nvT1n6s4W2rkf
// qmQya9LbkWFX8699B3OYLh1cCpYu/wfjrq2eWVc9f2kA2tqjCxyWYVZYMDCFYy7urv+PMHxI0cdpZ0zN
// dxaYH9lSHq4wlaFS2P0Z5cVqlB5iorHejVy6SoqWdbD9dN4fkRsPaedQIRmY5uyllWT56VDfrwA8ejY6
// uvALPcYTW53fKB9UvaP7M01c2CRWbeiBLrR76FHOUxl+M+kPMUN+TclZM59XXaNHIQdGhiXlKRF0yUjK
// N1Anw5QKCRWgZIyUkME2iBy/zZR5DuksdUzmFjZaLrgs8ds/jY1KchE9+baVJTW/U0ySUHjBGJn68BH+
// NNAFv2Op7vF9OYFi4dzV0j9COsM6u+ztP+p/7DDHolgcNRJr0NINgRNM28rK8sP2Mofr1DFHdQxZHwn0
// Ty6k59DGJJTwNAID51KmLHzy7UXzv1k5nc41rLm4d+T37rb73gB72s8tUZoDgVo6xwlOyQT/nYEUCVvk
// W5zeWlL4Y6elMteGsbKOlNfI3JfAYwcXeXuk1tkNQxr++28AAAD//4gW9tF0QAAA
// ` + "`" + `,
// 		},
// 		"test-inner/symlinked.css": {
// 			Path:     "test-inner/symlinked.css",
// 			Filename: "symlinked.css",
// 			Filesize: 31,
// 			FMode:    os.FileMode(420),
// 			Modtime:  1446483796,
// 			Compressed: ` + "`" + `
// H4sIAAAJbogA/youSMxTqObiTMksLshJrLTKzMvJzEvVTcrJT8625qoFBAAA//8GIFmZHwAAAA==
// ` + "`" + `,
// 		},
// 		"test-inner/test.css": {
// 			Path:     "test-inner/test.css",
// 			Filename: "test.css",
// 			Filesize: 18,
// 			FMode:    os.FileMode(420),
// 			Modtime:  1446483733,
// 			Compressed: ` + "`" + `
// H4sIAAAJbogA/0rKT6ms5uLMTSxKz8yzMrDmqgUEAAD///n2jcUSAAAA
// ` + "`" + `,
// 		},
// 	},
// 	)
// }
// `

// 	Equal(t, string(b), expected)
// }
