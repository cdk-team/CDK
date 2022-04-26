
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

package util

import (
	"bytes"
	"context"
	"github.com/cdk-team/CDK/pkg/errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// ref https://docs.docker.com/engine/api/v1.24/
func UnixHttpSend(method string, unixPath string, uri string, data string) (string, error) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", unixPath)
			},
		},
	}

	var response *http.Response
	var err error

	switch method {
	case "post":
		response, err = httpc.Post(uri, "application/json", strings.NewReader(data))
	case "get":
		response, err = httpc.Get(uri)
	}

	if err != nil {
		return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "Unix HTTP Request failed."}
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, response.Body)
	return buf.String(), nil
}

func HttpSendJson(method string, url string, data string) (string, error) {
	req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "HTTP Request failed."}
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "HTTP Request failed."}
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
