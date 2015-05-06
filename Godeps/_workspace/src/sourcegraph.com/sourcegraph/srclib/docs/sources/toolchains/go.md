# Go Toolchain
## Installation

This toolchain is not a standalone program; it provides additional functionality
to editor plugins and other applications that use [srclib](https://srclib.org).

Browse the code at [srclib-go](https://sourcegraph.com/sourcegraph/srclib-go) and [file issues on GitHub](https://github.com/sourcegraph/srclib-go).

<iframe src="http://ghbtns.com/github-btn.html?user=sourcegraph&repo=srclib-go&type=watch&count=true&size=large"
  allowtransparency="true" frameborder="0" scrolling="0" width="170" height="30"></iframe>

First,
[install the `src` program (see srclib installation instructions)](../install.md#install-srclib).

Then run:

```bash
# download and fetch dependencies
go get -v sourcegraph.com/sourcegraph/srclib-go
cd $GOPATH/sourcegraph.com/sourcegraph/srclib-go

# build the srclib-go program in .bin/srclib-go (this is currently required by srclib to discover the program)
make

# link this toolchain in your SRCLIBPATH (default ~/.srclib) to enable it
src toolchain add sourcegraph.com/sourcegraph/srclib-go
```

To verify that installation succeeded, run:

```
src toolchain list
```

You should see this srclib-go toolchain in the list.

Now that this toolchain is installed, any program that relies on srclib (such as
editor plugins) will support Go.
