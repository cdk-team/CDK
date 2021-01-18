package evaluate

import (
	"io/ioutil"
	"log"
)

var RouteLocalNetProcPath = "/proc/sys/net/ipv4/conf/all/route_localnet"

func CheckRouteLocalNetworkValue() {
	data, err := ioutil.ReadFile(RouteLocalNetProcPath)
	if err != nil {
		log.Printf("err found while open %s: %v\n", RouteLocalNetProcPath, err)
		return
	}
	log.Printf("net.ipv4.conf.all.route_localnet = %s\n", string(data))
	if string(data) == "1" {
		log.Println("You may be able to access the localhost service of the current container node or other nodes.")
		// from: https://github.com/kubernetes/kubernetes/issues/92315
		log.Println("CVE-2020-8558: The Kubelet and kube-proxy components in versions 1.1.0-1.16.10, 1.17.0-1.17.6, and 1.18.0-1.18.3 were found to contain a security issue which allows adjacent hosts to reach TCP and UDP services bound to 127.0.0.1 running on the node or in the node's network namespace. Node setting allows for neighboring hosts to bypass localhost boundary.")
	}
}
