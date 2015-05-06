# Frequently asked questions
## Questions About Toolchains
### Why are there 2 different execution schemes for tools: running directly as a program and running inside a Docker container?

Running toolchains as normal programs means that it'll pick up context from your
local system, such as interpreter/compiler versions, dependencies, etc. This is
desirable when you're using `src` to analyze code that you're editing locally,
and when you don't want the analysis output to be reproducible by other people.

Running toolchains inside Docker means that analysis occurs in a reproducible
environment, with fixed versions of interpreters/compilers and dependencies.
This is desirable when you want to share analysis output with others, such as
when you're uploading it to an external service (like
[Sourcegraph](https://sourcegraph.com)).

### Why can't a toolchain's tools use different Docker images from each other?

The Docker image built for a toolchain should be capable of running all of the
toolchain's functionality. It would add a lot of complexity to either:

* allow toolchains to contain multiple Dockerfiles (some of which would probably
  be `FROM` others); or
* allow tools to generate new Dockerfiles (and then build them with `docker
  build - < context.tar`) or run sub-Docker containers.

If a tool truly can't reuse the scanner's Dockerfile, then move it to a separate
toolchain.
