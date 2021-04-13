// +build thin no_netcat_tool

package netcat

import "github.com/cdk-team/CDK/conf"

func RunVendorNetcat() {
	print(conf.ThinIgnoreTool)
	return
}