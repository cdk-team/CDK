package conf

import "time"

// TCP port scanner
type TCPScannerConfS struct {
	Timeout     time.Duration
	MaxParallel int64
	PortList    map[string]string
}

var TCPScannerConf = TCPScannerConfS{
	Timeout:     500 * time.Millisecond,
	MaxParallel: 50,
	PortList: map[string]string{
		"ssh":                 "22",
		"http":                "80",
		"https":               "443",
		"docker-api":          "2375",
		"etcd":                "2379",
		"cAdvisor":            "4194",
		"k8s-api-server":      "6443",
		"http-1":              "8080",
		"https-1":             "8443",
		"kubelet-auth":        "10250",
		"kubelet-read":        "10255",
		"dashboard":           "30000",
		"nodeport-service":    "30001-32767", //default NodePort service port range：30000-32767。
		"tiller,weave,calico": "44134",
	},
}
