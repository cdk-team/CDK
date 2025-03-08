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

package kubectl

import (
	"fmt"
	"github.com/cdk-team/CDK/conf"
	"log"
	"strings"
)

var kcurlBanner = `kcurl - send HTTP request to K8s api-server.

Usage:
  ./cdk kcurl (<token_path>|anonymous|default) (get|post) <url> [<post_data>]

Options:
  token_path  connect api-server with user-specified service-account token.
  anonymous   connect api-server using system:anonymous service-account.
  default     connect api-server using pod default service-account token.

Example: 
  ./cdk kcurl default get 'https://192.168.0.234:6443/api/v1/nodes'
  ./cdk kcurl /var/run/secrets/kubernetes.io/serviceaccount/token get 'https://192.168.0.234:6443/api/v1/nodes'
  ./cdk kcurl anonymous post 'https://192.168.0.234:6443/api/v1/nodes' '{"apiVersion":"v1",...}'
`

func KubectlToolApi(args []string) {

	var opts = K8sRequestOption{}

	// err break
	if len(args) != 3 && len(args) != 4 {
		fmt.Println(kcurlBanner)
		return
	}

	switch args[0] {
	case "default":
		opts.TokenPath = conf.K8sSATokenDefaultPath
	case "anonymous":
		opts.TokenPath = ""
		opts.Anonymous = true
	default:
		opts.TokenPath = args[0]
	}

	switch strings.ToUpper(args[1]) {
	case "POST":
		opts.Method = "POST"
	case "GET":
		opts.Method = "GET"
	default: // err break
		fmt.Println(kcurlBanner)
		return
	}

	if len(args) == 3 {
		opts.Url = args[2]
	} else {
		opts.Url = args[2]
		opts.PostData = args[3]
	}

	resp, err := ServerAccountRequest(opts)
	if err != nil {
		log.Println("failed to get api-server response")
		fmt.Println(err)
	} else {
		log.Println("api-server response:")
		fmt.Println(resp)
	}
}
