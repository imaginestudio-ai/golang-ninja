
# A Quick Introduction to Go

ImagineDevOps  that you are a developer and you want to create a command-line utility. Similarly, imagine that you have a REST API and you want to create a RESTful server that implements that REST API. The first thought that will come to your mind will most likely be which programming language to use.

The most common answer to this question is to use the programming language you know best. However, this book is here to help you consider Go for all these and many more tasks and projects. In this chapter, we begin by explaining what Go is, and continue with the history of Go, and how to run Go code. We will explain some core characteristics of Go, such as how to define variables, control the flow of your programs, and get user input, and we will apply some of these concepts by creating a command-line phone book application.

We will cover the following topics:

-   Introducing Go
-   Hello World!
-   Running Go code
-   Important characteristics of Go
-   Developing the `which(1)` utility in Go
-   Logging information
-   Overview of Go generics
-   Developing a basic phone book application

Bookmark

# Introducing Go

Go is an open-source systems programming language initially developed as an internal Google project that went public back in 2009. The spiritual fathers of Go are Robert Griesemer, Ken Thomson, and Rob Pike.

Although the official name of the language is Go, it is sometimes (wrongly) referred to as _Golang_. The official reason for this is that [go.org](http://go.org) was not available for registration and [golang.org](http://golang.org) was chosen instead. The practical reason for this is that when you are querying a search engine for Go-related information, the word _Go_ is usually interpreted as a verb. Additionally, the official Twitter hashtag for Go is _#golang_.

Although Go is a general-purpose programming language, it is primarily used for writing system tools, command-line utilities, web services, and software that work over networks. Go can also be used for teaching programming and is a good candidate as your first programming language because of its lack of verbosity and clear ideas and principles. Go can help you develop the following kinds of applications:

-   Professional web services
-   Networking tools and servers such as Kubernetes and Istio
-   Backend systems
-   System utilities
-   Powerful command-line utilities such as `docker` and `hugo`
-   Applications that exchange data in JSON format
-   Applications that process data from relational databases, NoSQL databases, or other popular data storage systems
-   Compilers and interpreters for programming languages you design
-   Database systems such as CockroachDB and key/value stores such as etcd

There are many things that Go does better than other programming languages, including the following:

-   The default behavior of the Go compiler can catch a large set of silly errors that might result in bugs.
-   Go uses fewer parentheses than C, C++, or Java, and no semicolons, which makes the look of Go source code more human-readable and less error-prone.
-   Go comes with a rich and reliable standard library.
-   Go has support for concurrency out of the box through goroutines and channels.
-   Goroutines are really lightweight. You can easily run thousands of goroutines on any modern machine without any performance issues.
-   Unlike C, Go supports functional programming.
-   Go code is backward compatible, which means that newer versions of the Go compiler accept programs that were created using a previous version of the language without any modifications. This compatibility guarantee is limited to major versions of Go. For example, there is no guarantee that a Go 1.x program will compile with Go 2.x.

Now that we know what Go can do and what Go is good at, let's discuss the history of Go.

## The history of Go

As mentioned earlier, Go started as an internal Google project that went public back in 2009. Griesemer, Thomson, and Pike designed Go as a language for professional programmers who want to build reliable, robust, and efficient software that is easy to manage. They designed Go with simplicity in mind, even if simplicity meant that Go was not going to be a programming language for everyone.

The figure that follows shows the programming languages that directly or indirectly influenced Go. As an example, Go syntax looks like C whereas the package concept was inspired by Modula-2.

![Graphical user interface, application
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_01_01.png)

Figure 1.1: The programming languages that influenced Go

The deliverable was a programming language, its tools, and its standard library. What you get with Go, apart from its syntax and tools, is a pretty rich standard library and a type system that tries to save you from easy mistakes such as implicit type conversions, unused variables and unused packages. The Go compiler catches most of these easy mistakes and refuses to compile until you do something about them. Additionally, the Go compiler can find difficult to catch mistakes such as race conditions.

If you are going to install Go for the first time, you can start by visiting [https://golang.org/dl/](https://golang.org/dl/). However, there is a big chance that your UNIX variant has a ready-to-install package for the Go programming language, so you might want to get Go by using your favorite package manager.

## Why UNIX and not Windows?

You might ask why we're talking about UNIX all the time and not discussing Microsoft Windows as well. There are two main reasons for this. The first reason is that most Go programs will work on Windows machines without any code changes because **Go is portable by design**—this means that you should not worry about the operating system you are using.

However, you might need to make small changes to the code of some system utilities for them to work in Windows. Additionally, there are still going to be some libraries that only work on Windows machines and some that only work on non-Windows machines. The second reason is that many services written in Go are executed in a Docker environment—**Docker images use the Linux operating system**, which means that you should program your utilities having the Linux operating system in mind.

As far as user experience is concerned, UNIX and Linux are very similar. The main difference is that Linux is open-source software whereas UNIX is proprietary software.

## The advantages of Go

Go comes with some important advantages for developers, starting with the fact that it was designed and is currently maintained by real programmers. Go is also easy to learn, especially if you are already familiar with programming languages such as C, Python, or Java. On top of that, Go code is good-looking, at least to me, which is great, especially when you are programming applications for a living and you have to look at code on a daily basis. Go code is also easy to read and offers support for Unicode out of the box, which means that you can make changes to existing Go code easily. Lastly, Go has reserved only 25 keywords, which makes it much easier to remember the language. Can you do that with C++?

Go also comes with concurrency capabilities using a simple concurrency model that is implemented using **goroutines** and **channels**. Go manages OS threads for you and has a powerful runtime that allows you to spawn lightweight units of work (_goroutines_) that communicate with each other using _channels_. Although Go comes with a rich standard library, there are really handy Go packages such as `cobra` and `viper` that allow Go to develop complex command-line utilities such as `docker` and `hugo`. This is greatly supported by the fact that Go's executable binaries are _statically linked_, which means that once they are generated, they do not depend on any shared libraries and include all required information.

Due to its simplicity, Go code is predictable and does not have strange side effects, and although Go supports **pointers**, it does not support pointer arithmetic like C, unless you use the `unsafe` package, which is the root of many bugs and security holes. Although Go is not an object-oriented programming language, Go interfaces are very versatile and allow you to mimic some of the capabilities of object-oriented languages such as **polymorphism**, **encapsulation**, and **composition**.

Additionally, the latest Go versions offer support for **generics**, which simplifies your code when working with multiple data types. Last but not least, Go comes with support for **garbage collection**, which means that no manual memory management is needed.

Although Go is a very practical and competent programming language, it is not perfect:

-   Although this is a personal preference rather than an actual technical shortcoming, Go has no direct support for object-oriented programming, which is a popular programming paradigm.
-   Although goroutines are lightweight, they are not as powerful as OS threads. Depending on the application you are trying to implement, there might exist some rare cases where goroutines will not be appropriate for the job. However, in most cases, designing your application with goroutines and channels in mind will solve your problems.
-   Although garbage collection is fast enough most of the time and for almost all kinds of applications, there are times when you need to handle memory allocation manually—Go cannot do that. In practice, this means that Go will not allow you to perform any memory management manually.

However, there are many cases where you can choose Go, including the following:

-   Creating complex command-line utilities with multiple commands, sub-commands, and command-line parameters
-   Building highly concurrent applications
-   Developing servers that work with APIs and clients that interact by exchanging data in myriad formats including JSON, XML, and CSV
-   Developing WebSocket servers and clients
-   Developing gRCP servers and clients
-   Developing robust UNIX and Windows system tools
-   Learning programming

In the following sections, we will cover a number of concepts and utilities in order to build a solid foundation of knowledge, before building a simplified version of the `which(1)` utility. At the end of the chapter, we'll develop a naive phone book application that will keep evolving as we explain more Go features in the chapters that follow.

But first, we'll present the `go doc` command, which allows you to find information about the Go standard library, its packages, and their functions. Then, we'll show how to execute Go code using the `Hello World!` program as an example.

## The go doc and godoc utilities

The Go distribution comes with a plethora of tools that can make your life as a programmer easier. Two of these tools are the `go doc` subcommand and `godoc` utility, which allow you to see the documentation of existing Go functions and packages without needing an internet connection. However, if you prefer viewing the Go documentation online, you can visit [https://pkg.go.dev/](https://pkg.go.dev/). As `godoc` is not installed by default, you might need to install it by running `go install golang.org/x/tools/cmd/godoc@latest`.

The `go doc` command can be executed as a normal command-line application that displays its output on a terminal, and `godoc` as a command-line application that starts a web server. In the latter case, you need a web browser to look at the Go documentation. The first utility is similar to the UNIX `man(1)` command, but for Go functions and packages.

The number after the name of a UNIX program or system call refers to the section of the manual a manual page belongs to. Although most of the names can be found only once in the manual pages, which means that putting the section number is not required, there are names that can be located in multiple sections because they have multiple meanings, such as `crontab(1)` and `crontab(5)`. Therefore, if you try to retrieve the manual page of a name with multiple meanings without stating its section number, you will get the entry that has the smallest section number.

So, in order to find information about the `Printf()` function of the `fmt` package, you should execute the following command:

```markup
$ go doc fmt.Printf
```

Similarly, you can find information about the entire `fmt` package by running the following command:

```markup
$ go doc fmt
```

The second utility requires executing `godoc` with the `-http` parameter:

```markup
$ godoc -http=:8001
```

The numeric value in the preceding command, which in this case is `8001`, is the port number the HTTP server will listen to. As we have omitted the IP address, `godoc` is going to listen to all network interfaces.

You can choose any port number that is available provided that you have the right privileges. However, note that port numbers `0`\-`1023` are restricted and can only be used by the root user, so it is better to avoid choosing one of those and pick something else, provided that it is not already in use by a different process.

You can omit the equals sign in the presented command and put a space character in its place. So, the following command is completely equivalent to the previous one:

```markup
$ godoc -http :8001
```

After that, you should point your web browser to the `http://localhost:8001/` URL in order to get the list of available Go packages and browse their documentation. If you are using Go for the first time, you will find the Go documentation very handy for learning the parameters and the return values of the functions you want to use—as you progress in your Go journey, you will use the Go documentation for learning the gory details of the functions and variables that you want to use.


# Hello World!

The following is the Go version of the Hello World program. Please type it and save it as `hw.go`:

```markup
package main
import (
    "fmt"
)
func main() {
    fmt.Println("Hello World!")
}
```

Each Go source code begins with a `package` declaration. In this case, the name of the package is `main`, which has a special meaning in Go. The `import` keyword allows you to include functionality from existing packages. In our case, we only need some of the functionality of the `fmt` package that belongs to the standard Go library. Packages that are not part of the standard Go library are imported using their full internet path. The next important thing if you are creating an executable application is a `main()` function. Go considers this the entry point to the application and begins the execution of the application with the code found in the `main()` function of the `main` package.

`hw.go` is a Go program that runs on its own. Two characteristics make `hw.go` an autonomous source file that can generate an executable binary: the name of the package, which should be `main`, and the presence of the `main()` function—we discuss Go functions in more detail in the next subsection but we will learn even more about functions and methods, which are functions attached to specific data types, in _Chapter 5_, _Go Packages and Functions_.

## Introducing functions

Each Go function definition begins with the `func` keyword followed by its name, signature and implementation. As happens with the `main` package, you can name your functions anything you want—there is a global Go rule that also applies to function and variable names and is valid for all packages except `main`: **everything that begins with a lowercase letter is considered private and is accessible in the current package only**. 

The statement "Each Go function definition begins with the func keyword followed by its name, signature, and implementation" means that when defining a function in the Go programming language, there is a specific syntax and structure that needs to be followed.

Here's a breakdown of the components mentioned:

 - `func` keyword: The `func` keyword is used to indicate the start of a function declaration in Go.

- Function name: Following the func keyword, you provide the name of the function. The name should be a valid identifier and should follow Go's naming conventions.

- Function signature: The function signature consists of the function parameters and their types, as well as the return type (if any). It defines the inputs the function expects and the outputs it produces. The signature is enclosed in parentheses `()`. For example, `func add(a int, b int) int ` declares a function named `add` that takes two `int` parameters and returns an `int` value.

- Function implementation: The implementation of the function is the actual block of code that is executed when the function is called. It is enclosed in curly braces `{}`. The code inside the braces defines the logic and operations performed by the function.

Here's an example of a simple Go function that adds two integers:

go
```
func add(a int, b int) int {
    return a + b
}
```
In the above example, the function is named `add`, it takes two `int` parameters `a` and `b`, and it returns their sum as an `int` value.

By following the defined structure of a Go function, you can create reusable code blocks that encapsulate specific functionality and can be called from other parts of your program.


You might now ask how functions are organized and delivered. Well, the answer is in packages—the next subsection sheds some light on that.

## Introducing packages

Go programs are organized in packages—even the smallest Go program should be delivered as a package. The `package` keyword helps you define the name of a new package, which can be anything you want with just one exception: if you are creating an executable application and not just a package that will be shared by other applications or packages, you should name your package `main`. You will learn more about developing Go packages in _Chapter 5_, _Go Packages and Functions_.

Packages can be used by other packages. In fact, reusing existing packages is a good practice that saves you from having to write lots of code or implement existing functionality from scratch.

The `import` keyword is used for importing other Go packages in your Go programs in order to use some or all of their functionality. A Go package can either be a part of the rich Standard Go library or come from an external source. Packages of the standard Go library are imported by name (`os`) without the need for a hostname and a path, whereas external packages are imported using their full internet paths, like `github.com/spf13/cobra`.

In Go, a package is a unit of code organization. It consists of one or more Go source files that work together to provide a set of related functions, types, and other declarations. Packages can be imported and used by other code, allowing for modularity and reusability.

The `package` keyword is used to define the name of a new package. It appears at the top of each Go source file and is followed by the name of the package. 


#### For example:

```
package mypackage
```

The name you choose for your package can be anything you want, as long as it follows Go's naming conventions for identifiers. You can give your package a meaningful name that reflects its purpose or functionality.

However, there is one special case related to executable applications. When you are creating an executable program in Go, the package that contains the entry point of the program must be named `main`.

The `main` package serves as the entry point for Go programs. It contains the `main()` function, which is the starting point of execution when the program is run. The `main()` function is mandatory in the `main` package and should have the following signature:

```
func main() {
    // Program logic goes here
}
```
The `main` package cannot be imported by other code because it represents the top-level program rather than a reusable package. When you compile a Go file with the `main` package, it produces an executable file that can be run independently.

Here's an example of a simple Go program with the `main` package and the `main()` function:


```
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

In this example, the package name is `main`, and the `main()` function serves as the entry point of the program. The code inside the `main()` function is executed when the program is run, printing "Hello, World!" to the console.

By following this convention, Go distinguishes between packages that are intended to be imported and used by other code and the `main` package, which represents an executable program with an entry point.

#### Comparing Go Packages to Javascript framework configuration files

In a Vue 3 program, there are several configuration files that are essential for the project's setup and serve purposes similar to the main package in Go. These configuration files help define the project's structure, dependencies, build process, and entry points. Here are the main configuration files in a Vue 3 program:

- `package.json`: This file contains metadata about the project and lists the project's dependencies, devDependencies, and scripts. It serves as the entry point for managing dependencies, running scripts, and defining project-level configurations.

- `vue.config.js` (optional): This file is used to customize the Vue CLI configuration. It allows you to modify various aspects of the build process, configure webpack, specify plugin options, define development server settings, and more. It acts as the main configuration file for your Vue project's build and development setup.

These configuration files collectively define the project's structure, build process, and settings, similar to how the main package in Go defines the entry point and configuration for an executable program. They provide flexibility and customization options for your Vue 3 project, allowing you to tailor it to your specific needs and requirements.

# Running Go code

You now need to know how to execute `hw.go` or any other Go application. As will be explained in the two subsections that follow, there are two ways to execute Go code: as a compiled language using `go build` or as a scripting language using `go run`. So, let us find out more about these two ways of running Go code.

## Compiling Go code

In order to compile Go code and create a binary executable file, you need to use the `go build` command. What `go build` does is create an executable file for you to distribute and execute manually. This means that the `go build` command requires an additional step for running your code.

The generated executable is automatically named after the source code filename without the `.go` file extension. Therefore, because of the `hw.go` source file, the executable will be called `hw`. In case this is not what you want, `go build` supports the `-o` option, which allows you to change the filename and the path of the generated executable file. As an example, if you want to name the executable file as `helloWorld`, you should execute `go build -o helloWorld hw.go` instead. If no source files are provided, `go build` looks for a `main` package in the current directory.

After that, you need to execute the generated executable binary file on your own. In our case, this means executing either `hw` or `helloWorld`. This is shown in the next output:

```markup
$ go build hw.go
$ ./hw
Hello World!
```

Now that we know how to compile Go code, let us continue with using Go as a scripting language.

## Using Go like a scripting language

The `go run` command builds the named Go package, which in this case is the `main` package implemented in a single file, creates a temporary executable file, executes that file, and deletes it once it is done—to our eyes, this looks like using a scripting language. In our case, we can do the following:

```markup
$ go run hw.go
Hello World!
```

If you want to test your code, then using `go run` is a better choice. However, if you want to create and distribute an executable binary, then `go build` is the way to go.

## Important formatting and coding rules

You should know that Go comes with some strict formatting and coding rules that help the developer avoid beginner mistakes and bugs—once you learn these few rules and Go idiosyncrasies as well as the implications they have for your code, you will be free to concentrate on the actual functionality of your code. Additionally, the Go compiler is here to help you follow these rules with its expressive error messages and warnings. Last, Go offers standard tooling (`gofmt`) that can format your code for you so you never have to think about it.

The following is a list of important Go rules that will help you while reading this chapter:

-   Go code is delivered in packages and you are free to use the functionality found in existing packages. However, if you are going to import a package, you should use some of this functionality—there are some exceptions to this rule that mainly have to with initializing connections, but they are not important for now.
-   You either use a variable or you do not declare it at all. This rule helps you avoid errors such as misspelling an existing variable or function name.
-   There is only one way to format curly braces in Go.
-   Coding blocks in Go are embedded in curly braces even if they contain just a single statement or no statements at all.
-   Go functions can return multiple values.
-   You cannot automatically convert between different data types, even if they are of the same kind. As an example, you cannot implicitly convert an integer to a floating point.

Go has more rules but these rules are the most important ones and will keep you going for most of the book. You are going to see all these rules in action in this chapter as well as other chapters. For now, let's consider the only way to format curly braces in Go because this rule applies everywhere.

Look at the following Go program named `bad-curly.go`:

```markup
package main
import (
    "fmt"
)
func main() 
{
    fmt.Println("Go has strict rules for curly braces!")
}
```

Although it looks just fine, if you try to execute it, you will be fairly disappointed, because the code will not compile and therefore you will get the following syntax error message:

```markup
$ go run curly.go
# command-line-arguments
./curly.go:7:6: missing function body
./curly.go:8:1: syntax error: unexpected semicolon or newline before {
```

The official explanation for this error message is that Go requires the use of semicolons as statement terminators in many contexts, and the compiler automatically inserts the required semicolons when it thinks that they are necessary. Therefore, putting the opening curly brace (`{`) in its own line will make the Go compiler insert a semicolon at the end of the previous line (`func main()`), which is the main cause of the error message. The correct way to write the previous code is the following in `good-curly.go`:

```markup
package main
import (
    "fmt"
)
func main() {
    fmt.Println("Go has strict rules for curly braces!")
}
```

After learning about this global rule, let us continue by presenting some important characteristics of Go.

Bookmark

# Important characteristics of Go

This big section discusses important and essential Go features including variables, controlling program flow, iterations, getting user input, and Go concurrency. We begin by discussing variables, variable declaration, and variable usage.

## Defining and using variables

ImagineDevOps  that you wanted to perform some basic mathematical calculations with Go. In that case, you need to define variables to keep your input and your results.

Go provides multiple ways to declare new variables in order to make the variable declaration process more natural and convenient. You can declare a new variable using the `var` keyword followed by the variable name, followed by the desired data type (we will cover data types in detail in _Chapter 2_, _Basic Go Data Types_). If you want, you can follow that declaration with `=` and an initial value for your variable. If there is an initial value given, you can omit the data type and the compiler will guess it for you.

This brings us to a very important Go rule: **if no initial value is given to a variable, the Go compiler will automatically initialize that variable to the zero value of its data type**.

There is also the `:=` notation, which can be used instead of a `var` declaration. `:=` defines a new variable by inferring the data of the value that follows it. The official name for `:=` is **short assignment statement** and it is very frequently used in Go, especially for getting the return values from functions and `for` loops with the `range` keyword.

The short assignment statement can be used in place of a `var` declaration with an implicit type. You rarely see the use of `var` in Go; the `var` keyword is mostly used for declaring global or local variables without an initial value. The reason for the former is that every statement that exists outside of the code of a function must begin with a keyword such as `func` or `var`.

This means that the short assignment statement cannot be used outside of a function environment because it is not available there. Last, you might need to use `var` when you want to be explicit about the data type. For example, when you want `int8` or `int32` instead of `int`.

Therefore, although you can declare local variables using either `var` or `:=`, only `const` (when the value of a variable is not going to change) and `var` work for **global variables**, which are variables that are defined outside of a function and are not embedded in curly braces. Global variables can be accessed from anywhere in a package without the need to explicitly pass them to a function and can be changed unless they were defined as constants using the `const` keyword.

### Printing variables

Programs tend to display information, which means that they need to print data or send it somewhere for other software to store or process it. For printing data on the screen, Go uses the functionality of the `fmt` package. If you want Go to take care of the printing, then you might want to use the `fmt.Println()` function. However, there are times that you want to have full control over how data is going to get printed. In such cases, you might want to use `fmt.Printf()`.

`fmt.Printf()` is similar to the C `printf()` function and requires the use of control sequences that specify the data type of the variable that is going to get printed. Additionally, the `fmt.Printf()` function allows you to format the generated output, which is particularly convenient for floating point values because it allows you to specify the digits that will be displayed in the output (`%.2f` displays 2 digits after the decimal point). Lastly, the `\n` character is used for printing a newline character and therefore creating a new line, as `fmt.Printf()` does not automatically insert a newline—this is not the case with `fmt.Printf()`, which inserts a newline.

The following program illustrates how you can declare new variables, how to use them, and how to print them—type the following code into a plain text file named `variables.go`:

```markup
package main
import (
    "fmt"
    "math"
)
var Global int = 1234
var AnotherGlobal = -5678
func main() {
    var j int
    i := Global + AnotherGlobal
    fmt.Println("Initial j value:", j)
    j = Global
    // math.Abs() requires a float64 parameter
    //so we type cast it appropriately
    k := math.Abs(float64(AnotherGlobal))
    fmt.Printf("Global=%d, i=%d, j=%d k=%.2f.\n", Global, i, j, k)
}
```

Personally, I prefer to make global variables stand out by either beginning them with an uppercase letter or using all capital letters.

This program contains the following:

-   A global `int` variable named `Global`.
-   A second global variable named `AnotherGlobal`—Go automatically infers its data type from its value, which in this case is an integer.
-   A local variable named `j` that is of type `int`, which, as you will learn in the next chapter, is a special data type. `j` does not have an initial value, which means that Go automatically assigns the zero value of its data type, which in this case is `0`.
-   Another local variable named `i`—Go infers its data type from its value. As it is the sum of two `int` values, it is also an `int`.
-   As `math.Abs()` requires a `float64` parameter, you cannot pass `AnotherGlobal` to it because `AnotherGlobal` is an `int` variable. The `float64()` type cast converts the value of `AnotherGlobal` to `float64`. Note that `AnotherGlobal` continues to be `int`.
-   Lastly, `fmt.Printf()` formats and prints our output.

Running `variables.go` produces the following output:

```markup
Initial j value: 0
Global=1234, i=-4444, j=1234 k=5678.00.
```

This example demonstrated another important Go rule that was also mentioned previously: Go does not allow **implicit data conversions** like C.

As you saw in `variables.go` when using the `math.Abs()` function that expects a `float64` value, an `int` value cannot be used when a `float64` value is expected even if this particular conversion is straightforward and error-free. The Go compiler refuses to compile such statements. You should convert the `int` value to a `float64` explicitly using `float64()` for things to work properly.

For conversions that are not straightforward (for example, `string` to `int`), there exist specialized functions that allow you to catch issues with the conversion in the form of an `error` variable that is returned by the function.

## Controlling program flow

So far, we have seen Go variables but how do we change the flow of a Go program based on the value of a variable or some other condition? Go supports the `if/else` and `switch` control structures. Both control structures can be found in most modern programming languages, so if you have already programmed in another programming language, you should already be familiar with `if` and `switch`. `if` statements use no parenthesis for embedding the conditions that need to be examined because Go does not use parentheses in general. As expected, `if` has support for `else` and `else if` statements.

To demonstrate the use of `if`, let's use a very common pattern in Go that is used almost everywhere. This pattern says that if the value of an `error` variable as returned from a function is `nil`, then everything is OK with the function execution. Otherwise, there is an error condition somewhere that needs special care. This pattern is usually implemented as follows:

```markup
err := anyFunctionCall()
if err != nil {
    // Do something if there is an error
}
```

`err` is the variable that holds the `error` value as returned from a function and `!=` means that the value of the `err` variable is not `nil`. You will see similar code multiple times in Go programs.

Lines beginning with `//` are single-line comments. If you put `//` in the middle of a line, then everything after `//` is considered a comment. This rule does not apply if `//` is inside a string value.

The `switch` statement has two different forms. In the first form, the `switch` statement has an expression that is being evaluated, whereas in the second form, the `switch` statement has no expression to evaluate. In that case, expressions are evaluated in each `case` statement, which increases the flexibility of `switch`. The main benefit you get from `switch` is that when used properly, it simplifies complex and hard-to-read `if-else` blocks.

Both `if` and `switch` are illustrated in the following code, which is designed to process user input given as a command-line argument—please type it and save it as `control.go`. For learning purposes, we present the code of `control.go` in pieces in order to explain it better:

```markup
package main
import (
    "fmt"
    "os"
    "strconv"
)
```

This first part contains the expected preamble with the imported packages. The implementation of the `main()` function starts next:

```markup
func main() {
    if len(os.Args) != 2 {
        fmt.Println("Please provide a command line argument")
        return
    }
    argument := os.Args[1]
```

This part of the program makes sure that you have a single command-line argument to process, which is accessed as `os.Args[1]`, before continuing. We will cover this in more detail later, but you can refer to _Figure 1.2_ for more information about the `os.Args` slice.

```markup
    // With expression after switch
    switch argument {
    case "0":
        fmt.Println("Zero!")
    case "1":
        fmt.Println("One!")
    case "2", "3", "4":
        fmt.Println("2 or 3 or 4")
        fallthrough
    default:
        fmt.Println("Value:", argument)
    }
```

Here you see a `switch` block with four branches. The first three require exact `string` matches and the last one matches everything else. The order of the `case` statements is important because only the first match is executed. The `fallthrough` keyword tells Go that after this branch is executed, it will continue with the next branch, which in this case is the `default` branch:

```markup
    value, err := strconv.Atoi(argument)
    if err != nil {
        fmt.Println("Cannot convert to int:", argument)
        return
    }
```

As command-line arguments are initialized as string values, we need to convert user input into an integer value using a separate call, which in this case is a call to `strconv.Atoi()`. If the value of the `err` variable is `nil`, then the conversion was successful, and we can continue. Otherwise, an error message is printed onscreen and the program exits.

The following code shows the second form of `switch`, where the condition is evaluated at each `case` branch:

```markup
    // No expression after switch
    switch {
    case value == 0:
        fmt.Println("Zero!")
    case value > 0:
        fmt.Println("Positive integer")
    case value < 0:
        fmt.Println("Negative integer")
    default:
        fmt.Println("This should not happen:", value)
    }
}
```

This gives you more flexibility but requires more thinking when reading the code. In this case, the `default` branch should not be executed, mainly because any valid integer value would be caught by the other three branches. Nevertheless, the `default` branch is there, which is a good practice because it can catch unexpected values.

Running `control.go` generates the next output:

```markup
$ go run control.go 10
Value: 10
Positive integer
$ go run control.go 0
Zero!
Zero!
```

Each one of the two `switch` blocks in `control.go` creates one line of output.

## Iterating with for loops and range

This section is all about iterating in Go. Go supports `for` loops as well as the `range` keyword for iterating over all the elements of arrays, slices, and (as you will see in _Chapter 3_, _Composite Data Types_) maps. An example of Go simplicity is the fact that Go provides support for the `for` keyword only, instead of including direct support for `while` loops. However, depending on how you write a `for` loop, it can function as a `while` loop or an infinite loop. Moreover, `for` loops can implement the functionality of JavaScript's `forEach` function when combined with the `range` keyword.

You need to put curly braces around a `for` loop even if it contains a single statement or no statements at all.

You can also create `for` loops with variables and conditions. A `for` loop can be exited with a `break` keyword and you can skip the current iteration with the `continue` keyword. When used with `range`, `for` loops allow you to visit all the elements of a slice or an array without knowing the size of the data structure. As you will see in _Chapter 3_, _Composite Data Types_, `for` and `range` allow you to iterate over the elements of a map in a similar way.

The following program illustrates the use of `for` on its own and with the `range` keyword—type it and save it as `forLoops.go` in order to execute it afterward:

```markup
package main
import "fmt"
func main() {
    // Traditional for loop
    for i := 0; i < 10; i++ {
        fmt.Print(i*i, " ")
    }
    fmt.Println()
}
```

The previous code illustrates a traditional `for` loop that uses a local variable named `i`. This prints the squares of `0`, `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, and `9` onscreen. The square of `10` is not printed because it does not satisfy the `10 < 10` condition.

The following code is idiomatic Go:

```markup
    i := 0
    for ok := true; ok; ok = (i != 10) {
        fmt.Print(i*i, " ")
        i++
    }
    fmt.Println()
```

You might use it, but it is sometimes hard to read, especially for people that are new to Go. The following code shows how a `for` loop can simulate a `while` loop, which is not supported directly:

```markup
    // For loop used as while loop
    i := 0
    for {
        if i == 10 {
            break
        }
        fmt.Print(i*i, " ")
        i++
    }
    fmt.Println()
```

The `break` keyword in the `if` condition exits the loop early and acts as the loop exit condition.

Lastly, given a slice named `aSlice`, you iterate over all its elements with the help of `range`, which returns two ordered values: the index of the current element in the slice and its value. If you want to ignore either of these return values, which is not the case here, you can use `_` in the place of the value that you want to ignore. If you just need the index, you can leave out the second value from `range` entirely without using `_`.

```markup
    // This is a slice but range also works with arrays
    aSlice := []int{-1, 2, 1, -1, 2, -2}
    for i, v := range aSlice {
        fmt.Println("index:", i, "value: ", v)
    }
```

If you run `forLoops.go`, you get the following output:

```markup
$ go run forLoops.go
0 1 4 9 16 25 36 49 64 81
0 1 4 9 16 25 36 49 64 81
0 1 4 9 16 25 36 49 64 81
index: 0 value:  -1
index: 1 value:  2
index: 2 value:  1
index: 3 value:  -1
index: 4 value:  2
index: 5 value:  -2
```

The previous output illustrates that the first three `for` loops are equivalent and therefore produce the same output. The last six lines show the index and the value of each element found in `aSlice`.

Now that we know about `for` loops, let us see how to get user input.

## Getting user input

Getting user input is an important part of every program. This section presents two ways of getting user input, which are reading from standard input and using the command-line arguments of the program.

### Reading from standard input

The `fmt.Scanln()` function can help you read user input while the program is already running and store it to a `string` variable, which is passed as a pointer to `fmt.Scanln()`. The `fmt` package contains additional functions for reading user input from the console (`os.Stdin`), from files or from argument lists.

The following code illustrates reading from standard input—type it and save it as `input.go`:

```markup
package main
import (
    "fmt"
)
func main() {
    // Get User Input
    fmt.Printf("Please give me your name: ")
    var name string
    fmt.Scanln(&name)
    fmt.Println("Your name is", name)
}
```

While waiting for user input, it is good to let the user know what kind of information they have to give, which is the purpose of the `fmt.Printf()` call. The reason for not using `fmt.Println()` instead is that `fmt.Println()` automatically adds a newline character at the end, which is not what we want here.

Executing `input.go` generates the following kind of output and user interaction:

```markup
$ go run input.go
Please give me your name: Mihalis
Your name is Mihalis
```

### Working with command-line arguments

Although typing user input when needed might look like a nice idea, this is not usually how real software works. Usually, user input is given in the form of command-line arguments to the executable file. By default, command-line arguments in Go are stored in the `os.Args` slice. Go also offers the `flag` package for parsing command-line arguments, but there are better and more powerful alternatives.

The figure that follows shows the way command-line arguments work in Go, which is the same as in the C programming language. It is important to know that the `os.Args` slice is properly initialized by Go and is available to the program when referenced. The `os.Args` slice contains `string` values:

![Text
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_01_02.png)

Figure 1.2: How the os.Args slice works

The first command-line argument stored in the `os.Args` slice is always the name of the executable. If you are using `go run`, you will get a temporary name and path, otherwise, it will be the path of the executable as given by the user. The remaining command-line arguments are what comes after the name of the executable—the various command-line arguments are automatically separated by space characters unless they are included in double or single quotes.

The use of `os.Args` is illustrated in the code that follows, which is to find the minimum and the maximum numeric values of its input while ignoring invalid input such as characters and strings. Type the code and save it as `cla.go` (or any other filename you want):

```markup
package main
import (
    "fmt"
    "os"
    "strconv"
)
```

As expected, `cla.go` begins with its preamble. The `fmt` package is used for printing output whereas the `os` package is required because `os.Args` is a part of it. Lastly, the `strconv` package contains functions for converting strings to numeric values. Next, we make sure that we have at least one command-line argument:

```markup
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Need one or more arguments!")
        return
    }
```

Remember that the first element in `os.Args` is always the path of the executable file so `os.Args` is never totally empty. Next, the program checks for errors in the same way we have looked at in previous examples. You will learn more about errors and error handling in _Chapter 2_, _Basic Go Data Types_:

```markup
    var min, max float64
    for i := 1; i < len(arguments); i++ {
        n, err := strconv.ParseFloat(arguments[i], 64)
        if err != nil {
            continue
        }
```

In this case, we use the `error` variable returned by `strconv.ParseFloat()` to make sure that the call to `strconv.ParseFloat()` was successful and we have a valid numeric value to process. Otherwise, we should continue to the next command-line argument. The `for` loop is used for iterating over all available command-line arguments except the first one, which uses an index value of `0`. This is another popular technique for working with all command-line arguments.

The following code is used for properly initializing the value of `min` and `max` after the first command-line argument has been processed:

```markup
        if i == 1 {
            min = n
            max = n
            continue
        }
```

We are using `i == 1` as the test of if this is the first iteration. In this case, it is, so we are processing the first command-line argument. The next code checks whether the current value is our new minimum or maximum—this is where the logic of the program is implemented:

```markup
        if n < min {
            min = n
        }
        if n > max {
            max = n
        }
    }
    fmt.Println("Min:", min)
    fmt.Println("Max:", max)
}
```

The last part of the program is about printing your findings, which are the minimum and the maximum numeric values of all valid command-line arguments. The output you get from `cla.go` depends on its input:

```markup
$ go run cla.go a b 2 -1
Min: -1
Max: 2
```

In this case, `a` and `b` are invalid, and the only valid input are `-1` and `2`, which are the minimum value and the maximum value, respectively.

```markup
$ go run cla.go a 0 b -1.2 10.32
Min: -1.2
Max: 10.32
```

In this case, `a` and `b` are invalid input and therefore ignored.

```markup
$ go run cla.go
Need one or more arguments!
```

In the final case, as `cla.go` has no input to process, it prints a help message. If you execute the program with no valid input values, for example `go run cla.go a b c`, then both `Min` and `Max` values are going to be zero.

The next subsection shows a technique for differentiating between different data types using `error` variables.

## Using error variables to differentiate between input types

Now let me show you a technique that uses `error` variables to differentiate between various kinds of user input. For this technique to work, you should go from more specific cases to more generic ones. If we are talking about numeric values, you should first examine whether a string is a valid integer before examining whether the same string is a floating-point value because every valid integer is also a valid floating-point value.

This is illustrated in the next code excerpt:

```markup
    var total, nInts, nFloats int
    invalid := make([]string, 0)
    for _, k := range arguments[1:] {
        // Is it an integer?
        _, err := strconv.Atoi(k)
        if err == nil {
            total++
            nInts++
            continue
        }
```

First, we create three variables for keeping a count of the total number of valid values examined, the total number of integer values found, and the total number of floating-point values found, respectively. The `invalid` variable, which is a slice, is used for keeping all non-numeric values.

Once again, we need to iterate over all the command-line arguments except the first one, which has an index value of `0`, because this is the path of the executable file. We ignore the path of the executable using `arguments[1:]` instead of just `arguments`—selecting a continuous part of a slice is discussed in the next chapter.

The call to `strconv.Atoi()` determines whether we are processing a valid `int` value or not. If so, we increase the `total` and `nInts` counters:

```markup
        // Is it a float
        _, err = strconv.ParseFloat(k, 64)
        if err == nil {
            total++
            nFloats++
            continue
        }
```

Similarly, if the examined string represents a valid floating-point value, the call to `strconv.ParseFloat()` is going to be successful and the program will update the relevant counters. Lastly, if a value is not numeric, it is added to the `invalid` slice with a call to `append()`:

```markup
        // Then it is invalid
        invalid = append(invalid, k)
    }
```

This is a common practice for keeping unexpected input in applications. The previous code can be found as `process.go` in the GitHub repository of the book—not presented here is extra code that warns you when your invalid input is more than the valid one. Running `process.go` produces the next kind of output:

```markup
$ go run process.go 1 2 3
#read: 3 #ints: 3 #floats: 0
```

In this case, we process `1`, `2`, and `3`, which are all valid integer values.

```markup
$ go run process.go 1 2.1 a    
#read: 2 #ints: 1 #floats: 1
```

In this case, we have a valid integer, `1`, a floating-point value, `2.1`, and an invalid value, `a`.

```markup
$ go run process.go a 1 b
#read: 1 #ints: 1 #floats: 0
Too much invalid input: 2
a
b
```

If the invalid input is more than the valid one, then `process.go` prints an extra error message.

The next subsection discusses the concurrency model of Go.

## Understanding the Go concurrency model

This section is a quick introduction to the Go concurrency model. The Go concurrency model is implemented using goroutines and channels. A **goroutine** is the smallest executable Go entity. In order to create a new goroutine, you have to use the `go` keyword followed by a predefined function or an anonymous function—both methods are equivalent as far as Go is concerned.

Note that you can only execute functions or anonymous functions as goroutines.

A **channel** in Go is a mechanism that, among other things, allows goroutines to communicate and exchange data. If you are an amateur programmer or you're hearing about goroutines and channels for the first time, do not panic. Goroutines and channels, as well as pipelines and sharing data among goroutines, will be explained in much more detail in _Chapter 7_, _Go Concurrency_.

Although it is easy to create goroutines, there are other difficulties when dealing with concurrent programming including goroutine synchronization and sharing data between goroutines—this is a Go mechanism for avoiding side effects when running goroutines. As `main()` runs as a goroutine as well, you do not want `main()` to finish before the other goroutines of the program because when `main()` exits, the entire program along with any goroutines that have not finished yet will terminate. Although goroutines are do not share any variables, they can share memory. The good thing is that there are various techniques for the `main()` function to wait for goroutines to exchange data through channels or, less frequently in Go, using shared memory.

Type the following Go program, which synchronizes goroutines using `time.Sleep()` calls (this is not the right way to synchronize goroutines—we will discuss the proper way to synchronize goroutines in _Chapter 7_, _Go Concurrency_), into your favorite editor and save it as `goRoutines.go`:

```markup
package main
import (
    "fmt"
    "time"
)
func myPrint(start, finish int) {
    for i := start; i <= finish; i++ {
        fmt.Print(i, " ")
    }
    fmt.Println()
    time.Sleep(100 * time.Microsecond)
}
func main() {
    for i := 0; i < 5; i++ {
        go myPrint(i, 5)
    }
    time.Sleep(time.Second)
}
```

The preceding naively implemented example creates 4 goroutines and prints some values on the screen using the `myPrint()` function—the `go` keyword is used for creating goroutines. Running `goRoutines.go` generates the next output:

```markup
$ go run goRoutines.go
2 3 4 5
0 4 1 2 3 1 2 3 4 4 5
5
3 4 5
5
```

However, if you run it multiple times, you'll most likely get a different output each time:

```markup
1 2 3 4 5 
4 2 5 3 4 5 
3 0 1 2 3 4 5 
4 5
```

This happens because goroutines are initialized in random order and start running in random order. The Go scheduler is responsible for the execution of goroutines just like the OS scheduler is responsible for the execution of the OS threads. _Chapter 7_, _Go Concurrency_, discusses Go concurrency in more detail and presents the solution to that randomness issue with the use of a `sync.WaitGroup` variable—however, keep in mind that Go concurrency is everywhere, which is the main reason for including this section here. Therefore, as some error messages generated by the compiler talk about goroutines, you should not think that these goroutines were created by you.

The next section shows a practical example, which is developing a Go version of the `which(1)` utility, which locates a program file in the user's `PATH` value.

Bookmark

# Developing the which(1) utility in Go

Go can work with your operating system through a set of packages. A good way of learning a new programming language is by trying to implement simple versions of traditional UNIX utilities. In this section, you'll see a Go version of the `which(1)` utility, which will help you understand the way Go interacts with the underlying OS and reads environment variables.

The presented code, which will implement the functionality of `which(1)`, can be divided into three logical parts. The first part is about reading the input argument, which is the name of the executable file that the utility will be searching for. The second part is about reading the `PATH` environment variable, splitting it, and iterating over the directories of the `PATH` variable. The third part is about looking for the desired binary file in these directories and determining whether it can be found or not, whether it is a regular file, and whether it is an executable file. If the desired executable file is found, the program terminates with the help of the `return` statement. Otherwise, it will terminate after the `for` loop ends and the `main()` function exits.

Now let us see the code, beginning with the logical preamble that usually includes the package name, the `import` statements, and other definitions with a global scope:

```markup
package main
import (
    "fmt"
    "os"
    "path/filepath"
)
```

The `fmt` package is used for printing onscreen, the `os` package is for interacting with the underlying operating system, and the `path/filepath` package is used for working with the contents of the `PATH` variable that is read as a long string, depending on the number of directories it contains.

The second logical part of the utility is the following:

```markup
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide an argument!")
        return
    }
    file := arguments[1]
    path := os.Getenv("PATH")
    pathSplit := filepath.SplitList(path)
    for _, directory := range pathSplit {
```

First, we read the command-line arguments of the program (`os.Args`) and save the first command-line argument into the `file` variable. Then, we get the contents of the `PATH` environment variable and split it using `filepath.SplitList()`, which offers a **portable** way of separating a list of paths. Lastly, we iterate over all the directories of the `PATH` variable using a `for` loop with `range` as `filepath.SplitList()` returns a slice.

The rest of the utility contains the following code:

```markup
        fullPath := filepath.Join(directory, file)
        // Does it exist?
        fileInfo, err := os.Stat(fullPath)
        if err == nil {
            mode := fileInfo.Mode()
            // Is it a regular file?
            if mode.IsRegular() {
                // Is it executable?
                if mode&0111 != 0 {
                    fmt.Println(fullPath)
                    return
                }
            }
        }
    }
}
```

We construct the full path that we examine using `filepath.Join()` that is used for concatenating the different parts of a path using an **OS-specific** separator—this makes `filepath.Join()` work in all supported operating systems. In this part, we also get some lower-level information about the file—remember that in UNIX everything is a file, which means that we want to make sure that we are dealing with a regular file that is also executable.

In this first chapter, we are including the entire code of the presented source files. However, starting from _Chapter 2_, _Basic Go Types_, this will not be the case. This serves two purposes: the first one is that you get to see the code that really matters and the second one is that we save book space.

Executing `which.go` generates the following kind of output:

```markup
$ go run which.go which
/usr/bin/which
$ go run which.go doesNotExist
```

The last command could not find the `doesNotExist` executable—according to the UNIX philosophy and the way UNIX pipes work, utilities generate no output onscreen if they have nothing to say. However, an exit code of `0` means success whereas a **non-zero** exit code usually means failure.

Although it is useful to print error messages onscreen, there are times that you need to keep all error messages together and be able to search them when it is convenient for you. In this case, you need to use one or more log files.

Bookmark

# Logging information

All UNIX systems have their own log files for writing logging information that comes from running servers and programs. Usually, most system log files of a UNIX system can be found under the `/var/log` directory. However, the log files of many popular services, such as Apache and Nginx, can be found elsewhere, depending on their configuration.

Logging and putting logging information in log files is a practical way of examining data and information from your software asynchronously either locally or at a central log server or using server software such as Elasticsearch, Beats, and Grafana Loki.

Generally speaking, using a log file to write some information used to be considered a better practice than writing the same output on screen for two reasons: firstly, because the output does not get lost as it is stored on a file, and secondly, because you can search and process log files using UNIX tools, such as `grep(1)`, `awk(1)`, and `sed(1)`, which cannot be done when messages are printed on a terminal window. However, this is not true anymore.

As we usually run our services via `systemd`, programs should log to `stdout` so `systemd` can put logging data in the journal. [https://12factor.net/logs](https://12factor.net/logs) offers more information about app logs. Additionally, in cloud native applications, we are encouraged to simply log to `stderr` and let the container system redirect the `stderr` stream to the desired destination.

The UNIX logging service has support for two properties named **logging level** and **logging facility**. The logging level is a value that specifies the severity of the log entry. There are various logging levels, including `debug`, `info`, `notice`, `warning`, `err`, `crit`, `alert`, and `emerg`, in reverse order of severity. The `log` package of the standard Go library does not support working with logging levels. The logging facility is like a category used for logging information. The value of the logging facility part can be one of `auth`, `authpriv`, `cron`, `daemon`, `kern`, `lpr`, `mail`, `mark`, `news`, `syslog`, `user`, `UUCP`, `local0`, `local1`, `local2`, `local3`, `local4`, `local5`, `local6`, or `local7` and is defined inside `/etc/syslog.conf`, `/etc/rsyslog.conf`, or another appropriate file depending on the server process used for system logging on your UNIX machine. This means that if a logging facility is not defined correctly, it will not be handled; therefore, the log messages you send to it might get ignored and therefore lost.

The `log` package sends log messages to standard error. Part of the `log` package is the `log/syslog` package, which allows you to send log messages to the `syslog` server of your machine. Although by default `log` writes to standard error, the use of `log.SetOutput()` modifies that behavior. The list of functions for sending logging data includes `log.Printf()`, `log.Print()`, `log.Println()`, `log.Fatalf()`, `log.Fatalln()`, `log.Panic()`, `log.Panicln()` and `log.Panicf()`.

Logging is for application code, not library code. If you are developing libraries, do not put logging in them.

In order to write to system logs, you need to call the `syslog.New()` function with the appropriate parameters. Writing to the main system log file is as easy as calling `syslog.New()` with the `syslog.LOG_SYSLOG` option. After that you need to tell your Go program that all logging information goes to the new logger—this is implemented with a call to the `log.SetOutput()` function. The process is illustrated in the following code—type it on your favorite plain text editor and save it as `systemLog.go`:

```markup
package main
import (
    "log"
    "log/syslog"
)
func main() {
    sysLog, err := syslog.New(syslog.LOG_SYSLOG, "systemLog.go")
    if err != nil {
        log.Println(err)
        return
    } else {
        log.SetOutput(sysLog)
        log.Print("Everything is fine!")
    }
}
```

After the call to `log.SetOutput()`, all logging information goes to the `syslog` logger variable that sends it to `syslog.LOG_SYSLOG`. Custom text for the log entries coming from that program is specified as the second parameter to the `syslog.New()` call.

Usually, you want to store logging data on user-defined files because they group relevant information, which makes them easier to process and inspect.

Running `systemLog.go` generates no output—however, if you look at the system logs of a macOS Big Sur machine, for example, you will find entries like the following inside `/var/log/system.log`:

```markup
Dec  5 16:20:10 iMac systemLog.go[35397]: 2020/12/05 16:20:10 Everything is fine!
Dec  5 16:43:18 iMac systemLog.go[35641]: 2020/12/05 16:43:18 Everything is fine!
```

The number inside the brackets is the process ID of the process that wrote the log entry—in our case, `35397` and `35641`.

Similarly, if you execute `journalctl -xe` on a Linux machine, you can see entries similar to the next:

```markup
Dec 05 16:33:43 thinkpad systemLog.go[12682]: 2020/12/05 16:33:43 Everything is fine!
Dec 05 16:46:01 thinkpad systemLog.go[12917]: 2020/12/05 16:46:01 Everything is fine!
```

The output on your own operating system might be slightly different but the general idea is the same.

Bad things happen all the time, even to good people and good software. So, the next subsection covers the Go way of dealing with bad situations in your programs.

## log.Fatal() and log.Panic()

The `log.Fatal()` function is used when something erroneous has happened and you just want to exit your program as soon as possible after reporting that bad situation. The call to `log.Fatal()` terminates a Go program at the point where `log.Fatal()` was called after printing an error message. In most cases, this custom error message can be `Not enough arguments`, `Cannot access file`, or similar. Additionally, it returns back a non-zero exit code, which in UNIX indicates an error.

There are situations where a program is about to fail for good and you want to have as much information about the failure as possible—`log.Panic()` implies that something really unexpected and unknown, such as not being able to find a file that was previously accessed or not having enough disk space, has happened. Analogous to the `log.Fatal()` function, `log.Panic()` prints a custom message and immediately terminates the Go program.

Have in mind that `log.Panic()` is equivalent to a call to `log.Print()` followed by a call to `panic()`. `panic()` is a built-in function that stops the execution of the current function and begins panicking. After that, it returns to the caller function. On the other hand, `log.Fatal()` calls `log.Print()` and then `os.Exit(1)`, which is an immediate way of terminating the current program.

Both `log.Fatal()` and `log.Panic()` are illustrated in the `logs.go` file, which contains the next Go code:

```markup
package main
import (
    "log"
    "os"
)
func main() {
    if len(os.Args) != 1 {
        log.Fatal("Fatal: Hello World!")
    }
    log.Panic("Panic: Hello World!")
}
```

If you call `logs.go` without any command-line arguments, it calls `log.Panic()`. Otherwise, it calls `log.Fatal()`. This is illustrated in the next output from an Arch Linux system:

```markup
$ go run logs.go 
2020/12/03 18:39:26 Panic: Hello World!
panic: Panic: Hello World!
goroutine 1 [running]:
log.Panic(0xc00009ef68, 0x1, 0x1)
        /usr/lib/go/src/log/log.go:351 +0xae
main.main()
        /home/mtsouk/Desktop/mGo3rd/code/ch01/logs.go:12 +0x6b
exit status 2
$ go run logs.go 1
2020/12/03 18:39:30 Fatal: Hello World!
exit status 1
```

So, the output of `log.Panic()` includes additional low-level information that, hopefully, will help you to resolve difficult situations that happened in your Go code.

## Writing to a custom log file

Most of the time, and especially on applications and services that are deployed to production, you just need to write your logging data in a log file of your choice. This can be for many reasons, including writing debugging data without messing with the system log files, or keeping your own logging data separate from system logs in order to transfer it or store it in a database or software like Elasticsearch. This subsection teaches you how to write to a custom log file that is usually application-specific.

Writing to files and file I/O are both covered in _Chapter 6_, _Telling a UNIX System What to Do_—however, saving information to files is very handy when troubleshooting and debugging Go code, which is why this is covered in the first chapter of the book.

The path of the log file that is used is hardcoded into the code using a global variable named `LOGFILE`. For the purposes of this chapter and for preventing your file system from getting full in case something goes wrong, that log file resides inside the `/tmp` directory, which is not the usual place for storing data because usually, the `/tmp` directory is emptied after each system reboot.

Additionally, at this point, this will save you from having to execute `customLog.go` with root privileges and from putting unnecessary files into your precious system directories.

Type the following code and save it as `customLog.go`:

```markup
package main
import (
    "fmt"
    "log"
    "os"
    "path"
)
func main() {
    LOGFILE := path.Join(os.TempDir(), "mGo.log")
    f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// The call to os.OpenFile() creates the log file for writing, 
// if it does not already exist, or opens it for writing 
// by appending new data at the end of it (os.O_APPEND)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()
```

The `defer` keyword tells Go to execute the statement just before the current function returns. This means that `f.Close()` is going to be executed just before `main()` returns. We'll go into more detail on `defer` in _Chapter 5_, _Go Packages and Functions_.

```markup
    iLog := log.New(f, "iLog ", log.LstdFlags)
    iLog.Println("Hello there!")
    iLog.Println("Mastering Go 3rd edition!")
}
```

The last three statements create a new log file based on an opened file (`f`) and write two messages to it using `Println()`.

If you ever decide to use the code of `customLog.go` in a real application, you should change the path stored in `LOGFILE` into something that makes more sense.

Running `customLog.go` generates no output. However, what is really important is what has been written in the custom log file:

```markup
$ cat /tmp/mGo.log
iLog 2020/12/05 17:31:07 Hello there!
iLog 2020/12/05 17:31:07 Mastering Go 3rd edition!
```

## Printing line numbers in log entries

In this subsection, you'll learn how to print the filename as well as the line number in the source file where the statement that wrote a log entry is located.

The desired functionality is implemented with the use of `log.Lshortfile` in the parameters of `log.New()` or `SetFlags()`. The `log.Lshortfile` flag adds the filename as well as the line number of the Go statement that printed the log entry in the log entry itself. If you use `log.Llongfile` instead of `log.Lshortfile`, then you get the full path of the Go source file—usually, this is not necessary, especially when you have a really long path.

Type the following code and save it as `customLogLineNumber.go`:

```markup
package main
import (
    "fmt"
    "log"
    "os"
    "path"
)
func main() {
    LOGFILE := path.Join(os.TempDir(), "mGo.log")
    f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()
    LstdFlags := log.Ldate | log.Lshortfile
    iLog := log.New(f, "LNum ", LstdFlags)
    iLog.Println("Mastering Go, 3rd edition!")
    iLog.SetFlags(log.Lshortfile | log.LstdFlags)
    iLog.Println("Another log entry!")
}
```

In case you are wondering, you are allowed to change the format of the log entries during program execution—this means that when there is a reason, you can print more analytical information in the log entries. This is implemented with multiple calls to `iLog.SetFlags()`.

Running `customLogLineNumber.go` generates no output but writes the following entries in the file path that is specified by the value of the `LOGFILE` global variable:

```markup
$ cat /tmp/mGo.log 
LNum 2020/12/05 customLogLineNumber.go:24: Mastering Go, 3rd edition!
LNum 2020/12/05 17:33:23 customLogLineNumber.go:27: Another log entry!
```

You will most likely get a different output on your own machine, which is the expected behavior.

Bookmark

# Overview of Go generics

This section discusses Go generics, which is a forthcoming Go feature. Currently, generics and Go are under discussion by the Go community. However, one way or another, it is good to know how generics work, its philosophy, and what the generics discussions are about.

Go generics has been one of the most requested additions to the Go programming language. At the time of writing, it is said that generics is going to be part of Go 1.18.

The main idea behind generics in Go, as well as any other programming language that supports generics, is not having to write special code for supporting multiple data types when performing the same task.

Currently, Go supports multiple data types in functions such as `fmt.Println()` using the empty interface and reflection—both interfaces and reflection are discussed in _Chapter 4_, _Reflection and Interfaces_.

However, demanding every programmer to write lots of code and implement lots of functions and methods for supporting multiple custom data types is not the optimal solution—generics comes into play for providing an alternative to the use of interfaces and reflection for supporting multiple data types. The following code showcases how and where generics can be useful:

```markup
package main
import (
    "fmt"
)
func Print[T any](s []T) {
    for _, v := range s {
        fmt.Print(v, " ")
    }
    fmt.Println()
}
func main() {
    Ints := []int{1, 2, 3}
    Strings := []string{"One", "Two", "Three"}
    Print(Ints)
    Print(Strings)
}
```

In this case, we have a function named `Print()` that uses generics through a generics variable, which is specified by the use of `[T any]` after the function name and before the function parameters. Due to the use of `[T any]`, `Print()` can accept any slice of any data type and work with it. However, `Print()` does not work with input other than slices and that is fine because if your application supports slices of different data types, this function can still save you from having to implement **multiple** functions for supporting each distinct slice. This is the general idea behind generics.

In _Chapter 4_, _Reflection and Interfaces_, you will learn about the _empty interface_ and how it can be used for accepting data of any data type. However, the empty interface requires extra code for working with specific data types.

We end this section by stating some useful facts about generics:

-   You do not need to use generics in your programs all the time.
-   You can continue working with Go as before even if you use generics.
-   You can fully replace generics code with non-generics code. The question is are you willing to write the extra code required for this?
-   I believe that generics should be used when they can create simpler code and designs. It is better to have **repetitive straightforward** code than optimal abstractions that slow down your applications.
-   There are times that you need to limit the data types that are supported by a function that uses generics—this is not a bad thing as all data types do not share the same capabilities. Generally speaking, generics can be useful when processing data types that share some characteristics.

You need time to get used to generics and use generics at its full potential. Take your time. We will cover generics in more depth in _Chapter 13_, _Go Generics_.

Bookmark

# Developing a basic phone book application

In this section, to utilize the skills you've picked up so far, we will develop a basic phone book application in Go. Despite its limitations, the presented application is a command-line utility that searches a slice of structures that is statically defined (_hardcoded_) in the Go code. The utility offers support for two commands named `search` and `list` that search for a given surname and return its full record if the surname is found, and lists all available records, respectively.

The implementation has many shortcomings, including the following:

-   If you want to add or delete any data, you need to change the source code
-   You cannot present the data in a sorted form, which might be OK when you have 3 entries but might not work with more than 40 entries
-   You cannot export your data or load it from an external file
-   You cannot distribute the phone book application as a binary file because it uses hardcoded data
    
    The chapters that follow enhance the functionality of the phone book application in order to be fully functional, versatile, and powerful.
    

The code of `phoneBook.go` can be briefly described as follows:

-   There exists a new user-defined data type for holding the records of the phone book that is a Go structure with three fields named `Name`, `Surname`, and `Tel`. Structures group a set of values into a **single data type**, which allows you to pass and receive this set of values as a single entity.
-   There exists a global variable that holds the data of the phone book, which is a slice of structures named `data`.
-   There exist two functions that help you implement the functionality of the `search` and `list` commands.
-   The contents of the `data` global variable are defined in the `main()` function using multiple `append()` calls. You can change, add, or delete the contents of the `data` slice according to your needs.
-   Lastly, the program can only serve one task at a time. This means that to perform multiple queries, you have to run the program multiple times.

Let us now see `phoneBook.go` in more detail, beginning with its preamble:

```markup
package main
import (
    "fmt"
    "os"
)
```

After that, we have a section where we declare a Go structure named `Entry` as well as a global variable named `data`:

```markup
type Entry struct {
    Name    string
    Surname string
    Tel     string
}
var data = []Entry{}
```

After that, we define and implement two functions for supporting the functionality of the phone book:

```markup
func search(key string) *Entry {
    for i, v := range data {
        if v.Surname == key {
            return &data[i]
        }
    }
    return nil
}
func list() {
    for _, v := range data {
        fmt.Println(v)
    }
}
```

The `search()` function performs a linear search on the `data` slice. Linear search is slow, but it does the job for now considering that the phone book does not contain lots of entries. The `list()` function just prints the contents of the `data` slice using a `for` loop with `range`. As we are not interested in displaying the index of the element that we print, we ignore it using the `_` character and just print the structure that holds the actual data.

Lastly, we have the implementation of the `main()` function. The first part of it is as follows:

```markup
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        exe := path.Base(arguments[0])
        fmt.Printf("Usage: %s search|list <arguments>\n", exe)
        return
    }
```

The `exe` variable holds the path to the executable file—it is a nice and professional touch to print the name of the executable binary in the instructions of the program.

```markup
    data = append(data, Entry{"Mihalis", "Tsoukalos", "2109416471"})
    data = append(data, Entry{"Mary", "Doe", "2109416871"})
    data = append(data, Entry{"John", "Black", "2109416123"})
```

In this part, we check whether we were given any command arguments. If not (`len(arguments) == 1`), the program prints a message and exits by calling `return`. Otherwise, it puts the desired data in the `data` slice before continuing.

The rest of the `main()` function implementation is as follows:

```markup
    // Differentiate between the commands
    switch arguments[1] {
    // The search command
    case "search":
        if len(arguments) != 3 {
            fmt.Println("Usage: search Surname")
            return
        }
        result := search(arguments[2])
        if result == nil {
            fmt.Println("Entry not found:", arguments[2])
            return
        }
        fmt.Println(*result)
    // The list command
    case "list":
        list()
    // Response to anything that is not a match
    default:
        fmt.Println("Not a valid option")
    }
}
```

This code uses a `case` block, which is really handy when you want to write readable code and avoid using multiple and nested `if` blocks. That `case` block differentiates between the two supported commands by examining the value of `arguments[1]`. If the given command is not recognized, the `default` branch is executed instead. For the `search` command, `arguments[2]` is also examined.

Working with `phoneBook.go` looks as follows:

```markup
$ go build phoneBook.go
$ ./phoneBook list
{Mihalis Tsoukalos 2109416471}
{Mary Doe 2109416871}
{John Black 2109416123}
$ ./phoneBook search Tsoukalos
{Mihalis Tsoukalos 2109416471}
$ ./phoneBook search Tsouk
Entry not found: Tsouk
$ ./phoneBook
Usage: ./phoneBook search|list <arguments>
```

The first command lists the contents of the phone book whereas the second command searches for a given surname (`Tsoukalos`). The third command searches for something that does not exist in the phone book and the last command builds `phoneBook.go` and runs the generated executable without any arguments, which prints the instructions of the program.

Despite its shortcomings, `phoneBook.go` has a clean design that you can easily extend and works as expected, which is a great starting point. The phone book application will keep improving in the chapters that follow as we learn more advanced concepts.

Bookmark

# Exercises

-   Our version of `which(1)` stops after finding the first occurrence of the desired executable. Make the necessary changes to `which.go` in order to find all possible occurrences of the desired executable.
-   The current version of `which.go` processes the first command-line argument only. Make the necessary changes to `which.go` in order to accept and search the `PATH` variable for multiple executable binaries.
-   Read the documentation of the `fmt` package at [https://golang.org/pkg/fmt/](https://golang.org/pkg/fmt/).

Bookmark

# Summary

If you are using Go for the first time, the information in this chapter will help you understand the advantages of Go, how Go code looks, and some important characteristics of Go such as variables, iterations, flow control, and the Go concurrency model. If you already know Go, then this chapter is a good reminder of where Go excels and the kinds of software where it is advised to use Go. Lastly, we built a basic phone book application with the techniques that we have learned so far.

The next chapter discusses the basic data types of Go in more detail.

Bookmark

# Additional resources

-   The official Go website: [https://golang.org/](https://golang.org/)
-   The Go Playground: [https://play.golang.org/](https://play.golang.org/)
-   The `log` package: [https://golang.org/pkg/log/](https://golang.org/pkg/log/)
-   Elasticsearch Beats: [https://www.elastic.co/beats/](https://www.elastic.co/beats/)
-   Grafana Loki: [https://grafana.com/oss/loki/](https://grafana.com/oss/loki/)
-   Microsoft Visual Studio: [https://visualstudio.microsoft.com/](https://visualstudio.microsoft.com/)
-   The Standard Go library: [https://golang.org/pkg/](https://golang.org/pkg/)