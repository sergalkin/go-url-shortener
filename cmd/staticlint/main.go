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
