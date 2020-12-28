package plugin

import "fmt"

type PluginInterface interface {
	Desc() string
	Run() bool
}

type TaskInterface interface {
	Exec() bool
	Desc() string
}

var Exploits map[string]PluginInterface
var Tasks map[string]TaskInterface

func init() {
	Exploits = make(map[string]PluginInterface)
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

func RegisterExploit(name string, plugin PluginInterface) {
	Exploits[name] = plugin
}

func RunSingleTask(name string) {
	Tasks[name].Exec()
}

func RegisterTask(name string, task TaskInterface) {
	Tasks[name] = task
}