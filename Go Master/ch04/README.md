# Reflection and Interfaces

Do you remember the phone book application from the previous chapter? You might wonder what happens if you want to sort user-defined data structures, such as phone book records, based on your own criteria, such as a surname or first name. What happens when you want to sort different datasets that share some common behavior without having to implement sorting from scratch for each one of the different data types using multiple functions? Now imagine that you have a utility like the phone book application that can process two different formats of CSV data files based on the given input file. Each kind of CSV record is stored in a different Go structure, which means that each kind of CSV record might be sorted differently. How do you implement that without having to write two different command-line utilities? Lastly, imagine that you want to write a utility that sorts really unusual data. For example, imagine that you want to sort a slice that holds various kinds of 3D shapes based on their volume. Can this be performed easily and in a way that makes sense?

The answer to all these questions and concerns is the use of _interfaces_. However, interfaces are not just about data manipulation and sorting. Interfaces are about expressing abstractions and identifying and defining behaviors that can be shared among different data types. Once you have implemented an interface for a data type, a new world of functionality becomes available to the variables and the values of that type, which can save you time and increase your productivity. Interfaces work with _methods on types_ or _type methods_, which are like functions attached to given data types, which in Go are usually structures. Remember that once you implement the required type methods of an interface, that interface is **satisfied implicitly**, which is also the case with the _empty interface_ that is explained in this chapter.

Another handy Go feature is _reflection_, which allows you to examine the structure of a data type at execution time. However, as reflection is an advanced Go feature, you do not need to use it on a regular basis.

This chapter covers:

-   Reflection
-   Type methods
-   Interfaces
-   Working with two different CSV file formats
-   Object-oriented programming in Go
-   Updating the phone book application

Just Imagine

# Reflection

We begin this chapter with reflection, which is an advanced Go feature, not because it is an easy subject but because it is going to help you understand how Go works with different data types, including interfaces, and why they are needed.

You might be wondering how you can find out the names of the fields of a structure at execution time. In such cases, you need to use reflection. Apart from enabling you to print the fields and the values of a structure, reflection also allows you to explore and manipulate unknown structures like the ones created from decoding JSON data.

The two main questions that I asked myself when I was introduced to reflection for the first time were the following:

-   Why was reflection included in Go?
-   When should I use reflection?

To answer the first question, reflection allows you to **dynamically** learn the type of an arbitrary object along with information about its structure. Go provides the `reflect` package for working with reflection. Remember when we said in a previous chapter that `fmt.Println()` is clever enough to understand the data types of its parameters and act accordingly? Well, behind the scenes, the `fmt` package uses reflection to do that.

As far as the second question is concerned, reflection allows you to handle and work with data types that do not exist at the time at which you write your code but might exist in the future, which is when we use an existing package with user-defined data types.

Additionally, reflection might come in handy when you have to work with data types that do not implement a common interface and therefore have an uncommon or unknown behavior—this does not mean that they have bad or erroneous behavior, just uncommon behavior such as a user-defined structure.

The introduction of _generics_ in Go might make the use of reflection less frequent in some cases, because with generics you can work with different data types more easily and without the need to know their exact data types in advance. However, nothing beats reflection for fully exploring the structure and the data types of a variable. We talk about reflection compared to generics in _Chapter 13_, _Go Generics_.

The most useful parts of the `reflect` package are two data types named `reflect.Value` and `reflect.Type`. Now, `reflect.Value` is used for storing values of any type, whereas `reflect.Type` is used for representing Go types. There exist two functions named `reflect.TypeOf()` and `reflect.ValueOf()` that return the `reflect.Type` and `reflect.Value` values, respectively. Note that `reflect.TypeOf()` returns the actual type of variable—if we are examining a structure, it returns the name of the structure.

As structures are really important in Go, the `reflect` package offers the `reflect.NumField()` method for listing the number of fields in a structure as well as the `Field()` method for getting the `reflect.Value` value of a specific field of a structure.

The `reflect` package also defines the `reflect.Kind` data type, which is used for representing the _specific_ data type of a variable: `int`, `struct`, etc. The documentation of the `reflect` package lists all possible values of the `reflect.Kind` data type. The `Kind()` function returns the kind of a variable.

Last, the `Int()` and `String()` methods return the integer and string value of a `reflect.Value`, respectively.

Reflection code can look unpleasant and hard to read sometimes. Therefore, according to the Go philosophy, you should rarely use reflection unless it is absolutely necessary because despite its cleverness, it does not create clean code.

## Learning the internal structure of a Go structure

The next utility shows how to use reflection to discover the internal structure and fields of a Go structure variable. Type it and save it as `reflection.go`.

```markup
package main
import (
    "fmt"
    "reflect"
)
type Secret struct {
    Username string
    Password string
}
type Record struct {
    Field1 string
    Field2 float64
    Field3 Secret
}
func main() {
    A := Record{"String value", -12.123, Secret{"Mihalis", "Tsoukalos"}}
```

We begin by defining a `Record` structure variable that contains another structure value (`Secret{"Mihalis", "Tsoukalos"}`).

```markup
    r := reflect.ValueOf(A)
    fmt.Println("String value:", r.String())
```

This returns the `reflect.Value` of the `A` variable.

```markup
    iType := r.Type()
```

Using `Type()` is how we get the data type of a variable—in this case variable `A`.

```markup
    fmt.Printf("i Type: %s\n", iType)
    fmt.Printf("The %d fields of %s are\n", r.NumField(), iType)
    for i := 0; i < r.NumField(); i++ {
```

The previous `for` loop allows you to visit all the fields of a structure and examine their characteristics.

```markup
        fmt.Printf("\t%s ", iType.Field(i).Name)
        fmt.Printf("\twith type: %s ", r.Field(i).Type())
        fmt.Printf("\tand value _%v_\n", r.Field(i).Interface())
```

The previous `fmt.Printf()` statements return the name, the data type, and the value of the fields.

```markup
        // Check whether there are other structures embedded in Record
        k := reflect.TypeOf(r.Field(i).Interface()).Kind()
        // Need to convert it to string in order to compare it
        if k.String() == "struct" {
```

In order to check the data type of a variable with a string, we need to convert the data type into a `string` variable first.

```markup
            fmt.Println(r.Field(i).Type())
        }
        // Same as before but using the internal value
        if k == reflect.Struct {
```

You can also use the internal representation of a data type during checking. However, this makes less sense than using a `string` value.

```markup
            fmt.Println(r.Field(i).Type())
        }
    }
}
```

Running `reflection.go` produces the following output:

```markup
$ go run reflection.go
String value: <main.Record Value>
i Type: main.Record
The 3 fields of main.Record are
        Field1  with type: string       and value _String value_
        Field2  with type: float64      and value _-12.123_
        Field3  with type: main.Secret  and value _{Mihalis Tsoukalos}_
main.Secret
main.Secret
```

`main.Record` is the full unique name of the structure as defined by Go—`main` is the package name and `Record` is the `struct` name. This happens so that Go can differentiate between the elements of different packages.

The presented code does not modify any values of the structure. If you were to make changes to the values of the structure fields, you would use the `Elem()` method and pass the structure as a pointer to `ValueOf()`—remember that pointers allow you to make changes to the actual variable. There exist methods that allow you to modify an existing value. In our case, we are going to use `SetString()` for modifying a `string` field and `SetInt()` for modifying an `int` field.

This technique is illustrated in the next subsection.

## Changing structure values using reflection

Learning about the internal structure of a Go structure is handy, but what is more practical is being able to change values in the Go structure, which is the subject of this subsection.

Type the following Go code and save it as `setValues.go`—it can also be found in the GitHub repository of the book.

```markup
package main
import (
    "fmt"
    "reflect"
)
type T struct {
    F1 int
    F2 string
    F3 float64
}
func main() {
    A := T{1, "F2", 3.0}
```

`A` is the variable that is examined in this program.

```markup
    fmt.Println("A:", A)
    r := reflect.ValueOf(&A).Elem()
```

With the use of `Elem()` and a pointer to variable `A`, variable `A` can be modified if needed.

```markup
    fmt.Println("String value:", r.String())
    typeOfA := r.Type()
    for i := 0; i < r.NumField(); i++ {
        f := r.Field(i)
        tOfA := typeOfA.Field(i).Name
        fmt.Printf("%d: %s %s = %v\n", i, tOfA, f.Type(), f.Interface())
        k := reflect.TypeOf(r.Field(i).Interface()).Kind()
        if k == reflect.Int {
            r.Field(i).SetInt(-100)
        } else if k == reflect.String {
            r.Field(i).SetString("Changed!")
        }
    }
```

We are using `SetInt()` for modifying an integer value and `SetString()` for modifying a `string` value. Integer values are set to `-100` and string values are set to `Changed!`.

```markup
    fmt.Println("A:", A)
}
```

Running `setValues.go` creates the next output:

```markup
$ go run setValues.go
A: {1 F2 3}
String value: <main.T Value>
0: F1 int = 1
1: F2 string = F2
2: F3 float64 = 3
A: {-100 Changed! 3}
```

The first line of output shows the initial version of `A` whereas the last line shows the final version of `A` with the modified fields. The main use of such code is _dynamically_ changing the values of the fields of a structure.

## The three disadvantages of reflection

Without a doubt, reflection is a powerful Go feature. However, as with all tools, reflection should be used sparingly for three main reasons:

-   The first reason is that extensive use of reflection will make your programs hard to read and maintain. A potential solution to this problem is good documentation, but developers are notorious for not having the time to write proper documentation.
-   The second reason is that the Go code that uses reflection makes your programs slower. Generally speaking, Go code that works with a particular data type is always faster than Go code that uses reflection to dynamically work with any Go data type. Additionally, such dynamic code makes it difficult for tools to refactor or analyze your code.
-   The last reason is that reflection errors cannot be caught at build time and are reported at runtime as panics, which means that reflection errors can potentially crash your programs. This can happen months or even years after the development of a Go program! One solution to this problem is extensive testing before a dangerous function call. However, this adds even more Go code to your programs, which makes them even slower.

Now that we know about reflection and what it can do for us, it is time to begin the discussion about type methods, which are necessary for using interfaces.

Just Imagine

# Type methods

A **type method** is a function that is attached to a specific data type. Although type methods (or methods on types) are in reality functions, they are defined and used in a slightly different way.

The methods on types feature gives some object-oriented capabilities to Go, which is very handy and is used extensively in Go. Additionally, interfaces require type methods to work.

Defining new type methods is as simple as creating new functions, provided that you follow certain rules that associate the function with a data type.

## Creating type methods

So, imagine that you want to do calculations with 2x2 matrices. A very natural way of implementing that is by defining a new data type and defining type methods for adding, subtracting, and multiplying 2x2 matrices using that new data type. To make it even more interesting and generic, we are going to create a command-line utility that accepts the elements of two 2x2 matrices as command-line arguments, which are eight integer values in total, and performs all three calculations between them using the defined type methods.

Having a data type called `ar2x2`, you can create a type method named `FunctionName` for it as follows:

```markup
func (a ar2x2) FunctionName(parameters) <return values> {
    ...
}
```

The `(a ar2x2)` part is what makes the `FunctionName()` function a type method because it associates `FunctionName()` with the `ar2x2` data type. **No other data type can use that function**. However, you are free to implement `FunctionName()` for other data types or as a regular function. If you have a `ar2x2` variable named `varAr`, you can invoke `FunctionName()` as `varAr.FunctionName(...)`, which looks like selecting the field of a structure variable.

You are not obligated to develop type methods if you do not want to. In fact, each type method can be **rewritten as a regular function**. Therefore, `FunctionName()` can be rewritten as follows:

```markup
func FunctionName(a ar2x2, parameters...) <return values> {
    ...
}
```

Have in mind that under the hood, the Go compiler does turn methods into regular function calls with the `self` value as the first parameter. However, **interfaces require the use of type methods** to work.

The expressions used for selecting a field of a structure or a type method of a data type, which would replace the ellipsis after the variable name above, are called **selectors**.

Performing calculations between matrices of a given size is one of the rare cases where using an array instead of a slice makes more sense because you do not have to modify the size of the matrices. Some might argue that using a slice instead of an array pointer is a better practice—you are allowed to use what makes more sense to you.

Most of the time, and when there is such a need, the results of a type method are saved in the variable that invoked the type method—in order to implement that for the `ar2x2` data type, we pass a pointer to the array that invoked the type method, like `func (a *ar2x2)`.

The next subsection illustrates type methods in action.

## Using type methods

This subsection shows the use of type methods using the `ar2x2` data type as an example. The `Add()` function and the `Add()` method use the exact same algorithm for adding two matrices. The only difference between them is the way they are being called and the fact that the function returns an array whereas the method saves the result to the calling variable.

Although adding and subtracting matrices is a straightforward process—you just add or subtract each element of the first matrix with the element of the second matrix that is located at the same position—matrix multiplication is a more complex process. This is the main reason that both addition and subtraction use `for` loops, which means that the code can also work with bigger matrices, whereas multiplication uses static code that cannot be applied to bigger matrices without major changes.

If you are defining type methods for a structure, you should make sure that the names of the type methods do not conflict with any field name of the structure because the Go compiler will reject such ambiguities.

Type the following code and save it as `methods.go`.

```markup
package main
import (
    "fmt"
    "os"
    "strconv"
)
type ar2x2 [2][2]int
// Traditional Add() function
func Add(a, b ar2x2) ar2x2 {
    c := ar2x2{}
    for i := 0; i < 2; i++ {
        for j := 0; j < 2; j++ {
            c[i][j] = a[i][j] + b[i][j]
        }
    }
    return c
}
```

Here, we have a traditional function that adds two `ar2x2` variables and returns their result.

```markup
// Type method Add()
func (a *ar2x2) Add(b ar2x2) {
    for i := 0; i < 2; i++ {
        for j := 0; j < 2; j++ {
            a[i][j] = a[i][j] + b[i][j]
        }
    }
}
```

Here we have a type method named `Add()` that is attached to the `ar2x2` data type. The result of the addition is not returned. What happens is that the `ar2x2` variable that called the `Add()` method is going to be modified and hold the result—this is the reason for using a pointer when defining the type method. If you do not want that behavior, you should modify the signature and the implementation of the type method to fit your needs.

```markup
// Type method Subtract()
func (a *ar2x2) Subtract(b ar2x2) {
    for i := 0; i < 2; i++ {
        for j := 0; j < 2; j++ {
            a[i][j] = a[i][j] - b[i][j]
        }
    }
}
```

The previous method subtracts `ar2x2` `b` from `ar2x2` `a` and the result is saved in `a`.

```markup
// Type method Multiply()
func (a *ar2x2) Multiply(b ar2x2) {
    a[0][0] = a[0][0]*b[0][0] + a[0][1]*b[1][0]
    a[1][0] = a[1][0]*b[0][0] + a[1][1]*b[1][0]
    a[0][1] = a[0][0]*b[0][1] + a[0][1]*b[1][1]
    a[1][1] = a[1][0]*b[0][1] + a[1][1]*b[1][1]
}
```

As we are working with small arrays, we do the multiplications without using `for` loops.

```markup
func main() {
    if len(os.Args) != 9 {
        fmt.Println("Need 8 integers")
        return
    }
    k := [8]int{}
    for index, i := range os.Args[1:] {
        v, err := strconv.Atoi(i)
        if err != nil {
            fmt.Println(err)
            return
        }
        k[index] = v
    }
    a := ar2x2{{k[0], k[1]}, {k[2], k[3]}}
    b := ar2x2{{k[4], k[5]}, {k[6], k[7]}}
```

The `main()` function gets the input and creates two 2x2 matrices. After that, it performs the desired calculations with these two matrices.

```markup
    fmt.Println("Traditional a+b:", Add(a, b))
    a.Add(b)
    fmt.Println("a+b:", a)
    a.Subtract(a)
    fmt.Println("a-a:", a)
    a = ar2x2{{k[0], k[1]}, {k[2], k[3]}}
```

We calculate `a+b` using two different ways: using a regular function and using a type method. As both `a.Add(b)` and `a.Subtract(a)` change the value of `a`, we have to initialize `a` before using it again.

```markup
    a.Multiply(b)
    fmt.Println("a*b:", a)
    a = ar2x2{{k[0], k[1]}, {k[2], k[3]}}
    b.Multiply(a)
    fmt.Println("b*a:", b)
}
```

Last, we calculate `a*b` and `b*a` to show that they are different because the commutative property does not apply to matrix multiplication.

Running `methods.go` produces the next output:

```markup
$ go run methods.go 1 2 0 0 2 1 1 1
Traditional a+b: [[3 3] [1 1]]
a+b: [[3 3] [1 1]]
a-a: [[0 0] [0 0]]
a*b: [[4 6] [0 0]]
b*a: [[2 4] [1 2]]
```

The input here is two 2x2 matrices, `[[1 2] [0 0]]` and `[[2 1] [1 1]]`, and the output is their calculations.

Now that we know about type methods, it is time to begin exploring interfaces as interfaces cannot be implemented without type methods.

Just Imagine

# Interfaces

An **interface** is a Go mechanism for defining behavior that is implemented using a set of methods. Interfaces play a key role in Go and can simplify the code of your programs when they have to deal with multiple data types that perform the same task—recall that `fmt.Println()` works for almost all data types. But remember, interfaces should not be unnecessarily complex. If you decide to create your own interfaces, then you should begin with a common behavior that you want to be used by multiple data types.

Interfaces work with **methods on types** (or **type methods**), which are like functions attached to given data types, which in Go are usually structures (although we can use any data type we want).

As you already know, once you implement the required type methods of an interface, that interface is **satisfied implicitly**.

The **empty interface** is defined as just `interface{}`. As the empty interface has no methods, it means that it is already **implemented by all data types**.

Once you implement the methods of an interface for a data type, that interface is satisfied **automatically** for that data type.

In a more formal way, a Go interface type defines (or describes) the **behavior** of other types by specifying a set of _methods_ that need to be implemented for supporting that behavior. For a data type to satisfy an interface, it needs to implement **all the type methods** required by that interface. Therefore, interfaces are **abstract types** that specify a set of methods that need to be implemented so that another type can be considered an instance of the interface. So, an interface is two things: **a set of methods and a type**. Have in mind that small and well-defined interfaces are usually the most popular ones.

As a rule of thumb, only create a new interface when you want to share a common behavior between two or more concrete data types. This is basically **duck typing**.

The biggest **advantage** you get from interfaces is that if needed, you can pass a variable of a data type that implements a particular interface to any function that expects a parameter of that specific interface, which saves you from having to write separate functions for each supported data type. However, Go offers an alternative to this with the recent addition of generics.

Interfaces can also be used for providing a kind of polymorphism in Go, which is an object-oriented concept. _Polymorphism_ offers a way of accessing objects of different types in the same uniform way when they share a common behavior.

Lastly, interfaces can be used for _composition_. In practice, this means that you can combine existing interfaces and create new ones that offer the combined behavior of the interfaces that were brought together. The next figure shows interface composition in a graphical way.

![Graphical user interface, application, Word
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_04_01.png)

Figure 4.1: Interface composition

Put simply, the previous figure illustrates that because of its definition, satisfying interface `ABC` requires satisfying `InterfaceA`, `InterfaceB`, and `InterfaceC`. Additionally, any `ABC` variable can be used instead of an `InterfaceA` variable, an `InterfaceB` variable, or an `InterfaceC` variable because it supports all these three behaviors. Last, only `ABC` variables can be used where an `ABC` variable is expected. There is nothing prohibiting you from including additional methods in the definition of the `ABC` interface if the combination of existing interfaces does not describe the desired behavior accurately.

When you combine existing interfaces, it is better that the interfaces do not contain methods with the same name.

What you should keep in mind is that there is no need for an interface to be impressive and require the implementation of a large number of methods. In fact, the fewer methods an interface has, the more generic and widely used it can be, which improves its usefulness and therefore its usage.

The subsection that follows illustrates the use of `sort.Interface`.

## The sort.Interface interface

The `sort` package contains an interface named `sort.Interface` that allows you to sort slices according to your needs and your data, provided that you implement `sort.Interface` for the custom data types stored in your slices. The `sort` package defines the `sort.Interface` as follows:

```markup
type Interface interface {
    // Len is the number of elements in the collection.
    Len() int
    // Less reports whether the element with
    // index i should sort before the element with index j.
    Less(i, j int) bool
    // Swap swaps the elements with indexes i and j.
    Swap(i, j int)
}
```

What we can understand from the definition of `sort.Interface` is that in order to implement `sort.Interface`, we need to implement the following three type methods:

-   `Len() int`
-   `Less(i, j int) bool`
-   `Swap(i, j int)`

The `Len()` method returns the length of the slice that will be sorted and helps the interface to process all slice elements whereas the `Less()` method, which compares and sorts elements in pairs, defines how elements are going to be compared and therefore sorted. The return value of `Less()` is `bool`, which means that `Less()` only cares about whether the element at index `i` is bigger or not than the element at index `j` in the way that the two elements are being compared. Lastly, the `Swap()` method is used for swapping two elements of the slice, which is required for the sorting algorithm to work.

The following code, which can be found as `sort.go`, illustrates the use of `sort.Interface`.

```markup
package main
import (
    "fmt"
    "sort"
)
type S1 struct {
    F1 int
    F2 string
    F3 int
}
// We want to sort S2 records based on the value of F3.F1,
// Which is equivalent to S1.F1 as F3 is an S1 structure
type S2 struct {
    F1 int
    F2 string
    F3 S1
}
```

The `S2` structure includes a field named `F3` that is of the `S1` data type, which is also a structure.

```markup
type S2slice []S2
```

You need to have a slice because all sorting operations work on slices. It is for this slice, which should be a new data type that in this case is called `S2slice`, that you are going to implement the three type methods of the `sort.Interface`.

```markup
// Implementing sort.Interface for S2slice
func (a S2slice) Len() int {
    return len(a)
}
```

Here is the implementation of `Len()` for the `S2slice` data type. It is usually that simple.

```markup
// What field to use when comparing
func (a S2slice) Less(i, j int) bool {
    return a[i].F3.F1 < a[j].F3.F1
}
```

Here is the implementation of `Less()` for the `S2slice` data type. This method defines the way elements get sorted. In this case, by using a field of the embedded data structure (`F3.F1`).

```markup
func (a S2slice) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}
```

This is the implementation of the `Swap()` type method that defines the way to swap slice elements during sorting. It is usually that simple.

```markup
func main() {
    data := []S2{
        S2{1, "One", S1{1, "S1_1", 10}},
        S2{2, "Two", S1{2, "S1_1", 20}},
        S2{-1, "Two", S1{-1, "S1_1", -20}},
    }
    fmt.Println("Before:", data)
    sort.Sort(S2slice(data))
    fmt.Println("After:", data)
    // Reverse sorting works automatically
    sort.Sort(sort.Reverse(S2slice(data)))
    fmt.Println("Reverse:", data)
}
```

Once you have implemented `sort.Interface`, you'll see that `sort.Reverse()`, which is used for reverse sorting your slice, works automatically.

Running `sort.go` generates the following output:

```markup
$ go run sort.go
Before: [{1 One {1 S1_1 10}} {2 Two {2 S1_1 20}} {-1 Two {-1 S1_1 -20}}]
After: [{-1 Two {-1 S1_1 -20}} {1 One {1 S1_1 10}} {2 Two {2 S1_1 20}}]
Reverse: [{2 Two {2 S1_1 20}} {1 One {1 S1_1 10}} {-1 Two {-1 S1_1 -20}}]
```

The first line shows the elements of the slice as initially stored. The second line shows the sorted version whereas the last line shows the reverse sorted version.

## The empty interface

As mentioned before, the **empty interface** is defined as just `interface{}` and is already implemented by **all data types**. Therefore, variables of any data type can be put in the place of a parameter of the empty interface data type. Therefore, a function with an `interface{}` parameter can accept variables of any data type in this place. However, if you intend to work with `interface{}` function parameters without examining their data type inside the function, you should process them with statements that work on all data types, otherwise your code may crash or misbehave.

The program that follows defines two structures named `S1` and `S2` but just a single function named `Print()` for printing both of them. This is allowed because `Print()` requires an `interface{}` parameter that can accept both `S1` and `S2` variables. The `fmt.Println(s)` statement inside `Print()` can work with both `S1` and `S2`.

If you create a function that accepts one or more `interface{}` parameters and you run a statement that can only be applied to a limited number of data types, things will not work out well. As an example, not all `interface{}` parameters can be multiplied by `5` or be used in `fmt.Printf()` with the `%d` control string.

The source code of `empty.go` is as follows:

```markup
package main
import "fmt"
type S1 struct {
    F1 int
    F2 string
}
type S2 struct {
    F1 int
    F2 S1
}
func Print(s interface{}) {
    fmt.Println(s)
}
func main() {
    v1 := S1{10, "Hello"}
    v2 := S2{F1: -1, F2: v1}
    Print(v1)
    Print(v2)
```

Although `v1` and `v2` are of different data types, `Print()` can work with both of them.

```markup
    // Printing an integer
    Print(123)
    // Printing a string
    Print("Go is the best!")
}
```

`Print()` can also work with integers and strings.

Running `empty.go` produces the following output:

```markup
{10 Hello}
{-1 {10 Hello}}
123
Go is the best!
```

Using the empty interface is easy as soon as you realize that you can pass any type of variable in the place of an `interface{}` parameter and you can return any data type as an `interface{}` return value. However, with great power comes great responsibility—you should be very careful with `interface{}` parameters and their return values, because in order to use their real values you have to be sure about their underlying data type. We'll discuss this in the next section.

## Type assertions and type switches

A **type assertion** is a mechanism for working with the underlying concrete value of an interface. This mainly happens because interfaces are virtual data types without their own values—interfaces just define behavior and do not hold data of their own. But what happens when you do not know the data type before attempting a type assertion? How can you differentiate between the supported data types and the unsupported ones? How can you choose a different action for each supported data type? The answer is by using _type switches_. **Type switches** use `switch` blocks for data types and allow you to differentiate between type assertion values, which are data types, and process each data type the way you want. On the other hand, in order to use the empty interface in type switches, you need to use **type assertions**.

You can have type switches for all kinds of interfaces and data types in general.

Therefore, the real work begins once you enter the function, because this is where you need to define the supported data types and the actions that take place for each supported data type.

Type assertions use the `x.(T)` notation, where `x` is an interface type and `T` is a type, and help you extract the value that is hidden behind the empty interface. For a type assertion to work, `x` should not be `nil` and the dynamic type of `x` should be identical to the `T` type.

The following code can be found as `typeSwitch.go`:

```markup
package main
import "fmt"
type Secret struct {
    SecretValue string
}
type Entry struct {
    F1 int
    F2 string
    F3 Secret
}
func Teststruct(x interface{}) {
    // type switch
    switch T := x.(type) {
    case Secret:
        fmt.Println("Secret type")
    case Entry:
        fmt.Println("Entry type")
    default:
        fmt.Printf("Not supported type: %T\n", T)
    }
}
```

This is a type switch that supports the `Secret` and `Entry` data types.

```markup
func Learn(x interface{}) {
    switch T := x.(type) {
    default:
        fmt.Printf("Data type: %T\n", T)
    }
}
```

The `Learn()` function prints the data type of its input parameter.

```markup
func main() {
    A := Entry{100, "F2", Secret{"myPassword"}}
    Teststruct(A)
    Teststruct(A.F3)
    Teststruct("A string")
    Learn(12.23)
    Learn('€')
}
```

The last part of the code calls the desired functions to explore variable `A`. Running `typeSwitch.go` produces the following output:

```markup
$ go run typeSwitch.go
Entry type
Secret type
Not supported type: string
Data type: float64
Data type: int32
```

As you can see, we have managed to execute different code based on the data type of the variable passed to `TestStruct()` and `Learn()`.

Strictly speaking, type assertions allow you to perform two main tasks:

-   Checking whether an interface value keeps a particular type. When used this way, a type assertion returns two values: the underlying value and a `bool` value. The underlying value is what you might want to use. However, it is the value of the `bool` variable that tells you whether the type assertion was successful or not and therefore whether you can use the underlying value or not. Checking whether a variable named `aVar` is of the `int` type requires the use of the `aVar.(int)` notation, which returns two values. If successful, it returns the real `int` value of `aVar` and `true`. Otherwise, it returns `false` as the second value, which means that the type assertion was not successful and that the real value could not be extracted.
-   Using the concrete value stored in an interface or assigning it to a new variable. This means that if there is a `float64` variable in an interface, a type assertion allows you to get that value.

The functionality offered by the `reflect` package helps Go identify the underlying data type and the real value of an `interface{}` variable.

So far, we have seen a variation of the first case where we extract the data type stored in an empty interface variable. Now, we are going to learn how to extract the real value stored in an empty interface variable.

As explained, trying to extract the concrete value from an interface using a type assertion can have two outcomes:

-   If you use the correct concrete data type, you get the underlying value without any issues
-   If you use an incorrect concrete data type, your program will panic

All these are illustrated in `assertions.go`, which contains the next code as well as lots of code comments that explain the process.

```markup
package main
import (
    "fmt"
)
func returnNumber() interface{} {
    return 12
}
func main() {
    anInt := returnNumber()
```

The `returnNumber()` function returns an `int` value that is wrapped in an empty interface.

```markup
    number := anInt.(int)
    number++
    fmt.Println(number)
```

In the previous code, we get the `int` value wrapped in an empty interface variable (`anInt`).

```markup
    // The next statement would fail because there
    // is no type assertion to get the value:
    // anInt++
    // The next statement fails but the failure is under 
    // control because of the ok bool variable that tells
    // whether the type assertion is successful or not
    value, ok := anInt.(int64)
    if ok {
        fmt.Println("Type assertion successful: ", value)
    } else {
        fmt.Println("Type assertion failed!")
    }
    // The next statement is successful but 
    // dangerous because it does not make sure that
    // the type assertion is successful.
    // It just happens to be successful
    i := anInt.(int)
    fmt.Println("i:", i)
    // The following will PANIC because anInt is not bool
    _ = anInt.(bool)
}
```

The last statement panics the program because the `anInt` variable does not hold a `bool` value. Running `assertions.go` generates the next output:

```markup
$ go run assertions.go
13
Type assertion failed!
i: 12
panic: interface conversion: interface {} is int, not bool
goroutine 1 [running]:
main.main()
        /Users/mtsouk/Desktop/mGo3rd/code/ch04/assertions.go:39 +0x192
```

The reason for the panic is written onscreen: `panic: interface conversion: interface {} is int, not bool`. What else can the Go compiler do to help you?

Next we discuss the `map[string]interface{}` map and its use.

## The map\[string\]interface{} map

You have a utility that processes its command-line arguments; if everything goes as expected, then you get the supported types of command-line arguments and everything goes smoothly. But what happens when something unexpected occurs? In that case, the `map[string]interface{}` map is here to help and this subsection shows how!

Remember that the biggest advantage you get from using a `map[string]interface{}` map or any map that stores an `interface{}` value in general, is that you still have your data in its original state and data type. If you use `map[string]string` instead, or anything similar, then any data you have is going to be converted into a `string`, which means that you are going to lose information about the original data type and the structure of the data you are storing in the map.

Nowadays, web services work by exchanging JSON records. If you get a JSON record in a supported format, then you can process it as expected and everything will be fine. However, there are times when you might get an erroneous record or a record in an unsupported JSON format. In these cases, using `map[string]interface{}` for storing these unknown JSON records (_arbitrary data_) is a good choice because `map[string]interface{}` is good at storing JSON records of an unknown type. We are going to illustrate that using a utility named `mapEmpty.go` that processes arbitrary JSON records given as command-line arguments. We process the input JSON record in two ways that are similar but not identical. There is no real difference between the `exploreMap()` and `typeSwitch()` functions apart from the fact that `typeSwitch()` generates a much richer output. The code of `mapEmpty.go` is as follows:

```markup
package main
import (
    "encoding/json"
    "fmt"
    "os"
)
var JSONrecord = `{
    "Flag": true,
    "Array": ["a","b","c"],
    "Entity": {
      "a1": "b1",
      "a2": "b2",
      "Value": -456,
      "Null": null
    },
    "Message": "Hello Go!"
  }`
```

This global variable holds the default value of `JSONrecord`, in case there is no user input.

```markup
func typeSwitch(m map[string]interface{}) {
    for k, v := range m {
        switch c := v.(type) {
        case string:
            fmt.Println("Is a string!", k, c)
        case float64:
            fmt.Println("Is a float64!", k, c)
        case bool:
            fmt.Println("Is a Boolean!", k, c)
        case map[string]interface{}:
            fmt.Println("Is a map!", k, c)
            typeSwitch(v.(map[string]interface{}))
        default:
            fmt.Printf("...Is %v: %T!\n", k, c)
        }
    }
    return
}
```

The `typeSwitch()` function uses a type switch for differentiating between the values in its input map. If a map is found, then we _recursively_ call `typeSwitch()` on the new map in order to examine it even more.

The `for` loop allows you to examine all the elements of the `map[string]interface{}` map.

```markup
func exploreMap(m map[string]interface{}) {
    for k, v := range m {
        embMap, ok := v.(map[string]interface{})
        // If it is a map, explore deeper
        if ok {
            fmt.Printf("{\"%v\": \n", k)
            exploreMap(embMap)
            fmt.Printf("}\n")
        } else {
            fmt.Printf("%v: %v\n", k, v)
        }
    }
}
```

The `exploreMap()` function inspects the contents of its input map. If a map is found, then we call `exploreMap()` on the new map _recursively_ in order to examine it on its own.

```markup
func main() {
    if len(os.Args) == 1 {
        fmt.Println("*** Using default JSON record.")
    } else {
        JSONrecord = os.Args[1]
    }
    JSONMap := make(map[string]interface{})
    err := json.Unmarshal([]byte(JSONrecord), &JSONMap)
```

As you will learn in _Chapter 6_, _Telling a UNIX System What To Do_, `json.Unmarshal()` processes JSON data and converts it into a Go value. Although this value is usually a Go structure, in this case we are using a map as specified by the `map[string]interface{}` variable. Strictly speaking, the second parameter of `json.Unmarshal()` is of the empty interface data type, which means that its data type can be anything.

```markup
    if err != nil {
        fmt.Println(err)
        return
    }
    exploreMap(JSONMap)
    typeSwitch(JSONMap)
}
```

`map[string]interface{}` is extremely handy for storing JSON records when you do not know their schema in advance. In other words, `map[string]interface{}` is good at storing arbitrary JSON data of unknown schema.

Running `mapEmpty.go` produces the following output:

```markup
$ go run mapEmpty.go 
*** Using default JSON record.
Message: Hello Go!
Flag: true
Array: [a b c]
{"Entity": 
Value: -456
Null: <nil>
a1: b1
a2: b2
}
Is a Boolean! Flag true
...Is Array: []interface {}!
Is a map! Entity map[Null:<nil> Value:-456 a1:b1 a2:b2]
Is a string! a2 b2
Is a float64! Value -456
...Is Null: <nil>!
Is a string! a1 b1
Is a string! Message Hello Go!
$ go run mapEmpty.go '{"Array": [3, 4], "Null": null, "String": "Hello Go!"}'
Array: [3 4]
Null: <nil>
String: Hello Go!
...Is Array: []interface {}!
...Is Null: <nil>!
Is a string! String Hello Go!
$ go run mapEmpty.go '{"Array":"Error"' 
unexpected end of JSON input
```

The first run is without any command-line parameters, which means that it uses the default value of `JSONrecord` and therefore outputs the hardcoded data. The other two executions use user data. First, valid data, and then data that does not represent a valid JSON record. The error message in the third execution is generated by `json.Unmarshal()` as it cannot understand the schema of the JSON record.

## The error data type

As promised, we are revisiting the `error` data type because it is an interface defined as follows:

```markup
type error interface {
    Error() string
}
```

So, in order to satisfy the `error` interface you just need to implement the `Error() string` type method. This does not change the way we use errors to find out whether the execution of a function or method was successful or not but shows how important interfaces are in Go as they are being used transparently all the time. However, the crucial question is _when_ you should implement the `error` interface on your own instead of using the default one. The answer to that question is when you want to give more context to an error condition.

Now, let us talk about the `error` interface in a more practical situation. When there is nothing more to read from a file, Go returns an `io.EOF` error, which, strictly speaking, is not an error condition but a logical part of reading a file. If a file is totally empty, you still get `io.EOF` when you try to read it. However, this might cause problems in some situations and you might need to have a way of differentiating between a totally empty file and a file that has been read fully and there is nothing more to read. One way of dealing with that issue is with the help of the `error` interface.

The code example that is presented here is connected to File I/O. Putting it here might generate some questions about reading files in Go—however, I feel that this is the appropriate place to put it because it is connected to errors and error handling more than it is connected to file reading.

The code of `errorInt.go` without the `package` and `import` blocks is as follows:

```markup
type emptyFile struct {
    Ended bool
    Read  int
}
```

This is a new data type that is used in the program.

```markup
// Implement error interface
func (e emptyFile) Error() string {
    return fmt.Sprintf("Ended with io.EOF (%t) but read (%d) bytes", e.Ended, e.Read)
}
```

This is the implementation of the `error` interface for `emptyFile`.

```markup
// Check values
func isFileEmpty(e error) bool {
    // Type assertion
    v, ok := e.(emptyFile)
```

This is a type assertion for getting an `emptyFile` structure from the `error` variable.

```markup
    if ok {
        if v.Read == 0 && v.Ended == true {
            return true
        }
    }
    return false
}
```

This is a method for checking whether a file is empty or not. The `if` condition translates as: if you have read 0 bytes (`v.Read == 0`) and you have reached the end of the file (`v.Ended == true`), then the file is empty.

If you are dealing with multiple `error` variables, you should add a type switch to the `isFileEmpty()` function after the type assertion.

```markup
func readFile(file string) error {
    var err error
    fd, err := os.Open(file)
    if err != nil {
        return err
    }
    defer fd.Close()
    reader := bufio.NewReader(fd)
    n := 0
    for {
        line, err := reader.ReadString('\n')
        n += len(line)
```

We read the input file line by line—you are going to learn more about File I/O in _Chapter 6__, Telling a UNIX System What to Do_.

```markup
        if err == io.EOF {
            // End of File: nothing more to read
            if n == 0 {
                return emptyFile{true, n}
            }
```

If we have reached the end of a file (`io.EOF`) and we have read `0` characters, then we are dealing with an empty file. This kind of context is added to the `emptyFile` structure and returned as an `error` value.

```markup
            break
        } else if err != nil {
            return err
        }
    }
    return nil
}
func main() {
    flag.Parse()
    if len(flag.Args()) == 0 {
        fmt.Println("usage: errorInt <file1> [<file2> ...]")
        return
    }
    for _, file := range flag.Args() {
        err := readFile(file)
        if isFileEmpty(err) {
            fmt.Println(file, err)
```

This is where we check the error message of the `readFile()` function. The order we do the checking in is important because only the first match is executed. This means that we have to go from more specific cases to more generic conditions.

```markup
        } else if err != nil {
            fmt.Println(file, err)
        } else {
            fmt.Println(file, "is OK.")
        }
    }
}
```

Running `errorInt.go` produces the next output:

```markup
$ go run errorInt.go /etc/hosts /tmp/doesNotExist /tmp/empty /tmp /tmp/Empty.txt
/etc/hosts is OK.
/tmp/doesNotExist open /tmp/doesNotExist: no such file or directory
/tmp/empty open /tmp/empty: permission denied
/tmp read /tmp: is a directory
/tmp/Empty.txt Ended with io.EOF (true) but read (0) bytes
```

The first file (`/etc/hosts`) was read without any issues, whereas the second file (`/tmp/doesNotExist`) could not be found. The third file (`/tmp/empty`) was there but we did not have the required file permissions to read it, whereas the fourth file (`/tmp`) was in reality a directory. The last file (`/tmp/Empty.txt`) was there but was empty, which is the error situation that we wanted to catch.

## Writing your own interfaces

After learning about using existing interfaces, we will write another command-line utility that sorts 3D shapes according to their volumes. This task requires learning the following tasks:

-   Creating new interfaces
-   Combining existing interfaces
-   Implementing `sort.Interface` for 3D shapes

Creating your own interfaces is easy. For reasons of simplicity, we include our own interface in the `main` package. However, this is rarely the case as we usually want to share our interfaces, which means that interfaces are usually included in Go packages other than `main`.

The following code excerpt defines a new interface:

```markup
type Shape2D interface {
    Perimeter() float64
}
```

This interface has the following properties:

-   It is called `Shape2D`
-   It requires the implementation of a single method named `Perimeter()` that returns a `float64` value

Apart from being user-defined, there is nothing special about that interface compared to the built-in Go interfaces—you can use it as you do all other existing interfaces. So, in order for a data type to satisfy the `Shape2D` interface, it needs to implement a type method named `Perimeter()` that returns a `float64` value.

### Using a Go interface

The code that follows presents the simplest way of using an interface, which is by calling its method directly, as if it was a function, to get a result. Although this is allowed, it is rarely the case as we usually create functions that accept interface parameters in order for these functions to be able to work with multiple data types.

The code uses a handy technique for quickly finding out whether a given variable is of a given data type that was presented earlier in `assertions.go`. In this case, we examine whether a variable is of the `Shape2D` interface by using the `interface{}(a).(Shape2D)` notation, where `a` is the variable that is being examined and `Shape2D` is the data type against the variable being checked.

The next program is called `Shape2D.go`—its most interesting parts are the following:

```markup
type Shape2D interface {
    Perimeter() float64
}
```

This is the definition of the `Shape2D` interface that requires the implementation of the `Perimeter()` type method.

```markup
type circle struct {
    R float64
}
func (c circle) Perimeter() float64 {
    return 2 * math.Pi * c.R
}
```

This is where the `circle` type implements the `Shape2D` interface with the implementation of the `Perimeter()` type method.

```markup
func main() {
    a := circle{R: 1.5}
    fmt.Printf("R %.2f -> Perimeter %.3f \n", a.R, a.Perimeter())
    _, ok := interface{}(a).(Shape2D)
    if ok {
        fmt.Println("a is a Shape2D!")
    }
}
```

As stated before, the `interface{}(a).(Shape2D)` notation checks whether the `a` variable satisfies the `Shape2D` interface without using its underlying value (`circle{R: 1.5}`).

Running `Shape2D.go` creates the next output:

```markup
R 1.50 -> Perimeter 9.425 
a is a Shape2D!
```

### Implementing sort.Interface for 3D shapes

In this section we will create a utility for sorting various 3D shapes based on their volume, which clearly illustrates the power and versatility of Go interfaces. This time, we will use a single slice for storing all kinds of structures that all satisfy a given interface. The fact that Go considers interfaces as data types allows us to create slices with elements that satisfy a given interface without getting any error messages.

This kind of scenario can be useful in various cases because it illustrates how to store elements with different data types that all satisfy a common interface on the same slice and how to sort them using `sort.Interface`. Put simply, the presented utility sorts different structures with different numbers and names of fields that all share a common behavior through an interface implementation. The dimensions of the shapes are created using random numbers, which means that each time you execute the utility, you get a different output.

The name of the interface is `Shape3D` and requires the implementation of the `Vol() float64` type method. This interface is satisfied by the `Cube`, `Cuboid`, and `Sphere` data types. The `sort.Interface` interface is implemented for the `shapes` data type, which is defined as a slice of `Shape3D` elements.

All floating-point numbers are randomly generated using the `rF64(min, max float64) float64` function. As floating-point numbers have a lot of decimal points, printing is implemented using a separate function named `PrintShapes()` that uses an `fmt.Printf("%.2f ", v)` statement to specify the number of decimal points that are displayed onscreen—in this case, we print the first two decimal points of each floating-point value.

As you might recall, once you have implemented `sort.Interface`, you can also sort your data in reverse order using `sort.Reverse()`.

Type the following code on your favorite editor and save it as `sortShapes.go`. The code illustrates how to sort 3D shapes based on their volume.

```markup
package main
import (
    "fmt"
    "math"
    "math/rand"
    "sort"
    "time"
)
const min = 1
const max = 5
func rF64(min, max float64) float64 {
    return min + rand.Float64()*(max-min)
}
```

The `rF64()` function generates `float64` random values.

```markup
type Shape3D interface {
    Vol() float64
}
```

The definition of the `Shape3D` interface.

```markup
type Cube struct {
    x float64
}
type Cuboid struct {
    x float64
    y float64
    z float64
}
type Sphere struct {
    r float64
}
func (c Cube) Vol() float64 {
    return c.x * c.x * c.x
}
```

`Cube` implementing the `Shape3D` interface.

```markup
func (c Cuboid) Vol() float64 {
    return c.x * c.y * c.z
}
```

`Cuboid` implementing the `Shape3D` interface.

```markup
func (c Sphere) Vol() float64 {
    return 4 / 3 * math.Pi * c.r * c.r * c.r
}
```

`Sphere` implementing the `Shape3D` interface.

```markup
type shapes []Shape3D
```

This is the data type that uses `sort.Interface`.

```markup
// Implementing sort.Interface
func (a shapes) Len() int {
    return len(a)
}
func (a shapes) Less(i, j int) bool {
    return a[i].Vol() < a[j].Vol()
}
func (a shapes) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}
```

The previous three functions implement `sort.Interface`.

```markup
func PrintShapes(a shapes) {
    for _, v := range a {
        switch v.(type) {
        case Cube:
            fmt.Printf("Cube: volume %.2f\n", v.Vol())
        case Cuboid:
            fmt.Printf("Cuboid: volume %.2f\n", v.Vol())
        case Sphere:
            fmt.Printf("Sphere: volume %.2f\n", v.Vol())
        default:
            fmt.Println("Unknown data type!")
        }
    }
    fmt.Println()
}
func main() {
    data := shapes{}
    rand.Seed(time.Now().Unix())
```

The `PrintShapes()` function is used for customizing the output.

```markup
    for i := 0; i < 3; i++ {
        cube := Cube{rF64(min, max)}
        cuboid := Cuboid{rF64(min, max), rF64(min, max), rF64(min, max)}
        sphere := Sphere{rF64(min, max)}
        data = append(data, cube)
        data = append(data, cuboid)
        data = append(data, sphere)
    }
    PrintShapes(data)
    // Sorting
    sort.Sort(shapes(data))
    PrintShapes(data)
    // Reverse sorting
    sort.Sort(sort.Reverse(shapes(data)))
    PrintShapes(data)
}
```

The following code produces shapes with randomly generated dimensions using the `rF64()` function.

Running `sortShapes.go` produces the following output:

```markup
Cube: volume 105.27
Cuboid: volume 34.88
Sphere: volume 212.31
Cube: volume 55.76
Cuboid: volume 28.84
Sphere: volume 46.50
Cube: volume 52.41
Cuboid: volume 36.90
Sphere: volume 257.03
```

This is the unsorted output of the program:

```markup
Cuboid: volume 28.84
Cuboid: volume 34.88
Cuboid: volume 36.90
Sphere: volume 46.50
Cube: volume 52.41
...
Sphere: volume 257.03
```

This is the sorted output of the program from smaller to bigger shapes:

```markup
Sphere: volume 257.03
...
Cuboid: volume 28.84
```

This is the reversed sorted output of the program from bigger to smaller shapes.

The next section shows a technique for differentiating between two CSV file formats in your programs.

Just Imagine

# Working with two different CSV file formats

In this section we are going to implement a separate command-line utility that works with two different CSV formats. The reason we are doing this is that there are times when you will need your utilities to be able to work with multiple data formats.

Remember that the records of each CSV format are stored using their own Go structure under a different variable name. As a result, we need to implement `sort.Interface` for both CSV formats and therefore for both slice variables.

The two supported formats are the following:

-   _Format 1_: name, surname, telephone number, time of last access
-   _Format 2_: name, surname, area code, telephone number, time of last access

As the two CSV formats that are going to be used have a different number of fields, the utility determines the format that is being used by the number of fields found in the first record that was read and acts accordingly. After that, the data will be sorted using `sort.Sort()`—the data type of the slice that keeps the data helps Go determine the sort implementation that is going to be used without any help from the developer.

The main benefit you get from functions that work with empty interface variables is that you can add support for additional data types easily at a later time without the need to implement additional functions and without breaking existing code.

What follows is the implementation of the most important functions of the utility beginning with `readCSVFile()` because the logic of the utility is found in the `readCSVFile()` function.

```markup
func readCSVFile(filepath string) error {
.
.
.
```

The code that has to do with reading the input file and making sure that it exists is omitted for brevity.

```markup
    var firstLine bool = true
    var format1 = true
```

The first line of the CSV file determines its format—therefore, we need a flag variable for specifying whether we are dealing with the first line (`firstLine`) or not. Additionally, we need a second variable for specifying the format we are working with (`format1` is that variable).

```markup
    for _, line := range lines {
        if firstLine {
            if len(line) == 4 {
                format1 = true
            } else if len(line) == 5 {
                format1 = false
```

The first format has four fields whereas the second format has five fields.

```markup
            } else {
                return errors.New("Unknown File Format!")
            }
            firstLine = false
        }
```

If the first line of the CSV file has neither four nor five fields, then we have an error, and the function returns with a custom error message.

```markup
        if format1 {
            if len(line) == 4 {
                temp := F1{
                    Name:       line[0],
                    Surname:    line[1],
                    Tel:        line[2],
                    LastAccess: line[3],
                }
                d1 = append(d1, temp)
            }
```

If we are working with `format1`, we add data to the `d1` global variable.

```markup
        } else {
            if len(line) == 5 {
                temp := F2{
                    Name:       line[0],
                    Surname:    line[1],
                    Areacode:   line[2],
                    Tel:        line[3],
                    LastAccess: line[4],
                }
                d2 = append(d2, temp)
            }
```

If we are working with `format2`, we add data to the `d2` global variable.

```markup
        }
    }
    return nil
}
```

The `sortData()` function accepts an empty interface parameter. The code of the function determines the data type of the slice that is passed as an empty interface to that function using a type switch. After that, a type assertion allows you to use the actual data stored under the empty interface parameter. Its full implementation is as follows:

```markup
func sortData(data interface{}) {
    // type switch
    switch T := data.(type) {
    case Course1:
        d := data.(Course1)
        sort.Sort(Course1(d))
        list(d)
    case Course2:
        d := data.(Course2)
        sort.Sort(Course2(d))
        list(d)
    default:
        fmt.Printf("Not supported type: %T\n", T)
    }
}
```

The type switch does the job of determining the data type we are working with, which can be either `Course1` or `Course2`. If you want to look at the implementation of `sort.Interface`, you should view the `sortCSV.go` source code file.

Lastly, `list()` prints the data of the data variable used using the technique found in `sortData()`. Although the code that handles `Course1` and `Course2` is the same as in `sortData()`, you still need a type assertion to get the data from the empty interface variable.

```markup
func list(d interface{}) {
    switch T := d.(type) {
    case Course1:
        data := d.(Course1)
        for _, v := range data {
            fmt.Println(v)
        }
    case Course2:
        data := d.(Course2)
        for _, v := range data {
            fmt.Println(v)
        }
    default:
        fmt.Printf("Not supported type: %T\n", T)
    }
}
```

Running `sortCSV.go` produces the following kind of output:

```markup
$ go run sortCSV.go /tmp/csv.file
{Jane Doe 0800123456 1609310777}
{Dimitris Tsoukalos 2109416871 1609310731}
{Dimitris Tsoukalos 2109416971 1609310734}
{Mihalis Tsoukalos 2109416471 1609310706}
{Mihalis Tsoukalos 2109416571 1609310717}
```

The program correctly found out the format of `/tmp/csv.file` and worked with it even though it supports two CSV formats. Trying to work with an unsupported format generates the next output:

```markup
$ go run sortCSV.go /tmp/differentFormat.csv
Unknown File Format!
```

This means that the code successfully understands that we are dealing with an unsupported format.

The next section explores the limited object-oriented capabilities of Go.

Just Imagine

# Object-oriented programming in Go

As Go does not support all object-oriented features, it cannot replace an object-oriented programming language fully. However, it can **mimic some object-oriented concepts**.

First of all, a Go structure with its type methods is like an object with its methods. Second, interfaces are like abstract data types that define behaviors and objects of the same class, which is similar to **polymorphism**. Third, Go supports **encapsulation**, which means it supports hiding data and functions from the user by making them private to the structure and the current Go package. Lastly, combining interfaces and structures is like **composition** in object-oriented terminology.

If you really want to develop applications using the object-oriented methodology, then choosing Go might not be your best option. As I am not really into **Java**, I would suggest looking at **C++** or **Python** instead. The general rule here is to choose the best tool for your job.

You have already seen some of these points earlier in this chapter—the next chapter discusses how to define private fields and functions. The example that follows, which is named `objO.go`, illustrates composition and polymorphism as well as embedding an anonymous structure into an existing one to get all its fields.

```markup
package main
import (
    "fmt"
)
type IntA interface {
    foo()
}
type IntB interface {
    bar()
}
type IntC interface {
    IntA
    IntB
}
```

The `IntC` interface combines interfaces `IntA` and `IntB`. If you implement `IntA` and `IntB` for a data type, then this data type implicitly satisfies `IntC`.

```markup
func processA(s IntA) {
    fmt.Printf("%T\n", s)
}
```

This function works with data types that satisfy the `IntA` interface.

```markup
type a struct {
    XX int
    YY int
}
// Satisfying IntA
func (varC c) foo() {
    fmt.Println("Foo Processing", varC)
}
```

Structure `c` satisfying `IntA` as it implements `foo()`.

```markup
// Satisfying IntB
func (varC c) bar() {
    fmt.Println("Bar Processing", varC)
}
```

Structure `c` satisfying `IntB`. As structure `c` satisfies both `IntA` and `IntB`, it implicitly satisfies `IntC`, which is a composition of the `IntA` and `IntB` interfaces.

```markup
type b struct {
    AA string
    XX int
}
// Structure c has two fields
type c struct {
    A a
    B b
}
```

This structure has two fields named `A` and `B` that are of the `a` and `b` data types, respectively.

```markup
// Structure compose gets the fields of structure a
type compose struct {
    field1 int
    a
}
```

This new structure uses an anonymous structure (`a`), which means that it gets the fields of that anonymous structure.

```markup
// Different structures can have methods with the same name
func (A a) A() {
    fmt.Println("Function A() for A")
}
func (B b) A() {
    fmt.Println("Function A() for B")
}
func main() {
    var iC c = c{a{120, 12}, b{"-12", -12}}
```

Here we define a `c` variable that is composed of an `a` structure and a `b` structure.

```markup
    iC.A.A()
    iC.B.A()
```

Here we access a method of the `a` structure (`A.A()`) and a method of the `b` structure (`B.A()`).

```markup
    // The following will not work
    // iComp := compose{field1: 123, a{456, 789}}
    // iComp := compose{field1: 123, XX: 456, YY: 789}
    iComp := compose{123, a{456, 789}}
    fmt.Println(iComp.XX, iComp.YY, iComp.field1)
```

When using an anonymous structure inside another structure, as we do with `a{456, 789}`, you can access the fields of the anonymous structure, which is the `a{456, 789}` structure, directly as `iComp.XX` and `iComp.YY`.

```markup
    iC.bar()
    processA(iC)
}
```

Although `processA()` works with `IntA` variables, it can also work with `IntC` variables because the `IntC` interface satisfies `IntA`!

All the code in `objO.go` is pretty simplistic compared to the code of a real object-oriented programming language that supports abstract classes and inheritance. However, it is more than adequate for generating types and elements with a structure in them, as well as for having different data types with the same method names.

Running `objO.go` produces the next output:

```markup
$ go run objO.go
Function A() for A
Function A() for B
456 789 123
Bar Processing {{120 12} {-12 -12}}
main.c
```

The first two lines of the output show that two different structures can have a method with the same name. The third line proves that when using an anonymous structure inside one other structure, you can access the fields of the anonymous structure directly. The fourth line is the output of the `iC.bar()` call, where `iC` is a `c` variable accessing a method from the `IntB` interface. The last line is the output of `processA(iC)` that requires an `IntA` parameter and prints the real data type of its parameter, which in this case is `main.c`.

Evidently, although Go is not an object-oriented programming language, it can mimic some of the characteristics of object-oriented programming. Moving on, the last section of this chapter is about updating the phone book application by reading an environment variable and sorting its output.

Just Imagine

# Updating the phone book application

The functionality that is added to this new version of the phone book utility is the following:

-   The CSV file path can be optionally given as an environment variable named `PHONEBOOK`
-   The `list` command sorts the output based on the surname field

Although we could have given the path of the CSV file as a command-line argument instead of the value of an environment variable, it would have complicated the code, especially if that argument was made optional. More advanced Go packages such as `viper`, which is presented in _Chapter 6_, _Telling a UNIX System What to Do_, simplify the process of parsing command-line arguments with the use of command-line options such as `-f` followed by a file path or `--filepath`.

The current default value of `CSVFILE` is set to my home directory on a macOS Big Sur machine—you should change that default value to fit your needs or use a proper value for the `PHONEBOOK` environment variable.

Last, if the `PHONEBOOK` environment variable is not set, then the utility uses a default value for the CSV file path. Generally speaking, not having to recompile your software for user-defined data is considered a good practice.

## Setting up the value of the CSV file

The value of the CSV file is set in the `setCSVFILE()` function, which is defined as follows:

```markup
func setCSVFILE() error {
    filepath := os.Getenv("PHONEBOOK")
    if filepath != "" {
        CSVFILE = filepath
    }
```

Here is where we read the `PHONEBOOK` environment variable. The rest of the code is about making sure that we can use that file path or the default one if `PHONEBOOK` is not set.

```markup
    _, err := os.Stat(CSVFILE)
    if err != nil {
        fmt.Println("Creating", CSVFILE)
        f, err := os.Create(CSVFILE)
        if err != nil {
            f.Close()
            return err
        }
        f.Close()
    }
```

If the specified file does not exist, it is created using `os.Create()`.

```markup
    fileInfo, err := os.Stat(CSVFILE)
    mode := fileInfo.Mode()
    if !mode.IsRegular() {
        return fmt.Errorf("%s not a regular file", CSVFILE)
    }
```

Then, we make sure that the specified path belongs to a regular file that can be used for saving data.

```markup
    return nil
}
```

In order to simplify the implementation of the `main()` function, we moved the code related to the existence of and access to the CSV file path to `setCSVFILE()`.

The first time we set the `PHONEBOOK` environment variable and executed the phone book application, we got the following output—you should get something similar.

```markup
$ export PHONEBOOK="/tmp/csv.file"
$ go run phoneCourse.go list        
Creating /tmp/csv.file
```

As `/tmp/csv.file` does not exist, `phoneCourse.go` creates it from scratch. This verifies that the Go code of the `setCSVFILE()` function works as expected.

Now that we know where to get and write our data, it is time to learn how to sort it using `sort.Interface`, which is the subject of the subsection that follows.

## Using the sort package

The first thing to decide when trying to sort data is the field that is going to be used for sorting. After that, we need to decide what we are going to do when two or more records have the same value in the main field used for sorting.

The code related to sorting using `sort.Interface` is the following:

```markup
type PhoneCourse []Entry
```

You need to have a separate data type—`sort.Interface` is implemented for this data type.

```markup
var data = PhoneCourse{}
```

As you have a separate data type for implementing `sort.Interface`, the data type of the `data` variable needs to change and become `PhoneCourse`. Then `sort.Interface` is implemented for `PhoneCourse`.

```markup
// Implement sort.Interface
func (a PhoneCourse) Len() int {
    return len(a)
}
```

The `Len()` function has a standard implementation.

```markup
// First based on surname. If they have the same
// surname take into account the name.
func (a PhoneCourse) Less(i, j int) bool {
    if a[i].Surname == a[j].Surname {
        return a[i].Name < a[j].Name
    }
```

The `Less()` function is the place to define how you are going to sort the elements of the slice. What we say here is that if the entries that are compared, which are Go structures, have the same `Surname` field value, then compare these entries using their `Name` field values.

```markup
    return a[i].Surname < a[j].Surname
}
```

If the entries have different values in the `Surname` field, then compare them using the `Surname` field.

```markup
func (a PhoneCourse) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}
```

The `Swap()` function has a standard implementation. After implementing the desired interface, we need to tell our code to sort the data, which happens in the implementation of the `list()` function:

```markup
func list() {
    sort.Sort(PhoneCourse(data))
    for _, v := range data {
        fmt.Println(v)
    }
}
```

Now that we know how sorting is implemented, it is time to use the utility. First, we add some entries:

```markup
$ go run phoneCourse.go insert Mihalis Tsoukalos 2109416471
$ go run phoneCourse.go insert Mihalis Tsoukalos 2109416571
$ go run phoneCourse.go insert Dimitris Tsoukalos 2109416871
$ go run phoneCourse.go insert Dimitris Tsoukalos 2109416971
$ go run phoneCourse.go insert Jane Doe 0800123456
```

Last, we print the contents of the phone book using the `list` command:

```markup
$ go run phoneCourse.go list                      
{Jane Doe 0800123456 1609310777}
{Dimitris Tsoukalos 2109416871 1609310731}
{Dimitris Tsoukalos 2109416971 1609310734}
{Mihalis Tsoukalos 2109416471 1609310706}
{Mihalis Tsoukalos 2109416571 1609310717}
```

As `Dimitris` comes before `Mihalis` alphabetically, all relevant entries come first as well, which means that our sorting works as expected.

Just Imagine

# Exercises

-   Create a slice of structures using a structure that you created and sort the elements of the slice using a field from the structure
-   Integrate the functionality of `sortCSV.go` in `phonebook.go`
-   Add support for a `reverse` command to `phonebook.go` in order to list its entries in reverse order
-   Use the empty interface and a function that allows you to differentiate between two different structures that you create

Just Imagine

# Summary

In this chapter, we learned about _interfaces_, which are like contracts, and also about _type methods_, type assertion, and reflection. Although reflection is a very powerful Go feature, it might slow down your Go programs because it adds a layer of complexity at runtime. Furthermore, your Go programs could crash if you use reflection carelessly.

The last section of this chapter discussed writing Go code that follows the principles of object-oriented programming. If you are going to remember just one thing from this chapter, it should be that Go is not an object-oriented programming language, but it can mimic some of the functionality offered by object-programming languages, such as Java, Python, and C++.

The next chapter discusses Go packages, functions, and automation using GitHub and GitLab CI/CD systems.

Just Imagine

# Additional resources

-   The documentation of the `reflect` package: [https://golang.org/pkg/reflect/](https://golang.org/pkg/reflect/)
-   The documentation of the `sort` package: [https://golang.org/pkg/sort/](https://golang.org/pkg/sort/)
-   Working with errors in Go 1.13: [https://blog.golang.org/go1.13-errors](https://blog.golang.org/go1.13-errors)
-   The implementation of the `sort` package: [https://golang.org/src/sort/](https://golang.org/src/sort/)