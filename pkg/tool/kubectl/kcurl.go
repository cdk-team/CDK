package kubectl

import (
	"fmt"
	"github.com/Xyntax/CDK/conf"
	"log"
)

var kcurlBanner = `kcurl - send HTTP request to K8s api-server.

Usage:
  ./cdk kcurl (<token_path>|anonymous|default) (get|post) <url> ["post_args"]

Options:
  token_path  connect api-server with user-specified service-account token.
  anonymous   connect api-server using system:anonymous service-account
  default     connect api-server using pod default service-account token.

Example: 
  ./cdk kcurl default get "https://192.168.0.234:6443/api/v1/nodes"
  ./cdk kcurl /var/run/secrets/kubernetes.io/serviceaccount/token get "https://192.168.0.234:6443/api/v1/nodes"
  ./cdk kcurl anonymous post "https://192.168.0.234:6443/api/v1/nodes" "k1=v1&k2=v2"
`


func KubectlToolApi(args []string) {

	var opts = K8sRequestOption{}

	// err break
	if len(args) != 3 && len(args) != 4 {
		fmt.Println(kcurlBanner)
		return
	}

	switch args[0] {
	case "default":
		opts.TokenPath = conf.K8sSATokenDefaultPath
	case "anonymous":
		opts.TokenPath = ""
		opts.Anonymous = true
	default:
		opts.TokenPath = args[0]
	}

	switch args[1] {
	case "post":
		opts.Method = "post"
	case "get":
		opts.Method = "get"
	default: // err break
		fmt.Println(kcurlBanner)
		return
	}

	if len(args) == 3 {
		opts.Url = args[2]
	} else {
		opts.Url = args[2]
		fmt.Println("post data:", opts.Args)
		opts.Args = args[3]
	}

	resp, err := ServerAccountRequest(opts)
	if err != nil {
		log.Println("failed to get api-server response")
		fmt.Println(err)
	} else {
		if len(resp) == 0{
			log.Println("empty server response body, try `https://` instead of `http://` in url.")
		}
		log.Println("api-server response:")
		fmt.Println(resp)
	}
}
