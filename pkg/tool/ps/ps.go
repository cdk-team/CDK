package ps

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"log"
)

func RunPs() {
	ps, err := process.Processes()
	if err != nil {
		log.Fatal("get process list failed.")
	}
	for _, p := range ps {
		pexe, _ := p.Exe()
		ppid, _ := p.Ppid()
		user, _ := p.Username()
		fmt.Printf("%v\t%v\t%v\t%v\n", user, p.Pid, ppid, pexe)
	}
}
