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

package etcdctl

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cdk-team/CDK/pkg/errors"
	"github.com/tidwall/gjson"
)

type EtcdRequestOption struct {
	Endpoint  string
	Api       string
	PostData  string
	TlsConfig *tls.Config
	Anonymous bool
	Silent    bool
}

func DoRequest(opt EtcdRequestOption) (map[string]string, error) {
	// http client
	if opt.TlsConfig == nil || len(opt.TlsConfig.Certificates) == 0 || opt.TlsConfig.RootCAs == nil {
		opt.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: opt.TlsConfig,
		},
		Timeout: time.Duration(5) * time.Second,
	}

	request, err := http.NewRequest("POST", opt.Endpoint+opt.Api, bytes.NewBuffer([]byte(opt.PostData)))
	if err != nil {
		return nil, &errors.CDKRuntimeError{Err: err, CustomMsg: "err found while generate post request in net.http ."}
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if resp != nil {
		defer resp.Body.Close()
	} else if err != nil {
		return nil, &errors.CDKRuntimeError{Err: err, CustomMsg: "err found in post request."}
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.CDKRuntimeError{Err: err, CustomMsg: "err found in post request."}
	}

	kvs := gjson.Get(string(content), "kvs").Array()
	ret := make(map[string]string, len(kvs))
	for _, k := range kvs {
		name, err := base64.StdEncoding.DecodeString(k.Get("key").String())
		if err != nil {
			fmt.Println("base64 decode failed:", err.Error())
			continue
		}

		ret[string(name)] = ""
		if !opt.Silent {
			fmt.Println(string(name))
		}

		if k.Get("value").Exists() {
			v, _ := base64.StdEncoding.DecodeString(k.Get("value").String())
			if !opt.Silent {
				fmt.Println(string(v))
			}
			ret[string(name)] = string(v)
		}
	}
	return ret, nil
}

func GenerateQuery(key string) (query string) {
	b64key := base64.StdEncoding.EncodeToString([]byte(strings.TrimSuffix(key, "\n")))
	if key == "/" {
		bzero := base64.StdEncoding.EncodeToString([]byte{0})
		query = fmt.Sprintf("{\"range_end\": \"%s\", \"key\": \"%s\", \"keys_only\":true}", bzero, b64key)
	} else {
		query = fmt.Sprintf("{\"key\": \"%s\"}", b64key)
	}
	return
}
