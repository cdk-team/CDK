
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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsDir(t *testing.T) {
	type testCase struct {
		dir      string
		expected bool
	}
	testGroup := map[string]testCase{
		"1": {
			dir:      "/Users/xy/go",
			expected: true,
		},
		"2": {
			dir:      "/Users/xy/",
			expected: true,
		},
		"3": {
			dir:      "/etc/hosts",
			expected: false,
		},
	}

	for key, v := range testGroup {
		t.Run(key, func(t *testing.T) {
			assert.Equal(t, v.expected, IsDir(v.dir), "That should be equal")
		})
	}

}

func TestIsSoftLink(t *testing.T) {
	type testCase struct {
		dir      string
		expected bool
	}
	testGroup := map[string]testCase{
		"1": {
			dir:      "/Users/xy/go",
			expected: true,
		},
		"2": {
			dir:      "/Users/xy/",
			expected: false,
		},
		"3": {
			dir:      "/etc/hosts",
			expected: false,
		},
	}

	for key, v := range testGroup {
		t.Run(key, func(t *testing.T) {
			assert.Equal(t, v.expected, IsSoftLink(v.dir), "That should be equal")
		})
	}

}
