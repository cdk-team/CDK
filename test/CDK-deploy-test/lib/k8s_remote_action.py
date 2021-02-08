import subprocess
import time
from lib.ssh_remote_action import output_err
from lib.conf import CDK, K8S, DEV_PATH


def k8s_pod_upload():
    print('[upload] CDK binary to K8s pod:{}'.format(K8S.TARGET_POD))
    cmd = r'{} --kubeconfig={} cp {} {}:{}'.format(DEV_PATH.KUBECTL_PATH, K8S.KUBE_CONFIG, CDK.BIN_PATH, K8S.TARGET_POD,
                                                   K8S.REMOTE_POD_PATH)
    ret = subprocess.Popen(cmd, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)

    time.sleep(1)
    cmd1 = r'{} --kubeconfig={} exec {} -- ls {}'.format(DEV_PATH.KUBECTL_PATH, K8S.KUBE_CONFIG, K8S.TARGET_POD,
                                                         K8S.REMOTE_POD_PATH)

    ret1 = subprocess.Popen(cmd1, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    if K8S.REMOTE_POD_PATH in str(ret1.stdout.read()):
        return
    else:
        print(str(ret1.stdout.read()))
        print(str(ret1.stderr.read()))
        raise Exception("Upload cdk binary to K8s failed.\nCMD: " + cmd)


def check_pod_exec(cmd, white_list, black_list, verbose=False):
    # OCI runtime exec failed: exec failed: container_linux.go:344: starting container process caused "text file busy"
    time.sleep(1)

    cmd_parsed = r'{} --kubeconfig={} exec {} -- {} {}'.format(DEV_PATH.KUBECTL_PATH, K8S.KUBE_CONFIG, K8S.TARGET_POD,
                                                               K8S.REMOTE_POD_PATH, cmd)
    print('[TEST] [{}] {}'.format('K8s Pod', cmd_parsed))

    ret = subprocess.Popen(cmd_parsed, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                           shell=True)

    stdout = str(ret.stdout.read())
    stderr = str(ret.stderr.read())

    if verbose:
        print('stdout\n', stdout)
        print('stderr\n', stderr)

    for pattern in white_list:
        if pattern not in stdout + stderr:
            output_err('K8s Pod', cmd_parsed, pattern, 'white')

    for pattern in black_list:
        if pattern in stdout + stderr:
            output_err('K8s Pod', cmd_parsed, pattern, 'black')


if __name__ == '__main__':
    k8s_pod_upload()
