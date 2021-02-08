import os
from lib.conf import CDK

version = 'cdk_v0.10'

cmd = '''
cd {}; 
rm ../../cdk_release_binary/cdk_* 2>&1; 
gox -os "linux darwin" -arch "386 amd64 arm arm64 mips mips64 mips64le mipsle" 2>&1; 
mv cdk_* ../../cdk_release_binary/;
cd ../../cdk_release_binary/ && tar -zcvf {}_release.tar.gz cdk_*;
'''.strip().format(CDK.BUILD_PATH,version)


def gox_release():
    print("check cdk version")
    print("check cdk banner in .go")
    print("check cdk banner in readme.md")
    print("check cdk banner in github wiki")
    os.system(cmd)


if __name__ == '__main__':
    gox_release()
