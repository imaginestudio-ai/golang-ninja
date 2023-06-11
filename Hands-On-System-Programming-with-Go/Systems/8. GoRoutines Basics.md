# Goroutines - Basic Features

In the previous chapter, you learned about Unix signal handling as well as adding support for pipes and creating graphical images in Go.

The subject of this really important chapter is goroutines. Go uses goroutines and **channels** in order to program concurrent applications in its own way while providing support for traditional concurrency techniques. Everything in Go is executed using goroutines; when a program starts its execution, its single goroutine automatically calls the main() function in order to begin the actual execution of the program.

In this chapter, we will present the easy parts of goroutines using easy to follow code examples. However, in [Chapter 10](https://subscription.imaginedevops.io/book/programming/9781787125643/10)_,_ _Goroutines - Advanced Features_, that is coming next, we will talk about more important and advanced techniques related to goroutines and channels; so, make sure that you fully understand this chapter before reading the next one.

Therefore, this chapter will tell you about the following:

-   Creating goroutines
-   Synchronizing goroutines
-   About channels and how to use them
-   Reading and writing to channels
-   Creating and using pipelines
-   Changing the Go code of the wc.go utility from _[](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_[Chapter 6](https://subscription.imaginedevops.io/book/programming/9781787125643/6), _File Input and Output_, in order to use goroutines in the new implementation
-   Improving the goroutine version of wc.go even further

Just Imagine

# About goroutines

A **goroutine** is the minimum Go entity that can be executed concurrently. Note that the use of the word _minimum_ is very important here because goroutines are not autonomous entities. Goroutines live in threads that live in Unix processes. Putting it simply, processes can be autonomous and exist on their own, whereas both goroutines and threads cannot. So, in order to create a goroutine, you will need to have a process with at least one thread. The good thing is that goroutines are lighter than threads, which are lighter than processes. Everything in Go is executed using goroutines, which makes perfect sense since Go is a concurrent programming language by design. As you have just learned, when a Go program starts its execution, its single goroutine calls the main() function, which starts the actual program execution.

You can define a new goroutine using the go keyword followed by a function name or the full definition of an anonymous function. The go keyword starts the function argument to it in a new goroutine and allows the invoking function to continue on by itself.

However, as you will see, you cannot control or make any assumptions about the order your goroutines are going to get executed because this depends on the scheduler of the operating system as well as the load of the operating system.

# Concurrency and parallelism

A very common misconception is that **concurrency** and **parallelism** refer to the same thing, which is far from true! Parallelism is the simultaneous execution of multiple things, whereas concurrency is a way of structuring your components so that they can be independently executed when possible.

Only when you build things concurrently you can safely execute them in parallel: when and if your operating system and your hardware permit it. The Erlang programming language did this a long time ago, long before CPUs had multiple cores and computers had lots of RAM.

In a valid concurrent design, adding concurrent entities makes the whole system run faster because more things can run in parallel. So, the desired parallelism comes from a better concurrent expression and implementation of the problem. The developer is responsible for taking concurrency into account during the design phase of a system and benefit from a potential parallel execution of the components of the system. So, the developer should not think about parallelism, but about breaking things into independent components that solve the initial problem when combined.

Even if you cannot run your functions in parallel on a Unix machine, a valid concurrent design will still improve the design and the maintainability of your programs. In other words, concurrency is better than parallelism!

Just Imagine

# The sync Go packages

The sync Go package contains functions that can help you synchronize goroutines; the most important functions of sync are sync.Add, sync.Done, and sync.Wait. The synchronization of goroutines is a mandatory task for every programmer.

Note that the synchronization of goroutines has nothing to do with shared variables and shared state. Shared variables and shared state have to do with the method you want to use for performing concurrent interactions.

# A simple example

In this subsection, we will present a simple program that creates two goroutines. The name of the sample program will be aGoroutine.go and will be presented in three parts; the first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "time" 
) 
 
func namedFunction() { 
   time.Sleep(10000 * time.Microsecond) 
   fmt.Println("Printing from namedFunction!") 
} 
```

Apart from the expected package and import statements, you can see the implementation of a function named namedFunction() that sleeps for a while before printing a message on the screen.

The second part of aGoroutine.go contains the following Go code:

```markup
func main() { 
   fmt.Println("Chapter 09 - Goroutines.") 
   go namedFunction() 
```

Here, you create a goroutine that executes the namedFunction() function. The last part of this naive program is the following:

```markup
   go func() { 
         fmt.Println("An anonymous function!") 
   }() 
 
   time.Sleep(10000 * time.Microsecond) 
   fmt.Println("Exiting...") 
} 
```

Here, you create another goroutine that executes an anonymous function that contains a single fmt.Println() statement.

As you can see, goroutines that run this way are totally isolated from each other and cannot exchange any kind of data, which is not always the operational style that is desired.

If you forget to call the time.Sleep() function in the main() function, or if time.Sleep() sleeps for a small amount of time, then main() will finish too early and the two goroutines will not have enough time to start and therefore finish their jobs; as a result, you will not see all the expected output on your screen!

Executing aGoroutine.go will generate the following output:

```markup
$ go run aGoroutine.go
Chapter 09 - Goroutines.
Printing from namedFunction!
Exiting... 
```

# Creating multiple goroutines

This subsection will show you how to create many goroutines and the problems that arise from having to handle more goroutines. The name of the program will be moreGoroutines.go and will be presented in three parts.

The first part of moreGoroutines.go is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "time" 
) 
```

The second part of the program has the following Go code:

```markup
func main() { 
   fmt.Println("Chapter 09 - Goroutines.") 
 
   for i := 0; i < 10; i++ { 
         go func(x int) { 
               time.Sleep(10) 
               fmt.Printf("%d ", x) 
         }(i) 
   } 
```

This time, the anonymous function takes a parameter named x, which has the value of the i variable. The for loop that uses the i variable creates ten goroutines, one by one.

The last part of the program is the following:

```markup
   time.Sleep(10000) 
   fmt.Println("Exiting...") 
} 
```

Once again, if you put a smaller value as the parameter to time.Sleep(), you will see different results when you execute the program.

Executing moreGoroutines.go will generate a somehow strange output:

```markup
$ go run moreGoroutines.go
Chapter 09 - Goroutines.
1 7 Exiting...
2 3
```

However, the big surprise comes when you execute moreGoroutines.go multiple times:

```markup
$ go run moreGoroutines.go
Chapter 09 - Goroutines.
Exiting...
$ go run moreGoroutines.go
Chapter 09 - Goroutines.
3 1 0 9 2 Exiting...
4 5 6 8 7
$ go run moreGoroutines.go
Chapter 09 - Goroutines.
2 0 1 8 7 3 6 5 Exiting...
4
```

As you can see, all previous outputs of the program are different from the first one! So, not only the output is not coordinated and there is not always enough time for all goroutines to get executed; you cannot be sure about the order the goroutines will get executed. However, although you cannot do anything about the latter problem because the order that goroutines get executed depends on various parameters that the developer cannot control, the next subsection will teach you how to synchronize goroutines and give them enough time to finish without having to call time.Sleep().

# Waiting for goroutines to finish their jobs

This subsection will demonstrate to you the correct way to make a calling function that wait for its goroutines to finish their jobs. The name of the program will be waitGR.go and will be presented in four parts; the first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "sync" 
) 
```

There is nothing special here apart from the absence of the time package and the addition of the sync package.

The second part has the following Go code:

```markup
func main() { 
   fmt.Println("Waiting for Goroutines!") 
 
   var waitGroup sync.WaitGroup 
   waitGroup.Add(10) 
```

Here, you create a new variable with a type of sync.WaitGroup, which waits for a group of goroutines to finish. The number of goroutines that belong to that group is defined by one or multiple calls to the sync.Add() function.

Calling sync.Add() before the Go statement in order to prevent race conditions is important.

Additionally, the sync.Add(10) call tells our program that we will wait for ten goroutines to finish.

The third part of the program is the following:

```markup
   var i int64 
   for i = 0; i < 10; i++ { 
 
         go func(x int64) { 
               defer waitGroup.Done() 
               fmt.Printf("%d ", x) 
         }(i) 
   } 
```

Here, you create the desired number of goroutines using a for loop, but you could have used multiple sequential Go statements. When each goroutine finishes its job, the sync.Done() function is executed: the use of the defer keyword right after the function definition tells the anonymous function to automatically call sync.Done() just before it finishes.

The last part of waitGR.go is the following:

```markup
   waitGroup.Wait() 
   fmt.Println("\nExiting...") 
} 
```

The good thing here is that there is no need to call time.Sleep() because sync.Wait() does the necessary waiting for us.

Once again, it should be noted here that you should not make any assumptions about the order the goroutines will get executed in which is also verified by the following output:

```markup
$ go run waitGR.go
Waiting for Goroutines!
9 0 5 6 7 8 2 1 3 4
Exiting...
$ go run waitGR.go
Waiting for Goroutines!
9 0 5 6 7 8 3 1 2 4
Exiting...
$ go run waitGR.go
Waiting for Goroutines!
9 5 6 7 8 1 0 2 3 4
Exiting...
```

If you call waitGroup.Add() more times than needed, you will get the following error message when you execute waitGR.go:

```markup
Waiting for Goroutines!
fatal error: all goroutines are asleep - deadlock!
    
goroutine 1 [semacquire]:
sync.runtime_Semacquire(0xc42000e28c)
      /usr/local/Cellar/go/1.8.3/libexec/src/runtime/sema.go:47 +0x34
sync.(*WaitGroup).Wait(0xc42000e280)
      /usr/local/Cellar/go/1.8.3/libexec/src/sync/waitgroup.go:131 +0x7a
main.main()
      /Users/mtsouk/ch/ch9/code/waitGR.go:22 +0x13c
exit status 2
9 0 1 2 6 7 8 3 4 5
```

This happens because when you tell your program to wait for n+1 goroutines by calling sync.Add(1) n+1 times, your program cannot have only n goroutines (or less)! Putting it simply, this will make sync.Wait() to wait indefinitely for one or more goroutines to call sync.Done() without any luck, which is obviously a deadlock situation that prevents your program from finishing.

# Creating a dynamic number of goroutines

This time, the number of goroutines that will be created will be given as a command-line argument: the name of the program will be dynamicGR.go and will be presented in four parts.

The first part of dynamicGR.go is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
   "path/filepath" 
   "strconv" 
   "sync" 
) 
```

The second part of dynamicGR.go contains the following Go code:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Printf("usage: %s integer\n",filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   numGR, _ := strconv.ParseInt(os.Args[1], 10, 64) 
   fmt.Printf("Going to create %d goroutines.\n", numGR) 
   var waitGroup sync.WaitGroup 
 
   var i int64 
   for i = 0; i < numGR; i++ { 
         waitGroup.Add(1) 
```

As you can see, the waitGroup.Add(1) statement is called just before you create a new goroutine.

The third part of the Go code of dynamicGR.go is the following:

```markup
         go func(x int64) { 
               defer waitGroup.Done() 
               fmt.Printf(" %d ", x) 
         }(i) 
   } 
```

In the preceding part, each simplistic goroutine is created.

The last part of the program is the following:

```markup
   waitGroup.Wait() 
   fmt.Println("\nExiting...") 
} 
```

Here, you just tell the program to wait for all goroutines to finish using the waitGroup.Wait() statement.

The execution of dynamicGR.go requires an integer parameter, which is the number of goroutines you want to create:

```markup
$ go run dynamicGR.go 15
Going to create 15 goroutines.
 0  2  4  1  3  5  14  10  8  9  12  11  6  13  7
Exiting...
$ go run dynamicGR.go 15
Going to create 15 goroutines.
 5  3  14  4  10  6  7  11  8  9  12  2  13  1  0
Exiting...
$ go run dynamicGR.go 15
Going to create 15 goroutines.
 4  2  3  6  5  10  9  7  0  12  11  1  14  13  8
Exiting...
```

As you can imagine, the more goroutines you want to create, the more diverse outputs you will have because there is no way to control the order that the goroutines of a program are going to be executed.

# About channels

A **channel**, putting it simply, is a communication mechanism that allows goroutines to exchange data. However, some rules exist here. First, each channel allows the exchange of a particular data type, which is also called the **element type** of the channel, and second, for a channel to operate properly, you will need to use some Go code to receive what is sent via the channel.

You should declare a new channel using the chan keyword and you can close a channel using the close() function. Additionally, as each channel has its own type, the developer should define it.

Last, a very important detail: when you are using a channel as a function parameter, you can specify its direction, that is, whether it will be used for writing or reading. In my opinion, if you know the purpose of a channel in advance, use this capability because it will make your program more robust as well as safer: otherwise, just do not define the purpose of the channel function parameter. As a result, if you declare that a channel function parameter will be used for reading only and you try to write to it, you will get an error message that will most likely save you from nasty bugs.

The error message you will get when you try to read from a write channel will be similar to the following:

```markup
# command-line-arguments
./writeChannel.go:13: invalid operation: <-c (receive from send-only type chan<- int)
```

# Writing to a channel

In this subsection, you will learn how to write to a channel. The presented program will be called writeChannel.go and you will see it in three parts.

The first part has the expected preamble:

```markup
package main 
 
import ( 
   "fmt" 
   "time" 
) 
```

As you can understand, the use of channels does not require any extra Go packages.

The second part of writeChannel.go is the following:

```markup
func writeChannel(c chan<- int, x int) { 
   fmt.Println(x) 
   c <- x 
   close(c) 
   fmt.Println(x) 
} 
```

Although the writeChannel() function writes to the channel, the data will be lost because currently nobody reads the channel in the program.

The last part of the program contains the following Go code:

```markup
func main() { 
   c := make(chan int) 
   go writeChannel(c, 10) 
   time.Sleep(2 * time.Second) 
} 
```

Here, you can see the definition of a channel variable named c with the help of the chan keyword that is used for the int data.

Executing writeChannel.go will create the following output:

```markup
 $ go run writeChannel.go
 10
```

This is not what you expected to see! The cause of this unpredicted output is that the second fmt.Println(x) statement was not executed. The reason for this is pretty simple: the c <- x statement is blocking the execution of the rest of the writeChannel() function because nobody is reading from the c channel.

# Reading from a channel

This subsection will improve the Go code of writeChannel.go by allowing you to read from a channel. The presented program will be called readChannel.go and be presented in four parts.

The first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "time" 
) 
```

The second part of readChannel.go has the following Go code:

```markup
func writeChannel(c chan<- int, x int) { 
   fmt.Println(x) 
   c <- x 
   close(c) 
   fmt.Println(x) 
} 
```

Once again, note that if nobody collects the data written to a channel, the function that sent it will stall while waiting for someone to read its data. However, in [Chapter 10](https://subscription.imaginedevops.io/book/programming/9781787125643/10)_,_ _Goroutines - Advanced Features_, you will see a very pretty solution to this problem.

The third part has the following Go code:

```markup
func main() { 
   c := make(chan int) 
   go writeChannel(c, 10) 
   time.Sleep(2 * time.Second) 
   fmt.Println("Read:", <-c) 
   time.Sleep(2 * time.Second) 
```

Here, the <-c statement in the fmt.Println() function is used for reading a single value from the channel: the same statement can be used for storing the value of a channel into a variable. However, if you do not store the value you read from a channel, it will be lost.

The last part of readChannel.go is the following:

```markup
   _, ok := <-c 
   if ok { 
         fmt.Println("Channel is open!") 
   } else { 
         fmt.Println("Channel is closed!") 
   } 
} 
```

Here, you see a technique that allows you to find out whether the channel that you want to read from is closed or not. However, if the channel was open, the presented Go code will discard the read value of the channel because of the use of the \_ character in the assignment.

Executing readChannel.go will create the following output:

```markup
$ go run readChannel.go
10
Read: 10
10
Channel is closed!
$ go run readChannel.go
10
10
Read: 10
Channel is closed!
```

# Explaining h1s.go

In [Chapter 8](https://subscription.imaginedevops.io/book/programming/9781787125643/8)_,_ _Processes and Signals_, you saw how Go handles Unix signals using many examples including h1s.go. However, now that you understand more about goroutines and channels, it is time to explain the Go code of h1s.go a little more.

As you already know that h1s.go uses channels and goroutines, it should be clear now that the anonymous function that is executed as a goroutine reads from the sigs channel using an infinite for loop. This means that each time there is a signal that interests us, the goroutine will read it from the sigs channel and handle it.

Just Imagine

# Pipelines

Go programs rarely use a single channel. One very common technique that uses multiple channels is called a **pipeline**. So, a pipeline is a method for connecting goroutines so that the output of a goroutine becomes the input of another with the help of channels. The benefits of using pipelines are as follows:

-   One of the benefits you get from using pipelines is that there is a constant flow in your program because nobody waits for everything to be completed in order to start the execution of goroutines and channels of the program
-   Additionally, you are using less variables and therefore less memory space because you do not have to save everything
-   Last, the use of pipelines simplifies the design of the program and improves its maintainability

The code of pipelines.go, which works with a pipeline of integers, will be presented in five parts; the first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
   "path/filepath" 
   "strconv" 
) 
```

The second part contains the following Go code:

```markup
func genNumbers(min, max int64, out chan<- int64) { 
 
   var i int64 
   for i = min; i <= max; i++ { 
         out <- i 
   } 
   close(out) 
} 
```

Here, you define a function that takes three arguments: two integers and one output channel. The output channel will be used for writing data that will be read in another function: this is how a pipeline is created.

The third part of the program is the following:

```markup
func findSquares(out chan<- int64, in <-chan int64) { 
   for x := range in { 
         out <- x * x 
   } 
   close(out) 
} 
```

This time, the function takes two arguments that are both channels. However, out is an output channel, whereas in is an input channel used for reading data.

The fourth part contains the definition of another function:

```markup
func calcSum(in <-chan int64) { 
   var sum int64 
   sum = 0 
   for x2 := range in { 
         sum = sum + x2 
   } 
   fmt.Printf("The sum of squares is %d\n", sum) 
} 
```

The last function of pipelines.go takes just one argument, which is a channel used for reading data.

The last part of pipelines.go is the implementation of the main() function:

```markup
func main() { 
   if len(os.Args) != 3 { 
         fmt.Printf("usage: %s n1 n2\n", filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
   n1, _ := strconv.ParseInt(os.Args[1], 10, 64) 
   n2, _ := strconv.ParseInt(os.Args[2], 10, 64) 
 
   if n1 > n2 { 
         fmt.Printf("%d should be smaller than %d\n", n1, n2) 
         os.Exit(10) 
   } 
 
   naturals := make(chan int64) 
   squares := make(chan int64) 
   go genNumbers(n1, n2, naturals) 
   go findSquares(squares, naturals) 
   calcSum(squares) 
} 
```

Here, the main() function firstly reads its two command-line arguments and creates the necessary channel variables (naturals and squares). Then, it calls the functions of the pipeline: note that the last function of the channel is not being executed as a goroutine.

The following figure shows a graphical representation of the pipeline used in pipelines.go in order to the way this particular pipeline works:

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/e6d2874d-12a1-4441-b5a1-f2af5c0056fe.png)

A graphical representation of the pipeline structure used in pipelines.go

Running pipelines.go generates the following output:

```markup
$ go run pipelines.go
usage: pipelines n1 n2
exit status 1
$ go run pipelines.go 3 2
3 should be smaller than 2
exit status 10
$ go run pipelines.go 3 20
The sum of squares is 2865
$ go run pipelines.go 1 20
The sum of squares is 2870
$ go run pipelines.go 20 20
The sum of squares is 400
```

Just Imagine

# A better version of wc.go

As we talked about in _[](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_[Chapter 6](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_,_ _File Input and Output_, in this chapter, you will learn how to create a version of wc.go that uses goroutines. The name of the new utility will be dWC.go and will be presented in four parts. Note that the current version of dWC.go considers each command-line argument as a file.

The first part of the utility is the following:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "io" 
   "os" 
   "path/filepath" 
   "regexp" 
   "sync" 
) 
```

The second part has the following Go code:

```markup
func count(filename string) { 
   var err error 
   var numberOfLines int = 0 
   var numberOfCharacters int = 0 
   var numberOfWords int = 0 
 
   f, err := os.Open(filename) 
   if err != nil { 
         fmt.Printf("%s\n", err) 
         return 
   } 
   defer f.Close() 
 
   r := bufio.NewReader(f) 
   for { 
         line, err := r.ReadString('\n') 
 
         if err == io.EOF { 
               break 
         } else if err != nil { 
               fmt.Printf("error reading file %s\n", err) 
         } 
         numberOfLines++ 
         r := regexp.MustCompile("[^\\s]+") 
         for range r.FindAllString(line, -1) { 
               numberOfWords++ 
         } 
         numberOfCharacters += len(line) 
   } 
 
   fmt.Printf("\t%d\t", numberOfLines) 
   fmt.Printf("%d\t", numberOfWords) 
   fmt.Printf("%d\t", numberOfCharacters) 
   fmt.Printf("%s\n", filename) 
} 
```

The count() function does all the processing without returning any information to the main() function: it just prints the lines, words, and characters of its input file and exits. Although the current implementation of the count() function does the desired job, it is not the correct way to design a program because there is no way to control its output of the program.

The third part of the utility is the following:

```markup
func main() { 
   if len(os.Args) == 1 { 
         fmt.Printf("usage: %s <file1> [<file2> [... <fileN]]\n", 
               filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
```

The last part of dWC.go is the following:

```markup
   var waitGroup sync.WaitGroup 
   for _, filename := range os.Args[1:] { 
         waitGroup.Add(1) 
         go func(filename string) { 
               count(filename) 
               defer waitGroup.Done() 
         }(filename) 
   } 
   waitGroup.Wait() 
} 
```

As you can see, each input file is being processed by a different goroutine. As expected, you cannot make any assumptions about the order the input files will be processed.

Executing dWC.go will generate the following output:

```markup
$ go run dWC.go /tmp/swtag.log /tmp/swtag.log doesnotExist
open doesnotExist: no such file or directory
          48    275   3571  /tmp/swtag.log
          48    275   3571  /tmp/swtag.log
  
```

Here, you can see that although the doesnotExist filename is the last command-line argument, it is the first one in the output of dWC.go!

Although dWC.go uses goroutines, there is no cleverness in it because goroutines run without communicating with each other and without performing any other tasks. Additionally, the output might get scrambled because there is no guarantee that the fmt.Printf() statements of the count() function will not get interrupted.

As a result, the forthcoming section as well as some of the techniques that will be presented in [Chapter 10](https://subscription.imaginedevops.io/book/programming/9781787125643/10)_,_ _Goroutines - Advanced Features_, will improve dWC.go.

# Calculating totals

The current version of dWC.go cannot calculate totals, which can be easily solved by processing the output of dWC.go with awk:

```markup
$ go run dWC.go /tmp/swtag.log /tmp/swtag.log | awk '{sum1+=$1; sum2+=$2; sum3+=$3} END {print "\t", sum1, "\t", sum2, "\t", sum3}'
       96    550   7142
  
```

Still, this is far from being perfect and elegant!

The main reason that the current version of dWC.go cannot calculate totals is that its goroutines have no way of communicating with each other. This can be easily solved with the help of channels and pipelines. The new version of dWC.go will be called dWCtotal.go and will be presented in five parts.

The first part of dWCtotal.go is the following:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "io" 
   "os" 
   "path/filepath" 
   "regexp" 
) 
 
type File struct { 
   Filename   string 
   Lines      int 
   Words      int 
   Characters int 
   Error      error 
} 
```

Here, a new struct type is defined. The new structure is called File and has four fields and an additional field for keeping error messages. This is the correct way for a pipeline to circulate multiple values. One might argue that a better name for the File structure would have been Counts, Results, FileCounts, or FileResults.

The second part of the program is the following:

```markup
func process(files []string, out chan<- File) { 
   for _, filename := range files { 
         var fileToProcess File 
         fileToProcess.Filename = filename 
         fileToProcess.Lines = 0 
         fileToProcess.Words = 0 
         fileToProcess.Characters = 0 
         out <- fileToProcess 
   } 
   close(out) 
} 
```

A better name of the process() function would have been beginProcess() or processResults(). You can try to make that change on your own throughout the dWCtotal.go program.

The third part of dWCtotal.go has the following Go code:

```markup
func count(in <-chan File, out chan<- File) { 
   for y := range in { 
         filename := y.Filename 
         f, err := os.Open(filename) 
         if err != nil { 
               y.Error = err 
               out <- y 
               continue 
         } 
         defer f.Close() 
         r := bufio.NewReader(f) 
         for { 
               line, err := r.ReadString('\n') 
               if err == io.EOF { 
                     break 
               } else if err != nil { 
                     fmt.Printf("error reading file %s", err) 
                     y.Error = err 
                     out <- y 
                     continue 
               } 
               y.Lines = y.Lines + 1 
               r := regexp.MustCompile("[^\\s]+") 
               for range r.FindAllString(line, -1) { 
                     y.Words = y.Words + 1 
               } 
               y.Characters = y.Characters + len(line) 
         } 
         out <- y 
   } 
   close(out) 
} 
```

Although the count() function still calculates the counts, it does not print them. It just sends the counts of lines, words, and characters as well as the filename to another channel using a struct variable of the File type.

There exists one very important detail here, which is the last statement of the count() function: in order to properly end a pipeline, you should close all involved channels, starting from the first one. Otherwise, the execution of the program will fail with an error message similar to the following one:

```markup
fatal error: all goroutines are asleep - deadlock!
```

However, as far as closing the channels of a pipeline is concerned, you should also be careful about closing channels too early, especially when there are splits in a pipeline.

The fourth part of the program contains the following Go code:

```markup
func calculate(in <-chan File) { 
   var totalWords int = 0 
   var totalLines int = 0 
   var totalChars int = 0 
   for x := range in { 
         totalWords = totalWords + x.Words 
         totalLines = totalLines + x.Lines 
         totalChars = totalChars + x.Characters 
         if x.Error == nil { 
               fmt.Printf("\t%d\t", x.Lines) 
               fmt.Printf("%d\t", x.Words) 
               fmt.Printf("%d\t", x.Characters) 
               fmt.Printf("%s\n", x.Filename) 
         } 
   } 
 
   fmt.Printf("\t%d\t", totalLines) 
   fmt.Printf("%d\t", totalWords) 
   fmt.Printf("%d\ttotal\n", totalChars) 
} 
```

There is nothing special here: the calculate() function does the dirty job of printing the output of the program.

The last part of dWCtotal.go is the following:

```markup
func main() { 
   if len(os.Args) == 1 { 
         fmt.Printf("usage: %s <file1> [<file2> [... <fileN]]\n", 
               filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   files := make(chan File)   values := make(chan File) 
 
   go process(os.Args[1:], files) 
   go count(files, values) 
   calculate(values) 
} 
```

Since the files channel is only used for passing around filenames, it could have been a string channel instead of a File channel. However, this way the code is more consistent.

Now dWCtotal.go automatically generates totals even if it has to process just one file:

```markup
$ go run dWCtotal.go /tmp/swtag.log
      48    275   3571  /tmp/swtag.log
      48    275   3571  total
$ go run dWCtotal.go /tmp/swtag.log /tmp/swtag.log doesNotExist
      48    275   3571  /tmp/swtag.log
      48    275   3571  /tmp/swtag.log
      96    550   7142  total
```

Note that both dWCtotal.go and dWC.go implement the same core functionality, which is counting the words, characters, and lines of a file: it is the way the information is handled that is different because dWCtotal.go uses a pipeline and not isolated goroutines.

[Chapter 10](https://subscription.imaginedevops.io/book/programming/9781787125643/10)_,_ _Goroutines - Advanced Features_, will use other techniques to implement the functionality of dWCtotal.go.

# Doing some benchmarking

In this section, we will compare the performance of wc.go from _[](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_[Chapter 6](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_,_ _File Input and Output__,_ with the performance of wc(1), dWC.go and dWCtotal.go. In order for the results to be more accurate, all three utilities will process relatively big files:

```markup
$ wc /tmp/*.data
  712804 3564024 9979897 /tmp/connections.data
  285316  855948 4400685 /tmp/diskSpace.data
  712523 1425046 8916670 /tmp/memory.data
 1425500 2851000 5702000 /tmp/pageFaults.data
  285658  840622 4313833 /tmp/uptime.data
 3421801 9536640 33313085 total
  
```

So, the time(1) utility will measure the following commands:

```markup
$ time wc /tmp/*.data /tmp/*.data
$ time wc /tmp/uptime.data /tmp/pageFaults.data
$ time ./dWC /tmp/*.data /tmp/*.data
$ time ./dWC /tmp/uptime.data /tmp/pageFaults.data
$ time ./dWCtotal /tmp/*.data /tmp/*.data
$ time ./dWCtotal /tmp/uptime.data /tmp/pageFaults.data
$ time ./wc /tmp/uptime.data /tmp/pageFaults.data
$ time ./wc /tmp/*.data /tmp/*.data
```

The following figure shows a graphical representation of the real field from the output of the time(1) utility when used to measure the aforementioned commands:

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/de995173-8729-4a8b-8960-775d6436074c.png)

Plotting the real field of the time(1) utility

The original wc(1) utility is by far the fastest of all. Additionally, dWC.go is faster than both dWCtotal.go and wc.go. Apart from dWC.go, the remaining two Go versions have the same performance.

Just Imagine

# Exercises

1.  Create a pipeline that reads text files, finds the number of occurrences of a given word, and calculates the total number of occurrences of the word in all files.
2.  Try to make dWCtotal.go faster.
3.  Create a simple Go program that plays ping pong using channels. You should define the total number of pings and pongs using a command-line argument.

Just Imagine

# Summary

In this chapter, we talked about creating and synchronizing goroutines as well as about creating and using pipelines and channels to allow goroutines to communicate with each other. Additionally, we developed two versions of the wc(1) utility that use goroutines to process their input files.

Make sure that you fully understand the concepts of this chapter before continuing with the next chapter because in the next chapter, we will talk about more advanced features related to goroutines and channels including shared memory, buffered channels, the select keyword, the GOMAXPROCS environment variable, and signal channels.