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

package util

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

)

// from https://stackoverflow.com/questions/40682760/what-syscall-method-could-i-use-to-get-the-default-network-gateway
const (
    file  = "/proc/net/route"
    line  = 1    // line containing the gateway addr. (first line: 0)
    sep   = "\t" // field separator
    field = 2    // field containing hex gateway address (first field: 0)
)

// GetGateway returns the default gateway for the system.
func GetGateway() (string, error) {

    file, err := os.Open(file)
    if err != nil {
        return "", err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {

        // jump to line containing the agteway address
        for i := 0; i < line; i++ {
            scanner.Scan()
        }

        // get field containing gateway address
        tokens := strings.Split(scanner.Text(), sep)
        gatewayHex := "0x" + tokens[field]

        // cast hex address to uint32
        d, _ := strconv.ParseInt(gatewayHex, 0, 64)
        d32 := uint32(d)

        // make net.IP address from uint32
        ipd32 := make(net.IP, 4)
        binary.LittleEndian.PutUint32(ipd32, d32)

        // format net.IP to dotted ipV4 string
        ip := net.IP(ipd32).String()
        
        return ip, nil
    }

	return "", fmt.Errorf("no default gateway found")
}


