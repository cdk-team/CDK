package evaluate

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/cdk-team/CDK/pkg/util/capability"
)

func GetProcCapabilities() bool {
	data, err := ioutil.ReadFile("/proc/self/status")
	if err != nil {
		log.Println(err)
		return false
	}

	pattern := regexp.MustCompile(`(?i)capeff:\s*?[a-z0-9]+\s`)
	params := pattern.FindStringSubmatch(string(data))

	for _, matched := range params {
		log.Println("Capabilities:")
		fmt.Printf("\t%s", matched)

		// make capabilities readable
		lst := strings.Split(matched, ":")
		if len(lst) == 2 {
			capStr := strings.TrimSpace(lst[1])
			caps, err := capability.CapHexParser(capStr)
			fmt.Printf("\tCap decode: 0x%s = %s\n", capStr, capability.CapListToString(caps))

			if err != nil {
				log.Printf("[-] capability.CapHexParser: %v\n", err)
				return false
			}
			for _, c := range caps {
				switch c {
				case "CAP_DAC_READ_SEARCH":
					fmt.Println("[!] CAP_DAC_READ_SEARCH enabled. You can read files from host. Use 'cdk run cap-dac-read-search' ... for exploitation.")
				case "CAP_SYS_MODULE":
					fmt.Println("[!] CAP_SYS_MODULE enabled. You can escape the container via loading kernel module. More info at https://xcellerator.github.io/posts/docker_escape/.")
				case "CAP_SYS_ADMIN":
					fmt.Println("Critical - SYS_ADMIN Capability Found.")
				}
			}
		}

		if strings.Contains(matched, "3fffffffff") {
			fmt.Println("Critical - Possible Privileged Container Found.")
			return true
		}
	}
	return false
}
