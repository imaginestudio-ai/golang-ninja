# Using Containers for Distribution

In this chapter, we’ll explore the world of containerization and examine the many reasons why you should use Docker containers for testing and distributing your applications. The term _containerization_ refers to a style of software packaging that makes it simple to deploy and run in any setting. First, we’ll go over the basics of Docker, covered by a simple application that can be built into an image and run as a container. Then, we return to our audiofile application, for a more advanced example, to learn how to create multiple Docker containers that can be composed and run together. These examples give you not only an understanding of the basic flags used for running containers but also some advanced flags that show you how to run containers with mapped network stacks, volumes, and ports.

We also explain how to use Docker containers for integration testing, which increases your confidence, because, let’s face it, mocking API responses can cover only so much. A good mix of unit and integration tests gives you not just the coverage but also confidence that the overall system works.

Finally, we will discuss some of the disadvantages of adopting Docker. Consider the increased complexity of administering containerized applications, as well as the additional overhead of operating several containers on a single host. Docker as an external dependency may be a disadvantage in itself. This chapter will help you determine when to use, and not to use, containers for your application.

By the end of this chapter, you will have a strong grasp of how to utilize Docker containers and how they might assist your development, testing, and deployment workflow. You will be able to containerize your application, test it with Docker, and release it with Docker Hub. Specifically, we’ll cover the following topics:

-   Why use containers?
-   Testing with containers
-   Distributing with containers

Bookmark

# Technical requirement

For this chapter, you will need to do the following:

-   Download and install Docker Desktop at [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)[](https://www.docker.com/products/docker-desktop/%0A)
-   Install the Docker Compose plugin

You can also find the code examples on GitHub at [https://github.com/ImagineDevOps DevOps/Building-Modern-CLI-Applications-in-Go/tree/main/Chapter13](https://github.com/ImagineDevOps DevOps/Building-Modern-CLI-Applications-in-Go/tree/main/Chapter13)

Bookmark

# Why use containers?

First, let’s talk about what a container is. A **container** is a standardized software unit that allows the transport of a program from one computing environment quickly and reliably to another by bundling the application’s code and all its dependencies into a single encapsulation. Simply put, containers let you package all your dependencies into a single container so that it can run on any machine. Containers are isolated from one another and bundle their own system libraries and settings, so they don’t conflict with other containers or the host system. This makes them a lightweight and portable alternative to **virtual machines** (**VMs**). Popular containerization tools include **Docker** and **Kubernetes**.

## Benefiting from containers

Let’s break down some of the benefits of using containers in your Go project:

-   **Portability**: Containers make it possible to support consistency of behavior across various environments, lowering the possibility of errors and incompatibilities.
-   **Isolation**: They offer a degree of isolation from the host system and other containers, which increases their level of security and reduces their propensity for conflicts.
-   **Lightweight**: Compared to VMs, containers are smaller and start up faster, which increases their operating efficiency.
-   **Scalability**: They can be easily scaled up or down, enabling effective resource use. For example, if you utilize containers for your application, then you can create multiple identical containers running your application deployed across multiple servers.
-   **Versioning**: Containers can be versioned, making it simple to revert to earlier iterations as needed.
-   **Modularity**: Because containers can be created and managed separately, they are simple to update and maintain.
-   **Cost-effective**: By lowering the number of systems you need to run your applications, containers can help you save money on infrastructure and maintenance.

Creating and running command-line applications is made simple and reliable by containers. Regardless of the host machine’s configuration, this means that the application will always be built and run in the same manner. Application development and deployment across different operating systems are made significantly simpler by including all necessary dependencies and runtime environments within the container image. Finally, containers make it simple to duplicate development environments, enabling multiple developers or teams to work together in the same area while guaranteeing that the application is developed and executed uniformly across various environments.

Additionally, using containers makes it simpler to integrate applications with **continuous integration and continuous deployment** (**CI/CD**) pipelines. Since all the necessary dependencies exist within the container’s image, the pipeline can more reliably and easily build and run the application, eliminating the need to configure the pipeline’s host machine’s development environment.

Finally, the consistency of an isolated environment with containers is another benefit that makes it easier to distribute your command-line application while guaranteeing that it will operate exactly as expected. Users no longer need to configure their environment for the application, making containers, while also lightweight, a great way to distribute across various environments and platforms.

As you can clearly see, there are a variety of situations where containers can prove useful, including command-line application development and testing! Now, let’s discuss when you may not want to use containers.

## Deciding not to use containers

While containers are often helpful, there are some circumstances in which they might not be the best option:

-   **High-performance computing**: High-performance computing and other tasks that need direct access to the host system’s resources might not be good candidates for containers because of the additional overhead they cause.
-   **Requiring high levels of security**: Containers share the host’s kernel and might not offer as much isolation as a VM. VMs may be a better option if your workload demands a high level of security.
-   **Neglecting container-native features**: You may not see the benefit of using containers if you do not plan on using any of the native features included for scaling, rolling updates, service discovery, and load balancing.
-   **Inflexible applications**: If an application requires a very specific operating environment in order to function properly, it might not even be possible to containerize it, as there are limited operating systems and platforms that are supported.
-   **Team inertia**: If you or your team are unwilling to learn about containers and container orchestration, then it will be difficult to incorporate a new tool.

Nevertheless, it’s important to note that these situations are not always the case and that there are some solutions available, including the use of VMs, particular security features of container orchestration platforms, specialized container runtimes such as **gVisor** or **Firecracker**, and others.

In the following examples and within the next sections, we will be using Docker to show how easy it can be to start using Docker and use it to create a consistent environment for testing and distribution.

In the `Chapter-13` GitHub repository, we go over a very simple example for building an image and running a container. The `main.go` file is simple:

```markup
func main() {
    var helloFlag bool
    flag.BoolVar(&helloFlag, "hello", false, "Print 'Hello,
      World!'")
    flag.Parse()
    if helloFlag {
        fmt.Println("Hello, World!")
    }
}
```

Passing in the `hello` flag to the built application will print out `"``Hello, World!"`.

## Building a simple Docker image

To start, software can be packaged as an **image**, a small, self-contained executable that contains the program’s source code, libraries, configuration files, runtime, and environment variables. Images are the building blocks of containers and are used to create and run them.

Let’s build a Docker image for this very simple application. To do so, we’ll need to create a **Dockerfile**. You can create a file named `Dockerfile` that will automatically be recognized when you run the command-line Docker commands, or create a file with the `.dockerfile` extension, which will require the `–f` or `--file` flag for passing in the filename.

A Dockerfile contains instructions for building a Docker image, as depicted in the following diagram. Each instruction creates a new layer within the image. The layers are combined to create the final image. There are many different kinds of instructions you can put in the Dockerfile. For example, you can tell Docker to copy files into the base image, set environment variables, run commands, and specify the executables to run when a container is initialized:

![Figure 13.1 – Visual of a Dockerfile transformed into an image with layers by the build command](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_13.1_B18883.jpg)

Figure 13.1 – Visual of a Dockerfile transformed into an image with layers by the build command

For our base image, let’s visit Docker Hub’s website at [https://hub.docker.com](https://hub.docker.com) and search for the official Golang Docker base image for Go `v1.19`. We see that we can use the `golang` image with tag `1.19`. The `FROM` instruction is the first line of the Dockerfile and it sets the base image to use:

```markup
FROM golang:1.19
```

Then, copy all the files:

```markup
COPY . .
```

Build the `hello` `world` application:

```markup
RUN go build main.go
```

Finally, run the application while passing the `hello` flag:

```markup
CMD ["./main", "--hello"]
```

Altogether, the Dockerfile contains the preceding instructions with some descriptive comments indicated by `#` as the first character in line.

To build a Docker image from a Dockerfile, we call the `docker build` command. The command takes the following syntax:

```markup
docker build [options] path| url | -
```

When run, the command does the following:

-   Reads the instructions specified within the Dockerfile and performs them in order
-   Each instruction creates a new layer in the image, and the final image combines them all
-   Tags the new image with the specified or generated name, and—optionally—a tag in the `name:tag` format

The `options` parameter can be used to pass in different options to the command, which can include build-time variables, targets, and more. The `path | url | -` argument specifies the location of the Dockerfile.

Let’s try building this image from the Dockerfile we created for our hello world application. Within the root of the repository, run the following command:

`docker build --``tag hello-world:latest`

After running the command, you should see similar output to this:

```markup
[+] Building 2.4s (8/8) FINISHED
=> [internal] load build definition from Dockerfile          0.0s
=> => transferring dockerfile: 238B                          0.0s
=> [internal] load .dockerignore                             0.0s
=> => transferring context: 2B                               0.0s
=> [internal] load metadata for docker.io/library/golang:1.19 1.2s
=> [internal] load build context                             0.0s
=> => transferring context: 2.25kB                           0.0s
=> CACHED [1/4] FROM docker.io/library/golang:1.19@sha256:bb9811fad43a7d6fd217324 0.0s
=> [2/4] COPY . .                                            0.0s
=> [3/4] RUN go build main.go                                1.0s
=> exporting to image                                        0.1s
=> => exporting layers                                       0.0s
=> => writing image sha256:91f97dc0109218173ccae884981f700c83848aaf524266de20f950   0.0s
=> => naming to docker.io/library/hello-world:latest         0.0s
```

From about midway through the output, you’ll see that the layers of the image are built, concluding with the final image tagged as `hello-world:latest`.

You can view the images that exist, by running the following command in your terminal:

```markup
% docker images
REPOSITORY         TAG        IMAGE ID        CREATED        SIZE
hello-world    latest    91f97dc01092    18 minutes ago  846MB
```

Now that we’ve successfully built our Docker image for this simple hello world application, let’s follow up by running it within a container.

## Running a simple Docker container

When you run a Docker container, Docker Engine takes an existing image and creates a new running instance of it. This container exists within an isolated environment that has its own filesystem, network interfaces, and process space. However, the image is a necessary starting point for creating—or running—the container.

Note

When a container is running, it can make changes to the filesystem, such as creating or modifying files. However, these changes are not saved in the image and will be lost when the container is stopped. If you want to save the changes, you can create a new image of the container using the `docker` `commit` command.

To create and run a Docker container from an image, we call the `docker run` command. The command takes the following syntax:

`docker run [options] image[:tag] [``command] [arg...]`

The `docker run` command checks if the image exists locally; if not, then it will pull it from Docker Hub. Docker Engine then creates a new container from this image, with all layers or instructions applied. We’ll break this down here:

![Figure 13.2 – Visual of an image used to create a container with the run command](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_13.2_B18883.jpg)

Figure 13.2 – Visual of an image used to create a container with the run command

As mentioned, when `docker run` is called, the following steps occur:

1.  Docker checks if the requested image exists locally; if not, it retrieves it from a registry, such as Docker Hub.
2.  From the image, it creates a new container.
3.  It starts the container and executes the commands specified within the instructions of the Dockerfile.
4.  It attaches the terminal to the container’s process in order to display any output from the commands.

The `options` parameter can be used to pass in different options to the command, which can include mapping ports, setting environment variables, and more. The `image[:tag]` argument specifies the image to use for creating the container. Finally, the `command` and `[arg...]` arguments are used to specify any commands to run within the container.

In each example that we call the `docker run` command, we pass in the `--rm` flag, which tells Docker to automatically remove the container when it exits. This will save you from accidentally ending up with many gigabytes of stopped containers sitting in the background.

Try to run an image from the `hello-world:latest` image that we created for our hello world application. Within the root of the repository, run the following command and see the text output:

```markup
% docker run --rm hello-world:latest
Hello, World!
```

We did it! A simple Dockerfile for a simple hello world application. Within the next two sections, we’ll return to the audiofile command-line application example and use this new skill of building images and running containers for testing and distribution.

Bookmark

# Testing with containers

So far, within our command-line application journey, we’ve built tests and mocked the service output. The benefit of using containers, besides having a consistent and isolated environment for tests to run on any host machine, is that you can use them for running integration tests that provide more reliable test coverage for your application.

## Creating the integration test file

We created a new `integration_test.go` file to handle the configuration and execution of integration tests, but we don’t want it to run with all the other tests. To specify its uniqueness, let’s tag it with `int`, short for integration. At the top of the file, we add the following build tag:

```markup
//go:build int && pro
```

We include the `pro` build tag because we are testing all the available features.

## Writing the integration tests

First, let’s write the `ConfigureTest()` function to prepare for our integration tests:

```markup
func ConfigureTest() {
    getClient = &http.Client{
        Timeout: 15 * time.Second,
    }
    viper.SetDefault("cli.hostname", "localhost")
    viper.SetDefault("cli.port", 8000)
    utils.InitCLILogger()
}
```

In the preceding code, you can see that we use an actual client, not the mocked client that is currently used within unit tests. We use `viper` to set the hostname and port for the API we connect as localhost on port `8000`. Finally, we initialize the logger files so that we don’t run into any panics while logging.

For the integration test, let’s use a specific workflow:

1.  **Upload audio**: First, we want to make sure an audio file exists within the local storage.
2.  **Get audio by id**: From the previous step, we can retrieve the audiofile ID returned and use this to retrieve the audio metadata from storage.
3.  **List all audio**: We list all the audio metadata and confirm that the previously uploaded audio exists within the list.
4.  **Search audio by value**: Search for that uploaded audio based on metadata we know exists within the description.
5.  **Delete audio by id**: Finally, delete the initial audio file we uploaded by the ID we retrieved from _step 1_.

The order is specific as the latter steps within the workflow depend on the first.

The integration tests are like the unit tests, but paths to real files are passed in, and the actual API is called. Within the `integration_tests.go` file exists a `TestWorkflow` function that calls the commands in the order listed previously. Since the code is similar to the unit tests, let’s just go over the first two command calls, and then move straight into using Docker to execute the integration tests!

Before any methods are tested, the integration test is configured by calling the `ConfigureTest` function:

```markup
ConfigureTest()
fmt.Println("*** Testing upload ***")
b := bytes.NewBufferString("")
rootCmd.SetOut(b)
rootCmd.SetArgs([]string{"upload", "--filename",
  "../audio/algorithms.mp3"})
err := rootCmd.Execute()
if err != nil {
    fmt.Println("err: ", err)
}
uploadResponse, err := ioutil.ReadAll(b)
if err != nil {
    t.Fatal(err)
}
id := string(uploadResponse)
if id == "" {
    t.Fatalf("expected id returned")
}
```

In the preceding code, we then use `rootCmd` to call the `upload` command with the filename set to `../audio/algorithms.mp3`. We execute the command and read the response back as a byte slice that is then converted to a string and stored in the `id` variable. This `id` variable is then used for the following tests. We run the `get` command and pass in the same `id` variable to retrieve the audiofile metadata for the previously uploaded audio:

```markup
fmt.Println("*** Testing get ***")
rootCmd.SetArgs([]string{"get", "--id", id, "--json"})
err = rootCmd.Execute()
if err != nil {
    fmt.Println("err: ", err)
}
getResponse, err := ioutil.ReadAll(b)
if err != nil {
    t.Fatal(err)
}
var audio models.Audio
json.Unmarshal(getResponse, &audio)
if audio.Id != id {
    t.Fatalf("expected matching audiofile returned")
}
```

We continue testing the `list`, `search`, and `delete` commands similarly and ensure that the specific metadata with a matching `id` variable is returned each time. When the tests are done, we try to run the integration test. Without the API running locally, running the following command fails miserably:

```markup
go test ./cmd -tags "int pro"
```

Before we try again, let’s build a Dockerfile to run the API within a contained environment.

## Writing the Dockerfiles

In the real world, our API might be hosted on some external website. However, we are currently running on `localhost`, and running it within a container will allow users to easily run it no matter which machine they use. In this section, we will create two Dockerfiles: one for the CLI and another for the API.

### Writing the API Dockerfile

First, we’ll create an `api.Dockerfile` file to hold all the instructions to build the image and run the container for the audiofile API:

```markup
FROM golang:1.19
# Set the working directory
WORKDIR /audiofile
# Copy the source code
COPY . .
# Download the dependencies
RUN go mod download
# Expose port 8000
EXPOSE 8000
# Build the audiofile application with the pro tag so all
# features are available
RUN go build -tags "pro" -o audiofile main.go
RUN chmod +x audiofile
# Start the audiofile API
CMD ["./audiofile", "api"]
```

Let’s build this image. The `–f` flag allows you to specify the `api.Dockerfile` file to use, and the `–t` flag allows you to name and tag the image:

```markup
% docker build -f api.Dockerfile -t audiofile:api .
```

After the command executes, you can run the `docker images` command to confirm its creation:

```markup
% docker images
REPOSITORY        TAG        IMAGE ID        CREATED        SIZE
audiofile      api         12afba7f3fb7        9 minutes ago  1.75GB
```

Now that we see that the image has been built successfully, let’s run the container and test it out!

Run the following command to run the container:

```markup
% docker run -p 8000:8000 --rm audiofile:api
Starting API at http://localhost:8000
Press Ctrl-C to stop.
```

You’ll see the preceding output if the API is started successfully. We have the audiofile API running in a container within your host. Remember that any commands will check against the flat file storage, pointing to the `audiofile` directory created under the `home` directory. Any audio files uploaded, processed, and with metadata stored within the container will not be saved unless we commit the changes. Since we are just running integration tests, this won’t be necessary.

Note

The `–p` flag within the `docker run` command allows you to specify the port mapping between the host and container. The syntax is `-p host_port:container_port`. This maps the host’s port to the container’s port.

Within a separate terminal, let’s run the integration tests again and see them pass:

```markup
% go test ./cmd -tags "int pro"
ok      github.com/marianina8/audiofile/cmd     0.909s
```

Success! We’ve now run integration tests connecting to the audiofile API within a container.

### Writing the CLI Dockerfile

Now, for running the CLI integration tests within a container, we’ll create a `cli.Dockerfile` file. It will hold all the instructions to build the image and run the container for the integration tests:

```markup
FROM golang:1.19
# Set the working directory
WORKDIR /audiofile
# Copy the source code
COPY . .
# Download the dependencies
RUN go mod download
# Execute `go test -v ./cmd -tags int pro` when the
# container is running
CMD ["go", "test", "-v", "./cmd", "-tags", "int pro"]
```

The preceding comments clarify each instruction, but let’s break down the Docker instructions:

1.  Specify and pull from the base image as `golang:1.19`.
2.  Set the working directory to `/audiofile`.
3.  Copy over all the source code to the working directory.
4.  Download all the Go dependencies.
5.  Execute `go test –v ./cmd -tags` `int pro`.

Let’s build this image:

```markup
% docker build -f cli.Dockerfile -t audiofile:cli .
```

Then, while ensuring the `audiofile:api` container is already running, run the `audiofile:cli` container:

```markup
% docker run --rm --network host audiofile:cli
```

You’ll see that the integration tests run successfully.

Note

The `--network host` flag within the `docker run` command is used to connect a container to the host’s network stack. It means that the container will have access to the host’s network interfaces, IP address, and ports. Be careful with security if the container runs any service.

Now, we’ve created two containers for the API and CLI, but rather than having to run each separately within two separate terminals, it’d be easier to use **Docker Compose**. Docker Compose is a plugin for Docker Engine that allows you to define and run multiple Docker applications all with a single file, `docker-compose.yml`, starting and stopping the entire application with a single `stop/start` command.

### Writing the Docker Compose file

Inside the `docker-compose.yml` Docker Compose file, we define both containers that need to be run, while specifying any parameters we’ve previously set via flags for the `docker` `run` command:

```markup
version: '3'
services:
  cli:
    build:
      context: .
      dockerfile: cli.Dockerfile
    image: audiofile:cli
    network_mode: host
    depends_on:
      - api
  api:
    build:
      context: .
      dockerfile: api.Dockerfile
    image: audiofile:api
    ports:
    - "8000:8000"
```

Let’s explain the preceding file. First, there are two services defined: `cli` and `api`. Beneath each service are a set of similar keys:

-   The `build` key, which is used to specify the context and location of the Dockerfile.
-   The `context` key is used to specify where to look for the Dockerfile. Both are set to `.`, which tells the Docker Compose service to look in the current directory.
-   The `dockerfile` key allows us to specify the name of the Dockerfile—in this case, `cli.Dockerfile` for the `cli` service and `api.Dockerfile` for the `api` service.
-   The `image` key allows us to give a name and tag the image.

For the `cli` service, we’ve added some further keys:

-   The `network_mode` key is used to specify the networking mode for a service. When it is set to `host`, like it is for the `cli` service, it means to use the host machine’s network stack (like the `–network host` flag used when calling `docker run` for the CLI).
-   The `depends_on` key allows us to specify the order of which services should be running first. In this case, the `api` service must be running first
-   For the `api` service, there’s an additional key:
-   The `ports` key is used to specify port mappings between the host machine and the container. Its syntax is `` `host_port:container_port` `` and is like the `–p` or `--publish` flag when calling the `docker` `run` command.

Now that we’ve got the Docker Compose file completed, we just have one simple command, `docker-compose up`, to run the integration tests within a containerized environment:

```markup
% docker-compose up
[+] Running 3/2
 Network audiofile_default  Created                        0.1s
 Container audiofile-api-1  Created                        0.0s
 Container audiofile-cli-1  Created                        0.0s
Attaching to audiofile-api-1, audiofile-cli-1
audiofile-api-1  | Starting API at http://localhost:8000
audiofile-api-1  | Press Ctrl-C to stop.
audiofile-cli-1  | === RUN   TestWorkflow
audiofile-cli-1  | --- PASS: TestWorkflow (1.14s)
…
audiofile-cli-1  | ok   github.com/marianina8/audiofile/cmd     1.163s
```

Now, no matter which platform you’re running the containers on, the results will be consistent while running the tests within a container. Integration testing provides more comprehensive testing as it will catch bugs that might exist within the **end-to-end** (**E2E**) flow from the command to the API to the filesystem and back. We can therefore increase our confidence with tests that can ensure our CLI and API are more stable and reliable as a whole. In the next section, we’ll discuss how to distribute your CLI application with containers.

Bookmark

# Distributing with containers

There are various advantages to running a CLI inside a container as opposed to directly on the host. Utilizing a container makes the setup and installation of the program easier. This can be helpful if the application needs numerous dependencies or libraries that are challenging to install. Additionally, regardless of the language or tools used to construct the program, adopting a container enables a more dependable and uniform method of distribution. Using a container as a distribution method can be a flexible solution for the majority of applications that can operate in a Linux environment, even though there may be language-specific alternatives. Finally, distributing through containers will be useful for developers unfamiliar with the Go language but who already have the Docker toolbox installed on their machines.

## Building a new image to run as an executable

To build an image that can run as an executable, we must create an `ENTRYPOINT` instruction on the image to specify the main executable. Let’s create a new Dockerfile, `dist.Dockerfile`, which contains the following instructions:

```markup
FROM golang:1.19
# Set the working directory
WORKDIR /audiofile
# Copy the source code
COPY . .
# Download the dependencies
RUN go mod download
# Expose port 8000
EXPOSE 8000
# Build the audiofile application with the pro tag so all
# features are available
RUN go build -tags "pro" -o audiofile main.go
# Start the audiofile API
ENTRYPOINT ["./audiofile"]
```

Since these instructions are mostly similar to the other Dockerfiles explained in the previous sections, we won’t go into a detailed explanation. The only instruction to note is the `ENTRYPOINT` instruction, which specifies `./audiofile` as the main executable.

We can build this image with the following command:

```markup
% docker build -f dist.Dockerfile -t audiofile:dist .
```

After confirming that the image is successfully built, we are now ready to run the container and interact with it as an executable.

## Interacting with your container as an executable

To interact with your container as an executable, you can configure your container to use an interactive TTY (terminal) with the `ENTRYPOINT` command in Docker. The `-i` and `–t` options stand for _interactive_ and _TTY_ respectively, and when the two flags work together, you can interact with the `ENTRYPOINT` command in a terminal-like environment. Remember to have the API running first. Now, let’s show how it’ll look when we run the container for the `audiofile:dist` image:

```markup
% docker run --rm --network host -ti audiofile:dist help
A command line interface allows you to interact with the Audiofile service.
Basic commands include: get, list, and upload.
Usage:
  audiofile [command]
Available Commands:
...
Use "audiofile [command] --help" for more information about a command.
```

Just by typing `help` at the end of the `docker run` command passes in `help` as input to the main executable, or `ENTRYPOINT`: `./audiofile`. As expected, the help text is output.

The `docker run` command uses a few additional commands; the `–network host` flag uses the host’s network stack for the container, and the `–rm` command tells Docker to automatically remove the container when it exits.

You can run any of the commands by just replacing the word `help` with the name of the other command. To run `upload`, for example, run this command:

```markup
% docker run --rm --network host -ti audiofile:dist upload –filename audio/algorithms.mp3
```

You can now interact with your command-line application through a container passing in commands and not have to worry if it will respond any differently based on the host machine. As previously mentioned, any filesystem changes, or files uploaded, as in preceding the file, are not saved when the container exists. There is a way to run the API so that local file storage maps to a container path.

## Mapping host machine to container file paths

As mentioned, you can map a host machine path to a Docker container file path so that files on the host computer may be accessed from inside the container. This can be helpful for things such as giving the container access to data volumes or application configuration files.

The `-v` or `—volume` option can be used to translate a host machine path to a container path when executing a container. This flag’s syntax is `host path:container path`. For instance, the `docker run -v /app/config:/etc/config imageName:tag` command would be used to map the host machine’s `/app/config` directory to the container’s `/``etc/config` directory.

It’s crucial to remember that both the host path and container path need to be present in the container image before the container can be executed. You must construct the container path before starting the container if it does not already exist in the container image.

If you dig into the audiofile API that is running on your local host, you’ll see that the flat file storage is mapped to the `/audiofile` folder existing under the host’s `home` directory. On my macOS instance, if I wanted to run the audiofile API within a Docker container but be able to read from and access or upload data to the flat file storage, then I would need to map the `audiofile` directory under my `HOME` directory to an appropriate location. This `docker run` command would do it:

```markup
docker run -p 8000:8000 --rm -v $HOME/audiofile:/root/audiofile  audiofile:api
```

Run the preceding command first and then run the CLI container, or modify the `docker-compose.yml` file’s API service to include the following:

```markup
    volumes:
      - "${HOME}/audiofile:/root/audiofile"
```

Either way, when you run the container for integration tests or as an executable, you’ll be interacting with your local storage mapped to the `/root/audiofile` directory within the container. If you’ve been playing around with the audiofile CLI and uploading directory, then when you start the container up and run the `list` command, you’ll see preexisting metadata instead of an empty list returned.

Mapping a path from your host to a container is an option that you can share with your users when instructing them how to use the audiofile application.

## Reducing image size by using multi-stage builds

By running the `docker images` command, you’ll see that some of the images built are quite large. To reduce the size of these images, you may need to rewrite your Dockerfile to use multi-stage builds. A **multi-stage build** is a process of dividing up the build into multiple stages, in which it is possible to remove unnecessary dependencies, artifacts, and configurations from the final image. This is especially useful when building images for large applications where you can save on deployment time as well as infrastructure costs.

A way that single-stage and multi-stage builds differ is that multi-stage builds allow you to use multiple `FROM` statements, each defining a new stage of the build process. You can selectively copy artifacts, or builds, from one stage or another, allowing you to take what you need and discard the rest, essentially allowing you to remove anything unnecessary and clean up space.

Let’s consider the `dist.Dockerfile` file and rewrite it. In our multi-stage build process, let’s define our stages:

-   **Stage 1**: Build our application
-   **Stage 2**: Copy the executable, expose the port, and create an entry point

First, we create a new file, `dist-multistage.Dockerfile`, with the following instructions:

```markup
# Stage 1
FROM golang:1.19 AS build
WORKDIR /audiofile
COPY . .
RUN go mod download
RUN go build -tags "pro" -o audiofile main.go
# Stage 2
FROM alpine:latest
COPY --from=build /audiofile/audiofile .
EXPOSE 8000
ENTRYPOINT ["./audiofile"]
```

In _Stage 1_, we copy all the code files, download all dependencies, then build the application—basically, all as in the original instructions within `dist.Dockerfile`, but without the `EXPOSE` and `ENTRYPOINT` instructions. One thing to note is that we’ve named the stage `build`, with the following line:

```markup
FROM golang:1.19 AS build
```

In _Stage 2_, we copy over just the compiled binary from the `build` stage and nothing else. To do this, we run the following instruction:

```markup
COPY --from=build /audiofile/audiofile .
```

The command allows us to copy a file or directory from a previous stage, `build`, to the current stage. The `--from=build` option specifies the stage name to copy the file from. `/audiofile/audiofile` is the path of the file in the `build` stage, and `.` at the end of the command specifies the destination directory, the root directory, of the current stage.

Let’s try building it and comparing the new size against the original size:

```markup
REPOSITORY TAG            IMAGE ID        CREATED        SIZE
audiofile    dist            1361cbc7be3e    2 minutes ago    1.78GB
audiofile     dist-multistage    ab5640f99ef2    5 minutes ago    24MB
```

That’s a big difference! Using multi-stage builds will help you to save on deployment time and infrastructure costs, so it’s definitely worth the time writing your Dockerfiles using this process.

## Distributing your Docker image

There are many methods for making your Docker images accessible to others. **Docker Hub**, a public registry where you can post your images and make them readily available to others, is a popular alternative. Another alternative is to use **GitHub Packages** to store and distribute your Docker images alongside other sorts of packages. There are other cloud-based registries such as **Amazon Elastic Container Registry** (**ECR**), **Google Container Registry** (**GCR**), and **Azure Container Registry** (**ACR**) that provide extra services such as image scanning (for OS vulnerabilities, for example) and signing.

It’s a good idea to give instructions on how to utilize your image and run a container in the README file of the repository where the image is located. People who are interested in utilizing your image will be able to readily access instructions on how to retrieve the image, run a container using the image, and any other pertinent facts.

There are several advantages to publishing a Docker image, including simple distribution, versioning, deployment, collaboration, and scalability. Your image may be rapidly and readily distributed to others, making it simple for others to utilize and operate your application. Versioning helps you to maintain track of several versions of your image so that you may revert to an earlier version if necessary. Easy deployment allows you to deploy your application to several environments with little modifications. Sharing images via a registry facilitates collaboration with other developers on a project. And scalability is simple to accomplish by using the same image to create as many containers as you need, making it simple to grow your application.

In this chapter, we’ll publish to Docker Hub as an example for our audiofile CLI project.

### Publishing your Docker image

To publish an image to Docker Hub, you’ll first need to create an account on the website. Once you have an account, you can sign in and create a new repository to store your image. After that, you can use the Docker command-line tool to log in to your Docker Hub account, tag your image with the repository name, and push the image to the repository. Here is an example of the commands you would use to do this:

```markup
docker login --username=your_username
docker tag your_image your_username/your_repository:your_tag
docker push your_username/your_repository:your_tag
```

1.  Let’s try this with our audiofile API and CLI images. First, I will log in with my username and password:
    
    ```markup
    % docker login --username=marianmontagnino
    
    Password:
    
    Login Succeeded
    
    Logging in with your password grants your terminal complete access to your account.
    
    For better security, log in with a limited-privilege personal access token. Learn more at https://docs.docker.com/go/access-tokens/
    ```
    
2.  Next, I tag my CLI image:
    
    ```markup
     % docker tag audiofile:dist  marianmontagnino/audiofile:latest
    ```
    
3.  Finally, I publish the image to Docker Hub:
    
    ```markup
    % docker push marianmontagnino/audiofile:latest
    
    The push refers to repository [docker.io/marianmontagnino/audiofile]
    
    c0f557e70e4f: Pushed
    
    98f8be277d74: Pushed
    
    6c199763ccbe: Pushed
    
    8f2f7ffa843f: Pushed
    
    10bb928a2e24: Pushed
    
    f1ce3f3654c3: Mounted from library/golang
    
    3685241d2bbb: Mounted from library/golang
    
    dddbac67c6fa: Mounted from library/golang
    
    85f9ebffaf4d: Mounted from library/golang
    
    72235aad06ad: Mounted from library/golang
    
    5d37ad02a8e2: Mounted from library/golang
    
    ea8ab45f064e: Mounted from library/golang
    
    latest: digest: sha256:b7b3f58da01d360fc1a3f2e2bd617a44d3f7be d6b6625464c9d787b8a71ead2e size: 2851
    ```
    

Let’s confirm in Docker Hub to make sure that the container exists:

![Figure 13.3 – Screenshot of the Docker Hub website showing the audiofile image tagged with latest](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_13.3_B18883.jpg)

Figure 13.3 – Screenshot of the Docker Hub website showing the audiofile image tagged with latest

It’s a great idea to include instructions for running a container using the image within a README file of the repository where the image is stored. This makes it easy for people who want to use the image to learn how to pull the image and run the container properly. As an example, here are sample instructions for our previously uploaded audiofile CLI image:

To run the audiofile CLI container, ensure that the audiofile API container is running first. Next, run the `docker` command:

```markup
% docker run --rm --network host -ti marianmontagnino/audiofile:latest help
```

You’ll see that the help text is output. Let’s update the instructions on the Docker Hub repository.

### Updating the README

From the Docker Hub repository, where our image is stored (in this example, the audiofile repository), we can scroll down to the bottom of the page to see a **README** section:

![Figure 13.4 – Screenshot of the README section in the repository on Docker Hub](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_13.4_B18883.jpg)

Figure 13.4 – Screenshot of the README section in the repository on Docker Hub

Click **here** to edit the repository description. Add the instructions we discussed previously, then click the **Update** button:

![Figure 13.5 – Screenshot of the updated README section](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_13.5_B18883.jpg)

Figure 13.5 – Screenshot of the updated README section

Follow these instructions to similarly publish the audiofile API image to your Docker Hub repository. Now that the images exist in a public repository on Docker Hub, they are available to share and distribute to other users.

## Depending on Docker

The fact that users must have Docker installed on their computers is one of the key disadvantages of utilizing Docker to deploy a CLI. However, if your program has complicated dependencies or is designed to operate on various platforms, this Docker dependency may be easier to handle. Using Docker may assist in avoiding difficulties with many libraries and unexpected interactions with various system setups.

Bookmark

# Summary

We’ve gone into the realm of containerization and examined the numerous advantages of utilizing Docker containers for your applications in this chapter. The fundamentals of creating and running a simple Docker image and container are explained, as well as some more sophisticated instances using our audiofile application, which requires the construction of multiple containers that can be composed and run together.

Clearly, utilizing Docker for integration testing boosts your trust in the whole system, and we discussed how to run integration tests using Docker Compose.

At the same time, we’ve acknowledged some of Docker’s drawbacks, such as the increased complexity of maintaining containerized applications, the additional burden of operating several containers on a single host, and the external dependency of Docker itself.

Overall, this chapter has given you a strong knowledge of when to utilize Docker containers for command-line applications—for testing and distribution. Now, you can ensure that the application runs consistently across any host machine. It is up to you, though, to decide if the upsides outweigh the downsides of having an external dependency and some level of complexity.

In the next chapter, [_Chapter 14_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_14.xhtml#_idTextAnchor359), _Publishing your Go Binary as a Homebrew Formula with GoReleaser_, we’ll take distribution to a next level. We will get your application available on the official Homebrew repository to further increase the distribution of your application.

Bookmark

# Questions

1.  Which command is used to create and run a container from an image?
2.  Which `docker run` flag is used to attach a host machine path to a container path?
3.  Which Docker command is used to see all created containers?

Bookmark

# Further reading

-   _Docker: Up and Running: Shipping Reliable Containers in Production_ by _Sean P. Kane_ and _Karl Matthias_
-   _Continuous Delivery with Docker and Jenkins: Delivering software at scale_ by _Rafal Leszko_
-   _Docker in Action_ by _Jeff Nickoloff_ and _Stephen Kuenzli_

Bookmark

# Answers

1.  The `docker` `run` command.
2.  The `-v`, or `--volume`, flag is used to attach a host machine path to a container path during execution.
3.  `docker ps` or `docker` `container ls`.