package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"lesiw.io/boolset"
)

func main() { singlechecker.Main(boolset.Analyzer) }
