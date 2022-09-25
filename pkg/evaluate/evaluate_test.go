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
	"testing"

	"github.com/cdk-team/CDK/pkg/util"
)

func TestDumpCgroup(t *testing.T) {
	fmt.Printf("\n[Information Gathering - Cgroups]\n")
	DumpCgroup()
}

func TestFindSidFiles(t *testing.T) {
	fmt.Printf("\n[Information Gathering - SIDs]\n")
	FindSidFiles()
}

func TestKernelExploitSuggester(t *testing.T) {
	util.PrintH2("Exploit Pre - Kernel Exploits")
	kernelExploitSuggester()
}
