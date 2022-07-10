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

package cli_test

import (
	"bytes"
	// "io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cdk-team/CDK/pkg/cli"
	_ "github.com/cdk-team/CDK/pkg/exploit" // register all exploits
)

type testArgsCase struct {
	name string
	args []string
	successStr string
}

func doParseCDKMainWithTimeout() {

	result := make(chan bool, 1)

	go func ()  {
		result <- cli.ParseCDKMain()
	}()

	select {
	case <-time.After(time.Second * 2):
		log.Println("check run ok, timeout in 2s, and return.")
		return
	case <-result:
		return 
	}

}

func TestParseCDKMain(t *testing.T) {


	// ./cdk eva 2>&1 | head
	// ./cdk run test-poc | head
	// ./cdk ifconfig | head

	tests := []testArgsCase{
		{
			name: "./cdk eva",
			args: []string{"./cdk_cli_path", "eva"},
			successStr: "current user",
		},
		{
			name: "./cdk run test-poc",
			args: []string{"./cdk_cli_path", "run", "test-poc"},
			successStr: "run success",
		},
		{
			name: "./cdk ifconfig",
			args: []string{"./cdk_cli_path", "ifconfig"},
			successStr: "GetLocalAddresses",
		},
	}

	for _, tt := range tests {

		// fmt.Print and log.Print to buffer, and check output
		var buf bytes.Buffer
		log.SetOutput(&buf)

		// hook fmt.X to buffer, hook os.Stdout
		// oldStdout := os.Stdout
		// r, w, _ := os.Pipe()
		// os.Stdout = w

		// hook os.Args
		args := tt.args
		os.Args = args

		t.Run(tt.name, func(t *testing.T) {
			doParseCDKMainWithTimeout()
			// out, _ := ioutil.ReadAll(r)

			// check success string in buf and out
			// if !bytes.Contains(buf.Bytes(), []byte(tt.successStr)) && !bytes.Contains(out, []byte(tt.successStr)) {
			// 	t.Errorf(("parse cdk main failed, name: %s, args: %v, buf: %s, out: %s"), tt.name, tt.args, buf.String()[:1000], string(out)[:1000])
			// }

			if !bytes.Contains(buf.Bytes(), []byte(tt.successStr)) {

				// get sub string from buf, lenght is 1000
				str := buf.String()
				if len(str) > 1000 {
					str = str[:1000]
				}

				t.Errorf(("parse cdk main failed, name: %s, args: %v, buf: %s"), tt.name, tt.args, str)
			}


		})

		// return to os.Stdout default
		// os.Stdout = oldStdout
		// w.Close()
	}
}
