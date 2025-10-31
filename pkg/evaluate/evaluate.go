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

package evaluate

import "log"

// CallBasics executes the baseline evaluation profile.
func CallBasics() {
	if err := NewEvaluator().RunProfile(ProfileBasic, nil); err != nil {
		log.Printf("basic evaluation failed: %v", err)
	}
}

// CallAddedFunc executes the additional evaluation profile.
func CallAddedFunc() {
	if err := NewEvaluator().RunProfile(ProfileAdditional, nil); err != nil {
		log.Printf("additional evaluation failed: %v", err)
	}
}
