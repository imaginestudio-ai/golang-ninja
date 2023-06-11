# Building an Audio Metadata CLI

_Hands-on learning_ is one of the best ways to learn. So, in this chapter, we will build out a few of our example audio metadata CLI use cases from start to finish. The code is available online and can be explored alongside this chapter or independently. Forking the GitHub repo and playing around with the code, adding in new use cases and tests, are encouraged as these are excellent ways to learn before diving into some of the ways to refine your CLI in the following chapters.

Although this example covered in this chapter is not built on an empty code base – it is built on top of an existing REST API – it’s worth noting that the implementation of commands does not necessarily rely on an API. This is only an example and it’s encouraged that you use your imagination in this chapter on how commands could be implemented if not relying on an API. This chapter will give you an experimental code base and you’ll learn about the following topics:

-   Defining the components
-   Implementing use cases
-   Testing and mocking

Bookmark

# Technical requirements

Download the following code to follow along:

[https://github.com/ImagineDevOps DevOps/Building-Modern-CLI-Applications-in-Go/tree/main/Chapter03/audiofile](https://github.com/ImagineDevOps DevOps/Building-Modern-CLI-Applications-in-Go/tree/main/Chapter03/audiofile)[](https://github.com/ImagineDevOps DevOps/Building-Modern-CLI-Applications-in-Go/tree/main/Chapter03)

Install the latest version of VS Code with the latest Go tools.

Bookmark

# Defining the components

The following is the folder structure for our audio metadata CLI. The main folders in this structure were described in the last chapter. Here, we will go into further detail on what each folder contains, and the files and code that exist within them, in order from top to bottom:

```markup
   |--cmd
   |----api
   |----cli
   |------command
   |--extractors
   |----tags
   |----transcript
   |--internal
   |----interfaces
   |--models
   |--services
   |----metadata
   |--storage
   |--vendor
```

## cmd/

As previously mentioned in [_Chapter 2_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_02.xhtml#_idTextAnchor036), _Structuring Go Code for CLI Applications_ in the _Commonly used program layouts for robust applications_ section, the `cmd` folder is the main entry point for the different applications of the project.

## cmd/api/

The `main.go` file, which in found in the `cmd/api/` folder, will start to run the audio metadata API locally on the machine. It takes in a port number as an optional flag, defaulting to `8000`, and passes the port number into a `Run` method within the `services` method that starts the metadata service:

```markup
package main
import (
    metadataService "audiofile/services/metadata"
    "flag"
    "fmt"
)
func main() {
    var port int
    flag.IntVar(&port, "p", 8000, "Port for metadata
      service")
    flag.Parse()
    fmt.Printf("Starting API at http://localhost:%d\n",
      port)
    metadataService.Run(port)
}
```

We make use of the `flag` package, which implements simple command-line flag parsing. There are different flag types that can be defined, such as `String`, `Bool`, and `Int`. In the preceding example, a `-p` flag is defined to override the default port of `8000`. `flag.Parse()` is called after all the flags are defined to parse the line into defined flags. There are a few syntactical methods allowed for passing flags to the command using Go’s `flag` package. The value `8080` will be parsed either way:

```markup
-p=8080
-p 8080  // this works for non-boolean flags only
```

Sometimes, a flag does not require an argument and is enough on its own for the code to know exactly what to do:

```markup
-p
```

Action can be taken on the flag that’s passed in, but the variable will contain the default value, `8000`, when defined.

To start the API from the project’s root directory, run `go run cmd/api/main.go` and you will see the following output:

```markup
audiofile go run cmd/api/main.go
Starting API at http://localhost:8000
```

## cmd/cli/

This `main.go` file, in the `cmd/cli/` folder, runs the CLI, and like many other CLIs, this one will utilize the API by making calls to it. Since the API will need to be running for the CLI to work, run the API first in a separate terminal or in the background. The `cmd/cli/main.go` file contains the following code:

```markup
package main
import (
    "audiofile/internal/command"
    "audiofile/internal/interfaces"
    "fmt"
    "net/http"
    "os"
)
func main() {
    client := &http.Client{}
    cmds := []interfaces.Command{
        command.NewGetCommand(client),
        command.NewUploadCommand(client),
        command.NewListCommand(client),
    }
    parser := command.NewParser(cmds)
    if err := parser.Parse(os.Args[1:]); err != nil {
        os.Stderr.WriteString(fmt.Sprintf("error: %v",
          err.Error()))
        os.Exit(1)
    }
}
```

Within the `main.go` file, the commands are added to a slice of interface `Command` type. Each command is defined and added:

```markup
command.NewGetCommand(client),
command.NewUploadCommand(client),
command.NewListCommand(client),
```

Each command takes the `client` variable, a default `http.Client` as a parameter to use to make HTTP requests to the audio metadata API endpoints. Passing in the `client` command allows it to be easily mocked for testing, which we will discuss in the next section.

The commands are then passed into a `NewParser` method, which creates a pointer to `command.Parser`:

```markup
parser := command.NewParser(cmds)
```

This `Parse` function receives all arguments after the application name via the `os.Args[1:]` parameter value. For example, say the command line is called as follows:

```markup
 ./audiofile-cli upload -filename recording.m4v
```

Then, the first argument, `os.Args[0]`, returns the following value:

```markup
 audiofile-cli
```

To explain this further, let’s look at the `Command` struct and the fields present within it:

![Figure 3.1 – Command struct and flag.FlagSet entities](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_3.01.jpg)

Figure 3.1 – Command struct and flag.FlagSet entities

Let us look at the `GetCommand` struct depicted in the figure:

```markup
type GetCommand struct {
    fs *flag.FlagSet
    client interfaces.Client
    id string
}
```

Each of the commands has a flag set, which contains a name for the command and error handling, a client, and an ID.

The arguments to a Go program are stored in the `os.Args` slice, which is a collection of strings. The name of the executable being run is stored in the first element of the `os.Args` slice (i.e., `os.Args[0]`), while the arguments passed to the executable are stored in the subsequent elements (`os.Args[1:]`).

When you see the code, `parser.Parse(os.Args[1:])`, it means you’re passing the remainder of the command-line arguments to `parse.Parse` function, skipping the first argument (the name of the program). All the arguments on the command line, besides the program’s name, will be passed to the function in this case.

That means when we pass in `os.Args[1:]`, we are passing into `parse.Parse` all the arguments after the program name:

```markup
 upload –filename recording.m4v
```

Parse takes `args`, a string list, and returns an `error` type. The function converts command-line parameters into executable commands.

Let’s walk through the code alongside the following flow chart:

-   It checks for less than 1 args. If so, `help()` returns n`il.`
-   `Args[0]` is assigned to subcommand if the slice has at least one item. This shows the user’s command.
-   The function then cycles over the `Parser` struct’s `p.commands` property. It checks each command’s name (obtained by executing the `Name()` method) against the `subcommand` variable.
-   The function executes the command’s `ParseFlags` method with the rest of the `args` slice if a match is found (`args[1:]`). Finally, the function runs the command and returns the result.
-   If no match is found, the method returns an unknown subcommand error message using the `fmt.Errorf` function.

Essentially, the code finds and executes a command from command line arguments. Then, the matching command is run.

![Figure 3.2 – Flow diagram for the Parse method](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_3.02.jpg)

Figure 3.2 – Flow diagram for the Parse method

A command exists for each API endpoint. For example, `UploadCommand` will call the `/upload` endpoint, `ListCommand` will call the `/list` endpoint, and `GetCommand` will call the `/get` endpoint of the REST API.

Within the `Parse` method, the length of `args` is checked. If no arguments are passed, then help is printed and the program returns `nil`:

```markup
audiofile ./audiofile-cli
usage: ./audiofile-cli <command> [<args>]
These are a few Audiofile commands:
    get      Get metadata for a particular audio file by id
    list     List all metadata
    upload   Upload audio file
```

## cmd/cli/command

In the `cmd/cli/command` folder, there are commands to match each of the audiofile API endpoints. In the next section, we will code the `upload`, `list`, and `get` commands to implement a couple of the use cases described in the previous chapter. Rather than defining the code for one of these commands here, I’ll provide a structure used to define a random command that satisfies the `Command` interface:

```markup
package command
import (
    "github.com/marianina8/ audiofile/internal/cli"
    "github.com/marianina8/ audiofile/internal/interfaces"
    "flag"
    "fmt"
)
func NewRandomCommand(client interfaces.Client)
    *RandomCommand {
    gc := &RandomCommand{
        fs: flag.NewFlagSet("random",
           flag.ContinueOnError),
        client: client,
    }
    gc.fs.StringVar(&gc.flag, "flag", "", "string flag for
      random command")
    return gc
}
type RandomCommand struct {
    fs *flag.FlagSet
    flag string
}
func (cmd *RandomCommand) Name() string {
    return cmd.fs.Name()
}
func (cmd *RandomCommand) ParseFlags(flags []string) error {
    return cmd.fs.Parse(flags)
}
func (cmd *RandomCommand) Run() error {
    fmt.Println(rand.Intn(100))
    return nil
}
```

The `upload`, `get`, and `list` commands follow the same structure, but the implementation of the constructor and `Run` methods differ.

Also, in the `cmd/cli/command` folder, there is a parser of the struct type with a method to parse the arguments, match them with the commands, and parse any flags found after the subcommand. The `NewParser` function creates a new instance of the `Parser` struct. It takes a slice of type `[]interfaces.Command` as input and returns a pointer to a `Parser` struct. This initialization method provides an easy way to set up the struct with a set of desired commands. The following is the code inside `parser.go`:

```markup
package command
import (
    "github.com/marianina8/audiofile/internal/interfaces"
    "fmt"
)
type Parser struct {
    commands []interfaces.Command
}
func NewParser(commands []interfaces.Command) *Parser {
    return &Parser{commands: commands}
}
func (p *Parser) Parse(args []string) error {
    if len(args) < 1 {
        help()
        return nil
    }
    subcommand := args[0]
    for _, cmd := range p.commands {
        if cmd.Name() == subcommand {
            cmd.ParseFlags(args[1:])
            return cmd.Run()
        }
    }
    return fmt.Errorf("Unknown subcommand: %s", subcommand)
}
```

The code checks the number of arguments passed to the `Parse` method. If the number of arguments is less than 1, a `help` function from a separate `help.go` file is called to print the help text to guide the user on proper usage:

```markup
func help() {
    help := `usage: ./audiofile-cli <command> [<flags>]
These are a few Audiofile commands:
    get      Get metadata for a particular audio file by id
    list     List all metadata
    upload   Upload audio file
    `
    fmt.Println(help)
}
```

## extractors/

This folder contains implementations for the different extractors of audio metadata. In this case, subfolders exist for the `tags` and `transcript` implementations.

## extractors/tags

The `tags` package is implemented within the `extractors/tags` folder. Tags metadata may include title, album, artists, composer, genre, release year, lyrics, and any additional comments. The code is available within the GitHub repository and utilizes the `github.com/dhowden/tag` Go package.

## extractors/transcript

The `transcript` package is implemented within the `extractors/transcript` folder. Like the other extraction package, the code can be found in the GitHub repository. However, transcript analysis is handled by AssemblyAI, a third-party API, and requires an API key, which can be set locally to `ASSEMBLY_API_KEY`.

## internal/interfaces

The `internal/interfaces` folder holds interfaces utilized by the application. It includes both the `Command` and `Storage` interfaces. Interfaces provide a way for developers to create multiple types that meet the same interface specifications allowing flexibility and modularity in the design of the application. The `storage.go` file defines the storage interface:

```markup
package interfaces
import (
    "audiofile/models"
)
type Storage interface {
    Upload(bytes []byte, filename string) (string, string,
      error)
    SaveMetadata(audio *models.Audio) error
    List() ([]*models.Audio, error)
    GetByID(id string) (*models.Audio, error)
    Delete(id string, tag string) error
}
```

The preceding interface satisfies all possible use cases. Specific implementations can be defined in the `storage` folder. If you choose to define the storage type within a configuration, you can easily swap out implementations and switch from one storage type to another. In this example, we define flat file storage with an implementation of each method to satisfy the interface.

First utilized in the `cmd/cli/main.go` file, the `Command` interface is defined by the following code in `internal/interfaces/command.go`:

```markup
type Command interface {
   ParseFlags([]string) error
   Run() error
   Name() string
}
```

Notice how each of the commands in the `cmd/cli/command/` folder implements the preceding interface.

## models/

The `models` folder contains a structs shared across the different applications. The first struct defined for the `audiofile` application is `Audio`:

```markup
type Audio struct {
    Id       string
    Path     string
    Metadata Metadata
    Status   string
    Error    []error
}
```

The `Id` variable contains the unique **identifier(ID)**, of the `Audio` file. The path the stored local copy of the audio file. The `Metadata` variable contains the data extracted from the audio file. In the following example, tags and speech-to-text transcript data are being stored:

```markup
type Metadata struct {
    Tags         Tags         `json:"tags"`
    Transcript   string       `json:"transcript"`
}
```

It’s not necessary to know the struct for each extraction type. The most important thing is the main entity type, `Audio`, and its value field, `Metadata`.

## services/metadata

Although multiple services could be implemented in the `services` folder, we’re currently only utilizing one API service, the audio metadata service. The only method that exists in the `metadata.go` file is the `CreateMetadataServer` method, which is called in the `metadata` package, and the `Run` method, which is called from the `cmd/api/main.go` file. This file also contains the struct for `MetadataService`:

```markup
type MetadataService struct {
    Server *http.Server
    Storage interfaces.Storage
}
```

`CreateMetadataService` takes an argument, a port of the `int` type, to define the server’s port running on localhost. It also takes an argument, `storage`, which is an implementation of the `Storage` interface. The handlers that declare each endpoint of the API server are also defined. This function returns a pointer to `MetadataService`:

```markup
func CreateMetadataService(port int, storage
   interfaces.Storage) *MetadataService {
    mux := http.NewServeMux()
    metadataService := &MetadataService{
        Server: &http.Server{
            Addr:    fmt.Sprintf(":%v", port),
            Handler: mux,
        },
        Storage: storage,
    }
    mux.HandleFunc("/upload",
      metadataService.uploadHandler)
    mux.HandleFunc("/request",
      metadataService.getByIDHandler)
    mux.HandleFunc("/list", metadataService.listHandler)
    return metadataService
}
```

The `Run` method, which takes an argument, `port`, defined by the value of the `p` flag or the default value of `8000`, calls the `CreateMetadataService` method and initiates running the server by calling the `ListenAndServer` method on the server. Any error with starting the API will be returned immediately:

```markup
func Run(port int) {
    flatfileStorage := storage.FlatFile{}
    service:= CreateMetadataService(port, flatfileStorage)
    err := service.Server.ListenAndServe()
    if err != nil {
        fmt.Println("error starting api: ", err)
    }
}
```

Implementations of each of the handlers will be discussed in the next section when handling a few use cases.

## storage/

In the `storage` folder, there is the `flatfile.go` file, which implements a method of storing metadata locally to a flat file organized via ID on the local disk. The code implementation of this will not be discussed in this book because it goes beyond the scope of focus on the CLI. However, you can view the code in the GitHub repository.

## vendor/

The `vendor` directory holds all direct and indirect dependencies.

Bookmark

# Implementing use cases

Remember the use cases defined in the previous chapter? Let’s try to implement a couple of them:

-   UC-01 Upload audio
-   UC-02 Request metadata

## Uploading audio

In this use case, an authenticated user uploads an audio file by giving the location of the file on their device for the purpose of extracting its metadata. Under the hood, the upload process will save a local copy and run the metadata extraction process on the audio file. A unique ID for the audio file is returned immediately.

Before we begin to implement this use case, let’s consider what the command for uploading may look like. Suppose we’ve settled on the following final command structure:

```markup
./audiofile-cli upload -filename <filepath>
```

Since `/cmd/cli/main.go` is already defined, we’ll just need to make sure that the `upload` command exists and satisfies the `command` interface, with the `ParseFlags`, `Run`, and `Name` methods. In the `internal/command` folder, we define the `upload` command in the `upload.go` file within the `command` package:

```markup
package command
import (
    "github.com/marianina8/audiofile/internal/interfaces"
    "bytes"
    "flag"
    "fmt"
    "io"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
)
func NewUploadCommand(client interfaces.Client)
    *UploadCommand {
    gc := &UploadCommand{
        fs:     flag.NewFlagSet("upload",
                  flag.ContinueOnError),
        client: client,
    }
    gc.fs.StringVar(&gc.filename, "filename", "", "full
      path of filename to be uploaded")
    return gc
}
type UploadCommand struct {
    fs       *flag.FlagSet
    client   interfaces.Client
    filename string
}
func (cmd *UploadCommand) Name() string {
    return cmd.fs.Name()
}
func (cmd *UploadCommand) ParseFlags(flags []string)
  error {
    if len(flags) == 0 {
        fmt.Println("usage: ./audiofile-cli
          upload -filename <filename>")
        return fmt.Errorf("missing flags")
    }
    return cmd.fs.Parse(flags)
}
func (cmd *UploadCommand) Run() error {
    // implementation for upload command
    return nil
}
```

The `NewUploadCommand` method implements our desired command structure by defining a new flag set for the `upload` command:

```markup
flag.NewFlagSet("upload", flag.ContinueOnError)
```

This method call passes the string, `upload`, into the method’s `name` parameter and flag. `ContinueOnError` in the `flag.ErrorHandling` parameter defines how the application should react if an error occurs when parsing the flag. The different, and mostly self-explanatory, options for handling errors upon parsing include the following:

-   `flag.ContinueOnError`
-   `flag.ExitOnError`
-   `flag.PanicOnError`

Now that we’ve defined and added the `upload` command, we can test it out. Upon testing, you’ll see that the `upload` command runs without an error but exits immediately with no response. Now, we are ready to implement the `Run` method of the `upload` command.

When we first started implementing a CLI for the audiofile application, an API already existed. We discussed how this API starts and runs `MetadataServer`, which handles requests to a few existing endpoints. For this use case, we are concerned with the `http://localhost/upload` endpoint.

With this in mind, let’s delve deeper into the documentation for the upload endpoint of this REST API so we will know exactly how to construct a `curl` command.

### Uploading audio

In order to upload audio, we’ll need to know how to communicate with the API to handle certain tasks. Here are the details required to design a request to handle uploading audio:

-   **Method**: `POST`
-   **Endpoint**: `http://localhost/upload`
-   **Header**: `Content-Type: multipart/form-data`
-   **Form data**: `Key ("file") Value (bytes) Name(base` `of filename)`

Make sure that the API is running, and then test out the endpoint using `curl`. Immediately, the ID of the uploaded file is returned:

```markup
curl --location --request POST 'http://localhost/upload' \
--form 'file=@"recording.mp3"'
8a6dc954-d6df-4fc0-882e-14eb1581d968%
```

After successfully testing out the API endpoint, we can write the Go code that handles the same functionality as the previous `curl` command within the `Run` method of `UploadCommand`.

The new `Run` method can now be defined. The method supplies the filename that’s been passed into the `upload` command as a flag parameter and saves the bytes of that file to a multipart form `POST` request to the `http://localhost/upload` endpoint:

```markup
func (cmd *UploadCommand) Run() error {
    if cmd.filename == "" {
        return fmt.Errorf("missing filename")
    }
    fmt.Println("Uploading", cmd.filename, "...")
    url := "http://localhost/upload"
    method := "POST"
    payload := &bytes.Buffer{}
    multipartWriter := multipart.NewWriter(payload)
    file, err := os.Open(cmd.filename)
    if err != nil {
        return err
    }
    defer file.Close()
    partWriter, err := multipartWriter
      .CreateFormFile("file", filepath.Base(cmd.filename))
    if err != nil {
        return err
    }
    _, err = io.Copy(partWriter, file)
    if err != nil {
        return err
    }
    err = multipartWriter.Close()
    if err != nil {
        return err
    }
    client := cmd.client
    req, err := http.NewRequest(method, url, payload)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type",
      multipartWriter.FormDataContentType())
    res, err := client.Do(req)
    if err != nil {
        return err
    }
    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return err
    }
    fmt.Println("Audiofile ID: ", string(body))
    return err
}
```

The first CLI command, `upload`, has been implemented! Let’s implement another use case, requesting metadata by ID.

## Requesting metadata

In the requesting metadata use case, an authenticated user requests audio metadata by the audio file’s ID. Under the hood, the request metadata process will, within the flat file storage implementation, search for the `metadata.json` file corresponding with the audio file and decode its contents into the `Audio` struct.

Before implementing the request metadata use case, let’s consider what the command for requesting metadata will look like. The final command structure will look like this:

```markup
./audiofile-cli get -id <ID>
```

For simplification, `get` is the command to request metadata. Let’s define the new `get` command, and in `/cmd/cli/main.go`, confirm that it is present in the list of commands to recognize when the application is run. The structure for defining the `get` command is similar to that of the first command, `upload`:

```markup
package command
import (
    "github.com/marianina8/audiofile/internal/interfaces"
    "bytes"
    "flag"
    "fmt"
    "io"
    "net/http"
    "net/url"
)
func NewGetCommand(client interfaces.Client) *GetCommand {
    gc := &GetCommand{
        fs:     flag.NewFlagSet("get",
                  flag.ContinueOnError),
        client: client,
    }
    gc.fs.StringVar(&gc.id, "id", "", "id of audiofile
      requested")
    return gc
}
type GetCommand struct {
    fs     *flag.FlagSet
    client interfaces.Client
    id     string
}
func (cmd *GetCommand) Name() string {
    return cmd.fs.Name()
}
func (cmd *GetCommand) ParseFlags(flags []string) error {
    if len(flags) == 0 {
        fmt.Println("usage: ./audiofile-cli get -id <id>")
        return fmt.Errorf("missing flags")
    }
    return cmd.fs.Parse(flags)
}
func (cmd *GetCommand) Run() error {
    // implement get command
    return nil
}
```

The `NewGetCommand` method implements our desired command structure by defining a new flag set for the `get` command, `flag.NewFlagSet("get", flag.ContinueOnError)`. This method receives the string, `get`, in the method’s `name` parameter and `flag.ContinueOnError` in the `flag.ErrorHandling` parameter.

Let’s delve deeper into the documentation for the get endpoint of this REST API so we will know exactly how to construct a curl command.

### Requesting metadata

In order to request audio metadata, we’ll need to know how to communicate with the API to handle this task. Here are the details required to design a request for audio metadata:

-   **Method**: `GET`
-   **Endpoint**: `http://localhost/get`
-   **Query parameter**: `id` – ID of audio file

Make sure that the API is running, and then test out the `get` endpoint using `curl`. Immediately, the metadata of the requested audio file is returned in JSON format. This data could be returned in different formats, and we could add an additional flag to determine the format of the returned metadata:

```markup
curl --location --request GET
'http://localhost/request?id=270c3952-0b48-4122-bf2a-
 e4a005303ecb'
{audiofile metadata in JSON format}
```

After confirming that the API endpoint works as expected, we can write the Go code that handles the same functionality as the preceding `curl` command within the `Run` method of `GetCommand`. The new `Run` method can now be defined:

```markup
func (cmd *GetCommand) Run() error {
    if cmd.id == "" {
        return fmt.Errorf("missing id")
    }
    params := "id=" + url.QueryEscape(cmd.id)
    path := fmt.Sprintf("http://localhost/request?%s",
      params)
    payload := &bytes.Buffer{}
    method := "GET"
    client := cmd.client
    req, err := http.NewRequest(method, path, payload)
    if err != nil {
        return err
    }
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    b, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("error reading response: ",
          err.Error())
        return err
    }
    fmt.Println(string(b))
    return nil
}
```

Now that the request metadata use case has been implemented, let’s compile the code and test out the first couple of CLI commands: `upload`, for uploading and processing audio metadata, and `get`, for requesting metadata by audiofile ID.

Giving the CLI a more specific name, `audiofile-cli`, let’s generate the build by running the following command:

```markup
go build -o audiofile-cli cmd/cli/main.go
```

Bookmark

# Testing a CLI

Now that we have successfully built the CLI application, we can do some testing to make sure that it’s working. We can test out the commands we’ve created and then write out proper tests to make sure any future changes don’t break the current functionality.

## Manual testing

To upload an audio file, we’ll run the following command:

```markup
./audiofile-cli upload -filename audio/beatdoctor.mp3
```

The result is as expected:

```markup
Uploading audio/beatdoctor.mp3 ...
Audiofile ID:  8a6a8942-161e-4b10-bf59-9d21785c9bd9
```

Now that we have the audiofile ID, we can immediately get the metadata, which will change as the metadata updates after each extraction process. The command for requesting metadata is as follows:

```markup
./audiofile-cli get -id=8a6a8942-161e-4b10-bf59-
9d21785c9bd9
```

The result is the populated `Audio` struct in JSON format:

```markup
{
    "Id": "8a6a8942-161e-4b10-bf59-9d21785c9bd9",
    "Path": "/Users/marian/audiofile/8a6a8942-161e-4b10-
    bf59-9d21785c9bd9/beatdoctor.mp3",
    "Metadata": {
        "tags": {
            "title": "Shot In The Dark",
            "album": "Best Bytes Volume 4",
            "artist": "Beat Doctor",
            "album_artist": "Toucan Music (Various
              Artists)",
            "genre": "Electro House",
            "comment": "URL: http://freemusicarchive.org/
            music/Beat_Doctor/Best_Bytes_Volume_4/
            09_beat_doctor_shot_in_the_dark\r\nComments:
            http://freemusicarchive.org/\r\nCurator: Toucan
            Music\r\nCopyright: Attribution-NonCommercial
            3.0 International: http://creativecommons.org/
            licenses/by-nc/3.0/"
        },
        "transcript": "This is Sharon."
    },
    "Status": "Complete",
    "Error": null
}
```

The results are as expected. However, not all audio passed into the CLI will return the same data. This is just an example. Some audio may not have any tags at all and transcription will be skipped if you don’t have the `ASSEMBLYAI_API_KEY` environment variable set with an AssemblyAI API key. Ideally, API keys should not be set as environment variables, which can be leaked easily, but this is a temporary option. In [_Chapter 4_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_04.xhtml#_idTextAnchor087), _Popular Frameworks for Building CLIs_, you will learn about Viper, which is a configuration library that pairs perfectly with the Cobra CLI framework.

## Testing and mocking

Now, we can start writing some unit tests. In the `main.go` file, there is a root function that parses the arguments passed into the application. Using VS Code and the extension for Go support, you can right-click on a function and see an option for generating unit tests, **Go: Generate Unit Tests** **For Function**.

![Figure 3.3 –Screenshot of VS Code menu of Go options](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_3.03.jpg)

Figure 3.3 –Screenshot of VS Code menu of Go options

Select the `Parse` function in the `commands` package and then click on the option to generate the following table-driven unit tests inside the `parser_test.go` file, we can see the test function for the parsing functionality:

```markup
func TestParser_Parse(t *testing.T) {
    type fields struct {
        commands []interfaces.Command
    }
    type args struct {
        args []string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := &Parser{
                commands: tt.fields.commands,
            }
            if err := p.Parse(tt.args.args); (err != nil)
              != tt.wantErr {
                t.Errorf("Parser.Parse() error = %v,
                  wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

This provides a great template for us to implement some tests given different argument and flag combinations utilized in the method. When running the tests, we don’t want the client to call the REST endpoints, so we mock the client and fake responses. We do all this inside the `parser_test.go` file. Since each of the commands takes in a client, we can easily mock the interface. This is done in the file using the following code:

```markup
type MockClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}
func (m *MockClient) Do(req *http.Request) (*http.Response,
    error) {
    if strings.Contains(req.URL.String(), "/upload") {
        return &http.Response{
            StatusCode: 200,
            Body:       io.NopCloser
              (strings.NewReader("123")),
        }, nil
    }
    if strings.Contains(req.URL.String(), "/request") {
        value, ok := req.URL.Query()["id"]
        if !ok || len(value[0]) < 1 {
             return &http.Response{
    StatusCode: 500,
    Body: io.NopCloser(strings.NewReader("url param 'id' is 
    missing")),
    }, fmt.Errorf("url param 'id' is missing")
        }
        if value[0] != "123" {
            return &http.Response{
                StatusCode: 500,
                Body:       io.NopCloser
               (strings.NewReader("audiofile id does not
                 exist")),
            }, fmt.Errorf("audiofile id does not exist")
        }
        file, err := os.ReadFile("testdata/audio.json")
        if err != nil {
            return nil, err
        }
        return &http.Response{
            StatusCode: 200,
            Body:       io.NopCloser
              (strings.NewReader(string(file))),
        }, nil
    }
    return nil, nil
}
```

The `MockClient` interface is satisfied by `http.DefaultClient`. The `Do` method is mocked. Within the `Do` method, we check which endpoint is being called (`/upload` or`/get`) and respond with the mock response. In the preceding example, any call to the `/upload` endpoint responds with an `OK` status and a string, `123`, representing the ID of the audio file. A call to the `/get` endpoint checks the IDs passed in as a URL parameter. If the ID matches the audiofile ID of `123`, then the mocked client will return a successful response with the audio JSON in the body of the response. If there is a request for any ID other than `123`, then a status code of 500 is returned with an error message that the ID does not exist.

Now that the mocked client is complete, we fill in success and failure cases for each command, `upload` and `get`, within the `Parse` function’s unit tests:

```markup
func TestParser_Parse(t *testing.T) {
    mockClient := &MockClient{}
    type fields struct {
        commands []interfaces.Command
    }
    type args struct {
        args []string
    }
```

The `tests` variable contains an array of data that contains the name of the test, the fields or commands available, the string arguments potentially passed into the command-line application, and a `wantErr` Boolean value that is set depending on whether we expect an error to be returned in the test or not. Let’s go over each test:

```markup
   tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
```

The first test, named `upload – failure – does not exist`, simulates the following command:

```markup
./audiofile-cli upload -filename doesNotExist.mp3
```

The filename, `doesNotExist.mp3`, is a file that does not exist in the root folder. Within the `Run()` method of the `upload` command, the file is opened. This is where the error occurs and the output is an error message, `file does` `not exist`:

```markup
        {
            name: "upload - failure - does not exist",
            fields: fields{
                commands: []interfaces.Command{
                    NewUploadCommand(mockClient),
                },
            },
            args: args{
                args: []string{"upload", "-filename",
                  "doesNotExist.mp3"},
            },
            wantErr: true, // error = open
              doesNotExist.mp3: no such file or directory
        },
```

The test named `upload – success – uploaded` checks the successful case of a file being uploaded to storage with an audiofile ID being returned in response. In order to get this test to work, there is a `testdata` folder in the `command` package, and within it exists a small audio file to test with, simulating the following command:

```markup
./audiofile-cli upload -filename testdata/exists.mp3
```

This file is successfully opened and sent to the `/upload` endpoint. The mocked client’s `Do` function sees that the request is to the `/upload` endpoint and sends an `OK` status along with the audiofile ID of `123` within the body of the response and no error. This matches the `wantErr` value of `false`:

```markup
        {
            name: "upload - success - uploaded",
            fields: fields{
                commands: []interfaces.Command{
                    NewUploadCommand(mockClient),
                },
            },
            args: args{
                args: []string{"upload", "-filename", "
                  testdata/exists.mp3"},
            },
            wantErr: false,
        },
```

After uploading, we can now _get_ the metadata associated with the audiofile. The next test case, `get – failure – id does not exist`, tests a request for an audiofile ID that does not exist. Instead of passing in `123`, that is, that ID of an audiofile that exists, we pass in an ID that does not exist, simulating the following command via the CLI:

```markup
./audiofile-cli get -id 567
```

`wantErr` is set to `true` and we get the expected error, `audiofile id does not exist`. The response from the `/request` endpoint returns the error message in the body of the response.

```markup
        {
            name: "get - failure - id does not exist",
            fields: fields{
                commands: []interfaces.Command{
                    NewGetCommand(mockClient),
                },
            },
            args: args{
                args: []string{"get", "-id", "567"},
            },
            wantErr: true, // error = audiofile id does not
              exist
        },
```

The test named `get – success – requested` checks whether the `get` command was successful in retrieving an ID of an audiofile that exists. The ID passed is `"123"`, and in the mocked client, you can see that when that specific ID is passed into the request, the API endpoint returns a 200 success code with the body of the audiofile metadata.

This is simulated with the following command:

```markup
./audiofile-cli get -id 123
        {
            name: "get - success - requested",
            fields: fields{
                commands: []interfaces.Command{
                    NewGetCommand(mockClient),
                },
            },
            args: args{
                args: []string{"get", "-id", "123"},
            },
            wantErr: false,
        },
    }
```

The following code loops through the previously described `tests` array to run each test with the arguments passed into the command and checks whether the final `wantErr` value matches the expected error:

```markup
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := &Parser{
                commands: tt.fields.commands,
            }
            if err := p.Parse(tt.args.args); (err != nil)
              != tt.wantErr {
                t.Errorf("Parser.Parse() error = %v,
                  wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

To run these tests, from the repository type the following:

```markup
go test ./cmd/cli/command -v
```

This will execute all the preceding tests and print the following output:

```markup
--- PASS: TestParser_Parse (0.06s)
    --- PASS: TestParser_Parse/upload_-_failure_-
    _does_not_exist (0.00s)
    --- PASS: TestParser_Parse/upload_-_success_-_uploaded
    (0.06s)
    --- PASS: TestParser_Parse/get_-_failure_-
    _id_does_not_exist (0.00s)
    --- PASS: TestParser_Parse/get_-_success_-_requested
    (0.00s)
PASS
ok      github.com/marianina8/audiofile/cmd/cli/command
(cached)
```

It’s important to test success and failure cases for all the commands. Although this was just a starting example; more test cases could be added. For example, in the previous chapter, we discussed the upload use case in more detail. You could test it with large files that exceed the limit, or whether the file passed into the `upload` command is an audio file. In the state that the current implementation is in, a large file would successfully upload. Since this is not what we want, we can modify the `UploadCommand` `Run` method to check the size of the file before calling the request to the `/upload` endpoint. However, this is just an example and hopefully gives you an idea of how a CLI can be built alongside an existing API.

Bookmark

# Summary

Throughout this chapter, we have gone through an example of building an audio metadata CLI. Going through each of the different components that make up this CLI has helped us to determine how a CLI could be structured and how files are structured, whether as part of an existing code base or as a new CLI.

We learned how to implement the first two main use cases of the CLI, uploading audio and getting audio metadata. The details provided on the structure of the commands gave you an idea of how commands could be built out without the use of any additional parsing packages. You also learned how to implement a use case, test your CLI, and mock a client interface.

While this chapter gave you an idea of how to build a CLI, some commands such as nested subcommands and flag combinations can get complicated. In the next chapter, we’ll discuss how to use some popular frameworks to help parse complicated commands and improve the CLI development process overall. You’ll see how these frameworks can exponentially speed up the development of a new CLI!

Bookmark

# Questions

1.  What are the benefits of using a storage interface? If you were to use a different storage option, how easy would it be to swap out for the current flat file storage implementation?
2.  What’s the difference between an argument and a flag? In the following real-world example, what qualifies as an argument or a flag?
    
    ```markup
    ./audiofile-cli upload -filename music.mp3
    ```
    
3.  Suppose you’d like to create an additional test for when a user runs the `get` command without passing in any arguments or flags:
    
    ```markup
    ./audiofile-cli get
    ```
    

What would an additional entry to the `tests` array look like?

Bookmark

# Answers

1.  Interfaces benefit us when writing modular code that’s decoupled and reduces dependency across different parts of the code base. Since we have an interface, it’s much easier to swap out the implementation. In the existing code, you’d swap the implementation type in the `Run` method of the `metadata` package.
2.  In this example:
    
    ```markup
     ./audiofile-cli upload -filename music.mp3
    ```
    

`upload`, `-filename`, and `music.mp3` are all considered arguments. However, flags are specific arguments that are specifically marked by a specific syntax. In this case, `-filename` is a flag.

1.  An additional test for when a user runs the `get` command without passing in any arguments or flags would look like this:
    
    ```markup
         {
            name: "get - failure - missing required id flag",
            fields: fields{
                commands: cmds,
            },
            args: args{
                args: []string{"get"},
            },
            wantErr: true,
         },
    ```