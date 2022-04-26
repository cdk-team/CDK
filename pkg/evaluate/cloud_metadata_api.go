
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
	"github.com/cdk-team/CDK/conf"
	"github.com/idoubi/goz"
	"log"
	"strings"
)

func CheckCloudMetadataAPI() {
	for _, apiInstance := range conf.CloudAPI {
		cli := goz.NewClient(goz.Options{
			Timeout: 1,
		})
		resp, err := cli.Get(apiInstance.API)
		if err != nil {
			log.Printf("failed to dial %s API.", apiInstance.CloudProvider)
			continue
		}
		r, _ := resp.GetBody()
		if strings.Contains(r.String(), apiInstance.ResponseMatch) {
			fmt.Printf("\t%s Metadata API available in %s\n", apiInstance.CloudProvider, apiInstance.API)
			fmt.Printf("\tDocs: %s\n", apiInstance.DocURL)
		} else {
			log.Printf("failed to dial %s API.", apiInstance.CloudProvider)
		}
	}
}
