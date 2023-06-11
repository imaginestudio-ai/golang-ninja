# Custom Builds and Testing CLI Commands

With any Golang application, you’ll need to build and test. However, it is increasingly important as the project and its user base grow. Build tags with Boolean logic give you the ability to create targeted builds and testing and further stabilize your project with each new feature.

Given a deeper understanding of build tags and how to use them, we will use a real-world example, the audio file CLI, to integrate levels (free and pro) and enable a profiling feature.

Build tags are not only used as input when building but also when testing. We will spend the latter half of this chapter on testing. We will learn specifically how to mock an HTTP client that our CLI is using, configure tests locally, write tests for individual commands, and run them. In this chapter, we will cover the following topics in detail:

-   What are build tags and how can you use them?
-   Building with tags
-   Testing CLI commands


# What are build tags and how can you use them?

**Build tags** are indicators of when a code file should be included within the build process. In Go, they are defined by a single line at the top, or near the top, of any source file, not just a Go file. They must precede the package clause and be followed by a blank line. They have the following syntax:

```markup
//go:build [tag]
```

This line can only be defined once in a file. More than one definition would generate an error. However, when more than one tag is used, they interact using Boolean logic. In [_Chapter 7_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_07.xhtml#_idTextAnchor143), _Developing for Different Platforms_, we briefly touched on tags and their logic. The other method for handling the development of different platforms uses a series of `if-else` statements that check the operating system at runtime. Another method is to include the operating system in the filename. For example, if there’s a filename ending in `_windows.go`, we indicate to the compiler to only include this file when building for `windows`.

Tags can help separate code to include when compiling for different operating systems using `$GOOS` and `$GOARCH`. Valid combinations of operating systems and the architecture can be found here: [https://go.dev/doc/install/source#environment](https://go.dev/doc/install/source#environment).

Besides targeting platforms, build tags can be customized to separate featured code or integration tests. Often, integration tags receive a specific tag, as they often take a longer time to run. Separating unit tests from integration tests adds a level of control when testing your application.

These build constraints, when used together, can powerfully compile different versions of your code. As mentioned, they are evaluated together using Boolean logic. Expressions contain build tags combined using the `||`, `&&`, and `!` operators and parentheses. To learn more about build constraints, run the following command in your terminal:

```markup
go help buildconstraint
```

As an example, the following build tags constrain a file to build when the `linux` or `openbsd` tags are satisfied and when `amd64` is satisfied and `cgo` is not:

```markup
//go:build (linux  || openbsd) && amd64 && !cgo
```

Run `go` `env` in your terminal to see which tags are satisfied automatically when building your application. You’ll see the target operating system (`$GOOS`) and architecture (`$GOARCH`) and `unix` if the operating system is Unix or Unix-like. The `cgo` field is determined by the `CGO_ENABLED` environment variable, the term for each Go major release, and any additional tags given by the `–``tags` flag.

As mentioned earlier, you can create your own pro and free versions based on tags placed at the top of code files, `//go:build pro` or `//go:build free`. Integration test files can be tagged with `//go:build int`, for example. However you want to customize your builds, you can do so with the power of tags and Boolean logic. Now, in the next section, let’s use tags in our code to do just that.

Bookmark

# How to utilize build tags

As mentioned, we can use build tags to separate builds based on the operating system and architecture. Within the audio file repository, we’re already doing so with the following files associated with the `play` and `bug` commands. For the `bug` command, we have the following files:

-   `bug_darwin.go //` only builds on Darwin systems
-   `bug_linux.go //` only builds on Linux systems
-   `bug_windows.go //` only builds on Windows platforms

Each of those files contains a function that is specifically coded for the targeted platform. The file suffixes have similar functionality to the build tags. You can choose a file suffix that matches the exact platform and architecture. However, build tags are preferred when you want to target more than one platform and architecture. Inside the files is the matching build tag, used as an example, but duplicates functionality. Inside `bug_darwin.go`, for example, at the top of the file is the following:

```markup
//go:build darwin
```

Since we already have these build tags set up throughout the repo to target platforms where needed, let’s explore a few other ways to utilize build tags.

## Creating a pro, free, and dev version

Suppose the command-line interface utilized build tags to create different levels of access to the application’s features. This could be for admin or basic level users or restricted by the level of permissions, but it could also be, especially if the CLI was for external customers, a pro and free level version of your application.

First, it’s important to decide which commands will be available for each version. Let’s give this a try with the audio file application:

![Table 11.1 – List of commands included in the free or pro level](https://static.packt-cdn.com/products/9781804611654/graphics/image/Table_11.1_B18883.jpg)

Table 11.1 – List of commands included in the free or pro level

Let’s also include a dev version; this simply allows the API to be run locally. In a real-world scenario, the application would be configured to call a public API, and storage could be done in a database. This gives us another build tag to create.

Now, let’s use build tags to distinguish the free, pro, and dev versions. The dev version build tag is placed at the top of the `cmd/api.go` file, making the API command only available when the `dev` tag is specified:

```markup
//go:build dev
```

Then, the tag to distinguish the pro version is as follows:

```markup
//go:build !free && pro
```

There are a few files, as previously mentioned, that already have build tags to target platforms. This build tag means that the file will be available in the free, pro, and dev versions:

```markup
//go:build darwin
```

The preceding build tags utilize Boolean logic to state that the file should be included in the build process when both the `darwin` and `free` tags are defined.

Let’s break down the tags here with the Boolean logic syntax examples:

![Table 11.2 – Boolean logic examples
](https://static.packt-cdn.com/products/9781804611654/graphics/image/Table_11.2_B18883.jpg)

Table 11.2 – Boolean logic examples

This Boolean logic included within the build tag will allow developers to build for any combination of platforms and versions.

## Adding build tags to enable pprof

Another way to utilize build tags is to enable profiling on your API service. `pprof` is a tool for visualizing and analyzing profile data. The tool reads a collection of samples in `proto`, or protocol buffer, format and then creates reports that help visualize and analyze the data. This tool can generate text and graphical reports.

Note

To learn more about how to use this tool, visit [https://pkg.go.dev/net/http/pprof](https://pkg.go.dev/net/http/pprof).

For this case, we’ll define a build tag called `pprof` to appropriately match its usage. Within the `services/metadata/metadata.go` file, we define the metadata service used to extract information from the audio files uploaded via the command-line interface. The `CreateMetadataService` function creates the metadata service and defines all the endpoints with matching handlers. To enable profiling, we will add this new block of code:

```markup
if profile {
    mux.HandleFunc("/debug/pprof/", pprof.Index)
    mux.HandleFunc("/debug/pprof/{action}", pprof.Index)
    mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
}
```

At the top of the file, after the inputs, we’ll define the variable that it’s dependent on:

```markup
var (
    profile = false
)
```

However, we need some way to set the `profile` variable to `true`. To do so, we create a new file: `services/metadata/pprof.go`. This file contains the following content:

```markup
//go:build profile && (free || pro)
package metadata
func init() {
    profile = true
}
```

As you can see, whether building the `free`, `pro`, or `dev` version, if the `profile` build tag is added as tag input, then the `init` function will be called to set the `profile` variable to `true`. Now, we have another idea of how to use build tags – to set Boolean variables that act as feature flags. Now that we’ve changed the necessary files to include the build tags, let’s use these as inputs to the build commands.

Bookmark

# Building with tags

By now, we have built our applications using `Makefile`, which contains the following command specific to building a Darwin application:

```markup
build-darwin:
    go build -tags darwin -o bin/audiofile main.go
    chmod +x bin/audiofile
```

For the Darwin build, we can additionally build a version for a free and pro version and also a profile version to enable `pprof`.

## Building a free version

To build a `free` version for the Darwin operating system, we need to modify the preceding `make` command and create a new one:

```markup
build-darwin-free:
    go build -tags "darwin free" -o bin/audiofile main.go
    chmod +x bin/audiofile
```

In the `build-darwin-free` command, we pass in the two build tags: `darwin` and `free`. This will include files such as `bug_darwin.go` and `play_darwin.go`, which contain the following line at the top of the Go file:

```markup
//go:build darwin
```

Similarly, the files will be included in the build when we build the `pro` version.

## Building a pro version

To build a `pro` version for the Darwin operating system, we need to add a new `build` command:

```markup
build-darwin-pro:
    go build -tags "darwin pro" -o bin/audiofile main.go
    chmod +x bin/audiofile
```

In the `build-darwin-pro` command, we pass in the two build tags: `darwin` and `pro`.

## Building to enable pprof on the pro version

To build a `pro` version that has `pprof` enabled, we add the following `build` command:

```markup
build-darwin-pro-profile:
    go build -tags "darwin pro profile" -o bin/audiofile main.go
    chmod +x bin/audiofile
```

In the `build-darwin-pro-profile` command, we pass three build tags: `darwin`, `pro`, and `profile`. This will include the `services/metadata/pprof.go` file, which includes the line at the top of the file:

```markup
//go:build profile
```

Similarly, the files will be included in the build when we build for the free version.

At this point, we’ve learned what build tags are, the different ways to use build tags within your code, and, finally, how to build applications targeted to specific uses using build tags. Specifically, while build tags can be used to define different levels of features available (free versus pro), you can also enable profiling or any other debug tooling using build tags. Now that we have understood how to build our command-line application for different targets, let’s learn how to test our CLI commands.

Bookmark

# Testing CLI commands

While building your command-line application, it’s important to also build testing around it so you can ensure that the application works as expected. There are a few things that typically need to be done, including the following:

1.  Mock the HTTP client
2.  Handle test configuration
3.  Create a test for each command

We’ll go over the code for each of these steps that exist in the audio file repository for [_Chapter 11_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_11.xhtml#_idTextAnchor258).

## Mocking the HTTP client

To mock the HTTP client, we’ll need to create an interface to mimic the client’s `Do` method, as well as a function that returns this interface, which is both satisfied by the real and mocked client.

In the `cmd/client.go` file, we’ve written some code to handle all of this:

```markup
type AudiofileClient interface {
    Do(req *http.Request) (*http.Response, error)
}
var (
    getClient = GetHTTPClient()
)
func GetHTTPClient() AudiofileClient {
    return &http.Client{
        Timeout: 15 * time.Second,
    }
}
```

We can now easily create a mock client by replacing the `getClient` variable with a function that returns a mocked client. If you look at each command’s code, it uses the `getClient` variable. For example, the `upload.go` file calls the `Do` method with the following line:

```markup
resp, err := getClient.Do(req)
```

When the application runs, this returns the actual HTTP client with a 15-second timeout. However, in each test, we’ll set the `getClient` variable to a mocked HTTP client.

The mocked HTTP client is set in the `cmd/client_test.go` file. First, we define the type:

```markup
type ClientMock struct {
}
```

Then, to satisfy the `AudiofileClient` interface previously defined, we implement the `Do` method:

```markup
func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
```

Some of the requests, including `list`, `get`, and `search` endpoints, will return data that is stored in JSON files under the `cmd/testfiles` folder. We read these files and store them in the corresponding byte slices: `listBytes`, `getBytes`, and `searchBytes`:

```markup
listBytes, err := os.ReadFile("./testfiles/list.json")
if err != nil {
    return nil, fmt.Errorf("unable to read testfile/list.json")
}
getBytes, err := os.ReadFile("./testfiles/get.json")
if err != nil {
    return nil, fmt.Errorf("unable to read testfile/get.json")
}
searchBytes, err := os.ReadFile("./testfiles/search.json")
if err != nil {
    return nil, fmt.Errorf("unable to read testfile/search.json")
}
```

The data read from these files is used within the response. Since the `Do` method receives the request, we can create a switch case for each request endpoint and then handle the response individually. You can create more detailed cases to handle errors, but in this case, we are only returning the successful case. For the first case, the `/request` endpoint, we return `200 OK`, but the body of the response also contains the string value from `getBytes`. You can see the actual data in the `./``testfiles/get.json` file:

```markup
    switch req.URL.Path {
         case "/request":
             return &http.Response{
                 Status:  "OK",
                 StatusCode: http.StatusOK,
                 Body: ioutil.NopCloser(bytes.NewBufferString(string(getBytes))),
      ContentLength: int64(len(getBytes)),
      Request: req,
      Header: make(http.Header, 0),
  }, nil
```

For the `/upload` endpoint, we return `200 OK`, but the body of the response also contains the `"123"` string value:

```markup
         case "/upload":
             return &http.Response{
                 Status:  "OK",
                 StatusCode: http.StatusOK,
      Body: ioutil.NopCloser(bytes.NewBufferString("123")),
      ContentLength: int64(len("123")),
      Request: req,
      Header: make(http.Header, 0),
   }, nil
```

For the `/list` endpoint, we return `200 OK`, but the body of the response also contains the string value from `listBytes`. You can see the actual data in the `./``testfiles/list.json` file:

```markup
        case "/list":
            return &http.Response{
                Status:  "OK",
                StatusCode: http.StatusOK,
                Body: ioutil.NopCloser(bytes.
                      NewBufferString(string(listBytes))),
                      ContentLength: int64(len(listBytes)),
                      Request: req,
                      Header: make(http.Header, 0),
 }, nil
```

For the `/delete` endpoint, we return `200 OK`, but the body of the response also contains `"successfully deleted audio with` `id: 456"`:

```markup
        case "/delete":
            return &http.Response{
                Status:  "OK",
                StatusCode: http.StatusOK,
                Body: ioutil.NopCloser(bytes.
                      NewBufferString("successfully deleted 
                        audio with id: 456")),
                      ContentLength: int64(len("successfully 
                                     deleted audio with id: 
                                     456")),
                      Request: req,
                      Header: make(http.Header, 0),
}, nil
```

For the `/search` endpoint, we return `200 OK`, but the body of the response also contains the string value from `searchBytes`. You can see the actual data in the `./``testfiles/search.json` file:

```markup
        case "/search":
            return &http.Response{
                Status:  "OK",
                StatusCode: http.StatusOK,
                Body: ioutil.NopCloser(bytes.
                NewBufferString(string(searchBytes))),
                ContentLength: int64(len(list searchBytes 
                Bytes)),
                Request: req,
                Header: make(http.Header, 0),
}, nil
}
return &http.Response{}, nil
}
```

Finally, if the request path doesn’t match any of the endpoints in the `switch` statement, then an empty response is returned.

## Handling test configuration

We handle the test configuration in the `cmd/root_test.go` file:

```markup
var Logger *zap.Logger
var Verbose *zap.Logger
func ConfigureTest() {
    getClient = &ClientMock{}
    viper.SetDefault("cli.hostname", "testHostname")
    viper.SetDefault("cli.port", 8000)
    utils.InitCLILogger()
}
```

Within the `ConfigureTest` function, we set the `getClient` variable to a pointer to the `ClientMock` type. Because the `viper` configuration values are checked when the command is called, we set some default values for the CLI’s hostname and port to random test values. Finally, in this file, the regular logger, `Logger`, and verbose logger, `Verbose`, are both defined and then later initialized by the `utils.InitCLILogger()` method call.

## Creating a test for a command

Now that we have the mocked client, configuration, and loggers set up, let’s create a test for the commands. Before I dive into the code for each, it’s important to mention the line of code that’s reused at the start of each test:

```markup
ConfigureTest()
```

The preceding section discusses the details of this function, but it prepares each state with a mocked client, default configuration values, and initialized loggers. In our examples, we use the `testing` package, which provides support for automated tests in Go. It is designed to be used in concert with the `go test` command, which executes any function in your code defined with the following format:

```markup
func TestXxx(*testing.T)
```

`Xxx` can be replaced with anything else, but the first character needs to be capital. The name itself is used to identify the type of test that is being executed. I won’t go into each individual test, just three as examples. To view the entirety of tests, visit the audio file repository for this chapter.

### Testing the bug command

The function for testing the `bug` command is defined here. It takes a single parameter, which is a pointer to the `testing.T` type, and fits the function format defined in the last section. Let’s break down the code:

```markup
func TestBug(t *testing.T) {
    ConfigureTest()
    b := bytes.NewBufferString("")
    rootCmd.SetOut(b)
    rootCmd.SetArgs([]string{"bug", "unexpected"})
    err := rootCmd.Execute()
    if err != nil {
        fmt.Println("err: ", err)
    }
    actualBytes, err := ioutil.ReadAll(b)
    if err != nil {
        t.Fatal(err)
    }
    expectedBytes, err := os.ReadFile("./testfiles/bug.txt")
    if err != nil {
        t.Fatal(err)
    }
    if strings.TrimSpace(string(actualBytes)) != strings.
       TrimSpace(string(expectedBytes)) {
        t.Fatal(string(actualBytes), "!=", 
          string(expectedBytes))
    }
}
```

In this function, we first define the output buffer, `b`, which we can later read for comparison to the expected output. We set the arguments using the `SetArgs` method and pass in an unexpected argument. The command is executed with the `rootCmd.Execute()` method and the actual result is read from the buffer and saved in the `actualBytes` variable. The expected output is stored within the `./testfiles/bug.txt` file and is read into the `expectedBytes` variable. We compare these values to ensure that they are equal. Since we passed in an unexpected argument, the command usage is printed out. This test is designed to pass; however, if the trimmed strings are not equal, the test fails.

### Testing the get command

The function for testing the `get` command is defined here. Similarly, the function definition fits the format to be picked up in the `go test` command. Remember the mocked client and that the `get` command calls the `/request` endpoint. The response body contains the value found in the `./testfiles/get.json` file. Let’s break down the code:

```markup
func TestGet(t *testing.T) {
    ConfigureTest()
    b := bytes.NewBufferString("")
    rootCmd.SetOut(b)
```

We pass in the following arguments to mimic the `audiofile get –id 123 –``json` call:

```markup
    rootCmd.SetArgs([]string{"get", "--id", "123", "--json"})
```

We execute the root command with the preceding arguments:

```markup
    err := rootCmd.Execute()
    if err != nil {
        fmt.Println("err: ", err)
    }
```

We read the actual data output from `rootCmd`’s execution and store it in the `actualBytes` variable:

```markup
    actualBytes, err := ioutil.ReadAll(b)
    if err != nil {
        t.Fatal(err)
    }
```

We read the expected data output from the `./``testfiles/get.json` file:

```markup
    expectedBytes, err := os.ReadFile("./testfiles/get.json")
    if err != nil {
        t.Fatal(err)
    }
```

Then, the data of both `actualBytes` and `expectedBytes` is unmarshalled into the `models.Audio` struct and then compared:

```markup
    var audio1, audio2 models.Audio
    json.Unmarshal(actualBytes, &audio1)
    json.Unmarshal(expectedBytes, &audio2)
    if !(audio1.Id == audio2.Id &&
    audio1.Metadata.Tags.Album == audio2.Metadata.Tags.Album &&
    audio1.Metadata.Tags.AlbumArtist == audio2.Metadata.Tags.AlbumArtist &&
    audio1.Metadata.Tags.Artist == audio2.Metadata.Tags.Artist &&
    audio1.Metadata.Tags.Comment == audio2.Metadata.Tags.Comment &&
    audio1.Metadata.Tags.Composer == audio2.Metadata.Tags.Composer &&
    audio1.Metadata.Tags.Genre == audio2.Metadata.Tags.Genre &&
    audio1.Metadata.Tags.Lyrics == audio2.Metadata.Tags.Lyrics &&
    audio1.Metadata.Tags.Year == audio2.Metadata.Tags.Year) {
        t.Fatalf("expected %q got %q", string(expectedBytes), string(actualBytes))
    }
}
```

This test was designed to succeed, but if the data is not as expected, then the test fails.

### Testing the upload command

The function for testing the `upload` command is defined here. Again, the function definition fits the format to be picked up in the `go test` command. Remember the mocked client and that the `upload` command calls the `/upload` endpoint with a mocked response body containing the `"123"` value. Let’s break down the code:

```markup
func TestUpload(t *testing.T) {
    ConfigureTest()
    b := bytes.NewBufferString("")
    rootCmd.SetOut(b)
    rootCmd.SetArgs([]string{"upload", "--filename", "list.
                   go"})
    err := rootCmd.Execute()
    if err != nil {
        fmt.Println("err: ", err)
    }
    expected := "123"
    actualBytes, err := ioutil.ReadAll(b)
    if err != nil {
        t.Fatal(err)
    }
    actual := string(actualBytes)
    if !(actual == expected) {
        t.Fatalf("expected \"%s\" got \"%s\"", expected, 
                actual)
    }
}
```

`rootCmd`’s arguments are set to mimic the following command call:

```markup
audiofile upload –filename list.go
```

The file type and data are not validated because that happens on the API side, which is mocked. However, since we know the body of the response contains the `123` value, we set the expected variable to `123`. The `actual` value, which contains the output of the command execution, is then later compared to the expected one. The test is designed for success, but if the values are not equal, then the test fails.

We’ve now gone over several examples of how to test a CLI Cobra command. You can now create your own tests for your CLI, by mocking your own HTTP client and creating tests for each individual command. We haven’t done so in this chapter, but it’s good to know that build tags can also be used to separate different kinds of tests – for example, integration tests and unit tests.

## Running the tests

To test your commands, you can run `go test` and pass in a few additional flags:

-   `-v` for verbose mode
-   `-tags` for any files you want to specifically target

In our test, we want to target just the `pro` build tag because that will cover all commands. We add two additional `Makefile` commands, one to run tests in verbose mode and one that doesn’t:

```markup
test:
  go test ./... -tags pro
test-verbose:
  go test –v ./... -tags pro
```

After saving the `Makefile` from the terminal, you can execute the command:

```markup
make test
```

The following output is expected:

```markup
go test ./cmd -tags pro
ok      github.com/marianina8/audiofile/cmd
```

We now know how to run the tests utilizing build tags as well. This should be all the tools needed to run your own CLI testing.

Bookmark

# Summary

In this chapter, you learned what build tags are and how to use them for different purposes. Build tags can be used for generating builds of different levels, separating our specific tests, or adding debug features. You also learned how to generate builds with the build tags that you added to the top of your files and how to utilize the Boolean logic of tags to quickly determine whether files will or won’t be included.

You also learned how to test your Cobra CLI commands with Golang’s default `testing` package. Some necessary tools were also included, such as learning how to mock an HTTP client. Together with the build tags, you can now not only build targeted applications with tags but also run tests with the same tags to target specific tests. In the next chapter, [_Chapter 12_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_12.xhtml#_idTextAnchor291), _Cross-Compilation Across Different Platforms_, we will learn how to use these tags and compile for the different major operating systems: `darwin`, `linux`, and `windows`.

Bookmark

# Questions

1.  Where does the build tag go in a Golang file and what is the syntax?
2.  What flag is used for both `go build` and `go test` to pass in the build tags?
3.  What build tag could you place on an integration test Golang file and how would you run `go test` with the tag?

Bookmark

# Answers

1.  It’s placed at the top of the file, before the package declaration, followed by a single empty line. The syntax is: `//``go:build [tag]`.
2.  The `–tags` flag is used to pass in build tags for both the `go build` and `go` `test` methods.
3.  You could add the `//go:build int` build tag at the top of any integration test file, and then modify the test file to run this command: `go test ./cmd -tags "``pro int"`.

Bookmark

# Further reading

-   Read more about the `build` package at [https://pkg.go.dev/go/build](https://pkg.go.dev/go/build), and read more about the `testing` package at [https://pkg.go.dev/testing](https://pkg.go.dev/testing)