//go:build !no_probe_tool
// +build !no_probe_tool

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

package probe

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/cdk-team/CDK/conf"
	"golang.org/x/sync/semaphore"
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

func (ps *PortScanner) Start() {
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
		for _, p := range ps.portRange {
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
	timeout := time.Duration(conf.TCPScannerConf.Timeout) * time.Millisecond

	TCPPScan(ipRange, portFromTo, conf.TCPScannerConf.MaxParallel, timeout)
}

func TCPScanToolAPI(ipRange string, portRange string, parallel int64, timeoutMS int) {
	portFromTo, _ := GetTaskPortListByString(portRange)
	timeout := time.Duration(timeoutMS) * time.Millisecond

	TCPPScan(ipRange, portFromTo, parallel, timeout)
}

func TCPPScan(ipRange string, portRange []FromTo, parallel int64, timeout time.Duration) {

	ps := &PortScanner{
		ipRange:   ipRange,
		portRange: portRange,
		lock:      semaphore.NewWeighted(parallel),
		timeout:   timeout,
	}

	startTime := time.Now()
	log.Printf("scanning %v with user-defined ports, max parallels:%v, timeout:%v\n", ps.ipRange, parallel, ps.timeout)
	ps.Start()

	endTime := time.Now()
	useTime := int64(endTime.Sub(startTime).Seconds() * 1000)
	log.Printf("scanning use time:%vms\n", useTime)
	log.Printf("ending; @args is ips: %v, max parallels:%v, timeout:%v\n", ps.ipRange, conf.TCPScannerConf.MaxParallel, ps.timeout)

}
