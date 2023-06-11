# Go Generics

This chapter is about _generics_ and how to use the new syntax to write generic functions and define generic data types. Currently, generics are under development, but the official release is pretty close and we have a good idea of what features generics are going to have and how generics are going to work.

The new generics syntax is coming to Go 1.18, which, according to the Go development cycle, is going to be officially released in February 2022.

Let me make something clear from the beginning: you do not have to use Go generics if you do not want to and you can still write wonderful, efficient, maintainable, and correct software in Go! Additionally, the fact that you can use generics and support lots of data types, if not all available data types, does not mean that you should do that. Always **support the required data types**, no more, no less, but do not forget to keep an eye on the future of your data and the possibility of supporting data types that were not known at the time of writing your code.

This chapter covers:

-   Introducing generics
-   Constraints
-   Defining new data types with generics
-   Interfaces versus generics
-   Reflection versus generics
-   Concluding remarks: what does the future look like for Go developers?

Just Imagine

# Introducing generics

Generics are a feature that gives you the capability of not precisely specifying the data type of one or more function parameters, mainly because you want to make your functions as generic as possible. In other words, generics allow functions to process several data types without the need to write special code, as is the case with the empty interface or interfaces in general. However, when working with interfaces in Go, you have to write extra code to determine the data type of the interface variable you are working with, which is not the case with generics.

Let me begin by presenting a small code example that implements a function that clearly shows a case where generics can be handy and save you from having to write lots of code:

```markup
func PrintSlice[T any](s []T) {
    for _, v := range s {
        fmt.Println(v)
    }
}
```

So, what do we have here? There is a function named `PrintSlice()` that accepts a slice of any data type. This is denoted by the use of `[]T` in the function signature in combination with the `[T any]` part. The `[T any]` part tells the compiler that the data type `T` is going to be determined at execution time. We are also free to use multiple data types using the `[T, U, W any]` notation—after which we should use the `T, U, W` data types in the function signature.

The `any` keyword tells the compiler that there are **no constraints** about the data type of `T`. We are going to discuss constraints in a while—for now, just learn the syntax of generics.

Now imagine writing separate functions to implement the functionality of `PrintSlice()` for slices of integers, strings, floating-point numbers, complex values, and so on. So, we have found a profound case where using generics simplifies the code and our programming efforts. However, not all cases are so obvious, and we should be very careful about overusing `any`.

But what happens if you want to try generics before the official release? There is a solution, which is to visit [https://go2goplay.golang.org/](https://go2goplay.golang.org/) and place and run your code there.

The initial screen of [https://go2goplay.golang.org/](https://go2goplay.golang.org/) is presented in the next figure. Just like the regular Go Playground, the upper part is where you write the code, whereas the bottom part is where you get the results of your code or potential error messages after pressing the **Run** button:

![Graphical user interface, text, email
Description automatically generated](https://static.packt-cdn.com/products/9781801079310/graphics/Images/B17194_13_01.png)

Figure 13.1: The Go Playground for trying generics

The rest of this chapter is going to execute the Go code in [https://go2goplay.golang.org/](https://go2goplay.golang.org/). However, I am going to present the code as usual and the output of the code by copying and pasting it from [https://go2goplay.golang.org/](https://go2goplay.golang.org/).

The following (`hw.go`) is code that uses generics, to help you understand more about them before going into more advanced examples:

```markup
package main
import (
    "fmt"
)
func PrintSlice[T any](s []T) {
    for _, v := range s {
        fmt.Print(v, " ")
    }
    fmt.Println()
}
```

`PrintSlice()` is similar to the function that we saw earlier in this chapter. However, it prints the elements of each slice in the same line and prints a new line with the help of `fmt.Println()`.

```markup
func main() {
    PrintSlice([]int{1, 2, 3})
    PrintSlice([]string{"a", "b", "c"})
    PrintSlice([]float64{1.2, -2.33, 4.55})
}
```

Here we call `PrintSlice()` with three different data types: `int`, `string`, and `float64`. The Go compiler is not going to complain about that. Instead, it is going to execute the code as if we had three separate functions, one for each data type.

Therefore, running `hw.go` produces the next output:

```markup
1 2 3 
a b c 
1.2 -2.33 4.55
```

So, each slice is printed as expected using a single generic function.

With that information in mind, let us begin by discussing generics and constraints.

Just Imagine

# Constraints

Let us say that you have a function that works with generics that multiplies two numeric values. Should this function work with all data types? Can this function work with all data types? Can you multiply two strings or two structures? The solution for avoiding that kind of issue is the use of **constraints**.

Forget about multiplication for a while and think about something simpler. Let us say that we want to compare variables for equality—is there a way to tell Go that we only want to work with values that can be compared? Go 1.18 is going to come with predefined constraints—one of them is called `comparable` and includes data types that can be compared for equality or inequality.

The code of `allowed.go` illustrates the use of the `comparable` constraint:

```markup
package main
import (
    "fmt"
)
func Same[T comparable](a, b T) bool {
    if a == b {
        return true
    }
    return false
}
```

The `Same()` function uses the **predefined** `comparable` constraint instead of `any`. In reality, the `comparable` constraint is just a predefined interface that includes all data types that can be compared with `==` or `!=`. We do not have to write any extra code for checking our input as the function signature makes sure that we are going to deal with acceptable data types only.

```markup
func main() {
    fmt.Println("4 = 3 is", Same(4,3))
    fmt.Println("aa = aa is", Same("aa","aa"))
    fmt.Println("4.1 = 4.15 is", Same(4.1,4.15))
}
```

The `main()` function calls `Same()` three times and prints its results.

Running `allowed.go` produces the next output:

```markup
4 = 3 is false
aa = aa is true
4.1 = 4.15 is false
```

As only `Same("aa","aa")` is `true`, we get the respective output.

If you try to run a statement that includes `Same([]int{1,2},[]int{1,3})`, which tries to compare two slices, the Go Playground is going to generate the following error message:

```markup
type checking failed for main
prog.go2:19:31: []int does not satisfy comparable
```

This happens because we cannot directly compare two slices—this kind of functionality should be implemented manually.

The next subsection shows how to create your own constraints.

## Creating constraints

This subsection presents an example where we define the data types that are allowed to be passed as parameters to a generic function using an _interface_. The code of `numeric.go` is as follows:

```markup
package main
import (
    "fmt"
)
type Numeric interface {
    type int, int8, int16, int32, int64, float64
}
```

In here, we define a new interface called `Numeric` that specifies the list of supported data types. You can use any data type you want as long as it can be used with the generic function that you are going to implement. In this case, we could have added `string` or `uint` to the list of supported data types.

```markup
func Add[T Numeric](a, b T) T {
    return a + b
}
```

This is the definition of the generic function that uses the `Numeric` constraint.

```markup
func main() {
    fmt.Println("4 + 3 =", Add(4,3))
    fmt.Println("4.1 + 3.2 =", Add(4.1,3.2))
}
```

The previous code is the implementation of the `main()` function with the calls to `Add()`.

Running `numeric.go` produces the next output:

```markup
4 + 3 = 7
4.1 + 3.2 = 7.3
```

Nevertheless, Go rules still apply. Therefore, if you try to call `Add(4.1,3)`, you are going to get the next error message:

```markup
type checking failed for main
prog.go2:16:33: default type int of 3 does not match inferred type float64 for T
```

The reason for this error is that the `Add()` function expects two parameters of **the same data type**. However, `4.1` is a `float64` whereas `3` is an `int`, so not the same data type.

The next section shows how to use generics when defining new data types.

Just Imagine

# Defining new data types with generics

In this section we are going to create a new data type with the use of generics, which is presented in `newDT.go`. The code of `newDT.go` is the following:

```markup
package main
import (
    "fmt"
    "errors"
)
type TreeLast[T any] []T
```

The previous statement declares a new data type named `TreeLast` that uses generics.

```markup
func (t TreeLast[T]) replaceLast(element T) (TreeLast[T], error) {
    if len(t) == 0 {
        return t, errors.New("This is empty!")
    }
    
    t[len(t) - 1] = element
    return t, nil
}
```

`replaceLast()` is a method that operates on `TreeLast` variables. Apart from the function signature, there is nothing else that shows the use of generics.

```markup
func main() {
    tempStr := TreeLast[string]{"aa", "bb"}
    fmt.Println(tempStr)
    tempStr.replaceLast("cc")
    fmt.Println(tempStr)
```

In this first part of `main()`, we create a `TreeLast` variable with the `aa` and `bb` `string` values and we replace the `bb` value with `cc`, using a call to `replaceLast("cc")`.

```markup
    tempInt := TreeLast[int]{12, -3}
    fmt.Println(tempInt)
    tempInt.replaceLast(0)
    fmt.Println(tempInt)
}
```

The second part of `main()` does a similar thing to the first part using a `TreeLast` variable populated with `int` values. So `TreeLast` works with both `string` and `int` values without any issues.

Running `newDT.go` produces the next output:

```markup
[aa bb]
[aa cc]
```

The preceding is the output related to the `TreeLast[string]` variable.

```markup
[12 -3]
[12 0]
```

The final output is related to the `TreeLast[int]` variable.

The next subsection is about using generics in Go structures.

## Using generics in Go structures

In this section, we are going to implement a linked list that works with generics—this is one of the cases where the use of generics simplifies things because it allows you to implement the linked list once while being able to work with multiple data types.

The code of `structures.go` is the following:

```markup
package main
import (
    "fmt"
)
type node[T any] struct {
    Data T
    next *node[T]
}
```

The `node` structure uses generics in order to support nodes that can store all kinds of data. This does not mean that the `next` field of a `node` can point to another `node` with a `Data` field with a different data type. The rule that a linked list contains elements of the same data type still applies—it just means that if you want to create three linked lists, one for storing `string` values, one for storing `int` values, and a third one for storing JSON records of a given `struct` data type, you do not need to write any extra code to do so.

```markup
type list[T any] struct {
    start *node[T]
}
```

This is the definition of the root node of a linked list of `node` nodes. Both `list` and `node` must share the same data type, `T`. However, as stated before, this does not prevent you from creating multiple linked lists of various data types.

You can still replace `any` with a constraint in both the definition of `node` and `list` if you want to restrict the list of allowed data types.

```markup
func (l *list[T]) add(data T) {
    n := node[T]{
        Data: data,
        next: nil,
    }
```

The `add()` function is generic in order to be able to work with all kinds of nodes. Apart from the signature of `add()`, all the remaining code is not associated with the use of generics.

```markup
    if l.start == nil {
        l.start = &n
        return
    }
    
    if l.start.next == nil {
        l.start.next = &n
        return
    }
```

These two `if` blocks have to do with the adding of a new node to the linked list.

```markup
    temp := l.start
    l.start = l.start.next
    l.add(data)
    l.start = temp
}
```

The last part of `add()` has to do with defining the proper associations between nodes when adding a new node to the list.

```markup
func main() {
    var myList list[int]
```

First, we define a linked list of `int` values in `main()`, which is the linked list that we are going to work with.

```markup
    fmt.Println(myList)
```

The initial value of `myList` is `nil`, as the list is empty and does not contain any nodes.

```markup
    myList.add(12)
    myList.add(9)
    myList.add(3)
    myList.add(9)
```

In this first part, we add four elements to the linked list.

```markup
    // Print all elements
    for {
        fmt.Println("*", myList.start)
        if myList.start == nil {
            break
        }
        myList.start = myList.start.next
    }
}
```

The last part of `main()` is about printing all the elements of the list by traversing it with the help of the `next` field, which points to the next node in the list.

Running `structures.go` produces the next output:

```markup
{<nil>}
* &{12 0xc00010a060}
* &{9 0xc00010a080}
* &{3 0xc00010a0b0}
* &{9 <nil>}
* <nil>
```

Let us discuss the output a little more. The first line shows that the value of the empty list is `nil`. The first node of the list holds the value `12` and a memory address (`0xc00010a060`) that points to the second node. This goes on until we reach the last node, which holds the value of `9`, which appears twice in this linked list, and points to `nil`, because it is the last node. Therefore, the use of generics makes the linked list able to work with multiple data types.

The next section discusses the differences between using interfaces and generics to support multiple data types.

Just Imagine

# Interfaces versus generics

This section presents a program that increments a numeric value by one using interfaces and generics so that you can compare the implementation details.

The code of `interfaces.go` illustrates the two techniques and contains the next code:

```markup
package main
import (
    "fmt"
)
type Numeric interface {
    type int, int8, int16, int32, int64, float64
}
```

This is where we define a constraint named `Numeric` for limiting the permitted data types.

```markup
func Print(s interface{}) {
    // type switch
    switch s.(type) {
```

The `Print()` function uses the empty interface for getting input and a type switch to work with that input parameter.

Put simply, we are using a _type switch_ to differentiate between the supported data types—in this case, the supported data types are just `int` and `float64`, which has to do with the implementation of the type switch. However, adding more data types requires code changes, which is not the most efficient solution.

```markup
    case int:
        fmt.Println(s.(int)+1)
```

This branch is how we handle the `int` case.

```markup
    case float64:
        fmt.Println(s.(float64)+1)
```

This branch is how we handle the `float64` case.

```markup
    default:
        fmt.Println("Unknown data type!")
    }
}
```

The `default` branch is how we handle all unsupported data types.

The biggest issue with `Print()` is that due to the use of the empty interface, it accepts all kinds of input. As a result, the function signature does not help us limit the allowed data types. The second issue with `Print()` is that we need to specifically handle each case—handling more cases means writing more code.

On the other hand, the compiler does not have to guess many things with that code, which is not the case with generics, where the compiler and the runtime have more work to do. This kind of work introduces delays in the execution time.

```markup
func PrintGenerics[T any](s T) {
    fmt.Println(s)
}
```

`PrintGenerics()` is a generic function that can handle all available data types simply and elegantly.

```markup
func PrintNumeric[T Numeric](s T) {
    fmt.Println(s+1)
}
```

The `PrintNumeric()` function supports all numeric data types with the use of the `Numeric` constraint. No need to specifically add code for supporting each distinct data type, as happens with `Print()`.

```markup
func main() {
    Print(12)
    Print(-1.23)
    Print("Hi!")
```

The first part of `main()` uses `Print()` with various types of input: an `int` value, a `float64` value, and a `string` value, respectively.

```markup
    PrintGenerics(1)
    PrintGenerics("a")
    PrintGenerics(-2.33)
```

As stated before, `PrintGenerics()` works with all data types, including `string` values.

```markup
    PrintNumeric(1)
    PrintNumeric(-2.33)
}
```

The last part of `main()` uses `PrintNumeric()` with numeric values only, due to the use of the `Numeric` constraint.

Running `interfaces.go` produces the next output:

```markup
13
-0.22999999999999998
Unknown data type!
```

The preceding three lines of the output are from the `Print()` function, which uses the empty interface.

```markup
1
a
-2.33
```

The previous three lines of the output are from the `PrintGenerics()` function, which uses generics and supports all available data types. As a result, it cannot increase the value of its input because we do not know for sure that we are dealing with a numeric value. Therefore, it just prints the given input.

```markup
2
-1.33
```

The last two lines are generated by the two `PrintNumeric()` calls, which operate using the `Numeric` constraint.

So in practice, when you have to support multiple data types, the use of generics might be a better choice than using interfaces.

The next section discusses the use of reflection as a way of bypassing the use of generics.

Just Imagine

# Reflection versus generics

In this section, we develop a utility that prints the elements of a slice in two ways: first, using reflection, and second, using generics.

The code of `reflection.go` is as follows:

```markup
package main
import (
    "fmt"
    "reflect"
)
func PrintReflection(s interface{}) {
    fmt.Println("** Reflection")
    val := reflect.ValueOf(s)
    if val.Kind() != reflect.Slice {
        return
    }
    for i := 0; i < val.Len(); i++ {
        fmt.Print(val.Index(i).Interface(), " ")
    }
    fmt.Println()
}
```

Internally, the `PrintReflection()` function works with slices only. However, as we cannot express that in the function signature, we need to accept an empty interface parameter. Additionally, we have to write more code to get the desired output.

In more detail, first, we need to make sure that we are processing a slice (`reflect.Slice`) and second, we have to print the slice elements using a `for` loop, which is pretty ugly.

```markup
func PrintSlice[T any](s []T) {
    fmt.Println("** Generics")
    for _, v := range s {
        fmt.Print(v, " ")
    }
    fmt.Println()
}
```

Once again, the implementation of the generic function is simpler and therefore easier to understand. Moreover, the function signature specifies that only slices are accepted as function parameters—we do not have to perform any additional checks for that as this is a job for the Go compiler. Last, we use a simple `for` loop with `range` to print the slice elements.

```markup
func main() {
    PrintSlice([]int{1, 2, 3})
    PrintSlice([]string{"a", "b", "c"})
    PrintSlice([]float64{1.2, -2.33, 4.55})
    PrintReflection([]int{1, 2, 3})
    PrintReflection([]string{"a", "b", "c"})
    PrintReflection([]float64{1.2, -2.33, 4.55})
}
```

The `main()` function calls `PrintSlice()` and `PrintReflection()` with various kinds of input to test their operation.

Running `reflection.go` generates the next output:

```markup
** Generics
1 2 3
** Generics
a b c
** Generics
1.2 -2.33 4.55
```

The first six lines are produced by taking advantage of generics and print the elements of a slice of `int` values, a slice of `string` values, and a slice of `float64` values.

```markup
** Reflection
1 2 3
** Reflection
a b c
** Reflection
1.2 -2.33 4.55
```

The last six lines of the output produce the same output, but this time using reflection. There is no difference in the output—all differences are in the code found in the implementations of `PrintReflection()` and `PrintSlice()` for printing the output. As expected, generics code is simpler and shorter than Go code that uses reflection, especially when you must support lots of different data types.

Just Imagine

# Exercises

-   Create a `PrintMe()` method in `structures.go` that prints all the elements of the linked list
-   Create two extra functions in `reflection.go` in order to support the printing of strings using reflection and generics
-   Implement the `delete()` and `search()` functionality using generics for the linked list found in `structures.go`
-   Implement a doubly-linked list using generics starting with the code found in `structures.go`

Just Imagine

# Summary

This chapter presented generics and gave you the rationale behind the invention of generics. Additionally, it presented the Go syntax for generics as well as some issues that might come up if you use generics carelessly. It is expected that there are going to be changes to the Go standard library in order to support generics and that there is going to be a new package named `slices` to take advantage of the new language features.

Although a function with generics is more flexible, code with generics usually runs slower than code that works with predefined static data types. So, the price you pay for flexibility is execution speed. Similarly, Go code with generics has a bigger compilation time than equivalent code that does not use generics. Once the Go community begins working with generics in real-world scenarios, the cases where generics offer the highest productivity are going to become much more evident. At the end of the day, programming is about understanding the cost of your decisions. Only then can you consider yourself a programmer. So, understanding the cost of using generics instead of interfaces, reflection, or other techniques is important.

So, what does the future look like for Go developers? In short, it looks wonderful! You should already be enjoying programming in Go, and you should continue to do so as the language evolves. If you want to know the latest and greatest of Go as it is being discussed by the team, you should definitely visit the official GitHub place of the Go team at [https://github.com/golang](https://github.com/golang).

Go helps you to create great software! So, go and create great software!

Just Imagine

# Additional resources

-   Google I/O 2012—Meet the Go team: [https://youtu.be/sln-gJaURzk](https://youtu.be/sln-gJaURzk)
-   Meet the authors of Go: [https://youtu.be/3yghHvvZQmA](https://youtu.be/3yghHvvZQmA)
-   This is a video of Brian Kernighan interviewing Ken Thompson—not directly related to Go: [https://youtu.be/EY6q5dv\_B-o](https://youtu.be/EY6q5dv_B-o)
-   Brian Kernighan on successful language design—not directly related to Go: [https://youtu.be/Sg4U4r\_AgJU](https://youtu.be/Sg4U4r_AgJU)
-   Brian Kernighan: UNIX, C, AWK, AMPL, and Go Programming from the Lex Fridman Podcast: [https://youtu.be/O9upVbGSBFo](https://youtu.be/O9upVbGSBFo)
-   Why Generics? [https://blog.golang.org/why-generics](https://blog.golang.org/why-generics)
-   The Next Step for Generics: [https://blog.golang.org/generics-next-step](https://blog.golang.org/generics-next-step)
-   A Proposal for Adding Generics to Go: [https://blog.golang.org/generics-proposal](https://blog.golang.org/generics-proposal)
-   Proposal for the `slices` package: [https://github.com/golang/go/issues/45955](https://github.com/golang/go/issues/45955)