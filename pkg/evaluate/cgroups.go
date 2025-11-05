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

	"github.com/cdk-team/CDK/pkg/util"
)

func DumpCgroup() {

	cgroupLst, err := util.GetCgroup(1)
	sLst := make([]string, 0, len(cgroupLst))

	if err != nil {
		log.Printf("/proc/1/cgroup error: %v\n", err)
		return
	}

	log.Println("/proc/1/cgroup file content:")
	for _, v := range cgroupLst {
		fmt.Printf("\t%s\n", v.OriginalInfo)
		sLst = append(sLst, v.OriginalInfo)
	}

	cgroupLstSelf, err := util.GetCgroup(0)
	if err != nil {
		log.Printf("/proc/self/cgroup error: %v\n", err)
		return
	}

	log.Println("/proc/self/cgroup file added content (compare pid 1) :")
	for _, v := range cgroupLstSelf {
		if !util.StringContains(sLst, v.OriginalInfo) {
			fmt.Printf("\t%s\n", v.OriginalInfo)
		}
	}

}

func init() {
	RegisterSimpleCheck(CategoryCgroups, "cgroups.dump", "Dump cgroup configuration", DumpCgroup)
}
