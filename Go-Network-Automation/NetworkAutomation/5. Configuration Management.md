# Configuration Management

Configuration management is a process that helps us enforce the desired configuration state on an IT system. It’s a way to make sure a network device, in our context, performs as expected as we roll out new settings. As this becomes a mundane task we perform repeatedly, it’s no surprise network configuration management is the most common network automation use case according to the NetDevOps 2020 Survey (_Further reading_).

In the previous chapter, we discussed common configuration management tasks, along with some helpful tools and libraries that can help you write programs to automate those tasks in Go. In this chapter, we will focus on a few concrete examples, taking a closer look at how Go can help us connect and interact with network devices from different networking vendors using standard protocols. We will cover four areas in this chapter:

-   Before we introduce any new examples, we will define a three-node multi-vendor virtual network lab to test the code examples in this chapter and later chapters of this book.
-   Next, we will explore how we can use Go and SSH to interact with network devices.
-   Then, we will repeat the exercise following the same program structure as with SSH but using HTTP to contrast these different options.
-   Finally, we will extract and parse the resulting operational state to verify that our configuration changes have been successful.

Note that we have deliberately avoided talking about YANG-based APIs here as we will cover them extensively in the last few chapters of this book.

In this chapter, we will cover the following topics:

-   Environment setup
-   Interacting with network devices via SSH
-   Interacting with network devices via HTTP
-   State validation

Just Imagine

# Technical requirements

You can find the code examples for this chapter in the book’s GitHub repository: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go), under the `ch06` folder.

Important Note

We recommend that you execute the Go programs in this chapter in a virtual lab environment. Refer to the _Appendix_ for prerequisites and instructions on how to build it.

Just Imagine

# Environment setup

One of the easiest and safest ways to learn and experiment with network automation is to build a lab environment. Thanks to the progress we’ve had in the last decade, today, we have access to virtualized and containerized network devices from different networking vendors and plenty of tools that can help us build a virtual topology from them.

In this book, we will use one of those tools: **Containerlab**. This tool, which is written in Go, allows you to build arbitrary network topologies from container images. The fact that you can create and run topologies based on a plain YAML file in a matter of seconds makes it a strong choice to run quick tests. Please refer to the _Appendix_ for installation instructions and recommendations for host operating systems.

## Creating the topology

Throughout the rest of this book, we will work with a base network topology consisting of three containerized network devices running different **network operating** **systems** (**NOSes**):

-   `srl`: Running Nokia’s **Service Router Linux** (**SR Linux**)
-   `cvx`: Running NVIDIA’s Cumulus Linux
-   `ceos`: Running Arista’s EOS

The following diagram depicts the device interconnections. They all come up with their default (blank) configuration:

![Figure 6.1 – Test topology](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_06_01.jpg)

Figure 6.1 – Test topology

We can describe this topology with the following YAML file, which is a representation that **Containerlab** can interpret and translate into a running topology:

```markup
name: netgo
topology:
  nodes:
    srl:
      kind: srl
      image: ghcr.io/nokia/srlinux:21.6.4
    ceos:
      kind: ceos
      image: ceos:4.26.4M
    cvx:
      kind: cvx
      image: networkop/cx:5.0.0
      runtime: docker
  links:
    - endpoints: ["srl:e1-1", "ceos:eth1"]
    - endpoints: ["cvx:swp1", "ceos:eth2"]
```

You can find this YAML file, like the rest of the code examples, in this book’s GitHub repository, specifically in the `topo-base` directory. If you go through the _Appendix_ to learn more about Containerlab or you have it running already, you can bring up the entire lab with the following command:

```markup
topo-base$ sudo containerlab deploy -t topo.yml --reconfigure
```

Once the lab is up, you can access each device by its hostname using the credentials shown in the following table:

<table id="table001-1" class="No-Table-Style _idGenTablePara-1"><colgroup><col> <col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">Device</span></p></td><td class="No-Table-Style"><p><span class="No-Break">Username</span></p></td><td class="No-Table-Style"><p><span class="No-Break">Password</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">clab-netgo-srl</span></p></td><td class="No-Table-Style"><p><span class="No-Break">admin</span></p></td><td class="No-Table-Style"><p><span class="No-Break">admin</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">clab-netgo-ceos</span></p></td><td class="No-Table-Style"><p><span class="No-Break">admin</span></p></td><td class="No-Table-Style"><p><span class="No-Break">admin</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">clab-netgo-cvx</span></p></td><td class="No-Table-Style"><p><span class="No-Break">cumulus</span></p></td><td class="No-Table-Style"><p><span class="No-Break">cumulus</span></p></td></tr></tbody></table>

Table 6.1 – Containerlab access credentials

For example, to access NVIDIA’s device via SSH, you would execute `ssh cumulus@clab-netgo-cvx`:

```markup
⇨  ssh cumulus@clab-netgo-cvx
cumulus@clab-netgo-cvx's password: cumulus
Linux cvx 5.14.10-300.fc35.x86_64 #1 SMP Thu Oct 7 20:48:44 UTC 2021 x86_64
Welcome to NVIDIA Cumulus (R) Linux (R)
cumulus@cvx:mgmt:~$ exit
```

If you want to learn more about Containerlab or run this lab setup in the cloud, check out the instructions in the _Appendix_ of this book.

Just Imagine

# Interacting with network devices via SSH

**Secure Shell** (**SSH**) is the predominant protocol that network engineers use to securely access and configure network devices via a **command-line interface** (**CLI**) that transports unstructured data to display to end users. This interface simulates a computer terminal, so we’ve used it traditionally for human interactions.

One of the first steps network engineers take when they embark on the journey of automating mundane tasks is to create scripts that run a set of CLI commands for them in sequence to achieve an outcome. Otherwise, they would run the commands themselves interactively via an SSH pseudo-terminal.

While this gives us speed, this is not the only benefit of network automation. As we cover different technologies through the rest of this book, other benefits, such as reliability, repeatability, and consistency, to name a few, become a common theme. For now, we will start by crafting an SSH connection to a network device in Go and send configuration commands line by line, to then take advantage of a higher-level package in Go that abstracts away the connection details of the different networking vendors, making the development experience simpler for network engineers.

## Describing the network device configurations

The first task we want to do with Go is to configure each of the devices of the three-node topology we defined in the preceding section. As a learning exercise, we will create three different Go programs to configure each device independently so that you can contrast the different approaches. While each program is unique, they all follow the same design structure. One program uses SSH to connect and configure a device, another one uses Scrapligo, and the last one uses HTTP, as we’ll cover in the next section.

To make the code examples meaningful, but at the same time not overly complicated, we have limited the device configurations to apply to the following sections:

-   A unique IPv4 address on each of the transit links
-   A **Border Gateway Protocol** (**BGP**) peering established between those IPs
-   A unique loopback address that is also redistributed into BGP

The goal of these settings is to establish reachability between all three loopback interfaces.

In real-life automation systems, developers strive to find a common data model you can use to represent device configurations for any vendor. The two main examples of this are IETF and OpenConfig YANG models. We will do the same in this case by defining a standard schema for the input data we will use for all three network devices but using Go directly to define the data structures instead of the YANG modeling language. This schema has just enough information to meet the goal of establishing end-to-end reachability:

```markup
type Model struct {
    Uplinks  []Link `yaml:"uplinks"`
    Peers    []Peer `yaml:"peers"`
    ASN      int    `yaml:"asn"`
    Loopback Addr   `yaml:"loopback"`
}
type Link struct {
    Name   string `yaml:"name"`
    Prefix string `yaml:"prefix"`
}
type Peer struct {
    IP  string `yaml:"ip"`
    ASN int    `yaml:"asn"`
}
type Addr struct {
    IP string `yaml:"ip"`
}
```

In each of the programs, we supply the parameters to the data model to generate the device’s configuration via the `input.yml` file, which is available in the program’s folder. For the first example, this file looks as follows:

```markup
# input.yml
asn: 65000
 
loopback: 
  ip: "198.51.100.0"
 
uplinks:
  - name: "ethernet-1/1"
    prefix: "192.0.2.0/31"
 
peers:
  - ip: "192.0.2.1"
    asn: 65001
```

After we open this file for reading, we deserialize this information into an instance of a `Model` type – which represents the data model – with the `Decode` method. The following output represents these steps:

```markup
func main() {
    src, err := os.Open("input.yml")
    // process error
    defer src.Close()
    d := yaml.NewDecoder(src)
    var input Model
    err = d.Decode(&input)
    // process error
}
```

Then, we pass the input variable (of the `Model` type) to a config generator function (`devConfig`), which transforms this information into syntax that the target device can understand. The result of this transformation is a vendor-specific configuration serialized into bytes that you can transfer to the remote device.

A transport library establishes the connection to the remote device using default credentials, which you can overwrite via command-line flags. The session we have created has an `io.Writer` element that we can use to send the configuration to the remote device:

![Figure 6.2 – Program structure](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_06_02.jpg)

Figure 6.2 – Program structure

Now that we’re familiar with the structure of the program, let’s explore different implementations of it to learn more about the Go packages that are available to communicate with network devices, starting with SSH and Scrapligo.

## Using Go’s SSH package to access network devices

The first device from the topology that we are configuring is the containerized Nokia **SR Linux**. Although this NOS supports a variety of interfaces, including structured APIs such as gNMI and NETCONF, in this case, we are configuring it interactively via SSH, using the same commands that a human operator would use. We will execute these commands as a multi-line string, which we can craft using Go’s `text/template` template package.

Go’s SSH package, `golang.org/x/crypto/ssh`, belongs to a set of packages that are still part of the Go project but developed outside the main Go tree under looser compatibility requirements. Although this is not the only SSH Go client, other packages tend to reuse parts of this package, so they become higher-level abstractions.

As described in the general program design, we use the `Model` data structure to hold the device configuration inputs and merge them with the `srlTemplate` template to produce a valid device configuration as a buffer of bytes:

```markup
const srlTemplate = `
enter candidate
{{- range $uplink := .Uplinks }}
set / interface {{ $uplink.Name }} subinterface 0 ipv4 address {{ $uplink.Prefix }}
set / network-instance default interface {{ $uplink.Name }}.0
{{- end }}
...
`
```

The `srlTemplate` constant has a template that starts by looping (using the `range` keyword) over the uplinks of a `Model` instance. For each `Link`, it takes the `Name` and `Prefix` properties of it to create a couple of CLI commands we can place in a buffer. In the following code, we are running the `Execute` method to pass the inputs via the `in` variable and put the binary representation of interactive CLI commands on `b`, which we later expect to send to the remote device (`cfg`):

```markup
func devConfig(in Model)(b bytes.Buffer, err error){
    t, err := template.New("config").Parse(srlTemplate)
    // process error
    err = t.Execute(&b, in)
    // process error
    return b, nil
}
func main() {
    /* ... <omitted for brevity > ... */
    var input Model
    err = d.Decode(&input)
    // process error
    cfg, err := devConfig(input)
    /* ... <continues next > ... */
}
```

We have hardcoded the authentication credentials to the correct values to fit the lab, but you can override them if necessary. We use these arguments to establish initial connectivity with the `srl` network device:

```markup
func main() {
    /* ... <continues from before > ... */
    settings := &ssh.ClientConfig{
        User: *username,
        Auth: []ssh.AuthMethod{
            ssh.Password(*password),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    conn, err := ssh.Dial(
        "tcp",
        fmt.Sprintf("%s:%d", *hostname, sshPort),
        settings,
    )
    // process error
    defer conn.Close()
    /* ... <continues next > ... */
}
```

If the authentication credentials are correct and there are no connectivity problems, the `ssh.Dial` function returns a connection handler (`conn`), representing a single SSH connection. This connection acts as a single transport for potentially various channels. One such channel is a pseudo-terminal session used for interactive communication with the remote device, but it may also include extra channels that you can use for port forwarding.

The following code snippet spawns a new terminal session and sets the expected terminal parameters, such as terminal height, width, and **TeleTYpe** (**TTY**) speed. The `ssh.Session` type provides functions to retrieve standard input and standard output pipes that connect to the remote terminal:

```markup
func main() {
    /* ... <continues from before > ... */
    session, err := conn.NewSession()
    // process error
    defer session.Close()
    modes := ssh.TerminalModes{
        ssh.ECHO:          1,
        ssh.TTY_OP_ISPEED: 115200,
        ssh.TTY_OP_OSPEED: 115200,
    }
    if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
        log.Fatal("request for pseudo terminal failed: ", err)
    }
    stdin, err := session.StdinPipe()
    // process error
    stdout, err := session.StdoutPipe()
    // process error
    session.Shell()
    /* ... <continues next > ... */
}
```

In conformance with the rest of the Go packages, standard input and standard output pipes implement the `io.Writer` and `io.Reader` interfaces, respectively. This means you can use them to write data in to and read output from the remote network device. We will go back to the `cfg` buffer with the CLI config and use the `WriteTo` method to send this config over to the target node:

```markup
func main() {
    /* ... <continues from before > ... */
    log.Print("connected. configuring...")
    cfg.WriteTo(stdin)
}
```

This is the expected output of this program:

```markup
ch06/ssh$ go run main.go 
go: downloading golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce
go: downloading gopkg.in/yaml.v2 v2.4.0
2022/02/07 21:11:44 connected. configuring...
2022/02/07 21:11:44 disconnected. dumping output...
enter candidate
set / interface ethernet-1/1 subinterface 0 ipv4 address 192.0.2.0/31
set / network-instance default interface ethernet-1/1.0
...
set / network-instance default protocols bgp ipv4-unicast admin-state enable
commit now
quit
Using configuration file(s): []
Welcome to the srlinux CLI.
Type 'help' (and press <ENTER>) if you need any help using this.
--{ running }--[  ]--                                                           
A:srl#                                                                          
--{ running }--[  ]--                                                           
A:srl# enter candidate                                                          
--{ candidate shared default }--[  ]--                                          
A:srl# set / interface ethernet-1/1 subinterface 0 ipv4 address 192.0.2.0/31    
--{ * candidate shared default }--[  ]-- 
.......                                
--{ * candidate shared default }--[  ]--                                        
A:srl# commit now                                                               
All changes have been committed. Leaving candidate mode.
--{ + running }--[  ]--                                                         
A:srl# quit
```

You can find the complete example in the `ch06/ssh` folder (_Further reading_).

## Automating routine SSH tasks

Common network elements, such as routers and switches, display data for people rather than computers via the CLI. We rely on screen scraping to let our programs consume this human-readable data. One popular screen-scraping Python library, whose name comes from _scrape cli_, is Scrapli.

Scrapli has a version in Go, which we will explore in the following example, called Scrapligo. The goal of this package is to offer the next layer of abstraction on top of SSH and hide away some transport complexities while providing several convenient functions and supporting the CLI flavors of different networking vendors.

To show `scrapligo` in action, we will configure another network device in the topology: Arista’s cEOS (`ceos`). Just like we did with `srl`, we will use a list of CLI commands to push the desired network state so that the initial steps of parsing and instantiating a string from a template are the same. What changes is the template, which uses Arista EOS’s syntax:

```markup
const ceosTemplate = `
...
!
router bgp {{ .ASN }}
  router-id {{ .Loopback.IP }}
{{- range $peer := .Peers }}  
  neighbor {{ $peer.IP }} remote-as {{ $peer.ASN }}
{{- end }}
  redistribute connected
!
`
```

The difference starts when we get to the SSH connection setup. We create a device driver (`GetNetworkDriver`) to connect to the remote device with the device hostname and authentication credentials. The platform definition comes from the `platform` package of `scrapligo`. From then on, it only takes a single method call on this driver to open an SSH connection to the remote device:

```markup
func main() {
    /* ... <omitted for brevity > ... */
    conn, err := platform.NewPlatform(
        *nos,
        *hostname,
        options.WithAuthNoStrictKey(),
        options.WithAuthUsername(*username),
        options.WithAuthPassword(*password),
    )
    // process error  
    driver, err := conn.GetNetworkDriver()
    // process error  
 
    err = driver.Open()
    // process error  
    defer driver.Close()
    /* ... <continues next > ... */
}
```

One of the extra features that `scrapli` offers is the `cscrapligocfg` package, which defines a high-level API to work with a remote network device’s configuration. This API understands different CLI flavors, it can sanitize a configuration before sending it to the device, and it can generate configuration diffs for us. But, most importantly, this package allows for a single function call to load the entire device configuration as a string, taking care of things such as privilege escalation and configuration merging or replacement. We will do this here with the `LoadConfig` method:

```markup
func main() {
    /* ... <continues from before > ... */
    conf, err := cfg.NewCfg(driver, *nos)
    // process error
 
    // sanitize config by removing keywords like "!" and "end"
    err = conf.Prepare()
    // process error
 
    response, err = conf.LoadConfig(config.String(), false)
    // process error
}
```

These are all the steps you need to configure the device in this case. After you run the program with `go run`, you can `ssh` to the device to check that the configuration is now there:

```markup
ch06/scrapli$ go run main.go 
2022/02/14 17:06:16 Generated config: 
!
configure
!
ip routing
!
interface Ethernet1
  no switchport
  ip address 192.0.2.1/31
!
...
```

Normally, to get a response coming back from a device, we need to read the response buffer carefully until we see a command-line prompt, as it normally ends with an **end-of-file** (**EOF**). Although we don’t show it here, `scrapligo` can do this for us by reading the received buffer and converting the response into a string.

Another popular Go SSH package that provides a high-level API to execute commands at scale is `yahoo/vssh`. We won’t cover it here, but you can find an example in the `ch06/vssh` directory of this book’s repository (_Further reading_) to configure the network devices of the topology.

Just Imagine

# Interacting with network devices via HTTP

Over the last decade, networking vendors have begun to include **application programming interfaces** (**APIs**) to manage their devices as a supplement to the CLI. It’s not uncommon to find network devices with a robust RESTful API that gives you read and write access to it.

A RESTful API is a stateless client-server communication architecture that runs over HTTP. The request and responses generally transport structured data (JSON, XML, and so on), but they might as well carry plain text. This makes the RESTful API a better-suited interface for machine-to-machine interactions.

## Using Go’s HTTP package to access network devices

The remaining device to configure is NVIDIA’s Cumulus Linux (`cvx`). We will use its OpenAPI-based RESTful API to configure it. We will encode the configuration in a JSON message and send it over an HTTP connection with Go’s `net/http` package.

As in the SSH examples, we normally load the input data and transform it into the shape the target device expects with the `devConfig` function, but in this case, it’s a JSON payload. Because of this, we no longer need templates to build the network device configuration, as we can now use data structures in Go to encode and decode data from JSON or any other encoding format.

The data structures represent the configuration data model of the target device. Ideally, this data model would match the one we defined previously, so we don’t need to define anything else. But that’s not what we see in the field, where all the network vendors have proprietary data models. The good news is that both IETF and OpenConfig offer vendor-agnostic models; we’ll explore these later in [_Chapter 8_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_08.xhtml#_idTextAnchor182), _Network APIs_. For now, these are some of the data structures we will use for this device’s configuration:

```markup
type router struct {
    Bgp
}
type bgp struct {
    ASN      int
    RouterID string
    AF       map[string]addressFamily
    Enabled  string
    Neighbor map[string]neighbor
}
type neighbor struct {
    RemoteAS int
    Type     string
}
```

Inside the main function, we parse the program flags and use them to store the HTTP connection settings inside a data structure with all the details required to build an HTTP request, including any non-default transport settings for an HTTP client. We do this entirely for convenience purposes as we want to pass these details to different functions:

```markup
type cvx struct {
    url   string
    token string
    httpC http.Client
}
func main() {
    /* ... <omitted for brevity > ... */
    device := cvx{
        url:   fmt.Sprintf("https://%s:%d", *hostname, defaultNVUEPort),
        token: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", *username, *password))),
        httpC: http.Client{
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
            },
        },
    }
    /* ... <continues next > ... */
}
```

Now, we can send the configuration over and make it a candidate config on the target device. We can later apply this configuration on the device by referencing the revision ID we associate our desired configuration with. Let’s look at the steps to do this that showcase different attributes to consider when working with HTTP.

First, we will create a new revision ID, which we include as a query parameter `(?rev=<revisionID>`) in the URL to connect to the device API. Now, the `addr` is variable the target device URL that contains `device hostname` and `revisionID`:

```markup
func main() {
    /* ... <continues from before > ... */
    // create a new candidate configuration revision
    revisionID, err := createRevision(device)
    // process error
    addr, err := url.Parse(device.url + "/nvue_v1/")
    // process error
    params := url.Values{}
    params.Add("rev", revisionID)
    addr.RawQuery = params.Encode()
    /* ... <continues next > ... */
}
```

With the URL linked to the revision ID, we put together the PATCH request for the configuration change. This points to `addr` and `cfg`, which is the JSON device configuration that the `devConfig` function returns. We also add an HTTP `Authorization` header with the encoded username and password and signal that the payload is a JSON message:

```markup
func main() {
    /* ... <continues from before > ... */
    req, err := http.NewRequest("PATCH", addr.String(), &cfg)
    // process error
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", "Basic "+device.token)
    /* ... <continues next > ... */
}
```

Once we have the HTTP request built, we can pass it to the device HTTP client’s method, `Do`, which serializes everything into a binary format, sets up a TCP session, and sends the HTTP request over it.

Finally, to apply the candidate configuration changes, we must make another PATCH request inside the `applyRevision` function:

```markup
func main() {
    /* ... <continues from before > ... */
    res, err := device.httpC.Do(req)
    // process error
    defer res.Body.Close()
    // Apply candidate revision
    if err := applyRevision(device, revisionID); err != nil {
        log.Fatal(err)
    }
}
```

You can find the code for this example in the `ch06/http` directory of this book’s GitHub repository (_Further reading_). This is what you should see when you run this program:

```markup
ch06/http$ go run main.go 
2022/02/14 16:42:26 generated config {
 "interface": {
  "lo": {
   "ip": {
    "address": {
     "198.51.100.2/32": {}
...
 "router": {
  "bgp": {
   "autonomous-system": 65002,
   "router-id": "198.51.100.2"
  }
 },
 "vrf": {
  "default": {
   "router": {
    "bgp": {
...
     "enable": "on",
     "neighbor": {
      "192.0.2.2": {
       "remote-as": 65001,
       "type": "numbered"
      },
      "203.0.113.4": {
       "remote-as": 65005,
       "type": "numbered"
      }
...
}
2022/02/14 16:42:27 Created revisionID: changeset/cumulus/2022-02-14_16.42.26_K4FJ
{
  "state": "apply",
  "transition": {
    "issue": {},
    "progress": ""
  }
}
```

Just like with SSH, we rarely use `net/http` directly in our programs to interact with a REST API and normally use a higher-level package instead.

## Getting config inputs from other systems via HTTP

Until this point, the data to generate a particular device configuration has come from a static file that is present in the program’s folder. These values are network device vendor-agnostic.

In real-world network automation systems, these values can come from other systems. For example, an **IP address management** (**IPAM**) tool can allocate IP addresses dynamically via a REST API call for a particular device, which you can use to build its configuration. The collection of systems that supply these parameters becomes what some refer to as the _source of truth_. Nautobot is an infrastructure resource modeling application that falls into this category.

This also highlights the fact that to automate networks, we not only need to interact with network devices but also integrate with other systems such as Nautobot. This is why we are dedicating this example to exploring how to Go use to interact with a free public instance of Nautobot available for anyone at [https://demo.nautobot.com/](https://demo.nautobot.com/).

The Go client package for Nautobot is automatically generated from its OpenAPI specification, which means its structure might be familiar to you if you have already worked with other OpenAPI-derived packages, which is an advantage of machine-generated code.

In the following example, we are using the auto-generated Nautobot Go package to define a Nautobot API client pointing to [https://demo.nautobot.com/](https://demo.nautobot.com/) with an API token:

```markup
func main() {
    token, err := NewSecurityProviderNautobotToken("...")
    // process error
 
    c, err := nb.NewClientWithResponses(
        "https://demo.nautobot.com/api/",
        nb.WithRequestEditorFn(token.Intercept),
    )
    /* ... <continues next > ... */
}
```

The `c` client allows us to interact with the remote Nautobot instance. In this example, we want to add one of the lab topology nodes (`ceos`) to the **data center infrastructure management** (**DCIM**) resource collection of the Nautobot instance. The device details are in the `device.json` file:

```markup
{
    "name": "ams01-ceos-02",
    "device_type": {
        "slug": "ceos"
    },
    "device_role": {
        "slug": "router"
    },
    "site": {
        "slug": "ams01"
    }
}
```

Before we can add the device to Nautobot, we must make sure the device type, device role, and site we are referencing in the `device.json` file exist by name already in Nautobot. The `createResources` function takes care of this. Then, we get the IDs of these resources (device type, device role, and site) with the `getDeviceIDs` function, to associate the new device with its type, role, and site:

```markup
func main() {
    /* ... <continues from before > ... */
    err = createResources(c)
    // process error
 
    dev, err := os.Open("device.json")
    // process error
    defer dev.Close()
 
    d := json.NewDecoder(dev)
 
    var device nb.Device
    err = d.Decode(&device)
    // process error
 
    found, devWithIDs, err := getDeviceIDs(c, device)
    /* ... <continues next > ... */
}
```

If the device is not already in Nautobot, we can create it with the auto-generated `DcimDevicesCreateWithResponse` function:

```markup
func main() {
    /* ... <continues from before > ... */
    created, err := c.DcimDevicesCreateWithResponse(
        context.TODO(),
        nb.DcimDevicesCreateJSONRequestBody(*devWithIDs))
    check(err)
}
```

After running the program with `go run nautobot` from the `ch06/nautobot` folder, you should see the following in the Nautobot graphical interface at [https://demo.nautobot.com/](https://demo.nautobot.com/):

![Figure 6.3 – Nautobot screenshot](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_06_03.jpg)

Figure 6.3 – Nautobot screenshot

The data that we pass to these Dcim functions ends up in HTTP requests, just like the ones we built manually earlier in this chapter. Here, we don’t deal with URL queries, HTTP paths, or JSON payloads directly as the package abstracts away all that from us. This allows the developers to focus more on business value and less on implementation details. It makes the API easier to consume.

The focus of this chapter so far has been more on pushing configurations down to network devices and less on reading the state of the network after this operation. While configuration management’s primary focus is on producing and deploying configurations in the correct format, state validation can play a key role in verifying your configuration changes have been successful. In the next section, we will learn how to retrieve and parse operational data from a remote device.

Just Imagine

# State validation

The way network devices model and store their state internally is often different from their configuration data model. Traditional CLI-first network devices display the state in a tabular format to the end user, making it easier for network operators to interpret and reason about it. In API-enabled network operating systems, they can present the state in a structured format, making the data friendlier for automation, but we still need to prepare the right data model for deserialization.

In this section, we will look at three different methods you could use to read the state from a network device through a code example that gathers operational data from the devices we just configured with `crypto/ssh`, `net/http`, and `scrapligo` in the preceding sections of this chapter. For each network device, we will use one of these resources to get the data in the format we need:

-   **RESTful API calls**: To retrieve and parse data from an HTTP interface
-   **Regular expressions**: To parse plain text received via SSH
-   **TextFSM templates**: To simplify parsing tabular data

## Checking routing information

At this point, you should have a three-node topology running. Each network device has a loopback address we redistribute into BGP. Arista cEOS’s loopback address is `198.51.100.1/32`, for example. The goal of the next program is to verify the setup. We retrieve the routing table information from every device to check whether all three IPv4 loopback addresses are present. This way, we can verify our configuration intent – established end-to-end reachability between all devices.

The program has two building blocks:

-   `GetRoutes`: A method that connects to the network device, gets the information we need, and puts it in a common format
-   `checkRoutes`: A function that reads the routes from `GetRoutes` and compares them to the list of loopback addresses we expect to see (`expectedRoutes`)

One caveat is that the API type a network device supports to access its operational data remotely may vary, from the transport protocol to the format of the textual representation of the data. In our example, this translates into different implementation details of `GetRoutes` per networking vendor. Here, we take it a bit to the extreme for educational purposes and make the implementation per vendor completely different from one another to showcase REST APIs, regular expressions, and TextFSM independently:

![Figure 6.4 – Checking routing information](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_06_04.jpg)

Figure 6.4 – Checking routing information

Each network device has its own data structure. For example, we create SRL for SR Linux. The `SRL`, `CVX`, and `CEOS` types implement the `Router` interface, as each one has a `GetRoutes` method that contains the implementation details for that specific vendor.

In the main program, a user only needs to initialize the devices with the authentication details, so it creates a variable of the type we created for that device. Then, it can run the route collection tasks concurrently by firing off a goroutine for each device that runs the device type’s `GetRoutes` method. The `Router` interface successfully hides away the implementation details of a particular vendor from the user, as the call is always the same `router.GetRoutes`:

```markup
type Router interface {
    GetRoutes(wg *sync.WaitGroup)
}
 
func main() {
     cvx := CVX{
     Hostname: "clab-netgo-cvx",
      Authentication: Authentication{
      Username: "cumulus",
     Password: "cumulus",
     },
    }
    srl := SRL{
     Hostname: "clab-netgo-srl",
     Authentication: Authentication{
      Username: "admin",
      Password: "admin",
     },
    }
    ceos := CEOS{
     Hostname: "clab-netgo-ceos",
     Authentication: Authentication{
      Username: "admin",
      Password: "admin",
     },
    }
 
    log.Printf("Checking reachability...")
 
    devices := []Router{cvx, srl, ceos}
 
    var wg sync.WaitGroup
    for _, router := range devices {
        wg.Add(1)
        go router.GetRoutes(&wg)
    }
    wg.Wait()
}
```

Because all `GetRoutes` instances run in the background in their own goroutine, we added a `wg` wait group to make sure we don’t finish the main goroutine until we have collected and verified all the devices. Before the end of each `GetRoutes` method, we call the `expectedRoutes` function to process the routes we get from that device.

We verify the parsed state (routes) by checking that each `expectedRoutes`, which contains a unique set of loopback addresses, is present in each device’s routing table. For every IPv4 prefix received, we check whether it’s present in `expectedRoutes` and change a boolean flag to signal this. If, by the end of this, we have prefixes in `expectedRoutes` with a Boolean value of `false`, it means they were not present in the device’s routing table, and we create a log message:

```markup
func checkRoutes(device string, in []string, wg *sync.WaitGroup) {
    defer wg.Done()
    log.Printf("Checking %s routes", device)
    expectedRoutes := map[string]bool{
        "198.51.100.0/32": false,
        "198.51.100.1/32": false,
        "198.51.100.2/32": false,
    }
    for _, route := range in {
        if _, ok := expectedRoutes[route]; ok {
            log.Print("Route ", route,
                        " found on ", device)
            expectedRoutes[route] = true
        }
    }
    for route, found := range expectedRoutes {
        if !found {
            log.Print("! Route ", route, 
                        " NOT found on ", device)
        }
    }
}
```

Following this, we examine each of the `GetRoutes` method implementations. As with the rest of the examples, you can find the complete program in the `ch06/state` folder of this book’s GitHub repository (_Further reading_).

### Parsing command outputs with regular expressions

We use regular expressions to parse and extract information from unstructured data. The Go standard library includes the `regexp` package, which understands the RE2 syntax. This is a regular expression library designed with safety as one of its primary goals. One of the main consequences of that decision is the lack of back-references and look-around operations, which are unsafe and can lead to denial of service exploits.

In this case, the `GetRoutes` method uses `scrapligo` to connect and sends a `show` command to extract the routing table information from an SRL device type in this case. One way to parse this information is to iterate over the output line by line while matching expected patterns with regular expressions, close to what we did for the `ch05/closed-loop` example (_Further reading_):

```markup
func (r SRL) GetRoutes(wg *sync.WaitGroup) {
    lookupCmd := "show network-instance default route-table ipv4-unicast summary"
 
    conn, err := platform.NewPlatform(
        "nokia_srl",
        r.Hostname,
        options.WithAuthNoStrictKey(),
        options.WithAuthUsername(r.Username),
        options.WithAuthPassword(r.Password),
        options.WithTermWidth(176),
    )
    // process error
 
    driver, err := conn.GetNetworkDriver()
    // process error
    err = driver.Open()
    // process error 
    defer driver.Close()
 
    resp, err := driver.SendCommand(lookupCmd)
    // process error
 
    ipv4Prefix := regexp.
            MustCompile(`(\d{1,3}\.){3}\d{1,3}\/\d{1,2}`)
 
    out := []string{}
    for _, match := range ipv4Prefix.FindAll(
    resp.RawResult, -1) {
        out = append(out, string(match))
    }
    go checkRoutes(r.Hostname, out, wg)
}
```

To make things a bit simpler, we assume that anything that matches the IPv4 address pattern in the entire output is a prefix installed in the routing table. This way, instead of reading and parsing a tabular data structure, we tell our program to find all text occurrences that match the IPv4 route pattern and put them on a string slice (`out`) that we pass to the `checkRoutes` function for further processing.

### Parsing semi-formatted command outputs with templates

Parsing various output formats with regular expressions can be tedious and error-prone. This is why Google created `TextFSM`, initially as a Python library, to implement a template-based parsing of semi-formatted text. They designed it specifically to parse information from network devices and it has a wide range of community-developed templates maintained in **ntc-templates** (_Further reading_).

We will use one of these community templates to parse the `ip` route command’s output in the implementation of `GetRoutes` for Arista cEOS. Scrapligo embeds a Go port of TextFSM and can conveniently parse the response using the `TextFsmParse` function:

```markup
func (r CEOS) GetRoutes(wg *sync.WaitGroup) {
    template := "https://raw.githubusercontent.com/networktocode/ntc-templates/master/ntc_templates/templates/arista_eos_show_ip_route.textfsm"
    lookupCmd := "sh ip route"
    conn, err := core.NewEOSDriver(
        r.Hostname,
        base.WithAuthStrictKey(false),
        base.WithAuthUsername(r.Username),
        base.WithAuthPassword(r.Password),
    )
    // process error
    err = conn.Open()
    // process error
    defer conn.Close()
    resp, err := conn.SendCommand(lookupCmd)
    // process error
    parsed, err := resp.TextFsmParse(template)
    // process error
    out := []string{}
    for _, match := range parsed {
        out = append(out, fmt.Sprintf(
                "%s/%s", match["NETWORK"], match["MASK"]))
    }
    go checkRoutes(r.Hostname, out, wg)
}
```

The `parsed` variable that stores the parsed data is a slice that contains `map[string]interface{}` values, where keys correspond to the TextFSM values defined in a template. Thus, just by looking at the `show ip route` template, we can extract the network and mask (prefix length) information and append it to a string slice (`out`) that we pass to the `checkRoutes` function for further processing.

### Getting JSON-formatted data with REST API requests

Thus far in this chapter, we’ve seen two different ways of interacting with a REST API – one using the `net/http` package and another using an auto-generated high-level package (`nautobot`). But you also have other options, such as `go-resty`, which builds on top of `net/http` to offer an improved user experience when interacting with REST API endpoints.

In the following implementation of `GetRoutes`, we are taking advantage of `go-resty` to build the required HTTP headers for authentication, extend the URL with query parameters, and unmarshal a response into a user-defined data structure (`routes`):

```markup
Code Block 1:
func (r CVX) GetRoutes(wg *sync.WaitGroup) {
client := resty.NewWithClient(&http.Client{
Transport: &http.Transport{
TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
},
})
client.SetBaseURL("https://" + r.Hostname + ":8765" )
client.SetBasicAuth(r.Username, r.Password)
var routes map[string]interface{}
_, err := client.R().
SetResult(&routes).
SetQueryParams(map[string]string{
"rev": "operational",
}).
Get("/nvue_v1/vrf/default/router/rib/ipv4/route")
// process error
out := []string{}
for route := range routes {
out = append(out, route)
}
go checkRoutes(r.Hostname, out, wg)
}
```

We have created a REST API client to request the routing table information (`...rib/ipv4/route`) from the target device (type CVX). We decoded the JSON payload response with the routing table prefixes as keys into the `routes` variable of the `map[string]interface{}` type. Next, we looped through `routes` to append all keys to a string slice (`out`) we can pass to the `checkRoutes` function.

## Validating end-to-end reachability

You can run this program to check whether all three routers in the topology can reach one another from the `ch06/state` folder (_Further reading_). Make sure all the devices have the configs from the examples that used `crypto/ssh`, `net/http`, and `scrapligo` to configure them earlier in this chapter. The expected output should look as follows:

```markup
ch06/state$ go run main.go 
2022/03/10 17:06:30 Checking reachability...
2022/03/10 17:06:30 Collecting CEOS routes
2022/03/10 17:06:30 Collecting CVX routes
2022/03/10 17:06:30 Collecting SRL routes
2022/03/10 17:06:30 Checking clab-netgo-cvx routes
2022/03/10 17:06:30 Route 198.51.100.0/32 found on clab-netgo-cvx
2022/03/10 17:06:30 Route 198.51.100.1/32 found on clab-netgo-cvx
2022/03/10 17:06:30 Route 198.51.100.2/32 found on clab-netgo-cvx
2022/03/10 17:06:31 Checking clab-netgo-ceos routes
2022/03/10 17:06:31 Route 198.51.100.0/32 found on clab-netgo-ceos
2022/03/10 17:06:31 Route 198.51.100.1/32 found on clab-netgo-ceos
2022/03/10 17:06:31 Route 198.51.100.2/32 found on clab-netgo-ceos
2022/03/10 17:06:34 Checking clab-netgo-srl routes
2022/03/10 17:06:34 Route 198.51.100.0/32 found on clab-netgo-srl
2022/03/10 17:06:34 Route 198.51.100.1/32 found on clab-netgo-srl
2022/03/10 17:06:34 Route 198.51.100.2/32 found on clab-netgo-srl
```

If any of the routes were not present on any of the devices, we would’ve seen messages such as these:

```markup
2022/03/10 15:59:55 ! Route 198.51.100.0/32 NOT found on clab-netgo-cvx
2022/03/10 15:59:55 ! Route 198.51.100.1/32 NOT found on clab-netgo-cvx
```

Just Imagine

# Summary

Configuration generation, deployment, reporting, and compliance remain the most popular network automation operations. This is where the immediate benefits of introducing automation are greatest and most visible, making it the first logical step into the world of automation and DevOps. Configuration management is one of those repetitive tasks network engineers spend most of their time on, so it’s a natural fit for automation. But sending a new configuration to a device is just part of a broader process that should consider failure handling, from syntax errors in the configuration to how to recover properly if the connection to a remote device drops. In this context, you can abstract some repetitive tasks with reusable code that offers generic functionality to reduce the time and effort to automate your use cases. This is what automation frameworks offer, which we will discuss in the next chapter.

Just Imagine

# Further reading

To learn more about the topics that were covered in this chapter, take a look at the following resources:

-   NetDevOps 2020 Survey: [https://dgarros.github.io/netdevops-survey/reports/2020](https://dgarros.github.io/netdevops-survey/reports/2020)
-   `topo` directory: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/topo-base/topo.yml](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/topo-base/topo.yml)
-   `ch06/ssh` folder: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch06/ssh](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch06/ssh)
-   `ch06/vssh` directory: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch06/vssh](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch06/vssh)
-   `ch06/http` directory: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch06/http](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch06/http)
-   demo.nautobot.com: [https://demo.nautobot.com/](https://demo.nautobot.com/)
-   `ch06/state` directory: https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch06/ssh
-   `ch05/closed-loop` example: https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch05/closed-loop/main.go#L138
-   ntc-templates: [https://github.com/networktocode/ntc-templates](https://github.com/networktocode/ntc-templates)