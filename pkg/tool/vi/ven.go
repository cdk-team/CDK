// +build !thin,!no_vi_tool


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

package vi

import (
	"github.com/bkthomps/Ven/screen"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"log"
	"os"
)

func RunVendorVi() {

	if len(os.Args) != 2 {
		print("Usage: ./cdk vi <file_name>\n")
		return
	}
	userArg := os.Args[1]
	tCellScreen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	encoding.Register()
	quit := make(chan struct{})
	s := &screen.Screen{}
	s.Init(tCellScreen, quit, userArg)
	<-quit
	tCellScreen.Fini()
}
