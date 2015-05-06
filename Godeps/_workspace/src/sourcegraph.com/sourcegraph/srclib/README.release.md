# srclib release process

We use [Equinox](https://equinox.io) to release cross-compiled Go binaries for a
variety of platforms.

Releases are signed using a private key. Committers with release privileges
should already have this private key. Releasers should also have
`$EQUINOX_ACCOUNT` and `$EQUINOX_SECRET` set in their environment.

To issue a new release:

```
make release V=x.y.z
```

where x.y.z is the version number of the release.


## Cross-compiling

All releases are cross-compiled for multiple platforms. See the `make
upload-release` target definition for a list of these platforms.

Note: to cross-compile Go binaries, you'll have to perform a one-time setup step
in your GOROOT:

```
cd $GOROOT/src
GOOS=darwin GOARCH=amd64 ./make.bash --no-clean
GOOS=linux GOARCH=386 ./make.bash --no-clean
GOOS=linux GOARCH=amd64 ./make.bash --no-clean
```

You can omit the command that contains your current `GOOS` and `GOARCH`, as
you've already bootstrapped that combo.

Users of `src` can check for updates with `src version` and update the program
with `src selfupdate`.
