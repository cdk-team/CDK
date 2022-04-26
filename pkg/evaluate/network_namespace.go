
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
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func CheckNetNamespace() {
	p := "/proc/net/unix"
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Println("err found while open", p)
		return
	}

	if strings.Contains(string(data), "/systemd/") || strings.Contains(string(data), "/run/user/") {
		fmt.Println("\thost unix-socket found, seems container started with --net=host privilege.")
	} else {
		fmt.Println("\tcontainer net namespace isolated.")
	}

	pattern := regexp.MustCompile(`\@/containerd-shim/[\w\d\/-]*?\.sock`)
	params := pattern.FindAllStringSubmatch(string(data), -1)
	for _, matched := range params {
		fmt.Printf("\tfound containerd-shim socket in: %s\n", matched)
	}

}
