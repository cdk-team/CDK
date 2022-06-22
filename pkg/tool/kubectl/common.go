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

package kubectl

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/cdk-team/CDK/conf"
	"github.com/cdk-team/CDK/pkg/errors"
	"github.com/cdk-team/CDK/pkg/util"
)

// MaybeSuccessfulStatuscodeList from https://www.w3.org/Protocols/HTTP/HTRESP.html
var MaybeSuccessfulStatuscodeList = []int{
	100, // RFC 7231, 6.2.1
	101, // RFC 7231, 6.2.2
	102, // RFC 2518, 10.1
	103, // RFC 8297

	200, // RFC 7231, 6.3.1
	201, // RFC 7231, 6.3.2
	202, // RFC 7231, 6.3.3
	203, // RFC 7231, 6.3.4
	204, // RFC 7231, 6.3.5
	205, // RFC 7231, 6.3.6
	206, // RFC 7233, 4.1
	207, // RFC 4918, 11.1
	208, // RFC 5842, 7.1
	226, // RFC 3229, 10.4.1
}

type K8sRequestOption struct {
	TokenPath string
	Token     string
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
	var tokenErr error
	if opts.Anonymous {
		opts.Token = ""
	} else if opts.TokenPath != "" {
		opts.Token, tokenErr = GetServiceAccountToken(opts.TokenPath)
	} else if opts.Token == "" {
		opts.Token, tokenErr = GetServiceAccountToken(conf.K8sSATokenDefaultPath)
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
	if len(opts.Token) > 0 {
		token := strings.TrimSpace(opts.Token)
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

	res := string(content)

	// Fix a bug reported by the author of crossc2 on whc2021.
	// When the DeployBackdoorDaemonset call fails and returns an error, it will still feedback true.
	if !util.IntContains(MaybeSuccessfulStatuscodeList, resp.StatusCode) {
		errMsg := fmt.Sprintf("err found in post request, error response code: %v.", resp.Status)
		return res, &errors.CDKRuntimeError{
			Err:       err,
			CustomMsg: errMsg,
		}
	}

	return res, nil
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
