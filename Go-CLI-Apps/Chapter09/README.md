# Building for Humans versus Machines

Thinking about your end user while you develop your command-line application will make you a more empathic developer. Consider not just how you feel about the way certain **command-line interfaces** (**CLIs**) behave but also how you could improve the experience for yourself and others. Much goes into usability and it’s not possible to cram it all into a single chapter, so we suggest following up with the suggested article and book in the _Further_ _reading_ section.

One of the first points to consider when building your command-line interface is that while it will be primarily used by humans, it can also be called within scripts, and the output from your program could be used as input into another application, such as **grep** or **awk**. Within this chapter, we’ll go over how to build for both and how to tell when you’re outputting to one versus the other.

The second point is the use of ASCII art to increase information density. Whether you’re outputting data as a table, or adding color or emojis, the idea is to make information jump out of the terminal in a way that the end user can quickly understand the data presented to them.

Finally, consistency also increases clarity for your users. When your CLI uses consistency within flag names and positional arguments across different commands and subcommands, your user can feel more confident in the steps they need to take when navigating your CLI. By the end of the chapter, you’ll hopefully have more to consider when building your CLI and be prompted to make usability improvements. Within this chapter, we’ll cover the following topics:

-   Building for humans versus machines
-   Increasing information density with ASCII art
-   Being consistent across CLIs


# Building for humans versus machines

CLIs have a long history where their interactions were tailored for other programs and machines. Their design was more similar to functions within a program than a graphical interface. Because of this, many Unix programs today still operate under the assumption that they will be interacting with another program.

Today, however, CLIs are more often used by humans than other machines while still carrying an outdated interaction design. It’s time that we built CLIs for their primary user—the human.

In this section, we will compare the machine-first design to the human-first design and learn how to check whether you are outputting to the TTY. As we can recall from [_Chapter 1_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_01.xhtml#_idTextAnchor014), _Understanding CLI Standards_, **TTY** is short for **TeleTYpewriter**, which evolved into the input and output device to interact with large mainframes. In today’s world, desktop environments for operating systems, or **OSs** for short, provide a terminal window. This terminal window is a virtual teletypewriter. They are often called **pseudo-teletypes**, or **PSY** for short. It’s also an indication that a human is on the other end, versus a program.

## Is it a TTY?

First, let’s understand devices. **Devices** can be anything from hard drives, RAM disks, DVD players, keyboards, mouses, printers, tape drivers, to TTYs. A **device driver** provides the interface between the operating system and the device; it provides an API that the operating system understands and accepts.

![Figure 8.1 – Figure showing communication from OS to the TTY device via a device driver](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_8.1._B18883.jpg)

Figure 8.1 – Figure showing communication from OS to the TTY device via a device driver

On Unix-based OSs, there are two major device drivers:

-   **Block** – interfaces for devices such as hard drives, RAM disks, and DVD players
-   **Character** – interfaces for the keyboard, mouse, printers, tape drivers, TTYs, and so on

If you check that the standard input, **stdin**, or standard output, **stdout**, is a **character** device, then you can assume that you are receiving input from or sending output to a human.

### Is it a TTY on a Unix or Linux operating system?

In a terminal, if you type the `tty` command, it will output the file name connected to **stdin**. Effectively, it is the number of the terminal window.

Let’s run the command in our Unix terminal window and see what the result is:

```markup
mmontagnino@Marians-MacCourse-Pro marianina8 % tty
/dev/ttys014
```

There is a shorthand silent, `-s`, flag that can be used to suppress output. However, the application still returns an exit code:

-   Exit code 0 – standard input is coming from a TTY
-   Exit code 1 – standard input is not coming from a TTY
-   Exit code 2 – syntax error from invalid parameters
-   Exit code 3 – a write error

In Unix, typing `&&` after a command means that the second command will only execute if the first command runs successfully, with exit code 0. So, let’s try this code to see if we’re running in a TTY:

```markup
mmontagnino@Marians-MacCourse-Pro marianina8 % tty -s && echo "this is a tty"
this is a tty
```

Since we ran those commands in a terminal, the result is `this is` `a tty`.

### Programmatically check on a Unix or Linux operating system

There are a few ways to do this programmatically. We can use the code located in the `Chapter-8/isatty.go` file:

```markup
func IsaTTY() {
  fileInfo, _ := os.Stdout.Stat()
  if (fileInfo.Mode() & os.ModeCharDevice) != 0 {
    fmt.Println("Is a TTY")
  } else {
    fmt.Println("Is not a TTY")
  }
}
```

The preceding code grabs the file info from the standard output, **stdout**, file with the following code:

```markup
fileInfo, _ := os.Stdout.Stat()
```

Then, we check the result of a bitwise operation, `&`, between `fileInfo.Mode()` and `os.ModeCharDevice`. The bitwise operator, `&`, copies a bit to the result if it exists in both operands.

Let’s take a quite simple example: `7&6` within a truth table. `7` values are represented by binary `111` and `6` values are represented by `110`.

![Figure 8.2 – Truth table to show the & operation calculation](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_8.2._B18883.jpg)

Figure 8.2 – Truth table to show the & operation calculation

The `&` operation checks each bit and whether they are the same, and if so, carry a bit over, or 1. If the bits differ, no bit is carried over, or 0. The resulting value is `110`.

Now, in our more complicated example, the following code, `fileInfo.Mode() & os.ModeCharDevice`, performs a bitwise operation between `fileInfo.Mode()` and `os.ModeCharDevice`. Let’s look at what this operation looks like when the code standard output is connected to a terminal:

<table id="table001-2" class="No-Table-Style"><colgroup><col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style" colspan="2"><p><strong class="bold">Is </strong><span class="No-Break"><strong class="bold">a TTY</strong></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">Code</span></p></td><td class="No-Table-Style"><p><span class="No-Break">Value</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">fileInfo.Mode()</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">Dcrw--w----</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">os.ModeCharDevice</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">c---------</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><code class="literal">fileInfo.Mode()</code> <code class="literal">&amp; </code><span class="No-Break"><code class="literal">os.ModeCharDevice</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">c---------</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p>(<code class="literal">fileInfo.Mode()</code> <code class="literal">&amp;</code> <code class="literal">os.ModeCharDevice)</code> <code class="literal">!= 0</code></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">TRUE</code></span></p></td></tr></tbody></table>

Figure 8.3 – The code next to its value when standard output is connected to a TTY

In _Figure 8__.3_, the file mode of the standard output is defined by the `fileInfo.Mode()` method call; its value is **Dcrw--w----**. If you look at the documentation for the **os** package at [https://pkg.go.dev/os](https://pkg.go.dev/os), you will see that the `os.ModeDevice`, **D**, bit is set to indicate that the file is a device file, followed by the `os.ModeCharDevice`, **c**, bit set to indicate that it is a Unix character device. When we do a bitwise operation against the mode of `stdin` against `os.ModCharDevice`, we see that the same bits are carried over and the result does not equal zero, hence `(fileInfo.Mode() & os.ModeCharDevice) != 0` is **true**, and the device is a TTY.

What would this code look like if the output were piped into another process? Let’s look:

<table id="table002" class="No-Table-Style"><colgroup><col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style" colspan="2"><p><strong class="bold">Is not </strong><span class="No-Break"><strong class="bold">a TTY</strong></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">Code</span></p></td><td class="No-Table-Style"><p><span class="No-Break">Value</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">fileInfo.Mode()</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">prw-rw----</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">os.ModeCharDevice</code></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">c---------</code></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><code class="literal">fileInfo.Mode()</code> <code class="literal">&amp;</code> <span class="No-Break"><code class="literal">os.ModeCharDevice</code></span></p></td><td class="No-Table-Style"><p><code class="literal">----------</code></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p>(<code class="literal">fileInfo.Mode()</code> <code class="literal">&amp;</code> <code class="literal">os.ModeCharDevice)</code> <code class="literal">!= 0</code></p></td><td class="No-Table-Style"><p><span class="No-Break"><code class="literal">FALSE</code></span></p></td></tr></tbody></table>

Figure 8.4 – The code next to its value when standard output is not connected to a TTY

Now the standard output’s value is **prw-rw----**. The `os.ModeNamedPipe`, **p**, bit is set to indicate that it is connected to a **pipe**, a redirection to another process. When the bitwise operation is performed against `os.ModeCharDevice`, we see that no bits are copied over, hence `(fileInfo.Mode() & os.ModeCharDevice) != 0` is **false**, and the device is not a TTY.

### Programmatically check on any operating system

We suggest using a package that has already gone through the trouble of determining the code for a larger set of operating systems to check whether standard output is sent to a TTY. The most popular package we found was [github.com/mattn/go-isatty](https://github.com/mattn/go-isatty), which we used in the `Chapter-8/utils/isatty.go` file:

```markup
package utils
import (
  "fmt"
  "os"
  isatty "github.com/mattn/go-isatty"
)
func IsaTTY() {
  if isatty.IsTerminal(os.Stdout.Fd()) ||  isatty.
     IsCygwinTerminal(os.Stdout.Fd()) {
    fmt.Println("Is a TTY")
  } else {
    fmt.Println("Is not a TTY")
  }
}
```

Now that we know whether we are outputting to a TTY, which indicates that there is a human on the other end, versus not a TTY, we can tailor our output accordingly.

## Designing for a machine

As aforementioned, CLIs were originally designed for machines first. It is important to understand what it exactly means to design another program. Although we would want to tailor our applications toward a human-first design, there will be times when we would need to output in a way that can easily be passed as input to the `grep` or `awk` command, because other applications will expect streams of either plain or JSON text.

Users will be using your CLI in many unexpected ways. Some of those ways are often within a bash script that pipes the output of your command as input into another application. If your application, as it should, outputs in the human-readable format first, it needs to also output in machine-readable format when the standard input is not connected to a TTY terminal. In the latter case, make sure any color and ASCII art, in the form of progress bars, for example, are disabled. The text should also be single-lined tabular data that can easily be integrated with the `grep` and `awk` tools.

Also, it is important that you offer several persistent flags for your users to output in machine-readable output when necessary:

-   `--plain`, for outputting plain text with one record of data per line
-   `--json`, for outputting JSON text that can be piped to and from the curl command
-   `--quiet`, `-q`, or `--silent`, `-s`, for suppressing nonessential output

Provide plain text when it does not impact usability. In other cases, offer the optional previous flags to give the user the ability to pipe its output easily into the input of another.

## Designing for a human

The modern command-line application is designed for its primary consumer—the human. This may seemingly complicate the interface because there’s a bit more to consider. The way data is output and how quickly the data is returned can affect how a user perceives the quality and robustness of your CLI. We’ll go over some key areas of design:

-   Conversation as the norm
-   Empathy
-   Personalization
-   Visual language

Let’s go into each in more detail so we can fully understand how this impacts a human-centred design.

### Conversation as the norm

Since your CLI will be responding to a human and not another program, interaction should flow like a conversation. As an application leans toward a conversational language, the user will feel more at ease. Consider your application as the guide, as well, toward usage of the CLI.

When a user runs a command and is missing important flags or arguments, then your application can prompt for these values. Prompts, or surveys, are a way to include a conversational back-and-forth flow of asking questions and receiving answers from the user. However, prompts should not be a requirement as flags and arguments should be available options for your commands. We will be going over prompts in more detail in [_Chapter 10_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_10.xhtml#_idTextAnchor225), _Interactivity with Prompts and_ _Terminal Dashboards_.

If your application contains a state, then communicate the current state similar to how `git` provides a `status` command and notifies the user when any commands change the state. Similarly, if your application provides workflows, typically defined by a chain of commands, then you can suggest commands to run next.

Being succinct is important when communicating with your user. Just like in conversation, if we muddle our words with too much extraneous information, people can become confused about the point we are trying to make. By communicating what’s important, but keeping it brief, our users will get the most important information quickly.

Context is important. If you are communicating with an end user versus a developer, that makes a difference. In that case, unless you are in verbose mode, there’s no reason to output anything only a developer would understand.

If the user is doing anything dangerous, ask for confirmation and match the level of confirmation with the level of danger that can be invoked by the command:

-   **Mild**:
    -   Example: deleting a file
    -   Confirmation:
        -   If the command is a `delete` command, you don’t need to confirm
        -   If not a `delete` command, prompt for confirmation
-   **Moderate**:
    -   Example: deleting a directory, remote resource, or bulk modification that cannot easily be reverted
    -   Confirmation:
        -   Prompt for confirmation.
        -   Provide a **dry run** operation. A **dry run** operation is used to see the results of the operation without actually making any modifications to the data.
-   **Severe**:
    -   Example: deleting something complex, such as an entire remote application or server
    -   Confirmation:
        -   Prompt for confirmation along with asking them to either type something non-trivial, such as the name of the resource they are deleting, or use a flag such as `–confirm="name-of-resource"` so it is still scriptable

In general, we want to make it increasingly more difficult for the user to do something more difficult. It is a way of guiding the user away from any accidents.

Any user input should always be validated early on to prevent anything unnecessarily bad from happening. Make the error returned understandable to the user who passed in bad data.

In a conversation, any confidential information must be secured. Make sure that any passwords are protected and provide secure methods for users to submit their credentials. For example, consider only accepting sensitive data via files only. You can offer a `–password-file` flag that allows the user to pass in a file or data via standard input. This method provides a discreet method for passing in secret data.

Be transparent in conversation. Any actions that cross the boundaries of the program should be stated explicitly. This includes reading or writing files that the user did not pass in as arguments unless these files are storing an internal state within a cache. This may also include any actions when talking to a remote server.

Finally, response time is more important than speed. Print something to the user in under 100 milliseconds. If you are making a network request, print out something before the request is made so it doesn’t look like the application is hanging or appearing broken. This will make your application appear more robust to its end user.

Let’s revisit our audio metadata CLI project. Under [_Chapter 8_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_08.xhtml#_idTextAnchor166)’s `audiofile` repo, we’ll make some changes to create a conversational flow where it might be missing.

#### Example 1: Prompt for information when a flag is missing

Using the Cobra CLI, if a flag is required, it would automatically return an error if the flag were missing when the command is called. Based on some of the guidelines mentioned in this section, rather than just returning an error, let’s prompt for missing data instead. In the `audiofile` code for [_Chapter 8_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_08.xhtml#_idTextAnchor166), in the `utils/ask.go` file, we create two functions using the survey package [github.com/AlecAivazis/survey/v2](https://github.com/AlecAivazis/survey/v2) as follows:

```markup
func AskForID() (string, error) {
  id := ""
  prompt := &survey.Input{
    Message: "What is the id of the audiofile?",
  }
  survey.AskOne(prompt, &id)
  if id == "" {
    return "", fmt.Errorf("missing required argument: id")
  }
  return id, nil
}
func AskForFilename() (string, error) {
  file := ""
  prompt := &survey.Input{
    Message: "What is the filename of the audio to upload
      for metadata extraction?",
    Suggest: func(toComplete string) []string {
      files, _ := filepath.Glob(toComplete + "*")
      return files
    },
  }
  survey.AskOne(prompt, &file)
  if file == "" {
    return "", fmt.Errorf("missing required argument:
      file")
  }
  return file, nil
}
```

These two functions can now be called when checking the flags that are passed and whether the values are still empty. For example, in the `cmd/get.go` file, we check for the `id` flag value and if it’s still empty, prompt the user for the `id`:

```markup
id, _ := cmd.Flags().GetString("id")
if id == "" {
  id, err = utils.AskForID()
  if err != nil {
    return nil, err
  }
}
```

Running this gives the user the following experience:

```markup
mmontagnino@Marians-MBP audiofile % ./bin/audiofile get
? What is the id of the audiofile?
```

Similarly, in the `cmd/upload.go` file, we check for the filename flag value and if it’s still empty, prompt the user for the filename. Because the prompt allows the user to drill down suggested files, we now get the following experience:

```markup
mmontagnino@Marians-MBP audiofile % ./bin/audiofile upload
? What is the filename of the audio to upload for metadata extraction? [tab for suggestions]
```

Then, press the Tab key for suggestions to reveal a drill-down menu:

```markup
mmontagnino@Marians-MBP audiofile % ./bin/audiofile upload
? What is the filename of the audio to upload for metadata extraction? audio/beatdoctor.mp3 [Use arrows to move, enter to select, type to continue]
 audio/algorithms.mp3
> audio/beatdoctor.mp3
 audio/nightowl.mp3
```

Providing a prompt helps to guide the user and for them to understand how to run the command works.

#### Example 2: Confirm deletion

Another way we can help to guide users toward safely using the CLI and protecting them from making any mistakes is to ask the user for confirmation when doing something dangerous. Although it is not necessary to do so during an explicit delete operation, we created a confirmation function that can be used with a configurable message in any type of dangerous situation. The function exists under the `utils/confirm.go` file:

```markup
func Confirm(confirmationText string) bool {
  confirmed := false
  prompt := &survey.Confirm{
    Message: confirmationText,
  }
  survey.AskOne(prompt, &confirmed)
  return confirmed
}
```

#### Example 3: Notify users when making a network request

Before any HTTP request is made, notifying the user helps them to understand what’s going on, especially if the request hangs or becomes unresponsive. We’ve added a message prior to each network request in each command. The `get` command now has the following line prior to the client running the `Do` method:

```markup
fmt.Printf("Sending request: %s %s %s...\n",
           http.MethodGet, path, payload)
resp, err := client.Do(req)
if err != nil {
  return nil, err
}
```

### Empathy

There are some simple modifications you can make to your command-line application to empathize with your users:

-   Be helpful:
    -   Provide help text and documentation
    -   Suggest commands
    -   Rewrite errors in an understandable way
-   Invite user feedback and bug submission

In [_Chapter 9_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_09.xhtml#_idTextAnchor190), _The Empathic Side of Development_, we will go through the ways in which you can help guide your users toward success using help text, documentation, widespread support, and providing an effortless way for users to provide feedback and submit bugs.

#### Example 1: Offering command suggestions

The Cobra CLI offers some empathy when a user mistypes a command. Let’s look at the following example where the user mistypes `upload` as `upolad`:

```markup
mmontagnino@Marians-MacCourse-Pro audiofile % ./bin/audiofile upolad
Error: unknown command "upolad" for "audiofile"
Did you mean this?
        upload
Run 'audiofile --help' for usage.
```

#### Example 2 – Offer an effortless way to submit bugs

In [_Chapter 9_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_09.xhtml#_idTextAnchor190)_, The Empathic Side of Development,_ we define a bug command that will launch the default browser and navigate to the GitHub repository's new issue page to file a bug report:

```markup
mmontagnino@Marians-MacCourse-Pro audiofile % ./bin/audiofile bug --help
Bug opens the default browser to start a bug report which will include useful system information.
Usage:
  audiofile bug [flags]
Examples:
audiofile bug
```

#### Example 3: Print usage command is used incorrectly

Suppose a user does not input a value to search for when running the search command. The CLI application will prompt for a value to search for. If a value is not passed in by the user, the CLI will output the proper usage of the command:

```markup
mmontagnino@Marians-MacCourse-Pro audiofile % ./bin/audiofile search
?  What value are you searching for?
Error: missing required argument (value)
Usage:
  audiofile search [flags]
Flags:
  -h, --help           help for search
      --json           return json format
      --plain          return plain format
      --value string   string to search for in metadata
```

### Personalization

In general, make the default the right thing for most users, but also allow users to personalize their experience with your CLI. The configuration gives the users a chance to personalize their experience with your CLI and make it more their own.

#### Example 1: Technical configuration with Viper

Using `audiofile` as an example, let’s create a simple configuration setup with Viper to offer the user the ability to change any defaults to their liking. The configurations that we’ve created are for the API and CLI applications. For the API, we’ve defined the `configs/api.json` file, which contains the following:

```markup
{
  "api": {
    "port": 8000
  }
}
```

The API will always execute locally to where it’s being executed. Then, for the CLI, we’ve defined a similar simple file, `configs/cli.json`, containing the following:

```markup
{
  "cli": {
    "hostname": "localhost",
    "port": 8000
  }
}
```

If the API is running on an external host with a different port, then these values can be modified within the configuration. For the CLI to point to the new hostname, we’ll need to update any references within the CLI commands to use the value in the configuration. For example, in the `cmd/get.go` file, the path is defined as:

```markup
path := fmt.Sprintf("http://%s:%d/request?%s",
    viper.Get("cli.hostname"), viper.GetInt("cli.port"),
    params)
```

To initialize these values and provide defaults if any required values are missing from the configuration, we run a `Configure` function defined in `cmd/root.go`:

```markup
func Configure() {
  viper.AddConfigPath("./configs")
  viper.SetConfigName("cli")
  viper.SetConfigType("json")
  viper.ReadInConfig()
  viper.SetDefault("cli.hostname", "localhost")
  viper.SetDefault("cli.port", 8000)
}
```

A similar code exists within the `cmd/api.go` file to gather some of the same information. Now that this is set up, if there are any changes the user wants to make to the hostname, log level, or port, there is only one configuration file to modify.

#### Example 2: Environment variable configuration

Suppose there is an environment variable specific to the application that allows users to define the foreground and background color to use. This environment variable could be named `AUDIOFILE_COLOR_MODE`. Using the Viper configuration again, values for the foreground and background texts may be used to overwrite default settings. While this is not implemented within our CLI, the Viper configuration may look like the following:

```markup
{
  "cli": {
    "colormode": {
      "foreground": "white",
      "background": "black",
    }
  }
}
```

#### Example 3: Storage location

Sometimes users want the location of certain output, logging, for example, to be stored in a particular area. Providing details within Viper can allow defaults to be overwritten. Again, this is not currently implemented within our CLI, but if we were to provide this option within our configuration, it may look like this:

```markup
{
  "api": {
    "local_storage": "/Users/mmontagnino/audiofile"
  }
}
```

Any other new configuration values can be added with a similar approach. Providing the ability to configure your application is the start for personalization. Think of the many ways you can configure your CLI: color settings, disabling prompts or ASCII art, default formatting, and more.

### Pagination

Use a pager when you are outputting a lot of text, but be careful because sometimes the implementation can be error-prone.

#### Pagination for Unix or Linux

On a Unix or Linux machine, you may use the `less` command for pagination. Calling the `less` command with a sensible set of options, such as `less -FIRX`, pagination does not occur if the contents fit on a single screen, case is ignored when searching, color and formatting are enabled, and the content is kept on the screen when `less` quits. We will use this as an example within the next section when outputting table data, and in preparation, within the `utils` package, we add the following files: `pager_darwin.go` and `pager_linux.go`, with a `Pager` function. In our case, though, we use the `–r` flag only because we want to continue displaying colors in the table:

```markup
func Pager(data string) error {
  lessCmd := exec.Command("less", "-r")
  lessCmd.Stdin = strings.NewReader(data)
  lessCmd.Stdout = os.Stdout
  lessCmd.Stderr = os.Stderr
  err := lessCmd.Run()
  if err != nil {
    return err
  }
  return nil
}
```

#### Pagination for Windows

On a Windows machine, we use the `more` command instead. Within the `utils` package, we add the `pager_windows.go` file following with a `Pager` function:

```markup
func Pager(data string) error {
    moreCmd := exec.Command("cmd", "/C", "more")
    moreCmd.Stdin = strings.NewReader(data)
    moreCmd.Stdout = os.Stdout
    moreCmd.Stderr = os.Stderr
    err := moreCmd.Run()
    if err != nil {
        return err
    }
    return nil
}
```

Now you know how to handle the pagination of output on the three major operating systems. This will also help users when you are outputting a large amount of data to scroll through the output easily.

### Visual language

Depending on the data, it might be easier for the users to see it in plain text, table format, or in JSON format. Remember to provide the user with options to return data in the format they prefer with the `–plain` or `–``json` flag.

Note

Sometimes, for all of the data to appear within a user’s window, some lines may be wrapped within a cell. This will break scripts.

There are many visual cues that can be displayed to the user to increase information density. For example, if something is going to take a long time, use a progress bar and provide an estimate of the time remaining. If there is a success or failure, utilize color codes to provide an additional level of information for the user to consume.

We now know how to determine whether we are outputting to a human via a terminal or to another application, so knowing the difference allows us to output data appropriately. Let’s continue to the next section to discuss fun examples to provide data with ASCII visualizations to improve information density.

Just Imagine

# Increasing information density with ASCII art

As the title of this section states, you can increase information density using ASCII art. For example, running the `ls` command shows file permissions in a way a user can easily scan with their eyes and understand with pattern recognition. Also, using a highlighter pen when studying in a textbook to literally highlight a sentence or group of words makes certain phrases jump out as more important. In this section, we’ll talk about some common uses for ASCII art to increase the understanding of the importance of shared information.

## Displaying information with tables

Probably the clearest way that data can be displayed to users is in a table format. Just like the `ls` format, patterns can jump out more easily in a table format. Sometimes records can contain data that is longer than the width of the screen and lines become wrapped. This can break scripts that might be relying on one record per line.

Let’s take our audiofile as an example and instead of returning the JSON output, use the package to return the data cleanly in a table. We can keep the ability to return JSON output for when the user decides to require it using the `–``json` flag.

The simplest way of outputting data as a table with the `pterm` package is using the default table. Next to the models, there currently exists a `JSON()` method that will take the struct and then output it in JSON format. Similarly, we add a `Table()` method on the pointer to the struct. In the `models/audio.go` file, we add the following bit of code for the header table:

```markup
var header = []string{
  "ID",
  "Path",
  "Status",
  "Title",
  "Album",
  "Album Artist",
  "Composer",
  "Genre",
  "Artist",
  "Lyrics",
  "Year",
  "Comment",
}
```

This defines the header for the audio table. We then add some code to transform an `audio` struct into a row:

```markup
func row(audio Audio) []string {
  return []string{
    audio.Id,
    audio.Path,
    audio.Status,
    audio.Metadata.Tags.Title,
    audio.Metadata.Tags.Album,
    audio.Metadata.Tags.AlbumArtist,
    audio.Metadata.Tags.Composer,
    audio.Metadata.Tags.Genre,
    audio.Metadata.Tags.Artist,
    audio.Metadata.Tags.Lyrics,
    strconv.Itoa(audio.Metadata.Tags.Year),
    strings.Replace(audio.Metadata.Tags.Comment, "\r\n",
        "", -1),
  }
}
```

Now we use the `pterm` package to create the table from the header row and function to convert an audio item into a row, each of type `[]string`. The `Table()` method for `Audio` and `AudioList` structs are defined below:

```markup
func (list *AudioList) Table() (string, error) {
  data := pterm.TableData{header}
  for _, audio := range *list {
    data = append(
      data,
      row(audio),
    )
  }
  return pterm.DefaultTable.WithHasHeader()
     .WithData(data).Srender()
}
func (audio *Audio) Table() (string, error) {
  data := pterm.TableData{header, row(*audio)}
  return pterm.DefaultTable.WithHasHeader().WithData(data).
    Srender()
}
```

All the data in this example is output one record per line. If you decide on a different implementation and this is not the case for your code, make sure you add the `–plain` flag as an optional flag where once it is called, it will print one record per line. Doing this will ensure that scripts do not break on the output of the command. Regardless, depending on the size of the data and terminal, you may notice that the data wraps around and might be hard to read. If you are running Unix, run the `tput rmam` command to remove line wrapping from `terminal.app` and then `tput smam` to add line wrapping back in. On Windows, there will be a setting under your console properties. Either way, this should make viewing the table data a bit easier!

If a lot of data is returned within the table, then it’s important to add paging for increased usability. As mentioned in the last section, we’ve added a `Pager` function to the `utils` package. Let’s modify the code so that it checks whether the data is being output to a terminal, and if so, page the data using the `Pager` function. In the `utils/print.go` file, within the `Print` function, we paginate the JSON formatted data, for example, as follows:

```markup
if jsonFormat {
    if IsaTTY() {
        err = Pager(string(b))
        if err != nil {
            return b, fmt.Errorf("\n paging: %v\n ", err)
        }
    } else {
        return b, fmt.Errorf("not a tty")
    }
}
```

If the output is returned to a terminal, then we paginate, otherwise we return the bytes with an error that informs the calling function it is not a terminal. For example, the `cmd/list.go` file calls the preceding `Print` function:

```markup
formatedBytes, err := utils.Print(b, jsonFormat)
if err != nil {
    fmt.Fprintf(cmd.OutOrStdout(), string(formatedBytes))
}
```

When it receives the error, then it just prints the string value to standard output.

## Clarifying with emojis

A picture is worth a thousand words. So much information can be shared just by adding an emoji. For example, think of the simple green checkbox, ![](https://static.packt-cdn.com/products/9781804611654/graphics/image/02.png), that is so often used on Slack or in GitHub to signal approval. Then, there is the opposite case with a red x, ![](https://static.packt-cdn.com/products/9781804611654/graphics/image/03.png), to symbolize that something went wrong.

Emojis are letters that exist within the UTF-8 (Unicode) character set, which covers almost all the characters and symbols of the world. There are websites that will share this Unicode emoji mapping. Visit `https://unicode.org/emoji/charts/full-emoji-list.html` to view the full character list.

### Example 1 – Green checkmark for successful operations

In our audiofile, we add the emoji to the output to the `upload` command. At the top of the file, we add the emoji constant with a UTF-8 character code:

```markup
const (
  checkMark = "\U00002705"
)
```

Then, we use it in the following output:

```markup
fmt.Println(checkMark, " Successfully uploaded!")
fmt.Println(checkMark, " Audiofile ID: ", string(body))
```

Running the upload command after a new recompile and run shows the emoji next to the output, indicating a successful upload. The green checkmark assures the user that everything ran as expected and that there were no errors:

```markup
 Successfully uploaded!
 Audiofile ID: b91a5155-76e9-4a70-90ea-d659c66d39e2
```

### Example 2 – Magnifying glass for search operations

We’ve also added a magnifying glass, ![](https://static.packt-cdn.com/products/9781804611654/graphics/image/013.png), in a similar way when the user runs the search command without the `--value` flag. The new prompt looks like this:

```markup
?  What value are you searching for?
```

### Example 3 – Red for error messages

If there is an invalid operation or an error message, you could also add a red x to symbolize when something goes wrong:

```markup
 Error message!
```

Emojis not only add a fun element to your CLI but also a very valuable one. The little emoji is another way to increase information density and get important points across to the user.

## Using color with intention

Adding color highlights important information for the end user. Don’t overdo it, though; if you end up with multiple different colors frequently used throughout, it’s hard for anything to jump out as important. So use it sparingly, but also intentionally.

An obvious color choice for errors is red, and success is green. Some packages make adding color to your CLI easy. One such package we will use in our examples is `https://github.com/fatih/color`.

Within the audiofile, we look at a few examples where we could integrate colors. For example, the ID for the table that we just listed out. We import the library and then use it to change the color of the `ID` field:

```markup
var IdColor = color.New(color.FgGreen).SprintFunc()
func row(audio Audio) []string {
  return []string{
    IdColor(audio.Id),
    ...
  }
}
```

In the `utils/ask.go` file, we define an `error` function that can be used within the three ask prompts:

```markup
var (
  missingRequiredArumentError =
    func(missingArg string) error {
    return fmt.Errorf(errorColor(fmt.Sprintf("missing
      required argument (%s)", missingArg)))
  }
)
```

The `fmt.Errorf` function receives the `errorColor` function, which is defined within a new `utils/errors.go` file:

```markup
package utils
import "github.com/fatih/color"
var errorColor = color.New(color.BgRed,
  color.FgWhite).SprintFunc()
```

Together, we recompile code and try to run it again, purposely omitting required flags from commands. We see that the command errors out and prints the error with a red background and white foreground, defined by the `color.BgRed` and `color.FgWhite` values. There are many ways to add color. In the `color` package we’re using, the prefix `Fg` stands for foreground and the prefix `Bg` stands for background.

Use colors intentionally, and you will visually transfer the most important information easily to the end user.

## Spinners and progress bars

Spinners and progress bars signify that the command is still processing; the only difference is that progress bars visually display progress. Since it is common to build concurrency into applications, you can also show multiple progress bars running simultaneously. Think about how the Docker CLI often shows multiple files being downloaded simultaneously. This helps the user understand that there’s something happening, progress is made, and nothing is stalling.

### Example 1 – Spinner while playing music

There are different ways that you can add spinners to your Golang project. In the audiofile project, we’ll show a quick way to add a spinner using the `github.com/pterm/pterm` package. In the audiofile project, for each play command distinct for each operating system, we add some code to start and stop the spinner. Let’s look at `play_darwin.go`, for example:

```markup
func play(audiofilePath string) error {
    cmd := exec.Command("afplay", audiofilePath)
    if err := cmd.Start(); err != nil {
        return err
    }
    spinnerInfo := &pterm.SpinnerPrinter{}
    if utils.IsaTTY() {
        spinnerInfo, _ = pterm.DefaultSpinner.Start("Enjoy the 
          music...")
    }
    err := cmd.Wait()
    if err != nil {
        return err
    }
    if utils.IsaTTY() {
        spinnerInfo.Stop()
    }
    return nil
}
```

Running the `play` command for any audio file shows the following output:

```markup
▀ Enjoy the music... (3m54s)
```

It’s hard to capture the spinner in the previous line, but the black box spins around in a circle while the music plays.

### Example 2 – Progress bar when uploading a file

Next, within the `upload` command, we can show code to display the progress of uploading a file. Since the API only uses local flat file storage, the upload goes so quickly it’s hard to see the change in the progress bar, but you can add some `time.Sleep` calls in between each increment to see the progress appear more gradually. Within the `cmd/upload.go` file, we’ve added several statements to create the progress bar and then increment the progress along with title updates:

```markup
p, _ := pterm.DefaultProgressbar.WithTotal(4).WithTitle("Initiating upload...").Start()
```

This first line initiates the progress bar, and then to update the progress bar, the following lines are used:

```markup
pterm.Success.Println("Created multipart writer")
p.Increment()
p.UpdateTitle("Sending request...")
```

Notice that when we first define the progress bar, we call the `WithTotal` method, which takes the total number of steps. This means that for each step where `p.Increment()` is called, the progress bar progresses by 25 percent or 100 divided by the total number of steps. When running a spinner, it’s great to add the visualizer to let the user know that the application is currently running a command that might take a while:

```markup
Process response... [4/4] ███████████             65% | 5s
```

The progress bar gives the user a quick visual of how quickly the command is progressing. It’s a great visual indicator for any command that will take a long time and can be clearly split into multiple steps for progression. Again, spinners and progress bars should not be displayed unless the output is being displayed to the terminal or TTY. Make sure you add a check for TTY before outputting the progress bar or spinner.

## Disabling colors

There are different reasons why color may be disabled for a CLI. A few of these things include:

-   The standard out or standard error pipe is not connected to a TTY or interactive terminal. There is one exception to this. If the CLI is running within a CI environment, such as Jenkins, then color is usually supported, and it is recommended to keep color on.
-   The `NO_COLOR` or `MYAPP_NO_COLOR` environment variable is set to true. This can be defined and set to disable color for all programs that check it or specifically for your program.
-   The `TERM` environment variable is set to dumb.
-   The user passes in the `–``no-color` flag.

Some percentage of your users may be colorblind. Allowing your users to swap out one color for another is a nice way to consider this specific part of your user base. This could be done within the configuration file or application. Allowing them to specify a color and then overwrite it with a preferred color will again allow the user to customize the CLI. This customization will provide users with an improved experience.

Including ASCII art within your application increases information density—a visual indicator that easily helps users to understand some important information. It adds clarity and conciseness. Now let’s discuss a way to make your CLI more intuitive through consistency.

Just Imagine

# Being consistent across CLIs

Learning about command-line syntax, flags, and environment variables requires an upfront cost that pays off in the long run with efficiency if programs are consistent across the board. For example, terminal conventions are ingrained into our fingertips. Reusing these conventions by following preexisting patterns helps to make a CLI more intuitive and guessable. This is what makes users efficient.

There are times when preexisting patterns break usability. As mentioned earlier, a lot of Unix commands don’t return any output by default, which can cause confusion for people who are new to using the terminal or CLI. In this case, it’s fine to break the pattern for the benefit of increased usability.

There are specific topics to consider when maintaining consistency with the larger community of CLIs, but also within the application itself:

-   Naming
-   Positional versus flag arguments
-   Flag naming
-   Usage

## Naming

Use consistent command, subcommand, and flag names to help users intuit your command-line application. Some modern command-line applications, such as the AWS command-line application, will use Unix commands to stay consistent. For example, look at this AWS command:

```markup
aws s3 ls s3://mybucket --summarize
```

The previous command uses the `ls` command to list `S3` objects in the `S3` bucket. It’s important to use common, and non-ambiguous, command names outside of reusing shell commands in your CLI. Take the following as examples that can be logically grouped by type:

![](https://static.packt-cdn.com/products/9781804611654/graphics/image/Table_8.1_B18883.jpg)

Table 8.1 – Example grouping commands by type

These are common names across CLIs. You can also consider integrating some common Unix commands:

-   `cp` (copy)
-   `ls` (list)
-   `mv` (move)

These common command names remove confusion from a long list of ambiguous or unique names. One common confusion is the difference between the update and upgrade commands. It’s best to use one or the other as keeping both will only confuse your users. Also, for the command names that are used often, follow the standard shorthand for these popular commands as well. For example:

-   `-``v`, `--version`
-   `-``h`, `--help`
-   `-``a`, `--all`
-   `-``p`, `--port`

Rather than listing all examples, just consider some of the most common command-line applications you use. Think about which command names make sense for consistency across the board. This will benefit not only your application but the community of command-line applications as a whole as further standards are solidified.

## Positional versus flag arguments

It’s important to stay consistent with arguments and their position. For example, in the AWS CLI, the `s3` argument is consistently next to its arguments:

```markup
aws s3 ls s3://<target-bucket>
aws s3 cp <local-file> <s3-target-location>/<local-file>
```

The consistent position of specific arguments will build a clear pattern that users will follow intuitively.

If flags, that we had mentioned before, are available with one command, they can be available for another command where they make sense. Rather than changing the flag name for each command, stay consistent between commands. Do the same with subcommands. Let’s look at some examples from the GitHub CLI:

```markup
gh codespace list --json
gh issue list –json
```

The GitHub CLI keeps the list subcommand consistent across different commands and reuses the `–json` flag, which has the same behavior across the application.

Note

Required arguments are usually better as positional rather than flags.

## Flag naming

Not only is it important to stay consistent on the position of arguments and the flag names across different commands, but it’s also important to be consistent within the naming. For example, there are flags that can be defined in camel case, `–camelCase`, snake case, `--SnakeCase`, or with dashes, `--flag-with-dashes`. Staying consistent with the way you are naming your flags in your application is also important!

## Usage

In previous chapters, we discussed the grammar of a command and how applications can be defined with a consistent structure: **noun-verb** or **verb-noun**. Staying consistent with the structure also lends to a more intuitive design.

When building your command-line application, if you think about how to stay consistent across other programs and internal to your application, you will create a more intuitive and easier to learn command-line application where your users feel naturally supported.

Just Imagine

# Summary

In this chapter, you learned some specific points to consider when building for a machine versus a human. Machines like simple text and have certain expectations of the data that is returned from other applications. Machine output can sometimes break usability. Designing for humans first, we talked about how we can easily switch to machine-friendly output when needed with the use of some popular flags: `--json`, `--plain`, and `--silence`.

Much goes into a usable design, and we went over some of the ways you can increase the usability of your CLI—from using color with intention, outputting data in tables, paging through long text, and being consistent. All of the aforementioned elements will help the user feel more comfortable and guided when using your CLI, which is one of the main goals we want to achieve. We can summarize with a quick table what a good CLI design looks like versus a bad CLI design:

![Figure 8.5 – Good versus bad CLI design](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_8.5._B18883.jpg)

Figure 8.5 – Good versus bad CLI design

In the next chapter, [_Chapter 9_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_09.xhtml#_idTextAnchor190), _Empathic Side of Development_, we will continue discussing how to develop for humans by incorporating more empathy.

Just Imagine

# Questions

1.  What common flags in scripts can be used with a command-line application to keep the output stable?
2.  What flag should you check to see if the end user does not want color set within the terminal? And what common flag can be used to disable color from the output?
3.  Think about how there could be two commands with similar names and how this adds ambiguity. What ambiguous commands have you come across in your experience of CLIs?

Just Imagine

# Further reading

-   _The Anti-Mac_ _Interface_: [https://www.nngroup.com/articles/anti-mac-interface/](https://www.nngroup.com/articles/anti-mac-interface/)
-   _The Humane Interface: New Directions for Designing Interactive Systems_ by Jef Raskin

Just Imagine

# Answers

1.  `--json` and `--plain` flags keep data consistent and reduce the risk of breaking scripts.
2.  Either the `TERM=dumb`, `NO_COLOR`, or `MYAPP_NO_COLOR` environment variables. The most common flag for disabling color is the `–``no-color` flag.
3.  Update versus upgrade are commonly confused, as well as name and host.

Just Imagine

Previous Chapter