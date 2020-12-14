package lib

import (
	"fmt"
	"github.com/Xyntax/CDK/conf"
	"github.com/Xyntax/CDK/pkg/evaluate"
	"github.com/Xyntax/CDK/pkg/kubectl"
	"github.com/Xyntax/CDK/pkg/netcat"
	"github.com/Xyntax/CDK/pkg/network"
	"github.com/Xyntax/CDK/pkg/probe"
	"github.com/Xyntax/CDK/pkg/ps"
	"github.com/Xyntax/CDK/pkg/util"
	"github.com/Xyntax/CDK/pkg/vi"
	"github.com/docopt/docopt-go"
	"log"
	"os"
	"strconv"
)

func PassInnerArgs() {
	os.Args = os.Args[1:]
}

func ParseCDKMain() {

	if len(os.Args) == 1 {
		docopt.PrintHelpAndExit(nil, BannerContainer)
	}

	// nc needs -v and -h , parse it outside
	if os.Args[1] == "nc" {
		// https://github.com/jiguangin/netcat
		PassInnerArgs()
		netcat.RunVendorNetcat()
		return
	}


	if Args["evaluate"].(bool) {

		fmt.Printf("\n[Information Gathering - System Info]\n")
		evaluate.BasicSysInfo()

		fmt.Printf("\n[Information Gathering - Services]\n")
		evaluate.SearchSensitiveEnv()
		evaluate.SearchSensitiveService()

		fmt.Printf("\n[Information Gathering - Commands and Capabilities]\n")
		evaluate.SearchAvailableCommands()
		evaluate.GetProcCapabilities()

		fmt.Printf("\n[Information Gathering - Mounts]\n")
		evaluate.MountEscape()

		fmt.Printf("\n[Information Gathering - Net Namespace]\n")
		evaluate.CheckNetNamespace()

		fmt.Printf("\n[Discovery - K8s API Server]\n")
		evaluate.CheckK8sAnonymousLogin()

		fmt.Printf("\n[Discovery - K8s Service Account]\n")
		evaluate.CheckK8sServiceAccount(conf.K8sSATokenDefaultPath)

		fmt.Printf("\n[Discovery - Cloud Provider Metadata API]\n")
		evaluate.CheckCloudMetadataAPI()

		if Args["--full"].(bool) {
			fmt.Printf("\n[Information Gathering - Sensitive Files]\n")
			evaluate.SearchLocalFilePath()
		}
		return
	}

	if Args["run"].(bool) {
		if Args["--list"].(bool) {
			ListAllPlugin()
			os.Exit(0)
		}
		name := Args["<exploit>"].(string)
		if Plugins[name] == nil {
			fmt.Printf("\nInvalid script name: %s , available scripts:\n", name)
			ListAllPlugin()
			return
		}
		RunSinglePlugin(name)
		return
	}

	if Args["<tool>"] != nil {
		args := Args["<args>"].([]string)

		switch Args["<tool>"] {
		case "vi":
			PassInnerArgs()
			vi.RunVendorVi()
		case "kcurl":
			kubectl.KubectlMain(args)
		case "ucurl":
			if len(args) != 4 {
				log.Fatal("invalid input args, Example: ./cdk ucurl get /var/run/docker.sock http://127.0.0.1/info \"\"")
			}
			ans := util.UnixHttpSend(args[0], args[1], args[2], args[3])
			log.Println("response:")
			fmt.Println(ans)
		case "ifconfig":
			network.GetLocalAddresses()
		case "ps":
			ps.RunPs()
		case "probe":
			if len(args) != 4 {
				log.Println("Invalid input args.")
				log.Println("usage: cdk probe <ip> <port> <parallels> <timeout-ms>")
				log.Fatal("example: cdk probe 192.168.1.0-255 22,80,100-110 50 1000")
			}
			parallel, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				log.Println("err found when parse input arg <parallel>")
				log.Fatal(err)
			}
			timeout, err := strconv.Atoi(args[3])
			if err != nil {
				log.Println("err found when parse input arg <timeout-ms>")
				log.Fatal(err)
			}
			probe.TCPScanToolAPI(args[0], args[1], parallel, timeout)
		default:
			docopt.PrintHelpAndExit(nil, BannerContainer)
		}
	}
}
