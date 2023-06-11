# Working with gRPC

This chapter is about working with gRPC in Go. **gRPC**, which stands for **gRPC Remote Procedure Calls**, is an alternative to RESTful services that was developed by Google. The main advantage of gRPC is that it is faster than working with REST and JSON messages. Additionally, creating clients for gRPC services is also faster due to the available tooling. Last, as gRPC uses the binary data format, it is lighter than RESTful services that work with the JSON format.

The process for creating a gRPC server and client has three main steps. The first step is creating the **interface definition language** (**IDL**) file, the second step is the development of the gRPC server, and the third step is the development of the gRPC client that can interact with the gRPC server.

This chapter covers the following topics:

-   Introduction to gRPC
-   Defining an interface definition language file
-   Developing a gRPC server
-   Developing a gRPC client

Just Imagine

# Introduction to gRPC

gRPC is an open source **remote procedure call** (**RPC**) system that was developed at Google back in 2015, is built on top of HTTP/2, allows you to create services easily, and uses _protocol buffers_ as the IDL which specifies the format of the interchanged messages and the service interface.

gRPC clients and servers can be written in any programming language without the need to have clients written in the same programming language as their servers. This means that you can develop a client in Python, even if the gRPC server is implemented in Go. The list of supported programming languages includes, but is not limited to, Python, Java, C++, C#, PHP, Ruby, and Kotlin.

The advantages of gRPC include the following:

-   The use of binary format for data exchange makes gRPC faster than services that work with data in plain text format
-   The command-line tools provided make your work simpler and faster
-   Once you have defined the functions and the messages of a gRPC service, creating servers and clients for it is simpler than RESTful services
-   gRPC can be used for streaming services
-   You do not have to deal with the details of data exchange because gRPC takes care of the details

The list of advantages should not make you think that gRPC is a panacea that does not have any flaws—always use the best tool or technology for the job.

The next section discusses protocol buffers, which are the foundation of gRPC services.

## Protocol buffers

A **protocol buffer** (**protobuf**) is basically a **method for serializing structured data**. A part of protobuf is the IDL. As protobuf uses binary format for data exchange, it takes up less space than plain text serialization formats. However, data needs to be encoded and decoded to be machine-usable and human-readable, respectively. Protobuf has its own data types that are translated to natively supported data types of the programming language used.

Generally speaking, the IDL file is the center of each gRPC service because it defines the format of the data that is exchanged as well as the service interface. You cannot have a gRPC service without having a protobuf file at hand. Strictly speaking, a protobuf file includes the definition of services, methods of services, and the format of the messages that are going to be exchanged—it is not an exaggeration to say that if you want to understand a gRPC service, you should start by looking at its definition file. The next section shows the protobuf file that is going to be used in our gRPC service.

Just Imagine

# Defining an interface definition language file

The gRPC service that we are developing is going to support the following functionality:

-   The server should return its date and time to the client
-   The server should return a randomly generated password of a given length to the client
-   The server should return random integers to the client

Before we begin developing the gRPC client and server for our service, we need to define the IDL file. We need a separate GitHub repository to host the files related to the IDL file, which is going to be [https://github.com/mactsouk/protoapi](https://github.com/mactsouk/protoapi).

Next, we are going to present the IDL file, which is called `protoapi.proto`:

```markup
syntax = "proto3";
```

The presented file uses the `proto3` version of the protocol buffers language—there is also an older version of the language named `proto2`, which has some minor syntax differences. If you do not specify that you want to use `proto3`, then the protocol buffer compiler assumes you want to use `proto2`. The definition of the version must be in the first non-empty, non-comment line in the `.proto` file.

```markup
option go_package = "./;protoapi";
```

The gRPC tools are going to generate Go code from that `.proto` file. The previous line specifies that the name of the Go package that is going to be created is `protoapi`. The output is going to be written in the current directory as `protoapi.proto` due to the use of `./`.

```markup
service Random {
    rpc GetDate (RequestDateTime) returns (DateTime);
    rpc GetRandom (RandomParams) returns (RandomInt);
    rpc GetRandomPass (RequestPass) returns (RandomPass);
}
```

This block specifies the name of the gRPC service (`Random`) as well as the supported methods. Additionally, it specifies the messages that need to be exchanged for an interaction. So, for `GetDate`, the client needs to send a `RequestDateTime` message and expects to get a `DateTime` message back.

These messages are defined in the same `.proto` file.

```markup
// For random number
```

All `.proto` files support C- and C++-type comments. This means that you can use `// text` and `/* text */` comments in your `.proto` files.

```markup
message RandomParams {
    int64 Seed = 1;
    int64 Place = 2;
}
```

A random number generator starts with a seed value, which in our case is specified by the client and sent to the server with a `RandomParams` message. The `Place` field specifies the place of the random number that is going to be returned in the sequence of randomly generated integers.

```markup
message RandomInt {
    int64 Value = 1;
}
```

The previous two messages are related to the `GetRandom` method. `RandomParams` is for setting the parameters of the request, whereas `RandomInt` is for storing a random number that is generated by the server. All message fields are of the `int64` data type.

```markup
// For date time
message DateTime {
    string Value = 1;
}
message RequestDateTime {
    string Value = 2;
}
```

The previous two messages are for supporting the operation of the `GetDate` method. The `RequestDateTime` message is a dummy one in the sense that it does not hold any useful data—we just need to have a message that the client sends to the server—you can store any kind of information in the `Value` field of `RequestDateTime`. The information returned by the server is stored as a `string` value in a `DateTime` message.

```markup
// For random password
message RequestPass {
    int64 Seed = 1;
    int64 Length = 8;
}
message RandomPass {
    string Password = 1;
}
```

Lastly, the previous two messages are for the operation of `GetRandomPass`.

So, the IDL file:

-   specifies that we are using `proto3`.
-   defines the name of the service, which is `Random`.
-   specifies that the name of the generated Go package is going to be `protoapi`.
-   defines that the gRPC service is going to support three methods: `GetDate`, `GetRandom`, and `GetRandomPass`. It also defines the names of the messages that are going to be exchanged in these three method calls.
-   defines the format of six messages that are used for data exchange.

The next important step is converting that file into a format that can be used by Go. You might need to download some extra tools in order to process `protoapi.proto`, or any other `.proto` file, and generate the relevant Go `.pb.go` files. The name of the protocol buffer compiler binary is `protoc`—on my macOS machine, I had to install `protoc` using the `brew install protobuf` command. Similarly, I also had to install `protoc-gen-go-grpc` and `protoc-gen-go` packages using Homebrew—the last two packages are Go-related.

On a Linux machine, you need to install `protobuf` using your favorite package manager and `protoc-gen-go` using the `go install github.com/golang/protobuf/protoc-gen-go@latest` command. Similarly, you should install the `protoc-gen-go-grpc` executable by running `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`.

Starting with Go 1.16, `go install` is the recommended way of building and installing packages in module mode. The use of `go get` is deprecated. However, when using `go install`, do not forget to add `@latest` after the package name to install the latest version.

So, the conversion process requires the next step:

```markup
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protoapi.proto
```

After that, we have a file named `protoapi_grpc.pb.go` and a file named `protoapi.pb.go`—both located in the root directory of the GitHub repository. The `protoapi.pb.go` source code file contains the messages, whereas `protoapi_grpc.pb.go` contains the services.

The first ten lines of `protoapi_grpc.pb.go` are as follows:

```markup
// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
package protoapi
```

As discussed earlier, the name of the package is `protoapi`.

```markup
import (
        context "context"
        grpc "google.golang.org/grpc"
        codes "google.golang.org/grpc/codes"
        status "google.golang.org/grpc/status"
)
```

This is the `import` block—the reason for having `context "context"` is that `context` used to be an external Go package that was not a part of the standard Go library.

The first lines of `protoapi.pb.go` are the following:

```markup
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
//      protoc-gen-go v1.27.1
//      protoc        v3.17.3
// source: protoapi.proto
package protoapi
```

Both `protoapi_grpc.pb.go`, and `protoapi.pb.go` are part of the `protoapi` Go package, which means that we only need to include them once in our code.

The next section is about the development of the gRPC server.

Just Imagine

# Developing a gRPC server

In this section, we are going to create a gRPC server based on the `api.proto` file presented in the previous section. As gRPC needs external packages, we are going to need a GitHub repository to host the files, which is going to be [https://github.com/mactsouk/grpc](https://github.com/mactsouk/grpc).

The code of `gServer.go` (located in the `server` directory) that is related to gRPC (some functions were omitted for brevity) is the following:

```markup
package main
import (
    "context"
    "fmt"
    "math/rand"
    "net"
    "os"
    "time"
```

This program uses `math/rand` instead of the more secure `crypto/rand` for generating random numbers because we need a seed value to be able to reproduce random number sequences.

```markup
    "github.com/mactsouk/protoapi"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)
```

The `import` block includes external Google packages as well as `github.com/mactsouk/protoapi`, which is the one that we created earlier. The `protoapi` package contains structures, interfaces, and functions that are specific to the gRPC service that is being developed, whereas the external Google packages contain generic code related to gRPC.

```markup
type RandomServer struct {
    protoapi.UnimplementedRandomServer
}
```

This structure is named after the name of the gRPC service. This structure is going to implement the interface required by the gRPC server. The use of `protoapi.UnimplementedRandomServer` is required for the implementation of the interface—this is standard practice.

```markup
func (RandomServer) GetDate(ctx context.Context, r *protoapi.RequestDateTime) (*protoapi.DateTime, error) {
    currentTime := time.Now()
    response := &protoapi.DateTime{
        Value: currentTime.String(),
    }
    return response, nil
}
```

This is the first method of the interface named after the `GetDate` function found in the `service` block of `protoapi.proto`. This method requires no input from the client, so it ignores the `r` parameter.

```markup
func (RandomServer) GetRandom(ctx context.Context, r *protoapi.RandomParams) (*protoapi.RandomInt, error) {
    rand.Seed(r.GetSeed())
    place := r.GetPlace()
```

The `GetSeed()` and `GetPlace()` get methods are implemented by `protoc`, are related to the fields of `protoapi.RandomParams`, and should be used in order to read data from the client message.

```markup
    temp := random(min, max)
    for {
        place--
        if place <= 0 {
            break
        }
        temp = random(min, max)
    }
    response := &protoapi.RandomInt{
        Value: int64(temp),
    }
    return response, nil
}
```

The server constructs a `protoapi.RandomInt` variable that is going to be returned to the client. We end the implementation of the second method of the interface here.

```markup
func (RandomServer) GetRandomPass(ctx context.Context, r *protoapi.RequestPass) (*protoapi.RandomPass, error) {
    rand.Seed(r.GetSeed())
    temp := getString(r.GetLength())
```

The `GetSeed()` and `GetLength()` get methods are implemented by `protoc`, are related to the fields of `protoapi.RequestPass`, and should be used in order to read the data from the client message.

```markup
    response := &protoapi.RandomPass{
        Password: temp,
    }
    return response, nil
}
```

In the last part of `GetRandomPass()`, we construct the response (`protoapi.RandomPass`) in order to send it to the client.

```markup
var port = ":8080"
func main() {
    if len(os.Args) == 1 {
        fmt.Println("Using default port:", port)
    } else {
        port = os.Args[1]
    }
```

The first part of `main()` is about specifying the TCP port that is going to be used for the service.

```markup
    server := grpc.NewServer()
```

The previous statement creates a new gRPC server that is not attached to any specific gRPC service.

```markup
    var randomServer RandomServer
    protoapi.RegisterRandomServer(server, randomServer)
```

The previous statements call `protoapi.RegisterRandomServer()` to create a gRPC server for our specific service.

```markup
    reflection.Register(server)
```

It is not mandatory to call `reflection.Register()`, but it helps when you want to list the available services found on a server—in our case, it could have been omitted.

```markup
    listen, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println(err)
        return
    }
```

The previous code starts a TCP service that listens to the desired TCP port.

```markup
    fmt.Println("Serving requests...")
    server.Serve(listen)
}
```

The last part of the program is about telling the gRPC server to begin serving client requests. This happens by calling the `Serve()` method and using the network parameters stored in `listen`.

The `curl(1)` utility does not work with binary data and therefore cannot be used for testing a gRPC server. However, there is an alternative for testing gRPC services, which is the `grpcurl` utility ([https://github.com/fullstorydev/grpcurl](https://github.com/fullstorydev/grpcurl)).

Now that we have the gRPC server ready, let us continue by developing the client that can help us test the operation of the gRPC server.

Just Imagine

# Developing a gRPC client

This section presents the development of the gRPC client based on the `api.proto` file presented earlier. The main purpose of the client is to test the functionality of the server. However, what is really important is the implementation of three helper functions, each one corresponding to a different RPC call, because these three functions allow you to interact with the gRPC server. The purpose of the `main()` function of `gClient.go` is to use these three helper functions.

So, the code of `gClient.go` is the following:

```markup
package main
import (
    "context"
    "fmt"
    "math/rand"
    "os"
    "time"
    "github.com/mactsouk/protoapi"
    "google.golang.org/grpc"
)
var port = ":8080"
func AskingDateTime(ctx context.Context, m protoapi.RandomClient) (*protoapi.DateTime, error) {
```

You can name the `AskingDateTime()` function anything you want. However, the signature of the function must contain a `context.Context` parameter, as well as a `RandomClient` parameter in order to be able to call `GetDate()` later on. The client does not need to implement any of the functions of the IDL—it just has to call them.

```markup
    request := &protoapi.RequestDateTime{
        Value: "Please send me the date and time",
    }
```

We first construct a `protoapi.RequestDateTime` variable that holds the data for the client request.

```markup
    return m.GetDate(ctx, request)
}
```

Then, we call the `GetDate()` method to send the client request to the server. This is handled by the code in the `protoapi` module—we just call `GetDate()` with the correct parameters. This is where the implementation of the first helper function ends—although it is not mandatory to have such a helper function, it makes the code cleaner.

```markup
func AskPass(ctx context.Context, m protoapi.RandomClient, seed int64, length int64) (*protoapi.RandomPass, error) {
    request := &protoapi.RequestPass{
        Seed:   seed,
        Length: length,
    }
```

The `AskPass()` helper function is for calling the `GetRandomPass()` gRPC method in order to get a random password from the server process. First, the function constructs a `protoapi.RequestPass` variable with the given values for `Seed` and `Length`, which are parameters of `AskPass()`.

```markup
    return m.GetRandomPass(ctx, request)
}
```

Then, we call `GetRandomPass()` to send the client request to the server and get the response. Finally, the function returns.

Due to the way gRPC works and the tools that simplify things, the implementation of `AskPass()` is short. Doing the same using RESTful services requires more code.

```markup
func AskRandom(ctx context.Context, m protoapi.RandomClient, seed int64, place int64) (*protoapi.RandomInt, error) {
    request := &protoapi.RandomParams{
        Seed:  seed,
        Place: place,
    }
    return m.GetRandom(ctx, request)
}
```

The last helper function, `AskRandom()`, operates in an analogous way. We construct the client message (`protoapi.RandomParams`), send it to the server by calling `GetRandom()`, and get the server response as returned by `GetRandom()`.

```markup
func main() {
    if len(os.Args) == 1 {
        fmt.Println("Using default port:", port)
    } else {
        port = os.Args[1]
    }
    conn, err := grpc.Dial(port, grpc.WithInsecure())
    if err != nil {
        fmt.Println("Dial:", err)
        return
    }
```

The gRPC client needs to connect to the gRPC server using `grpc.Dial()`. However, we are not done yet as we need to specify the gRPC service the client is going to connect to—this is going to happen in a while. The `grpc.Insecure()` function that is passed as a parameter to `grpc.Dial()` returns a `DialOption` value that disables security for the client connection.

```markup
    rand.Seed(time.Now().Unix())
    seed := int64(rand.Intn(100))
```

Due to the different seed values generated and sent to the gRPC server each time the client code gets executed, we are going to get different random values and passwords back from the gRPC server.

```markup
    client := protoapi.NewRandomClient(conn)
```

Next, we need to create a gRPC client by calling `protoapi.NewRandomClient()` and passing the TCP connection to `protoapi.NewRandomClient()`. This `client` variable is going to be used for all interactions with the server. The name of the function that is called depends on the name of the gRPC service—this allows you to differentiate among the different gRPC services that a machine might support.

```markup
    r, err := AskingDateTime(context.Background(), client)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Server Date and Time:", r.Value)
```

First, we call the `AskingDateTime()` helper function to get the date and time from the gRPC server.

```markup
    length := int64(rand.Intn(20))
    p, err := AskPass(context.Background(), client, 100, length+1)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Random Password:", p.Password)
```

Then, we call `AskPass()` to get a randomly generated password. The length of the password is specified by the `length := int64(rand.Intn(20))` statement.

```markup
    place := int64(rand.Intn(100))
    i, err := AskRandom(context.Background(), client, seed, place)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Random Integer 1:", i.Value)
```

Then, we test `AskRandom()` with different parameters to make sure that it is going to return different values back.

```markup
    k, err := AskRandom(context.Background(), client, seed, place-1)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Random Integer 2:", k.Value)
}
```

With both the server and client completed, we are ready to test them.

## Testing the gRPC server with the client

Now that we have developed both the server and the client, we are ready to use them. First, we should run `gServer.go` as follows:

```markup
$ go run gServer.go
Using default port: :8080
Serving requests..
```

The server process does not produce any other output.

Then, we execute `gClient.go` without any command-line parameters. The output you get from `gClient.go` should be similar to the following:

```markup
$ go run gClient.go
Using default port: :8080
Server Date and Time: 2021-07-05 08:32:19.654905 +0300 EEST m=+2.197816168
Random Password: $1!usiz|36
Random Integer 1: 92
Random Integer 2: 78
```

Apart from the first line, which is the client execution from the UNIX shell, and the second line, which is about the TCP port that is used for connecting to the gRPC server, the next line of the output shows the date and time as returned by the gRPC server. Then we have a random password as generated by the server, as well as two random integers.

If we execute `gClient.go` more than once, we are going to get a different output:

```markup
$ go run gClient.go
Using default port: :8080
Server Date and Time: 2021-07-05 08:32:23.831445 +0300 EEST m=+6.374535148
Random Password: $1!usiz|36N}DO*}{
Random Integer 1: 10
Random Integer 2: 68
```

The fact that the gRPC server returned different values proves that the gRPC server works as expected.

gRPC can do more things than the ones presented in this chapter, such as exchange arrays of messages and streaming—RESTful servers cannot be used for data streaming. However, a discussion of these is beyond the scope of this book.

Just Imagine

# Exercises

-   Convert `gClient.go` into a command-line utility using `cobra`.
-   Try to convert `gServer.go` into a RESTful server.
-   Create a RESTful service that uses gRPC for data exchange. Define the REST API you are going to support, but use gRPC for the communication between the RESTful server and the gRPC server—in this case, the RESTful server is going to act as a gRPC client to the gRPC server.
-   Create your own gRPC service that implements integer addition and subtraction.
-   How easy or difficult would it be to convert the Phone Course application into a gRPC service?
-   Implement a gRPC service that calculates the length of a string.

Just Imagine

# Summary

gRPC is fast, easy to use and understand, and exchanges data in binary format. This chapter taught you how to define the methods and the messages of a gRPC service, how to translate them into Go code, and how to develop a server and a client for that gRPC service.

So, should you use gRPC or stick with RESTful services? Only you can answer that question. You should go with what feels more natural to you. However, if you are still in doubt and cannot decide, begin by developing a RESTful service and then implement the same service using gRPC. After that, you should be ready to choose.

The last chapter of the book is about **generics**, which is a Go feature that is currently under development and is going to be officially included in Go in 2022. However, nothing prohibits us from discussing generics and showing some Go code to better understand generics.

Just Imagine

# Additional resources

-   gRPC: [https://grpc.io/](https://grpc.io/)
-   Protocol Buffers 3 Language Guide: [https://developers.google.com/protocol-buffers/docs/proto3](https://developers.google.com/protocol-buffers/docs/proto3)
-   The `grpcurl` utility: [https://github.com/fullstorydev/grpcurl](https://github.com/fullstorydev/grpcurl)
-   Johan Brandhorst's website: [https://jbrandhorst.com/page/about/](https://jbrandhorst.com/page/about/)
-   The documentation of the `google.golang.org/grpc` package can be found at [https://pkg.go.dev/google.golang.org/grpc](https://pkg.go.dev/google.golang.org/grpc)
-   Go and gRPC tutorial: [https://grpc.io/docs/languages/go/basics/](https://grpc.io/docs/languages/go/basics/)
-   Protocol buffers: [https://developers.google.com/protocol-buffers](https://developers.google.com/protocol-buffers)