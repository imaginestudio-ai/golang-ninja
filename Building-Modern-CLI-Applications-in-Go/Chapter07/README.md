# Developing for Different Platforms

One of the main reasons Go is such a powerful language for building a command-line application is how easy it is to develop an application that can be run on multiple machines. Go provides several packages that allow developers to write code that interacts with the computer independent of the specific operating system. These packages include `os`, `time`, `path`, and `runtime`. In the first section, we will discuss some commonly used functions in each of these packages and then provide some simple examples to pair with the explanations.

To further drill down the importance of these files, we will revisit the `audiofile` code and implement a couple of new features that utilize some of the methods that exist in these packages. After all, the best way to learn is by implementing new features with the new functions and methods you’ve learned about.

We will then learn how to use the `runtime` library to check the operating system the application is running on and then use that to switch between codes. By learning about build tags, what they are, and how to use them, we will learn about a cleaner way to switch between code blocks to implement a new feature that can be run on three different operating systems: Darwin, Windows, and Linux. By the end of the chapter, you’ll feel more confident when building your application, knowing that the code you are writing will work seamlessly, independent of the platform.

In this chapter, we will cover the following key topics:

-   Packages for platform-independent functionality
-   Implementing independent or platform-specific code
-   Build tags for targeted platforms


# Packages for platform-independent functionality

When you are building a **command-line interface** (**CLI**) that will be shared with the public, it’s important that the code is platform-independent to support users who are running the CLI on different operating systems. Golang has supportive packages that provide platform-independent interfaces to operating system functionality. A few of these packages include `os`, `time`, and `path`. Another useful package is the `runtime` package, which helps when detecting the operating system the application is running on, among other things. We will review each of these packages with some simple examples to show how to apply some of the available methods.

## The os package

The **operating system** (**os**) package has a Unix-like design but applies uniformly across all operating systems. Think of all the operating system commands you can run in a shell, including external commands. The `os` package is your go-to package. We discussed calling external commands in the previous chapter; now we will discuss this at a higher level and focus on the commands in certain groups: environmental, file, and process operations.

### Environmental operations

As the name suggests, the `os` package contains functions that give us information about the environment in which the application is running, as well as change the environment for future method calls. These common operations are for the following working directories:

-   `func Chdir(dir string) error`: This changes the current working directory
-   `func Getwd() (dir string, err error)`: This gets the current working directory

There are also operations for the environment, as follows:

-   `func Environ() []string`: This lists environment keys and values
-   `func Getenv(key string) string`: This gets environment variables by key
-   `func Setenv(key, value string) error`: This sets environment variables by key and value
-   `func Unsetenv(key string) error`: This unsets an environment variable by key
-   `func Clearenv()`: This clears environment variables
-   `func ExpandEnv(s string) string`: This expands values of environment variable keys in strings to their values

The [_Chapter 7_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_07.xhtml#_idTextAnchor143) code exists on GitHub in the `environment.go` file, where we have provided some sample code demonstrating using these operations:

```markup
func environment() {
    dir, err := os.Getwd()
    if err != nil {
        fmt.Println("error getting working directory:", err)
    }
    fmt.Println("retrieved working directory: ", dir)
    fmt.Println("setting WORKING_DIR to", dir)
    err = os.Setenv("WORKING_DIR", dir)
    if err != nil {
        fmt.Println("error setting working directory:", err)
    }
    fmt.Println(os.ExpandEnv("WORKING_DIR=${WORKING_DIR}"))
    fmt.Println("unsetting WORKING_DIR")
    err = os.Unsetenv("WORKING_DIR")
    if err != nil {
        fmt.Println("error unsetting working directory:", err)
    }
    fmt.Println(os.ExpandEnv("WORKING_DIR=${WORKING_DIR}"))
    fmt.Printf("There are %d environment variables:\n", len(os.
        Environ()))
    for _, envar := range os.Environ() {
        fmt.Println("\t", envar)
    }
}
```

To briefly describe the preceding code, we first get the working directory, then set it to the `WORKING_DIR` environment variable. To show the change, we utilize `os.ExpandEnv` to print the key-value pair. We then unset the `WORKING_DIR` environment variable. Again, we show it is unset by using `os.ExpandEnv` to print out the key-value pair. The `os.ExpandEnv` variable will print an empty string if the environment variable is unset. Finally, we print out the count of the environment variables and then range through all to print them. Running the preceding code will produce the following output:

```markup
retrieved working directory:  /Users/mmontagnino/Code/src/github.com/marianina8/Chapter-7
setting WORKING_DIR to /Users/mmontagnino/Code/src/github.com/marianina8/Chapter-7
WORKING_DIR=/Users/mmontagnino/Code/src/github.com/marianina8/Chapter-7
There are 44 environment variables.
key=WORKING_DIR, value=/Users/mmontagnino/Code/src/github.com/marianina8/Chapter-7
unsetting WORKING_DIR
WORKING_DIR=
```

If you run this code on your machine rather than Linux, Unix, or Windows, the resulting output will be similar. Try for yourself.

Notes on running the following examples

To run the [_Chapter 7_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_07.xhtml#_idTextAnchor143) examples, you’ll first need to run the install command to install the sleep command to your GOPATH. On Unix-like systems, run the `make install` command followed by the `make run` command. On Linux systems, run the `./build-linux.sh` script followed by the `./run-linux.sh` script. On Windows, run `.\build-windows.ps1` followed by the `.\run-windows.ps1` Powershell script.

### File operations

The `os` package also offers a wide variety of file operations that can be applied universally across different operating systems. Many functions and methods can be applied to files, so rather than going over each by name, I will group the functionality and name a few of each:

-   The following can be used to change file, directory, and link permissions and owners:
    -   `func Chmod(name string, mode` `FileMode) error`
    -   `func Chown(name string uid, gid` `int) error`
    -   `func Lchown(name string uid, gid` `int) error`
-   The following can be used to create pipes, files, directories, and links:
    -   `func Pipe() (r *File, w *File,` `err error)`
    -   `func Create(name string) (*``File, error)`
    -   `func Mkdir(name string, perm` `FileMode) error`
    -   `func Link(oldname, newname` `string) error`
-   The following are used to read from files, directories, and links:
    -   `func ReadFile(name string) ([]``byte, error)`
    -   `func ReadDir(name string) ([]``DirEntry, error)`
    -   `func Readlink(name string) (``string, error)`
-   The following retrieve user-specific data:
    -   `func UserCacheDir() (``string, error)`
    -   `func UserConfigDir() (``string, error)`
    -   func UserHomeDir() (string, error)
-   The following are used to write to files:
    -   func (f \*File) Write(b \[\]byte) (n int, err error)
    -   func (f \*File) WriteString(s string) (n int, err error)
    -   `func WriteFile(name string, data []byte, perm` `FileMode) error`
-   The following are used for file comparison:
    -   `func SameFile(fi1, fi2` `FileInfo) bool`

There is a `file.go` file within the [_Chapter 7_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_07.xhtml#_idTextAnchor143) code on GitHub in which we have some sample code using these operations. Within the file are multiple functions, the first, `func createFiles() error`, handles the creation of three files to play around with:

```markup
func createFiles() error {
    filename1 := "file1"
    filename2 := "file2"
    filename3 := "file3"
    f1, err := os.Create(filename1)
    if err != nil {
        return fmt.Errorf("error creating %s: %v\n", filename1, 
          err)
    }
    defer f1.Close()
    f1.WriteString("abc")
    f2, err := os.Create(filename2)
    if err != nil {
        return fmt.Errorf("error creating %s: %v\n", filename2, 
          err)
    }
    defer f2.Close()
    f2.WriteString("123")
    f3, err := os.Create(filename3)
    if err != nil {
        return fmt.Errorf("error creating %s: %v", filename3, 
          err)
    }
    defer f3.Close()
    f3.WriteString("xyz")
    return nil
}
```

The `os.Create` method allows file creation to work seamlessly on different operating systems. The next function, `file()`, utilizes these files to show how to use methods that exist within the `os` package. The `file()` function primarily gets or changes the current working directory and runs different functions, including the following:

-   `func createExamplesDir() (string, error)`: This creates an `examples` directory in the user’s home directory
-   `func printFiles(dir string) error`: This prints the files/directories under the directory represented by `dir string`
-   `func sameFileCheck(f1, f2 string) error`: This checks whether two files, represented by the `f1` and `f2` strings are the same file

Let’s first show the `file()` function to get the overall gist of what is going on:

```markup
originalWorkingDir, err := os.Getwd()
if err != nil {
    fmt.Println("getting working directory: ", err)
}
fmt.Println("working directory: ", originalWorkingDir)
examplesDir, err := createExamplesDir()
if err != nil {
    fmt.Println("creating examples directory: ", err)
}
err = os.Chdir(examplesDir)
if err != nil {
    fmt.Println("changing directory error:", err)
}
fmt.Println("changed working directory: ", examplesDir)
workingDir, err := os.Getwd()
if err != nil {
    fmt.Println("getting working directory: ", err)
}
fmt.Println("working directory: ", workingDir)
createFiles()
err = printFiles(workingDir)
if err != nil {
    fmt.Printf("Error printing files in %s\n", workingDir)
}
err = os.Chdir(originalWorkingDir)
if err != nil {
    fmt.Println("changing directory error: ", err)
}
fmt.Println("working directory: ", workingDir)
symlink := filepath.Join(originalWorkingDir, "examplesLink")
err = os.Symlink(examplesDir, symlink)
if err != nil {
    fmt.Println("error creating symlink: ", err)
}
fmt.Printf("created symlink, %s, to %s\n", symlink, examplesDir)
err = printFiles(symlink)
if err != nil {
    fmt.Printf("Error printing files in %s\n", workingDir)
}
file := filepath.Join(examplesDir, "file1")
linkedFile := filepath.Join(symlink, "file1")
err = sameFileCheck(file, linkedFile)
if err != nil {
    fmt.Println("unable to do same file check: ", err)
}
// cleanup
err = os.Remove(symlink)
if err != nil {
    fmt.Println("removing symlink error: ", err)
}
err = os.RemoveAll(examplesDir)
if err != nil {
    fmt.Println("removing directory error: ", err)
}
```

Let’s walk through the preceding code. First, we get the current working directory and print it out. Then, we call the `createExamplesDir()` function and change direction into it.

We then get the current working directory after we change it to ensure it’s now the `examplesDir` value. Next, we call the `createFiles()` function to create those three files inside the `examplesDir` folder and call the `printFiles()` function to list the files in the `examplesDir` working directory.

We change the working directory back to the original working directory and create a `symlink` to the `examplesDir` folder under the home directory. We print the files existing under the `symlink` to see that they are equal.

Next, we take `file0` from `examplesDir` and `file0` from `symlink` and compare them within the `sameFileCheck` function to ensure they are equal.

Finally, we run some cleanup functions to remove the `symlink` and `examplesDir` folders.

The `file` function utilizes many methods available in the `os` package, from getting the working directory to changing it, creating a `symlink`, and removing files and directories. Showing the separate function call code will give more uses of the `os` package. First, let’s show the code for `createExamplesDir`:

```markup
func createExamplesDir() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("getting user's home directory: 
          %v\n", err)
    }
    fmt.Println("home directory: ", homeDir)
    examplesDir := filepath.Join(homeDir, "examples")
    err = os.Mkdir(examplesDir, os.FileMode(int(0777)))
    if err != nil {
        return "", fmt.Errorf("making directory error: %v\n", 
          err)
    }
    fmt.Println("created: ", examplesDir)
    return examplesDir, nil
}
```

The preceding code uses the `os` package when getting the user’s home directory with the `os.UserHomeDir` method and then creates a new folder with the `os.Mkdir` method. The next function, `printFiles`, gets the files to print from the `os.ReadDir` method:

```markup
func printFiles(dir string) error {
    files, err := os.ReadDir(dir)
    if err != nil {
        return fmt.Errorf("read directory error: %s\n", err)
    }
    fmt.Printf("files in %s:\n", dir)
    for i, file := range files {
        fmt.Printf(" %v %v\n", i, file.Name())
    }
    return nil
}
```

Lastly, `sameFileCheck` takes two files represented by strings, `f1` and `f2`. To get the file info for each file, the `os.Lstat` method is called on the file string. `os.SameFile` takes this file info and returns a `boolean` value to symbolize the result – `true` if the files are the same and `false` if not:

```markup
func sameFileCheck(f1, f2 string) error {
    fileInfo0, err := os.Lstat(f1)
    if err != nil {
        return fmt.Errorf("getting fileinfo: %v", err)
    }
    fileInfo0Linked, err := os.Lstat(f2)
    if err != nil {
        return fmt.Errorf("getting fileinfo: %v", err)
    }
    isSameFile := os.SameFile(fileInfo0, fileInfo0Linked)
    if isSameFile {
        fmt.Printf("%s and %s are the same file.\n", fileInfo0.
            Name(), fileInfo0Linked.Name())
    } else {
    fmt.Printf("%s and %s are NOT the same file.\n", fileInfo0.
        Name(), fileInfo0Linked.Name())
    }
    return nil
}
```

This concludes the code samples utilizing methods from the `os` package related to file operations. Next, we will discuss some operations related to processes running on the machine.

### Process operations

When calling external commands, we can get a **process ID** (**pid**), associated with the process. Within the `os` package, we can perform actions on the process, send the process signals, or wait for the process to complete and then receive a process state with information regarding the process that was completed. In the [_Chapter 7_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_07.xhtml#_idTextAnchor143) code, we have a `process()` function, which utilizes some of the following methods for processes and process states:

-   `func Getegid() int`: This returns the effective group ID of the caller. Note, this is not supported in Windows, the concept of group IDs is specific to Unix-like or Linux systems. For example, this will return `–1` on Windows.
-   `func Geteuid() int`: This returns the effective user ID of the caller. Note, this is not supported in Windows, the concept of user IDs is specific to Unix-like or Linux systems. For example, this will return `-1` on Windows.
-   `func Getpid() int`: This gets the process ID of the caller.
-   `func FindProcess(pid int) (*Process, error)`: This returns the process associated with the `pid`.
-   `func (p *Process) Wait() (*ProcessState, error)`: This returns the process state when the process completes.
-   `func (p *ProcessState) Exited() bool`: This returns `true` if the process exited.
-   `func (p *ProcessState) Success() bool`: This returns `true` if the process exited successfully.
-   `func (p *ProcessState) ExitCode() int`: This returns the exit code of the process.
-   `func (p *ProcessState) String() string`: This returns the process state in string format.

The code is as follows and starts with several print line statements that return the caller’s effective group, user, and process ID. Next, a `cmd` sleep command is defined. The command is started and from the `cmd` value, and we get the pid:

```markup
func process() {
    fmt.Println("Caller group id:", os.Getegid())
    fmt.Println("Caller user id:", os.Geteuid())
    fmt.Println("Process id of caller", os.Getpid())
    cmd := exec.Command(filepath.Join(os.Getenv("GOPATH"), 
           "bin", "sleep"))
    fmt.Println("running sleep for 1 second...")
    if err := cmd.Start(); err != nil {
        panic(err)
    }
    fmt.Println("Process id of sleep", cmd.Process.Pid)
    this, err := os.FindProcess(cmd.Process.Pid)
    if err != nil {
        fmt.Println("unable to find process with id: ", cmd.
            Process.Pid)
    }
    processState, err := this.Wait()
    if err != nil {
        panic(err)
    }
    if processState.Exited() && processState.Success() {
        fmt.Println("Sleep process ran successfully with exit 
            code: ", processState.ExitCode())
    } else {
        fmt.Println("Sleep process failed with exit code: ", 
            processState.ExitCode())
    }
    fmt.Println(processState.String())
}
```

From the process' pid, we then can find the process using the `os.FindProcess` method. We call the `Wait()` method in the process to get `os.ProcessState`. This `Wait()` method, like the `cmd.Wait()` method, waits for the process to complete. Once completed, the process state is returned. We can check whether the process state is exited with the `Exited()` method and whether it was successful with the `Success()` method. If so, we print that the process ran successfully along with the exit code, which we get from the `ExitCode()` method. Finally, the process state can be printed cleanly with the `String()` method.

## The time package

Operating systems provide access to time via two different types of internal clocks:

-   **A wall clock**: This is used for telling the time and is subject to variations due to clock synchronization with the **Network Time** **Protocol** (**NTP**)
-   **A monotonic clock**: This is used for measuring time and is not subject to variations due to clock synchronization

To be more specific on the variations, if the wall clock notices that it is moving faster or slower than the NTP, it will adjust its clock rate. The monotonic clock will not adjust. When measuring durations, it’s important to use the monotonic clock. Luckily with Go, the `Time` struct contains both the wall and monotonic clocks, and we don’t need to specify which is used. Within the [_Chapter 7_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_07.xhtml#_idTextAnchor143) code, there is a `timer.go` file, which shows how to get the current time and duration, regardless of the operating system:

```markup
func timer() {
    start := time.Now()
    fmt.Println("start time: ", start)
    time.Sleep(1 * time.Second)
    elapsed := time.Until(start)
    fmt.Println("elapsed time: ", elapsed)
}
```

When running the following code, you’ll see a similar output:

```markup
start time:  2022-09-24 23:47:38.964133 -0700 PDT m=+0.000657043
elapsed time:  -1.002107875s
```

Also, many of you have also seen that there is a `time.Now().Unix()` method. It returns to the epoch time, or time that has elapsed since the Unix epoch, January 1, 1970, UTC. These methods will work similarly regardless of the operating system and architecture they are run on.

## The path package

When developing a command-line application for different operating systems, you’ll most likely have to deal with handling file or directory path names. In order to handle these appropriately across different operating systems, you’ll need to use the `path` package. Because this package does not handle Windows paths with drive letters or backslashes, as we used in the previous examples, we’ll use the `path/filepath` package.

The `path/filepath` package uses either forward or back slashes depending on the operating system. Just for fun, within the [_Chapter 7_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_07.xhtml#_idTextAnchor143) `walking.go` file, I’ve used the `filepath` package to walk through a directory. Let’s look at the code:

```markup
func walking() {
    workingDir, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    dir1 := filepath.Join(workingDir, "dir1")
    filepath.WalkDir(dir1, func(path string, d fs.DirEntry, err 
      error) error {
        if !d.IsDir() {
            contents, err := os.ReadFile(path)
            if err != nil {
                return err
            }
            fmt.Printf("%s -> %s\n", d.Name(), 
                string(contents))
        }
        return nil
    })
}
```

We get the current working directory with `os.Getwd()`. Then create a path for the `dir1` directory that can be used for any operating system using the `filepath.Join` method. Finally, we walk the directory using `filepath.WalkDir` and print out the filename and its contents.

## The runtime package

The final package to discuss within this section is the `runtime` package. It’s mentioned because it’s used to easily determine the operating system the code is running on and therefore execute blocks of code, but there’s so much information you can get from the `runtime` system:

-   `GOOS`: This returns the running application's operating system target
-   `GOARCH:` This returns the running application’s architecture target
-   `func GOROOT() string`: This returns the root of the Go tree
-   `Compiler`: This returns the name of the compiler toolchain that built the binary
-   `func NumCPU() int`: This returns the number of logical CPUs usable by the current process
-   `func NumGoroutine() int`: This returns the number of goroutines that currently exist
-   `func Version() string`: This returns the Go tree’s version string

This package will provide you with enough information to understand the `runtime` environment. Within the [_Chapter 7_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_07.xhtml#_idTextAnchor143) code in the `checkRuntime.go` file is the `checkRuntime` function, which puts each of these into practice:

```markup
func checkRuntime() {
    fmt.Println("Operating System:", runtime.GOOS)
    fmt.Println("Architecture:", runtime.GOARCH)
    fmt.Println("Go Root:", runtime.GOROOT())
    fmt.Println("Compiler:", runtime.Compiler)
    fmt.Println("No. of CPU:", runtime.NumCPU())
    fmt.Println("No. of Goroutines:", runtime.NumGoroutine())
    fmt.Println("Version:", runtime.Version())
    debug.PrintStack()
}
```

Running the code will provide a similar output to the following:

```markup
Operating System: darwin
Architecture: amd64
Go Root: /usr/local/go
Compiler: gc
No. of CPU: 10
No. of Goroutines: 1
Version: go1.19
goroutine 1 [running]:
runtime/debug.Stack()
        /usr/local/go/src/runtime/debug/stack.go:24 +0x65
runtime/debug.PrintStack()
        /usr/local/go/src/runtime/debug/stack.go:16 +0x19
main.checkRuntime()
        /Users/mmontagnino/Code/src/github.com/marianina8/Chapter-7/checkRuntime.go:17 +0x372
main.main()
        /Users/mmontagnino/Code/src/github.com/marianina8/Chapter-7/main.go:9 +0x34
```

Now that we have learned about some of the packages required for building a command-line application that runs across multiple operating systems and architectures, in the next section, we’ll return to the `audiofile` CLI from previous chapters and implement a few new functions and show how the methods and functions we’ve learned in this section can come into play.

Just Imagine

# Implementing independent or platform-specific code

The best way to learn is to put what has been learned into practice. In this section, we’ll revisit the `audiofile` CLI to implement a few new commands. In the code for the new features we’ll implement, the focus will be on the use of the `os` and `path`/`filepath` packages.

## Platform-independent code

Let’s now implement a few new features for the `audiofile` CLI that will run independently of the operating system:

-   `Delete`: This deletes stored metadata by ID
-   `Search`: This searches stored metadata for a specific search string

The creation of each of these new feature commands was initiated with the cobra-CLI; however, the platform-specific code is isolated in the `storage/flatfile.go` file, which is the flat file storage for the storage interface.

First, let’s show the `Delete` method:

```markup
func (f FlatFile) Delete(id string) error {
    dirname, err := os.UserHomeDir()
    if err != nil {
        return err
    }
    audioIDFilePath := filepath.Join(dirname, "audiofile", id)
    err = os.RemoveAll(audioIDFilePath)
    if err != nil {
        return err
    }
    return nil
}
```

The flat file storage is stored under the user’s home directory under the `audiofile` directory. Then, as each new audio file and matching metadata is added, it is stored within its unique identifier ID. From the `os` package, we use `os.UserHomeDir()` to get the user’s home directory and then use the `filepath.Join` method to create the required path to delete all the metadata and files associated with the ID independent of the operating system. Make sure you have some audiofiles stored locally in the flat file storage. If not, add a few files. For example, use the `audio/beatdoctor.mp3` file and upload using the following command:

```markup
./bin/audiofile upload --filename audio/beatdoctor.mp3
```

The ID is returned after a successful upload:

```markup
Uploading audio/beatdoctor.mp3 ...
Audiofile ID:  a5d9ab11-6f5f-4da0-9307-a3b609b0a6ba
```

You can ensure that the data has been added by running the `list` command:

```markup
./bin/audiofile list
```

The `audiofile` metadata is returned, so we have double-checked its existence in storage:

```markup
    {
        "Id": "a5d9ab11-6f5f-4da0-9307-a3b609b0a6ba",
        "Path": "/Users/mmontagnino/audiofile/a5d9ab11-6f5f-4da0-9307-a3b609b0a6ba/beatdoctor.mp3",
        "Metadata": {
            "tags": {
                "title": "Shot In The Dark",
                "album": "Best Bytes Volume 4",
                "artist": "Beat Doctor",
                "album_artist": "Toucan Music (Various Artists)",
                "composer": "",
                "genre": "Electro House",
                "year": 0,
                "lyrics": "",
                "comment": "URL: http://freemusicarchive.org/music/Beat_Doctor/Best_Bytes_Volume_4/09_beat_doctor_shot_in_the_dark\r\nComments: http://freemusicarchive.org/\r\nCurator: Toucan Music\r\nCopyright: Attribution-NonCommercial 3.0 International: http://creativecommons.org/licenses/by-nc/3.0/"
            },
            "transcript": ""
        },
        "Status": "Complete",
        "Error": null
    },
```

Now, we can delete it:

```markup
./bin/audiofile delete --id a5d9ab11-6f5f-4da0-9307-a3b609b0a6ba
success
```

Then confirm that it’s been deleted by trying to get the audio by ID:

```markup
./bin/audiofile get --id a5d9ab11-6f5f-4da0-9307-a3b609b0a6ba
Error: unexpected response: 500 Internal Server Error
Usage:
  audiofile get [flags]
Flags:
  -h, --help        help for get
      --id string   audiofile id
unexpected response: 500 Internal Server Error%
```

Looks like an unexpected error has occurred, and we haven’t properly implemented how to handle this when searching for metadata for a file that has been deleted. We’ll need to modify the `services/metadata/handler_getbyid.go` file. At line 20, where we call the `GetById` method and handle the error, let’s return `200` instead of `500` after confirming the error is related to a folder not being found. It’s not necessarily an error that the user is searching for an ID that does not exist:

```markup
audio, err := m.Storage.GetByID(id)
if err != nil {
    if strings.Contains(err.Error(), "not found") ||     strings.Contains(err.Error(), "no such file or directory") {
        io.WriteString(res, "id not found")
        res.WriteHeader(200)
        return
    }
    res.WriteHeader(500)
    return
}
```

Let’s try it again:

```markup
./bin/audiofile get --id a5d9ab11-6f5f-4da0-9307-a3b609b0a6ba
id not found
```

That’s much better! Now let’s implement the search functionality. The implementation again is isolated to the `storage/flatfile.go` file where you will find the `Search` method:

```markup
func (f FlatFile) Search(searchFor string) ([]*models.Audio, error) {
    dirname, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
    audioFilePath := filepath.Join(dirname, "audiofile")
    matchingAudio := []*models.Audio{}
    err = filepath.WalkDir(audioFilePath, func(path string, 
          d fs.DirEntry, err error) error {
        if d.Name() == "metadata.json" {
            contents, err := os.ReadFile(path)
            if err != nil {
                return err
            }
            if strings.Contains(strings.
               ToLower(string(contents)), strings.
               ToLower(searchFor)) {
                data := models.Audio{}
                err = json.Unmarshal(contents, &data)
                if err != nil {
                    return err
                }
                matchingAudio = append(matchingAudio, &data)
            }
        }
        return nil
    })
    return matchingAudio, err
}
```

Like most of the methods existing in the storage, we start by getting the user’s home directory with the `os.UserHomeDir()` method and then, again, use `filepath.Join` to get the root `audiofile` path directory, which we will be walking. The `filepath.WalkDir` method is called starting at `audioFilePath`. We check each of the `metadata.json` files to see whether the `searchFor` string exists within the contents. The method returns a slice of `*models.Audio` and if the `searchFor` string is found within the contents, the audio is appended onto the slice that will be returned later.

Let’s give this a try with the following command and see that the expected metadata is returned:

```markup
./bin/audiofile search --value "Beat Doctor"
```

Now that we’ve created a few new commands to show how the `os` package and `path/filepath` packages can be used in a real-life example, let’s try to write some code that can run specifically on one operating system or another.

## Platform-specific code

Suppose your command-line application requires an external application that exists on the operating system, but the application required differs between operating systems. For the `audiofile` command-line application, suppose we want to create a command to play the audio file via the command line. Each operating system will need to use a different command to play the audio, as follows:

-   macOS: `afplay <filepath>`
-   Windows: `start <filepath>`
-   Linux: `aplay <filepath>`

Again, we use the Cobra-CLI to create the new `play` command. Let’s look at each different function that would need to be called for each operating system to play the audio file. First is the code for macOS:

```markup
func darwinPlay(audiofilePath string) {
    cmd := exec.Command("afplay", audiofilePath)
    if err := cmd.Start(); err != nil {
       panic(err)
    }
    fmt.Println("enjoy the music!")
    err := cmd.Wait()
    if err != nil {
       panic(err)
    }
}
```

We create a command to use the `afplay` executable and pass in the `audiofilePath`. Next is the code for Windows:

```markup
func windowsPlay(audiofilePath string) {
    cmd := exec.Command("cmd", "/C", "start", audiofilePath)
    if err := cmd.Start(); err != nil {
        return err
    }
    fmt.Println("enjoy the music!")
    err := cmd.Wait()
    if err != nil {
        return err
    }
}
```

This is a very similar function, except it uses the `start` executable in Windows to play the audio. Last is the code for Linux:

```markup
func linuxPlay(audiofilePath string) {
    cmd := exec.Command("aplay", audiofilePath)
    if err := cmd.Start(); err != nil {
        panic(err)
    }
    fmt.Println("enjoy the music!")
    err := cmd.Wait()
    if err != nil {
        panic(err)
    }
}
```

Again, the code is practically identical except for the application which is called to play the audio. In another case, this code could be more specific for the operating system, require different arguments, and even require a full path specific to the operating system. Regardless, we are ready to use these functions within the `play` command’s `RunE` field. The full `play` command is as follows:

```markup
var playCmd = &cobra.Command{
    Use: "play",
    Short: "Play audio file by id",
    RunE: func(cmd *cobra.Command, args []string) error {
        b, err := getAudioByID(cmd)
        if err != nil {
            return err
        }
        audio := models.Audio{}
        err = json.Unmarshal(b, &audio)
        if err != nil {
            return err
        }
        switch runtime.GOOS {
        case "darwin":
            darwinPlay(audio.Path)
            return nil
        case "windows":
            windowsPlay(audio.Path)
            return nil
        case "linux":
            linuxPlay(audio.Path)
            return nil
        default:
            fmt.Println(`Your operating system isn't supported 
                for playing music yet.
                Feel free to implement your additional use 
                case!`)
        }
        return nil
    },
}
```

The important part of this code is that we have created a switch case for the `runtime.GOOS` value, which tells us what operating system the application is running on. Depending on the operating system, a different method is called to start a process to play the audio file. Let’s recompile and try the play method with one of the stored audio file IDs:

```markup
./bin/audiofile play --id bf22c5c4-9761-4b47-aab0-47e93d1114c8
enjoy the music!
```

The final section of this chapter will show us how to implement this differently, if we’d like to, using build tags.

Just Imagine

# Build tags for targeted platforms

Built tags, or build constraints, can be used for many purposes, but in this section, we will be discussing how to use build tags to identify which files should be included in a package when building for specific operating systems. Build tags are given in a comment at the top of a file:

```markup
//go:build
```

Build tags are passed in as flags when running `go build`. There could be more than one tag on a file, and they follow on from the comment with the following syntax:

```markup
//go:build [tags]
```

Each tag is separated by a space. Suppose we want to indicate that this file will only be included in a build for the Darwin operating system, then we would add this to the top of the file:

```markup
//go:build darwin
```

Then when building the application, we would use something like this:

```markup
go build –tags darwin
```

This is just a super quick overview of how build tags can be used to constrain files specific to operating systems. Before we go into an implementation of this, let’s discuss the `build` package in a bit more detail.

## The build package

The `build` package gathers information about Go packages. In the _Chapter07_ code repository, there is a `buildChecks.go` file, which uses the `build` package to get information about the current package. Let’s see what information this code can give us:

```markup
func buildChecks() {
    ctx := build.Context{}
    p1, err := ctx.Import(".", ".", build.AllowBinary)
    if err != nil {
        fmt.Println("err: ", err)
    }
    fmt.Println("Dir:", p1.Dir)
    fmt.Println("Package name: ", p1.Name)
    fmt.Println("AllTags: ", p1.AllTags)
    fmt.Println("GoFiles: ", p1.GoFiles)
    fmt.Println("Imports: ", p1.Imports)
    fmt.Println("isCommand: ", p1.IsCommand())
    fmt.Println("IsLocalImport: ", build.IsLocalImport("."))
    fmt.Println(ctx)
}
```

We first create the `context` variable and then call the `Import` method. The `Import` method is defined in the documentation as follows:

```markup
func (ctxt *Context) Import(path string, srcDir string, mode ImportMode) (*Package, error)
```

It returns the details about the Go package named by the `path` and `srcDir` source directory parameters. In this case, the `main` package is returned from the package, then we can check all the variables and methods that exist to get more information on the package. Running this method locally will return something like this:

```markup
Dir: .
Package name:  main
AllTags:  [buildChecks]
GoFiles:  [checkRuntime.go environment.go file.go main.go process.go timer.go walking.go]
Imports:  [fmt io/fs os os/exec path/filepath runtime runtime/debug strings time]
isCommand/main package:  true
IsLocalImport:  true
```

Most of the values we are checking are self-explanatory. `AllTags` returns all tags that exist within the `main` package. `GoFiles` returns all the files included in the `main` package. `Imports` are all the unique imports that exist within the package. `IsCommand()` returns `true` if the package is considered a command to be installed, or if it is the main package. Finally, the `IsLocalImport` method checks whether an import file is local. This is a fun extra detail to interest you more about what the `build` package could potentially offer you.

## Build tags

Now that we have learned a little bit more about the `build` package, let’s use it for the main purpose of this chapter, building packages for specific operating systems. Build tags should be named intentionally, and since we are using them for a specific purpose, we can name each build tag by an operating system:

```markup
//go:build darwin
//go:build linux
//go:build windows
```

Let’s revisit the audio file code. Remember how in the `play` command, we check the `runtime` operating system and then call a specific method. Let’s rewrite this code using build tags.

### Example in the audio file

Let’s first simplify the command’s code to the following:

```markup
var playCmd = &cobra.Command{
    Use: "play",
    Short: "Play audio file by id",
    Long: `Play audio file by id`,
    RunE: func(cmd *cobra.Command, args []string) error {
        b, err := getAudioByID(cmd)
        if err != nil {
            return err
        }
        audio := models.Audio{}
        err = json.Unmarshal(b, &audio)
        if err != nil {
            return err
        }
        return play(audio.Path)
    },
}
```

Basically, we’ve simplified the code greatly by removing the operating system switch statement and the three functions that implement the play feature for each operating system. Instead, we’ve taken the code and created three new files: `play_darwin.go`, `play_windows.go`, and `play_linux.go`. Within each of these files is a build tag for each operating system. Let’s take the Darwin file, `play_darwin.go`, for example:

```markup
//go:build darwin
package cmd
import (
    "fmt"
    "os/exec"
)
func play(audiofilePath string) error {
    cmd := exec.Command("afplay", audiofilePath)
    if err := cmd.Start(); err != nil {
        return err
    }
    fmt.Println("enjoy the music!")
    err := cmd.Wait()
    if err != nil {
        return err
    }
    return nil
}
```

Notice that the `play` function has been renamed to match the function called in the `play` command in `play.go`. Since only one of the files gets included in the build, there’s no confusion as to which `play` function is called. We ensure that only one gets called within the `make` file, which is how we are currently running the application. In `Makefile`, I’ve designated a command to build specifically for Darwin:

```markup
build-darwin:
    go build -tags darwin -o bin/audiofile main.go
    chmod +x bin/audiofile
```

A Go file containing the `play` function is created for Windows and Linux. The specific tags for each operating system will similarly need to be passed into the `-tags` flag when building your application. In later chapters, we will discuss cross-compiling, which is the next step. But before we do, let’s leave this chapter by reviewing a list of OS-level differences to keep in mind while developing for multiple platforms.

## OS-level differences

Since you’ll be building your application for the main operating systems, it’s important to know the differences between them and know what to look out for. Let’s dive in with the following list:

-   **Filesystem**:
    -   Windows uses a different filesystem than Linux and Unix, so be mindful of the file paths when accessing files in your Go code.
    -   File paths in Windows use backslashes, (`\`), as directory separators, while Linux and Unix use forward slashes (`/`).
-   **Permissions**:
    -   Unix-like systems use file modes to manage permissions, where permissions are assigned to files and directories.
    -   Windows uses an **access control list** (**ACL**) to manage permissions, where permissions are assigned to specific users or groups for a file or directory in a more flexible and granular manner.
    -   In general, it’s a good practice to carefully consider user and group permissions when developing any command-line application, regardless of the operating system it will be running on.
-   **Executing commands**:
    -   The `exec` package in Go provides a convenient way to run commands in the same manner as in the terminal. However, it’s important to note that the command and its arguments must be passed in the correct format for each operating system.
    -   On Windows, you need to specify the file extension (for example, `.exe`, `.bat`, etc.) to run an executable file.
-   **Environmental variables**:
    -   Environmental variables can be used to configure your application, but their names and values may be different between Windows and Linux/Unix.
    -   On Windows, environmental variable names are case-insensitive, while on Linux/Unix, they are case-sensitive.
-   **Line endings**:
    -   Windows uses a different line ending character than Linux/Unix, so be careful when reading or writing files in your Go code. Windows uses a carriage return (`\r`) followed by a line feed (`\n`), while Linux/Unix uses only a line feed (`\n`).
-   **Signal handling**:
    -   In Unix systems, the `os/signal` package provides a way to handle signals sent to your application. However, this package is not supported on Windows.
    -   To handle signals in a cross-platform way, you can use the `os/exec` package instead.
-   **User input**:
    -   The way user input is read may also be different between Windows and Linux/Unix. On Windows, you may need to use the `os.Stdin` property, while on Linux/Unix you can use `os.Stdin` or the `bufio` package to read user input.
-   **Console colors**:
    -   On Windows, the console does not support ANSI escape codes for changing text color, so you will need to use a different approach for coloring text in the console.
    -   There are libraries available in Go, such as `go-colorable`, that provide a platform-independent way to handle console colors.
-   **Standard streams**:
    -   Standard streams, such as `os.Stdin`, `os.Stdout`, and `os.Stderr` may behave differently between Windows and Linux/Unix. It’s important to test your code on both platforms to make sure it works as expected.

These are some of the differences to be aware of when developing a command-line application in Go for different operating systems. It’s important to thoroughly test your application on each platform to ensure it behaves as expected.

Just Imagine

# Summary

The more operating systems your application supports, the more complicated it will get. Hopefully armed with the knowledge of some supportive packages for developing independently of the platform, you’ll feel confident that your application will run similarly across different operating systems. Also, by checking the `runtime` operating system and even separating code into separate operating system-specific files with build tags, you have at least a couple of options for defining how to organize your code. This chapter goes more in-depth than may be necessary, but hopefully, it inspires you.

Building for multiple operating systems will expand the usage of your command-line application. Not only can you reach Linux or Unix users but also Darwin and Windows users as well. If you want to grow your user base, then building an application to support more operating systems is an easy way to do so.

In the next chapter, [_Chapter 8_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_08.xhtml#_idTextAnchor166), _Building for Humans Versus Machines_, we’ll learn how to build a CLI that outputs according to who is receiving it: a machine or human. We’ll also learn how to structure the language for clarity and name commands for consistency with the rest of the CLIs in the community.

Just Imagine

# Questions

1.  What are the two different clocks that exist within an operating system? And does the `time.Time` struct in Go store one or the other clock, or both? Which should be used for calculating duration?
2.  Which package constant can be used to determine the `runtime` operating system?
3.  Where is the build tag comment set within a Go file – at the top, bottom, or above the defined function?

Just Imagine

# Answers

1.  The wall clock and monotonic clock. The `time.Time` struct stores both time values. The monotonic clock value should be used when calculating duration.
2.  `runtime.GOOS`
3.  At the top first line of the Go file.

Just Imagine

# Further reading

-   Visit the online documentation for the packages discussed at [https://pkg.go.dev/](https://pkg.go.dev/).