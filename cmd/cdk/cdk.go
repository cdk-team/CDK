package main

import (
	_ "github.com/Xyntax/CDK/pkg/exploit"  // register all exploits
	_ "github.com/Xyntax/CDK/pkg/task" // register all task
	"github.com/Xyntax/CDK/pkg/cli"
)

func main() {
	cli.ParseCDKMain()
}
