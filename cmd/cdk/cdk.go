package main

import (
	"github.com/cdk-team/CDK/pkg/cli"
	_ "github.com/cdk-team/CDK/pkg/exploit" // register all exploits
	_ "github.com/cdk-team/CDK/pkg/task"    // register all task
)

func main() {
	cli.ParseCDKMain()
}
