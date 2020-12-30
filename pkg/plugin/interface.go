package plugin

import "fmt"

type ExploitInterface interface {
	Desc() string
	Run() bool
}

type TaskInterface interface {
	Exec() bool
	Desc() string
}

var Exploits map[string]ExploitInterface
var Tasks map[string]TaskInterface

func init() {
	Exploits = make(map[string]ExploitInterface)
	Tasks = make(map[string]TaskInterface)
}


func ListAllExploit() {
	for name, plugin := range Exploits {
		fmt.Println(name, "\t", plugin.Desc())
	}
}

func RunSingleExploit(name string) {
	Exploits[name].Run()
}

func RegisterExploit(name string, exploit ExploitInterface) {
	Exploits[name] = exploit
}

func RunSingleTask(name string) {
	Tasks[name].Exec()
}

func RegisterTask(name string, task TaskInterface) {
	Tasks[name] = task
}