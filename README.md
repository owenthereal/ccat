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
$ brew install jingweno/ccat/ccat
```

Reference: [jingweno/homebrew-ccat](https://github.com/jingweno/homebrew-ccat)

### Arch Linux

```
$ pacaur -S ccat
$ pacaur -S ccat-git
```
The ccat package will reflect the current release snapshot, while the ccat-git will be based on the current source available in the master branch of the git repo. You can use any AUR helper in place of pacaur [AUR Helpers](https://wiki.archlinux.org/index.php/AUR_helpers)

### From source

Prerequisites:
- [Git](http://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Go](https://golang.org/doc/install) 1.4+

```
$ go get -u github.com/jingweno/ccat
```

## Usage

```
$ ccat FILE1 FILE2 ...
$ ccat --bg=dark FILE1 FILE 2 ... # dark background
$ ccat # read from standard input
$ curl https://raw.githubusercontent.com/jingweno/ccat/master/main.go | ccat
```

## Demo

![demo](https://dl.dropboxusercontent.com/u/1079131/ccat.gif)

## Roadmap

- [ ] nicer default color scheme
- [ ] customizable color scheme
- [ ] ?

## FAQ

### Why not pygments?

You could use Python's [pygments](http://pygments.org/) to achieve pretty much the same thing:

```
$ alias ccat="pygmentize -g"
$ ccat FILE1
```

`ccat` is a \*nix alternative to pygments: no interpreter, one binary, native speed, POSIX standard etc..

### Why not GNU source-highlight?

You could also use [GNU source-highlight]() to perform the same task:

```
$ alias ccat="source-highlight -fesc -o STDOUT -i"
$ ccat FILE1
```

`ccat` is an alternative to source-highlight. Written in go.

## License

[MIT](https://github.com/jingweno/ccat/blob/master/LICENSE)

## Credits

Thanks to [Sourcegraph](https://github.com/sourcegraph) who built [this](https://github.com/sourcegraph/syntaxhighlight) awesome syntax-highlighting package.
