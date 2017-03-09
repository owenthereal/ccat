# ccat
[![OpenCollective](https://opencollective.com/ccat/backers/badge.svg)](#backers) 
[![OpenCollective](https://opencollective.com/ccat/sponsors/badge.svg)](#sponsors)

`ccat` is the colorizing `cat`. It works similar to `cat` but displays content with syntax highlighting.

## Supported Languages

* JavaScript
* Java
* Ruby
* Python
* Go
* C
* JSON

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
$ ccat FILE1 FILE2 ... --html # output in HTML
$ ccat --bg=dark FILE1 FILE 2 ... # dark background
$ ccat -G String="_darkblue_" -G Plaintext="darkred" FILE # set color codes
$ ccat --palette # show palette
$ ccat # read from standard input
$ curl https://raw.githubusercontent.com/jingweno/ccat/master/main.go | ccat
```

It's recommended to alias `ccat` to `cat`:

```
alias cat=ccat
```

The overhead of `ccat` comparing to `cat` is mimimum:

```
$ wc -l main.go
123 main.go
$ time cat main.go > /dev/null
cat main.go > /dev/null  0.00s user 0.00s system 61% cpu 0.005 total
$ time ccat main.go > /dev/null
ccat main.go > /dev/null  0.00s user 0.00s system 78% cpu 0.007 total
```

You can always invoke `cat` after aliasing `ccat` by typing `\cat`.

## Demo

[![demo](https://asciinema.org/a/21858.png)](https://asciinema.org/a/21858)

## Roadmap

- [ ] nicer default color scheme
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

## Support

### Backers
Support us with a monthly donation and help us continue our activities. [[Become a backer](https://opencollective.com/ccat#backer)]

<a href="https://opencollective.com/ccat/backer/0/website" target="_blank"><img src="https://opencollective.com/ccat/backer/0/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/1/website" target="_blank"><img src="https://opencollective.com/ccat/backer/1/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/2/website" target="_blank"><img src="https://opencollective.com/ccat/backer/2/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/3/website" target="_blank"><img src="https://opencollective.com/ccat/backer/3/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/4/website" target="_blank"><img src="https://opencollective.com/ccat/backer/4/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/5/website" target="_blank"><img src="https://opencollective.com/ccat/backer/5/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/6/website" target="_blank"><img src="https://opencollective.com/ccat/backer/6/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/7/website" target="_blank"><img src="https://opencollective.com/ccat/backer/7/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/8/website" target="_blank"><img src="https://opencollective.com/ccat/backer/8/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/9/website" target="_blank"><img src="https://opencollective.com/ccat/backer/9/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/10/website" target="_blank"><img src="https://opencollective.com/ccat/backer/10/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/11/website" target="_blank"><img src="https://opencollective.com/ccat/backer/11/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/12/website" target="_blank"><img src="https://opencollective.com/ccat/backer/12/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/13/website" target="_blank"><img src="https://opencollective.com/ccat/backer/13/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/14/website" target="_blank"><img src="https://opencollective.com/ccat/backer/14/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/15/website" target="_blank"><img src="https://opencollective.com/ccat/backer/15/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/16/website" target="_blank"><img src="https://opencollective.com/ccat/backer/16/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/17/website" target="_blank"><img src="https://opencollective.com/ccat/backer/17/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/18/website" target="_blank"><img src="https://opencollective.com/ccat/backer/18/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/19/website" target="_blank"><img src="https://opencollective.com/ccat/backer/19/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/20/website" target="_blank"><img src="https://opencollective.com/ccat/backer/20/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/21/website" target="_blank"><img src="https://opencollective.com/ccat/backer/21/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/22/website" target="_blank"><img src="https://opencollective.com/ccat/backer/22/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/23/website" target="_blank"><img src="https://opencollective.com/ccat/backer/23/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/24/website" target="_blank"><img src="https://opencollective.com/ccat/backer/24/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/25/website" target="_blank"><img src="https://opencollective.com/ccat/backer/25/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/26/website" target="_blank"><img src="https://opencollective.com/ccat/backer/26/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/27/website" target="_blank"><img src="https://opencollective.com/ccat/backer/27/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/28/website" target="_blank"><img src="https://opencollective.com/ccat/backer/28/avatar.svg"></a>
<a href="https://opencollective.com/ccat/backer/29/website" target="_blank"><img src="https://opencollective.com/ccat/backer/29/avatar.svg"></a>


### Sponsors
Become a sponsor and get your logo on our README on Github with a link to your site. [[Become a sponsor](https://opencollective.com/ccat#sponsor)]

<a href="https://opencollective.com/ccat/sponsor/0/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/1/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/2/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/3/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/3/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/4/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/4/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/5/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/5/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/6/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/6/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/7/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/7/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/8/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/8/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/9/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/9/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/10/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/10/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/11/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/11/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/12/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/12/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/13/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/13/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/14/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/14/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/15/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/15/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/16/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/16/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/17/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/17/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/18/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/18/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/19/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/19/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/20/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/20/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/21/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/21/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/22/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/22/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/23/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/23/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/24/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/24/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/25/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/25/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/26/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/26/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/27/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/27/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/28/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/28/avatar.svg"></a>
<a href="https://opencollective.com/ccat/sponsor/29/website" target="_blank"><img src="https://opencollective.com/ccat/sponsor/29/avatar.svg"></a>
