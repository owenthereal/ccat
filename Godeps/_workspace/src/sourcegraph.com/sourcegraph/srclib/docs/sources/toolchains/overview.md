# Language toolchains

A **toolchain** is a program that implements functionality for analyzing
projects and source code, according to the specifications defined in this
document.

A **tool** is a subcommand of a toolchain that executes one particular action
(e.g., runs one type of analysis on source code). Tools accept command-line
arguments and produce output on stdout. Common operations implemented by tools
include:

* Scanning: runs before all other tools and finds all *source units* of the
  language (e.g., Python packages, Ruby gems, etc.) in a directory tree.
  Scanners also determine which other tools to call on each source unit.
* Dependency resolution: resolves raw dependencies to git/hg clone URIs,
  subdirectories, and commit IDs if possible (e.g., `foo@0.2.1` to
  github.com/alice/foo commit ID abcd123).
* Graphing: performs type checking/inference and static analysis (called
  "graphing") on the language's source units and dumps data about all
  definitions and references.

Toolchains may contain any number of tool subcommands. For example, a toolchain
could implement only a Go scanner, and the scanner would specify that the Go
packages it finds should be graphed by a tool in a different toolchain. Or a
single toolchain could implement the entire set of Go source analysis
operations.

srclib ships with a default set of toolchains for some popular programming languages.
Repository authors and srclib users may install third-party toolchains to add
features or override the default toolchains.

A **toolchain path** is its repository's clone URI joined with the toolchain's
path within that repository. For example, a toolchain defined in the root
directory of the repository "github.com/alice/srclib-python" would have the
toolchain path "github.com/alice/srclib-python".

A tool is identified by its toolchain's path and the name of the operation it
performs. For example, "github.com/alice/srclib-python scan".

Repository authors can choose which toolchains and tools to use in their
project's Srcfile. If none are specified, the defaults apply.


# Toolchain discovery

The **SRCLIBPATH environment variable** lists places to look for srclib toolchains.
The value is a colon-separated string of paths. If it is empty, `$HOME/.srclib`
is used.

If DIR is a directory listed in SRCLIBPATH, the directory
"DIR/github.com/foo/bar" defines a toolchain named "github.com/foo/bar".

Toolchain directories must contain a Srclibtoolchain file describing and configuring the
toolchain. To see all available toolchains, run `src info toolchains`.

## Tool discovery

A toolchain's tools are described in its Srclibtoolchain file. To see all
available tools (provided by all available toolchains), run `src info tools`.


# Running tools

There are 2 modes of execution for srclib tools:

1.  As a normal **installed program** on your system: to produce analysis
    that relies on locally installed compiler/interpreter and dependency
    versions. (Used when you'll consume the analysis results locally, such as
    during editing of local code.)

    An installed tool is an executable program located at "TOOLCHAIN/.bin/NAME",
    where TOOLCHAIN is the toolchain path and NAME is the last component in the
    toolchain path. For example, the installed tool for "github.com/foo/bar"
    would be at "SRCLIBPATH/github.com/foo/bar/.bin/bar".

2.  Inside a **Docker container**: to produce analysis independent of your local
    configuration and versions. (Used when other people or services will reuse
    the analysis results, such as on [Sourcegraph](https://sourcegraph.com).)

    A Docker-containerized tool is a directory (under SRCLIBPATH) that contains a
    Dockerfile. There is no installation necessary for these tools; the `src`
    program knows how to build and run their Docker container.

    When the Docker container runs, the project's source code is always
    volume-mounted at `/src` (in the container).

Tools may support either or both of these execution modes. Their behavior should
be the same, if possible, regardless of the execution mode.

<!---
TODO(sqs): Clarify this. What does "should be the same" mean?
--->

Toolchains are not typically invoked directly by users; the `src` program invokes
them as part of higher-level commands. However, it is possible to invoke them
directly. To run a tool, run:

```
src tool TOOLCHAIN TOOL [ARG...]
```

For example:

```
src tool github.com/alice/srclib-python scan ~/my-python-project
```


# Toolchain & tool specifications

Toolchains and their tools must conform to the protocol described below. The
protocol is the same whether the tool is run directly or inside a Docker
container. This includes four subcommands, listed below:

## info
This command should display a human-readable info describing
the version, author, etc. (free-form)

## scan (scanners)

Tools that perform the `scan` operation are called **scanners**. They scan a
directory tree and produce a JSON array of source units (in Go,
`[]*unit.SourceUnit`) they encounter.

**Arguments:** none; scanners scan the tree rooted at the current directory (typically the root directory of a repository)

**Stdin:** JSON object representation of repository config (typically `{}`)

**Options:**

* `--repo URI`: the URI of the repository that contains the directory tree being
  scanned
* `--subdir DIR`: the path of the current directory (in which the scanner is
  run), relative to the root directory of the repository being scanned (this is
  typically the root, `"."`, as it is most useful to scan the entire
  repository)

**Stdout:** `[]*unit.SourceUnit`. For a more detailed description, [read the scanner output spec](scanner-output.md).

See the `scan.Scan` function for an implementation of the calling side of this
protocol.

Scanners sometimes need the option values to produce the correct paths. (E.g.,
the Go scanner requires these to generate the right import paths when running in
Docker, since the system GOPATH is not available in the container.)



## depresolve (dependency resolvers)

Tools that perform the `dep` operation are called **dependency resolvers**. They
resolve "raw" dependencies, such as the name and version of a dependency
package, into a full specification of the dependency's target.

**Arguments:** none

**Stdin:** JSON object representation of a source unit (`*unit.SourceUnit`)

**Options:** none

**Stdout:** `[]*dep.Resolution` JSON array with each item corresponding to the
same-index entry in the source unit's `Dependencies` field. For a more
detailed description, [read the dependency resolution output sepc](dependency-resolution-output.md).

## graph  (graphers)

Tools that perform the `graph` operation are called **graphers**. Depending on
the programming language of the source code they analyze, they perform a
combination of parsing, static analysis, semantic analysis, and type inference.
Graphers perform these operations on a source unit and have read access to all
of the source unit's files.

**Arguments:** none

**Stdin:** JSON object representation of a source unit (`*unit.SourceUnit`)

**Options:** none

**Stdout:** JSON graph output (`grapher.Output`). field. For a more
detailed description, [read the grapher output spec](grapher-output.md).

<!---
TODO(sqs): Can we provide the output of `dep` to the `graph` tool? Usually
graphers have to resolve all of the same deps that `dep` would have to. But
we're already providing a full JSON object on stdin, so making it an array or
sending another object would slightly complicate things.
--->

# Available Toolchains

<!--- Stolen from overview.md. --->

<ul>
  <li><a href="go.md">Go</a></li>
  <li><a href="java.md">Java</a></li>
  <li><a href="python.md">Python</a></li>
  <li><a href="javascript.md">JavaScript</a></li>
  <li><a href="haskell.md">Haskell</a></li>
  <li><a href="ruby.md">Ruby (WIP)</a></li>
  <li><a href="php.md">PHP (WIP)</a></li>

</ul>
