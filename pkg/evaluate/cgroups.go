package evaluate

import (
	"fmt"
	"io/ioutil"
	"log"
	"bufio"
	"strings"
)

var MainPIDCgroup = "/proc/1/cgroup"

func DumpMainCgroup() {

	data, err := ioutil.ReadFile(MainPIDCgroup)
	if err != nil {
		log.Printf("err found while open %s: %v\n", MainPIDCgroup, err)
		return
	}
	log.Println("/proc/1/cgroup file content:")

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		fmt.Printf("\t%s\n", scanner.Text())
	}

}
