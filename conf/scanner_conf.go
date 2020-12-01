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
		//"test": "1-3",
		"http":             "80",
		"https":            "443",
		"ssh":              "22",
		"docker-api":       "2375",
		"http-1":           "8080",
		"https-1":          "8443",
		"k8s-api-server":   "6443",
		"kubelet-auth":     "10250",
		"kubelet-read":     "10255",
		"nodeport-service": "30000-32767", //default NodePort service port range：30000-32767。
	},
}
