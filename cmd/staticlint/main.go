/*
Staticlint - multi-checker for checking code using multiple static analyzers.
Included static analyzers:
	1) all standard golang.org/x/tools/go/analysis/passes;
	2) all SA checks from http://staticcheck.io;
	3) ST1003, ST1006, ST1023 from http://staticcheck.io;
	4) go-critic, body close analyzer;
	5) custom-made ExitAnalyzer, that prevents from using os.Exit in function main inside main package.
How to use:
	main.go [-flag] [package]
The flags are:
	-builtin
		Run standard go lang static analyzer.
	-static
		Run static checks from https://staticheck.io.
	-extra
		Run additional checks from go-critic, body close analyzer.
*/
package main

import (
	"flag"
	"log"

	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/sergalkin/go-url-shortener.git/cmd/staticlint/multicheck"
)

var (
	builtIn = flag.Bool("builtin", false, "Run standard go lang static analyser")
	static  = flag.Bool("static", false, "Run static checks from https://staticheck.io")
	extra   = flag.Bool("extra", false, "Run additional checks from go-critic, bodyclose analyzer")
)

func init() {
	flag.Parse()

	if !*builtIn && !*static && !*extra {
		log.Fatal("No checks defined")
	}
}

func main() {
	l := *multicheck.NewWithOptions(
		multicheck.WithBuiltin(*builtIn),
		multicheck.WithStatic(*static),
		multicheck.WithExtra(*extra),
	)

	multichecker.Main(l...)
}
