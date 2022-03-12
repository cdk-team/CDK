package evaluate

import (
	"fmt"
	"log"

	"github.com/cdk-team/CDK/pkg/util"
)

func DumpCgroup() {

	cgroupLst, err := util.GetCgroup(1)
	sLst := make([]string, len(cgroupLst))

	if err != nil {
		log.Printf("/proc/1/cgroup error: %v\n", err)
		return
	}

	log.Println("/proc/1/cgroup file content:")
	for _, v := range cgroupLst {
		fmt.Printf("\t%s\n", v.OriginalInfo)
		sLst = append(sLst, v.OriginalInfo)
	}

	cgroupLstSelf, err := util.GetCgroup(0)
	if err != nil {
		log.Printf("/proc/self/cgroup error: %v\n", err)
		return
	}

	log.Println("/proc/self/cgroup file added content (compare pid 1) :")
	for _, v := range cgroupLstSelf {
		if !util.StringContains(sLst, v.OriginalInfo) {
			fmt.Printf("\t%s\n", v.OriginalInfo)
		}
	}

}
