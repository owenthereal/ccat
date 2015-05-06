# Data Model Overview

### [Source unit](../toolchains/scanner-output.md)

A source unit is code that Sourcegraph processes as one logical unit. It is
conceptually similar to a compilation unit.

A source unit consists of a type, a set of files and directories, along with
metadata specific to the type of source unit. Typically a source unit is the
level of code structure that creates a package or binary product that other
projects can depend on.

Examples include:

* A Go package
* A Node.js package
* A Ruby gem
* A Python pip package

A file or directory may appear in more than one source unit. For example, a
frontend JavaScript library and a node.js package might include the same
JavaScript file (to make it usable in the browser and in node.js, respectively).
Or 2 Ruby gems in the same repository might refer to the same `.rb` file.

Examples of things that would NOT be good source units and the reasons why:

* Things that you can import, include, or require in source code are not
  necessarily good source units. For example, a Python package can be imported
  in source code (`import mypkg`), but the information necessary to build a
  Python package is not contained within the package itself (it's usually in
  setup.py). And when an external project wants to depend on a Python package,
  they first must specify the dependency to the Python pip package containing
  that Python package. So, the pip package is the right level to use as the
  source unit.


NOTE: We haven't yet found a succinct way to describe a source unit. Modify this
document as we come up with clearer explanations.

## [Defs and refs](../toolchains/grapher-output.md)

For each source unit, a graph is generated from the associated code, consisting
of **definitions** and **references**. A definition is the original point in
code where an object is defined, be it a function, module, or class. References
are anywhere that def is then later used in code--references then link back to
the original definition. References can link to definitions in other source
units and repositories.

See [grapher output specs](../toolchains/grapher-output.md) for more
information.

## [Toolchain](../toolchains/overview.md)

A primary tool called 'src' will serve as a harness for all of the individual
language toolchains. It will also serve as the API endpoint, for users to query
for data. Each language-specific toolchain is composed of five parts - any
language that implements these can be automatically used with srclib.

1. **Scanner** - Traverses the directory tree looking for source units. Invokes
   the grapher and dependency lister on each unit.
2. **Grapher** - Produces the code graph through static analysis.
3. **Dependency Resolver** - Resolves raw dependency names to VCS repository
   urls.
4. **Formatter** - Exports functions that the API uses to transform raw
   language-specific data into a generic, useable structure. This is the only
   tool that must be written in Go.

This separation allows certain parts to be implemented before others, allowing
incremental benefits to accrue before the toolchain is complete. For example,
the grapher could be implemented before proper dependency resolution, allowing
jump to definition and documentation only within a local codebase. The
dependency resolution could be implemented without the grapher, allowing
dependency tracking, without code graphing.

The src tool, after invoking the various elements of the toolchain, will extract
other information, such as the repository URL and the commit ID, and use git/hg
blame to find authorship information for definitions and references.
