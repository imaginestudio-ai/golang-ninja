# Cross-Compilation across Different Platforms

This chapter introduces the user to cross-compilation, a powerful feature of Go, across different platforms. While build automation tools exist, understanding how to cross-compile provides essential knowledge for debugging and customization when necessary. This chapter will explain the different operating systems and architectures that Go can compile and how to determine which is needed. After Go is installed in your environment, there is a command, `go env`, with which you can see all the Go-related environment variables. We will discuss the two major ones used for building: `GOOS` and `GOARCH`.

We will give examples of how to build or install an application for each major operating system: Linux, macOS, and Windows. You will learn how to determine the Go operating system and architecture settings based on your environment and the available architectures for each major operating system.

This chapter ends with an example script to automate cross-compilation across the major operating systems and architectures. A script to run on the Darwin, Linux, or Windows environments is provided. In this chapter, we will cover the following topics in detail:

-   Manual compilation versus build automation tools
-   Using `GOOS` and `GOARCH`
-   Compiling for Linux, macOS, and Windows
-   Scripting to compile for multiple platforms

# Manual compilation versus build automation tools

In [_Chapter 14_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_14.xhtml#_idTextAnchor359), _Publishing Your Go Binary as a Homebrew Formula with GoReleaser_, we will delve into a fantastic open source tool, **GoReleaser**, which automates the process of building and releasing Go binaries. Despite its power and usefulness, it’s crucial to know how to manually compile your Go code. You see, not all projects can be built and released with GoReleaser. For instance, if your application requires unique build flags or dependencies, manual compilation may be necessary. Moreover, understanding how to manually compile your code is essential for addressing issues that may crop up during the build process. In essence, tools such as GoReleaser can make the process a lot smoother, but having a good grasp of the manual compile process is vital to ensure that your **command-line interface (CLI)** applications can be built and released in various scenarios.

Bookmark

# Using GOOS and GOARCH

When developing your command-line application, it is important to maximize the audience by developing for as many platforms as possible. However, you may also want to target just a particular set of operating systems and architectures. In the past, it was much more difficult to deploy to platforms that differed from the one you were developing on. In fact, developing on a macOS platform and deploying it on a Windows machine involved setting up a Windows build machine to build the binary. The tooling would have to be synchronized, and there would be other deliberations that made collaborative testing and distribution cumbersome.

Luckily, Golang has solved this by building support for multiple platforms directly into the language’s toolchain. As discussed in [_Chapter 7_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_07.xhtml#_idTextAnchor143), _Developing for Different Platforms_, and [_Chapter 11_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_11.xhtml#_idTextAnchor258), _Custom Builds and Testing CLI Commands_, we learned how to write platform-independent code and use the `go build` command and build tags to target specific operating systems and architectures. You may also use environment variables to target the operating system and architecture as well.

First, it’s good to know which operating systems and architectures are available for distribution. To find out, within your terminal, run the following command:

```markup
go tool dist list
```

The list is output in the following format: `GOOS`/`GOARCH`. `GOOS` is a local environment variable that defines the operating system to compile for and stands for **Go Operating System**. `GOARCH`, pronounced “gore-ch,” is a local environment variable that defines the architecture to compile for and stands for **Go Architecture**.

![Figure 12.1 – List of supported operating systems and architectures](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_12.1_B18883.jpg)

Figure 12.1 – List of supported operating systems and architectures

You can also call the preceding command with the `–json` flag to view more details. For example, for `linux/arm64`, you can see that it’s supported by `Cgo` from the `"CgoSupported"` field, but also that it is a first-class **port**, another word for a `GOOS/GOARCH` pair, indicated by the `"``FirstClass"` field:

```markup
{
"GOOS": "linux",
"GOARCH": "arm64",
"CgoSupported": true,
"FirstClass": true
},
```

A first-class port has the following properties:

-   Releases are blocked by broken builds
-   Official binaries are provided
-   Installation is documented

Next, determine your local operating system and architecture settings by running the following command within your terminal:

```markup
go env GOOS GOARCH
```

Currently, running this command on my macOS machine with an AMD64 architecture gives the following output:

```markup
darwin
amd64
```

The first environment variable, `GOOS`, is set to `darwin`, and the second environment variable, `GOARCH`, is set to `amd64`. We now know what `GOOS` and `GOARCH` are within the Go environment, the possible values, and also what values are set on your machine. Let’s learn how to use these environment variables.

You can use these two environment variables for compiling. Let’s generate a build to target the `darwin/amd64` port. You’ll do so by setting the `GOOS` or `GOARCH` environment variables and then running the `go build` command, or more specifically along with the `build` command:

```markup
GOOS=darwin GOARCH=amd64 go build
```

Let’s try this out with the audio file CLI and learn all the ways to compile for the three main operating systems: Linux, macOS, and Windows.

Bookmark

# Compiling for Linux, macOS, and Windows

There are several different ways to compile our command-line application for different operating systems and we’ll go over examples of each of these. First, you can compile by building or installing your application:

-   **Building** – Compiles the executable file and then moves it to the current folder or the filename indicated by the `–o` (output) flag
-   **Installing** – Compiles the executable file and then installs it to the `$GOPATH/bin` folder or `$GOBIN` if it is set and caches all non-main packages, which are imported to the `$``GOPATH/pkg` folder

## Building using tags

In our previous chapter, [_Chapter 11_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_11.xhtml#_idTextAnchor258), _Custom Builds and Testing CLI Commands_, we learned to build specifically for the macOS or Darwin operating system. To better understand how to use the `build` command, we run `go build –help` to see the usage:

```markup
mmontagnino@Marians-MacBook-Pro audiofile % go build -help
usage: go build [-o output] [build flags] [packages]
Run 'go help build' for details
```

Running `go help build` will reveal the build flags available. However, in these examples, we only use the `tags` flag. Within the `Makefile`, we already have the following commands:

```markup
build-darwin-free:
    go build -tags "darwin free" -o bin/audiofile main.go
    chmod +x bin/audiofile
build-darwin-pro:
    go build -tags "darwin pro" -o bin/audiofile main.go
    chmod +x bin/audiofile
build-darwin-pro-profile:
    go build -tags "darwin pro profile" -o bin/audiofile main.go
    chmod +x bin/audiofile
```

In these commands, we compile the application and output it to the `bin/audiofile` filename. To specify the Darwin operating system, we pass in the Darwin build tag to specify the files associated with the Darwin operating system. We’ll need to modify the output files to a folder that specifies Darwin, but also for other specifics such as the free versus the pro version since we’ll be building for other operating systems and levels. Let’s modify these.

### Building applications for a Darwin operating system using tags

The new `Makefile` commands to compile the application for the Darwin operating system are now as follows:

```markup
build-darwin-free:
    go build -tags "darwin free" -o builds/free/darwin/audiofile main.go
    chmod +x builds/free/darwin/audiofile
build-darwin-pro:
    go build -tags "darwin pro" -o builds/pro/darwin/audiofile main.go
    chmod +x builds/pro/darwin/audiofile
build-darwin-pro-profile:
    go build -tags "darwin pro profile" -o builds/profile/darwin/audiofile main.go
    chmod +x builds/profile/darwin/audiofile
```

We’ve swapped out the `bin/audiofile` output to something more specific. The free version for Darwin now outputs to `builds/free/darwin/audiofile`, the pro version outputs to `builds/pro/darwin/audiofile`, and the profile version outputs to `builds/profile/darwin/audiofile`. Let’s continue with the next operating system, Linux.

We can do the same for Linux and Windows, like so:

```markup
build-linux-free:
    go build -tags "linux free" -o builds/free/linux/audiofile main.go
    chmod +x builds/free/linux/audiofile
build-linux-pro:
   go build -tags "linux pro" -o builds/pro/linux/audiofile main.go
   chmod +x builds/pro/linux/audiofile
build-linux-pro-profile:
   go build -tags "linux pro profile" -o builds/profile/linux/audiofile main.go
   chmod +x builds/profile/linux/audiofile
build-windows-free:
    go build -tags "windows free" -o builds/free/windows/ audiofile.exe main.go
build-windows-pro:
    go build -tags "windows pro" -o builds/pro/windows/audiofile.exe main.go
build-windows-pro-profile:
    go build -tags "windows pro profile" -o builds/profile/windows/audiofile.exe main.go
```

The free Windows version is output to `builds/free/windows/audiofile.exe`, the pro Windows version is output to `builds/pro/windows/audiofile.exe`, and the Windows profile version is output to `builds/profile/windows/audiofile.exe`. Now, suppose we don’t want to run each of the individual commands one by one, as there are so many to run! We can write a command to build all versions using tags.

### Building applications for all operating systems using tags

Let’s add a new `Makefile` command to build all the operating systems. Basically, we write one command that calls all other commands:

```markup
build-all: build-darwin-free build-darwin-pro build-darwin-pro-profile build-linux-free build-linux-pro build-linux-pro-profile build-windows-free build-windows-pro build-windows-pro-profile
```

Let’s try running this command via the terminal:

```markup
make build-all
```

If you’re running on Darwin, you’ll see the following output:

```markup
mmontagnino@Marians-MacBook-Pro audiofile % make build-all
go build -tags "darwin free" -o builds/free/darwin/audiofile main.go
chmod +x builds/free/darwin/audiofile
go build -tags "darwin pro" -o builds/pro/darwin/audiofile main.go
chmod +x builds/pro/darwin/audiofile
go build -tags "darwin pro profile" -o builds/profile/darwin/audiofile main.go
chmod +x builds/profile/darwin/audiofile
go build -tags "linux free" -o builds/free/linux/audiofile main.go
# internal/goos
/usr/local/go/src/internal/goos/zgoos_linux.go:7:7: GOOS redeclared in this block
        /usr/local/go/src/internal/goos/zgoos_darwin.go:7:7: other declaration of GOOS
/usr/local/go/src/internal/goos/zgoos_linux.go:9:7: IsAix redeclared in this block
        /usr/local/go/src/internal/goos/zgoos_darwin.go:9:7: other declaration of IsAix
/usr/local/go/src/internal/goos/zgoos_linux.go:10:7: IsAndroid redeclared in this block
...
/usr/local/go/src/internal/goos/zgoos_linux.go:17:7: too many errors
make: *** [build-linux-free] Error 2
```

I’ve removed part of the error message; however, the most important message is `GOOS redeclared in this block`. This error message comes up when the operating system is set but conflicts with the `GOOS` environment variable. For example, the command that failed used the operating build tag to specify a Linux build:

```markup
go build -tags "linux free" -o builds/free/linux/audiofile main.go
```

However, running `go env | grep GOOS` in my macOS terminal shows the value of the `GOOS` environment variable:

```markup
GOOS="darwin"
```

Let’s modify the build commands to set the `GOOS` environment variable so it matches the output type based on the build tag.

## Building using the GOOS environment variable

The Linux builds have been modified to set the `GOOS` environment variable to Linux by prepending `GOOS=linux` before the `build` command:

```markup
build-linux-free:
    GOOS=linux go build -tags "linux free" -o builds/free/linux/audiofile main.go
    chmod +x builds/free/linux/audiofile
build-linux-pro:
    GOOS=linux go build -tags "linux pro" -o builds/pro/linux/audiofile main.go
    chmod +x builds/pro/linux/audiofile
build-linux-pro-profile:
    GOOS=linux go build -tags "linux pro profile" -o builds/profile/linux/audiofile main.go
    chmod +x builds/profile/linux/audiofile
```

The Windows builds have been modified to set the `GOOS` environment variable to Windows by prepending `GOOS=windows` before the `build` command:

```markup
build-windows-free:
    GOOS=windows go build -tags "windows free" -o builds/free/windows/audiofile.exe main.go
build-windows-pro:
    GOOS=windows go build -tags "windows pro" -o builds/pro/windows/audiofile.exe main.go
build-windows-pro-profile:
    GOOS=windows go build -tags "windows pro profile" -o builds/profile/windows/audiofile.exe main.go
```

Now, let’s try the `build-all` command again. It runs successfully and we can see all the files generated by the `build` command by running `find –type –f ./builds` in the repo:

```markup
mmontagnino@Marians-MacBook-Pro audiofile % find ./builds -type f
./builds/pro/linux/audiofile
./builds/pro/darwin/audiofile
./builds/pro/windows/audiofile.exe
./builds/free/linux/audiofile
./builds/free/darwin/audiofile
./builds/free/windows/audiofile.exe
./builds/profile/linux/audiofile
./builds/profile/darwin/audiofile
./builds/profile/windows/audiofile.exe
```

## Building using the GOARCH environment variable

Many different possible architecture values can be associated with a single operating system. Rather than creating a command for each, we’ll start with just one example:

```markup
build-darwin-amd64-free:
    GOOS=darwin GOARCH=amd64 go build -tags "darwin free" -o builds/free/darwin/audiofile main.go
    chmod +x builds/free/darwin/audiofile
```

This example specifies the operating system, the `GOOS` environment variable, as `darwin`, and then the architecture, the `GOARCH` environment variable, as `amd64`.

There’d be too many commands to create if we were to create a `build` command for each architecture of each major operating system. We’ll save this for a script within the last section of this chapter.

## Installing using tags and GOOS env va

-   As mentioned earlier, another way to compile your command-line application is by installing it. The `install` command compiles the application, like the `go build` command, but also with the additional step of moving the compiled application to the `$GOPATH/bin` folder or `$GOBIN` value. To learn more about the `install` command, we run the following `go install –``help` command:

```markup
mmontagnino@Marians-MacBook-Pro audiofile % go install -help
usage: go install [build flags] [packages]
Run 'go help install' for details
```

-   The same flags for building are available for installing. Again, we will use the `tags` flag only. Let’s first run the `install` command on the macOS system:

```markup
go install -tags "darwin pro" github.com/marianina8/audiofile
```

However, running `go env | grep GOPATH` in my macOS terminal shows the value of the `GOOS` environment variable:

```markup
mmontagnino@Marians-MacBook-Pro audiofile % go env | grep GOPATH
GOPATH="/Users/mmontagnino/Code"
```

Confirm that the audio file CLI executable exists in the `$GOPATH/bin` or `/``Users/mmontagnino/Code/bin` folder.

As mentioned, we can use build tags to separate builds based on the operating system and architecture. Within the audio file repository, we’re already doing so with the following files associated with the `play` and `bug` commands. For the `bug` command, we have the following files. Now, let’s add some `install` commands within the `Makefile` now that we understand how to use the build tags and `GOOS` environment variables.

### install commands for the Darwin operating system

The `install` commands for the Darwin operating system include passing in the specific tags, including `darwin`, and the levels, defined by tags, to install:

```markup
install-darwin-free:
    go install -tags "darwin free" github.com/marianina8/audiofile
install-darwin-pro:
    go install -tags "darwin pro" github.com/marianina8/audiofile
install-darwin-pro-profile:
    go install -tags "darwin pro profile" github.com/marianina8/audiofile
```

### install commands for the Linux operating system

The `install` commands for the Linux operating system include passing in the specific tags, including `linux`, and the package to install. To ensure the commands do not error out with conflicting `GOOS` settings, we set the matching environment variable, `GOOS`, to `linux`:

```markup
install-linux-free:
    GOOS=linux go install -tags "linux free" github.com/marianina8/audiofile
install-linux-pro:
    GOOS=linux go install -tags "linux pro" github.com/marianina8/audiofile
install-linux-pro-profile:
    GOOS=linux go install -tags "linux pro profile" github.com/marianina8/audiofile
```

### install commands for the Windows operating system

The `install` commands for the Windows operating system include passing in the specific tags, including `windows`, and the package to install. To ensure the commands do not error out with conflicting `GOOS` settings, we set the matching environment variable, `GOOS`, to `windows`:

```markup
install-windows-free:
    GOOS=windows go install -tags "windows free" github.com/marianina8/audiofile
install-windows-pro:
    GOOS=windows go install -tags "windows pro" github.com/marianina8/audiofile
install-windows-pro-profile:
    GOOS=windows go install -tags "windows pro profile" github.com/marianina8/audiofile
```

Remember that for your `Makefile`, you’ll need to change the location of the package if you have forked the repo under your own account. Run the `make` command for the operating system you need and confirm that the application is installed by checking the `$GOPATH/bin` or `$``GOBIN` folder.

## Installing using tags and GOARCH env var

While many different possible architecture values can be associated with a single operating system, let’s start with just one example of installing with `GOARCH` `env var`:

```markup
install-linux-amd64-free:
    GOOS=linux GOARCH=amd64 go install -tags "linux free" github.com/marianina8/audiofile
```

This example specifies the operating system, the `GOOS` environment variable, as `linux`, and then the architecture, the `GOARCH` environment variable, as `amd64`. Rather than creating a command for each pair of operating systems and architectures, again, we’ll save this for a script within the last section of this chapter.

Bookmark

# Scripting to compile for multiple platforms

We’ve learned several different ways to compile for operating systems using the `GOOS` and `GOARCH` environment variables and using build tags. The `Makefile` can fill up rather quickly with all the different combinations of `GOOS`/`GOARCH` pairs and scripting may provide a better solution if you want to generate builds for many more specific architectures.

## Creating a bash script to compile in Darwin or Linux

Let’s start by creating a bash script. Let’s name it `build.sh`. To create the file, I simply type the following:

```markup
touch build.sh
```

The preceding command creates the file when it does not exist. The file extension is `.sh`, which, while unnecessary to add, clearly indicates that the file is a bash script type. Next, we want to edit it. If using `vi`, use the following command:

```markup
vi build.sh
```

Otherwise, edit the file using the editor of your choice.

### Adding the shebang

The first line of a bash script is called the **shebang**. It is a character sequence that indicates the program loader’s first instruction. It defines which interpreter to run when reading, or interpreting, the script. The first line to indicate to use the bash interpreter is as follows:

```markup
#!/bin/bash
```

The shebang consists of a couple of elements:

-   `#!` instructs the program loader to load an interpreter for the code
-   `/bin/bash` indicates the bash or interpreter’s location

These are some typical shebangs for different interpreters:

<table id="table001-3" class="No-Table-Style"><colgroup><col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Interpreter</strong></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Shebang</strong></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">Bash</span></p></td><td class="No-Table-Style"><p><code class="literal">#!/</code><span class="No-Break"><code class="literal">bin/bash</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">Bourne shell</span></p></td><td class="No-Table-Style"><p><code class="literal">#!/</code><span class="No-Break"><code class="literal">bin/sh</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">Powershell</span></p></td><td class="No-Table-Style"><p><code class="literal">#!/</code><span class="No-Break"><code class="literal">user/bin/pwsh</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p>Other <span class="No-Break">scripting languages</span></p></td><td class="No-Table-Style"><p><code class="literal">#!/</code><span class="No-Break"><code class="literal">user/bin/env &lt;interpreter&gt;</code></span></p></td></tr></tbody></table>

Table 12.1 – Shebang lines for different interpreters

### Adding comments

To add comments to your bash script, simply start the comment with the `#` symbol and the pound sign, followed by comment text. This text can be used by you and other developers to document information that might not be easily understood from the code alone. It could also just add some details on the usage of the script, who the author is, and so on.

### Adding print lines

In a bash file, to print lines out, simply use the `echo` command. These print lines will help you to understand exactly where your application is within its running process. Use these lines with intention and they will give you and your users some useful insight that can even make debugging easier.

### Adding code

Within the bash script, we’ll generate builds for all the differing build tags for each operating system and architecture pair. Let’s first start to see which architecture values are available for Darwin:

```markup
go tool dist list | grep darwin
```

The values returned are as follows:

```markup
darwin/amd64
darwin/arm64
```

Let’s generate the different Darwin builds – free, pro, and profile versions – for all architectures with the following code:

```markup
# Generate darwin builds
darwin_archs=(amd64 arm64)
for darwin_arch in ${darwin_archs[@]}
do
    echo "building for darwin/${darwin_arch} free version..."
    env GOOS=darwin GOARCH=${darwin_arch} go build -tags free -o builds/free/darwin/${darwin_arch}/audiofile main.go
    echo "building for darwin/${darwin_arch} pro version..."
    env GOOS=darwin GOARCH=${darwin_arch} go build -tags pro -o builds/pro/darwin/${darwin_arch}/audiofile main.go
    echo "building for darwin/${darwin_arch} profile version..."
    env GOOS=darwin GOARCH=${darwin_arch} go build -tags profile -o builds/profile/darwin/${darwin_arch}/audiofile main.go
done
```

Next, let’s do the same with Linux, first grabbing the architecture values available:

```markup
go tool dist list | grep linux
```

The values returned are as follows:

```markup
linux/386        linux/mips64le
linux/amd64    linux/mipsle
linux/arm        linux/ppc64
linux/arm64    linux/ppc64le
linux/loong64    linux/riscv64
linux/mips        linux/s390x
linux/mips64
```

Let’s generate the different Linux builds – the free, pro, and profile versions – for all architectures with the following code:

```markup
# Generate linux builds
linux_archs=(386 amd64 arm arm64 loong64 mips mips64 mips64le mipsle ppc64 ppc64le riscv64 s390x)
for linux_arch in ${linux_archs[@]}
do
    echo "building for linux/${linux_arch} free version..."
    env GOOS=linux GOARCH=${linux_arch} go build -tags free -o builds/free/linux/${linux_arch}/audiofile main.go
    echo "building for linux/${linux_arch} pro version..."
    env GOOS=linux GOARCH=${linux_arch} go build -tags pro -o builds/pro/linux/${linux_arch}/audiofile main.go
    echo "building for linux/${linux_arch} profile version..."
    env GOOS=linux GOARCH=${linux_arch} go build -tags profile -o builds/profile/linux/${linux_arch}/audiofile main.go
done
```

Next, let’s do the same with Windows, first grabbing the architecture values available:

```markup
go tool dist list | grep windows
```

The values returned are as follows:

```markup
windows/386
windows/amd64
windows/arm
windows/arm64
```

Finally, let’s generate the different Windows builds – the free, pro, and profile versions – for all architectures with the following code:

```markup
# Generate windows builds
windows_archs=(386 amd64 arm arm64)
for windows_arch in ${windows_archs[@]}
do
    echo "building for windows/${windows_arch} free version..."
    env GOOS=windows GOARCH=${windows_arch} go build -tags free -o builds/free/windows/${windows_arch}/audiofile.exe main.go
    echo "building for windows/${windows_arch} pro version..."
    env GOOS=windows GOARCH=${windows_arch} go build -tags pro -o builds/pro/windows/${windows_arch}/audiofile.exe main.go
    echo "building for windows/${windows_arch} profile version..."
    env GOOS=windows GOARCH=${windows_arch} go build -tags profile -o builds/profile/windows/${windows_arch}/audiofile.exe main.go
done
```

Here’s the code when run from the Darwin/macOS or Linux terminal:

```markup
./build.sh
```

We can check that the executable files have been generated. The full list is quite long, and they have been organized within the following nested folder structure:

```markup
/builds/{level}/{operating-system}/{architecture}/{audiofile-executable}
```

![Figure 12.2 – Screenshot of generated folders from the build bash script](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_12.2_B18883.jpg)

Figure 12.2 – Screenshot of generated folders from the build bash script

A script to generate these builds will need to be different if run on Windows, for example. If you are running your application on Darwin or Linux, try running the build script and see the generated builds populate. You can now share these builds with other users running on a different platform. Next, we’ll create a PowerShell script to generate the same builds to run in Windows.

## Creating a PowerShell script in Windows

Let’s start by creating a PowerShell script. Let’s name it `build.ps1`. Create the file by typing the following command within PowerShell:

```markup
notepad build.ps1
```

The preceding command asks to create the file when it does not exist. The file extension is `.ps1`, which indicates that the file is a PowerShell script type. Next, we want to edit it. You may use Notepad or another editor of your choice.

Unlike a bash script, a PowerShell script does not require a shebang. To learn more about how to write a PowerShell script, you can review the documentation here: [https://learn.microsoft.com/en-us/powershell/](https://learn.microsoft.com/en-us/powershell/).

### Adding comments

To add comments to your PowerShell script, simply start the comment with a `#` symbol and a pound sign, followed by comment text.

### Adding print lines

In a PowerShell file, to print lines out, simply use the `Write-Output` command:

```markup
Write-Output "building for windows/amd64..."
```

Writing output will help you to understand exactly where your application is within its running process, make it easier to debug, and give the user a sense that something is running. Having no output at all is not only boring but also communicates nothing to the user.

### Adding code

Within the PowerShell script, we’ll generate builds for all the differing build tags for each operating system and architecture pair. Let’s start by seeing which architecture values are available for Darwin via a Windows command:

```markup
PS C:\Users\mmontagnino\Code\src\github.com\marianina8\audiofile> go tool dist list | Select-String darwin
```

Using the `Select-String` command, we can return only the values that contain `darwin`. These values are returned:

```markup
darwin/amd64
darwin/arm64
```

We can run a similar command for Linux:

```markup
PS C:\Users\mmontagnino\Code\src\github.com\marianina8\audiofile> go tool dist list | Select-String linux
```

And a command for Windows:

```markup
PS C:\Users\mmontagnino\Code\src\github.com\marianina8\audiofile> go tool dist list | Select-String windows
```

The same values are returned within the previous sections, so I won’t print them out. However, now that we know how to get the architecture for each operating system, we can add the code to generate the builds for all of them.

The code to generate Darwin builds is as follows:

```markup
# Generate darwin builds
$darwin_archs="amd64","arm64"
foreach ($darwin_arch in $darwin_archs)
{
    Write-Output "building for darwin/$($darwin_arch) free version..."
    $env:GOOS="darwin";$env:GOARCH=$darwin_arch; go build -tags free -o .\builds\free\darwin\$darwin_arch\audiofile main.go
    Write-Output "building for darwin/$($darwin_arch) pro version..."
    $env:GOOS="darwin";$env:GOARCH=$darwin_arch; go build -tags pro -o .\builds\pro\darwin\$darwin_arch\audiofile main.go
    Write-Output "building for darwin/$($darwin_arch) profile version..."
    $env:GOOS="darwin";$env:GOARCH=$darwin_arch; go build -tags profile -o .\builds\profile\darwin\$darwin_arch\audiofile main.go
}
```

The code to generate Linux builds is as follows:

```markup
# Generate linux builds
$linux_archs="386","amd64","arm","arm64","loong64","mips","mips64","mips64le","mipsle","ppc64","ppc64le","riscv64","s390x"
foreach ($linux_arch in $linux_archs)
{
    Write-Output "building for linux/$($linux_arch) free version..."
    $env:GOOS="linux";$env:GOARCH=$linux_arch; go build -tags free -o .\builds\free\linux\$linux_arch\audiofile main.go
    Write-Output "building for linux/$($linux_arch) pro version..."
    $env:GOOS="linux";$env:GOARCH=$linux_arch; go build -tags pro -o .\builds\pro\linux\$linux_arch\audiofile main.go
    Write-Output "building for linux/$($linux_arch) profile version..."
    $env:GOOS="linux";$env:GOARCH=$linux_arch; go build -tags profile -o .\builds\profile\linux\$linux_arch\audiofile main.go
}
```

Finally, the code to generate Windows builds is as follows:

```markup
# Generate windows builds
$windows_archs="386","amd64","arm","arm64"
foreach ($windows_arch in $windows_archs)
{
    Write-Output "building for windows/$($windows_arch) free version..."
    $env:GOOS="windows";$env:GOARCH=$windows_arch; go build -tags free -o .\builds\free\windows\$windows_arch\audiofile.exe main.go
    Write-Output "building for windows/$($windows_arch) pro version..."
    $env:GOOS="windows";$env:GOARCH=$windows_arch; go build -tags pro -o .\builds\pro\windows\$windows_arch\audiofile.exe main.go
    Write-Output "building for windows/$($windows_arch) profile version..."
    $env:GOOS="windows";$env:GOARCH=$windows_arch; go build -tags profile -o .\builds\profile\windows\$windows_arch\audiofile.exe main.go
}
```

Each section generates a build for one of the three major operating systems and all the available architectures. To run the script from PowerShell, just run the following script:

```markup
./build.ps1
```

The following will be the output for each port:

```markup
building for $GOOS/$GOARCH [free/pro/profile] version...
```

Check the `builds` folder to see all the ports generated successfully. The full list is quite long, and they have been organized within the following nested folder structure:

```markup
/builds/{level}/{operating-system}/{architecture}/{audiofile-executable}
```

Now, we can generate builds for all operating systems and architectures from a PowerShell script, which can be run on Windows. If you run any of the major operating systems – Darwin, Linux, or Windows – you can now generate a build for your own platform or anyone else who would like to use your application.

Bookmark

# Summary

In this chapter, you learned what the `GOOS` and `GOARCH` environment variables are and how you can use them, as well as build tags, to customize builds based on the operating system, architecture, and levels. These environment variables help you to learn more about the environment you’re building in and possibly understand why a build may have trouble executing on another platform.

There are also two ways to compile an application – building or installing. In this chapter, we discussed how to build or install the application and what the difference is. The same flags are available for each command, but we discussed how to build or install on each of the major operating systems using the `Makefile`. However, this also showed how large the `Makefile` can become!

Finally, we learned how to create a simple script to run in Darwin, Linux, or Windows to generate all the builds needed for all the major operating systems. You learned how to write both a bash and PowerShell script to generate builds. In the next chapter, [_Chapter 13_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_13.xhtml#_idTextAnchor331), _Using Containers for Distribution_, we will learn how to run these compiled applications on containers made from different operating system images. Finally, in [_Chapter 14_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_14.xhtml#_idTextAnchor359), _Publishing Your Go Binary as a Homebrew Formula with GoReleaser_, you’ll explore the tools required to automate the process of building and releasing your Go binaries across a range of operating systems and architectures. By learning how to use GoReleaser, you can significantly accelerate the process of releasing and deploying your application. This way, you can concentrate on developing new features and addressing bugs instead of getting bogged down with the build and compile process. Ultimately, using GoReleaser can save you valuable time and energy that you can use to make your application even better.

Bookmark

# Questions

1.  What Go environment variables define the operating system and the architecture?
2.  What additional security do you get from building with a first-class port?
3.  What command would you run on Linux to find the port values for the Darwin operating system?

Bookmark

# Answers

1.  `GOOS` is the Golang operating system, and `GOARCH` is the Golang architecture value.
2.  There are several reasons why a first-class port is more secure: releases are blocked by broken builds, official binaries are provided, and installation is documented.
3.  `go tool dist list |` `grep darwin`.

Bookmark

# Further reading

-   Read more about compiling at [https://go.dev/doc/tutorial/compile-install](https://go.dev/doc/tutorial/compile-install)[](https://go.dev/doc/tutorial/compile-install%0A)
-   Read more about Go environment variables at [https://pkg.go.dev/cmd/go](https://pkg.go.dev/cmd/go)