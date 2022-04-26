
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
	"github.com/cdk-team/CDK/pkg/util"
	"os"
	"path/filepath"
	"strings"
)

func SearchLocalFilePath() {

	filepath.Walk(conf.SensitiveFileConf.StartDir, func(path string, info os.FileInfo, err error) error {
		for _, name := range conf.SensitiveFileConf.NameList {
			currentPath := strings.ToLower(path)
			//if util.IsSoftLink(currentPath) && util.IsDir(currentPath) {
			//	fmt.Println("skip", currentPath)
			//	return filepath.SkipDir // skip soft link or it will run into container runtime filesystem
			//}
			if strings.Contains(currentPath, name) {
				fmt.Printf("\t%s - %s\n", name, path)
				if util.IsDir(currentPath) {
					return filepath.SkipDir // stop dive if sensitive dir found
				}
				return nil
			}
		}
		return nil
	})

}
