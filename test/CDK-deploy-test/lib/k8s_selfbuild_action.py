import time
from fabric import Connection
from invoke import UnexpectedExit
from lib.conf import CDK, SELFBUILD_K8S


def get_remote_conn():
    connect_kwargs = {'key_filename': SELFBUILD_K8S.KEY_PATH}
    return Connection(SELFBUILD_K8S.HOST, SELFBUILD_K8S.USER, connect_kwargs=connect_kwargs)


conn = get_remote_conn()


def output_err(env, cmd, pattern, type):
    print('[ERROR] {} while running cmd: {}\nexcepted {} pattern:\n{}'.format(env, cmd, type, pattern))


def update_remote_bin():
    print('[upload] CDK binary to self-build k8s master node')
    try:
        conn.put(CDK.BIN_PATH, SELFBUILD_K8S.REMOTE_HOST_PATH)
        conn.run('chmod a+x {}'.format(SELFBUILD_K8S.REMOTE_HOST_PATH))
    except Exception as e:
        print('errors while update cdk binary.')
        print(e)
        exit(1)


def k8s_master_ssh_cmd(cmd_parsed, white_list, black_list, verbose=False):
    print('[TEST] [{}] {}'.format('Selfbuild k8s master node', cmd_parsed))

    try:
        result = conn.run(cmd_parsed, hide=bool(1 - verbose))
        for pattern in white_list:
            if pattern not in result.stdout + result.stderr:
                output_err('Selfbuild K8s Master Node', cmd_parsed, pattern, 'white')

        for pattern in black_list:
            if pattern in result.stdout + result.stderr:
                output_err('Selfbuild K8s Master Node', cmd_parsed, pattern, 'black')
        return result.stdout + result.stderr
    except UnexpectedExit as e:
        print('invoke UnexpectedExit')
        print(e)


def selfbuild_k8s_pod_upload():
    # upload cdk to master node via ssh
    update_remote_bin()

    # cp cdk from master node to target pod via (kubectl in master node).
    cmd = r'kubectl cp {} {}:{}'.format(SELFBUILD_K8S.REMOTE_HOST_PATH, SELFBUILD_K8S.TARGET_POD,
                                        SELFBUILD_K8S.REMOTE_POD_PATH)
    k8s_master_ssh_cmd(cmd, [], [], True)

    time.sleep(1)
    # check if upload success
    cmd1 = r'kubectl exec {} ls {}'.format(SELFBUILD_K8S.TARGET_POD, SELFBUILD_K8S.REMOTE_POD_PATH)
    resp = k8s_master_ssh_cmd(cmd1, [], [], False)
    if SELFBUILD_K8S.REMOTE_POD_PATH in str(resp):
        return
    else:
        raise Exception("Upload cdk binary to self-build k8s failed.\nCMD: " + cmd)


def check_selfbuild_k8s_pod_exec(cmd, white_list, black_list, verbose=False):
    # OCI runtime exec failed: exec failed: container_linux.go:344: starting container process caused "text file busy"
    time.sleep(1)
    cmd_parsed = r'kubectl exec {} -- {} {}'.format(
        SELFBUILD_K8S.TARGET_POD,
        SELFBUILD_K8S.REMOTE_POD_PATH,
        cmd
    )
    # print('[TEST] [{}] {}'.format('Selfbuild K8s Pod', cmd_parsed))
    k8s_master_ssh_cmd(cmd_parsed, white_list, black_list, verbose)
