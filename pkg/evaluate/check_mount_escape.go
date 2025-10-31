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

package evaluate

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/cdk-team/CDK/pkg/errors"
	"github.com/cdk-team/CDK/pkg/util"
)

// The checkClose function calls close on a Closer and panics with a
// runtime error if the Closer returns an error
func checkClose(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(&errors.CDKRuntimeError{Err: err})
	}
}

func MountEscape() {

	mounts, _ := util.GetMountInfo()

	for _, m := range mounts {

		// TODO: why so may null byte in the Mounts
		// [Information Gathering - Mounts]
		// :    -
		if m.Major == "" {
			continue
		}

		// ? why match those mount points?
		if strings.Contains(m.Device, "/") || strings.Contains(m.Fstype, "ext") {
			matched, _ := regexp.MatchString("/kubelet/|/dev/[\\w-]*?\\blog$|/etc/host[\\w]*?$|/etc/[\\w]*?\\.conf$", m.Root)
			if !matched {
				m.Root = util.RedBold.Sprint(m.Root)
				m.Fstype = util.RedBold.Sprint(m.Fstype)
			}
		}

		// find lxcfs mount point for escape exploit
		if m.Device == "lxcfs" && util.StringContains(m.Opts, "rw") {
			fmt.Printf("Find mounted lxcfs with rw flags, run `%s` or `%s` to escape container!\n", util.RedBold.Sprint("cdk run lxcfs-rw"), util.RedBold.Sprint("cdk run lxcfs-rw-cgroup"))
			m.Device = util.RedBold.Sprint(m.Device)
			m.MountPoint = util.RedBold.Sprint(m.Device)
		}

		fmt.Println(m.String())

	}
}

func init() {
	RegisterSimpleCheck(CategoryMounts, "mounts.escape", "Inspect mount escape opportunities", MountEscape)
}
