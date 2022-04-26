
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
	"github.com/cdk-team/CDK/conf"
	gops "github.com/mitchellh/go-ps"
	"log"
	"regexp"
)

func SearchSensitiveService() {
	processList, err := gops.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
	}
	for _, proc := range processList {
		ans, err := regexp.MatchString(conf.SensitiveProcessRegex, proc.Executable())
		if err != nil {
			log.Println(err)
		} else if ans {
			log.Printf("service found in process:\n\t%d\t%d\t%s\n", proc.Pid(), proc.PPid(), proc.Executable())
		}
	}
}
