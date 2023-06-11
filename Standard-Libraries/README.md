


# Go Standard Libraries

## About the Course
The book begins with exploring the functionalities available for interaction with the environment and the operating system. We will explore common string operations, date/time manipulations, and numerical problems. We’ll then move on to working with the database, accessing the filesystem, and performing I/O operations. From a networking perspective, we will touch client and server-side solutions. The basics of concurrency are also covered before we wrap up with a few tips and tricks.

By the end of the book, you will have a good overview of the features of the Golang standard library and what you can achieve with them. Finally, you will be proficient in implementing solutions with powerful standard libraries.
## Instructions and Navigation
All of the code is organized into folders. Each folder starts with a number followed by the application name. For example, Chapter02.



The code will look like the following:
```
package main
import (
  "log"
  "runtime"
)
```

Although the Go programming platform is cross-platform, the recipes in the book usually assumes a Unix-based operating system, or at least that some common Unix utilities are available for execution. For Windows users, the Cygwin or GitBash utilities could be useful. The sample code works best with this setup:

* Unix-based environment
* A version of Go equal to or greater than 1.9.2
* An internet connection
* Read, write, and execute permissions on a folder where the sample code will be created and executed
