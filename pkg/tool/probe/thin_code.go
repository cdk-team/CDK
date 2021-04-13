// +build no_probe_tool

package probe

import "github.com/cdk-team/CDK/conf"

func TCPScanToolAPI(ipRange string, portRange string, parallel int64, timeoutMS int) {
	print(conf.ThinIgnoreTool)
	return
}

func TCPScanExploitAPI(ipRange string) {
	print(conf.ThinIgnoreTool)
	return
}
