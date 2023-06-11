# Working with TCP/IP and WebSocket

This chapter teaches you how to work with the lower-level protocols of TCP/IP, which are TCP and UDP, with the help of the `net` package so that we can develop TCP/IP servers and clients. Additionally, this chapter illustrates how to develop servers and clients for the _WebSocket_ protocol, which is based on HTTP, as well as UNIX domain sockets, for programming services that work on the local machine only.

In more detail, this chapter covers:

-   TCP/IP
-   The `net` package
-   Developing a TCP client
-   Developing a TCP server
-   Developing a UDP client
-   Developing a UDP server
-   Developing concurrent TCP servers
-   Working with UNIX domain sockets
-   Creating a WebSocket server
-   Creating a WebSocket client

Just Imagine

# TCP/IP

TCP/IP is a family of protocols that help the internet operate. Its name comes from its two most well-known protocols: TCP and IP. TCP stands for _Transmission Control Protocol_. TCP software transmits data between machines using segments, which are also called _TCP packets_. The main characteristic of TCP is that it is a reliable protocol, which means that it makes sure that a packet was delivered without requiring any extra code from the programmer. If there is no proof of packet delivery, TCP resends that particular packet. Among other things, TCP packets can be used for establishing connections, transferring data, sending acknowledgments, and closing connections.

When a TCP connection is established between two machines, a full-duplex virtual circuit, similar to a telephone call, is created between those two machines. The two machines constantly communicate to make sure that data is sent and received correctly. If the connection fails for some reason, the two machines try to find the problem and report to the relevant application. The TCP header of each packet includes the **source port** and **destination port** fields. These two fields, plus the source and destination IP addresses, are combined to uniquely identify every single TCP connection. All these details are handled by TCP/IP as long as you provide the required details without any extra effort.

When creating TCP/IP server processes, remember that port numbers **0-1024** have restricted access and can only be used by the root user, which means that you need administrative privileges to use a port in that range. Running a process with root privileges is a security risk and must be avoided.

IP stands for _Internet Protocol_. The main characteristic of IP is that it is not a reliable protocol by nature. IP encapsulates the data that travels over a TCP/IP network because it is responsible for delivering packets from the source host to the destination host according to the IP addresses. IP must find an addressing method for sending a packet to its destination effectively. Although there are dedicated devices, called routers, that perform IP routing, every TCP/IP device has to perform some basic routing. The first version of the IP protocol is now called **IPv4** to differentiate it from the latest version of the IP protocol, which is called **IPv6**. The main problem with IPv4 is that it is about to run out of available IP addresses, which is the main reason for creating the IPv6 protocol. This happened because an IPv4 address is represented using 32 bits only, which allows a total number of 2<sup class="Superscript--PACKT-">32</sup> (4,294,967,296) different IP addresses. On the other hand, IPv6 uses 128 bits to define each one of its addresses. The format of an IPv4 address is `10.20.32.245` (four parts separated by dots), while the format of an IPv6 address is `3fce:1706:4523:3:150:f8ff:fe21:56cf` (eight parts separated by colons).

**UDP** (**User Datagram Protocol**) is based on IP, which means that it is also unreliable. UDP is simpler than TCP, mainly because UDP is not reliable by design. As a result, UDP messages can be lost, duplicated, or arrive out of order. Furthermore, packets can arrive faster than the recipient can process them. So, UDP is used when speed is more important than reliability.

This chapter implements both TCP and UDP software—TCP and UDP services are the basis of the internet, and it is handy to know how to develop TCP/IP servers and clients in Go. But first, let us talk about the `nc(1)` utility.

## The nc(1) command-line utility

The `nc(1)` utility, which is also called `netcat(1)`, comes in very handy when you want to test TCP/IP servers and clients. Actually, `nc(1)` is a utility for everything that involves TCP and UDP as well as IPv4 and IPv6, including opening TCP connections, sending and receiving UDP messages, and acting as a TCP server.

You can use `nc(1)` as a client for a TCP service that runs on a machine with the `10.10.1.123` IP address and listens to port number `1234`, as follows:

```markup
$ nc 10.10.1.123 1234
```

The `-l` option tells `netcat(1)` to act as a server, which means that `netcat(1)` starts listening for incoming connections at the given port number. By default, `nc(1)` uses the TCP protocol. However, if you execute `nc(1)` with the `-u` flag, then `nc(1)` uses the UDP protocol, either as a client or as a server. Finally, the `-v` and `-vv` options tell `netcat(1)` to generate verbose output, which can come in handy when you want to troubleshoot network connections.

Just Imagine

# The net package

The `net` package of the Go Standard Library is all about TCP/IP, UDP, domain name resolution, and UNIX domain sockets. The `net.Dial()` function is used to connect to a network as a client, whereas the `net.Listen()` function is used to tell a Go program to accept incoming network connections and thus act as a server. The return value of both `net.Dial()` and `net.Listen()` is of the `net.Conn` data type, which implements the `io.Reader` and `io.Writer` interfaces—this means that you can both read and write to a `net.Conn` connection using code related to file I/O. The first parameter of both `net.Dial()`and `net.Listen()` is the network type, but this is where their similarities end.

The `net.Dial()` function is used for connecting to a remote server. The first parameter of the `net.Dial()` function defines the network protocol that is going to be used, while the second parameter defines the server address, which must also include the port number. Valid values for the first parameter are `tcp`, `tcp4` (IPv4-only), `tcp6` (IPv6-only), `udp`, `udp4` (IPv4-only), `udp6` (IPv6-only), `ip`, `ip4` (IPv4-only), `ip6` (IPv6-only), `unix` (UNIX sockets), `unixgram`, and `unixpacket`. On the other hand, valid values for `net.Listen()` are `tcp`, `tcp4`, `tcp6`, `unix`, and `unixpacket`.

Execute the `go doc net.Listen` and `go doc net.Dial` commands for more detailed information regarding these two functions.

Just Imagine

# Developing a TCP client

This section presents two equivalent ways of developing TCP clients.

## Developing a TCP client with net.Dial()

First, we are going to present the most widely used way, which is implemented in `tcpC.go`:

```markup
package main
import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)
```

The `import` block contains packages such as `bufio` and `fmt` that also work with file I/O operations.

```markup
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide host:port.")
        return
    }
```

First, we read the details of the TCP server we want to connect to.

```markup
    connect := arguments[1]
    c, err := net.Dial("tcp", connect)
    if err != nil {
        fmt.Println(err)
        return
    }
```

With the connection details, we call `net.Dial()`—its first parameter is the protocol we want to use, which in this case is `tcp`, and its second parameter is the connection details. A successful `net.Dial()` call returns an open connection (a `net.Conn` interface), which is a generic stream-oriented network connection.

```markup
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">> ")
        text, _ := reader.ReadString('\n')
        fmt.Fprintf(c, text+"\n")
        message, _ := bufio.NewReader(c).ReadString('\n')
        fmt.Print("->: " + message)
        if strings.TrimSpace(string(text)) == "STOP" {
            fmt.Println("TCP client exiting...")
            return
        }
    }
}
```

The last part of the TCP client keeps reading user input until the word `STOP` is given as input—in this case, the client **waits for the server response before terminating** after `STOP` because this is how the `for` loop is constructed. This mainly happens because the server might have a useful answer for us, and we do not want to miss that. All given user input is sent (written) to the open TCP connection using `fmt.Fprintf()`, whereas `bufio.NewReader()` is used for reading data from the TCP connection, just like you would do with a regular file.

Using `tcpC.go` to connect to a TCP server, which in this case is implemented with `nc(1)`, produces the next kind of output:

```markup
$ go run tcpC.go localhost:1234
>> Hello!
->: Hi from nc -l 1234
>> STOP
->: Bye!
TCP client exiting...
```

Lines beginning with `>>` denote user input, whereas lines beginning with `->` signify server messages. After sending `STOP`, we wait for the server response and then the client ends the TCP connection. The previous code demonstrates how to create a proper TCP client in Go with some extra logic in it (the `STOP` keyword).

The next subsection shows a different way of creating a TCP client.

## Developing a TCP client that uses net.DialTCP()

This subsection presents an alternative way to develop a TCP client. The difference lies in the Go functions that are being used for establishing the TCP connection, which are `net.DialTCP()` and `net.ResolveTCPAddr()`, and not in the functionality of the client.

The code of `otherTCPclient.go` is as follows:

```markup
package main
import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)
```

Although we are working with TCP/IP connections, we need packages such as `bufio` because UNIX treats network connections as files, so we are basically working with I/O operations over networks.

```markup
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide a server:port string!")
        return
    }
```

We need to read the details of the TCP server we want to connect to, including the desired port number. The utility cannot operate with default parameters when working with TCP/IP unless we are developing a very specialized TCP client.

```markup
    connect := arguments[1]
    tcpAddr, err := net.ResolveTCPAddr("tcp4", connect)
    if err != nil {
        fmt.Println("ResolveTCPAddr:", err)
        return
    }
```

The `net.ResolveTCPAddr()` function is specific to TCP connections, hence its name, and resolves the given address to a `*net.TCPAddr` value, which is a structure that represents the address of a TCP endpoint—in this case, the endpoint is the TCP server we want to connect to.

```markup
    conn, err := net.DialTCP("tcp4", nil, tcpAddr)
    if err != nil {
        fmt.Println("DialTCP:", err)
        return
    }
```

With the TCP endpoint at hand, we call `net.DialTCP()` to connect to the server. Apart from the use of `net.ResolveTCPAddr()` and `net.DialTCP()`, the rest of the code that has to do with the TCP client and TCP server interaction is exactly the same.

```markup
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">> ")
        text, _ := reader.ReadString('\n')
        fmt.Fprintf(conn, text+"\n")
        message, _ := bufio.NewReader(conn).ReadString('\n')
        fmt.Print("->: " + message)
        if strings.TrimSpace(string(text)) == "STOP" {
            fmt.Println("TCP client exiting...")
            conn.Close()
            return
        }
    }
}
```

Lastly, an infinite `for` loop is used for interacting with the TCP server. The TCP client reads user data, which is sent to the server. After that, it reads data from the TCP server. Once again, the `STOP` keyword ends the TCP connection on the client side using the `Close()` method.

Working with `otherTCPclient.go` and interacting with a TCP server process produces the next kind of output:

```markup
$ go run otherTCPclient.go localhost:1234
>> Hello!
->: Hi from nc -l 1234
>> STOP
->: Thanks for connection!
TCP client exiting...
```

The interaction is the same as with `tcpC.go`—we have just learned a different way of developing TCP clients. If you want my opinion, I prefer the implementation found in `tcpC.go` because it uses more generic functions. However, this is just personal taste.

The next section shows how to program TCP servers.

Just Imagine

# Developing a TCP server

This section presents two ways of developing a TCP server that can interact with TCP clients, just as we did with the TCP client.

## Developing a TCP server with net.Listen()

The TCP server presented in this section, which uses `net.Listen()`, returns the current date and time to the client in a single network packet. In practice, this means that after accepting a client connection, the server gets the time and date from the operating system and sends that data back to the client. The `net.Listen()` function listens for connections, whereas the `net.Accept()` method waits for the next connection and returns a generic `Conn` variable with the client information. The code of `tcpS.go` is as follows:

```markup
package main
import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
    "time"
)
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide port number")
        return
    }
```

The TCP server should know about the port number it is going to use—this is given as a command-line argument.

```markup
    PORT := ":" + arguments[1]
    l, err := net.Listen("tcp", PORT)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer l.Close()
```

The `net.Listen()` function listens for connections and is what makes that particular program a server process. If the second parameter of `net.Listen()` contains a port number without an IP address or a hostname, `net.Listen()` listens to all available IP addresses of the local system, which is the case here.

```markup
    c, err := l.Accept()
    if err != nil {
        fmt.Println(err)
        return
    }
```

We just call `Accept()` and wait for a client connection—`Accept()` blocks until a connection comes. There is something unusual with this particular TCP server: it can only serve the first TCP client that is going to connect to it because the `Accept()` call is outside of the `for` loop and therefore is called only once. Each individual client should be specified by a different `Accept()` call.

Correcting that is left as an exercise for the reader.

```markup
    for {
        netData, err := bufio.NewReader(c).ReadString('\n')
        if err != nil {
            fmt.Println(err)
            return
        }
        if strings.TrimSpace(string(netData)) == "STOP" {
            fmt.Println("Exiting TCP server!")
            return
        }
        fmt.Print("-> ", string(netData))
        t := time.Now()
        myTime := t.Format(time.RFC3339) + "\n"
        c.Write([]byte(myTime))
    }
}
```

This endless `for` loop keeps interacting with the same TCP client until the word `STOP` is sent from the client. As it happened with the TCP clients, `bufio.NewReader()` is used for reading data from the network connection, whereas `Write()` is used for sending data to the TCP client.

Running `tcpS.go` and interacting with a TCP client produces the next kind of output:

```markup
$ go run tcpS.go 1234
-> Hello!
-> Have to leave now!
EOF
```

The server connection ended automatically with the client connection because the `for` loop concluded when `bufio.NewReader(c).ReadString('\n')` had nothing more to read. The client was `nc(1)`, which produced the next output:

```markup
$ nc localhost 1234
Hello!
2021-04-12T08:53:32+03:00
Have to leave now!
2021-04-12T08:53:51+03:00
```

In order to exit `nc(1)`, we need to press Ctrl + D, which is `EOF` (End Of File) in UNIX.

So, we now know how to develop a TCP server in Go. As it happened with the TCP client, there is an alternative way to develop a TCP server, which is presented in the next subsection.

## Developing a TCP server that uses net.ListenTCP()

This time, this alternative version of the TCP server implements the echo service. Put simply, the TCP server sends back to the client the data that was received by the client.

The code of `otherTCPserver.go` is as follows:

```markup
package main
import (
    "fmt"
    "net"
    "os"
    "strings"
)
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide a port number!")
        return
    }
    SERVER := "localhost" + ":" + arguments[1]
    s, err := net.ResolveTCPAddr("tcp", SERVER)
    if err != nil {
        fmt.Println(err)
        return
    }
```

The previous code gets the TCP port number value as a command-line argument, which is used in `net.ResolveTCPAddr()`—this is required to define the TCP port number the TCP server is going to listen to.

That function only works with TCP, hence its name.

```markup
    l, err := net.ListenTCP("tcp", s)
    if err != nil {
        fmt.Println(err)
        return
    }
```

Similarly, `net.ListenTCP()` only works with TCP and is what makes that program a TCP server ready to accept incoming connections.

```markup
    buffer := make([]byte, 1024)
    conn, err := l.Accept()
    if err != nil {
        fmt.Println(err)
        return
    }
```

As before, due to the place where `Accept()` is called, this particular implementation can work with a single client only. This is used for reasons of simplicity. The concurrent TCP server that is developed later on in this chapter puts the `Accept()` call inside the endless `for` loop.

```markup
    for {
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println(err)
            return
        }
        if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
            fmt.Println("Exiting TCP server!")
            conn.Close()
            return
        }
```

You need to use `strings.TrimSpace()` in order to remove any space characters from your input and compare the result with `STOP`, which has a special meaning in this implementation. When the `STOP` keyword is received from the client, the server closes the connection using the `Close()` method.

```markup
        fmt.Print("> ", string(buffer[0:n-1]), "\n")
        _, err = conn.Write(buffer)
        if err != nil {
            fmt.Println(err)
            return
        }
    }
}
```

All previous code is for interacting with the TCP client until the client decides to close the connection.

Running `otherTCPserver.go` and interacting with a TCP client produces the next kind of output:

```markup
$ go run otherTCPserver.go 1234
> Hello from the client!
Exiting TCP server!
```

The first line that begins with `>` is the client message, whereas the second line is the server output when getting the `STOP` message from the client. Therefore, the TCP server processes client requests as programmed and exits when it gets the `STOP` message, which is the expected behavior.

The next section is about developing UDP clients.

Just Imagine

# Developing a UDP client

This section demonstrates how to develop a UDP client that can interact with UDP services. The code of `udpC.go` is as follows:

```markup
package main
import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide a host:port string")
        return
    }
    CONNECT := arguments[1]
```

This is how we get the UDP server details from the user.

```markup
    s, err := net.ResolveUDPAddr("udp4", CONNECT)
    c, err := net.DialUDP("udp4", nil, s)
```

The previous two lines declare that we are using UDP and that we want to connect to the UDP server that is specified by the return value of `net.ResolveUDPAddr()`. The actual connection is initiated using `net.DialUDP()`.

```markup
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
    defer c.Close()
```

This part of the program finds the details of the UDP server by calling the `RemoteAddr()` method.

```markup
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">> ")
        text, _ := reader.ReadString('\n')
        data := []byte(text + "\n")
        _, err = c.Write(data)
```

Data is read from the user using `bufio.NewReader(os.Stdin)` and is written to the UDP server using `Write()`.

```markup
        if strings.TrimSpace(string(data)) == "STOP" {
            fmt.Println("Exiting UDP client!")
            return
        }
```

If the input read from the user is the `STOP` keyword, then the connection is terminated.

```markup
        if err != nil {
            fmt.Println(err)
            return
        }
        buffer := make([]byte, 1024)
        n, _, err := c.ReadFromUDP(buffer)
```

Data is read from the UDP connection using the `ReadFromUDP()` method.

```markup
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Printf("Reply: %s\n", string(buffer[0:n]))
    }
}
```

The `for` loop is going to keep going forever until the `STOP` keyword is received as input or the program is terminated in some other way.

Working with `udpC.go` is as simple as follows—the client side is implemented using `nc(1)`:

```markup
$ go run udpC.go localhost:1234
The UDP server is 127.0.0.1:1234
```

`127.0.0.1:1234` is the value of `c.RemoteAddr().String()`, which shows the details of the UDP server we have connected to.

```markup
>> Hello!
Reply: Hi from the server
```

Our client sent `Hello!` to the UDP server and received `Hi from the server` back.

```markup
>> Have to leave now :)
Reply: OK - bye from nc -l -u 1234
```

Our client sent `Have to leave now :)` to the UDP server and received `OK - bye from nc -l -u 1234` back.

```markup
>> STOP
Exiting UDP client!
```

Finally, after sending the `STOP` keyword to the server, the client prints `Exiting UDP client!` and terminates—the message is defined in the Go code and can be anything you want.

The next section is about programming a UDP server.

Just Imagine

# Developing a UDP server

This section shows how to develop a UDP server, which generates and returns random numbers to its clients. The code for the UDP server (`udpS.go`) is as follows:

```markup
package main
import (
    "fmt"
    "math/rand"
    "net"
    "os"
    "strconv"
    "strings"
    "time"
)
func random(min, max int) int {
    return rand.Intn(max-min) + min
}
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide a port number!")
        return
    }
    PORT := ":" + arguments[1]
```

The UDP port number the server is going to listen to is provided as a command-line argument.

```markup
    s, err := net.ResolveUDPAddr("udp4", PORT)
    if err != nil {
        fmt.Println(err)
        return
    }
```

The `net.ResolveUDPAddr()` function creates a UDP endpoint that is going to be used to create the server.

```markup
    connection, err := net.ListenUDP("udp4", s)
    if err != nil {
        fmt.Println(err)
        return
    }
```

The `net.ListenUDP("udp4", s)` function call makes this process a server for the `udp4` protocol using the details specified by its second parameter.

```markup
    defer connection.Close()
    buffer := make([]byte, 1024)
```

The `buffer` variable stores a `byte` slice and is used to read data from the UDP connection.

```markup
    rand.Seed(time.Now().Unix())
    for {
        n, addr, err := connection.ReadFromUDP(buffer)
        fmt.Print("-> ", string(buffer[0:n-1]))
```

The `ReadFromUDP()` and `WriteToUDP()` methods are used to read data from a UDP connection and write data to a UDP connection, respectively. Additionally, due to the way UDP operates, the UDP server can serve multiple clients.

```markup
        if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
            fmt.Println("Exiting UDP server!")
            return
        }
```

The UDP server terminates when any one of the clients sends the `STOP` message. Aside from this, the `for` loop is going to keep running forever.

```markup
        data := []byte(strconv.Itoa(random(1, 1001)))
        fmt.Printf("data: %s\n", string(data))
```

A `byte` slice is stored in the `data` variable and used to write the desired data to the client.

```markup
        _, err = connection.WriteToUDP(data, addr)
        if err != nil {
            fmt.Println(err)
            return
        }
    }
}
```

Working with `udpS.go` is as simple as the following:

```markup
$ go run udpS.go 1234
-> Hello from client!
data: 395
```

Lines beginning with `->` show data coming from a client. Lines beginning with `data` show random numbers generated by the UDP server—in this case, `395`.

```markup
-> Going to terminate the connection now.
data: 499
```

The previous two lines show another interaction with a UDP client.

```markup
-> STOP
Exiting UDP server!
```

Once the UDP server receives the `STOP` keyword from the client, it closes the connection and exits.

On the client side, which uses `udpC.go`, we have the next interaction:

```markup
$ go run udpC.go localhost:1234
The UDP server is 127.0.0.1:1234
>> Hello from client!
Reply: 395
```

The client sends the `Hello from client!` message to the server and receives `395`.

```markup
>> Going to terminate the connection now.
Reply: 499
```

The client sends `Going to terminate the connection now.` to the server and receives the `499` random number.

```markup
>> STOP
Exiting UDP client!
```

When the client gets `STOP` as user input, it terminates the UDP connection and exits.

The next section shows how to develop a concurrent TCP server that uses goroutines for serving its clients.

Just Imagine

# Developing concurrent TCP servers

This section teaches a pattern for developing concurrent TCP servers, which are servers that are using separate goroutines to serve their clients following a successful `Accept()` call. Therefore, such servers can serve multiple TCP clients at the same time. This is how real-world production servers and services are implemented.

The code of `concTCP.go` is as follows:

```markup
package main
import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
)
var count = 0
func handleConnection(c net.Conn) {
    fmt.Print(".")
```

The previous statement is not required—it just informs us that a new client has been connected.

```markup
    for {
        netData, err := bufio.NewReader(c).ReadString('\n')
        if err != nil {
            fmt.Println(err)
            return
        }
        temp := strings.TrimSpace(string(netData))
        if temp == "STOP" {
            break
        }
        fmt.Println(temp)
        counter := "Client number: " + strconv.Itoa(count) + "\n"
        c.Write([]byte(string(counter)))
    }
```

The `for` loop makes sure that `handleConnection()` is not going to exit automatically. Once again, the `STOP` keyword stops the current client connection—however, the server process, as well as all other active client connections, are going to keep running.

```markup
    c.Close()
}
```

This is the end of the function that is executed as a goroutine to serve clients. All you need in order to serve a client is a `net.Conn` parameter with the client details. After reading client data, the server sends back to the client a message indicating the number of the client that is being served so far.

```markup
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide a port number!")
        return
    }
    PORT := ":" + arguments[1]
    l, err := net.Listen("tcp4", PORT)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer l.Close()
    for {
        c, err := l.Accept()
        if err != nil {
            fmt.Println(err)
            return
        }
        go handleConnection(c)
        count++
    }
}
```

Each time a new client connects to the server, the `count` variable is increased. Each TCP client is served by a separate goroutine that executes the `handleConnection()` function. This frees the server process and allows it to accept new connections. Put simply, while multiple TCP clients are served, the TCP server is free to interact with more TCP clients. As before, new TCP clients are connected using the `Accept()` function.

Working with `concTCP.go` produces the next kind of output:

```markup
$ go run concTCP.go 1234
.Hello
.Hi from  nc localhost 1234
```

The first line of output is from the first TCP client, whereas the second line is from the second TCP client. This means that the concurrent TCP server works as expected. Therefore, when you want to be able to serve multiple TCP clients in your TCP services, you can use the technique and code presented as a template for developing your TCP servers.

The next section shows how to work with UNIX domain sockets, which are really fast for interactions on the local machine only.

Just Imagine

# Working with UNIX domain sockets

A _UNIX Domain Socket_ or Inter-Process Communication (IPC) socket is a data communications endpoint that allows you to exchange data between processes that run **on the same machine**. You might ask, why use UNIX domain sockets instead of TCP/IP connections for processes that exchange data on the same machine? First, because UNIX domain sockets are faster than TCP/IP connections and second, because UNIX domain sockets require fewer resources than TCP/IP connections. So, you can use UNIX domain sockets when both the clients and the server are on the same machine.

## A UNIX domain socket server

This section illustrates how to develop a UNIX domain socket server. Although we do not have to deal with TCP ports and network connections, the code presented is very similar to the code of the TCP server as found in `tcpS.go` and `concTCP.go`. The presented server implements the echo service.

The source code of `socketServer.go` is as follows:

```markup
package main
import (
    "fmt"
    "net"
    "os"
)
func echo(c net.Conn) {
```

The `echo()` function is used for serving client requests, hence the use of the `net.Conn` parameter that holds the client's details:

```markup
    for {
        buf := make([]byte, 128)
        n, err := c.Read(buf)
        if err != nil {
            fmt.Println("Read:", err)
            return
        }
```

We read data from the socket connection using `Read()` inside a `for` loop.

```markup
        data := buf[0:n]
        fmt.Print("Server got: ", string(data))
        _, err = c.Write(data)
        if err != nil {
            fmt.Println("Write:", err)
            return
        }
    }
}
```

In this second part of `echo()`, we send back to the client the data that the client sent to us. The `buf[0:n]` notation makes sure that we are going to send back the same amount of data that was read even if the size of the buffer is bigger.

This function serves all client connections—as you are going to see in a while, it is executed as a goroutine, which is the main reason that it does not return any values.

You cannot tell whether this function serves TCP/IP connections or UNIX socket domain connections, which mainly happens because UNIX treats all connections as files.

```markup
func main() {
    if len(os.Args) == 1 {
        fmt.Println("Need socket path")
        return
    }
    socketPath := os.Args[1]
```

This is the part where we specify the socket file that is going to be used by the server and its clients. In this case, the path to the socket file is given as a command-line argument.

```markup
    _, err := os.Stat(socketPath)
    if err == nil {
        fmt.Println("Deleting existing", socketPath)
        err := os.Remove(socketPath)
        if err != nil {
            fmt.Println(err)
            return
        }
    }
```

If the socket file already exists, you should **delete it** before the program continues—`net.Listen()` creates that file again.

```markup
    l, err := net.Listen("unix", socketPath)
    if err != nil {
        fmt.Println("listen error:", err)
        return
    }
```

What makes this a UNIX domain socket server is the use of `net.Listen()` with the `"unix"` parameter. In this case, we need to provide `net.Listen()` with the path of the socket file.

```markup
    for {
        fd, err := l.Accept()
        if err != nil {
            fmt.Println("Accept error:", err)
            return
        }
        go echo(fd)
    }
}
```

Each client connection is handled by a goroutine—in this sense, this is a concurrent UNIX domain socket server that can work with multiple clients! So, if you need to serve thousands of domain socket clients on a production server, this is the way to go!

In the next section, we are going to see the server in action as it is going to interact with the UNIX domain socket client that we are going to create.

## A UNIX domain socket client

This subsection shows a UNIX domain socket client implementation, which can be used to communicate with a domain socket server, such as the one developed in the previous subsection. The relevant code can be found in `socketClient.go`:

```markup
package main
import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
    "time"
)
func main() {
    if len(os.Args) == 1 {
        fmt.Println("Need socket path")
        return
    }
    socketPath := os.Args[1]
```

This is the part where we get from the user the socket file that is going to be used—the socket file should already exist and be handled by the UNIX domain socket server.

```markup
    c, err := net.Dial("unix", socketPath)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer c.Close()
```

The `net.Dial()` function is used for connecting to the socket.

```markup
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">> ")
        text, _ := reader.ReadString('\n')
        _, err = c.Write([]byte(text))
```

Here, we convert user input into a `byte` slice and send it to the server using `Write()`.

```markup
        if err != nil {
            fmt.Println("Write:", err)
            break
        }
        buf := make([]byte, 256)
        n, err := c.Read(buf[:])
        if err != nil {
            fmt.Println(err, n)
            return
        }
        fmt.Print("Read:", string(buf[0:n]))
```

This `fmt.Print()` statement prints as many characters from the `buf` slice as the number of characters read from the `Read()` method using the `buf[0:n]` notation.

```markup
        if strings.TrimSpace(string(text)) == "STOP" {
            fmt.Println("Exiting UNIX domain socket client!")
            return
        }
```

If the word `STOP` is given as input, then the client returns and therefore closes the connection to the server. Generally speaking, it is always good to have a way of gracefully exiting such a utility.

```markup
        time.Sleep(5 * time.Second)
    }
}
```

The `time.Sleep()` call is used to delay the `for` loop and emulate the operation of a real program.

Working with both `socketServer.go` and `socketClient.go`, provided that the server is executed first, generates the next kind of output:

```markup
$ go run socketServer.go /tmp/packt.socket
Server got: Hello!
Server got: STOP
Read: EOF
```

Although the client connection ended, the server continues to run and waits for more client requests.

On the client side, we have the following:

```markup
$ go run socketClient.go /tmp/packt.socket
>> Hello!
Read: Hello!
>> STOP
Read: STOP
Exiting UNIX domain socket client!
```

In the previous two sections, we learned how to create UNIX domain socket clients and servers that are faster than TCP/IP servers but work on the same machine only.

The sections that follow concern the WebSocket protocol.

Just Imagine

# Creating a WebSocket server

The WebSocket protocol is a computer communications protocol that provides **full-duplex** (transmission of data in two directions simultaneously) communication channels over a single TCP connection. The WebSocket protocol is defined in RFC 6455 ([https://tools.ietf.org/html/rfc6455](https://tools.ietf.org/html/rfc6455)) and uses `ws://` and `wss://` instead of `http://` and `https://`, respectively. Therefore, the client should begin a WebSocket connection by using a URL that starts with `ws://`.

In this section, we are going to develop a small yet fully functional WebSocket server using the `gorilla/websocket` ([https://github.com/gorilla/websocket](https://github.com/gorilla/websocket)) module. The server implements the **Echo service**, which means that it automatically returns its input to the client.

The `golang.org/x/net/websocket` package offers another way of developing WebSocket clients and servers. However, according to its documentation, [golang.org/x/net/websocket](http://golang.org/x/net/websocket) lacks some features and it is advised that you use [https://godoc.org/github.com/gorilla/websocket](https://godoc.org/github.com/gorilla/websocket), the one used here, or [https://godoc.org/nhooyr.io/websocket](https://godoc.org/nhooyr.io/websocket) instead.

The advantages of the WebSocket protocol include the following:

-   A WebSocket connection is a full-duplex, bidirectional communications channel. This means that a server does not need to wait to read from a client to send data to the client and vice versa.
-   WebSocket connections are raw TCP sockets, which means that they do not have the overhead required to establish an HTTP connection.
-   WebSocket connections can also be used for sending HTTP data. However, plain HTTP connections cannot work as WebSocket connections.
-   WebSocket connections live until they are killed, so there is no need to reopen them all the time.
-   WebSocket connections can be used for real-time web applications.
-   Data can be sent from the server to the client at any time, without the client even requesting it.
-   WebSocket is part of the HTML5 specification, which means that it is supported by all modern web browsers.

Before showing the server implementation, it would be good for you to know that the `websocket.Upgrader` method of the `gorilla/websocket` package **upgrades** an HTTP server connection to the WebSocket protocol and allows you to define the parameters of the upgrade. After that, your HTTP connection is a WebSocket connection, which means that you will not be allowed to execute statements that work with the HTTP protocol.

The next subsection shows the implementation of the server.

## The implementation of the server

This subsection presents the implementation of the WebSocket server that implements the `echo` service, which can be really handy when testing network connections.

The GitHub repository used for keeping the code can be found at [https://github.com/mactsouk/ws](https://github.com/mactsouk/ws). If you want to follow along with this section, you should download that repository and put it inside `~/go/src`—in my case, it was put inside `~/go/src/github.com/mactsouk`, in the `ws` folder.

The GitHub repository contains a `Dockerfile` file for producing a Docker image from the WebSocket server source file.

The implementation of the WebSocket server can be found in `ws.go`, which contains the next code:

```markup
package main
import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    "github.com/gorilla/websocket"
)
```

This is the external package used for working with the WebSocket protocol.

```markup
var PORT = ":1234"
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
```

This is where the parameters of `websocket.Upgrader` are defined. They are going to be used shortly.

```markup
func rootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome!\n")
    fmt.Fprintf(w, "Please use /ws for WebSocket!")
}
```

This is a regular HTTP handler function.

```markup
func wsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Connection from:", r.Host)
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("upgrader.Upgrade:", err)
        return
    }
    defer ws.Close()
```

A WebSocket server application calls the `Upgrader.Upgrade` method in order to get a WebSocket connection from an HTTP request handler. Following a successful call to `Upgrader.Upgrade`, the server begins working with the WebSocket connection and the WebSocket client.

```markup
    for {
        mt, message, err := ws.ReadMessage()
        if err != nil {
            log.Println("From", r.Host, "read", err)
            break
        }
        log.Print("Received: ", string(message))
        err = ws.WriteMessage(mt, message)
        if err != nil {
            log.Println("WriteMessage:", err)
            break
        }
    }
}
```

The `for` loop in `wsHandler()` handles all incoming messages for `/ws`—you can use any technique you want. Additionally, in the presented implementation, **only the client** is allowed to close an existing WebSocket connection unless there is a network issue, or the server process is killed.

Last, remember that in a WebSocket connection, you cannot use `fmt.Fprintf()` statements for sending data to the WebSocket client—if you use any of these, or any other call that can implement the same functionality, the WebSocket connection fails and you are not going to be able to send or receive any data. Therefore, the only way to send and receive data in a WebSocket connection implemented with `gorilla/websocket` is through `WriteMessage()` and `ReadMessage()` calls, respectively. Of course, you can always implement the desired functionality on your own by working with raw network data, but implementing this goes beyond the scope of this book.

```markup
func main() {
    arguments := os.Args
    if len(arguments) != 1 {
        PORT = ":" + arguments[1]
    }
```

If there is not a command-line argument, use the default port number stored in `PORT`.

```markup
    mux := http.NewServeMux()
    s := &http.Server{
        Addr:         PORT,
        Handler:      mux,
        IdleTimeout:  10 * time.Second,
        ReadTimeout:  time.Second,
        WriteTimeout: time.Second,
    }
```

These are the details of the HTTP server that also handles WebSocket connections.

```markup
    mux.Handle("/", http.HandlerFunc(rootHandler))
    mux.Handle("/ws", http.HandlerFunc(wsHandler))
```

The endpoint used for WebSocket can be anything you want—in this case, it is `/ws`. Additionally, you can have multiple endpoints that work with the WebSocket protocol.

```markup
    log.Println("Listening to TCP Port", PORT)
    err := s.ListenAndServe()
    if err != nil {
        log.Println(err)
        return
    }
}
```

The code presented uses `log.Println()` instead of `fmt.Println()` for printing messages—as this is a server process, using `log.Println()` is a much better choice than `fmt.Println()` because logging information is sent to files that can be examined at a later time. However, during development, you might prefer `fmt.Println()` calls and avoid writing to your log files because you can see your data on screen immediately without having to look elsewhere.

The server implementation is short, yet fully functional. The single most important call in the code is `Upgrader.Upgrade` because this is what upgrades an HTTP connection to a WebSocket connection.

Getting and running the code from GitHub requires the following steps—most of the steps have to do with module initialization and downloading the required packages:

```markup
$ cd ~/go/src/github.com/mactsouk/
$ git clone https://github.com/mactsouk/ws.git
$ cd ws
$ go mod init
$ go mod tidy
$ go mod download
$ go run ws.go
```

To test that server, you need to have a client. As we have not developed our own client so far, we are going to test the WebSocket server using two other means.

### Using websocat

`websocat` is a command-line utility that can help you test WebSocket connections. However, as `websocat` is not installed by default, you need to install it on your machine using your package manager of choice. You can use it as follows, provided that there is a WebSocket server at the desired address:

```markup
$ websocat ws://localhost:1234/ws
Hello from websocat!
```

This is what we type and send to the server.

```markup
Hello from websocat!
```

This is what we get back from the WebSocket server, which implements the `echo` service—different WebSocket servers implement different functionality.

```markup
Bye!
```

Again, the previous line is user input given to `websocat`.

```markup
Bye!
```

And the last line is the data sent back from the server. The connection was closed by pressing Ctrl + D on the `websocat` client.

Should you wish for a more verbose output from `websocat`, you can execute it with the `-v` flag:

```markup
$ websocat -v ws://localhost:1234/ws
[INFO  websocat::lints] Auto-inserting the line mode
[INFO  websocat::stdio_threaded_peer] get_stdio_peer (threaded)
[INFO  websocat::ws_client_peer] get_ws_client_peer
[INFO  websocat::ws_client_peer] Connected to ws
Hello from websocat!
Hello from websocat!
Bye!
Bye!
[INFO  websocat::sessionserve] Forward finished
[INFO  websocat::ws_peer] Received WebSocket close message
[INFO  websocat::sessionserve] Reverse finished
[INFO  websocat::sessionserve] Both directions finished
```

In both cases, the output from our WebSocket server should be similar to the following:

```markup
$ go run ws.go
2021/04/10 20:54:30 Listening to TCP Port :1234
2021/04/10 20:54:42 Connection from: localhost:1234
2021/04/10 20:54:57 Received: Hello from websocat!
2021/04/10 20:55:03 Received: Bye!
2021/04/10 20:55:03 From localhost:1234 read websocket: close 1005 (no status)
```

The next subsection illustrates how to test the WebSocket server using HTML and JavaScript.

### Using JavaScript

The second way of testing the WebSocket server is by creating a web page with some HTML and JavaScript code. This technique gives you more control over what is happening, but requires more code and familiarity with HTML and JavaScript.

The HTML page with the JavaScript code that makes it act like a WebSocket client is the following:

```markup
<!DOCTYPE html>
<meta charset="utf-8">
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Testing a WebSocket Server</title>
  </head>
  <body>
    <h2>Hello There!</h2>
    <script>
        let ws = new WebSocket("ws://localhost:1234/ws");
```

This is the single most important JavaScript statement because this is where you specify the address of the WebSocket server, the port number, and the endpoint you want to connect to.

```markup
        console.log("Trying to connect to server.");
        ws.onopen = () => {
            console.log("Connected!");
            ws.send("Hello From the Client!")
        };
```

The `ws.onopen` event is used for making sure that the WebSocket connection is open, whereas the `send()` method is used for sending messages to the WebSocket server.

```markup
        ws.onmessage = function(event) {
          console.log(`[message] Data received from server: ${event.data}`);
          ws.close(1000, "Work complete");
        };
```

The `onmessage` event is triggered each time the WebSocket server sends a new message—however, in our case, the connection is closed as soon as the first message from the server is received.

Lastly, the `close()` JavaScript method is used for closing a WebSocket connection—in our case, the `close()` call is included in the `onmessage` event. Calling `close()` triggers the `onclose` event, which contains the code that follows:

```markup
        ws.onclose = event => {
            if (event.wasClean) {
              console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
            }
            console.log("Socket Closed Connection: ", event);
        };
        ws.onerror = error => {
            console.log("Socket Error: ", error);
        };
    </script>
  </body>
</html>
```

You can see the output of the JavaScript code by visiting the JavaScript console on your favorite web browser, which in this case is Google Chrome. The following screenshot shows the generated output.

![Text
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_09_01.png)

Figure 9.1: Interacting with the WebSocket server using JavaScript

For the WebSocket interaction defined in `test.html`, the WebSocket server generated the following output in the command line:

```markup
2021/04/10 21:43:22 Connection from: localhost:1234
2021/04/10 21:43:22 Received: Hello From the Client!
2021/04/10 21:43:22 From localhost:1234 read websocket: close 1000 (normal): Work complete
```

Both ways verify that the WebSocket server works as expected: the client can connect to the server, the server sends data that is received by the client, and the client closes the connection with the server successfully. So, it is time to develop our own WebSocket client in Go.

Just Imagine

# Creating a WebSocket client

This subsection shows how to program a WebSocket client in Go. The client reads user data that sends it to the server and reads the server response. The `client` directory of [https://github.com/mactsouk/ws](https://github.com/mactsouk/ws) contains the implementation of the WebSocket client—I find it more convenient to include both implementations in the same repository.

As with the WebSocket server, the `gorilla/websocket` package is going to help us develop the WebSocket client.

We are going to see `gorilla` in the next chapter when working with RESTful services.

The code of `./client/client.go` is as follows:

```markup
package main
import (
    "bufio"
    "fmt"
    "log"
    "net/url"
    "os"
    "os/signal"
    "syscall"
    "time"
    "github.com/gorilla/websocket"
)
var SERVER = ""
var PATH = ""
var TIMESWAIT = 0
var TIMESWAITMAX = 5
var in = bufio.NewReader(os.Stdin)
```

The `in` variable is just a shortcut for `bufio.NewReader(os.Stdin)`.

```markup
func getInput(input chan string) {
    result, err := in.ReadString('\n')
    if err != nil {
        log.Println(err)
        return
    }
    input <- result
}
```

The `getInput()` function, which is executed as a goroutine, gets user input that is transferred to the `main()` function via the `input` channel. Each time the program reads some user input, the old goroutine ends and a new `getInput()` goroutine begins in order to get new input.

```markup
func main() {
    arguments := os.Args
    if len(arguments) != 3 {
        fmt.Println("Need SERVER + PATH!")
        return
    }
    SERVER = arguments[1]
    PATH = arguments[2]
    fmt.Println("Connecting to:", SERVER, "at", PATH)
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)
```

The WebSocket client handles UNIX interrupts with the help of the `interrupt` channel. When the appropriate signal is caught (`syscall.SIGINT`), the WebSocket connection with the server is closed with the help of the `websocket.CloseMessage` message. This is how professional tools work!

```markup
    input := make(chan string, 1)
    go getInput(input)
    URL := url.URL{Scheme: "ws", Host: SERVER, Path: PATH}
    c, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
    if err != nil {
        log.Println("Error:", err)
        return
    }
    defer c.Close()
```

The WebSocket connection begins with a call to `websocket.DefaultDialer.Dial()`. Everything that goes to the `input` channel is transferred to the WebSocket server using the `WriteMessage()` method.

```markup
    done := make(chan struct{})
    go func() {
        defer close(done)
        for {
            _, message, err := c.ReadMessage()
            if err != nil {
                log.Println("ReadMessage() error:", err)
                return
            }
            log.Printf("Received: %s", message)
        }
    }()
```

Another goroutine, which this time is implemented using an anonymous Go function, is responsible for reading data from the WebSocket connection using the `ReadMessage()` method.

```markup
    for {
        select {
        case <-time.After(4 * time.Second):
            log.Println("Please give me input!", TIMESWAIT)
            TIMESWAIT++
            if TIMESWAIT > TIMESWAITMAX {
                syscall.Kill(syscall.Getpid(), syscall.SIGINT)
            }
```

The `syscall.Kill(syscall.Getpid(), syscall.SIGINT)` statement sends the interrupt signal to the program using Go code. According to the logic of `client.go`, the interrupt signal makes the program close the WebSocket connection with the server and terminate its execution. This only happens if the current number of timeout periods is bigger than a predefined global value.

```markup
        case <-done:
            return
        case t := <-input:
            err := c.WriteMessage(websocket.TextMessage, []byte(t))
            if err != nil {
                log.Println("Write error:", err)
                return
            }
            TIMESWAIT = 0
```

If you get user input, the current number of the timeout periods (`TIMESWAIT`) is reset and new input is read.

```markup
            go getInput(input)
        case <-interrupt:
            log.Println("Caught interrupt signal - quitting!")
            err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
```

Just before we close the client connection, we send `websocket.CloseMessage` to the server in order to do the closing the right way.

```markup
            if err != nil {
                log.Println("Write close error:", err)
                return
            }
            select {
            case <-done:
            case <-time.After(2 * time.Second):
            }
            return
        }
    }
}
```

As `./client/client.go` is in a separate directory to `ws.go`, we need to run the next commands in order to collect the required dependencies and run it:

```markup
$ cd client
$ go mod init
$ go mod tidy
$ go mod download
```

Interacting with the WebSocket server produces the next kind of output:

```markup
$ go run client.go localhost:1234 ws
Connecting to: localhost:1234 at ws
Hello there!
2021/04/10 21:30:33 Received: Hello there!
```

The previous two lines show user input as well as the server response.

```markup
2021/04/10 21:30:37 Please give me input! 0
2021/04/10 21:30:41 Please give me input! 1
2021/04/10 21:30:45 Please give me input! 2
2021/04/10 21:30:49 Please give me input! 3
2021/04/10 21:30:53 Please give me input! 4
2021/04/10 21:30:57 Please give me input! 5
2021/04/10 21:30:57 Caught interrupt signal - quitting!
2021/04/10 21:30:57 ReadMessage() error: websocket: close 1000 (normal)
```

The last lines show how the automatic timeout process works.

The WebSocket server generated the following output for the previous interaction:

```markup
2021/04/10 21:30:29 Connection from: localhost:1234
2021/04/10 21:30:33 Received: Hello there!
2021/04/10 21:30:57 From localhost:1234 read websocket: close 1000 (normal)
```

However, if a WebSocket server cannot be found at the address provided, the WebSocket client produces the next output:

```markup
$ go run client.go localhost:1234 ws
Connecting to: localhost:1234 at ws
2021/04/09 10:29:23 Error: dial tcp [::1]:1234: connect: connection refused
```

The `connection refused` message indicates that there is no process listening to port `1234` at `localhost`.

WebSocket gives you an alternative way of creating services. As a rule of thumb, WebSocket is better when we want to exchange lots of data, and we want the connection to remain open all the time and exchange data in full-duplex. However, if you are not sure about what to use, begin with a TCP/IP service and see how it goes before upgrading it to the WebSocket protocol.

Just Imagine

# Exercises

-   Develop a concurrent TCP server that generates random numbers in a predefined range.
-   Develop a concurrent TCP server that generates random numbers in a range that is given by the TCP client. This can be used as a way of randomly picking values from a set.
-   Add UNIX signal processing to the concurrent TCP server developed in this chapter to gracefully stop the server process when a given signal is received.
-   Develop a UNIX domain socket server that generates random numbers. After that, program a client for that server.
-   Develop a WebSocket server that creates a variable number of random integers that are sent to the client. The number of random integers is specified by the client at the initial client message.

Just Imagine

# Summary

This chapter was all about the `net` package, TCP/IP, TCP, UDP, UNIX sockets, and WebSocket, which implement pretty low-level connections. TCP/IP is what governs the internet. Additionally, WebSocket is handy when you must transfer lots of data. Lastly, UNIX domain sockets are preferred when the data exchange between the server and its various clients takes place on the same machine. Go can help you create all kinds of concurrent servers and clients. You are now ready to begin developing and deploying your own services!

The next chapter is about REST APIs, exchanging JSON data over HTTP, and developing RESTful clients and servers—Go is widely used for developing RESTful clients and servers.

Just Imagine

# Additional resources

-   The WebSocket protocol: [https://tools.ietf.org/rfc/rfc6455.txt](https://tools.ietf.org/rfc/rfc6455.txt)
-   Wikipedia WebSocket: [https://en.wikipedia.org/wiki/WebSocket](https://en.wikipedia.org/wiki/WebSocket)
-   Gorilla WebSocket package: [https://github.com/gorilla/websocket](https://github.com/gorilla/websocket)
-   Gorilla WebSocket docs: [https://www.gorillatoolkit.org/pkg/websocket](https://www.gorillatoolkit.org/pkg/websocket)
-   The `websocket` package: [https://pkg.go.dev/golang.org/x/net/websocket](https://pkg.go.dev/golang.org/x/net/websocket)