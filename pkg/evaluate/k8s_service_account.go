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

	"github.com/cdk-team/CDK/conf"
	"github.com/cdk-team/CDK/pkg/tool/kubectl"
)

func CheckPrivilegedK8sServiceAccount(tokenPath string) bool {
	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: "",
			Server:    "",
			Api:       "/apis",
			Method:    "get",
			PostData:  "",
			Anonymous: false,
		})
	if err != nil {
		fmt.Println(err)
		return false
	}
	if len(resp) > 0 && strings.Contains(resp, "APIGroupList") {
		fmt.Println("\tservice-account is available")

		// check if the current service-account can list namespaces
		log.Println("trying to list namespaces")
		resp, err := kubectl.ServerAccountRequest(
			kubectl.K8sRequestOption{
				TokenPath: "",
				Server:    "",
				Api:       "/api/v1/namespaces",
				Method:    "get",
				PostData:  "",
				Anonymous: false,
			})
		if err != nil {
			fmt.Println(err)
			return false
		}
		if len(resp) > 0 && strings.Contains(resp, "kube-system") {
			fmt.Println("\tsuccess, the service-account have a high authority.")
			fmt.Println("\tnow you can make your own request to takeover the entire k8s cluster with `./cdk kcurl` command\n\tgood luck and have fun.")
			return true
		} else {
			fmt.Println("\tfailed")
			fmt.Println("\tresponse:" + resp)
			return false
		}
	} else {
		fmt.Println("\tservice-account is not available")
		fmt.Println("\tresponse:" + resp)
		return false
	}
}

func init() {
	RegisterSimpleCheck(
		CategoryK8sServiceAccount,
		"k8s.privileged_service_account",
		"Check Kubernetes service account privileges",
		func() {
			CheckPrivilegedK8sServiceAccount(conf.K8sSATokenDefaultPath)
		},
	)
}
