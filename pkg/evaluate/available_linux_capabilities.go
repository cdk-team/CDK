package evaluate

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/cdk-team/CDK/pkg/util"
	"github.com/cdk-team/CDK/pkg/util/capability"
)

func GetProcCapabilities() bool {
	data, err := ioutil.ReadFile("/proc/self/status")
	if err != nil {
		log.Println(err)
		return false
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	log.Println("Capabilities hex of Caps(CapInh|CapPrm|CapEff|CapBnd|CapAmb):")

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Cap") {
			fmt.Printf("\t%s\n", line)
		}
	}

	pattern := regexp.MustCompile(`(?i)capeff:\s*?[a-z0-9]+\s`)
	params := pattern.FindStringSubmatch(string(data))

	for _, matched := range params {

		// make capabilities readable
		lst := strings.Split(matched, ":")
		if len(lst) == 2 {

			capStr := strings.TrimSpace(lst[1])
			caps, err := capability.CapHexParser(capStr)

			fmt.Printf("\tCap decode: 0x%s = %s\n", capStr, capability.CapListToString(caps))
			fmt.Printf("\tAdd capability list: %s\n", capability.CapListToString(getAddCaps(caps)))

			if err != nil {
				log.Printf("[-] capability.CapHexParser: %v\n", err)
				return false
			}

			fmt.Printf("[*] Maybe you can exploit the Capabilities below:\n")
			for _, c := range caps {
				switch c {
				case "CAP_DAC_READ_SEARCH":
					fmt.Println("[!] CAP_DAC_READ_SEARCH enabled. You can read files from host. Use 'cdk run cap-dac-read-search' ... for exploitation.")
				case "CAP_SYS_MODULE":
					fmt.Println("[!] CAP_SYS_MODULE enabled. You can escape the container via loading kernel module. More info at https://xcellerator.github.io/posts/docker_escape/.")
				case "CAP_SYS_ADMIN":
					fmt.Println("Critical - SYS_ADMIN Capability Found. Try 'cdk run rewrite-cgroup-devices/mount-cgroup/...'.")
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

func getAddCaps(currentCaps []string) []string {
	var addCaps []string
	for _, c := range currentCaps {
		if !util.StringContains(capability.DockerDefaultCaps, c) {
			addCaps = append(addCaps, c)
		}
	}
	return addCaps
}
