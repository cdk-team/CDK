

class SERVER: # your remote server for test
    HOST = '39.104.80.49'
    USER = 'root'
    KEY_PATH = '/Users/xy/.ssh/lezhen-cdk.pem'

class CDK:
    # local source-code dir to run `go build`
    BUILD_PATH = '/usr/local/go/bin/src/github.com/Xyntax/CDK/cmd/cdk'
    # build command
    BUILD_CMD = 'cd {} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cdk.go'.format(BUILD_PATH)
    # binary after build
    BIN_PATH = '/usr/local/go/bin/src/github.com/Xyntax/CDK/cmd/cdk/cdk'

    # you can keep it unchanged
    REMOTE_HOST_PATH = '/root/cdk-fabric'
    REMOTE_CONTAINER_PATH = '/cdk-fabric'