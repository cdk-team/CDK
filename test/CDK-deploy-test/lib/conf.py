

class SERVER:
    HOST = '39.104.80.49'
    USER = 'root'
    KEY_PATH = '/Users/xy/.ssh/lezhen-cdk.pem'

class CDK:
    BUILD_PATH = '/usr/local/go/bin/src/github.com/Xyntax/CDK/cmd/cdk'
    BUILD_CMD = 'cd {} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cdk.go'.format(BUILD_PATH)
    BIN_PATH = '/usr/local/go/bin/src/github.com/Xyntax/CDK/cmd/cdk/cdk'
    REMOTE_HOST_PATH = '/root/cdk-fabric'
    REMOTE_CONTAINER_PATH = '/cdk-fabric'