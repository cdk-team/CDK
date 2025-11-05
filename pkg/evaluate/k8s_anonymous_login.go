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

package evaluate

import (
	"fmt"
	"log"
	"strings"

	"github.com/cdk-team/CDK/pkg/tool/kubectl"
)

func CheckK8sAnonymousLogin() bool {

	// check if api-server allows system:anonymous request
	log.Println("checking if api-server allows system:anonymous request.")

	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: "",
			Server:    "", // default
			Api:       "/",
			Method:    "get",
			PostData:  "",
			Anonymous: true,
		})
	if err != nil {
		fmt.Println(err)
	}

	if strings.Contains(resp, "/api") {
		fmt.Println("\tcongrats, api-server allows anonymous request.")
		log.Println("trying to list namespaces")

		// check if system:anonymous can list namespaces
		resp, err := kubectl.ServerAccountRequest(
			kubectl.K8sRequestOption{
				TokenPath: "",
				Server:    "", // default
				Api:       "/api/v1/namespaces",
				Method:    "get",
				PostData:  "",
				Anonymous: true,
			})
		if err != nil {
			fmt.Println(err)
		}
		if len(resp) > 0 && strings.Contains(resp, "kube-system") {
			fmt.Println("\tsuccess, the system:anonymous role have a high authority.")
			fmt.Println("\tnow you can make your own request to takeover the entire k8s cluster with `./cdk kcurl` command\n\tgood luck and have fun.")
			return true
		} else {
			fmt.Println("\tfailed.")
			fmt.Println("\tresponse:" + resp)
			return true
		}
	} else {
		fmt.Println("\tapi-server forbids anonymous request.")
		fmt.Println("\tresponse:" + resp)
		return false
	}
}

func init() {
	RegisterSimpleCheck(
		CategoryK8sAPIServer,
		"k8s.anonymous_login",
		"Attempt anonymous Kubernetes API login",
		func() {
			CheckK8sAnonymousLogin()
		},
	)
}
