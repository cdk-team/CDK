
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

package capability

import (
	"fmt"
	"strconv"
	"strings"
)

/**
 ** POSIX-draft defined capabilities.
 ** Refer to the code https://github.com/torvalds/linux/blob/master/include/uapi/linux/capability.h by torvalds.
 ** Code written by neargle on 2021-02-08, last commit id 9d3a39a5f1e45827b008fff1ee9cf3cac3409665.
 **/

/* In a system with the [_POSIX_CHOWN_RESTRICTED] option defined, this
   overrides the restriction of changing file ownership and group
   ownership. */

var CAP_CHOWN = 0

/* Override all DAC access, including ACL execute access if
   [_POSIX_ACL] is defined. Excluding DAC access covered by
   CAP_LINUX_IMMUTABLE. */

var CAP_DAC_OVERRIDE = 1

/* Overrides all DAC restrictions regarding read and search on files
   and directories, including ACL restrictions if [_POSIX_ACL] is
   defined. Excluding DAC access covered by CAP_LINUX_IMMUTABLE. */

var CAP_DAC_READ_SEARCH = 2

/* Overrides all restrictions about allowed operations on files, where
   file owner ID must be equal to the user ID, except where CAP_FSETID
   is applicable. It doesn't override MAC and DAC restrictions. */

var CAP_FOWNER = 3

/* Overrides the following restrictions that the effective user ID
   shall match the file owner ID when setting the S_ISUID and S_ISGID
   bits on that file; that the effective group ID (or one of the
   supplementary group IDs) shall match the file owner ID when setting
   the S_ISGID bit on that file; that the S_ISUID and S_ISGID bits are
   cleared on successful return from chown(2) (not implemented). */

var CAP_FSETID = 4

/* Overrides the restriction that the real or effective user ID of a
   process sending a signal must match the real or effective user ID
   of the process receiving the signal. */

var CAP_KILL = 5

/* Allows setgid(2) manipulation */
/* Allows setgroups(2) */
/* Allows forged gids on socket credentials passing. */

var CAP_SETGID = 6

/* Allows set*uid(2) manipulation (including fsuid). */
/* Allows forged pids on socket credentials passing. */

var CAP_SETUID = 7

/**
 ** Linux-specific capabilities
 **/

/* Without VFS support for capabilities:
 *   Transfer any capability in your permitted set to any pid,
 *   remove any capability in your permitted set from any pid
 * With VFS support for capabilities (neither of above, but)
 *   Add any capability from current's capability bounding set
 *       to the current process' inheritable set
 *   Allow taking bits out of capability bounding set
 *   Allow modification of the securebits for a process
 */

var CAP_SETPCAP = 8

/* Allow modification of S_IMMUTABLE and S_APPEND file attributes */

var CAP_LINUX_IMMUTABLE = 9

/* Allows binding to TCP/UDP sockets below 1024 */
/* Allows binding to ATM VCIs below 32 */

var CAP_NET_BIND_SERVICE = 10

/* Allow broadcasting, listen to multicast */

var CAP_NET_BROADCAST = 11

/* Allow interface configuration */
/* Allow administration of IP firewall, masquerading and accounting */
/* Allow setting debug option on sockets */
/* Allow modification of routing tables */
/* Allow setting arbitrary process / process group ownership on
   sockets */
/* Allow binding to any address for transparent proxying (also via NET_RAW) */
/* Allow setting TOS (type of service) */
/* Allow setting promiscuous mode */
/* Allow clearing driver statistics */
/* Allow multicasting */
/* Allow read/write of device-specific registers */
/* Allow activation of ATM control sockets */

var CAP_NET_ADMIN = 12

/* Allow use of RAW sockets */
/* Allow use of PACKET sockets */
/* Allow binding to any address for transparent proxying (also via NET_ADMIN) */

var CAP_NET_RAW = 13

/* Allow locking of shared memory segments */
/* Allow mlock and mlockall (which doesn't really have anything to do
   with IPC) */

var CAP_IPC_LOCK = 14

/* Override IPC ownership checks */

var CAP_IPC_OWNER = 15

/* Insert and remove kernel modules - modify kernel without limit */
var CAP_SYS_MODULE = 16

/* Allow ioperm/iopl access */
/* Allow sending USB messages to any device via /dev/bus/usb */

var CAP_SYS_RAWIO = 17

/* Allow use of chroot() */

var CAP_SYS_CHROOT = 18

/* Allow ptrace() of any process */

var CAP_SYS_PTRACE = 19

/* Allow configuration of process accounting */

var CAP_SYS_PACCT = 20

/* Allow configuration of the secure attention key */
/* Allow administration of the random device */
/* Allow examination and configuration of disk quotas */
/* Allow setting the domainname */
/* Allow setting the hostname */
/* Allow calling bdflush() */
/* Allow mount() and umount(), setting up new smb connection */
/* Allow some autofs root ioctls */
/* Allow nfsservctl */
/* Allow VM86_REQUEST_IRQ */
/* Allow to read/write pci config on alpha */
/* Allow irix_prctl on mips (setstacksize) */
/* Allow flushing all cache on m68k (sys_cacheflush) */
/* Allow removing semaphores */
/* Used instead of CAP_CHOWN to "chown" IPC message queues, semaphores
   and shared memory */
/* Allow locking/unlocking of shared memory segment */
/* Allow turning swap on/off */
/* Allow forged pids on socket credentials passing */
/* Allow setting readahead and flushing buffers on block devices */
/* Allow setting geometry in floppy driver */
/* Allow turning DMA on/off in xd driver */
/* Allow administration of md devices (mostly the above, but some
   extra ioctls) */
/* Allow tuning the ide driver */
/* Allow access to the nvram device */
/* Allow administration of apm_bios, serial and bttv (TV) device */
/* Allow manufacturer commands in isdn CAPI support driver */
/* Allow reading non-standardized portions of pci configuration space */
/* Allow DDI debug ioctl on sbpcd driver */
/* Allow setting up serial ports */
/* Allow sending raw qic-117 commands */
/* Allow enabling/disabling tagged queuing on SCSI controllers and sending
   arbitrary SCSI commands */
/* Allow setting encryption key on loopback filesystem */
/* Allow setting zone reclaim policy */
/* Allow everything under CAP_BPF and CAP_PERFMON for backward compatibility */

var CAP_SYS_ADMIN = 21

/* Allow use of reboot() */

var CAP_SYS_BOOT = 22

/* Allow raising priority and setting priority on other (different
   UID) processes */
/* Allow use of FIFO and round-robin (realtime) scheduling on own
   processes and setting the scheduling algorithm used by another
   process. */
/* Allow setting cpu affinity on other processes */
/* Allow setting realtime ioprio class */
/* Allow setting ioprio class on other processes */

var CAP_SYS_NICE = 23

/* Override resource limits. Set resource limits. */
/* Override quota limits. */
/* Override reserved space on ext2 filesystem */
/* Modify data journaling mode on ext3 filesystem (uses journaling
   resources) */
/* NOTE: ext2 honors fsuid when checking for resource overrides, so
   you can override using fsuid too */
/* Override size restrictions on IPC message queues */
/* Allow more than 64hz interrupts from the real-time clock */
/* Override max number of consoles on console allocation */
/* Override max number of keymaps */
/* Control memory reclaim behavior */

var CAP_SYS_RESOURCE = 24

/* Allow manipulation of system clock */
/* Allow irix_stime on mips */
/* Allow setting the real-time clock */

var CAP_SYS_TIME = 25

/* Allow configuration of tty devices */
/* Allow vhangup() of tty */

var CAP_SYS_TTY_CONFIG = 26

/* Allow the privileged aspects of mknod() */

var CAP_MKNOD = 27

/* Allow taking of leases on files */

var CAP_LEASE = 28

/* Allow writing the audit log via unicast netlink socket */

var CAP_AUDIT_WRITE = 29

/* Allow configuration of audit via unicast netlink socket */

var CAP_AUDIT_CONTROL = 30

/* Set or remove capabilities on files */

var CAP_SETFCAP = 31

/* Override MAC access.
   The base kernel enforces no MAC policy.
   An LSM may enforce a MAC policy, and if it does and it chooses
   to implement capability based overrides of that policy, this is
   the capability it should use to do so. */

var CAP_MAC_OVERRIDE = 32

/* Allow MAC configuration or state changes.
   The base kernel requires no MAC configuration.
   An LSM may enforce a MAC policy, and if it does and it chooses
   to implement capability based checks on modifications to that
   policy or the data required to maintain it, this is the
   capability it should use to do so. */

var CAP_MAC_ADMIN = 33

/* Allow configuring the kernel's syslog (printk behaviour) */

var CAP_SYSLOG = 34

/* Allow triggering something that will wake the system */

var CAP_WAKE_ALARM = 35

/* Allow preventing system suspends */

var CAP_BLOCK_SUSPEND = 36

/* Allow reading the audit log via multicast netlink socket */

var CAP_AUDIT_READ = 37

/*
 * Allow system performance and observability privileged operations
 * using perf_events, i915_perf and other kernel subsystems
 */

var CAP_PERFMON = 38

/*
 * CAP_BPF allows the following BPF operations:
 * - Creating all types of BPF maps
 * - Advanced verifier features
 *   - Indirect variable access
 *   - Bounded loops
 *   - BPF to BPF function calls
 *   - Scalar precision tracking
 *   - Larger complexity limits
 *   - Dead code elimination
 *   - And potentially other features
 * - Loading BPF Type Format (BTF) data
 * - Retrieve xlated and JITed code of BPF programs
 * - Use bpf_spin_lock() helper
 *
 * CAP_PERFMON relaxes the verifier checks further:
 * - BPF progs can use of pointer-to-integer conversions
 * - speculation attack hardening measures are bypassed
 * - bpf_probe_read to read arbitrary kernel memory is allowed
 * - bpf_trace_printk to print kernel memory is allowed
 *
 * CAP_SYS_ADMIN is required to use bpf_probe_write_user.
 *
 * CAP_SYS_ADMIN is required to iterate system wide loaded
 * programs, maps, links, BTFs and convert their IDs to file descriptors.
 *
 * CAP_PERFMON and CAP_BPF are required to load tracing programs.
 * CAP_NET_ADMIN and CAP_BPF are required to load networking programs.
 */
var CAP_BPF = 39

/* Allow checkpoint/restore related operations */
/* Allow PID selection during clone3() */
/* Allow writing to ns_last_pid */

var CAP_CHECKPOINT_RESTORE = 40

type Cap struct {
	int
}

var DockerDefaultCaps = []string{
	"CAP_CHOWN",
	"CAP_DAC_OVERRIDE",
	"CAP_FOWNER",
	"CAP_FSETID",
	"CAP_KILL",
	"CAP_SETGID",
	"CAP_SETUID",
	"CAP_SETPCAP",
	"CAP_NET_BIND_SERVICE",
	"CAP_NET_RAW",
	"CAP_SYS_CHROOT",
	"CAP_MKNOD",
	"CAP_AUDIT_WRITE",
	"CAP_SETFCAP",
}

// CAPStringsList cap strings with order
var CAPStringsList = []string{
	"CAP_CHOWN",
	"CAP_DAC_OVERRIDE",
	"CAP_DAC_READ_SEARCH",
	"CAP_FOWNER",
	"CAP_FSETID",
	"CAP_KILL",
	"CAP_SETGID",
	"CAP_SETUID",
	"CAP_SETPCAP",
	"CAP_LINUX_IMMUTABLE",
	"CAP_NET_BIND_SERVICE",
	"CAP_NET_BROADCAST",
	"CAP_NET_ADMIN",
	"CAP_NET_RAW",
	"CAP_IPC_LOCK",
	"CAP_IPC_OWNER",
	"CAP_SYS_MODULE",
	"CAP_SYS_RAWIO",
	"CAP_SYS_CHROOT",
	"CAP_SYS_PTRACE",
	"CAP_SYS_PACCT",
	"CAP_SYS_ADMIN",
	"CAP_SYS_BOOT",
	"CAP_SYS_NICE",
	"CAP_SYS_RESOURCE",
	"CAP_SYS_TIME",
	"CAP_SYS_TTY_CONFIG",
	"CAP_MKNOD",
	"CAP_LEASE",
	"CAP_AUDIT_WRITE",
	"CAP_AUDIT_CONTROL",
	"CAP_SETFCAP",
	"CAP_MAC_OVERRIDE",
	"CAP_MAC_ADMIN",
	"CAP_SYSLOG",
	"CAP_WAKE_ALARM",
	"CAP_BLOCK_SUSPEND",
	"CAP_AUDIT_READ",
	"CAP_PERFMON",
	"CAP_BPF",
	"CAP_CHECKPOINT_RESTORE",
}

func CapHexParser(capHex string) ([]string, error) {
	var capTextList []string

	numb, err := strconv.ParseUint(capHex, 16, 64)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(CAPStringsList); i++ {
		flag := numb & 0x1
		if flag == uint64(1) {
			capTextList = append(capTextList, CAPStringsList[i])
		}
		numb = numb >> 1
	}

	return capTextList, nil
}

func CapListToString(lst []string) string {
	return strings.Join(lst[:], ",")
}

func CapHexToText(capHex string) string {

	lst, err := CapHexParser(capHex)
	if err != nil {
		fmt.Printf("CapHexParser error: %s", err)
		return ""
	}
	return strings.Join(lst[:], ",")
}
