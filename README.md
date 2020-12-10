
# CDK - Zero Dependency Container Penetration Toolkit

English | [简体中文](https://github.com/Xyntax/CDK/wiki/CDK-Home-CN)

![png](https://static.cdxy.me/20201203170308_NwzGiT_Screenshot.jpeg)

## Legal Disclaimer

Usage of CDK for attacking targets without prior mutual consent is illegal.
CDK is for security testing purposes only.

## Overview

CDK is an open-sourced container penetration toolkit, designed for offering stable exploitation in different slimmed containers without any OS dependency. It comes with useful net-tools and many powerful PoCs/EXPs helps you to escape container and takeover K8s cluster easily.

Currently still under development, submit [issues](https://github.com/Xyntax/CDK/issues) or mail <i@cdxy.me> if you need any help. 

## Installation

Download latest release in: https://github.com/Xyntax/CDK/releases/

Drop executable files into target container and start testing.

## Usage
```
Container DucK
Zero-dependency docker/k8s penetration toolkit by <i@cdxy.me>
Find tutorial, configuration and use-case in https://github.com/Xyntax/CDK/wiki

Usage:
  cdk evaluate [--full]
  cdk run (--list | <exploit> [<args>...])
  cdk <tool> [<args>...]

Evaluate:
  cdk evaluate                              Gather information to find weekness inside container.
  cdk evaluate --full                       Enable file scan during information gathering.

Exploit:
  cdk run --list                            List all available exploits.
  cdk run <exploit> [<args>...]             Run single exploit, docs in https://github.com/Xyntax/CDK/wiki

Tool:
  vi <file>                                 Edit files in container like "vi" command.
  ps                                        Show process information like "ps -ef" command.
  nc [options]                              Create TCP tunnel.
  ifconfig                                  Show network information.
  kcurl	(get|post) <url> <data>             Make request to K8s api-server.
  ucurl (get|post) <socket> <uri> <data>    Make request to docker unix socket.
  probe <ip> <port> <parallel> <timeout-ms> TCP port scan, example: cdk probe 10.0.1.0-255 80,8080-9443 50 1000

Options:
  -h --help     Show this help msg.
  -v --version  Show version.
```

## Features

CDK have three modules:

1. Evaluate: gather information inside container to find potential weakness. 
2. Exploit: for container escaping, persistance and lateral movement
3. Tool: network-tools and APIs for TCP/HTTP requests, tunnels and K8s cluster management.

### Evaluate Module

Usage 
```
cdk evaluate [--full]
```
This command will run the scripts below without local file scanning, using `--full` to enable all.

|Tactics|Script|Supported|Usage/Example|
|---|---|---|---|
|Information Gathering|OS Basic Info|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-System-Info)|
|Information Gathering|Available Capabilities|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-Commands-and-Capabilities)|
|Information Gathering|Available Linux Commands|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-Commands-and-Capabilities)|
|Information Gathering|Mounts|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-Mounts)|
|Information Gathering|Net Namespace|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-Net-Namespace)|
|Information Gathering|Sensitive ENV|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-Services)|
|Information Gathering|Sensitive Process|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-Services)|
|Information Gathering|Sensitive Local Files|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-Sensitive-Files)|
|Discovery|K8s Api-server Info|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-K8s-API-Server)|
|Discovery|K8s Service-account Info|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-K8s-Service-Account)|
|Discovery|Cloud Provider Metadata API|✔|[link](https://github.com/Xyntax/CDK/wiki/Evaluate:-Cloud-Provider-Metadata-API)|

### Exploit Module

List all available exploits:
```
cdk run --list
```

Run targeted exploit:
```
cdk run <script-name> [options]
```

|Tactic|Technique|CDK Exploit Name|Supported|Doc|
|---|---|---|---|---|
|Escaping|docker-runc CVE-2019-5736|runc-pwn|✔||
|Escaping|docker-cp CVE-2019-14271||||
|Escaping|containerd-shim CVE-2020-15257||||
|Escaping|dirtycow CVE-2016-5159||||
|Escaping|docker.sock PoC (DIND attack)|docker-sock-check|✔|[link](https://github.com/Xyntax/CDK/wiki/Exploit:-docker-sock-check)|
|Escaping|docker.sock Backdoor Image Deploy|docker-sock-deploy|✔|[link](https://github.com/Xyntax/CDK/wiki/Exploit:-docker-sock-deploy)|
|Escaping|Device Mount Escaping|mount-disk|✔|[link](https://github.com/Xyntax/CDK/wiki/Exploit:-mount-disk)|
|Escaping|Cgroups Escaping|mount-cgroup|✔|[link](https://github.com/Xyntax/CDK/wiki/Exploit:-mount-cgroup)|
|Escaping|Procfs Escaping|mount-procfs|✔|[link](https://github.com/Xyntax/CDK/wiki/Exploit:-mount-procfs)|
|Escaping|Ptrace Escaping PoC|check-ptrace|✔|[link](https://github.com/Xyntax/CDK/wiki/Exploit:-check-ptrace)|
|Discovery|K8s Component Probe|service-probe|✔|[link](https://github.com/Xyntax/CDK/wiki/Exploit:-service-probe)|
|Lateral Movement|K8s Service Account Control||||
|Lateral Movement|Attack K8s api-server||||
|Lateral Movement|Attack K8s Kubelet||||
|Lateral Movement|Attack K8s Dashboard||||
|Lateral Movement|Attack K8s Helm||||
|Lateral Movement|Attack K8s Etcd||||
|Lateral Movement|Attack Private Docker Registry||||
|Remote Control|Reverse Shell|reverse-shell|✔|[link](https://github.com/Xyntax/CDK/wiki/Exploit:-reverse-shell)|
|Credential Access|Access Key Scanning|ak-leakage|✔|[link](https://github.com/Xyntax/CDK/wiki/Exploit:-ak-leakage)|
|Credential Access|Dump K8s Secrets|k8s-secret-dump|✔||
|Credential Access|Dump K8s Config|k8s-configmap-dump|✔||
|Persistence|Deploy WebShell||||
|Persistence|Deploy Backdoor Pod||||
|Persistence|Deploy Shadow K8s api-server||||
|Persistence|Deploy K8s CronJob||||
|Defense Evasion|Disable K8s Audit||||


### Tool Module

Running commands like in Linux, little different in input-args, see the usage link.
```
cdk nc [options]
cdk ps
```

|Command|Description|Supported|Usage/Example|
|---|---|---|---|
|nc|TCP Tunnel|✔|[link](https://github.com/Xyntax/CDK/wiki/Tool:-nc)|
|ps|Process Information|✔|[link](https://github.com/Xyntax/CDK/wiki/Tool:-ps)|
|ifconfig|Network Information|✔|[link](https://github.com/Xyntax/CDK/wiki/Tool:-ifconfig)|
|vi|Edit Files|✔|[link](https://github.com/Xyntax/CDK/wiki/Tool:-vi)|
|kcurl|Request to K8s api-server|✔|[link](https://github.com/Xyntax/CDK/wiki/Tool:-kcurl)|
|dcurl|Request to Docker HTTP API|||
|ucurl|Request to Docker Unix Socket|✔|[link](https://github.com/Xyntax/CDK/wiki/Tool:-ucurl)|
|rcurl|Request to Docker Registry API|||
|probe|IP/Port Scanning|✔|[link](https://github.com/Xyntax/CDK/wiki/Tool:-probe)|

## TODO

1. Echo loader for delivering CDK into target container via Web RCE. 
2. EDR defense evasion.
3. Compile optimization.
4. Dev docs