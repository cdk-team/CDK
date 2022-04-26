
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
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"

	"github.com/cdk-team/CDK/pkg/errors"
)

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// ReadLines reads a whole file into memory
// and returns a slice of its lines.
// from https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func FileExist(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !fileInfo.IsDir()
}

func IsSoftLink(FilePath string) bool {
	fileInfo, err := os.Lstat(FilePath)
	if err != nil {
		return false
	}
	if sys := fileInfo.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			nlink := uint64(stat.Nlink)
			if nlink == 1 { // soft link ==1; hard link == 2
				return true
			}
		}
	}
	return false
}

func IsDir(FilePath string) bool {
	fileInfo, err := os.Stat(FilePath)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func RewriteFile(path string, content string, perm os.FileMode) {
	cmdFile, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		log.Fatal("overwrite file:", path, "err: "+err.Error())
	} else {
		n, _ := cmdFile.Seek(0, os.SEEK_END)
		_, err = cmdFile.WriteAt([]byte(content), n)
		log.Println("overwrite file:", path, "success.")
		defer cmdFile.Close()
	}
}

func WriteFile(path string, content string) error {
	var d = []byte(content)
	err := ioutil.WriteFile(path, d, 0666)
	if err != nil {
		return err
	}
	return nil
}

func WriteFileAdd(path string, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func WriteShellcodeToCrontab(header string, filePath string, shellcode string) error {
	shellcode = fmt.Sprintf("\n%s\n* * * * * root %s", header, shellcode)
	err := WriteFileAdd(filePath, shellcode)
	if err != nil {
		return &errors.CDKRuntimeError{Err: err, CustomMsg: "err found while writing shellcode to host crontab from container."}
	}
	return nil
}
