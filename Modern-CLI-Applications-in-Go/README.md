# Modern CLI Applications in Go


This is the code repository for [Building Modern CLI Applications in Go](https://www.imaginedevops.io/product/building-modern-cli-applications-in-go/9781804611654?utm_source=github&utm_medium=repository&utm_campaign=9781804611654), published by ImagineDevOps .

**Develop next-level CLIs to improve user experience, increase platform usage, and maximize production**

## What is this book about?
Although graphical user interfaces (GUIs) are intuitive and user-friendly, nothing beats a command-line

This book covers the following exciting features:
* Master the Go code structure, testing, and other essentials
* Add a colorful dashboard to your CLI using engaging ASCII banners
* Use Cobra, Viper, and other frameworks to give your CLI an edge
* Handle inputs, API commands, errors, and timeouts like a pro
* Target builds for specific platforms the right way using build tags
* Build with empathy, using easy bug submission and traceback
* Containerize, distribute, and publish your CLIs quickly and easily




## Instructions and Navigations
All of the code is organized into folders. For example, Chapter04.

The code will look like the following:
```
func init() {
    audioCmd.Flags().StringP("filename", "f", "", "audiofile")
    uploadCmd.AddCommand(audioCmd)
}
```

**Following is what you need for this book:**
This book is for beginner- and intermediate-level Golang developers who take an interest in developing CLIs and enjoy learning by doing. You'll need an understanding of basic Golang programming concepts, but will require no prior knowledge of CLI design and development. This book helps you join a community of CLI developers and distribute within the popular Homebrew package management tool.

With the following software and hardware list you can run all code files present in the book (Chapter 1-14).
### Software and Hardware List
| Chapter | Software required | OS required |
| -------- | ------------------------------------ | ----------------------------------- |
| 1-14 | Go 1.19 | Windows, Mac OS X, and Linux (Any) |
| 1-14 | Cobra CLI | Windows, Mac OS X, and Linux (Any) |
| 1-14 | Docker | Windows, Mac OS X, and Linux (Any) |
| 1-14 | Docker Compose | Windows, Mac OS X, and Linux (Any) |
| 1-14 | GoReleaser CLI | Windows, Mac OS X, and Linux (Any) |

