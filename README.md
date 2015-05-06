# ccat

`ccat` is the colorizing `cat`. It works similar to `cat` but displays content with syntax highlighting.

## Supported Languages

* JavaScript
* Java
* Ruby
* Python
* Go
* C

## Installation

### OSX

```
$ brew tap jingweno/ccat
$ brew install ccat
```

Reference: [jingweno/homebrew-ccat](https://github.com/jingweno/homebrew-ccat)

### From source

Prerequisites:
- [Git](http://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Go](https://golang.org/doc/install) 1.4+

```
$ go get -u github.com/jingweno/ccat
```

## Examples

```
$ ccat FILE1 FILE2 ...
$ ccat # read from standard input
```

## Demo

![demo](https://dl.dropboxusercontent.com/u/1079131/ccat.gif)

## Credits

Thanks to [Sourcegraph](https://github.com/sourcegraph) who built [this](https://github.com/sourcegraph/syntaxhighlight) awesome syntax-highlighting package.

## FAQ

### Why not pygments?

You could use [pygments](http://pygments.org/) to achieve pretty much the same thing:

```
$ alias ccat="pygmentize -g"
$ ccat FILE1
```

`ccat` is a \*nix alternative to pygments: no interpreter, one binary, native speed, POSIX standard etc..

## License

[MIT](https://github.com/jingweno/ccat/blob/master/LICENSE)
