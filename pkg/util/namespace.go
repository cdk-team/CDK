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
	"io/ioutil"
	"strings"

	"github.com/cdk-team/CDK/pkg/errors"
)

// CheckUnpriUserNS checks if the current host enable unprivileged user namespace.
// reference:
//
//	https://blog.trailofbits.com/2019/07/19/understanding-docker-container-escapes/
//	https://unit42.paloaltonetworks.com/cve-2022-0492-cgroups/
//
// exceptional case:
//
//	the sysctl files(/proc/sys/kernel/unprivileged_userns_clone) only exist in Debian, Ubuntu.
//	We can not check the sysctl file in other distros, test in CentOS Linux release 8.4.2105 (Core).
func CheckUnpriUserNS() error {

	data, err := ioutil.ReadFile("/proc/sys/kernel/unprivileged_userns_clone")
	if err != nil {
		return &errors.CDKRuntimeError{Err: err, CustomMsg: "check prerequisites error."}
	}

	if strings.TrimSuffix(string(data), "\n") != "1" {
		return &errors.CDKRuntimeError{Err: nil, CustomMsg: "host os does NOT enable unprivileged user namespace."}
	}

	return nil
}
