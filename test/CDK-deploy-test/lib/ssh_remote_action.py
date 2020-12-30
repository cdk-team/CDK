import os
from fabric import Connection
from lib.conf import CDK, SERVER
from invoke.exceptions import UnexpectedExit


def get_remote_conn():
    connect_kwargs = {'key_filename': SERVER.KEY_PATH}
    return Connection(SERVER.HOST, SERVER.USER, connect_kwargs=connect_kwargs)


conn = get_remote_conn()


def output_err(env, cmd, pattern, type):
    print('[ERROR] {} while running cmd: {}\nexcepted {} pattern:\n{}'.format(env, cmd, type, pattern))


def test():
    result = conn.run("uname -s", hide=True)
    msg = "Ran {0.command!r} on {0.connection.host}, got stdout:\n{0.stdout}"
    print(msg.format(result))


def update_remote_bin():
    print('[compile and upload]')
    try:
        os.system(CDK.BUILD_CMD)
        conn.put(CDK.BIN_PATH, CDK.REMOTE_HOST_PATH)
        conn.run('chmod a+x {}'.format(CDK.REMOTE_HOST_PATH))
    except Exception as e:
        print('errors while compile and update.')
        print(e)
        exit(1)


def check_host_exec(cmd_parsed, white_list, black_list, verbose=False):
    print('[TEST] [{}] {}'.format('ECS', cmd_parsed))
    try:
        result = conn.run(cmd_parsed, hide=bool(1 - verbose))
        for pattern in white_list:
            if pattern not in result.stdout + result.stderr:
                output_err('ECS', cmd_parsed, pattern, 'white')

        for pattern in black_list:
            if pattern in result.stdout + result.stderr:
                output_err('ECS', cmd_parsed, pattern, 'black')

    except UnexpectedExit as e:
        pass
        # print('invoke UnexpectedExit')
        # print(e)


def check_host_evaluate(cmd, white_list, black_list, verbose=False):
    cmd_parsed = "{} {}".format(CDK.REMOTE_HOST_PATH, cmd)
    print('[TEST] [{}] {}'.format('ECS', cmd_parsed))

    try:
        result = conn.run(cmd_parsed, hide=bool(1 - verbose))
        for pattern in white_list:
            if pattern not in result.stdout + result.stderr:
                output_err('ECS', cmd_parsed, pattern, 'white')

        for pattern in black_list:
            if pattern in result.stdout + result.stderr:
                output_err('ECS', cmd_parsed, pattern, 'black')

    except UnexpectedExit as e:
        pass
        # print('invoke UnexpectedExit')
        # print(e)


def inside_container_cmd(image, docker_args, cmd, white_list, black_list, verbose=False):
    # docker run -v /root/cdk_linux_amd64:/cdk_linux_amd64 --rm --net=host ubuntu /bin/bash -c "/cdk_linux_amd64 cmd"
    success = True

    cmd_parsed = "docker run -v {}:{} --rm {} {} /bin/sh -c \"{} {}\"".format(
        CDK.REMOTE_HOST_PATH,
        CDK.REMOTE_CONTAINER_PATH,
        docker_args,
        image,
        CDK.REMOTE_CONTAINER_PATH,
        cmd
    )
    print('[TEST] [{}] {}'.format(image, cmd_parsed))

    try:
        result = conn.run(cmd_parsed, hide=bool(1 - verbose))
        for pattern in white_list:
            if pattern not in result.stdout + result.stderr:
                output_err(image, cmd_parsed, pattern, 'white')
                success = False
        for pattern in black_list:
            if pattern in result.stdout + result.stderr:
                output_err(image, cmd_parsed, pattern, 'black')
    except UnexpectedExit as e:
        pass
        # print('invoke UnexpectedExit')
        # print(e)
        # return
