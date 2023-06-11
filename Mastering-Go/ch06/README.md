# Telling a UNIX System What to Do

This chapter teaches you about **systems programming** in Go. Systems programming involves working with files and directories, process control, signal handling, network programming, system files, configuration files, and file input and output (I/O). If you recall from _Chapter 1_, _A Quick Introduction to Go_, the reason for writing system utilities with Linux in mind is that often Go software is executed in a Docker environment—Docker images use the Linux operating system, which means that you might need to **develop your utilities with the Linux operating system in mind**. However, as Go code is portable, most system utilities work on Windows machines without any changes or with minor modifications. Among other things, this chapter implements two utilities, one that finds cycles in UNIX file systems and another that converts JSON data to XML data and vice versa. Additionally, in this chapter we are going to improve the phone book application with the help of the `cobra` package.

**Important note**: Starting with Go 1.16, the `GO111MODULE` environment variable defaults to `on`—this affects the use of Go packages that do not belong to the Go standard library. In practice, this means that you must put your code under `~/go/src`. You can go back to the previous behavior by setting `GO111MODULE` to `auto`, but you do not want to do that—modules are the future. The reason for mentioning this in this chapter is that both `viper` and `cobra` prefer to be treated as Go modules instead of packages, which changes the development process but not the code.

This chapter covers:

-   `stdin`, `stdout`, and `stderr`
-   UNIX processes
-   Handling UNIX signals
-   File input and output
-   Reading plain text files
-   Writing to a file
-   Working with JSON
-   Working with XML
-   Working with YAML
-   The `viper` package
-   The `cobra` package
-   Finding cycles in a UNIX file system
-   New to Go 1.16
-   Updating the phone book application

Bookmark

# stdin, stdout, and stderr

Every UNIX operating system has three files open all the time for its processes. Remember that UNIX considers everything, even a printer or your mouse, as a file. UNIX uses **file descriptors**, which are positive integer values, as an internal representation for accessing open files, which is much prettier than using long paths. So, by default, all UNIX systems support three special and standard filenames: `/dev/stdin`, `/dev/stdout`, and `/dev/stderr`, which can also be accessed using file descriptors `0`, `1`, and `2`, respectively. These three file descriptors are also called **standard input**, **standard output**, and **standard error**, respectively. Additionally, file descriptor `0` can be accessed as `/dev/fd/0` on a macOS machine and as both `/dev/fd/0` and `/dev/pts/0` on a Debian Linux machine.

Go uses `os.Stdin` for accessing standard input, `os.Stdout` for accessing standard output, and `os.Stderr` for accessing standard error. Although you can still use `/dev/stdin`, `/dev/stdout`, and `/dev/stderr` or the related file descriptor values for accessing the same devices, it is better, safer, and more portable to stick with `os.Stdin`, `os.Stdout`, and `os.Stderr`.

Bookmark

# UNIX processes

As Go servers, utilities, and Docker images are mainly executed on Linux, it is good to know about Linux processes and threads.

Strictly speaking, a **process** is an execution environment that contains instructions, user data and system data parts, and other types of resources that are obtained during runtime. On the other hand, a **program** is a binary file that contains instructions and data that are used for initializing the instruction and user data parts of a process. Each running UNIX process is uniquely identified by an unsigned integer, which is called the **process ID** of the process.

There are three process categories: **user processes**, **daemon processes**, and **kernel processes**. User processes run in user space and usually have no special access rights. Daemon processes are programs that can be found in the user space and run in the background without the need for a terminal. Kernel processes are executed in kernel space only and can fully access all kernel data structures.

The C way of creating new processes involves the calling of the `fork(2)` system call. The return value of `fork(2)` allows the programmer to differentiate between a parent and a child process. Although you can fork a new process in Go using the `exec` package, Go does not allow you to control threads—Go offers **goroutines**, which the user can create on top of threads that are created and handled by the Go runtime.

Bookmark

# Handling UNIX signals

UNIX signals offer a very handy way of _interacting asynchronously with your applications_. However, UNIX signal handling requires the use of Go channels that are used exclusively for this task. So, it would be good to talk a little about the concurrency model of Go, which requires the use of goroutines and channels for signal handling.

A **goroutine** is the smallest executable Go entity. In order to create a new goroutine you **have to** use the `go` keyword followed by a predefined function or an anonymous function—the methods are equivalent. A **channel** in Go is a mechanism that among other things allows goroutines to communicate and exchange data. If you are an amateur programmer or are hearing about goroutines and channels for the first time, do not panic. Goroutines and channels are explained in much more detail in _Chapter 7_, _Go Concurrency_.

In order for a goroutine or a function to terminate the entire Go application, it should call `os.Exit()` instead of `return`. However, most of the time, you should exit a goroutine or a function using `return` because you just want to exit that specific goroutine or function and not stop the entire application.

The presented program handles `SIGINT`, which is called `syscall.SIGINT` in Go, and `SIGINFO` separately and uses a `default` case in a `switch` block for handling the remaining cases (other signals). The implementation of that `switch` block allows you to differentiate between the various signals according to your needs.

There exists a _dedicated channel_ that receives all signals, as defined by the `signal.Notify()` function. Go channels can have a capacity—the capacity of this particular channel is `1` in order to be able to receive and keep one signal at a time. This makes perfect sense as a signal can terminate a program and there is no need to try to handle another signal at the same time. There is usually an anonymous function that is executed as a goroutine and performs the signal handling and nothing else. The main task of that goroutine is to listen to the channel for data. Once a signal is received, it is sent to that channel, read by the goroutine, and stored into a variable—at this point the channel can receive more signals. That variable is processed by a `switch` statement.

Some signals cannot be caught, and the operating system cannot ignore them. So, the `SIGKILL` and `SIGSTOP` signals cannot be blocked, caught, or ignored and the reason for this is that they allow privileged users as well as the UNIX kernel to terminate any process they desire.

Create a text file by typing the following code—a good filename for it would be `signals.go`.

```markup
package main
import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)
func handleSignal(sig os.Signal) {
    fmt.Println("handleSignal() Caught:", sig)
}
```

`handleSignal()` is a separate function for handling signals. However, you can also handle signals **inline**, in the branches of a `switch` statement.

```markup
func main() {
    fmt.Printf("Process ID: %d\n", os.Getpid())
    sigs := make(chan os.Signal, 1)
```

We create a channel with data of type `os.Signal` because all channels must have a type.

```markup
    signal.Notify(sigs)
```

The previous statement means _handle all signals that can be handled_.

```markup
    start := time.Now()
    go func() {
        for {
            sig := <-sigs
```

Wait until you read data (`<-`) from the `sigs` channel and store it in the `sig` variable.

```markup
            switch sig {
```

Depending on the read value, act accordingly. This is how you differentiate between signals.

```markup
            case syscall.SIGINT:
                duration := time.Since(start)
                fmt.Println("Execution time:", duration)
```

For the handling of `syscall.SIGINT`, we calculate the time that has passed since the beginning of the program execution, and print it on screen.

```markup
            case syscall.SIGINFO:
                handleSignal(sig)
```

The code of the `syscall.SIGINFO` case calls the `handleSignal()` function—it is up to the developer to decide on the details of the implementation. On Linux machines, you should replace `syscall.SIGINFO` with another signal such as `syscall.SIGUSR1` or `syscall.SIGUSR2` because `syscall.SIGINFO` **is not available on Linux** ([https://github.com/golang/go/issues/1653](https://github.com/golang/go/issues/1653)).

```markup
                // do not use return here because the goroutine exits
                // but the time.Sleep() will continue to work!
                os.Exit(0)
            default:
                fmt.Println("Caught:", sig)
            }
```

If there is not a match, the `default` case handles the rest of the values and just prints a message.

```markup
        }
    }()
    for {
        fmt.Print("+")
        time.Sleep(10 * time.Second)
    }
}
```

The endless `for` loop at the end of the `main()` function is for emulating the operation of a real program. Without an endless `for` loop, the program exits almost immediately.

Running `signals.go` and interacting with it creates the following kind of output:

```markup
$ go run signals.go
Process ID: 74252
+Execution time: 9.989863093s
+Caught: user defined signal 1
+signal: killed
```

The second line of output was generated by pressing Ctrl + C on the keyboard, which on UNIX machines sends the `syscall.SIGINT` signal to the program. The third line of output was caused by executing `kill -USR1 74252` on a different terminal. The last line in the output was generated by the `kill -9 74252` command. As the `KILL` signal, which is also represented by the number `9`, cannot be handled, it terminates the program, and the shell prints the `killed` message.

## Handling two signals

If you want to handle a limited number of signals, instead of all of them, you should replace the `signal.Notify(sigs)` statement with the next statement:

```markup
signal.Notify(sigs, syscall.SIGINT, syscall.SIGINFO)
```

After that you need to make the appropriate changes to the code of the goroutine responsible for signal handling in order to identify and handle `syscall.SIGINT` and `syscall.SIGINFO`—the current version (`signals.go`) already handles both of them.

Now, we need to learn how to read and write files in Go.

Bookmark

# File I/O

This section discusses file I/O in Go, which includes the use of the `io.Reader` and `io.Writer` interfaces, buffered and unbuffered I/O, as well as the `bufio` package.

The `io/ioutil` package ([https://golang.org/pkg/io/ioutil/](https://golang.org/pkg/io/ioutil/)) is deprecated in Go version 1.16. Existing Go code that uses the functionality of `io/ioutil` will continue to work but it is better to stop using that package.

## The io.Reader and io.Writer interfaces

This subsection presents the definitions of the popular `io.Reader` and `io.Writer` interfaces because these two interfaces are the basis of file I/O in Go—the former allows you to read from a file whereas the latter allows you to write to a file. The definition of the `io.Reader` interface is the following:

```markup
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

This definition, which should be revisited when we want one of our data types to satisfy the `io.Reader` interface, tells us the following:

-   The `Reader` interface requires the implementation of a single method
-   The parameter of `Read()` is a byte slice
-   The return values of `Read()` are an integer and an error

The `Read()` method takes a byte slice as input, which is going to be filled with data **up to its length**, and returns the number of bytes read as well as an `error` variable.

The definition of the `io.Writer` interface is the following:

```markup
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

The previous definition, which should be revisited when we want one of our data types to satisfy the `io.Writer` interface and to write to a file, reveals the following information:

-   The interface requires the implementation of a single method
-   The parameter of `Write()` is a byte slice
-   The return values of `Write()` are an integer and an `error` value

The `Write()` method takes a byte slice, which contains the data that you want to write, as input and returns the number of bytes written and an `error` variable.

## Using and misusing io.Reader and io.Writer

The code that follows showcases the use of `io.Reader` and `io.Writer` for custom data types, which in this case are two Go structures named `S1` and `S2`.

For the `S1` structure, the presented code implements both interfaces in order to read user data from the terminal and print data to the terminal, respectively. Although this is redundant as we already have `fmt.Scanln()` and `fmt.Printf()`, it is a good exercise that shows how versatile and flexible both interfaces are. In a different situation, you could have used `io.Writer` for writing to a log service, or keeping a second backup copy of the written data, or anything else that fits your needs. However, this is also an example of interfaces allowing you to do crazy or, if you prefer, unusual things—it is up to the developer to create the desired functionality using the appropriate Go concepts and features!

The `Read()` method is using `fmt.Scanln()` to get user input from the terminal whereas the `Write()` method is printing the contents of its buffer parameter as many times as the value of the `F1` field of the structure using `fmt.Printf()`!

For the `S2` structure, the presented code implements the `io.Reader` interface only in the traditional way. The `Read()` method reads the `text` field of the `S2` structure, which is a byte slice. When there is nothing left to read, the `Read()` method returns the expected `io.EOF` error, which in reality is not an error but an expected situation. Along with the `Read()` method there exist two helper methods, named `eof()`, which declares that there is nothing more to read, and `readByte()`, which reads the `text` field of the `S2` structure byte by byte. After the `Read()` method is done, the `text` field of the `S2` structure, which is used as a buffer, is emptied.

With this implementation, the `io.Reader` for `S2` can be used for reading in a traditional way, which in this case is with `bufio.NewReader()` and multiple `Read()` calls—the number of `Read()` calls depends on the size of the buffer that is used, which in this case is a byte slice with 2 places for data.

Type the following code and save it as `ioInterface.go`:

```markup
package main
import (
    "bufio"
    "fmt"
    "io"
)
```

The previous part shows that we are using the `io` and `bufio` packages for working with files.

```markup
type S1 struct {
    F1 int
    F2 string
}
type S2 struct {
    F1   S1
    text []byte
}
```

These are the two structures we are going to work with.

```markup
// Using pointer to S1 for changes to be persistent when the method exits
func (s *S1) Read(p []byte) (n int, err error) {
    fmt.Print("Give me your name: ")
    fmt.Scanln(&p)
    s.F2 = string(p)
    return len(p), nil
}
```

In the preceding code, we are implementing the `io.Reader()` interface for `S1`.

```markup
func (s *S1) Write(p []byte) (n int, err error) {
    if s.F1 < 0 {
        return -1, nil
    }
    for i := 0; i < s.F1; i++ {
        fmt.Printf("%s ", p)
    }
    fmt.Println()
    return s.F1, nil
}
```

The previous method implements the `io.Writer` interface for `S1`.

```markup
func (s S2) eof() bool {
    return len(s.text) == 0
}
func (s *S2) readByte() byte {
    // this function assumes that eof() check was done before
    temp := s.text[0]
    s.text = s.text[1:]
    return temp
}
```

The previous function is an implementation of `bytes.Buffer.ReadByte` from the standard library.

```markup
func (s *S2) Read(p []byte) (n int, err error) {
    if s.eof() {
        err = io.EOF
        return
    }
    l := len(p)
    if l > 0 {
        for n < l {
```

The previous functions read from the given buffer until the buffer is empty. This is where we implement `io.Reader` for `S2`.

```markup
            p[n] = s.readByte()
            n++
            if s.eof() {
                s.text = s.text[0:0]
                break
            }
        }
    }
    return
}
```

When all data is read, the relevant structure field is emptied. The previous method implements `io.Reader` for `S2`. However, the operation of `Read()` is supported by `eof()` and `readByte()`, which are also user-defined.

Recall that Go allows you to name the return values of a function—in that case, a `return` statement without any additional arguments automatically returns the current value of each named return variable in the order they appear in the function signature. The `Read()` method uses that feature.

```markup
func main() {
    s1var := S1{4, "Hello"}
    fmt.Println(s1var)
```

We initialize an `S1` variable that is named `s1var`.

```markup
    buf := make([]byte, 2)
    _, err := s1var.Read(buf)
```

The previous line is reading for the `s1var` variable using a buffer with `2` bytes.

```markup
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Read:", s1var.F2)
    _, _ = s1var.Write([]byte("Hello There!"))
```

In the previous line, we call the `Write()` method for `s1var` in order to write the contents of a byte slice.

```markup
    s2var := S2{F1: s1var, text: []byte("Hello world!!")}
```

In the previous code, we initialize an `S2` variable that is named `s2var`.

```markup
    // Read s2var.text
    r := bufio.NewReader(&s2var)
```

We now create a reader for `s2var`.

```markup
    for {
        n, err := r.Read(buf)
        if err == io.EOF {
            break
```

We keep reading from `s2var` until there is an `io.EOF` condition.

```markup
        } else if err != nil {
            fmt.Println("*", err)
            break
        }
        fmt.Println("**", n, string(buf[:n]))
    }
}
```

Running `ioInterface.go` produces the next output:

```markup
$ go run ioInterface.go
{4 Hello}
```

The first line of the output shows the contents of the `s1var` variable.

```markup
Give me your name: Mike
```

Calling the `Read()` method of the `s1var` variable.

```markup
Read: Mike
Hello There! Hello There! Hello There! Hello There!
```

The previous line is the output of `s1var.Write([]byte("Hello There!"))`.

```markup
** 2 He
** 2 ll
** 2 o 
** 2 wo
** 2 rl
** 2 d!
** 1 !
```

The last part of the output illustrates the reading process using a buffer with a size of 2. The next section discusses buffered and unbuffered operations.

## Buffered and unbuffered file I/O

Buffered file I/O happens when there is a buffer for temporarily storing data before reading data or writing data. Thus, instead of reading a file byte by byte, you read many bytes at once. You put the data in a buffer and wait for someone to read it in the desired way.

Unbuffered file I/O happens when there is no buffer to temporarily store data before actually reading or writing it—this can affect the performance of your programs.

The next question that you might ask is how to decide when to use buffered and when to use unbuffered file I/O. When dealing with critical data, unbuffered file I/O is generally a better choice because buffered reads might result in out-of-date data and buffered writes might result in data loss when the power of your computer is interrupted. However, most of the time, there is no definitive answer to that question. This means that you can use whatever makes your tasks easier to implement. However, keep in mind that **buffered readers can also improve performance** by reducing the number of system calls needed to read from a file or socket, so there can be a real performance impact on what the programmer decides to use.

There is also the `bufio` package. As the name suggests, `bufio` is about buffered I/O. Internally, the `bufio` package implements the `io.Reader` and `io.Writer` interfaces, which it wraps in order to create the `bufio.Reader` and `bufio.Writer` objects, respectively. The `bufio` package is very popular for working with plain text files and you are going to see it in action in the next section.

Bookmark

# Reading text files

In this section you will learn how to read plain text files, as well as using the `/dev/random` UNIX device, which offers you a way of getting random numbers.

## Reading a text file line by line

The function for reading a file line by line is found in `byLine.go` and is named `lineByLine()`. The technique for reading a text file line by line is also used when reading a plain text file word by word as well as when reading a plain text file character by character because you usually process plain text files line by line. The presented utility prints every line that it reads, which makes it a simplified version of the `cat(1)` utility.

First, you create a new reader to the desired file using a call to `bufio.NewReader()`. Then, you use that reader with `bufio.ReadString()` in order to read the input file line by line. The trick is done by the parameter of `bufio.ReadString()`, which is a character that tells `bufio.ReadString()` to keep reading until that character is found. Constantly calling `bufio.ReadString()` when that parameter is the newline character (`\n`) results in reading the input file line by line.

The implementation of `lineByLine()` is as follows:

```markup
func lineByLine(file string) error {
    f, err := os.Open(file)
    if err != nil {
        return err
    }
    defer f.Close()
    r := bufio.NewReader(f)
```

After making sure that you can open the given file for reading (`os.Open()`), you create a new reader using `bufio.NewReader()`.

```markup
    for {
        line, err := r.ReadString('\n')
```

`bufio.ReadString()` returns two values: the string that was read and an `error` variable.

```markup
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Printf("error reading file %s", err)
            break
        }
        fmt.Print(line)
```

The use of `fmt.Print()` instead of `fmt.Println()` for printing the input line shows that the newline character is included in each input line.

```markup
    }
    return nil
}
```

Running `byLine.go` generates the following kind of output:

```markup
$ go run byLine.go ~/csv.data
Dimitris,Tsoukalos,2101112223,1600665563
Mihalis,Tsoukalos,2109416471,1600665563
Jane,Doe,0800123456,1608559903
```

The previous output shows the contents of `~/csv.data` presented line by line with the help of `byLine.go`. The next subsection shows how to read a plain text file word by word.

## Reading a text file word by word

Reading a plain text file word by word is the single most useful function that you want to perform on a file because you usually want to process a file on a per-word basis—it is illustrated in this subsection using the code found in `byWord.go`. The desired functionality is implemented in the `wordByWord()` function. The `wordByWord()` function uses **regular expressions** to separate the words found in each line of the input file. The regular expression defined in the `regexp.MustCompile("[^\\s]+")` statement states that we use whitespace characters to separate one word from another.

The implementation of the `wordByWord()` function is as follows:

```markup
func wordByWord(file string) error {
    f, err := os.Open(file)
    if err != nil {
        return err
    }
    defer f.Close()
    r := bufio.NewReader(f)
    for {
        line, err := r.ReadString('\n')
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Printf("error reading file %s", err)
            return err
        }
        r := regexp.MustCompile("[^\\s]+")
```

This is the place where you define the regular expression you want to use.

```markup
        words := r.FindAllString(line, -1)
```

This is where you apply the regular expression to split the `line` variable into fields.

```markup
        for i := 0; i < len(words); i++ {
            fmt.Println(words[i])
        }
```

This `for` loop just prints the fields of the `words` slice. If you want to know the number of words found in the input line, you can just find the value of the `len(words)` call.

```markup
    }
    return nil
}
```

Running `byWord.go` produces the following kind of output:

```markup
$ go run byWord.go ~/csv.data
Dimitris,Tsoukalos,2101112223,1600665563
Mihalis,Tsoukalos,2109416471,1600665563
Jane,Doe,0800123456,1608559903
```

As `~/csv.data` does not contain any whitespace characters, each line is considered a single word!

## Reading a text file character by character

In this subsection, you learn how to read a text file character by character, which is a pretty rare requirement unless you want to develop a text editor. You take each line that you read and split it using a `for` loop with `range`, which returns two values. You discard the first, which is the location of the current character in the line variable, and you use the second. However, that value is a rune, which means that you have to convert it into a character using `string()`.

The implementation of `charByChar()` is as follows:

```markup
func charByChar(file string) error {
    f, err := os.Open(file)
    if err != nil {
        return err
    }
    defer f.Close()
    r := bufio.NewReader(f)
    for {
        line, err := r.ReadString('\n')
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Printf("error reading file %s", err)
            return err
        }
        for _, x := range line {
            fmt.Println(string(x))
        }
```

Note that, due to the `fmt.Println(string(x))` statement, each character is printed in a distinct line, which means that the output of the program is going to be large. If you want a more compressed output, you should use the `fmt.Print()` function instead.

```markup
    }
    return nil
}
```

Running `byCharacter.go` and filtering it with `head(1)`, without any parameters, produces the following kind of output:

```markup
$ go run byCharacter.go ~/csv.data | head
D
...
,
T
```

The use of the `head(1)` utility without any parameters limits the output to just 10 lines.

The next section is about reading from `/dev/random`, which is a UNIX system file.

## Reading from /dev/random

In this subsection, you learn how to read from the `/dev/random` system device. The purpose of the `/dev/random` system device is to generate random data, which you might use for testing your programs or, in this case, as the seed for a random number generator. Getting data from `/dev/random` can be a little bit tricky, and this is the main reason for specifically discussing it here.

The code of `devRandom.go` is the following:

```markup
package main
import (
    "encoding/binary"
    "fmt"
    "os"
)
```

You need `encoding/binary` because you are reading binary data from `/dev/random` that you convert into an integer value.

```markup
func main() {
    f, err := os.Open("/dev/random")
    defer f.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
    var seed int64
    binary.Read(f, binary.LittleEndian, &seed)
    fmt.Println("Seed:", seed)
}
```

There are two representations named **little endian** and **big endian** that have to do with the **byte order** in the internal representation. In our case, we are using little endian. The _endian-ness_ has to do with the way different computing systems order multiple bytes of information.

A real-world example of endian-ness is how different languages read text in different ways: European languages tend to be read from left to right, whereas Arabic texts are read from right to left.

In a big endian representation, bytes are read from left to right, while little endian reads bytes from right to left. For the `0x01234567` value, which requires 4 bytes for storing, the big endian representation is `01 | 23 | 45 | 67` whereas the little endian representation is `67 | 45 | 23 | 01`.

Running `devRandom.go` creates the following kind of output:

```markup
$ go run devRandom.go
Seed: 422907465220227415
```

This means that the `/dev/random` device is a good place for getting random data including a seed value for your random number generator.

## Reading a specific amount of data from a file

This subsection teaches you how to read a specific amount of data from a file. The presented utility can come in handy when you want to see a small part of a file. The numeric value that is given as a command-line argument specifies the size of the buffer that is going to be used for reading. The most important code of `readSize.go` is the implementation of the `readSize()` function:

```markup
func readSize(f *os.File, size int) []byte {
    buffer := make([]byte, size)
    n, err := f.Read(buffer)
```

All the magic happens in the definition of the `buffer` variable because this is where we define the maximum amount of data that we want to read. Therefore, each time you invoke `readSize()`, the function is going to read from `f` at most `size` characters.

```markup
    // io.EOF is a special case and is treated as such
    if err == io.EOF {
        return nil
    }
    if err != nil {
        fmt.Println(err)
        return nil
    }
    return buffer[0:n]
}
```

The remaining code is about error conditions; `io.EOF` is a special and expected condition that should be treated separately and return the read characters as a byte slice to the caller function.

Running `readSize.go` produces the following kind of output:

```markup
$ go run readSize.go 12 readSize.go
package main
```

In this case, we read 12 characters from `readSize.go` itself because of the `12` parameter.

Now that we know how to read files, it is time to learn how to write to files.

Bookmark

# Writing to a file

So far, we have seen ways to read files. This subsection shows how to write data to files in four different ways and how to append data to an existing file. The code of `writeFile.go` is as follows:

```markup
package main
import (
    "bufio"
    "fmt"
    "io"
    "os"
)
func main() {
    buffer := []byte("Data to write\n")
    f1, err := os.Create("/tmp/f1.txt")
```

`os.Create()` returns an `*os.File` value associated with the file path that is passed as a parameter. Note that if the file already exists, `os.Create()` truncates it.

```markup
    if err != nil {
        fmt.Println("Cannot create file", err)
        return
    }
    defer f1.Close()
    fmt.Fprintf(f1, string(buffer))
```

The `fmt.Fprintf()` function, which requires a `string` variable, helps you write data to your own files using the format you want. The only requirement is having an `io.Writer` to write to. In this case, a valid `*os.File` variable, which satisfies the `io.Writer` interface, does the job.

```markup
    f2, err := os.Create("/tmp/f2.txt")
    if err != nil {
        fmt.Println("Cannot create file", err)
        return
    }
    defer f2.Close()
    n, err := f2.WriteString(string(buffer))
```

The `os.WriteString()` method writes the contents of a string to a valid `*os.File` variable.

```markup
    fmt.Printf("wrote %d bytes\n", n)
    f3, err := os.Create("/tmp/f3.txt")
```

Here we create a temporary file on our own. Later on in this chapter you are going to learn about using `os.CreateTemp()` for creating temporary files.

```markup
    if err != nil {
        fmt.Println(err)
        return
    }
    w := bufio.NewWriter(f3)
```

This function returns a `bufio.Writer`, which satisfies the `io.Writer` interface.

```markup
    n, err = w.WriteString(string(buffer))
    fmt.Printf("wrote %d bytes\n", n)
    w.Flush()
    f := "/tmp/f4.txt"
    f4, err := os.Create(f)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f4.Close()
    for i := 0; i < 5; i++ {
        n, err = io.WriteString(f4, string(buffer))
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Printf("wrote %d bytes\n", n)
    }
    // Append to a file
    f4, err = os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
```

`os.OpenFile()` provides a _better way_ to create or open a file for writing. `os.O_APPEND` is saying that if the file already exists, you should append to it instead of truncating it. `os.O_CREATE` is saying that if the file does not already exist, it should be created. Last, `os.O_WRONLY` is saying that the program should open the file for writing only.

```markup
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f4.Close()
    // Write() needs a byte slice
    n, err = f4.Write([]byte("Put some more data at the end.\n"))
```

The `Write()` method gets its input from a byte slice, which is the Go way of writing. All previous techniques used strings, which is not the best way, especially when working with binary data. However, using strings instead of byte slices is more practical as it is more convenient to manipulate `string` values than the elements of a byte slice, especially when working with Unicode characters. On the other hand, using `string` values increases allocation and can cause a lot of garbage collection pressure.

```markup
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("wrote %d bytes\n", n)
}
```

Running `writeFile.go` generates some information output about the bytes written on disk. What is really interesting is seeing the files created in the `/tmp` folder:

```markup
$ ls -l /tmp/f?.txt
-rw-r--r--  1 mtsouk  wheel   14 Feb 27 19:44 /tmp/f1.txt
-rw-r--r--  1 mtsouk  wheel   14 Feb 27 19:44 /tmp/f2.txt
-rw-r--r--  1 mtsouk  wheel   14 Feb 27 19:44 /tmp/f3.txt
-rw-r--r--  1 mtsouk  wheel  101 Feb 27 19:44 /tmp/f4.txt
```

The previous output shows that the same amount of information has been written in `f1.txt`, `f2.txt`, and `f3.txt`, which means that the presented writing techniques are equivalent.

The next section shows how to work with JSON data in Go.

Bookmark

# Working with JSON

The Go standard library includes `encoding/json`, which is for working with JSON data. Additionally, Go allows you to add support for **JSON fields** in Go structures using **tags**, which is the subject of the _Structures and JSON_ subsection. Tags control the encoding and decoding of JSON records to and from Go structures. But first we should talk about marshaling and unmarshaling JSON records.

## Using Marshal() and Unmarshal()

Both the marshaling and unmarshaling of JSON data are important procedures for working with JSON data using Go structures. **Marshaling** is the process of converting a Go structure into a JSON record. You usually want that for transferring JSON data via computer networks or for saving it on disk. **Unmarshaling** is the process of converting a JSON record given as a byte slice into a Go structure. You usually want that when receiving JSON data via computer networks or when loading JSON data from disk files.

The number one bug when converting JSON records into Go structures and vice versa is not making the required fields of your Go structures **exported**. When you have issues with marshaling and unmarshaling, begin your debugging process from there.

The code in `encodeDecode.go` illustrates both the marshaling and unmarshaling of JSON records using hardcoded data for simplicity:

```markup
package main
import (
    "encoding/json"
    "fmt"
)
type UseAll struct {
    Name    string `json:"username"`
    Surname string `json:"surname"`
    Year    int    `json:"created"`
}
```

What the previous metadata tells us is that the `Name` field of the `UseAll` structure is translated to `username` in the JSON record, and **vice versa**, the `Surname` field is translated to `surname`, and **vice versa**, and the `Year` structure field is translated to `created` in the JSON record, and **vice versa**. This information has to do with the marshaling and unmarshaling of JSON data. Other than this, you treat and use `UseAll` as a regular Go structure.

```markup
func main() {
    useall := UseAll{Name: "Mike", Surname: "Tsoukalos", Year: 2021}
    // Regular Structure
    // Encoding JSON data -> Convert Go Structure to JSON record with fields
    t, err := json.Marshal(&useall)
```

The `json.Marshal()` function requires a pointer to a structure variable—its real data type is an empty interface variable—and returns a byte slice with the encoded information and an `error` variable.

```markup
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Printf("Value %s\n", t)
    }
    // Decoding JSON data given as a string
    str := `{"username": "M.", "surname": "Ts", "created":2020}`
```

JSON data usually comes as a string.

```markup
    // Convert string into a byte slice
    jsonRecord := []byte(str)
```

However, as `json.Unmarshal()` requires a byte slice, you need to convert that string into a byte slice before passing it to `json.Unmarshal()`.

```markup
    // Create a structure variable to store the result
    temp := UseAll{}
    err = json.Unmarshal(jsonRecord, &temp)
```

The `json.Unmarshal()` function requires the byte slice with the JSON record and a pointer to the Go structure variable that is going to store the JSON record and returns an `error` variable.

```markup
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Printf("Data type: %T with value %v\n", temp, temp)
    }
}
```

Running `encodeDecode.go` produces the next output:

```markup
$ go run encodeDecode.go
Value {"username":"Mike","surname":"Tsoukalos","created":2021}
Data type: main.UseAll with value {M. Ts 2020}
```

The next subsection illustrates how to define the JSON tags in a Go structure in more detail.

## Structures and JSON

Imagine that you have a Go structure that you want to convert into a JSON record without including any empty fields—the next code illustrates how to perform that task with the use of `omitempty`:

```markup
// Ignoring empty fields in JSON
type NoEmpty struct {
    Name    string `json:"username"`
    Surname string `json:"surname"`
    Year    int    `json:"creationyear,omitempty"`
}
```

Last, imagine that you have some sensitive data on some of the fields of a Go structure that you do not want to include in the JSON records. You can do that by including the `"-"` special value in the desired `json:` structure tags. This is shown in the next code excerpt:

```markup
// Removing private fields and ignoring empty fields
type Password struct {
    Name     string `json:"username"`
    Surname  string `json:"surname,omitempty"`
    Year     int    `json:"creationyear,omitempty"`
    Pass     string `json:"-"`
}
```

So, the `Pass` field is going to be ignored when converting a `Password` structure into a JSON record using `json.Marshal()`.

These two techniques are illustrated in `tagsJSON.go`, which can be found in the `ch06` directory of the GitHub repository of this book. Running `tagsJSON.go` produces the next output:

```markup
$ go run tagsJSON.go
noEmptyVar decoded with value {"username":"Mihalis","surname":""}
password decoded with value {"username":"Mihalis"}
```

For the first line of output, we have the following: the value of `noEmpty`, which is converted into a `NoEmpty` structure variable named `noEmptyVar`, is `NoEmpty{Name: "Mihalis"}`. The `noEmpty` structure has the default values for the `Surname` and `Year` fields. However, as they are not specifically defined, `json.Marshal()` ignores the `Year` field because it has the `omitempty` tag but does not ignore the `Surname` field, which has the empty string value.

For the second line of output: the value of the `password` variable is `Password{Name: "Mihalis", Pass: "myPassword"}`. When the `password` variable is converted into a JSON record, the `Pass` field is not included in the output. The remaining two fields of the `Password` structure, `Surname` and `Year`, are omitted because of the `omitempty` tag. So, what is left is the `username` field.

So far, we have seen working with single JSON records. But what happens when you have multiple records to process? The next subsection answers this question and many more!

## Reading and writing JSON data as streams

Imagine that you have a slice of Go structures that represent JSON records that you want to process. Should you process the records one by one? It can be done but does it look efficient? It does not! The good thing is that Go supports the processing of multiple JSON records as streams instead of individual records, which is faster and more efficient. This subsection teaches how to perform that using the `JSONstreams.go` utility, which contains the following two functions:

```markup
// DeSerialize decodes a serialized slice with JSON records
func DeSerialize(e *json.Decoder, slice interface{}) error {
    return e.Decode(slice)
}
```

The `DeSerialize()` function is used for reading input in the form of JSON records, decoding it, and putting it into a slice. The function writes the slice, which is of the `interface{}` data type and is given as a parameter, and gets its input from the buffer of the `*json.Decoder` parameter. The `*json.Decoder` parameter along with its buffer is defined in the `main()` function in order to avoid allocating it all the time and therefore losing the performance gains and efficiency of using this type—the same applies to the use of `*json.Encoder` that follows:

```markup
// Serialize serializes a slice with JSON records
func Serialize(e *json.Encoder, slice interface{}) error {
    return e.Encode(slice)
}
```

The `Serialize()` function accepts two parameters, a `*json.Encoder` and a slice of any data type, hence the use of `interface{}`. The function processes the slice and writes the output to the buffer of the `json.Encoder`—this buffer is passed as a parameter to the encoder at the time of its creation.

Both the `Serialize()` and `DeSerialize()` functions can work with any type of JSON record due to the use of `interface{}`.

The `JSONstreams.go` utility generates random data. Running `JSONstreams.go` creates the next output:

```markup
$ go run JSONstreams.go
After Serialize:[{"key":"XVLBZ","value":16},{"key":"BAICM","value":89}]
After DeSerialize:
0 {XVLBZ 16}
1 {BAICM 89}
```

The input slice of structures, which is generated in `main()`, is serialized as seen in the first line of the output. After that it is deserialized into the original slice of structures.

## Pretty printing JSON records

This subsection illustrates how to **pretty print** JSON records, which means printing JSON records in a pleasant and readable format without knowing the format of the Go structure that holds the JSON records. As there exist two ways to read JSON records, individually and as a stream, there exist two ways to pretty print JSON data: as single JSON records and as a stream. Therefore, we are going to implement two separate functions named `prettyPrint()` and `JSONstream()`, respectively.

The implementation of the `prettyPrint()` function is the following:

```markup
func PrettyPrint(v interface{}) (err error) {
    b, err := json.MarshalIndent(v, "", "\t")
    if err == nil {
        fmt.Println(string(b))
    }
    return err
}
```

All the work is done by `json.MarshalIndent()`, which **applies indent** to format the output.

Although both `json.MarshalIndent()` and `json.Marshal()` produce a JSON text result (_byte slice_), only `json.MarshalIndent()` allows applying customizable indent, whereas `json.Marshal()` generates a more compact output.

For pretty printing streams of JSON data, you should use the `JSONstream()` function:

```markup
func JSONstream(data interface{}) (string, error) {
  buffer := new(bytes.Buffer)
  encoder := json.NewEncoder(buffer)
  encoder.SetIndent("", "\t")
```

The `json.NewEncoder()` function returns a new `Encoder` that writes to a writer that is passed as a parameter to `json.NewEncoder()`. An `Encoder` writes JSON values to an output stream. Similarly to `json.MarshalIndent()`, the `SetIndent()` method allows you to apply a customizable indent to a stream.

```markup
  err := encoder.Encode(data)
  if err != nil {
    return "", err
  }
  return buffer.String(), nil
}
```

After we are done configuring the encoder, we are free to process the JSON stream using `Encode()`.

These two functions are illustrated in `prettyPrint.go`, which generates JSON records using random data. Running `prettyPrint.go` produces the following kind of output:

```markup
Last record: {BAICM 89}
{
        "key": "BAICM",
        "value": 89
}
[
        {
                "key": "XVLBZ",
                "value": 16
        },
        {
                "key": "BAICM",
                "value": 89
        }
]
```

The previous output shows the beautified output of a single JSON record followed by the beautified output of a slice with two JSON records—all JSON records are represented as Go structures.

The next section is about working with XML in Go.

Bookmark

# Working with XML

This section briefly describes how to work with XML data in Go using records. The idea behind XML and Go is the same as with JSON and Go. You put tags in Go structures in order to specify the XML tags and you can still serialize and deserialize XML records using `xml.Unmarshal()` and `xml.Marshal()`, which are found in the `encoding/xml` package. However, there exist some differences that are illustrated in `xml.go`:

```markup
package main
import (
    "encoding/xml"
    "fmt"
)
type Employee struct {
    XMLName   xml.Name `xml:"employee"`
    ID        int      `xml:"id,attr"`
    FirstName string   `xml:"name>first"`
    LastName  string   `xml:"name>last"`
    Height    float32  `xml:"height,omitempty"`
    Address
    Comment string `xml:",comment"`
}
```

This is where the structure for the XML data is defined. However, there is additional information regarding the name and the type of each XML element. The `XMLName` field provides the name of the XML record, which in this case will be `employee`.

A field with the tag `",comment"` is a comment and it is formatted as such in the output. A field with the tag `attr` appears as an attribute to the provided field name (which is `id` in this case) in the output. The `name>first` notation tells Go to embed the `first` tag inside a tag called `name`.

Lastly, a field with the `omitempty` option is omitted from the output if it is empty. An **empty value** is any of `0`, `false`, a `nil` pointer or interface, and any array, slice, map, or string with a length of zero.

```markup
type Address struct {
    City, Country string
}
func main() {
    r := Employee{ID: 7, FirstName: "Mihalis", LastName: "Tsoukalos"}
    r.Comment = "Technical Writer + DevOps"
    r.Address = Address{"SomeWhere 12", "12312, Greece"}
    output, err := xml.MarshalIndent(&r, "  ", "    ")
```

As is the case with JSON, `xml.MarshalIndent()` is for beautifying the output.

```markup
    if err != nil {
        fmt.Println("Error:", err)
    }
    output = []byte(xml.Header + string(output))
    fmt.Printf("%s\n", output)
}
```

The output of `xml.go` is the following:

```markup
<?xml version="1.0" encoding="UTF-8"?>
  <employee id="7">
      <name>
          <first>Mihalis</first>
          <last>Tsoukalos</last>
      </name>
      <City>SomeWhere 12</City>
      <Country>12312, Greece</Country>
      <!--Technical Writer + DevOps-->
  </employee>
```

The previous output shows the XML version of the Go structure given as input to the program.

In the next section we develop a utility that converts JSON records to XML records and vice versa.

## Converting JSON to XML and vice versa

As promised, we are going to produce a utility that converts records between the JSON and XML formats. The input is given as a command-line argument. The utility tries to guess the format of the input starting from XML. If `xml.Unmarshal()` fails, then the utility tries using `json.Unmarshal()`. If there is not a match, then the user is informed about the error condition. On the other hand, if `xml.Unmarshal()` is successful, the data is stored into an `XMLrec` variable and is then converted into a `JSONrec` variable. The same happens with `json.Unmarshal()` in the case where the `xml.Unmarshal()` call is unsuccessful.

The logic of the utility can be found in the Go structures:

```markup
type XMLrec struct {
    Name    string `xml:"username"`
    Surname string `xml:"surname,omitempty"`
    Year    int    `xml:"creationyear,omitempty"`
}
type JSONrec struct {
    Name    string `json:"username"`
    Surname string `json:"surname,omitempty"`
    Year    int    `json:"creationyear,omitempty"`
}
```

Both structures store the same data. However, the former (`XMLrec`) is for storing XML data whereas the latter (`JSONrec`) is for storing JSON data.

Running `JSON2XML.go` produces the next kind of output:

```markup
$ go run JSON2XML.go '<XMLrec><username>Mihalis</username></XMLrec>'
<XMLrec><username>Mihalis</username></XMLrec>
{"username":"Mihalis"}
```

So, we give an XML record as input, which is converted into a JSON record.

The next output illustrates the reverse process:

```markup
$ go run JSON2XML.go '{"username": "Mihalis"}'
{"username": "Mihalis"}
<XMLrec><username>Mihalis</username></XMLrec>
```

In the previous output the input is a JSON record, and the output is an XML record.

The next section discusses working with YAML files in Go.

Bookmark

# Working with YAML

In this section, we briefly discuss how to work with YAML files in Go. The Go standard library does not include support for YAML files, which means that you should look at external packages for YAML support. There exist three main packages that allow you to work with YAML from Go:

-   [https://github.com/kylelemons/go-gypsy](https://github.com/kylelemons/go-gypsy)
-   [https://github.com/go-yaml/yaml](https://github.com/go-yaml/yaml)
-   [https://github.com/goccy/go-yaml](https://github.com/goccy/go-yaml)

Choosing one is a matter of personal preference. We are going to work with `go-yaml` in this section using the code found in `yaml.go`. Due to the use of Go modules, `yaml.go` is developed in `~/go/src/github.com/mactsouk/yaml`—you can also find it in the GitHub repository of this book. The most important part of it is the next:

```markup
var yamlfile = `
image: Golang
matrix:
  docker: python
  version: [2.7, 3.9]
`
```

The `yamlfile` variable contains the YAML data; you usually read the data from a file—we are just using that to save some space.

```markup
type Mat struct {
    DockerImage string    `yaml:"docker"`
    Version     []float32 `yaml:",flow"`
}
```

The `Mat` structure defines two fields and their associations with the YAML file. The `Version` field is a slice of `float32` values. As there is no name for the `Version` field, the name is going to be `version`. The `flow` keyword says that the marshal is using a flow style, which is useful for structs, sequences, and maps.

```markup
type YAML struct {
    Image  string
    Matrix Mat
}
```

The YAML structure embeds a `Mat` structure and contains a field named `Image`, which is associated with `image` in the YAML file. The `main()` function contains the expected `yaml.Unmarshal()` and `yaml.Marshal()` calls.

Once you have the source file at the desired place, run the next commands—if you need to run any extra commands, the `go` binary is nice enough to help you:

```markup
$ go mod init
$ go mod tidy
```

The `go mod init` command initializes and writes a new `go.mod` file in the current directory whereas the `go mod tidy` command synchronizes `go.mod` with the source code.

If you want to play it safe and you are using packages that do not belong to the standard library, then developing inside `~/go/src`, committing to a GitHub repository, and using Go modules for all dependencies might be the best option. However, this does not mean that you must develop your own packages in the form of Go modules.

Running `yaml.go` produces the next output:

```markup
$ go run yaml.go
After Unmarshal (Structure):
{Golang {python [2.7 3.9]}}
After Marshal (YAML code):
image: Golang
matrix:
  docker: python
  version: [2.7, 3.9]
```

The previous output shows how the `{Golang {python [2.7 3.9]}}` text is converted into a YAML file and vice versa. Now that we know about working with JSON, XML, and YAML data in Go, we are ready to learn about the `viper` package.

Bookmark

# The viper package

**Flags** are specially formatted strings that are passed into a program to control its behavior. Dealing with flags on your own might become very frustrating if you want to support multiple flags and options. Go offers the `flag` package for working with command-line options, parameters, and flags. Although `flag` can do many things, it is not as capable as other external Go packages. Thus, if you are developing simple UNIX system command-line utilities, you might find the `flag` package very interesting and useful. But you are not reading this book to create simple command-line utilities! Therefore, I'll skip the `flag` package and introduce you to an external package named `viper`, which is a powerful Go package that supports a plethora of options. `viper` uses the `pflag` package instead of `flag`, which is also illustrated in the code we will look at in the following sections.

All `viper` projects follow a pattern. First, you initialize `viper` and then you define the elements that interest you. After that, you get these elements and read their values in order to use them. The desired values can be taken either directly, as happens when you are using the `flag` package from the standard Go library, or indirectly using configuration files. When using formatted configuration files in the JSON, YAML, TOML, HCL, or Java properties format, `viper` does all the parsing for you, which saves you from having to write and debug lots of Go code. `viper` also allows you to extract and save values in Go structures. However, this requires that the fields of the Go structure match the keys of the configuration file.

The home page of `viper` is on GitHub ([https://github.com/spf13/viper](https://github.com/spf13/viper)). Please note that you are not obliged to use every capability of `viper` in your tools—just the features that you want. The general rule is to use the features of Viper that simplify your code. Put simply, if your command-line utility requires too many command-line parameters and flags, then it would be better to use a configuration file instead.

## Using command-line flags

The first example shows how to write a simple utility that accepts two values as command-line parameters and prints them on screen for verification. This means that we are going to need two command-line flags for these parameters.

Starting from **Go version 1.16**, using modules is the default behavior, which the `viper` package needs to use. So, you need to put `useViper.go`, which is the name of the source file, inside `~/go` for things to work. As my GitHub username is `mactsouk`, I had to run the following commands:

```markup
$ mkdir ~/go/src/github.com/mactsouk/useViper
$ cd ~/go/src/github.com/mactsouk/useViper
$ vi useViper.go
$ go mod init
$ go mod tidy
```

You can either edit `useViper.go` on your own or copy it from the GitHub repository of this book. Keep in mind that the last two commands should be executed when `useViper.go` is ready and includes all required external packages.

The implementation of `useViper.go` is as follows:

```markup
package main
import (
    "fmt"
    "github.com/spf13/pflag"
    "github.com/spf13/viper"
)
```

We need to import both the `pflag` and `viper` packages as we are going to use the functionality from both of them.

```markup
func aliasNormalizeFunc(f *pflag.FlagSet, n string) pflag.NormalizedName {
    switch n {
    case "pass":
        n = "password"
        break
    case "ps":
        n = "password"
        break
    }
    return pflag.NormalizedName(n)
}
```

The `aliasNormalizeFunc()` function is used for creating additional aliases for a flag—in this case an alias for the `--password` flag. According to the existing code, the `--password` flag can be accessed as either `--pass` or `--ps`.

```markup
func main() {
    pflag.StringP("name", "n", "Mike", "Name parameter")
```

In the preceding code, we create a new flag called `name` that can also be accessed as `-n`. Its default value is `Mike` and its description that appears in the usage of the utility is `Name parameter`.

```markup
    pflag.StringP("password", "p", "hardToGuess", "Password")
    pflag.CommandLine.SetNormalizeFunc(aliasNormalizeFunc)
```

We create another flag named `password` that can also be accessed as `-p` and has a default value of `hardToGuess` and a description. Additionally, we register a **normalization function** for generating aliases for the `password` flag.

```markup
    pflag.Parse()
    viper.BindPFlags(pflag.CommandLine)
```

The `pflag.Parse()` call should be used after all command-line flags are defined—its purpose is to parse the command-line flags into the defined flags.

Additionally, the `viper.BindPFlags()` call makes all flags available to the `viper` package—strictly speaking, we say that the `viper.BindPFlags()` call binds an existing set of `pflag` flags (`pflag.FlagSet`) to `viper`.

```markup
    name := viper.GetString("name")
    password := viper.GetString("password")
```

The previous commands show how to read the values of two `string` command-line flags.

```markup
    fmt.Println(name, password)
    // Reading an Environment variable
    viper.BindEnv("GOMAXPROCS")
    val := viper.Get("GOMAXPROCS")
    if val != nil {
        fmt.Println("GOMAXPROCS:", val)
    }
```

The `viper` package can work with environment variables. We first need to call `viper.BindEnv()` to tell `viper` the environment variable that interests us and then we can read its value by calling `viper.Get()`. If `GOMAXPROCS` is not already set, the `fmt.Println()` call will not get executed.

```markup
    // Setting an Environment variable
    viper.Set("GOMAXPROCS", 16)
    val = viper.Get("GOMAXPROCS")
    fmt.Println("GOMAXPROCS:", val)
}
```

Similarly, we can change the current value of an environment variable using `viper.Set()`.

The good thing is that `viper` automatically provides usage information:

```markup
$ go run useViper.go --help
Usage of useViper:
  -n, --name string       Name parameter (default "Mike")
  -p, --password string   Password (default "hardToGuess")
pflag: help requested
exit status 2
```

Using `useViper.go` without any command-line arguments produces the next kind of output. Remember that we are inside `~/go/src/github.com/mactsouk/useViper`:

```markup
$ go run useViper.go
Mike hardToGuess
GOMAXPROCS: 16
```

However, if we provide values for the command-line flags, the output is going to be slightly different:

```markup
$ go run useViper.go -n mtsouk -p hardToGuess
mtsouk hardToGuess
GOMAXPROCS: 16
```

In this second case, we used the shortcuts for the command-line flags because it is faster.

The next subsection discusses the use of JSON files for storing configuration information.

## Reading JSON configuration files

The `viper` package can read JSON files to get its configuration, and this subsection illustrates how. Using text files for storing configuration details can be very helpful when writing complex applications that require lots of data and setup. This is illustrated in `jsonViper.go`.

Once again, we need to put `jsonViper.go` inside `~/go/src/github.com/mactsouk/jsonViper`—please adjust that command to fit your own GitHub username, although if you do not create a GitHub repository, you can use `mactsouk`. The code of `jsonViper.go` is as follows:

```markup
package main
import (
    "encoding/json"
    "fmt"
    "os"
    "github.com/spf13/viper"
)
type ConfigStructure struct {
    MacPass     string `mapstructure:"macos"`
    LinuxPass   string `mapstructure:"linux"`
    WindowsPass string `mapstructure:"windows"`
    PostHost    string `mapstructure:"postgres"`
    MySQLHost   string `mapstructure:"mysql"`
    MongoHost   string `mapstructure:"mongodb"`
}
```

There is an _important point here_: although we are using a JSON file to store the configuration, the Go structure uses `mapstructure` instead of `json` for the fields of the JSON configuration file.

```markup
var CONFIG = ".config.json"
func main() {
    if len(os.Args) == 1 {
        fmt.Println("Using default file", CONFIG)
    } else {
        CONFIG = os.Args[1]
    }
    viper.SetConfigType("json")
    viper.SetConfigFile(CONFIG)
    fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
    viper.ReadInConfig()
```

The previous four statements declare that we are using a JSON file, let `viper` know the path to the configuration file, print the configuration file used, and read and parse that configuration file.

Keep in mind that `viper` does not check whether the configuration file actually exists and is readable. If the file cannot be found or read, `viper.ReadInConfig()` acts like processing an empty configuration file.

```markup
    if viper.IsSet("macos") {
        fmt.Println("macos:", viper.Get("macos"))
    } else {
        fmt.Println("macos not set!")
    }
```

The `viper.IsSet()` call checks whether a key named `macos` can be found in the configuration. If it is set, it reads its value using `viper.Get("macos")` and prints it on screen.

```markup
    if viper.IsSet("active") {
        value := viper.GetBool("active")
        if value {
            postgres := viper.Get("postgres")
            mysql := viper.Get("mysql")
            mongo := viper.Get("mongodb")
            fmt.Println("P:", postgres, "My:", mysql, "Mo:", mongo)
        }
    } else {
        fmt.Println("active is not set!")
    }
```

In the aforementioned code, we check whether the `active` key can be found before reading its value. If its value is equal to `true` then we read the values from three more keys named `postgres`, `mysql`, and `mongodb`.

As the `active` key should hold a Boolean value, we use `viper.GetBool()` for reading it.

```markup
    if !viper.IsSet("DoesNotExist") {
        fmt.Println("DoesNotExist is not set!")
    }
```

As expected, trying to read a key that does not exist fails.

```markup
    var t ConfigStructure
    err := viper.Unmarshal(&t)
    if err != nil {
        fmt.Println(err)
        return
    }
```

The call to `viper.Unmarshal()` allows you to put the information from the JSON configuration file into a properly defined Go structure—this is optional but handy.

```markup
    PrettyPrint(t)
}
```

The implementation of the `PrettyPrint()` function was presented in `prettyPrint.go` earlier on in this chapter.

Now you need to download the dependencies of `jsonViper.go`:

```markup
$ go mod init
$ go mod tidy # This command is not always required
```

The contents of the current directory are as follows:

```markup
$ ls -l
total 44
-rw-r--r-- 1 mtsouk users    85 Feb 22 18:46 go.mod
-rw-r--r-- 1 mtsouk users 29678 Feb 22 18:46 go.sum
-rw-r--r-- 1 mtsouk users  1418 Feb 22 18:45 jsonViper.go
-rw-r--r-- 1 mtsouk users   189 Feb 22 18:46 myConfig.json
```

The contents of the `myConfig.json` file used for testing are as follows:

```markup
{
    "macos": "pass_macos",
    "linux": "pass_linux",
    "windows": "pass_windows",
    "active": true,
    "postgres": "machine1",
    "mysql": "machine2",
    "mongodb": "machine3"
}
```

Running `jsonViper.go` on the preceding JSON file produces the next output:

```markup
$ go run jsonViper.go myConfig.json
Using config: myConfig.json
macos: pass_macos
P: machine1 My: machine2 Mo: machine3
DoesNotExist is not set!
{
  "MacPass": "pass_macos",
  "LinuxPass": "pass_linux",
  "WindowsPass": "pass_windows",
  "PostHost": "machine1",
  "MySQLHost": "machine2",
  "MongoHost": "machine3"
}
```

The previous output is generated by `jsonViper.go` when parsing `myConfig.json` and trying to find the desired information.

The next section discusses a Go package for creating powerful and professional command-line utilities such as `docker` and `kubectl`.

Bookmark

# The cobra package

`cobra` is a very handy and popular Go package that allows you to develop command-line utilities with commands, subcommands, and aliases. If you have ever used `hugo`, `docker`, or `kubectl` you are going to realize immediately what the `cobra` package does, as all these tools are developed using `cobra`. Commands can have one or more aliases, which is very handy when you want to please both amateur and experienced users. `cobra` also supports **persistent flags** and **local flags**, which are flags that are available to all commands and flags that are available to given commands only, respectively. Also, by default, `cobra` uses `viper` for parsing its command-line arguments.

All `cobra` projects follow the same development pattern. You use the `cobra` utility, then you create commands, and then you make the desired changes to the generated Go source code files in order to implement the desired functionality. Depending on the complexity of your utility, you might need to make lots of changes to the created files. Although `cobra` saves you lots of time, you still have to write the code that implements the desired functionality for each command.

You need to take some extra steps in order to download the `cobra` binary the right way:

```markup
$ GO111MODULE=on go get -u -v github.com/spf13/cobra/cobra
```

The previous command downloads the `cobra` binary and the required dependencies using Go modules even if you are using a Go version older than 1.16.

It is not necessary to know about all of the supported environment variables such as `GO111MODULE`, but sometimes they can help you resolve tricky problems with your Go installation. So, if you want to learn about your current Go environment, you can use the `go env` command.

For the purposes of this section, we are going to need a GitHub repository—this is optional, but it is the only way for the readers of this book to have access to the presented code.

The path of the GitHub repository is [https://github.com/mactsouk/go-cobra](https://github.com/mactsouk/go-cobra). The first thing to do is place the files of the GitHub repository in the right place. Everything is going to be much easier if you put it inside `~/go`; the exact place depends on the GitHub repository, because the Go compiler will not have to search for the Go files.

In our case, we are going to put it under `~/go/src/github.com/mactsouk` because `mactsouk` is my GitHub username. This requires running the next commands:

```markup
$ cd ~/go/src/github.com
$ mkdir mactsouk # only required if the directory is not there
$ cd mactsouk
$ git clone git@github.com:mactsouk/go-cobra.git
$ cd go-cobra
$ ~/go/bin/cobra init --pkg-name github.com/mactsouk/go-cobra
Using config file: /Users/mtsouk/.cobra.yaml
Your Cobra application is ready at
/Users/mtsouk/go/src/github.com/mactsouk/go-cobra
$ go mod init 
go: creating new go.mod: module github.com/mactsouk/go-cobra
```

As the `cobra` package works better with modules, we define the project dependencies using Go modules. In order to specify that a Go project uses Go modules, you should execute `go mod init`. This command creates two files named `go.sum` and `go.mod`.

```markup
$ go run main.go 
go: finding module for package github.com/spf13/cobra
go: finding module for package github.com/mitchellh/go-homedir
go: finding module for package github.com/spf13/viper
go: found github.com/mitchellh/go-homedir in github.com/mitchellh/go-homedir v1.1.0
go: found github.com/spf13/cobra in github.com/spf13/cobra v1.1.3
go: found github.com/spf13/viper in github.com/spf13/viper v1.7.1
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.
```

All lines beginning with `go:` have to do with Go modules and will appear only once. The last lines are the default message of a `cobra` project—we are going to modify that message later on. You are now ready to begin working with the `cobra` tool.

## A utility with three commands

This subsection illustrates the use of the `cobra add` command, which is used for adding new commands to a `cobra` project. The names of the commands are `one`, `two`, and `three`:

```markup
$ ~/go/bin/cobra add one
Using config file: /Users/mtsouk/.cobra.yaml
one created at /Users/mtsouk/go/src/github.com/mactsouk/go-cobra
$ ~/go/bin/cobra add two 
$ ~/go/bin/cobra add three
```

The previous commands create three new files in the `cmd` folder named `one.go`, `two.go`, and `three.go`, which are the initial implementations of the three commands.

The first thing you should usually do is delete unwanted code from `root.go` and change the messages of the utility and each command as described in the `Short` and `Long` fields. However, if you want you can leave the source files unchanged.

The next subsection enriches the utility by adding command-line flags to the commands.

## Adding command-line flags

We are going to create two global command-line flags and one command-line flag that is attached to a given command (`two`) and is not supported by the other two commands. Global command-line flags are defined in the `./cmd/root.go` file. We are going to define two global flags named `directory`, which is a string, and `depth`, which is an unsigned integer.

Both global flags are defined in the `init()` function of `./cmd/root.go`.

```markup
rootCmd.PersistentFlags().StringP("directory", "d", "/tmp", "Path to use.")
rootCmd.PersistentFlags().Uint("depth", 2, "Depth of search.")
viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))
viper.BindPFlag("depth", rootCmd.PersistentFlags().Lookup("depth"))
```

We use `rootCmd.PersistentFlags()` to define global flags followed by the data type of the flag. The name of the first flag is `directory` and its shortcut is `d` whereas the name of the second flag is `depth` and has no shortcut—if you want to add a shortcut to it, you should use the `UintP()` method instead. After defining the two flags, we pass their control to `viper` by calling `viper.BindPFlag()`. The first flag is a `string` whereas the second one is a `uint` value. As both of them are available in the `cobra` project, we call `viper.GetString("directory")` to get the value of the `directory` flag and `viper.GetUint("depth")` to get the value of the `depth` flag.

Last, we add a command-line flag that is only available to the `two` command using the next line in the `./cmd/two.go` file:

```markup
twoCmd.Flags().StringP("username", "u", "Mike", "Username value")
```

The name of the flag is `username` and its shortcut is `u`. As this is a local flag available to the `two` command only, we can get its value by calling `cmd.Flags().GetString("username")` inside the `./cmd/two.go` file only.

The next subsection creates command aliases for the existing commands.

## Creating command aliases

In this subsection we continue building on the code from the previous subsection by creating aliases for existing commands. This means that commands `one`, `two`, and `three` will also be accessible as `cmd1`, `cmd2`, and `cmd3` respectively.

In order to do that, you need to add an extra field named `Aliases` in the `cobra.Command` structure of each command—the data type of the `Aliases` field is `string` slice. So, for the `one` command, the beginning of the `cobra.Command` structure in `./cmd/one.go` will look as follows:

```markup
var oneCmd = &cobra.Command{
    Use:     "one",
    Aliases: []string{"cmd1"},
    Short:   "Command one",
```

You should make similar changes to `./cmd/two.go` and `./cmd/three.go`. Please keep in mind that the **internal name** of the `one` command is `oneCmd`—the other commands have analogous internal names.

If you accidentally put the `cmd1` alias, or any other alias, in multiple commands, the Go compiler will not complain. However, only its first occurrence gets executed.

The next subsection enriches the utility by adding subcommands for the `one` and `two` commands.

## Creating subcommands

This subsection illustrates how to create two subcommands for the command named `three`. The names of the two subcommands will be `list` and `delete`. The way to create them using the `cobra` utility is the following:

```markup
$ ~/go/bin/cobra add list -p 'threeCmd'
Using config file: /Users/mtsouk/.cobra.yaml
list created at /Users/mtsouk/go/src/github.com/mactsouk/go-cobra
$ ~/go/bin/cobra add delete -p 'threeCmd'
Using config file: /Users/mtsouk/.cobra.yaml
delete created at /Users/mtsouk/go/src/github.com/mactsouk/go-cobra
```

The previous commands create two new files inside `./cmd` named `delete.go` and `list.go`. The `-p` flag is followed by the **internal name** of the command you want to associate the subcommands with. The internal name of the `three` command is `threeCmd`. You can verify that these two commands are associated with the `three` command as follows:

```markup
$ go run main.go three delete
delete called
$ go run main.go three list
list called
```

If you run `go run main.go two list`, Go considers `list` as a command-line argument of `two` and it will not execute the code in `./cmd/list.go`. The final version of the `go-cobra` project has the following structure and contains the following files, as generated by the `tree(1)` utility:

```markup
$ tree
.
├── LICENSE
├── README.md
├── cmd
│   ├── delete.go
│   ├── list.go
│   ├── one.go
│   ├── root.go
│   ├── three.go
│   └── two.go
├── go.mod
├── go.sum
└── main.go
1 directory, 11 files
```

At this point you might wonder what happens when you want to create two subcommands with the same name for two different commands. In that case, you create the first subcommand and rename its file before creating the second one.

As there is no point in presenting long listings of code, you can find the code of the `go-cobra` project at [https://github.com/mactsouk/go-cobra](https://github.com/mactsouk/go-cobra). The `cobra` package is also illustrated in the final section where we radically update the phone book application.

Bookmark

# Finding cycles in a UNIX file system

This section implements a practical UNIX command-line utility that can find cycles (loops) in UNIX file systems. The idea behind the utility is that with UNIX **symbolic links**, there is a possibility to create cycles in our file system. This can perplex backup software such as `tar(1)` or utilities such as `find(1)` and can create security-related issues. The presented utility, which is called `FScycles.go`, tries to inform us about such situations.

The idea behind the solution is that we keep every visited directory path in a map and if a path appears for the second time, then we have a cycle. The map is called `visited` and is defined as `map[string]int`.

If you are wondering why we are using a string and not a byte slice or some other kind of slice as the key for the `visited` map, it is because maps cannot have slices as keys because **slices are not comparable**.

The output of the utility depends on the root path used for initialing the search process—that path is given as a command-line argument to the utility.

The `filepath.Walk()` function does not traverse symbolic links by design in order to avoid cycles. However, in our case, we want to traverse symbolic links to directories in order to discover loops. We solve that issue in a while.

The utility uses `IsDir()`, which is a function that helps you to identify directories—we are only interested in directories because only directories and symbolic links to directories can create cycles in file systems. Last, the utility uses `os.Lstat()` because it can handle symbolic links. Additionally, `os.Lstat()` returns information about the symbolic link without following it, which is not the case with `os.Stat()`—in this case we do not want to automatically follow symbolic links.

The important code of `FScycles.go` can be found in the implementation of `walkFunction()`:

```markup
func walkFunction(path string, info os.FileInfo, err error) error {
    fileInfo, err := os.Stat(path)
    if err != nil {
        return nil
    }
    fileInfo, _ = os.Lstat(path)
    mode := fileInfo.Mode()
```

First, we make sure that the path actually exists, and then we call `os.Lstat()`.

```markup
    // Find regular directories first
    if mode.IsDir() {
        abs, _ := filepath.Abs(path)
        _, ok := visited[abs]
        if ok {
            fmt.Println("Found cycle:", abs)
            return nil
        }
        visited[abs]++
        return nil
    }
```

If a regular directory is already visited, then we have a cycle. The `visited` map keeps track of all visited directories.

```markup
    // Find symbolic links to directories
    if fileInfo.Mode()&os.ModeSymlink != 0 {
        temp, err := os.Readlink(path)
        if err != nil {
            fmt.Println("os.Readlink():", err)
            return err
        }
        newPath, err := filepath.EvalSymlinks(temp)
        if err != nil {
            return nil
        }
```

The `filepath.EvalSymlinks()` function is used for finding out where symbolic links point to. If that destination is another directory, then the code that follows makes sure that it is going to be visited as well using an additional call to `filepath.Walk()`.

```markup
        linkFileInfo, err := os.Stat(newPath)
        if err != nil {
            return err
        }
        linkMode := linkFileInfo.Mode()
        if linkMode.IsDir() {
            fmt.Println("Following...", path, "-->", newPath)
```

The `linkMode.IsDir()` statement makes sure that only directories are being followed.

```markup
            abs, _ := filepath.Abs(newPath)
```

The call to `filepath.Abs()` returns the absolute path of the path that is given as a parameter. The keys of the `visited` slice are values returned by `filepath.Abs()`.

```markup
            _, ok := visited[abs]
            if ok {
                fmt.Println("Found cycle!", abs)
                return nil
            }
            visited[abs]++
            err = filepath.Walk(newPath, walkFunction)
            if err != nil {
                return err
            }
            return nil
        }
    }
    return nil
}
```

Running `FScycles.go` produces the following kind of output:

```markup
$ go run FScycles.go ~
Following... /home/mtsouk/.local/share/epiphany/databases/indexeddb/v0 --> /home/mtsouk/.local/share/epiphany/databases/indexeddb
Found cycle! /home/mtsouk/.local/share/epiphany/databases/indexeddb
```

The previous output tells us that there is a cycle in the home directory of the current user—once we identify the loop, we should remove it on our own.

The remaining sections of this chapter discuss some new features that came with Go version 1.16.

Bookmark

# New to Go 1.16

Go 1.16 came with some new features including embedding files in Go binaries as well as the introduction of the `os.ReadDir()` function, the `os.DirEntry` type, and the `io/fs` package.

As these features are related to systems programming, they are included and explored in the current chapter. We begin by presenting the embedding of files into Go binary executables.

## Embedding files

This section presents a feature that first appeared in Go 1.16 that allows you to **embed static assets** into Go binaries. The allowed data types for keeping an embedded file are `string`, `[]byte`, and `embed.FS`. This means that a Go binary may contain a file that you do not have to manually download when you execute the Go binary! The presented utility embeds two different files that it can retrieve based on the given command-line argument.

The code that follows, which is saved as `embedFiles.go`, illustrates this new Go feature:

```markup
package main
import (
    _ "embed"
    "fmt"
    "os"
)
```

You need the `embed` package in order to embed any files in your Go binaries. As the `embed` package is not directly used, you need to put `_` in front of it so that the Go compiler won't complain.

```markup
//go:embed static/image.png
var f1 []byte
```

You need to begin a line with `//go:embed`, which denotes a Go comment but is treated in a special way, followed by the path to the file you want to embed. In this case we embed `static/image.png`, which is a binary file. The next line should define the variable that is going to hold the data of the embedded file, which in this case is a byte slice named `f1`—using a byte slice is recommended for binary files because we are going to directly use that byte slice to save that binary file.

```markup
//go:embed static/textfile
var f2 string
```

In this case we save the contents of a plain text file, which is `static/textfile`, in a `string` variable named `f2`.

```markup
func writeToFile(s []byte, path string) error {
    fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer fd.Close()
    n, err := fd.Write(s)
    if err != nil {
        return err
    }
    fmt.Printf("wrote %d bytes\n", n)
    return nil
}
```

The `writeToFile()` function is used for storing a byte slice into a file and is a helper function that can be used in other cases as well.

```markup
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Print select 1|2")
        return
    }
    fmt.Println("f1:", len(f1), "f2:", len(f2))
```

This statement prints the lengths of the `f1` and `f2` variables to make sure that they represent the size of the embedded files.

```markup
    switch arguments[1] {
    case "1":
        filename := "/tmp/temporary.png"
        err := writeToFile(f1, filename)
        if err != nil {
            fmt.Println(err)
            return
        }
    case "2":
        fmt.Print(f2)
    default:
        fmt.Println("Not a valid option!")
    }
}
```

The `switch` block is responsible for returning the desired file to the user—in the case of `static/textfile`, the file contents are printed on the screen. For the binary file, we decided to store it as `/tmp/temporary.png`.

This time we are going to compile `embedFiles.go` to make things more realistic, because it is the executable binary file that holds the embedded files. We build the binary file using `go build embedFiles.go`. Running `embedFiles` produces the following kind of output:

```markup
$ ./embedFiles 2
f1: 75072 f2: 14
Data to write
$ ./embedFiles 1
f1: 75072 f2: 14
wrote 75072 bytes
```

The next output verifies that `temporary.png` is located at the right path (`/tmp/temporary.png`):

```markup
$ ls -l /tmp/temporary.png 
-rw-r--r--  1 mtsouk  wheel  75072 Feb 25 15:20 /tmp/temporary.png
```

Using the embedding functionality, we can create a utility that embeds its own source code and prints it on screen when it gets executed! This is a fun way of using embedded files. The source code of `printSource.go` is the following:

```markup
package main
import (
    _ "embed"
    "fmt"
)
//go:embed printSource.go
var src string
func main() {
    fmt.Print(src)
}
```

As before, the file that is being embedded is defined in the `//go:embed` line. Running `printSource.go` prints the aforementioned code on screen.

## ReadDir and DirEntry

This section discusses `os.ReadDir()` and `os.DirEntry`. However, it begins by discussing the deprecation of the `io/ioutil` package—the functionality of the `io/ioutil` package has been transferred to other packages. So, we have the following:

-   `os.ReadDir()`, which is a new function, returns `[]DirEntry`. This means that it cannot directly replace `ioutil.ReadDir()`, which returns `[]FileInfo`. Although neither `os.ReadDir()` nor `os.DirEntry` offers any new functionality, they make things faster and simpler, which is important.
-   The `os.ReadFile()` function directly replaces `ioutil.ReadFile()`.
-   The `os.WriteFile()` function can directly replace `ioutil.WriteFile()`.
-   Similarly, `os.MkdirTemp()` can replace `ioutil.TempDir()` without any changes. However, as the `os.TempDir()` name was already taken, the new function name is different.
-   The `os.CreateTemp()` function is the same as `ioutil.TempFile()`. Although the name `os.TempFile()` was not taken, the Go people decided to name it `os.CreateTemp()` in order to be on par with `os.MkdirTemp()`.

Both `os.ReadDir()` and `os.DirEntry` can be found as `fs.ReadDir()` and `fs.DirEntry` in the `io/fs` package for working with the file system interface found in `io/fs`.

The `ReadDirEntry.go` utility illustrates the use of `os.ReadDir()`. Additionally, we are going to see `fs.DirEntry` in combination with `fs.WalkDir()` in action in the next section—`io/fs` only supports `WalkDir()`, which uses `DirEntry` by default. Both `fs.WalkDir()` and `filepath.WalkDir()` are using `DirEntry` instead of `FileInfo`. This means that in order to see any performance improvements when walking directory trees, you need to change `filepath.Walk()` calls to `filepath.WalkDir()` calls.

The presented utility calculates the size of a directory tree using `os.ReadDir()` with the help of the next function:

```markup
func GetSize(path string) (int64, error) {
    contents, err := os.ReadDir(path)
    if err != nil {
        return -1, err
    }
    var total int64
    for _, entry := range contents {
        // Visit directory entries
        if entry.IsDir() {
```

If we are processing a directory, we need to keep digging.

```markup
            temp, err := GetSize(filepath.Join(path, entry.Name()))
            if err != nil {
                return -1, err
            }
            total += temp
            // Get size of each non-directory entry
        } else {
```

If it is a file, then we just need to get its size. This involves calling `Info()` to get general information about the file and then `Size()` to get the size of the file:

```markup
            info, err := entry.Info()
            if err != nil {
                return -1, err
            }
            // Returns an int64 value
            total += info.Size()
        }
    }
    return total, nil
}
```

Running `ReadDirEntry.go` produces the next output, which indicates that the utility works as expected:

```markup
$ go run ReadDirEntry.go /usr/bin
Total Size: 1170983337
```

Last, keep in mind that both `ReadDir` and `DirEntry` are copied from the Python programming language.

The next section introduces us to the `io/fs` package.

## The io/fs package

This section illustrates the functionality of the `io/fs` package, which was first introduced in Go 1.16. As `io/fs` offers a unique kind of functionality, we begin this section by explaining what `io/fs` can do. Put simply, `io/fs` offers a **read-only** file system interface named `FS`. Note that `embed.FS` implements the `fs.FS` interface, which means that `embed.FS` can take advantage of some of the functionality offered by the `io/fs` package. This means that your applications can create their own internal file systems and work with their files.

The code example that follows, which is saved as `ioFS.go`, creates a file system using `embed` by putting all the files of the `./static` folder in there. `ioFS.go` supports the following functionality: list all files, search for a filename, and extract a file using `list()`, `search()`, and `extract()`, respectively. We begin by presenting the implementation of `list()`:

```markup
func list(f embed.FS) error {
    return fs.WalkDir(f, ".", walkFunction)
}
```

All the magic happens in the `walkFunction()` function, which is implemented as follows:

```markup
func walkFunction(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err
    }
    fmt.Printf("Path=%q, isDir=%v\n", path, d.IsDir())
    return nil
}
```

The `walkFunction()` function is pretty compact as all functionality is implemented by Go.

Then, we present the implementation of the `extract()` function:

```markup
func extract(f embed.FS, filepath string) ([]byte, error) {
    s, err := fs.ReadFile(f, filepath)
    if err != nil {
        return nil, err
    }
    return s, nil
}
```

The `ReadFile()` function is used for retrieving a file, which is identified by its file path, from the `embed.FS` file system as a byte slice, which is returned from the `extract()` function.

Last, we have the implementation of the `search()` function, which is based on `walkSearch()`:

```markup
func walkSearch(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err
    }
    if d.Name() == searchString {
```

`searchString` is a global variable that holds the search string. When a match is found, the matching path is printed on screen.

```markup
        fileInfo, err := fs.Stat(f, path)
        if err != nil {
            return err
        }
        fmt.Println("Found", path, "with size", fileInfo.Size())
        return nil
    }
```

Before printing a match, we make a call to `fs.Stat()` in order to get more details about it.

```markup
    return nil
}
```

The `main()` function specifically calls these three functions. Running `ioFS.go` produces the next kind of output:

```markup
$ go run ioFS.go
Path=".", isDir=true
Path="static", isDir=true
Path="static/file.txt", isDir=false
Path="static/image.png", isDir=false
Path="static/textfile", isDir=false
Found static/file.txt with size 14
wrote 14 bytes
```

Initially, the utility lists all files in the file system (lines beginning with `Path`). Then, it verifies that `static/file.txt` can be found in the file system. Last, it verifies that the writing of 14 bytes into a new file was successful as all 14 bytes have been written.

So, it turns out that Go version 1.16 introduced important additions and performance improvements.

Bookmark

# Updating the phone book application

In this section we will change the format that the phone application uses for storing its data. This time, the phone book application uses JSON all over. Additionally, it uses the `cobra` package for implementing the supported commands. As a result, all relevant code resides on its own GitHub repository and not in the `ch06` directory of the GitHub repository of this book. The path of the GitHub repository is [https://github.com/mactsouk/phonebook](https://github.com/mactsouk/phonebook)—you can `git clone` that directory if you want but try to create your own version if you have the time.

When developing real applications, do not forget to `git commit` and `git push` your changes from time to time to ensure that you keep a history of the development phase in GitHub or GitLab. Among other things, this is a good way to keep backups!

## Using cobra

First, you need to create an empty GitHub repository and clone it:

```markup
$ cd ~/go/src/github.com/mactsouk
$ git clone git@github.com:mactsouk/phonebook.git
$ cd phonebook
```

The output of the `git clone` command is not important, so it is omitted.

The first task after cloning the GitHub repository, which at this point is almost empty, is to run the `cobra init` command with the appropriate parameters.

```markup
$ ~/go/bin/cobra init --pkg-name github.com/mactsouk/phonebook
Using config file: /Users/mtsouk/.cobra.yaml
Your Cobra application is ready at
/Users/mtsouk/go/src/github.com/mactsouk/phonebook
```

Then, you should create the structure of the application using the `cobra` binary. Once you have the structure, it is easy to know what you have to implement. The structure of the application is based on the supported commands.

```markup
$ ~/go/bin/cobra add list
Using config file: /Users/mtsouk/.cobra.yaml
list created at /Users/mtsouk/go/src/github.com/mactsouk/phonebook
$ ~/go/bin/cobra add delete
$ ~/go/bin/cobra add insert
$ ~/go/bin/cobra add search
```

At this point the structure of the project should be the following:

```markup
$ tree
.
├── LICENSE
├── README.md
├── cmd
│   ├── delete.go
│   ├── insert.go
│   ├── list.go
│   ├── root.go
│   └── search.go
├── go.mod
├── go.sum
└── main.go
1 directory, 10 files
```

After that you should declare that we want to use Go modules by executing the next command:

```markup
$ go mod init
go: creating new go.mod: module github.com/mactsouk/phonebook
```

If needed you can run `go mod tidy` after `go mod init`. At this point executing `go run main.go` should download all required package dependencies and generate the default `cobra` output.

The next subsection discusses the storing of JSON data on disk.

## Storing and loading JSON data

This functionality of the `saveJSONFile()` helper function is implemented in `./cmd/root.go` using the following function:

```markup
func saveJSONFile(filepath string) error {
    f, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer f.Close()
    err = Serialize(&data, f)
    if err != nil {
        return err
    }
    return nil
}
```

So, basically, all you have to do is serialize the slice of structures using `Serialize()` and save the result into a file. Next, we need to load the JSON data from the file.

The loading functionality is implemented in `./cmd/root.go` using the `readJSONFile()` helper function. All you have to do is read the data file with the JSON data and put that data into a slice of structures by deserializing it.

## Implementing the delete command

The `delete` command deletes existing entries from the phone book application—it is implemented in `./cmd/delete.go`:

```markup
var deleteCmd = &cobra.Command{
    Use:   "delete",
    Short: "delete an entry",
    Long:  `delete an entry from the phone book application.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Get key
        key, _ := cmd.Flags().GetString("key")
        if key == "" {
            fmt.Println("Not a valid key:", key)
            return
        }
```

First, we read the appropriate command-line flag (`key`) in order to be able to identify the record that is going to be deleted.

```markup
        // Remove data
        err := deleteEntry(key)
        if err != nil {
            fmt.Println(err)
            return
        }
    },
}
```

Then, we call the `deleteEntry()` helper function to actually delete the key. After a successful deletion, `deleteEntry()` calls `saveJSONFile()` for changes to take effect.

The next subsection discusses the `insert` command.

## Implementing the insert command

The `insert` command requires user input, which means that it should support local command-line flags for doing so. As each record has three fields, the command requires three command-line flags. Then it calls the `insert()` helper function for writing the data to disk. Please refer to the `./cmd/insert.go` source file for the details of the implementation of the `insert` command.

## Implementing the list command

The `list` command lists the contents in the phone book application. It requires no command-line arguments and is basically implemented with the `list()` function:

```markup
func list() {
    sort.Sort(PhoneBook(data))
    text, err := PrettyPrintJSONstream(data)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(text)
    fmt.Printf("%d records in total.\n", len(data))
}
```

The function sorts the data before calling `PrettyPrintJSONstream()` in order to beautify the generated output.

## Implementing the search command

The `search` command is used for looking into the phone book application for a given phone number. It is implemented with the `search()` function found in `./cmd/search.go` that looks into the index map for a given key. If the key is found, then the respective record is returned.

Apart from the JSON-related operations and changes due to the use of JSON and `cobra`, all other Go code is almost the same as the version of the phone book application from _Chapter 4_, _Reflection and Interfaces_.

Working with the phone book utility produces the next kind of output:

```markup
$ go run main.go list
[
        {
                "name": "Mastering",
                "surname": "Go",
                "tel": "333123",
                "lastaccess": "1613503772"
        }
]
1 records in total.
```

This is the output from the `list` command. Adding an entry is as simple as running the next command:

```markup
$ go run main.go insert -n Mike -s Tsoukalos -t 9416471
```

Running the `list` command verifies the success of the `insert` command:

```markup
$ go run main.go list
[
        {
                "name": "Mastering",
                "surname": "Go",
                "tel": "333123",
                "lastaccess": "1613503772"
        },
        {
                "name": "Mike",
        "surname": "Tsoukalos",
        "tel": "9416471",
        "lastaccess": "1614602404"
        }
]
2 records in total.
```

Then you can delete that entry by running `go run main.go delete --key 9416471`. As stated before, the keys of the application are the phone numbers, which means that we delete entries based on phone numbers. However, nothing prohibits you from implementing deletion based on other properties.

If a command is not found, then you are going to get the next kind of output:

```markup
$ go run main.go doesNotExist
Error: unknown command "doesNotExist" for "phonebook"
Run 'phonebook --help' for usage.
Error: unknown command "doesNotExist" for "phonebook"
exit status 1
```

As the `doesNotExist` command is not supported by the command-line application, `cobra` prints a descriptive error message (`unknown command`).

Bookmark

# Exercises

-   Use the functionality of `byCharacter.go`, `byLine.go`, and `byWord.go` in order to create a simplified version of the `wc(1)` UNIX utility.
-   Create a full version of the `wc(1)` UNIX utility using the `viper` package for processing command-line options.
-   Create a full version of the `wc(1)` UNIX utility using commands instead of command-line options with the help of the `cobra` package.
-   Modify `JSONstreams.go` to accept user data or data from a file.
-   Modify `embedFiles.go` in order to save the binary file at a user-selected location.
-   Modify `ioFS.go` in order to get the desired command as well as the search string as a command-line argument.
-   Make `ioFS.go` a `cobra` project.
-   The `byLine.go` utility uses `ReadString('\n')` to read the input file. Modify the code to use `Scanner` ([https://golang.org/pkg/bufio/#Scanner](https://golang.org/pkg/bufio/#Scanner)) for reading.
-   Similarly, `byWord.go` uses `ReadString('\n')` to read the input file—modify the code to use `Scanner` instead.
-   Modify the code of `yaml.go` in order to read the YAML data from an external file.

Bookmark

# Summary

This chapter was all about systems programming and file I/O in Go and included topics such as signal handling, working with command-line arguments, reading and writing plain text files, working with JSON data, and creating powerful command-line utilities using `cobra`.

This is one of the most important chapters in this book because you cannot create any real-world utility without interacting with the operating system as well as the file system.

The next chapter is about concurrency in Go, with the main subjects being goroutines, channels, and data sharing.

Bookmark

# Additional resources

-   The `viper` package: [https://github.com/spf13/viper](https://github.com/spf13/viper)
-   The `cobra` package: [https://github.com/spf13/cobra](https://github.com/spf13/cobra)
-   The documentation of `encoding/json`: [https://golang.org/pkg/encoding/json](https://golang.org/pkg/encoding/json)
-   The documentation of `io/fs`: [https://golang.org/pkg/io/fs/](https://golang.org/pkg/io/fs/)
-   Endian-ness: [https://en.wikipedia.org/wiki/Endianness](https://en.wikipedia.org/wiki/Endianness)
-   Go 1.16 release notes: [https://golang.org/doc/go1.16](https://golang.org/doc/go1.16)