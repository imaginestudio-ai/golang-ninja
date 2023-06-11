# Calling External Processes and Handling Errors and Timeouts

Many command-line applications interact with other external commands or API services. This chapter will guide you through how to call these external processes and how to handle timeouts and other errors when they occur. The chapter will start with a deep dive into the `os/exec` package, which contains everything you need to create commands that call external processes that give you multiple options for creating and running commands. You’ll learn how to retrieve data from the standard output and standard error pipes, as well as creating additional file descriptors for similar usage.

Another external process involves calling external API service endpoints. The `net/http` package is discussed and is where we start defining the client, then create the requests that it executes. We will discuss the different ways requests can be both created and executed.

Timeouts and other errors can occur when calling either type of process. We will end the chapter by looking at how to capture when timeouts and errors occur in our code. It’s important to be mindful that these things can happen and so it’s important to write code that can handle them. The specific action taken upon error is dependent on the use case, so we’ll discuss the code to capture these cases only. To summarize, we’ll be covering the following topics:

-   Calling external processes
-   Interacting with REST APIs
-   Handling the expected – timeouts and errors

# Calling external processes

Within your command-line application, you may need to call some external processes. Sometimes, there are Golang libraries offered for third-party tools that function as a wrapper. For example, Go CV, [https://gocv.io/](https://gocv.io/), is a Golang wrapper offered for OpenCV, an open source computer vision library. Then, there’s GoFFmpeg, [https://github.com/xfrr/goffmpeg](https://github.com/xfrr/goffmpeg), which is a wrapper offered for FFmpeg, a library for recording, converting, and streaming audio and video files. Often, you need to install an underlying tool, such as OpenCV or FFmpeg, and then the library interacts with it. Calling these external processes then means importing the wrapper package and calling its methods within your code. Often, when you dive into the code, you’ll find that these libraries provide a wrapper for the C code.

Besides importing a wrapper for an external tool, you may call external applications using the `os/exec` Golang library. This is the main purpose of the library and in this section, we will be digging into how to use it to call external applications.

First, let’s review each of the variables, types, and functions that exist within the `os/exec` package with an example of each.

## The os/exec package

By digging deeper into the `exec` package, you will find that it is a wrapper for the `os.StartProcess` method, making it easier to handle the remapping of standard in and standard out, connecting the input and output with pipes, and handling other modifications.

For clarity, it’s important to note that this package does not invoke the operating system’s shell and so doesn’t handle tasks handled typically by the shell: expanding glob patterns, pipelines, or redirections. If it is necessary to expand glob patterns, then you can call the shell directly and make sure to escape values to make it safe, or you can also use the path or file path’s `Glob` function. To expand any environment variables that exist in a string, use the `os` package’s `ExpandEnv` function.

In the following subsections, we’ll start to discuss the different variables, types, functions, and methods that exist within the `os/exec` package.

### Variables

`ErrNotFound` is the error variable returned when an executable file is not found in the application’s `$``PATH` variables.

### Types

`Cmd` is a struct that represents an external command. Defining a variable of this type is just in preparation for the command to be run. Once this variable, of the `Cmd` type, is run via either the `Run`, `Output`, or `CombinedOutput` method, it cannot be reused. There are several fields on this `Cmd` struct that we can also elaborate upon:

-   `Path string` This is the only required field. It is the path of the command to run; if the path is relative, then it will be relative to the value stored in the `Dir` field.
-   `Args []string` This field holds the arguments for the command. `Args[0]` represents the command. `Path` and `Args` are set when the command is run, but if `Args` is `nil` or empty, then just `{Path}` is used during execution.
-   `Env []string` The `Env` field represents the environment for the command to run. Each value in the slice must be in the following format: `"key=value"`. If the value is empty or `nil`, then the command uses the current environment. If the slice has duplicate key values, then the last value for the duplicate key is used.
-   `Dir string` The `Dir` field represents the working directory of the command. If it’s not set, then the current directory is used.
-   `Stdin io.Reader` The `Stdin` field specifies the command process’ standard input. If the data is `nil`, then the process reads from `os.DevNull`, the null device. However, if the standard input is `*os.File`, then the contents are piped. During execution, a goroutine reads from standard input and then sends that data to the command. The `Wait` method will not complete until the goroutine starts copying. If it does not complete, then it could be because of an **end-of-file** (**EOF**), read, or write-to-pipe error.
-   `Stdout io.Writer` The `Stdout` field specifies the command process’ standard output. If the standard output is `nil`, then the process connects to the `os.DevNull` null device. If the standard output is `*os.File`, then output is sent to it instead. During execution, a goroutine reads from the command process and sends data to the writer.
-   `Stderr io.Writer` The `Stderr` field specifies the command process’ standard error output. If the standard error is `nil`, then the process connects to the `os.DevNull` null device. If the standard error is `*os.File`, then error output is sent to it instead. During execution, a goroutine reads from the command process and sends data to the writer.
-   `ExtraFiles []*os.File` The `ExtraFiles` field specifies additional files inherited by the command process. It doesn’t include standard input, standard output, or standard error, so if not empty, entry _x_ becomes the _3+x_ file descriptor. This field is not supported on Windows.
-   `SysProcAttr *syscall.SysProcAttr` `SysProcAttr` holds system-specific attributes that are passed down to `os.StartProcess` as an `os.ProcAttr`’s `Sys` field.
-   `Process *os.Process` The `Process` field holds the underlying process once the command is run.
-   `ProcessState *os.ProcessState` The `ProcessState` field contains information about the process. It becomes available after the wait or run method is called.

### Methods

The following are the methods that exist on the `exec.Cmd` object:

-   `func (c *Cmd) CombinedOutput() ([]byte, error)` The `CombinedOutput` method returns both the standard output and standard error into 1-byte string output.
-   `func (c *Cmd) Output ([]byte, error)` The `Output` method returns just the standard output. If an error occurs, it will usually be of the `*ExitError` type, and if the command’s standard error, `c.Stderr`, is `nil`, `Output` populates `ExitError.Stderr`.
-   `func (c *Cmd) Run() error` The `Run` method starts executing the command and then waits for it to complete. If there was no problem copying standard input, standard output, or standard error and the command exits with a zero status, then the error returned will be `nil`. If the command exits with an error, it will usually be of the `*ExitError` type, but could be other error types as well.
-   `func (c *Cmd)` `Start() error`
-   The `Start` method will start executing the command and not wait for it to complete. If the `Start` method runs successfully, then the `c.Process` field will be set. The `c.Wait` field will then return the exit code and release resources once complete.
-   f`unc (c* Cmd) StderrPipe() (io.ReadCloser, error)` `StderrPipe` returns a pipe that is connected to the command’s standard error. There won’t be a need to ever close the pipe because the `Wait` method will close the pipe once the command exits. Do not call the `Wait` method until all reads from the standard error pipe have completed. Do not use this command with the `Run` method for the same reason.
-   `func (c* Cmd) StdinPipe() (io.WriteCloser, error`) `StdinPipe` returns a pipe that is connected to the command’s standard input. The pipe will be closed after `Wait`, and the command exits. However, sometimes the command will not run until the standard input pipe is closed, and thus you can call the `Close` method to close the pipe sooner.
-   `func (c *Cmd) StdoutPipe() (io.ReadCloser, error`) The `StdoutPipe` method returns a pipe that is connected to the command’s standard output. There’s no need to close the pipe because `Wait` will close the pipe once the command exits. Again, do not call `Wait` until all reads from the standard output pipe have completed. Do not use this command with the `Run` method for the same reason.
-   `func (c *Cmd) String() string` The `String` method returns a human-readable description of the command, `c`, for debugging purposes. The specific output may differ between Go version releases. Also, do not use this as input to a shell, as it’s not suitable for that purpose.
-   `func (c *Cmd) Wait() error` The `Wait` method waits for any copying to standard input, for standard output or standard error to complete, and for the command to exit. To utilize the `Wait` method, the command must have been started by the `Start` method and not the `Run` method. If there are no errors with copying from pipes and the process exits with a `0` exit status code, then the error returned will be `nil`. If the command’s `Stdin`, `Stdout`, or `Stderr` field is not set to `*os.File`, then `Wait` also ensures that the respective input-output loop process completes as well.

`Error` is a struct that represents an error returned from the `LookPath` function when it fails to recognize the file as an executable. There are a couple of fields and methods of this specific error type that we will define in detail.

The following are the methods that exist on the `Error` type:

-   `func (e *Error) Unwrap() error` If the error returned is a chain of errors, then you can utilize the `Unwrap` method to _unwrap_ it and determine what kind of error it is.

`ExitError` is a struct that represents an error when a command exits unsuccessfully. `*os.ProcessState` is embedded into this struct, so all values and fields will also be available to the `ExitError` type. Finally, there are a few fields of this type that we can define in more detail:

-   `Stderr []byte` This field holds a set of the standard error output responses if not collected from the `Cmd.Output` method. `Stderr` may only contain the prefix and suffix of the error output if it’s sufficiently long. The middle will contain text about the number of omitted bytes. For debugging purposes, and if you want to include the entirety of the error messages, then redirect to `Cmd.Stderr`.

The following is the method that exists on the `ExitError` type:

-   `func (e *ExitError) Error() string` The `Error` method returns the exit error represented as a string.

### Functions

The following are functions that exist within the `os/exec` package:

-   `func LookPath(file string) (string, error)` The `LookPath` function checks to see whether the file is an executable and can be found. If the file is a relative path, then it is relative to the current directory.
-   `func Command(name string, arg ...string) *Cmd` The `Command` function returns the `Cmd` struct with just the path and args set. If name has path separators, then the `LookPath` function is used to confirm the file is found and executable. Otherwise, `name` is used directly as the path. This function behaves slightly differently on Windows. For example, it will execute the whole command line as a single string, including quoted args, then handle its own parsing.
-   `func CommandContext(ctx context.Context, name string, arg ...string) *Cmd` Similar to the `Command` function, but receives context. If the context is executed before the command completes, then it will kill the process by calling `os.Process.Kill`.

Now that we’ve really dived deep into the `os/exec` package and the structs, functions, and methods needed to execute functions, let’s actually use them in code to execute a function externally. Let’s create commands using the `Cmd` struct, but also with the `Command` and `CommandContext` functions. We can then take one example command and run it using either the `Run`, `Output`, or `CombinedOutput` method. Finally, we will handle some errors typically returned from these methods.

Note

If you want to follow along with the examples coming up, within the `Chapter-6` repository, install the necessary applications. In Windows, use the `.\build-windows.p1` PowerShell script. In Darwin, use the `make install` command. Once the applications are installed, run `go` `run main.go`.

## Creating commands using the Cmd struct

There are several different ways of creating commands. The first way is with the `Cmd` struct within the `exec` package.

### Using the Cmd struct

We first define the `cmd` variable with an unset `Cmd` structure. The following code resides in `/examples/command.go` within the `CreateCommandUsingStruct` function:

```markup
cmd := exec.Cmd{}
```

Each field is set separately. The path is set using `filepath.Join`, which is safe for use across different operating systems:

```markup
cmd.Path = filepath.Join(os.Getenv("GOPATH"), "bin", "uppercase")
```

Each field is set separately. The `Args` field contains the command name in the `Args[0]` position, followed by the rest of the arguments to be passed in:

```markup
cmd.Args = []string{"uppercase", "hack the planet"}
```

The following three file descriptors are set – `Stdin`, `Stdout`, and `Stderr`:

```markup
cmd.Stdin = os.Stdin // io.Reader
cmd.Stdout = os.Stdout // io.Writer
cmd.Stderr = os.Stderr // io.Writer
```

However, there’s a `writer`, file descriptor that’s passed into the `ExtraFiles` field. This specific field is inherited by the command process. It’s important to note that a pipe won’t work if you don’t pass the writer in `ExtraFiles`, because the child must get the writer to be able to write to it:

```markup
reader, writer, err := os.Pipe()
if err != nil {
    panic(err)
}
cmd.ExtraFiles = []*os.File{writer}
if err := cmd.Start(); err != nil {
    panic(err)
}
```

Within the actual uppercase command that’s called, there’s code in `cmd/uppercase/uppercase.go` that takes the first argument after the command name and changes the case to uppercase. The new uppercased text is then encoded into the pipe or extra file descriptor:

```markup
input := os.Args[1:]
output := strings.ToUpper(strings.Join(input, ""))
pipe := os.NewFile(uintptr(3), "pipe")
err := json.NewEncoder(pipe).Encode(output)
if err != nil {
    panic(err)
}
```

Back to the `CreateCommandUsingStruct` function, the value that’s encoded into the pipe can now be read via the `read` file descriptor of the pipe and then output with the following code:

```markup
var data string
decoder := json.NewDecoder(reader)
if err := decoder.Decode(&data); err != nil {
    panic(err)
}
fmt.Println(data)
```

We now know one way of creating a command using the `Cmd` struct. Everything could have been defined at once at the same time as the command was initialized and depends on your preference.

### Using the Command function

Another way to create a command is with the `exec.Command` function. The following code resides in `/examples/command.go` within `CreateCommandUsingCommandFunction`:

```markup
cmd := exec.Command(filepath.Join(os.Getenv("GOPATH"), "bin", "uppercase"), "hello world")
reader, writer, err := os.Pipe()
if err != nil {
    panic(err)
}
```

The `exec.Command` function takes the file path to the command as the first argument. A slice of strings representing the arguments is optionally passed for the remaining parameters. The rest of the function is the same. Because `exec.Command` does not take any additional parameters, we similarly define the `ExtraFiles` field outside the original variable initialization.

## Running the command

Now that we know how to create commands, there are multiple different ways to run or start running a command. While each of these methods has already been described in detail earlier in this section, we’ll now share an example of using each.

### Using the Run method

The `Run` method, as mentioned earlier, starts the command process, and then waits for its completion. The code for this is called from the `main.go` file but can be found under `/examples/running.go`. In this example, we call a different command called `lettercount`, which counts the letters in a string and then prints out the result:

```markup
cmd := exec.Command(filepath.Join(os.Getenv("GOPATH"), "bin", "lettercount"), "four")
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
var count int
```

Again, we use the `ExtraFiles` field to pass in an additional file descriptor to write the result to:

```markup
reader, writer, err := os.Pipe()
if err != nil {
    panic(err)
}
cmd.ExtraFiles = []*os.File{writer}
if err := cmd.Run(); err != nil {
    panic(err)
}
if err := json.NewDecoder(reader).Decode(&count); err != nil {
    panic(err)
}
```

The result is finally printed with the following code:

```markup
fmt.Println("letter count: ", count)
```

### Using the Start command

The `Start` method is like the `Run` method; however, it doesn’t wait for the process to complete. You can find the code that uses the `Start` command in `examples/running.go`. For the most part, it’s identical, but you’ll be replacing the code block containing `cmd.Run` with the following:

```markup
if err := cmd.Start(); err != nil {
    panic(err)
}
err = cmd.Wait()
if err != nil {
    panic(err)
}
```

It’s very important to call the `cmd.Wait` method because it releases resources taken by the command process.

### Using the Output command

As the method name suggests, the `Output` method returns anything that’s been piped into the standard out pipe. The most common way to push from a command to the standard output pipe is through any of the print methods in the `fmt` package. An additional line is added to the end of the `main` function for the `lettercount` command:

```markup
fmt.Printf("successfully counted the letters of \"%v\" as %d\n", input, len(runes))
```

The only difference within the code that utilizes this `Output` method, which can be found in the `examples/running.go` file under the `OutputMethod` function, is this line of code:

```markup
out, err := cmd.Output()
```

The `out` variable is a byte slice that can later be cast to a string to be printed out. This variable captures the standard out and when the function is run, the output displayed is as follows:

```markup
output: successfully counted the letters of "four" as 4
```

### Using the CombinedOutput command

As the method name suggests, the `CombinedOutput` method returns a combined output of the standard output and standard error piped data. Add a line toward the end of the `lettercount` command’s `main` function:

```markup
fmt.Fprintln(os.Stderr, "this is where the errors go")
```

The only big difference between the calls from the previous function and the current function, `CombinedOutputMethod`, is this line:

```markup
CombinedOutput, err := cmd.CombinedOutput()
```

Similarly, it returns a byte slice, but now contains the combined output of standard error and standard output.

### Executing commands on Windows

Alongside the examples are similar files that end with `_windows.go`. The major thing to note, in the previous examples, is that `ExtraFiles` is not supported on Windows. These Windows-specific and simple examples execute an external `ping` command to `google.com`. Let’s take a look at one:

```markup
func CreateCommandUsingCommandFunction() {
    cmd := exec.Command("cmd", "/C", "ping", "google.com")
    output, err := cmd.CombinedOutput()
    if err != nil {
        panic(err)
    }
    fmt.Println(string(output))
}
```

Like the commands we’ve written for Darwin, we can create commands using the `exec.Command` function or the struct and call `Run`, `Start`, `Wait`, `Output`, and `CombinedOutput` just the same.

Also, for pagination, `less` is used on Linux and UNIX machines, but `more` is used on Windows. Let’s quickly show this code:

```markup
func Pagination() {
    moreCmd := exec.Command("cmd", "/C", "more")
    moreCmd.Stdin = strings.NewReader(blob)
    moreCmd.Stdout = os.Stdout
    moreCmd.Stderr = os.Stderr
    err := moreCmd.Run()
    if err != nil {
        panic(err)
    }
}
var (
    blob = `
    …
    `
)
```

Similarly, we can pass in the name and all arguments using the `exec.Command` method. We also pass the long text into the `moreCmd.Stdin` field.

So, the `os/exec` package offers different ways to create and run external commands. Whether you create a quick command using the `exec.Command` method or directly create one with the `exec.Cmd` struct and then run the `Start` command, you have options. Finally, you can either retrieve the standard output and error output separately or together. Knowing all about the `os/exec` package will make it easy to successfully run external commands from your Go command-line application.

Just Imagine

# Interacting with REST APIs

Often, if a company or user has already created an API, the command-line application will send requests to either the REST API or the gRPC endpoints. Let’s first talk about using REST API endpoints. It is important to understand the `net/http` package. It’s quite a large package with many types, methods, and functions, many of which are used for development on the server side. In this context, the command-line application will be the client of the API, so we won’t discuss each in detail. We’ll go into a few basic use cases from the client side though.

## Get request

Let’s revisit the code from [_Chapter 3_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_03.xhtml#_idTextAnchor061), _Building an Audio Metadata CLI_. Within the `Run` command of the CLI command code, found in the `/cmd/cli/command/get.go` file, is a snippet of code that calls the corresponding API request endpoint using the `GET` method:

```markup
params := "id=" + url.QueryEscape(cmd.id)
path := fmt.Sprintf("http://localhost/request?%s", params)
payload := &bytes.Buffer{}
method := "GET"
client := cmd.client
```

Notice that in the preceding code, we take the field value, `id`, which has been set on the `cmd` variable, and pass it into the HTTP request as a parameter. Consider the flags and arguments to be passed which are to be used as parameters for your HTTP request. The following code executes the request:

```markup
req, err := http.NewRequest(method, path, payload)
if err != nil {
    return err
}
resp, err := client.Do(req)
if err != nil {
    return err
}
defer resp.Body.Close()
```

Finally, the response is read into a byte string and printed. Prior to accessing the body of the response, check whether the response or body is `nil`. This can save you from some future headaches:

```markup
b, err := io.ReadAll(resp.Body)
if err != nil {
    return err
}
fmt.Println(string(b))
return nil
```

However, in reality, there will be much more done with the response:

1.  **Check the response status code**: If the response is `200` `OK`, then we can return the output as it was a successful response. Otherwise, in the next section, _Handling the expected – timeouts and errors_, we’ll discuss how to handle other responses.
2.  **Log the response**: We may, ideally, log the response if it doesn’t contain any sensitive data. This detailed information can be written to a log file or output when in verbose mode.
3.  **Store the response**: Sometimes, the response may be stored in a local database or cache.
4.  **Transform the data**: This returned data may also be unmarshaled into a local data struct. The struct types of data returned must be defined and, preferably, would utilize the same struct models defined within the API. In that case, if `Content-Type` in the header is set to `application/json`, we would unmarshal the JSON response into the struct.

Currently, in the audiofile application, we transform the data into an `Audio` struct like this:

```markup
var audio Audio
If err := json.Unmarshal(b, &audio); err != nil {
    fmt.Println("error unmarshalling JSON response"
}
```

But what if the response body isn’t in JSON format and the content type is something else? In a perfect world, we’d have API documentation that informs us of what to expect so we can handle it accordingly. Alternatively, you can check to confirm the type first using the following:

```markup
contentType := http.DetectContentType(b) // b are the bytes from reading the resp.Body
```

A quick search on the internet for HTTP content types will return a large list. In the preceding example, the audio company might have decided to return a `Content-Type` value of `audio/wave`. In that case, we could either download or stream the result. There are also different HTTP method types defined as constants within the `net/http` package:

-   `MethodGet`: Used when requesting data
-   `MethodPost`: Used for inserting data
-   `MethodPut`: Request is idempotent, used for inserting or updating an entire resource
-   `MethodPatch`: Similar to `MethodPut`, but sends only partial data to update without modifying the entire resource
-   `MethodDelete`: Used for deleting or removing data
-   `MethodConnect`: Used when talking to a proxy, when the URI begins with `https://`
-   `MethodOptions`: Used to describe the communication options, or allowed methods, with the target
-   `MethodTrace`: Used for debugging by providing a message loop-back along the path of the target

There are many possibilities for the method types and content types of data returned. In the preceding `Get` example, we use a client’s `Do` method to call the method. Another option is to use the `http.Get` method. If we use that method, then we would use this code instead to execute the request:

```markup
resp, err := http.Get(path)
if err != nil {
    return err
}
defer resp.Body.Close()
```

Similarly, rather than using the `client.Do` method for a post or to post a form, there are specific `http.Post` and `http.PostForm` methods that can be used instead. There are times when one method works better for what you are doing. At this point, it’s just important to understand your options.

## Pagination

Suppose there is a large amount of data being returned by the request. Rather than overloading the client by receiving the data all at once, often pagination is an option. There are two fields that can be passed in as parameters to the call:

-   `Limit`: The number of objects to be returned
-   `Page`: The cursor for multiple pages of results returned

We can define these internally and then formulate the path as follows:

```markup
path := fmt.Sprintf("http://localhost/request?limit=%d&page=%d", limit, page)
```

Make sure, if you’re using an external API, to construct their documentation with the proper parameters for pagination and usage. This is just a general example. In fact, there are several other ways of doing pagination. You can send additional requests in a loop, incrementing the page until all data is retrieved.

From the command side, however, you could return all the data after pagination, but you can also handle pagination on the CLI side. A way to handle it on the client side after a large amount of data is collected from an HTTP `Get` request is to pipe the data. This data can be piped into the operating system’s pager command. For UNIX, `less` is the pager command. We create the command and then pipe the string output to the `Stdin` pipe. This code can be found in the `examples/pagination.go` file. Similar to the other examples we’ve shared when creating a command, we create a pipe and pass in the writer as an extra file descriptor to the command so that data may be written out:

```markup
pagesCmd := exec.Command(filepath.Join(os.Getenv("GOPATH"), "bin", "pages"))
reader, writer, err := os.Pipe()
if err != nil {
    panic(err)
}
pagesCmd.Stdin = os.Stdin
pagesCmd.Stdout = os.Stdout
pagesCmd.Stderr = os.Stderr
pagesCmd.ExtraFiles = []*os.File{writer}
if err := pagesCmd.Run(); err != nil {
    panic(err)
}
```

Again, the data from the reader is decoded into the `data` `string` variable:

```markup
var data string
decoder := json.NewDecoder(reader)
if err := decoder.Decode(&data); err != nil {
    panic(err)
}
```

This string is then passed into the `Strings.NewReader` method and defined as the input for the `less` UNIX command:

```markup
lessCmd := exec.Command("/usr/bin/less")
lessCmd.Stdin = strings.NewReader(data)
lessCmd.Stdout = os.Stdout
err = lessCmd.Run()
if err != nil {
    panic(err)
}
```

When the command is run, the data is output as pages. The user then can press the spacebar to continue to the next page or use any of the command keys to navigate the data output.

## Rate limiting

Often, when dealing with third-party APIs, there’s a limit to how many requests can be handled within a particular time. This is commonly known as **rate limiting**. For a single command, you might require multiple requests to an HTTP endpoint and so you might prefer to limit how often you’re sending these requests. Most public APIs will inform users of their rate limits, but there are times when you’ll hit the rate limit of an API unexpectedly. We’ll discuss how to limit your requests to stay within the limits.

There is a useful library, `x/time/rate`, that can be used to define the limit, which is how often something should be executed, and limiters that control the process from executing within the limit. Let’s use some example code, supposing we want to execute something every five seconds.

The code for this particular example is located in the `examples/limiting.go` file. To reiterate, this is just an example and there are different ways to use `runner`. We’re going to cover just a basic use case. We start by defining a struct that contains a function, `Run`, and the `limiter` field, which controls how often it will run. The `Limit()` function will use the `runner` struct to call a function within a rate limit:

```markup
type runner struct {
    Run func() bool
    limiter *rate.Limiter
}
func Limit() {
    thing := runner{}
    start := time.Now()
```

After defining `thing` as a `runner` instance, we get the start time and then define the function of `thing`. If the call is allowed within the time, because it does not exceed the limit, we print the current timestamp and return a `false` variable. We exit the function when at least 30 seconds have passed:

```markup
    thing.Run = func() bool {
        if thing.limiter.Allow() {
            fmt.Println(time.Now()) // or call request
            return false
        }
        if time.Since(start) > 30*time.Second {
            return true
        }
        return false
    }
```

We define the limiter for `thing`. We’ve used a customer variable, which we’ll look at in more detail shortly. Simply, the `NewLimiter` method takes two variables. The first parameter is the limit, one event every five seconds, and the second parameter allows bursts for, at most, a single token:

```markup
    thing.limiter = rate.NewLimiter(forEvery(1, 5*time.
    Second),     1)
```

For those not familiar with the difference between a limit and a burst, a burst defines the number of concurrent requests the API can handle. The rate limit is the number of requests allowed per the defined time.

Next, inside a `for` loop, we call the `Run` function and only break when it returns `true`, which should be after 30 seconds have passed:

```markup
    for {
        if thing.Run() {
            break
        }
    }
}
```

As mentioned, the `forEvery` function, which returns a rate limit, is passed into the `NewLimiter` method. It simply calls the `rate.Every` method, which takes the minimum time interval between events and converts it into a limit:

```markup
func forEvery(eventCount int, duration time.Duration) rate.Limit {
    return rate.Every(duration / time.Duration(eventCount))
}
```

We run this code and the timestamps are output. Notice that they are output every five seconds:

```markup
2022-09-11 18:45:44.356917 -0700 PDT m=+0.000891459
2022-09-11 18:45:49.356877 -0700 PDT m=+5.000891042
2022-09-11 18:45:54.356837 -0700 PDT m=+10.000891084
2022-09-11 18:45:59.356797 -0700 PDT m=+15.000891084
2022-09-11 18:46:04.356757 -0700 PDT m=+20.000891167
2022-09-11 18:46:09.356718 -0700 PDT m=+25.000891167
```

There are other ways of handling limiting requests, such as using a `time.Sleep(d Duration)` method after the code that is called inside a loop. I suggest using the `rate` package because it is great for not only limiting executions but also handling bursts. It has a lot more functionality that can be used for more complex situations when you are sending requests to an external API.

You’ve now learned how to send requests to external APIs and how to handle the response, and when you receive a successful response, how to transform and paginate the results. Also, because rate limiting is commonly required for APIs, we’ve discussed how to do that. Since this section has only handled the case of success, let’s consider how to handle the case of failure in the following section.

Just Imagine

# Handling the expected – timeouts and errors

When building a CLI that calls external commands or sends HTTP requests to an external API, with data that is passed in by the user, it’s a good idea to expect the unexpected. In a perfect world, you can guard against bad data. I’m sure you are familiar with the phrase _garbage in_, _garbage out._ You can create tests that also ensure that your code is covered for as many bad cases as you can think of. However, timeouts and errors happen. It’s the nature of software, and as you come across them within your development and also in production, you can modify your code to handle new cases.

## Timeouts with external command processes

Let’s first discuss how to handle timeouts when calling external commands. The timeout code exists within the `examples/timeout.go` file. The following is the entire method, which calls the `timeout` command. If you take a look at the `timeout` command code, located within `cmd/timeout/timeout.go`, you’ll see that it contains a basic infinite loop. This command will time out, but we need to handle the timeout with the following code:

```markup
func Timeout() {
    errChan := make(chan error, 1)
    cmd := exec.Command(filepath.Join(os.Getenv("GOPATH"), 
           "bin", "timeout"))
    if err := cmd.Start(); err != nil {
        panic(err)
    }
    go func() {
        errChan <- cmd.Wait()
    }()
    select {
        case <-time.After(time.Second * 10):
            fmt.Println("timeout command timed out")
            return
        case err := <-errChan:
            if err != nil {
                fmt.Println("timeout error:", err)
            }
    }
}
```

We first define an error channel, `errChan`, which will receive any error returned from the `cmd.Wait()` method. The command, `cmd`, is then defined, and next `cmd`’s `Start` method is called to initiate the external process. Within a Go function, we wait for the command to return using the `cmd.Wait()` method. `errChan` will only receive the error value once the command has exited and the copying to standard input and standard error has completed. Within the following `select` block, we wait to receive from two different channels. The first case waits for the time returned after 10 seconds. The second case waits for the command to complete and receive the error value. This code allows us to gracefully handle any timeout issues.

## Errors or panics with external command processes

First, let’s define the difference between errors and panics. Errors occur when the application can be recovered but is in an abnormal state. If a panic occurs, then something unexpected happened. For example, we try to access a field on a `nil` pointer or attempt to access an index that is out of bounds for an array. We can start by handling errors.

There are a couple of errors that exist within the `os/exec` package:

-   `exec.ErrDot`: Error when the file path of the command failed to resolve within the current directory, `"."`, hence the name `ErrDot`
-   `exec.ErrNotFound`: Error when the executable fails to resolve in the defined file path

You can check for the type to handle each error uniquely.

### Handling errors when a command’s path cannot be found

The following code exists within the `examples/error.go` file in the `HandlingDoesNotExistErrors` function:

```markup
cmd := exec.Command("doesnotexist", "arg1")
if errors.Is(cmd.Err, exec.ErrDot) {
    fmt.Println("path lookup resolved to a local directory")
}
if err := cmd.Run(); err != nil {
    if errors.Is(err, exec.ErrNotFound) {
        fmt.Println("executable failed to resolve")
    }
}
```

When checking the type of the command, use the `errors.Is` method, rather than checking whether `cmd.Err == exec.ErrDot` because the error is not returned directly. The `errors.Is` method checks the error chain for any occurrence of the specific error type.

### Handling other errors

Also, within the `examples/error.go` file is handling an error thrown by the command process itself. This second method, `HandlingOtherMethods`, sets the command’s standard error to a buffer that we can later use if an error is returned from the command. Let’s take a look at the code:

```markup
cmd := exec.Command(filepath.Join(os.Getenv("GOPATH"), "bin", "error"))
var out bytes.Buffer
var stderr bytes.Buffer
cmd.Stdout = &out
cmd.Stderr = &stderr
if err := cmd.Run(); err != nil {
    fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    return
}
fmt.Println(out.String())
```

When an error is encountered, we print not only the error, `exit status 1`, but also any data that has been piped into the standard error pipe, which should give the users more detail on why the error occurred.

To further understand how this code works, let’s take a look at the error command implementation that exists in the `cmd/error/error.go` file:

```markup
func main() {
    if len(os.Args) != 0 { // not passing in any arguments in this example throws an error
        fmt.Fprintf(os.Stderr, "missing arguments\n")
        os.Exit(1)
    }
    fmt.Println("executing command with no errors")
}
```

Since we are not passing any arguments into the command function, after we check the length of `os.Args`, we print to the standard error pipe the reason we are exiting with a non-zero exit code. This is a very simple way to handle errors in an effective manner. When calling this external process, we just return the errors, but as we’ve all probably experienced, error messages can be a bit cryptic. In later chapters, we will talk about how we can rewrite these to be more human-readable and provide a few examples.

In [_Chapter 4_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_04.xhtml#_idTextAnchor087), _Popular Frameworks for Building CLIs_, we discussed the use of the `RunE` function within the Cobra Command struct, which allows us to return an error value when the command is run. If you are calling an external process within the `RunE` method, then you can capture and return the error to the user, after rewriting it to a more human-readable format, of course!

Panics are handled differently than errors, but it is a good practice to, within your own code, provide a way to recover from a panic gracefully. You can see this code initiated within the `examples/panic.go` file within the `Panic` method. This calls the `panic` command, located in `cmd/panic/panic.go`. This command simply panics and then recovers. It returns the panic message to the standard error pipe, prints the stack, and exits with a non-zero exit code:

```markup
defer func() {
    if panicMessage := recover(); panicMessage != nil {
        fmt.Fprintf(os.Stderr, "(panic) : %v\n", panicMessage)
        debug.PrintStack()
        os.Exit(1)
    }
}()
panic("help!")
```

On the side that runs this command, we handle it just like any other error by capturing the error and printing data piped into the standard error.

## Timeouts and other errors with HTTP requests

Similarly, you could also experience errors when sending requests to an external API server. To be clear, timeouts are considered errors as well. The code for this example is located within `examples/http.go`, which contains two functions:

-   `HTTPTimeout()`
-   `HTTPError()`

Before we dig into the previous methods, let’s talk about the code that needs to be running in order for these methods to execute properly.

The `cmd/api/` folder contains the code for defining the handlers and starting an HTTP server locally. The `mux.HandleFunc` method defines the request pattern and matches it to the `handler` function. The server is defined by its address, which runs on localhost, port `8080`, and the `Handler`, `mux`. Finally, the `server.ListenAndServe()` method is called on the defined server:

```markup
func main() {
    mux := http.NewServeMux()
    server := &http.Server{
        Addr: ":8080",
        Handler: mux,
    }
    mux.HandleFunc("/timeout", timeoutHandler)
    mux.HandleFunc("/error", errorHandler)
    err := server.ListenAndServe()
    if err != nil {
        fmt.Println("error starting api: ", err)
        os.Exit(1)
    }
}
```

The timeout handler is defined simply. It waits two seconds before sending the response by using the `time.After(time.Second*2)` channel:

```markup
func timeoutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("got /timeout request")
    <-time.After(time.Second * 2)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("this took a long time"))
}
```

The error handler returns a status code of `http.StatusInternalServerError`:

```markup
func errorHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("got /error request")
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("internal service error"))
}
```

In a separate terminal, run the `make install` command inside the root of the repository to start the API server. Now, let’s look at the code that calls each endpoint and show how we handle it. Let’s first discuss the first type of error – the timeout:

-   `HTTPTimeout`: Inside the `examples/http.go` file resides the `HTTPTimeout` method. Let’s walk through the code together:
    -   First, we _define the client_ using the `http.Client` struct, specifying the timeout as one second. Remember that as the timeout handler on the API returns a response after two seconds, the request is sure to timeout:

```markup
client := http.Client{
    Timeout: 1 * time.Second,
}
```

-   Next, we _define the request_: a `GET` method to the `/timeout` endpoint. We pass in an empty body:

```markup
body := &bytes.Buffer{}
req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/timeout", body)
if err != nil {
    panic(err)
}
```

-   The client `Do` method is called with the request variable passed in as a parameter. We wait for the server to respond within a second and if not, an error is returned. Any errors returned by the client’s `Do` method will be of the `*url.Error` type. You can access the different fields to this error type, but in the following code, we check whether the error’s `Timeout` method returns `true`. In this statement, we can act however we’d like. We can return the error for now. We can back off and retry or we can exit. It depends on what your specific use case is:

```markup
resp, err := client.Do(req)
if err != nil {
    urlErr := err.(*url.Error)
    if urlErr.Timeout() {
        fmt.Println("timeout: ", err)
        return
    }
}
defer resp.Body.Close()
```

When this method is executed, the output is printed:

```markup
timeout:  Get "http://localhost:8080/timeout": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

A timeout is just one error, but there are many others you might encounter. Since the client `Do` method returns a particular error type in the `net/url` package, let’s discuss that. Inside the `net/url` package exists the `url.Error` type definition:

```markup
type Error struct {
    Op  string // Operation
    URL string // URL
    Err error // Error
}
```

The error contains the `Timeout()`method, which returns `true` when a request times out, and it is important to note that when the response status is anything other than `200 OK`, the error is not set. However, the status code indicates an error response. Error responses can be split into two different categories:

-   **Client error responses** (status codes range from `400` to `499`) indicate an error on the client’s side. A few examples of this include `Bad Request (400)`, `Unauthorized (401)`, and `Not` `Found (404)`.
-   **Server error messages** (status codes range from `500` to `599`) indicate an error on the server side. A few common examples of this include `Internal Server Error (500)`, `Bad Gateway (502)`, and `Service` `Unavailable (503)`.

`HTTPErrors`: Some sample code of how this can be handled exists within the `examples/http.go` file within the `HTTPErrors` method. Again, it’s important to make sure that the API server is running before executing this code:

-   The code within the method starts by calling a `GET` request to the `/``error` endpoint:

```markup
resp, err := http.Get("http://localhost:8080/error")
```

-   If the error is not `nil`, then we cast it to the `url.Error` type to access the fields and methods within it. For example, we check whether `urlError` is a timeout or a temporary network error. If it is neither, then we can output as much information as we know about the error to standard output. This additional information can help us to determine what steps to take next:

```markup
if err != nil {
    urlErr := err.(*url.Error)
    if urlErr.Timeout() {
         // a timeout is a type of error
        fmt.Println("timeout: ", err)
        return
    }
    if urlErr.Temporary() {
        // a temporary network error, retry later
        fmt.Println("temporary: ", err)
        return
    }
    fmt.Printf("operation: %s, url: %s, error: %s\n", urlErr.
        Op,        urlErr.URL, urlErr.Error())
    return
}
```

-   Since the status code error response isn’t considered a Golang error, the response body might have some useful information. If it’s not `nil`, then we can read the status code:

```markup
if resp != nil {
    defer resp.Body.Close()
```

-   We initially check that `StatusCode` doesn’t equal `http.StatusOK`. From there, we can check for particular error messages and take the appropriate action. In this example, we only check for three different types of error responses, but you can check for whichever ones make sense for what you’re doing:

```markup
if resp.StatusCode != http.StatusOK {
        // action for when status code is not okay
        switch resp.StatusCode {
        case http.StatusBadRequest:
            fmt.Printf("bad request: %v\n", resp.Status)
        case http.StatusInternalServerError:
            fmt.Printf("internal service error: %v\n", resp.
                Status)
        default:
            fmt.Printf("unexpected status code: %v\n", resp.
                StatusCode)
        }
    }
```

-   Finally, a client or server error status does not necessarily mean that the response body is `nil`. We can output the response body in case there’s any useful information we can further gather:

```markup
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("err:", err)
    }
    fmt.Println("response body:", string(data))
}
```

This concludes the section for handling HTTP timeouts and other errors. Although the examples are simple, they give you the necessary information and guidance to handle timeouts, temporary networks, and other errors.

Just Imagine

# Summary

Over the course of this chapter, you’ve learned about the `os/exec` package in depth. This included learning about the different ways to create commands: using the `command` struct or the `Command` method. Not only have we created commands, but we’ve also passed file descriptors to them to receive information back. We learned about the different ways to run a command using the `Run` or `Start` method and the multiple ways of retrieving data from the standard output, standard error types, and other file descriptors.

In this chapter, we also discussed the `net/http` and `net/url` packages, which are important to be comfortable with when creating HTTP requests to external API servers. Several examples taught us how to create requests with the methods on `http.Client`, including `Do`, `Get`, `Post`, and `PostForm`.

It’s important to learn how to build robust code, and handling errors gracefully is part of the process. We need to know how to capture errors first, so we discussed how to detect some common errors that can occur when running an external process or sending a request to an external API server. Capturing and handling other errors gives us confidence that our code is ready to take appropriate action when they occur. Finally, we now know how to check for different status codes when the response is not okay.

With all the information learned in this chapter, we should now be more confident in building a CLI that interacts with external commands or sends requests to external APIs. In the next chapter, we’ll learn how to write code that can run on multiple different architectures and operating systems.

Just Imagine

# Questions

1.  What method in the `time` package do we use to receive the time after a particular duration via a channel?
2.  What is the error type returned from `http.Client`’s `Do` method?
3.  When an HTTP request receives a response with a status code other than `StatusOK`, is the error returned from the request populated?

Just Imagine

# Answers

1.  `time.After(d Duration) <-``chan Time`
2.  `*``url.Error`
3.  No

Just Imagine

# Further reading

-   Visit the online documentation for `net/http` `at` `h`[ttps://pkg.go.dev/net/http](https://pkg.go.dev/net/http), and for net/url at h[ttps://pkg.go.dev/net/url](https://pkg.go.dev/net/url)