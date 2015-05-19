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
$ brew install ccat
```

### Arch Linux

```
$ pacaur -S ccat
$ pacaur -S ccat-git
```
The ccat package will reflect the current release snapshot, while the ccat-git will be based on the current source available in the master branch of the git repo. You can use any AUR helper in place of pacaur [AUR Helpers](https://wiki.archlinux.org/index.php/AUR_helpers)

### Standalone

`ccat` can be easily installed as an executable.
Download the latest [compiled binaries](https://github.com/jingweno/ccat/releases) and put it in your executable path.

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

## Alternatives

`ccat` is designed to be distributed in one binary, run at native speed
and follow the POSIX standards. There're alternatives out there.
Use them at your own risk :):

* [pygments](http://pygments.org/)
* [source-highlight](https://www.gnu.org/software/src-highlite/)

## License

[MIT](https://github.com/jingweno/ccat/blob/master/LICENSE)

## Credits

Thanks to [Sourcegraph](https://github.com/sourcegraph) who built [this](https://github.com/sourcegraph/syntaxhighlight) awesome syntax-highlighting package.
