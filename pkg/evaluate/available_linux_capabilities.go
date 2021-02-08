package evaluate

import (
	"fmt"
	"github.com/cdk-team/CDK/pkg/util/capability"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func GetProcCapabilities() bool {
	data, err := ioutil.ReadFile("/proc/self/status")
	if err != nil {
		log.Println(err)
		return false
	}

	pattern := regexp.MustCompile("(?i)capeff:\\s*?[a-z0-9]+\\s")
	params := pattern.FindStringSubmatch(string(data))
	for _, matched := range params {
		log.Println("Capabilities:")
		fmt.Printf("\t%s", matched)

		// make capabilities readable
		lst := strings.Split(matched, ":")
		if len(lst) == 2 {
			capStr := strings.TrimSpace(lst[1])
			fmt.Printf("\tCap decode: 0x%s = %s\n", capStr, capability.CapHexToText(capStr))
		}

		if strings.Contains(matched, "3fffffffff") {
			fmt.Println("Critical - Possible Privileged Container Found.")
			return true
		}
	}
	return false
}
