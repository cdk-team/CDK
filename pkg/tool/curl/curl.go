package curl

import (
	"fmt"
	"github.com/idoubi/goz"
	"log"
)


func RunCurlCmd(args []string){
	cli := goz.NewClient()

	resp, err := cli.Get("http://127.0.0.1:8091/get")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%T", resp)

}