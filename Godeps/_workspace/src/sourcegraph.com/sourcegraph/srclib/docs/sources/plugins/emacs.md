# Emacs Plugin

<div class="embed-responsive embed-responsive-16by9">
<iframe class="embed-responsive-item" src="//www.youtube.com/embed/cm59qQD6khs" frameborder="0" allowfullscreen></iframe>
</div>

## Features
- Documentation lookups
- Type information
- Find usages (across all open-source projects globally)

<br>

## Installation Instructions
First, make sure you've [installed srclib](../install.md), along with the toolchains for the programming
languages you're using. Once srclib is installed, you can install the emacs plugin by navigating to your `.emacs.d`
directory and cloning the repository.


```bash
cd ~/.emacs.d
git clone https://github.com/sourcegraph/emacs-sourcegraph-mode.git
```

To install the plugin, append the following code to `~/.emacs.d/init.el`.
```lisp
(add-to-list 'load-path "~/.emacs.d/emacs-sourcegraph-mode")
(require 'sourcegraph-mode)
```

Sourcegraph-mode can be enabled in a buffer with M-x, then typing `sourcegraph-mode`.

Now, in any file (with `sourcegraph-mode` enabled), run `sourcegraph-describe`
(or C-M-.) to see docs, type info, and examples.

[//]: # (Issue #134:  sourcegraph mode for the default CLI emacs under OS X doesn't work)
> *Note that the default CLI emacs app in OS X doesn't enable json-mode. Do not use it. If you do, there will be always one error message to load the .emacs file.*

<br>

## Contribute on GitHub
<iframe src="http://ghbtns.com/github-btn.html?user=sourcegraph&repo=emacs-sourcegraph-mode&type=watch&count=true&size=large"
  allowtransparency="true" frameborder="0" scrolling="0" width="170" height="30"></iframe>

<br>
