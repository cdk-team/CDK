
# CDK - Zero Dependency Container Penetration Toolkit

English | [简体中文](https://github.com/cdk-team/CDK/wiki/CDK-Home-CN)

![png](https://static.cdxy.me/20201203170308_NwzGiT_Screenshot.jpeg)

## Legal Disclaimer

Usage of CDK for attacking targets without prior mutual consent is illegal.
CDK is for security testing purposes only.

## Overview

CDK is an open-sourced container penetration toolkit, designed for offering stable exploitation in different slimmed containers without any OS dependency. It comes with useful net-tools and many powerful PoCs/EXPs helps you to escape container and takeover K8s cluster easily.

Currently still under development, submit [issues](https://github.com/cdk-team/CDK/issues) or mail <i@cdxy.me> if you need any help. 

## Installation

Download latest release in: https://github.com/cdk-team/CDK/releases/

Drop executable files into target container and start testing.

## Usage
```
Usage:
  cdk evaluate [--full]
  cdk run (--list | <exploit> [<args>...])
  cdk auto-escape <cmd>
  cdk <tool> [<args>...]

Evaluate:
  cdk evaluate                              Gather information to find weakness inside container.
  cdk evaluate --full                       Enable file scan during information gathering.

Exploit:
  cdk run --list                            List all available exploits.
  cdk run <exploit> [<args>...]             Run single exploit, docs in https://github.com/cdk-team/CDK/wiki

Auto Escape:
  cdk auto-escape <cmd>                     Escape container in different ways then let target execute <cmd>.

Tool:
  vi <file>                                 Edit files in container like "vi" command.
  ps                                        Show process information like "ps -ef" command.
  nc [options]                              Create TCP tunnel.
  ifconfig                                  Show network information.
  kcurl <path> (get|post) <uri> <data>      Make request to K8s api-server.
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
|Information Gathering|OS Basic Info|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-System-Info)|
|Information Gathering|Available Capabilities|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-Commands-and-Capabilities)|
|Information Gathering|Available Linux Commands|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-Commands-and-Capabilities)|
|Information Gathering|Mounts|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-Mounts)|
|Information Gathering|Net Namespace|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-Net-Namespace)|
|Information Gathering|Sensitive ENV|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-Services)|
|Information Gathering|Sensitive Process|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-Services)|
|Information Gathering|Sensitive Local Files|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-Sensitive-Files)|
|Information Gathering|Kube-proxy Route Localnet(CVE-2020-8558)|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-check-net.ipv4.conf.all.route_localnet)|
|Discovery|K8s Api-server Info|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-K8s-API-Server)|
|Discovery|K8s Service-account Info|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-K8s-Service-Account)|
|Discovery|Cloud Provider Metadata API|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-Cloud-Provider-Metadata-API)|

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
|Escaping|containerd-shim CVE-2020-15257|shim-pwn|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-shim-pwn)|
|Escaping|docker.sock PoC (DIND attack)|docker-sock-check|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-docker-sock-check)|
|Escaping|docker.sock Backdoor Image Deploy|docker-sock-deploy|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-docker-sock-deploy)|
|Escaping|Device Mount Escaping|mount-disk|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-mount-disk)|
|Escaping|Cgroups Escaping|mount-cgroup|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-mount-cgroup)|
|Escaping|Procfs Escaping|mount-procfs|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-mount-procfs)|
|Escaping|Ptrace Escaping PoC|check-ptrace|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-check-ptrace)|
|Escaping|Exploit lxcfs|lxcfs-rw|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-lxcfs-rw)|
|Escaping|Rewrite Cgroup(devices.allow)|rewrite-cgroup-devices|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-rewrite-cgroup-devices)|
|Discovery|K8s Component Probe|service-probe|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-service-probe)|
|Discovery|Dump Istio Sidecar Meta|istio-check|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-check-istio)|
|Remote Control|Reverse Shell|reverse-shell|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-reverse-shell)|
|Credential Access|Access Key Scanning|ak-leakage|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-ak-leakage)|
|Credential Access|Dump K8s Secrets|k8s-secret-dump|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-secret-dump)|
|Credential Access|Dump K8s Config|k8s-configmap-dump|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-configmap-dump)|
|Persistence|Deploy WebShell|webshell-deploy|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-webshell-deploy)|
|Persistence|Deploy Backdoor Pod|k8s-backdoor-daemonset|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-backdoor-daemonset)|
|Persistence|Deploy Shadow K8s api-server|k8s-shadow-apiserver|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-shadow-apiserver)|
|Persistence|K8s MITM Attack (CVE-2020-8554)|k8s-mitm-clusterip|✔|[link](https://github.com/cdk-team/CDK/wiki/Evaluate:-k8s-mitm-clusterip)|
|Persistence|Deploy K8s CronJob|k8s-cronjob|✔|[link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-cronjob)|

### Tool Module

Running commands like in Linux, little different in input-args, see the usage link.
```
cdk nc [options]
cdk ps
```

|Command|Description|Supported|Usage/Example|
|---|---|---|---|
|nc|TCP Tunnel|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-nc)|
|ps|Process Information|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-ps)|
|ifconfig|Network Information|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-ifconfig)|
|vi|Edit Files|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-vi)|
|kcurl|Request to K8s api-server|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-kcurl)|
|dcurl|Request to Docker HTTP API|||
|ucurl|Request to Docker Unix Socket|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-ucurl)|
|rcurl|Request to Docker Registry API|||
|probe|IP/Port Scanning|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-probe)|

## Developer Docs

* [run test in container.](https://github.com/cdk-team/CDK/wiki/Run-Test)

## Contributing to CDK
  
First off, thanks for taking the time to contribute!   
  
By reporting any issue, ideas or PRs, your GitHub ID will be listed here.

* https://github.com/cdk-team/CDK/blob/main/thanks.md

#### Bug Reporting

Bugs are tracked as [GitHub Issues](https://github.com/cdk-team/CDK/issues). Create an issue with the current CDK version, error msg and the environment. Describe the exact steps which reproduce the problem.

#### Suggesting Enhancements

Enhancement suggestions are tracked as [GitHub Discussions](https://github.com/cdk-team/CDK/discussions). You can publish any thoughts here to discuss with developers directly.

#### Pull Requests

Fix problems or maintain CDK's quality:

* Describe the current CDK version, environment, problem and exact steps that reproduce the problem.
* Running screenshots or logs before and after you fix the problem.

New feature or exploits:

* Explain why this enhancement would be useful to other users.
* Please enable a sustainable environment for us to review contributions.
* Screenshots about how this new feature works. 
* If you are committing a new evaluate/exploit scripts, please add a simple doc to your PR message, here is an [example](https://github.com/cdk-team/CDK/wiki/Exploit:-docker-sock-deploy). 



## TODO

1. Echo loader for delivering CDK into target container via Web RCE. 
2. EDR defense evasion.
3. Compile optimization.
4. Dev docs
