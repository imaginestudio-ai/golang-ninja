# Goroutines - Advanced Features

This is the second chapter of this book that deals with goroutines: the most important feature of the Go programming language: as well as channels that greatly improve what goroutines can do, and we will continue this from where we stopped it in [Chapter 9](https://subscription.imaginedevops.io/book/programming/9781787125643/9)_,_ _Goroutines - Basic Features_.

Thus, you will learn how to use various types of channels, including buffered channels, signal channels, nil channels, and channels of channels! Additionally, you will learn how you can utilize shared memory and mutexes with goroutines as well as how to time out a program when it is taking too long to finish.

Specifically, this chapter will discuss the following topics:

-   Buffered channels
-   The select keyword
-   Signal channels
-   Nil channels
-   Channel of channels
-   Timing out a program and avoiding waiting forever for it to end
-   Shared memory and goroutines
-   Using sync.Mutex in order to guard shared data
-   Using sync.RWMutex in order to protect your shared data
-   Changing the code of dWC.go from [Chapter 9](https://subscription.imaginedevops.io/book/programming/9781787125643/9), _Goroutines - Basic Features_, in order to add support for buffered channels and mutexes to it

Just Imagine

# The Go scheduler

In the previous chapter, we said that the kernel scheduler is responsible for the order your goroutines will be executed in, which is not completely accurate. The kernel scheduler is responsible for the execution of the threads your programs have. The Go runtime has its own scheduler that is responsible for the execution of the goroutines using a technique known as **m:n scheduling**, where _m_ goroutines are executed using _n_ operating system threads using multiplexing. As the Go scheduler has to deal with the goroutines of a single program, its operation is much cheaper and faster than the operation of the kernel scheduler.

Just Imagine

# The sync Go package

Once again, we will use functions and data types from the sync package in this chapter. Particularly, you will learn about the usefulness of the sync.Mutex and sync.RWMutex types and the functions supporting them.

Just Imagine

# The select keyword

A select statement in Go is like a switch statement for channels and allows a goroutine to wait on multiple communication operations. Therefore, the main advantage you get from using the select keyword is that the same function can deal with multiple channels using a single select statement! Additionally, you can have nonblocking operations on channels.

The name of the program that will be used for illustrating the select keyword will be useSelect.go and will be presented in five parts. The useSelect.go program allows you to generate the number of random you want, which is defined in the first command-line argument, up to a certain limit, which is the second command-line argument.

The first part of useSelect.go is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "math/rand" 
   "os" 
   "path/filepath" 
   "strconv" 
   "time" 
) 
```

The second part of useSelect.go is the following:

```markup
func createNumber(max int, randomNumberChannel chan<- int, finishedChannel chan bool) { 
   for { 
         select { 
         case randomNumberChannel <- rand.Intn(max): 
         case x := <-finishedChannel: 
               if x { 
                     close(finishedChannel) 
                     close(randomNumberChannel) 
                     return 
               } 
         } 
   } 
}
```

Here, you can see how the select keyword allows you to listen to and coordinate two channels (randomNumberChannel and finishedChannel) at the same time. The select statement waits for a channel to unblock and then executes on that.

The for loop of the createNumber() function will not end on this own. Therefore, createNumber() will keep generating random numbers for as long as the randomNumberChannel branch of the select statement is used. The createNumber() function will exit when it gets the Boolean value true in the finishedChannel channel.

A better name for the finishedChannel channel would have been done or even noMoreData.

The third part of the program contains the following Go code:

```markup
func main() { 
   rand.Seed(time.Now().Unix()) 
   randomNumberChannel := make(chan int) 
   finishedChannel := make(chan bool) 
 
   if len(os.Args) != 3 { 
         fmt.Printf("usage: %s count max\n", filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   n1, _ := strconv.ParseInt(os.Args[1], 10, 64) 
   count := int(n1) 
   n2, _ := strconv.ParseInt(os.Args[2], 10, 64) 
   max := int(n2) 
 
   fmt.Printf("Going to create %d random numbers.\n", count) 
```

There is nothing special here: you just read the command-line arguments before starting the desired goroutine.

The fourth part of useSelect.go is where you will start the desired goroutine and create a for loop in order to generate the desired number of random numbers:

```markup
   go createNumber(max, randomNumberChannel, finishedChannel) 
   for i := 0; i < count; i++ { 
         fmt.Printf("%d ", <-randomNumberChannel) 
   } 
 
   finishedChannel <- false 
   fmt.Println() 
   _, ok := <-randomNumberChannel 
   if ok { 
         fmt.Println("Channel is open!") 
   } else { 
         fmt.Println("Channel is closed!") 
   } 
```

Here, you also send a message to finishedChannel and check whether the randomNumberChannel channel is open or closed after sending the message to finishedChannel. As you sent false to finishedChannel, the finishedChannel channel will remain open. Note that a message sent to a closed channel panics, whereas a message received from a closed channel returns the zero value immediately.

Note that once you close a channel, you cannot write to this channel. However, you can still read from that channel!

The last part of useSelect.go has the following Go code:

```markup
   finishedChannel <- true   _, ok = <-randomNumberChannel 
   if ok { 
         fmt.Println("Channel is open!") 
   } else { 
         fmt.Println("Channel is closed!") 
   } 
} 
```

Here, you sent the true value to finishedChannel, so your channels will close and the createNumber() goroutine will exit.

Running useSelect.go will create the following output:

```markup
$ go run useSelect.go 2 100
Going to create 2 random numbers.
19 74
Channel is open!
Channel is closed!
```

As you will see in the bufChannels.go program that explains buffered channels, the select statement can also save you from overflowing a buffered channel.

Just Imagine

# Signal channels

A **signal channel** is a channel that is used just for signaling. Signal channels will be illustrated using the signalChannel.go program with a rather unusual example that will be presented in five parts. The program executes four goroutines: when the first one is finished, it sends a signal to a signal channel by closing it, which will unblock the second goroutine. When the second goroutine finishes its job, it closes another channel that unblocks the remaining two goroutines. Note that signal channels are not the same as channels that carry the os.Signal values.

The first part of the program is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "time" 
) 
 
func A(a, b chan struct{}) { 
   <-a 
   fmt.Println("A!") 
   time.Sleep(time.Second) 
   close(b) 
} 
```

The A() function is blocked by the channel defined in the a parameter. This means that until this channel is closed, the A() function cannot continue its execution. The last statement of the function closes the channel that is stored in the b variable, which will be used for unblocking other goroutines.

The second part of the program is the implementation of the B() function:

```markup
func B(b, c chan struct{}) { 
   <-b 
   fmt.Println("B!") 
   close(c) 
} 
```

Similarly, the B() function is blocked by the channel stored in the b argument, which means that until the b channel is closed, the B() function will be waiting in its first statement.

The third part of signalChannel.go is the following:

```markup
func C(a chan struct{}) { 
   <-a 
   fmt.Println("C!") 
} 
```

Once again, the C() function is blocked by the channel stored in its a argument.

The fourth part of the program is the following:

```markup
func main() { 
   x := make(chan struct{}) 
   y := make(chan struct{}) 
   z := make(chan struct{})
```

Defining a signal channel as an empty struct with no fields is a very common practice because empty structures take no memory space. In such a case, you could have used a bool channel instead.

The last part of signalChannel.go has the following Go code:

```markup
   go A(x, y) 
   go C(z) 
   go B(y, z) 
   go C(z) 
 
   close(x) 
   time.Sleep(2 * time.Second) 
} 
```

Here, you start four goroutines. However, until you close the a channel, all of them will be blocked! Additionally, A() will finish first and unblock B() that will unblock the two C() goroutines. So, this technique allows you to define the order of execution of your goroutines.

If you execute signalChannel.go, you will get the following output:

```markup
$ go run signalChannel.go
A!
B!
C!
C!
```

As you can see, the goroutines are being executed in the desired order despite the A() function taking more time to execute than the others due to the time.Sleep() function call.

Just Imagine

# Buffered channels

**Buffered channels** allow the Go scheduler to put jobs in the queue quickly in order to be able to serve more requests. Moreover, you can use buffered channels as **semaphores** in order to limit throughput. The technique works as follows: incoming requests are forwarded to a channel, which processes one request at a time. When the channel is done, it sends a message to the original caller saying that it is ready to process a new request. So, the capacity of the buffer of the channel restricts the number of simultaneous requests it can keep and process: this can be easily implemented using a for loop with a call to time.Sleep() at its end.

Buffered channels will be illustrated in bufChannels.go, which will be presented in four parts.

The first part of the program is the following:

```markup
package main 
 
import ( 
   "fmt" 
) 
```

The preamble proves that you do not need any extra packages for supporting buffered channels in your Go program.

The second part of the program has the following Go code:

```markup
func main() { 
   numbers := make(chan int, 5) 
```

Here, you create a new channel named numbers with 5 places, which is denoted by the last parameter of the make statement. This means that you can write five integers to that channel without having to read any one of them in order to make space for the others. However, you cannot put six integers on a channel with five integer places!

The third part of bufChannels.go is the following:

```markup
   counter := 10 
   for i := 0; i < counter; i++ { 
         select { 
         case numbers <- i: 
         default: 
               fmt.Println("Not enough space for", i) 
         } 
   } 
```

Here, you try to put 10 integers to a buffered channel with 5 places. However, the use of the select statement allows you to know whether you have enough space for storing all the integers or not and act accordingly!

The last part of bufChannels.go is the following:

```markup
   for i := 0; i < counter*2; i++ { 
         select { 
         case num := <-numbers: 
               fmt.Println(num) 
         default:
               fmt.Println("Nothing more to be done!")    
               break 
         } 
   } 
} 
```

Here, you also use a select statement while trying to read 20 integers from a channel. However, as soon as reading from the channel fails, the for loop exits using a break statement. This happens because when there is nothing left to read from the numbers channel, the num := <-numbers statement will block, which makes the case statement to go to the default branch.

As you can see from the code, there is no goroutine in bufChannels.go, which means that buffered channels can work on their own.

Executing bufChannels.go will generate the following output:

```markup
$ go run bufChannels.go
Not enough space for 5
Not enough space for 6
Not enough space for 7
Not enough space for 8
Not enough space for 9
0
1
2
3
4
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
Nothing more to be done!
```

Just Imagine

# About timeouts

Can you imagine waiting forever for something to perform an action? Neither can I! So, in this section you will learn how to implement **timeouts** in Go with the help of the select statement.

The program with the sample code will be named timeOuts.go and will be presented in four parts; the first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "time" 
) 
```

The second part of timeOuts.go is the following:

```markup
func main() { 
   c1 := make(chan string) 
   go func() { 
         time.Sleep(time.Second * 3) 
         c1 <- "c1 OK" 
   }() 
```

The time.Sleep() statement in the goroutine is used for simulating the time it will take for the goroutine to do its real job.

The third part of timeOuts.go has the following code:

```markup
   select { 
   case res := <-c1: 
         fmt.Println(res) 
   case <-time.After(time.Second * 1): 
         fmt.Println("timeout c1") 
   } 
```

This time the use of time.After() is required for declaring the time you want to wait before timing out. The wonderful thing here is that if the time of time.After() expires without the select statement having received any data from the c1 channel, the case branch of time.After() will get executed.

The last part of the program will have the following Go code:

```markup
   c2 := make(chan string) 
   go func() { 
         time.Sleep(time.Second * 3) 
         c2 <- "c2 OK" 
   }() 
 
   select { 
   case res := <-c2: 
         fmt.Println(res) 
   case <-time.After(time.Second * 4): 
         fmt.Println("timeout c2") 
   } 
} 
```

In the previous code, you see an operation that does not time out because it is completed within the desired time, which means that the first branch of the select block will get executed instead of the second one that signifies the timeout.

The execution of timeOuts.go will generate the following output:

```markup
$ go run timeOuts.go
timeout c1
c2 OK
```

# An alternative way to implement timeouts

The technique of this subsection will let you not wait for any stubborn goroutines to finish their jobs. Therefore, this subsection will show you how to time out goroutines with the help of the timeoutWait.go program that will be presented in four parts. Despite the code differences between timeoutWait.go and timeOuts.go, the general idea is exactly the same.

The first part of timeoutWait.go contains the expected preamble:

```markup
package main 
 
import ( 
   "fmt" 
   "sync" 
   "time" 
) 
```

The second part of timeoutWait.go is the following:

```markup
func timeout(w *sync.WaitGroup, t time.Duration) bool { 
   temp := make(chan int) 
   go func() { 
         defer close(temp) 
         w.Wait() 
   }() 
 
   select { 
   case <-temp: 
         return false 
   case <-time.After(t): 
         return true 
   } 
} 
```

Here, you declare a function that does the entire job. The core of the function is the select block that works the same way as in timeOuts.go. The anonymous function of timeout() will successfully end when the w.Wait() statement returns, which will happen when the appropriate number of sync.Done() calls have been executed, which means that all goroutines will be finished. In this case, the first case of the select statement will be executed.

Note that the temp channel is needed in the select block and nowhere else. Additionally, the element type of the temp channel could have been anything, including bool.

The third part of timeOuts.go has the following code:

```markup
func main() { 
   var w sync.WaitGroup 
   w.Add(1) 
 
   t := 2 * time.Second 
   fmt.Printf("Timeout period is %s\n", t) 
 
   if timeout(&w, t) { 
         fmt.Println("Timed out!") 
   } else { 
         fmt.Println("OK!") 
   } 
```

The last fragment of the program has the following Go code:

```markup
   w.Done() 
   if timeout(&w, t) { 
         fmt.Println("Timed out!") 
   } else { 
         fmt.Println("OK!") 
   } 
} 
```

After the anticipated w.Done() call has been executed, the timeout() function will return true, which will prevent the timeout from happening.

As mentioned at the beginning of this subsection, timeoutWait.go actually prevents your program from having to wait indefinitely for one or more goroutines to end.

Executing timeoutWait.go will create the following output:

```markup
$ go run timeoutWait.go
Timeout period is 2s
Timed out!
OK!
```

Just Imagine

# Channels of channels

In this section, we will talk about creating and using a channel of channels. Two possible reasons to use such a channel are as follows:

-   For acknowledging that an operation finished its job
-   For creating many worker processes that will be controlled by the same channel variable

The name of the naive program that will be developed in this section is cOfC.go and will be presented in four parts.

The first part of the program is the following:

```markup
package main 
 
import ( 
   "fmt" 
) 
 
var numbers = []int{0, -1, 2, 3, -4, 5, 6, -7, 8, 9, 10} 
```

The second part of the program is the following:

```markup
func f1(cc chan chan int, finished chan struct{}) { 
   c := make(chan int) 
   cc <- c 
   defer close(c) 
 
   total := 0 
   i := 0 
   for { 
         select { 
         case c <- numbers[i]: 
               i = i + 1 
               i = i % len(numbers) 
               total = total + 1 
         case <-finished: 
               c <- total 
               return 
         } 
   } 
} 
```

The f1() function returns integer numbers that belong to the numbers variable. When it is about to end, it also returns the number of integers it has sent back to the caller function using the c <- total statement.

As you cannot use a channel of channels directly, you should first read from it (cc <- c) and get a channel that you can actually use. The handy thing here is that although you can close the c channel, the channel of channels (cc) will be still up and running.

The third part of cOfC.go is the following:

```markup
func main() { 
   c1 := make(chan chan int) 
   f := make(chan struct{}) 
 
   go f1(c1, f) 
   data := <-c1 
```

In this Go code, you can see that you can declare a channel of channels using the chan keyword two consecutive times.

The last part of cOfC.go has the following Go code:

```markup
   i := 0 
   for integer := range data { 
         fmt.Printf("%d ", integer) 
         i = i + 1 
         if i == 100 { 
               close(f) 
         } 
   } 
   fmt.Println() 
} 
```

Here, you limit the number of integers that will be created by closing the f channel when you have the number of integers you want.

Executing cOfC.go will generate the following output:

```markup
$ go run cOfC.go
0 -1 2 3 -4 5 6 -7 8 9 10 0 -1 2 3 -4 5 6 -7 8 9 10 0 -1 2 3 -4 5 6 -7 8 9 10 0 -1 2 3 -4 5 6 -7 8 9 10 0 -1 2 3 -4 5 6 -7 8 9 10 0 -1 2 3 -4 5 6 -7 8 9 10 0 -1 2 3 -4 5 6 -7 8 9 10 0 -1 2 3 -4 5 6 -7 8 9 10 0 -1 2 3 -4 5 6 -7 8 9 10 0 100
```

A channel of channels is an advanced Go feature that you probably will not need to use in your system software. However, it is good to know that it exists.

Just Imagine

# Nil channels

This section will talk about **nil channels**, which are a special sort of channel that will always block. The name of the program will be nilChannel.go and will be presented in four parts.

The first part of the program contains the expected preamble:

```markup
package main 
 
import ( 
   "fmt" 
   "math/rand" 
   "time" 
) 
```

The second portion contains the implementation of the addIntegers() function:

```markup
func addIntegers(c chan int) { 
   sum := 0 
   t := time.NewTimer(time.Second) 
 
   for { 
         select { 
         case input := <-c: 
               sum = sum + input 
         case <-t.C: 
               c = nil 
               fmt.Println(sum) 
         } 
   } 
} 
```

The addIntegers() function stops after the time defined in the time.NewTimer() function passes and will go to the relevant branch of the case statement. There, it makes c a nil channel, which means that the channel will stop receiving new data and that the function will just wait there.

The third part of nilChannel.go is the following:

```markup
func sendIntegers(c chan int) { 
   for { 
         c <- rand.Intn(100) 
   } 
} 
```

Here, the sendIntegers() function keeps generating random numbers and sends them to the c channel as long as the c channel is open. However, here you also have a goroutine that is never cleaned up.

The last part of the program has the following Go code:

```markup
func main() { 
   c := make(chan int) 
   go addIntegers(c) 
   go sendIntegers(c) 
   time.Sleep(2 * time.Second) 
} 
```

Executing nilChannel.go will generate the following output:

```markup
$ go run nilChannel.go
162674704
$ go run nilChannel.go
165021841
```

Just Imagine

# Shared memory

Shared memory is the traditional way that threads use for communicating with each other. Go comes with built-in synchronization features that allow a single goroutine to own a shared piece of data. This means that other goroutines must send messages to this single goroutine that owns the shared data, which prevents the corruption of the data! Such a goroutine is called a **monitor goroutine**. In Go terminology, this is s_haring by communicating instead of communicating by sharing_.

This technique will be illustrated in the sharedMem.go program, which will be presented in five parts. The first part of sharedMem.go has the following Go code:

```markup
package main 
 
import ( 
   "fmt" 
   "math/rand" 
   "sync" 
   "time" 
) 
```

The second part is the following:

```markup
var readValue = make(chan int) 
var writeValue = make(chan int) 
 
func SetValue(newValue int) { 
   writeValue <- newValue 
} 
 
func ReadValue() int { 
   return <-readValue 
} 
```

The ReadValue() function is used for reading the shared variable, whereas the SetValue() function is used for setting the value of the shared variable. Also, the two channels used in the program need to be global variables in order to avoid passing them as arguments to all the functions of the program. Note that these global variables are usually wrapped up in a Go library or a struct with methods.

The third part of sharedMem.go is the following:

```markup
func monitor() { 
   var value int 
   for { 
         select { 
         case newValue := <-writeValue: 
               value = newValue 
               fmt.Printf("%d ", value) 
         case readValue <- value: 
         } 
   } 
} 
```

The logic of sharedMem.go can be found in the implementation of the monitor() function. When you have a read request, the ReadValue() function attempts to read from the readValue channel. Then, the monitor() function returns the current value that is kept in the value parameter. Similarly, when you want to change the stored value, you call SetValue(), which writes to the writeValue channel that is also handled by the select statement. Once again, the select block plays a key role because it orchestrates the operations of the monitor() function.

The fourth portion of the program has the following Go code:

```markup
func main() { 
   rand.Seed(time.Now().Unix()) 
   go monitor() 
   var waitGroup sync.WaitGroup 
 
   for r := 0; r < 20; r++ { 
         waitGroup.Add(1) 
         go func() { 
               defer waitGroup.Done() 
               SetValue(rand.Intn(100)) 
         }() 
   } 
```

The last part of the program is the following:

```markup
   waitGroup.Wait() 
   fmt.Printf("\nLast value: %d\n", ReadValue()) 
} 
```

Executing sharedMem.go will generate the following output:

```markup
$ go run sharedMem.go
33 45 67 93 33 37 23 85 87 23 58 61 9 57 20 61 73 99 42 99
Last value: 99
$ go run sharedMem.go
71 66 58 83 55 30 61 73 94 19 63 97 12 87 59 38 48 81 98 49
Last value: 49
```

If you want to share more values, you can define a new structure that will hold the desired variables with the data types you prefer.

# Using sync.Mutex

**Mutex** is an abbreviation for **mutual exclusion**; the Mutex variables are mainly used for thread synchronization and for protecting shared data when multiple writes can occur at the same time. A mutex works like a buffered channel of capacity 1 that allows at most one goroutine to access a shared variable at a time. This means that there is no way for two or more goroutines to try to update that variable simultaneously. Although this is a perfectly valid technique, the general Go community prefers to use the monitor goroutine technique presented in the previous section.

In order to use sync.Mutex, you will have to declare a sync.Mutex variable first. You can lock that variable using the Lock method and release it using the Unlock method. The sync.Lock() method gives you exclusive access over the shared variable for a region of code that finishes when you call the Unlock() method and is called a **critical section**.

Each critical section of a program cannot be executed without locking it first using sync.Lock(). However, if a lock has already been taken, everybody should wait for its release first. Although multiple functions might wait to get a lock, only one of them will get it when it will be released.

You should try to make critical sections as small as possible; in other words, do not delay releasing a lock because other goroutines might want to use it. Additionally, forgetting to unlock Mutex will most likely result in a deadlock.

The name of the Go program with the code for illustrating the use of sync.Mutex will be mutexSimple.go and will be presented in five chunks.

The first part of mutexSimple.go contains the expected preamble:

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

The second part of the program is the following:

```markup
var aMutex sync.Mutex 
var sharedVariable string = "" 
 
func addDot() { 
   aMutex.Lock() 
   sharedVariable = sharedVariable + "." 
   aMutex.Unlock() 
} 
```

Note that a critical section is not always obvious and you should be very careful when specifying it. Also note that a critical section cannot be embedded in another critical section when both critical sections use the same Mutex variable! Putting it simply, avoid, at almost all costs, spreading mutexes across functions because that makes really hard to see whether you are embedding or not!

Here, addDot() adds a dot character at the end of the sharedVariable string. However, as the string should be altered simultaneously by multiple goroutines, you use a sync.Mutex variable to protect it. As the critical section contains just one command, the waiting period for getting access to the mutex will be fairly small, if not instantaneous. However, in a real-world situation, the waiting period might be much longer, especially on software such as database servers where many things happen simultaneously by thousands of processes: you can simulate that by adding a call to time.Sleep() in the critical section.

Note that it is the responsibility of the developer to associate a mutex with one or more shared variables!

The third code segment of mutexSimple.go is the implementation of another function that uses the mutex:

```markup
func read() string { 
   aMutex.Lock() 
   a := sharedVariable 
   aMutex.Unlock() 
   return a 
} 
```

Although locking the shared variable while reading it is not absolutely necessary, this kind of locking prevents the shared variable from changing while you are reading it. This might look like a small issue here but imagine reading the balance of your bank account instead!

The fourth part is where you define the number of goroutines that you will start:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Printf("usage: %s n\n", filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   numGR, _ := strconv.ParseInt(os.Args[1], 10, 64) 
   var waitGroup sync.WaitGroup 
```

The final part of mutexSimple.go contains the following Go code:

```markup
   var i int64 
   for i = 0; i < numGR; i++ { 
         waitGroup.Add(1) 
         go func() { 
               defer waitGroup.Done() 
               addDot() 
         }() 
   } 
   waitGroup.Wait() 
   fmt.Printf("-> %s\n", read()) 
   fmt.Printf("Length: %d\n", len(read())) 
} 
```

Here, you start the desired number of goroutines. Each goroutine calls the addDot() function that accesses the shared variable: and you wait for them to finish before reading the value of the shared variable using the read() function.

The output you will get from executing mutexSimple.go will be similar to the following:

```markup
$ go run mutexSimple.go 20
-> ....................
Length: 20
$ go run mutexSimple.go 30
-> ..............................
Length: 30
```

# Using sync.RWMutex

Go offers another type of mutex, called sync.RWMutex, that allows multiple readers to hold the lock but only a single writer - sync.RWMutex is an extension of sync.Mutex that adds two methods named sync.RLock and sync.RUnlock, which are used for locking and unlocking for reading purposes. Locking and unlocking a sync.RWMutex for exclusive writing should be done with Lock() and Unlock(), respectively.

This means that either one writer can hold the lock or multiple readers: not both! You will most likely use such a mutex when most of the goroutines want to read a variable and you do not want goroutines to wait in order to get an exclusive lock.

In order to demystify sync.RWMutex a little, you should discover that the sync.RWMutex type is a Go structure currently defined as follows:

```markup
type RWMutex struct { 
   w           Mutex 
   writerSem   uint32 
   readerSem   uint32  
   readerCount int32 
   readerWait  int32 
}                
```

So, there is nothing to be afraid of here! Now, it is time to see a Go program that uses sync.RWMutex. The program will be named mutexRW.go and will be presented in five parts.

The first part of mutexRW.go contains with the expected preamble as well as the definition of a global variable and a new struct type:

```markup
package main 
 
import ( 
   "fmt" 
   "sync" 
   "time" 
) 
 
var Password = secret{counter: 1, password: "myPassword"} 
 
type secret struct { 
   sync.RWMutex 
   counter  int 
   password string 
} 
```

The secret structure embeds sync.RWMutex and therefore it can call all the methods of sync.RWMutex.

The second part of mutexRW.go has the following Go code:

```markup
func Change(c *secret, pass string) { 
   c.Lock() 
   fmt.Println("LChange") 
   time.Sleep(20 * time.Second) 
   c.counter = c.counter + 1 
   c.password = pass 
   c.Unlock() 
} 
```

This function makes changes to one of its arguments, which means that it requires an exclusive lock, hence the use of the Lock() and Unlock() functions.

The third part of the sample code is the following:

```markup
func Show(c *secret) string { 
   fmt.Println("LShow") 
   time.Sleep(time.Second) 
   c.RLock() 
   defer c.RUnlock() 
   return c.password 
} 
 
func Counts(c secret) int { 
   c.RLock() 
   defer c.RUnlock() 
   return c.counter 
} 
```

Here, you can see the definition of two functions that use an sync.RWMutex for reading. This means that multiple instances of them can get the sync.RWMutex lock.

The fourth portion of the program is the following:

```markup
func main() { 
   fmt.Println("Pass:", Show(&Password)) 
   for i := 0; i < 5; i++ { 
         go func() { 
               fmt.Println("Go Pass:", Show(&Password)) 
         }() 
   } 
```

Here, you start five goroutines in order to make things more interesting and random.

The last part of mutexRW.go is the following:

```markup
   go func() { 
         Change(&Password, "123456") 
   }() 
 
   fmt.Println("Pass:", Show(&Password)) 
   time.Sleep(time.Second) 
   fmt.Println("Counter:", Counts(Password)) 
} 
```

Although shared memory and the use of a mutex are still a valid approach to concurrent programming, using goroutines and channels is a more modern way that follows the Go philosophy. Therefore, if you can solve a problem using channels and pipelines, you should prefer that way instead of using shared variables.

Executing mutexRW.go will generate the following output:

```markup
$ go run mutexRW.go
LShow
Pass: myPassword
LShow
LShow
LShow
LShow
LShow
LShow
LChange
Go Pass: 123456
Go Pass: 123456
Pass: 123456
Go Pass: 123456
Go Pass: 123456
Go Pass: 123456
Counter: 2
```

If the implementation of Change() was using a RLock() call as well as a RUnlock() call, which would have been totally wrong, then the output of the program would have been the following:

```markup
$ go run mutexRW.go
LShow
Pass: myPassword
LShow
LShow
LShow
LShow
LShow
LShow
LChange
Go Pass: myPassword
Pass: myPassword
Go Pass: myPassword
Go Pass: myPassword
Go Pass: myPassword
Go Pass: myPassword
Counter: 1
```

Put simply, you should be fully aware of the locking mechanism you are using and the way it works. In this case, it is the timing that is deciding what Counts() will return: the timing depends on the time.Sleep() call of the Change() function that emulates the processing that will happen in a real function. The problem is that the use of RLock() and RUnlock() in Change() allows multiple goroutines to read the shared variable and therefore get the wrong output from the Counts() function.

Just Imagine

# The dWC.go utility revisited

In this section, we will change the implementation of the dWC.go utility developed in the previous chapter.

The first version of the program will use a buffered channel whereas the second version of the program will use shared memory for keeping the counts for each file you process.

# Using a buffered channel

The name of this implementation will be WCbuffered.go and will be presented in five parts.

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
) 
 
type File struct { 
   Filename   string 
   Lines      int 
   Words      int 
   Characters int 
   Error      error 
} 
```

The File structure will keep the counts for each input file. The second chunk of WCbuffered.go has the following Go code:

```markup
func monitor(values <-chan File, count int) { 
   var totalWords int = 0 
   var totalLines int = 0 
   var totalChars int = 0 
   for i := 0; i < count; i++ { 
         x := <-values 
         totalWords = totalWords + x.Words 
         totalLines = totalLines + x.Lines 
         totalChars = totalChars + x.Characters 
         if x.Error == nil { 
               fmt.Printf("\t%d\t", x.Lines) 
               fmt.Printf("%d\t", x.Words) 
               fmt.Printf("%d\t", x.Characters) 
               fmt.Printf("%s\n", x.Filename) 
         } else { 
               fmt.Printf("\t%s\n", x.Error) 
         } 
   } 
   fmt.Printf("\t%d\t", totalLines) 
   fmt.Printf("%d\t", totalWords) 
   fmt.Printf("%d\ttotal\n", totalChars) 
} 
```

The monitor() function collects all the information and prints it. The for loop inside monitor() makes sure that it will collect the right amount of data.

The third part of the program contains the implementation of the count() function:

```markup
func count(filename string, out chan<- File) { 
   var err error 
   var nLines int = 0 
   var nChars int = 0 
   var nWords int = 0 
 
   f, err := os.Open(filename) 
   defer f.Close() 
   if err != nil { 
         newValue := File{ 
Filename: filename, 
Lines: 0, 
Characters: 0, 
Words: 0, 
Error: err } 
         out <- newValue 
         return 
   } 
 
   r := bufio.NewReader(f) 
   for { 
         line, err := r.ReadString('\n') 
 
         if err == io.EOF { 
               break 
         } else if err != nil { 
               fmt.Printf("error reading file %s\n", err) 
         } 
         nLines++ 
         r := regexp.MustCompile("[^\\s]+") 
         for range r.FindAllString(line, -1) { 
               nWords++ 
         } 
         nChars += len(line) 
   } 
   newValue := File { 
Filename: filename, 
Lines: nLines, 
Characters: nChars, 
Words: nWords, 
Error: nil }   out <- newValue 
} 
```

When the count() function is done, it sends the information to the buffered channel, so there is nothing special here.

The fourth portion of WCbuffered.go is the following:

```markup
func main() { 
   if len(os.Args) == 1 { 
         fmt.Printf("usage: %s <file1> [<file2> [... <fileN]]\n", 
               filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   values := make(chan File, len(os.Args[1:])) 
```

Here, you create a buffered channel named values with as many places as the number of files you will process.

The last portion of the utility is the following:

```markup
   for _, filename := range os.Args[1:] {         go func(filename string) { 
               count(filename, values) 
         }(filename) 
   } 
   monitor(values, len(os.Args[1:])) 
} 
```

# Using shared memory

The good thing with shared memory and mutexes is that, in theory, they usually take a very small amount of the code, which means that the rest of the code can work concurrently without any other delays. However, only after you have implemented something can you see what really happens!

The name of this implementation will be WCshared.go and will be presented in five parts: the first part of the utility is the following:

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
 
type File struct { 
   Filename   string 
   Lines      int 
   Words      int 
   Characters int 
   Error      error 
} 
 
var aM sync.Mutex 
var values = make([]File, 0) 
```

The values slice will be the shared variable of the program whereas the name of the mutex variable will be aM.

The second chunk of WCshared.go has the following Go code:

```markup
func count(filename string) { 
   var err error 
   var nLines int = 0 
   var nChars int = 0 
   var nWords int = 0 
 
   f, err := os.Open(filename) 
   defer f.Close() 
   if err != nil { 
         newValue := File{Filename: filename, Lines: 0, Characters: 0, Words: 0, Error: err} 
         aM.Lock() 
         values = append(values, newValue) 
         aM.Unlock() 
         return 
   } 
 
   r := bufio.NewReader(f) 
   for { 
         line, err := r.ReadString('\n') 
 
         if err == io.EOF { 
               break 
         } else if err != nil { 
               fmt.Printf("error reading file %s\n", err) 
         } 
         nLines++ 
         r := regexp.MustCompile("[^\\s]+") 
         for range r.FindAllString(line, -1) { 
               nWords++ 
         } 
         nChars += len(line) 
   } 
 
   newValue := File{Filename: filename, Lines: nLines, Characters: nChars, Words: nWords, Error: nil} 
   aM.Lock() 
   values = append(values, newValue) 
   aM.Unlock() 
} 
```

So, just before the count() function exits, it adds an element to the values slice using a critical section.

The third part of WCshared.go is the following:

```markup
func main() { 
   if len(os.Args) == 1 { 
         fmt.Printf("usage: %s <file1> [<file2> [... <fileN]]\n", 
               filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
```

Here, you just deal with the command-line arguments of the utility.

The fourth part of WCshared.go contains the following Go code:

```markup
   var waitGroup sync.WaitGroup 
   for _, filename := range os.Args[1:] { 
         waitGroup.Add(1) 
         go func(filename string) { 
               defer waitGroup.Done() 
               count(filename) 
         }(filename) 
   } 
 
   waitGroup.Wait()
```

Here, you just start the desired number of goroutines and wait for them to finish their jobs.

The last code slice of the utility is the following:

```markup
   var totalWords int = 0 
   var totalLines int = 0 
   var totalChars int = 0 
   for _, x := range values { 
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

When all goroutines are done, it is time to process the contents of the shared variable, calculate totals, and print the desired output. Note that in this case, there is no shared variable of any kind and therefore there is no need for a mutex: you just wait to gather all results and print them.

# More benchmarking

This section will measure the performance of WCbuffered.go and WCshared.go using the handy time(1) utility. However, this time, instead of presenting a graph, I will give you the actual output of the time(1) utility:

```markup
$ time go run WCshared.go /tmp/*.data /tmp/*.data
real  0m31.836s
user  0m31.659s
sys   0m0.165s
$ time go run WCbuffered.go /tmp/*.data /tmp/*.data
real  0m31.823s
user  0m31.656s
sys   0m0.171s
```

As you can see, both utilities performed equally well, or equally badly if you prefer! However, apart from the speed of a program, what also matters is the clarity of its design and how easy it is to make code changes to it! Additionally, the presented way also times the compile times of both utilities, which might make the results less accurate.

The reason that both programs can easily generate totals is that they both have a control point. For the WCshared.go utility, the control point is the shared variable, whereas for WCbuffered.go, the control point is the buffered channel that collects the desired information inside the monitor() function.

Just Imagine

# Detecting race conditions

If you use the \-race flag when running or building a Go program, you will turn on the Go **race detector**, which makes the compiler create a modified version of the typical executable file. This modified version can record the accesses to shared variables as well as all synchronization events that take place, including calls to sync.Mutex, sync.WaitGroup, and so on. After doing some analysis of the events, the race detector prints a report that can help you identify potential problems so that you can correct them.

In order to showcase the operation of the race detector, we will use the code of the rd.go program, which will be presented in four parts. For this particular program, the **data race** will happen because two or more goroutines access the same variable concurrently and at least one of them changes the value of the variable in some way.

Note that the main() program is also a goroutine in Go!

The first part of the program is the following:

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

Nothing special here, so if there is a problem with the program, it is not in the preamble.

The second part of rd.go is the following:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) != 2 { 
         fmt.Printf("usage: %s number\n", filepath.Base(arguments[0])) 
         os.Exit(1) 
   } 
   numGR, _ := strconv.ParseInt(os.Args[1], 10, 64) 
   var waitGroup sync.WaitGroup 
   var i int64 
```

Once again, there is no problem in this particular code.

The third segment of rd.go has the following Go code:

```markup
   for i = 0; i < numGR; i++ { 
         waitGroup.Add(1) 
         go func() { 
               defer waitGroup.Done() 
               fmt.Printf("%d ", i) 
         }() 
   } 
```

This code is very suspicious because you try to print the value of a variable that keeps changing all the time because of the for loop.

The last part of rd.go is the following:

```markup
   waitGroup.Wait() 
   fmt.Println("\nExiting...") 
} 
```

There is nothing special in the last chunk of code.

Enabling the Go race detector for rd.go will generate the following output:

```markup
$ go run -race rd.go 10==================WARNING: DATA RACE
Read at 0x00c420074168 by goroutine 6:
  main.main.func1()
      /Users/mtsouk/Desktop/goCourse/ch/ch10/code/rd.go:25 +0x6c
    
Previous write at 0x00c420074168 by main goroutine:
  main.main()
      /Users/mtsouk/Desktop/goCourse/ch/ch10/code/rd.go:21 +0x30c
    
Goroutine 6 (running) created at:
  main.main()
      /Users/mtsouk/Desktop/goCourse/ch/ch10/code/rd.go:26 +0x2e2
==================
==================
WARNING: DATA RACE
Read at 0x00c420074168 by goroutine 7:
 main.main.func1()
     /Users/mtsouk/Desktop/goCourse/ch/ch10/code/rd.go:25 +0x6c
    
Previous write at 0x00c420074168 by main goroutine:
 main.main()
     /Users/mtsouk/Desktop/goCourse/ch/ch10/code/rd.go:21 +0x30c
    
Goroutine 7 (running) created at:
  main.main()
      /Users/mtsouk/Desktop/goCourse/ch/ch10/code/rd.go:26 +0x2e2
==================
2 3 4 4 5 6 7 8 9 10
Exiting...
Found 2 data race(s)
exit status 66 
```

So, the race detector found two data races. The first one happens when number 1 was not printed at all and the second when number 4 was printed two times. Additionally, number 0 was not printed despite being the initial value of i. Last, you should not get number 10 in the output but you did get it because the last value of i is indeed 10. Note that the main.main.func1() notation found in the preceding output means that Go talks about an anonymous function.

Put simply, what the previous two messages tell you is that there is something wrong with the i variable because it keeps changing while the goroutines of the program try to read it. Additionally, you cannot deterministically tell what will happen first.

Running the same program without the race detector will generate the following output:

```markup
$ go run rd.go 10
10 10 10 10 10 10 10 10 10 10
Exiting...
```

The problem with rd.go can be found in the anonymous function. As the anonymous function takes no arguments, it uses the current value of i, which cannot be determined with any certainty as it depends on the operating system and the Go scheduler: this is where the race situation happens! So, have in mind that one of the easiest places to have a race condition is inside a goroutine spawned from an anonymous function! As a result, if you have to solve such as situation, start by converting the anonymous function into regular functions with defined arguments!

Programs that use the race detector are slower and need more RAM than the same programs without the race detector. Last, if the race detector has nothing to report, it will generate no output.

Just Imagine

# About GOMAXPROCS

The GOMAXPROCS environment variable (and Go function) allows you to limit the number of operating system threads that can execute user-level Go code simultaneously.

Starting with Go version 1.5, the default value of GOMAXPROCS should be the number of cores available on your Unix system.

Although using a GOMAXPROCS value that is smaller than the number of the cores a Unix machine has might affect the performance of a program, specifying a GOMAXPROCS value that is bigger than the number of the available cores will not make your program run faster!

The code of goMaxProcs.go allows you to determine the value of the GOMAXPROCS - it will be presented in two parts.

The first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "runtime" 
) 
func getGOMAXPROCS() int {
   return runtime.GOMAXPROCS(0) 
} 
```

The second part is the following:

```markup
func main() { 
   fmt.Printf("GOMAXPROCS: %d\n", getGOMAXPROCS()) 
} 
```

Executing goMaxProcs.go on an Intel i7 machine with hyper threading support and the latest Go version gives the following output:

```markup
$ go run goMaxProcs.go 
GOMAXPROCS: 8 
```

However, if you execute goMaxProcs.go on a Debian Linux machine that runs an older Go version and has an older processor, it will generate the following output:

```markup
$ go version 
go version go1.3.3 linux/amd64 
$ go run goMaxProcs.go 
GOMAXPROCS: 1 
```

The way to change the value of GOMAXPROCS on the fly is as follows:

```markup
$ export GOMAXPROCS=80; go run goMaxProcs.go 
GOMAXPROCS: 80 
```

However, putting a value bigger than 256 will not work:

```markup
$ export GOMAXPROCS=800; go run goMaxProcs.go 
GOMAXPROCS: 256 
```

Last, have in mind that if you are running a concurrent program such as dWC.go using a single core, the concurrent version of the program might not be faster than the version of the program without goroutines! In some situations, this happens because the use of goroutines as well as the various calls to the sync.Add, sync.Wait, and sync.Done functions slows down the performance of a program. This can be verified by the following output:

```markup
$ export GOMAXPROCS=8; time go run dWC.go /tmp/*.data
    
real  0m10.826s
user  0m31.542s
sys   0m5.043s
$ export GOMAXPROCS=1; time go run dWC.go /tmp/*.data
    
real  0m15.362s
user  0m15.253s
sys   0m0.103s
$ time go run wc.go /tmp/*.data
    
real  0m15.158sexit
user  0m15.023s
sys   0m0.120s
```

Just Imagine

# Exercises

1.  Read carefully the documentation page of the sync package that can be found at [https://golang.org/pkg/sync/](https://golang.org/pkg/sync/).
2.  Try to implement dWC.go using a different shared memory technique than the one used in this chapter.
3.  Implement a struct data type that holds your account balance and make functions that read the amount of money you have and make changes to the money. Create an implementation that uses sync.RWMutex and another one that uses sync.Mutex.
4.  What would happen to mutexRW.go if you used Lock() and Unlock() everywhere instead of RLock() and RUnlock()?
5.  Try to implement traverse.go from _[](https://subscription.imaginedevops.io/book/programming/9781787125643/5)_[Chapter 5](https://subscription.imaginedevops.io/book/programming/9781787125643/5)_,_ _Files and Directories_ using goroutines.
6.  Try to create an implementation of improvedFind.go from [](https://subscription.imaginedevops.io/book/programming/9781787125643/5)[Chapter 5](https://subscription.imaginedevops.io/book/programming/9781787125643/5)_,_ _Files and Directories_ using goroutines.

Just Imagine

# Summary

This chapter talked about some advanced Go features related to goroutines, channels, and concurrent programming. However, the moral of this chapter is that channels can do many things and can be used in many situations, which means that the developer must be able to choose the appropriate technique to implement a task based on their experience.

The subject of the next chapter will be web development in Go and it will contain very interesting material, including sending and receiving JSON data, developing web servers and web clients, as well as talking to a MongoDB database from your Go code.