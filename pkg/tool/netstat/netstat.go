package netstat

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
	"log"
	"sort"
)

func RunNetstat() {
	log.Printf("[+] run netstat, using RunNestat()")
	stats, err := net.Connections("all")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("ipType\t\tconnection\tlocalAddr\t\t\tstatus\t\t\tremoteAddr\t\t\tpid\n")
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Type < stats[j].Type
	})
	for _, stat := range stats {
		switch stat.Family {
		case 2:
			switch stat.Type {
			case 1:
				fmt.Printf("ipv4\t\ttcp\t\t%-16s\t\t%-13s\t\t%-16s\t\t%d\n", fmt.Sprintf("%s:%d", stat.Laddr.IP, stat.Laddr.Port), stat.Status, fmt.Sprintf("%s:%d", stat.Raddr.IP, stat.Raddr.Port), stat.Pid)
			case 2:
				fmt.Printf("ipv4\t\tudp\t\t%-16s\t\t%-13s\t\t%-16s\t\t%d\n", fmt.Sprintf("%s:%d", stat.Laddr.IP, stat.Laddr.Port), stat.Status, fmt.Sprintf("%s:%d", stat.Raddr.IP, stat.Raddr.Port), stat.Pid)
			}
		case 23:
			switch stat.Type {
			case 1:
				fmt.Printf("ipv6\t\ttcp\t\t%-16s\t\t%-13s\t\t%-16s\t\t%d\n", fmt.Sprintf("%s:%d", stat.Laddr.IP, stat.Laddr.Port), stat.Status, fmt.Sprintf("%s:%d", stat.Raddr.IP, stat.Raddr.Port), stat.Pid)
			case 2:
				fmt.Printf("ipv6\t\tudp\t\t%-16s\t\t%-13s\t\t%-16s\t\t%d\n", fmt.Sprintf("%s:%d", stat.Laddr.IP, stat.Laddr.Port), stat.Status, fmt.Sprintf("%s:%d", stat.Raddr.IP, stat.Raddr.Port), stat.Pid)
			}

		}
	}
}
