package evaluate

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

// https://github.com/kubernetes/dns/blob/master/docs/specification.md
func DNSBasedServiceDiscovery() {
	dnsNames := []string{"any.any.svc.cluster.local.", "any.any.any.svc.cluster.local."}

	var results []*net.SRV
	for _, name := range dnsNames {
		_, srvs, err := net.LookupSRV("", "", name)
		if err != nil {
			fmt.Printf("error when requesting coreDNS: %s\n", err.Error())
			continue
		}

		results = append(results, srvs...)
	}

	sort.Slice(results, func(i, j int) bool {
		switch strings.Compare(results[i].Target, results[j].Target) {
		case -1:
			return true
		case 0:
			return results[i].Port < results[j].Port
		case 1:
			return false
		}
		return false
	})

	for _, srv := range results {
		fmt.Println(srv.Target, srv.Port)
	}
}
