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
	"fmt"
	"log"
)

const Colorful = true

// fmt.Printf(util.GreenBold.Sprint("\n[Information Gathering - System Info]\n"))
func PrintH2(title string) {
	fmt.Printf(BlueBold.Sprint("\n[  ") + GreenBold.Sprint(title) + BlueBold.Sprint("  ]\n"))
}

func PrintItemKey(key string, color bool) {
	key = key + "\n"
	if color {
		log.Printf(YellowBold.Sprint(key))
	} else {
		log.Printf(key)
	}
}

func PrintItemValue(value string, color bool) {
	value = "\t" + value + "\n"
	if color {
		fmt.Printf(RedBold.Sprint(value))
	} else {
		fmt.Printf(value)
	}
}

func PrintItemValueWithKeyOneLine(key, value string, color bool) {
	if color {
		log.Printf("%s: %s", key, GreenBold.Sprint(value))
	} else {
		log.Printf("%s: %s", key, value)
	}
}

func PrintOrignal(out string) {
	fmt.Printf("%s\n", out)
}

