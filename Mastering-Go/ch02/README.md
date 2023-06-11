

Data is stored and used in variables and all Go variables should have a data type that is determined either implicitly or explicitly. Knowing the built-in data types of Go allows you to understand how to manipulate simple data values and construct more complex data structures when simple data types are not enough or not efficient for a given job.

This chapter is all about the basic data types of Go and the data structures that allow you to **group data of the same data type**. But let us begin with something more practical: imagine that you want to read data as command-line arguments of a utility. How can you be sure that what you have read was what you expected? How can you handle error situations? What about reading not just numbers and strings but dates and times from the command line? Do you have to write your own parser for working with dates and times?

This chapter will answer all these questions and many more by implementing the following three utilities:

-   A command-line utility that parses dates and times
-   A utility that generates random numbers and random strings
-   A new version of the phone book application that contains randomly generated data

This chapter covers:

-   The `error` data type
-   Numeric data types
-   Non-numeric data types
-   Go constants
-   Grouping similar data
-   Pointers
-   Generating random numbers
-   Updating the phone book application

We begin this chapter with the `error` data type, because errors play a key role in Go.

Bookmark

# The error data type

Go provides a special data type for representing error conditions and error messages named `error`—in practice, this means that Go treats errors as values. In order to program successfully in Go, you should be aware of the error conditions that might occur with the functions and methods you are using and handle them accordingly.

As you already know from the previous chapter, Go follows the next convention about `error` values: if the value of an `error` variable is `nil`, then there was no error. As an example, let us consider `strconv.Atoi()`, which is used for converting a `string` value into an `int` value (`Atoi` stands for ASCII to Int). As specified by its signature, `strconv.Atoi()` returns `(int, error)`. Having an `error` value of `nil` means that the conversion was successful and that you can use the `int` value if you want. Having an `error` value that is not `nil` means that the conversion was unsuccessful and that the `string` input is not a valid `int` value.

If you want to learn more about `strconv.Atoi()`, you should execute `go doc strconv.Atoi` in your terminal window.

You might wonder what happens if you want to create your own error messages. Is this possible? Should you wish to return a custom error, you can use `errors.New()` from the `errors` package. This usually happens inside a function other than `main()` because `main()` does not return anything to any other function. Additionally, a good place to define your custom errors is inside the Go packages you create.

You will most likely work with errors in your programs without needing the functionality of the `errors` package. Additionally, you do not need to define custom error messages unless you are creating big applications or packages.

If you want to format your error messages in the way `fmt.Printf()` works, you can use the `fmt.Errorf()` function, which simplifies the creation of custom error messages—the `fmt.Errorf()` function returns an `error` value just like `errors.New()`.

And now we should talk about something important: you should have a global error handling tactic in each application that should not change. In practice, this means the following:

-   All error messages should be handled at the same level, which means that all errors should either be returned to the calling function or be handled at the place they occurred.
-   It should be clearly documented how to handle critical errors. This means that there will be situations where a critical error should terminate the program and other times where a critical error might just create a warning message onscreen.
-   It is considered a good practice to send all error messages to the _log service_ of your machine because this way the error messages can be examined at a later time. However, this is not always true, so exercise caution when setting this up—for example, cloud native apps do not work that way.

The `error` data type is actually defined as an _interface_—interfaces are covered in _Chapter 4_, _Reflection and Interfaces_.

Type the following code in your favorite text editor and save it as `error.go` in the directory where you put the code for this chapter. Using `ch02` as the directory name is a good idea.

```markup
package main
import (
    "errors"
    "fmt"
    "os"
    "strconv"
)
```

The first part is the preamble of the program—`error.go` uses the `fmt`, `os`, `strconv`, and `errors` packages.

```markup
// Custom error message with errors.New()
func check(a, b int) error {
    if a == 0 && b == 0 {
        return errors.New("this is a custom error message")
    }
    return nil
}
```

The preceding code implements a function named `check()` that returns an `error` value. If both input parameters of `check()` are equal to `0`, the function returns a custom error message using `errors.New()`—otherwise it returns `nil`, which means that everything is OK.

```markup
// Custom error message with fmt.Errorf()
func formattedError(a, b int) error {
    if a == 0 && b == 0 {
        return fmt.Errorf("a %d and b %d. UserID: %d", a, b, os.Getuid())
    }
    return nil
}
```

The previous code implements `formattedError()`, which is a function that returns a formatted error message using `fmt.Errorf()`. Among other things, the error message prints the user ID of the user that executed the program with a call to `os.Getuid()`. When you want to create a custom error message, using `fmt.Errorf()` gives you more control over the output.

```markup
func main() {
    err := check(0, 10)
    if err == nil {
        fmt.Println("check() ended normally!")
    } else {
        fmt.Println(err)
    }
    err = check(0, 0)
    if err.Error() == "this is a custom error message" {
        fmt.Println("Custom error detected!")
    }
    err = formattedError(0, 0)
    if err != nil {
        fmt.Println(err)
    }
    i, err := strconv.Atoi("-123")
    if err == nil {
        fmt.Println("Int value is", i)
    }
    i, err = strconv.Atoi("Y123")
    if err != nil {
        fmt.Println(err)
    }
}
```

The previous code is the implementation of the `main()` function where you can see the use of the `if err != nil` statement multiple times as well as the use of `if err == nil`, which is used to make sure that everything was OK before executing the desired code.

Running `error.go` produces the next output:

```markup
$ go run error.go
check() ended normally!
Custom error detected!
a 0 and b 0. UserID: 501
Int value is -123
strconv.Atoi: parsing "Y123": invalid syntax
```

Now that you know about the `error` data type, how to create custom errors, and how to use `error` values, we'll continue with the basic data types of Go that can be logically divided into two main categories: numeric data types and non-numeric data types. Go also supports the `bool` data type, which can have a value of `true` or `false` only.

Bookmark

# Numeric data types

Go supports integer, floating-point, and complex number values in various versions depending on the memory space they consume—this saves memory and computing time. Integer data types can be either signed or unsigned, which is not the case for floating point numbers.

The table that follows lists the numeric data types of Go.

<table id="table001" class="No-Table-Style _idGenTablePara-1"><colgroup><col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Heading--PACKT-">Data Type</p></td><td class="No-Table-Style"><p class="Table-Column-Heading--PACKT-">Description</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">int8</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">8-bit signed integer</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">int16</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">16-bit signed integer</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">int32</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">32-bit signed integer</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">int64</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">64-bit signed integer</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">int</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">32- or 64-bit signed integer</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">uint8</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">8-bit unsigned integer</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">uint16</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">16-bit unsigned integer</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">uint32</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">32-bit unsigned integer</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">uint64</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">64-bit unsigned integer</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">uint</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">32- or 64-bit unsigned integer</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">float32</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">32-bit floating-point number</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">float64</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">64-bit floating-point number</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">complex64</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Complex number with <code class="Code-In-Text--PACKT-">float32</code> parts</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">Complex128</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Complex number with <code class="Code-In-Text--PACKT-">float64</code> parts</p></td></tr></tbody></table>

The `int` and `uint` data types are special as they are the most efficient sizes for signed and unsigned integers on a given platform and can be either 32 or 64 bits each—their size is defined by Go itself. The `int` data type is the most widely used data type in Go due to its versatility.

The code that follows illustrates the use of numeric data types—you can find the entire program as `numbers.go` inside the `ch02` directory of the book GitHub repository.

```markup
func main() {
    c1 := 12 + 1i
    c2 := complex(5, 7)
    fmt.Printf("Type of c1: %T\n", c1)
    fmt.Printf("Type of c2: %T\n", c2)
```

The previous code creates two complex variables in two different ways—both ways are perfectly valid and equivalent. Unless you are into mathematics, you will most likely not use complex numbers in your programs. However, the existence of complex numbers shows how modern Go is.

```markup
    var c3 complex64 = complex64(c1 + c2)
    fmt.Println("c3:", c3)
    fmt.Printf("Type of c3: %T\n", c3)
    cZero := c3 - c3
    fmt.Println("cZero:", cZero)
```

The previous code continues to work with complex numbers by adding and subtracting two pairs of them. Although `cZero` is equal to zero, it is still a complex number and a `complex64` variable.

```markup
    x := 12
    k := 5
    fmt.Println(x)
    fmt.Printf("Type of x: %T\n", x)
    div := x / k
    fmt.Println("div", div)
```

In this part, we define two integer variables named `x` and `k`—their data type is identified by Go based on their initial values. Both are of type `int`, which is what Go prefers to use for storing integer values. Additionally, when you divide two integer values, you get an integer result even when the division is not perfect. This means that if this is not what you want, you should take extra care—this is shown in the next code excerpt:

```markup
    var m, n float64
    m = 1.223
    fmt.Println("m, n:", m, n)
    y := 4 / 2.3
    fmt.Println("y:", y)
    divFloat := float64(x) / float64(k)
    fmt.Println("divFloat", divFloat)
    fmt.Printf("Type of divFloat: %T\n", divFloat)
}
```

The previous code works with `float64` values and variables. As `n` does not have an initial value, it is **automatically assigned** with the zero value of its data type, which is `0` for the `float64` data type.

Additionally, the code presents a technique for dividing integer values and getting a floating-point result, which is the use of `float64()`: `divFloat := float64(x) / float64(k)`. This is a type conversion where two integers (`x` and `k`) are converted to `float64` values. As the division between two `float64` values is a `float64` value, we get the result in the desired data type.

Running `numbers.go` creates the following output:

```markup
$ go run numbers.go
Type of c1: complex128
Type of c2: complex128
c3: (17+8i)
Type of c3: complex64
cZero: (0+0i)
12
Type of x: int
div 2
m, n: 1.223 0
y: 1.7391304347826086
divFloat 2.4
Type of divFloat: float64
```

The output shows that both `c1` and `c2` are `complex128` values, which is the preferred complex data type for the machine on which the code was executed. However, `c3` is a `complex64` value because it was created using `complex64()`. The value of `n` is `0` because the `n` variable was not initialized, which means that Go automatically assigned the zero value of its data type to `n`.

After learning about numeric data types, it is time to learn about non-numeric data types, which is the subject of the next section.

Bookmark

# Non-numeric data types

Go has support for **Strings**, **Characters**, **Runes**, **Dates**, and **Times**. However, Go does not have a dedicated `char` data type. We begin by explaining the string-related data types.

For Go, dates and times are the same thing and are represented by the same data type. However, it is up to you to determine whether a time and date variable contains valid information or not.

## Strings, Characters, and Runes

Go supports the `string` data type for representing strings. A Go string is just a collection of bytes and can be accessed as a whole or as an array. A single byte can store any ASCII character—however, multiple bytes are usually needed for storing a single Unicode character.

Nowadays, supporting Unicode characters is a common requirement—Go is designed with Unicode support in mind, which is the main reason for having the rune data type. A rune is an `int32` value that is used for representing a single Unicode code point, which is an integer value that is used for representing single Unicode characters or, less frequently, providing formatting information.

Although a rune is an `int32` value, you cannot compare a rune with an `int32` value. Go considers these two data types as totally different.

You can create a new byte slice from a given string by using a `[]byte("A String")` statement. Given a byte slice variable `b`, you can convert it into a string using the `string(b)` statement. When working with byte slices that contain Unicode characters, the number of bytes in a byte slice is not always connected to the number of characters in the byte slice, because most Unicode characters require more than one byte for their representation. As a result, when you try to print each single byte of a byte slice using `fmt.Println()` or `fmt.Print()`, the output is not text presented as characters but integer values. If you want to print the contents of a byte slice as text, you should either print it using `string(byteSliceVar)` or using `fmt.Printf()` with `%s` to tell `fmt.Printf()` that you want to print a string. You can initialize a new byte slice with the desired string by using a statement such as `[]byte("My Initialization String")`.

We will cover byte slices in more detail in the _Byte slices_ section.

You can define a rune using single quotes: `r := '€'` and you can print the integer value of the bytes that compose it as `fmt.Println(r)`—in this case, the integer value is `8364`. Printing it as a single Unicode character requires the use of the `%c` control string in `fmt.Printf()`.

As strings can be accessed as arrays, you can iterate over the runes of the string using a `for` loop or point to a specific character if you know its place in the string. The length of the string is the same as the number of characters found in the string, which is usually not true for byte slices because Unicode characters usually require more than one byte.

The following Go code illustrates the use of strings and runes and how you can work with strings in your code. You can find the entire program as `text.go` in the `ch02` directory of the GitHub repository of the book.

The first part of the program defines a string literal that contains a Unicode character. Then it accesses its first character as if the string was an array.

```markup
func main() {
    aString := "Hello World! €"
    fmt.Println("First character", string(aString[0]))
```

The next part is about working with runes.

```markup
    // Runes
    // A rune
    r := '€'
    fmt.Println("As an int32 value:", r)
    // Convert Runes to text
    fmt.Printf("As a string: %s and as a character: %c\n", r, r)
    // Print an existing string as runes
    for _, v := range aString {
        fmt.Printf("%x ", v)
    }
    fmt.Println()
```

First, we define a rune named `r`. What makes this a rune is the use of single quotes around the `€` character. The rune is an `int32` value and is printed as such by `fmt.Println()`. The `%c` control string in `fmt.Printf()` prints a rune as a character.

Then we iterate over `aString` as a slice or an array using a `for` loop with `range` and print the contents of `aString` as runes.

```markup
    // Print an existing string as characters
    for _, v := range aString {
        fmt.Printf("%c", v)
    }
    fmt.Println()
}
```

Lastly, we iterate over `aString` as a slice or an array using a `for` loop with `range` and print the contents of `aString` as characters.

Running `text.go` produces the following output:

```markup
$ go run text.go
First character H
As an int32 value: 8364
As a string: %!s(int32=8364) and as a character: €
48 65 6c 6c 6f 20 57 6f 72 6c 64 21 20 20ac
Hello World! €
```

The first line of the output shows that we can access a `string` as an array whereas the second line verifies that a rune is an integer value. The third line shows what to expect when you print a `rune` as a `string` and as a character—the correct way is to print it as a character. The fifth line shows how to print a string as runes and the last line shows the output of processing a string as characters using `range` and a `for` loop.

### Converting from int to string

You can convert an integer value into a string in two main ways: using `string()` and using a function from the `strconv` package. However, the two methods are fundamentally different. The `string()` function converts an integer value into a Unicode code point, which is a single character, whereas functions such as `strconv.FormatInt()` and `strconv.Itoa()` convert an integer value into a string value with the same representation and the same number of characters.

This is illustrated in the `intString.go` program—its most important statements are the following. You can find the entire program in the GitHub repository of the book.

```markup
    input := strconv.Itoa(n)
    input = strconv.FormatInt(int64(n), 10)
    input = string(n)
```

Running `intString.go` generates the following kind of output:

```markup
$ go run intString.go 100
strconv.Itoa() 100 of type string
strconv.FormatInt() 100 of type string
string() d of type string
```

The data type of the output is always `string`, however, `string()` converted `100` into `d` because the ASCII representation of `d` is `100`.

### The unicode package

The `unicode` standard Go package contains various handy functions for working with Unicode code points. One of them, which is called `unicode.IsPrint()`, can help you to identify the parts of a string that are printable using runes.

The following code excerpt illustrates the functionality of the `unicode` package:

```markup
    for i := 0; i < len(sL); i++ {
        if unicode.IsPrint(rune(sL[i])) {
            fmt.Printf("%c\n", sL[i])
        } else {
            fmt.Println("Not printable!")
        }
    }
```

The `for` loop iterates over the contents of a string defined as a list of runes (`"\x99\x00ab\x50\x00\x23\x50\x29\x9c"`) while `unicode.IsPrint()` examines whether the character is printable or not—if it returns `true` then a rune is printable.

You can find this code excerpt inside the `unicode.go` source file at the `ch02` directory in the GitHub repository of the book. Running `unicode.go` produces the following output:

```markup
Not printable!
Not printable!
a
b
P
Not printable!
#
P
)
Not printable!
```

This utility is very handy for filtering your input or filtering data before printing it on screen, storing it in log files, transferring it on a network, or storing it in a database.

### The strings package

The `strings` standard Go package allows you to manipulate UTF-8 strings in Go and includes many powerful functions. Many of these functions are illustrated in the `useStrings.go` source file, which can be found in the `ch02` directory of the book GitHub repository.

If you are working with text and text processing, you definitely need to learn all the gory details and functions of the `strings` package, so make sure that you experiment with all these functions and create many examples that will help you to clarify things.

The most important parts of `useStrings.go` are the following:

```markup
import (
    "fmt"
    s "strings"
    "unicode"
)
var f = fmt.Printf
```

As we are going to use the `strings` package multiple times, we create a convenient alias for it named `s`. We do the same for the `fmt.Printf()` function where we create a global alias using a variable named `f`. These two shortcuts make code less populated with long, repeated lines of code. You can use it when learning Go but this is not recommended in any kind of production software, as it makes code less readable.

The first code excerpt is the following.

```markup
f("EqualFold: %v\n", s.EqualFold("Mihalis", "MIHAlis"))
f("EqualFold: %v\n", s.EqualFold("Mihalis", "MIHAli"))
```

The `strings.EqualFold()` function compares two strings without considering their case and returns `true` when they are the same and `false` otherwise.

```markup
f("Index: %v\n", s.Index("Mihalis", "ha"))
f("Index: %v\n", s.Index("Mihalis", "Ha"))
```

The `strings.Index()` function checks whether the string of the second parameter can be found in the string that is given as the first parameter and returns the index where it was found for the first time. On an unsuccessful search, it returns `-1`.

```markup
    f("Prefix: %v\n", s.HasPrefix("Mihalis", "Mi"))
    f("Prefix: %v\n", s.HasPrefix("Mihalis", "mi"))
    f("Suffix: %v\n", s.HasSuffix("Mihalis", "is"))
    f("Suffix: %v\n", s.HasSuffix("Mihalis", "IS"))
```

The `strings.HasPrefix()` function checks whether the given string, which is the first parameter, begins with the string that is given as the second parameter. In the previous code, the first call to `strings.HasPrefix()` returns `true`, whereas the second returns `false`.

Similarly, the `strings.HasSuffix()` function checks whether the given string ends with the second string. Both functions take into account the case of the input string and the case of the second parameter.

```markup
    t := s.Fields("This is a string!")
    f("Fields: %v\n", len(t))
    t = s.Fields("ThisIs a\tstring!")
    f("Fields: %v\n", len(t))
```

The handy `strings.Fields()` function splits the given string around one or more white space characters as defined by the `unicode.IsSpace()` function and returns a slice of substrings found in the input string. If the input string contains white characters only, it returns an empty slice.

```markup
    f("%s\n", s.Split("abcd efg", ""))
    f("%s\n", s.Replace("abcd efg", "", "_", -1))
    f("%s\n", s.Replace("abcd efg", "", "_", 4))
    f("%s\n", s.Replace("abcd efg", "", "_", 2))
```

The `strings.Split()` function allows you to split the given string according to the desired separator string—the `strings.Split()` function returns a **string slice**. Using `""` as the second parameter of `strings.Split()` allows you to process a string character by character.

The `strings.Replace()` function takes four parameters. The first parameter is the string that you want to process. The second parameter contains the string that, if found, will be replaced by the third parameter of `strings.Replace()`. The last parameter is the maximum number of replacements that are allowed to happen. If that parameter has a negative value, then there is no limit to the number of replacements that can take place.

```markup
    f("SplitAfter: %s\n", s.SplitAfter("123++432++", "++"))
    trimFunction := func(c rune) bool {
        return !unicode.IsLetter(c)
    }
    f("TrimFunc: %s\n", s.TrimFunc("123 abc ABC \t .", trimFunction))
```

The `strings.SplitAfter()` function splits its first parameter string into substrings based on the separator string that is given as the second parameter to the function. The separator string is included in the returned slice.

The last lines of code define a **trim function** named `trimFunction` that is used as the second parameter to `strings.TrimFunc()` in order to filter the given input based on the return value of the trim function—in this case, the trim function keeps all letters and nothing else due to the `unicode.IsLetter()` call.

Running `useStrings.go` produces the next output:

```markup
To Upper: HELLO THERE!
To Lower: hello there
THis WiLL Be A Title!
EqualFold: true
EqualFold: false
Prefix: true
Prefix: false
Suffix: true
Suffix: false
Index: 2
Index: -1
Count: 2
Count: 0
Repeat: ababababab
TrimSpace: This is a line.
TrimLeft: This is a      line. 
TrimRight:      This is a        line.
Compare: 1
Compare: 0
Compare: -1
Fields: 4
Fields: 3
[a b c d   e f g]
_a_b_c_d_ _e_f_g_
_a_b_c_d efg
_a_bcd efg
Join: Line 1+++Line 2+++Line 3
SplitAfter: [123++ 432++ ]
TrimFunc: abc ABC
```

Visit the documentation page of the `strings` package at [https://golang.org/pkg/strings/](https://golang.org/pkg/strings/) for the complete list of available functions. You will see the functionality of the `strings` package in other places in this book.

Enough with strings and text; the next section is about working with dates and times in Go.

## Times and dates

Often, we need to work with date and time information to store the time an entry was last used in a database or the time an entry was inserted into a database, which brings us to another interesting topic: working with dates and times in Go.

The king of working with times and dates in Go is the `time.Time` data type, which represents an instant in time with _nanosecond precision_. Each `time.Time` value is associated with a location (time zone).

If you are a UNIX person, you might already know about the UNIX epoch time and wonder how to get it in Go. The `time.Now().Unix()` function returns the popular UNIX epoch time, which is the number of seconds that have elapsed since 00:00:00 UTC, January 1, 1970. If you want to convert the UNIX time to the equivalent `time.Time` value, you can use the `time.Unix()` function. If you are not a UNIX person, then you might not have heard about the UNIX epoch time before but now you know what it is!

The `time.Since()` function calculates the time that has passed since a given time and returns a `time.Duration` variable—the `duration` data type is defined as `type Duration int64`. Although a `Duration` is, in reality, an `int64` value, you cannot compare or convert a `duration` to an `int64` value implicitly because **Go does not allow implicit data type conversions**.

The single most important topic about Go and dates and times is the way Go parses a string in order to convert it into a date and a time. The reason that this is important is usually such input is given as a string and not as a valid date variable. The function used for parsing is `time.Parse()` and its full signature is `Parse(layout, value string) (Time, error)`, where `layout` is the parse string and `value` is the input that is being parsed. The `time.Time` value that is returned is a moment in time with nanosecond precision and contains both date and time information.

The next table shows the most widely used strings for parsing dates and times.

<table id="table002" class="No-Table-Style _idGenTablePara-1"><colgroup><col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Heading--PACKT-">Parse Value</p></td><td class="No-Table-Style"><p class="Table-Column-Heading--PACKT-">Meaning (examples)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">05</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">12-hour value (12pm, 07am)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">15</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">24-hour value (23, 07)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">04</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Minutes (55, 15)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">05</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Seconds (5, 23)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">Mon</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Abbreviated day of week (Tue, Fri)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">Monday</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Day of week (Tuesday, Friday)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">02</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Day of month (15, 31)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">2006</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Year with 4 digits (2020, 2004)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">06</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Year with the last 2 digits (20, 04)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">Jan</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Abbreviated month name (Feb, Mar)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">January</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Full month name (July, August)</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">MST</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Time zone (EST, UTC)</p></td></tr></tbody></table>

The previous table shows that if you want to parse the `30 January 2020` string and convert it into a Go date variable, you should match it against the `02 January 2006` string—you cannot use anything else in its place when matching a string with the `30 January 2020` format. Similarly, if you want to parse the `15 August 2020 10:00` string, you should match it against the `02 January 2006 15:04` string. The documentation of the `time` package ([https://golang.org/pkg/time/](https://golang.org/pkg/time/)) contains even more detailed information about parsing dates and times—however, the ones presented here should be more than enough for regular use.

### A utility for parsing dates and times

On a rare occasion, a situation can happen when we do not know anything about our input. If you do not know the exact format of your input, then you need to try matching your input against multiple Go strings without being sure that you are going to succeed in the end. This is the approach that the example uses. The Go matching strings for dates and times can be tried in any order.

If you are matching a string that only contains the date, then your time will be set to `00:00` by Go and will most likely be incorrect. Similarly, when matching the time only, your date will be incorrect and should not be used.

The formatting strings can be also used for printing dates and times in the desired format. So in order to print the current date in the `01-02-2006` format, you should use `time.Now().Format("01-02-2006")`.

The code that follows illustrates how to work with epoch time in Go and showcases the parsing process—create a text file, type the following code, and save it as `dates.go`.

```markup
package main
import (
    "fmt"
    "os"
    "time"
)
```

This is the expected preamble of the Go source file.

```markup
func main() {
    start := time.Now()
    if len(os.Args) != 2 {
        fmt.Println("Usage: dates parse_string")
        return
    }
    dateString := os.Args[1]
```

This is how we get user input that is stored in the `dateString` variable. If the utility gets no input, there is no point in continuing its operation.

```markup
    // Is this a date only?
    d, err := time.Parse("02 January 2006", dateString)
    if err == nil {
        fmt.Println("Full:", d)
        fmt.Println("Time:", d.Day(), d.Month(), d.Year())
    }
```

The first test is for matching a date only using the `02 January 2006` format. If the match is successful, you can access the individual fields of a variable that holds a valid date using `Day()`, `Month()`, and `Year()`.

```markup
    // Is this a date + time?
    d, err = time.Parse("02 January 2006 15:04", dateString)
    if err == nil {
        fmt.Println("Full:", d)
        fmt.Println("Date:", d.Day(), d.Month(), d.Year())
        fmt.Println("Time:", d.Hour(), d.Minute())
    }
```

This time we try to match a string using `"02 January 2006 15:04"`, which contains a date and a time value. If the match is successful, you can access the fields of a valid time using `Hour()` and `Minute()`.

```markup
    // Is this a date + time with month represented as a number?
    d, err = time.Parse("02-01-2006 15:04", dateString)
    if err == nil {
        fmt.Println("Full:", d)
        fmt.Println("Date:", d.Day(), d.Month(), d.Year())
        fmt.Println("Time:", d.Hour(), d.Minute())
    }
```

This time we try to match against the `"02-01-2006 15:04"` format, which contains both a date and a time. Note that it is compulsory that the string that is being examined contains the `-` and the `:` characters as specified in the `time.Parse()` call and that `"02-01-2006 15:04"` is different from `"02/01/2006 1504"`.

```markup
    // Is it time only?
    d, err = time.Parse("15:04", dateString)
    if err == nil {
        fmt.Println("Full:", d)
        fmt.Println("Time:", d.Hour(), d.Minute())
    }
```

The last match is for time only using the `"15:04"` format. Note that the `:` should exist in the string that is being examined.

```markup
    t := time.Now().Unix()
    fmt.Println("Epoch time:", t)
    // Convert Epoch time to time.Time
    d = time.Unix(t, 0)
    fmt.Println("Date:", d.Day(), d.Month(), d.Year())
    fmt.Printf("Time: %d:%d\n", d.Hour(), d.Minute())
    duration := time.Since(start)
    fmt.Println("Execution time:", duration)
}
```

The last part of `dates.go` shows how to work with UNIX epoch time. You get the current date and time in epoch time using `time.Now().Unix()` and you can convert that to a `time.Time` value using a call to `time.Unix()`.

Lastly, you can calculate the time duration between the current time and a time in the past using a call to `time.Since()`.

Running `dates.go` creates the following kind of output, depending on its input:

```markup
$ go run dates.go 
Usage: dates parse_string
$ go run dates.go 14:10
Full: 0000-01-01 14:10:00 +0000 UTC
Time: 14 10
Epoch time: 1607964956
Date: 14 December 2020
Time: 18:55
Execution time: 163.032µs
$ go run dates.go "14 December 2020"
Full: 2020-12-14 00:00:00 +0000 UTC
Time: 14 December 2020
Epoch time: 1607964985
Date: 14 December 2020
Time: 18:56
Execution time: 180.029µs
```

If a command-line argument such as `14 December 2020` contains space characters, you should put it in double quotes for the UNIX shell to treat it as a single command-line argument. Running `go run dates.go 14 December 2020` does not work.

Now that we know how to work with dates and times, it is time to learn more about time zones.

### Working with different time zones

The presented utility accepts a date and a time and converts them into different time zones. This can be particularly handy when you want to preprocess log files from different sources that use different time zones in order to convert these different time zones into a common one.

Once again, you need `time.Parse()` in order to convert a valid input into a `time.Time` value before doing the conversions. This time the input string contains the time zone and is parsed by the `"02 January 2006 15:04 MST"` string.

In order to convert the parsed date and time into New York time, the program uses the following code:

```markup
    loc, _ = time.LoadLocation("America/New_York")
    fmt.Printf("New York Time: %s\n", now.In(loc))
```

This technique is used multiple times in `convertTimes.go`.

Running `convertTimes.go` generates the following output:

```markup
$ go run convertTimes.go "14 December 2020 19:20 EET"
Current Location: 2020-12-14 19:20:00 +0200 EET
New York Time: 2020-12-14 12:20:00 -0500 EST
London Time: 2020-12-14 17:20:00 +0000 GMT
Tokyo Time: 2020-12-15 02:20:00 +0900 JST
$ go run convertTimes.go "14 December 2020 20:00 UTC"
Current Location: 2020-12-14 22:00:00 +0200 EET
New York Time: 2020-12-14 15:00:00 -0500 EST
London Time: 2020-12-14 20:00:00 +0000 GMT
Tokyo Time: 2020-12-15 05:00:00 +0900 JST
$ go run convertTimes.go "14 December 2020 25:00 EET"
parsing time "14 December 2020 25:00": hour out of range
```

In the last execution of the program, the code has to parse `25` as the hour of the day, which is wrong and generates the `hour out of range` error message.

Bookmark

# Go constants

Go supports **constants**, which are variables that cannot change their values. Constants in Go are defined with the help of the `const` keyword. Generally speaking, constants can be either **global or local variables**.

However, you might need to rethink your approach if you find yourself defining too many constant variables with a local scope. The main benefit you get from using constants in your programs is the guarantee that their value will not change during program execution. Strictly speaking, the value of a constant variable is defined at compile time, not at runtime—this means that it is included in the binary executable. Behind the scenes, Go uses Boolean, string, or number as the type for storing constant values because this gives Go more flexibility when dealing with constants.

The next subsection discusses the constant generator iota, which is a handy way of creating sequences of constants.

## The constant generator iota

The **constant generator iota** is used for declaring a sequence of related values that use incrementing numbers without the need to explicitly type each one of them.

The concepts related to the `const` keyword, including the constant generator iota, are illustrated in the `constants.go` file.

```markup
package main
import (
    "fmt"
)
type Digit int
type Power2 int
const PI = 3.1415926
const (
    C1 = "C1C1C1"
    C2 = "C2C2C2"
    C3 = "C3C3C3"
)
```

In this part, we declare two new types named `Digit` and `Power2` that will be used in a while, and four new constants named `PI`, `C1`, `C2`, and `C3`.

A Go **type** is a way of defining a new **named type** that uses the same underlying type as an existing type. This is mainly used for differentiating between different types that might use the same kind of data. The `type` keyword can be used for defining _structures_ and _interfaces_.

```markup
func main() {
    const s1 = 123
    var v1 float32 = s1 * 12
    fmt.Println(v1)
    fmt.Println(PI)
    const (
        Zero Digit = iota
        One
        Two
        Three
        Four
    )
```

The previous code defines a constant named `s1`. Here you also see the definition of a _constant generator iota_ based on `Digit`, which is equivalent to the next declaration of four constants:

```markup
const (
    Zero = 0
    One = 1
    Two = 2
    Three = 3
    Four = 4
)
```

Although we are defining constants inside `main()`, constants can be normally found outside of `main()` or any other function or method.

The last part of `constants.go` is as follows.

```markup
    fmt.Println(One)
    fmt.Println(Two)
    const (
        p2_0 Power2 = 1 << iota
        _
        p2_2
        _
        p2_4
        _
        p2_6
    )
    fmt.Println("2^0:", p2_0)
    fmt.Println("2^2:", p2_2)
    fmt.Println("2^4:", p2_4)
    fmt.Println("2^6:", p2_6)
}
```

There is another constant generator iota here that is a little different than the previous one. Firstly, you can see the use of the underscore character in a `const` block with a constant generator iota, which allows you to skip unwanted values. Secondly, the value of `iota` always increments and can be used in expressions, which is what occurred in this case.

Now let us see what really happens inside the `const` block. For `p2_0`, `iota` has the value of `0` and `p2_0` is defined as `1`. For `p2_2`, `iota` has the value of `2` and `p2_2` is defined as the result of the expression `1 << 2`, which is `00000100` in binary representation. The decimal value of `00000100` is `4`, which is the result and the value of `p2_2`. Analogously, the value of `p2_4` is `16` and the value of `p2_6` is `64`.

Running `constants.go` produces the next output:

```markup
$ go run constants.go
1476
3.1415926
1
2
2^0: 1
2^2: 4
2^4: 16
2^6: 64
```

Having data is good but what happens when you have lots of similar data? Do you need to have lots of variables to hold this data or is there a better way to do so? Go answers these questions by introducing arrays and slices.

Bookmark

# Grouping similar data

There are times when you want to keep multiple values of the same data type under a single variable and access them using an index number. The simplest way to do that in Go is by using arrays or slices.

Arrays are the most widely used data structures and can be found in almost all programming languages due to their simplicity and speed of access. Go provides an alternative to arrays that is called a slice. The subsections that follow help you understand the differences between arrays and slices so that you know which data structure to use and when.

The quick answer is that you can use slices instead of arrays almost anywhere in Go but we are also demonstrating arrays because they can still be useful and because slices are implemented by Go using arrays!

## Arrays

Arrays in Go have the following characteristics and limitations:

-   When defining an array variable, you must define its size. Otherwise, you should put `[...]` in the array declaration and let the Go compiler find out the length for you. So you can create an array with 4 `string` elements either as `[4]string{"Zero", "One", "Two", "Three"}` or as `[...]string{"Zero", "One", "Two", "Three"}`. If you put nothing in the square brackets, then a slice is going to be created instead. The (valid) indexes for that particular array are `0`, `1`, `2`, and `3`.
-   You cannot change the size of an array after you have created it.
-   When you pass an array to a function, what is happening is that Go creates a copy of that array and passes that copy to that function—therefore any changes you make to an array inside a function are lost when the function returns.

As a result, arrays in Go are not very powerful, which is the main reason that Go has introduced an additional data structure named **slice** that is similar to an array but is dynamic in nature and is explained in the next subsection. However, data in both arrays and slices is accessed the same way.

## Slices

Slices in Go are more powerful than arrays mainly because they are dynamic, which means that they can grow or shrink after creation if needed. Additionally, any changes you make to a slice inside a function also affect the original slice. But how does this happen? Strictly speaking, _all parameters in Go are passed by value_—there is no other way to pass parameters in Go.

In reality, a slice value is a _header_ that contains **a pointer to an underlying array** where the elements are actually stored, the length of the array, and its capacity—the capacity of a slice is explained in the next subsection. Note that the slice value does not include its elements, just a pointer to the underlying array. So, when you pass a slice to a function, Go makes a copy of that header and passes it to the function. This copy of the slice header includes the pointer to the underlying array. That slice header is defined in the `reflect` package ([https://golang.org/pkg/reflect/#SliceHeader](https://golang.org/pkg/reflect/#SliceHeader)) as follows:

```markup
type SliceHeader struct {
    Data uintptr
    Len  int
    Cap  int
}
```

A side effect of passing the slice header is that it is faster to pass a slice to a function because Go does not need to make a copy of the slice and its elements, just the slice header.

You can create a slice using `make()` or like an array without specifying its size or using `[...]`. If you do not want to initialize a slice, then using `make()` is better and faster. However, if you want to initialize it at the time of creation, then `make()` cannot help you. As a result, you can create a slice with three `float64` elements as `aSlice := []float64{1.2, 3.2, -4.5}`. Creating a slice with space for three `float64` elements with `make()` is as simple as executing `make([]float64, 3)`. Each element of that slice has a value of `0`, which is the zero value of the `float64` data type.

Both slices and arrays can have many dimensions—creating a slice with two dimensions with `make()` is as simple as writing `make([][]int, 2)`. This returns a slice with two dimensions where the first dimension is `2` (rows) and the second dimension (columns) is unspecified and should be **explicitly specified** when adding data to it.

If you want to define and initialize a slice with two dimensions at the same time, you should execute something similar to `twoD := [][]int{{1, 2, 3}, {4, 5, 6}}`.

You can find the length of an array or a slice using `len()`. As you will find out in the next subsection, slices have an additional property named _capacity_. You can add new elements to a full slice using the `append()` function. `append()` automatically allocates the required memory space.

The example that follows clarifies many things about slices—feel free to experiment with it. Type the following code and save it as `goSlices.go`.

```markup
package main
import "fmt"
func main() {
    // Create an empty slice
    aSlice := []float64{}
    // Both length and capacity are 0 because aSlice is empty
    fmt.Println(aSlice, len(aSlice), cap(aSlice))
    // Add elements to a slice
    aSlice = append(aSlice, 1234.56)
    aSlice = append(aSlice, -34.0)
    fmt.Println(aSlice, "with length", len(aSlice))
```

The `append()` commands add two new elements to `aSlice`. You should save the return value of `append()` to an existing variable or a new one.

```markup
    // A slice with length 4
    t := make([]int, 4)
    t[0] = -1
    t[1] = -2
    t[2] = -3
    t[3] = -4
    // Now you will need to use append
    t = append(t, -5)
    fmt.Println(t)
```

Once a slice has no place left for more elements, you should add new elements to it using `append()`.

```markup
    // A 2D slice
    // You can have as many dimensions as needed
    twoD := [][]int{{1, 2, 3}, {4, 5, 6}}
    // Visiting all elements of a 2D slice
    // with a double for loop
    for _, i := range twoD {
            for _, k := range i {
                fmt.Print(k, " ")
            }
            fmt.Println()
    }
```

The previous code shows how to create a 2D slice variable named `twoD` and initialize it at the same time.

```markup
    make2D := make([][]int, 2)
    fmt.Println(make2D)
    make2D[0] = []int{1, 2, 3, 4}
    make2D[1] = []int{-1, -2, -3, -4}
    fmt.Println(make2D)
}
```

The previous part shows how to create a 2D slice with `make()`. What makes the `make2D` a 2D slice is the use of `[][]int` in `make()`.

Running `goSlices.go` produces the next output:

```markup
$ go run goSlices.go 
[] 0 0
[1234.56 -34] with length 2
[-1 -2 -3 -4 -5]
1 2 3 
4 5 6 
[[] []]
[[1 2 3 4] [-1 -2 -3 -4]]
```

### About slice length and capacity

Both arrays and slices support the `len()` function for finding out their length. However, slices also have an additional property called **capacity** that can be found using the `cap()` function.

The capacity of a slice is really important when you want to select a part of a slice or when you want to reference an array using a slice. Both subjects will be discussed over the next few sections.

The capacity shows how much a slice can be expanded without the need to allocate more memory and change the underlying array. Although after slice creation the capacity of a slice is handled by Go, a developer can define the capacity of a slice at creation time using the `make()` function—after that the capacity of the slice doubles each time the length of the slice is about to become bigger than its current capacity. The first argument of `make()` is the type of the slice and its dimensions, the second is its initial length and the third, which is optional, is the capacity of the slice. Although the data type of a slice cannot change after creation, the other two properties can change.

Writing something like `make([]int, 3, 2)` generates an error message because at any given time the capacity of a slice (`2`) cannot be smaller than its length (`3`).

But what happens when you want to append a slice or an array to an existing slice? Should you do that element by element? Go supports the `...` operator, which is used for exploding a slice or an array into multiple arguments before appending it to an existing slice.

The figure that follows illustrates with a graphical representation how length and capacity work in slices.

![A picture containing diagram
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_02_01.png)

Figure 2.1: How slice length and capacity are related

For those of you that prefer code, here is a small Go program that showcases the length and capacity properties of slices. Type it and save it as `capLen.go`.

```markup
package main
import "fmt"
func main() {
    // Only length is defined. Capacity = length
    a := make([]int, 4)
```

In this case, the capacity of `a` is the same as its length, which is `4`.

```markup
    fmt.Println("L:", len(a), "C:", cap(a))
    // Initialize slice. Capacity = length
    b := []int{0, 1, 2, 3, 4}
    fmt.Println("L:", len(b), "C:", cap(b))
```

Once again, the capacity of slice `b` is the same as its length, which is `5`.

```markup
    // Same length and capacity
    aSlice := make([]int, 4, 4)
    fmt.Println(aSlice)
```

This time the capacity of slice `aSlice` is the same as its length, not because Go decided to do so but because we specified it.

```markup
    // Add an element
    aSlice = append(aSlice, 5)
```

When you add a new element to slice `aSlice`, its capacity is doubled and becomes `8`.

```markup
    fmt.Println(aSlice)
    // The capacity is doubled
    fmt.Println("L:", len(aSlice), "C:", cap(aSlice))
    // Now add four elements
    aSlice = append(aSlice, []int{-1, -2, -3, -4}...)
```

The `...` operator expands `[]int{-1, -2, -3, -4}` into multiple arguments and `append()` appends each argument one by one to `aSlice`.

```markup
    fmt.Println(aSlice)
    // The capacity is doubled
    fmt.Println("L:", len(aSlice), "C:", cap(aSlice))
}
```

Running `capLen.go` produces the next output:

```markup
$ go run capLen.go 
L: 4 C: 4
L: 5 C: 5
[0 0 0 0]
[0 0 0 0 5]
L: 5 C: 8
[0 0 0 0 5 -1 -2 -3 -4]
L: 9 C: 16
```

Setting the correct capacity of a slice, if known in advance, will make your programs faster because Go will not have to allocate a new underlying array and have all the data copied over.

Working with slices is good but what happens when you want to work with a continuous part of an existing slice? Is there a practical way to select a part of a slice? Fortunately, the answer is yes—the next subsection sheds some light on selecting a _continuous part_ of a slice.

### Selecting a part of a slice

Go allows you to select parts of a slice, provided that all desired elements are next to each other. This can be pretty handy when you select a range of elements and you do not want to give their indexes one by one. In Go you select a part of a slice by defining two indexes, the first one is the beginning of the selection whereas the second one is the end of the selection, without including the element at that index, separated by `:`.

If you want to process all the command-line arguments of a utility apart from the first one, which is its name, you can assign it to a new variable (`arguments := os.Args`) for ease of use and use the `arguments[1:]` notation to skip the first command-line argument.

However, there is a variation where you can add a third parameter that controls the capacity of the resulting slice. So, using `aSlice[0:2:4]` selects the first 2 elements of a slice (at indexes `0` and `1`) and creates a new slice with a maximum capacity of `4`. The resulting capacity is defined as the result of the `4-0` subtraction where `4` is the maximum capacity and `0` is the first index—if the first index is omitted, it is automatically set to `0`. In this case, the capacity of the result slice will be `4` because `4-0` equals `4`.

If we would have used `aSlice[2:4:4]`, we would have created a new slice with the `aSlice[2]` and `aSlice[3]` elements and with a capacity of `4-2`. Lastly, **the resulting capacity cannot be bigger** than the capacity of the original slice because in that case, you would need a different underlying array.

Type the following code using your favorite editor and save it as `partSlice.go`.

```markup
package main
import "fmt"
func main() {
    aSlice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    fmt.Println(aSlice)
    l := len(aSlice)
    // First 5 elements
    fmt.Println(aSlice[0:5])
    // First 5 elements
    fmt.Println(aSlice[:5])
```

In this first part, we define a new slice named `aSlice` that has 10 elements. Its capacity is the same as its length. Both `0:5` and `:5` notations select the first 5 elements of the slice, which are the elements found at indexes `0`, `1`, `2`, `3`, and `4`.

```markup
    // Last 2 elements
    fmt.Println(aSlice[l-2 : l])
    // Last 2 elements
    fmt.Println(aSlice[l-2:])
```

Given the length of the slice (`l`), we can select the last two elements of the slice either as `l-2 : l` or as `l-2:`.

```markup
    // First 5 elements
    t := aSlice[0:5:10]
    fmt.Println(len(t), cap(t))
    // Elements at indexes 2,3,4
    // Capacity will be 10-2
    t = aSlice[2:5:10]
    fmt.Println(len(t), cap(t))
```

Initially, the capacity of `t` will be `10-0`, which is `10`. In the second case, the capacity of `t` will be `10-2`.

```markup
    // Elements at indexes 0,1,2,3,4
    // New capacity will be 6-0
    t = aSlice[:5:6]
    fmt.Println(len(t), cap(t))
}
```

The capacity of `t` is now `6-0` and its length is going to be `5` because we have selected the first 5 elements of slice `aSlice`.

Running `partSlice.go` generates the next output:

```markup
$ go run partSlice.go 
[0 1 2 3 4 5 6 7 8 9]
```

The previous line is the output of `fmt.Println(aSlice)`.

```markup
[0 1 2 3 4]
[0 1 2 3 4]
```

The previous two lines are generated from `fmt.Println(aSlice[0:5])` and `fmt.Println(aSlice[:5])`.

```markup
[8 9]
[8 9]
```

Analogously, the previous two lines are generated from `fmt.Println(aSlice[l-2 : l])` and `fmt.Println(aSlice[l-2:])`.

```markup
5 10
3 8
5 6
```

The last three lines print the length and the capacity of `aSlice[0:5:10]`, `aSlice[2:5:10]` and `aSlice[:5:6]`.

### Byte slices

A **byte slice** is a slice of the `byte` data type (`[]byte`). Go knows that most `byte` slices are used to store strings and so makes it easy to switch between this type and the `string` type. There is nothing special in the way you can access a byte slice compared to the other types of slices. What is special is that Go uses _byte slices_ for performing file I/O operations because they allow you to determine with precision the amount of data you want to read or write to a file. This happens because bytes are a universal unit among computer systems.

As Go does not have a `char` data type, it uses `byte` and `rune` for storing character values. A single `byte` can only store a single ASCII character whereas a `rune` can store Unicode characters. However, a rune can occupy multiple bytes.

The small program that follows illustrates how you can convert a `byte` slice into a `string` and vice versa, which you need for most File I/O operations—type it and save it as `byteSlices.go`.

```markup
package main
import "fmt"
func main() {
    // Byte slice
    b := make([]byte, 12)
    fmt.Println("Byte slice:", b)
```

An empty byte slice contains zeros—in this case, 12 zeros.

```markup
    b = []byte("Byte slice €")
    fmt.Println("Byte slice:", b)
```

In this case, the size of `b` is the size of the string `"Byte slice €"`, without the double quotes—`b` now points to a different memory location than before, which is where `"Byte slice €"` is stored. This is how you convert a `string` into a `byte` slice.

As Unicode characters like € need more than one byte for their representation, the length of the `byte` slice might not be the same as the length of the string that it stores.

```markup
    // Print byte slice contents as text
    fmt.Printf("Byte slice as text: %s\n", b)
    fmt.Println("Byte slice as text:", string(b))
```

The previous code shows how to print the contents of a byte slice as text using two techniques. The first one is by using the `%s` control string and the second one using `string()`.

```markup
    // Length of b
    fmt.Println("Length of b:", len(b))
}
```

The previous code prints the real length of the byte slice.

Running `byteSlices.go` produces the next output:

```markup
$ go run byteSlices.go 
Byte slice: [0 0 0 0 0 0 0 0 0 0 0 0]
Byte slice: [66 121 116 101 32 115 108 105 99 101 32 226 130 172]
Byte slice as text: Byte slice €
Byte slice as text: Byte slice €
Length of b: 14
```

The last line of the output proves that although the `b` byte slice has 12 characters, it has a size of `14`.

### Deleting an element from a slice

There is no default function for deleting an element from a slice, which means that if you need to delete an element from a slice, you must write your own code. Deleting an element from a slice can be tricky, so this subsection presents two techniques for doing so. The first technique virtually divides the original slice into two slices, split at the index of the element that needs to be deleted. Neither of the two slices includes the element that is going to be deleted. After that, we concatenate these two slices and creates a new one. The second technique copies the last element at the place of the element that is going to be deleted and creates a new slice by excluding the last element from the original slice.

The next figure shows a graphical representation of the two techniques for deleting an element from a slice.

![A picture containing text, sign
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_02_02.png)

Figure 2.2: Deleting an element from a slice

The following program shows the two techniques that can be used for deleting an element from a slice. Create a text file by typing the following code—save it as `deleteSlice.go`.

```markup
package main
import (
    "fmt"
    "os"
    "strconv"
)
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Need an integer value.")
        return
    }
    index := arguments[1]
    i, err := strconv.Atoi(index)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Using index", i)
    aSlice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
    fmt.Println("Original slice:", aSlice)
    // Delete element at index i
    if i > len(aSlice)-1 {
        fmt.Println("Cannot delete element", i)
        return
    }
    // The ... operator auto expands aSlice[i+1:] so that
    // its elements can be appended to aSlice[:i] one by one
    aSlice = append(aSlice[:i], aSlice[i+1:]...)
    fmt.Println("After 1st deletion:", aSlice)
```

Here we logically divide the original slice into two slices. The two slices are split at the index of the element that needs to be deleted. After that, we concatenate these two slices with the help of `...`. Next, we see the second technique in action.

```markup
    // Delete element at index i
    if i > len(aSlice)-1 {
        fmt.Println("Cannot delete element", i)
        return
    }
    // Replace element at index i with last element
    aSlice[i] = aSlice[len(aSlice)-1]
    // Remove last element
    aSlice = aSlice[:len(aSlice)-1]
    fmt.Println("After 2nd deletion:", aSlice)
}
```

We replace the element that we want to delete with the last element using the `aSlice[i] = aSlice[len(aSlice)-1]` statement and then we remove the last element with the `aSlice = aSlice[:len(aSlice)-1]` statement.

Running `deleteSlice.go` produces the following kind of output, depending on your input:

```markup
$ go run deleteSlice.go 1
Using index 1
Original slice: [0 1 2 3 4 5 6 7 8]
After 1st deletion: [0 2 3 4 5 6 7 8]
After 2nd deletion: [0 8 3 4 5 6 7]
```

As the slice has 9 elements, you can delete the element at index value `1`.

```markup
$ go run deleteSlice.go 10
Using index 10
Original slice: [0 1 2 3 4 5 6 7 8]
Cannot delete element 10
```

As the slice has only 9 elements, you cannot delete an element with an index value of `10` from the slice.

### How slices are connected to arrays

As mentioned before, behind the scenes, each slice is implemented using an **underlying array**. The length of the underlying array is the same as the capacity of the slice and there exist pointers that connect the slice elements to the appropriate array elements.

You can understand that by connecting an existing array with a slice, Go allows you to reference an array or a part of an array using a slice. This has some strange capabilities including the fact that the changes to the slice affect the referenced array! However, when the capacity of the slice changes, the connection to the array ceases to exist! This happens because when the capacity of a slice changes, so does the underlying array, and the connection between the slice and the original array does not exist anymore.

Type the following code and save it as `sliceArrays.go`.

```markup
package main
import (
    "fmt"
)
func change(s []string) {
    s[0] = "Change_function"
}
```

This is a function that changes the first element of a slice.

```markup
func main() {
    a := [4]string{"Zero", "One", "Two", "Three"}
    fmt.Println("a:", a)
```

Here we define an array named `a` with 4 elements.

```markup
    var S0 = a[0:1]
    fmt.Println(S0)
    S0[0] = "S0"
```

Here we connect `S0` with the first element of the array `a` and we print it. Then we change the value of `S0[0]`.

```markup
    var S12 = a[1:3]
    fmt.Println(S12)
    S12[0] = "S12_0"
    S12[1] = "S12_1"
```

In this part, we associate `S12` with `a[1]` and `a[2]`. Therefore `S12[0]` = `a[1]` and `S12[1]` = `a[2]`. Then, we change the values of both `S12[0]` and `S12[1]`. These two changes will also change the contents of `a`. Put simply, `a[1]` takes the new value of `S12[0]` and `a[2]` takes the new value of `S12[1]`.

```markup
    fmt.Println("a:", a)
```

And we print variable `a`, which has not changed at all in a direct way. However, due to the connections of `a` with `S0` and `S12`, the contents of `a` have changed!

```markup
    // Changes to slice -> changes to array
    change(S12)
    fmt.Println("a:", a)
```

As the slice and the array are connected, any changes you make to the slice will also affect the array even if the changes take place inside a function.

```markup
    // capacity of S0
    fmt.Println("Capacity of S0:", cap(S0), "Length of S0:", len(S0))
    // Adding 4 elements to S0
    S0 = append(S0, "N1")
    S0 = append(S0, "N2")
    S0 = append(S0, "N3")
    a[0] = "-N1"
```

As the capacity of `S0` changes, it is no longer connected to the same underlying array (`a`).

```markup
    // Changing the capacity of S0
    // Not the same underlying array anymore!
    S0 = append(S0, "N4")
    fmt.Println("Capacity of S0:", cap(S0), "Length of S0:", len(S0))
    // This change does not go to S0
    a[0] = "-N1-"
    // This change does not go to S12
    a[1] = "-N2-"
```

However, array `a` and slice `S12` are still connected because the capacity of `S12` has not changed.

```markup
    fmt.Println("S0:", S0)
    fmt.Println("a: ", a)
    fmt.Println("S12:", S12)
}
```

Lastly, we print the final versions of `a`, `S0`, and `S12`.

Running `sliceArrays.go` produces the following output:

```markup
$ go run sliceArrays.go 
a: [Zero One Two Three]
[Zero]
[One Two]
a: [S0 S12_0 S12_1 Three]
a: [S0 Change_function S12_1 Three]
Capacity of S0: 4 Length of S0: 1
Capacity of S0: 8 Length of S0: 5
S0: [-N1 N1 N2 N3 N4]
a:  [-N1- -N2- N2 N3]
S12: [-N2- N2]
```

Let us now discuss the use of the `copy()` function in the next subsection.

### The copy() function

Go offers the `copy()` function for copying an existing array to a slice or an existing slice to another slice. However, the use of `copy()` can be tricky because the destination slice is not auto-expanded if the source slice is bigger than the destination slice. Additionally, if the destination slice is bigger than the source slice, then `copy()` does not empty the elements from the destination slice that did not get copied. This is better illustrated in the figure that follows.

![A picture containing text, sign
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_02_03.png)

Figure 2.3: The use of the copy() function

The following program illustrates the use of `copy()`—type it in your favorite text editor and save it as `copySlice.go`.

```markup
package main
import "fmt"
func main() {
    a1 := []int{1}
    a2 := []int{-1, -2}
    a5 := []int{10, 11, 12, 13, 14}
    fmt.Println("a1", a1)
    fmt.Println("a2", a2)
    fmt.Println("a5", a5)
    // copy(destination, input)
    // len(a2) > len(a1)
    copy(a1, a2)
    fmt.Println("a1", a1)
    fmt.Println("a2", a2)
```

Here we run the `copy(a1, a2)` command. In this case, the `a2` slice is bigger than `a1`. After `copy(a1, a2)`, `a2` remains the same, which makes perfect sense as `a2` is the input slice, whereas the first element of `a2` is copied to the first element of `a1` because `a1` has space for a single element only.

```markup
    // len(a5) > len(a1)
    copy(a1, a5)
    fmt.Println("a1", a1)
    fmt.Println("a5", a5)
```

In this case, `a5` is bigger than `a1`. Once again, after `copy(a1, a5)`, `a5` remains the same whereas `a5[0]` is copied to `a1[0]`.

```markup
    // len(a2) < len(a5) -> OK
    copy(a5, a2)
    fmt.Println("a2", a2)
    fmt.Println("a5", a5)
}
```

In this last case, `a2` is shorter than `a5`. This means that the entire `a2` is copied into `a5`. As the length of `a2` is 2, only the first 2 elements of `a5` change.

Running `copySlice.go` produces the next output:

```markup
$ go run copySlice.go 
a1 [1]
a2 [-1 -2]
a5 [10 11 12 13 14]
a1 [-1]
a2 [-1 -2]
```

The `copy(a1, a2)` statement does not alter the `a2` slice, just `a1`. As the size of `a1` is `1`, only the first element from `a2` is copied.

```markup
a1 [10]
a5 [10 11 12 13 14]
```

Similarly, `copy(a1, a5)` alters `a1` only. As the size of `a1` is `1`, only the first element from `a5` is copied to `a1`.

```markup
a2 [-1 -2]
a5 [-1 -2 12 13 14]
```

Last, `copy(a5, a2)` alters `a5` only. As the size of `a5` is `5`, only the first two elements from `a5` are altered and become the same as the first two elements of `a2`, which has a size of `2`.

### Sorting slices

There are times when you want to present your information sorted and you want Go to do the job for you. In this subsection, we'll see how to sort slices of various standard data types using the functionality offered by the `sort` package.

The `sort` package can sort slices of built-in data types without the need to write any extra code. Additionally, Go provides the `sort.Reverse()` function for sorting in the reverse order than the default. However, what is really interesting is that `sort` allows you to write your own sorting functions for custom data types by implementing the `sort.Interface` interface—you will learn more about the `sort.Interface` interface and interfaces in general in _Chapter 4,_ _Reflection and Interfaces_.

So, you can sort a slice of integers saved as `sInts` by typing `sort.Ints(sInts)`. When sorting a slice of integers in reverse order using `sort.Reverse()`, you need to pass the desired slice to `sort.Reverse()` using `sort.IntSlice(sInts)` because the `IntSlice` type implements the `sort.Interface` internally, which allows you to sort in a different way than usual. The same applies to the other standard Go data types.

Create a text file with the code that illustrates the use of `sort` and name it `sortSlice.go`.

```markup
package main
import (
    "fmt"
    "sort"
)
func main() {
    sInts := []int{1, 0, 2, -3, 4, -20}
    sFloats := []float64{1.0, 0.2, 0.22, -3, 4.1, -0.1}
    sStrings := []string{"aa", "a", "A", "Aa", "aab", "AAa"}
    fmt.Println("sInts original:", sInts)
    sort.Ints(sInts)
    fmt.Println("sInts:", sInts)
    sort.Sort(sort.Reverse(sort.IntSlice(sInts)))
    fmt.Println("Reverse:", sInts)
```

As `sort.Interface` knows how to sort integers, it is trivial to sort them in reverse order. Sorting in reverse order is as simple as calling the `sort.Reverse()` function.

```markup
    fmt.Println("sFloats original:", sFloats)
    sort.Float64s(sFloats)
    fmt.Println("sFloats:", sFloats)
    sort.Sort(sort.Reverse(sort.Float64Slice(sFloats)))
    fmt.Println("Reverse:", sFloats)
    fmt.Println("sStrings original:", sStrings)
    sort.Strings(sStrings)
    fmt.Println("sStrings:", sStrings)
    sort.Sort(sort.Reverse(sort.StringSlice(sStrings)))
    fmt.Println("Reverse:", sStrings)
}
```

The same rules apply when sorting floating point numbers and strings.

Running `sortSlice.go` produces the next output:

```markup
$ go run sortSlice.go
sInts original: [1 0 2 -3 4 -20]
sInts: [-20 -3 0 1 2 4]
Reverse: [4 2 1 0 -3 -20]
sFloats original: [1 0.2 0.22 -3 4.1 -0.1]
sFloats: [-3 -0.1 0.2 0.22 1 4.1]
Reverse: [4.1 1 0.22 0.2 -0.1 -3]
sStrings original: [aa a A Aa aab AAa]
sStrings: [A AAa Aa a aa aab]
Reverse: [aab aa a Aa AAa A]
```

The output illustrates how the original slices were sorted in both normal and reverse order.

Bookmark

# Pointers

Go has support for pointers but not for pointer arithmetic, which is the cause of many bugs and errors in programming languages like C. A **pointer** is the memory address of a variable. You need to **dereference** a pointer in order to get its value—dereferencing is performed using the `*` character in front of the pointer variable. Additionally, you can get the memory address of a normal variable using an `&` in front of it.

The next diagram shows the difference between a pointer to an `int` and an `int` variable.

![Graphical user interface
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_02_04.png)

Figure 2.4: An int variable and a pointer to an int

If a pointer variable points to an existing regular variable, then any changes you make to the stored value using the pointer variable will modify the regular variable.

The format and the values of memory addresses might be different between different machines, different operating systems, and different architectures.

You might ask, what is the point of using pointers since there is no support for pointer arithmetic. The main benefit you get from pointers is that passing a variable to a function as a pointer (we can call that _by reference_) does not discard any changes you make to the value of that variable inside that function when the function returns. There exist times where you want that functionality because it simplifies your code, but the price you pay for that simplicity is being extra careful with what you do with a pointer variable. Remember that slices are passed to functions without the need to use a pointer—it is Go that passes the pointer to the underlying array of a slice and there is no way to change that behavior.

Apart from reasons of simplicity, there exist three more reasons for using pointers:

-   Pointers allow you to share data between functions. However, when sharing data between functions and goroutines, you should be extra careful with race condition issues.
-   Pointers are also very handy when you want to tell the difference between the zero value of a variable and a value that is not set (`nil`). This is particularly useful with structures because pointers (and therefore **pointers to structures**, which are fully covered in the next chapter), can have the `nil` value, which means that you can compare a pointer to a structure with the `nil` value, which is not allowed for normal structure variables.
-   Having support for pointers and, more specifically, pointers to structures allows Go to support data structures such as linked lists and binary trees, which are widely used in computer science. Therefore, you are allowed to define a structure field of a `Node` structure as `Next *Node`, which is a pointer to another `Node` structure. Without pointers, this would have been difficult to implement and may be too slow.

The following code illustrates how you can use pointers in Go—create a text file named `pointers.go` and type the presented code.

```markup
package main
import "fmt"
type aStructure struct {
    field1 complex128
    field2 int
}
```

This is a structure with two fields named `field1` and `field2`.

```markup
func processPointer(x *float64) {
    *x = *x * *x
}
```

This is a function that gets a pointer to a `float64` variable as input. As we are using a pointer, all changes to the function parameter inside the function are persistent. Additionally, there is no need to return something.

```markup
func returnPointer(x float64) *float64 {
    temp := 2 * x
    return &temp
}
```

This is a function that requires a `float64` parameter as input and returns a pointer to a `float64`. In order to return the memory address of a regular variable, you need to use `&` (`&temp`).

```markup
func bothPointers(x *float64) *float64 {
    temp := 2 * *x
    return &temp
}
```

This is a function that requires a pointer to a `float64` as input and returns a pointer to a `float64` as output. The `*x` notation is used for getting the value stored in the memory address stored in `x`.

```markup
func main() {
    var f float64 = 12.123
    fmt.Println("Memory address of f:", &f)
```

To get the memory address of a regular variable named `f`, you should use the `&f` notation.

```markup
    // Pointer to f
    fP := &f
    fmt.Println("Memory address of f:", fP)
    fmt.Println("Value of f:", *fP)
    // The value of f changes
    processPointer(fP)
    fmt.Printf("Value of f: %.2f\n", f)
```

`fP` is now a pointer to the memory address of the `f` variable. Any changes to the value stored in the `fP` memory address have an effect on the `f` value as well. However, this is only true for as long as `fP` points to the memory address of the `f` variable.

```markup
    // The value of f does not change
    x := returnPointer(f)
    fmt.Printf("Value of x: %.2f\n", *x)
```

The value of `f` does not change because the function only uses its value.

```markup
    // The value of f does not change
    xx := bothPointers(fP)
    fmt.Printf("Value of xx: %.2f\n", *xx)
```

In this case, the value of `f`, as well as the value stored in the `fP` memory address, does not change because the `bothPointers()` function does not make any changes to the value stored in the `fP` memory address.

```markup
    // Check for empty structure
    var k *aStructure
```

The `k` variable is a pointer to an `aStructure` structure. As `k` points to nowhere, Go makes it point to `nil`, which is the **zero value for pointers**.

```markup
    // This is nil because currently k points to nowhere
    fmt.Println(k)
    // Therefore you are allowed to do this:
    if k == nil {
        k = new(aStructure)
    }
```

As `k` is `nil`, we are allowed to assign it to an empty `aStructure` value with `new(aStructure)` without losing any data. Now, `k` is no longer `nil` but both fields of `aStructure` have the zero values of their data types.

```markup
    fmt.Printf("%+v\n", k)
    if k != nil {
        fmt.Println("k is not nil!")
    }
}
```

Just make sure that `k` is not `nil`—you might consider that check redundant, but it does not hurt to double-check.

Running `pointers.go` generates the following kind of output:

```markup
Memory address of f: 0xc000014090
Memory address of f: 0xc000014090
Value of f: 12.123
Value of f: 146.97
Value of x: 293.93
Value of xx: 293.93
<nil>
&{field1:(0+0i) field2:0}
k is not nil!
```

We revisit pointers in the next chapter where we discuss structures. Next, we discuss generating random numbers and random strings.

Bookmark

# Generating random numbers

Random number generation is an art as well as a research area in computer science. This is because computers are purely logical machines, and it turns out that using them to generate random numbers is extremely difficult! Go can help you with that using the functionality of the `math/rand` package. Each random number generator needs a **seed** to start producing numbers. The seed is used for initializing the entire process and is extremely important because if you always start with the same seed, you will always get the same sequence of pseudo-random numbers. This means that everybody can regenerate that sequence, and that particular sequence will not be random after all. However, this feature is really useful for testing purposes. In Go, the `rand.Seed()` function is used for initializing a random number generator.

If you are really interested in random number generation, you should start by reading the second volume of _The Art of Computer Programming_ by Donald E. Knuth (Addison-Wesley Professional, 2011).

The following function, which is part of `randomNumbers.go` found in `ch02` in the book's GitHub repository, is what generates random numbers in the `[min, max)` range.

```markup
func random(min, max int) int {
    return rand.Intn(max-min) + min
}
```

The `random()` function does all of the work, which is generating pseudo-random numbers in a given range from `min` to `max-1` by calling `rand.Intn()`. `rand.Intn()` generates non-negative random integers from `0` up to the value of its single parameter minus `1`.

The `randomNumbers.go` utility accepts four command-line parameters but can also work with fewer parameters by using default values. By default, `randomNumbers.go` produces 100 random integers from `0` up to and including `99`.

```markup
$ go run randomNumbers.go 
Using default values!
39 75 78 89 39 28 37 96 93 42 60 69 50 9 69 27 22 63 4 68 56 23 54 14 93 61 19 13 83 72 87 29 4 45 75 53 41 76 84 51 62 68 37 11 83 20 63 58 12 50 8 31 14 87 13 97 17 60 51 56 21 68 32 41 79 13 79 59 95 56 24 83 53 62 97 88 67 59 49 65 79 10 51 73 48 58 48 27 30 88 19 16 16 11 35 45 72 51 41 28
```

In the next output, we define each of the parameters manually (the last parameter of the utility is the seed value):

```markup
$ go run randomNumbers.go 1 5 10 10
3 1 4 4 1 1 4 4 4 3
$ go run randomNumbers.go 1 5 10 10
3 1 4 4 1 1 4 4 4 3
$ go run randomNumbers.go 1 5 10 11
1 4 2 1 3 2 2 4 1 3
```

The first two times the seed value was `10`, so we got the same output. The third time the value of the seed was `11`, which generated a different output.

## Generating random strings

ImagineDevOps  that you want to generate random strings that can be used as difficult to guess passwords or for testing purposes. Based on random number generation, we create a utility that produces random strings. The utility is implemented as `genPass.go` and can be found in the `ch02` directory of the book's GitHub repository. The core functionality of `genPass.go` is found in the next function.

```markup
func getString(len int64) string {
    temp := ""
    startChar := "!"
    var i int64 = 1
    for {
        myRand := random(MIN, MAX)
        newChar := string(startChar[0] + byte(myRand))
        temp = temp + newChar
        if i == len {
            break
        }
        i++
    }
    return temp
}
```

As we only want to get printable ASCII characters, we limit the range of pseudo-random numbers that can be generated. The total number of printable characters in the ASCII table is 94. This means that the range of the pseudo-random numbers that the program can generate should be from 0 to 94, without including 94. Therefore, the values of the `MIN` and `MAX` global variables, which are not shown here, are `0` and `94`, respectively.

The `startChar` variable holds the first ASCII character that can be generated by the utility, which, in this case, is the exclamation mark, which has a decimal ASCII value of `33`. Given that the program can generate pseudo-random numbers up to `94`, the maximum ASCII value that can be generated is `93 + 33`, which is equal to `126`, which is the ASCII value of `~`. All generated characters are kept in the `temp` variable, which is returned once the `for` loop exits. The `string(startChar[0] + byte(myRand))` statement converts the random integers into characters in the desired range.

The `genPass.go` utility accepts a single parameter, which is the length of the generated password. If no parameter is given, `genPass.go` produces a password with 8 characters, which is the default value of the `LENGTH` variable.

Running `genPass.go` produces the following kind of output:

```markup
$ go run genPass.go
Using default values...
!QrNq@;R
$ go run genPass.go 20
sZL>{F~"hQqY>r_>TX?O
```

The first program execution uses the default value for the length of the generated string whereas the second program execution creates a random string with 20 characters.

## Generating secure random numbers

If you intend to use these pseudo-random numbers for security-related work, it is important that you use the `crypto/rand` package, which implements a cryptographically secure pseudo-random number generator. You do not need to define a seed when using the `crypto/rand` package.

The following function that is part of the `cryptoRand.go` source code showcases how secure random numbers are generated with the functionality of `crypto/rand`.

```markup
func generateBytes(n int64) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return nil, err
    }
    return b, nil
}
```

The `rand.Read()` function randomly generates numbers that occupy the entire `b` byte slice. You need to decode that byte slice using `base64.URLEncoding.EncodeToString(b)` in order to get a valid string without any control or unprintable characters. This conversion takes place in the `generatePass()` function, which is not shown here.

Running `cryptoRand.go` creates the following kind of output:

```markup
$ go run cryptoRand.go   
Using default values!
Ce30g--D
$ go run cryptoRand.go 20
AEIePSYb13KwkDnO5Xk_
```

The output is not different from the one generated by `genPass.go`, it is just that the random numbers are generated more securely, which means that they can be used in applications where security is important.

Now that we know how to generate random numbers and random strings, we are going to revisit the phone book application and use these techniques to populate the phone book with random data.

Bookmark

# Updating the phone book application

In this last section of the book, we are going to create a function that populates the phone book application from the previous chapter with random data, which is pretty handy when you want to put lots of data in an application for testing purposes.

I have used this handy technique in the past in order to put sample data on Kafka topics.

The biggest change in this version of the phone book application is that the searching is based on the telephone number because it is easier to search random numbers instead of random strings. But this is a small code change in the `search()` function—this time `search()` uses `v.Tel == key` instead of `v.Surname == key` in order to try to match the `Tel` field.

The `populate()` function of `phoneBook.go` (as found in the `ch02` directory) does all the work—the implementation of `populate()` is the following.

```markup
func populate(n int, s []Entry) {
    for i := 0; i < n; i++ {
        name := getString(4)
        surname := getString(5)
        n := strconv.Itoa(random(100, 199))
        data = append(data, Entry{name, surname, n})
    }
}
```

The `getString()` function generates letters from `A` to `Z` and nothing else in order to make the generated strings more readable. There is no point in using special characters in names and surnames. The generated telephone numbers are in the 100 to 198 range, which is implemented using a call to `random(100, 199)`. The reason for this is that it is easier to search for a three-digit number. Feel free to experiment with the generated names, surnames, and telephone numbers.

Working with `phoneBook.go` generates the following kind of output:

```markup
$ go run phoneBook.go search 123  
Data has 100 entries.
{BHVA QEEQL 123}
$ go run phoneBook.go search 1234
Data has 100 entries.
Entry not found: 1234
$ go run phoneBook.go list
Data has 100 entries.
{DGTB GNQKI 169}
{BQNU ZUQFP 120}
...
```

Although these randomly generated names and surnames are not perfect, they are more than enough for testing purposes. In the next chapter, we'll learn how to work with CSV data.

Bookmark

# Exercises

-   Create a function that concatenates two arrays into a new slice.
-   Create a function that concatenates two arrays into a new array.
-   Create a function that concatenates two slices into a new array.

Bookmark

# Summary

In this chapter, we learned about the basic data types of Go, including numerical data types, strings, and errors. Additionally, we learned how to group similar values using arrays and slices. Lastly, we learned about the differences between arrays and slices and why slices are more versatile than arrays, as well as pointers and generating random numbers and strings in order to provide random data to the phone book application.

The next chapter discusses a couple of more complex composite data types of Go, _maps_ and _structures_. Maps can use keys of different data types whereas structures can group multiple data types and create new ones that you can access as single entities. As you will see in later chapters, structures play a key role in Go.

Bookmark

# Additional resources

-   The `sort` package documentation: [https://golang.org/pkg/sort/](https://golang.org/pkg/sort/)
-   The `time` package documentation: [https://golang.org/pkg/time/](https://golang.org/pkg/time/)
-   The `crypto/rand` package documentation: [https://golang.org/pkg/crypto/rand/](https://golang.org/pkg/crypto/rand/)
-   The `math/rand` package documentation: [https://golang.org/pkg/math/rand/](https://golang.org/pkg/math/rand/)