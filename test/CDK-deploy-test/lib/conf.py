class SERVER:  # your remote server for test
    HOST = '118.195.140.100'
    USER = 'root'
    PASS = ''
    KEY_PATH = '/Users/xy/.ssh/id_rsa'


class DEV_PATH:
    KUBECTL_PATH = '/Users/xy/Desktop/lezhen-test-case/k8s/kubectl'
    GO_BINARY = '/Users/xy/go/go1.16beta1/bin/go'


class CDK:
    # local source-code dir to run `go build`
    BUILD_PATH = '/Users/xy/go/CDK/cmd/cdk'
    # build command
    BUILD_CMD = 'cd {} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 {} build cdk.go'.format(BUILD_PATH, DEV_PATH.GO_BINARY)
    # binary after build
    BIN_PATH = '/Users/xy/go/CDK/cmd/cdk/cdk'
    # you can keep it unchanged
    REMOTE_HOST_PATH = '/root/cdk-fabric'
    REMOTE_CONTAINER_PATH = '/cdk-fabric'


class K8S:
    KUBE_CONFIG = '/Users/xy/.kube/cdk.config'
    # KUBE_CONFIG = '/Users/xy/.kube/config'
    # upload cdk to target pod then check command output using kubectl
    TARGET_POD = 'myappnew'
    # you can keep it unchanged
    REMOTE_POD_PATH = '/cdk-fabric'


class SELFBUILD_K8S:
    # Master node SSH
    HOST = '118.195.140.100'
    USER = 'root'
    PASS = ''
    KEY_PATH = '/Users/xy/.ssh/id_rsa'
    REMOTE_HOST_PATH = '/root/cdk-fabric'
    # upload cdk to target pod then check command output using kubectl
    TARGET_POD = 'myappnew'
    # you can keep it unchanged
    REMOTE_POD_PATH = '/cdk-fabric'
    KUBERNETES_SERVICE_PORT = '443'
    KUBERNETES_SERVICE_HOST = '172.16.252.1'
