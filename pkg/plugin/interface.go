package plugin

import (
	"fmt"
	"os"
	"text/tabwriter"
)

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
	writer := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	for name, plugin := range Exploits {
		str := fmt.Sprintf("%s\t %s", name, plugin.Desc())
		fmt.Fprintln(writer, str)
	}

	writer.Flush()
}

func RunSingleExploit(name string) {
	Exploits[name].Run()
}

func RegisterExploit(name string, exploit ExploitInterface) {
	Exploits[name] = exploit
}

func RunSingleTask(name string) {
	// fmt.Printf("[+] Running exploit: %s.\n", name)
	// fmt.Printf("[+] %s\n", Tasks[name].Desc())
	// Can not call cli.Args here, because it will cause "import cycle".
	// fmt.Printf("[+] Args: %v.\n", cli.Args["<args>"])
	Tasks[name].Exec()
}

func RegisterTask(name string, task TaskInterface) {
	Tasks[name] = task
}
