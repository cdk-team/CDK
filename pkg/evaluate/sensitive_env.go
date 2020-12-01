package evaluate

import (
	"github.com/Xyntax/CDK/conf"
	"log"
	"os"
	"regexp"
)

func SearchSensitiveEnv() {
	for _, env := range os.Environ() {
		ans, err := regexp.MatchString(conf.SensitiveEnvRegex, env)
		if err != nil {
			log.Println(err)
		} else if ans {
			log.Printf("sensitive env found:\n\t%s", env)
		}
	}
}
