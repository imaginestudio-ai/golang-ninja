# Appendix : Building a Testing Environment

Every chapter of this book includes Go code examples to illustrate some points we make in the text. You can find all these Go programs in this book’s GitHub repository (see the _Further reading_ section of this chapter). While you don’t have to execute them all, we believe that manually running the code and observing the result may help reinforce the learned material and explain the finer details.

The first part of this book, _Chapters 1_ to _5_, includes relatively short code examples you can run in the Go Playground (_Further reading_) or on any computer with Go installed. For instructions on how to install Go, you can refer to _Chapter 1_ or follow the official download and installation procedure (_Further reading_).

The rest of the book, starting from [_Chapter 6_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_06.xhtml#_idTextAnchor144), assumes you can interact with a virtual topology, which we run in containers with the help of `containerlab` (_Further reading_). This _Appendix_ documents the process of building a testing environment that includes the compatible version of `containerlab` and other related dependencies, to make sure you get a seamless experience running examples from any chapter of this book.

Just Imagine

# What is a testing environment?

The primary goal is to build an environment with the right set of hardware and software that meets the minimum requirements to execute the code examples. We base the requirements on the assumption that you’re deploying a **virtual machine** (**VM**), as we realize you might not deploy this on a dedicated bare-metal server.

When it comes to deploying a VM for testing (testbed), you have two options, both of which we discuss later:

-   You can deploy this VM in a self-hosted environment, such as VMware or **Kernel-based Virtual** **Machine** (**KVM**).
-   You could use a cloud-hosted environment—for example, **Amazon Web** **Services** (**AWS**).

From the hardware perspective, we assume that the underlying CPU architecture is 64-bit x86, and our recommendation is to give the VM at least 2 vCPUs and 4 GB of RAM and ideally double that to make things a bit faster.

We describe all software provisioning and configuration in an Ansible playbook included in this book’s GitHub repository (_Further reading_). We highly recommend you use the automated approach we have prepared for you to install all the dependencies to run the code examples in the book.

You can still install these packages on top of any Linux distribution—for example, **Windows Subsystem for Linux version 2** (**WSL 2**). In case you want to do the installation manually, we include a full list of dependencies here:

<table id="table001-4" class="No-Table-Style _idGenTablePara-1"><colgroup><col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Package</strong></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Version</strong></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">Go</span></p></td><td class="No-Table-Style"><p><span class="No-Break">1.18.1</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">containerlab</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break">0.25.1</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">Docker</span></p></td><td class="No-Table-Style"><p><span class="No-Break">20.10.14</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><code class="literal">ansible-core</code> (only required for <a href="https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_07.xhtml#_idTextAnchor161"><span class="No-Break"><em class="italic">Chapter 7</em></span></a><span class="No-Break">)</span></p></td><td class="No-Table-Style"><p><span class="No-Break">2.12.5</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p>Terraform (only required for <a href="https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_07.xhtml#_idTextAnchor161"><span class="No-Break"><em class="italic">Chapter 7</em></span></a><span class="No-Break">)</span></p></td><td class="No-Table-Style"><p><span class="No-Break">1.1.9</span></p></td></tr></tbody></table>

Table 12.1 – Software dependencies

## Step 1 – building a testing environment

In the following section, we describe the two automated ways of building a testing environment. If you are unsure which option is right for you, we recommend you pick the first one, as it has minimal external dependencies and is completely managed by a cloud service provider. This is also the only option that we (the authors of this book) can test and verify, and hence it should give you the most consistent experience.

### Option 1 – cloud-hosted

We have picked AWS as the cloud service provider because of its popularity and general familiarity in our industry. Inside this book’s GitHub repository (_Further reading_), we have included an Ansible playbook that completely automates all tasks required to create a VM in AWS. You are free to use any other cloud provider but you will have to do the provisioning manually.

The testing environment is a single Linux VM in AWS running `containerlab` to create container-based network topologies. The next diagram illustrates what the AWS environment looks like:

![Figure 12.1 – Target environment](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_12_01.jpg)

Figure 12.1 – Target environment

To conform with the hardware requirements stated earlier, we recommend you run at least a `t2.medium`\-, ideally a `t2.large`\-sized VM (**Elastic Compute Cloud** (**EC2**) instance). But the AWS Free Tier plan (_Further reading_) does not cover these instance types, so you should expect to incur some charges associated with the running of the VM. We assume you are familiar with the costs and billing structure of AWS and use financial common sense when working with a cloud-hosted environment.

Before you run the playbook, you need to make sure you meet the following requirements:

1.  Create an AWS account (AWS Free Tier (_Further reading_)).
2.  Create an AWS access key (AWS Programmatic access (_Further reading_)).
3.  A Linux OS with the following packages:
    -   Git
    -   Docker
    -   GNU Make

With all this in place, you can go ahead and clone the book’s GitHub repository (_Further reading_) with the `git` `clone` command:

```markup
$ git clone https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go
```

After you clone the repository, change directory to it.

#### Input variables

Before you can start the deployment, you need to supply your AWS account credentials (`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`). You do this by exporting a pair of environment variables containing the key ID and secret values, as follows. Check out AWS Programmatic access (_Further reading_) for instructions on how to create an access key:

```markup
$ export AWS_ACCESS_KEY_ID='…'
$ export AWS_SECRET_ACCESS_KEY='…'
```

Besides these required variables, there are other three optional input variables that you can adjust to fine-tune your deployment environment:

<table id="table002-1" class="No-Table-Style _idGenTablePara-1"><colgroup><col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Name</strong></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Values</strong></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">AWS_DISTRO</code></span></p></td><td class="No-Table-Style"><p><code class="literal">fedora</code> or <code class="literal">ubuntu</code> (<span class="No-Break">default: </span><span class="No-Break"><code class="literal">fedora</code></span><span class="No-Break">)</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">AWS_REGION</code></span></p></td><td class="No-Table-Style"><p>One of the AWS Regions (<span class="No-Break">default: </span><span class="No-Break"><code class="literal">us-east-1</code></span><span class="No-Break">)</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">VM_SIZE</code></span></p></td><td class="No-Table-Style"><p>One of the AWS instance types (<span class="No-Break">default: </span><span class="No-Break"><code class="literal">t2.large</code></span><span class="No-Break">)</span></p></td></tr></tbody></table>

Table 12.2 – Testing VM options

If you choose to change any of these default values, you can do this the same way as the AWS access key. Here’s an example:

```markup
$ export AWS_DISTRO=ubuntu
$ export AWS_REGION=eu-west-2
```

In that scenario, we selected Ubuntu as the Linux distribution of the VM and London (`eu-west-2`) as the AWS Region for deployment.

#### Deployment process

Once you have set all the required input variables, you can deploy the testing environment. From within the book repository directory, run the `make env-build` command, which deploys the VM and installs all the required software packages:

```markup
Network-Automation-with-Go$ make env-build
AWS_ACCESS_KEY_ID is AKIAVFPUEFZCFVFGXXXX
AWS_SECRET_ACCESS_KEY is **************************
Using /etc/ansible/ansible.cfg as config file
PLAY [Create EC2 instance] *************************************************************************************************************************************************************
TASK [Gathering Facts] *****************************************************************************************************************************************************************
ok: [localhost]
### ... <omitted for brevity > ... ###
TASK [Print out instance information for the user] *************************************************************************************************************************************
ok: [testbed] => {}
MSG:
['SSH: ssh -i lab-state/id_rsa fedora@ec2-54-86-51-96.compute-1.amazonaws.com\n', 'To upload cEOS image: scp -i lab-state/id_rsa ~/Downloads/cEOS64-lab-4.28.0F.tar fedora@ec2-54-86-51-96.compute-1.amazonaws.com:./network-automation-with-go\n']
PLAY RECAP *****************************************************************************************************************************************************************************
localhost                  : ok=28   changed=9    unreachable=0    failed=0    skipped=3    rescued=0    ignored=0   
testbed                    : ok=36   changed=24   unreachable=0    failed=0    skipped=11   rescued=0    ignored=0
```

Assuming that the playbook has completed successfully, you can see the VM access details in the logs, as the preceding output shows. You can also view the connection details at any time after you’ve deployed the environment by running the `make` `env-show` command:

```markup
Network-Automation-with-Go$ make env-show
fedora@ec2-54-86-51-96.compute-1.amazonaws.com
```

Now, you can use this information to connect to the provisioned VM. The playbook generates an **Secure Shell** (**SSH**) private key (`lab-state/id_rsa`), so don’t forget to always use it for SSH authentication:

```markup
Network-Automation-with-Go$ ssh -i lab-state/id_rsa fedora@ec2-54-86-51-96.compute-1.amazonaws.com 
fedora@testbed:~$  go version
go version go1.18.1 linux/amd64
fedora@testbed:~$  ls network-automation-with-go/
LICENSE  Makefile  README.md  ch01  ch02  ch03  ch04  ch05  ch06  ch07  ch08  ch09  ch10  ch12  lab-state  topo-base  topo-full
```

You can connect to the VM and check the Go version installed and take a look at the files of the book’s repository.

### Option 2 – self-hosted

Another option is to create a VM in a private environment. This environment could be your personal computer running a hypervisor such as VirtualBox, an ESXi server, an OpenStack cluster, or something else as long as it can allocate the CPU and memory the VM requires to run the lab topology. The OS on the VM has to be either Ubuntu 22.04 or Fedora 35.

Once you have built the VM with SSH enabled, make sure you can SSH to the IP address of the VM and access it with its credentials. Then, change the Ansible inventory file (`inventory`) in the `ch12/testbed` folder (_Further reading_) of your personal computer’s copy of this book’s GitHub repository to point to your VM. It should look something like this:

```markup
# inventory
[local-vm]
192.168.122.18
[local-vm:vars]
ansible_user=fedora
ansible_password=fedora
ansible_sudo_pass=fedora
```

Include at least the IP address (`ansible_host`) to reach the VM, and the `ansible_user`, `ansible_password`, or `ansible_ssh_private_key_file` user credentials.

In the same `ch12/testbed` folder (_Further reading_), there is an Ansible playbook that calls the `configure_instance` role. Use this playbook to auto-configure your VM to run the book examples, like so:

```markup
# configure-local-vm.yml
- name: Configure Instance(s)
  hosts: local-vm
  gather_facts: true
  vars_files:
    - ./vars/go_inputs.yml
    - ./vars/clab_inputs.yml
    - ./vars/aws_common.yml
  roles:
    - {role: configure_instance, become: true}
```

The playbook filename is `configure-local-vm.yml` and the inventory filename is `inventory`, so from the `ch12/testbed` folder (_Further reading_), run `ansible-playbook configure-local-vm.yml -i inventory -v` to get the VM ready to go.

## Step 2 – uploading container images

Not all networking vendors make it simple to access their container-based **network OSes** (**NOSes**). If you can’t pull the image directly from a container registry such as Docker Hub, you might need to download the image from their website and upload it to the test VM. The only container image in the book that we can’t pull from a public registry at the time of writing is Arista’s **cEOS** image. Here, we describe the process of uploading this image into the testing environment.

The first thing you need to do is download the image from [arista.com](http://arista.com) (_Further reading_). You should select the 64-bit cEOS image from the 4.28(F) train—for example, `cEOS64-lab-4.28.0F.tar`. You can copy the image to the test VM with the `scp` command using the generated SSH private key:

```markup
Network-Automation-with-Go$ scp -i lab-state/id_rsa ~/Downloads/cEOS64-lab-4.28.0F.tar fedora@ec2-54-86-51-96.compute-1.amazonaws.com:./network-automation-with-go
cEOS64-lab-4.28.0F.tar                        100%  434MB  26.6MB/s   00:16
```

Then, SSH to the instance and import the image with the `docker` command:

```markup
Network-Automation-with-Go$ ssh -i lab-state/id_rsa fedora@ec2-54-86-51-96.compute-1.amazonaws.com
fedora@testbed:~$  cd network-automation-with-go 
fedora@testbed:~$  docker import cEOS64-lab-4.28.0F.tar ceos:4.28
sha256:dcdc721054804ed4ea92f970b5923d8501c28526ef175242cfab0d1 58ac0085c
```

You can now use this image (`ceos:4.28`) in the `image` section of one or more routers in the topology file.

## Step 3 – iInteracting with the testing environment

We recommend you start with a fresh build of a virtual network topology at the beginning of _Chapters 6_ through _8_. To orchestrate the topologies, we use `containerlab`, which is available in the testing VM. `containerlab` offers a quick way to run arbitrary network topologies based on their definition provided in a human-readable YAML file.

Important Note

`containerlab` is written in Go and serves as a great example of an interactive CLI program that orchestrates local container resources.

You can find the following `base` topology definition file in the `topo-base` directory of this book’s GitHub repository (_Further reading_):

```markup
name: netgo
topology:
  nodes:
    srl:
      kind: srl
      image: ghcr.io/nokia/srlinux:21.6.4
    ceos:
      kind: ceos
      image: ceos:4.28.0F
      startup-config: ceos-startup
    cvx:
      kind: cvx
      image: networkop/cx:5.0.0
      runtime: docker
  links:
    - endpoints: ["srl:e1-1", "ceos:eth1"]
    - endpoints: ["cvx:swp1", "ceos:eth2"]
```

This YAML file defines a three-node topology, as the next diagram shows. One node runs Nokia SR Linux, another NVIDIA Cumulus Linux, and the last one runs Arista cEOS. In this scenario, all network devices come up with their default startup configurations, and throughout each chapter, we describe how to establish full end-to-end reachability between all three of them:

![Figure 12.2 – “Base” network topology](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_12_02.jpg)

Figure 12.2 – “Base” network topology

The next two chapters (_Chapters 9_ and _10_) rely on a slightly different version of the preceding topology. Unlike the `base` topology, the `full` topology comes up fully configured and includes an extra set of nodes to emulate physical servers attached to the network devices:

![Figure 12.3 – “Full” network topology](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_12_03.jpg)

Figure 12.3 – “Full” network topology

These end hosts run different applications that interact with the existing network topology.

Just Imagine

# Launching a virtual network topology

You can use a `containerlab` binary to deploy the test topology. For convenience, we included a couple of `make` targets that you can use:

-   `make lab-base` to create the `base` topology used in _Chapters 6_ through _8_
-   `make lab-full` to create the `full` topology used in _Chapters 9_ and _10_

Here’s an example of how you can create the `base` topology from inside the test VM:

```markup
fedora@testbed network-automation-with-go$ make lab-base
...
+---+-----------------+--------------+--------------
| # | Name            | Container ID | Image
+---+-----------------+--------------+--------------
| 1 | clab-netgo-ceos | fe422727f351 | ceos:4.28.0F
| 2 | clab-netgo-cvx  | 85e5b9135e1b | cx:5.0.0
| 3 | clab-netgo-srl  | 00106bef1d4e |srlinux:21.6.4
+---+-----------------+--------------+--------------
```

You now have `clab-netgo-ceos`, `clab-netgo-cvx` and `clab-netgo-srl` routers ready to go.

## Connecting to the devices

`containerlab` uses Docker to run the containers. This means we can use standard Docker capabilities to connect to the devices—for example, you can use the `docker exec` command to start any process inside a container:

```markup
fedora@testbed:~$  docker exec -it clab-netgo-srl sr_cli
Welcome to the srlinux CLI.                      
A:srl# show version | grep Software
Software Version  : v21.6.4
```

`sr_cli` in the preceding example is the CLI process for an SR Linux device. The following table displays the “default shell” process for each virtual network device:

<table id="table003" class="No-Table-Style _idGenTablePara-1"><colgroup><col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">NOS</strong></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Command</strong></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">Cumulus Linux</span></p></td><td class="No-Table-Style"><p><code class="literal">bash</code> <span class="No-Break">or </span><span class="No-Break"><code class="literal">vtysh</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">SR Linux</span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">sr_cli</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">EOS</span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">Cli</code></span></p></td></tr></tbody></table>

Table 12.3 – Device default shells

You can also use SSH to connect to the default shell. The next table provides the hostname and the corresponding credentials you can use to connect to each device:

<table id="table004" class="No-Table-Style _idGenTablePara-1"><colgroup><col> <col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Device</strong></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Username</strong></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Password</strong></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">clab-netgo-srl</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">admin</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">admin</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">clab-netgo-ceos</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">admin</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">admin</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">clab-netgo-cvx</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">cumulus</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">cumulus</code></span></p></td></tr></tbody></table>

Table 12.4 – Device credentials

Here’s how you can connect to Arista cEOS and Cumulus Linux, for example:

```markup
fedora@testbed:~$  ssh admin@clab-netgo-ceos
(admin@clab-netgo-ceos) Password: admin
ceos>en
ceos#exit
fedora@testbed:~$
fedora@testbed:~$  ssh cumulus@clab-netgo-cvx
cumulus@clab-netgo-cvx's password: cumulus
Welcome to NVIDIA Cumulus (R) Linux (R)
cumulus@cvx:mgmt:~$
```

Once you finish the chapter, you can destroy the topology.

## Destroying the network topology

You can clean up both virtual network topologies using the `make` `cleanup` command:

```markup
fedora@testbed:~/network-automation-with-go$ make cleanup
```

The `make cleanup` command only cleans up the virtual network topology while all the cloud resources are still running.

## Step 4 – cleaning up of the cloud-hosted environment

Once you’re done working with the cloud-hosted testing environment, you can clean it up so that you don’t pay for something you might no longer need. You can do this using another Ansible playbook that makes sure all the AWS resources you created before are now wiped out:

```markup
etwork-Automation-with-Go$ make env-delete
AWS_ACCESS_KEY_ID is AKIAVFPUEFZCFVFGXXXX
AWS_SECRET_ACCESS_KEY is **************************
PLAY [Delete EC2 instance] *************************************************************************************************************************************************************
TASK [Gathering Facts] *****************************************************************************************************************************************************************
ok: [localhost]
### ... <omitted for brevity > ... ###
TASK [Cleanup state files] *************************************************************************************************************************************************************
changed: [localhost] => (item=.region)
changed: [localhost] => (item=.vm)
PLAY RECAP *****************************************************************************************************************************************************************************
localhost                  : ok=21   changed=8    unreachable=0    failed=0    skipped=3    rescued=0    ignored=0
```

Just Imagine

# Further reading

-   Course’s GitHub repository: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go)
-   Go Playground: [https://play.golang.org/](https://play.golang.org/)
-   Official download and install procedure: [https://golang.org/doc/install#install](https://golang.org/doc/install#install)
-   `containerlab`: [https://containerlab.dev/](https://containerlab.dev/)
-   AWS Free Tier: [https://aws.amazon.com/free/](https://aws.amazon.com/free/)
-   AWS Programmatic access: [https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys)
-   `ch12/testbed`: https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch12/testbed
-   `ch12/testbed/inventory`: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch12/testbed/inventory](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch12/testbed/inventory)
-   Arista: [https://www.arista.com/en/support/software-download](https://www.arista.com/en/support/software-download)
-   Beginner’s Guide—Downloading Python: [https://wiki.python.org/moin/BeginnersGuide/Download](https://wiki.python.org/moin/BeginnersGuide/Download)
-   Installing Ansible with `pip`: [https://docs.ansible.com/ansible/latest/installation\_guide/intro\_installation.html#installing-ansible-with-pip](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html#installing-ansible-with-pip)
-   _Getting Started - Installing_ _Git_: [https://git-scm.com/book/en/v2/Getting-Started-Installing-Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
-   Installing `pip`—_Supported_ _Methods_: [https://pip.pypa.io/en/stable/installation/#supported-methods](https://pip.pypa.io/en/stable/installation/#supported-methods)
-   Get Arista cEOS: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch12/testbed/get\_arista\_ceos.md](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch12/testbed/get_arista_ceos.md)
-   AWS access keys: [https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys)
-   AWS Regions: [https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html)
-   AWS instance types: [https://aws.amazon.com/ec2/instance-types/](https://aws.amazon.com/ec2/instance-types/)