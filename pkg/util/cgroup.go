
/*
Copyright 2022 The Authors of https://github.com/CDK-TEAM/CDK .

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

const mountInfoPath string = "/proc/self/mountinfo"
const hostDeviceFlag string = "/etc/hosts"
const cgroupInfoPath string = "/proc/self/cgroup"

type MountInfo struct {
	Device            string
	Fstype            string
	Root              string
	MountPoint        string
	Opts              []string
	Major             string
	Minor             string
	SuperBlockOptions []string
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

		// former Part
		// https://man7.org/linux/man-pages/man5/proc.5.html
		fields := strings.Fields(parts[0])
		// mountID = fields[0] ; parentID = fields[1]
		blockId := strings.Split(fields[2], ":")
		if len(blockId) != 2 {
			return nil, fmt.Errorf("found invalid mountinfo line in file %s: %s ", mountInfoPath, r)
		}
		mi.Major = blockId[0]
		mi.Minor = blockId[1]
		mi.Root = fields[3]
		mi.MountPoint = fields[4]
		mi.Opts = strings.Split(fields[5], ",")

		// latter part
		fields = strings.Fields(parts[1])
		mi.Fstype = fields[0]
		mi.Device = fields[1]
		mi.SuperBlockOptions = strings.Split(fields[2], ",")

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

// get kernel version
func GetKernelVersion() ([]int, error) {
	utsInfo := &unix.Utsname{}
	err := unix.Uname(utsInfo)
	if err != nil {
		return nil, err
	}
	relStr := string(utsInfo.Release[:])
	relIdx := strings.Index(relStr, "-")
	if relIdx == -1 {
		return nil, errors.New("unknown internal error when executing uname")
	}
	ret := make([]int, 3)
	for _, v := range strings.Split(relStr[:relIdx], ".") {
		verData, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, verData)
	}
	return ret, nil
}

// get cgroup version V1/V2
// hybrid mode will not work in container
func GetCgroupVersion() (int, error) {
	// detect by /sys/fs/cgroup/cgroup.controllers
	// others:
	// or /proc/filesystems
	// or directly try to mount cgroup2 with none
	_, err := os.Stat("/sys/fs/cgroup/cgroup.controllers")
	if err == nil {
		return 2, nil
	}
	if strings.Contains(err.Error(), "no such file or directory") {
		return 1, nil
	}
	return -1, err
}

type CgroupInfo struct {
	HierarchyID   int
	ControllerLst string // split by "," but should not be split
	CgroupPath    string
	OriginalInfo  string
}

func GetAllCGroup() ([]CgroupInfo, error) {
	return GetCgroup(0)
}

// GetCgroup returns the cgroup info of the process
// param pid: 0 = self, 1 = container main process
func GetCgroup(pid int) ([]CgroupInfo, error) {
	var cginfo []CgroupInfo
	var pidStr string

	if pid == 0 {
		pidStr = "self"
	} else {
		pidStr = fmt.Sprint(pid)
	}

	cgroupInfoPath := fmt.Sprintf("/proc/%s/cgroup", pidStr)
	datafd, err := os.Open(cgroupInfoPath)
	if err != nil {
		return nil, err
	}
	defer datafd.Close()

	sc := bufio.NewScanner(datafd)
	for sc.Scan() {
		// Sample "9:devices:/docker/fc1413683c2976fa292c0b1e011224706c1ecc151bad9ceabc9cfcb8dce4ddbb"
		originalInfo := sc.Text()
		singleCG := strings.Split(strings.TrimSuffix(originalInfo, "\n"), ":")
		hID, err := strconv.Atoi(singleCG[0])
		if err != nil {
			return nil, err
		}
		cginfo = append(cginfo, CgroupInfo{hID, singleCG[1], singleCG[2], originalInfo})
	}

	return cginfo, nil
}

func GetAllCGroupSubSystem() ([]string, error) {
	cgSyses, err := GetAllCGroup()
	if err != nil {
		return nil, err
	}
	var syses []string
	for _, v := range cgSyses {
		syses = append(syses, v.ControllerLst)
	}
	return syses, nil
}
