package kubectl

import (
	"errors"
	"fmt"
	"github.com/Xyntax/CDK/conf"
	"github.com/idoubi/goz"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

type K8sRequestOption struct {
	TokenPath string
	Server string
	Api string
	Method string
	Args string
	Anonymous bool
}

func ApiServerAddr() (string, error) {
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		text := "err: cannot find kubernetes api host in ENV"
		return "", errors.New(text)
	}
	return "https://" + net.JoinHostPort(host, port), nil
}

func GetServiceAccountToken(tokenPath string) (string, error) {
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

/*
curl -s https://192.168.0.234:6443/api/v1/nodes?watch  --header "Authorization: Bearer xxx" --cacert ca.crt
*/
//https://github.com/kubernetes/client-go/blob/66db2540991da169fb60fce735064a55bfc52b71/rest/config.go#L483
func ServerAccountRequest(opts K8sRequestOption) string {

	// parse token
	var token string
	var tokenErr error
	if opts.Anonymous {
		token=""
	}else if opts.TokenPath == "" {
		token,tokenErr = GetServiceAccountToken(conf.K8sSATokenDefaultPath)
	}else{
		token,tokenErr = GetServiceAccountToken(opts.TokenPath)
	}
	if tokenErr != nil {
		log.Println(tokenErr)
		return ""
	}

	// parse url
	var url string
	var server string
	var urlErr error
	if opts.Server == "" {
		server,urlErr = ApiServerAddr()
		url = server+opts.Api
	} else {
		url = opts.Server+opts.Api
		urlErr = nil
	}
	if urlErr != nil {
		log.Println(urlErr)
		return ""
	}

	cli := goz.NewClient()
	switch strings.ToLower(opts.Method) {
	case "get":
		if len(token) > 0 {
			resp, err := cli.Get(url, goz.Options{
				Headers: map[string]interface{}{
					"Authorization": "Bearer " + string(token),
				},
				//FormParams: map[string]interface{}{
				//	"key1": "value1",
				//	"key2": []string{"value21", "value22"},
				//	"key3": "333",
				//},
			})
			if err != nil {
				log.Fatalln(err)
			}
			r, _ := resp.GetBody()
			return r.String()
		} else {
			resp, err := cli.Get(url)
			if err != nil {
				log.Fatalln(err)
			}
			r, _ := resp.GetBody()
			return r.String()
		}

	case "post":
		if len(token) > 0 {
			resp, err := cli.Post(url, goz.Options{
				Headers: map[string]interface{}{
					"Authorization": "Bearer " + string(token),
				},
				//FormParams: map[string]interface{}{
				//	"key1": "value1",
				//	"key2": []string{"value21", "value22"},
				//	"key3": "333",
				//},
			})
			if err != nil {
				log.Fatalln(err)
			}
			r, _ := resp.GetBody()
			return r.String()
		} else {
			resp, err := cli.Post(url)
			if err != nil {
				log.Fatalln(err)
			}
			r, _ := resp.GetBody()
			return r.String()
		}

	default:
		fatalWithUsage()
	}
	return ""
}

func fatalWithUsage() {
	log.Println("invalid input args")
	fmt.Printf("\nusage:\n\t./cdk kcurl get \"https://192.168.0.234:6443/api/v1/nodes\"\n")
	fmt.Printf("\t./cdk kcurl post \"https://192.168.0.234:6443/api/v1/nodes\" \"k1=v1&k2=v2\"\n")
	log.Fatal()
}

func KubectlMain(args []string) {
	// kubectl get https://1/2 id=1 --anonymous
	fmt.Println(args)
	//method := args[0]
	//url := args[1]
	//
	//var opts K8sRequestOption
	//
	//{
	//	tokenPath: "",
	//	server:    "",
	//	api:       "",
	//	method:    "",
	//	args:      "",
	//	anonymous: false,
	//}
	//
	//opts.anonymous = false
	//for _, arg.():= range args{
	//	if arg == "--anonymous"{
	//		opts.anonymous=true
	//	}
	//}
	//
	//switch len(args) {
	//case 3:
	//	postArgs := args[2]
	//	fmt.Println("post data:", postArgs)
	//	log.Println("response:")
	//	fmt.Println(ServerAccountRequest())
	//case 2:
	//	log.Println("response:")
	//	fmt.Println(ServerAccountRequest(token, method, url, ""))
	//default:
	//	fatalWithUsage()
	//}
}
