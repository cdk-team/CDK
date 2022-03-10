package util

import (
	"io/ioutil"
	"strings"

	"github.com/cdk-team/CDK/pkg/errors"
)

// CheckUnpriUserNS checks if the current host enable unprivileged user namespace.
// reference:
//   https://blog.trailofbits.com/2019/07/19/understanding-docker-container-escapes/
//   https://unit42.paloaltonetworks.com/cve-2022-0492-cgroups/
// exceptional case:
//  the sysctl files(/proc/sys/kernel/unprivileged_userns_clone) only exist in Debian, Ubuntu.
//  We can not check the sysctl file in other distros, test in CentOS Linux release 8.4.2105 (Core).
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
