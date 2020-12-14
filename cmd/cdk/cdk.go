package main

import (
	_ "github.com/Xyntax/CDK/pkg/evaluate" // register all scripts
	_ "github.com/Xyntax/CDK/pkg/exploit"  // register all scripts
	"github.com/Xyntax/CDK/pkg/lib"
)

func main() {
	lib.ParseCDKMain()
}
