# Defining the Command-Line Process

At the core of a command-line application is its ability to process user input and return a result that either a user can easily comprehend or that another process can read as standard input. In [_Chapter 1_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_01.xhtml#_idTextAnchor014), _Understanding CLI Standards_, we discussed the anatomy of a command-line application, but this chapter will go into detail on each aspect of its anatomy, breaking down the different types of input: subcommands, arguments, and flags. Additionally, other inputs will be discussed: `stdin`, signals, and control characters.

Just as there are many types of input that a command-line application can receive, there are many types of methods for processing data. This chapter won’t leave you hanging – examples of processing for each input type will follow.

Finally, it’s just as important to understand how to return the result, either data if successful or an error on failure, in a way that both humans and computers can easily interpret.

This chapter will cover how to output the data for each end user and the best practices for CLI success. We will cover the following topics:

-   Receiving the input and user interaction
-   Processing data
-   Returning the resulting output and defining best practices


# Receiving the input and user interaction

The primary methods for receiving input via a command-line application are through its subcommands, arguments, and options, also known as **flags**. However, additional input can come in the form of `stdin`, signals, and control characters. In this section, we’ll break down each different input type and when and how to interact with the user.

## Defining subcommands, arguments, and flags

Before we start characterizing the main types of input, let’s reiterate the structural pattern that explains the generalized location for each input type in terms of its predictability and familiarity. There’s an excellent description of the pattern within the **Cobra Framework documentation**. This is one of the best explanations because it compares the structure to natural language and, just like speaking and writing, the syntax needs to be properly interpreted:

```markup
APPNAME NOUN VERB –ADJECTIVE
```

Note

The **argument** is the noun and the **command or subcommand(s)** is the verb. Like any modifier, the **flag** is an adjective and adds description.

Note

Most other programming languages suggest using two dashes instead of one. Go is unique in the fact that the single dash and double dash are equivalent to the internal flag package. It is important to note, however, that the Cobra CLI flag does differentiate between single and double dashes, where a single dash is for a short-form flag, and the double dash is for a long-form flag.

In the preceding example, the command and argument, or `NOUN VERB`, can also be ordered as `VERB NOUN`. However, `NOUN VERB` is more commonly used. Use what makes sense to you:

```markup
APPNAME ARGUMENT <COMMAND | SUBCOMMANDS> --FLAG
```

You might run up against limitations depending on your command-line parser. However, if possible, make arguments, flags, and subcommands order-independent. Now, let’s define each in more detail next and use **Cobra** to create a command that utilizes each input type.

### Commands and subcommands

At a very basic level, a command is a specific instruction given to a command-line application. In the pattern we just looked at, these are verbs. Think of the way we naturally speak. If we were to talk to a dog, we’d give it commands such as “_roll over_,” “_speak_,” or “_stay_.” Since you define the application, you can choose the verbs to define instructions. However, the most important thing to remember when choosing a command (and subcommand) is for names to be clear and consistent.

Ambiguity can cause a lot of stress for a new user. Suppose you have two commands: `yarn update` and `yarn upgrade`. For a developer who is using `yarn` for the first time, do you think it’s clear how these commands are different? Clarity is paramount. Not only does it make your application easier to use but it also puts your developer at ease.

As you gain a broad view of your application, you can intuitively determine more clear and more concise language when defining your commands. If your application feels a bit complex, you can utilize subcommands for simplification, and whenever possible, use familiar words for both commands and subcommands.

Let’s use the **Docker** application as an example of how subcommands are clearly defined. Docker has a list of management commands such as the following:

-   `container` to manage containers
-   `image` to manage images

You’ll notice that when you run `docker` `container` or `docker` `image`, the usage is printed out, along with a list of subcommands, and you’ll also notice that there are several subcommands used across these two commands. They remain consistent.

Users of Docker know that the action (`ls`, `rm`, or `inspect`, for example) is related to the subject (`image` or `container`). The command follows the expected pattern of `"APPNAME ARGUMENT COMMAND"` – `docker` `image` `ls` and `docker` `container` `ls` too. Notice that `docker` also uses familiar Unix commands – `ls` and `rm`. Always use a familiar command where you can.

Using the Cobra CLI, let’s make two commands, with one as a subcommand of the other. Here’s the first command we’ll add:

```markup
cobra-cli add command
command created at /Users/marian/go/src/github.com/
  marianina8/application
```

Then, add the subcommand:

```markup
cobra-cli add subcommand
subcommand created at /Users/marian/go/src/github.com/
  marianina8/application
```

Then, create it as a subcommand by modifying the default line to run `AddCommand` on `commandCmd`:

```markup
func init() {
    commandCmd.AddCommand(subcommandCmd)
}
```

The Cobra CLI makes it incredibly easy not only to create commands but also subcommands as well. Now, when the command is called with the subcommand, we get confirmation that the subcommand is called:

```markup
./application command subcommand
subcommand called
```

Now, let us understand arguments.

### Arguments

Arguments are nouns – things – that are acted upon by the command. They are positional to the command and usually come before the command. The order is not strict; just be consistent with the order throughout your application. However, the very first argument is the application name.

Multiple arguments are okay for actions against multiple files, or multiple strings of input. Take, for example, the `rm` command and removing multiple files. For example, `rm arg1.txt arg2.txt arg3.txt` would act on (by removing) the multiple files listed after the command. Allow globbing where it makes sense. If a user wants to remove all the text files in the current directory, then an example of `rm *.txt` would also be expected to work. Now, consider the `mv` command, which requires two arguments for the source and target files. An example of `mv old.txt new.txt` will move `old.txt`, the source, to the target, `new.txt`. Globs may also be used with this command.

Note

Having multiple arguments for _different_ things might mean rethinking the way that you’re structuring your command. It could also mean that you could be utilizing flags here instead.

Again, familiarity plays in your favor. Use the standard name if there is one and your users will thank you. Here are examples of some common arguments: `history`, `tag`, `volume`, `log`, and `service`.

Let’s modify the subcommand’s generated `Run` field to identify and print out its arguments:

```markup
Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 0 {
        fmt.Println("subcommand called")
    } else {
        fmt.Println("subcommand called with arguments: ",
          args)
    }
},
```

Now, when we run the same subcommand with arguments, the following output is printed out:

```markup
  ./application command subcommand argument1 argument2
subcommand called with arguments:  [argument1 argument2]
```

Interestingly, flags can provide more clarity over arguments. In general, it does require more typing, but flags can make it more clear what’s going on. Another additional benefit is if you decide to make changes to how you receive input, it’s much easier to add or remove a flag than it is to modify an existing command, which can break things.

### Flags

Flags are adjectives that add a description to an action or command. They are named parameters and can be denoted in different ways, with or without a user-specified value:

-   A **hyphen with a single-letter** **name** (`-h`)
-   A **double-hyphen with a multiple-letter** **name** (`--help`)
-   A **double-hyphen with a multiple-letter name and a user-specified value** (`--file audio.txt`, or `–-file=audio.txt`)

It’s important to have full-length versions of all flags – single letters are only useful for commonly used flags. If you use single letters for all available flags, there may be more than one flag that starts with that same letter, and that single letter would make sense intuitively for more than one flag. This can add confusion, so it’s best not to clutter the list of single-letter flags.

Single-letter flags may also be concatenated together. For example, take the `ls` command. You can run `ls -l -h -F` or `ls -lhF` and the result is the same. Obviously, this depends on the command-line parser used, but because CLI applications typically allow you to concatenate single-letter flags, it’s a good idea to allow this as well.

Finally, the flag order is typically not strict, so whether a user runs `ls –lhF`, `ls –hFl`, or `ls –Flh`, the result is the same.

As an example, we can add a couple of flags to the root command, one local and another persistent, meaning that it’s available to the command and all subcommands. In `commandCmd`, within the `init()` function, the following lines do just that:

```markup
commandCmd.Flags().String("localFlag", "", "a local string
  flag")
commandCmd.PersistentFlags().Bool("persistentFlag", false,
  "a persistent boolean flag")
```

In `commandCmd`’s `Run` field, we add these lines:

```markup
localFlag, _ := cmd.Flags().GetString("localFlag")
if localFlag != "" {
    fmt.Printf("localFlag is set to %s\n", localFlag)
}
```

In `subcommandCmd`’s `Run` field, we also add the following lines:

```markup
persistentFlag, _ := cmd.Flags().GetBool("persistentFlag")
fmt.Printf("persistentFlag is set to %v\n", persistentFlag)
```

Now, when we compile the code and run it again, we can test out both flags. Notice that there are multiple ways of passing in flags and in both cases, the results are the same:

```markup
  ./application command --localFlag=”123”
command called
localFlag is set to 123
  ./application command --localFlag “123”
command called
localFlag is set to 123
```

The persistent flag, although defined at the `commandCmd` level, is available within `subcommandCmd`, and when the flag is missing, the default value is used:

```markup
  ./application command subcommand
subcommand called
persistentFlag is set to false
  ./application command subcommand --persistentFlag
subcommand called
persistentFlag is set to true
```

Now, we’ve covered the most common methods of receiving input to your CLI: commands, arguments, and flags. The next methods of input include piping, signal and control characters, and direct user interaction. Let’s dive into these now.

## Piping

In Unix, piping redirects the standard output of one command-line application into the standard input of another. It is represented by the ‘`|`’ character, which combines two or more commands. The general structure is `cmd1 | cmd2 |cmd3 | .... | cmdN`, the standard output from `cmd1` is the standard input for `cmd2`, and so on.

Creating a simple command-line application that does one thing and one thing well follows the Unix philosophy. It reduces the complexity of a single CLI, so you’ll see many examples of different applications that can be chained together by pipes. Here are a few examples:

-   `cat file.txt | grep "word" |` `sort`
-   `sort list.txt |` `uniq`
-   `find . -type f –name main.go |` `grep audio`

As an example, let’s create a command that takes in standard input from a common application. Let’s call it `piper`:

```markup
cobra-cli add piper
piper created at /Users/marian/go/src/github.com/
  marianina8/application
```

For the newly generated `piperCmd`’s `Run` field, add the following lines:

```markup
reader := bufio.NewReader(os.Stdin)
s, _ := reader.ReadString('\n')
fmt.Printf("piped in: %s\n", s)
```

Now, compile and run the `piper` command with some piped-in input:

```markup
  echo “hello world” | ./application piper
piper called
piped in: hello world
```

Now, suppose your command has a standard output that is written to a broken pipe; the kernel will raise a `SIGPIPE` signal. This is received as input to the command-line application, which can then output an error regarding the broken pipe. Besides receiving signals from the kernel, other signals, such as `SIGINT`, can be triggered by users who press control character key combinations such as _Ctrl + C_ that interrupt the application. This is just one type of signal and control character, but more will be discussed in the following section.

## Signals and control characters

As the name implies, signals are another way to communicate specific and actionable input by signaling to a command-line application. Sometimes, these signals can be from the kernel, or from users that press control characters key combinations and trigger signals to the application. There are two different types of signals:

-   **Synchronous signals** – triggered by errors that occur when the program executes. These signals include `SIGBUS`, `SIGFPE`, and `SIGSEGV`.
-   **Asynchronous signals** – triggered from the kernel or another application. These signals include `SIGHUP`, `SIGINT`, `SIGQUIT`, and `SIGPIPE`.

Note

A few signals, such as `SIGKILL` and `SIGSTOP`, may not be caught by a program, so utilizing the `os/signal` package for custom handling will not affect the result.

There is a lot to discuss in depth on signals, but the main point is that they are just another method of receiving input. We’ll stay focused on how this data is received by the command-line application. The following is a table explaining some of the most commonly used signals, control character combinations, and their descriptions:

![Figure 5.1 – Table of signals with related key combinations and descriptions](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_5.01.jpg)

Figure 5.1 – Table of signals with related key combinations and descriptions

The following are two function calls added to `rootCmd` to handle exiting your application with grace when a `SIGINT` or `SIGTSTP` signal is received. The `Execute` function that calls `rootCmd` now looks like this:

```markup
func Execute() {
    SetupInterruptHandler()
    SetupStopHandler()
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}
```

The `SetupInterruptHandler` code is as follows:

```markup
func SetupInterruptHandler() {
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGINT)
    go func() {
        <-c
        fmt.Println("\r- Wake up! Sleep has been
          interrupted.")
        os.Exit(0)
    }()
}
```

Similarly, the `SetupStopHandler` code is as follows:

```markup
func SetupStopHandler() {
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTSTP)
    go func() {
        <-c
        fmt.Println("\r- Wake up! Stopped sleeping.")
        os.Exit(0)
    }()
}
```

Now, we’ll need a command to interrupt or stop the application. Let’s use the Cobra CLI and add a `sleep` command:

```markup
  cobra-cli add sleep
sleep created at /Users/marian/go/src/github.com/
  marianina8/application
```

The `Run` field of `sleepCmd` is changed to run an infinite loop that prints out some Zs (`Zzz`) until a signal interrupts the `sleep` command and wakes it up:

```markup
Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("sleep called")
    for {
        fmt.Println("Zzz")
        time.Sleep(time.Second)
    }
},
```

By running the `sleep` command and then using _Ctrl + C_, we get the following output:

```markup
  ./application sleep
sleep called
Zzz
Zzz
- Wake up!  Sleep has been interrupted.
- Wake up!  Stopped sleeping.
```

Trying again but now using _Ctrl + Z_, we get the following output:

```markup
  ./application sleep
sleep called
Zzz
Zzz
Zzz
- Wake up!  Stopped sleeping.
```

You can utilize signals to interrupt or quit your application gracefully or take action when an alarm is triggered. While commands, arguments, and flags are the most commonly known types of input for command-line applications, it is important to consider handling these signal inputs to create a more robust application. If a terminal hangs and `SIGHUP` is received, your application can save information on the last state and handle cleanup where necessary. In this case, while it’s not as common, it’s just as important.

## User interaction

Although your user input can be in the form of commands, arguments, and flags, user interaction is more of a back-and-forth interaction between the user and the application. Suppose a user misses a required flag for a particular subcommand – your application can prompt the user and receive the value for that flag via standard input. Sometimes, rather than utilizing the more standard input of commands, arguments, and flags, an interactive command-line application can be built instead.

An interactive CLI would prompt for input and then receive it through `stdin`. There are some useful packages for building interactive and accessible prompts in Go. For the following examples, we’ll use the [https://github.com/AlecAivazis/survey](https://github.com/AlecAivazis/survey) package. There are multiple fun ways to prompt input using the `survey` package. A `survey` command will ask questions that need to be stored in a variable. Let’s define it as `qs`, a slice of the `*``survey.Question` type:

```markup
var qs = []*survey.Question{}
```

`survey` can prompt the user for different types of input, as defined here:

-   **Simple** **text input**

At a very basic level, users can receive basic text input:

```markup
{
    Name: "firstname",
    Prompt: &survey.Input{Message: "What is your first
      name?"},
    Validate: survey.Required,
    Transform: survey.Title,
},
Output:
  ? What is your first name?
```

-   **Suggesting options**

This terminal option allows you to give the user suggestions for the prompted question:

```markup
{
    Name: "favoritecolor",
    Prompt: &survey.Select{
    Message: "What's your favorite color?",
    Options: []string{"red", "orange", "yellow",
      "green", "blue", "purple", "black", "brown",
        "white"},
    Default: "white",
},
Output:
  ? What is your favorite color? [tab for suggestions]
```

Hitting the _Tab_ key will show the available options:

```markup
? What is your favorite color? [Use arrows to
  navigate, enter to select, type to complement
    answer]
red
orange
yellow
green
blue
purple
black
brown
white
```

-   **Inputting** **multiple lines**

When receiving input, sometimes, pressing the _Return_ key will immediately pass any text received before directly as input to the program. Utilizing the `survey` package allows you to enter multiple lines before receiving input:

```markup
{
    Name: "story",
    Prompt: &survey.Multiline{
    Message: "Tell me a story.",
    },
},
Output:
  ? Tell me a story [Enter 2 empty lines to finish]
A long line time ago in a faraway town, there lived a
  princess who lived in a castle far away from the
    city.  She was always sleeping, until one day…
```

-   **Protecting** **password input**

To keep data private, when inputting private information, the `survey` package will replace the characters with `*` symbols:

```markup
{
    Name: "secret",
    Prompt: &survey.Password{
    Message: "Tell me a secret",
    },
},
Output:
? Tell me a secret: ************
```

-   **Confirming with Yes** **or No**

Users can respond with a simple yes or no to the command prompt:

```markup
{
    Name: "good",
    Prompt: &survey.Confirm{
    Message: "Are you having a good day?",
    },
},
Output:
? Are you having a good day? (Y/n)
```

Now, let us see how to select from a checkbox option.

-   **Selecting from a** **checkbox option**

Multiple options can be selected within a vertical checkbox option. Navigating the options is done with the up and down arrows, and selecting is done with the spacebar:

```markup
{
    Name: "favoritepies",
    Prompt: &survey.MultiSelect{
    Message: "What pies do you like:",
    Options: []string{"Pumpkin", "Lemon Meringue",
      "Cherry", "Apple", "Key Lime", "Pecan", "Boston
        Cream", "Rhubarb", "Blackberry"},
    },
},
Output:
? What pies do you like: [Use arrows to move, space to
select, <right> to all, <left> to none, type to
filter]
> [ ] Pumpkin
  [ ] Lemon Meringue
  [ ] Cherry
  [ ] Apple
  [ ] Key Lime
  [ ] Pecan
….
```

Create a new `survey` command with the following:

`cobra-cli` `add survey`

The `Run` field of `surveyCmd` creates a struct that receives all the answers to questions asked:

```markup
Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("survey called")
    answers := struct {
        FirstName string
        FavoriteColor string
        Story string
        Secret string
        Good bool
        FavoritePies []string
    }{}
```

The `Ask` method then takes in the questions, `qs`, and then receives all the answers to the questions asked into a pointer to the `answers` struct:

```markup
    err := survey.Ask(qs, &answers)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
```

Finally, the results are printed out:

```markup
    fmt.Println("*********** SURVEY RESULTS ***********")
    fmt.Printf("First Name: %s\n", answers.FirstName)
    fmt.Printf("Favorite Color: %s\n",
        answers.FavoriteColor)
    fmt.Printf("Story: %s\n", answers.Story)
    fmt.Printf("Secret: %s\n", answers.Secret)
    fmt.Printf("It's a good day: %v\n", answers.Good)
    fmt.Printf("Favorite Pies: %s\n", answers.FavoritePies)
},
```

Testing out the `survey` command, we get the following:

```markup
  ./application survey
survey called
? What is your first name? Marian
? What's your favorite color? white
? Tell me a story.
I went to the dodgers game last night and
they lost, but I still had fun!
? Tell me a secret ********
? Are you having a good day? Yes
? What pies do you prefer: Pumpkin, Lemon Meringue, Key
    Lime, Pecan, Boston Cream
*********** SURVEY RESULTS ***********
First Name: Marian
Favorite Color: white
Story: I went to the dodgers game last night and
they lost, but I still had fun!
Secret: a secret
It's a good day: true
Favorite Pies: [Pumpkin Lemon Meringue Key Lime Pecan
    Boston Cream]
```

Although these examples are just a selection of the many input prompts provided by the `survey` package, you can visit the GitHub page to view examples of all the possible options. Playing around with prompts reminds me of early text-based RPG games that used them to prompt the gamer’s character. Having learned about the many different types of input, whether user-based, from the kernel, or from other piped applications, let’s discuss how to process this incoming data.

Just Imagine

# Processing data

**Data processing** is when raw data is fed into a process, analyzed, and then used to generate useful information or output. At a very general level, this can include sorting data, searching or querying for data, and converting data from one type of input into another. For a CLI, the input can be received in the different ways discussed in the previous section. When receiving arguments using the Cobra framework, all the values are read in as string input. However, given a string of `123`, we can do a type check by utilizing the `strconv` package’s `Atoi` method, which converts an ASCII string into an integer:

```markup
val, err := strconv.Atoi("123")
```

If the string value cannot be converted because it isn’t a string representation of an integer, then an error will be thrown. If the string is a representation of an integer, then the integer value will be stored in the `val` variable.

The `strconv` package can be used to check, with conversion, many other types, including Boolean, float, and `uint` values as well.

Flags, on the other hand, can have predefined types. Within the Cobra framework, the `pflag` package is used, which is just an extension of the standard go `flag` package. For example, when a flag is defined, you can define it specifically as a `String`, `Bool`, `Int`, or custom type. The preceding `123` value, if read in as an `Int` flag, could be defined with the following lines of code:

```markup
var intValue int
flag.IntVar(&intValue, "flagName", 123, "help message")
```

This can be done similarly for `String` and `Bool` flags. You can even create a flag with a custom, specific interface using the `Var` method:

```markup
var value Custom
flag.Var(&value, "name", "help message")
```

Just ensure that the `Custom` struct satisfies the following interface defined within the `pflag` package for custom flags:

```markup
// (The default value is represented as a string.)
type Value interface {
    String() string
    Set(string) error
    Type() string
}
```

I defined the `Custom` struct as the following:

```markup
type Custom struct {
    Value string
}
```

Therefore, the `Set` method is simply defined as follows:

```markup
func (c *Custom) Set(value string) error {
    c.Value = value
    return nil
}
```

Passing the value into the flag was handled by `flag: --name="custom value`. The `String` method is then used to print the value:

```markup
fmt.Println(cmd.Flag("name").Value.String())
```

It looks like this:

```markup
custom value
```

Besides passing in string values that can be converted into different types, oftentimes, a path to a file is passed in. There are multiple ways of reading data from files. Let’s list each, along with a method to handle this way of reading in a file and a pro and a con for each:

-   **In its entirety, all at once**: The `os.ReadFile` method reads the entire file and returns its contents. It does not error when encountering the **end of** **file** (**EOF**):
    
    ```markup
    func all(filename string) {
    
        content, err := os.ReadFile(filename)
    
        if err != nil {
    
            fmt.Printf("Error reading file: %s\n", err)
    
            return
    
        }
    
        fmt.Printf("content: %s\n", content)
    
    }
    ```
    
-   **Pros**: Faster performance
-   **Cons**: Consumes more memory in a shorter amount of time
-   **In predefined chunks**: The `file.Read` method reads in the buffer at its predetermined size and returns the bytes, which can be printed after being cast as a string. Unlike the `ioutil.ReadFile` method, `file.Read` from the buffer will error when it reaches the EOF:
    
    ```markup
    func chunk(file *os.File) {
    
        const size = 8 // chunk size
    
        buff := make([]byte, size)
    
        fmt.Println("content: ")
    
        for {
    
            // read content to buffer of size, 8 bytes
    
            read8Bytes, err := file.Read(buff)
    
            if err != nil {
    
                if err != io.EOF {
    
                    fmt.Println(err)
    
                }
    
                break
    
            }
    
            // print content from buffer
    
            fmt.Println(string(buff[:read8Bytes]))
    
        }
    ```
    
-   **Pros**: Easy to implement, consumes little memory
-   **Cons**: If the chunks are not properly chosen, you may have inaccurate results, increased complexity when comparing or analyzing the data, and potential error propagation.
-   **Line by line**: By default, a new scanner will split the text up by lines, so it’s not necessary to define the `split` function. The `scanner.Text()` method reads into the next token that delimits each scan – in the following example, line by line. Finally, `scanner.Scan()` does not return an error when it encounters the EOF:
    
    ```markup
    func line(file *os.File) {
    
        scanner := bufio.NewScanner(file)
    
        lineCount := 0
    
        for scanner.Scan() {
    
            fmt.Printf("%d: %s\n", lineCount,
    
              scanner.Text())
    
            lineCount++
    
        }
    
        if err := scanner.Err(); err != nil {
    
            fmt.Printf("error scanning line by line:
    
                %s\n", err)
    
        }
    
    }
    ```
    
-   **Pros**: Easy to implement – an intuitive way to read in data and output data.
-   **Cons**: Processing an extremely large file may cause memory constraints. Increased complexity may cause inaccurate results, if the data is not well suited to line by line processing.
-   **Word by word** To overwrite the default `Split` function, pass `bufio.ScanWords` into the `Split` function. This will then define the tokens between each word and scan between each token. Again, scanning in this way will not encounter an error at the EOF either:
    
    ```markup
    func word(file *os.File) {
    
        scanner := bufio.NewScanner(file)
    
        scanner.Split(bufio.ScanWords)
    
        wordCount := 0
    
        for scanner.Scan() {
    
            fmt.Printf("%d: %s\n", wordCount,
    
                scanner.Text())
    
          wordCount++
    
        }
    
        if err := scanner.Err(); err != nil {
    
            fmt.Printf("error scanning by words: %s\n",
    
                err)
    
        }
    
    }
    ```
    
-   **Pros**: Easy to implement – an intuitive way to read data and output data
-   **Cons**: Inefficient and time consuming for large files. Increased complexity may cause inaccurate results, if the data is not well suited to word by word processing

Choosing the way to handle processing the data received from the file depends on the use case. Additionally, there are three main types of data processing: batch, online, and real-time.

As you can guess from the name, batch processing takes similar tasks that are collected, or batched, and then runs them simultaneously. Online processing requires internet connectivity to reach an API endpoint to fully process data and return a result. Real-time processing is the execution of data in such a short period that the data is instantaneously output.

Examples of different use cases requiring a specific type of processing vary. Bank transactions, billing, and reporting often use batch processing.

A CLI that utilizes an API behind the scenes would often require internet access to handle online processing. Real-time processing is used when timeliness is of utmost importance, often in manufacturing, fraud detection, and computer vision tools.

Once the data has been processed, the result must be returned to the user or receiving process. In the next section, we will discuss the details of returning the output and defining the best practices for returning data.

Just Imagine

# Returning the resulting output and defining best practices

When returning output from a process, it’s important to know to who or what you’re returning data. It’s incredibly important to return output that’s human-readable. However, to determine whether you’re returning data to a human or a machine, check whether you’re writing to a TTY. Remember TTY? You can refer to [_Chapter 1_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_01.xhtml#_idTextAnchor014), _Understanding CLI Standards_, in which we discussed the history of the CLI interface and the term TTY, short for teletypewriter or teletype.

If writing to a TTY, we can check whether the `stdout` file descriptor refers to a terminal or not, and change the output depending on the result.

Let’s check out this block of code, which checks whether the `stdout` file descriptor is writing to a TTY or not:

```markup
if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() &
    os.ModeCharDevice) != 0 {
    fmt.Println("terminal")
} else {
    fmt.Println("not a terminal")
}
```

Let’s call it within the `Run` method of a command called `tty` using the following command:

```markup
./application tty
```

Then, the output is as follows,:

```markup
terminal
```

However, if we pipe the output to a file by calling `./application tty > file.txt`, then the contents of the file are as follows:

```markup
not a terminal
```

Certainly, it makes sense to add colored ASCII text when returning output to a human, but that’s often useless and extraneous information for output to a machine process.

When writing output, always put humans first, specifically in terms of usability. However, if the machine-readable output does not affect usability, then output in machine-readable output. Because streams of text are universal input in Unix, it’s typical for programs to be linked together by pipes. The output is typically lines of text, and programs expect input as lines of text as well. A user should expect to write output that can easily be grepped. You cannot know for sure where the output will be sent to and which other processes may be consuming the output. Always check whether the output is sent to a terminal and print for another program if it’s not. However, if using a machine-readable output breaks usability, but the human-readable output cannot be easily processed by another machine process, default to human-readable output and then define the `–plain` flag to display this output as machine-readable output. Clean lines of text in tabular format are easily integrated with `grep` and `awk`. This gives the user the choice to define the format of the output.

Beyond defining the output for humans versus machines, it’s standard to add a flag to define a specific format for the data returned. The `–json` flag is used when requesting data to be returned in JSON format and the `–xml` flag is used to request XML format. There’s a Unix tool, `jq`, that can be integrated with a program’s JSON output. In fact, this tool can manipulate any data returned in JSON format. Many tools within the Unix ecosystem take advantage of this and you can too.

Historically, because many of the older Unix programs were written for scripts or other programs, often, no output is returned on success. This can be confusing for users. Success cannot always be assumed, so it’s ideal to display output on success. There’s no reason to elaborate, so keep it brief and informative. Defining a `–quit` (or `–q`) flag can suppress unnecessary information if necessary.

Sometimes, a CLI can keep track of the state. The **GitHub CLI** is probably the best and most common example that many of you have already experienced. It does an excellent job of informing users of state changes and the current state using `git status`. This information needs to be transparent to the user, as it can often confirm the result of an action expected to change the state. The user understands their possible next steps by knowing the state.

Some of these next steps may also be suggested to the user. In fact, it’s ideal to give users suggestions because it feels like they are being guided along, rather than left alone in the wild with a new CLI application. When a user first interacts with a CLI, it’s best to make the learning experience similar to a guided adventure. Let’s give a quick example in regard to GitHub’s CLI. Consider when you have to merge the main branch into your current branch. Now and then, there’ll be conflicts after the merge, and the CLI guides you when you check `git status`:

```markup
On branch {branch name}
Your branch and 'origin/{branch name}' have diverged
And have 1 and 1 different commits each, respectively.
    (use "git pull" to merge the remote branch into yours)
You have unmerged paths.
    (fix conflicts and run "git commit")
    (use "git merge –abort" to abort the merge)
Unmerged paths:
    (use "git add <file>..." to mark resolution)
             Both modified:     merge.json
```

Note

The response reminds the user of their current branch and state, as well as suggesting different options that the user could take. Not all CLIs handle the state, but when you do, it’s best to make it well known and provide users with a clear path forward.

If there’s any communication with a remote server, reading or writing of files (except for a cache), or any other actions that cross the boundary of the program’s internals, communicate those actions to the user. I love HomeBrew’s `install` command on their CLI. It’s clear exactly what’s going on behind the scenes when you use `brew install` for an application.

When a file is being downloaded or created, it’s clearly stated:

```markup
==> Downloading https://ghcr.io/v2/homebrew/core/dav1d/manifests/1.0.0
###########################################################
############# 100.0%
```

And look how hashtags are used to designate progress – they utilize ASCII characters in a way that increases information density. I love the cold glass of beer icon next to files existing in the `Cellar` folder. It makes you think of all the brew formulas existing inside a beer cellar. **Emojis** are worth a thousand words.

When an error is evoked, the text is displayed in red, intending to evoke a sense of urgency and alertness. Color, if used, must be used intentionally. A green failure, or red success, is confusing for users. I’m certain, just like utilizing ASCII art to increase information density, color has the same purpose. A green success cannot be mistaken easily for a failure, and vice versa. Make sure to make important information stand out by using colors infrequently. Too many colors will make it difficult for anything to stand out.

However, while color may excite some of us, it annoys others. There may be any number of reasons why someone may want to disable the color in their CLI. For whatever reason to continue in a black-and-white world, there are specific times color should not be used:

-   When piping to another program
-   When the `NO_COLOR` environment variable is set
-   When the `TERM` environment variable is set to `dumb`
-   When the `–no-color` flag is passed
-   When your app’s `MYAPP_NO_COLOR` environment variable is set

It goes without saying that if we don’t allow colors, we don’t allow animations either! Well, I won’t tell you what to do, just try it for yourself – pipe an animation to a file via `stdout`. I dare you! You might end up with some great ASCII art, but it will be busy and difficult to understand the data. The goal is clarity. With ASCII art, color intent, and animations to increase the information density, we need to understand at some point that we need to use clear words that are understood by all. Consider your wording from the perspective of someone who is using your CLI for the first time. Guide users with your words.

As for printing log output, only do so under the verbose mode, represented by the `–verbose` flag and `–v` for short. Don’t use the `stderr` file descriptor as a log file.

If a CLI outputs a lot of text at once, such as `git diff`, a pager is used. Thank goodness. This makes it so much easier to page through the output to review differences rather than receiving all the text at once. This is just one of the many ways that GitHub has delivered a very thoughtful CLI to its users.

Finally, make errors stand out – use red text or a red _x_ emoji to increase understanding if an error occurs. If colors are disabled, then use text to communicate that an error has occurred and offer some suggestions for the next steps to take – and, even better, an avenue toward support via email or a website.

Just Imagine

# Summary

In this chapter, you learned about the command-line process – receiving input, processing data, and returning the output. The most popular different types of input have been discussed: from **subcommands**, **arguments**, and **flags**, to **signals** and **control characters**.

We created an interactive survey to receive input from a user and discussed data processing. We also learned how to take the first steps of processing: converting argument string data, converting and checking the type, receiving data from typed and custom flags, and finally, reading data from a file.

We also covered a brief explanation regarding the different types of processing: batch, online, and real-time processing. Ultimately, the use case will lead you to understand what sort of input you’ll require, and whether running tasks in batches, over the internet, or in real time is required.

Returning the output is just as important as receiving it, if not more! This is your chance to create a more pleasant experience for your user. Now that you’re developing for humans first, you have the opportunity to put yourself in their shoes.

How would you want to receive data in a way that makes you feel assured, understanding failures and what to do next, and where to find help? Not all processes run successfully, so let’s at least make users feel that they’re on the path to success. In _Part 2_, [_Chapter 6_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_06.xhtml#_idTextAnchor123), _Calling External Processes, Handling Errors and Timeouts_, we will continue to discuss the command-line process in more detail, focusing on external processes and how to handle timeouts and errors and communicate them to the user effectively.

Just Imagine

# Questions

1.  Are arguments or flags preferred for CLI programs? Why?
2.  What key combination interrupts a computer process?
3.  What flag can be added to your CLI to modify the output into plain output that can easily be integrated with tools such as `grep` and `awk`?

Just Imagine

# Answers

1.  Flags are preferred for CLI programs because they make it much easier to add or remove functionality.
2.  _Ctrl +_ _C_.
3.  The `–plain` flag can be added to remove any unnecessary data from the output.

Just Imagine

# Further reading

-   What is a TTY? (https://unix.stackexchange.com/questions/4126/what-is-the-exact-difference-between-a-terminal-a-shell-a-tty-and-a-con/4132#4132)
-   NO\_COLOR (https://no-color.org/)
-   _12 Factor CLI_ _Apps_ (https://medium.com/@jdxcode/12-factor-cli-apps-dd3c227a0e46)