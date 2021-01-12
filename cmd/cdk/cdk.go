package main

import (
	_ "github.com/cdk-team/CDK/pkg/exploit"  // register all exploits
	_ "github.com/cdk-team/CDK/pkg/task" // register all task
	"github.com/cdk-team/CDK/pkg/cli"
)

func main() {
	cli.ParseCDKMain()
}
