# Composite Data Types

Go offers support for maps and structures, which are composite data types and the main subject of this chapter. The reason that we present them separately from arrays and slices is that both maps and structures are more flexible and powerful than arrays and slices. The general idea is that if an array or a slice cannot do the job, you might need to look at maps. If a map cannot help you, then you should consider creating and using a structure.

You have already seen structures in _Chapter 1_, _A Quick Introduction to Go_, where we created the initial version of the phone book application. However, in this chapter we are going to learn more about structures as well as maps. This knowledge will allow us to read and save data in CSV format using structures and **create an index** for quickly searching a slice of structures, based on a given key, by using a map.

Last, we are going to apply some of these Go features to improve the phone book application we originally developed in _Chapter 1_, _A Quick Introduction to Go_. The new version of the phone book application loads and saves its data from disk, which means that it is no longer needed to hardcode your data.

This chapter covers:

-   Maps
-   Structures
-   Pointers and structures
-   Regular expressions and pattern matching
-   Improving the phone book application

Maps can use keys of different data types whereas structures can group multiple data types and create new ones. So, without further ado, let us begin by presenting maps.

Just Imagine

# Maps

Both arrays and slices limit you to using positive integers as indexes. Maps are powerful data structures because they allow you to use indexes of various data types as keys to look up your data as long as these keys are **comparable**. A practical rule of thumb is that you should use a map when you are going to need indexes that are not positive integer numbers or when the integer indexes have big gaps.

Although `bool` variables are comparable, it makes no sense to use a `bool` variable as the key to a Go map because it only allows for two distinct values. Additionally, although floating point values are comparable, precision issues caused by the internal representation of such values might create bugs and crashes, so you might want to avoid using floating point values as keys to Go maps.

You might ask, why do we need maps and what are their advantages? The following list will help clarify things:

-   Maps are very versatile. Later in this chapter we will create a _database index_ using a map, which allows us to search and access slice elements based on a given key or, in more advanced situations, a combination of keys.
-   Although this is not always the case, working with maps in Go is fast, as you can access all elements of a map in **linear time**. Inserting and retrieving elements from a map is fast and does not depend on the cardinality of the map.
-   Maps are easy to understand, which leads to clear designs.

You can create a new `map` variable using either `make()` or a map literal. Creating a new map with `string` keys and `int` values using `make()` is as simple as writing `make(map[string]int)` and assigning its return value to a variable. On the other hand, if you decide to create a map using a map literal, you need to write something like the following:

```markup
m := map[string]int {
    "key1": -1
    "key2": 123
}
```

The map literal version is faster when you want to add data to a map at the time of creation.

You should make no assumptions about the order of the elements inside a map. Go randomizes keys when iterating over a map—this is done on purpose and is an intentional part of the language design.

You can find the length of a map, which is the number of keys in the map, using the `len()` function, which also works with arrays and slices; and you can delete a key and value pair from a map using the `delete()` function, which accepts two arguments: the name of the map and the name of the key, in that order.

Additionally, you can tell whether a key `k` exists on a map named `aMap` by the second return value of the `v, ok := aMap[k]` statement. If `ok` is set to `true`, then `k` exists, and its value is `v`. If it does not exist, `v` will be set to the zero value of its data type, which depends on the definition of the map. If you try to get the value of a key that does not exist in a map, Go will not complain about it and returns the zero value of the data type of the value.

Now, let us discuss a special case where a map variable has the `nil` value.

## Storing to a nil map

You are allowed to assign a map variable to `nil`. In that case, you will not be able to use that variable until you assign it to a new map variable. Put simply, if you try to store data on a `nil` map, your program will crash. This is illustrated in the next bit of code, which is the implementation of the `main()` function of the `nilMap.go` source file that can be found in the `ch03` directory of the GitHub repository of this book.

```markup
func main() {
    aMap := map[string]int{}
    aMap["test"] = 1
```

This works because `aMap` points somewhere, which is the return value of `map[string]int{}`.

```markup
    fmt.Println("aMap:", aMap)
    aMap = nil
```

At this point `aMap` points to `nil`, which is a synonym for nothing.

```markup
    fmt.Println("aMap:", aMap)
    if aMap == nil {
        fmt.Println("nil map!")
        aMap = map[string]int{}
    }
```

Testing whether a map points to `nil` before using it is a good practice. In this case, `if aMap == nil` allows us to determine whether we can store a key/pair value to `aMap` or not—we cannot and if we try it, the program will crash. We correct that by issuing the `aMap = map[string]int{}` statement.

```markup
    aMap["test"] = 1
    // This will crash!
    aMap = nil
    aMap["test"] = 1
}
```

In this last part of the program, we illustrate how your program will crash if you try to store on a `nil` map—never use such code in production!

In real-world applications, if a function accepts a map argument, then it should check that the map is not `nil` before working with it.

Running `nilMap.go` produces this output:

```markup
$ go run nilMap.go
aMap: map[test:1]
aMap: map[]
nil map!
panic: assignment to entry in nil map
goroutine 1 [running]:
main.main()
        /Users/mtsouk/Desktop/mGo3rd/code/ch03/nilMap.go:21 +0x225
```

The reason the program crashed is shown in the program output: `panic: assignment to entry in nil map`.

## Iterating over maps

When `for` is combined with the `range` keyword it implements the functionality of `foreach` loops found in other programming languages and allows you to iterate over all the elements of a map without knowing its size or its keys. When `range` is applied on a map, it returns **key and value pairs** in that order.

Type the following code and save it as `forMaps.go`.

```markup
package main
import "fmt"
func main() {
    aMap := make(map[string]string)
    aMap["123"] = "456"
    aMap["key"] = "A value"
    // range works with maps as well
    for key, v := range aMap {
        fmt.Println("key:", key, "value:", v)
    }
```

In this case we use both the key and the value that returned from `range`.

```markup
    for _, v := range aMap {
        fmt.Print(" # ", v)
    }
    fmt.Println()
}
```

In this case, as we are only interested in the values returned by the map, we ignore the keys.

As you already know, you should make no assumptions about the order that the key and value pairs of a map will be returned in from a `for` and `range` loop.

Running `forMaps.go` produces this output:

```markup
$ go run forMaps.go
key: key value: A value
key: 123 value: 456
 # 456 # A value
```

Having covered maps, it is time to learn about Go structures.

Just Imagine

# Structures

Structures in Go are both very powerful and very popular and are used for organizing and grouping various types of data under the same name. Structures are the more versatile data types in Go and they can even be associated with functions, which are called methods.

Structures, as well as other user-defined data types, are usually defined outside the `main()` function or any other package function so that they have a global scope and are available to the entire Go package. Therefore, unless you want to make clear that a type is only useful within the current local scope and is not expected to be used elsewhere, you should write the definitions of new data types outside functions.

## Defining new structures

When you define a new structure, you group a set of values into a single data type, which allows you to pass and receive this set of values as a single entity. A structure has **fields**, and each field has its own data type, which can even be another structure or slice of structures. Additionally, as a structure is a new data type, it is defined using the `type` keyword followed by the name of the structure and ending with the `struct` keyword, which signifies that we are defining a new structure.

The following code defines a new structure named `Entry`:

```markup
type Entry struct {
    Name    string
    Surname string
    Year    int
}
```

The `type` keyword allows you to define new data types or create aliases for existing ones. Therefore, you are allowed to say `type myInt int` and define a new data type called `myInt` that is an alias for `int`. However, Go considers `myInt` and `int` as totally different data types that you cannot compare directly even though they hold the same kind of values. Each structure defines a new data type, hence the use of the `type` keyword.

For reasons that will become evident in _Chapter 5_, _Go Packages and Functions_, the fields of a structure usually begin with an uppercase letter—this depends on what you want to do with the fields. The `Entry` structure has three fields named `Name`, `Surname`, and `Year`. The first two fields are of the `string` data type, whereas the last field holds an `int` value.

These three fields can be accessed with the dot notation as `V.Name`, `V.Surname`, and `V.Year`, where `V` is the name of the variable holding the instance of the `Entry` structure. A **structure literal** named `p1` can be defined as `p1 := aStructure{"fmt", 12, -2}`.

There exist two ways to work with structure variables. The first one is as **regular variables** and the second one is as **pointer variables** that point to the memory address of a structure. Both ways are equally good and are usually embedded into separate functions because they allow you to initialize some or all of the fields of structure variables properly and/or do any other tasks you want before using the structure variable. As a result, there exist two main ways to create a new structure variable using a function. The first one returns a regular structure variable whereas the second one returns a pointer to a structure. Each one of these two ways has two variations. The first variation returns a structure instance that is initialized by the Go compiler, whereas the second variation returns a structure instance that is initialized by the user.

The order in which you put the fields in the definition of a structure type is significant for the **type identity** of the defined structure. Put simply, two structures with the same fields will not be considered identical in Go if their fields are not in the same order.

## Using the new keyword

Additionally, you can create new structure instances using the `new()` keyword: `pS := new(Entry)`. The `new()` keyword has the following properties:

-   It allocates the proper memory space, which depends on the data type, and then it zeroes it
-   It always **returns a pointer** to the allocated memory
-   It works for all data types except _channel_ and _map_

All these techniques are illustrated in the code that follows. Type the following code in your favorite text editor and save it as `structures.go`.

```markup
package main
import "fmt"
type Entry struct {
    Name    string
    Surname string
    Year    int
}
// Initialized by Go
func zeroS() Entry {
    return Entry{}
}
```

Now is a good time to remind you of an important Go rule: **If no initial value is given to a variable, the Go compiler automatically initializes that variable to the zero value of its data type**. For structures, this means that a structure variable without an initial value is initialized to the zero values of each one of the data types of its fields. Therefore, the `zeroS()` function returns a zero-initialized `Entry` structure.

```markup
// Initialized by the user
func initS(N, S string, Y int) Entry {
    if Y < 2000 {
        return Entry{Name: N, Surname: S, Year: 2000}
    }
    return Entry{Name: N, Surname: S, Year: Y}
}
```

In this case the user initializes the new structure variable. Additionally, the `initS()` function checks whether the value of the `Year` field is smaller than `2000` and acts accordingly. If it is smaller than `2000`, then the value of the `Year` field becomes `2000`. This condition is specific to the requirements of the application you are developing—what this shows is that the place where you initialize a structure is good for checking your input.

```markup
// Initialized by Go - returns pointer
func zeroPtoS() *Entry {
    t := &Entry{}
    return t
}
```

The `zeroPtoS()` function returns a pointer to a zero-initialized structure.

```markup
// Initialized by the user - returns pointer
func initPtoS(N, S string, Y int) *Entry {
    if len(S) == 0 {
        return &Entry{Name: N, Surname: "Unknown", Year: Y}
    }
    return &Entry{Name: N, Surname: S, Year: Y}
}
```

The `initPtoS()` function also returns a pointer to a structure but also checks the length of the user input. Again, this checking is application-specific.

```markup
func main() {
    s1 := zeroS()
    p1 := zeroPtoS()
    fmt.Println("s1:", s1, "p1:", *p1)
    s2 := initS("Mihalis", "Tsoukalos", 2020)
    p2 := initPtoS("Mihalis", "Tsoukalos", 2020)
    fmt.Println("s2:", s2, "p2:", *p2)
    fmt.Println("Year:", s1.Year, s2.Year, p1.Year, p2.Year)
    pS := new(Entry)
    fmt.Println("pS:", pS)
}
```

The `new(Entry)` call returns a **pointer** to an `Entry` structure. Generally speaking, when you have to initialize lots of structure variables, it is considered a good practice to create a function for doing so as this is less error-prone.

Running `structures.go` creates the following output:

```markup
s1: {  0} p1: {  0}
s2: {Mihalis Tsoukalos 2020} p2: {Mihalis Tsoukalos 2020}
Year: 0 2020 0 2020
pS: &{  0}
```

As the zero value of a string is the empty string, `s1`, `p1`, and `pS` do not show any data for the `Name` and `Surname` fields.

The next subsection shows how to group structures of the same data type and use them as the elements of a slice.

## Slices of structures

You can create slices of structures in order to group and handle multiple structures under a single variable name. However, accessing a field of a given structure requires knowing the exact place of the structure in the slice.

For now, have a look at the following figure to better understand how a slice of structures works and how you can access the fields of a specific slice element.

![A picture containing shape
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_03_01.png)

Figure 3.1: A slice of structures

So, each slice element is a structure that is accessed using a slice index. Once we select the slice element we want, we can select its field.

As the whole process can be a little perplexing, the code of this subsection sheds some light and clarifies things. Type the following code and save it as `sliceStruct.go`. You can also find it by the same name in the `ch03` directory in the GitHub repository of the book.

```markup
package main
import (
    "fmt"
    "strconv"
)
type record struct {
    Field1 int
    Field2 string
}
func main() {
    S := []record{}
    for i := 0; i < 10; i++ {
        text := "text" + strconv.Itoa(i)
        temp := record{Field1: i, Field2: text}
        S = append(S, temp)
    }
```

You still need `append()` to add a new structure to the slice.

```markup
    // Accessing the fields of the first element
    fmt.Println("Index 0:", S[0].Field1, S[0].Field2)
    fmt.Println("Number of structures:", len(S))
    sum := 0
    for _, k := range S {
        sum += k.Field1
    }
    fmt.Println("Sum:", sum)
}
```

Running `sliceStruct.go` produces the following output:

```markup
Index 0: 0 text0
Number of structures: 10
Sum: 45
```

We revisit structures in the next chapter where we discuss reflection, as well as _Chapter 6_, _Telling a UNIX System What to Do_ where we learn how to work with JSON data using structures. For now, let us discuss regular expressions and pattern matching.

Just Imagine

# Regular expressions and pattern matching

**Pattern matching** is a technique for searching a string for some set of characters based on a specific search pattern that is based on regular expressions and grammars.

A **regular expression** is a sequence of characters that defines a search pattern. Every regular expression is compiled into a recognizer by building a generalized transition diagram called a **finite automaton**. A finite automaton can be either deterministic or nondeterministic. Nondeterministic means that more than one transition out of a state can be possible for the same input. A **recognizer** is a program that takes a string `x` as input and is able to tell whether `x` is a sentence of a given language.

A **grammar** is a set of production rules for strings in a formal language—the production rules describe how to create strings from the alphabet of the language that are valid according to the syntax of the language. A grammar does not describe the meaning of a string or what can be done with it in whatever context—it only describes its form. What is important here is to realize that grammars are at the heart of regular expressions because without a grammar, you cannot define or use a regular expression.

So, you might wonder why we're talking about regular expressions and pattern matching in this chapter. The reason is simple. In a while, you will learn how to store and read CSV data from plain text files, and you should be able to tell whether the data you are reading is valid or not.

## About Go regular expressions

We begin this subsection by presenting some common _match patterns_ used for constructing regular expressions.

<table id="table001-1" class="No-Table-Style _idGenTablePara-1"><colgroup><col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Heading--PACKT-">Expression</p></td><td class="No-Table-Style"><p class="Table-Column-Heading--PACKT-">Description</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">.</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Matches any character</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">*</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Means any number of times—cannot be used on its own</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">?</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Zero or one time—cannot be used on its own</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">+</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Means one or more times—cannot be used on its own</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">^</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">This denotes the beginning of the line</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">^</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">This denotes the end of the line</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">[]</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">[]</code> is for grouping characters</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">[A-Z]</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">This means all characters from capital <code class="Code-In-Text--PACKT-">A</code> to capital <code class="Code-In-Text--PACKT-">Z</code></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">\d</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Any digit in <code class="Code-In-Text--PACKT-">0-9</code></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">\D</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">A non-digit</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">\w</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Any word character: <code class="Code-In-Text--PACKT-">[0-9A-Za-z_]</code></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">\W</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">Any non-word character</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">\s</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">A whitespace character</p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-"><code class="Code-In-Text--PACKT-">\S</code></p></td><td class="No-Table-Style"><p class="Table-Column-Content--PACKT-">A non-whitespace character</p></td></tr></tbody></table>

The characters presented in the previous table are used for constructing and defining the grammar of the regular expression. The Go package responsible for defining regular expressions and performing pattern matching is called `regexp`. We use the `regexp.MustCompile()` function to create the regular expression and the `Match()` function to see whether the given string is a match or not.

The `regexp.MustCompile()` function parses the given regular expression and returns a `regexp.Regexp` variable that can be used for matching—`regexp.Regexp` is the representation of a **compiled regular expression**. The function panics if the expression cannot be parsed, which is good because you will know that your expression is invalid early in the process. The `re.Match()` method returns `true` if the given **byte slice** matches the `re` regular expression, which is a `regexp.Regexp` variable, and `false` otherwise.

Creating separate functions for pattern matching can be handy because it allows you to reuse the functions without worrying about the context of the program.

Keep in mind that although regular expressions and pattern matching look convenient and handy at first, they are the root of lots of bugs. My advice is to use the simplest regular expression that can solve your problem. However, if you can avoid using regular expressions at all, it would be much better in the long run!

## Matching names and surnames

The presented utility matches names and surnames, that is, strings that begin with an uppercase letter and continue with lowercase letters. The input should not contain any numbers or other characters.

The source code of the utility can be found in `nameSurRE.go`, which is located in the `ch03` directory. The function that supports the desired functionality is named `matchNameSur()` and is implemented as follows:

```markup
func matchNameSur(s string) bool {
    t := []byte(s)
    re := regexp.MustCompile(`^[A-Z][a-z]*$`)
    return re.Match(t)
}
```

The logic of the function is in the `` `^[A-Z][a-z]*$` `` regular expression, where `^` denotes the beginning of a line and `$` denotes the end of a line. What the regular expression does is match anything that begins with an uppercase letter (`[A-Z]`) and continues with any number of lowercase letters (`[a-z]*`). This means that `Z` is a match, but `ZA` is not a match because the second letter is uppercase. Similarly, `Jo+` is not a match because it contains a `+` character.

Running `nameSurRE.go` with various types of input produces the following output:

```markup
$ go run nameSurRE.go Z 
true
$ go run nameSurRE.go ZA
false
$ go run nameSurRE.go Mihalis
True
```

This technique can help you check user input.

## Matching integers

The presented utility matches both signed and unsigned integers—this is implemented in the way we define the regular expression. If we only want unsigned integers, then we should remove the `[-+]?` from the regular expression or replace it with `[+]?`.

The source code of the utility can be found in `intRE.go`, which is in the `ch03` directory. The `matchInt()` function that supports the desired functionality is implemented as follows:

```markup
func matchInt(s string) bool {
    t := []byte(s)
    re := regexp.MustCompile(`^[-+]?\d+$`)
    return re.Match(t)
}
```

As before, the logic of the function is found in the regular expression that is used for matching integers, which is `` `^[-+]?\d+$` ``. In plain English, what we say here is that we want to match something that begins with `–` or `+`, which is optional (`?`), and ends with any number of digits (`\d+`)—it is required that we have at least one digit before the end of the string that is examined (`$`).

Running `intRE.go` with various types of input produces the following output:

```markup
$ go run intRE.go 123
true
$ go run intRE.go /123
false
$ go run intRE.go +123.2
false
$ go run intRE.go +
false
$ go run intRE.go -123.2
false 
```

Later in this book, you will learn how to test Go code by writing testing functions—for now, we will do most of the testing manually.

## Matching the fields of a record

This example takes a different approach as we read an entire record and split it prior to doing any checking. Additionally, we make an extra check to make sure that the record we are processing contains the right number of fields. Each record should contain three fields: name, surname, and telephone number.

The full code of the utility can be found in `fieldsRE.go`, which is located in the `ch03` directory. The function that supports the desired functionality is implemented as follows:

```markup
func matchRecord(s string) bool {
    fields := strings.Split(s, ",")
    if len(fields) != 3 {
        return false
    }
    if !matchNameSur(fields[0]) {
        return false
    }
    if !matchNameSur(fields[1]) {
        return false
    }
    return matchTel(fields[2])
}
```

What the `matchRecord()` function does first is to separate the fields of the record based on the `,` character and then send each individual field to an appropriate function for further checking after making sure that the record has the right number of fields, which is a common practice. The field splitting is done using `strings.Split(s, ",")`, which returns a slice with as many elements as there are fields of the record.

If the checks of the first two fields are successful, then the function returns the return value of `matchTel(fields[2])` because it is that last check that determines the final result.

Running `fieldsRE.go` with various types of input produces the following output:

```markup
$ go run fieldsRE.go Name,Surname,2109416471
true
$ go run fieldsRE.go Name,Surname,OtherName 
false
$ go run fieldsRE.go One,Two,Three,Four
false
```

The first record is correct and therefore the `true` value is returned, which is not true for the second running where the phone number field is not correct. The last one failed because it contains four fields instead of three.

Just Imagine

# Improving the phone book application

It is time to update the phone book application. The new version of the phone book utility has the following improvements:

-   Support for the `insert` and `delete` commands
-   Ability to read data from a file and write it before it exits
-   Each entry has a last visited field that is updated
-   Has a database index that is implemented using a Go map
-   Uses regular expressions to verify the phone numbers read

## Working with CSV files

Most of the time you do not want to lose your data or have to begin without any data every time you execute your application. There exist many techniques for doing so—the easiest one is by saving your data locally. A very easy to work with format is CSV, which is what is explained here and used in the phone book application later on. The good thing is that Go provides a dedicated package for working with CSV data named `encoding/csv` ([https://golang.org/pkg/encoding/csv/](https://golang.org/pkg/encoding/csv/)). For the presented utility, both the input and output files are given as command-line arguments.

There exist two very popular Go interfaces named `io.Reader` and `io.Write` that are to do with reading from files and writing to files. Almost all reading and writing operations in Go use these two interfaces. The use of the same interface for readers allows readers to share some common characteristics but most importantly allows you to **create your own readers** and use them anywhere that Go expects an `io.Reader` reader. The same applies to writers that satisfy the `io.Write` interface. You will learn more about interfaces in _Chapter 4_, _Reflection and Interfaces_.

The main tasks that need to be implemented are the following:

-   Loading CSV data from disk and putting it into a slice of structures
-   Saving data to disk using CSV format

The `encoding/csv` package contains functions that can help you read and write CSV files. As we are dealing with small CSV files, we use `csv.NewReader(f).ReadAll()` to read the entire input file all at once. For bigger data files or if we wanted to check the input or make any changes to the input as we read it, it would have been better to read it line by line using `Read()` instead of `ReadAll()`.

Go assumes that the CSV file uses the comma character (`,`) for separating the different fields of each line. Should we wish to change that behavior, we should change the value of the `Comma` variable of the CSV reader or the writer depending on the task we want to perform. We change that behavior in the output CSV file, which separates its fields using the tab character.

For reasons of compatibility, it is better if the input and output CSV files are using the same field delimiter. We are just using the tab character as the field delimiter in the output file in order to illustrate the use of the `Comma` variable.

As working with CSV files is a new topic, there is a separate utility named `csvData.go` in the `ch03` directory of the GitHub repository of this book that illustrates the techniques for reading and writing CSV files. The source code of `csvData.go` is presented in chunks. First, we present the preamble of `csvData.go` that contains the `import` section as well as the `Record` structure and the `myData` global variable, which is a slice of `Record`.

```markup
package main
import (
    "encoding/csv"
    "fmt"
    "os"
)
type Record struct {
    Name       string
    Surname    string
    Number     string
    LastAccess string
}
var myData = []Record{}
```

Then we present the `readCSVFile()` function, which reads the plain text file with the CSV data.

```markup
func readCSVFile(filepath string) ([][]string, error) {
    _, err := os.Stat(filepath)
    if err != nil {
        return nil, err
    }
    f, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    // CSV file read all at once
    // lines data type is [][]string
    lines, err := csv.NewReader(f).ReadAll()
    if err != nil {
        return [][]string{}, err
    }
    return lines, nil
}
```

Note that we check whether the given file path exists and is associated with a regular file inside the function. There is no right or wrong decision about where to perform that checking—you just have to be consistent. The `readCSVFile()` function returns a `[][]string` slice that contains all the lines we have read. Additionally, have in mind that `csv.NewReader()` does separate the fields of each input line, which is the main reason for needing a slice with two dimensions to store the input.

After that, we illustrate the writing to a CSV file technique with the help of the `saveCSVFile()` function.

```markup
func saveCSVFile(filepath string) error {
    csvfile, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer csvfile.Close()
    csvwriter := csv.NewWriter(csvfile)
    // Changing the default field delimiter to tab
    csvwriter.Comma = '\t'
    for _, row := range myData {
        temp := []string{row.Name, row.Surname, row.Number, row.LastAccess}
        _ = csvwriter.Write(temp)
    }
    csvwriter.Flush()
    return nil
}
```

Note that change in the default value of `csvwriter.Comma`.

Last, we can see the implementation of the `main()` function.

```markup
func main() {
    if len(os.Args) != 3 {
        fmt.Println("csvData input output!")
        return
    }
    input := os.Args[1]
    output := os.Args[2]
    lines, err := readCSVFile(input)
    if err != nil {
        fmt.Println(err)
        return
    }
    // CSV data is read in columns - each line is a slice
    for _, line := range lines {
        temp := Record{
            Name:       line[0],
            Surname:    line[1],
            Number:     line[2],
            LastAccess: line[3],
        }
        myData = append(myData, temp)
        fmt.Println(temp)
    }
    err = saveCSVFile(output)
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

The `main()` function puts what you have read with `readCSVFile()` in the `myData` slice—remember that `lines` is a slice with two dimensions and that each row in `lines` is already separated into fields.

In this case, each line of input contains four fields. The contents of the CSV data file used as input are as follows:

```markup
$ cat ~/csv.data
Dimitris,Tsoukalos,2101112223,1600665563
Mihalis,Tsoukalos,2109416471,1600665563
Jane,Doe,0800123456,1608559903
```

Running `csvData.go` produces the following kind of output:

```markup
$ go run csvData.go ~/csv.data /tmp/output.data
{Dimitris Tsoukalos 2101112223 1600665563}
{Mihalis Tsoukalos 2109416471 1600665563}
{Jane Doe 0800123456 1608559903}
```

The contents of the output CSV file are the following:

```markup
$ cat /tmp/output.data
Dimitris        Tsoukalos       2101112223      1600665563
Mihalis Tsoukalos       2109416471      1600665563
Jane    Doe     0800123456      1608559903
```

The `output.data` file uses tab characters for separating the different fields of each record. The `csvData.go` utility can be handy for converting between different types of CSV files.

## Adding an index

This subsection explains how the database index is implemented with the help of a map. Indexing in databases is based on one or more keys that are unique. In practice, we index something that is unique and that we want to access quickly. In the database case, primary keys are unique by default and cannot be present in more than one record. In our case, phone numbers are used as primary keys, which means that the index is built based on the phone number field of the structure.

As a rule of thumb, you index a field that is going to be used for searching. There is no point in creating an index that is not going to be used for querying.

Let us now see what this means in practice. ImagineDevOps  that we have a slice named `S` with the following kind of data:

```markup
S[0]={0800123123, ...}
S[1]={0800123000, ...}
S[2]={2109416471, ...}
.
.
.
```

So, each slice element is a structure that can contain much more data apart from the telephone number. How can we create an index for it? The index, which is named `Index`, is going to have the following data and format:

```markup
Index["0800123123"] = 0
Index["0800123000"] = 1
Index["2109416471"] = 2
.
.
.
```

This means that if we want to look for the telephone number `0800123000`, we should see whether `0800123000` exists as a key in `Index`. If it is there, then we know that the value of `0800123000`, which is `Index["0800123000"]`, is the index of the slice element that contains the desired record. So, as we know which slice element to access, we do not have to search the entire slice. With that in mind, let us update the application.

## The improved version of the phone book application

It would be a shame to create a phone book application that has its entries hardcoded in the code. This time, the entries of the address book are read from an external file that contains data in CSV format. Similarly, the new version saves its data into the same CSV file, which you can read afterward.

Each entry of the phone book application is based on the following structure:

```markup
type Entry struct {
    Name       string
    Surname    string
    Tel        string
    LastAccess string
}
```

The key to the entries is the `Tel` field and therefore its values. In practice, this means that if you try to add an entry that uses an existing `Tel` value, the process fails. This also means that the application searches the phone book using the `Tel` values. Databases use primary keys to identify between unique records—the phone book application has a small database implemented as a slice of `Entry` structures. Last, phone numbers are saved without any `–` characters in them. The utility removes all `–` characters from phone numbers, if there are any, before saving them.

Personally, I prefer to explore the various parts of a bigger application by creating smaller programs that when combined implement some or all of the functionality of the bigger program. This helps me understand how the bigger application needs to be implemented. This makes it much easier for me to connect all the dots afterward and develop the final product.

As this is a real application implemented as a command-line utility, it should support commands for data manipulation and searching. The updated functionality of the utility is explained in the following list:

-   Data insertion using the `insert` command
-   Data deletion using the `delete` command
-   Data searching using the `search` command
-   Listing of the available records through the `list` command

To make the code simpler, the path of the CSV data file is hardcoded. Additionally, the CSV file is automatically read when the utility is executed but it is automatically updated when the `insert` and the `delete` commands are executed.

Although Go supports CSV, **JSON** is a far more popular format that is used for data exchange in web services. However, working with CSV data is simpler than working with data in JSON format. Working with JSON data is explored in _Chapter 6_, _Telling a UNIX System What to Do_.

As explained earlier, this version of the phone book application has support for _indexing_ to find the desired records faster without needing to make a linear search to the slice that holds your phone book data. The indexing technique that is used is not very fancy, but it makes searching really fast: provided that the searching process is based on phone numbers, we are going to create a map that associates each phone number with the index number of the record that contains that phone number in the slice of structures. This way, a simple and fast map lookup tells us whether a phone number already exists or not. If the phone number exists, we can access its record directly without having to search the entire slice of structures for it. The only downside of this technique, and every indexing technique, is that you must keep the map up to date all the time.

The previous process is called a high-level design of the application. For such a simple application, you do not have to be too analytic about the capabilities of the application—stating the supported commands and the location of the CSV data file is enough. However, for a RESTful server that implements a REST API, the design phase or the dependencies of the program are as important as the development phase itself.

The entire code of the updated phone book utility can be found in `ch03` as `phoneCourse.go`—as always, we are referring to the GitHub repository of the book. This is the last time we make this clarification—from now on, we will only tell you the name of the source file unless there is a specific reason to do otherwise.

The most interesting parts of the `phoneCourse.go` file are presented here, starting from the implementation of the `main()` function, which is presented in two parts. The first part is about getting a command to execute and having a valid CSV file to work with.

```markup
func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Usage: insert|delete|search|list <arguments>")
        return
    }
    // If the CSVFILE does not exist, create an empty one
    _, err := os.Stat(CSVFILE)
    // If error is not nil, it means that the file does not exist
    if err != nil {
        fmt.Println("Creating", CSVFILE)
        f, err := os.Create(CSVFILE)
        if err != nil {
            f.Close()
            fmt.Println(err)
            return
        }
        f.Close()
    }
```

If the file path specified by the `CSVFILE` global variable does not already exist, we have to create it for the rest of the program to use it. This is determined by the return value of the `os.Stat(CSVFILE)` call.

```markup
    fileInfo, err := os.Stat(CSVFILE)
    // Is it a regular file?
    mode := fileInfo.Mode()
    if !mode.IsRegular() {
        fmt.Println(CSVFILE, "not a regular file!")
        return
    }
```

Not only must the `CSVFILE` exist but it should be a regular UNIX file, which is determined by the call to `mode.IsRegular()`. If it is not a regular file, the utility prints an error message and exits.

```markup
    err = readCSVFile(CSVFILE)
    if err != nil {
        fmt.Println(err)
        return
    }
```

This is the place where we read the `CSVFILE`, even if it is empty. The contents of `CSVFILE` are kept in the `data` global variable that is defined as `[]Entry{}`, which is a slice of `Entry` variables.

```markup
    err = createIndex()
    if err != nil {
        fmt.Println("Cannot create index.")
        return
    }
```

This is where we create the index by calling `createIndex()`. The index is kept in the `index` global variable that is defined as `map[string]int`.

The second part of `main()` is about running the right command and understanding whether the command was executed successfully or not.

```markup
    // Differentiating between the commands
    switch arguments[1] {
    case "insert":
        if len(arguments) != 5 {
            fmt.Println("Usage: insert Name Surname Telephone")
            return
        }
        t := strings.ReplaceAll(arguments[4], "-", "")
        if !matchTel(t) {
            fmt.Println("Not a valid telephone number:", t)
            return
        }
```

You need to remove any `–` characters from the telephone number before storing it using `strings.ReplaceAll()`. If there are not any `–` characters, then no substitution takes place.

```markup
        temp := initS(arguments[2], arguments[3], t)
        // If it was nil, there was an error
        if temp != nil {
            err := insert(temp)
            if err != nil {
                fmt.Println(err)
                return
            }
        }
    case "delete":
        if len(arguments) != 3 {
            fmt.Println("Usage: delete Number")
            return
        }
        t := strings.ReplaceAll(arguments[2], "-", "")
        if !matchTel(t) {
            fmt.Println("Not a valid telephone number:", t)
            return
        }
        err := deleteEntry(t)
        if err != nil {
            fmt.Println(err)
        }
    case "search":
        if len(arguments) != 3 {
            fmt.Println("Usage: search Number")
            return
        }
        t := strings.ReplaceAll(arguments[2], "-", "")
        if !matchTel(t) {
            fmt.Println("Not a valid telephone number:", t)
            return
        }
        temp := search(t)
        if temp == nil {
            fmt.Println("Number not found:", t)
            return
        }
        fmt.Println(*temp)
    case "list":
        list()
    default:
        fmt.Println("Not a valid option")
    }
}
```

In this relatively big `switch` block, we can see what is executed for each given command. So, we have the following:

-   For the `insert` command, we execute the `insert()` function
-   For the `list` command, we execute the `list()` function, which is the only function that requires any arguments
-   For the `delete` command, we execute the `deleteEntry()` function
-   For the `search` command, we execute the `search()` function

Anything else is handled by the `default` branch. The index is created and updated using the `createIndex()` function, which is implemented as follows:

```markup
func createIndex() error {
    index = make(map[string]int)
    for i, k := range data {
        key := k.Tel
        index[key] = i
    }
    return nil
}
```

Put simply, you access the entire `data` slice and put the index and value pairs of the slice in a map using the value as the key for the map and the slice index as the value of the map.

The `delete` command is implemented as follows:

```markup
func deleteEntry(key string) error {
    i, ok := index[key]
    if !ok {
        return fmt.Errorf("%s cannot be found!", key)
    }
    data = append(data[:i], data[i+1:]...)
    // Update the index - key does not exist any more
    delete(index, key)
    err := saveCSVFile(CSVFILE)
    if err != nil {
        return err
    }
    return nil
}
```

The operation of the `deleteEntry()` function is simple. First, you search the index for the telephone number in order to find the place of the entry in the slice with the data—if it does not exist, then you simply create an error message with `fmt.Errorf("%s cannot be found!", key)` and the function returns. If the telephone number can be found, then you delete that relevant entry from the `data` slice with `append(data[:i], data[i+1:]...)`.

Then, you must update the index—taking care of the index is the price you have to pay for the extra speed that the index gives you. Also, after you delete an entry, you should save the updated data by calling `saveCSVFile(CSVFILE)` for changes to take effect.

Strictly speaking, as the current version of the phone book application processes one request at a time, you do not need to update the index because it is created from scratch each time you use the application. On database management systems, indexes are also saved on disk in order to avoid the major cost of creating them from scratch.

The `insert` command is implemented as follows:

```markup
func insert(pS *Entry) error {
    // If it already exists, do not add it
    _, ok := index[(*pS).Tel]
    if ok {
        return fmt.Errorf("%s already exists", pS.Tel)
    }
    data = append(data, *pS)
    // Update the index
    _ = createIndex()
    err := saveCSVFile(CSVFILE)
    if err != nil {
        return err
    }
    return nil
}
```

The index here helps you determine whether the telephone number you are trying to add already exists or not—as stated earlier, if you try to add an entry that uses an existing `Tel` value, the process fails. If this test passes, you add the new record in that `data` slice, update the index, and save the data to the CSV file.

The `search` command, which uses the index, is implemented as follows:

```markup
func search(key string) *Entry {
    i, ok := index[key]
    if !ok {
        return nil
    }
    data[i].LastAccess = strconv.FormatInt(time.Now().Unix(), 10)
    return &data[i]
}
```

Due to the index, searching for a telephone number is straightforward—the code just looks in the `index` for the desired telephone number. If it is present, the code returns that record—otherwise, the code returns `nil`, which is possible because the function returns a pointer to an `Entry` variable. Before returning the record, the `search()` function updates the `LastAccess` field of the structure that is about to be returned in order to know the last time it was accessed.

The initial contents of the CSV data file used as input are as follows:

```markup
$ cat ~/csv.data 
Dimitris,Tsoukalos,2101112223,1600665563
Mihalis,Tsoukalos,2109416471,1600665563
Mihalis,Tsoukalos,2109416771,1600665563
Efipanios,Savva,2101231234,1600665582
```

As long at the telephone number is unique, the name and surname fields can exist multiple times. Running `phoneCourse.go` produces the following kind of output:

```markup
$ go run phoneCourse.go list
{Dimitris Tsoukalos 2101112223 1600665563}
{Mihalis Tsoukalos 2109416471 1600665563}
{Mihalis Tsoukalos 2109416771 1600665563}
{Efipanios Savva 2101231234 1600665582}
$ go run phoneCourse.go delete 2109416771
$ go run phoneCourse.go search 2101231234
{Efipanios Savva 2101231234 1608559833}
$ go run phoneCourse.go search 210-1231-234
{Efipanios Savva 2101231234 1608559840}
```

Due to our code, `210-1231-234` is converted to `2101231234`.

```markup
$ go run phoneCourse.go delete 210-1231-234
$ go run phoneCourse.go search 210-1231-234
Number not found: 2101231234
$ go run phoneCourse.go insert Jane Doe 0800-123-456
$ go run phoneCourse.go insert Jane Doe 0800-123-456
0800123456 already exists
$ go run phoneCourse.go search 2101112223 
{Dimitris Tsoukalos 2101112223 1608559928}
```

The contents of the CSV file after the previous commands are the following:

```markup
$ cat ~/csv.data 
Dimitris,Tsoukalos,2101112223,1600665563
Mihalis,Tsoukalos,2109416471,1600665563
Jane,Doe,0800123456,1608559903
```

You can see that the utility converted `0800-123-456` into `0800123456`, which is the desired behavior.

Despite being much better than the previous version, the new version of the phone book utility is still not perfect. Here is a list of things that can be improved:

-   Ability to sort its output based on the telephone number
-   Ability to sort its output based on the surname
-   Ability to use JSON records and JSON slices for the data instead of CSV files

The phone book application will keep improving, starting from the next chapter, where sorting slices with structure elements is implemented.

Just Imagine

# Exercises

-   Write a Go program that converts an existing array into a map.
-   Write a Go program that converts an existing map into two slices—the first slice contains the keys of the map whereas the second slice contains the values. The values at index `n` of the two slices should correspond to a key and value pair that can be found in the original map.
-   Make the necessary changes to `nameSurRE.go` to be able to process multiple command-line arguments.
-   Change the code of `intRE.go` to process multiple command-line arguments and display totals of `true` and `false` results at the end.
-   Make changes to `csvData.go` to separate the fields of a record based on the `#` character.
-   Write a Go utility that converts `os.Args` into a slice of structures with fields for storing the index and the value of each command-line argument—you should define the structure that is going to be used on your own.
-   Make the necessary changes to `phoneCourse.go` in order to create the index based on the `LastAccess` field. Is this practical? Does it work? Why?
-   Make changes to `csvData.go` in order to separate the fields of a record with a character that is given as a command-line argument.

Just Imagine

# Summary

In this chapter we discussed the composite data types of Go, which are maps and structures. Additionally, we talked about working with CSV files as well as about using regular expressions and pattern matching in Go. We can now keep our data in proper structures, validate it using regular expressions, when this is possible, and store it in CSV files to achieve data persistency.

The next chapter is about type methods, which are functions attached to a data type, reflection, and interfaces. All these things will allow us to improve the phone book application.

Just Imagine

# Additional resources

-   The `encoding/csv` documentation: [https://golang.org/pkg/encoding/csv/](https://golang.org/pkg/encoding/csv/)
-   The `runtime` package documentation: [https://golang.org/pkg/runtime/](https://golang.org/pkg/runtime/)