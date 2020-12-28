package evaluate

import (
	"github.com/shirou/gopsutil/v3/host"
	"log"
	"os"
	"os/user"
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

	// os/kernel version
	kversion, _ := host.KernelVersion()
	platform, family, osversion, _ := host.PlatformInformation()
	log.Println(family, platform, osversion, "kernel:", kversion)

}
