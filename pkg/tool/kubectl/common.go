package kubectl

import (
	"github.com/Xyntax/CDK/conf"
	"github.com/Xyntax/CDK/pkg/errors"
	"github.com/idoubi/goz"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

type K8sRequestOption struct {
	TokenPath string
	Server    string
	Api       string
	Method    string
	Args      string
	Url       string
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
func ServerAccountRequest(opts K8sRequestOption) (string, error) {

	// parse token
	var token string
	var tokenErr error
	if opts.Anonymous {
		token = ""
	} else if opts.TokenPath == "" {
		token, tokenErr = GetServiceAccountToken(conf.K8sSATokenDefaultPath)
	} else {
		token, tokenErr = GetServiceAccountToken(opts.TokenPath)
	}
	if tokenErr != nil {
		return "", &errors.CDKRuntimeError{Err: tokenErr, CustomMsg: "load K8s service account token error."}
	}

	// parse url if opts.Url is ""
	if len(opts.Url) == 0 {
		var server string
		var urlErr error
		if opts.Server == "" {
			server, urlErr = ApiServerAddr()
			opts.Url = server + opts.Api
		} else {
			opts.Url = opts.Server + opts.Api
			urlErr = nil
		}
		if urlErr != nil {
			return "", &errors.CDKRuntimeError{Err: urlErr, CustomMsg: "err found while searching local K8s apiserver addr."}
		}
	}

	cli := goz.NewClient()
	switch strings.ToLower(opts.Method) {
	case "get":
		if len(token) > 0 {
			resp, err := cli.Get(opts.Url, goz.Options{
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
				return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "http request error."}
			}
			r, _ := resp.GetBody()
			return r.String(), nil
		} else {
			resp, err := cli.Get(opts.Url)
			if err != nil {
				return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "http request error."}
			}
			r, _ := resp.GetBody()
			return r.String(), nil
		}

	case "post":
		if len(token) > 0 {
			resp, err := cli.Post(opts.Url, goz.Options{
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
				return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "http request error."}
			}
			r, _ := resp.GetBody()
			return r.String(), nil
		} else {
			resp, err := cli.Post(opts.Url)
			if err != nil {
				return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "http request error."}
			}
			r, _ := resp.GetBody()
			return r.String(), nil
		}

	default:
		return "", errors.New("K8s request: invalid http method " + strings.ToLower(opts.Method))
	}
}
