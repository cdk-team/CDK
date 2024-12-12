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
	"bytes"
	"fmt"
	"github.com/cdk-team/CDK/pkg/tool/kubectl"
	"log"
	"os/exec"
	"strings"
)

func CheckPrivilegedK8sServiceAccount(tokenPath string, address string) bool {
	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: tokenPath + "/token",
			Server:    address,
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

func GetDefaultK8SAccountInfo() string {
	// 执行 df -T 命令来看 serviceaccount 保存路径
	cmd := exec.Command("df", "-T")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Command execution failed: %s\n", err)
		return ""
	}

	output := out.String()
	lines := strings.Split(output, "\n")

	var serviceAccountLines []string

	for _, line := range lines {
		if strings.Contains(line, "serviceaccount") {
			fmt.Println("\tk8s account service path fetch success" + line)
			serviceAccountLines = append(serviceAccountLines, line)
		}
	}

	return strings.Join(serviceAccountLines, "\n")
}

func GetKubernetesAddress() (string, error) {
	cmd := exec.Command("env")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %s", err)
	}

	output := out.String()
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "KUBERNETES_PORT_443_TCP_ADDR=") {
			return strings.TrimPrefix(line, "KUBERNETES_PORT_443_TCP_ADDR="), nil
		}
	}

	return "", fmt.Errorf("KUBERNETES_PORT_443_TCP_ADDR not found")
}
