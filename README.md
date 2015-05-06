# ccat

`ccat` is the colorizing `cat`. It works similar to `cat` but displays content with syntax highlighting.

# Installation

## OSX

```
$ brew tap jingweno/ccat
$ brew install ccat
```

Reference: [jingweno/homebrew-ccat](https://github.com/jingweno/homebrew-ccat)

## From source

Prerequisites:
- [Git](http://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Mercurial](http://mercurial.selenic.com/)
- [Go](https://golang.org/doc/install)

```
$ go install github.com/jingweno/ccat
```

# Examples

```
$ ccat FILE1 FILE2 ...
$ ccat # read from standard input
```

# Demo

![demo](https://dl.dropboxusercontent.com/u/1079131/ccat.gif)

# Credits

Thanks to [Sourcegraph](https://github.com/sourcegraph) who built [this](https://github.com/sourcegraph/syntaxhighlight) awesome syntax-highlighting package.
