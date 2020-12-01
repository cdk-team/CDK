package evaluate

import (
	"fmt"
	"github.com/Xyntax/CDK/conf"
	"github.com/idoubi/goz"
	"log"
	"strings"
)

func CheckCloudMetadataAPI() {
	for _, apiInstance := range conf.CloudAPI {
		cli := goz.NewClient(goz.Options{
			Timeout: 1,
		})
		resp, err := cli.Get(apiInstance.API)
		if err != nil {
			log.Printf("failed to dial %s API.", apiInstance.CloudProvider)
			continue
		}
		r, _ := resp.GetBody()
		if strings.Contains(r.String(), apiInstance.ResponseMatch) {
			fmt.Printf("\t%s Metadata API available in %s\n", apiInstance.CloudProvider, apiInstance.API)
			fmt.Printf("\tDocs: %s\n", apiInstance.DocURL)
		} else {
			log.Printf("failed to dial %s API.", apiInstance.CloudProvider)
		}
	}
}
