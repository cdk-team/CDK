package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const mountInfoPath string = "/proc/self/mountinfo"
const hostDeviceFlag string = "/etc/hosts"

type MountInfo struct {
	Device     string
	Fstype     string
	Root       string
	MountPoint string
	Opts       []string
	Major      string
	Minor      string
}

// find block device id
func FindTargetDeviceID(mi *MountInfo) bool {
	if mi.MountPoint == hostDeviceFlag {
		log.Printf("found host blockDeviceId Major: %s Minor: %s\n", mi.Major, mi.Minor)
		return true
	}
	return false
}

func GetMountInfo() ([]MountInfo, error) {
	f, err := os.Open(mountInfoPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}
	// 2346 2345 0:261 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw
	mountInfos := make([]MountInfo, len(ret))

	for _, r := range ret {
		parts := strings.Split(r, " - ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("found invalid mountinfo line in file %s: %s ", mountInfoPath, r)
		}
		mi := MountInfo{}
		fields := strings.Fields(parts[0])
		mi.Root = fields[3]
		mi.MountPoint = fields[4]
		mi.Opts = strings.Split(fields[5], ",")
		blockId := strings.Split(fields[2], ":")
		if len(blockId) != 2 {
			return nil, fmt.Errorf("found invalid mountinfo line in file %s: %s ", mountInfoPath, r)
		}
		mi.Major = blockId[0]
		mi.Minor = blockId[1]
		fields = strings.Fields(parts[1])
		mi.Device = fields[1]
		mi.Fstype = fields[0]
		mountInfos = append(mountInfos, mi)
	}

	return mountInfos, err
}

func MakeDev(major, minor string) int {
	ret1, err := strconv.ParseInt(major, 10, 64)
	if err != nil {
		log.Printf("convert major number to int64 err: %v\n", err)
		return 0
	}
	ret2, err := strconv.ParseInt(minor, 10, 64)
	if err != nil {
		log.Printf("convert minor number to int64 err: %v\n", err)
		return 0
	}

	return int(((ret1 & 0xfff) << 8) | (ret2 & 0xff) | ((ret1 &^ 0xfff) << 32) | ((ret2 & 0xfffff00) << 12))
}

// set all block device accessible
func SetBlockAccessible(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_SYNC, 0200)
	if err != nil {
		return fmt.Errorf("open devices.allow failed. %v\n", err)
	}
	defer f.Close()

	l, err := f.Write([]byte("a"))
	if err != nil {
		return fmt.Errorf("write devices.allow failed. %v\n", err)
	}

	if l != 1 {
		return fmt.Errorf("write \"a\" to devices.allow failed.\n")
	}
	log.Printf("set all block device accessible success.\n")

	return nil
}
