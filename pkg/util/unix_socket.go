package util

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"strings"
)

// ref https://docs.docker.com/engine/api/v1.24/
func UnixHttpSend(method string, unixPath string, uri string, data string) string {
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
		panic(err)
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, response.Body)
	return buf.String()
}
