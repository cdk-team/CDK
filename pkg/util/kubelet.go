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
	"fmt"
	"golang.org/x/net/route"
)


// GetGateway returns the default gateway for the system.
// from https://gist.github.com/abimaelmartell/dcbbff464dc0778165b2dcc5092f90e6
func GetGateway() (string, error) {
	var defaultRoute = [4]byte{0, 0, 0, 0}

	rib, _ := route.FetchRIB(0, route.RIBTypeRoute, 0)
	messages, err := route.ParseRIB(route.RIBTypeRoute, rib)

	if err != nil {
		return "", err
	}

	for _, message := range messages {
		route_message := message.(*route.RouteMessage)
		addresses := route_message.Addrs

		var destination, gateway *route.Inet4Addr
		ok := false

		if destination, ok = addresses[0].(*route.Inet4Addr); !ok {
			continue
		}

		if gateway, ok = addresses[1].(*route.Inet4Addr); !ok {
			continue
		}

		if destination == nil || gateway == nil {
			continue
		}

		if destination.IP == defaultRoute {
			gateway :=  gateway.IP
			str := fmt.Sprintf("%v.%v.%v.%v", gateway[0], gateway[1], gateway[2], gateway[3])
			return str, nil
		}
	}

	return "", fmt.Errorf("no default gateway found")
}


