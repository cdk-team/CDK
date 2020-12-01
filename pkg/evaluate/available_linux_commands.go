package evaluate

import (
	"github.com/Xyntax/CDK/conf"
	"log"
	"os/exec"
	"strings"
)

func SearchAvailableCommands() {
	ans := []string{}
	for _, cmd := range conf.LinuxCommandChecklist {
		_, err := exec.LookPath(cmd)
		if err == nil {
			ans = append(ans, cmd)
		}
	}
	log.Printf("available commands:\n\t%s\n", strings.Join(ans, ","))
}
