package evaluate

import (
	"fmt"
	"github.com/Xyntax/CDK/pkg/kubectl"
	"log"
	"strings"
)

func CheckK8sAnonymousLogin() bool {

	// check if api-server allows system:anonymous request
	log.Println("checking if api-server allows system:anonymous request.")

	resp := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: "",
			Server:    "", // default
			Api:       "/",
			Method:    "get",
			Args:      "",
			Anonymous: true,
		})

	if strings.Contains(resp, "/api") {
		fmt.Println("\tcongrats, api-server allows anonymous request.")
		log.Println("trying to list namespaces")

		// check if system:anonymous can list namespaces
		resp := kubectl.ServerAccountRequest(
			kubectl.K8sRequestOption{
				TokenPath: "",
				Server:    "", // default
				Api:       "/api/v1/namespaces",
				Method:    "get",
				Args:      "",
				Anonymous: true,
			})
		if len(resp) > 0 && strings.Contains(resp, "kube-system") {
			fmt.Println("\tsuccess, the system:anonymous role have a high authority.")
			fmt.Println("\tnow you can make your own request to takeover the entire k8s cluster with `./cdk kcurl` command\n\tgood luck and have fun.")
			return true
		} else {
			fmt.Println("\tfailed.")
			fmt.Println("\tresponse:" + resp)
			return true
		}
	} else {
		fmt.Println("\tapi-server forbids anonymous request.")
		fmt.Println("\tresponse:" + resp)
		return false
	}
}
