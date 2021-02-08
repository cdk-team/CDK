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
