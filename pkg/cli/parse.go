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

package cli

import (
	"fmt"
	"github.com/cdk-team/CDK/pkg/tool/netstat"

	"github.com/cdk-team/CDK/pkg/evaluate"
	"github.com/cdk-team/CDK/pkg/plugin"
	"github.com/cdk-team/CDK/pkg/tool/dockerd_api"
	"github.com/cdk-team/CDK/pkg/tool/etcdctl"
	"github.com/cdk-team/CDK/pkg/tool/kubectl"

	"log"
	"os"
	"strconv"

	"github.com/cdk-team/CDK/pkg/tool/netcat"
	"github.com/cdk-team/CDK/pkg/tool/network"
	"github.com/cdk-team/CDK/pkg/tool/probe"
	"github.com/cdk-team/CDK/pkg/tool/ps"
	"github.com/cdk-team/CDK/pkg/tool/vi"
	"github.com/docopt/docopt-go"
)

func PassInnerArgs() {
	os.Args = os.Args[1:]
}

func ParseCDKMain() bool {

	if len(os.Args) == 1 {
		docopt.PrintHelpAndExit(nil, BannerContainer)
	}

	// nc needs -v and -h , parse it outside
	if os.Args[1] == "nc" {
		// https://github.com/jiguangin/netcat
		PassInnerArgs()
		netcat.RunVendorNetcat()
		return true
	}

	// docopt argparse start
	parseDocopt()

	if Args["auto-escape"].(bool) {
		plugin.RunSingleTask("auto-escape")
		return true
	}

	// support for cdk eva(Evangelion) and cdk evaluate
	fok := Args["evaluate"]
	ok := Args["eva"]

	// docopt let fok = true, so we need to check it
	// fix #37 https://github.com/cdk-team/CDK/issues/37
	if ok.(bool) || fok.(bool) {

		fmt.Printf(BannerHeader)
		evaluate.CallBasics()

		if Args["--full"].(bool) {
			evaluate.CallAddedFunc()
		}
		return true
	}

	if Args["run"].(bool) {
		if Args["--list"].(bool) {
			plugin.ListAllExploit()
			os.Exit(0)
		}
		name := Args["<exploit>"].(string)
		if plugin.Exploits[name] == nil {
			fmt.Printf("\nInvalid script name: %s , available scripts:\n", name)
			plugin.ListAllExploit()
			return true
		}
		plugin.RunSingleExploit(name)
		return true
	}

	if Args["<tool>"] != nil {
		args := Args["<args>"].([]string)

		switch Args["<tool>"] {
		case "vi":
			PassInnerArgs()
			vi.RunVendorVi()
		case "kcurl":
			kubectl.KubectlToolApi(args)
		case "ectl":
			etcdctl.EtcdctlToolApi(args)
		case "ucurl":
			dockerd_api.UcurlToolApi(args)
		case "dcurl":
			dockerd_api.DcurlToolApi(args)
		case "ifconfig":
			network.GetLocalAddresses()
		case "ps":
			ps.RunPs()
		case "netstat":
			netstat.RunNetstat()
		case "probe":
			if len(args) != 4 {
				log.Println("Invalid input args.")
				log.Println("usage: cdk probe <ip> <port> <parallels> <timeout-ms>")
				log.Fatal("example: cdk probe 192.168.1.0-255 22,80,100-110 50 1000")
			}
			parallel, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				log.Println("err found when parse input arg <parallel>")
				log.Fatal(err)
			}
			timeout, err := strconv.Atoi(args[3])
			if err != nil {
				log.Println("err found when parse input arg <timeout-ms>")
				log.Fatal(err)
			}
			probe.TCPScanToolAPI(args[0], args[1], parallel, timeout)
		default:
			docopt.PrintHelpAndExit(nil, BannerContainer)
		}
	}

	return false
}
