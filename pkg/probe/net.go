package probe

import (
	"context"
	"fmt"
	"github.com/Xyntax/CDK/conf"
	"golang.org/x/sync/semaphore"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type PortScanner struct {
	ipRange   string
	portRange []FromTo
	lock      *semaphore.Weighted
	timeout   time.Duration
}

func ScanPort(ip string, port int, timeout time.Duration) bool {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(ip, port, timeout)
		}
		return false
	}

	_ = conn.Close()
	return true
}

func (ps *PortScanner) Start(portFromTo []FromTo) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	base, start, end, err := GetTaskIPList(ps.ipRange)
	if err != nil {
		log.Println("error found when gene ip list to scan task")
		log.Fatal(err)
	}

	// iterate ip in task list
	for ipExt := start; ipExt <= end; ipExt++ {
		ip := base + "." + fmt.Sprintf("%d", ipExt)
		// iterate port in task list
		for _, p := range portFromTo {
			// iterate port from A-B
			for port := p.From; port <= p.To; port++ {
				// lock down the context
				ps.lock.Acquire(context.TODO(), 1)
				wg.Add(1)
				go func(port int, p FromTo) {
					defer ps.lock.Release(1)
					defer wg.Done()
					if ScanPort(ip, port, ps.timeout) {
						fmt.Printf("open %s: %s:%d\n", p.Desc, ip, port)
					}
				}(port, p) // send all sync objects into args
			}
		}
	}
}

func TCPScanExploitAPI(ipRange string) {
	portFromTo, _ := GetTaskPortList()
	ps := &PortScanner{
		ipRange:   ipRange,
		portRange: portFromTo,
		lock:      semaphore.NewWeighted(conf.TCPScannerConf.MaxParallel),
		timeout:   conf.TCPScannerConf.Timeout,
	}
	log.Printf("scanning %v with pre-defined ports, max parallels:%v, timeout:%v\n", ps.ipRange, conf.TCPScannerConf.MaxParallel, ps.timeout)
	ps.Start(portFromTo)
}

func TCPScanToolAPI(ipRange string, portRange string, parallel int64, timeoutMS int) {
	portFromTo, _ := GetTaskPortListByString(portRange)

	ps := &PortScanner{
		ipRange:   ipRange,
		portRange: portFromTo,
		lock:      semaphore.NewWeighted(parallel),
		timeout:   time.Duration(timeoutMS) * time.Millisecond,
	}
	log.Printf("scanning %v with user-defined ports, max parallels:%v, timeout:%v\n", ps.ipRange, parallel, ps.timeout)
	ps.Start(portFromTo)
}
