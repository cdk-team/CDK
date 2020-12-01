package main

import (
	"github.com/Xyntax/CDK/pkg/lib"
	_ "github.com/Xyntax/CDK/pkg/exploit" // register all scripts
	_ "github.com/Xyntax/CDK/pkg/evaluate" // register all scripts
)

func main() {
	lib.ParseDocopt()
}
