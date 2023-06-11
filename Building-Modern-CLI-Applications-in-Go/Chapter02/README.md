# Structuring Go Code for CLI Applications

Programming is like any other creative process. It all begins with a blank slate. Unfortunately, when faced with a blank slate and minimal experience with programming applications from scratch, doubt can kick in – without knowing how to start, you may feel that it’s not possible at all.

This chapter is a guide on the first steps of creating a new application, beginning with some of the most popular ways to structure code, describing each, and weighing up their pros and cons. The concept of domain-driven design is discussed, as this can also influence the resulting structure of an application.

An example of an _audio metadata CLI application_ gives us an idea of what some real-world use cases or requirements could look like. Learning how to define an application’s use cases and requirements is a tedious but necessary step to ensuring a successful project that also meets the needs of all parties involved.

By the end of this chapter, you will have learned all the skills needed to build your application based on your specific use cases and requirements.

This chapter will cover the following topics:

-   Commonly used program layouts for robust applications
-   Determining use cases and requirements
-   Structuring an audio metadata CLI application



# Commonly used program layouts for robust applications

Along your programming journey, you may come across many different structures for applications. There is no standard programming layout for Go. Given all this freedom, however, the choice of the structure must be carefully made because it will dictate whether we understand and know how to maintain our application. The proper structure for the application will ideally also be simple, easy to test, and directly reflect the business design and how the code works.

When choosing a structure for your Go application, use your best judgment. Do not choose arbitrarily. Listen to the advice in context and learn to justify your choices. There’s no reason to choose a structure early, as your code will evolve over time and some structures work better for small applications while others are better for medium to large applications.

## Program layouts

Let’s dig into some common and emerging structural patterns that have been developed for the Go language so far. Understanding each option will help you choose the best design structure for your next application.

### Flat structure

This is the simplest structure to start with and is the most common when you are starting with an application, only have a small number of files, and are still learning about the requirements. It’s much easier to evolve a flat structure into a modular structure, so it’s best to keep it simple at the start and partition it out later as the project grows.

Let us now see some advantages and disadvantages of this structure:

-   **Pros**:
    -   It’s great for small applications and libraries
    -   There are no circular dependencies
    -   It’s easy to refactor into a modular structure
-   **Cons**:
    -   This can be complex and disorganized as the project grows
    -   Everything can be accessed and modified by everything else
-   **Example**
-   As the name implies, all the files reside in the root directory in a flat structure. There is no hierarchy or organization and this works well when there is a small number of files:

![Figure 2.1 – Example of a flat code structure](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_2.1_B18883.jpg)

Figure 2.1 – Example of a flat code structure

As your project grows, there are several different ways to group your code to keep it organized, each with its advantages and disadvantages.

### Grouping code by function

Code is separated by its similar functionality. In a _Go REST API_ project, as an example, Go files are commonly grouped by handlers and models.

Let us now see some advantages and disadvantages of this structure:

-   **Pros**:
    -   It’s easy to refactor your code into other modular structures
    -   It’s easy to organize
    -   It discourages a global state
-   **Cons**:
    -   Shared variables or functionality may not have a clear place to live
    -   It can be unclear where initialization occurs

To mitigate any confusion that can occur, it’s best to follow Go best practices. If you choose the **group-by-function** structure, use the `main.go` file to initialize the application from the project root. This structure, as implied by the name, separates code based on its function. The following figure is an example of groups by function and the types of code that would fall into these different categories:

![Figure 2.2 – Example of grouping by functionality](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_2.2_B18883.jpg)

Figure 2.2 – Example of grouping by functionality

-   **Example**: The following is an example of folder organization that follows the group-by-function structure. Similar to the example grouping, folders associated with handlers contain code for each type of handler, folders associated with extractors contain code for each particular extraction type, and storage is also organized by type:

![Figure 2.3 – Example of a group-by-function structure](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_2.3_B18883.jpg)

Figure 2.3 – Example of a group-by-function structure

### Grouping by module

Unfortunately, the title of this style of architecture is a bit redundant. To clarify, grouping by module means creating individual packages that each serve a function and contain everything necessary to accomplish these functions within them:

-   **Pros**:
    -   It’s easier to maintain
    -   There is faster development
    -   There is low coupling and high cohesion
-   **Cons**:
    -   It’s complex and harder to understand
    -   It must have strict rules to remain well organized
    -   It may cause stuttering in package method names
    -   It can be unclear how to organize aggregated functionality
    -   Circular dependencies may occur

The following is a visual representation of how packages can be grouped by module. In the following example, the code is grouped depending on the implementation of the extractor interface:

![Figure 2.4 – Visual representation of grouping by module](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_2.4_B18883.jpg)

Figure 2.4 – Visual representation of grouping by module

-   **Example**
-   The following is an example of an organizational structure in which code is grouped into specific module folders. In the following example, the code to extract, store, and define the type, tags, transcript, and other metadata is stored within a single defined folder:

![Figure 2.5 – Example of a group-by-module structure](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_2.5_B18883.jpg)

Figure 2.5 – Example of a group-by-module structure

### Grouping by context

This type of structure is typically driven by the domain or the specific subject for which the project is being developed. The common domain language used in communication between developers and domain experts is often referred to as ubiquitous language. It helps developers to understand the business and helps domain experts to understand the technical impact of changes.

**Hexagonal architecture**, also called **ports** and **adapters**, is a popular domain-driven design architecture that conceptually divides the functional areas of an application across multiple layers. The boundaries between these layers are interfaces, also called ports, which define how they communicate with each other, and the adapters exist between the layers. In this layered architecture, the outer layers can only talk to the inner layers, not the other way around:

-   **Pros**:
    -   There is increased communication between members of the business team and developers
    -   It’s flexible as business requirements change
    -   It’s easy to maintain
-   **Cons**:
    -   It requires domain expertise and for developers to understand the business first before implementation
    -   It’s costly since it requires longer initial development times
    -   It’s not suited to short-term projects

The following provides a typical visual representation of a hexagonal structure. The arrows point inward toward entities to distinguish that the outer layers have access to the inner layers, but not the other way around:

![Figure 2.6 – Visual representation of hexagonal architecture](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_2.6_B18883.jpg)

Figure 2.6 – Visual representation of hexagonal architecture

-   **Example**: The following is a folder structure organized by context. Services with individual business functions are separated into their respective folders:

![Figure 2.7 – Example of a group-by-context structure](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_2.7_B18883.jpg)

Figure 2.7 – Example of a group-by-context structure

That wraps up the different types of organizational structures for a Go application. There’s not necessarily a right or wrong folder structure to use for your application; however, the business structure, the size of the project, and your general preference can play a part in the final decision. This is an important decision to make, so think carefully before moving forward!

## Common folders

No matter the chosen structure, there are commonly named folders across existing Go projects. Following this pattern will help to increase understanding for maintainers and future developers of the application:

-   **cmd**: The `cmd` folder is the main entry point for the application. The directory name matches the name of the application.
-   **pkg**: The `pkg` folder contains code that may be used by external applications. Although there is debate on the usefulness of this folder, `pkg` is explicit, and being explicit makes understanding crystal clear. I am a proponent of keeping this folder solely due to its clarity.
-   **internal**: The `internal` folder contains private code and libraries that cannot be accessed by external applications.
-   **vendor**: The `vendor` folder contains the application dependencies. It is created by the `go mod vendor` command. It’s usually not committed to a code repository unless you’re creating a library; however, some people feel safer having a backup.
-   **api**: The `api` folder typically contains the code for an application’s _REST API_. It is also a place for Swagger specification, schema, and protocol definition files.
-   **web**: The `web` folder contains specific web assets and application components.
-   **configs**: The `configs` folder contains configuration files, including any `confd` or `consul-template` files.
-   **init**: The `init` folder contains any system initiation (start) and process management (stop/start) scripts with supervisor configurations.
-   **scripts**: The `scripts` folder contains scripts to perform various builds, installations, analyses, and operations. Separating these scripts will help to keep the **makefile** small and tidy.
-   **build**: The `build` folder contains files for packaging and continuous integration. Any cloud, container, or package configurations and scripts for packaging are usually stored under the `/build/package` folder, and continuous integration files are stored under `build/ci`.
-   **deployments** (or **deploy**): The `deployments` folder stores configuration and template files related to system and container orchestration.
-   **test**: There are different ways of storing test files. One method is to keep them all together under a `test` folder or to keep the test files right alongside the code files. This is a matter of preference.

Note

No matter what folders are contained within your project structure, use folder names that clearly indicate what is contained. This will help current and future maintainers and developers of the project find what they are looking for. Sometimes, it will be difficult to determine the best name for a package. Avoid overly used terms such as _util_, _common_, or _script_. Format the package name as all lowercase, do not use `snake_case` or `camelCase`, and consider the functional responsibility of the package and find a name that reflects it.

All of the aforementioned common folders and structural patterns described apply to building a CLI application. Depending on whether the CLI is a new feature of an existing application or not, you may be inheriting an existing structure. If there is an existing `cmd` folder, then it’s best to define an entry to your CLI there under a folder name that identifies your CLI application. If it is a new CLI application, start with a flat structure and grow into a modular structure from there. From the aforementioned examples, you can see how a flat structure naturally grows as the application extends to offer more features over time.

Just Imagine

# Determining use cases and requirements

Before building your CLI application, you’ll need to have an idea of the application’s purpose and responsibilities. The purpose of the application can be defined as an overarching description, but to start implementing, it’s necessary to break down the purpose into use cases and requirements. The goal of use cases and requirements is to drive effective discussion around what an application should do, with the result that everyone has a shared understanding of what is going to be built and continues these discussions as the application evolves.

## Use cases

**Use cases** are a way of documenting the functional requirements of a project. This step, at least for CLIs, is typically handled by an engineer after gathering some high-level requirements from their internal, or external, business customers. It’s important to have a clear picture of the application’s purpose and document the use cases before any technical implementation because the use cases themselves will not include any implementation-specific language or details about the interface. During any discussion of requirements with customers, topics related to implementation may crop up. Ideally, it’s best to steer the conversation back toward use case requirements and handle one thing at a time, staying focused on the right kind of discussion with the right people. The resulting use cases will reflect the goals of the application.

## Requirements

**Requirements** document all the nonfunctional restraints on how the system should perform the use cases. Although a system can still work without meeting these nonfunctional requirements, it may not meet user or consumer expectations as a result. There are some common categories for requirements, each discussed in detail next:

-   **Security**

Security requirements ensure that sensitive information is securely transmitted and that the application adheres to secure coding standards and best practices.

A few example security requirements include the following:

-   Exclude sensitive data related to sessions and systems from logs
-   Delete unused accounts
-   There are no default passwords in use for the application

-   **Capacity**

Capacity requirements deal with the size of data that must be handled by an application to achieve production goals. Determine the storage requirements of today and how your application will need to scale with increased volume demands. A few examples of capacity requirements include the following:

-   Storage space requirements for logging
-   Several concurrent users can use the application at any given time
-   The limit on the amount of data that can be passed into the application

-   **Compatibility**

Compatibility requirements determine the minimum hardware and operating system requirements for the application to run as expected. Examples include stating the following in your installation requirements:

-   The required architecture
-   All compatible and non-compatible hardware
-   CPU and memory requirements

-   **Reliability** **and availability**

Reliability and availability requirements define what happens during full or partial failure and set the standard for your application’s accessibility. A few examples would include the following:

-   Minimum allowed failures per transaction or time frame
-   Defining accessibility hours for your application

-   **Maintainability** **and manageability**

Maintainability requirements determine how easily the application can be fixed when a bug is discovered or enhanced when there are new feature requirements. Manageability requirements determine how easily an administrator can manage an application. Examples of maintainability requirements include the following:

-   Bugs must be detected quickly and fixed within an appropriate period
-   The application should maintain compatibility with the latest hardware and operating system versions

-   **Scalability**

Scalability requirements determine the highest workload under which your application can still perform as expected. It is mainly driven by two factors: early software decisions and the infrastructure. Scaling can be horizontal or vertical, where horizontal scaling involves adding more nodes to the system and vertical scaling means adding more memory or faster CPUs to a machine. A couple of examples include the following:

-   Several concurrently connected users can use the application with the expected results
-   The number of transactions per millisecond is limited

-   **Usability**

Usability requirements determine the quality of the user experience. A few simple examples include the following:

-   The application helps guide users toward the correct usage when they do the wrong thing
-   Help and documentation inform users about new arguments and flags to use
-   During a long operation, users are kept up to date on its progress

-   **Performance**

Performance requirements determine the responsiveness of an application. This includes the following:

-   The minimum required time for users to wait for specific operations to complete
-   Responsiveness to users’ actions

-   **Environment**

The environment requirements determine which environments the system will be expected to perform within. A few examples include the following:

-   The required environment variables that must be set
-   Dependencies on third-party software that need to be installed alongside applications

By taking the time to define the use cases and requirements, everyone involved will get a clear picture and have a shared understanding of the purpose and functionality of the application. A shared understanding will lead to a product that benefits in several ways, which we will discuss now.

## Disadvantages and benefits of use cases and requirements

Having functional and nonfunctional requirements mapped through use cases and requirements can greatly benefit the outcome of an application.

Here are some disadvantages of determining use cases and requirements:

-   It slows down the development process because requirements require time to be properly defined
-   Use cases and requirements may change over time

Next, we have some advantages of determining the use cases and requirements:

-   It provides the best possible outcome
-   Engaging in problem-solving discussions with your team determines potential issues, misuse, or misunderstanding
-   It defines the application’s goals, future targets, and estimated costs
-   You can prioritize each of the requirements

The goal is to gain a level of clarity that helps developers focus on solving the problem with the least amount of ambiguity. Beneficial discussions and collaborative time spent with the team determining the goals of the application are necessary aspects of the process that can be achieved in parallel with defining the use cases and requirements.

## Use cases, diagrams, and requirements for a CLI

Let’s discuss a theoretical scenario to illustrate how to build use cases and diagrams for a CLI. Suppose there is a large audio company with one particular team that focuses entirely on metadata extraction. This team provides audio metadata to their customers and other internal teams within the same audio company. Currently, they have an API available to anyone within the company’s internal network, but an operations team requests a CLI tool. The operations team recognizes the benefit of rapidly building scripts around a CLI application, which could open new opportunities for innovation for the team.

The existing customer-facing API use cases should be similar to the CLI since the implementation and the user interface are not a part of the documentation. Consider the use cases for the metadata team’s internal-facing CLI here. For record-keeping, we’ll number them and take the first several use cases as examples:

1.  Uploading audio
2.  Requesting metadata
3.  Extract metadata
4.  Processing speech to text
5.  Requesting speech-to-text transcripts
6.  Listing audio metadata in storage
7.  Searching audio metadata in storage
8.  Deleting audio from storage

For record-keeping, we’ll number them and take the first three use cases as examples.

### Use case 1 – uploading audio

An authenticated member of the operations team can upload audio by providing a file path. The upload process will automatically save uploads to storage and trigger audio processing to extract the metadata, and the application will respond with a unique ID to use when requesting the metadata.

This use case can be broken down into some common components:

-   **The actors** are the end users. This can be defined as a human or another machine process. The primary actor in this example use case is a member of the operations team, but since the team wants to use this CLI for scripting, another machine process is also an actor.
-   **Preconditions** are statements that must take place for the use case to occur. In this example, the member must be authenticated before any of the use cases can run successfully. The preconditions in _Figure 2__.8_ are represented by the solid line with an arrow pointing toward **Verify TLS Certificate**, which confirms through the **Certificate Management Client** that the user is authenticated.
-   **Triggers** are events that start another use case. These triggers can be either internal or external. In this example use case, the trigger is external – when a user runs the `upload` command. This use case triggers another use case, _Use case 3, Extract Metadata_, internally to extract metadata from the audio file and save it to storage. This is represented by the **Metadata Extractor** process box.
-   When everything happens as intended without exceptions or errors, the **basic flow** is activated. In _Figure 2__.8_, the basic flow is a solid line. The user uploads the audio and eventually returns an ID in response. Success!
-   The **alternative flow** shows variations of the basic flow, in which errors or exceptions happen. In _Figure 2__.8_, the alternative flow is a dotted line. The user uploads the audio, but an error occurs – for example, _the user is invalid_ or _the audio file does_ _not exist_.

Note

The use case diagram for uploading audio is illustrated with the basic flow in a solid line and the alternative flow in a dotted line.

![Figure 2.8: Use case diagram for uploading audio with a metadata CLI](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_2.8_B18883.jpg)

Figure 2.8: Use case diagram for uploading audio with a metadata CLI

Alongside the diagram (_Figure 2__.8_), we can write out the use case entirely as follows:

-   **Use case 1 –** **uploading audio:**
    -   **Name**: Uploading audio.
    -   **Description**: The **Actor** uploads audio by providing a file path. The application returns a unique ID to use for requesting the audio metadata. The actor represents a member of the operations team or another machine process.
    -   **Precondition**: The actor must be authenticated. In _Figure 2.8_, this is represented by the solid line from the **Upload Audio** box to the **Verify TLS** **Certificate** diamond.
    -   **Trigger**: The **Actor** triggers the upload command while passing in a valid file path as a flag – in _Figure 2.8_, the arrow pointing from **Actor** to **Upload Audio**.
-   **Basic flow:**

The actor runs the `upload` command in the CLI – in _Figure 2__.8_, the arrow pointing from **Actor** to **Upload Audio****:**

1.  Once **Preconditions** have been validated, the audio is vaildated. In _Figure 2__.8_, this is represented by the **Validate** **Audio** box.
2.  In _Figure 2__.8_, the validated audio moves to the **Process Metadata** step, which involves extracting the metadata, this is represented by the arrow pointed to the **Metadata Extractor** process box.
3.  The validated audio moves to the next step of **Upload Audio**, which saves the audio to the **database** (**DB**), represented by the **Upload to DB** box in _Figure 2__.8_.
4.  In _Figure 2__.8_, the **Return ID** box represents the **ID** being returned from the database, which is later passed down to the **Actor**.

-   **Alternative flow:**
    
    1.  **Error for an unauthenticated user**: An error is returned to the actor when TLS certification fails.
        -   **End use case**: In _Figure 2.8_, if the user is invalid, the error is returned, as represented by the dotted line from the **Invalid User** to **Error** box then arrow back to the **Actor**.
    2.  **Error for invalid audio**: An error is returned to the actor when audio fails to pass the validation process.
        -   **End use case**: In _Figure 2.8_, if the audio is invalid, an error is returned to the actor, represented by the **Failed Validation** to **Error** box then arrow back to the **Actor**.
    3.  **Error uploading the validated audio to storage**: An error is returned to the actor when audio upload to the database fails.
    
    -   **End use case**: In _Figure 2.8_, the dotted line returned from **Upload to DB** to the **Failed Upload** to **Error** box then arrow back to the **Actor**.

### Use case 2 – requesting metadata

An authenticated member of the operations team can retrieve audio metadata by providing an **ID** that was either returned after the upload or found by listing or searching for audio. The `get` command will output the requested audio metadata, with matching ID, in the specified format – either plain text or JSON.

Note

The use case diagram for requesting audio is illustrated with the basic flow in a solid line and the alternative flow in a dotted line.

![Figure 2.9: Use case diagram for the use case of requesting metadata](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_2.9_B18883.jpg)

Figure 2.9: Use case diagram for the use case of requesting metadata

With the preceding diagram (_Figure 2__.9_) at hand, let’s get into the use case as follows:

**Use case 2 –** **requesting metadata:**

-   **Name**: Requesting metadata.
-   **Description**: The **Actor** requests audio metadata by calling the `get` command and providing an **ID** for the audio. The application will output the requested audio metadata in plaintext or JSON format.
-   **Precondition**: The **Actor** must be authenticated. In _Figure 2.9_, this is represented by the solid line from the **Request Audio** box to the **Verify TLS** **Certificate** diamond.
-   **Trigger**: The **Actor** calls the `get` command while passing in the **ID** as an argument- in _Figure 2.9_, the arrow pointing from **Actor** to the **Request** **Audio** box.

Note that different formatting levels were used for preceding Use case 1 - make consistent throughout chapter for all use cases?

**Basic flow**

1.  The actor runs the `get` command in the CLI. In _Figure 2__.9_, the basic flow is represented in the solid line and starts with the arrow pointing from **Actor** to the **Request** **Audio** box.
2.  Once **Preconditions** have been validated, the audio metadata is retrieved by its **ID** from the database. In _Figure 2__.9_, this is represented by the solid line connecting **Request Metadata By ID** to **Database**.
3.  The **Database** returns the metadata successfully. In _Figure 2__.9_, this is represented by the line connecting **Request Metadata By ID** to the **Passed** box.
4.  Finally, the formatted metadata is returned to the **Actor**. In _Figure 2__.9_, this is represented by the solid line connecting **Passed** to **Actor**.

**Alternative flow**

1.  **Error for unauthenticated user**: An error is returned to the actor when TLS certification fails.
    -   **End use case**: In _Figure 2.9_, if the user is invalid, the error is returned, as represented by the dotted line from the **Invalid User** box to the **Error** box then the arrow back to the **Actor**.
2.  **Error for not found**: An error is returned if there is no matching metadata for the **ID**.

-   **End use case**: In _Figure 2.9_, the flow is represented by the dotted line from the **Failed** box to the **Error** box and then the arrow back to **Actor**.

### Use case 3 – extract metadata

Triggered by **Upload Audio**, metadata, including tags and transcript data, is extracted from the audio file and placed in storage.

Note

The use case diagram for requesting audio is illustrated with the basic flow in a solid line and the alternative flow in a dotted line.

![Figure 2.10: Use case diagram for processing metadata](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_2.10_B18883.jpg)

Figure 2.10: Use case diagram for processing metadata

With the preceding diagram (_Figure 2__.10_) in mind, let’s get into the matching use case:

-   **Use case 3 –** **extract metadata:**
    -   **Name**: Extract metadata
    -   **Description**: The metadata extraction process consists of extracting specific metadata, including album, artist, year, and speech-to-text transcription, and storing it in a metadata object in storage
    -   **Precondition**: Validated audio
    -   **Trigger**: Uploading audio
-   **Basic flow:**
    1.  Once the **Preconditions** are met, the **Metadata Extractor** process extracts tags and stores the data on the metadata object. In _Figure 2__.10_, this is represented by the solid line from the successfully validated audio, with the **Passed** box to the **Metadata Extractor** process box, then from the **Extract Tags** box to the **Passed** box.
    2.  Next, the transcript is extracted. In _Figure 2__.10_, this is represented by the solid line from the **Passed** box to the **Extract Transcript** box to the next **Passed** box.
    3.  The metadata extraction completes and updates the metadata object. In _Figure 2__.10_, this step is represented by the solid line from the **Passed** box to **Completed** and the solid line that runs to the **Metadata Object**.
    4.  The metadata object is stored. In _Figure 2__.10_, this is represented by the solid line from **Metadata Object** to **Database**.
-   **Alternative flow:**

Note that previous two use cases above used different formatting for "End use case" lines - check and make consistent throughout chapter for all use cases

1.  **Error extracting tag data**: In _Figure 2__.10_, this is represented by the dotted line from **Extract Tags** to **Error**.
    -   **Stores an error on the metadata object**: In _Figure 2.10_, this is represented by the dotted line from **Error** to **Metadata Object**.
    -   **End use case**: In _Figure 2.10_, this is represented by the solid line from **Metadata Object** to the **Database**.
2.  **Error extracting the transcript**: An error when extracting transcript metadata occurs. In _Figure 2__.10_, this is represented by the dotted line from **Extract Transcript** to **Error**.

-   **Stores an error on the metadata object**: In _Figure 2.10_, this is represented by the dotted line from **Error** to **Metadata Object**.
-   **End use case**: In _Figure 2.10_, this is the solid line from **Metadata Object** to **Database**.

It’s not necessary to write out the full documentation for each use case in order to understand the concept. Typically, the functional requirements as described by their use cases are reviewed by the stakeholders, and they are discussed to ensure there is agreement across the board.

## Requirements for a metadata CLI

Given our theoretical scenario of an internal team handling all audio metadata, a few nonfunctional requirements may also be requested and defined between the internal team and their customers. The requirements, for example, may include the following:

-   The application must run on Linux, macOS, and Windows
-   The ID returned from when the audio is uploaded must be returned immediately
-   The application must clearly state if the user misuses the application and uploads a file type other than audio

There are many more possible requirements for this metadata CLI application, but it is most important to understand what a requirement is, how to form your own, and how it differs from a use case. Use cases and requirements can be broken down into phases for more granularity, especially for scalability. Applications will grow over time and certain features will be added to match the growing requirements. To reference an earlier CLI guideline, _prototype first and optimize later_, it’s best just to get the application working first before optimizing. Depending on the type of issues encountered, whether slow processing, the inability to support a large number of concurrent users, or the inability to handle a certain number of transactions per minute, you will need to resolve them in different ways. For example, you can do load testing when optimizing for concurrent use, or use a memory cache along with a database to optimize the number of transactions handled per minute.

Building a simple prototype for your application can be done in parallel with defining use cases and requirements.

Just Imagine

# Structuring an audio metadata CLI application

The first step to building a CLI application is creating the folder structure, but if you aren’t starting from scratch, determine where the CLI application may be added. Suppose the existing structure for the audio metadata API application was built with a domain-driven architecture. To understand how it may be structured, let’s categorize the building blocks of the application.

## Bounded context

**Bounded context** brings a deeper meaning to entities and objects. In the case of our metadata application, consumers utilize the API to search for audio transcription. The operations team would like to search for audio metadata using a CLI. API consumers may be interested in both the metadata and audio transcription but other teams may be more focused on the results of audio transcription. Each team brings a different context to the metadata. However, since tags, album, artist, title, and transcription are all considered metadata, they can be encapsulated within a single entity.

## Language

The **language** used to delineate between different contexts is called **ubiquitous language**. Because teams have slightly different meanings for different terms, this language helps to describe the application in terms that are agreed upon by all involved parties.

For the metadata application, the term **metadata** encompasses all the data extracted from audio, including transcription, and **metadata extraction** is the process of extracting technical metadata and transcription from audio. The term **user** refers to any member of an internal team within the larger organization, and the term **audio** to any recorded sound within a specific limit on length.

## Entities and value objects

**Entities** are models of objects defined by the language. Value objects are fields that exist within an entity. For example, the main entities for the metadata CLI are audio and metadata. Metadata is a value object within the audio entity. Also, each extraction type may be its own value object within the Metadata entity. The list of entity and value objects for this audio metadata CLI application includes the following:

-   Audio
-   Metadata
-   Tags
-   Transcripts

## Aggregation

**Aggregation** is the merging of two separate entities. Suppose within the metadata team at an audio company users would like to make corrections to transcriptions, which is primarily handled by artificial intelligence. Although the transcription may be 95% accurate, there is a team of reviewers that can make corrections to transcriptions to reach 99-100% accuracy. There would be two microservices within the metadata application, one being metadata extraction and other being transcription review. A new aggregated entity may be required: **TranscriptionReview**.

## Service

The term service is generic, so this specifically refers to services within the context of the business domain. In the case of the metadata application, the domain services are the metadata service that extracts metadata from audio and a transcription review service that allows users to add corrections to transcription.

## Events

In the context of domain-driven design, events are domain-specific and notify other processes within the same domain of their occurrence. In this particular case, when a user uploads audio, they receive an ID back immediately. However, the metadata extraction process is triggered behind the scenes and rather than continuously polling on the request metadata command or endpoint to retrieve the status of the metadata object, an event can be sent to an event listener service. The CLI could have a command that continuously listens for process completion.

## Repository

A repository is a collection of the domain or entity objects. The repository has the responsibility of adding, updating, getting, and deleting objects. It makes aggregation possible. A repository is implemented within the domain layer, so there should be no knowledge of the specific database or storage – within the domain, the repository is only an interface. In the case of this metadata application, the repository can have different implementations – MongoDB, ElasticSearch, or flat file.

## Creating the structure

Understanding the components of a domain-driven design, specific to an audio metadata CLI, we can start structuring the folders specific to a metadata CLI. Here is an example layout:

```markup
/Users/username/go/src/github.com/audiocompany/audiofile
   |--cmd
   |----api
   |----cli
   |--extractors
   |----tags
   |----transcript
   |--internal
   |----interfaces
   |--models
   |--services
   |----metadata
   |--storage
   |--vendor
```

### Main folders

Each folder is is follows:

-   `cmd`: The command folder is the main entry point for two different applications that are a part of the audio metadata project: the API and CLI.
-   `extractors`: This folder will hold the packages that will extract metadata from the audio. Although this extractor list will grow, we can start with a few extractor packages: `tags` and `transcript`.
-   `models`: This folder will hold all the structs for the domain entities. The main entities to include are audio and metadata. Each of the extractors may also have its own data model and can be stored in this folder.
-   `services`: Three services have been defined in our previous discussion – the metadata (extraction) service, the transcript review service, and an event listener service, which will listen for processing events and output notifications. Existing and new services exist within this folder.
-   `storage`: The interface and individual implementations for storage exist within this folder.

Just Imagine

# Summary

Throughout this chapter, we have learned how to create a structure for a new application based on the unique requirements of the business domain. We looked at the most popular folder structures for applications and the pros and cons of each, and how to write documentation on use cases and nonfunctional requirements.

While this chapter provided an example layout and the main folders that exist within that example, remember that this is an example of a more developed project. Start simple, always with a flat structure, but start organizing for your future folder structure as you continue to build. Just bear in mind that your code structure will take time. Rome wasn’t built in a day.

After covering these topics, we then discussed a hypothetical real-world example of a company with a team focused entirely on audio metadata. We followed this up with some of the potential use cases for a CLI offering, which would be a fast and efficient alternative to the existing API.

Finally, we discussed a folder structure that could satisfy the requirements of the CLI and API audio metadata application. In [_Chapter 3_](https://subscription.imaginedevops.io/book/programming/9781804611654/2B18883_03.xhtml#_idTextAnchor061), _Building an Audio Metadata CLI_, we will build out the folder structure with the models, interfaces, and implementations to get the CLI application working. That concludes this chapter on how to structure your Go CLI application! Hopefully, it will help you get started.

Just Imagine

# Questions

1.  If you want to share packages with external applications or users, what common folder would these packages reside in?
2.  In ports-and-adapters, or hexagonal, architecture, what are the ports and what are the adapters?
3.  For listing audio, in a real-world example, how would you define the actors, preconditions, and triggers of this use case?

Just Imagine

# Answers

1.  The `pkg` folder contains code that may be used by external applications.
2.  Ports are the interfaces and the adapters are the implementations in a hexagonal architecture. Ports allow communication between different layers of the architecture while the adapters provide the actual implementation.
3.  The actors are the operations team members or any user of the CLI. A precondition of the use case is that the user must be authenticated first. The use case is triggered by either the API’s /list endpoint for the metadata service or running the CLI command for listing audio.

Just Imagine

# Further reading

-   Kat Zein – _How Do you Structure Your Go Apps_ from GopherCon 2018 ([https://www.youtube.com/watch?v=oL6JBUk6tj0](https://www.youtube.com/watch?v=oL6JBUk6tj0)) – an excellent talk about the most common folder structures for Go applications