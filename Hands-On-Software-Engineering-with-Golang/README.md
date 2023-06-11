

# Hands-on Software Engineering with Golang

This is the code repository for [Hands-on Software Engineering with Golang](https://www.imaginedevops.io/in/programming/hands-on-software-engineering-with-golang?utm_source=github&utm_medium=repository&utm_campaign=), published by ImagineDevOps .

**Move beyond basic programming to design and build reliable software with clean code**

## What is this book about?

This book distills the industry’s best practices for writing lean Go code that
is easy to test and maintain and explores their practical application on Links
‘R’ US: an example project that crawls web-pages and applies the PageRank
algorithm to assign an importance score to each one.

This book covers the following exciting features:

* Understand different stages of the software development life cycle and the role of a software engineer
* Create APIs using gRPC and leverage the middleware offered by the gRPC ecosystem
* Discover various approaches to managing package dependencies for your projects
* Build an end-to-end project from scratch and explore different strategies for scaling it
* Develop a graph processing system and extend it to run in a distributed manner
* Deploy Go services on Kubernetes and monitor their health using Prometheus

## Instructions
All of the code is organized into folders labelled after the chapter they
appear on. For example, Chapter02 contains the source code for the second book
chapter and so on.

The Makefile has been updated to manage dependencies via Go modules instead of
the dep tool. However, the dep tool will be used as a _fall-back_ for old Go
versions (that lack module support) or if the `GO111MODULE` environment
variable is set to `off` prior to running any of the Makefile targets.

Go 1.18+ is required for running the code/tests from the individual chapters.
The latest version of Go for your platform can be downloaded
[here](https://go.dev/dl/).

We also provide a PDF file that has color images of the screenshots/diagrams
used in this book. [Click here to download
it](https://static.packt-cdn.com/downloads/9781838554491_ColorImages.pdf).

### Intended audience
This Golang programming book is for developers and software engineers looking to use Go to design and build scalable distributed systems effectively. Knowledge of Go programming and basic networking principles is required.
