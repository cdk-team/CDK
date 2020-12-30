package probe

import (
	"fmt"
	"github.com/Xyntax/CDK/conf"
	"os/exec"
	"strconv"
	"strings"
)

type FromTo struct {
	Desc string
	From int
	To   int
}

func GetTaskPortList() ([]FromTo, int) {
	res := make([]FromTo, 0)
	tot := 0

	for desc, port := range conf.TCPScannerConf.PortList {
		from := 0
		to := 0
		fromTo := strings.Split(port, "-")
		from, _ = strconv.Atoi(fromTo[0])
		to = from
		if len(fromTo) == 2 {
			to, _ = strconv.Atoi(fromTo[1])
		}
		a := FromTo{
			Desc: desc,
			From: from,
			To:   to,
		}
		res = append(res, a)
		tot += 1 + to - from
	}
	return res, tot
}

func GetTaskPortListByString(s string) ([]FromTo, int) {
	res := make([]FromTo, 0)
	tot := 0

	for _, port := range strings.Split(s, ",") {
		from := 0
		to := 0
		fromTo := strings.Split(port, "-")
		from, _ = strconv.Atoi(fromTo[0])
		to = from
		if len(fromTo) == 2 {
			to, _ = strconv.Atoi(fromTo[1])
		}
		a := FromTo{
			Desc: "",
			From: from,
			To:   to,
		}
		res = append(res, a)
		tot += 1 + to - from
	}
	return res, tot
}

func GetTaskIPList(ip string) (base string, start, end int, err error) {
	fromTo := strings.Split(ip, "-")
	ipStart := fromTo[0]
	err = fmt.Errorf("Invalid IP Range (eg. 1.1.1.1-3)\n")

	tIp := strings.Split(ipStart, ".")
	if len(tIp) != 4 {
		return
	}
	start, _ = strconv.Atoi(tIp[3])
	end = start
	if len(fromTo) == 2 {
		end, _ = strconv.Atoi(fromTo[1])
	}
	if end == 0 {
		return
	}
	base = fmt.Sprintf("%s.%s.%s", tIp[0], tIp[1], tIp[2])
	err = nil
	return
}

func Ulimit() int64 {
	out, err := exec.Command("ulimit", "-n").Output()
	if err != nil {
		panic(err)
	}

	s := strings.TrimSpace(string(out))

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}
