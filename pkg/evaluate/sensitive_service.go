package evaluate

import (
	"github.com/cdk-team/CDK/conf"
	gops "github.com/mitchellh/go-ps"
	"log"
	"regexp"
)

func SearchSensitiveService() {
	processList, err := gops.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
	}
	for _, proc := range processList {
		ans, err := regexp.MatchString(conf.SensitiveProcessRegex, proc.Executable())
		if err != nil {
			log.Println(err)
		} else if ans {
			log.Printf("service found in process:\n\t%d\t%d\t%s\n", proc.Pid(), proc.PPid(), proc.Executable())
		}
	}
}
