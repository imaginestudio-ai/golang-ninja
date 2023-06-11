# Go Packages and Functions

The main focus of this chapter is Go **packages**, which are Go's way of organizing, delivering, and using code. The most common component of packages is **functions**, which are pretty flexible and powerful and are used for data processing and manipulation. Go also supports **modules**, which are packages with version numbers. This chapter will also explain the operation of `defer`, which is used for cleaning up and releasing resources.

Regarding the visibility of package elements, Go follows a simple rule that states that functions, variables, data types, structure fields, and so forth that begin with an uppercase letter are **public**, whereas functions, variables, types, and so on that begin with a lowercase letter are **private**. This is the reason why `fmt.Println()` is named `Println()` instead of just `println()`. The same rule applies not only to the name of a `struct` variable but to the fields of a `struct` variable—in practice, this means that you can have a `struct` variable with both private and public fields. However, this rule does not affect package names, which are allowed to begin with either uppercase or lowercase letters.

In summary, this chapter covers:

-   Go packages
-   Functions
-   Developing your own packages
-   Using GitHub to store Go packages
-   A package for working with a database
-   Modules
-   Creating better packages
-   Creating documentation
-   GitLab Runners and Go
-   GitHub Actions and Go
-   Versioning utilities

Just Imagine

# Go packages

Everything in Go is delivered in the form of packages. A Go package is a Go source file that begins with the `package` keyword, followed by the name of the package.

Note that packages can have structure. For example, the `net` package has several subdirectories, named `http`, `mail`, `rpc`, `smtp`, `textproto`, and `url`, which should be imported as `net/http`, `net/mail`, `net/rpc`, `net/smtp`, `net/textproto`, and `net/url`, respectively.

Apart from the packages of the Go standard library, there are external packages that can be imported using their full address and that should be downloaded on the local machine, before their first use. One such example is [https://github.com/spf13/cobra](https://github.com/spf13/cobra), which is stored in GitHub.

Packages are mainly used for grouping related functions, variables, and constants so that you can transfer them easily and use them in your own Go programs. Note that apart from the `main` package, Go packages are not autonomous programs and cannot be compiled into executable files on their own. As a result, if you try to execute a Go package as if it were an autonomous program, you are going to be disappointed:

```markup
$ go run aPackage.go
go run: cannot run non-main package
```

Instead, packages need to be called directly or indirectly from a `main` package in order to be used, as we have shown in previous chapters.

## Downloading Go packages

In this subsection, you will learn how to download external Go packages using [https://github.com/spf13/cobra](https://github.com/spf13/cobra) as an example. The `go get` command for downloading the `cobra` package is as follows:

```markup
$ go get github.com/spf13/cobra
```

Note that you can download the package without using `https://` in its address. The results can be found inside the `~/go` directory—the full path is `~/go/src/github.com/spf13/cobra`. As the `cobra` package comes with a binary file that helps you structure and create command-line utilities, you can find that binary file inside `~/go/bin` as `cobra`.

The following output, which was created with the help of the `tree(1)` utility, shows a high-level view with 3 levels of detail of the structure of `~/go` on my machine:

```markup
$ tree ~/go -L 3
/Users/mtsouk/go
├── bin
│   ├── cobra
│   ├── go-outline
│   ├── gocode
│   ├── gocode-gomod
│   ├── godef
│   ├── golint
│   ├── gopkgs
│   └── goreturns
├── pkg
│   ├── darwin_amd64
│   │   ├── github.com
│   │   ├── golang.org
│   │   ├── gonum.org
│   │   └── google.golang.org
│   ├── mod
│   │   ├── 9fans.net
│   │   ├── cache
│   │   ├── cloud.google.com
│   │   ├── github.com
│   │   ├── go.opencensus.io@v0.22.4
│   │   ├── golang.org
│   │   └── google.golang.org
│   └── sumdb
│       └── sum.golang.org
└── src
    ├── github.com
    │   ├── sirupsen
    │   └── spf13
    └── golang.org
        └── x
23 directories, 8 files
```

The `x` path, which is displayed last, is used by the Go team.

Basically, there are three main directories under `~/go` with the following properties:

-   The `bin` directory: This is where binary tools are placed.
-   The `pkg` directory: This is where reusable packages are put. The `darwin_amd64` directory, which can be found on macOS machines only, contains compiled versions of the installed packages. On a Linux machine, you can find a `linux_amd64` directory instead of `darwin_amd64`.
-   The `src` directory: This is where the source code of the packages is located. The underlying structure is based on the URL of the package you are looking for. So, the URL for the `github.com/spf13/viper` package is `~/go/src/github.com/spf13/viper`. If a package is downloaded as a module, then it will be located under `~/go/pkg/mod`.

Starting with Go 1.16, `go install` is the recommended way of building and installing packages in module mode. The use of `go get` is deprecated, but this chapter uses `go get` because it's commonly used online and is worth knowing about. However, most of the chapters in this book use `go mod init` and `go mod tidy` for downloading external dependencies for your own source files.

If you want to upgrade an existing package, you should execute `go get` with the `-u` option. Additionally, if you want to see what is happening behind the scenes, add the `-v` option to the `go get` command—in this case, we are using the `viper` package as an example, but we abbreviate the output:

```markup
$ go get -v github.com/spf13/viper
github.com/spf13/viper (download)
...
github.com/spf13/afero (download)
get "golang.org/x/text/transform": found meta tag get.metaImport{Prefix:"golang.org/x/text", VCS:"git", RepoRoot:"https://go.googlesource.com/text"} at //golang.org/x/text/transform?go-get=1
get "golang.org/x/text/transform": verifying non-authoritative meta tag
...
github.com/fsnotify/fsnotify
github.com/spf13/viper
```

What you can basically see in the output is the dependencies of the initial package being downloaded before the desired package—most of the time, you do not want to know that.

We will continue this chapter by looking at the most important package element: _functions_.

Just Imagine

# Functions

The main elements of packages are functions, which are the subject of this section.

Type methods and functions are implemented in the same way and sometimes, the terms functions and type methods are used interchangeably.

A piece of advice: functions must be as independent from each other as possible and must do one job (and only one job) well. So, if you find yourself writing functions that do multiple things, you might want to consider replacing them with multiple functions instead.

You should already know that all function definitions begin with the `func` keyword, followed by the function's signature and its implementation, and that functions accept none, one, or more arguments and return none, one, or more values back. The single-most popular Go function is `main()`, which is used in every executable Go program—the `main()` function accepts no parameters and returns nothing, but it is the starting point of every Go program. Additionally, when the `main()` function ends, the entire program ends as well.

## Anonymous functions

**Anonymous functions** can be defined inline without the need for a name, and they are usually used for implementing things that require a small amount of code. In Go, a function can return an anonymous function or take an anonymous function as one of its arguments. Additionally, anonymous functions can be attached to Go variables. Note that anonymous functions are called **lambdas** in functional programming terminology. Similar to that, a **closure** is a specific type of anonymous function that carries or _closes over_ variables that are in the same lexical scope as the anonymous function that was defined.

It is considered a good practice for anonymous functions to have a small implementation and a local focus. If an anonymous function does not have a local focus, then you might need to consider making it a regular function. When an anonymous function is suitable for a job, it is extremely convenient and makes your life easier; just do not use too many anonymous functions in your programs without having a good reason to. We will look at anonymous functions in action in a while.

## Functions that return multiple values

As you already know from functions such as `strconv.Atoi()`, functions can return multiple distinct values, which saves you from having to create a dedicated structure for returning and receiving multiple values from a function. However, if you have a function that returns more than 3 values, you should reconsider that decision and maybe redesign it to use a single structure or slice for grouping and returning the desired values as a single entity—this makes handling the returned values simpler and easier. Functions, anonymous functions, and functions that return multiple values are all illustrated in `functions.go`, as shown in the following code:

```markup
package main
import "fmt"
func doubleSquare(x int) (int, int) {
    return x * 2, x * x
}
```

This function returns two `int` values, without the need for having separate variables to keep them—the returned values are created on the fly. Note the compulsory use of **parentheses** when a function returns more than one value.

```markup
// Sorting from smaller to bigger value
func sortTwo(x, y int) (int, int) {
    if x > y {
        return y, x
    }
    return x, y
}
```

The preceding function returns two `int` values as well.

```markup
func main() {
    n := 10
    d, s := doubleSquare(n)
```

The previous statement reads the two return values of `doubleSquare()` and saves them in `d` and `s`.

```markup
    fmt.Println("Double of", n, "is", d)
    fmt.Println("Square of", n, "is", s)
    // An anonymous function
    anF := func(param int) int {
        return param * param
    }
```

The `anF` variable holds an **anonymous function** that requires a single parameter as input and returns a single value. The only difference between an anonymous function and a regular one is that the name of the anonymous function is `func()` and that there is no `func` keyword.

```markup
    fmt.Println("anF of", n, "is", anF(n))
    fmt.Println(sortTwo(1, -3))
    fmt.Println(sortTwo(-1, 0))
}
```

The last two statements print the return values of `sortTwo()`. Running `functions.go` produces the following output:

```markup
Double of 10 is 20
Square of 10 is 100
anF of 10 is 100
-3 1
-1 0
```

The subsection that follows illustrates functions that have named return values.

## The return values of a function can be named

Unlike C, Go allows you to name the return values of a Go function. Additionally, when such a function has a `return` statement without any arguments, the function automatically returns the current value of each named return value, in the order in which they were declared in the function signature.

The following function is included in `namedReturn.go`:

```markup
func minMax(x, y int) (min, max int) {
    if x > y {
        min = y
        max = x
        return min, max
```

This `return` statement returns the values stored in the `min` and `max` variables—both `min` and `max` are defined in the **function signature** and not in the function body.

```markup
    }
    min = x
    max = y
    return
}
```

This `return` statement is equivalent to `return min, max`, which is based on the function signature and the use of named return values.

Running `namedReturn.go` produces the following output:

```markup
$ go run namedReturn.go 1 -2
-2 1
-2 1
```

## Functions that accept other functions as parameters

Functions can accept other functions as parameters. The best example of a function that accepts another function as an argument can be found in the `sort` package. You can provide the `sort.Slice()` function with another function as an argument that specifies the way sorting is implemented. The signature of `sort.Slice()` is `func Slice(slice interface{}, less func(i, j int) bool)`. This means the following:

-   The `sort.Slice()` function does not return any data.
-   The `sort.Slice()` function requires two arguments, a slice of type `interface{}` and another function—the slice variable is modified inside `sort.Slice()`.
-   The function parameter of `sort.Slice()` is named `less` and should have the `func(i, j int) bool` signature—there is no need for you to name the anonymous function. The name `less` is required because all function parameters should have a name.
-   The `i` and `j` parameters of `less` are indexes of the `slice` parameter.

Similarly, there is another function in the `sort` package named `sort.SliceIsSorted()` that is defined as `func SliceIsSorted(slice interface{}, less func(i, j int) bool) bool`. `sort.SliceIsSorted()` returns a `bool` value and checks whether the `slice` parameter is sorted according to the rules of the second parameter, which is a function.

You are not obliged to use an anonymous function in either `sort.Slice()` or `sort.SliceIsSorted()`. You can define a regular function with the required signature and use that. However, using an anonymous function is more convenient.

The use of both `sort.Slice()` and `sort.SliceIsSorted()` is illustrated in the Go program that follows—the name of the source file is `sorting.go`:

```markup
package main
import (
    "fmt"
    "sort"
)
type Grades struct {
    Name    string
    Surname string
    Grade   int
}
func main() {
    data := []Grades{{"J.", "Lewis", 10}, {"M.", "Tsoukalos", 7},
        {"D.", "Tsoukalos", 8}, {"J.", "Lewis", 9}}
    isSorted := sort.SliceIsSorted(data, func(i, j int) bool {
        return data[i].Grade < data[j].Grade
    })
```

The `if else` block that follows checks the `bool` value of `sort.SliceIsSorted()` to determine whether the slice is sorted:

```markup
    if isSorted {
        fmt.Println("It is sorted!")
    } else {
        fmt.Println("It is NOT sorted!")
    }
    sort.Slice(data,
        func(i, j int) bool { return data[i].Grade < data[j].Grade })
    fmt.Println("By Grade:", data)
}
```

The call to `sort.Slice()` sorts the data according to the anonymous function that is passed as the second argument to `sort.Slice()`.

Running `sorting.go` produces the following output:

```markup
It is NOT sorted!
By Grade: [{M. Tsoukalos 7} {D. Tsoukalos 8} {J. Lewis 9} {J. Lewis 10}]
```

## Functions can return other functions

Apart from accepting functions as arguments, functions can also return anonymous functions, which can be handy when the returned function is not always the same but depends on the function's input or other external parameters. This is illustrated in `returnFunction.go`:

```markup
package main
import "fmt"
func funRet(i int) func(int) int {
    if i < 0 {
        return func(k int) int {
            k = -k
            return k + k
        }
    }
    return func(k int) int {
        return k * k
    }
}
```

The signature of `funRet()` declares that the function returns another function with the `func(int) int` signature. The implementation of the function is unknown, but it is going to be defined at runtime. Functions are returned using the `return` keyword. The developer should take care and save the returned function.

```markup
func main() {
    n := 10
    i := funRet(n)
    j := funRet(-4)
```

Note that `n` and `-4` are only used for determining the anonymous functions that are going to be returned from `funRet()`.

```markup
    fmt.Printf("%T\n", i)
    fmt.Printf("%T %v\n", j, j)
    fmt.Println("j", j, j(-5))
```

The first statement prints the signature of the function whereas the second statement prints the function signature and its memory address. The last statement also returns the memory address of `j`, because `j` is a pointer to the anonymous function and the value of `j(-5)`.

```markup
    // Same input parameter but DIFFERENT
    // anonymous functions assigned to i and j
    fmt.Println(i(10))
    fmt.Println(j(10))
}
```

Although both `i` and `j` are called with the same input (`10`), they are going to return different values because they store different anonymous functions.

Running `returnFunction.go` generates the following output:

```markup
func(int) int
func(int) int 0x10a8d40
j 0x10a8d40 10
100
-20
```

The first line of the output shows the data type of the `i` variable that holds the return value of `funRet(n)`, which is `func(int) int` as it holds a function. The second line of output shows the data type of `j`, as well as the memory address where the anonymous function is stored. The third line shows the memory address of the anonymous function stored in the `j` variable, as well as the return value of `j(-5)`. The last two lines are the return values of `i(10)` and `j(10)`, respectively.

So, in this subsection, we learned about functions returning functions. This makes Go a functional programming language, albeit not a pure one, and allows Go to benefit from the functional programming paradigm.

We are now going to examine variadic functions, which are functions with a variable number of parameters.

## Variadic functions

Variadic functions are functions that can accept a variable number of parameters—you already know about `fmt.Println()` and `append()`, which are both variadic functions that are widely used. In fact, most functions found in the `fmt` package are variadic.

The general ideas and rules behind variadic functions are as follows:

-   Variadic functions use the _pack operator_, which consists of a `...`, followed by a data type. So, for a variadic function to accept a variable number of `int` values, the pack operator should be `...int`.
-   The pack operator can only be used once in any given function.
-   The variable that holds the pack operation is a slice and, therefore, is accessed as a slice inside the variadic function.
-   The variable name that is related to the pack operator is always last in the list of function parameters.
-   When calling a variadic function, you should put a list of values separated by `,` in the place of the variable with the _pack operator_ or a slice with the _unpack operator_.

This list contains all the rules that you need to know in order to define and use variadic functions.

The pack operator can also be used with an empty interface. In fact, most functions in the `fmt` package use `...interface{}` to accept a variable number of arguments of all data types. You can find the source code of the latest implementation of `fmt` at [https://golang.org/src/fmt/](https://golang.org/src/fmt/).

However, there is a situation that needs special care here—I made that mistake when I was learning Go. If you try to pass `os.Args`, which is a slice of strings (`[]string`), as `...interface{}` to a variadic function), your code will not compile and will generate an error message similar to `cannot use os.Args (type []string) as type []interface {} in argument to <function_name>`. This happens because the two data types (`[]string` and `[]interface{}`) do not have the same representations in memory—this applies to all data types. In practice, this means that you cannot write `os.Args...` to pass each individual value of the `os.Args` slice to a variadic function.

On the other hand, if you just use `os.Args`, it will work, but this passes the entire slice as a _single entity_ instead of its individual values! This means that the `everything(os.Args, os.Args)` statement works but does not do what you want.

The solution to this problem is converting the slice of strings—or any other slice—into a slice of `interface{}`. One way to do that is by using the code that follows:

```markup
empty := make([]interface{}, len(os.Args[1:]))
for i, v := range os.Args {
    empty[i] = v
}
```

Now, you are allowed to use `empty...` as an argument to the variadic function. This is the only subtle point related to variadic functions and the pack operator.

As there is no standard library function to perform that conversion for you, you have to write your own code. Note that the conversion takes time because the code must visit all slice elements. The more elements the slice has, the more time the conversion will take. This topic is also discussed at [https://github.com/golang/go/wiki/InterfaceSlice](https://github.com/golang/go/wiki/InterfaceSlice).

We are now ready to see variadic functions in action. Type the following Go code using your favorite text editor and save it as `variadic.go`:

```markup
package main
import (
    "fmt"
    "os"
)
```

As variadic functions are built into the grammar of the language, you do not need anything extra to support variadic functions.

```markup
func addFloats(message string, s ...float64) float64 {
```

This is a variadic function that accepts a `string` and an unknown number of `float64` values. It prints the `string` variable and calculates the sum of the `float64` values.

```markup
    fmt.Println(message)
    sum := float64(0)
    for _, a := range s {
        sum = sum + a
    }
```

This `for` loop accesses the pack operator as a slice, so there is nothing special here.

```markup
    s[0] = -1000
    return sum
}
```

You can also access individual elements of the `s` slice.

```markup
func everything(input ...interface{}) {
    fmt.Println(input)
}
```

This is another variadic function that accepts an unknown number of `interface{}` values.

```markup
func main() {
    sum := addFloats("Adding numbers...", 1.1, 2.12, 3.14, 4, 5, -1, 10)
```

You can put the arguments of a variadic function inline.

```markup
    fmt.Println("Sum:", sum)
    s := []float64{1.1, 2.12, 3.14}
```

But you usually use a slice variable with the unpack operator.

```markup
    sum = addFloats("Adding numbers...", s...)
    fmt.Println("Sum:", sum)
    everything(s)
```

The previous code works because the content of `s` is not unpacked.

```markup
    // Cannot directly pass []string as []interface{}
    // You have to convert it first!
    empty := make([]interface{}, len(os.Args[1:]))
```

You can convert `[]string` into `[]interface{}` in order to use the unpack operator.

```markup
    for i, v := range os.Args[1:] {
        empty[i] = v
    }
    everything(empty...)
```

And now, we can unpack the contents of `empty`.

```markup
    arguments := os.Args[1:]
    empty = make([]interface{}, len(arguments))
    for i := range arguments {
        empty[i] = arguments[i]
    }
```

This is a slightly different way of converting `[]string` into `[]interface{}`.

```markup
    everything(empty...)
    // This will work!
    str := []string{"One", "Two", "Three"}
    everything(str, str, str)
}
```

The previous statement works because you are passing the entire `str` variable three times—not its contents. So, the slice contains three elements—each element is equal to the contents of the `str` variable.

Running `variadic.go` produces the following output:

```markup
$ go run variadic.go
Adding numbers...
Sum: 24.36
Adding numbers...
Sum: 6.36
[[-1000 2.12 3.14]]
[]
[]
[[One Two Three] [One Two Three] [One Two Three]]
```

The last line of the output shows that we have passed the `str` variable three times to the `everything()` function as three separate entities.

Variadic functions come in very handy when you want to have an unknown number of parameters in a function. The next subsection discusses the use of `defer`, which we have already used multiple times.

## The defer keyword

So far, we have seen `defer` in `ch03/csvData.go`, as well as in the implementations of the phone book application. But what does `defer` do? The `defer` keyword postpones the execution of a function until the surrounding function returns.

Usually, `defer` is used in file I/O operations to keep the function call that closes an opened file close to the call that opened it, so that you do not have to remember to close a file that you have opened just before the function exits.

It is very important to remember that **deferred functions** are executed in **last in, first out** (**LIFO**) order after the surrounding function has been returned. Putting it simply, this means that if you `defer` function `f1()` first, function `f2()` second, and function `f3()` third in the same surrounding function, then when the surrounding function is about to return, function `f3()` will be executed first, function `f2()` will be executed second, and function `f1()` will be the last one to get executed.

In this section, we will discuss the dangers of `defer`, when used carelessly using a simple program. The code for `defer.go` is as follows.

```markup
package main
import (
    "fmt"
)
func d1() {
    for i := 3; i > 0; i-- {
        defer fmt.Print(i, " ")
    }
}
```

In `d1()`, `defer` is executed inside the function body with just a `fmt.Print()` call. Remember that these calls to `fmt.Print()` are executed just before function `d1()` returns.

```markup
func d2() {
    for i := 3; i > 0; i-- {
        defer func() {
            fmt.Print(i, " ")
        }()
    }
    fmt.Println()
}
```

In `d2()`, `defer` is attached to an anonymous function that does not accept any parameters. In practice, this means that the anonymous function should get the value of `i` on its own—this is dangerous because the current value of `i` depends on when the anonymous function is executed.

The anonymous function is a **closure**, and that is why it has access to variables that would normally be out of scope.

```markup
func d3() {
    for i := 3; i > 0; i-- {
        defer func(n int) {
            fmt.Print(n, " ")
        }(i)
    }
}
```

In this case, the current value of `i` is passed to the anonymous function as a parameter that initializes the `n` function parameter. This means that there are no ambiguities about the value that `i` has.

```markup
func main() {
    d1()
    d2()
    fmt.Println()
    d3()
    fmt.Println()
}
```

The task of `main()` is to call `d1()`, `d2()`, and `d3()`.

Running `defer.go` produces the following output:

```markup
$ go run defer.go
1 2 3
0 0 0
1 2 3
```

You will most likely find the generated output complicated and challenging to understand, which proves that the operation and the results of the use of `defer` can be tricky if your code is not clear and unambiguous. Let me explain the results so that you get a better idea of how tricky `defer` can be if you do not pay close attention to your code.

Let's start with the first line of the output `(1 2 3)` that is generated by the `d1()` function. The values of `i` in `d1()` are `3`, `2`, and `1` in that order. The function that is deferred in `d1()` is the `fmt.Print()` statement; as a result, when the `d1()` function is about to return, you get the three values of the `i` variable of the `for` loop in reverse order. This is because deferred functions are executed in LIFO order.

Now, let me explain the second line of the output that is produced by the `d2()` function. It is really strange that we got three zeros instead of `1 2 3` in the output; however, there is a reason for that. After the `for` loop ended, the value of `i` is `0`, because it is that value of `i` that made the `for` loop terminate. However, the tricky point here is that the deferred anonymous function is evaluated after the `for` loop ends because it has no parameters, which means that it is evaluated three times for an `i` value of `0`, hence the generated output. This kind of confusing code is what might lead to the creation of nasty bugs in your projects, so try to avoid it.

Finally, we will talk about the third line of the output, which is generated by the `d3()` function. Due to the parameter of the anonymous function, each time the anonymous function is deferred, it gets and therefore uses the current value of `i`. As a result, each execution of the anonymous function has a different value to process without any ambiguities, hence the generated output.

After that, it should be clear that the best approach to using `defer` is the third one, which is exhibited in the `d3()` function, because you intentionally pass the desired variable in the anonymous function in an easy-to-read way. Now that we have learned about `defer`, it is time to discuss something completely different: how to develop your own packages.

Just Imagine

# Developing your own packages

At some point, you are going to need to develop your own packages to organize your code and distribute it if needed. As stated at the beginning of this chapter, everything that begins with an uppercase letter is considered public and can be accessed from outside its package, whereas all other elements are considered private. The only exception to this Go rule is package names—it is a best practice to use lowercase package names, even though uppercase package names are allowed.

Compiling a Go package can be done manually, if the package exists on the local machine, but it is also done automatically after you download the package from the internet, so there is no need to worry about it. Additionally, if the package you are downloading contains any errors, you will learn about them at downloading time.

However, if you want to compile a package that has been saved in the `post05.go` file (a combination of _PostgreSQL_ and _Chapter 05_) on your own, you can use the following command:

```markup
$ go build  -o post.a post05.go
```

So, the previous command compiles the `post05.go` file and saves its output in the `post.a` file:

```markup
$ file post.a
post.a: current ar archive
```

The `post.a` file is an `ar` archive.

The main reason for compiling Go packages on your own is to check for syntax or other kinds of errors in your code. Additionally, you can build Go packages as plugins ([https://golang.org/pkg/plugin/](https://golang.org/pkg/plugin/)) or shared libraries. Discussing more about these is beyond the scope of this book.

## The init() function

Each Go package can optionally have a private function named `init()` that is automatically executed at the beginning of execution time—`init()` runs when the package is initialized at the beginning of program execution. The `init()` function has the following characteristics:

-   `init()` takes no arguments.
-   `init()` returns no values.
-   The `init()` function is optional.
-   The `init()` function is called implicitly by Go.
-   You can have an `init()` function in the `main` package. In that case, `init()` is executed _before_ the `main()` function. In fact, all `init()` functions are always executed prior to the `main()` function.
-   A source file can contain multiple `init()` functions—these are executed in the order of declaration.
-   The `init()` function or functions of a package are executed only once, even if the package is imported multiple times.
-   Go packages can contain multiple files. Each source file can contain one or more `init()` functions.

The fact that the `init()` function is a private function by design means that it cannot be called from outside the package in which it is contained. Additionally, as the user of a package has no control over the `init()` function, you should think carefully before using an `init()` function in public packages or changing any global state in `init()`.

There are some exceptions where the use of `init()` makes sense:

-   For initializing network connections that might take time prior to the execution of package functions or methods.
-   For initializing connections to one or more servers prior to the execution of package functions or methods.
-   For creating required files and directories.
-   For checking whether required resources are available or not.

As the order of execution can be perplexing sometimes, in the next subsection, we will explain the order of execution in more detail.

## Order of execution

This subsection illustrates how Go code is executed. As an example, if a `main` package imports package `A` and package `A` depends on package `B`, then the following will take place:

-   The process starts with `main` package.
-   The `main` package imports package `A`.
-   Package `A` imports package `B`.
-   The global variables, if any, in package `B` are initialized.
-   The `init()` function or functions of package `B`, if they exist, run. This is the first `init()` function that gets executed.
-   The global variables, if any, in package `A` are initialized.
-   The `init()` function or functions of package `A`, if there are any, run.
-   The global variables in the `main` package are initialized.
-   The `init()` function or functions of `main` package, if they exist, run.
-   The `main()` function of the `main` package begins execution.

Notice that if the `main` package imports package `B` on its own, nothing is going to happen because everything related to package `B` is triggered by package `A`. This is because package `A` imports package `B` first.

The following diagram shows what is happening behind the scenes regarding the order of execution of Go code:

![Diagram, chat or text message
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_05_01.png)

Figure 5.1: Order of execution in Go

You can learn more about the order of execution in Go by reading the Go Language Specification document at [https://golang.org/ref/spec#Order\_of\_evaluation](https://golang.org/ref/spec#Order_of_evaluation) and about package initialization at [https://golang.org/ref/spec#Package\_initialization](https://golang.org/ref/spec#Package_initialization).

Just Imagine

# Using GitHub to store Go packages

This section will teach you how to create a GitHub repository where you can keep your Go package and make it available to the world.

First, you need to create the GitHub repository on your own. The easiest way to create a new GitHub repository is by visiting the GitHub website and going to the _Repositories_ tab, where you can see your existing repositories and create new ones. Press the **New** button and type in the necessary information for creating a new GitHub repository. If you made your repository public, everyone will be able to see it—if it is a private repository, only the people you choose are going to be able to look into it.

Having a clear `README.md` file in your GitHub repository that explains the way the Go package works is considered a very good practice.

Next, you need to clone the repository on your local computer. I usually clone it using the `git(1)` utility. As the name of the repository is `post05` and my GitHub username is `mactsouk`, the `git clone` command looks as follows:

```markup
$ git clone git@github.com:mactsouk/post05.git
```

Type `cd post05` and you are done! After that, you just have to write the code of the Go package and remember to `git commit` and `git push` the code to the GitHub repository.

The look of such a repository, after it has been used for a while, can be seen in _Figure 5.2_—you are going to learn more about the `post05` repository in a while:

![Graphical user interface, text, application, email, website
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_05_02.png)

Figure 5.2: A GitHub repository with a Go package

Using GitLab instead of GitHub for hosting your code does not require making any changes to the way you work.

If you want to use that package, you just need to `go get` the package using its URL and include it in your `import` block—we will see this when we actually use it in a program.

The next section presents a Go package that allows you to work with a database.

Just Imagine

# A package for working with a database

This section will develop a Go package for working with a given database schema stored on a Postgres database, with the end goal of demonstrating how to develop, store, and use a package. When interacting with specific schemas and tables in your application, you usually create separate packages with all the database-related functions—this also applies to NoSQL databases.

Go offers a generic package ([https://golang.org/pkg/database/sql/](https://golang.org/pkg/database/sql/)) for working with databases. However, each database requires a specific package that acts as the driver and allows Go to connect and work with this specific database.

The steps for creating the desired Go package are as follows:

-   Downloading the necessary external Go packages for working with PostgreSQL.
-   Creating package files.
-   Developing the required functions.
-   Using the Go package for developing utilities.
-   Using CI/CD tools for automation (this is optional).

You might be wondering why we would create such a package for working with a database and not write the actual commands in our programs when needed. The reasons for this include the following:

-   A Go package can be shared by all team members that work with the application.
-   A Go package allows people to use the database in ways that are documented.
-   The specialized functions you put in your Go package fit your needs a lot better.
-   People do not need full access to the database—they just use the package functions and the functionality they offer.
-   If you ever make changes to the database, people do not need to know about them, as long as the functions of the Go package remain the same.

Put simply, the functions you create can interact with a specific database schema, along with its tables and data—it would be almost impossible to work with an unknown database schema without knowing how the tables are connected to each other.

Apart from all these technical reasons, it is really fun to create Go packages that are shared among multiple developers!

Let's now continue by learning more about the database and its tables.

## Getting to know your database

You most likely need to download an additional package for working with a database server such as Postgres, MySQL, or MongoDB. In this case, we are using _PostgreSQL_ and therefore need to download a Go package that allows us to communicate with PostgreSQL. There are two main Go packages for connecting to PostgreSQL—we are going to use the `github.com/lib/pq` package here, but it is up to you to decide which package to use.

There is another Go package for working with PostgreSQL called `jackc/pgx` that can be found at [https://github.com/JackC/pgx](https://github.com/JackC/pgx).

You can download that package as follows:

```markup
$ go get github.com/lib/pq
```

To make things simpler, the PostgreSQL server is executed from a Docker image using a `docker-compose.yml` file, which has the following contents:

```markup
version: '3'
services:
  postgres:
    image: postgres
    container_name: postgres
    environment:
      - POSTGRES_USER=mtsouk
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=master
    volumes:
      - ./postgres:/var/lib/postgresql/data/
    networks:
      - psql
    ports:
      - "5432:5432"
volumes:
  postgres:
networks:
  psql:
    driver: bridge
```

The default port number the PostgreSQL server listens to is `5432`. As we connect to that PostgreSQL server from the same machine, the hostname that is going to be used is `localhost` or, if you prefer an IP address, `127.0.0.1`. If you are using a different PostgreSQL server, then you should change the connection details in the code that follows accordingly.

In PostgreSQL, a schema is a namespace that contains named database objects such as tables, views, and indexes. PostgreSQL automatically creates a schema called `public` for every new database.

The following Go utility, which is named `getSchema.go`, verifies that you can connect successfully to a PostgreSQL database and get a list of the available databases and tables in the given database and the `public` schema—all connection information is provided as command-line arguments:

```markup
package main
import (
    "database/sql"
    "fmt"
    "os"
    "strconv"
    _ "github.com/lib/pq"
)
```

The `lib/pq` package, which is the interface to the PostgreSQL database, is not used directly by the code. Therefore, you need to import the `lib/pq` package with `_` in order to prevent the Go compiler from creating an error message related to importing a package and not "using" it.

Most of the time, you do not need to import a package with `_`, but this is one of the exceptions. This kind of import is usually because the imported package has side effects, such as registering itself as the database handler for the `sql` package:

```markup
func main() {
    arguments := os.Args
    if len(arguments) != 6 {
        fmt.Println("Please provide: hostname port username password db")
        return
    }
```

Having a good help message for the information required by such a utility is very handy.

```markup
    host := arguments[1]
    p := arguments[2]
    user := arguments[3]
    pass := arguments[4]
    database := arguments[5]
```

This is where we collect the details of the database connection.

```markup
    // Port number SHOULD BE an integer
    port, err := strconv.Atoi(p)
    if err != nil {
        fmt.Println("Not a valid port number:", err)
        return
    }
    // connection string
    conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, database)
```

This is how you define the connection string with the details for the connection to the PostgreSQL database server. The connection string should be passed to the `sql.Open()` function for establishing the connection. So far, we have no connection.

```markup
    // open PostgreSQL database
    db, err := sql.Open("postgres", conn)
    if err != nil {
        fmt.Println("Open():", err)
        return
    }
    defer db.Close()
```

The `sql.Open()` function opens the database connection and keeps it open until the program ends, or until you execute `Close()` in order to properly close the database connection.

```markup
    // Get all databases
    rows, err := db.Query(`SELECT "datname" FROM "pg_database"
    WHERE datistemplate = false`)
    if err != nil {
        fmt.Println("Query", err)
        return
    }
```

In order to execute a `SELECT` query, you need to create it first. As the presented `SELECT` query contains no parameters, which means that it does not change based on variables, you can pass it to the `Query()` function and execute it. The _live_ outcome of the `SELECT` query is kept in the `rows` variable, which is a **cursor**. You do not get all the results from the database, as a query might return millions of records, but you get them one by one—this is the point of using a cursor.

```markup
    for rows.Next() {
        var name string
        err = rows.Scan(&name)
        if err != nil {
            fmt.Println("Scan", err)
            return
        }
        fmt.Println("*", name)
    }
    defer rows.Close()
```

The previous code shows how to process the results of a `SELECT` query, which can be from nothing to lots of rows. As the `rows` variable is a cursor, you advance from row to row by calling `Next()`. After that, you need to assign the values returned from the `SELECT` query into Go variables, in order to use them. This happens with a call to `Scan()`, which requires pointer parameters. If the `SELECT` query returns multiple values, you need to put multiple parameters in `Scan()`. Lastly, you must call `Close()` with `defer` for the `rows` variable in order to close the statement and free various types of used resources.

```markup
    // Get all tables from __current__ database
    query := `SELECT table_name FROM information_schema.tables WHERE 
        table_schema = 'public' ORDER BY table_name`
    rows, err = db.Query(query)
    if err != nil {
        fmt.Println("Query", err)
        return
    }
```

We are going to execute another `SELECT` query in the current database, as provided by the user. The definition of the `SELECT` query is kept in the `query` variable for simplicity and for creating easy to read code. The contents of the `query` variable are passed to the `db.Query()` method.

```markup
    // This is how you process the rows that are returned from SELECT
    for rows.Next() {
        var name string
        err = rows.Scan(&name)
        if err != nil {
            fmt.Println("Scan", err)
            return
        }
        fmt.Println("+T", name)
    }
    defer rows.Close()
}
```

Once again, we need to process the rows returned by the `SELECT` statement using the `rows` cursor and the `Next()` method.

Running `getSchema.go` generates the following kind of output:

```markup
$ go run getSchema.go localhost 5432 mtsouk pass go
* postgres
* master
* go
+T userdata
+T users
```

But what is the output telling us? Lines beginning with `*` show PostgreSQL databases, whereas lines beginning with `+T` show database tables—this is our decision. Therefore, this specific PostgreSQL installation contains three databases named `postgres`, `master`, and `go`. The `public` schema of the `go` database, which is specified by the last command-line argument, contains two tables named `userdata` and `users`.

The main advantage of the `getSchema.go` utility is that it is generic and can be used for learning more about PostgreSQL servers, which is the main reason that it requires so many command-line arguments to work.

Now that we know how to access and query a PostgreSQL database using Go, the next task should be creating a GitHub or GitLab repository for keeping and distributing the Go package we are about to develop.

## Storing the Go package

The first action we should take is creating a repository for storing the Go package. In our case, we are going to use a GitHub repository to keep the package. It is not a bad idea to keep the GitHub repository **private** during development, before exposing it to the rest of the world, especially when you are creating something critical.

Keeping the GitHub repository private does not affect the development process, but it might make sharing the Go package more difficult, so in some cases, it would be good to make it public.

For simplicity, we will use a public Go repository for the Go module, which is named `post05`—its full URL is [https://github.com/mactsouk/post05](https://github.com/mactsouk/post05).

In order to use that package on your machines, you should `go get` it first. However, during development, you should begin with `git clone git@github.com:mactsouk/post05.git` to get the contents of the GitHub repository and make changes to it.

## The design of the Go package

The following diagram shows the database schema that the Go package works on. Remember that when working with a specific database and schema, you need to "include" the schema information in your Go code. Put simply, the Go code should know about the schema it works on:

![Graphical user interface, application
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_05_03.png)

Figure 5.3: The two database tables the Go package works on

This is a simple schema that allows us to keep user data and update it. What connects the two tables is the **user ID**, which should be unique. Additionally, the `Username` field on the `Users` table should also be unique as two or more users cannot share the same username.

This schema already exists in the PostgreSQL database server, which means that the Go code assumes that the relevant tables are in the right place and are stored in the correct PostgreSQL database. Apart from the `Users` table, there is also a table named `Userdata` that holds information about each user. Once a record is entered in the `Users` table, it cannot be changed. What can change, however, is the data stored in the `Userdata` table.

If you want to create the `Users` database and the two tables in a database called `go`, you can execute the following statements, which are saved in a file named `create_tables.sql` with the `psgl` utility:

```markup
DROP DATABASE IF EXISTS go;
CREATE DATABASE go;
DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Userdata;
\c go;
CREATE TABLE Users (
    ID SERIAL,
    Username VARCHAR(100) PRIMARY KEY
);
CREATE TABLE Userdata (
    UserID Int NOT NULL,
    Name VARCHAR(100),
    Surname VARCHAR(100),
    Description VARCHAR(200)
);
```

The command-line utility for working with Postgres is called `psql`. The `psql` command for executing the code of `create_tables.sql` is as follows:

```markup
$ psql -h localhost -p 5432 -U mtsouk master < create_tables.sql
```

Now that we have the necessary infrastructure up and running, let's begin discussing the Go package. The tasks that the Go package should perform to make our lives easier are as follows:

-   Create a new user
-   Delete an existing user
-   Update an existing user
-   List all users

Each of these tasks should have one or more Go functions or methods to support it, which is what we are going to implement in the Go package:

-   A function to initiate the Postgres connection—the connection details should be given by the user and the package should be able to use them. However, the helper function to initiate the connection can be private.
-   There should exist default values in some of the connection details.
-   A function that checks whether a given username exists—this is a helper function that might be private.
-   A function that inserts a new user into the database.
-   A function that deletes an existing user from the database.
-   A function for updating an existing user.
-   A function for listing all users.

Now that we know the overall structure and functionality of the Go package, we should begin implementing it.

## The implementation of the Go package

In this subsection, we will implement the Go package for working with the Postgres database and the given database schema. We will present each function separately—if you combine all these functions, then you have the functionality of the entire package.

During package development, you should regularly commit your changes to the GitHub or GitLab repository.

The first element that you need in your Go package is one or more **structures** that can hold data from the database tables. Most of the times, you need as many structures as there are database tables—we will begin with that and see how it goes. Therefore, we will define the following structures:

```markup
type User struct {
    ID       int
    Username string
}
type Userdata struct {
    ID          int
    Name        string
    Surname     string
    Description string
}
```

If you think about this, you should see that there is no point in creating two separate Go structures. This is because the `User` structure holds no real data, and there is no point in passing multiple structures to the functions that process data for the `Users` and `Userdata` PostgreSQL tables. Therefore, we can create a single Go structure for holding all the data that has been defined, as follows:

```markup
type Userdata struct {
    ID          int
    Username    string
    Name        string
    Surname     string
    Description string
}
```

I have decided to name the structure after the database table for simplicity—however, in this case, this is not completely accurate as the `Userdata` structure has more fields than the `Userdata` database table. The thing is that we do not need everything from the `Userdata` database table.

The preamble of the package is as follows:

```markup
package post05
import (
    "database/sql"
    "errors"
    "fmt"
    "strings"
    _ "github.com/lib/pq"
)
```

For the first time in this book, you will see a package name different than `main`, which in this case is `post05`. As the package communicates with PostgreSQL, we import the `github.com/lib/pq` package and we use `_` in front of the package's path. As we discussed earlier, this happens because the imported package is registering itself as the database handler for the `sql` package, but it is not being directly used in the code. It is only being used through the `sql` package.

Next, you should have variables to hold the connection details. In the case of the `post05` package, this can be implemented with the following global variables:

```markup
// Connection details
var (
    Hostname = ""
    Port     = 2345
    Username = ""
    Password = ""
    Database = ""
)
```

Apart from the `Port` variable, which has an initial value, the other global variables have the default value of their data type, which is `string`. All these variables must be properly initialized by the Go code that uses the `post05` package and should be accessible from outside the package, which means that their first letter should be in uppercase.

The `openConnection()` function, which is **private** and only accessed within the scope of the package, is defined as:

```markup
func openConnection() (*sql.DB, error) {
    // connection string
    conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)
    // open database
    db, err := sql.Open("postgres", conn)
    if err != nil {
        return nil, err
    }
    return db, nil
}
```

You have already seen the previous code in the `getSchema.go` utility. You create the connection string and you pass it to `sql.Open()`.

Now, let's consider the `exists()` function, which is also private:

```markup
// The function returns the User ID of the username
// -1 if the user does not exist
func exists(username string) int {
    username = strings.ToLower(username)
    db, err := openConnection()
    if err != nil {
        fmt.Println(err)
        return -1
    }
    defer db.Close()
    userID := -1
    statement := fmt.Sprintf(`SELECT "id" FROM "users" where username = '%s'`, username)
    rows, err := db.Query(statement)
```

This is where we define the query that shows whether the provided username exists in the database or not. As all our data is kept in the database, we need to interact with the database all the time.

```markup
    for rows.Next() {
        var id int
        err = rows.Scan(&id)
        if err != nil {
            fmt.Println("Scan", err)
            return -1
        }
```

If the `rows.Scan(&id)` call is executed without any errors, then we know that a result has been returned, which is the desired user ID.

```markup
        userID = id
    }
    defer rows.Close()
    return userID
}
```

The last part of `exists()` sets or closes the query to free resources and returns the ID value of the username that is given as a parameter to `exists()`.

During development, I include many `fmt.Println()` statements in the package code for debugging purposes. However, I have removed most of them in the final version of the Go package and replaced them with `error` values. These `error` values are passed to the program that uses the functionality of the package, which is responsible for deciding what to do with the error messages and error conditions. You can also use logging for this—the output can go to standard output or even `/dev/null` when not needed.

Here is the implementation of the `AddUser()` function:

```markup
// AddUser adds a new user to the database
// Returns new User ID
// -1 if there was an error
func AddUser(d Userdata) int {
    d.Username = strings.ToLower(d.Username)
```

All usernames are converted into lowercase in order to avoid duplicates. This is a design decision.

```markup
    db, err := openConnection()
    if err != nil {
        fmt.Println(err)
        return -1
    }
    defer db.Close()
    userID := exists(d.Username)
    if userID != -1 {
        fmt.Println("User already exists:", Username)
        return -1
    }
    insertStatement := `insert into "users" ("username") values ($1)`
```

This is how we construct a query that accepts parameters. The presented query requires one value that is named `$1`.

```markup
    _, err = db.Exec(insertStatement, d.Username)
```

This is how you pass the desired value, which is `d.Username`, into the `insertStatement` variable.

```markup
    if err != nil {
        fmt.Println(err)
        return -1
    }
    userID = exists(d.Username)
    if userID == -1 {
        return userID
    }
    insertStatement = `insert into "userdata" ("userid", "name", "surname", "description") values ($1, $2, $3, $4)`
```

The presented query needs 4 values that are named `$1`, `$2`, `$3`, and `$4`.

```markup
    _, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)
    if err != nil {
        fmt.Println("db.Exec()", err)
        return -1
    }
```

As we need to pass 4 variables to `insertStatement`, we will put 4 values in the `db.Exec()` call.

```markup
    return userID
}
```

This is the end of the function that adds a new user to the database. The implementation of the `DeleteUser()` function is as follows.

```markup
// DeleteUser deletes an existing user
func DeleteUser(id int) error {
    db, err := openConnection()
    if err != nil {
        return err
    }
    defer db.Close()
    // Does the ID exist?
    statement := fmt.Sprintf(`SELECT "username" FROM "users" where id = %d`, id)
    rows, err := db.Query(statement)
```

Here, we double-check whether the given user ID exists or not in the `users` table.

```markup
    var username string
    for rows.Next() {
        err = rows.Scan(&username)
        if err != nil {
            return err
        }
    }
    defer rows.Close()
    if exists(username) != id {
        return fmt.Errorf("User with ID %d does not exist", id)
    }
```

If the previously returned username exists and has the same user ID as the parameter to `DeleteUser()`, then you can continue the deletion process, which contains two steps: first, deleting the relevant user data from the `userdata` table, and, second, deleting the data from the `users` table.

```markup
    // Delete from Userdata
    deleteStatement := `delete from "userdata" where userid=$1`
    _, err = db.Exec(deleteStatement, id)
    if err != nil {
        return err
    }
    // Delete from Users
    deleteStatement = `delete from "users" where id=$1`
    _, err = db.Exec(deleteStatement, id)
    if err != nil {
        return err
    }
    return nil
}
```

Now, let's examine the implementation of the `ListUsers()` function.

```markup
func ListUsers() ([]Userdata, error) {
    Data := []Userdata{}
    db, err := openConnection()
    if err != nil {
        return Data, err
    }
    defer db.Close()
```

Once again, we need to open a connection to the database before executing a database query.

```markup
    rows, err := db.Query(`SELECT  
        "id","username","name","surname","description"
        FROM "users","userdata"
        WHERE users.id = userdata.userid`)
    if err != nil {
        return Data, err
    }
    for rows.Next() {
        var id int
        var username string
        var name string
        var surname string
        var description string
        err = rows.Scan(&id, &username, &name, &surname, &description)
        temp := Userdata{ID: id, Username: username, Name: name, Surname: surname, Description: description}
```

At this point, we will store the data we've received from the `SELECT` query into a `Userdata` structure. This is added to the slice that is going to be returned from the `ListUsers()` function. This process continues until there is nothing left to read.

```markup
        Data = append(Data, temp)
        if err != nil {
            return Data, err
        }
    }
    defer rows.Close()
    return Data, nil
}
```

After updating the contents of `Data` using `append()`, we end the query, and the function returns the list of available users, as stored in `Data`.

Lastly, let's examine the `UpdateUser()` function:

```markup
// UpdateUser is for updating an existing user
func UpdateUser(d Userdata) error {
    db, err := openConnection()
    if err != nil {
        return err
    }
    defer db.Close()
    userID := exists(d.Username)
    if userID == -1 {
        return errors.New("User does not exist")
    }
```

First, we need to make sure that the given username exists in the database—the update process is based on the username.

```markup
    d.ID = userID
    updateStatement := `update "userdata" set "name"=$1, "surname"=$2, "description"=$3 where "userid"=$4`
    _, err = db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
    if err != nil {
        return err
    }
    return nil
}
```

The update statement stored in `updateStatement` that is executed using the desired parameters with the help of `db.Exec()` updates the user data.

Now that we know the details of how to implement each function in the `post05` package, it is time to begin using that package!

## Testing the Go package

In order to test the package, we must develop a command-line utility called `postGo.go`.

As `postGo.go` uses an external package, even if we develop that package, you should not forget to download the latest version of that external package using `go get` or `go get -u`.

As `postGo.go` is used for testing purposes only, we hardcoded most of the data apart from the username of the user we put into the database. All usernames are randomly generated.

The code of `postGo.go` is as follows:

```markup
package main
import (
    "fmt"
    "math/rand"
    "time"
    "github.com/mactsouk/post05"
)
```

As the `post05` package works with Postgres, there is no need to import `lib/pq` here:

```markup
var MIN = 0
var MAX = 26
func random(min, max int) int {
    return rand.Intn(max-min) + min
}
func getString(length int64) string {
    startChar := "A"
    temp := ""
    var i int64 = 1
    for {
        myRand := random(MIN, MAX)
        newChar := string(startChar[0] + byte(myRand))
        temp = temp + newChar
        if i == length {
            break
        }
        i++
    }
    return temp
}
```

Both the `random()` and `getString()` functions are helper functions for generating random strings that are used as usernames.

```markup
func main() {
    post05.Hostname = "localhost"
    post05.Port = 5432
    post05.Username = "mtsouk"
    post05.Password = "pass"
    post05.Database = "go"
```

This is where you define the connection parameters to the Postgres server, as well as the database you are going to work in (`go`). As all these variables are in the `post05` package, they are accessed as such.

```markup
    data, err := post05.ListUsers()
    if err != nil {
        fmt.Println(err)
        return
    }
    for _, v := range data {
        fmt.Println(v)
    }
```

We begin by listing existing users.

```markup
    SEED := time.Now().Unix()
    rand.Seed(SEED)
    random_username := getString(5)
```

Then, we generate a random string that is used as the username. All randomly generated usernames are 5 characters long because of the `getString(5)` call. You can change that value if you want.

```markup
    t := post05.Userdata{
        Username:    random_username,
        Name:        "Mihalis",
        Surname:     "Tsoukalos",
        Description: "This is me!"}
    id := post05.AddUser(t)
    if id == -1 {
        fmt.Println("There was an error adding user", t.Username)
    }
```

The previous code adds a new user to the database—the user data, including the username, is kept in a `post05.Userdata` structure. That `post05.Userdata` structure is passed to the `post05.AddUser()` function, which returns the user ID of the new user.

```markup
    err = post05.DeleteUser(id)
    if err != nil {
        fmt.Println(err)
    }
```

Here, we delete the user that we created using the user ID value returned by `post05.AddUser(t)`.

```markup
    // Trying to delete it again!
    err = post05.DeleteUser(id)
    if err != nil {
        fmt.Println(err)
    }
```

If you try to delete the same user again, the process fails because the user does not exist.

```markup
    id = post05.AddUser(t)
    if id == -1 {
        fmt.Println("There was an error adding user", t.Username)
    }
```

Here, we add the same user again—however, as user ID values are generated by Postgres, this time, the user is going to have a different user ID value than before.

```markup
    t = post05.Userdata{
        Username:    random_username,
        Name:        "Mihalis",
        Surname:     "Tsoukalos",
        Description: "This might not be me!"}
```

Here, we update the `Description` field of the `post05.Userdata` structure before passing it to `post05.UpdateUser()`, in order update the information stored in the database.

```markup
    err = post05.UpdateUser(t)
    if err != nil {
        fmt.Println(err)
    }
}
```

Working with `postGo.go` creates the following kind of output:

```markup
$ go run postGo.go
{4 mhmxz Mihalis Tsoukalos This might not be me!}
{6 wsdlg Mihalis Tsoukalos This might not be me!}
User with ID 7 does not exist
```

The previous output confirms that `postGo.go` works as expected as it can connect to the database, add a new user, and delete an existing one. This also means that the `post05` package works as expected. Now that we know how to create Go packages, let's briefly discuss Go modules.

Just Imagine

# Modules

A Go module is like a Go package with a version—however, Go modules can consist of multiple packages. Go uses **semantic versioning** for versioning modules. This means that versions begin with the letter `v`, followed by the `major.minor.patch` version numbers. Therefore, you can have versions such as v1.0.0, v1.0.5, and v2.0.2. The `v1`, `v2`, and `v3` parts signify the major version of a Go package that is usually not backward compatible. This means that if your Go program works with v1, it will not necessarily work with v2 or v3—it might work, but you cannot count on it. The second number in a version is about features. Usually, v1.1.0 has more features than v1.0.2 or v1.0.0, while being compatible with all older versions. Lastly, the third number is just about bug fixes without having any new features. Note that semantic versioning is also used for Go versions.

Go modules were introduced in Go v1.11 but were finalized in Go v1.13.

If you want to learn more about modules, visit and read [https://blog.golang.org/using-go-modules](https://blog.golang.org/using-go-modules), which has five parts, as well as [https://golang.org/doc/modules/developing](https://golang.org/doc/modules/developing). Just remember that **a Go module is similar but not identical to a regular Go package with a version**, and that a module can consist of multiple packages.

Just Imagine

# Creating better packages

This section provides handy advice that can help you develop better Go packages. Here are several good rules to follow to create high-class Go packages:

-   The first unofficial rule of a successful package is that its elements must be connected in some way. Thus, you can create a package for supporting cars, but it would not be a good idea to create a single package for supporting cars and bicycles and airplanes. Put simply, it is better to split the functionality of a package unnecessarily into multiple packages than to add too much functionality to a single Go package.
-   A second practical rule is that you should use your own packages first for a reasonable amount of time before giving them to the public. This helps you discover silly bugs and make sure that your packages operate as expected. After that, give them to some fellow developers for additional testing before making them publicly available. Additionally, you should always write tests for any package you intend others to use.
-   Next, make sure your package has a clear and useful API so that any consumer can be productive with it quickly.
-   Try and limit the public API of your packages to only what is absolutely necessary. Additionally, give your functions descriptive but not very long names.
-   Interfaces, and in future Go versions, **generics**, can improve the usefulness of your functions, so when you think it is appropriate, use an interface instead of a single type as a function parameter or return type.
-   When updating one of your packages, try not to break things and create incompatibilities with older versions unless it is absolutely necessary.
-   When developing a new Go package, try to use multiple files in order to group similar tasks or concepts.
-   Do not create a package that already exists from scratch. Make changes to the existing package and maybe create your own version of it.
-   Nobody wants a Go package that prints logging information on the screen. It would be more professional to have a flag for turning on logging when needed. The Go code of your packages should be in harmony with the Go code of your programs. This means that if you look at a program that uses your packages and your function names stand out in the code in a bad way, it would be better to change the names of your functions. As the name of a package is used almost everywhere, try to use concise and expressive package names.
-   It is more convenient if you put new Go type definitions near where they are used the first time because nobody, including yourself, wants to search source files for definitions of new data types.
-   Try to create test files for your packages, because packages with test files are considered more professional than ones without them; small details make all the difference and give people confidence that you are a serious developer! Notice that writing tests for your packages is not optional and that you should avoid using packages that do not include tests. You will learn more about testing in _Chapter 11_, _Code Testing and Profiling_.

Always remember that apart from the fact that the actual Go code in a package should be bug-free, the next most important element of a successful package is its documentation, as well as some code examples that clarify its use and showcase the idiosyncrasies of the functions of the package. The next section discusses creating documentation in Go.

Just Imagine

# Generating documentation

This section discusses how to create **documentation** for your Go code using the code of the `post05` package as an example. The new package is renamed and is now called `document`.

Go follows a simple rule regarding documentation: in order to document a function, a method, a variable, or even the package itself, you can write comments, as usual, that should be located directly before the element you want to document, without any empty lines in between. You can use one or more single-line comments, which are lines beginning with `//`, or _block_ comments, which begin with `/*` and end with `*/`—everything in-between is considered a comment.

It is highly recommended that each Go package you create has a block comment preceding the `package` declaration that introduces developers to the package, and also explains what the package does.

Instead of presenting the entire code of the `post05` package, we will only present the important part, which means that function implementations are going to be `return` statements only. The new version of `post05.go` is called `document.go` and comes with the following code and comments:

```markup
/*
The package works on 2 tables on a PostgreSQL data base server.
The names of the tables are:
    * Users
    * Userdata
The definitions of the tables in the PostgreSQL server are:
    CREATE TABLE Users (
        ID SERIAL,
        Username VARCHAR(100) PRIMARY KEY
    );
    CREATE TABLE Userdata (
        UserID Int NOT NULL,
        Name VARCHAR(100),
        Surname VARCHAR(100),
        Description VARCHAR(200)
    );
    This is rendered as code
This is not rendered as code
*/
package document
```

This is the first block of documentation that is located right before the name of the package. This is the appropriate place to document the functionality of the package, as well as other essential information. In this case, we are presenting the SQL create commands that fully describe the database tables we are going to work on. Another important element is specifying the database server this package interacts with. Other information that you can put at the beginning of a package is the author, the license, and the version of the package.

If a line in a block comment begins with a tab, then it is rendered differently in the graphical output, which is good for differentiating between various kinds of information in the documentation:

```markup
// BUG(1): Function ListUsers() not working as expected
// BUG(2): Function AddUser() is too slow
```

The `BUG` keyword is special when writing documentation. Go knows that bugs are part of the code and therefore should be documented as well. You can write any message you want after a `BUG` keyword, and you can place them anywhere you want—preferably close to the bugs they describe.

```markup
import (
    "database/sql"
    "fmt"
    "strings"
)
```

The `github.com/lib/pq` package was removed from the `import` block to make the file size smaller.

```markup
/*
This block of global variables holds the connection details to the Postgres server
    Hostname: is the IP or the hostname of the server
    Port: is the TCP port the DB server listens to
    Username: is the username of the database user
    Password: is the password of the database user
    Database: is the name of the Database in PostgreSQL
*/
var (
    Hostname = ""
    Port     = 2345
    Username = ""
    Password = ""
    Database = ""
)
```

The previous code shows a way of documenting lots of variables at once—in this case, global variables. The good thing with this way is that you do not have to put a comment before each global variable and make the code less readable. The only downside of this method is that you should remember to update the comments, should you wish to make any changes to the code. However, documenting multiple variables at once might not end up rendering correctly in web-based `godoc` pages. For that reason, you might want to document each field directly.

```markup
// The Userdata structure is for holding full user data
// from the Userdata table and the Username from the
// Users table
type Userdata struct {
    ID          int
    Username    string
    Name        string
    Surname     string
    Description string
}
```

The previous excerpt shows how to document a Go structure—this is especially useful when you have lots of structures in a source file and you want to have a quick look at them.

```markup
// openConnection() is for opening the Postgres connection
// in order to be used by the other functions of the package.
func openConnection() (*sql.DB, error) {
    var db *sql.DB
    return db, nil
}
```

When documenting a function, it is good to begin the first line of the comments with the function name. Apart from that, you can write any information that you consider important in the comments.

```markup
// The function returns the User ID of the username
// -1 if the user does not exist
func exists(username string) int {
    fmt.Println("Searching user", username)
    return 0
}
```

In this case, we will explain the return values of the `exists()` function as they have a special meaning.

```markup
// AddUser adds a new user to the database
//
// Returns new User ID
// -1 if there was an error
func AddUser(d Userdata) int {
    d.Username = strings.ToLower(d.Username)
    return -1
}
/*
    DeleteUser deletes an existing user if the user exists.
    It requires the User ID of the user to be deleted.
*/
func DeleteUser(id int) error {
    fmt.Println(id)
    return nil
}
```

You can use block comments anywhere you want, not only at the beginning of a package.

```markup
// ListUsers lists all users in the database
// and returns a slice of Userdata.
func ListUsers() ([]Userdata, error) {
    // Data holds the records returned by the SQL query
    Data := []Userdata{}
    return Data, nil
}
```

When you request the documentation of the `Userdata` structure, Go automatically presents the functions that use `Userdata` as input or output, or both.

```markup
// UpdateUser is for updating an existing user
// given a Userdata structure.
// The user ID of the user to be updated is found
// inside the function.
func UpdateUser(d Userdata) error {
    fmt.Println(d)
    return nil
}
```

We are not done yet because we need to see the documentation somehow. There are two ways to see the documentation of the package. The first one involves using `go get`, which also means creating a GitHub repository of the package, as we did with `post05`. However, as this is for testing purposes, we are going to do things the easy way: we are going to copy it in `~/go/src` and access it from there. As the package is called `document`, we are going to create a directory with the same name inside `~/go/src`. After that, we are going to copy `document.go` in `~/go/src/document` and we are done—for more complex packages, the process is going to be more complex as well. In such cases, it would be better to `go get` the package from its repository.

Either way, the `go doc` command is going to work just fine with the `document` package:

```markup
$ go doc document
package document // import "document"
The package works on 2 tables on a PostgreSQL data base server.
The names of the tables are:
    * Users
    * Userdata
The definitions of the tables in the PostgreSQL server are:
        CREATE TABLE Users (
            ID SERIAL,
            Username VARCHAR(100) PRIMARY KEY
        );
        CREATE TABLE Userdata (
            UserID Int NOT NULL,
            Name VARCHAR(100),
            Surname VARCHAR(100),
            Description VARCHAR(200)
        );
        This is rendered as code
This is not rendered as code
var Hostname = "" ...
func AddUser(d Userdata) int
func DeleteUser(id int) error
func UpdateUser(d Userdata) error
type Userdata struct{ ... }
    func ListUsers() ([]Userdata, error)
BUG: Function ListUsers() not working as expected
BUG: Function AddUser() is too slow
```

If you want to see information about a specific function, you should use `go doc`, as follows:

```markup
$ go doc document ListUsers
package document // import "document"
func ListUsers() ([]Userdata, error)
    ListUsers lists all users in the database and returns a slice of Userdata.
```

Additionally, we can use the web version of the Go documentation, which can be accessed after running the `godoc` utility and going to the `Third Party` section—by default, the web server initiated by `godoc` listens to port number `6060` and can be accessed at `http://localhost:6060`.

A part of the documentation page for the `document` package is shown in _Figure 5.4_:

![Graphical user interface, application
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_05_04.png)

Figure 5.4: Viewing the documentation of a user-developed Go package

Go automatically puts the bugs at the end of the text and graphical output.

In my personal opinion, rendering the documentation is much better when using the graphical interface, making it better when you do not know what you are looking for. On the other hand, using `go doc` from the command line is much faster and allows you to process the output using traditional UNIX command-line tools.

The next section will briefly present the CI/CD systems of GitLab and GitHub, starting with GitLab Runners, which can be helpful for automating package development and deployment.

Just Imagine

# GitLab Runners and Go

When developing Go packages and documentation, you want to be able to test the results and find out bugs as quickly as possible. When everything works as expected, you might want to publish your results to the world automatically, without spending more time on this. One of the best solutions for this is using a CI/CD system for automating tasks. This section will briefly illustrate the use of GitLab Runners for automating Go projects.

In order to follow this section, you need to have a GitLab account, create a dedicated GitLab repository, and store the relevant files there.

We will begin with a GitLab repository that contains the following files:

-   `hw.go`: This is a sample program that is used to make sure that everything works.
-   `.gitignore`: It is not necessary to have such a file, but it is very handy for ignoring some files and directories.
-   `usePost05.go`: This is a sample Go file that uses an external package—please refer to the [https://gitlab.com/mactsouk/runners-go/](https://gitlab.com/mactsouk/runners-go/) repository for its contents.
-   `README.md`: This file is automatically displayed on the repository web page and is usually used for explaining the purpose of the repository.

There is also a directory called `.git` that contains information and metadata about the repository.

## The initial version of the configuration file

The first version of the configuration file is for making sure that everything is fine with our GitLab setup. The name of the configuration file is `.gitlab-ci.yml` and is a YAML file that should be located in the root directory of the GitLab repository. This initial version of the `.gitlab-ci.yml` configuration file compiles `hw.go` and creates a binary file, which is executed in a different stage than the one it was created in. This means that we should create an artifact for keeping and transferring that binary file:

```markup
$ cat .gitlab-ci.yml
image: golang:1.15.7
stages:
    - download
    - execute
compile:
    stage: download
    script:
        - echo "Getting System Info"
        - uname -a
        - mkdir bin
        - go version
        - go build -o ./bin/hw hw.go
    artifacts:
        paths:
            - bin/
execute:
    stage: execute
    script:
        - echo "Executing Hello World!"
        - ls -l bin
        - ./bin/hw
```

The important thing about the previous configuration file is that we are using an image that comes with Go already installed, which saves us from having to install it from scratch and allows us to specify the Go version we want to use.

However, if you ever want to install extra software, you can do that based on the Linux distribution being used. After saving the file, you need to push the changes to GitLab for the pipeline to start running. In order to see the results, you should click on the **CI/CD** option on the left bar of the GitLab UI.

_Figure 5.5_ shows information about a specific workflow based on the aforementioned YAML file. Everything that is in green is a good thing, whereas red is used in error situations. If you want to learn more information about a specific stage, you can press its button and see a more detailed output:

![Graphical user interface, text, application
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_05_05.png)

Figure 5.5: Looking at the progress of a GitLab pipeline

As everything looks fine, we are ready to create the final version of `.gitlab-ci.yml`. Note that if there is an error in the workflow, you will most likely receive an email the email address you used when registering for GitLab. If everything is fine, no email will be sent.

## The final version of the configuration file

The final version of the CI/CD configuration file compiles `usePost05.go`, which imports `post05`. This is used for illustrating how external packages are downloaded. The contents of `.gitlab-ci.yml` is as follows:

```markup
image: golang:1.15.7
stages:
    - download
    - execute
compile:
    stage: download
    script:
        - echo "Compiling usePost05.go"
        - mkdir bin
        - go get -v -d ./...
```

The `go get -v -d ./...` command is the Go way of downloading all the package dependencies of a project. After that, you are free to build your project and generate your executable file:

```markup
        - go build -o ./bin/usePost05 usePost05.go
    artifacts:
        paths:
            - bin/
```

The `bin` directory, along with its contents, will be available to the `execute` state:

```markup
execute:
    stage: execute
    script:
        - echo "Executing usePost05"
        - ls -l bin
        - ./bin/usePost05
```

Pushing it to GitLab automatically triggers its execution. This is illustrated in _Figure 5.6_, which shows in more detail the progress of the `compile` stage:

![Text
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_05_06.png)

Figure 5.6: Seeing a detailed view of a stage

What we can see here is that all the required packages are being downloaded and that `usePost05.go` is compiled without any issues. As we do not have a PostgreSQL instance available, we cannot try interacting with PostgreSQL, but we can execute `usePost05.go` and see the values of the `Hostname` and `Port` global variables. This is illustrated in _Figure 5.7_:

![Text
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_05_07.png)

Figure 5.7: Seeing more details about the execute stage

So far, we have seen how to use GitLab Runners to automate Go package development and testing. Next, we are going to create a CI/CD scenario in GitHub using GitHub Actions as another method of automating software publication.

Just Imagine

# GitHub Actions and Go

This section will use GitHub Actions to push a Docker image that contains a Go executable file in Docker Hub.

In order to follow this section, you must have a GitHub account, create a dedicated GitHub repository, and store the related files there.

We will begin with a GitHub repository that contains the following files:

-   `.gitignore`: This is an optional file that's used for ignoring files and directories during `git push` operations.
-   `usePost05.go`: This is the same file as before.
-   `Dockerfile`: This file is used for creating a Docker image with the Go executable. Please refer to [https://github.com/mactsouk/actions-with-go](https://github.com/mactsouk/actions-with-go) for its contents.
-   `README.md`: As before, this is a Markdown file that contains information about the repository.

In order to set up GitHub Actions, we need to create a directory named `.github` and then create another directory named `workflows` in it. The `.github/workflows` directory contains YAML files with the pipeline configuration.

_Figure 5.8_ shows the overview screen of the workflows of the selected GitHub repository:

![Graphical user interface, text, application, email
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_05_08.png)

Figure 5.8: Displaying the Workflows associated with a given GitHub repository

To push an image to Docker Hub, you need to log in. As this process requires using a password, which is sensitive information, the next subsection illustrates how to store secrets in GitHub.

## Storing secrets in GitHub

The credentials for connecting to Docker Hub are stored in GitHub using the secrets feature that exists in almost all CI/CD systems. However, the exact implementation might differ.

You can also use HashiCorp Vault as a central point for storing passwords and other sensitive data. Unfortunately, presenting HashiCorp Vault is beyond the scope of this book.

In your GitHub repository, go to the **Settings** tab and select **Secrets** from the left column. You will see your existing secrets, if any, and an **Add new secret** link, which you need to click on. Do this process twice to store your Docker Hub username and password.

_Figure 5.9_ shows the secrets associated with the GitHub repository used in this section—the presented secrets hold the username and password that were used for connecting to Docker Hub:

![A picture containing chart
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_05_09.png)

Figure 5.9: The secrets of the mactsouk/actions-with-go repository

## The final version of the configuration file

The final version of the configuration file compiles the Go code, puts it in a Docker image, as described by `Dockerfile`, connects with Docker Hub using the specified credentials, and pushes the Docker image to Docker Hub using the provided data. This is a very common way of automation when creating Docker images. The contents of `go.yml` is as follows:

```markup
name: Go + PostgreSQL
on: [push]
```

This line in the configuration file specifies that this pipeline is triggered on push operations only.

```markup
jobs:
  build:
    runs-on: ubuntu-18.04
```

This is the Linux image that is going to be used. As the Go binary is built inside a Docker image, we do not need to install Go on the Linux VM.

```markup
    steps:
    # Checks-out your repository under $GITHUB_WORKSPACE, so your job can access it
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        stable: 'false'
        go-version: '1.15.7'
    - name: Publish Docker Image
      env:
         USERNAME: ${{ secrets.USERNAME }}
         PASSWORD: ${{ secrets.PASSWORD }}
```

This is how you access the secrets that were stored previously and how you store them.

```markup
         IMAGE_NAME: gopost
      run: |
        docker images
        docker build -t "$IMAGE_NAME" .
        docker images
        echo "$PASSWORD" | docker login --username "$USERNAME" --password-stdin
        docker tag "${IMAGE_NAME}" "$USERNAME/${IMAGE_NAME}:latest"
        docker push "$USERNAME/${IMAGE_NAME}:latest"
        echo "* Running Docker Image"
        docker run ${IMAGE_NAME}:latest
```

This time, most of the work is performed by the `docker build` command because the Go executable is built inside a Docker image. The following screenshot shows some of the output of the pipeline defined by `go.yml`:

![Text
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_05_10.png)

Figure 5.10: The pipeline for pushing Docker images to Docker Hub

Automation saves you development time, so try to automate as many things as possible, especially when you're making your software available to the world.

Just Imagine

# Versioning utilities

One of the most difficult tasks is to automatically and uniquely version command-line utilities, especially when using a CI/CD system. This section presents a technique that uses a GitHub value to version a command-line utility on your local machine.

You can apply the same technique to GitLab—just search for the available GitLab variables and values and choose one that fits your needs.

This technique is used by both the `docker` and `kubectl` utilities, among others:

```markup
$ docker version
Client: Docker Engine - Community
 Cloud integration: 1.0.4
 Version:           20.10.0
 API version:       1.41
 Go version:        go1.13.15
 Git commit:        7287ab3
 Built:             Tue Dec  8 18:55:43 2020
 OS/Arch:           darwin/amd64
...
```

The previous output shows that `docker` uses the Git commit value for versioning—we are going to use a slightly different value that is longer than the one used by `docker`.

The utility that is used, which is saved as `gitVersion.go`, is implemented as follows:

```markup
package main
import (
    "fmt"
    "os"
)
var VERSION string
```

`VERSION` is the variable that is going to be set at runtime using the Go linker.

```markup
func main() {
    if len(os.Args) == 2 {
        if os.Args[1] == "version" {
            fmt.Println("Version:", VERSION)
        }
    }
}
```

The previous code says that if there is a command-line argument and its value is `version`, print the version message with the help of the `VERSION` variable.

What we need to do is tell the Go linker that we are going to define the value of the `VERSION` variable. This happens with the help of the `-ldflags` flag, which stands for linker flags—this passes values to the `cmd/link` package, which allows us to change values in imported packages at build time. The `-X` value that is used requires a key/value pair, where the key is a variable name and the value is the value that we want to set for that key. In our case, the key has the `main.Variable` form because we change the value of a variable in the `main` package. As the name of the variable in `gitVersion.go` is `VERSION`, the key is `main.VERSION`.

But first, we need to decide on the GitHub value that we are going to use as the version string. The `git rev-list HEAD` command returns a full list of commits for the current repository from the latest to the oldest. We only need the last one—the most recent—which we can get using `git rev-list -1 HEAD` or `git rev-list HEAD | head -1`. So, we need to assign that value to an environment variable and pass that environment variable to the Go compiler. As this value changes each time you make a commit and you always want to have the latest value, you should reevaluate it each time you execute `go build`—this will be shown in a while.

In order to provide `gitVersion.go` with the value of the desired environment variable, we should execute it as follows:

```markup
$ export VERSION=$(git rev-list -1 HEAD)
$ go build -ldflags "-X main.VERSION=$VERSION" gitVersion.go
```

This works on both `bash` and `zsh` shells. If you are using a different shell, you should make sure that you are defining an environment variable the right way.

If you want to execute the two commands at the same time, you can do the following:

```markup
$ export VERSION=$(git rev-list -1 HEAD) && go build -ldflags "-X main.VERSION=$VERSION" gitVersion.go
```

Running the generated executable, which is called `gitVersion`, produces the next output:

```markup
$ ./gitVersion version
Version: 99745c8fbaff94790b5818edf55f15d309f5bfeb
```

Your output is going to be different because your GitHub repository is going to be different. As GitHub generates random and unique values, you won't have the same version number twice!

Just Imagine

# Exercises

1.  Can you write a function that sorts three `int` values? Try to write two versions of the function: one with named returned values and another without named return values. Which one do you think is better?
2.   Rewrite the `getSchema.go` utility so that it works with the `jackc/pgx` package.
3.  Rewrite the `getSchema.go` utility so that it works with MySQL databases.
4.  Use GitLab CI/CD to push Docker images to Docker Hub.

Exercise 1: Here's an example of a function that sorts three integer values using two different approaches: one with named return values and another without named return values.

```
// Function with named return values
func sortThreeValues(a, b, c int) (sortedA, sortedB, sortedC int) {
    if a <= b && a <= c {
        sortedA = a
        if b <= c {
            sortedB = b
            sortedC = c
        } else {
            sortedB = c
            sortedC = b
        }
    } else if b <= a && b <= c {
        sortedA = b
        if a <= c {
            sortedB = a
            sortedC = c
        } else {
            sortedB = c
            sortedC = a
        }
    } else {
        sortedA = c
        if a <= b {
            sortedB = a
            sortedC = b
        } else {
            sortedB = b
            sortedC = a
        }
    }
    return
}

// Function without named return values
func sortThreeValuesAlt(a, b, c int) (int, int, int) {
    if a <= b && a <= c {
        if b <= c {
            return a, b, c
        }
        return a, c, b
    } else if b <= a && b <= c {
        if a <= c {
            return b, a, c
        }
        return b, c, a
    }
    if a <= b {
        return c, a, b
    }
    return c, b, a
}
```

Both versions of the function achieve the same result, which is sorting the three input integers in ascending order. The difference lies in how the return values are specified.

The first version uses named return values (sortedA, sortedB, and sortedC). This allows us to assign the sorted values directly to these named variables within the function body. The function ends with a bare return statement, which automatically returns the named values.

The second version does not use named return values. Instead, it directly returns the sorted values in the appropriate order using the return statement.

Regarding which version is better, it depends on personal preference and the specific use case. Named return values can make the code more readable by providing clear documentation of the returned values. It can also help improve code maintainability if the number of return values is large or if the function signature changes in the future. However, the second version without named return values can be more concise and straightforward.

Ultimately, it's up to you to decide which approach best suits your needs based on factors such as code readability, maintainability, and personal coding style.

Exercise 2: Here's an example of how you can modify the getSchema.go utility to work with the jackc/pgx package in Go:

```
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

func main() {
	connConfig, err := pgx.ParseConfig("postgres://username:password@localhost:5432/database_name")
	if err != nil {
		log.Fatal("Error parsing connection config:", err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer conn.Close(context.Background())

	schemaName := "public" // Change this to the desired schema name

	tables, err := getTables(conn, schemaName)
	if err != nil {
		log.Fatal("Error retrieving tables:", err)
	}

	for _, table := range tables {
		columns, err := getColumns(conn, schemaName, table)
		if err != nil {
			log.Println("Error retrieving columns for table", table, ":", err)
			continue
		}

		fmt.Println("Table:", table)
		for _, column := range columns {
			fmt.Println("- Column:", column)
		}
		fmt.Println()
	}
}

func getTables(conn *pgx.Conn, schemaName string) ([]string, error) {
	var tables []string
	rows, err := conn.Query(context.Background(), `SELECT table_name FROM information_schema.tables WHERE table_schema = $1`, schemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return tables, nil
}

func getColumns(conn *pgx.Conn, schemaName, tableName string) ([]string, error) {
	var columns []string
	rows, err := conn.Query(context.Background(), `SELECT column_name FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2`, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName string
		err := rows.Scan(&columnName)
		if err != nil {
			return nil, err
		}
		columns = append(columns, columnName)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return columns, nil
}
```

In this modified version, the database/sql package is replaced with the github.com/jackc/pgx/v4 package, which provides the pgx PostgreSQL driver. The connection and query operations are updated to use the pgx package's functions and types.

Make sure to replace the username, password, localhost, and database_name placeholders in the connection string with your actual database credentials.

The getTables and getColumns functions are updated to use the pgx.Conn type for executing queries and retrieving results. The queries themselves remain the same, as they are written in standard SQL.

The rest of the code remains mostly unchanged, and it will retrieve the tables and columns from the specified schema using the pgx package.

Please note that you need to have the pgx package installed. You can install it by running the following command:

```
go get github.com/jackc/pgx/v4
```

Remember to update the connection string and adjust the code as per your specific requirements.

Exercise 3:

 Here's an example of how you can modify the getSchema.go utility to work with MySQL databases in Go:

```
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database_name")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	schemaName := "public" // Change this to the desired schema name

	tables, err := getTables(db, schemaName)
	if err != nil {
		log.Fatal("Error retrieving tables:", err)
	}

	for _, table := range tables {
		columns, err := getColumns(db, schemaName, table)
		if err != nil {
			log.Println("Error retrieving columns for table", table, ":", err)
			continue
		}

		fmt.Println("Table:", table)
		for _, column := range columns {
			fmt.Println("- Column:", column)
		}
		fmt.Println()
	}
}

func getTables(db *sql.DB, schemaName string) ([]string, error) {
	var tables []string
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return tables, nil
}

func getColumns(db *sql.DB, schemaName, tableName string) ([]string, error) {
	var columns []string
	rows, err := db.Query("SHOW COLUMNS FROM " + tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName, columnType string
		err := rows.Scan(&columnName, &columnType)
		if err != nil {
			return nil, err
		}
		columns = append(columns, columnName+" ("+columnType+")")
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return columns, nil
}
```

In this modified version, we import the github.com/go-sql-driver/mysql package, which provides the MySQL driver for Go's database/sql package.

Make sure to replace the username, password, localhost, database_name, and any other required connection details in the connection string.

The getTables function uses the SHOW TABLES query to retrieve the list of tables in the database.

The getColumns function uses the SHOW COLUMNS FROM query to retrieve the columns of a specific table.

The rest of the code remains mostly unchanged, and it will retrieve the tables and columns from the specified schema using the MySQL driver.

Please note that you need to have the go-sql-driver/mysql package installed. You can install it by running the following command:

```
go get github.com/go-sql-driver/mysql
```

Remember to update the connection string and adjust the code as per your specific MySQL configuration.

Exercise 4:  To use GitLab CI/CD to push Docker images to Docker Hub, you need to set up a pipeline in your GitLab project that builds and pushes the Docker image whenever changes are made to the code. 

Here's a step-by-step guide to help you get started:

- Set up your GitLab repository:
- Create a new repository on GitLab or navigate to an existing repository.

- Make sure your codebase contains a Dockerfile that specifies the image you want to build.

- Set up Docker Hub repository:

- Create a Docker Hub account if you don't have one.

- Create a new repository on Docker Hub where you want to push your Docker images.

- Configure GitLab CI/CD:

In your GitLab repository, navigate to the CI/CD settings (Settings -> CI/CD).
Add the following environment variables under "Variables":
        
- DOCKER_USERNAME: Your Docker Hub username.
- DOCKER_PASSWORD: Your Docker Hub password or access token.

Create a GitLab CI/CD pipeline configuration file:

In the root of your GitLab repository, create a file named .gitlab-ci.yml.

Add the following content to configure the pipeline:

```
    image: docker:stable

    services:
      - docker:dind

    stages:
      - build

    build:
      stage: build
      script:
        - docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
        - docker build -t $DOCKER_USERNAME/your-image-name .
        - docker push $DOCKER_USERNAME/your-image-name
```

- Replace your-image-name with the name you want to give to your Docker image.

- Commit and push the .gitlab-ci.yml file to your GitLab repository. This will trigger the pipeline.

- GitLab CI/CD will automatically execute the pipeline:

It will pull the docker:stable image to use as the base image. It will spin up a Docker service. It will execute the build stage defined in the pipeline configuration file. The script will log in to Docker Hub using the provided credentials. It will build the Docker image using the Dockerfile in your repository. It will push the built image to your Docker Hub repository.

Monitor the pipeline execution:
        
- Go to your GitLab repository and navigate to the "CI/CD" section.
- You will be able to see the progress and status of the pipeline. Once the pipeline completes successfully, your Docker image should be available in your Docker Hub repository.

That's it! You have set up GitLab CI/CD to automatically build and push Docker images to Docker Hub whenever changes are made to your codebase.

Note: Make sure to keep your Docker Hub credentials secure. Storing them as environment variables in the GitLab CI/CD settings ensures they are not exposed in your code repository.

Remember to customize the pipeline configuration and script as per your specific requirements, such as specifying the Docker image name, tags, and additional build steps if needed.

Just Imagine

# Summary

This chapter presented two primary topics: functions and packages. Functions are first-class citizens in Go, which makes them powerful and handy. Remember that everything that begins with an uppercase letter is public. The only exception to this rule is package names. Private variables, functions, data type names, and structure fields can be strictly used and called internally in a package, whereas public ones are available to everyone. Additionally, we learned more about the `defer` keyword. Also, memorize that Go packages are not like Java classes—a Go package can be as big as it needs to be. Regarding Go modules, keep in mind that a Go module is multiple packages with a version.

Finally, this chapter discussed creating documentation, GitHub Actions and GitLab Runners, how the two CI/CD systems can help you automate boring processes and how to assign unique version numbers to your utilities.

The next chapter discusses system programming in general, as well as file I/O in more detail.

Just Imagine

# Additional resources

-   New module changes in Go 1.16: [https://blog.golang.org/go116-module-changes](https://blog.golang.org/go116-module-changes)
-   How do you structure your Go apps? Talk by Kat Zien from GopherCon UK 2018: [https://www.youtube.com/watch?v=1rxDzs0zgcE](https://www.youtube.com/watch?v=1rxDzs0zgcE)
-   PostgreSQL: [https://www.postgresql.org/](https://www.postgresql.org/)
-   PostgreSQL Go package: [https://github.com/lib/pq](https://github.com/lib/pq)
-   PostgreSQL Go package: [https://github.com/jackc/pgx](https://github.com/jackc/pgx)
-   HashiCorp Vault: [https://www.vaultproject.io/](https://www.vaultproject.io/)
-   The documentation of database/sql: [https://golang.org/pkg/database/sql/](https://golang.org/pkg/database/sql/%20)
-   You can learn more about GitHub Actions environment variables at [https://docs.github.com/en/actions/reference/environment-variables](https://docs.github.com/en/actions/reference/environment-variables)
-   GitLab CI/CD variables: [https://docs.gitlab.com/ee/ci/variables/](https://docs.gitlab.com/ee/ci/variables/)
-   The documentation of the `cmd/link` package: [https://golang.org/cmd/link/](https://golang.org/cmd/link/)
-   [golang.org](http://golang.org) moving to [go.dev](http://go.dev): [https://go.dev/blog/tidy-web](https://go.dev/blog/tidy-web)