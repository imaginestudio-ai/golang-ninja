# Code Testing and Profiling

The topics of this chapter are both practical and important, especially if you are interested in improving the performance of your Go programs and discovering bugs. This chapter primarily addresses code _optimization_, code _testing_, and code _profiling_.

Code optimization is the process where one or more developers try to make certain parts of a program run faster, be more efficient, or use fewer resources. Put simply, code optimization is about eliminating the bottlenecks of a program that matter. Code testing is about making sure that your code does what you want it to do. In this chapter, we are experiencing the Go way of code testing. The best time to write testing code is during development, as this can help to reveal bugs in the code as early as possible. Code profiling relates to measuring certain aspects of a program to get a detailed understanding of the way the code works. The results of code profiling may help you to decide which parts of your code need to change.

Have in mind that when writing code, we should focus on its **correctness** as well as other desirable properties such as readability, simplicity, and maintainability, not its **performance**. Once we are sure that the code is correct, then we might need to focus on its performance. A good trick on performance is to execute the code on machines that are going to be a bit slower than the ones that are going to be used in production.

This chapter covers:

-   Optimizing code
-   Benchmarking code
-   Profiling code
-   The `go tool trace` utility
-   Tracing a web server
-   Testing Go code
-   Cross-compilation
-   Using `go:generate`
-   Creating example functions

Just Imagine

# Optimizing code

Code optimization is both an art and a science. This means that there is no deterministic way to help you optimize your code and that you should use your brain and try many things if you want to make your code faster. However, the general principle regarding code optimization is **first make it correct, then make it fast**. Always remember what Donald Knuth said about optimization:

> "The real problem is that programmers have spent far too much time worrying about efficiency in the wrong places and at the wrong times; premature optimization is the root of all evil (or at least most of it) in programming."

Also, remember what the late Joe Armstrong, one of the developers of Erlang, said about optimization:

> "Make it work, then make it beautiful, then if you really, really have to, make it fast. 90 percent of the time, if you make it beautiful, it will already be fast. So really, just make it beautiful!"

If you are really into code optimization, you might want to read _Compilers: Principles, Techniques, and Tools_ by Alfred V. Aho, Monica S. Lam, Ravi Sethi, and Jeffrey D. Ullman (Pearson Education Limited, 2014), which focuses on compiler construction. Additionally, all volumes in _The Art of Computer Programming_ series by Donald Knuth (Addison-Wesley Professional, 1998) are great resources for all aspects of programming if you have the time to read them.

The section that follows is about benchmarking Go code, which helps you determine what is faster and what is slower in your code—this makes it a perfect place to begin.

Just Imagine

# Benchmarking code

Benchmarking measures the performance of a function or program, allowing you to compare implementations and to understand the performance impact of code changes. Using that information, you can easily reveal the part of the code that needs to be rewritten to improve its performance. It goes without saying that you should not benchmark Go code on a busy machine that is currently being used for other, more important, purposes unless you have a very good reason to do so! Otherwise, you might interfere with the benchmarking process and get inaccurate results, but most importantly, you might generate performance issues on the machine.

Most of the time, the load of the operating system plays a key role in the performance of your code. Let me tell you a story here: a Java utility I developed for a project performs lots of computations and finishes in 6,242 seconds when running on its own. It took about a day for four instances of the same Java command-line utility to run on the same Linux machine! If you think about it, running them one after the other would have been faster than running them at the same time!

Go follows certain conventions regarding benchmarking. The most important convention is that the name of a benchmark function must begin with `Benchmark`. After the `Benchmark` word, we can put an underscore or an uppercase letter. Therefore, both `BenchmarkFunctionName()` and `Benchmark_functionName()` are valid benchmark functions whereas `Benchmarkfunctionname()` is not. The same rule applies to testing functions that begin with `Test`. Although we are allowed to put the testing and benchmarking code on the same file with the other code, it should be avoided. By convention such functions are put in files that end with `_test.go`. Once the benchmarking or the testing code is correct, the `go test` subcommand does all the dirty work for you, which includes scanning all `*_test.go` files for special functions, generating a proper temporary `main` package, calling these special functions, getting the results, and generating the final output.

Starting from Go 1.17, we can shuffle the execution order of **both tests and benchmarks** with the help of the `shuffle` parameter (`go test -shuffle=on`). The `shuffle` parameter accepts a value, which is the seed for the random number generator, and can be useful when you want to replay an execution order. Its default value is `off`. The logic behind that capability is that sometimes the order in which tests and benchmarks are executed affects their results.

## Rewriting the main() function for better testing

There exists a clever way that you can rewrite each `main()` function in order to make testing and benchmarking a lot easier. The `main()` function has a restriction, which is that you cannot call it from test code—this technique presents a solution to that problem using the code found in `main.go`. The `import` block is omitted to save space.

```markup
func main() {
    err := run(os.Args, os.Stdout)
    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }
}
```

As we cannot have an executable program without a `main()` function, we have to create a minimalistic one. What `main()` does is to call `run()`, which is our own customized version of `main()`, send `os.Args` to it, and collect the return value of `run()`.

```markup
func run(args []string, stdout io.Writer) error {
    if len(args) == 1 {
        return errors.New("No input!")
    }
    // Continue with the implementation of run()
    // as you would have with main()
    return nil
}
```

As discussed before, the `run()` function, or any other function that is called by `main()` in the same way, replaces `main()` with the additional benefit of being able to be called by test functions. Put simply, the `run()` function contains the code that would have been located in `main()`—the only difference is that `run()` returns an error variable, which is not possible with `main()`, which can only return exit codes to the operating system. You might say that this creates a slightly bigger stack because of the extra function call but the benefits are more important than this additional memory usage. Although you can omit the second parameter (`stdout io.Writer`), which is used for redirecting the generated output, the first one is important because it allows you to pass the command-line arguments to `run()`.

Running `main.go` produces the next output:

```markup
$ go run main.go 
No input!
$ go run main.go some input
```

There is nothing special in the way `main.go` operates. The good thing is that you can call `run()` from anywhere you want, including the code you write for testing, and pass the desired parameters to `run()`! It is good to have that technique in mind because it might save you.

The subject of the next subsection is benchmarking buffered writing.

## Benchmarking buffered writing and reading

In this section, we are going to test whether the size of the buffer plays a key role in the performance of write operations. This also gives us the opportunity to discuss **table tests** as well as the use of the `testdata` folder, which is reserved by Go for storing files that are going to be used during benchmarking—both table tests and the `testdata` folder can also be used in testing functions.

Benchmark functions use `testing.B` variables whereas testing functions use `testing.T` variables. It is easy to remember.

The code of `table_test.go` is the following:

```markup
package table
import (
    "fmt"
    "os"
    "path"
    "strconv"
    "testing"
)
var ERR error
var countChars int
func benchmarkCreate(b *testing.B, buffer, filesize int) {
    filename := path.Join(os.TempDir(), strconv.Itoa(buffer))
    filename = filename + "-" + strconv.Itoa(filesize)
    var err error
    for i := 0; i < b.N; i++ {
        err = Create(filename, buffer, filesize)
    }
    ERR = err
```

And now some **important information** regarding benchmarking: each benchmark function is executed for **at least** one second by default—this duration also includes the execution time of the functions that are called by a benchmark function. If the benchmark function returns in a time that is less than one second, the value of `b.N` is increased, and the function runs again as many times in total as the value of `b.N`. The first time the value of `b.N` is 1, then it becomes 2, then 5, then 10, then 20, then 50, and so on. This happens because the faster the function, the more times Go needs to run it to get _accurate results_.

The reason for storing the return value of `Create()` in a variable named `err` and using another global variable named `ERR` afterward is tricky. We want to prevent the compiler from doing any optimizations that might exclude the function that we want to measure from being executed because its results are never used.

```markup
    err = os.Remove(filename)
    if err != nil {
        fmt.Println(err)
    }
    ERR = err
}
```

Neither the signature nor the name of `benchmarkCreate()` makes it a benchmark function. This is a helper function that allows you to call `Create()`, which creates a new file on disk and its implementation can be found in `table.go`, with the proper parameters. Its implementation is valid and it can be used by benchmark functions.

```markup
func BenchmarkBuffer4Create(b *testing.B) {
    benchmarkCreate(b, 4, 1000000)
}
func BenchmarkBuffer8Create(b *testing.B) {
    benchmarkCreate(b, 8, 1000000)
}
func BenchmarkBuffer16Create(b *testing.B) {
    benchmarkCreate(b, 16, 1000000)
}
```

These are three correctly defined benchmark functions that all call `benchmarkCreate()`. Benchmark functions require a single `*testing.B` variable and return no values. In this case, the numbers at the end of the function name indicate the size of the buffer.

```markup
func BenchmarkRead(b *testing.B) {
    buffers := []int{1, 16, 96}
    files := []string{"10.txt", "1000.txt", "5k.txt"}
```

This is the code that defines the array structures that are going to be used in the table tests. This saves us from having to implement `3x3 = 9` separate benchmark functions.

```markup
    for _, filename := range files {
        for _, bufSize := range buffers {
            name := fmt.Sprintf("%s-%d", filename, bufSize)
            b.Run(name, func(b *testing.B) {
                for i := 0; i < b.N; i++ {
                    t := CountChars("./testdata/"+filename, bufSize)
                    countChars = t
                }
            })
        }
    }
}
```

The `b.Run()` method, which allows you to run one or more **sub-benchmarks** within a benchmark function, accepts two parameters. First, the name of the sub-benchmark, which is displayed onscreen, and second, the function that implements the sub-benchmark. This is the proper way to run multiple benchmarks with the use of table tests. Just remember to define a proper name for each sub-benchmark because this is going to be displayed onscreen.

Running the benchmarks generates the next output:

```markup
$ go test -bench=. *.go
```

There are two important points here: first, the value of the `-bench` parameter specifies the benchmark functions that are going to be executed. The `.` value used is a regular expression that matches all valid benchmark functions. The second point is that if you omit the `-bench` parameter, no benchmark functions are going to be executed.

```markup
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-4790K CPU @ 4.00GHz
BenchmarkBuffer4Create-8            78212      12862 ns/op
BenchmarkBuffer8Create-8           145448       7929 ns/op
BenchmarkBuffer16Create-8          222421       5074 ns/op
```

The previous three lines are the results from the `BenchmarkBuffer4Create()`, `BenchmarkBuffer8Create()`, and `BenchmarkBuffer16Create()` benchmark functions respectively, and indicate their performance.

```markup
BenchmarkRead/10.txt-1-8            78852       17268 ns/op
BenchmarkRead/10.txt-16-8           84225       14161 ns/op
BenchmarkRead/10.txt-96-8           92056       14966 ns/op
BenchmarkRead/1000.txt-1-8           2821      395419 ns/op
BenchmarkRead/1000.txt-16-8         21147       56148 ns/op
BenchmarkRead/1000.txt-96-8         58035       20362 ns/op
BenchmarkRead/5k.txt-1-8              600     1901952 ns/op
BenchmarkRead/5k.txt-16-8            4893      239557 ns/op
BenchmarkRead/5k.txt-96-8           19892       57309 ns/op
```

The previous results are from the table tests with the 9 sub-benchmarks.

```markup
PASS
ok    command-line-arguments        44.756s
```

So, what does this output tell us? First, the `-8` at the end of each benchmark function signifies the number of goroutines used for its execution, which is essentially the value of the `GOMAXPROCS` environment variable. Similarly, you can see the values of `GOOS` and `GOARCH`, which show the operating system and the architecture of the machine. The second column in the output displays the number of times that the relevant function was executed. Faster functions are executed more times than slower functions. As an example, `BenchmarkBuffer4Create()` was executed `78212` times, while `BenchmarkBuffer16Create()` was executed `222421` times because it is faster! The third column in the output shows the average time of each run and is measured in nanoseconds per benchmark function execution (`ns/op`). The bigger the value of the third column, the slower the benchmark function. A large value in the third column is an indication that a function might need to be optimized.

Should you wish to include memory allocation statistics in the output, you can include `-benchmem` in the command:

```markup
BenchmarkBuffer4Create-8    91651  11580 ns/op  304 B/op    5 allocs/op
BenchmarkBuffer8Create-8    170814 6202 ns/op   304 B/op    5 allocs/op
```

The generated output is like the one without the `-benchmem` command-line parameter but includes two additional columns. The fourth column shows the amount of memory that was allocated on average in each execution of the benchmark function. The fifth column shows the number of allocations used to allocate the memory value of the fourth column.

So far, we have learned how to create benchmark functions to test the performance of our own functions to better understand potential bottlenecks that might need to be optimized. You might ask, _how often do we need to create benchmark functions?_ The answer is simple: when something runs slower than needed and/or when you want to choose between two or more implementations.

The next subsection shows how to compare benchmark results.

## The benchstat utility

Now imagine that you have benchmarking data, and you want to compare it with the results that were produced in another computer or with a different configuration. The `benchstat` utility can help you here. The utility can be found in the [golang.org/x/perf/cmd/benchstat](http://golang.org/x/perf/cmd/benchstat) package and can be downloaded using `go get -u golang.org/x/perf/cmd/benchstat`. Go puts all binary files in `~/go/bin` and `benchstat` is no exception.

The `benchstat` utility replaces the `benchcmp` utility that can be found at [https://pkg.go.dev/golang.org/x/tools/cmd/benchcmp](https://pkg.go.dev/golang.org/x/tools/cmd/benchcmp).

So, imagine that we have two benchmark results for `table_test.go` saved in `r1.txt` and `r2.txt`—you **should remove all lines** from the `go test` output that do not contain benchmarking results, which leaves all lines that begin with `Benchmark`. You can use `benchstat` as follows:

```markup
$ ~/go/bin/benchstat r1.txt r2.txt
name                old time/op  new time/op  delta
Buffer4Create-8     10.5µs ± 0%   0.8µs ± 0%   ~     (p=1.000 n=1+1)
Buffer8Create-8     6.88µs ± 0%  0.79µs ± 0%   ~     (p=1.000 n=1+1)
Buffer16Create-8    5.01µs ± 0%  0.78µs ± 0%   ~     (p=1.000 n=1+1)
Read/10.txt-1-8     15.0µs ± 0%   4.0µs ± 0%   ~     (p=1.000 n=1+1)
Read/10.txt-16-8    12.2µs ± 0%   2.6µs ± 0%   ~     (p=1.000 n=1+1)
Read/10.txt-96-8    11.9µs ± 0%   2.6µs ± 0%   ~     (p=1.000 n=1+1)
Read/1000.txt-1-8    381µs ± 0%   174µs ± 0%   ~     (p=1.000 n=1+1)
Read/1000.txt-16-8  54.0µs ± 0%  22.6µs ± 0%   ~     (p=1.000 n=1+1)
Read/1000.txt-96-8  19.1µs ± 0%   6.2µs ± 0%   ~     (p=1.000 n=1+1)
Read/5k.txt-1-8     1.81ms ± 0%  0.89ms ± 0%   ~     (p=1.000 n=1+1)
Read/5k.txt-16-8     222µs ± 0%   108µs ± 0%   ~     (p=1.000 n=1+1)
Read/5k.txt-96-8    51.5µs ± 0%  21.5µs ± 0%   ~     (p=1.000 n=1+1)
```

If the value of the `delta` column is `~`, as it happens to be here, it means that there was no significant change in the results. The previous output shows no differences between the two results. Discussing more about `benchstat` is beyond the scope of the book. Type `benchstat -h` to learn more about the supported parameters.

The next subsection touches on a sensitive subject, which is **incorrectly defined benchmark functions**.

## Wrongly defined benchmark functions

You should be very careful when defining benchmark functions because you might define them wrongly. Look at the Go code of the following benchmark function:

```markup
func BenchmarkFiboI(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = fibo1(i)
    }
}
```

The `BenchmarkFibo()` function has a valid name and the correct signature. The bad news is that this benchmark function is **logically wrong** and is not going to produce any results. The reason for this is that as the `b.N` value grows in the way described earlier; the runtime of the benchmark function also increases because of the `for` loop. This fact prevents `BenchmarkFiboI()` from converging to a stable number, which prevents the function from completing and therefore returning any results. For analogous reasons, the next benchmark function is also wrongly implemented:

```markup
func BenchmarkfiboII(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = fibo1(b.N)
    }
}
```

On the other hand, there is nothing wrong with the implementation of the following two benchmark functions:

```markup
func BenchmarkFiboIV(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = fibo1(10)
    }
}
func BenchmarkFiboIII(b *testing.B) {
    _ = fibo1(b.N)
}
```

Correct benchmark functions are a tool for identifying bottlenecks on your code that you should put in your own projects, especially when working with file I/O or CPU-intensive operations—as I am writing this, I have been waiting **3 days** for a Python program to finish its operation to test the performance of the brute force method of a mathematical algorithm. Enough with benchmarking. The next section discusses code profiling.

Just Imagine

# Profiling code

Profiling is a process of dynamic program analysis that measures various values related to program execution to give you a better understanding of the program behavior. In this section, we are going to learn how to profile Go code to understand it better and improve its performance. Sometimes, code profiling can even reveal bugs in the code such an endless loop or functions that never return.

The `runtime/pprof` standard Go package is used for profiling all kinds of applications apart from HTTP servers. The high-level `net/http/pprof` package should be used when you want to profile a web application written in Go. You can see the help page of the `pprof` tool by executing `go tool pprof -help`.

This next section is going to illustrate how to profile a command-line application, and the following subsection shows the profiling of an HTTP server.

## Profiling a command-line application

The code of the application is saved as `profileCla.go` and collects CPU and memory profiling data. What is interesting is the implementation of `main()` because this is where the collection of the profiling data takes place:

```markup
func main() {
    cpuFilename := path.Join(os.TempDir(), "cpuProfileCla.out")
    cpuFile, err := os.Create(cpuFilename)
    if err != nil {
        fmt.Println(err)
        return
    }
    pprof.StartCPUProfile(cpuFile)
    defer pprof.StopCPUProfile()
```

The previous code is about collecting CPU profiling data. `pprof.StartCPUProfile()` starts the collecting, which is stopped with the `pprof.StopCPUProfile()` call. All data is saved into a file named `cpuProfileCla.out` under the `os.TempDir()` directory—this depends on the OS used and makes the code **portable**. The use of `defer` means that `pprof.StopCPUProfile()` is going to get called just before `main()` exits.

```markup
    total := 0
    for i := 2; i < 100000; i++ {
        n := N1(i)
        if n {
            total = total + 1
        }
    }
    fmt.Println("Total primes:", total)
    total = 0
    for i := 2; i < 100000; i++ {
        n := N2(i)
        if n {
            total = total + 1
        }
    }
    fmt.Println("Total primes:", total)
    for i := 1; i < 90; i++ {
        n := fibo1(i)
        fmt.Print(n, " ")
    }
    fmt.Println()
    for i := 1; i < 90; i++ {
        n := fibo2(i)
        fmt.Print(n, " ")
    }
    fmt.Println()
    runtime.GC()
```

All the previous code performs lots of CPU-intensive calculations for the CPU profiler to have data to collect—this is where your actual code usually goes.

```markup
    // Memory profiling!
    memoryFilename := path.Join(os.TempDir(), "memoryProfileCla.out")
    memory, err := os.Create(memoryFilename)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer memory.Close()
```

We create a second file for collecting memory-related profiling data.

```markup
    for i := 0; i < 10; i++ {
        s := make([]byte, 50000000)
        if s == nil {
            fmt.Println("Operation failed!")
        }
        time.Sleep(50 * time.Millisecond)
    }
    err = pprof.WriteHeapProfile(memory)
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

The `pprof.WriteHeapProfile()` function writes the memory data into the specified file. Once again, we allocate lots of memory for the memory profiler to have data to collect.

Running `profileCla.go` is going to create two files in the folder returned by `os.TempDir()`—usually, we save them in a different folder. Feel free to change the code of `profileCla.go` and put the profiling files at a different place. So, what do we do next? We should use `go tool pprof` to process these files:

```markup
$ go tool pprof /path/ToTemporary/Directory/cpuProfileCla.out
(pprof) top
Showing nodes accounting for 5.65s, 98.78% of 5.72s total
Dropped 47 nodes (cum <= 0.03s)
Showing top 10 nodes out of 18
      flat  flat%   sum%        cum   cum%
     3.27s 57.17% 57.17%      3.65s 63.81%  main.N2 (inline)
```

The `top` command returns a summary of the top 10 entries.

```markup
(pprof) top10 -cum
Showing nodes accounting for 5560ms, 97.20% of 5720ms total
Dropped 47 nodes (cum <= 28.60ms)
Showing top 10 nodes out of 18
      flat  flat%   sum%        cum   cum%
      80ms  1.40%  1.40%     5660ms 98.95%  main.main
         0     0%  1.40%     5660ms 98.95%  runtime.main
```

The `top10 –cum` command returns the cumulative time for each function.

```markup
(pprof) list main.N1
list main.N1
Total: 5.72s
ROUTINE ======================== main.N1 in /Users/mtsouk/ch11/profileCla.go
     1.72s      1.83s (flat, cum) 31.99% of Total
         .          .     35:func N1(n int) bool {
         .          .     36:  k := math.Floor(float64(n/2 + 1))
      50ms       60ms     37:  for i := 2; i < int(k); i++ {
     1.67s      1.77s     38:         if (n % i) == 0 {
```

Last, the `list` command shows information about a given function. The previous output shows that the `if (n % i) == 0` statement is responsible for most of the time it takes `N1()` to run.

We are not showing the full output of these commands for brevity. Try the profile commands on your own in your own code to see their full output. Visit [https://blog.golang.org/pprof](https://blog.golang.org/pprof) from the Go blog to learn more about profiling.

You can also create PDF output of the profiling data from the shell of the Go profiler using the `pdf` command. Personally, most of the time, I begin with this command because it gives me a rich overview of the collected data.

Now, let us discuss how to profile an HTTP server, which is the subject of the next subsection.

## Profiling an HTTP server

As discussed, the `net/http/pprof` package should be used when you want to collect profiling data for a Go application that runs an HTTP server. To that end, importing `net/http/pprof` **installs various handlers** under the `/debug/pprof/` path. You are going to see more on this in a short while. For now, it is enough to remember that the `net/http/pprof` package should be used to profile web applications, whereas `runtime/pprof` should be used to profile all other kinds of applications.

The technique is illustrated in `profileHTTP.go`, which comes with the following code:

```markup
package main
import (
    "fmt"
    "net/http"
    "net/http/pprof"
    "os"
    "time"
)
```

As discussed earlier, you should import the `net/http/pprof` package.

```markup
func myHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
    fmt.Printf("Served: %s\n", r.Host)
}
func timeHandler(w http.ResponseWriter, r *http.Request) {
    t := time.Now().Format(time.RFC1123)
    Body := "The current time is:"
    fmt.Fprintf(w, "%s %s", Body, t)
    fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
    fmt.Printf("Served time for: %s\n", r.Host)
}
```

The previous two functions implement two handlers that are going to be used in our naïve HTTP server. `myHandler()` is the default handler function whereas `timeHandler()` returns the current time and date on the server.

```markup
func main() {
    PORT := ":8001"
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Using default port number: ", PORT)
    } else {
        PORT = ":" + arguments[1]
        fmt.Println("Using port number: ", PORT)
    }
    r := http.NewServeMux()
    r.HandleFunc("/time", timeHandler)
    r.HandleFunc("/", myHandler)
```

Up to this point, there is nothing special as we just register the handler functions.

```markup
    r.HandleFunc("/debug/pprof/", pprof.Index)
    r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    r.HandleFunc("/debug/pprof/profile", pprof.Profile)
    r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    r.HandleFunc("/debug/pprof/trace", pprof.Trace)
```

All previous statements install the handlers for the HTTP profiler—you can access them using the hostname and port number of the web server. You do not have to use all handlers.

```markup
    err := http.ListenAndServe(PORT, r)
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

Last, you start the HTTP server as usual.

What is next? First, you run the HTTP server (`go run profileHTTP.go`). After that, you run the next command to collect profiling data **while interacting with the HTTP server**:

```markup
$ go tool pprof http://localhost:8001/debug/pprof/profile
Fetching profile over HTTP from http://localhost:8001/debug/pprof/profile
Saved profile in /Users/mtsouk/pprof/pprof.samples.cpu.004.pb.gz
Type: cpu
Time: Jun 18, 2021 at 12:30pm (EEST)
Duration: 30s, Total samples = 10ms (0.033%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) %
```

The previous output shows the initial screen of the HTTP profiler—the available commands are the same as when profiling a command-line application.

You can either exit the shell and analyze your data later using `go tool pprof` or continue giving profiler commands. This is the general idea behind profiling HTTP servers in Go.

The next subsection discusses the web interface of the Go profiler.

## The web interface of the Go profiler

The good news is that starting with Go version 1.10, `go tool pprof` comes with a web user interface that you can start as `go tool pprof -http=[host]:[port] aProfile.out`—do not forget to put the correct values to `-http`.

A part of the web interface of the profiler is seen in the next figure, which shows how program execution time was spent.

![Diagram
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_11_01.png)

Figure 11.1: The web interface of the Go profiler

Feel free to browse the web interface and see the various options that are offered. Unfortunately, talking more about profiling is beyond the scope of this chapter. As always, if you are really interested in code profiling, experiment with it as much as possible.

The next section is about code tracing.

Just Imagine

# The go tool trace utility

Code tracing is a process that allows you to learn information such as the operation of the garbage collector, the lifetime of goroutines, the activity of each logical processor, and the number of operating system threads used. The `go tool trace` utility is a tool for viewing the data stored in trace files, which can be generated in any one of the following three ways:

-   With the `runtime/trace` package
-   With the `net/http/pprof` package
-   With the `go test -trace` command

This section illustrates the use of the first technique using the code of `traceCLA.go`:

```markup
package main
import (
    "fmt"
    "os"
    "path"
    "runtime/trace"
    "time"
)
```

The `runtime/trace` package is required for collecting all kinds of tracing data—there is no point in selecting specific tracing data as all tracing data is interconnected.

```markup
func main() {
    filename := path.Join(os.TempDir(), "traceCLA.out")
    f, err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer f.Close()
```

As it happened with profiling, we need to create a file to store tracing data. In this case the file is called `traceCLA.out` and is stored inside the temporary directory of your operating system.

```markup
    err = trace.Start(f)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer trace.Stop()
```

This part is all about acquiring data for `go tool trace`, and it has nothing to do with the purpose of the program. We start the tracing process using `trace.Start()`. When we are done, we call the `trace.Stop()` function. The `defer` call means that we want to terminate tracing when the `main()` function returns.

```markup
    for i := 0; i < 3; i++ {
        s := make([]byte, 50000000)
        if s == nil {
            fmt.Println("Operation failed!")
        }
    }
    for i := 0; i < 5; i++ {
        s := make([]byte, 100000000)
        if s == nil {
            fmt.Println("Operation failed!")
        }
        time.Sleep(time.Millisecond)
    }
}
```

All the previous code is about allocating memory to trigger the operation of the garbage collector and generate more tracing data—you can learn more about the Go garbage collector in _Appendix A_, _Go Garbage Collector_. The program is executed as usual. However, when it finishes, it populates `traceCLA.out` with tracing data. After that, we should process the tracing data as follows:

```markup
$ go tool trace /path/ToTemporary/Directory/traceCLA.out
```

The last command automatically starts a web server and opens the web interface of the trace tool on your default web browser—you can run it on your own computer to play with the web interface of the trace tool.

The `View trace` link shows information about the goroutines of your program and the operation of the garbage collector.

Have in mind that although `go tool trace` is very handy and powerful, it cannot solve every kind of performance problem. There are times where `go tool pprof` is more appropriate, especially when we want to reveal where our code spends most of its time.

As it happens with profiling, collecting tracing data for an HTTP server is a slightly different process, which is explained in the next subsection.

## Tracing a web server from a client

This section shows how to trace a web server application using `net/http/httptrace`. The package allows you to trace the phases of an HTTP request from a client. The code of `traceHTTP.go` that interacts with web servers is as follows:

```markup
package main
import (
    "fmt"
    "net/http"
    "net/http/httptrace"
    "os"
)
```

As expected, we need to import `net/http/httptrace` before being able to enable HTTP tracing.

```markup
func main() {
    if len(os.Args) != 2 {
        fmt.Printf("Usage: URL\n")
        return
    }
    URL := os.Args[1]
    client := http.Client{}
    req, _ := http.NewRequest("GET", URL, nil)
```

Up to this point, we prepare the client request to the web server as usual.

```markup
    trace := &httptrace.ClientTrace{
        GotFirstResponseByte: func() {
            fmt.Println("First response byte!")
        },
        GotConn: func(connInfo httptrace.GotConnInfo) {
            fmt.Printf("Got Conn: %+v\n", connInfo)
        },
        DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
            fmt.Printf("DNS Info: %+v\n", dnsInfo)
        },
        ConnectStart: func(network, addr string) {
            fmt.Println("Dial start")
        },
        ConnectDone: func(network, addr string, err error) {
            fmt.Println("Dial done")
        },
        WroteHeaders: func() {
            fmt.Println("Wrote headers")
        },
    }
```

The preceding code is all about tracing HTTP requests. The `httptrace.ClientTrace` structure defines the events that interest us, which are `GotFirstResponseByte`, `GotConn`, `DNSDone`, `ConnectStart`, `ConnectDone`, and `WroteHeaders`. When such an event occurs, the relevant code is executed. You can find more information about the supported events and their purpose in the documentation of the `net/http/httptrace` package.

```markup
    req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
    fmt.Println("Requesting data from server!")
    _, err := http.DefaultTransport.RoundTrip(req)
    if err != nil {
        fmt.Println(err)
        return
    }
```

The `httptrace.WithClientTrace()` function returns a new `context` value based on the given parent context while `http.DefaultTransport.RoundTrip()` wraps the request with the `context` value in order to keep track of the request.

Have in mind that Go HTTP tracing has been designed to trace the events of a single `http.Transport.RoundTrip`.

```markup
    _, err = client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

The last part sends the client request to the server for the tracing to begin.

Running `traceHTTP.go` generates the next output:

```markup
$ go run traceHTTP.go https://www.golang.org/
Requesting data from server!
DNS Info: {Addrs:[{IP:2a00:1450:4001:80e::2011 Zone:} {IP:142.250.185.81 Zone:}] Err:<nil> Coalesced:false}
```

In this first part, we see that the IP address of the server has been resolved, which means that the client is ready to begin interacting with the HTTP server.

```markup
Dial start
Dial done
Got Conn: {Conn:0xc000078000 Reused:false WasIdle:false IdleTime:0s}
Wrote headers
First response byte!
Got Conn: {Conn:0xc000078000 Reused:true WasIdle:false IdleTime:0s}
Wrote headers
First response byte!
DNS Info: {Addrs:[{IP:2a00:1450:4001:80e::2011 Zone:} {IP:142.250.185.81 Zone:}] Err:<nil> Coalesced:false}
Dial start
Dial done
Got Conn: {Conn:0xc0000a1180 Reused:false WasIdle:false IdleTime:0s}
Wrote headers
First response byte!
```

The previous output helps you understand the progress of the connection in more detail and is handy when troubleshooting. Unfortunately, talking more about tracing is beyond the scope of this book. The next subsection shows how to visit all the routes of a web server to make sure that they are properly defined.

## Visiting all routes of a web server

The `gorilla/mux` package offers a `Walk` function that can be used for visiting all the registered routes of a router—this can be very handy when you want to make sure that every route is registered and is working.

The code of `walkAll.go`, which contains lots of empty handler functions because its purpose is not to test handling functions but to visit them, is as follows (nothing prohibits you from using the same technique on a fully implemented web server):

```markup
package main
import (
    "fmt"
    "net/http"
    "strings"
    "github.com/gorilla/mux"
)
```

As we are using an external package, the running of `walkAll.go` should take place somewhere in `~/go/src`.

```markup
func handler(w http.ResponseWriter, r *http.Request) {
    return
}
```

This empty handler function is shared by all endpoints for reasons of simplicity.

```markup
func (h notAllowedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    handler(rw, r)
}
```

The `notAllowedHandler` handler also calls the `handler()` function.

```markup
type notAllowedHandler struct{}
func main() {
    r := mux.NewRouter()
    r.NotFoundHandler = http.HandlerFunc(handler)
    notAllowed := notAllowedHandler{}
    r.MethodNotAllowedHandler = notAllowed
    // Register GET
    getMux := r.Methods(http.MethodGet).Subrouter()
    getMux.HandleFunc("/time", handler)
    getMux.HandleFunc("/getall", handler)
    getMux.HandleFunc("/getid", handler)
    getMux.HandleFunc("/logged", handler)
    getMux.HandleFunc("/username/{id:[0-9]+}", handler)
    // Register PUT
    // Update User
    putMux := r.Methods(http.MethodPut).Subrouter()
    putMux.HandleFunc("/update", handler)
    // Register POST
    // Add User + Login + Logout
    postMux := r.Methods(http.MethodPost).Subrouter()
    postMux.HandleFunc("/add", handler)
    postMux.HandleFunc("/login", handler)
    postMux.HandleFunc("/logout", handler)
    // Register DELETE
    // Delete User
    deleteMux := r.Methods(http.MethodDelete).Subrouter()
    deleteMux.HandleFunc("/username/{id:[0-9]+}", handler)
```

The previous part is about defining the routes and the HTTP methods that we want to support.

```markup
    err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
```

This is how we call the `Walk()` method.

```markup
        pathTemplate, err := route.GetPathTemplate()
        if err == nil {
            fmt.Println("ROUTE:", pathTemplate)
        }
        pathRegexp, err := route.GetPathRegexp()
        if err == nil {
            fmt.Println("Path regexp:", pathRegexp)
        }
        qT, err := route.GetQueriesTemplates()
        if err == nil {
            fmt.Println("Queries templates:", strings.Join(qT, ","))
        }
        qRegexps, err := route.GetQueriesRegexp()
        if err == nil {
            fmt.Println("Queries regexps:", strings.Join(qRegexps, ","))
        }
        methods, err := route.GetMethods()
        if err == nil {
            fmt.Println("Methods:", strings.Join(methods, ","))
        }
        fmt.Println()
        return nil
    })
```

For each visited route, the program collects the desired information. Feel free to remove some of the `fmt.Println()` calls if it does not help your purpose to reduce output.

```markup
    if err != nil {
        fmt.Println(err)
    }
    http.Handle("/", r)
}
```

So, the general idea behind `walkAll.go` is that you assign an empty handler to each route that you have in your server and then you call `mux.Walk()` for visiting all routes. Enabling Go modules and running `walkAll.go` generates the next output:

```markup
$ go mod init
$ go mod tidy
$ go run walkAll.go
Queries templates: 
Queries regexps: 
Methods: GET
ROUTE: /time
Path regexp: ^/time$
Queries templates: 
Queries regexps: 
Methods: GET
```

The output shows the HTTP methods that each route supports as well as the format of the path. So, the `/time` endpoint works with `GET` and its path is `/time` because the value of `Path regexp` means that `/time` is between the beginning (`^`) and the end of the path (`$`).

```markup
ROUTE: /getall
Path regexp: ^/getall$
Queries templates: 
Queries regexps: 
Methods: GET
ROUTE: /getid
Path regexp: ^/getid$
Queries templates: 
Queries regexps: 
Methods: GET
ROUTE: /logged
Path regexp: ^/logged$
Queries templates: 
Queries regexps: 
Methods: GET
ROUTE: /username/{id:[0-9]+}
Path regexp: ^/username/(?P<v0>[0-9]+)$
Queries templates: 
Queries regexps: 
Methods: GET
```

In the case of `/username`, the output includes the regular expressions associated with that endpoint that is used for selecting the value of the `id` variable.

```markup
Queries templates: 
Queries regexps: 
Methods: PUT
ROUTE: /update
Path regexp: ^/update$
Queries templates: 
Queries regexps: 
Methods: PUT
Queries templates: 
Queries regexps: 
Methods: POST
ROUTE: /add
Path regexp: ^/add$
Queries templates: 
Queries regexps: 
Methods: POST
ROUTE: /login
Path regexp: ^/login$
Queries templates: 
Queries regexps: 
Methods: POST
ROUTE: /logout
Path regexp: ^/logout$
Queries templates: 
Queries regexps: 
Methods: POST
Queries templates: 
Queries regexps: 
Methods: DELETE
ROUTE: /username/{id:[0-9]+}
Path regexp: ^/username/(?P<v0>[0-9]+)$
Queries templates: 
Queries regexps: 
Methods: DELETE
```

Although visiting the routes of a web server is a kind of testing, it is not the official Go way of testing. The main thing to look for in such output is the absence of an endpoint, the use of the wrong HTTP method, or the absence of a parameter from an endpoint.

The next section discusses the testing of Go code in more detail.

Just Imagine

# Testing Go code

The subject of this section is the testing of Go code by **writing test functions**. Software testing is a very large subject and cannot be covered in a single section of a chapter in a book. So, this section tries to present as much practical information as possible.

Go allows you to write tests for your Go code to detect bugs. However, software testing can only show **the presence** of one or more bugs, **not the absence** of bugs. This means that you can never be 100% sure that your code has no bugs!

Strictly speaking, this section is about **automated testing**, which involves writing extra code to verify whether the real code—that is, the production code—works as expected or not. Thus, the result of a test function is either `PASS` or `FAIL`. You will see how this works shortly. Although the Go approach to testing might look simple at first, especially if you compare it with the testing practices of other programming languages, it is very efficient and effective because it does not require too much of the developer's time.

You should always put the testing code in a different source file. There is no need to create a huge source file that is hard to read and maintain. Now, let us present testing by revisiting the `matchInt()` function from _Chapter 3_, _Composite Data Types_.

## Writing tests for ./ch03/intRE.go

In this subsection, we write tests for the `matchInt()` function, which was implemented in `intRE.go` back in _Chapter 3_, _Composite Data Types_. First, we create a new file named `intRE_test.go`, which is going to contain all tests. Then, we rename the package from `main` to `testRE` and remove the `main()` function—this is an optional action. After that, we must decide what we are going to test and how. The main steps in testing include writing tests for expected input, unexpected input, empty input, and edge cases. All these are going to be seen in the code. Additionally, we are going to generate random integers, convert them to strings, and use them as input for `matchInt()`. Generally speaking, a good way to test functions that works with numeric values is by using random numbers, or random values in general, as input and see how your code behaves and handles these values.

The two test functions of `intRE_test.go` are the following:

```markup
func Test_matchInt(t *testing.T) {
    if matchInt("") {
        t.Error(`matchInt("") != true`)
    }
```

The `matchInt("")` call should return `false`, so if it returns `true`, it means that the function does not work as expected.

```markup
    if matchInt("00") == false {
        t.Error(`matchInt("00") != true`)
    }
```

The `matchInt("00")` call should also return `true` because `00` is a valid integer, so if it returns `false`, it means that the function does not work as expected.

```markup
    if matchInt("-00") == false {
        t.Error(`matchInt("-00") != true`)
    }
    if matchInt("+00") == false {
        t.Error(`matchInt("+00") != true`)
    }
}
```

This first test function uses static input to test the correctness of `matchInt()`. As discussed earlier, a testing function accepts a single `*testing.T` parameter and returns no values.

```markup
func Test_with_random(t *testing.T) {
    SEED := time.Now().Unix()
    rand.Seed(SEED)
    n := strconv.Itoa(random(-100000, 19999))
    if matchInt(n) == false {
        t.Error("n = ", n)
    }
}
```

The second test function uses random but valid input to test `matchInt()`. Therefore, the given input should always pass the test. Running the two test functions with `go test` creates the next output:

```markup
$ go test -v *.go
=== RUN   Test_matchInt
--- PASS: Test_matchInt (0.00s)
=== RUN   Test_with_random
--- PASS: Test_with_random (0.00s)
PASS
ok    command-line-arguments    0.410s
```

So, all tests passed, which means that everything is fine with `matchInt()`.

The next subsection discusses the use of the `TempDir()` method.

## The TempDir function

The `TempDir()` method works with both testing and benchmarking. Its purpose is to create a temporary directory that is going to be used during testing or benchmarking. Go automatically removes that temporary directory when the test and its subtests or the benchmarks are about to finish with the help of the `CleanUp()` method—this is arranged by Go and you do not need to use and implement `CleanUp()` on your own. The exact place where the temporary directory is going to be created depends on the operating system used. On macOS, it is under `/var/folders` whereas on Linux it is under `/tmp`. We are going to illustrate `TempDir()` in the next subsection where we also talk about `Cleanup()`.

## The Cleanup() function

Although we present the `Cleanup()` method in a testing scenario, `Cleanup()` works for both testing and benchmarking. Its name reveals its purpose, which is to clean up some things that we have created when testing or benchmarking a package. However, it is us who need to tell `Cleanup()` what to do—the parameter of `Cleanup()` is a function that does the cleaning up. That function is usually implemented inline as an anonymous function, but you can also create it elsewhere and call it by its name.

The `cleanup.go` file contains a dummy function named `Foo()`—as it contains no code, there is no point in presenting it. On the other hand, all important code can be found in `cleanup_test.go`:

```markup
func myCleanUp() func() {
    return func() {
        fmt.Println("Cleaning up!")
    }
}
```

The `myCleanUp()` function is going to be used as a parameter to `CleanUp()` and should have that specific signature. Apart from the signature, you can put any kind of code in the implementation of `myCleanUp()`.

```markup
func TestFoo(t *testing.T) {
    t1 := path.Join(os.TempDir(), "test01")
    t2 := path.Join(os.TempDir(), "test02")
```

These are the paths of two directories that we are going to create.

```markup
    err := os.Mkdir(t1, 0755)
    if err != nil {
        t.Error("os.Mkdir() failed:", err)
        return
    }
```

We create a directory with `os.Mkdir()` and we specify its path. Therefore, it is our duty to delete that directory when it is no longer needed.

```markup
    defer t.Cleanup(func() {
        err = os.Remove(t1)
        if err != nil {
            t.Error("os.Mkdir() failed:", err)
        }
    })
```

After `TestFoo()` finishes, `t1` is going to be deleted by the code of the anonymous function that is passed as a parameter to `t.CleanUp()`.

```markup
    err = os.Mkdir(t2, 0755)
    if err != nil {
        t.Error("os.Mkdir() failed:", err)
        return
    }
}
```

We create another directory with `os.Mkdir()`—however, in this case we are not deleting that directory. Therefore, after `TestFoo()` finishes, `t2` is not going to be deleted.

```markup
func TestBar(t *testing.T) {
    t1 := t.TempDir()
```

Because of the use of the `t.TempDir()` method, the value (directory path) of `t1` is assigned by the operating system. Additionally, that directory path is going to be automatically deleted when the test function is about to finish.

```markup
    fmt.Println(t1)
    t.Cleanup(myCleanUp())
}
```

Here we use `myCleanUp()` as the parameter to `Cleanup()`. This is handy when you want to perform the same cleanup multiple times. Running the tests creates the next output:

```markup
$ go test -v *.go
=== RUN   TestFoo
--- PASS: TestFoo (0.00s)
=== RUN   TestBar
/var/folders/sk/ltk8cnw50lzdtr2hxcj5sv2m0000gn/T/TestBar2904465158/01
```

This is the temporary directory that was created with `TempDir()` on a macOS machine.

```markup
Cleaning up!
--- PASS: TestBar (0.00s)
PASS
ok    command-line-arguments        0.096s
```

Checking whether the directories created by `TempDir()` are there shows that they have been successfully deleted. On the other hand, the directory stored in the `t2` variable of `TestFoo()` has not been deleted. Running the same tests again is going to fail because the `test02` file already exists and cannot be created:

```markup
$ go test -v *.go
=== RUN   TestFoo
    cleanup_test.go:33: os.Mkdir() failed: mkdir /var/folders/sk/ltk8cnw50lzdtr2hxcj5sv2m0000gn/T/test02: file exists
--- FAIL: TestFoo (0.00s)
=== RUN   TestBar
/var/folders/sk/ltk8cnw50lzdtr2hxcj5sv2m0000gn/T/TestBar2113309096/01
Cleaning up!
--- PASS: TestBar (0.00s)
FAIL
FAIL  command-line-arguments        0.097s
FAIL
```

The `/var/folders/sk/ltk8cnw50lzdtr2hxcj5sv2m0000gn/T/test02: file exists` error message shows the root of the problem. The solution is to clean up your tests.

The next subsection discusses the use of the `testing/quick` package.

## The testing/quick package

There are times where you need to create testing data without human intervention. The Go standard library offers the `testing/quick` package, which can be used for **black box testing** (a software testing method that checks the functionality of an application or function without any prior knowledge of its internal working) and is somewhat related to the `QuickCheck` package found in the **Haskell** programming language—both packages implement utility functions to help you with black box testing. With the help of `testing/quick`, Go generates random values of built-in types that you can use for testing, which saves you from having to generate all these values manually.

The code of `quickT.go` is the following:

```markup
package quickT
type Point2D struct {
    X, Y int
}
func Add(x1, x2 Point2D) Point2D {
    temp := Point2D{}
    temp.X = x1.X + x2.X
    temp.Y = x1.Y + x2.Y
    return temp
}
```

The previous code implements a single function that adds two `Point2D` variables—this is the function that we are going to test.

The code of `quickT_test.go` is as follows:

```markup
package quickT
import (
    "testing"
    "testing/quick"
)
var N = 1000000
func TestWithItself(t *testing.T) {
    condition := func(a, b Point2D) bool {
        return Add(a, b) == Add(b, a)
    }
    err := quick.Check(condition, &quick.Config{MaxCount: N})
    if err != nil {
        t.Errorf("Error: %v", err)
    }
}
```

The call to `quick.Check()` **automatically generates** random numbers based on the signature of its first argument, which is a function defined earlier. There is no need to create these random numbers on your own, which makes the code easy to read and write. The actual tests happen in the `condition` function.

```markup
func TestThree(t *testing.T) {
    condition := func(a, b, c Point2D) bool {
        return Add(Add(a, b), c) == Add(a, b)
    }
```

This implementation is **wrong on purpose**. To correct the implementation, we should replace `Add(Add(a, b), c) == Add(a, b)` with `Add(Add(a, b), c) == Add(c, Add(a, b))`. We did that to see the output that is generated when a test fails.

```markup
    err := quick.Check(condition, &quick.Config{MaxCount: N})
    if err != nil {
        t.Errorf("Error: %v", err)
    }
}
```

Running the created tests generates the next output:

```markup
$ go test -v *.go
=== RUN   TestWithItself
--- PASS: TestWithItself (0.86s)
```

As expected, the first test was successful.

```markup
=== RUN   TestThree
    quickT_test.go:28: Error: #1: failed on input quickT.Point2D{X:761545203426276355, Y:-915390795717609627}, quickT.Point2D{X:-3981936724985737618, Y:2920823510164787684}, quickT.Point2D{X:-8870190727513030156, Y:-7578455488760414673}
--- FAIL: TestThree (0.00s)
FAIL
FAIL    command-line-arguments  1.153s
FAIL
```

However, as expected, the second test generated an error. The good thing is that the input that caused the error is presented onscreen so that you can see the input that caused your function to fail.

The next subsection tells us how to time out tests that take too long to finish.

## Timing out tests

If the `go test` tool takes too long to finish or for some reason it never ends, there is the `-timeout` parameter that can help you.

To illustrate that, we are using the code from the previous subsection as well as the `-timeout` and `-count` command-line flags. While the former specifies the maximum allowed time duration for the tests, the latter specifies the number of times the tests are going to be executed.

Running `go test -v *.go -timeout 1s` tells `go test` that all tests should take at most one second to finish—on my machine, the tests did take less than a second to finish. However, running the following generates a different output:

```markup
$ go test -v *.go -timeout 1s -count 2
=== RUN   TestWithItself
--- PASS: TestWithItself (0.87s)
=== RUN   TestThree
    quickT_test.go:28: Error: #1: failed on input quickT.Point2D{X:-312047170140227400, Y:-5441930920566042029}, quickT.Point2D{X:7855449254220087092, Y:7437813460700902767}, quickT.Point2D{X:4838605758154930957, Y:-7621852714243790655}
--- FAIL: TestThree (0.00s)
=== RUN   TestWithItself
panic: test timed out after 1s
```

The output is longer than the presented one—the rest of the output has to do with goroutines being terminated before they have finished. The key thing here is that the `go test` command timed out the process due to the use of `-timeout 1s`.

## Testing code coverage

In this section, we are going to learn how to find information about the code coverage of our programs to discover blocks of code or single code statements that are not being executed by testing functions.

Among other things, seeing the code coverage of programs can reveal issues and bugs in the code, so do not underestimate its usefulness. However, the code coverage test complements unit testing without replacing it. The only thing to remember is that you should make sure that the testing functions do try to cover all cases and therefore try to run all available code. If the testing functions do not try to cover all cases, then the issue might be with them, not the code that is being tested.

The code of `coverage.go`, which has some intentional issues in order to show how unreachable code is identified, is as follows:

```markup
package coverage
import "fmt"
func f1() {
    if true {
        fmt.Println("Hello!")
    } else {
        fmt.Println("Hi!")
    }
}
```

The issue with this function is that the first branch of `if` is always true and therefore the `else` branch is not going to get executed.

```markup
func f2(n int) int {
    if n >= 0 {
        return 0
    } else if n == 1 {
        return 1
    } else {
        return f2(n-1) + f2(n-2)
    }
}
```

There exist two issues with `f2()`. The first one is that it does not work well with negative integers and the second one is that all positive integers are handled by the first `if` branch. Code coverage can only help you with the second issue. The code of `coverage_test.go` is the following—these are regular test functions that try to run all available code:

```markup
package coverage
import "testing"
func Test_f1(t *testing.T) {
    f1()
}
```

This test function naively tests the operation of `f1()`.

```markup
func Test_f2(t *testing.T) {
    _ = f2(123)
}
```

The second test function checks the operation of `f2()` by running `f2(123)`.

First, we should run `go test` as follows—the code coverage task is done by the `-cover` flag:

```markup
$ go test -cover *.go
ok    command-line-arguments    0.420s    coverage: 50.0% of statements
```

The previous output shows that we have `50%` code coverage, which is not a good thing! However, we are not done yet as we can generate a test coverage report. The next command generates the code coverage report:

```markup
$ go test -coverprofile=coverage.out *.go
```

The contents of `coverage.out` are as follows—yours might vary a little depending on your username and the folder used:

```markup
$ cat coverage.out 
mode: set
/Users/mtsouk/Desktop/coverage.go:5.11,6.10 1 1
/Users/mtsouk/Desktop/coverage.go:6.10,8.3 1 1
/Users/mtsouk/Desktop/coverage.go:8.8,10.3 1 0
/Users/mtsouk/Desktop/coverage.go:13.20,14.12 1 1
/Users/mtsouk/Desktop/coverage.go:14.12,16.3 1 1
/Users/mtsouk/Desktop/coverage.go:16.8,16.19 1 0
/Users/mtsouk/Desktop/coverage.go:16.19,18.3 1 0
/Users/mtsouk/Desktop/coverage.go:18.8,20.3 1 0
```

The format and the fields in each line of the coverage file are `name.go:line.column,line.column numberOfStatements count`. The **last field** is a flag that tells you whether the statements specified by `line.column,line.column` are covered or not. So, when you see `0` in the last field, it means that the code is not covered.

Last, the HTML output can be seen in your favorite web browser by running `go tool cover -html=coverage.out`. If you used a different filename than `coverage.out`, modify the command accordingly. The next figure shows the generated output—if you are reading the printed version of the book, you might not be able to see the colors. Red lines denote code that is not being executed whereas green lines show code that was executed by the tests.

![Graphical user interface, text
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_11_02.png)

Figure 11.2: Code coverage report

Some of the code is marked as `not tracked` (gray in color) because this is code that cannot be processed by the code coverage tool. The generated output clearly shows the code issues with both `f1()` and `f2()`. You just have to correct them now!

The next section discusses unreachable code and how to discover it.

## Finding unreachable Go code

Sometimes, a wrongly implemented `if` or a misplaced `return` statement can create blocks of code that are unreachable, that is, blocks of code that are not going to be executed at all. As this is a **logical kind of error**, which means that it is not going to get caught by the compiler, we need to find a way of discovering it.

Fortunately, the `go vet` tool, which examines Go source code and reports suspicious constructs, can help with that—the use of `go vet` is illustrated with the help of the `cannotReach.go` source code file, which contains the next two functions:

```markup
func S2() {
    return
    fmt.Println("Hello!")
}
```

There is a logical error here because `S2()` returns before printing the desired message.

```markup
func S1() {
    fmt.Println("In S1()")
    return
    fmt.Println("Leaving S1()")
}
```

Similarly, `S1()` returns without giving the `fmt.Println("Leaving S1()")` statement a chance to be executed.

Running `go vet` on `cannotReach.go` creates the next output:

```markup
$ go vet cannotReach.go
# command-line-arguments
./cannotReach.go:9:2: unreachable code
./cannotReach.go:16:2: unreachable code
```

The first message points to the `fmt.Println()` statement of `S2()` and the second one to the second `fmt.Println()` statement of `S1()`. In this case, `go vet` did a great job. However, `go vet` is not particularly sophisticated and cannot catch every possible type of logical error. If you need a more advanced tool, have a look at `staticcheck` ([https://staticcheck.io/](https://staticcheck.io/)), which can also be integrated with Microsoft Visual Studio Code ([https://code.visualstudio.com/](https://code.visualstudio.com/))—the next figure shows how Visual Studio Code signifies unreachable code with the help of `staticcheck`.

![Text
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_11_03.png)

Figure 11.3: Viewing unreachable code in Visual Studio Code

As a rule of thumb, it does not hurt to include `go vet` in your workflow. You can find more information about the capabilities of `go vet` by running `go doc cmd/vet`.

The next section illustrates how to test an HTTP server with a database backend.

Just Imagine

# Testing an HTTP server with a database backend

An HTTP server is a different kind of animal because it should already run for tests to get executed. Thankfully, the `net/http/httptest` package can help—you do not need to run the HTTP server on your own as the `net/http/httptest` package does the work for you, but you need to have the database server up and running. We are going to test the REST API server we have developed in _Chapter 10_, _Working with REST APIs_—we are going to copy the `server_test.go` file with the test code in the [https://github.com/mactsouk/rest-api](https://github.com/mactsouk/rest-api) GitHub repository of the server.

To create `server_test.go`, we do not have to change the implementation of the REST API server.

The code of `server_test.go`, which holds the test functions for the HTTP service, is the following:

```markup
package main
import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "strconv"
    "strings"
    "testing"
    "time"
    "github.com/gorilla/mux"
)
```

The only reason for including `github.com/gorilla/mux` is the use of `mux.SetURLVars()` later on.

```markup
func TestTimeHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/time", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(TimeHandler)
    handler.ServeHTTP(rr, req)
    status := rr.Code
    if status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}
```

The `http.NewRequest()` function is used for defining the HTTP request method, the endpoint, and for sending data to the endpoint when needed. The `http.HandlerFunc(TimeHandler)` call defines the handler function that is being tested.

```markup
func TestMethodNotAllowed(t *testing.T) {
    req, err := http.NewRequest("DELETE", "/time", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(MethodNotAllowedHandler)
```

We are testing `MethodNotAllowedHandler` in this test function.

```markup
    handler.ServeHTTP(rr, req)
    status := rr.Code
    if status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}
```

We know that this interaction is going to fail as we are testing `MethodNotAllowedHandler`. Therefore, we expect to get an `http.StatusNotFound` response code back—if we get a different code, the test function is going to fail.

```markup
func TestLogin(t *testing.T) {
    UserPass := []byte(`{"Username": "admin", "Password": "admin"}`)
```

Here we store the desired fields of a `User` structure in a `byte` slice. For the tests to work, the `admin` user should have `admin` as the password because this is what is used in the code—modify `server_test.go` in order to have the correct password for the `admin` user, or any other user with admin privileges, of your installation.

```markup
    req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(UserPass))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")
```

The previous lines of code construct the desired request.

```markup
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(LoginHandler)
    handler.ServeHTTP(rr, req)
```

`NewRecorder()` returns an initialized `ResponseRecorder` that is used in `ServeHTTP()`—`ServeHTTP()` is the method that performs the request. The response is saved in the `rr` variable.

There is also a test function for the `/logout` endpoint, which is not presented here as it is almost identical to `TestLogin()`. However, running the tests in random order might create issues with testing because `TestLogin()` should always get executed before `TestLogout()`.

```markup
    status := rr.Code
    if status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
        return
    }
}
```

If the status code is `http.StatusOK`, it means that the interaction worked as expected.

```markup
func TestAdd(t *testing.T) {
    now := int(time.Now().Unix())
    username := "test_" + strconv.Itoa(now)
    users := `[{"Username": "admin", "Password": "admin"}, {"Username":"` + username + `", "Password": "myPass"}]`
```

For the `Add()` handler, we need to pass an array of JSON records, which is constructed here. As we do not want to create the same username every time, we append the current timestamp to the `_test` string.

```markup
    UserPass := []byte(users)
    req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(UserPass))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")
```

This is where we construct the slice of JSON records (`UserPass`) and create the request.

```markup
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(AddHandler)
    handler.ServeHTTP(rr, req)
    // Check the HTTP status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
        return
    }
}
```

If the server response is `http.StatusOK`, then the request is successful and the test passes.

```markup
func TestGetUserDataHandler(t *testing.T) {
    UserPass := []byte(`{"Username": "admin", "Password": "admin"}`)
    req, err := http.NewRequest("GET", "/username/1", bytes.NewBuffer(UserPass))
```

Although we use `/username/1` in the request, this does not add any value in the `Vars` map. Therefore, we need to use the `SetURLVars()` function for changing the values in the `Vars` map—this is illustrated next.

```markup
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")
    vars := map[string]string{
        "id": "1",
    }
    req = mux.SetURLVars(req, vars)
```

The `gorilla/mux` package provides the `SetURLVars()` function for testing purposes—this function allows you to add elements to the `Vars` map. In this case, we need to set the value of the `id` key to `1`. You can add as many key/value pairs as you want.

```markup
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(GetUserDataHandler)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
        return
    }
expected := `{"ID":1,"Username":"admin","Password":"admin",
"LastLogin":0,"Admin":1,"Active":0}`
```

This is the record we expect to get back from our request. As we cannot guess the value of `LastLogin` in the server response, we replace it with `0`, hence the use of `0` here.

```markup
    serverResponse := rr.Body.String()
    result := strings.Split(serverResponse, "LastLogin")
    serverResponse = result[0] + `LastLogin":0,"Admin":1,"Active":0}`
```

As we do not want to use the value of `LastLogin` from the server response, we are changing it to `0`.

```markup
    if serverResponse != expected {
        t.Errorf("handler returned unexpected body: got %v but wanted %v",
            rr.Body.String(), expected)
    }
}
```

The last part of the code contains the standard Go way of checking whether we have received the expected answer or not.

Creating tests for HTTP services is easy once you understand the presented examples. This mainly happens because most of the code is repeated among test functions.

Running the tests generates the next output:

```markup
$ go test -v server_test.go main.go handlers.go
=== RUN   TestTimeHandler
2021/06/17 08:59:15 TimeHandler Serving: /time from
--- PASS: TestTimeHandler (0.00s)
```

This is the output from visiting the `/time` endpoint. Its result is `PASS`.

```markup
=== RUN   TestMethodNotAllowed
2021/06/17 08:59:15 Serving: /time from  with method DELETE
--- PASS: TestMethodNotAllowed (0.00s)
=== RUN   TestLogin
```

This is the output from visiting the `/time` endpoint with the `DELETE` HTTP method. Its result is `PASS` because we were expecting this request to fail as it uses the wrong HTTP method.

```markup
2021/06/17 08:59:15 LoginHandler Serving: /login from
2021/06/17 08:59:15 Input user: {0 admin admin 0 0 0}
2021/06/17 08:59:15 Found user: {1 admin admin 1620922454 1 0}
2021/06/17 08:59:15 Logging in: {1 admin admin 1620922454 1 0}
2021/06/17 08:59:15 Updating user: {1 admin admin 1623909555 1 1}
2021/06/17 08:59:15 Affected: 1
2021/06/17 08:59:15 User updated: {1 admin admin 1623909555 1 1}
--- PASS: TestLogin (0.01s)
```

This is the output from `TestLogin()` that tests the `/login` endpoint. All lines beginning with the date and time are generated by the REST API server and show the progress of the request.

```markup
=== RUN   TestLogout
2021/06/17 08:59:15 LogoutHandler Serving: /logout from
2021/06/17 08:59:15 Found user: {1 admin admin 1620922454 1 1}
2021/06/17 08:59:15 Logging out: admin
2021/06/17 08:59:15 Updating user: {1 admin admin 1620922454 1 0}
2021/06/17 08:59:15 Affected: 1
2021/06/17 08:59:15 User updated: {1 admin admin 1620922454 1 0}
--- PASS: TestLogout (0.01s)
```

This is the output from `TestLogout()` that tests the `/logout` endpoint, which also has the `PASS` result.

```markup
=== RUN   TestAdd
2021/06/17 08:59:15 AddHandler Serving: /add from
2021/06/17 08:59:15 [{0 admin admin 0 0 0} {0 test_1623909555 myPass 0 0 0}]
--- PASS: TestAdd (0.01s)
```

This is the output from the `TestAdd()` test function. The name of the new user that is created is `test_1623909555` and it should be different each time the test is executed.

```markup
=== RUN   TestGetUserDataHandler
2021/06/17 08:59:15 GetUserDataHandler Serving: /username/1 from
2021/06/17 08:59:15 Found user: {1 admin admin 1620922454 1 0}
--- PASS: TestGetUserDataHandler (0.00s)
PASS
ok    command-line-arguments        (cached)
```

Last, this is the output from the `TestGetUserDataHandler()` test function that was also executed without any issues.

The next subsection discusses fuzzing, which offers a different way of testing.

Just Imagine

# Fuzzing

As software engineers, we do not worry when things go as expected but when unexpected things happen. One way to deal with the unexpected is fuzzing. _Fuzzing_ (or _fuzz testing_) is a testing technique that provides invalid, unexpected, or random data on programs that require input.

The advantages of fuzzing include the following:

-   Making sure that the code can handle invalid or random input
-   Bugs that are discovered with fuzzing are usually severe and might indicate security risks
-   Attackers often use fuzzing for locating vulnerabilities, so it is good to be prepared

Fuzzing is going to be officially included in the Go language in a future Go release, but do not expect it in 2021. It is most likely going to be officially released with Go version 1.18 or Go version 1.19. The `dev.fuzz` branch at GitHub ([https://github.com/golang/go/tree/dev.fuzz](https://github.com/golang/go/tree/dev.fuzz)) contains the latest implementation of fuzzing. This branch is going to exist until the relevant code is merged to the master branch. With fuzzing comes the `testing.F` data type, in the same way that we use `testing.T` for testing and `testing.B` for benchmarking. If you want to try fuzzing in Go, begin by visiting [https://blog.golang.org/fuzz-beta](https://blog.golang.org/fuzz-beta).

The next section discusses a handy Go feature, which is cross-compilation.

Just Imagine

# Cross-compilation

Cross-compilation is the process of generating a binary executable file for a different architecture than the one on which we are working without having access to other machines. The main benefit that we receive from cross-compilation is that we do not need a second or third machine to create and distribute executable files for different architectures. This means that we basically need just a single machine for our development. Fortunately, Go has built-in support for cross-compilation.

To cross-compile a Go source file, we need to set the `GOOS` and `GOARCH` environment variables to the target operating system and architecture, respectively, which is not as difficult as it sounds.

You can find a list of available values for the `GOOS` and `GOARCH` environment variables at [https://golang.org/doc/install/source](https://golang.org/doc/install/source). Keep in mind, however, that not all `GOOS` and `GOARCH` combinations are valid.

The code of `crossCompile.go` is the following:

```markup
package main
import (
    "fmt"
    "runtime"
)
func main() {
    fmt.Print("You are using ", runtime.GOOS, " ")
    fmt.Println("on a(n)", runtime.GOARCH, "machine")
    fmt.Println("with Go version", runtime.Version())
}
```

Running it on a macOS machine with Go version 1.16.5 generates the next output:

```markup
$ go run crossCompile.go
You are using darwin on a(n) amd64 machine
with Go version go1.16.5
```

Compiling `crossCompile.go` **for the Linux OS** that runs on a machine with an `amd64` processor is as simple as running the next command on a macOS machine:

```markup
$ env GOOS=linux GOARCH=amd64 go build crossCompile.go
$ file crossCompile
crossCompile: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, Go BuildID=GHF99KZkGUrFADRlsS7l/ty-Ka44KVhMItrIvMZ6l/rdRP5mt_yw2AEox_8uET/HqP0KyUBaOB87LY7gvVu, not stripped
```

Transferring that file to an Arch Linux machine and running it generates the next output:

```markup
$ ./crossCompile 
You are using linux on a(n) amd64 machine
with Go version go1.16.5
```

One thing to notice here is that the cross-compiled binary file of `crossCompile.go` prints the Go version of the machine used for compiling it—this makes perfect sense as the target machine might not even have Go installed on it.

Cross-compilation is a great Go feature that can come in handy when you want to generate multiple versions of your executables through a CI/CD system and distribute them.

The next section discusses `go:generate`.

Just Imagine

# Using go:generate

Although `go:generate` is not directly connected to testing or profiling, it is a handy and advanced Go feature and I believe that this chapter is the perfect place for discussing it as it can also help you with testing. The `go:generate` directive is associated with the `go generate` command, was added in Go 1.4 in order to help with automation, and allows you to run commands described by directives within existing files.

The `go generate` command supports the `-v`, `-n`, and `-x` flags. The `-v` flag prints the names of packages and files as they are processed whereas the `-n` flag prints the commands that would be executed. Last, the `-x` flag prints commands as they are executed—this is great for debugging `go:generate` commands.

The main reasons that you might need to use `go:generate` are the following:

-   You want to download dynamic data from the Internet or some other source prior to the execution of the Go code.
-   You want to execute some code prior to running the Go code.
-   You want to generate a version number or other unique data before code execution.
-   You want to make sure that you have sample data to work with. For example, you can put data into a database using `go:generate`.

As using `go:generate` is not considered a good practice because it hides things from the developer and creates additional dependencies, I try to avoid it when I can, and I usually can. On the other hand, if you really need it, you are going to know it!

The use of `go:generate` is illustrated in `goGenerate.go`, which has the following content:

```markup
package main
import "fmt"
//go:generate ./echo.sh
```

This executes the `echo.sh` script, which should be available in the current directory.

```markup
//go:generate echo GOFILE: $GOFILE
//go:generate echo GOARCH: $GOARCH
//go:generate echo GOOS: $GOOS
//go:generate echo GOLINE: $GOLINE
//go:generate echo GOPACKAGE: $GOPACKAGE
```

`$GOFILE`, `$GOARCH`, `$GOOS`, `$GOLINE`, and `$GOPACKAGE` are special variables and are translated at the time of execution.

```markup
//go:generate echo DOLLAR: $DOLLAR
//go:generate echo Hello!
//go:generate ls -l
//go:generate ./hello.py
```

This executes the `hello.py` Python script, which should be available in the current directory.

```markup
func main() {
    fmt.Println("Hello there!")
}
```

The `go generate` command is not going to run the `fmt.Println()` statement or any other statements found in a Go source file. Last, have in mind that `go generate` is not executed automatically and must be run explicitly.

Working with `goGenerate.go` from within `~/go/src/` generates the next output:

```markup
$ go mod init
$ go mod tidy
$ go generate
Hello world!
GOFILE: goGenerate.go
GOARCH: amd64
GOOS: darwin
GOLINE: 9
GOPACKAGE: main
```

This is the output of the `$GOFILE`, `$GOARCH`, `$GOOS`, `$GOLINE`, and `$GOPACKAGE` variables, which shows the values of these variables defined at runtime.

```markup
DOLLAR: $
```

There is also a special variable named `$DOLLAR` for printing a dollar character in the output because `$` has a special meaning in the OS environment.

```markup
Hello!
total 32
-rwxr-xr-x  1 mtsouk  staff   32 Jun  2 18:18 echo.sh
-rw-r--r--  1 mtsouk  staff   45 Jun  2 16:15 go.mod
-rw-r--r--  1 mtsouk  staff  381 Jun  2 18:18 goGenerate.go
-rwxr-xr-x  1 mtsouk  staff   52 Jun  2 18:18 hello.py
drwxr-xr-x  5 mtsouk  staff  160 Jun  2 17:07 walk
```

This is the output of the `ls -l` command that shows the files found in the current directory at the time of the code execution. This can be used for testing whether some necessary files are present at the time of execution or not.

```markup
Hello from Python!
```

Last is the output of a naïve Python script.

Running `go generate` with `-n` shows the commands that are going to be executed:

```markup
$ go generate -n
./echo.sh
echo GOFILE: goGenerate.go
echo GOARCH: amd64
echo GOOS: darwin
echo GOLINE: 9
echo GOPACKAGE: main
echo DOLLAR: $
echo Hello!
ls -l
./hello.py
```

So, `go:generate` can help you work with the OS before program execution. However, as it hides things from the developer, its usage should be limited.

The last section of this chapter talks about example functions.

Just Imagine

# Creating example functions

Part of the documentation process is generating example code that showcases the use of some or all the functions and data types of a package. _Example functions_ have many benefits, including the fact that they are executable tests that are executed by `go test`. Therefore, if an example function contains an `// Output:` line, the `go test` tool checks whether the calculated output matches the values found after the `// Output:` line. Although we should include example functions in Go files that end with `_test.go`, we do not need to import the `testing` Go package for example functions. Moreover, the name of each example function must begin with `Example`. Lastly, **example functions take no input parameters and return no results**.

We are going to illustrate example functions using the code of `exampleFunctions.go` and `exampleFunctions_test.go`. The content of `exampleFunctions.go` is as follows:

```markup
package exampleFunctions
func LengthRange(s string) int {
    i := 0
    for _, _ = range s {
        i = i + 1
    }
    return i
}
```

The previous code presents a regular package that contains a single function named `LengthRange()`. The contents of `exampleFunctions_test.go`, which includes the example functions, are the following:

```markup
package exampleFunctions
import "fmt"
func ExampleLengthRange() {
    fmt.Println(LengthRange("Mihalis"))
    fmt.Println(LengthRange("Mastering Go, 3rd edition!"))
    // Output:
    // 7
    // 7
}
```

What the comment lines say is that the expected output is `7` and `7`, which is obviously wrong. This is going to be seen after we run `go test`:

```markup
$ go test -v exampleFunctions*
=== RUN   ExampleLengthRange
--- FAIL: ExampleLengthRange (0.00s)
got:
7
26
want:
7
7
FAIL
FAIL    command-line-arguments  0.410s
FAIL
```

As expected, there is an error in the generated output—the second generated value is `26` instead of the expected `7`. If we make the necessary corrections, the output is going to look as follows:

```markup
$ go test -v exampleFunctions*
=== RUN   ExampleLengthRange
--- PASS: ExampleLengthRange (0.00s)
PASS
ok      command-line-arguments  1.157s
```

Example functions can be a great tool both for learning the capabilities of a package and for testing the correctness of functions, so I suggest that you include both test code and example functions in your Go packages. As a bonus, your test functions **appear in the documentation of the package**, if you decide to generate package documentation.

Just Imagine

# Exercises

-   Implement a simple version of `ab(1)` ([https://httpd.apache.org/docs/2.4/programs/ab.html](https://httpd.apache.org/docs/2.4/programs/ab.html)) on your own using goroutines and channels for testing the performance of web services.
-   Write test functions for the `phoneBook.go` application from _Chapter 3_, _Composite Data Types_.
-   Create test functions for a package that calculates numbers in the Fibonacci sequence. Do not forget to implement that package.
-   Try to find the value of `os.TempDir()` in various operating systems.
-   Create three different implementations of a function that copies binary files and benchmark them to find the faster one. Can you explain why this function is faster?

Just Imagine

# Summary

This chapter discussed `go:generate`, code profiling and tracing, benchmarking, and testing Go code. You might find the Go way of testing and benchmarking boring, but this happens because **Go is boring and predictable** in general and that is a good thing! Remember that writing bug-free code is important whereas writing the fastest code possible is not always that important.

Most of the time, you need to be able to write **fast enough** code. So, **spend more time writing tests than benchmarks** unless your code runs really slowly. We have also learned how to find unreachable code and how to cross-compile Go code.

Although the discussions of the Go profiler and `go tool trace` are far from complete, you should understand that with topics such as profiling and code tracing, nothing can replace experimenting and trying new techniques on your own!

The next chapter is about creating gRPC services in Go.

Just Imagine

# Additional resources

-   The `generate` package: [https://golang.org/pkg/cmd/go/internal/generate/](https://golang.org/pkg/cmd/go/internal/generate/)
-   Generating code: [https://blog.golang.org/generate](https://blog.golang.org/generate)
-   Look at the code of `testing` at [https://golang.org/src/testing/testing.go](https://golang.org/src/testing/testing.go)
-   About `net/http/httptrace`: [https://golang.org/pkg/net/http/httptrace/](https://golang.org/pkg/net/http/httptrace/)
-   Introducing HTTP Tracing by Jaana Dogan: [https://blog.golang.org/http-tracing](https://blog.golang.org/http-tracing)
-   GopherCon 2019: Dave Cheney - Two Go Programs, Three Different Profiling Techniques: [https://youtu.be/nok0aYiGiYA](https://youtu.be/nok0aYiGiYA)