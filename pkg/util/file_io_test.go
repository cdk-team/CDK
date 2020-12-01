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
