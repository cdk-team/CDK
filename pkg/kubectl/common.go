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
func ServerAccountRequest(token string, method string, api string, args string) string {

	cli := goz.NewClient()
	switch strings.ToLower(method) {
	case "get":
		if len(token) > 0 {
			resp, err := cli.Get(api, goz.Options{
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
			resp, err := cli.Get(api)
			if err != nil {
				log.Fatalln(err)
			}
			r, _ := resp.GetBody()
			return r.String()
		}

	case "post":
		if len(token) > 0 {
			resp, err := cli.Post(api, goz.Options{
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
			resp, err := cli.Post(api)
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
	method := args[0]
	url := args[1]
	token, err := GetServiceAccountToken(conf.K8sSATokenDefaultPath)
	if err != nil {
		log.Fatal("reading service-account token failed, err:", err)
	}

	switch len(args) {
	case 3:
		postArgs := args[2]
		fmt.Println("post data:", postArgs)
		log.Println("response:")
		fmt.Println(ServerAccountRequest(token, method, url, postArgs))
	case 2:
		log.Println("response:")
		fmt.Println(ServerAccountRequest(token, method, url, ""))
	default:
		fatalWithUsage()
	}
}
