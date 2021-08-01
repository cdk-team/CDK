package kubectl

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/cdk-team/CDK/conf"
	"github.com/cdk-team/CDK/pkg/errors"
)

type K8sRequestOption struct {
	TokenPath string
	Server    string
	Api       string
	Method    string
	PostData  string
	Url       string
	Anonymous bool
}

func ApiServerAddr() (string, error) {
	protocol := ""
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		text := "err: cannot find kubernetes api host in ENV"
		return "", errors.New(text)
	}
	if port == "8080" || port == "8001" {
		protocol = "http://"
	} else {
		protocol = "https://"
	}
	return protocol + net.JoinHostPort(host, port), nil
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

	// http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	var request *http.Request
	opts.Method = strings.ToUpper(opts.Method)

	request, err := http.NewRequest(opts.Method, opts.Url, bytes.NewBuffer([]byte(opts.PostData)))
	if err != nil {
		return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "err found while generate post request in net.http ."}
	}

	// set request header
	if opts.Method == "POST" {
		request.Header.Set("Content-Type", "application/json")
	}
	// auth token
	if len(token) > 0 {
		token = strings.TrimSpace(token)
		request.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(request)
	if err != nil {
		return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "err found in post request."}
	}
	//defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "err found in post request."}
	}
	return string(content), nil
}

func GetServerVersion(serverAddr string) (string, error) {
	opts := K8sRequestOption{
		TokenPath: "",
		Server:    serverAddr,
		Api:       "/version",
		Method:    "GET",
		PostData:  "",
		Anonymous: true,
	}
	resp, err := ServerAccountRequest(opts)
	if err != nil {
		return "", err
	}
	// use regexp to find gitVersion
	versionPattern := regexp.MustCompile(`"gitVersion":.*?"(.*?)"`)
	results := versionPattern.FindStringSubmatch(resp)
	if len(results) != 2 {
		return "", errors.New("field gitVersion not found in response")
	}
	return results[1], nil
}
