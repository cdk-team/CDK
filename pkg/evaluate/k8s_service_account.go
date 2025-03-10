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
	"bufio"
	"fmt"
	"github.com/cdk-team/CDK/pkg/tool/kubectl"
	"log"
	"os"
	"strings"
)

var (
	k8sAccountInfoPath string
	kubernetesAddress  string
)

func CheckPrivilegedK8sServiceAccount(tokenPath string) bool {

	// fetch mount info
	file, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		fmt.Printf("error opening /proc/self/mountinfo: %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "serviceaccount") {
			fmt.Println("find serviceaccount successfully")
			k8sAccountInfoPath = line
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading /proc/self/mountinfo: %v\n", err)
	}

	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: k8sAccountInfoPath,
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
				TokenPath: k8sAccountInfoPath,
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
