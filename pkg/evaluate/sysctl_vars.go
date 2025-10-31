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

package evaluate

import (
	"io/ioutil"
	"log"
	"strings"
)

var RouteLocalNetProcPath = "/proc/sys/net/ipv4/conf/all/route_localnet"

func CheckRouteLocalNetworkValue() {
	data, err := ioutil.ReadFile(RouteLocalNetProcPath)
	if err != nil {
		log.Printf("err found while open %s: %v\n", RouteLocalNetProcPath, err)
		return
	}
	log.Printf("net.ipv4.conf.all.route_localnet = %s", string(data))
	if strings.TrimSpace(string(data)) == "1" {
		// from: https://github.com/kubernetes/kubernetes/issues/92315
		log.Println("You may be able to access the localhost service of the current container node or other nodes.")
	}
}

func init() {
	RegisterSimpleCheck(CategorySysctl, "sysctl.route_localnet", "Check route_localnet sysctl value", CheckRouteLocalNetworkValue)
}
