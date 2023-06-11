# Interactivity with Prompts and Terminal Dashboards

One powerful way to increase usability for users is to integrate interactivity with either prompts or terminal dashboards. Prompts are useful because they create a conversational approach while requesting input. Dashboards are useful because they allow developers to create a graphical interface from ASCII characters. That graphical interface, via a dashboard, can create powerful visual cues to allow users to navigate through different commands.

This chapter will give you examples of how to build user surveys from a series of prompts, and a terminal dashboard – whether learning about the Termdash library, designing the mockup, or implementing it for the audio file CLI.

Interactivity is fun. It’s the more human and empathetic approach to a command-line interface. However, remember to disable interactivity if you are not outputting to a terminal. This chapter will cover the basics of surveys and dive deep into the terminal dashboard. By the end of this chapter, you’ll have everything you need to create your own survey or dashboard. We will cover the following:

-   Guiding users with prompts
-   Designing a useful terminal dashboard
-   Implementing a terminal dashboard

# Guiding users with prompts

There are many ways to simply prompt the user, but if you want to create a whole survey that can retrieve information using a variety of different prompts – text input, multi-select, single-select, multi-line text, password, and more – it might be useful to use a preexisting library to handle this for you. Let’s create a generic customer survey using the `survey` package.

To show you how to use this package, I’ll create a survey that can prompt the user for different types of input:

-   **Text input** – for example, an email address
-   **Select** – for example, a user’s experience with the CLI
-   **Multiselect** – for example, any issues encountered
-   **Multiline** – for example, open-ended feedback

In the `Chapter-10` repository, a survey has been written to handle these four prompts. The questions, stored in the `qs` variable, are defined as a slice of `*survey.Question`:

```markup
questions := []*survey.Question{
    {
        Name: "email",
        Prompt: &survey.Input{
          Message: "What is your email address?"
   },
        Validate: survey.Required,
        Transform: survey.Title,
    },
    {
        Name: "rating",
        Prompt: &survey.Select{
            Message: "How would you rate your experience with 
                     the CLI?",
            Options: []string{"Hated it", "Disliked", "Decent", 
                             "Great", "Loved it"},
       },
    },
    {
        Name: "issues",
            Prompt: &survey.MultiSelect{
            Message: "Have you encountered any of these 
                     issues?",
            Options: []string{"audio player issues", "upload 
                             issues", "search issues", "other 
                             technical issues"},
        },
    },
    {
        Name: "suggestions",
        Prompt: &survey.Multiline{
            Message: "Please provide any other feedback or 
                     suggestions you may have.",
        },
    },
}
```

We’ll need an `answers` struct to store all the results from the prompts:

```markup
results := struct {
    Email string
    Rating string
    Issues []string
    Suggestions string
}{}
```

And finally, the method that asks the questions and stores the results:

```markup
err := survey.Ask(questions, &results)
if err != nil {
    fmt.Println(err.Error())
    return
}
```

Now that we’ve created the survey, we can try it out:

```markup
mmontagnino@Marians-MacCourse-Pro Chapter-10 % go run main.go
? What is your email? mmontagnino@gmail.com
? How would you rate your experience with the CLI? Great
? Have you encountered any of these issues? audio player issues, search issues
? Please provide any other feedback or suggestions you may have. [Enter 2 empty lines to finish]I want this customer survey embedded into the CLI and email myself the results!
```

Prompting the user is an easy way to integrate interactivity into your command-line application. However, there are even more colorful and fun ways to interact with your users. In the next section, we’ll discuss the terminal dashboard, the `termdash` package in detail, and how to mock up and implement a terminal dashboard.

Just Imagine

# Designing a useful terminal dashboard

Command-line interfaces don’t have to be limited to text. With **termdash**, a popular Golang library, you can build a terminal dashboard providing users with a user interface to visually see progress, alerts, text, and more. Colorful widgets placed within a clean dashboard that’s been neatly laid out can increase information density and present a lot of information to the user in a very user-friendly manner. In this section, we’ll learn about the library and the different layout choices and widget options. At the end of the chapter, we’ll design a terminal dashboard that we can implement in our **audio file** command-line interface.

## Learning about Termdash

Termdash is a Golang library that provides a customizable and cross-platform, terminal-based dashboard. On the project’s GitHub page, a fun and colorful demo provides an example of all possible widgets demonstrated within a dynamic layout. From the demo, you can see that you can go all out on a fancy dashboard. To do so, you’ll need to understand how to lay out a dashboard, interact with keyboard and mouse events, add widgets, and fine-tune the appearance with alignment and color. Within this section, we will break down the layers of a Termdash interface and the widgets that can be organized within it.

A Termdash dashboard consists of four main layers:

-   The terminal layer
-   The infrastructure layer
-   The container layer
-   The widgets layer

Let’s take a deep dive into each of them.

### The terminal layer

Think of the terminal layer of a dashboard as a 2D grid of cells that exist within a buffer. Each cell contains either an ASCII or Unicode character with the option to customize the foreground color, the color of text, the background color, or the color of the non-character space within the cell. Interactions with the mouse and keyboard happen on this layer as well.

Two terminal libraries can be used to interact at the cell level of a terminal:

-   **tcell**: Inspired by **termbox** and has many new improvements
-   **termbox**: No longer supported, although it is still an option

The following examples will utilize the `tcell` package to interact with the terminal. To start, create a new `tcell` instance to interact via the terminal API:

```markup
terminalLayer, err := tcell.New()
if err != nil {
   return err
}
defer terminalLayer.Close()
```

Notice that in this example, `tcell` has two methods: `New` and `Close`. `New` creates a new `tcell` instance in order to interact with the terminal and `Close` closes the terminal. It’s a good practice to defer closing access to `tcell` right after creation. Although there are no options passed into the `New` method, there are a few optional methods that can be called:

-   `ColorMode` sets the color mode when initializing a terminal
-   `ClearStyle` sets the foreground and background color when a terminal is cleared

An example of initializing a cell in `ColorMode` to access all 256 available terminal colors would look like this:

```markup
terminalLayer, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
if err != nil {
   return err
}
defer terminalLayer.Close()
```

`ClearStyle`, by default, will use `ColorDefault` if no specific `ClearStyle` is set. This `ColorDefault` is usually the default foreground and background colors of the terminal emulator, which are typically black and white. To set a terminal to use a yellow foreground and navy background style when the terminal is cleared, the `New` method, which accepts a slice of options, would be modified in the following way:

```markup
terminalLayer, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256), tcell.ClearStyle(cell.ColorYellow, cell.ColorNavy))
if err != nil {
   return err
}
defer terminalLayer.Close()
```

Now that we’ve created a new `tcell` that gives us access to the Terminal API, let’s discuss the next layer – infrastructure.

### The infrastructure layer

The infrastructure of a terminal dashboard provides the organization of the structure. The three main elements of the infrastructure layer include alignment, line style, and Termdash.

#### Alignment

Alignment is provided by the `align` package, which provides two alignment options – `align.Horizonal`, which includes predefined values of `left`, `center`, and `right` and `align.Vertical` with predefined values of `top`, `middle`, and `bottom`.

#### Line style

The line style defines the style of the line drawn on the terminal either when drawing boxes or borders.

The package exposes the options available via `LineStyle`. The `LineStyle` type represents a style that follows the Unicode options.

#### Termdash

Termdash provides the developer with the main entry point. Its most important purpose is to start and stop the dashboard application, control screen refreshing, process any runtime errors, and subscribe and listen for keyboard and mouse events. The `termdash.Run` method is the simplest way to start a Termdash application. The terminal may run until the context expires, a keyboard shortcut is called, or it times out. The simplest way to get started with the dashboard is with the following minimal code example, which creates a new `tcell` for the terminal layer, and a new **container** for the container layer. A container is another module within the `termdash` package, which we will dive into in the next section. We create context with a 2-minute timeout and then call the `Run` method of the `termdash` package:

```markup
if terminalLayer, err := tcell.New()
if err != nil {
   return err
}
defer terminalLayer.Close()
containerLayer, err := container.New(terminalLayer)
if err != nil {
   return err
}
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()
if err := termdash.Run(ctx, terminalLayer, containerLayer); err != nil {
   return err
}
```

In the preceding code example, the dashboard will run until the context expires, in 60 seconds.

Screen redrawing, or refreshing, for your Terminal dashboard can be done in a few ways: periodic, time-based redraws or manually triggered redraws. Only one method may be used, as using one means the other method is ignored. Besides that, the screen will refresh each time an input event occurs. The `termdash.RedrawInterval` method is an option that can be passed into the `Run` method to tell the dashboard application to redraw, or refresh, the screen at a particular interval. The `Run` method can be modified with the option to refresh every 5 seconds:

```markup
termdash.Run(ctx, terminalLayer, containerLayer, termdash.RedrawInterval(5*time.Second))
```

The dashboard may also be redrawn using a controller, which can be triggered manually. This option means that the dashboard is drawn only once and unlike the `Run` method, the user maintains control of the main goroutine. An example of this code, using the previously defined `tcell` and `container` variables defined earlier, can be passed into a new controller to be drawn manually:

```markup
termController, err := termdash.NewController(terminalLayer, containerLayer)
if err != nil {
    return err
}
defer termController.Close()
if err := termController.Redraw(); err != nil {
    return fmt.Errorf("error redrawing dashboard: %v", err)
}
```

The Termdash API provides a `termdash.ErrorHandler` option, which tells the dashboard how to handle errors gracefully. Without providing an implementation for this error handler, the dashboard will panic on all runtime errors. Errors can occur when processing or retrieving events, subscribing to an event, or when a container fails to draw itself.

An error handler is a callback method that receives an error and handles the error appropriately. It can be defined as a variable and, in the simplest case, just prints the runtime error:

```markup
errHandler := func(err error) {
   fmt.Printf("runtime error: %v", err)
}
```

When starting a Termdash application using the `Run` or `NewController` method, the error handler may be passed in as an option using the `termdash.ErrorHandler` method. For example, the `Run` method can be modified with a new option:

```markup
termdash.Run(ctx, terminalLayer, containerLayer, termdash.ErrorHandler(errHandler))
```

While the `NewController` method can be modified similarly:

```markup
termdash.NewController(terminalLayer, containerLayer, termdash.ErrorHandler(errHandler))
```

Through the `termdash` package, you can also subscribe to keyboard and mouse events. Typically, the container and certain widgets subscribe to keyboard and mouse events. Developers can also subscribe to certain mouse and keyboard events to take global action. For example, a developer may want the terminal to run a specific function when a specific key is set. `termdash.KeyboardSubscriber` is used to implement this functionality. With the following code, the user subscribes to the letters `q` and `Q` and responds to the keyboard events by running code to quit the dashboard:

```markup
keyboardSubscriber := func(k *terminalapi.Keyboard) {
    switch k.Key {
      case 'q':
      case 'Q':
          cancel()
    }
}
if err := termdash.Run(ctx, terminalLayer, containerLayer, termdash.KeyboardSubscriber(keyboardSubscriber)); err != nil {
return fmt.Errorf("error running termdash with keyboard subscriber: %v", err)
}
```

Another option is to call the `Run` method with the option to listen to mouse events using `termdash.MouseSubscriber`. Similarly, the following code can be called to do something random when the mouse button is clicked within the dashboard:

```markup
mouseClick := func(m *terminalapi.Mouse) {
    switch m.Button {
        case mouse.ButtonRight:
        // when the left mouse button is clicked - cancel
        cancel()
        case mouse.ButtonLeft:
        // when the left mouse button is clicked
        case mouse.ButtonMiddle:
        // when the middle mouse button is clicked
    }
}
if err := termdash.Run(ctx, terminalLayer, containerLayer, termdash.MouseSubscriber(mouseClick)); err != nil {
    return fmt.Errorf("error running termdash with mouse subscriber: %v", err)
}
```

### The container layer

The container layer provides options for dashboard layouts, container styles, keyboard focus, and margin and padding. It also provides a method for placing a widget within a container.

From the previous examples, we see that a new container is called using the `container.New` function. We’ll provide some new examples of how to organize your container and set it up with different layouts.

There are two main layout options:

-   Binary tree
-   Grid layouts

The **binary tree layout** organizes containers in a binary tree structure where each container is a node in a tree, which, unless empty, may contain either two sub-containers or a widget. Sub-containers can be split further with the same rules. There are two kinds of splits:

-   **Horizontal splits**, created with the `container.SplitHorizontal` method, will create top and bottom sub-containers specified by `container.Top` and `container.Bottom`
-   **Vertical splits**, created with the `container.SplitVertical` method, will create left and right sub-containers, specified by `container.Left` and `container.Right`

The `container.SplitPercent` option specifies the percentage of container split to use when spitting either vertically or horizontally. When the split percentage is not specified, a default of 50% is used. The following is a simple example of a binary tree layout using all the methods described:

```markup
    terminalLayer, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256),
        tcell.ClearStyle(cell.ColorYellow, cell.ColorNavy))
    if err != nil {
        return fmt.Errorf("tcell.New => %v", err)
    }
    defer terminalLayer.Close()
leftContainer := container.Left(
container.Border(linestyle.Light),
)
rightContainer :=
container.Right(
container.SplitHorizontal(
container.Top(
container.Border(linestyle.Light),
),
container.Bottom(
container.SplitVertical(
     container.Left(
     container.Border(linestyle.Light),
     ),
     container.Right(
     container.Border(linestyle.Light),
     ),
     ),
      ),
    )
)
containerLayer, err := container.New(
terminalLayer,
container.SplitVertical(
leftContainer,
rightContainer,
container.SplitPercent(60),
),
)
```

Notice how we drill down when splitting up the terminal into containers. First, we split vertically to divide the terminal into left and right portions. Then, we split the right portion horizontally. The bottom-right horizontally split portion is split vertically. Running this code will present the following dashboard:

![Figure 10.1 – Dashboard showing a container split using the binary layout](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.1_B18883.jpg)

Figure 10.1 – Dashboard showing a container split using the binary layout

Notice that the container to the left takes up about 60% percent of the full width. The other splits do not define a percentage and take up 50% of the container.

The other option for a dashboard is to use a **grid layout**, which organizes the layout into rows and columns. Unlike the binary tree layout, the grid layout requires a grid builder object. Rows, columns, or widgets are then added to the grid builder object.

Columns are defined using either the `grid.ColWidhPerc` function, which defines a column with a specified width percentage of the parent’s width, or `grid.ColWidthPercWithOpts`, which is an alternative that allows developers to additionally specify options when representing the column.

Rows are defined using either the `grid.RowHeightPerc` function, which defines a row with a specified height percentage of the parent’s height, or `grid.RowHeightPercWithOpts`, which is an alternative that allows developers to additionally specify options when representing the row.

To add a widget within the grid layout, utilize the `grid.Widget` method. The following is a simple example of a layout implemented by the `grid` package. The code uses all the related methods and adds an ellipses text widget within each cell:

```markup
    t, err := tcell.New()
    if err != nil {
        return fmt.Errorf("error creating tcell: %v", err)
    }
    rollingText, err := text.New(text.RollContent())
    if err != nil {
        return fmt.Errorf("error creating rolling text: %v", 
          err)
    }
    err = rollingText.Write("...")
    if err != nil {
        return fmt.Errorf("error writing text: %v", err)
    }
    builder := grid.New()
    builder.Add(
        grid.ColWidthPerc(60,
            grid.Widget(rollingText,
                container.Border(linestyle.Light),
            ),
        ),
    )
    builder.Add(
        grid.RowHeightPerc(50,
            grid.Widget(rollingText,
                container.Border(linestyle.Light),
            ),
        ),
    )
    builder.Add(
        grid.ColWidthPerc(20,
            grid.Widget(rollingText,
                container.Border(linestyle.Light),
            ),
        ),
    )
    builder.Add(
        grid.ColWidthPerc(20,
            grid.Widget(rollingText,
                container.Border(linestyle.Light),
            ),
        ),
    )
    gridOpts, err := builder.Build()
    if err != nil {
        return fmt.Errorf("error creating builder: %v", err)
    }
    c, err := container.New(t, gridOpts...)
```

Running the code generates the following dashboard:

![Figure 10.2 – Dashboard showing the container created using the grid layout](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.2_B18883.jpg)

Figure 10.2 – Dashboard showing the container created using the grid layout

Notice that the column width percentage equals 100%; anything more would cause a compilation error.

There is also the option of a dynamic layout that allows you to switch between different layouts on the dashboard. Using the `container.ID` option, you can identify a container with some text, which can be referenced later so there’s a way to identify which container will be dynamically updated using the `container.Update` method:

```markup
    t, err := tcell.New()
    if err != nil {
        return fmt.Errorf("error creating tcell: %v", err)
    }
    defer t.Close()
    b1, err := button.New("button1", func() error {
        return nil
    })
    if err != nil {
        return fmt.Errorf("error creating button: %v", err)
    }
    b2, err := button.New("button2", func() error {
        return nil
    })
    if err != nil {
        return fmt.Errorf("error creating button: %v", err)
    }
    c, err := container.New(
        t,
        container.PlaceWidget(b1),
        container.ID("123"),
    )
    if err != nil {
        return fmt.Errorf("error creating container: %v", err)
    }
    update := func(k *terminalapi.Keyboard) {
        if k.Key == 'u' || k.Key == 'U' {
            c.Update(
                "123",
                container.SplitVertical(
                    container.Left(
                        container.PlaceWidget(b1),
                    ),
                    container.Right(
                        container.PlaceWidget(b2),
                    ),
                ),
            )
        }
    }
    ctx, cancel := context.WithTimeout(context.Background(), 
      5*time.Second)
    defer cancel()
    if err := termdash.Run(ctx, t, c, termdash.
       KeyboardSubscriber(update)); err != nil {
        return fmt.Errorf("error running termdash: %v", err)
    }
```

In this code, the container ID is set to `123`. Originally, the widget contained just one button. The `update` method replaces the single button with a container split vertically, with one button on the left and another on the right. When running this code, pressing the _u_ key runs the update on the layout.

The original layout shows a single button:

![Figure 10.3 – Layout showing a single button](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.3_B18883.jpg)

Figure 10.3 – Layout showing a single button

After pressing the _u_ or _U_ key, the layout updates:

![Figure 10.4 – Layout showing two buttons after pressing the u key again](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.4_B18883.jpg)

Figure 10.4 – Layout showing two buttons after pressing the u key again

The container layer can be further configured using margin and padding settings. The margin is the space outside of the container’s border while the padding is the space between the inside of the container’s border and its content. The following image provides the best visual representation of margins and padding:

![Figure 10.5 – Margin and padding](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.5_B18883.jpg)

Figure 10.5 – Margin and padding

The margin and padding can be set with either absolute or relative values. An absolute margin can be set with the following options:

-   `container.MarginTop`
-   `container.MarginRight`
-   `container.MarginBottom`
-   `container.MarginLeft`

Absolute padding can be set with the following options:

-   `container.PaddingTop`
-   `container.PaddingRight`
-   `container.PaddingBottom`
-   `container.PaddingLeft`

Relative values for the margin and padding are set with percentages. The margin and padding’s top and bottom percentage values are relative to the container’s height:

-   `container.MarginTopPercent`
-   `container.MarginBottomPercent`
-   `container.PaddingTopPercent`
-   `container.PaddingBottomPercent`

The margin and padding’s right and left percentage values are relative to the container’s width:

-   `container.MarginRightPercent`
-   `container.MarginLeftPercent`
-   `container.PaddingRightPercent`
-   `container.PaddingLeftPercent`

Another form of placement within containers is alignment. The following methods are available from the align API to align content within the container:

-   `container.AlignHorizontal`
-   `container.AlignVertical`

Let’s put it all together in a simple example that extends upon the binary tree code example:

```markup
b, err := button.New("click me", func() error {
    return nil
})
if err != nil {
    return err
}
leftContainer :=
container.Left(
     container.Border(linestyle.Light),
           container.PlaceWidget(b),
           container.AlignHorizontal(align.HorizontalLeft),
     )
rightContainer :=
         container.Right(
             container.SplitHorizontal(
                 container.Top(
                    container.Border(linestyle.Light),
                    container.PlaceWidget(b),
                    container.AlignVertical(align.VerticalTop),
                 ),
                 container.Bottom(
                   container.SplitVertical(
                        container.Left(
                          container.Border(linestyle.Light),
                               container.PlaceWidget(b),
                               container.PaddingTop(3),
                               container.PaddingBottom(3),
                               container.PaddingRight(3),
                               container.PaddingLeft(3),
                         ),
                         container.Right(
                           container.Border(linestyle.Light),
                             container.PlaceWidget(b),
                             container.MarginTop(3),
                             container.MarginBottom(3),
                             container.MarginRight(3),
                             container.MarginLeft(3),
                        ),
                    ),
                ),
           ),
                )
containerLayer, err := container.New(
        terminalLayer,
        container.SplitVertical(
            leftContainer,
            rightContainer,
            container.SplitPercent(60),
        ),
    )
```

The resulting layout appears as follows:

![Figure 10.6 – Container showing different alignments for a button, with different margins and padding](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.6_B18883.jpg)

Figure 10.6 – Container showing different alignments for a button, with different margins and padding

You can also define a key to change the focus to the next or previous container using the `container.KeyFocusNext` and `container.KeyFocusPrevious` options.

### The widget layer

In several of the previous examples, we showed code that placed a widget in either a grid or binary tree container layout and also customized the alignment, margin, and padding. However, besides a simple button or text, there are different widget options, and the demo on the GitHub page shows an example of each:

![Figure 10.7 – Termdash sample screenshot showing all the widgets in a dashboard](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.7_B18883.jpg)

Figure 10.7 – Termdash sample screenshot showing all the widgets in a dashboard

Let’s do a quick example of each with a snippet of code to understand how each widget is created. To add each widget to a container, just use the `container.PlaceWidget` method that was used earlier for the simple text and button examples. Let’s go over a few other examples: a bar chart, donut, and gauge. For a detailed code of the other widgets, visit the very well-documented termdash wiki and check out the demo pages.

#### A bar chart

Here is some example code for creating a bar chart widget with individual values displayed relative to a `max` value:

```markup
    barChart, err := barchart.New()
    if err != nil {
        return err
    }
    values := []int{20, 40, 60, 80, 100}
    max := 100
    if err := barChart.Values(values, max); err != nil {
        return err
    }
```

The preceding code creates a new `barchart` instance and adds the values, a slice of `int`, plus the maximum `int` value. The resulting terminal dashboard looks like this:

![Figure 10.8 – Bar chart example](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.8_B18883.jpg)

Figure 10.8 – Bar chart example

Change the values of the `values` and `max` variables to see the chart change. The color of the bars can also be modified based on preference.

#### A donut

A donut, or progress circle chart, represents the completion of progress. Here is some example code for creating a donut chart to show percentages:

```markup
    greenDonut, err := donut.New(
        donut.CellOpts(cell.FgColor(cell.ColorGreen)),
        donut.Label("Green", cell.FgColor(cell.ColorGreen)),
    )
    if err != nil {
        return err
    }
    greenDonut.Percent(75)
```

The preceding code creates a new `donut` instance with options for the label and foreground color set to green. The resulting terminal dashboard looks like this:

![Figure 10.9 – Green donut at 75%](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.9_B18883.jpg)

Figure 10.9 – Green donut at 75%

Again, the color can be modified based on preference, and remember, since Termdash provides dynamic refreshing, the data can be automatically updated and redrawn, making it quite nice for showing progress.

#### A gauge

A gauge, or progress bar, is another way to measure the amount completed. The following is some sample code for showing how to create a progress gauge:

```markup
    progressGauge, err := gauge.New(
        gauge.Height(1),
        gauge.Border(linestyle.Light),
        gauge.BorderTitle("Percentage progress"),
    )
    if err != nil {
        return err
    }
    progressGauge.Percent(75)
```

This code creates a new instance of a gauge with options for a light border, a title, **Percentage progress**, and a slim height of `1`. The percentage, as with the donut, is 75%. The resulting terminal dashboard looks like this:

![Figure 10.10 – Gauge at 75% percent progress](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.10_B18883.jpg)

Figure 10.10 – Gauge at 75% percent progress

As mentioned before, because of dynamic redrawing, this is another great option for showing progress updates.

Now that we’ve shown examples of different widgets to include within a terminal dashboard, let’s sketch out a design using these widgets that we can later implement in our audio file command-line interface. Suppose we wanted to build a music player in a terminal dashboard. Here is a sample layout:

![Figure 10.11 – Terminal dashboard layout](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.11_B18883.jpg)

Figure 10.11 – Terminal dashboard layout

This layout can be created easily with the binary layout. The music library list section can be generated from a list of songs with number identifiers, which can be used in the text input section, where a song can be selected by ID. Any error messages associated with the input ID will be displayed right below. If the input is good, the selected song section will show rolling ASCII text with the song title, and the metadata section will display the text metadata of the selected song. Hitting the play button will start playing the selected song, and the stop button will stop it. Proceed to the next section where we’ll make this terminal dashboard a reality.

Just Imagine

# Implementing a terminal dashboard

When creating a terminal dashboard, you can create it as a separate standalone application or as a command that is called from the command-line application. In our specific example for the player terminal dashboard, we are going to call the dashboard when the `./bin/audiofile player` command is called.

First, from the audio file’s root repository, we’ll need to use `cobra-cli` to create the command:

```markup
cobra-cli add player
Player created at /Users/mmontagnino/Code/src/github.com/marianina8/audiofile
```

Now, we can create the code to generate the terminal dashboard, called within the `Run` field of the `player` command. Remember that the terminal dashboard consists of four main layers: the terminal, infrastructure, container, and widgets. Like a painting, we’ll start with the base layer: the terminal.

## Creating the terminal layer

The first thing you need to do is to create a terminal that provides access to any input and output. Termdash has a `tcell` package for creating a new `tcell`\-based terminal. Many terminals by default only support 16 colors, but other more modern terminals can support up to 256 colors. The following code specifically creates a new terminal with a 265-color mode.

```markup
t, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
```

After creating a terminal layer, we then create the infrastructure layer.

## Creating the infrastructure layer

The infrastructure layer handles the terminal setup, mouse and keyboard events, and containers. In our terminal dashboard player, we want to handle a few tasks:

-   Keyboard event to signal quitting
-   Running the terminal dashboard, which subscribes to this keyboard event

Let’s write the code to handle these two features required of the terminal dashboard.

### Subscribing to keyboard events

If we want to listen for key events, we create a keyboard subscriber to specify the keys to listen to:

```markup
quitter := func(k *terminalapi.Keyboard) {
    if k.Key == 'q' || k.Key == 'Q' {
        ...
    }
}
```

Now that we have defined a keyboard subscriber, we can use this as an input parameter to termdash’s `Run` method.

### Running the terminal

When running the terminal, you’ll need the terminal variable, container, and keyboard and mouse subscribers, as well as the timed redrawing interval and other options. The following code runs the `tcell`\-based terminal we created and the `quitter` keyboard subscriber, which listens for _q_ or _Q_ key events to quit the application:

```markup
if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(100*time.Millisecond)); err != nil {
    panic(err)
}
```

The `c` variable that’s passed into the `termdash.Run` method as the third parameter is the container. Let’s define the container now.

## Creating the container layer

When creating the container, it helps to look at the bigger picture of the layout and then narrow it down as you go. For example, when you first look at the planned layout, you’ll see the largest sections are made from left and right vertical splits.

![Figure 10.12 – Initial vertical split](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.12_B18883.jpg)

Figure 10.12 – Initial vertical split

As we begin to define the container, we’ll slowly drill down with more specifics, but we begin with the following:

-   **Vertical Split (Left)** – The music library
-   **Vertical Split (Right)** – All other widget

The final code reflects this drill-down process. Since we keep the left vertical split as the music library, we drill down with containers on the left, always starting with the larger containers and adding smaller ones within.

![Figure 10.13 – Horizontal split of right vertical space](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.13_B18883.jpg)

Figure 10.13 – Horizontal split of right vertical space

The next is a horizontal split that separates the left vertical split into the following:

-   **Horizontal Split (Top) 30%** – Text input, error messages, and the rolling song title text
-   **Horizontal Split (Bottom) 70%** – Metadata and play/stop buttons

Let’s take the top horizontal split and split it, again, horizontally:

-   **Horizontal Split (Top) 30%** – Text input and error message
-   **Horizontal Split (Bottom) 70%** – The rolling song title text

![Figure 10.14 – Horizontal split of top horizontal space](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.14_B18883.jpg)

Figure 10.14 – Horizontal split of top horizontal space

We split the earlier top part horizontally into the separated text input and error messages:

-   **Horizontal Split (Top) 60%** – Text input
-   **Horizontal Split (Bottom) 40%** – Error messages

![Figure 10.15 – Horizontal split of top horizontal space](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.15_B18883.jpg)

Figure 10.15 – Horizontal split of top horizontal space

Now, let’s drill down into the bottom 70% of the initial horizontal split of the right vertical container. Let’s split it up into two horizontal sections:

-   **Horizontal Split (Top) 80%** – The metadata section
-   **Horizontal Split (Bottom) 20%** – The button section (play/stop)

![Figure 10.16 – Horizontal split of bottom horizontal space](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.16_B18883.jpg)

Figure 10.16 – Horizontal split of bottom horizontal space

Finally, the last part to drill down to is the bottom horizontal split, which we will split vertically:

-   **Vertical Split (Left) 50%** – The play button
-   **Vertical Split (Right) 50%** – The stop button

![Figure 10.17 – Vertical split of bottom horizontal space](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.17_B18883.jpg)

Figure 10.17 – Vertical split of bottom horizontal space

The entire layout broken down with the container code shows this drill-down process – I’ve added comments for where the widgets will be placed for reference:

```markup
c, err := container.New(
    t,
    container.SplitVertical(
        container.Left(), // music library
        container.Right(
            container.SplitHorizontal(
                container.Top(
                    container.SplitHorizontal(
                        container.Top(
                            container.SplitHorizontal(
                                container.Top(), // text input
                                container.Bottom(), // error 
                                                    msgs
                                container.SplitPercent(60),
                            ),
                        ),
                        container.Bottom(), // rolling song 
                                            title
                        container.SplitPercent(30),
                    ),
                ),
                container.Bottom(
                    container.SplitHorizontal(
                        container.Top(), // metadata
                        container.Bottom(
                            container.SplitVertical(
                                container.Left(), // play 
                                                  button
                                container.Right(), // stop 
                                                   button
                            )
                        ),
                        container.SplitPercent(80),
                    ),
                ),
                container.SplitPercent(30),
            ),
        ),
    ),
)
```

Next, let’s create the widgets and place them within the appropriate containers to finalize the terminal dashboard.

## Creating the widgets layer

Going back to the original layout, all the different widgets we’ll need to implement are clear to see:

-   The music library list
-   Input text
-   Error messages
-   Rolling text – selected song (title by artist)
-   Metadata
-   The play button
-   The stop button

At this point, I am aware of which widget to use for each item on the list. However, if you have not yet decided, now is the time to determine the best Termdash widget to use for each item:

-   Text:
    -   Music library list
    -   Error messages
    -   Rolling text – selected song (title by artist), metadata
-   Text input:
    -   Input field
-   Button:
    -   The play button
    -   The stop button

Let’s create at least one of each type as an example. The full code is available in the `Chapter10` GitHub repository.

### Creating a text widget for the music library list

The music library list will take in the audio list and print the text in a section that will list the index of the song next to the title and artist. We define this widget with the following function:

```markup
func newLibraryContent(audioList *models.AudioList) (*text.Text, error) {
    libraryContent, err := text.New(text.RollContent(), text.
      WrapAtWords())
    if err != nil {
        panic(err)
    }
    for i, audiofile := range *audioList {
        libraryContent.Write(fmt.Sprintf("[id=%d] %s by %s\n", 
          i, audiofile.Metadata.Tags.Title, audiofile.Metadata.
          Tags.Artist))
    }
    return libraryContent, nil
}
```

The function is called in the `Run` function field like so:

```markup
libraryContent, err := newLibraryContent(audioList)
```

The error message and metadata items are also text widgets, so we’ll omit those code examples. Next, we’ll create the input text.

### Creating an input text widget for setting the current ID of a song

The input text section is where a user inputs the ID of the song displayed in the music library section. The input text is defined within the following function:

```markup
func newTextInput(audioList *models.AudioList, updatedID chan<- int, updateText, errorText chan<- string) *textinput.TextInput {
    input, _ := textinput.New(
        textinput.Label("Enter id of song: ", cell.
          FgColor(cell.ColorNumber(33))),
        textinput.MaxWidthCells(20),
        textinput.OnSubmit(func(text string) error {
            // set the id
            // set any error text
        return nil
    }),
    textinput.ClearOnSubmit(),
    )
    return input
}
```

### Creating a button to start playing the song associated with the input ID

The last type of widget is a button. There are two different buttons we need, but the following code is for the play button:

```markup
func newPlayButton(audioList *models.AudioList, playID <-chan int) (*button.Button, error) {
    playButton, err := button.New("Play", func() error {
        stopTheMusic()
        }
        go func() {
        if audiofileID <= len(*audioList)-1 && audiofileID >= 0 {
        pID, _ = play((*audioList)[audiofileID].Path, false, 
                     true)
        }}()
        return nil
    },
    button.FillColor(cell.ColorNumber(220)),
    button.GlobalKey('p'),
    )
    if err != nil {
        return playButton, fmt.Errorf("%v", err)
    }
    return playButton, nil
}
```

-   The function is called in the `Run` function field:

```markup
playButton, err := newPlayButton(audioList, playID)
```

-   Once all the widgets have been created, they are placed within the container in the appropriate places with the following line of code:

```markup
container.PlaceWidget(widget)
```

-   Once the widgets have been placed within the container, we can run the terminal dashboard with the following command:

```markup
./bin/audiofile player
```

-   Magically, the player terminal dashboards appear and we can select an ID to enter and play a song:

![Figure 10.18 – Audio file player terminal dashboard](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_10.18_B18883.jpg)

Figure 10.18 – Audio file player terminal dashboard

-   Voila! We’ve created a terminal dashboard to play the music in our audio file library. While you can view the metadata through the command-line application’s `get` and `list` commands and play music with the `play` command, the new player terminal dashboard allows you to view what exists in the audio file library in a more user-friendly fashion.

Just Imagine

# Summary

In this chapter, you learned how to create a survey with different interactive prompts and a terminal dashboard containing a variety of widgets. These are just examples that can hopefully inspire you in terms of interactivity within your own command-line application.

The survey example showed you how to use a variety of different types of prompts; you can prompt the user for their user experience, but as you’ve seen within the audio file CLI, you can also just prompt for missing information. These prompts can be input throughout your code in places where prompts may come in handy, or they can be strung along a list of other questions and you can create a more thorough survey for your users.

The player terminal dashboard gives you an example of how to create a terminal dashboard for a command-line interface. Consider the kind of data your users will be sending or retrieving from your command-line interface and let that guide you in your design of a more visual approach.

Just Imagine

# Questions

1.  What method is used to create the terminal layer?
2.  What method is used to place a widget inside a container?
3.  What’s the difference between the binary layout and the grid layout?

Just Imagine

# Answers

1.  `tcell.New()`
2.  `container.PlaceWidget(widget)`
3.  The grid layout allows you to split the container into horizontal rows and vertical columns. The binary layout allows you to split sub-containers horizontally or vertically.

Just Imagine

# Further reading

-   _The Big Course of Dashboards: Visualizing Your Data Using Real-World Business Scenarios_ by Wexler, Shaffer, and Cotgreave