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
