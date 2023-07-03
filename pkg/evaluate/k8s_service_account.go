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
	"encoding/json"
	"fmt"
	"github.com/cdk-team/CDK/pkg/tool/kubectl"
	"log"
	"regexp"
	"strings"
)

type SelfSubjectRulesReview struct {
	Status struct {
		ResourceRules []struct {
			Verbs     []string `json:"verbs"`
			Resources []string `json:"resources"`
		} `json:"resourceRules"`
	} `json:"status"`
}

var subjectRules = `
{
    "apiVersion": "authorization.k8s.io/v1",
    "kind": "SelfSubjectRulesReview",
    "spec": {
      "namespace": "$ns_current"
    }
}
`

var reviewResp SelfSubjectRulesReview

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
		if err != nil && !strings.Contains(resp, "403") {
			fmt.Println(err)
			return false
		}
		if len(resp) > 0 && strings.Contains(resp, "kube-system") {
			fmt.Println("\tsuccess, the service-account have a high authority.")
			fmt.Println("\tnow you can make your own request to takeover the entire k8s cluster with `./cdk kcurl` command\n\tgood luck and have fun.")
			return true
		}
		// get namespace name from resp
		re := regexp.MustCompile(`system:serviceaccount:(\S+):`)
		matches := re.FindStringSubmatch(resp)
		if len(matches) > 1 {
			log.Println("current namespace is:", matches[1])
			CheckSaPermissionForNs(matches[1])
			return true
		} else {
			//fail totally
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

// if attacker can get the namespace name, then check the sa permissions in that namespace
func CheckSaPermissionForNs(ns string) bool {
	log.Println("trying to list permissions of in namespace:", ns)
	subjectRules = strings.ReplaceAll(subjectRules, "$ns_current", ns)
	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: "",
			Server:    "",
			Api:       "/apis/authorization.k8s.io/v1/selfsubjectrulesreviews",
			Method:    "POSt",
			PostData:  subjectRules,
			Anonymous: false,
		})
	if err != nil {
		fmt.Println(err)
		return false
	}
	// Use json.Unmarshal to parse the JSON data
	if err := json.Unmarshal([]byte(resp), &reviewResp); err != nil {
		fmt.Println("Error:", err)
	}

	// Loop over resourceRules to get the verbs and resources
	for i, rule := range reviewResp.Status.ResourceRules {
		log.Printf("Namespace %s Resources[%d]: %s, And Permissons: %s", ns, i, strings.Join(rule.Resources, ","), rule.Verbs)
	}
	return true
}
