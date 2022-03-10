package util

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/cdk-team/CDK/pkg/errors"
)

func ByteToString(orig []byte) string {
	n := -1
	l := -1
	for i, b := range orig {
		// skip left side null
		if l == -1 && b == 0 {
			continue
		}
		if l == -1 {
			l = i
		}

		if b == 0 {
			break
		}
		n = i + 1
	}
	if n == -1 {
		return string(orig)
	}
	return string(orig[l:n])
}

func ContainsStrInSlice(elem string, sli []string) bool {
	// grabbed from https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
	for _, a := range sli {
		if a == elem {
			return true
		}
	}
	return false
}

func RandString(n int) string {
	// grabbed from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	const (
		letterBytes   = "abcde1fghij2klmno3pqrst4uvwxy5zABCD6EFGHI7JKLMN8OPQRS9TUVWX9YZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	sb := strings.Builder{}
	sb.Grow(n)
	rand.Seed(time.Now().UnixNano())
	// rand.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func RemoveDuplicateElement(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	temp := map[string]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// dataFromSliceOrFile returns data from the slice (if non-empty), or from the file,
// or an error if an error occurred reading the file
func dataFromSliceOrFile(data []byte, file string) ([]byte, error) {
	if len(data) > 0 {
		return data, nil
	}
	if len(file) > 0 {
		fileData, err := ioutil.ReadFile(file)
		if err != nil {
			return []byte{}, err
		}
		return fileData, nil
	}
	return nil, nil
}

// ShellExec run shell script by bash
func ShellExec(shellPath string) error {
	var command = shellPath
	if strings.HasPrefix(shellPath, "/") {
		command = shellPath
	} else {
		command = fmt.Sprintf("./%s .", shellPath)
	}
	cmd := exec.Command("/bin/bash", "-c", command)

	output, err := cmd.Output()
	if err != nil {
		return &errors.CDKRuntimeError{Err: err, CustomMsg: fmt.Sprintf("Execute Shell:%s failed", command)}
	}
	fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(output))
	return nil
}

// StringContains check string array contains a string
func StringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// IntContains check string array contains a int number
func IntContains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
