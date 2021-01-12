package evaluate

import (
	"fmt"
	"github.com/cdk-team/CDK/conf"
	"github.com/cdk-team/CDK/pkg/util"
	"os"
	"path/filepath"
	"strings"
)

func SearchLocalFilePath() {

	filepath.Walk(conf.SensitiveFileConf.StartDir, func(path string, info os.FileInfo, err error) error {
		for _, name := range conf.SensitiveFileConf.NameList {
			currentPath := strings.ToLower(path)
			//if util.IsSoftLink(currentPath) && util.IsDir(currentPath) {
			//	fmt.Println("skip", currentPath)
			//	return filepath.SkipDir // skip soft link or it will run into container runtime filesystem
			//}
			if strings.Contains(currentPath, name) {
				fmt.Printf("\t%s - %s\n", name, path)
				if util.IsDir(currentPath) {
					return filepath.SkipDir // stop dive if sensitive dir found
				}
				return nil
			}
		}
		return nil
	})

}
