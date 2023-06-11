# File Input and Output

In the previous chapter, we talked about manipulating files and directories as entities without looking at their contents. However, in this chapter, we will take a different approach and look into the contents of files: you might consider this chapter one of the most important chapters in this book because **file input** and **file output** are primary tasks of any operating system.

The main purpose of this chapter is to teach how the Go standard library permits us to open files, read their contents, process them if we like, create new files, and put the desired data into them. There are two main ways to read and write files: using the io package and using the functions of the bufio package. However, both packages work in a comparative way.

This chapter will tell you about the following:

-   Opening files for writing and reading
-   Using the io package for file input and output
-   Using the io.Writer and io.Reader interfaces
-   Using the bufio package for buffered input and output
-   Copying files in Go
-   Implementing a version of the wc(1) utility in Go
-   Developing a version of the dd(1) command in Go
-   Creating sparse files
-   The importance of byte slices in file input and output: byte slices were first mentioned in [Chapter 2](https://subscription.imaginedevops.io/book/programming/9781787125643/2), _Writing Programs in Go_
-   Storing structured data in files and reading them afterwards
-   Converting tabs into space characters and vice versa

This chapter will not talk about appending data to an existing file: you will have to wait until [Chapter 7](https://subscription.imaginedevops.io/book/programming/9781787125643/7), _Working with System Files_, to learn more about putting data at the end of a file without destroying its existing data.

Just Imagine

# About file input and output

File input and output includes everything that has to do with reading the data of a file and writing the desired data to a file. There is not a single operating system that does not offer support for files and therefore for file input and output.

As this chapter is pretty big, I will stop talking and start showing you practical Go code that will make things clearer. So, the first thing that you will learn in this chapter is byte slices, which are very important in applications that are concerned with file input and output.

# Byte slices

**Byte slices** are a kind of slices used for file reading and writing. Putting it simply, they are slices of bytes used as a buffer during file reading and writing operations. This section will present a small Go example where a byte slice is used for writing to a file and reading from a file. As you will see byte slices all over this chapter, make sure that you understand the presented example. The related Go code is saved as byteSlice.go and will be presented in three parts.

The first part is as follows:

```markup
package main 
 
import ( 
   "fmt" 
   "io/ioutil" 
   "os" 
) 
```

The second part of byteSlice.go is as follows:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Println("Please provide a filename") 
         os.Exit(1) 
   } 
   filename := os.Args[1] 
 
   aByteSlice := []byte("Mihalis Tsoukalos!\n") 
   ioutil.WriteFile(filename, aByteSlice, 0644) 
```

Here, you use the aByteSlice byte slice to save some text into a file that is identified by the filename variable. The last part of byteSlice.go is the following Go code:

```markup
   f, err := os.Open(filename) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
   defer f.Close() 
 
   anotherByteSlice := make([]byte, 100) 
   n, err := f.Read(anotherByteSlice) 
   fmt.Printf("Read %d bytes: %s", n, anotherByteSlice) 
} 
```

Here, you define another byte slice named anotherByteSlice with 100 places that will be used for reading from the file you created previously. Note that %s used in fmt.Printf() forces anotherByteSlice to be printed as a string: using Println() would have produced a totally different output.

Note that as the file is smaller, the f.Read() call will put less data into anotherByteSlice.

The size of anotherByteSlice denotes the maximum amount of data that can be stored into it after a single call to Read() or after any other similar operation that reads data from a file.

Executing byteSlice.go will generate the following output:

```markup
$ go run byteSlice.go usingByteSlices
Read 19 bytes: Mihalis Tsoukalos!
```

Checking the size of the usingByteSlices file will verify that the right amount of data was written to it:

```markup
$ wc usingByteSlices
   1   2  19 usingByteSlices
```

# About binary files

There is no difference between reading and writing binary and plain text files in Go. So, when processing a file, Go makes no assumptions about its format. However, Go offers a package named binary that allows you to make translations between different encodings such as **little endian** and **big endian**.

The readBinary.go file briefly illustrates how to convert an integer number to a little endian number and to a big endian number, which might be useful when the files you want to process contain certain kinds of data; this mainly happens when we are dealing with raw devices and raw packet manipulation: remember everything is a file! The source code of readBinary.go will be presented in two parts.

The first part is as follows:

```markup
package main 
 
import ( 
   "bytes" 
   "encoding/binary" 
   "fmt" 
   "os" 
   "strconv" 
) 
 
func main() { 
   if len(os.Args) != 2 { 
         fmt.Println("Please provide an integer") 
         os.Exit(1) 
   } 
   aNumber, _ := strconv.ParseInt(os.Args[1], 10, 64) 
```

There is nothing special in this part of the program. The second part is the following:

```markup
   buf := new(bytes.Buffer) 
   err := binary.Write(buf, binary.LittleEndian, aNumber) 
   if err != nil { 
         fmt.Println("Little Endian:", err) 
   } 
 
   fmt.Printf("%d is %x in Little Endian\n", aNumber, buf) 
   buf.Reset() 
   err = binary.Write(buf, binary.BigEndian, aNumber) 
   if err != nil { 
         fmt.Println("Big Endian:", err) 
   } 
   fmt.Printf("And %x in Big Endian\n", buf) 
} 
```

The second part contains all the important Go code: the conversions happen with the help of the binary.Write() method and the proper write parameter (binary.LittleEndian or binary.BigEndian). The bytes.Buffer variable is used for the io.Reader and io.Writer interfaces of the program. Lastly, the buf.Reset() statement resets the buffer in order to be used afterwards for storing the big endian.

Executing readBinary.go will generate the following output:

```markup
$ go run readBinary.go 1
1 is 0100000000000000 in Little Endian
And 0000000000000001 in Big Endian
```

You can find more information about the binary package by visiting its documentation page at [https://golang.org/pkg/encoding/binary/](https://golang.org/pkg/encoding/binary/).

Just Imagine

# Useful I/O packages in Go

The io package is for performing primitive file I/O operations, whereas the bufio package is for executing buffered I/O.

In buffered I/O, the operating system uses an intermediate buffer during file read and write operations in order to reduce the number of filesystem calls. As a result, buffered input and output is faster and more efficient.

Additionally, you can use some of the functions of the fmt package to write text to a file. Note that the flag package will be also used in this chapter as well as in all the forthcoming ones where the developed utilities need to support command-line flags.

# The io package

The io package offers functions that allow you to write to or read from files. Its use will be illustrated in the usingIO.go file, which will be presented in three parts. What the program does is read 8 bytes from a file and write them in a standard output.

The first part is the preamble of the Go program:

```markup
package main 
 
import ( 
   "fmt" 
   "io" 
   "os" 
) 
```

The second part is the following Go code:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Println("Please provide a filename") 
         os.Exit(1) 
   } 
 
   filename := os.Args[1] 
   f, err := os.Open(filename) 
   if err != nil { 
         fmt.Printf("error opening %s: %s", filename, err) 
         os.Exit(1) 
   } 
   defer f.Close() 
```

The program also uses the handy defer command that defers the execution of a function until the surrounding function returns. As a result, defer is used very frequently in file I/O operations because it saves you from having to remember to execute the Close() call after you are done working with a file or when you leave a function in any number of locations using a return statement or os.Exit().

The last part of the program is the following:

```markup
   buf := make([]byte, 8) 
   if _, err := io.ReadFull(f, buf); err != nil { 
         if err == io.EOF { 
               err = io.ErrUnexpectedEOF 
         } 
   } 
   io.WriteString(os.Stdout, string(buf)) 
   fmt.Println() 
} 
```

The io.ReadFull() function here reads from the reader of an open file and puts the data into a byte slice that has 8 places. You can also see here the use of the io.WriteString() function for printing data to a standard output (os.Stdout) that is also a file. However, this is not a very common practice as you can simply use fmt.Println() instead.

Executing usingIO.go generates the following output:

```markup
$ go run usingIO.go usingByteSlices
Mihalis
```

# The bufio package

The functions of the bufio package allow you to perform buffered file operations, which means that although its operations look similar to the ones found in io, they work in a slightly different way.

What bufio actually does is to wrap an io.Reader or io.Writer object into a new value that implements the required interface while providing buffering to the new value. One of the handy features of the bufio package is that it allows you to read a text file line by line, word by word, and character by character without too much effort.

Once again, an example will try to clarify things: the name of the Go file that showcases the use of bufio is bufIO.go and will be presented in four parts.

The first part is the expected preamble:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "os" 
) 
```

The second part is the following:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Println("Please provide a filename") 
         os.Exit(1) 
   } 
 
   filename := os.Args[1] 
```

Here, you just try to get the name of the file that you are going to use.

The third part of bufIO.go has the following Go code:

```markup
   f, err := os.Open(filename) 
   if err != nil { 
         fmt.Printf("error opening %s: %s", filename, err) 
         os.Exit(1) 
   } 
   defer f.Close() 
 
   scanner := bufio.NewScanner(f) 
```

The default behavior of bufio.NewScanner is to read its input line by line, which means that each time you call the Scan() method that reads the next token, a new line will be returned. The last part is where you actually call the Scan() method in order to read the full contents of the file:

```markup
   for scanner.Scan() { 
         line := scanner.Text() 
 
         if scanner.Err() != nil { 
               fmt.Printf("error reading file %s", err) 
               os.Exit(1) 
         } 
         fmt.Println(line) 
   } 
}
```

The Text() method returns the latest token from the Scan() method as a string, which in this case will be a line. However, if you ever get strange results while trying to read a file line by line, it will most likely be the way your file ends a line, which is usually the case with text files coming from Windows machines.

Executing bufIO.go and feeding wc(1) with its output can help you verify that bufIO.go works as expected:

```markup
$ go run bufIO.go inputFile | wc
      11      12      62
$ wc inputFile
      11      12      62 inputFile
```

Just Imagine

# File I/O operations

Now that you know the basics of the io and bufio packages, it is time to learn more detailed information about their usage and how they can help you work with files. But first, we will talk about the fmt.Fprintf() function.

# Writing to files using fmt.Fprintf()

The use of the fmt.Fprintf() function allows you to write formatted text to files in a way that is similar to the way the fmt.Printf() function works. Note that fmt.Fprintf() can write to any io.Writer interface and that our files will satisfy the io.Writer interface.

The Go code that illustrates the use of fmt.Fprintf() can be found in fmtF.go, which will be presented in three parts. The first part is the expected preamble:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
) 
```

The second part has the following Go code:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Println("Please provide a filename") 
         os.Exit(1) 
   } 
 
   filename := os.Args[1] 
   destination, err := os.Create(filename) 
   if err != nil { 
         fmt.Println("os.Create:", err) 
         os.Exit(1) 
   } 
   defer destination.Close() 
```

Note that the os.Create() function will truncate the file if it already exists.

The last part is the following:

```markup
   fmt.Fprintf(destination, "[%s]: ", filename) 
   fmt.Fprintf(destination, "Using fmt.Fprintf in %s\n", filename) 
} 
```

Here, you write the desired text data to the file that is identified by the destination variable using fmt.Fprintf() as if you were using the fmt.Printf() method.

Executing fmtF.go will generate the following output:

```markup
$ go run fmtF.go test
$ cat test
[test]: Using fmt.Fprintf in test 
```

In other words, you can create plain text files using fmt.Fprintf().

# About io.Writer and io.Reader

Both io.Writer and io.Reader are interfaces that embed the io.Write() and io.Read() methods, respectively. The use of io.Writer and io.Reader will be illustrated in readerWriter.go, which will be presented in four parts. The program computes the characters of its input file and writes the number of characters to another file: if you are dealing with Unicode characters that take more than one byte per character, you might consider that the program is reading bytes. The output filename has the name of the original file plus the .Count extension.

The first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "io" 
   "os" 
) 
```

The second part is the following:

```markup
func countChars(r io.Reader) int { 
   buf := make([]byte, 16) 
   total := 0 
   for { 
         n, err := r.Read(buf) 
         if err != nil && err != io.EOF { 
               return 0 
         } 
         if err == io.EOF { 
               break 
         } 
         total = total + n 
   } 
   return total 
} 
```

Once again, a byte slice is used during reading. The break statement allows you to exit the for loop. The third part is the following code:

```markup
func writeNumberOfChars(w io.Writer, x int) { 
   fmt.Fprintf(w, "%d\n", x) 
} 
```

Here you can see how you can write a number to a file using fmt.Fprintf(): I did not manage to do the same using a byte slice! Additionally, note that the presented code writes text to a file using an io.Writer variable (w).

The last part of readerWriter.go has the following Go code:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Println("Please provide a filename") 
         os.Exit(1) 
   } 
 
   filename := os.Args[1] 
   _, err := os.Stat(filename) 
   if err != nil { 
         fmt.Printf("Error on file %s: %s\n", filename, err) 
         os.Exit(1) 
   } 
 
   f, err := os.Open(filename) 
   if err != nil { 
         fmt.Println("Cannot open file:", err) 
         os.Exit(-1) 
   } 
   defer f.Close() 
 
   chars := countChars(f) 
   filename = filename + ".Count" 
   f, err = os.Create(filename) 
   if err != nil { 
         fmt.Println("os.Create:", err) 
         os.Exit(1) 
   } 
   defer f.Close() 
   writeNumberOfChars(f, chars) 
} 
```

The execution of readerWriter.go generates no output; so, it is up to you to check its correctness, which in this case happens with the help of wc(1):

```markup
$ go run readerWriter.go /tmp/swtag.log
$ wc /tmp/swtag.log
     119     635    7780 /tmp/swtag.log
$ cat /tmp/swtag.log.Count
7780
```

# Finding out the third column of a line

Now that you know how to read a file, it is time to present a modified version of the readColumn.go program you saw in [](https://subscription.imaginedevops.io/book/programming/9781787125643/3)[](https://subscription.imaginedevops.io/book/programming/9781787125643/3)[](https://subscription.imaginedevops.io/book/programming/9781787125643/3)[Chapter 3](https://subscription.imaginedevops.io/book/programming/9781787125643/3), _Advanced Go Features_. The new version is also named readColumn.go, but has two major improvements. The first is that you can provide the desired column as a command-line argument and the second is that it can read multiple files if it gets multiple command-line arguments.

The readColumn.go file will be presented in three parts. The first part of readColumn.go is the following:

```markup
package main 
 
import ( 
   "bufio" 
   "flag" 
   "fmt" 
   "io" 
   "os" 
   "strings" 
) 
```

The next part of readColumn.go contains the following Go code:

```markup
func main() { 
   minusCOL := flag.Int("COL", 1, "Column") 
   flag.Parse() 
   flags := flag.Args() 
 
   if len(flags) == 0 { 
         fmt.Printf("usage: readColumn <file1> [<file2> [... <fileN]]\n") 
         os.Exit(1) 
   } 
 
   column := *minusCOL 
 
   if column < 0 { 
         fmt.Println("Invalid Column number!") 
         os.Exit(1) 
   } 
```

As you will understand from the definition of the minusCOL variable, if the user does not use this flag, the program will print the contents of the first column of each file it reads.

The last part of readColumn.go is as follows:

```markup
   for _, filename := range flags { 
         fmt.Println("\t\t", filename) 
         f, err := os.Open(filename) 
         if err != nil { 
               fmt.Printf("error opening file %s", err) 
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
               } 
 
               data := strings.Fields(line) 
               if len(data) >= column { 
                     fmt.Println((data[column-1])) 
               } 
         } 
   } 
} 
```

The preceding code does not do anything that you have not seen before. The for loop is used for processing all command-line arguments. However, if a file fails to open for some reason, the program will not stop its execution, but it will continue processing the rest of the files if they exist. However, the program expects that its input files end in a newline and you might see strange results if an input file ends differently.

Executing readColumn.go generates the following output, which is abbreviated in order to save some book space:

```markup
$ go run readColumn.go -COL=3 pF.data isThereAFile up.data
            pF.data
            isThereAFile
error opening file open isThereAFile: no such file or directory
            up.data
0.05
0.05
0.05
0.05
0.05
0.05
```

In this case, there is no file named isThereAFile and the pF.data file does not have a third column. However, the program did its best and printed what it could!

Just Imagine

# Copying files in Go

Every operating system allows you to copy files because this is a very important and necessary operation. This section will show you how to copy files in Go now that you know how to read files!

# There is more than one way to copy a file!

Most programming languages offer more than one way to create a copy of a file and Go is no exception. It is up to the developer to decide which approach to implement.

The t_here is more than one way to do it_ rule applies to almost everything implemented in this book, but file copying is the most characteristic example of this rule because you can copy a file by reading it line by line, byte by byte, or all at once! However, this rule does not apply to the way Go likes to format its code!

# Copying text files

There is no point in treating the copying of text files in a special way unless you want to inspect or modify their contents. As a result, the three techniques presented here will not differentiate between plain text and binary file copying.

[Chapter 7](https://subscription.imaginedevops.io/book/programming/9781787125643/7)_, Working with System Files_, will talk about file permissions because there are times that you want to create new files with the file permissions you choose.

# Using io.Copy

This subsection will present a technique for copying files that uses the io.Copy() function. What is special about the io.Copy() function is the fact that is does not give you any flexibility in the process. The name of the program will be notGoodCP.go and will be presented in three parts. Note that a more appropriate filename for notGoodCP.go would have been copyEntireFileAtOnce.go or copyByReadingInputFileAllAtOnce.go!

The first part of the Go code of notGoodCP.go is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "io" 
   "os" 
) 
```

The second part is as follows:

```markup
func Copy(src, dst string) (int64, error) { 
   sourceFileStat, err := os.Stat(src) 
   if err != nil { 
         return 0, err 
   } 
 
   if !sourceFileStat.Mode().IsRegular() { 
         return 0, fmt.Errorf("%s is not a regular file", src) 
   } 
 
   source, err := os.Open(src) 
   if err != nil { 
         return 0, err 
   } 
   defer source.Close() 
 
   destination, err := os.Create(dst) 
   if err != nil { 
         return 0, err 
   } 
   defer destination.Close() 
   nBytes, err := io.Copy(destination, source) 
   return nBytes, err 
 
}
```

Here we define our own function that uses io.Copy() to make a copy of a file. The Copy() function checks whether the source file is a regular file before trying to copy it, which makes perfect sense.

The last part is the implementation of the main() function:

```markup
func main() { 
   if len(os.Args) != 3 { 
         fmt.Println("Please provide two command line arguments!") 
         os.Exit(1) 
   } 
 
   sourceFile := os.Args[1] 
   destinationFile := os.Args[2] 
   nBytes, err := Copy(sourceFile, destinationFile) 
 
   if err != nil { 
         fmt.Printf("The copy operation failed %q\n", err) 
   } else { 
         fmt.Printf("Copied %d bytes!\n", nBytes) 
   } 
} 
```

The best tool for testing whether a file is an exact copy of another file is the diff(1) utility, which also works with binary files. You can learn more about diff(1) by reading its main page.

Executing notGoodCP.go will generate the following results:

```markup
$ go run notGoodCP.go testFile aCopy
Copied 871 bytes!
$ diff aCopy testFile
$ wc testFile aCopy
      51     127     871 testFile
      51     127     871 aCopy
     102     254    1742 total
```

# Reading a file all at once!

The technique in this section will use the ioutil.WriteFile() and ioutil.ReadFile() functions. Note that ioutil.ReadFile() does not implement the io.Reader interface and therefore is a little restrictive.

The Go code for this section is named readAll.go and will be presented in three parts.

The first part has the following Go code:

```markup
package main 
 
import ( 
   "fmt" 
   "io/ioutil" 
   "os" 
) 
```

The second part is the following:

```markup
func main() { 
   if len(os.Args) != 3 { 
         fmt.Println("Please provide two command line arguments!") 
         os.Exit(1) 
   } 
 
   sourceFile := os.Args[1] 
   destinationFile := os.Args[2] 
```

The last part is as follows:

```markup
   input, err := ioutil.ReadFile(sourceFile) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
 
   err = ioutil.WriteFile(destinationFile, input, 0644) 
   if err != nil { 
         fmt.Println("Error creating the new file", destinationFile) 
         fmt.Println(err) 
         os.Exit(1) 
   } 
} 
```

Note that the ioutil.ReadFile() function reads the entire file, which might not be efficient when you want to copy huge files. Similarly, the ioutil.WriteFile() function writes all the given data to a file that is identified by its first argument.

The execution of readAll.go generates the following output:

```markup
$ go run readAll.go testFile aCopy
$ diff aCopy testFile
$ ls -l testFile aCopy
-rw-r--r--  1 mtsouk  staff  871 May  3 21:07 aCopy
-rw-r--r--@ 1 mtsouk  staff  871 May  3 21:04 testFile
$ go run readAll.go doesNotExist aCopy
open doesNotExist: no such file or directory
exit status 1
```

# An even better file copy program

This section will present a program that uses a more traditional approach, where a buffer is used for reading and copying to the new file.

Although traditional Unix command-line utilities are silent when there are no errors, it is not bad to print some kind of information, such as the number of bytes read, in your own tools. However, the right thing to do is to follow the Unix way.

There exist two main reasons that make cp.go better than notGoodCP.go. The first is that the developer has more control over the process in exchange for having to write more Go code and the second is that cp.go allows you to define the size of the buffer, which is the most important parameter in the copy operation.

The code of cp.go will be presented in five parts. The first part is the expected preamble along with a global variable that holds the size of the read buffer:

```markup
package main 
 
import ( 
   "fmt" 
   "io" 
   "os" 
   "path/filepath" 
   "strconv" 
) 
 
var BUFFERSIZE int64 
```

The second part is the following:

```markup
func Copy(src, dst string, BUFFERSIZE int64) error { 
   sourceFileStat, err := os.Stat(src) 
   if err != nil { 
         return err 
   } 
 
   if !sourceFileStat.Mode().IsRegular() { 
         return fmt.Errorf("%s is not a regular file.", src) 
   } 
 
   source, err := os.Open(src) 
   if err != nil { 
         return err 
   } 
   defer source.Close() 
```

As you can see here, the size of the buffer is given to the Copy() function as an argument. The other two command-line arguments are the input filename and the output filename.

The third part has the remaining Go code of the Copy() function:

```markup
   _, err = os.Stat(dst) 
   if err == nil { 
         return fmt.Errorf("File %s already exists.", dst) 
   } 
 
   destination, err := os.Create(dst) 
   if err != nil { 
         return err 
   } 
   defer destination.Close() 
 
   if err != nil { 
         panic(err) 
   } 
 
   buf := make([]byte, BUFFERSIZE) 
   for { 
         n, err := source.Read(buf) 
         if err != nil && err != io.EOF { 
               return err 
         } 
         if n == 0 { 
               break 
         } 
 
         if _, err := destination.Write(buf[:n]); err != nil { 
               return err 
         } 
   } 
   return err 
} 
```

There is nothing special here: you just keep calling source, Read() until you reach the end of the input file. Each time you read something, you call destination. Write() to save it to the output file. The buf\[:n\] notation allows you to read the first n characters from the buf slice.

The fourth part contains the following Go code:

```markup
func main() { 
   if len(os.Args) != 4 { 
         fmt.Printf("usage: %s source destination BUFFERSIZE\n", 
filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   source := os.Args[1] 
   destination := os.Args[2] 
   BUFFERSIZE, _ = strconv.ParseInt(os.Args[3], 10, 64) 
```

The filepath.Base() is used for getting the name of the executable file.

The last part is the following:

```markup
   fmt.Printf("Copying %s to %s\n", source, destination) 
   err := Copy(source, destination, BUFFERSIZE) 
   if err != nil { 
         fmt.Printf("File copying failed: %q\n", err) 
   } 
}
```

Executing cp.go will generate the following output:

```markup
$ go run cp.go inputFile aCopy 2048
Copying inputFile to aCopy
$ diff inputFile aCopy
```

If there is a problem with the copy operation, you will get a descriptive error message.

So, if the program cannot find the input file, it will print the following:

```markup
$ go run cp.go A /tmp/myCP 1024
Copying A to /tmp/myCP
File copying failed: "stat A: no such file or directory"
```

If the program cannot read the input file, you will get the following message:

```markup
$ go run cp.go inputFile /tmp/myCP 1024
Copying inputFile to /tmp/myCP
File copying failed: "open inputFile: permission denied"
```

If the program cannot create the output file, it will print the following error message:

```markup
$ go run cp.go inputFile /usr/myCP 1024
Copying inputFile to /usr/myCP
File copying failed: "open /usr/myCP: operation not permitted"
```

If the destination file already exists, you will get the following output:

```markup
$ go run cp.go inputFile outputFile 1024
Copying inputFile to outputFile
File copying failed: "File outputFile already exists."
```

# Benchmarking file copying operations

The size of the buffer you use in file operations is really important and affects the performance of your system tools, especially when you are dealing with very big files.

Although developing reliable software should be your main concern, you should not forget to make your systems software fast and efficient!

So, this section will try to see how the size of the buffer affects the file copying operations by executing cp.go with various buffer sizes and comparing its performance with readAll.go, notGoodCP.go as well as cp(1).

In the old Unix days when the amount of RAM on Unix machines was too small, using a large buffer was not recommended. However, nowadays, using a buffer with a size of 100 MB is not considered bad practice, especially when you know in advance that you are going to copy lots of big files such as the data files of a database server.

We will use three files with different sizes in our testing: these three files will be generated using the dd(1) utility, as shown here:

```markup
$dd if=/dev/urandom of=100MB count=100000 bs=1024
100000+0 records in
100000+0 records out
102400000 bytes transferred in 6.800277 secs (15058210 bytes/sec)
$ dd if=/dev/urandom of=1GB count=1000000 bs=1024
1000000+0 records in
1000000+0 records out
1024000000 bytes transferred in 68.887482 secs (14864820 bytes/sec)
$ dd if=/dev/urandom of=5GB count=5000000 bs=1024
5000000+0 records in
5000000+0 records out
5120000000 bytes transferred in 339.357738 secs (15087324 bytes/sec)
$ ls -l 100MB 1GB 5GB
-rw-r--r--  1 mtsouk  staff   102400000 May  4 10:30 100MB
-rw-r--r--  1 mtsouk  staff  1024000000 May  4 10:32 1GB
-rw-r--r--  1 mtsouk  staff  5120000000 May  4 10:38 5GB
```

The first file is 100 MB, the second is 1 GB, and the third is 5 GB in size.

Now, it is time for the actual testing using the time(1) utility. First, we will test the performance of notGoodCP.go and readAll.go:

```markup
$ time ./notGoodCP 100MB copy
Copied 102400000 bytes!
    
real  0m0.153s
user  0m0.003s
sys   0m0.084s
$ time ./notGoodCP 1GB copy
Copied 1024000000 bytes!
    
real  0m1.461s
user  0m0.029s
sys   0m0.833s
$ time ./notGoodCP 5GB copy
Copied 5120000000 bytes!
    
real  0m12.193s
user  0m0.161s
sys   0m5.251s
$ time ./readAll 100MB copy
    
real  0m0.249s
user  0m0.003s
sys   0m0.138s
$ time ./readAll 1GB copy
    
real  0m3.117s
user  0m0.639s
sys   0m1.644s
$ time ./readAll 5GB copy
    
real  0m28.918s
user  0m8.106s
sys   0m21.364s
```

Now, you will see the results from the cp.go program using four different buffer sizes, 16, 1024, 1048576, and 1073741824. First, let's copy the 100 MB file:

```markup
$ time ./cp 100MB copy 16
Copying 100MB to copy
    
real  0m13.240s
user  0m2.699s
sys   0m10.530s
$ time ./cp 100MB copy 1024
Copying 100MB to copy
    
real  0m0.386s
user  0m0.053s
sys   0m0.303s
$ time ./cp 100MB copy 1048576
Copying 100MB to copy
    
real  0m0.135s
user  0m0.001s
sys   0m0.075s
$ time ./cp 100MB copy 1073741824
Copying 100MB to copy
    
real  0m0.390s
user  0m0.011s
sys   0m0.136s
```

Then, we will copy the 1 GB file:

```markup
$ time ./cp 1GB copy 16
Copying 1GB to copy
    
real  2m10.054s
user  0m26.497s
sys   1m43.411s
$ time ./cp 1GB copy 1024
Copying 1GB to copy
    
real  0m3.520s
user  0m0.533s
sys   0m2.944s
$ time ./cp 1GB copy 1048576
Copying 1GB to copy
    
real  0m1.431s
user  0m0.006s
sys   0m0.749s
$ time ./cp 1GB copy 1073741824
Copying 1GB to copy
    
real  0m2.033s
user  0m0.012s
sys   0m1.310s
```

Next, we will copy the 5 GB file:

```markup
$ time ./cp 5GB copy 16Copying 5GB to copy
    
real  10m41.551s
user  2m11.695s
sys   8m29.248s
$ time ./cp 5GB copy 1024
Copying 5GB to copy
    
real  0m16.558s
user  0m2.415s
sys   0m13.597s
$ time ./cp 5GB copy 1048576
Copying 5GB to copy
    
real  0m7.172s
user  0m0.028s
sys   0m3.734s
$ time ./cp 5GB copy 1073741824
Copying 5GB to copy
    
real  0m8.612s
user  0m0.011s
sys   0m4.536s
```

Finally, let's present the results from the cp(1) utility that comes with macOS Sierra:

```markup
$ time cp 100MB copy
    
real  0m0.274s
user  0m0.002s
sys   0m0.105s
$ time cp 1GB copy
    
real  0m2.735s
user  0m0.003s
sys   0m1.014s
$ time cp 5GB copy
    
real  0m12.199s
user  0m0.012s
sys   0m5.050s
```

The following figure shows a graph with the values of the real fields from the output of the time(1) utility for all the aforementioned results:

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/a8e66124-879d-4896-b67f-28ae698552a5.png)

Benchmarking results for the various copy utilities

As you can see from the results, the cp(1) utility does a pretty good job. However, cp.go is more versatile because it allows you to define the size of the buffer. On the other hand, if you use cp.go with a small buffer size (16 bytes), then the entire process will be totally ruined! Additionally, it is interesting that readAll.go does a pretty decent job with relatively small files and it is slow only when copying the 5 GB file, which is not bad for such a small program: you can consider readAll.go as a quick and dirty solution!

Just Imagine

# Developing wc(1) in Go

The principal idea behind the code of the wc.go program is that you can read a text file line by line until there is nothing left to read. For each line you read, you find out the number of characters and the number of words it has. As you need to read your input line by line, the use of bufio is preferred instead of the plain io because it simplifies the code. However, trying to implement wc.go on your own using io would be a very educational exercise.

But first, you will see that the wc(1) utility generates the following output:

```markup
$ wc wc.go cp.go
      68     160    1231 wc.go
      45     112     755 cp.go
     113     272    1986 total
```

So, if wc(1) has to process more than one file, it automatically generates summary information.

In [Chapter 9](https://subscription.imaginedevops.io/book/programming/9781787125643/9), _Goroutines - Basic Features_, you will learn how to create a version of wc.go using Go routines. However, the core functionality of both versions will be exactly the same!

# Counting words

The trickiest part of the code implementation is word counting, which is implemented using regular expressions:

```markup
r := regexp.MustCompile("[^\\s]+") 
for range r.FindAllString(line, -1) { 
numberOfWords++ 
} 
```

Here, the provided regular expression separates the words of a line based on whitespace characters in order to count them afterwards!

# The wc.go code!

After this little introduction, it is time to see the Go code of wc.go, which will be presented in five parts. The first part is the expected preamble:

```markup
package main 
 
import ( 
   "bufio" 
   "flag" 
   "fmt" 
   "io" 
   "os" 
   "regexp" 
) 
```

The second part is the implementation of the countLines() function, which includes the core functionality of the program. Note that the name countLines() may have been a poor choice as countLines() also counts the words and the characters of a file:

```markup
func countLines(filename string) (int, int, int) { 
   var err error 
   var numberOfLines int 
   var numberOfCharacters int 
   var numberOfWords int 
   numberOfLines = 0 
   numberOfCharacters = 0 
   numberOfWords = 0 
 
   f, err := os.Open(filename) 
   if err != nil { 
         fmt.Printf("error opening file %s", err) 
         os.Exit(1) 
   } 
   defer f.Close() 
 
   r := bufio.NewReader(f) 
   for { 
         line, err := r.ReadString('\n') 
 
         if err == io.EOF { 
               break 
         } else if err != nil { 
               fmt.Printf("error reading file %s", err)                break 
         } 
 
         numberOfLines++ 
         r := regexp.MustCompile("[^\\s]+") 
         for range r.FindAllString(line, -1) { 
               numberOfWords++ 
         } 
         numberOfCharacters += len(line) 
   } 
 
   return numberOfLines, numberOfWords, numberOfCharacters 
} 
```

Lots of interesting things exist here. First of all, you can see the Go code presented in the previous section for counting the words of each line. Counting lines is easy because each time the bufio reader reads a new line, the value of the numberOfLines variable is increased by one. The ReadString() function tells the program to read until the first occurrence of '\\n' in the input: multiple calls to ReadString() mean that you are reading a file line by line.

Next, you can see that the countLines() function returns three integer values. Lastly, counting characters is implemented with the help of the len() function that returns the number of characters in a given string, which in this case is the line that was read. The for loop terminates when you get the io.EOF error message, which signifies that there is nothing left to read from the input file.

The third part of wc.go starts with the beginning of the implementation of the main() function, which also includes the configuration of the flag package:

```markup
func main() { 
   minusC := flag.Bool("c", false, "Characters") 
   minusW := flag.Bool("w", false, "Words") 
   minusL := flag.Bool("l", false, "Lines") 
 
   flag.Parse() 
   flags := flag.Args() 
 
   if len(flags) == 0 { 
         fmt.Printf("usage: wc <file1> [<file2> [... <fileN]]\n") 
         os.Exit(1) 
   } 
 
   totalLines := 0 
   totalWords := 0 
   totalCharacters := 0 
   printAll := false 
 
   for _, filename := range flag.Args() { 
```

The last for statement is for processing all the input files given to the program. The wc.go program supports three flags: the \-c flag is for printing the character count, the \-w flag is for printing the word count, and the \-l flag is for printing the line count.

The fourth part is the following:

```markup
         numberOfLines, numberOfWords, numberOfCharacters := countLines(filename) 
 
         totalLines = totalLines + numberOfLines 
         totalWords = totalWords + numberOfWords 
         totalCharacters = totalCharacters + numberOfCharacters 
 
         if (*minusC && *minusW && *minusL) || (!*minusC && !*minusW && !*minusL) { 
               fmt.Printf("%d", numberOfLines) 
               fmt.Printf("\t%d", numberOfWords) 
               fmt.Printf("\t%d", numberOfCharacters) 
               fmt.Printf("\t%s\n", filename) 
               printAll = true 
               continue 
         } 
 
         if *minusL { 
               fmt.Printf("%d", numberOfLines) 
         } 
 
         if *minusW { 
               fmt.Printf("\t%d", numberOfWords) 
         } 
 
         if *minusC { 
               fmt.Printf("\t%d", numberOfCharacters) 
         } 
 
         fmt.Printf("\t%s\n", filename) 
   } 
```

This part deals with printing the information on a per file basis depending on the command-line flags. As you can see, most of the Go code here is for handling the output according to the command-line flags.

The last part is the following:

```markup
   if (len(flags) != 1) && printAll { 
         fmt.Printf("%d", totalLines) 
         fmt.Printf("\t%d", totalWords) 
         fmt.Printf("\t%d", totalCharacters) 
         fmt.Println("\ttotal") 
return 
   } 
 
   if (len(flags) != 1) && *minusL { 
         fmt.Printf("%d", totalLines) 
   } 
 
   if (len(flags) != 1) && *minusW { 
         fmt.Printf("\t%d", totalWords) 
   } 
 
   if (len(flags) != 1) && *minusC { 
         fmt.Printf("\t%d", totalCharacters) 
   } 
 
   if len(flags) != 1 { 
         fmt.Printf("\ttotal\n") 
   } 
} 
```

This is where you print the total number of lines, words, and characters read according to the flags of the program. Once again, most of the Go code here is for modifying the output according to the command-line flags.

Executing wc.go will generate the following output:

```markup
$ go build wc.go
$ ls -l wc
-rwxr-xr-x  1 mtsouk  staff  2264384 Apr 29 21:10 wc
$ ./wc wc.go sparse.go notGoodCP.go
120   280   2319  wc.go
44    98    697   sparse.go
27    61    418   notGoodCP.go
191   439   3434  total
$ ./wc -l wc.go sparse.go
120   wc.go
44    sparse.go
164   total
$ ./wc -w -l wc.go sparse.go
120   280   wc.go
44    98    sparse.go
164   378   total
```

There is a subtle point here: using Go source files as command-line arguments to the go run wc.go command will fail. This will happen because the compiler will try to compile the Go source files instead of treating them as command-line arguments to the go run wc.go command. The following output proves this:

```markup
$ go run wc.go sparse.go
# command-line-arguments
./sparse.go:11: main redeclared in this block
      previous declaration at ./wc.go:49
$ go run wc.go wc.go
package main: case-insensitive file name collision:
"wc.go" and "wc.go"
$ go run wc.go cp.go sparse.go
# command-line-arguments
./cp.go:35: main redeclared in this block
      previous declaration at ./wc.go:49
./sparse.go:11: main redeclared in this block
      previous declaration at ./cp.go:35
```

Additionally, trying to execute wc.go on a Linux system with Go version 1.3.3 will fail with the following error message:

```markup
$ go version
go version go1.3.3 linux/amd64
$ go run wc.go
# command-line-arguments
./wc.go:40: syntax error: unexpected range, expecting {
./wc.go:46: non-declaration statement outside function body
./wc.go:47: syntax error: unexpected }
```

# Comparing the performance of wc.go and wc(1)

In this subsection, we will compare the performance of our version of wc(1) with the wc(1) version that comes with macOS Sierra 10.12.6. First, we will execute wc.go:

```markup
$ file wc
wc: Mach-O 64-bit executable x86_64
$ time ./wc *.data
672320      3361604     9413057     connections.data
269123      807369      4157790     diskSpace.data
672040      1344080     8376070     memory.data
1344533     2689066     5378132     pageFaults.data
269465      792715      4068250     uptime.data
3227481     8994834     31393299    total
    
real  0m17.467s
user  0m22.164s
sys   0m3.885s
```

Then, we will execute the macOS version of wc(1) to process the same files:

```markup
$ file `which wc`
/usr/bin/wc: Mach-O 64-bit executable x86_64
$ time wc *.data
672320 3361604 9413057 connections.data
269123  807369 4157790 diskSpace.data
672040 1344080 8376070 memory.data
1344533 2689066 5378132 pageFaults.data
269465  792715 4068250 uptime.data
3227481 8994834 31393299 total
    
real  0m0.086s
user  0m0.076s
sys   0m0.007s
```

Let's look at the good news here first; the two utilities generated exactly the same output, which means that our Go version of wc(1) works great and can process big text files!

Now, the bad news; wc.go is slow! It took wc(1) less than a second to process all five files, whereas it took wc.go nearly 18 seconds to perform the same task!

The general idea when developing software of any kind, on any platform, using any programming language, is that you should try to have a working version of it, which does not contain any bugs before trying to optimize it and not the other way round!

# Reading a text file character by character

Although reading a text file character by character is not needed for the development of the wc(1) utility, it would be good to know how to implement it in Go. The name of the file will be charByChar.go and will be presented in four parts.

The first part is the following Go code:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "io/ioutil" 
   "os" 
   "strings" 
) 
```

Although charByChar.go does not have many lines of Go code, it needs lots of Go standard packages, which is a naive indication that the task it implements is not trivial. The second part is as follows:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Not enough arguments!") 
         os.Exit(1) 
   } 
   input := arguments[1] 
```

The third part is the following:

```markup
   buf, err := ioutil.ReadFile(input) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
```

The last part has the following Go code:

```markup
   in := string(buf) 
   s := bufio.NewScanner(strings.NewReader(in)) 
   s.Split(bufio.ScanRunes) 
 
   for s.Scan() { 
         fmt.Print(s.Text()) 
   } 
} 
```

Here, ScanRunes is a split function that returns each character (rune) as a token. Then, the call to Scan() allows us to process each character one by one. There also exist ScanWords and ScanLines for getting words and lines, respectively. If you use fmt.Println(s.Text()) as the last statement in the program instead of fmt.Print(s.Text()), then each character will be printed on its own line and the task of the program will be more obvious.

Executing charByChar.go generates the following output:

```markup
$ go run charByChar.go testpackage main...
```

The wc(1) command can verify the correctness of the Go code of charByChar.go by comparing the input file with the output generated by charByChar.go:

```markup
$ go run charByChar.go test | wc
      32      54     439
$ wc test
      32      54     439 test
```

# Doing some file editing!

This section will present a Go program that converts tab characters to space characters in files and vice versa! This is the job that is usually done by a text editor, but it is good to know how to perform it on your own.

The code will be saved in tabSpace.go and will be presented in four parts.

Note that tabSpace.go reads text files line by line, but you can also develop a version that reads text file character by character.

In the current implementation, all the work is done with the help of regular expressions, pattern matching, and search and replace operations.

The first part is the expected preamble:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "io" 
   "os" 
   "path/filepath" 
   "strings" 
) 
```

The second part contains the following Go code:

```markup
func main() { 
   if len(os.Args) != 3 { 
         fmt.Printf("Usage: %s [-t|-s] filename!\n", filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
   convertTabs := false 
   convertSpaces := false 
   newLine := "" 
 
   option := os.Args[1] 
   filename := os.Args[2] 
   if option == "-t" { 
         convertTabs = true 
   } else if option == "-s" { 
         convertSpaces = true 
   } else { 
         fmt.Println("Unknown option!") 
         os.Exit(1) 
   } 
```

The third part contains the following Go code:

```markup
   f, err := os.Open(filename) 
   if err != nil { 
         fmt.Printf("error opening %s: %s", filename, err) 
         os.Exit(1) 
   } 
   defer f.Close() 
 
   r := bufio.NewReader(f) 
   for { 
         line, err := r.ReadString('\n') 
 
         if err == io.EOF { 
               break 
         } else if err != nil { 
               fmt.Printf("error reading file %s", err) 
               os.Exit(1) 
         } 
```

The last part is the following:

```markup
         if convertTabs == true { 
               newLine = strings.Replace(line, "\t", "    ", -1) 
         } else if convertSpaces == true { 
               newLine = strings.Replace(line, "    ", "\t", -1) 
         } 
 
         fmt.Print(newLine) 
   } 
} 
```

This part is where the magic happens using the appropriate strings.Replace() call. In its current implementation, each tab is replaced by four space characters and vice versa, but you can change that by modifying the Go code.

Once again, a big part of tabSpace.go relates to error handling because many strange things can happen when you try to open a file for reading!

According to the Unix philosophy, the output of tabSpace.go will be printed on the screen and will not be saved in a new text file. Using tabSpace.go with wc(1) can prove its correctness:

```markup
$ go run tabSpace.go -t cp.go > convert
$ wc convert cp.go
      76     192    1517 convert 
      76     192    1286 cp.go
     152     384    2803 total
$ go run tabSpace.go -s convert | wc
      76     192    1286
```

Just Imagine

# Interprocess communication

**Interprocess communication** (**IPC**), putting it simply, is allowing Unix processes to talk to each other. Various techniques exist that allow processes and programs to talk to each other. The single most popular technique used in Unix systems is the pipe, which exists since the early Unix days. [Chapter _8_](https://subscription.imaginedevops.io/book/programming/9781787125643/8), _Processes and Signals_, will talk more about implementing Unix pipes in Go. Another form of IPC is Unix domain sockets, which will also be discussed in [Chapter _8_](https://subscription.imaginedevops.io/book/programming/9781787125643/8), _Processes and Signals_.

[Chapter _12_](https://subscription.imaginedevops.io/book/programming/9781787125643/12), _Network Programming_, will talk about another form of Interprocess communication, which is network sockets. Shared memory also exists, but Go is against the use of shared memory as a means of communication. [Chapter _9_](https://subscription.imaginedevops.io/book/programming/9781787125643/9), _Goroutines - Basic Features_, and [Chapter 10](https://subscription.imaginedevops.io/book/programming/9781787125643/10), _Goroutines - Advanced Features_, will show various techniques that allow goroutines to communicate with others and share and exchange data.

Just Imagine

# Sparse files in Go

Large files that are created with the os.Seek() function may have holes in them and occupy fewer disk blocks than files with the same size, but without holes in them; such files are called sparse files. This section will develop a program that creates sparse files.

The Go code of sparse.go will be presented in three parts. The first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "log" 
   "os" 
   "path/filepath" 
   "strconv" 
) 
```

The second part of sparse.go has the following Go code:

```markup
func main() { 
   if len(os.Args) != 3 { 
         fmt.Printf("usage: %s SIZE filename\n", filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   SIZE, _ := strconv.ParseInt(os.Args[1], 10, 64) 
   filename := os.Args[2] 
 
   _, err := os.Stat(filename) 
   if err == nil { 
         fmt.Printf("File %s already exists.\n", filename) 
         os.Exit(1) 
   } 
```

The strconv.ParseInt() function is used for converting the command-line argument that defines the size of the sparse file from its string value to its integer value. Additionally, the os.Stat() call makes sure that you will not accidentally overwrite an existing file.

The last part is where the action takes place:

```markup
   fd, err := os.Create(filename) 
   if err != nil { 
         log.Fatal("Failed to create output") 
   } 
 
   _, err = fd.Seek(SIZE-1, 0) 
   if err != nil { 
         fmt.Println(err) 
         log.Fatal("Failed to seek") 
   } 
 
   _, err = fd.Write([]byte{0}) 
   if err != nil { 
         fmt.Println(err) 
         log.Fatal("Write operation failed") 
   } 
 
   err = fd.Close() 
   if err != nil { 
         fmt.Println(err) 
         log.Fatal("Failed to close file") 
   } 
} 
```

First, you try to create the desired sparse file using os.Create(). Then, you call fd.Seek() in order to make the file bigger without adding actual data. Lastly, you write a byte to it using fd.Write(). As you do not have anything more to do with the file, you call fd.Close() and you are done.

Executing sparse.go generates the following output:

```markup
$ go run sparse.go 1000 test
$ go run sparse.go 1000 test
File test already exists.
exit status 1
```

How can you tell whether a file is a sparse file or not? You will learn this in a while, but first, let's create some files:

```markup
$ go run sparse.go 100000 testSparse$ dd if=/dev/urandom  bs=1 count=100000 of=noSparseDD
100000+0 records in
100000+0 records out
100000 bytes (100 kB) copied, 0.152511 s, 656 kB/s
$ dd if=/dev/urandom seek=100000 bs=1 count=0 of=sparseDD
0+0 records in
0+0 records out
0 bytes (0 B) copied, 0.000159399 s, 0.0 kB/s
$ ls -l noSparseDD sparseDD testSparse
-rw-r--r-- 1 mtsouk mtsouk 100000 Apr 29 21:43 noSparseDD
-rw-r--r-- 1 mtsouk mtsouk 100000 Apr 29 21:43 sparseDD
-rw-r--r-- 1 mtsouk mtsouk 100000 Apr 29 21:40 testSparse
```

Note that some Unix variants will not create sparse files: the first such Unix variant that comes to mind is macOS that uses the HFS filesystem. Therefore, for better results, you can execute all these commands on a Linux machine.

So, how can you tell if any of these three files is a sparse file or not? The \-s flag of the ls(1) utility shows the number of filesystem blocks actually used by a file. So, the output of the ls -ls command allows you to detect if you are dealing with a sparse file or not:

```markup
$ ls -ls noSparseDD sparseDD testSparse
104 -rw-r--r-- 1 mtsouk mtsouk 100000 Apr 29 21:43 noSparseDD
      0 -rw-r--r-- 1 mtsouk mtsouk 100000 Apr 29 21:43 sparseDD
      8 -rw-r--r-- 1 mtsouk mtsouk 100000 Apr 29 21:40 testSparse
```

Now look at the first column of the output. The noSparseDD file, which was generated using the dd(1) utility, is not a sparse file. The sparseDD file is a sparse file generated using the dd(1) utility. Lastly, the testSparse is also a sparse file that was created using sparse.go.

Just Imagine

# Reading and writing data records

This section will teach you how to deal with writing and reading data records. What differentiates a record from other kinds of text data is that a record has a given structure with a specific number of fields: think of it as a row from a table in a relational database. Actually, records can be very useful for storing data in tables in case you want to develop your own database server in Go!

The Go code of records.go will save data in the CSV format and will be presented in four parts. The first part contains the following Go code:

```markup
package main 
 
import ( 
   "encoding/csv" 
   "fmt" 
   "os" 
) 
```

So, this is where you have to declare that you are going to read or write data in the CSV format. The second part is the following:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Println("Need just one filename!") 
         os.Exit(-1) 
   } 
 
   filename := os.Args[1] 
   _, err := os.Stat(filename) 
   if err == nil { 
         fmt.Printf("File %s already exists.\n", filename) 
         os.Exit(1) 
   } 
```

The third part of the program is as follows:

```markup
   output, err := os.Create(filename) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(-1) 
   } 
   defer output.Close() 
 
   inputData := [][]string{{"M", "T", "I."}, {"D", "T", "I."}, 
{"M", "T", "D."}, {"V", "T", "D."}, {"A", "T", "D."}} 
   writer := csv.NewWriter(output) 
   for _, record := range inputData { 
         err := writer.Write(record) 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(-1) 
         } 
   } 
   writer.Flush() 
```

You should be familiar with the operations in this part; the biggest difference from what you have seen so far in this chapter is that the writer is from the csv package.

The last part of records.go has the following Go code:

```markup
   f, err := os.Open(filename) 
   if err != nil { 
         fmt.Println(err) 
         return 
   } 
   defer f.Close() 
 
   reader := csv.NewReader(f) 
   reader.FieldsPerRecord = -1 
   allRecords, err := reader.ReadAll() 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
 
   for _, rec := range allRecords { 
         fmt.Printf("%s:%s:%s\n", rec[0], rec[1], rec[2]) 
   } 
} 
```

The reader reads the entire file at once to make the whole operation faster. However, if you are dealing with huge data files, you might need to read smaller parts of the file each time until you have read the complete file. The used reader is from the csv package.

Executing records.go will create the following output:

```markup
$ go run records.go recordsDataFile
M:T:I.                           
D:T:I.
M:T:D.
V:T:D.
A:T:D.
$ ls -l recordsDataFile
-rw-r--r--  1 mtsouk  staff  35 May  2 19:20 recordsDataFile
```

The CSV file, which is named recordsDataFile, contains the following data:

```markup
$ cat recordsDataFile
M,T,I.
D,T,I.
M,T,D.
V,T,D.
A,T,D.
```

Just Imagine

# File locking in Go

There are times that you do not want any other child of the same process to change a file or even access it because you are changing its data and you do not want the other processes to read incomplete or inconsistent data. Although you will learn more about file locking and go routines in [Chapter 9](https://subscription.imaginedevops.io/book/programming/9781787125643/9), _Goroutines - Basic Features_ and [Chapter 10](https://subscription.imaginedevops.io/book/programming/9781787125643/10), _Goroutines - Advanced Features_, this chapter will present a small Go example without a detailed explanation in order to give you an idea about how things work: you should wait until [Chapter 9](https://subscription.imaginedevops.io/book/programming/9781787125643/9), _Goroutines - Basic Features_ and [Chapter 10](https://subscription.imaginedevops.io/book/programming/9781787125643/10), _Goroutines - Advanced Features_, to learn more.

The presented technique will use Mutex, which is a general synchronization mechanism. The Mutex lock will allow us to lock a file from within the same Go process. As a result, this technique has nothing to do with the use of the flock(2) system call.

Various techniques exist for file locking. One of them is by creating an additional file that signifies that another program or process is using a given resource. The presented technique is more suitable for programs that use multiple go routines.

The file locking technique for writing will be illustrated in fileLocking.go, which will be presented in four parts. The first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "math/rand" 
   "os" 
   "sync" 
   "time" 
) 
 
var mu sync.Mutex 
 
func random(min, max int) int { 
   return rand.Intn(max-min) + min 
} 
```

The second part is the following:

```markup
func writeDataToFile(i int, file *os.File, w *sync.WaitGroup) { 
   mu.Lock() 
   time.Sleep(time.Duration(random(10, 1000)) * time.Millisecond) 
   fmt.Fprintf(file, "From %d, writing %d\n", i, 2*i) 
   fmt.Printf("Wrote from %d\n", i) 
   w.Done() 
mu.Unlock() 
} 
```

The locking of the file is done using the mu.Lock() statement and the unlocking of the file with the mu.Unlock() statement.

The third part is the following:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Println("Please provide one command line argument!") 
         os.Exit(-1) 
   } 
 
   filename := os.Args[1] 
   number := 3 
 
   file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
```

The last part is the following Go code:

```markup
   var w *sync.WaitGroup = new(sync.WaitGroup) 
   w.Add(number) 
 
   for r := 0; r < number; r++ { 
         go writeDataToFile(r, file, w) 
   } 
 
   w.Wait() 
} 
```

Executing fileLocking.go will create the following output:

```markup
$ go run fileLocking.go 123
Wrote from 0
Wrote from 2
Wrote from 1
$ cat /tmp/swtag.log
From 0, writing 0
From 2, writing 4
From 1, writing 2
```

The correct version of fileLocking.go has a call to mu.Unlock() at the end of the writeDataToFile() function, which allows all goroutines to use the file. If you remove that call to mu.Unlock() from the writeDataToFile() function, and execute fileLocking.go, you will get the following output:

```markup
$ go run fileLocking.go 123
Wrote from 2
fatal error: all goroutines are asleep - deadlock!
    
goroutine 1 [semacquire]:
sync.runtime_Semacquire(0xc42001024c)
      /usr/local/Cellar/go/1.8.1/libexec/src/runtime/sema.go:47 +0x34
sync.(*WaitGroup).Wait(0xc420010240)
      /usr/local/Cellar/go/1.8.1/libexec/src/sync/waitgroup.go:131 +0x7a
main.main()
     /Users/mtsouk/Desktop/goCourse/ch/ch6/code/fileLocking.go:47 +0x172
    
goroutine 5 [semacquire]:
sync.runtime_SemacquireMutex(0x112dcbc)
     /usr/local/Cellar/go/1.8.1/libexec/src/runtime/sema.go:62 +0x34
sync.(*Mutex).Lock(0x112dcb8)
      /usr/local/Cellar/go/1.8.1/libexec/src/sync/mutex.go:87 +0x9d
main.writeDataToFile(0x0, 0xc42000c028, 0xc420010240)
      /Users/mtsouk/Desktop/goCourse/ch/ch6/code/fileLocking.go:18 +0x3f
created by main.main
      /Users/mtsouk/Desktop/goCourse/ch/ch6/code/fileLocking.go:44 +0x151
    
goroutine 6 [semacquire]:
sync.runtime_SemacquireMutex(0x112dcbc)
      /usr/local/Cellar/go/1.8.1/libexec/src/runtime/sema.go:62 +0x34
sync.(*Mutex).Lock(0x112dcb8)
      /usr/local/Cellar/go/1.8.1/libexec/src/sync/mutex.go:87 +0x9d
main.writeDataToFile(0x1, 0xc42000c028, 0xc420010240)
      /Users/mtsouk/Desktop/goCourse/ch/ch6/code/fileLocking.go:18 +0x3f
created by main.main
      /Users/mtsouk/Desktop/goCourse/ch/ch6/code/fileLocking.go:44 +0x151exit status 2
$ cat 123
From 2, writing 4
```

The reason for getting this output is that apart from the first goroutine that will be able to execute the mu.Lock() statement, the rest of them cannot get Mutex. Therefore, they cannot write to the file, which means that they will never finish their jobs and wait forever, which is the reason that Go is generating the aforementioned error messages.

If you do not completely understand this example, you should wait until [Chapter 9](https://subscription.imaginedevops.io/book/programming/9781787125643/9), _Goroutines - Basic Features_ and [Chapter 10](https://subscription.imaginedevops.io/book/programming/9781787125643/10), _Goroutines - Advanced Features_.

Just Imagine

# A simplified Go version of the dd utility

The dd(1) tool can do many things, but this section will implement a small part of its functionality. Our version of dd(1) will include support for two command-line flags: one for specifying the block size in bytes (\-bs) and the other for specifying the total number of blocks that will be written (\-count). Multiplying these two values will give you the size of the generated file in bytes.

The Go code is saved as ddGo.go and will be presented to you in four parts. The first part is the expected preamble:

```markup
package main 
 
import ( 
   "flag" 
   "fmt" 
   "math/rand" 
   "os" 
   "time" 
) 
```

The second part contains the Go code of two functions:

```markup
func random(min, max int) int { 
   return rand.Intn(max-min) + min 
} 
 
func createBytes(buf *[]byte, count int) { 
   if count == 0 { 
         return 
   } 
   for i := 0; i < count; i++ { 
         intByte := byte(random(0, 9)) 
         *buf = append(*buf, intByte) 
   } 
} 
```

The first function is for getting random numbers and the second one is for creating a byte slice with the desired size filled with random numbers.

The third part of ddGo.go is the following:

```markup
func main() { 
   minusBS := flag.Int("bs", 0, "Block Size") 
   minusCOUNT := flag.Int("count", 0, "Counter") 
   flag.Parse() 
   flags := flag.Args() 
 
   if len(flags) == 0 { 
         fmt.Println("Not enough arguments!") 
         os.Exit(-1) 
   } 
 
   if *minusBS < 0 || *minusCOUNT < 0 { 
         fmt.Println("Count or/and Byte Size < 0!") 
         os.Exit(-1) 
   } 
 
   filename := flags[0] 
   rand.Seed(time.Now().Unix()) 
 
   _, err := os.Stat(filename) 
   if err == nil { 
         fmt.Printf("File %s already exists.\n", filename) 
         os.Exit(1) 
   } 
 
   destination, err := os.Create(filename) 
   if err != nil { 
         fmt.Println("os.Create:", err) 
         os.Exit(1) 
   } 
```

Here, you mainly deal with the command-line arguments of the program.

The last part is the following:

```markup
   buf := make([]byte, *minusBS) 
   buf = nil 
   for i := 0; i < *minusCOUNT; i++ { 
         createBytes(&buf, *minusBS) 
         if _, err := destination.Write(buf); err != nil { 
               fmt.Println(err) 
               os.Exit(-1) 
         } 
         buf = nil 
   } 
} 
```

The reason for emptying the buf byte slice each time you want to call createBytes() is that you do not want the buf byte slice to get bigger and bigger each time you call the createBytes() function. This happens because the append() function adds data at the end of a slice without touching the existing data.

In the first version of ddGo.go that I wrote, I forgot to empty the buf byte slice before each call to createBytes(). Consequently, the generated files were bigger than expected! It took me a while and a couple of fmt.Println(buf) statements to find out the reason for this unforeseen behavior!

The execution of ddGo.go will generate the files you want quite fast:

```markup
$ time go run ddGo.go -bs=10000 -count=5000 test3
    
real  0m1.655s
user  0m1.576s
sys   0m0.104s
$ ls -l test3
-rw-r--r--  1 mtsouk  staff  50000000 May  6 15:27 test3
```

Additionally, the use of random numbers makes the generated files of the same size different from each other:

```markup
$ go run ddGo.go -bs=100 -count=50 test1
$ go run ddGo.go -bs=100 -count=50 test2
$ ls -l test1 test2
-rw-r--r--  1 mtsouk  staff  5000 May  6 15:26 test1
-rw-r--r--  1 mtsouk  staff  5000 May  6 15:26 test2
$ diff test1 test2
Binary files test1 and test2 differ
```

Just Imagine

# Exercises

1.  Visit the documentation page of the bufio package that can be found at [https://golang.org/pkg/bufio/](https://golang.org/pkg/bufio/).
2.  Visit the documentation of the io package at [https://golang.org/pkg/io/](https://golang.org/pkg/io/).
3.  Try to make wc.go faster.
4.  Implement the functionality of tabSpace.go, but try to read your input text files character by character instead of line by line.
5.  Change the code of tabSpace.go in order to be able to get the number of spaces that will replace a tab as a command-line argument.
6.  Learn more information about the little endian and the big endian representations.

Just Imagine

# Summary

In this chapter, we talked about file input and output in Go. Among other things, we developed Go versions of the wc(1), dd(1), and cp(1) Unix command-line utilities while learning more about the io and bufio packages of the Go standard library, which allow you to read from and write to files.

In the next chapter, we will talk about another important subject, which is the Go way of working with the system files of a Unix machine. Additionally, you will learn how to read and change the Unix file permissions as well as how to find the owner and the group of a file. Also, we will talk about log files and how you can use pattern matching to acquire the information you want from log files.