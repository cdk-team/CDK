
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

package dockerd_api

import (
	"fmt"
	"github.com/cdk-team/CDK/pkg/util"
	"log"
)

func UcurlToolApi(args []string) {
	if len(args) != 4 {
		log.Fatal("invalid input args, Example: ./cdk ucurl get /var/run/docker.sock http://127.0.0.1/info \"\"")
	}
	ans, err := util.UnixHttpSend(args[0], args[1], args[2], args[3])
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("response:")
	fmt.Println(ans)
}

func DcurlToolApi(args []string) {
	if len(args) != 3 {
		log.Fatal("invalid input args, Example: ./cdk dcurl get http://127.0.0.1:2375/info \"\"")
	}
	ans, err := util.HttpSendJson(args[0], args[1], args[2])
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("response:")
	fmt.Println(ans)
}
