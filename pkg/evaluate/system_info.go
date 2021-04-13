package evaluate

import (
	"github.com/shirou/gopsutil/v3/host"
	"log"
	"os"
	"os/user"
	"io/ioutil"
)

func BasicSysInfo() {
	// current dir(pwd)
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("current dir:", dir)

	// current user(id)
	u, err := user.Current()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("current user:", u.Username, "uid:", u.Uid, "gid:", u.Gid, "home:", u.HomeDir)

	// hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("hostname:", hostname)

	// os/kernel version
	kversion, _ := host.KernelVersion()
	platform, family, osversion, _ := host.PlatformInformation()
	log.Println(family, platform, osversion, "kernel:", kversion)

}

func ASLR() {
	// ASLR off: /proc/sys/kernel/randomize_va_space = 0 
	var ASLRSetting = "/proc/sys/kernel/randomize_va_space"

	data, err := ioutil.ReadFile(ASLRSetting)
	if err != nil {
		log.Printf("err found while open %s: %v\n", RouteLocalNetProcPath, err)
		return
	}
	log.Printf("/proc/sys/kernel/randomize_va_space file content: %s", string(data))

	if string(data) == "0" {
		log.Println("ASLR is disabled.")
	} else {
		log.Println("ASLR is enabled.")
	}

}
