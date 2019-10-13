package main

import (
	"github.com/samwho/stringshasnfixargs/stringshasnfixargs"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(stringshasnfixargs.Analyzer)
}
