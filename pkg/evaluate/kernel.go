package evaluate

import (
	"os/exec"
	"strings"

	"github.com/cdk-team/CDK/conf"
	"github.com/cdk-team/CDK/pkg/util"
)

// kernelExploitSuggester
// use https://github.com/mzet-/linux-exploit-suggester to check kernel exploit
// run linux-exploit-suggester bash script to check kernel exploit
func kernelExploitSuggester() {
	script := conf.KernelExploitScript
	// check bash command available
	_, err := exec.LookPath("bash")
	if err != nil {
		return
	}
	
	// run command get output
	util.PrintItemValueWithKeyOneLine("refer", "https://github.com/mzet-/linux-exploit-suggester", false)
	output, err := exec.Command("bash", "-c", script).Output()
	if err != nil {
		return
	}

	// get all available exploit
	// sed "s,$(printf '\033')\\[[0-9;]*[a-zA-Z],,g" | grep -i "\[CVE" -A 10 | grep -Ev "^\-\-$" | sed -${E} "s,\[CVE-[0-9]+-[0-9]+\].*,${SED_RED},g"
	// ANSI escape code in output, reg can not match it
	indexs := make([]int, 0)
	lines := strings.Split(string(output), "\n")
	for index, line := range lines {
		if strings.Contains(line, "[CVE") {
			indexs = append(indexs, index)
		}
	}

	// print all exploit and after 10 lines
	for _, index := range indexs {
		for i := index; i < index+10; i++ {
			if i >= len(lines) {
				break
			}

			// do not print CVE number twice
			if i != index && strings.Contains(lines[i], "[CVE") {
				break
			}

			util.PrintOrignal(lines[i])
		}
	}
	
}
