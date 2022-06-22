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

package etcdctl

import (
	"fmt"
	"net/url"
	"strings"
)

var ectlBanner = `ectl - Unauthorized enumeration of ectd keys.

Usage:
  ./cdk ectl <endpoint> get <key>

Example: 
  ./cdk ectl http://172.16.5.4:2379 get /
`

func EtcdctlToolApi(args []string) {
	var opt = EtcdRequestOption{}
	// err break
	if len(args) != 3 {
		fmt.Println(ectlBanner)
		return
	}
	u, err := url.Parse(args[0])
	if err != nil {
		fmt.Println(ectlBanner)
		return
	}
	opt.Endpoint = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	opt.Api = u.Path

	switch strings.ToUpper(args[1]) {
	case "GET":
		opt.Api = "/v3/kv/range"
		opt.PostData = GenerateQuery(args[2])
		//opt.TlsConfig = &tls.Config{}
		DoRequest(opt)
	default: // err break
		fmt.Println(ectlBanner)
		return
	}
}
