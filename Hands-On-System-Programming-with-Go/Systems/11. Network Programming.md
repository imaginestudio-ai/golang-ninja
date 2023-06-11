# Network Programming

In the previous chapter, we talked about developing web applications, talking to databases, and dealing with JSON data in Go.

The topic of this chapter is the development of Go applications that work over TCP/IP networks. In addition, you will learn how to create TCP and UDP clients and servers. The central Go package of this chapter will be the net package: most of its functions are quite low level and require a good knowledge of TCP/IP and its family of protocols.

However, have in mind that network programming is a huge theme that cannot be covered in a single chapter. This chapter will give you the foundational directions for how to create TCP/IP applications in Go.

More analytically, this chapter will talk about the following topics:

-   How TCP/IP operates
-   The net Go standard package
-   Developing TCP clients and servers
-   Programing UDP clients and servers
-   Developing an RPC client
-   Implementing an RPC server
-   The Wireshark and tshark(1) network traffic analyzers
-   Unix sockets
-   Performing DNS lookups from Go programs

Just Imagine

# About network programming

**Network programming** is the development of applications that can operate over computer networks using TCP/IP, which is the dominant networking protocol. Therefore, without knowing the way TCP/IP and its protocols work, you cannot create network applications and develop TCP/IP servers.

The best two advices that I can give to developers of network applications, are to know the theory behind the task they want to perform and to know that networks fail all the time for several reasons. The nastiest types of network failures have to do with malfunctioning or misconfigured DNS servers, because such problems are challenging to find and difficult to correct.

# About TCP/IP

**TCP/IP** is a family of protocols that help the internet to operate. Its name comes from its two most well-known protocols: **TCP** and **IP**.

Every device that uses TCP/IP must have an IP address, which should be unique at least to its local network. It also needs a **network mask** (used for dividing big IP networks into smaller networks) that is related to its current network, one or more **DNS servers** (used for translating an IP address to a human-memorable format and vice versa) and, if you want to communicate with devices beyond your local network, the IP address of a device that will act as the **default gateway** (a network device that TCP/IP sends a network packet to when it cannot find where else to send it).

Each TCP/IP service, which in reality is a Unix process, listens to a port number that is unique to each machine. Note that port numbers 0-1023 are restricted and can only be used by the root user, so it is better to avoid using them and choose something else, provided that it is not already in use by a different process.

# About TCP

**TCP** stands for **Transmission** **Control** **Protocol**. TCP software transmits data between machines using segments, which are called TCP **packets**. The main characteristic of TCP is that it is a reliable protocol, which means that it attempts to make sure that a packet was delivered. If there is no proof of a packet delivery, TCP resends that particular packet. Among other things, a TCP packet can be used for establishing connections, transferring data, sending acknowledgments, and closing connections.

When a TCP connection is established between two machines, a full duplex virtual circuit, similar to the telephone call, is created between these two machines. The two machines constantly communicate to make sure that data are sent and received correctly. If the connection fails for some reason, the two machines try to find the problem and report to the relevant application.

TCP assigns a sequence number to each transmitted packet and expects a positive acknowledgment (ACK) from the receiving TCP stack. If the ACK is not received within a timeout interval, the data is retransmitted as the original packet is considered undelivered. The receiving TCP stack uses the sequence numbers to rearrange the segments when they arrive out of order, which also eliminates duplicate segments.

The TCP header of each packet includes **source port and destination port** fields. These two fields plus the source and destination IP addresses are combined to uniquely identify each TCP connection. The TCP header also includes a 6-bit flags field that is used to relay control information between TCP peers. The possible flags include SYN, FIN, RESET, PUSH, URG, and ACK. The SYN and ACK flags are used for the initial TCP 3-way handshake. The RESET flag signifies that the receiver wants to abort the connection.

# The TCP handshake!

When a connection is initiated, the client sends a TCP SYN packet to the server. The TCP header also includes a sequence number field that has an arbitrary value in the SYN packet. The server sends back a TCP \[SYN, ACK\] packet, which includes the sequence number of the opposite direction and an acknowledgment of the previous sequence number. Finally, in order to truly establish the TCP connection, the client sends a TCP ACK packet in order to acknowledge the sequence number of the server.

Although all these actions take place automatically, it is good to know what is happening behind the scenes!

# About UDP and IP

**IP** stands for **Internet Protocol**. The main characteristic of IP is that it is not a reliable protocol by nature. IP encapsulates the data that travels in a TCP/IP network because it is responsible for delivering packets from the source host to the destination host according to the IP addresses. IP has to find an addressing method to effectively send the packet to its destination. Although there exist dedicated devices called routers that perform IP routing, every TCP/IP device has to perform some basic routing.

**UDP** (short for **User Datagram Protocol**) is based on IP, which means that it is also unreliable. Generally speaking, UDP is simpler than TCP mainly because UDP is not reliable by design. As a result, UDP messages can be lost, duplicated, or arrive out of order. Furthermore, packets can arrive faster than the recipient can process them. So, UDP is used when speed is more important than reliability! An example for this is live video and audio applications where catching up is way more important than buffering and not losing any data.

So, when you do not need too many network packets to transfer the desired information, using a protocol that is based on IP might be more efficient than using TCP, even if you have to retransmit a network packet, because there is no traffic overhead from the TCP handshake.

# About Wireshark and tshark

**Wireshark** is a graphical application for analyzing network traffic of almost any kind. Nevertheless, there are times that you need something lighter that you can execute remotely without a graphical user interface. In such situations, you can use tshark, which is the command-line version of Wireshark.

In order to help you find the network data you really want, Wireshark and tshark have support for capture filters and display filters.

Capture filters are the filters that are applied during network data capturing; therefore, they make Wireshark discard network traffic that does not match the filter. Display filters are the filters that are applied after packet capturing; therefore, they just hide some network traffic without deleting it: you can always disable a display filter and get your hidden data back. Generally speaking, display filters are considered more useful and versatile than capture filters because, normally, you do not know in advance what you will capture or want to examine. Nevertheless, applying filters at capture time can save you time and disk space and that is the main reason for using them.

The following screenshot shows the traffic of a TCP handshake in more detail as captured by Wireshark. The client IP address is 10.0.2.15 and the destination IP address is 80.244.178.150. Additionally, a simple display filter (tcp && !http) makes Wireshark display fewer packets and makes the output less cluttered and therefore easier to read:

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/4cd7d321-edd4-4d49-8713-bc9cea9535f6.png)

The TCP handshake!

The same information can be seen in text format using tshark(1):

```markup
$ tshark -r handshake.pcap -Y '(tcp.flags.syn==1 ) || (tcp.flags == 0x0010 && tcp.seq==1 && tcp.ack==1)'
       18   5.144264    10.0.2.15 → 80.244.178.150 TCP 74 59897 → 80 [SYN] Seq=0 Win=29200 Len=0 MSS=1460 SACK_PERM=1 TSval=1585402 TSecr=0 WS=128
       19   5.236792 80.244.178.150 → 10.0.2.15    TCP 60 80 → 59897 [SYN, ACK] Seq=0 Ack=1 Win=65535 Len=0 MSS=1460
       20   5.236833    10.0.2.15 → 80.244.178.150 TCP 54 59897 → 80 [ACK] Seq=1 Ack=1 Win=29200 Len=0
```

The \-r parameter followed by an existing filename allows you to replay a previously captured data file on your screen, whereas a more complex display filter, which is defined after the \-Y parameter, does the rest of the job!

You can learn more about Wireshark at [https://www.wireshark.org/](https://www.wireshark.org/) and by looking at its documentation at [https://www.wireshark.org/docs/](https://www.wireshark.org/docs/).

# About the netcat utility

There are times that you will need to test a TCP/IP client or a TCP/IP server: the netcat(1) utility can help you with that by playing the role of the client or server in a TCP or UDP application.

You can use netcat(1) as a client for a TCP service that runs on a machine with the 192.168.1.123 IP address and listens to port number 1234, as follows:

```markup
$ netcat 192.168.1.123 1234
```

Similarly, you can use netcat(1) as a client for a UDP service that runs on a Unix machine named amachine.com and listens to port number 2345, as shown here:

```markup
$ netcat -vv -u amachine.com 2345
```

The \-l option tells netcat(1) to listen for incoming connections, which makes netcat(1) to act as a TCP or UDP server. If you try to use netcat(1) as a server with a port that is already in use, you will get the following output:

```markup
$ netcat -vv -l localhost -p 80
Can't grab 0.0.0.0:80 with bind : Permission denied
```

Just Imagine

# The net Go standard package

The most useful Go package for creating TCP/IP applications is the net Go standard package. The net.Dial() function is used for connecting to a network as a client, and the net.Listen() function is used for accepting connections as a server. The first parameter of both functions is the network type, but this is where the similarities end.

For the net.Dial() function, the network type can be one of tcp, tcp4 (IPv4-only), tcp6 (IPv6-only), udp, udp4 (IPv4-only), udp6 (IPv6-only), ip, ip4 (IPv4-only), ip6 (IPv6-only), Unix, Unixgram, or Unixpacket. For the net.Listen() function, the first parameter can be one of tcp, tcp4, tcp6, Unix, or Unixpacket.

The return value of the net.Dial() function is of the net.Conn interface type, which implements the io.Reader and io.Writer interfaces! This means that you already know how to access the variables of the net.Conn interface!

So, although the way you create a network connection is different from the way you create a text file, their access methods are the same because the net.Conn interface implements the io.Reader and io.Writer interfaces. Therefore, as network connections are treated as files, you might need to review _[](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_[Chapter 6](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_,_ _File Input and Output_, at this moment.

Just Imagine

# Unix sockets revisited

Back in _[](https://subscription.imaginedevops.io/book/programming/9781787125643/8)_[Chapter 8](https://subscription.imaginedevops.io/book/programming/9781787125643/8)_,_ _Processes and Signals_, we talked a little about Unix sockets and presented a small Go program that was acting as a Unix socket client. This section will also create a Unix socket server to make things even clearer. However, the Go code of the Unix socket client will be also explained here in more detail and will be enriched with error handling code.

# A Unix socket server

The Unix socket server will act as an Echo server, which means that it will send the received message back to the client. The name of the program will be socketServer.go and it will be presented to you in four parts.

The first part of socketServer.go is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "net" 
   "os" 
) 
```

The second part of the Unix socket server is the following:

```markup
func echoServer(c net.Conn) { 
   for { 
         buf := make([]byte, 1024) 
         nr, err := c.Read(buf) 
         if err != nil { 
               return 
         } 
 
         data := buf[0:nr] 
         fmt.Printf("->: %v\n", string(data)) 
         _, err = c.Write(data) 
         if err != nil { 
               fmt.Println(err) 
         } 
   } 
} 
```

This is where the function that serves incoming connections is implemented.

The third portion of the program has the following Go code:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide a socket file.") 
         os.Exit(100) 
   } 
   socketFile := arguments[1] 
 
   l, err := net.Listen("unix", socketFile) 
   if err != nil { 
         fmt.Println(err) 
os.Exit(100) 
   } 
```

Here, you can see the use of the net.Listen() function with the unix argument for creating the desired socket file.

Finally, the last part contains the following Go code:

```markup
   for { 
         fd, err := l.Accept() 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
         go echoServer(fd) 
   } 
} 
```

As you can see, each connection is first handled by the Accept() function and served by its own goroutine.

When socketServer.go serves a client, it generates the following output:

```markup
$ go run socketServer.go /tmp/aSocket
->: Hello Server!
```

If you cannot create the desired socket file, for instance, if it already exists, you will get an error message similar to the following:

```markup
$ go run socketServer.go /tmp/aSocket
listen unix /tmp/aSocket: bind: address already in use
exit status 100
```

# A Unix socket client

The name of the Unix socket client program is socketClient.go and will be presented in four parts.

The first part of the utility contains the expected preamble:

```markup
package main 
 
import ( 
   "fmt" 
   "io" 
   "log" 
   "net" 
   "os" 
   "time" 
) 
```

There is nothing special here, just the required Go packages. The second portion contains the definition of a Go function:

```markup
func readSocket(r io.Reader) { 
   buf := make([]byte, 1024) 
   for { 
         n, err := r.Read(buf[:]) 
         if err != nil { 
               fmt.Println(err) 
               return 
         } 
         fmt.Println("-> ", string(buf[0:n])) 
   } 
} 
```

The readSocket() function reads the data from a socket file using Read(). Note that, although socketClient.go just reads from the socket file, the socket is bisectional, which means that you can also write to it.

The third part has the following Go code:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide a socket file.") 
         os.Exit(100) 
   } 
   socketFile := arguments[1] 
 
   c, err := net.Dial("unix", socketFile) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
   defer c.Close() 
```

The net.Dial() function with the right first argument allows you to connect to the socket file before you try to read from it.

The last part of socketClient.go is the following:

```markup
   go readSocket(c) 
   for { 
         _, err := c.Write([]byte("Hello Server!")) 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
         time.Sleep(1 * time.Second) 
   } 
} 
```

In order to use socketClient.go, you must have another program dealing with the Unix socket file, which, in this case will be socketServer.go. So, if socketServer.go is already running, you will get the following output from socketClient.go:

```markup
$ go run socketClient.go /tmp/aSocket
->: Hello Server!
```

If you do not have enough Unix file permissions to read the desired socket file, then socketClient.go will fail with the following error message:

```markup
$ go run socketClient.go /tmp/aSocket
dial unix /tmp/aSocket: connect: permission denied
exit status 100
```

Similarly, if the socket file you want to read does not exist, socketClient.go will fail with the following error message:

```markup
$ go run socketClient.go /tmp/aSocket
dial unix /tmp/aSocket: connect: no such file or directory
exit status 100
```

Just Imagine

# Performing DNS lookups

There exist many types of DNS lookups, but two of them are the most popular. In the first type, you want to go from an IP address to a domain name and in the second type you want to go from a domain name to an IP address.

The following output shows an example of the first type of DNS lookup:

```markup
$ host 109.74.193.253
253.193.74.109.in-addr.arpa domain name pointer li140-253.members.linode.com.
```

The following output shows three examples of the second type of DNS lookup:

```markup
$ host www.mtsoukalos.eu
www.mtsoukalos.eu has address 109.74.193.253
$ host www.highiso.net
www.highiso.net has address 109.74.193.253
$ host -t a cnn.com
cnn.com has address 151.101.1.67
cnn.com has address 151.101.129.67
cnn.com has address 151.101.65.67
cnn.com has address 151.101.193.67
```

As you just saw in the aforementioned examples, an IP address can serve many hosts and a host name can have many IP addresses.

The Go standard library provides the net.LookupHost() and net.LookupAddr() functions that can answer DNS queries for you. However, none of them allow you to define the DNS server you want to query. While using standard Go libraries is ideal, there exist external Go libraries that allow you to choose the DNS server you desire, which is mainly required when troubleshooting DNS configurations.

# Using an IP address as input

The name of the Go utility that will return the hostname of an IP address will be lookIP.go and will be presented in three parts.

The first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "net" 
   "os" 
) 
```

The second part has the following Go code:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide an IP address!") 
         os.Exit(100) 
   } 
 
   IP := arguments[1] 
   addr := net.ParseIP(IP) 
   if addr == nil { 
         fmt.Println("Not a valid IP address!") 
         os.Exit(100) 
   } 
```

The net.ParseIP() function allows you to verify the validity of the given IP address and is pretty handy for catching illegal IP addresses such as 288.8.8.8 and 8.288.8.8.

The last part of the utility is the following:

```markup
   hosts, err := net.LookupAddr(IP) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   for _, hostname := range hosts { 
         fmt.Println(hostname) 
   } 
} 
```

As you can see, the net.LookupAddr() function returns a string slice with the list of names that match the given IP address.

Executing lookIP.go will generate the following output:

```markup
$ go run lookIP.go 288.8.8.8
Not a valid IP address!
exit status 100
$ go run lookIP.go 8.8.8.8
google-public-dns-a.google.com.
```

You can validate the output of dnsLookup.go using host(1) or dig(1):

```markup
$ host 8.8.8.8
8.8.8.8.in-addr.arpa domain name pointer google-public-dns-a.google.com.
```

# Using a host name as input

The name of this DNS utility will be lookHost.go and will be presented in three parts. The first part of the lookHost.go utility is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "net" 
   "os" 
) 
```

The second part of the program has the following Go code:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide an argument!") 
         os.Exit(100) 
   } 
 
   hostname := arguments[1] 
   IPs, err := net.LookupHost(hostname) 
```

Similarly, the net.LookupHost() function also returns a string slice with the desired information.

The third part of the program has the following code, which is for error checking and printing the output of net.LookupHost():

```markup
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   for _, IP := range IPs { 
         fmt.Println(IP) 
   } 
} 
```

Executing lookHost.go will generate the following output:

```markup
$ go run lookHost.go www.google
lookup www.google: no such host
exit status 100
$ go run lookHost.go www.google.com
2a00:1450:4001:81f::2004
172.217.16.164
```

The first line of the output is the IPv6 address, whereas the second output line is the IPv4 address of www.google.com.

You can verify the operation of lookHost.go by comparing its output with the output of the host(1) utility:

```markup
$ host www.google.com
www.google.com has address 172.217.16.164
www.google.com has IPv6 address 2a00:1450:4001:81a::2004
```

# Getting NS records for a domain

This subsection will present an additional kind of DNS lookup that returns the domain name servers for a given domain. This is very handy for troubleshooting DNS-related problems and finding out the status of a domain. The presented program will be named lookNS.go and will be presented in three parts.

The first part of the utility is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "net" 
   "os" 
) 
```

The second part has the following Go code:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide a domain!") 
         os.Exit(100) 
   } 
 
   domain := arguments[1] 
 
   NSs, err := net.LookupNS(domain) 
```

The net.LookupNS() function does all the work for us by returning a slice of NS elements.

The last part of the code is mainly for printing the results:

```markup
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   for _, NS := range NSs { 
         fmt.Println(NS.Host) 
   } 
} 
```

Executing lookNS.go will generate the following output:

```markup
$ go run lookNS.go mtsoukalos.eu
ns5.linode.com.
ns2.linode.com.
ns3.linode.com.
ns1.linode.com.
ns4.linode.com.
```

The reason that the following query will fail is that www.mtsoukalos.eu is not a domain but a single host, which means that it has no NS records associated with it:

```markup
$ go run lookNS.go www.mtsoukalos.eu
lookup www.mtsoukalos.eu on 8.8.8.8:53: no such host
exit status 100
```

You can use the host(1) utility to verify the previous output:

```markup
$ host -t ns mtsoukalos.eu
mtsoukalos.eu name server ns5.linode.com.
mtsoukalos.eu name server ns4.linode.com.
mtsoukalos.eu name server ns3.linode.com.
mtsoukalos.eu name server ns1.linode.com.
mtsoukalos.eu name server ns2.linode.com.
$ host -t ns www.mtsoukalos.eu
www.mtsoukalos.eu has no NS record
```

Just Imagine

# Developing a simple TCP server

This section will develop a TCP server that implements the **Echo** service. The Echo service is usually implemented using the UDP protocol due to its simplicity, but it can also be implemented with TCP. The Echo service usually uses port number 7, but our implementation will use other port numbers:

```markup
$ grep echo /etc/services
echo        7/tcp
echo        7/udp
```

The TCPserver.go file will hold the Go code of this section and will be presented in six parts. For reasons of simplicity, each connection is handled inside the main() function without calling a separate function. However, this is not the recommended practice.

The first part contains the expected preamble:

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

The second part of the TCP server is the following:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide port number") 
         os.Exit(100) 
   } 
```

The third part of TCPserver.go contains the following Go code:

```markup
   PORT := ":" + arguments[1] 
   l, err := net.Listen("tcp", PORT) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
   defer l.Close() 
```

What is important to remember here is that net.Listen() returns a Listener variable, which is a generic network listener for stream-oriented protocols. Additionally, the Listen() function can support more formats: check the documentation of the net package to find more information about that.

The fourth part of the TCP server has the following Go code:

```markup
   c, err := l.Accept() 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
```

Only after a successful call to Accept(), the TCP server can start interacting with TCP clients. Nonetheless, the current version of TCPserver.go has a very serious shortcoming: it can only serve a single TCP client, the first one that will connect to it.

The fifth portion of the TCPserver.go code is the following:

```markup
   for { 
         netData, err := bufio.NewReader(c).ReadString('\n') 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
```

Here, you read data from your client using bufio.NewReader().ReadString(). The aforementioned call allows you to read your input line by line. Additionally, the for loop allows you to keep reading data from the TCP client for as long as you wish.

The last part of the Echo TCP server is the following:

```markup
         fmt.Print("-> ", string(netData)) 
         c.Write([]byte(netData)) 
         if strings.TrimSpace(string(netData)) == "STOP" { 
               fmt.Println("Exiting TCP server!") 
               return 
         } 
   } 
} 
```

The current version of TCPserver.go stops when it receives the STOP string as input. Although TCP servers do not usually terminate in that style, this is a pretty handy way to terminate a TCP server process that will only serve a single client!

Next, we will test TCPserver.go with netcat(1):

```markup
$ go run TCPserver.go 1234
-> Hi!
-> STOP
Exiting TCP server!
```

The netcat(1) part is the following:

```markup
$ nc localhost 1234
Hi!
Hi!
STOP
STOP
```

Here, the first and third lines are our input, whereas the second and fourth lines are the responses from the Echo server.

If you try to use an improper port number, TCPserver.go will generate the following error message and exit:

```markup
$ go run TCPserver.go 123456
listen tcp: address 123456: invalid port
exit status 100
```

Just Imagine

# Developing a simple TCP client

In this section, we will develop a TCP client named TCPclient.go. The port number the client will try to connect to as well as the server address will be given as command-line arguments to the program. The Go code of the TCP client will be presented in five parts; the first part is the following:

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

The second part of TCPclient.go is the following:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide host:port.") 
         os.Exit(100) 
   } 
```

The third part of TCPclient.go has the following Go code:

```markup
   CONNECT := arguments[1] 
   c, err := net.Dial("tcp", CONNECT) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
```

Once again, you use the net.Dial() function to try to connect to the desired port of the desired TCP server.

The fourth part of the TCP client is the following:

```markup
   for { 
         reader := bufio.NewReader(os.Stdin) 
         fmt.Print(">> ") 
         text, _ := reader.ReadString('\n') 
         fmt.Fprintf(c, text+"\n") 
```

Here, you read data from the user that you will send to the TCP server using fmt.Fprintf().

The last part of TCPclient.go is the following:

```markup
         message, _ := bufio.NewReader(c).ReadString('\n') 
         fmt.Print("->: " + message) 
         if strings.TrimSpace(string(text)) == "STOP" { 
               fmt.Println("TCP client exiting...") 
               return 
         } 
   } 
} 
```

In this part, you get data from the TCP server using bufio.NewReader().ReadString(). The reason for using the strings.TrimSpace() function is to remove any spaces and newline characters from the variable you want to compare with the static string (STOP).

So, now it is time to verify that TCPclient.go works as expected using it to connect to TCPserver.go:

```markup
$ go run TCPclient.go localhost:1024
>> 123
->: 123
>> Hello server!
->: Hello server!
>> STOP
->: STOP
TCP client exiting...
```

If no process listens to the specified TCP port at the specified host, then you will get an error message similar to the following:

```markup
$ go run TCPclient.go localhost:1024
dial tcp [::1]:1024: getsockopt: connection refused
exit status 100
```

# Using other functions for the TCP server

In this subsection, we will develop the functionality of TCPserver.go using some slightly different functions. The name of the new TCP server will be TCPs.go and will be presented in four parts.

The first part of TCPs.go is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "net" 
   "os" 
) 
```

The second part of the TCP server is the following:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide a port number!") 
         os.Exit(100) 
   } 
 
   SERVER := "localhost" + ":" + arguments[1] 
```

So far, there are no differences from the code of TCPserver.go.

The differences start in the third part of TCPs.go, which is the following:

```markup
   s, err := net.ResolveTCPAddr("tcp", SERVER) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   l, err := net.ListenTCP("tcp", s) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
```

Here, you use the net.ResolveTCPAddr() and net.ListenTCP() functions. Is this version better than TCPserver.go? Not really. But the Go code might look a little clearer and this is a big advantage for some people. Additionally, net.ListenTCP() returns a TCPListener value that when used with net.AcceptTCP() instead of net.Accept() will return TCPConn, which offers more methods that allow you to change more socket options.

The last part of TCPs.go has the following Go code:

```markup
   buffer := make([]byte, 1024) 
 
   for { 
         conn, err := l.Accept() 
         n, err := conn.Read(buffer) 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
 
         fmt.Print("> ", string(buffer[0:n])) 
         _, err = conn.Write(buffer) 
 
         conn.Close() 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
   } 
} 
```

There is nothing special here. You still use Accept() to get and process client requests. However, this version uses Read() to get the client data all at once, which is great when you do not have to process lots of input.

The operation of TCPs.go is the same with the operation of TCPserver.go, so it will not be shown here.

If you try to create a TCP server using an invalid port number, TCPs.go will generate an informative error message, as shown here:

```markup
$ go run TCPs.go 123456
address 123456: invalid port
exit status 100
```

# Using alternative functions for the TCP client

Once again, we will implement TCPclient.go using some slightly different functions that are provided by the net Go standard package. The name of the new version will be TCPc.go and will be shown in four code segments.

The first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "net" 
   "os" 
) 
```

The second code segment of the program is the following:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide a server:port string!") 
         os.Exit(100) 
   } 
 
   CONNECT := arguments[1] 
   myMessage := "Hello from TCP client!\n" 
```

This time, we will send a static message to the TCP server.

The third part of TCPc.go is the following:

```markup
   tcpAddr, err := net.ResolveTCPAddr("tcp", CONNECT) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   conn, err := net.DialTCP("tcp", nil, tcpAddr) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
```

In this part, you see the use of net.ResolveTCPAddr() and net.DialTCP(), which is where the differences between TCPc.go and TCPclient.go exist.

The last part of the TCP client is the following:

```markup
   _, err = conn.Write([]byte(myMessage)) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   fmt.Print("-> ", myMessage) 
   buffer := make([]byte, 1024) 
 
   n, err := conn.Read(buffer) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   fmt.Print(">> ", string(buffer[0:n])) 
   conn.Close() 
} 
```

You might ask if you can use TCPc.go with TCPserver.go or TCPs.go with TCPclient.go. The answer is a definitive _yes_ because the implementation and the function names have nothing to do with the actual TCP/IP operations that take place.

Just Imagine

# Developing a simple UDP server

This section will also develop an Echo server. However, this time the Echo server will use the UDP protocol. The name of the program will be UDPserver.go and will be presented to you in five parts.

The first part contains the expected preamble:

```markup
package main 
 
import ( 
   "fmt" 
   "net" 
   "os" 
   "strings" 
) 
```

The second part is the following:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide a port number!") 
         os.Exit(100) 
   } 
   PORT := ":" + arguments[1] 
```

The third part of UDPserver.go is the following:

```markup
   s, err := net.ResolveUDPAddr("udp", PORT) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   connection, err := net.ListenUDP("udp", s) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
```

The UDP approach is similar to the TCP approach: you just call functions with different names.

The fourth part of the program has the following Go code:

```markup
   defer connection.Close() 
   buffer := make([]byte, 1024) 
 
   for { 
         n, addr, err := connection.ReadFromUDP(buffer) 
         fmt.Print("-> ", string(buffer[0:n])) 
         data := []byte(buffer[0:n]) 
         _, err = connection.WriteToUDP(data, addr) 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
```

In the UDP case, you use ReadFromUDP() to read from a UDP connection and WriteToUDP() to write to an UDP connection. Additionally, the UDP connection does not need to call a function similar to net.Accept().

The last part of the UDP server is the following:

```markup
         if strings.TrimSpace(string(data)) == "STOP" { 
               fmt.Println("Exiting UDP server!") 
               return 
         } 
   } 
} 
```

Once again, we will test UDPserver.go with netcat(1):

```markup
$ go run UDPserver.go 1234
-> Hi!
-> Hello!
-> STOP
Exiting UDP server!
```

Just Imagine

# Developing a simple UDP client

In this section, we will develop a UDP client, which we will name UDPclient.go and present in five parts.

As you will see, the code differences between the Go code of UDPclient.go and TCPc.go are basically the differences in the names of the functions used: the general idea is exactly the same.

The first part of the UDP client is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "net" 
   "os" 
) 
```

The second part of the utility contains the following Go code:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide a host:port string") 
         os.Exit(100) 
   } 
   CONNECT := arguments[1] 
```

The third part of UDPclient.go has the following Go code:

```markup
   s, err := net.ResolveUDPAddr("udp", CONNECT) 
   c, err := net.DialUDP("udp", nil, s) 
 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String()) 
   defer c.Close() 
```

Nothing special here: just the use of net.ResolveUDPAddr() and net.DialUDP() to connect to the UDP server.

The fourth part of the UDP client is the following:

```markup
   data := []byte("Hello UDP Echo server!\n") 
   _, err = c.Write(data) 
 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
```

This time, you send your data to the UDP server using Write(), although you will read from the UDP server using ReadFromUDP().

The last part of UDPclient.go is the following:

```markup
   buffer := make([]byte, 1024) 
   n, _, err := c.ReadFromUDP(buffer) 
   fmt.Print("Reply: ", string(buffer[:n])) 
} 
```

As we have UDPserver.go and we know that it works, we can test the operation of UDPclient.go using UDPserver.go:

```markup
$ go run UDPclient.go localhost:1234
The UDP server is 127.0.0.1:1234
Reply: Hello UDP Echo server!
```

If you execute UDPclient.go without a UDP server listening to the desired port, you will get the following output, which does not clearly state that it could not connect to an UDP server: it just shows an empty reply:

```markup
$ go run UDPclient.go localhost:1024
The UDP server is 127.0.0.1:1024
Reply:
```

Just Imagine

# A concurrent TCP server

In this section, you will learn how to develop a concurrent TCP server: each client connection will be assigned to a new goroutine that will serve the client request. Note that although TCP clients initially connect to the same port, they are served using a different port number than the main port number of the server: this is automatically handled by TCP and is the way TCP works.

Although creating a concurrent UDP server is also a possibility, it might not be absolutely necessary due to the way UDP works. However, if you have a really busy UDP service, then you might consider developing a concurrent UDP server.

The name of the program will be concTCP.go and will be presented in five parts. The good thing is that once you define a function to handle incoming connections, all you need is to execute that function as a goroutine, and the rest will be handled by Go!

The first part of concTCP.go is the following:

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
```

The second part of the concurrent TCP server is the following:

```markup
func handleConnection(c net.Conn) { 
   for { 
         netData, err := bufio.NewReader(c).ReadString('\n') 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
 
         fmt.Print("-> ", string(netData)) 
         c.Write([]byte(netData)) 
         if strings.TrimSpace(string(netData)) == "STOP" { 
               break 
         } 
   } 
   time.Sleep(3 * time.Second) 
   c.Close() 
} 
```

Here is the implementation of the function that handles each TCP request. The time delay at the end of it is used for giving you the necessary time to connect with another TCP client and prove that concTCP.go can serve multiple TCP clients.

The third part of the program contains the following Go code:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide a port number!") 
         os.Exit(100) 
   } 
 
   PORT := ":" + arguments[1] 
```

The fourth part of concTCP.go has the following Go code:

```markup
   l, err := net.Listen("tcp", PORT) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
   defer l.Close() 
```

So far, there is nothing special in the main() function because although concTCP.go will handle multiple requests, it only needs a single call to net.Listen().

The last chunk of Go code is the following:

```markup
   for { 
         c, err := l.Accept() 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
         go handleConnection(c) 
   } 
} 
```

All the differences in the way concTCP.go processes its requests can be found in the last lines of Go code. Each time the program accepts a new network request using Accept(), a new goroutine gets started and concTCP.go is immediately ready to accept more requests. Note that in order to terminate concTCP.go, you will have to press _Ctrl_ + _C_ because the STOP keyword is used for terminating each goroutine of the program.

Executing concTCP.go and connecting to it using various TCP clients, will generate the following output:

```markup
$ go run concTCP.go 1234
-> Hi!
-> Hello!
-> STOP
...
```

Just Imagine

# Remote procedure call (RPC)

**Remote Procedure Call** (**RPC**) is a client-server mechanism for interprocess communication. Note that the RPC client and the RPC server communicate using TCP/IP, which means that they can exist in different machines.

In order to develop the implementation of an RPC client or RPC server, you will need to follow some steps and call some functions in a given way. Neither of the two implementations is difficult; you just have to follow certain steps.

Also, visit the documentation page of the net/rpc Go standard package that can be found at https://golang.org/pkg/net/rpc/.

Note that the presented RPC example will use TCP for client-server interaction. However, you can also use HTTP for client-server communication.

# An RPC server

This subsection will present an RPC server named RPCserver.go. As you will see in the preamble of the RPCserver.go program, the RPC server imports a package named sharedRPC, which is implemented in the sharedRPC.go file: the name of the package is arbitrary. Its contents are the following:

```markup
package sharedRPC 
 
type MyInts struct { 
   A1, A2 uint 
   S1, S2 bool 
} 
type MyInterface interface { 
   Add(arguments *MyInts, reply *int) error 
   Subtract(arguments *MyInts, reply *int) error 
} 
```

So, here you define a new structure that holds the signs and the values of two unsigned integers and a new interface named MyInterface.

Then, you should install sharedRPC.go, which means that you should execute the following commands before you try to use the sharedRPC package in your programs:

```markup
$ mkdir ~/go
$ mkdir ~/go/src
$ mkdir ~/go/src/sharedRPC
$ export GOPATH=~/go
$ vi ~/go/src/sharedRPC/sharedRPC.go
$ go install sharedRPC
```

If you are on a macOS machine (darwin\_amd64) and you want to make sure that everything is OK, you can execute the following two commands:

```markup
$ cd ~/go/pkg/darwin_amd64/
$ ls -l sharedRPC.a
-rw-r--r--  1 mtsouk  staff  4698 Jul 27 11:49 sharedRPC.a
```

What you really must keep in mind is that, at the end of the day, what is being exchanged between an RPC server and an RPC client are function names and their arguments. Only the functions defined in the interface of sharedRPC.go can be used in an RPC interaction: the RPC server will need to implement the functions of the MyInterface interface. The Go code of RPCserver.go will be presented in five parts; the first part of the RPC server has the expected preamble, which also includes the sharedRPC package we made:

```markup
package main 
 
import ( 
   "fmt" 
   "net" 
   "net/rpc" 
   "os" 
   "sharedRPC" 
) 
```

The second part of RPCserver.go is the following:

```markup
type MyInterface int 
 
func (t *MyInterface) Add(arguments *sharedRPC.MyInts, reply *int) error { 
   s1 := 1 
   s2 := 1 
 
   if arguments.S1 == true { 
         s1 = -1 
   } 
 
   if arguments.S2 == true { 
         s2 = -1 
   } 
 
   *reply = s1*int(arguments.A1) + s2*int(arguments.A2) 
   return nil 
} 
```

Here is the implementation of the first function that will be offered to the RPC clients: you can have as many functions as you want, provided that they are included in the interface.

The third part of RPCserver.go has the following Go code:

```markup
func (t *MyInterface) Subtract(arguments *sharedRPC.MyInts, reply *int) error { 
   s1 := 1 
   s2 := 1 
 
   if arguments.S1 == true { 
         s1 = -1 
   } 
 
   if arguments.S2 == true { 
         s2 = -1 
   } 
 
   *reply = s1*int(arguments.A1) - s2*int(arguments.A2) 
   return nil 
} 
```

This is the second function that is offered to the RPC clients by this RPC server.

The fourth part of RPCserver.go contains the following Go code:

```markup
func main() { 
   PORT := ":1234" 
 
   myInterface := new(MyInterface) 
   rpc.Register(myInterface) 
 
   t, err := net.ResolveTCPAddr("tcp", PORT) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
   l, err := net.ListenTCP("tcp", t) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
```

As our RPC server uses TCP, you need to make calls to net.ResolveTCPAddr() and net.ListenTCP(). However, you will first need to call rpc.Register() in order to be able to serve the desired interface.

The last part of the program is the following:

```markup
   for { 
         c, err := l.Accept() 
         if err != nil { 
               continue 
         } 
         rpc.ServeConn(c) 
   } 
} 
```

Here, you accept a new TCP connection using Accept() as usual, but you serve it using rpc.ServeConn().

You will have to wait for the next section and the development of the RPC client in order to test the operation of RPCserver.go.

# An RPC client

In this section, we will develop an RPC client named RPCclient.go. The Go code of RPCclient.go will be presented in five parts; the first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "net/rpc" 
   "os" 
   "sharedRPC" 
) 
```

Note the use of the sharedRPC package in the RPC client.

The second part of RPCclient.go is the following:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide a host:port string!") 
         os.Exit(100) 
   } 
 
   CONNECT := arguments[1] 
```

The third part of the program has the following Go code:

```markup
   c, err := rpc.Dial("tcp", CONNECT) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   args := sharedRPC.MyInts{17, 18, true, false} 
   var reply int 
```

As the MyInts structure is defined in sharedRPC.go, you need to use it as sharedRPC.MyInts in the RPC client. Moreover, you call rpc.Dial() to connect to the RPC server instead of net.Dial().

The fourth part of the RPC client contains the following Go code:

```markup
   err = c.Call("MyInterface.Add", args, &reply) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
   fmt.Printf("Reply (Add): %d\n", reply) 
```

Here, you use the Call() function to execute the desired function in the RPC server. The result of the MyInterface.Add() function is stored in the reply variable, which was previously declared.

The last part of RPCclient.go is the following:

```markup
   err = c.Call("MyInterface.Subtract", args, &reply) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
   fmt.Printf("Reply (Subtract): %d\n", reply) 
} 
```

Here, you do the same thing as before for executing the MyInterface.Subtract() function.

As you can guess, you cannot test the RPC client without having an RCP server and vice versa: netcat(1) cannot be used for RPC.

First, you will need to start the RPCserver.go process:

```markup
$ go run RPCserver.go
```

Then, you will execute the RPCclient.go program:

```markup
$ go run RPCclient.go localhost:1234
Reply (Add): 1
Reply (Subtrack): -35
```

If the RPCserver.go process is not running and you try to execute RPCclient.go, you will get the following error message:

```markup
$ go run RPCclient.go localhost:1234
dial tcp [::1]:1234: getsockopt: connection refused
exit status 100
```

Of course, RPC is not for adding integers or natural numbers, but for doing much more complex operations that you want to control from a central point.

Just Imagine

# Exercises

1.  Read the documentation of the net package in order to find out about its list of available functions at [https://golang.org/pkg/net/](https://golang.org/pkg/net/).
2.  Wireshark is a great tool for analyzing network traffic of any kind: try to use it more.
3.  Change the code of socketClient.go in order to read the input from the user.
4.  Change the code of socketServer.go in order to return a random number to the client.
5.  Change the code of TCPserver.go in order to stop when it receives a given Unix signal from the user.
6.  Change the Go code of concTCP.go in order to keep track of the number of clients it has served and print that number before exiting.
7.  Add a quit() function to RPCserver.go that does what its name implies.
8.  Develop your own RPC example.

Just Imagine

# Summary

In this chapter, we introduced you to TCP/IP, and we talked about developing TCP and UDP servers and clients in Go and about creating RPC clients and servers.

At this point, there is no next chapter because this is the last chapter of this book! Congratulations for reading the whole book! You are now ready to start developing useful Unix command-line utilities in Go; so, go ahead and start programming your own tools immediately!