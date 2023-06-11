# Building Web Services

The core subject of this chapter is working with HTTP using the `net/http` package—remember that all web services require a web server in order to operate. Additionally, in this chapter, we are going to convert the phone book application into a web application that accepts HTTP connections and create a command-line client to work with it. Lastly, we'll illustrate how to create an **FTP** (**File Transfer Protocol**) server and how to export metrics from Go applications to Prometheus and work with the `runtime/metrics` package to get implementation-defined metrics exported by the Go runtime.

In more detail, this chapter covers:

-   The `net/http` package
-   Creating a web server
-   Updating the phone book application
-   Exposing metrics to Prometheus
-   Developing web clients
-   Creating a client for the phone book service
-   Creating file servers
-   Timing out HTTP connections

Just Imagine

# The net/http package

The `net/http` package offers functions that allow you to develop web servers and clients. For example, `http.Get()` and `http.NewRequest()` are used by clients for making HTTP requests, whereas `http.ListenAndServe()` is used for starting web servers by specifying the IP address and the TCP port the server listens to. Additionally, `http.HandleFunc()` defines supported URLs as well as the functions that are going to handle these URLs.

The next three subsections describe three important data structures of the `net/http` package—you can use these descriptions as a reference while reading this chapter.

## The http.Response type

The `http.Response` structure embodies the response from an HTTP request—both `http.Client` and `http.Transport` return `http.Response` values once the response headers have been received. Its definition can be found at [https://golang.org/src/net/http/response.go](https://golang.org/src/net/http/response.go):

```markup
type Response struct {
    Status     string // e.g. "200 OK"
    StatusCode int    // e.g. 200
    Proto      string // e.g. "HTTP/1.0"
    ProtoMajor int    // e.g. 1
    ProtoMinor int    // e.g. 0
    Header Header
    Body io.ReadCloser 
    ContentLength int64
    TransferEncoding []string
    Close bool
    Uncompressed bool
    Trailer Header 
    Request *Request
    TLS *tls.ConnectionState
}
```

You do not have to use all the structure fields, but it is good to know that they exist. However, some of them, such as `Status`, `StatusCode`, and `Body`, are more important than others. The Go source file, as well as the output of `go doc http.Response`, contains more information about the purpose of each field, which is also the case with most `struct` data types found in the standard Go library.

## The http.Request type

The `http.Request` structure represents an HTTP request as constructed by a client in order to be sent or received by an HTTP server. The public fields of `http.Request` are as follows:

```markup
type Request struct {
    Method string
    URL *url.URL
    Proto  string
    ProtoMajor int
    ProtoMinor int
    Header Header
    Body io.ReadCloser
    GetBody func() (io.ReadCloser, error)
    ContentLength int64
    TransferEncoding []string
    Close bool
    Host string
    Form url.Values
    PostForm url.Values
    MultipartForm *multipart.Form
    Trailer Header
    RemoteAddr string
    RequestURI string
    TLS *tls.ConnectionState
    Cancel <-chan struct{}
    Response *Response
}
```

The `Body` field holds the body of the request. After reading the body of a request, you are allowed to call `GetBody()`, which returns a new copy of the body—this is optional.

Let us now present the `http.Transport` structure.

## The http.Transport type

The definition of `http.Transport`, which gives you more control over your HTTP connections, is fairly long and complex:

```markup
type Transport struct {
    Proxy func(*Request) (*url.URL, error)
    DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
    Dial func(network, addr string) (net.Conn, error)
    DialTLSContext func(ctx context.Context, network, addr string) (net.Conn, error)
    DialTLS func(network, addr string) (net.Conn, error)
    TLSClientConfig *tls.Config
    TLSHandshakeTimeout time.Duration
    DisableKeepAlives bool
    DisableCompression bool
    MaxIdleConns int
    MaxIdleConnsPerHost int
    MaxConnsPerHost int
    IdleConnTimeout time.Duration
    ResponseHeaderTimeout time.Duration
    ExpectContinueTimeout time.Duration
    TLSNextProto map[string]func(authority string, c *tls.Conn) RoundTripper
    ProxyConnectHeader Header
    GetProxyConnectHeader func(ctx context.Context, proxyURL *url.URL, target string) (Header, error)
    MaxResponseHeaderBytes int64
    WriteBufferSize int
    ReadBufferSize int
    ForceAttemptHTTP2 bool
}
```

Note that `http.Transport` is pretty low-level, whereas `http.Client`, which is also used in this chapter, implements a high-level HTTP client—each `http.Client` contains a `Transport` field. If its value is `nil`, then `DefaultTransport` is used. You do not need to use `http.Transport` in all of your programs and you are not required to deal with all of its fields each time you use it. If you want to learn more about `DefaultTransport`, type `go doc http.DefaultTransport`.

Let us now learn how to develop a web server.

Just Imagine

# Creating a web server

This section presents a simple web server developed in Go in order to better understand the principles behind such applications.

Although a web server programmed in Go can do many things efficiently and securely, if what you really need is a powerful web server that supports modules, multiple websites, and virtual hosts, then you would be better off using a web server such as **Apache**, **Nginx**, or **Caddy** that is written in Go.

You might ask why the presented web server uses HTTP instead of **secure HTTP** (**HTTPS**). The answer to this question is simple: most Go web servers are deployed as Docker images and are hidden behind web servers such as Caddy and Nginx that provide the secure HTTP operation part using the appropriate security credentials. It does not make any sense to use the secure HTTP protocol along with the required security credentials without knowing how and under which domain name the application is going to be deployed. This is a common practice in microservices as well as regular web applications that are deployed in Docker images.

The `net/http` package offers functions and data types that allow you to develop powerful web servers and clients. The `http.Set()` and `http.Get()` methods can be used to make HTTP and HTTPS requests, whereas `http.ListenAndServe()` is used for creating web servers given the user-specified handler function or functions that handle incoming requests. As most web services require support for multiple endpoints, you end up needing multiple discrete functions for handling incoming requests, which also leads to the better design of your services.

The simplest way to define the supported endpoints, as well as the handler function that responds to each client request, is with the use of `http.HandleFunc()`, which can be called multiple times.

After this quick and somewhat theoretical introduction, it is time to begin talking about more practical topics, beginning with the implementation of a simple web server as illustrated in `wwwServer.go`:

```markup
package main
import (
    "fmt"
    "net/http"
    "os"
    "time"
)
func myHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
    fmt.Printf("Served: %s\n", r.Host)
}
```

This is a handler function that sends a message back to the client using the `w http.ResponseWriter`, which is also an interface that implements `io.Writer` and is used for sending the server response.

```markup
func timeHandler(w http.ResponseWriter, r *http.Request) {
    t := time.Now().Format(time.RFC1123)
    Body := "The current time is:"
    fmt.Fprintf(w, "<h1 align=\"center\">%s</h1>", Body)
    fmt.Fprintf(w, "<h2 align=\"center\">%s</h2>\n", t)
    fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
    fmt.Printf("Served time for: %s\n", r.Host)
}
```

This is another handler function called `timeHandler` that returns back the current time in HTML format. All `fmt.Fprintf()` calls send data back to the HTTP client whereas the output of `fmt.Printf()` is printed on the terminal the web server runs on. The first argument of `fmt.Fprintf()` is the `w http.ResponseWriter`, which implements `io.Writer` and therefore can accept data.

```markup
func main() {
    PORT := ":8001"
```

This is where you define the port number your web server listens to.

```markup
    arguments := os.Args
    if len(arguments) != 1 {
        PORT = ":" + arguments[1]
    }
    fmt.Println("Using port number: ", PORT)
```

If you do not want to use the predefined port number (`8001`), then you should provide `wwwServer.go` with your own port number as a command-line argument.

```markup
    http.HandleFunc("/time", timeHandler)
    http.HandleFunc("/", myHandler)
```

So, the web server supports the `/time` URL as well as `/`. The `/` path **matches every URL** not matched by other handlers. The fact that we associate `myHandler()` with `/` makes `myHandler()` the default handler function.

```markup
    err := http.ListenAndServe(PORT, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

The `http.ListenAndServe()` call begins the HTTP server using the predefined port number. As there is no hostname given in the `PORT` string, the web server is going to listen to all available network interfaces. The port number and the hostname should be separated with a colon (`:`), which should be there even if there is no hostname—in that case the server listens to all available network interfaces and, therefore, all supported hostnames. This is the reason that the value of `PORT` is `:8001` instead of just `8001`.

Part of the `net/http` package is the `ServeMux` type (`go doc http.ServeMux`), which is an **HTTP request multiplexer** that provides a slightly different way of defining handler functions and endpoints than the default one, which is used in `wwwServer.go`. So, if we do not create and configure our own `ServeMux` variable, then `http.HandleFunc()` uses `DefaultServeMux`, which is the default `ServeMux`. So, in this case we are going to implement the web service using the **default Go router**—this is the reason that the second parameter of `http.ListenAndServe()` is `nil`.

Running `wwwServer.go` and interacting with it using `curl(1)` produces the next output:

```markup
$ go run wwwServer.go
Using port number:  :8001
Served: localhost:8001
Served time for: localhost:8001
Served: localhost:8001
```

Note that as `wwwServer.go` does not terminate automatically, you need to stop it on your own.

On the `curl(1)` side, the interaction looks as follows:

```markup
$ curl localhost:8001
Serving: /
```

In this first case, we visit the `/` path of the web server and we are being served by `myHandler()`.

```markup
$ curl localhost:8001/time
<h1 align="center">The current time is:</h1><h2 align="center">Mon, 29 Mar 2021 08:26:27 EEST</h2>
Serving: /time
```

In this case we visit `/time` and we get HTML output back from `timeHandler()`.

```markup
$ curl localhost:8001/doesNotExist
Serving: /doesNotExist
```

In this last case, we visit `/doesNotExist`, which does not exist. As this cannot be matched by any other path, it is served by the default handler, which is the `myHandler()` function.

The next section is about making the phone book application a web application!

Just Imagine

# Updating the phone book application

This time the phone book application is going to work as a web service. The two main tasks that need to be performed are defining the API along with the endpoints and implementing the API. A third task that needs to be determined concerns **data exchange** between the application server and its clients. There exist four main approaches regarding data exchange:

-   Using plain text
-   Using HTML
-   Using JSON
-   Using a hybrid approach that combines plain text and JSON data

As JSON is explored in _Chapter 10_, _Working with REST APIs_, and HTML might not be the best option for a service because you need to separate the data from the HTML tags and parse the data, we are going to use the first approach. Therefore, the service is going to work with **plain text data**. We begin by defining the API that supports the operation of the phone book application.

## Defining the API

The API has support for the following URLs:

-   `/list`: This lists all available entries.
-   `/insert/name/surname/telephone/`: This inserts a new entry. Later on, we are going to see how to extract the desired information from a URL that contains user data.
-   `/delete/telephone/`: This deletes an entry based on the value of `telephone`.
-   `/search/telephone/`: This searches for an entry based on the value of `telephone`.
-   `/status`: This is an extra URL that returns the number of entries in the phone book.

The list of endpoints does not follow standard REST conventions—all these are going to be presented in _Chapter 10_, _Working with REST APIs_.

This time we not using the default Go router, which means that we define and configure our own `http.NewServeMux()` variable. This changes the way we provide handler functions: a handler function with the `func(http.ResponseWriter, *http.Request)` signature has to be converted into an `http.HandlerFunc` **type** and be used by the `ServeMux` type and its own `Handle()` method. Therefore, when using a different `ServeMux` than the default one, we should do that conversion explicitly by calling `http.HandlerFunc()`, which makes the `http.HandlerFunc` **type** act as an **adapter** that allows the use of ordinary functions as HTTP handlers provided that they have the required signature. This is not a problem when using the default Go router (`DefaultServeMux`) because the `http.HandleFunc()` **function** does that conversion automatically and internally.

To make things clearer, the `http.HandlerFunc` **type** has support for a **method** named `HandlerFunc()`—both the type and method are defined in the `http` package. The similarly named the `http.HandleFunc()` **function** (without an `r`) is used with the default Go router.

As an example, for the `/time` endpoint and the `timeHandler()` handler function, you should call `mux.Handle()` as `mux.Handle("/time", http.HandlerFunc(timeHandler))`. If you were using `http.HandleFunc()` and as a consequence `DefaultServeMux`, then you should have called `http.HandleFunc("/time", timeHandler)` instead.

The subject of the next subsection is the implementation of the HTTP endpoints.

## Implementing the handlers

The new version of the phone book is going to be created on a dedicated GitHub repository for storing and sharing it: [https://github.com/mactsouk/www-phone](https://github.com/mactsouk/www-phone). After creating the repository, you need to do the following:

```markup
$ cd ~/go/src/github.com/mactsouk # Replace with your own path
$ git clone git@github.com:mactsouk/www-phone.git
$ cd www-phone
$ touch handlers.go
$ touch www-phone.go
```

The `www-phone.go` file holds the code that defines the operation of the web server. Usually, handlers are put in a separate package, but for reasons of simplicity, we decided to put handlers in a separate file within the same package named `handlers.go`. The contents of the `handlers.go` file, which contains all functionality related to the serving of the clients, are the following:

```markup
package main
import (
    "fmt"
    "log"
    "net/http"
    "strings"
)
```

All required packages for `handlers.go` are imported even if some of them are already imported by `www-phone.go`. Note that the name of the package is `main`, which is also the case for `www-phone.go`.

```markup
const PORT = ":1234"
```

This is the default port number the HTTP server listens to.

```markup
func defaultHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Serving:", r.URL.Path, "from", r.Host)
    w.WriteHeader(http.StatusOK)
    Body := "Thanks for visiting!\n"
    fmt.Fprintf(w, "%s", Body)
}
```

This is the default handler, which serves all requests that are not a match for any of the other handlers.

```markup
func deleteHandler(w http.ResponseWriter, r *http.Request) {
    // Get telephone
    paramStr := strings.Split(r.URL.Path, "/")
```

This is the handler function for the `/delete` path, which begins by splitting the URL in order to read the desired information.

```markup
    fmt.Println("Path:", paramStr)
    if len(paramStr) < 3 {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintln(w, "Not found: "+r.URL.Path)
        return
    }
```

If we do not have enough parameters, we should send an error message back to the client with the desired HTTP code, which in this case is `http.StatusNotFound`. You can use any HTTP code you want as long as it makes sense. The `WriteHeader()` method sends back a header with the provided status code before writing the body of the response.

```markup
    log.Println("Serving:", r.URL.Path, "from", r.Host)
```

This is where the HTTP server sends data to log files—this mainly happens for debugging reasons.

```markup
    telephone := paramStr[2]
```

As the delete process is based on the telephone number, all that is required is a valid telephone number. This is where the parameter is read after splitting the provided URL.

```markup
    err := deleteEntry(telephone)
    if err != nil {
        fmt.Println(err)
        Body := err.Error() + "\n"
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "%s", Body)
        return
    }
```

Once we have a telephone number, we call `deleteEntry()` in order to delete it. The return value of `deleteEntry()` determines the result of the operation and, therefore, the response to the client.

```markup
    Body := telephone + " deleted!\n"
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "%s", Body)
}
```

At this point, we know that the delete operation was successful so we send a proper message to the client as well as the `http.StatusOK` status code. Type `go doc http.StatusOK` for the list of codes.

```markup
func listHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Serving:", r.URL.Path, "from", r.Host)
    w.WriteHeader(http.StatusOK)
    Body := list()
    fmt.Fprintf(w, "%s", Body)
}
```

The `list()` helper function that is used in the `/list` path cannot fail. Therefore, `http.StatusOK` is always returned when serving `/list`. However, sometimes the return value of `list()` can be empty.

```markup
func statusHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Serving:", r.URL.Path, "from", r.Host)
    w.WriteHeader(http.StatusOK)
    Body := fmt.Sprintf("Total entries: %d\n", len(data))
    fmt.Fprintf(w, "%s", Body)
}
```

The preceding code defines the handler function for the `/status` URL. It just returns information about the total number of entries found in the phone book. It can be used for verifying that the web service works fine.

```markup
func insertHandler(w http.ResponseWriter, r *http.Request) {
    // Split URL
    paramStr := strings.Split(r.URL.Path, "/")
    fmt.Println("Path:", paramStr)
```

As happened with `delete`, we need to split the given URL in order to extract the information. In this case, we need three elements as we are trying to insert a new entry into the phone book application.

```markup
    if len(paramStr) < 5 {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintln(w, "Not enough arguments: "+r.URL.Path)
        return
    }
```

Needing to extract three elements from the URL means that we require `paramStr` to have at least four elements in it, hence the `len(paramStr) < 5` condition.

```markup
    name := paramStr[2]
    surname := paramStr[3]
    tel := paramStr[4]
    t := strings.ReplaceAll(tel, "-", "")
    if !matchTel(t) {
        fmt.Println("Not a valid telephone number:", tel)
        return
    }
```

In the previous part, we get the desired data and make sure that the telephone number contains digits only—this happens with the use of the `matchTel()` helper function.

```markup
    temp := &Entry{Name: name, Surname: surname, Tel: t}
    err := insert(temp)
```

As the `insert()` helper function requires an `*Entry` value, we create one before calling it.

```markup
    if err != nil {
        w.WriteHeader(http.StatusNotModified)
        Body := "Failed to add record\n"
        fmt.Fprintf(w, "%s", Body)
    } else {
        log.Println("Serving:", r.URL.Path, "from", r.Host)
        Body := "New record added successfully\n"
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "%s", Body)
    }
    log.Println("Serving:", r.URL.Path, "from", r.Host)
}
```

This is the end of the handler for `/insert`. The last part of the implementation of `insertHandler()` deals with the return value of `insert()`. If there was not an error, then `http.StatusOK` is returned to the client. In the opposite case, `http.StatusNotModified` is returned to signify that there was not a change in the phone book. It is the job of the client to examine the status code of the interaction but it is the job of the server to send an appropriate status code back to the client.

```markup
func searchHandler(w http.ResponseWriter, r *http.Request) {
    // Get Search value from URL
    paramStr := strings.Split(r.URL.Path, "/")
    fmt.Println("Path:", paramStr)
    if len(paramStr) < 3 {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintln(w, "Not found: "+r.URL.Path)
        return
    }
    var Body string
    telephone := paramStr[2]
```

At this point, we extract the telephone number from the URL as we did with `/delete`.

```markup
    t := search(telephone)
    if t == nil {
        w.WriteHeader(http.StatusNotFound)
        Body = "Could not be found: " + telephone + "\n"
    } else {
        w.WriteHeader(http.StatusOK)
        Body = t.Name + " " + t.Surname + " " + t.Tel + "\n"
    }
    fmt.Println("Serving:", r.URL.Path, "from", r.Host)
    fmt.Fprintf(w, "%s", Body)
}
```

The last function of `handlers.go` ends here and is about the `/search` endpoint. The `search()` helper function checks whether the given input exists in the phone book records or not and acts accordingly. Additionally, the implementation of the `main()` function, which can be found in `www-phone.go`, is the following:

```markup
func main() {
    err := readCSVFile(CSVFILE)
    if err != nil {
        fmt.Println(err)
        return
    }
    err = createIndex()
    if err != nil {
        fmt.Println("Cannot create index.")
        return
    }
```

This first part of `main()` has to do with the initialization of the phone book application.

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

Here, we store the parameters of the HTTP server in the `http.Server` structure and use our own `http.NewServeMux()` instead of the default one.

```markup
    mux.Handle("/list", http.HandlerFunc(listHandler))
    mux.Handle("/insert/", http.HandlerFunc(insertHandler))
    mux.Handle("/insert", http.HandlerFunc(insertHandler))
    mux.Handle("/search", http.HandlerFunc(searchHandler))
    mux.Handle("/search/", http.HandlerFunc(searchHandler))
    mux.Handle("/delete/", http.HandlerFunc(deleteHandler))
    mux.Handle("/status", http.HandlerFunc(statusHandler))
    mux.Handle("/", http.HandlerFunc(defaultHandler))
```

This is the list of the supported URLs. Note that `/search` and `/search/` are both handled by the same handler function even though `/search` is going to fail as it does not include the required argument. On the other hand, `/delete/` is handled in a special way—this is going to be shown when testing the application. As we are using `http.NewServeMux()` and not the default Go router, we need to use `http.HandlerFunc()` when defining the handler functions.

```markup
    fmt.Println("Ready to serve at", PORT)
    err = s.ListenAndServe()
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

The `ListenAndServe()` method starts the HTTP server using the parameters defined previously in the `http.Server` structure. The rest of `www-phone.go` contains helper functions related to the operation of the phone book. Note that it is important to save and update the contents of the phone book application as often as possible because this is a live application, and you might lose data if it crashes.

The next command allows you to execute the application—you need to provide both files in `go run`:

```markup
$ go run www-phone.go handlers.go
Ready to serve at :1234
2021/03/29 17:13:49 Serving: /list from localhost:1234
2021/03/29 17:13:53 Serving: /status from localhost:1234
Path: [ search 2109416471]
Serving: /search/2109416471 from localhost:1234
Path: [ search]
2021/03/29 17:28:34 Serving: /list from localhost:1234
Path: [ search 2101112223]
Serving: /search/2101112223 from localhost:1234
Path: [ delete 2109416471]
2021/03/29 17:29:24 Serving: /delete/2109416471 from localhost:1234
Path: [ insert Mike Tsoukalos 2109416471]
2021/03/29 17:29:56 Serving: /insert/Mike/Tsoukalos/2109416471 from localhost:1234
2021/03/29 17:29:56 Serving: /insert/Mike/Tsoukalos/2109416471 from localhost:1234
Path: [ insert Mike Tsoukalos 2109416471]
2021/03/29 17:30:18 Serving: /insert/Mike/Tsoukalos/2109416471 from localhost:1234
```

On the client side, which is `curl(1)`, we have the next output:

```markup
$ curl localhost:1234/list
Dimitris Tsoukalos 2101112223
Jane Doe 0800123456
Mike Tsoukalos 2109416471
```

Here, we get all entries from the phone book application by visiting `/list`.

```markup
$ curl localhost:1234/status
Total entries: 3
```

Next, we visit `/status` and get back the expected output.

```markup
$ curl localhost:1234/search/2109416471
Mike Tsoukalos 2109416471
```

The previous command searches for an existing phone number—the server responds with its full record.

```markup
$ curl localhost:1234/delete/2109416471
2109416471 deleted!
```

The previous output shows that we have deleted the record with telephone number `2109416471`. In REST, this requires a `DELETE` method, but for reasons of simplicity we leave the details for _Chapter 10_, _Working with REST APIs_.

Now, let us try and visit `/delete` instead of `/delete/`:

```markup
$ curl localhost:1234/delete
<a href="/delete/">Moved Permanently</a>.
```

The presented message was generated by the Go router and tells us that we should try `/delete/` instead as `/delete` was moved permanently. This is the kind of message that we get by not specifically defining both `/delete` and `/delete/` in the routes.

Now, let us insert a new record:

```markup
$ curl localhost:1234/insert/Mike/Tsoukalos/2109416471
New record added successfully
```

In REST, this requires a `POST` method, but again, we will leave that for _Chapter 10_, _Working with REST APIs_.

If we try to insert the same record again, the response is going to be as follows:

```markup
$ curl localhost:1234/insert/Mike/Tsoukalos/2109416471
Failed to add record
```

Everything looks like it is working OK. We can now put the phone application online and interact with it using multiple HTTP requests as the `http` package uses multiple goroutines for interacting with clients—in practice, this means that the phone book application **runs concurrently**!

Later in this chapter we are going to create a command-line client for the phone book server. Additionally, _Chapter 11_, _Code Testing and Profiling_, shows how to test your code.

The next section shows how to expose metrics to Prometheus and how to build Docker images for server applications.

Just Imagine

# Exposing metrics to Prometheus

ImagineDevOps  that you have an application that writes files to disk and you want to get metrics for that application to better understand how the writing of multiple files has an effect on the general performance—you need to gather performance data for understanding the behavior of your application. Although the presented application uses the _gauge_ type of metric only because it is what is appropriate for the information that is sent to Prometheus, Prometheus accepts many types of data. The list of supported data types for metrics is the following:

-   _Counter_: This is a cumulative value that is used for representing increasing counters—the value of a counter can stay the same, go up, or be reset to zero but cannot decrease. Counters are usually used for representing cumulative values such as the number of requests served so far, the total number of errors, etc.
-   _Gauge_: This is a single numerical value that is allowed to increase or decrease. Gauges are usually used for representing values that can go up or down such as the number of requests, time durations, etc.
-   _Histogram_: A histogram is used for sampling observations and creating counts and buckets. Histograms are usually used for counting request durations, response times, etc.
-   _Summary_: A summary is like a histogram but can also calculate quantiles over sliding windows that work with times.

Both histograms and summaries are useful and handy for performing statistical calculations and properties. Usually, a counter or a gauge is all that you need for storing your system metrics.

This section shows how you can collect and expose a system to Prometheus. For reasons of simplicity, the presented application is going to generate random values. We begin by explaining the use of the `runtime/metrics` package, which provides Go runtime-related metrics.

## The runtime/metrics package

The `runtime/metrics` package makes metrics exported by the Go runtime available to the developer. Each metric name is specified by a path. As an example, the number of live goroutines is accessed as `/sched/goroutines:goroutines`. However, if you want to collect all available metrics, you should use `metrics.All()`—this saves you from having to write lots of code in order to collect all metrics manually.

Metrics are saved using the `metrics.Sample` data type. The definition of the `metrics.Sample` data structure is as follows:

```markup
type Sample struct {
    Name string
    Value Value
}
```

The `Name` value must correspond to the name of one of the metric descriptions returned by `metrics.All()`. If you already know the metric description, there is no need to use `metrics.All()`.

The use of the `runtime/metrics` package is illustrated in `metrics.go`. The presented code gets the value of `/sched/goroutines:goroutines` and prints it on screen:

```markup
package main
import (
    "fmt"
    "runtime/metrics"
    "sync"
    "time"
)
func main() {
    const nGo = "/sched/goroutines:goroutines"
```

The `nGo` variable holds the path of the metric we want to collect.

```markup
    // A slice for getting metric samples
    getMetric := make([]metrics.Sample, 1)
    getMetric[0].Name = nGo
```

After that we create a slice of type `metrics.Sample` in order to keep the metric value. The initial size of the slice is `1` because we are only collecting values for a single metric. We set the `Name` value to `/sched/goroutines:goroutines` as stored in `nGo`.

```markup
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            time.Sleep(4 * time.Second)
        }()
```

Here we manually create three goroutines to have relevant data to collect.

```markup
        // Get actual data
        metrics.Read(getMetric)
        if getMetric[0].Value.Kind() == metrics.KindBad {
            fmt.Printf("metric %q no longer supported\n", nGo)
        }
```

The `metrics.Read()` function collects the desired metrics based on the data in the `getMetric` slice.

```markup
        mVal := getMetric[0].Value.Uint64()
        fmt.Printf("Number of goroutines: %d\n", mVal)
    }
```

After reading the desired metric, we convert it into a numeric value (unsigned `int64` here) in order to use it in our program.

```markup
    wg.Wait()
    metrics.Read(getMetric)
    mVal := getMetric[0].Value.Uint64()
    fmt.Printf("Before exiting: %d\n", mVal)
}
```

The last lines of the code verify that after all goroutines have finished, the value of the metric is going to be `1`, which is the goroutine used for running the `main()` function.

Running `metrics.go` produces the next output:

```markup
$ go run metrics.go
Number of goroutines: 2
Number of goroutines: 3
Number of goroutines: 4
Before exiting: 1
```

We have created 3 goroutines and we already have a goroutine for running the `main()` function. Therefore, the maximum number of goroutines is indeed 4.

The subsections that follow illustrate how to make any metric you collect available to Prometheus.

## Exposing metrics

Collecting metrics is a totally different task from exposing them for Prometheus to collect them. This subsection shows how to make the metrics available for collection.

The code of `samplePro.go` is as follows:

```markup
package main
import (
    "fmt"
    "net/http"
    "math/rand"
    "time"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)
```

We need to use two external packages for communicating with Prometheus.

```markup
var PORT = ":1234"
var counter = prometheus.NewCounter(
    prometheus.CounterOpts{
        Namespace: "mtsouk",
        Name:      "my_counter",
        Help:      "This is my counter",
    })
```

This is how we define a new `counter` variable and specify the desired options. The `Namespace` field is very important as it allows you to group metrics in sets.

```markup
var gauge = prometheus.NewGauge(
    prometheus.GaugeOpts{
        Namespace: "mtsouk",
        Name:      "my_gauge",
        Help:      "This is my gauge",
    })
```

This is how we define a new `gauge` variable and specify the desired options.

```markup
var histogram = prometheus.NewHistogram(
    prometheus.HistogramOpts{
        Namespace: "mtsouk",
        Name:      "my_histogram",
        Help:      "This is my histogram",
    })
```

This is how we define a new `histogram` variable and specify the desired options.

```markup
var summary = prometheus.NewSummary(
    prometheus.SummaryOpts{
        Namespace: "mtsouk",
        Name:      "my_summary",
        Help:      "This is my summary",
    })
```

This is how we define a new `summary` variable and specify the desired options. However, as you are going to see, defining a metric variable is not enough. You also need to register it.

```markup
func main() {
    rand.Seed(time.Now().Unix())
    prometheus.MustRegister(counter)
    prometheus.MustRegister(gauge)
    prometheus.MustRegister(histogram)
    prometheus.MustRegister(summary)
```

In these four statements, you register the four metric variables. Now Prometheus knows about them.

```markup
    go func() {
        for {
            counter.Add(rand.Float64() * 5)
            gauge.Add(rand.Float64()*15 - 5)
            histogram.Observe(rand.Float64() * 10)
            summary.Observe(rand.Float64() * 10)
            time.Sleep(2 * time.Second)
        }
    }()
```

This goroutine runs for as long as the web server runs with the help of the endless `for` loop. In this goroutine, the metrics are updated every 2 seconds due to the use of the `time.Sleep(2 * time.Second)` statement—in this case using random values.

```markup
    http.Handle("/metrics", promhttp.Handler())
    fmt.Println("Listening to port", PORT)
    fmt.Println(http.ListenAndServe(PORT, nil))
}
```

As you already know, each URL is handled by a handler function that you usually implement on your own. However, in this case we are using the `promhttp.Handler()` handler function that comes with the `github.com/prometheus/client_golang/prometheus/promhttp` package—this saves us from having to write our own code. However, we still need to register the `promhttp.Handler()` handler function using `http.Handle()` before we start the web server. Note that the metrics are found under the `/metrics` path—Prometheus knows how to find that.

With `samplePro.go` running, getting the list of metrics that belong to the `mtsouk` namespace is as simple as running the next `curl(1)` command:

```markup
$ curl localhost:1234/metrics --silent | grep mtsouk
# HELP mtsouk_my_counter This is my counter
# TYPE mtsouk_my_counter counter
mtsouk_my_counter 19.948239343027772
```

This is the output from a `counter` variable. If the `| grep mtsouk` part is omitted, then you are going to get the list of all available metrics.

```markup
# HELP mtsouk_my_gauge This is my gauge
# TYPE mtsouk_my_gauge gauge
mtsouk_my_gauge 29.335329668135287
```

This is the output from a `gauge` variable.

```markup
# HELP mtsouk_my_histogram This is my histogram
# TYPE mtsouk_my_histogram histogram
mtsouk_my_histogram_bucket{le="0.005"} 0
mtsouk_my_histogram_bucket{le="0.01"} 0
mtsouk_my_histogram_bucket{le="0.025"} 0
. . .
mtsouk_my_histogram_bucket{le="5"} 4
mtsouk_my_histogram_bucket{le="10"} 9
mtsouk_my_histogram_bucket{le="+Inf"} 9
mtsouk_my_histogram_sum 44.52262035556937
mtsouk_my_histogram_count 9
```

This is the output from a `histogram` variable. Histograms contain _buckets_, hence the large number of output lines.

```markup
# HELP mtsouk_my_summary This is my summary
# TYPE mtsouk_my_summary summary
mtsouk_my_summary_sum 19.407554729772105
mtsouk_my_summary_count 9
```

The last lines of the output are for the `summary` data type.

So, the metrics are there and ready to be pulled by Prometheus—in practice, this means that every production Go application can export metrics that can be used for measuring its performance and discovering its bottlenecks. However, we are not done yet as we need to learn about building Docker images for Go applications.

### Creating a Docker image for a Go server

This section shows how to create a Docker image for a Go application. The main benefit you get from this is that you can deploy it in a Docker environment without worrying about compiling it and having the required resources—everything is included in the Docker image.

Still, you might ask, _"why not use a normal Go binary instead of a Docker image?"_ The answer is simple: Docker images can be put in `docker-compose.yml` files and can be deployed using Kubernetes. The same is not true about Go binaries.

You usually start with a base Docker image that already includes Go and you create the desired binary in there. The key point here is that `samplePro.go` uses an external package that should be downloaded in the Docker image before building the executable binary.

The process must start with `go mod init` and `go mod tidy`. The contents of the Dockerfile, which can be found as `dFilev2` in the GitHub repository of the book, are as follows:

```markup
# WITH Go Modules
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
```

As `golang:alpine` uses the latest Go version, which does not come with `git`, we install it manually.

```markup
RUN mkdir $GOPATH/src/server
ADD ./samplePro.go $GOPATH/src/server
```

If you want to use Go modules, you should put your code in `$GOPATH/src`.

```markup
WORKDIR $GOPATH/src/server
RUN go mod init
RUN go mod tidy
RUN go mod download
RUN mkdir /pro
RUN go build -o /pro/server samplePro.go
```

We download dependencies using various `go mod` commands. The building of the binary file is the same as before.

```markup
FROM alpine:latest
RUN mkdir /pro
COPY --from=builder /pro/server /pro/server
EXPOSE 1234
WORKDIR /pro
CMD ["/pro/server"]
```

In this second stage, we put the binary file into the desired location (`/pro`) and expose the desired port, which in this case is `1234`. The port number depends on the code in `samplePro.go`.

Building a Docker image using `dFilev2` is as simple as running the next command:

```markup
$ docker build -f dFilev2 -t go-app116 .
```

Once the Docker image has been created, there is no difference in the way you should use it in a `docker-compose.yml` file—a relevant entry in a `docker-compose.yml` file would look as follows:

```markup
  goapp:
    image: goapp
    container_name: goapp-int
    restart: always
    ports:
      - 1234:1234
    networks:
      - monitoring
```

The name of the Docker image is `goapp` whereas the internal name of the container would be `goapp-int`. So, if a different container from the `monitoring` network wants to access that container, it should use the `goapp-int` hostname. Last, the only open port is port number `1234`.

The next section illustrates how to expose metrics to Prometheus.

### Exposing the desired metrics

This section illustrates how to expose metrics from the `runtime/metrics` package to Prometheus. In our case we use `/sched/goroutines:goroutines` and `/memory/classes/total:bytes`. You already know about the former, which is the total number of goroutines. The latter metric is the amount of memory mapped by the Go runtime into the current process as read-write.

As the presented code uses an external package, it should be put inside `~/go/src` and Go modules should be enabled using `go mod init`.

The Go code of `prometheus.go` is as follows:

```markup
package main
import (
    "log"
    "math/rand"
    "net/http"
    "runtime"
    "runtime/metrics"
    "time"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)
```

The first external package is the Go client library for Prometheus and the second package is for using the default handler function (`promhttp.Handler()`).

```markup
var PORT = ":1234"
var n_goroutines = prometheus.NewGauge(
    prometheus.GaugeOpts{
        Namespace: "packt",
        Name:      "n_goroutines",
        Help:      "Number of goroutines"})
var n_memory = prometheus.NewGauge(
    prometheus.GaugeOpts{
        Namespace: "packt",
        Name:      "n_memory",
        Help:      "Memory usage"})
```

Here, we define the two Prometheus metrics.

```markup
func main() {
    rand.Seed(time.Now().Unix())
    prometheus.MustRegister(n_goroutines)
    prometheus.MustRegister(n_memory)
    const nGo = "/sched/goroutines:goroutines"
    const nMem = "/memory/classes/heap/free:bytes"
```

This is where you define the metrics you want to read from the `runtime/metrics` package.

```markup
    getMetric := make([]metrics.Sample, 2)
    getMetric[0].Name = nGo
    getMetric[1].Name = nMem
    http.Handle("/metrics", promhttp.Handler())
```

This is where you register the handler function for the `/metrics` path. We use `promhttp.Handler()`.

```markup
    go func() {
        for {
            for i := 1; i < 4; i++ {
                go func() {
                    _ = make([]int, 1000000)
                    time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
                }()
            }
```

Note that such a program should definitely have **at least two goroutines**: one for running the HTTP server and another one for collecting the metrics. Usually, the HTTP server is on the goroutine that runs the `main()` function and the metric collection happens in a user-defined goroutine.

The external `for` loop makes sure that the goroutine runs forever whereas the internal `for` loop creates additional goroutines so that the value of the `/sched/goroutines:goroutines` metric changes all the time.

```markup
            runtime.GC()
            metrics.Read(getMetric)
            goVal := getMetric[0].Value.Uint64()
            memVal := getMetric[1].Value.Uint64()
            time.Sleep(time.Duration(rand.Intn(15)) * time.Second)
            n_goroutines.Set(float64(goVal))
            n_memory.Set(float64(memVal))
        }
    }()
```

The `runtime.GC()` function tells the Go garbage collector to run and is called for changing the `/memory/classes/heap/free:bytes` metric. The two `Set()` calls update the values of the metrics.

You can read more about the operation of the Go garbage collector in _Appendix A_.

```markup
    log.Println("Listening to port", PORT)
    log.Println(http.ListenAndServe(PORT, nil))
}
```

The last statement runs the web server using the default Go router. Running `prometheus.go` from a directory inside `~/go/src/github.com/mactsouk` requires executing the next commands:

```markup
$ cd ~/go/src/github.com/mactsouk/Prometheus # use any path inside ~/go/src
$ go mod init
$ go mod tidy
$ go mod download
$ go run prometheus.go
2021/04/01 12:18:11 Listening to port :1234
```

Although `prometheus.go` generates no output apart from the previous line, the next subsection illustrates how to read the desired metrics from it using `curl(1)`.

## Reading metrics

You can get a list of the metrics from `prometheus.go` using `curl(1)` in order to make sure that the application works as expected. I always test the operation of such an application with `curl(1)` or some other similar utility such as `wget(1)` before trying to get the metrics with Prometheus.

```markup
$ curl localhost:1234/metrics --silent | grep packt
# HELP packt_n_goroutines Number of goroutines
# TYPE packt_n_goroutines gauge
packt_n_goroutines 5
# HELP packt_n_memory Memory usage
# TYPE packt_n_memory gauge
packt_n_memory 794624
```

The previous command assumes that `curl(1)` is executed on the same machine as the server and that the server listens to TCP port number `1234`. Next, we must enable Prometheus to pull the metrics. The easiest way for a Prometheus Docker image to be able to see the Go application with the metrics is to execute both as Docker images. There is an important point here: the `runtime/metrics` package was first introduced with Go version 1.16. This means that to build a Go source file that uses `runtime/metrics`, we need to use Go version 1.16 or newer, which means that we need to use modules for building the Docker image. Therefore, we are going to use the following Dockerfile:

```markup
FROM golang:alpine AS builder
```

This is the name of the base Docker image that is used for building the binary. `golang:alpine` always contains the latest Go version as long as you update it regularly.

```markup
RUN apk update && apk add --no-cache git
```

As `golang:alpine` does not come with `git`, we need to install it manually.

```markup
RUN mkdir $GOPATH/src/server
ADD ./prometheus.go $GOPATH/src/server
WORKDIR $GOPATH/src/server
RUN go mod init
RUN go mod tidy
RUN go mod download
```

The previous commands download the required dependencies before trying to build the binary.

```markup
RUN mkdir /pro
RUN go build -o /pro/server prometheus.go
FROM alpine:latest
RUN mkdir /pro
COPY --from=builder /pro/server /pro/server
EXPOSE 1234
WORKDIR /pro
CMD ["/pro/server"]
```

Building the desired Docker image, which is going to be named `goapp`, is as simple as running the next command:

```markup
$ docker build -f Dockerfile -t goapp .
```

As usual, the output of `docker images` verifies the successful creation of the `goapp` Docker image—in my case the relevant entry looks as follows:

```markup
goapp         latest       a1f0cd4bd8f5   5 seconds ago   16.9MB
```

Let us now discuss how to configure Prometheus to pull the desired metrics.

## Putting the metrics in Prometheus

To be able to pull the metrics, Prometheus needs a proper configuration file. The configuration file that is going to be used is as follows:

```markup
# prometheus.yml
scrape_configs:
  - job_name: GoServer
    scrape_interval: 5s
    static_configs:
       - targets: ['goapp:1234']
```

We tell Prometheus to connect to a host named `goapp` using port number `1234`. Prometheus pulls data every 5 seconds, according to the value of the `scrape_interval` field. You should put `prometheus.yml` in the `prometheus` directory, which should be in the same directory as the `docker-compose.yml` file that is presented next.

Prometheus as well as Grafana and the Go application are going to run as Docker containers using the next `docker-compose.yml` file:

```markup
version: "3"
services:
  goapp:
    image: goapp
    container_name: goapp
    restart: always
    ports:
      - 1234:1234
    networks:
      - monitoring
```

This is the part that deals with the Go application that collects the metrics. The Docker image name, as well as the internal hostname of the Docker container, is `goapp`. You should define the port number that is going to be open for connections. In this case both the internal and external port numbers are `1234`. The internal one is mapped to the external one. Additionally, you should put all Docker images under the same network, which in this case is called `monitoring` and is defined in a while.

```markup
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    restart: always
    user: "0"
    volumes:
      - ./prometheus/:/etc/prometheus/
```

This is how you pass your own copy of `prometheus.yml` to the Docker image to be used by Prometheus. So, `./prometheus/prometheus.yml` from the local machine can be accessed as `/etc/prometheus/prometheus.yml` from within the Docker image.

```markup
      - ./prometheus_data/:/prometheus/
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
```

This is where you tell Prometheus which configuration file to use.

```markup
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=200h'
      - '--web.enable-lifecycle'
    ports:
      - 9090:9090
    networks:
      - monitoring
```

This is where the definition of the Prometheus part of the scenario ends. The Docker image used is called `prom/prometheus:latest` and the internal name of it is `prometheus`. Prometheus listens to port number `9090`.

```markup
  grafana:
    image: grafana/grafana
    container_name: grafana
    depends_on:
      - prometheus
    restart: always
    user: "0"
    ports:
      - 3000:3000
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=helloThere
```

This is the current password of the `admin` user—you need that for connecting to Grafana.

```markup
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_PANELS_DISABLE_SANITIZE_HTML=true
      - GF_SECURITY_ALLOW_EMBEDDING=true
    networks:
      - monitoring
    volumes:
      - ./grafana_data/:/var/lib/grafana/
```

Last, we present the Grafana part. Grafana listens to port number `3000`.

```markup
volumes:
    grafana_data: {}
    prometheus_data: {}
```

The preceding two lines in combination with the two `volumes` fields allow both Grafana and Prometheus to save their data locally so that data is not lost each time you restart the Docker images.

```markup
networks:
  monitoring:
    driver: bridge
```

Internally, all three containers are known by the value of their `container_name` field. However, externally, you can connect to the open ports from your local machine as `http://localhost:port` or from another machine using `http://hostname:port`—the second way is not very secure and should be blocked by a firewall. Lastly, you need to run `docker-compose up` and you are done! The Go application begins exposing data and Prometheus begins collecting it.

The next figure shows the Prometheus UI (`http://hostname:9090`) displaying a simple plot of `packt_n_goroutines`:

![Graphical user interface, application, table
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_08_01.png)

Figure 8.1: The Prometheus UI

This output, which shows the values of the metrics in a graphical way, is very handy for debugging purposes, but it is far from being truly professional as Prometheus is not a visualization tool. The next subsection shows how to connect Prometheus with Grafana and create impressive plots.

## Visualizing Prometheus metrics in Grafana

There is no point in collecting metrics without doing something with them, and by something, I mean visualizing them. Prometheus and Grafana work very well together so we are going to use Grafana for the visualization part. The single most important task that you should perform in Grafana is connecting it with your Prometheus instance. In Grafana terminology, you should create a Grafana _data source_ that allows Grafana to get data from Prometheus.

The steps for creating a data source with our Prometheus installation are the following:

1.  First go to `http://localhost:3000` to connect to Grafana.
2.  The username of the administrator is `admin` whereas the password is defined in the value of the `GF_SECURITY_ADMIN_PASSWORD` parameter of the `docker-compose.yml` file.
3.  Then select **Add your first data source**. From the list of data sources, select **Prometheus**, which is usually at the top of the list.
4.  Put `http://prometheus:9090` in the URL field and then press the **Save & Test** button. Due to the internal network that exists between the Docker images, the Grafana container knows the Prometheus container by the `prometheus` hostname—this is the value of the `container_name` field. As you already know, you can also connect to Prometheus from your local machine using `http://localhost:9090`. We are done! The name of the data source is `Prometheus`.

After these steps, create a new dashboard from the initial Grafana screen and put a new panel on it. Select **Prometheus** as the data source of the panel, if it is not already selected. Then go to the **Metrics** drop-down menu and select the desired metrics. Click **Save** and you are done. Create as many panels as you want.

The next figure shows Grafana visualizing two metrics from Prometheus as exposed by `prometheus.go`.

![A picture containing text, monitor, electronics, screenshot
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_08_02.png)

Figure 8.2: Visualizing metrics in Grafana

Grafana has many more capabilities than the ones presented here—if you are working with system metrics and want to check the performance of your Go applications, Prometheus and Grafana are good choices.

After learning about HTTP servers, the next section shows how to develop HTTP clients.

Just Imagine

# Developing web clients

This section shows how to develop HTTP clients starting with a simplistic version and continuing with a more advanced one. In this simplistic version, all of the work is done by the `http.Get()` call, which is pretty convenient when you do not want to deal with lots of options and parameters. However, this type of call gives you no flexibility over the process. Notice that `http.Get()` returns an `http.Response` value. All this is illustrated in `simpleClient.go`:

```markup
package main
import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
)
func main() {
    if len(os.Args) != 2 {
        fmt.Printf("Usage: %s URL\n", filepath.Base(os.Args[0]))
        return
    }
```

The `filepath.Base()` function returns the last element of a path. When given `os.Args[0]` as its parameter, it returns the name of the executable binary file.

```markup
    URL := os.Args[1]
    data, err := http.Get(URL)
```

In the previous two statements we get the URL and get its data using `http.Get()`, which returns an `*http.Response` and an `error` variable. The `*http.Response` value contains all the information so you do not need to make any additional calls to `http.Get()`.

```markup
    if err != nil {
        fmt.Println(err)
        return
    }
    _, err = io.Copy(os.Stdout, data.Body)
```

The `io.Copy()` function reads from the `data.Body` reader, which contains the body of the server response, and writes the data to `os.Stdout`. As `os.Stdout` is always open, you do not need to open it for writing. Therefore, all data is written to standard output, which is usually the terminal window.

```markup
    if err != nil {
        fmt.Println(err)
        return
    }
    data.Body.Close()
}
```

Last, we close the `data.Body` reader to make the work of the garbage collection easier.

Working with `simpleClient.go` produces the next kind of output, which in this case is abbreviated:

```markup
$ go run simpleClient.go https://www.golang.org
<!DOCTYPE html>
<html lang="en">
<meta charset="utf-8">
<meta name="description" content="Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.">
...
</script>
```

Although `simpleClient.go` does the job of verifying that the given URL exists and is reachable, it offers no control over the process. The next subsection develops an advanced HTTP client that processes the server response.

## Using http.NewRequest() to improve the client

As the web client of the previous section is relatively simplistic and does not give you any flexibility, in this subsection, you learn how to read a URL without using the `http.Get()` function, and with more options. However, the extra flexibility comes at a cost as you must write more code.

The code of `wwwClient.go`, without the `import` block, is as follows:

```markup
package main
// For the import block go to the book GitHub repository
func main() {
    if len(os.Args) != 2 {
        fmt.Printf("Usage: %s URL\n", filepath.Base(os.Args[0]))
        return
    }
```

Although using `filepath.Base()` is not necessary, it makes your output more professional.

```markup
    URL, err := url.Parse(os.Args[1])
    if err != nil {
        fmt.Println("Error in parsing:", err)
        return
    }
```

The `url.Parse()` function parses a string into a `URL` structure. This means that if the given argument is not a valid URL, `url.Parse()` is going to notice. As usual, check the `error` variable.

```markup
    c := &http.Client{
        Timeout: 15 * time.Second,
    }
    request, err := http.NewRequest(http.MethodGet, URL.String(), nil)
    if err != nil {
        fmt.Println("Get:", err)
        return
    }
```

The `http.NewRequest()` function returns an `http.Request` object given a method, a URL, and an optional body. The `http.MethodGet` parameter defines that we want to retrieve the data using a `GET` HTTP method whereas `URL.String()` returns the string value of an `http.URL` variable.

```markup
    httpData, err := c.Do(request)
    if err != nil {
        fmt.Println("Error in Do():", err)
        return
    }
```

The `http.Do()` function sends an HTTP request (`http.Request`) using an `http.Client` and gets an `http.Response`. So, `http.Do()` does the job of `http.Get()` in a more detailed way.

```markup
    fmt.Println("Status code:", httpData.Status)
```

`httpData.Status` holds the HTTP status code of the response—this is really important because it allows you to understand what really happened with the request.

```markup
    header, _ := httputil.DumpResponse(httpData, false)
    fmt.Print(string(header))
```

The `httputil.DumpResponse()` function is used here to get the response from the server and is mainly used for debugging purposes. The second argument of `httputil.DumpResponse()` is a Boolean value that specifies whether the function is going to include the body or not in its output—in our case it is set to `false`, which excludes the response body from the output and only prints the header. If you want to do the same on the server side, you should use `httputil.DumpRequest()`.

```markup
    contentType := httpData.Header.Get("Content-Type")
    characterSet := strings.SplitAfter(contentType, "charset=")
    if len(characterSet) > 1 {
        fmt.Println("Character Set:", characterSet[1])
    }
```

Here we find out about the character set of the response by searching the value of `Content-Type`.

```markup
    if httpData.ContentLength == -1 {
        fmt.Println("ContentLength is unknown!")
    } else {
        fmt.Println("ContentLength:", httpData.ContentLength)
    }
```

Here, we try to get the content length from the response by reading `httpData.ContentLength`. However, if the value is not set, we print a relevant message.

```markup
    length := 0
    var buffer [1024]byte
    r := httpData.Body
    for {
        n, err := r.Read(buffer[0:])
        if err != nil {
            fmt.Println(err)
                break
        }
        length = length + n
    }
    fmt.Println("Calculated response data length:", length)
}
```

In the last part of the program, we use a technique for discovering the size of the server HTTP response on our own. If we wanted to display the HTML output on our screen, we could have printed the contents of the `r` buffer variable.

Working with `wwwClient.go` and visiting [https://www.golang.org](https://www.golang.org) produces the next output, which is the output of `fmt.Println("Status code:", httpData.Status)`:

```markup
$ go run wwwClient.go https://www.golang.org
Status code: 200 OK
```

Next, we see output of the `fmt.Print(string(header))` statement with the header data of the HTTP server response:

```markup
HTTP/2.0 200 OK
Alt-Svc: h3-29=":443"; ma=2592000,h3-T051=":443"; ma=2592000,h3-Q050=":443"; ma=2592000,h3-Q046=":443"; ma=2592000,h3-Q043=":443"; ma=2592000,quic=":443"; ma=2592000; v="46,43"
Content-Type: text/html; charset=utf-8
Date: Sat, 27 Mar 2021 19:19:25 GMT
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
Vary: Accept-Encoding
Via: 1.1 google
```

The last part of the output is about the character set of the interaction (`utf-8`) and the content length of the response (`9216`) as calculated by the code:

```markup
Character Set: utf-8
ContentLength is unknown!
EOF
Calculated response data length: 9216
```

The next section shows how to create a client for the phone book web service we developed earlier.

## Creating a client for the phone book service

In this subsection, we create a command-line utility that interacts with the phone book web service that was developed earlier in this chapter. This version of the phone book client is going to be created using the `cobra` package, which means that a dedicated GitHub or GitLab repository is required. In this case the repository can be found at [https://github.com/mactsouk/phone-cli](https://github.com/mactsouk/phone-cli). The first thing to do after running `git clone` is associating that repository with the directory that is going to be used for development:

```markup
$ cd ~/go/src/github.com/mactsouk
$ git clone git@github.com:mactsouk/phone-cli.git
$ cd phone-cli
$ ~/go/bin/cobra init --pkg-name github.com/mactsouk/phone-cli
$ go mod init
$ go mod tidy
$ go mod download
```

Next, we have to create the commands for the utility. The structure of the utility is implemented using the next `cobra` commands:

```markup
$ ~/go/bin/cobra add search
$ ~/go/bin/cobra add insert
$ ~/go/bin/cobra add delete
$ ~/go/bin/cobra add status
$ ~/go/bin/cobra add list
```

So, we have a command-line utility with five commands named `search`, `insert`, `delete`, `status`, and `list`. After that, we need to implement the commands and define their local parameters in order to interact with the phone book server.

Now let us see the implementations of the commands, starting from the implementation of the `init()` function of the `root.go` file because this is where the global command-line parameters are defined:

```markup
func init() {
    rootCmd.PersistentFlags().StringP("server", "S", "localhost", "Server")
    rootCmd.PersistentFlags().StringP("port", "P", "1234", "Port number")
    viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
    viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
}
```

So, we define two global parameters named `server` and `port`, which are the hostname and the port number, respectively. Both parameters have an alias and both parameters are handled by `viper`.

Let us now examine the implementation of the `status` command as found in `status.go`:

```markup
SERVER := viper.GetString("server")
PORT := viper.GetString("port")
```

All commands read the values of the `server` and `port` command-line parameters in order to get information about the server, and the `status` command is no exception.

```markup
// Create request
URL := "http://" + SERVER + ":" + PORT + "/status"
```

After that we construct the full URL of the request.

```markup
data, err := http.Get(URL)
if err != nil {
    fmt.Println(err)
    return
}
```

Then, we send a `GET` request to the server using `http.Get()`.

```markup
// Check HTTP Status Code
if data.StatusCode != http.StatusOK {
    fmt.Println("Status code:", data.StatusCode)
    return
}
```

After that we check the HTTP status code of the request to make sure that everything is OK.

```markup
// Read data
responseData, err := io.ReadAll(data.Body)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Print(string(responseData))
```

If everything is OK, we read the entire body of the server response, which is a `byte` slice, and print it on screen as a string. The implementation of `list` is almost identical to the implementation of `status`. The only differences are that the implementation is found in `list.go` and that the full URL is constructed as follows:

```markup
URL := "http://" + SERVER + ":" + PORT + "/list"
```

After that, let us see how the `delete` command is implemented in `delete.go`:

```markup
SERVER := viper.GetString("server")
PORT := viper.GetString("port")
number, _ := cmd.Flags().GetString("tel")
if number == "" {
    fmt.Println("Number is empty!")
    return
}
```

Apart from reading the values of the `server` and `port` global parameters, we read the value of the `tel` parameter. If `tel` has no value, the command returns.

```markup
// Create request
URL := "http://" + SERVER + ":" + PORT + "/delete/" + number
```

Once again, we construct the full URL of the request before connecting to the server.

```markup
// Send request to server
data, err := http.Get(URL)
if err != nil {
    fmt.Println(err)
    return
}
```

Then, we send the request to the server.

```markup
// Check HTTP Status Code
if data.StatusCode != http.StatusOK {
    fmt.Println("Status code:", data.StatusCode)
    return
}
```

If there was an error in the server response, the `delete` command terminates.

```markup
// Read data
responseData, err := io.ReadAll(data.Body)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Print(string(responseData))
```

If everything was fine, the server response text is printed on the screen.

The `init()` function of `delete.go` contains the definition of the local `tel` command-line parameter:

```markup
func init() {
    rootCmd.AddCommand(deleteCmd)
    deleteCmd.Flags().StringP("tel", "t", "", "Telephone number to delete")
}
```

This is a local flag available to the `delete` command only. Next, let us learn more about the `search` command and how it is implemented in `search.go`. The implementation is the same as in `delete` except for the full request URL:

```markup
URL := "http://" + SERVER + ":" + PORT + "/search/" + number
```

The `search` command also supports the `tel` command-line parameter for getting the telephone number to search for—this is defined in the `init()` function of `search.go`.

The last command that is presented is the `insert` command, which supports three local command-line parameters that are defined in the `init()` function in `insert.go`:

```markup
func init() {
    rootCmd.AddCommand(insertCmd)
    insertCmd.Flags().StringP("name", "n", "", "Name value")
    insertCmd.Flags().StringP("surname", "s", "", "Surname value")
    insertCmd.Flags().StringP("tel", "t", "", "Telephone value")
}
```

These three parameters are needed for getting the required user input. Note that the alias for `surname` is a **lowercase** `s` whereas the alias for `server`, which is defined in `root.go`, is an **uppercase** `S`. Commands and their aliases are user-defined—use common sense when selecting command names and aliases.

The command is implemented using the next code:

```markup
SERVER := viper.GetString("server")
PORT := viper.GetString("port")
```

First, we read the `server` and `port` global parameters.

```markup
number, _ := cmd.Flags().GetString("tel")
if number == "" {
    fmt.Println("Number is empty!")
    return
}
name, _ := cmd.Flags().GetString("name")
if number == "" {
    fmt.Println("Name is empty!")
    return
}
surname, _ := cmd.Flags().GetString("surname")
if number == "" {
    fmt.Println("Surname is empty!")
    return
}
```

Then, we get the values of the three local command-line parameters. If any one of them has an empty value, the command returns without sending the request to the server.

```markup
URL := "http://" + SERVER + ":" + PORT + "/insert/"
URL = URL + "/" + name + "/" + surname + "/" + number
```

Here, we create the server request in two steps for readability.

```markup
data, err := http.Get(URL)
if err != nil {
    fmt.Println("**", err)
    return
}
```

Then, we send the request to the server.

```markup
if data.StatusCode != http.StatusOK {
    fmt.Println("Status code:", data.StatusCode)
    return
}
```

Checking the HTTP status code is considered a good practice. Therefore, if everything is OK with the server response, we continue by reading the data. Otherwise, we print the status code, and we exit.

```markup
responseData, err := io.ReadAll(data.Body)
if err != nil {
    fmt.Println("*", err)
    return
}
fmt.Print(string(responseData))
```

After reading the body of the server response, which is stored in a `byte` slice, we print it on screen as a string using `string(responseData)`.

The client application generates the next kind of output:

```markup
$ go run main.go list
Dimitris Tsoukalos 2101112223
Jane Doe 0800123456
Mike Tsoukalos 2109416471
```

This is the output of the `list` command.

```markup
$ go run main.go status
Total entries: 3
```

The output of the `status` command informs us about the number of entries in the phone book.

```markup
$ go run main.go search --tel 0800123456
Jane Doe 0800123456
```

The previous output shows the use of the `search` command when successfully finding a number.

```markup
$ go run main.go search --tel 0800
Status code: 404
```

The previous output shows the use of the `search` command when not finding a number.

```markup
$ go run main.go delete --tel 2101112223
2101112223 deleted!
```

This is the output of the `delete` command.

```markup
$ go run main.go insert -n Michalis -s Tsoukalos -t 2101112223
New record added successfully
```

This is the operation of the `insert` command. If you try to insert the same number more than once, the server output is going to be `Status code: 304`.

The next section shows how to create an FTP server using `net/http`.

Just Imagine

# Creating file servers

Although a file server is not a web server per se, it is closely connected to web services because it is being implemented using similar Go packages. Additionally, file servers are frequently used for supporting the functionality of web servers and web services.

Go offers the `http.FileServer()` handler for doing so, as well as `http.ServeFile()`. The biggest difference between these two is that `http.FileServer()` is an `http.Handler` whereas `http.ServeFile()` is not. Additionally, `http.ServeFile()` is better at serving single files whereas `http.FileServer()` is better at serving entire directory trees.

A simple code example of `http.FileServer()` is presented in `fileServer.go`:

```markup
package main
import (
    "fmt"
    "log"
    "net/http"
)
var PORT = ":8765"
func defaultHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Serving:", r.URL.Path, "from", r.Host)
    w.WriteHeader(http.StatusOK)
    Body := "Thanks for visiting!\n"
    fmt.Fprintf(w, "%s", Body)
}
```

This is the expected default handler of the HTTP server.

```markup
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", defaultHandler)
    fileServer := http.FileServer(http.Dir("/tmp/"))
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))
```

`mux.Handle()` registers the file server as the handler for all URL paths that begin with `/static/`. However, when a match is found, we strip the `/static/` prefix before the file server tries to serve such a request because `/static/` is not part of the location where the actual files are located. As far as Go is concerned, `http.FileServer()` is just another handler.

```markup
    fmt.Println("Starting server on:", PORT)
    err := http.ListenAndServe(PORT, mux)
    fmt.Println(err)
}
```

Last, we start the HTTP server using `http.ListenAndServe()`.

Using `curl(1)` for visiting `/static/` produces the next kind of output, which is in HTML format:

```markup
$ curl http://localhost:8765/static/
<pre>
<a href="AlTest1.out">AlTest1.out</a>
<a href="adobegc.log">adobegc.log</a>
<a href="com.google.Keystone/">com.google.Keystone/</a>
<a href="data.csv">data.csv</a>
<a href="fseventsd-uuid">fseventsd-uuid</a>
<a href="powerlog/">powerlog/</a>
</pre>
```

You can also visit `http://localhost:8765/static/` in your web browser or an FTP client to browse the files and the directories of the FTP server.

The next subsection shows how to utilize `http.ServeFile()` to serve single files.

## Downloading the contents of the phone book application

In this subsection we create and implement an endpoint that allows us to download the contents of a single file. The code creates **a temporary file with a different filename for each request** with the contents of the phone book application. For reasons of simplicity, the presented code supports two HTTP endpoints: one for the default router and the other for the serving of the file. As we are serving a single file, we are going to use `http.ServeFile()`, which replies to a request with the contents of the specified file or directory.

Each temporary file is kept in the file system for _30 seconds_ before it is deleted. To simulate a real-world situation, the presented utility reads `data.csv`, puts it into a slice, and creates a file based on the contents of `data.csv`. The name of the utility is `getEntries.go`—its most important code is the implementation of the `getFileHandler()` function:

```markup
func getFileHandler(w http.ResponseWriter, r *http.Request) {
    var tempFileName string
    // Create temporary file name
    f, err := os.CreateTemp("", "data*.txt")
    tempFileName = f.Name()
```

The temporary path is created using `os.CreateTemp()` based on a given pattern and by adding a random string to the end. If the pattern includes an `*`, then the randomly generated string replaces the last `*`. The exact place where the file is created depends on the operating system being used.

```markup
    // Remove the file
    defer os.Remove(tempFileName)
```

We do not want to end up having lots of temporary files, so we delete the file when the handler function returns.

```markup
    // Save data to it
    err = saveCSVFile(tempFileName)
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintln(w, "Cannot create: "+tempFileName)
        return
    }
    fmt.Println("Serving ", tempFileName)
    http.ServeFile(w, r, tempFileName)
```

This is where the temporary file is sent to the client.

```markup
    // 30 seconds to get the file
    time.Sleep(30 * time.Second)
}
```

The `time.Sleep()` call delays the deletion of the temporary file for 30 seconds—you can define any delay period you like.

As far as the `main()` function is concerned, `getFileHandler()` is a regular handler function used in a `mux.HandleFunc("/getContents/", getFileHandler)` statement. Therefore, each time there is a client request for `/getContents/`, the contents of a file are returned to the HTTP client.

Running `getEntries.go` and visiting `/getContents/` produces the next kind of output:

```markup
$ curl http://localhost:8765/getContents/
Dimitris,Tsoukalos,2101112223,1617028128
Jane,Doe,0800123456,1608559903
Mike,Tsoukalos,2109416471,1617028196
```

As we are returning plain text data, the output is presented on screen.

_Chapter 10_, _Working with REST APIs_, presents a different way of creating a file server that supports both uploading and downloading files by using the `gorilla/mux` package.

The next section explains how to time out HTTP connections.

Just Imagine

# Timing out HTTP connections

This section presents techniques for timing out HTTP connections that take too long to finish and work either on the server or the client side.

## Using SetDeadline()

The `SetDeadline()` function is used by `net` to set the read and write deadlines of network connections. Due to the way the `SetDeadline()` function works, you need to call `SetDeadline()` before any read or write operation. Keep in mind that Go uses deadlines to implement timeouts, so you do not need to reset the timeout every time your application receives or sends any data. The use of `SetDeadline()` is illustrated in `withDeadline.go` and more specifically in the implementation of the `Timeout()` function:

```markup
var timeout = time.Duration(time.Second)
func Timeout(network, host string) (net.Conn, error) {
    conn, err := net.DialTimeout(network, host, timeout)
    if err != nil {
        return nil, err
    }
    conn.SetDeadline(time.Now().Add(timeout))
    return conn, nil
}
```

The `timeout` global variable defines the timeout period used in the `SetDeadline()` call.

The previous function is used in the next code inside `main()`:

```markup
t := http.Transport{
    Dial: Timeout,
}
client := http.Client{
        Transport: &t,
}
```

So, `http.Transport` uses `Timeout()` in the `Dial` field and `http.Client` uses `http.Transport`. When you call the `client.Get()` method with the desired URL, which is not shown here, `Timeout` is automatically being used because of the `http.Transport` definition. So, if the `Timeout` function returns before the server response is received, we have a timeout.

Using `withDeadline.go` produces the next kind of output:

```markup
$ go run withDeadline.go http://www.golang.org
Timeout value: 1s
<!DOCTYPE html>
...
```

The call was successful and took less than 1 second to finish, so there was no timeout.

```markup
$ go run withDeadline.go http://localhost:80
Timeout value: 1s
Get "http://localhost:80": read tcp 127.0.0.1:52492->127.0.0.1:80: i/o timeout
```

This time we have a timeout as the server took too long to answer.

Next, we show how to time out a connection using the `context` package.

## Setting the timeout period on the client side

This section presents a technique for timing out network connections that take too long to finish on the **client side**. So, if the client does not receive a response from the server in the desired time, it closes the connection. The `timeoutClient.go` source file, without the `import` block, illustrates the technique.

```markup
package main
// For the import block go to the book code repository
var myUrl string
var delay int = 5
var wg sync.WaitGroup
type myData struct {
    r   *http.Response
    err error
}
```

In the previous code we define global variables and a structure that are going to be used in the rest of the program.

```markup
func connect(c context.Context) error {
    defer wg.Done()
    data := make(chan myData, 1)
    tr := &http.Transport{}
    httpClient := &http.Client{Transport: tr}
    req, _ := http.NewRequest("GET", myUrl, nil)
```

This is where you initialize the variables of the HTTP connection. The `data` channel is used in the `select` statement that follows. Additionally, the `c context.Context` parameter comes with an embedded channel that is also used in the `select` statement.

```markup
    go func() {
        response, err := httpClient.Do(req)
        if err != nil {
            fmt.Println(err)
            data <- myData{nil, err}
            return
        } else {
            pack := myData{response, err}
            data <- pack
        }
    }()
```

The previous goroutine is used for interacting with the HTTP server. There is nothing special here as this is a regular interaction of an HTTP client with an HTTP server.

```markup
    select {
    case <-c.Done():
        tr.CancelRequest(req)
        <-data
        fmt.Println("The request was canceled!")
        return c.Err()
```

The code that this `select` block executes is based on whether the context is going to time out or not. If the context times out first, then the client connection is canceled using `tr.CancelRequest(req)`.

```markup
    case ok := <-data:
        err := ok.err
        resp := ok.r
        if err != nil {
            fmt.Println("Error select:", err)
            return err
        }
        defer resp.Body.Close()
        realHTTPData, err := io.ReadAll(resp.Body)
        if err != nil {
            fmt.Println("Error select:", err)
            return err
        }
        fmt.Printf("Server Response: %s\n", realHTTPData)
    }
    return nil
}
```

The second `select` branch deals with the data received from the HTTP server, which is handled in the usual way.

```markup
func main() {
    if len(os.Args) == 1 {
        fmt.Println("Need a URL and a delay!")
        return
    }
    myUrl = os.Args[1]
    if len(os.Args) == 3 {
        t, err := strconv.Atoi(os.Args[2])
        if err != nil {
            fmt.Println(err)
            return
        }
        delay = t
    }
```

The URL is read directly because it is already a string value whereas the delay period is converted into a numeric value using `strconv.Atoi()`.

```markup
    fmt.Println("Delay:", delay)
    c := context.Background()
    c, cancel := context.WithTimeout(c, time.Duration(delay)*time.Second)
    defer cancel()
```

The timeout period is defined by the `context.WithTimeout()` method. It is considered a good practice to use `context.Background()` in the `main()` function or the `init()` function of a package or in tests.

```markup
    fmt.Printf("Connecting to %s \n", myUrl)
    wg.Add(1)
    go connect(c)
    wg.Wait()
    fmt.Println("Exiting...")
}
```

The `connect()` function, which is also executed as a goroutine, either terminates normally or when the `cancel()` function is executed—the `cancel()` function is what calls the `Done()` method of `c`.

Working with `timeoutClient.go` and having a timeout situation generates the following kind of output:

```markup
$ go run timeoutClient.go http://localhost:80
Delay: 5
Connecting to http://localhost:80
Get "http://localhost:80": net/http: request canceled
The request was canceled!
Exiting...
```

The next subsection shows how to time out an HTTP request on the server side.

## Setting the timeout period on the server side

This section presents a technique for timing out network connections that take too long to finish on the server side. This is much more important than the client side because a server with too many open connections might not be able to process more requests unless some of the already open connections close. This usually happens for two reasons. The first reason is software bugs, and the second reason is when a server is experiencing a **Denial of Service** (**DoS**) attack!

The `main()` function in `timeoutServer.go` shows the technique:

```markup
func main() {
    PORT := ":8001"
    arguments := os.Args
    if len(arguments) != 1 {
        PORT = ":" + arguments[1]
    }
    fmt.Println("Using port number: ", PORT)
    m := http.NewServeMux()
    srv := &http.Server{
        Addr:         PORT,
        Handler:      m,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    }
```

This is where the timeout periods are defined. Note that you can define timeout periods for both reading and writing processes. The value of the `ReadTimeout` field specifies the maximum duration allowed to read the entire client request, including the body, whereas the value of the `WriteTimeout` field specifies the maximum time duration before timing out the sending of the client response.

```markup
    m.HandleFunc("/time", timeHandler)
    m.HandleFunc("/", myHandler)
    err := srv.ListenAndServe()
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

Apart from the parameters in the definition of `http.Server`, the rest of the code is as usual: it contains the handler functions and calls `ListenAndServe()` for starting the HTTP server.

Working with `timeoutServer.go` generates no output. However, if a client connects to it without sending any requests, the client connection is going to end after 3 seconds. The same is going to happen if it takes the client more than 3 seconds to receive the server response.

Just Imagine

# Exercises

-   Put all handlers from `www-phone.go` in a different Go package and modify `www-phone.go` accordingly. You need a different repository for storing the new package.
-   Modify `wwwClient.go` to save the HTML output to an external file.
-   Include the functionality of `getEntries.go` in the phone book application.
-   Implement a simple version of `ab(1)` using goroutines and channels. `ab(1)` is an Apache HTTP server benchmarking tool.

Just Imagine

# Summary

In this chapter, we learned how to work with HTTP, how to create Docker images from Go code, how to expose metrics to Prometheus, as well as how to develop HTTP clients and servers. We have also updated the phone book application into a web application and programmed a command-line client for it. Additionally, we learned how to time out HTTP connections and develop file servers.

We are now ready to begin developing powerful and concurrent HTTP applications—however, we are not done yet with HTTP. _Chapter 10_, _Working with REST APIs_, is going to connect the dots and show how to develop powerful RESTful servers and clients.

But first, we need to learn about working with TCP/IP, TCP, UDP, and WebSocket, which are the subjects of the next chapter.

Just Imagine

# Additional resources

-   Caddy server: [https://caddyserver.com/](https://caddyserver.com/)
-   Nginx server: [https://nginx.org/en/](https://nginx.org/en/)
-   Histograms in Prometheus: [https://prometheus.io/docs/practices/histograms/](https://prometheus.io/docs/practices/histograms/)
-   The `net/http` package: [https://golang.org/pkg/net/http/](https://golang.org/pkg/net/http/)
-   Official Docker Go images: [https://hub.docker.com/\_/golang/](https://hub.docker.com/_/golang/)