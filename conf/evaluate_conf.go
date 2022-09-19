
/*
Copyright 2022 The Authors of https://github.com/CDK-TEAM/CDK .

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conf

// check useful linux commands in container
var LinuxCommandChecklist = []string{
	"curl",
	"wget",
	"nc",
	"netcat",
	"kubectl",
	"docker",
	"find",
	"ps",
	"java",
	"python",
	"python3",
	"php",
	"node",
	"npm",
	"apt",
	"yum",
	"dpkg",
	"nginx",
	"httpd",
	"apache",
	"apache2",
	"ssh",
	"mysql",
	"mysql-client",
	"git",
	"svn",
	"vi",
	"capsh",
	"mount",
	"fdisk",
	"gcc",
	"g++",
	"make",
	"base64",
	"python2",
	"python2.7",
	"perl",
	"xterm",
	"sudo",
	"ruby",
}

var DefaultPathEnv = []string{
	"/usr/local/sbin",
	"/usr/local/bin",
	"/usr/sbin",
	"/usr/bin",
	"/sbin",
	"/bin",
	"/usr/games",
	"/usr/local/games",
	"/snap/bin",
}

// match ENV to find useful service
var SensitiveEnvRegex = "(?i)\\bssh_|k8s|kubernetes|docker|gopath"

// match process name to find useful service
var SensitiveProcessRegex = "(?i)ssh|ftp|http|tomcat|nginx|engine|php|java|python|perl|ruby|kube|docker|\\bgo\\b"

// match local file path to find sensitive file
// walk starts from StartDir and match substring(AbsFilePath,<names in NameList>)
type sensitiveFileRules struct {
	StartDir string
	NameList []string
}

var SensitiveFileConf = sensitiveFileRules{
	StartDir: "/",
	NameList: []string{
		`/docker.sock`,     // docker socket (http)
		`/containerd.sock`, // containerd socket (grpc)
		`/containerd/s/`,   // containerd-shim socket (grpc)
		`.kube/`,
		`.git/`,
		`.svn/`,
		`.pip/`,
		`/.bash_history`,
		`/.bash_profile`,
		`/.bashrc`,
		`/.ssh/`,
		`.token`,
		`/serviceaccount`,
		`.dockerenv`,
		`/config.json`,
	},
}

// Check cloud provider APIs in evaluate task
type cloudAPIS struct {
	CloudProvider string
	API           string
	ResponseMatch string
	DocURL        string
}

var CloudAPI = []cloudAPIS{
	{
		CloudProvider: "Alibaba Cloud",
		API:           "http://100.100.100.200/latest/meta-data/",
		ResponseMatch: "instance-id",
		DocURL:        "https://help.aliyun.com/knowledge_detail/49122.html",
	},
	{
		CloudProvider: "Azure",
		API:           "http://169.254.169.254/metadata/instance",
		ResponseMatch: "azEnvironment",
		DocURL:        "https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service",
	},
	{
		CloudProvider: "Google Cloud",
		API:           "http://metadata.google.internal/computeMetadata/v1/instance/disks/?recursive=true",
		ResponseMatch: "deviceName",
		DocURL:        "https://cloud.google.com/compute/docs/storing-retrieving-metadata",
	},
	{
		CloudProvider: "Tencent Cloud",
		API:           "http://metadata.tencentyun.com/latest/meta-data/",
		ResponseMatch: "instance-name",
		DocURL:        "https://cloud.tencent.com/document/product/213/4934",
	},
	{
		CloudProvider: "OpenStack",
		API:           "http://169.254.169.254/openstack/latest/meta_data.json",
		ResponseMatch: "availability_zone",
		DocURL:        "https://docs.openstack.org/nova/rocky/user/metadata-service.html",
	},
	{
		CloudProvider: "Amazon Web Services (AWS)",
		API:           "http://169.254.169.254/latest/meta-data/",
		ResponseMatch: "instance-id",
		DocURL:        "https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-retrieval.html",
	},
	{
		CloudProvider: "ucloud",
		API:           "http://100.80.80.80/meta-data/latest/uhost/",
		ResponseMatch: "uhost-id",
		DocURL:        "https://docs.ucloud.cn/uhost/guide/metadata/metadata-server",
	},
}
