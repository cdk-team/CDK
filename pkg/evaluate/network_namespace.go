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
