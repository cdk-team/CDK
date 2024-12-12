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
	"github.com/cdk-team/CDK/pkg/util"
)

// CallBasics is a function to call basic functions
func CallBasics() {
	util.PrintH2("Information Gathering - System Info")
	BasicSysInfo()
	FindSidFiles()

	util.PrintH2("Information Gathering - Services")
	SearchSensitiveEnv()
	SearchSensitiveService()

	util.PrintH2("Information Gathering - Commands and Capabilities")
	SearchAvailableCommands()
	GetProcCapabilities()

	util.PrintH2("Information Gathering - Mounts")
	MountEscape()

	util.PrintH2("Information Gathering - Net Namespace")
	CheckNetNamespace()

	util.PrintH2("Information Gathering - Sysctl Variables")
	CheckRouteLocalNetworkValue()

	util.PrintH2("Information Gathering - DNS-Based Service Discovery")
	DNSBasedServiceDiscovery()

	util.PrintH2("Discovery - K8s API Server")
	CheckK8sAnonymousLogin()

	util.PrintH2("Discovery - K8s Service Account")

	path := GetDefaultK8SAccountInfo()
	address, err := GetKubernetesAddress()

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("KUBERNETES_PORT_443_TCP_ADDR:", address)
	}

	CheckPrivilegedK8sServiceAccount(path, address)

	util.PrintH2("Discovery - Cloud Provider Metadata API")
	CheckCloudMetadataAPI()

	util.PrintH2("Exploit Pre - Kernel Exploits")
	kernelExploitSuggester()
}

func CallAddedFunc() {
	util.PrintH2("Information Gathering - Sensitive Files")
	SearchLocalFilePath()

	util.PrintH2("Information Gathering - ASLR")
	ASLR()

	util.PrintH2("Information Gathering - Cgroups")
	DumpCgroup()
}
