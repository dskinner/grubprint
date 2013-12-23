package main

import (
	"dasa.cc/food/usda"
	"flag"
)

var (
	dbinit = flag.Bool("dbinit", false, "init the database")
)

func init() {
	flag.Parse()
}

func main() {
	if *dbinit {
		usda.LoadAll()
	}
}
