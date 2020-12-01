package evaluate

import (
	"fmt"
	"github.com/Xyntax/CDK/conf"
	"github.com/Xyntax/CDK/pkg/util"
	"os"
	"path/filepath"
	"strings"
)

func SearchLocalFilePath() {
	// 1. check if system have <find> command
	//_, err := exec.LookPath("find")
	//
	//start := time.Now()
	//if err == nil {
	//	log.Println("using <find> command to search")
	//	for _, name := range findNameList {
	//		//log.Println("checking",name)
	//		execFindCommand(name)
	//	}
	//	//return
	//}
	//end := time.Now() // 打印耗时
	//fmt.Printf("执行耗时(s): %f\n", end.Sub(start).Seconds())

	// 2. full-disk scan
	//start1 := time.Now()
	filepath.Walk(conf.SensitiveFileConf.StartDir, func(path string, info os.FileInfo, err error) error {
		for _, name := range conf.SensitiveFileConf.NameList {
			currentPath := strings.ToLower(path)
			//if util.IsSoftLink(currentPath) && util.IsDir(currentPath) {
			//	fmt.Println("skip", currentPath)
			//	return filepath.SkipDir // skip soft link or it will run into container runtime filesystem
			//}
			if strings.Contains(currentPath, name) {
				fmt.Printf("\t%s - %s\n", name, path)
				if util.IsDir(currentPath){
					return filepath.SkipDir // stop dive if sensitive dir found
				}
				return nil
			}
		}
		return nil
	})
	//end1 := time.Now() // 打印耗时
	//fmt.Printf("执行耗时(s): %f\n", end1.Sub(start1).Seconds())
	// 3. match abs path only

}

//
//func execFindCommand(name string) {
//	//findCmd := fmt.Sprintf("/ -name %s 2>/dev/null", name)
//	cmd := exec.Command("find", startDir, "-name", name)
//
//	var output bytes.Buffer
//	cmd.Stdout = &output
//	e := cmd.Run()
//	if e != nil {
//		log.Println("run error :" + e.Error())
//	} else if output.String() != "" {
//		fmt.Printf("\t%s - %s\n", name, output.String())
//	}
//}
